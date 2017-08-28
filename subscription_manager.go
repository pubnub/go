package pubnub

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
)

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

type SubscriptionManager struct {
	sync.RWMutex

	subscriptionLock sync.Mutex

	listenerManager *ListenerManager
	stateManager    *StateManager
	pubnub          *PubNub

	messages        chan subscribeMessage
	ctx             Context
	subscribeCancel func()

	// Store the latest timetoken to subscribe with, null by default to get the
	// latest timetoken.
	timetoken int64

	// When changing the channel mix, store the timetoken for a later date
	storedTimetoken int64

	region int8

	subscriptionStateAnnounced bool

	filterExpression string
}

type SubscribeOperation struct {
	Channels         []string
	ChannelGroups    []string
	PresenceEnabled  bool
	Timetoken        int64
	FilterExpression string
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
	// manager.ctx, manager.subscribeCancel = context.WithCancel(context.Background())
	manager.messages = make(chan subscribeMessage, 1000)
	manager.Unlock()

	// go manager.startSubscribeLoop()
	// go manager.startSubscribeLoopWithRoutine()
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

	log.Println("subscribe reconnect")
	m.reconnect()
}

func (m *SubscriptionManager) adaptUnsubscribe(
	unsubscribeOperation *UnsubscribeOperation) {
	m.stateManager.adaptUnsubscribeBuilder(unsubscribeOperation)

	m.Lock()
	m.subscriptionStateAnnounced = false

	go func() {
		err := m.pubnub.Leave().Channels(unsubscribeOperation.Channels).
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

		tt := m.timetoken
		ctx := m.ctx
		filterExpr := m.filterExpression

		opts := &SubscribeOpts{
			pubnub:           m.pubnub,
			Channels:         combinedChannels,
			Groups:           combinedGroups,
			Timetoken:        tt,
			FilterExpression: filterExpr,
			ctx:              ctx,
			// 	// transport
			// 	// config/subkey
			// 	// config/uuid
			// 	// config/timeouts
		}

		// TODO: use context to be able to stop request
		res, err := executeRequest(opts)
		if err != nil {
			if strings.Contains(err.Error(), "context canceled") {
				m.listenerManager.announceStatus(&PNStatus{
					Category: CancelledCategory,
				})
				continue
			}

			// TODO: handle timeout
			// TODO: handle canceled
			// TODO: handle error
			break
		}

		announced := m.subscriptionStateAnnounced

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

		// TODO: fetch messages and if any, push them to the worker queue
		if len(envelope.Messages) > 0 {
			for message := range envelope.Messages {
				m.messages <- envelope.Messages[message]
			}
		}

		if m.storedTimetoken != -1 {
			m.timetoken = m.storedTimetoken
			m.storedTimetoken = -1
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

// TODO: how to stop/reconnect?
func (m *SubscriptionManager) startSubscribeLoopWithRoutine() {
	m.log("loop")
	m.stopSubscribeLoop()

	combinedChannels := m.stateManager.prepareChannelList(true)
	combinedGroups := m.stateManager.prepareGroupList(true)

	if len(combinedChannels) == 0 && len(combinedGroups) == 0 {
		m.listenerManager.announceStatus(&PNStatus{
			Category: DisconnectedCategory,
		})
		m.log("no channels left to subscribe")
		// return
	}

	// TODO: invoke subscribe with context
	// TODO: fields should be local and not exposed to users
	m.RLock()
	tt := m.timetoken
	ctx := m.ctx
	m.RUnlock()

	opts := &SubscribeOpts{
		pubnub:    m.pubnub,
		Channels:  combinedChannels,
		Groups:    combinedGroups,
		Timetoken: tt,
		ctx:       ctx,
		// 	// transport
		// 	// config/subkey
		// 	// config/uuid
		// 	// config/timeouts
	}

	m.subscriptionLock.Lock()
	defer m.subscriptionLock.Unlock()

	// TODO: use context to be able to stop request
	res, err := executeRequest(opts)
	log.Println("following")
	if err != nil {
		if strings.Contains(err.Error(), "context canceled") {
			m.listenerManager.announceStatus(&PNStatus{
				Category: CancelledCategory,
			})
			// return
		}

		// log.Println(err.StatusCode)

		// TODO: handle timeout
		// TODO: handle canceled
		// TODO: handle error
		// return
	}

	m.RLock()
	announced := m.subscriptionStateAnnounced
	m.RUnlock()

	if announced == false {
		m.listenerManager.announceStatus(&PNStatus{
			Category: ConnectedCategory,
		})

		m.Lock()
		m.subscriptionStateAnnounced = true
		m.Unlock()
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

	// TODO: fetch messages and if any, push them to the worker queue
	if len(envelope.Messages) > 0 {
		for message := range envelope.Messages {
			m.messages <- envelope.Messages[message]
		}
	}

	m.Lock()
	if m.storedTimetoken != -1 {
		m.timetoken = m.storedTimetoken
		m.storedTimetoken = -1
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
	m.Unlock()

	go m.startSubscribeLoopWithRoutine()
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
	// go m.startSubscribeLoopWithRoutine()
}

func (m *SubscriptionManager) disconnect() {
	m.log("disconnect")

	// m.stopHeartbeat()
	m.stopSubscribeLoop()
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