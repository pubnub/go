package pubnub

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
)

type SubscriptionManager struct {
	lock sync.RWMutex

	listenerManager *ListenerManager
	stateManager    *StateManager
	pubnub          *PubNub

	ctx Context

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

	messages := make(chan interface{}, 1000)

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
	unsubscribeOperation UnsubscribeOperation) {
	m.stateManager.adaptUnsubscribeBuilder(unsubscribeOperation)

	m.subscriptionStateAnnounced = false

	//TODO:
	//Leave

	if m.stateManager.isEmpty() {
		m.region = 0
		m.storedTimetoken = -1
		m.timetoken = 0
	} else {
		m.storedTimetoken = m.timetoken
		m.timetoken = 0
	}
}

// TODO: how to stop/reconnect?
func (m *SubscriptionManager) startSubscribeLoop(messages chan<- interface{}) {
	for true {
		log.Println("loop")
		combinedChannels := m.stateManager.prepareChannelList(true)
		combinedGroups := m.stateManager.prepareGroupList(true)

		if len(combinedChannels) == 0 && len(combinedGroups) == 0 {
			log.Println("stop")
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
			Transport: m.pubnub.GetSubscribeClient().Transport,
			// 	// transport
			// 	// config/subkey
			// 	// config/uuid
			// 	// config/timeouts
		}

		// TODO: use context to be able to stop request
		res, err := executeRequest(opts)
		if err != nil {
			// TODO: handle timeout
			// TODO: handle canceled
			// TODO: handle error
			return
		}

		if m.subscriptionStateAnnounced == false {
			// TODO: announce connect status event
		}

		var envelope subscribeEnvelope
		err = json.Unmarshal(res, &envelope)
		if err != nil {
			// TODO: send error to status
		}
		fmt.Printf("parsed: %#v\n", envelope)
		// TODO: fetch messages and if any, push them to the worker queue

		if len(envelope.Messages) > 0 {
			for message := range envelope.Messages {
				messages <- message
			}
		}

		if m.storedTimetoken != -1 {
			m.timetoken = m.storedTimetoken
			m.storedTimetoken = -1
		} else {
			tt, err := strconv.ParseInt(envelope.Metadata.Timetoken, 10, 64)
			if err != nil {
				// TODO: error
				log.Panicln("nil timetoken", envelope.Metadata)
			}

			m.timetoken = tt
		}

		m.region = envelope.Metadata.Region
	}
}

type subscribeEnvelope struct {
	Messages []interface{} `json:"m"`
	Metadata struct {
		Timetoken string `json:"t"`
		Region    int8   `json:"r"`
	} `json:"t"`
}

func subscribeMessageWorker(lm *ListenerManager, messages <-chan interface{}) {
	for message := range messages {
		fmt.Println(">>>", message)
		// TODO: parse
	}
}

func (m *SubscriptionManager) AddListener(listener *Listener) {
	m.listenerManager.addListener(listener)
}

func reconnect() {

}
