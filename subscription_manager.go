package pubnub

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
)

type SubscriptionManager struct {
	lock sync.RWMutex

	listenerManager *ListenerManager
	stateManager    *StateManager
	pubnub          *PubNub

	ctx             Context
	subscribeCancel func()

	// Store the latest timetoken to subscribe with, null by default to get the
	// latest timetoken.
	timetoken int64

	// When changing the channel mix, store the timetoken for a later date
	storedTimetoken int64

	region int8

	subscriptionStateAnnounced bool
}

type SubscribeOperation struct {
	Channels        []string
	ChannelGroups   []string
	PresenceEnabled bool
	Timetoken       int64
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

	manager.timetoken = 0
	manager.storedTimetoken = -1

	manager.subscriptionStateAnnounced = false

	manager.ctx, manager.subscribeCancel = context.WithCancel(context.Background())

	messages := make(chan subscribeMessage, 1000)

	go manager.startSubscribeLoop(messages)
	go subscribeMessageWorker(manager.listenerManager, messages)

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
	m.subscriptionStateAnnounced = false

	if subscribeOperation.Timetoken != 0 {
		m.timetoken = subscribeOperation.Timetoken
	}

	if m.timetoken != 0 {
		m.storedTimetoken = m.timetoken
	}

	m.timetoken = 0
}

func (m *SubscriptionManager) adaptUnsubscribe(
	unsubscribeOperation *UnsubscribeOperation) {
	m.stateManager.adaptUnsubscribeBuilder(unsubscribeOperation)

	m.subscriptionStateAnnounced = false

	err := m.pubnub.Leave(&LeaveOpts{
		Channels:      unsubscribeOperation.Channels,
		ChannelGroups: unsubscribeOperation.ChannelGroups,
	})

	if err != nil {
		m.listenerManager.announceStatus(&PNStatus{
			Category:              BadRequestCategory,
			ErrorData:             err,
			Error:                 true,
			Operation:             PNUnsubscribeOperation,
			AffectedChannels:      m.Channels,
			AffectedChannelGroups: m.ChannelGroups,
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

	if m.stateManager.isEmpty() {
		m.region = 0
		m.storedTimetoken = -1
		m.timetoken = 0
	} else {
		m.storedTimetoken = m.timetoken
		m.timetoken = 0
	}
	m.reconnect()
}

// TODO: how to stop/reconnect?
func (m *SubscriptionManager) startSubscribeLoop(messages chan<- subscribeMessage) {
	for true {
		combinedChannels := m.stateManager.prepareChannelList(true)
		combinedGroups := m.stateManager.prepareGroupList(true)

		if len(combinedChannels) == 0 && len(combinedGroups) == 0 {
			m.listenerManager.announceStatus(&PNStatus{
				Category: DisconnectedCategory,
			})
			break
		}

		// TODO: invoke subscribe with context
		// TODO: fields should be local and not exposed to users
		opts := &SubscribeOpts{
			pubnub:    m.pubnub,
			Channels:  combinedChannels,
			Groups:    combinedGroups,
			Timetoken: m.timetoken,
			ctx:       m.ctx,
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

		if m.subscriptionStateAnnounced == false {
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
				AffectedChannels:      m.Channels,
				AffectedChannelGroups: m.ChannelGroups,
			})
		}

		// TODO: fetch messages and if any, push them to the worker queue
		if len(envelope.Messages) > 0 {
			for message := range envelope.Messages {
				messages <- envelope.Messages[message]
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
					AffectedChannels:      m.Channels,
					AffectedChannelGroups: m.ChannelGroups,
				})
			}

			m.timetoken = tt
		}

		m.region = envelope.Metadata.Region
	}
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

	if strings.Contains(channel, "-pnpres") {
		//TODO: presencePayload
	} else {
		actualCh := ""
		subscribedCh := channel

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
			Timetoken:         publishMetadata.PublishTimetoken,
			Publisher:         payload.IssuingClientId,
			UserMetadata:      payload.UserMetadata,
		}

		lm.announceMessage(pnMessageResult)
	}
}

func (m *SubscriptionManager) AddListener(listener *Listener) {
	m.listenerManager.addListener(listener)
}

func (m *SubscriptionManager) reconnect() {
	if m.ctx != nil && m.subscribeCancel != nil {
		m.subscribeCancel()
		m.ctx, m.subscribeCancel = context.WithCancel(context.Background())
	}
}
