package pubnub

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"context"

	"github.com/pubnub/go/utils"
)

// Events:
// - ConnectedCategory - after connection established
// - DisconnectedCategory - after subscription loop stops for any reason (no
// channels left or error happend)

// Unsubscribe.
// When you unsubscirbe from channel or channel group the following events
// happens:
// - LoopStopCategory - immediately when no more channels or channel groups left
// to subscribe
// - PNUnsubscribeOperation - after leave request was fulfilled and server is
// notified about unsubscibed items
//
// Announcement
// Status, Message and Presence announcement happens in a distinct goroutine.
// It doesn't block subscribe loop.
// Keep in mind that each listener will receive the same pointer to a response
// object. You may wish to create a shallow copy of either the response or the
// response message by you own to not affect the other listeners.

// Heartbeat:
// - Heartbeat is enabled by default.
// - Default presence timeout is 300 seconds.
// - The first Heartbeat request will be scheduled to be executed after
// getHeartbeatInterval() seconds (default - 149).

type SubscriptionManager struct {
	sync.RWMutex

	subscriptionLock sync.Mutex

	listenerManager *ListenerManager
	stateManager    *StateManager
	pubnub          *PubNub
	transport       http.RoundTripper

	hbMutex sync.Mutex
	hbTimer *time.Ticker
	hbDone  chan bool

	messages chan subscribeMessage
	ctx      Context
	// hCtx            Context
	subscribeCancel func()
	heartbeatCancel func()

	// Store the latest timetoken to subscribe with, null by default to get the
	// latest timetoken.
	timetoken int64

	// When changing the channel mix, store the timetoken for a later date
	storedTimetoken int64

	region int8

	subscriptionStateAnnounced bool
	heartbeatStopCalled        bool

	filterExpression string
}

type SubscribeOperation struct {
	Channels         []string
	ChannelGroups    []string
	PresenceEnabled  bool
	Timetoken        int64
	FilterExpression string
	Transport        http.RoundTripper
}

type UnsubscribeOperation struct {
	Channels      []string
	ChannelGroups []string
}

type StateOperation struct {
	channels      []string
	channelGroups []string
	state         map[string]interface{}
}

// use Context for cancelation
// different PubNub instances should not affect each others subscriptions
func newSubscriptionManager(pubnub *PubNub) *SubscriptionManager {
	manager := &SubscriptionManager{}

	manager.pubnub = pubnub

	manager.listenerManager = newListenerManager()
	manager.stateManager = newStateManager()

	manager.Lock()
	manager.timetoken = 0
	manager.storedTimetoken = -1
	manager.subscriptionStateAnnounced = false
	manager.ctx, manager.subscribeCancel = context.WithCancel(context.Background())
	// manager.hCtx, manager.heartbeatCancel = context.WithCancel(context.Background())
	manager.messages = make(chan subscribeMessage, 1000)
	manager.Unlock()

	go subscribeMessageWorker(manager.listenerManager, manager.messages)

	// actions:
	// add channel
	// remove channel
	// add listener
	// remove listener
	// unsubscribe all
	// cancel
	// addListeners := func()
	return manager
}

func (m *SubscriptionManager) adaptState(stateOperation StateOperation) {
	m.stateManager.adaptStateBuilder(stateOperation)
}

func (m *SubscriptionManager) adaptSubscribe(
	subscribeOperation *SubscribeOperation) {
	m.stateManager.adaptSubscribeBuilder(subscribeOperation)

	log.Println("adapting a new subscription", subscribeOperation.Channels,
		subscribeOperation.PresenceEnabled)
	m.Lock()

	m.subscriptionStateAnnounced = false

	if subscribeOperation.Timetoken != 0 {
		m.timetoken = subscribeOperation.Timetoken
	}

	if m.timetoken != 0 {
		m.storedTimetoken = m.timetoken
	}

	if subscribeOperation.FilterExpression != "" {
		m.filterExpression = subscribeOperation.FilterExpression
	}

	m.timetoken = 0

	m.Unlock()

	m.reconnect()
}

func (m *SubscriptionManager) adaptUnsubscribe(
	unsubscribeOperation *UnsubscribeOperation) {
	m.stateManager.adaptUnsubscribeBuilder(unsubscribeOperation)

	m.Lock()
	m.subscriptionStateAnnounced = false

	go func() {
		_, err := m.pubnub.Leave().Channels(unsubscribeOperation.Channels).
			ChannelGroups(unsubscribeOperation.ChannelGroups).Execute()

		if err != nil {
			m.listenerManager.announceStatus(&PNStatus{
				Category:              BadRequestCategory,
				ErrorData:             err,
				Error:                 true,
				Operation:             PNUnsubscribeOperation,
				AffectedChannels:      unsubscribeOperation.Channels,
				AffectedChannelGroups: unsubscribeOperation.ChannelGroups,
			})
		} else {
			m.listenerManager.announceStatus(&PNStatus{
				Category:              AcknowledgmentCategory,
				StatusCode:            200,
				Operation:             PNUnsubscribeOperation,
				Uuid:                  m.pubnub.Config.Uuid,
				AffectedChannels:      unsubscribeOperation.Channels,
				AffectedChannelGroups: unsubscribeOperation.ChannelGroups,
			})
		}
	}()

	if m.stateManager.isEmpty() {
		m.region = 0
		m.storedTimetoken = -1
		m.timetoken = 0
	} else {
		m.storedTimetoken = m.timetoken
		m.timetoken = 0
	}

	m.Unlock()

	m.reconnect()
}

func (m *SubscriptionManager) startSubscribeLoop() {
	for {
		combinedChannels := m.stateManager.prepareChannelList(true)
		combinedGroups := m.stateManager.prepareGroupList(true)

		if len(combinedChannels) == 0 && len(combinedGroups) == 0 {
			m.listenerManager.announceStatus(&PNStatus{
				Category: DisconnectedCategory,
			})
			m.log("no channels left to subscribe")
			break
		}

		m.Lock()
		tt := m.timetoken
		ctx := m.ctx
		m.Unlock()

		opts := &SubscribeOpts{
			pubnub:    m.pubnub,
			Channels:  combinedChannels,
			Groups:    combinedGroups,
			Timetoken: tt,
			Transport: m.transport,
			Heartbeat: m.pubnub.Config.PresenceTimeout,
			ctx:       ctx,
		}

		// TODO: use context to be able to stop request
		res, _, err := executeRequest(opts)
		if err != nil {
			if strings.Contains(err.Error(), "timeout") {
				m.listenerManager.announceStatus(&PNStatus{
					Category: TimeoutCategory,
				})
				continue
			}

			if strings.Contains(err.Error(), "context canceled") {
				m.listenerManager.announceStatus(&PNStatus{
					Category: CancelledCategory,
				})
				break
			}

			if strings.Contains(err.Error(), "Forbidden") {
				m.listenerManager.announceStatus(&PNStatus{
					Category: AccessDeniedCategory,
				})
				m.unsubscribeAll()
				break
			}

			if strings.Contains(err.Error(), "400") ||
				strings.Contains(err.Error(), "Bad Request") {
				m.listenerManager.announceStatus(&PNStatus{
					Category: BadRequestCategory,
				})
				m.unsubscribeAll()
				break
			}

			m.listenerManager.announceStatus(&PNStatus{
				Category: UnknownCategory,
			})

			break
		}

		m.Lock()
		announced := m.subscriptionStateAnnounced
		m.Unlock()

		if announced == false {
			m.listenerManager.announceStatus(&PNStatus{
				Category: ConnectedCategory,
			})

			m.subscriptionStateAnnounced = true
		}

		var envelope subscribeEnvelope
		err = json.Unmarshal(res, &envelope)
		if err != nil {
			m.listenerManager.announceStatus(&PNStatus{
				Category:              BadRequestCategory,
				ErrorData:             err,
				Error:                 true,
				Operation:             PNSubscribeOperation,
				AffectedChannels:      combinedChannels,
				AffectedChannelGroups: combinedGroups,
			})
		}

		if len(envelope.Messages) > 0 {
			for _, message := range envelope.Messages {
				if m.pubnub.Config.CipherKey != "" {
					msg, _ := message.Payload.(string)
					decryptedMsg, err := utils.DecryptString(m.pubnub.Config.CipherKey, msg)

					if err != nil {
						m.listenerManager.announceStatus(&PNStatus{
							Category:              BadRequestCategory,
							ErrorData:             err,
							Error:                 true,
							Operation:             PNSubscribeOperation,
							AffectedChannels:      combinedChannels,
							AffectedChannelGroups: combinedGroups,
						})
						break
					}

					message.Payload = decryptedMsg

					m.messages <- message
				} else {
					m.messages <- message
				}
			}
		}

		if m.storedTimetoken != -1 {
			m.Lock()
			m.timetoken = m.storedTimetoken
			m.storedTimetoken = -1
			m.Unlock()
		} else {
			tt, err := strconv.ParseInt(envelope.Metadata.Timetoken, 10, 64)
			if err != nil {
				m.listenerManager.announceStatus(&PNStatus{
					Category:              BadRequestCategory,
					ErrorData:             err,
					Error:                 true,
					Operation:             PNSubscribeOperation,
					AffectedChannels:      combinedChannels,
					AffectedChannelGroups: combinedGroups,
				})
			}

			m.timetoken = tt
		}

		m.region = envelope.Metadata.Region
	}
}

func (m *SubscriptionManager) startHeartbeatTimer() {
	m.stopHeartbeat()
	m.log("heartbeat: new timer")

	m.hbMutex.Lock()
	m.hbDone = make(chan bool)
	m.hbTimer = time.NewTicker(time.Duration(
		m.pubnub.Config.HeartbeatInterval) * time.Second)

	go func() {
		defer m.hbMutex.Unlock()
		defer func() {
			m.hbDone = nil
		}()

		for {
			select {
			case <-m.hbTimer.C:
				m.performHeartbeatLoop()
			case <-m.hbDone:
				m.log("heartbeat: loop: after stop")
				return
			}
		}
	}()
}

func (m *SubscriptionManager) stopHeartbeat() {
	m.log("heartbeat: loop: stopping...")

	if m.hbTimer != nil {
		m.hbTimer.Stop()
		m.log("heartbeat: loop: timer stopped")
	}

	if m.hbDone != nil {
		m.hbDone <- true
		m.log("heartbeat: loop: done channel stopped")
	}
}

func (m *SubscriptionManager) performHeartbeatLoop() error {
	presenceChannels := m.stateManager.prepareChannelList(false)
	presenceGroups := m.stateManager.prepareGroupList(false)
	stateStorage := m.stateManager.createStatePayload()

	if len(presenceChannels) == 0 && len(presenceGroups) == 0 {
		m.log("heartbeat: no channels left")
		go m.stopHeartbeat()
		return nil
	}

	_, status, err := newHeartbeatBuilder(m.pubnub).
		Channels(presenceChannels).
		ChannelGroups(presenceGroups).
		State(stateStorage).
		Execute()

	if err != nil {
		errStatus := &PNStatus{
			Operation: PNHeartBeatOperation,
			Category:  BadRequestCategory,
			Error:     true,
			ErrorData: err,
		}

		m.listenerManager.announceStatus(errStatus)

		return err
	}

	hbStatus := &PNStatus{
		Category:   UnknownCategory,
		Error:      false,
		Operation:  PNHeartBeatOperation,
		StatusCode: status.StatusCode,
	}

	m.listenerManager.announceStatus(hbStatus)

	return nil
}

type subscribeEnvelope struct {
	Messages []subscribeMessage `json:"m"`
	Metadata struct {
		Timetoken string `json:"t"`
		Region    int8   `json:"r"`
	} `json:"t"`
}

type subscribeMessage struct {
	Shard             string      `json:"a"`
	SubscriptionMatch string      `json:"b"`
	Channel           string      `json:"c"`
	IssuingClientId   string      `json:"i"`
	SubscribeKey      string      `json:"k"`
	Flags             int         `json:"f"`
	Payload           interface{} `json:"d"`
	UserMetadata      interface{} `json:"u"`

	PublishMetaData publishMetadata `json:"p"`
}

type presenceEnvelope struct {
	Action    string
	Uuid      string
	Occupancy int
	Timestamp int64
	Data      interface{}
}

type publishMetadata struct {
	PublishTimetoken string `json:"t"`
	Region           int    `json:"r"`
}

type originationMetadata struct {
	Timetoken int64 `json:"t"`
	Region    int   `json:"r"`
}

func subscribeMessageWorker(lm *ListenerManager, messages <-chan subscribeMessage) {
	for message := range messages {
		processSubscribePayload(lm, message)
	}
}

func processSubscribePayload(lm *ListenerManager, payload subscribeMessage) {
	channel := payload.Channel
	subscriptionMatch := payload.SubscriptionMatch
	publishMetadata := payload.PublishMetaData

	if channel != "" && channel == subscriptionMatch {
		subscriptionMatch = ""
	}

	if strings.Contains(payload.Channel, "-pnpres") {
		var presencePayload map[string]interface{}
		var action, uuid, actualChannel, subscribedChannel string
		var occupancy int
		var timestamp int64
		var data interface{}
		var ok bool

		if presencePayload, ok = payload.Payload.(map[string]interface{}); !ok {
			lm.announceStatus(&PNStatus{
				Category:         UnknownCategory,
				ErrorData:        errors.New("Response presence parsing error"),
				Error:            true,
				Operation:        PNSubscribeOperation,
				AffectedChannels: []string{channel},
			})
			// return
		}

		action, _ = presencePayload["action"].(string)
		uuid, _ = presencePayload["uuid"].(string)
		occupancy, _ = presencePayload["occupancy"].(int)
		timestamp, _ = presencePayload["timestamp"].(int64)
		data = presencePayload["data"]
		timetoken, _ := strconv.ParseInt(publishMetadata.PublishTimetoken, 10, 64)

		strippedPresenceChannel := ""
		strippedPresenceSubscription := ""

		if channel != "" {
			strippedPresenceChannel = strings.Replace(channel, "-pnpres", "", -1)
		}

		if subscriptionMatch != "" {
			actualChannel = channel
			subscribedChannel = subscriptionMatch
			strippedPresenceSubscription = strings.Replace(subscriptionMatch, "-pnpres", "", -1)
		} else {
			subscribedChannel = channel
		}

		pnPresenceResult := &PNPresence{
			Event:             action,
			ActualChannel:     actualChannel,
			SubscribedChannel: subscribedChannel,
			Channel:           strippedPresenceChannel,
			Subscription:      strippedPresenceSubscription,
			State:             data,
			Timetoken:         timetoken,
			Occupancy:         occupancy,
			Uuid:              uuid,
			Timestamp:         timestamp,
		}

		lm.announcePresence(pnPresenceResult)
	} else {
		actualCh := ""
		subscribedCh := channel
		timetoken, _ := strconv.ParseInt(publishMetadata.PublishTimetoken, 10, 64)

		if subscriptionMatch != "" {
			actualCh = channel
			subscribedCh = subscriptionMatch
		}

		pnMessageResult := &PNMessage{
			Message:           payload.Payload,
			ActualChannel:     actualCh,
			SubscribedChannel: subscribedCh,
			Channel:           channel,
			Subscription:      subscriptionMatch,
			Timetoken:         timetoken,
			Publisher:         payload.IssuingClientId,
			UserMetadata:      payload.UserMetadata,
		}

		lm.announceMessage(pnMessageResult)
	}
}

func (m *SubscriptionManager) AddListener(listener *Listener) {
	m.listenerManager.addListener(listener)
}

func (m *SubscriptionManager) RemoveListener(listener *Listener) {
	m.listenerManager.removeListener(listener)
}

func (m *SubscriptionManager) reconnect() {
	m.log("reconnect")

	go m.startSubscribeLoop()
	m.startHeartbeatTimer()
}

func (m *SubscriptionManager) disconnect() {
	m.log("disconnect")

	m.stopHeartbeat()
	m.stopSubscribeLoop()
	// TODO: stop reconnection manager
}

func (m *SubscriptionManager) stopSubscribeLoop() {
	m.log("loop stop")

	if m.ctx != nil && m.subscribeCancel != nil {
		m.subscribeCancel()
		m.Lock()
		m.ctx, m.subscribeCancel = context.WithCancel(context.Background())
		m.Unlock()
	}
}

func (m *SubscriptionManager) getSubscribedChannels() []string {
	return m.stateManager.prepareChannelList(false)
}

func (m *SubscriptionManager) getSubscribedGroups() []string {
	return m.stateManager.prepareGroupList(false)
}

func (m *SubscriptionManager) unsubscribeAll() {
	m.adaptUnsubscribe(&UnsubscribeOperation{
		Channels:      m.stateManager.prepareChannelList(false),
		ChannelGroups: m.stateManager.prepareGroupList(false),
	})
}

func (m *SubscriptionManager) log(message string) {
	log.Printf("pubnub: subscribe: %s: %s: %s/%s\n",
		message,
		m.pubnub.Config.Uuid,
		m.stateManager.prepareChannelList(true),
		m.stateManager.prepareGroupList(true))
}
