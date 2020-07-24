package pubnub

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pubnub/go/utils"
)

// SubscriptionManager Events:
// - ConnectedCategory - after connection established
// - DisconnectedCategory - after subscription loop stops for any reason (no
// channels left or error happened)
// Unsubscribe
// When you unsubscribe from channel or channel group the following events
// happens:
// - LoopStopCategory - immediately when no more channels or channel groups left
// to subscribe
// - PNUnsubscribeOperation - after leave request was fulfilled and server is
// notified about unsubscibed items
// Announcement:
// Status, Message and Presence announcement happens in a distinct goroutine.
// It doesn't block subscribe loop.
// Keep in mind that each listener will receive the same pointer to a response
// object. You may wish to create a shallow copy of either the response or the
// response message by you own to not affect the other listeners.
// Heartbeat:
// - Heartbeat is enabled by default.
// - Default presence timeout is 0 seconds.
// - The first Heartbeat request will be scheduled to be executed after
// getHeartbeatInterval() seconds (default - 149).
type SubscriptionManager struct {
	sync.RWMutex
	subscriptionLock    sync.Mutex
	hbDataMutex         sync.RWMutex
	listenerManager     *ListenerManager
	stateManager        *StateManager
	pubnub              *PubNub
	reconnectionManager *ReconnectionManager
	transport           http.RoundTripper
	messages            chan subscribeMessage
	ctx                 Context
	subscribeCancel     func()
	heartbeatCancel     func()

	// Store the latest timetoken to subscribe with, null by default to get the
	// latest timetoken.
	timetoken int64

	// When changing the channel mix, store the timetoken for a later date
	storedTimetoken int64

	region int8

	subscriptionStateAnnounced   bool
	heartbeatStopCalled          bool
	exitSubscriptionManagerMutex sync.RWMutex
	exitSubscriptionManager      chan bool
	queryParam                   map[string]string
	channelsOpen                 bool
	requestSentAt                int64
}

// SubscribeOperation is the type to store the subscribe op params
type SubscribeOperation struct {
	Channels         []string
	ChannelGroups    []string
	PresenceEnabled  bool
	Timetoken        int64
	FilterExpression string
	State            map[string]interface{}
	QueryParam       map[string]string
}

// UnsubscribeOperation is the types to store unsubscribe op params
type UnsubscribeOperation struct {
	Channels      []string
	ChannelGroups []string
	QueryParam    map[string]string
}

// StateOperation is the types to store state op params
type StateOperation struct {
	channels      []string
	channelGroups []string
	state         map[string]interface{}
}

func newSubscriptionManager(pubnub *PubNub, ctx Context) *SubscriptionManager {
	manager := &SubscriptionManager{}

	manager.pubnub = pubnub

	manager.listenerManager = newListenerManager(ctx, pubnub)
	manager.stateManager = newStateManager()

	manager.Lock()
	manager.timetoken = 0
	manager.storedTimetoken = -1
	manager.subscriptionStateAnnounced = false
	manager.ctx, manager.subscribeCancel = contextWithCancel(backgroundContext)
	manager.messages = make(chan subscribeMessage, 1000)
	manager.reconnectionManager = newReconnectionManager(pubnub)
	manager.channelsOpen = true
	manager.Unlock()

	if manager.pubnub.Config.PNReconnectionPolicy != PNNonePolicy {

		manager.reconnectionManager.HandleReconnection(func() {
			go manager.reconnect()

			manager.Lock()
			manager.subscriptionStateAnnounced = true
			manager.Unlock()
			combinedChannels := manager.stateManager.prepareChannelList(true)
			combinedGroups := manager.stateManager.prepareGroupList(true)

			pnStatus := &PNStatus{
				Error:                 false,
				AffectedChannels:      combinedChannels,
				AffectedChannelGroups: combinedGroups,
				Category:              PNReconnectedCategory,
			}

			pubnub.Config.Log.Println("Status: ", pnStatus)

			manager.listenerManager.announceStatus(pnStatus)
		})
	}

	manager.reconnectionManager.HandleOnMaxReconnectionExhaustion(func() {
		combinedChannels := manager.stateManager.prepareChannelList(true)
		combinedGroups := manager.stateManager.prepareGroupList(true)

		pnStatus := &PNStatus{
			Error:                 false,
			AffectedChannels:      combinedChannels,
			AffectedChannelGroups: combinedGroups,
			Category:              PNReconnectionAttemptsExhausted,
		}
		pubnub.Config.Log.Println("Status: ", pnStatus)

		manager.listenerManager.announceStatus(pnStatus)

		manager.Disconnect()
	})

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

// Destroy closes the subscription manager, listeners and reconnection manager instances.
func (m *SubscriptionManager) Destroy() {
	if m.subscribeCancel != nil {
		m.subscribeCancel()
	}
	if m.channelsOpen {
		m.RLock()
		m.channelsOpen = false
		m.RUnlock()
		m.exitSubscriptionManagerMutex.RLock()
		if m.exitSubscriptionManager != nil {
			close(m.exitSubscriptionManager)
		}
		m.exitSubscriptionManagerMutex.RUnlock()
		if m.listenerManager.exitListener != nil {
			close(m.listenerManager.exitListener)
		}
		if m.listenerManager.exitListenerAnnounce != nil {
			close(m.listenerManager.exitListenerAnnounce)
		}
		if m.reconnectionManager.exitReconnectionManager != nil {
			m.reconnectionManager.stopHeartbeatTimer()
			close(m.reconnectionManager.exitReconnectionManager)
		}

	}

}

func (m *SubscriptionManager) adaptState(stateOperation StateOperation) {
	m.stateManager.adaptStateOperation(stateOperation)
}

func (m *SubscriptionManager) adaptSubscribe(
	subscribeOperation *SubscribeOperation) {
	m.stateManager.adaptSubscribeOperation(subscribeOperation)
	m.pubnub.Config.Log.Println("adapting a new subscription", subscribeOperation.Channels,
		subscribeOperation.PresenceEnabled)

	m.Lock()

	m.subscriptionStateAnnounced = false
	m.queryParam = subscribeOperation.QueryParam

	if subscribeOperation.Timetoken != 0 {
		m.timetoken = subscribeOperation.Timetoken
	}

	if m.timetoken != 0 {
		m.storedTimetoken = m.timetoken
	}

	m.timetoken = 0

	m.Unlock()

	m.reconnect()
}

func (m *SubscriptionManager) adaptUnsubscribe(
	unsubscribeOperation *UnsubscribeOperation) {
	m.pubnub.Config.Log.Println("before adaptUnsubscribeOperation")
	m.stateManager.adaptUnsubscribeOperation(unsubscribeOperation)
	m.pubnub.Config.Log.Println("after adaptUnsubscribeOperation")

	m.Lock()
	m.subscriptionStateAnnounced = false
	m.Unlock()

	go func() {
		announceAck := false
		if !m.pubnub.Config.SuppressLeaveEvents {
			_, err := m.pubnub.Leave().Channels(unsubscribeOperation.Channels).
				ChannelGroups(unsubscribeOperation.ChannelGroups).QueryParam(unsubscribeOperation.QueryParam).Execute()

			if err != nil {
				pnStatus := &PNStatus{
					Category:              PNBadRequestCategory,
					ErrorData:             err,
					Error:                 true,
					Operation:             PNUnsubscribeOperation,
					AffectedChannels:      unsubscribeOperation.Channels,
					AffectedChannelGroups: unsubscribeOperation.ChannelGroups,
				}
				m.pubnub.Config.Log.Println("Leave: err", err, pnStatus)
				m.listenerManager.announceStatus(pnStatus)
			} else {
				announceAck = true
			}
		} else {
			announceAck = true
		}

		if announceAck {
			pnStatus := &PNStatus{
				Category:              PNAcknowledgmentCategory,
				StatusCode:            200,
				Operation:             PNUnsubscribeOperation,
				UUID:                  m.pubnub.Config.UUID,
				AffectedChannels:      unsubscribeOperation.Channels,
				AffectedChannelGroups: unsubscribeOperation.ChannelGroups,
			}
			m.pubnub.Config.Log.Println("Leave: ack", pnStatus)
			m.listenerManager.announceStatus(pnStatus)
			m.pubnub.Config.Log.Println("After Leave: ack", pnStatus)
		}
	}()
	m.pubnub.Config.Log.Println("before storedTimetoken reset")
	m.Lock()
	if m.stateManager.isEmpty() {
		m.region = 0
		m.storedTimetoken = -1
		m.timetoken = 0
	} else {
		m.storedTimetoken = m.timetoken
		m.timetoken = 0
	}
	m.Unlock()
	m.pubnub.Config.Log.Println("after storedTimetoken reset")

	m.reconnect()
	m.pubnub.Config.Log.Println("after reconnect")
}

func (m *SubscriptionManager) startSubscribeLoop() {
	m.pubnub.Config.Log.Println("startSubscribeLoop")
	go subscribeMessageWorker(m)

	go m.reconnectionManager.startPolling()

	for {
		m.pubnub.Config.Log.Println("startSubscribeLoop looping...")
		combinedChannels := m.stateManager.prepareChannelList(true)
		combinedGroups := m.stateManager.prepareGroupList(true)

		if len(combinedChannels) == 0 && len(combinedGroups) == 0 {
			m.pubnub.Config.Log.Println("no channels left to subscribe")
			m.listenerManager.announceStatus(&PNStatus{
				Category: PNDisconnectedCategory,
			})

			m.reconnectionManager.stopHeartbeatTimer()

			break
		}

		m.Lock()
		tt := m.timetoken
		ctx := m.ctx
		m.Unlock()

		opts := &subscribeOpts{
			pubnub:           m.pubnub,
			Channels:         combinedChannels,
			ChannelGroups:    combinedGroups,
			Timetoken:        tt,
			Heartbeat:        m.pubnub.Config.PresenceTimeout,
			FilterExpression: m.pubnub.Config.FilterExpression,
			ctx:              ctx,
			QueryParam:       m.queryParam,
		}

		if s := m.stateManager.createStatePayload(); len(s) > 0 {
			opts.State = s
		}
		m.hbDataMutex.Lock()
		m.requestSentAt = time.Now().Unix()
		m.hbDataMutex.Unlock()

		res, _, err := executeRequest(opts)
		if err != nil {
			m.pubnub.Config.Log.Println(err.Error())

			if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "request canceled") {
				m.listenerManager.announceStatus(&PNStatus{
					Category: PNTimeoutCategory,
				})
				m.pubnub.Config.Log.Println("continue")
				continue
			} else {

				if strings.Contains(err.Error(), "context canceled") {
					pnStatus := &PNStatus{
						Category: PNCancelledCategory,
					}
					m.pubnub.Config.Log.Println("Status:", pnStatus)
					m.listenerManager.announceStatus(pnStatus)
					m.pubnub.Config.Log.Println("context canceled")
					break
				} else if strings.Contains(err.Error(), "Forbidden") ||
					strings.Contains(err.Error(), "403") {
					pnStatus := &PNStatus{
						Category: PNAccessDeniedCategory,
					}
					m.pubnub.Config.Log.Println("Status:", pnStatus)
					m.listenerManager.announceStatus(pnStatus)
					m.unsubscribeAll()
					break
				} else if strings.Contains(err.Error(), "400") ||
					strings.Contains(err.Error(), "Bad Request") {
					pnStatus := &PNStatus{
						Category: PNBadRequestCategory,
					}
					m.pubnub.Config.Log.Println("Status:", pnStatus)
					m.listenerManager.announceStatus(pnStatus)
					m.unsubscribeAll()
					break
				} else if strings.Contains(err.Error(), "530") || strings.Contains(err.Error(), "No Stub Matched") {
					pnStatus := &PNStatus{
						Category: PNNoStubMatchedCategory,
					}
					m.pubnub.Config.Log.Println("Status:", pnStatus)
					m.listenerManager.announceStatus(pnStatus)
					m.unsubscribeAll()
					break
				} else {
					pnStatus := &PNStatus{
						Category: PNUnknownCategory,
					}
					m.pubnub.Config.Log.Println("Status:", pnStatus)
					m.listenerManager.announceStatus(pnStatus)

					break
				}
			}

		}

		m.Lock()
		announced := m.subscriptionStateAnnounced

		if announced == false {

			m.listenerManager.announceStatus(&PNStatus{
				Category: PNConnectedCategory,
			})
			m.subscriptionStateAnnounced = true
		}
		m.Unlock()

		var envelope subscribeEnvelope
		err = json.Unmarshal(res, &envelope)
		if err != nil {
			pnStatus := &PNStatus{
				Category:              PNBadRequestCategory,
				ErrorData:             err,
				Error:                 true,
				Operation:             PNSubscribeOperation,
				AffectedChannels:      combinedChannels,
				AffectedChannelGroups: combinedGroups,
			}
			m.pubnub.Config.Log.Println("Unmarshal: err", err, pnStatus)

			m.listenerManager.announceStatus(pnStatus)
		}
		messageCount := len(envelope.Messages)
		if messageCount > 0 {
			if messageCount > m.pubnub.Config.MessageQueueOverflowCount {
				pnStatus := &PNStatus{
					Error:                 false,
					AffectedChannels:      combinedChannels,
					AffectedChannelGroups: combinedGroups,
					Category:              PNRequestMessageCountExceededCategory,
				}
				m.pubnub.Config.Log.Println("Status: ", pnStatus)

				m.listenerManager.announceStatus(pnStatus)
			}
			for _, message := range envelope.Messages {
				m.messages <- message
			}
		}

		m.Lock()
		if m.storedTimetoken != -1 {

			m.timetoken = m.storedTimetoken
			m.storedTimetoken = -1
		} else {
			tt, err := strconv.ParseInt(envelope.Metadata.Timetoken, 10, 64)
			if err != nil {

				pnStatus := &PNStatus{
					Category:              PNBadRequestCategory,
					ErrorData:             err,
					Error:                 true,
					Operation:             PNSubscribeOperation,
					AffectedChannels:      combinedChannels,
					AffectedChannelGroups: combinedGroups,
				}
				m.pubnub.Config.Log.Println("ParseInt: err", err, pnStatus)
				m.listenerManager.announceStatus(pnStatus)
			}

			m.timetoken = tt
		}

		m.region = envelope.Metadata.Region
		m.Unlock()
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
	Shard             string        `json:"a"`
	SubscriptionMatch string        `json:"b"`
	Channel           string        `json:"c"`
	IssuingClientID   string        `json:"i"`
	SubscribeKey      string        `json:"k"`
	Flags             int           `json:"f"`
	Payload           interface{}   `json:"d"`
	UserMetadata      interface{}   `json:"u"`
	MessageType       PNMessageType `json:"e"`
	SequenceNumber    int           `json:"s"`

	PublishMetaData publishMetadata `json:"p"`
}

type presenceEnvelope struct {
	Action    string
	UUID      string
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

func subscribeMessageWorker(m *SubscriptionManager) {
	m.Lock()
	if m.ctx == nil && m.subscribeCancel == nil {
		m.pubnub.Config.Log.Println("subscribeMessageWorker setting context")
		m.ctx, m.subscribeCancel = contextWithCancel(backgroundContext)
		m.pubnub.Config.Log.Println("subscribeMessageWorker after setting context")
	}

	m.pubnub.Config.Log.Println("subscribeMessageWorker")

	m.Unlock()
	m.pubnub.Config.Log.Println("acquiring lock exitSubscriptionManagerMutex")
	m.exitSubscriptionManagerMutex.Lock()
	if m.exitSubscriptionManager != nil {
		m.exitSubscriptionManager <- true
		m.pubnub.Config.Log.Println("close exitSubscriptionManager")
	}
	m.pubnub.Config.Log.Println("make channel exitSubscriptionManager")
	m.exitSubscriptionManager = make(chan bool)
	m.exitSubscriptionManagerMutex.Unlock()

SubscribeMessageWorkerLabel:
	for {
		m.pubnub.Config.Log.Println("subscribeMessageWorker looping...")
		combinedChannels := m.stateManager.prepareChannelList(true)
		combinedGroups := m.stateManager.prepareGroupList(true)

		if len(combinedChannels) == 0 && len(combinedGroups) == 0 {
			m.pubnub.Config.Log.Println("subscribeMessageWorker all channels unsubscribed")
			break
		}
		select {
		case <-m.exitSubscriptionManager:
			m.pubnub.Config.Log.Println("subscribeMessageWorker context done")
			break SubscribeMessageWorkerLabel
		case message := <-m.messages:
			m.pubnub.Config.Log.Println("subscribeMessageWorker messages")
			processSubscribePayload(m, message)
		}
	}
	m.pubnub.Config.Log.Println("subscribeMessageWorker after for")

}

func processPresencePayload(m *SubscriptionManager, payload subscribeMessage, channel, subscriptionMatch string, publishMeta publishMetadata) {
	var presencePayload map[string]interface{}
	var action, uuid, actualChannel, subscribedChannel string
	var occupancy int
	var timestamp int64
	var data interface{}
	var ok, hereNowRefresh bool

	if presencePayload, ok = payload.Payload.(map[string]interface{}); !ok {
		m.listenerManager.announceStatus(&PNStatus{
			Category:         PNUnknownCategory,
			ErrorData:        errors.New("Presence response parsing error"),
			Error:            true,
			Operation:        PNSubscribeOperation,
			AffectedChannels: []string{channel},
		})
	}

	action, _ = presencePayload["action"].(string)
	uuid, _ = presencePayload["uuid"].(string)
	occupancy, _ = presencePayload["occupancy"].(int)
	if presencePayload["timestamp"] != nil {
		m.pubnub.Config.Log.Println("presencePayload['timestamp'] type", reflect.TypeOf(presencePayload["timestamp"]).Kind())
		switch presencePayload["timestamp"].(type) {
		case int:
			timestamp = int64(presencePayload["timestamp"].(int))
			break
		case int64:
			timestamp = presencePayload["timestamp"].(int64)
			break
		case float64:
			timestamp = int64(presencePayload["timestamp"].(float64))
			break
		}

	}

	data = presencePayload["data"]
	if presencePayload["here_now_refresh"] != nil {
		hereNowRefresh = presencePayload["here_now_refresh"].(bool)
	}
	timetoken, _ := strconv.ParseInt(publishMeta.PublishTimetoken, 10, 64)

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
		UUID:              uuid,
		Timestamp:         timestamp,
		HereNowRefresh:    hereNowRefresh,
	}
	m.listenerManager.announcePresence(pnPresenceResult)
}

func processNonPresencePayload(m *SubscriptionManager, payload subscribeMessage, channel, subscriptionMatch string, publishMeta publishMetadata) {
	actualCh := ""
	subscribedCh := channel
	timetoken, _ := strconv.ParseInt(publishMeta.PublishTimetoken, 10, 64)

	if subscriptionMatch != "" {
		actualCh = channel
		subscribedCh = subscriptionMatch
	}
	var messagePayload interface{}

	switch payload.MessageType {
	case PNMessageTypeSignal:
		pnMessageResult := createPNMessageResult(payload.Payload, actualCh, subscribedCh, channel, subscriptionMatch, payload.IssuingClientID, payload.UserMetadata, timetoken)
		m.pubnub.Config.Log.Println("announceSignal,", pnMessageResult)
		m.listenerManager.announceSignal(pnMessageResult)
	case PNMessageTypeObjects:
		pnUUIDEvent, pnChannelEvent, pnMembershipEvent, eventType := createPNObjectsResult(payload.Payload, m, actualCh, subscribedCh, channel, subscriptionMatch)
		m.pubnub.Config.Log.Println("announceObjects,", pnUUIDEvent, pnChannelEvent, pnMembershipEvent, eventType)
		switch eventType {
		case PNObjectsUUIDEvent:
			m.pubnub.Config.Log.Println("pnUUIDEvent:", pnUUIDEvent)
			m.listenerManager.announceUUIDEvent(pnUUIDEvent)
		case PNObjectsChannelEvent:
			m.pubnub.Config.Log.Println("pnChannelEvent:", pnChannelEvent)
			m.listenerManager.announceChannelEvent(pnChannelEvent)
		case PNObjectsMembershipEvent:
			m.pubnub.Config.Log.Println("pnMembershipEvent:", pnMembershipEvent)
			m.listenerManager.announceMembershipEvent(pnMembershipEvent)
		}
	case PNMessageTypeMessageActions:
		pnMessageActionsEvent := createPNMessageActionsEventResult(payload.Payload, m, actualCh, subscribedCh, channel, subscriptionMatch, payload.IssuingClientID)
		m.pubnub.Config.Log.Println("PNMessageTypeMessageActions:", pnMessageActionsEvent)
		m.listenerManager.announceMessageActionsEvent(pnMessageActionsEvent)
	case PNMessageTypeFile:
		var err error
		messagePayload, err = parseCipherInterface(payload.Payload, m.pubnub.Config)
		if err != nil {
			pnStatus := &PNStatus{
				Category:         PNBadRequestCategory,
				ErrorData:        err,
				Error:            true,
				Operation:        PNSubscribeOperation,
				AffectedChannels: []string{channel},
			}
			m.pubnub.Config.Log.Println("DecryptString: err", err, pnStatus)
			m.listenerManager.announceStatus(pnStatus)

		}

		pnFilesEvent := createPNFilesEvent(messagePayload, m, actualCh, subscribedCh, channel, subscriptionMatch, payload.IssuingClientID, payload.UserMetadata, timetoken)
		m.pubnub.Config.Log.Println("PNMessageTypeFile:", PNMessageTypeFile)
		m.listenerManager.announceFile(pnFilesEvent)
	default:
		var err error
		messagePayload, err = parseCipherInterface(payload.Payload, m.pubnub.Config)
		if err != nil {
			pnStatus := &PNStatus{
				Category:         PNBadRequestCategory,
				ErrorData:        err,
				Error:            true,
				Operation:        PNSubscribeOperation,
				AffectedChannels: []string{channel},
			}
			m.pubnub.Config.Log.Println("DecryptString: err", err, pnStatus)
			m.listenerManager.announceStatus(pnStatus)

		}
		pnMessageResult := createPNMessageResult(messagePayload, actualCh, subscribedCh, channel, subscriptionMatch, payload.IssuingClientID, payload.UserMetadata, timetoken)
		m.pubnub.Config.Log.Println("announceMessage,", pnMessageResult)
		m.listenerManager.announceMessage(pnMessageResult)
	}
	m.pubnub.Config.Log.Println("after announceMessage")
}

func processSubscribePayload(m *SubscriptionManager, payload subscribeMessage) {
	channel := payload.Channel
	subscriptionMatch := payload.SubscriptionMatch
	publishMetadata := payload.PublishMetaData

	if channel != "" && channel == subscriptionMatch {
		subscriptionMatch = ""
	}

	if strings.Contains(payload.Channel, "-pnpres") {
		processPresencePayload(m, payload, channel, subscriptionMatch, publishMetadata)
	} else {
		processNonPresencePayload(m, payload, channel, subscriptionMatch, publishMetadata)
	}
}

func createPNFilesEvent(filePayload interface{}, m *SubscriptionManager, actualCh, subscribedCh, channel, subscriptionMatch, issuingClientID string, userMetadata interface{}, timetoken int64) *PNFilesEvent {
	var filesPayload map[string]interface{}
	var ok bool
	if filesPayload, ok = filePayload.(map[string]interface{}); !ok {
		m.listenerManager.announceStatus(&PNStatus{
			Category:         PNUnknownCategory,
			ErrorData:        errors.New("Files response parsing error"),
			Error:            true,
			Operation:        PNSubscribeOperation,
			AffectedChannels: []string{channel},
		})
		return nil
	}

	resp := PNFileMessageAndDetails{}
	resp.PNFile, resp.PNMessage = ParseFileInfo(filesPayload)
	resGetFile, _, _ := m.pubnub.GetFileURL().Channel(channel).ID(resp.PNFile.ID).Name(resp.PNFile.Name).Execute()

	if resGetFile != nil {
		resp.PNFile.URL = resGetFile.URL
	}

	pnFilesEvent := &PNFilesEvent{
		File:              resp,
		ActualChannel:     actualCh,
		SubscribedChannel: subscribedCh,
		Channel:           channel,
		Subscription:      subscriptionMatch,
		Timetoken:         timetoken,
		Publisher:         issuingClientID,
		UserMetadata:      userMetadata,
	}
	return pnFilesEvent
}

func createPNMessageActionsEventResult(maPayload interface{}, m *SubscriptionManager, actualCh, subscribedCh, channel, subscriptionMatch, issuingClientID string) *PNMessageActionsEvent {
	var messageActionsPayload map[string]interface{}
	var ok bool
	if messageActionsPayload, ok = maPayload.(map[string]interface{}); !ok {
		m.listenerManager.announceStatus(&PNStatus{
			Category:         PNUnknownCategory,
			ErrorData:        errors.New("Message Actions response parsing error"),
			Error:            true,
			Operation:        PNSubscribeOperation,
			AffectedChannels: []string{channel},
		})
		return nil
	}
	eventType := PNMessageActionsEventType(messageActionsPayload["event"].(string))
	var data map[string]interface{}
	resp := PNMessageActionsResponse{}

	if o, ok := messageActionsPayload["data"]; ok {
		data = o.(map[string]interface{})
		if d, ok := data["type"]; ok {
			resp.ActionType = d.(string)
		}
		if d, ok := data["value"]; ok {
			resp.ActionValue = d.(string)
		}
		if d, ok := data["actionTimetoken"]; ok {
			resp.ActionTimetoken = d.(string)
		}
		if d, ok := data["messageTimetoken"]; ok {
			resp.MessageTimetoken = d.(string)
		}
		resp.UUID = issuingClientID
	}

	pnMessageActionsEvent := &PNMessageActionsEvent{
		Event:             eventType,
		Data:              resp,
		ActualChannel:     actualCh,
		SubscribedChannel: subscribedCh,
		Channel:           channel,
		Subscription:      subscriptionMatch,
	}

	return pnMessageActionsEvent
}

func createPNObjectsResult(objPayload interface{}, m *SubscriptionManager, actualCh, subscribedCh, channel, subscriptionMatch string) (*PNUUIDEvent, *PNChannelEvent, *PNMembershipEvent, PNObjectsEventType) {
	var objectsPayload map[string]interface{}
	var ok bool
	if objectsPayload, ok = objPayload.(map[string]interface{}); !ok {
		m.listenerManager.announceStatus(&PNStatus{
			Category:         PNUnknownCategory,
			ErrorData:        errors.New("Objects response parsing error"),
			Error:            true,
			Operation:        PNSubscribeOperation,
			AffectedChannels: []string{channel},
		})
		return nil, nil, nil, PNObjectsNoneEvent
	}
	eventType := PNObjectsEventType(objectsPayload["type"].(string))
	event := PNObjectsEvent(objectsPayload["event"].(string))
	version := ""
	if d, ok := objectsPayload["version"]; ok {
		version = d.(string)
		if version == "1.0" {
			m.pubnub.Config.Log.Println("Ignoring version 1.0 event")
			return &PNUUIDEvent{}, &PNChannelEvent{}, &PNMembershipEvent{}, PNObjectsNoneEvent
		}
	} else {
		m.pubnub.Config.Log.Println("Ignoring non versioned event")
		return &PNUUIDEvent{}, &PNChannelEvent{}, &PNMembershipEvent{}, PNObjectsNoneEvent
	}
	var id, UUID, channelID, description, timestamp, updated, eTag, name, externalID, profileURL, email string
	var custom, data map[string]interface{}
	if o, ok := objectsPayload["data"]; ok {
		data = o.(map[string]interface{})
		if d, ok := data["uuid"]; ok {
			u := d.(map[string]interface{})
			UUID = u["id"].(string)
		}
		if d, ok := data["id"]; ok {
			id = d.(string)
		}
		if d, ok := data["channel"]; ok {
			ch := d.(map[string]interface{})
			channelID = ch["id"].(string)
		}
		if d, ok := data["name"]; ok {
			name = d.(string)
		}
		if d, ok := data["externalId"]; ok {
			externalID = d.(string)
		}
		if d, ok := data["profileUrl"]; ok {
			profileURL = d.(string)
		}
		if d, ok := data["email"]; ok {
			email = d.(string)
		}
		if d, ok := data["description"]; ok {
			description = d.(string)
		}
		if d, ok := data["timestamp"]; ok {
			timestamp = d.(string)
		}
		if d, ok := data["updated"]; ok {
			updated = d.(string)
		}
		if d, ok := data["eTag"]; ok {
			eTag = d.(string)
		}
		if d, ok := data["custom"]; ok {
			custom = d.(map[string]interface{})
		}

	}

	pnObjectsResult := &PNObjectsResponse{
		Event:       event,
		EventType:   eventType,
		ID:          id,
		Channel:     channel,
		Description: description,
		Timestamp:   timestamp,
		Updated:     updated,
		ETag:        eTag,
		Custom:      custom,
		Data:        data,
		Name:        name,
		ExternalID:  externalID,
		ProfileURL:  profileURL,
		Email:       email,
	}

	pnChannelEvent := &PNChannelEvent{
		Event:             pnObjectsResult.Event,
		ChannelID:         id,
		Description:       pnObjectsResult.Description,
		Timestamp:         pnObjectsResult.Timestamp,
		Name:              pnObjectsResult.Name,
		Updated:           pnObjectsResult.Updated,
		ETag:              pnObjectsResult.ETag,
		Custom:            pnObjectsResult.Custom,
		ActualChannel:     actualCh,
		SubscribedChannel: subscribedCh,
		Channel:           channel,
		Subscription:      subscriptionMatch,
	}

	pnUUIDEvent := &PNUUIDEvent{
		Event:             pnObjectsResult.Event,
		UUID:              id,
		Timestamp:         pnObjectsResult.Timestamp,
		Updated:           pnObjectsResult.Updated,
		ETag:              pnObjectsResult.ETag,
		Custom:            pnObjectsResult.Custom,
		Name:              pnObjectsResult.Name,
		ExternalID:        pnObjectsResult.ExternalID,
		ProfileURL:        pnObjectsResult.ProfileURL,
		Email:             pnObjectsResult.Email,
		ActualChannel:     actualCh,
		SubscribedChannel: subscribedCh,
		Channel:           channel,
		Subscription:      subscriptionMatch,
	}

	pnMembershipEvent := &PNMembershipEvent{
		Event:             pnObjectsResult.Event,
		UUID:              UUID,
		ChannelID:         channelID,
		Description:       pnObjectsResult.Description,
		Timestamp:         pnObjectsResult.Timestamp,
		Custom:            pnObjectsResult.Custom,
		ActualChannel:     actualCh,
		SubscribedChannel: subscribedCh,
		Channel:           pnObjectsResult.Channel,
		Subscription:      subscriptionMatch,
	}

	return pnUUIDEvent, pnChannelEvent, pnMembershipEvent, eventType
}

func createPNMessageResult(messagePayload interface{}, actualCh, subscribedCh, channel, subscriptionMatch, issuingClientID string, userMetadata interface{}, timetoken int64) *PNMessage {

	pnMessageResult := &PNMessage{
		Message:           messagePayload,
		ActualChannel:     actualCh,
		SubscribedChannel: subscribedCh,
		Channel:           channel,
		Subscription:      subscriptionMatch,
		Timetoken:         timetoken,
		Publisher:         issuingClientID,
		UserMetadata:      userMetadata,
	}

	return pnMessageResult

}

// parseCipherInterface handles the decryption in case a cipher key is used
// in case of error it returns data as is.
//
// parameters
// data: the data to decrypt as interface.
// cipherKey: cipher key to use to decrypt.
//
// returns the decrypted data as interface and error.
func parseCipherInterface(data interface{}, pnConf *Config) (interface{}, error) {
	if pnConf.CipherKey != "" {
		pnConf.Log.Println("reflect.TypeOf(data).Kind()", reflect.TypeOf(data).Kind(), data)
		switch v := data.(type) {
		case map[string]interface{}:

			if !pnConf.DisablePNOtherProcessing {
				//decrypt pn_other only
				msg, ok := v["pn_other"].(string)
				if ok {
					pnConf.Log.Println("v[pn_other]", v["pn_other"], v, msg)
					decrypted, errDecryption := utils.DecryptString(pnConf.CipherKey, msg, pnConf.UseRandomInitializationVector)
					if errDecryption != nil {
						pnConf.Log.Println(errDecryption, msg)
						return v, errDecryption
					} else {
						var intf interface{}
						err := json.Unmarshal([]byte(decrypted.(string)), &intf)
						if err != nil {
							pnConf.Log.Println("Unmarshal: err", err)
							return intf, err
						}
						v["pn_other"] = intf

						pnConf.Log.Println("reflect.TypeOf(v).Kind()", reflect.TypeOf(v).Kind(), v)
						return v, nil
					}
				}
				return v, nil
			}
			pnConf.Log.Println("return as is reflect.TypeOf(v).Kind()", reflect.TypeOf(v).Kind(), v)
			return v, nil
		case string:
			var intf interface{}
			decrypted, errDecryption := utils.DecryptString(pnConf.CipherKey, data.(string), pnConf.UseRandomInitializationVector)
			if errDecryption != nil {
				pnConf.Log.Println(errDecryption, intf)
				intf = data
				return intf, errDecryption
			}
			pnConf.Log.Println("reflect.TypeOf(intf).Kind()", reflect.TypeOf(decrypted).Kind(), decrypted)

			err := json.Unmarshal([]byte(decrypted.(string)), &intf)
			if err != nil {
				pnConf.Log.Println("Unmarshal: err", err)
				return intf, err
			}

			return intf, nil
		default:
			pnConf.Log.Println("returning as is", reflect.TypeOf(v).Kind())
			return v, nil
		}
	} else {
		pnConf.Log.Println("No Cipher, returning as is ", data)
		return data, nil
	}
}

// AddListener adds a new listener.
func (m *SubscriptionManager) AddListener(listener *Listener) {
	m.listenerManager.addListener(listener)
}

// RemoveListener removes the listener.
func (m *SubscriptionManager) RemoveListener(listener *Listener) {
	m.listenerManager.removeListener(listener)
}

// RemoveAllListeners removes all the listeners.
func (m *SubscriptionManager) RemoveAllListeners() {
	m.listenerManager.removeAllListeners()
}

// GetListeners gets all the listeners.
func (m *SubscriptionManager) GetListeners() map[*Listener]bool {
	listn := m.listenerManager.listeners
	return listn
}

func (m *SubscriptionManager) reconnect() {
	m.pubnub.Config.Log.Println("reconnect")
	m.reconnectionManager.stopHeartbeatTimer()
	m.pubnub.Config.Log.Println("after stopHeartbeatTimer")
	m.stopSubscribeLoop()

	combinedChannels := m.stateManager.prepareChannelList(true)
	combinedGroups := m.stateManager.prepareGroupList(true)

	if len(combinedChannels) == 0 && len(combinedGroups) == 0 {
		m.pubnub.Config.Log.Println("All channels or channel groups unsubscribed.")
	} else {
		go m.startSubscribeLoop()
		go m.pubnub.heartbeatManager.startHeartbeatTimer(false)
	}
}

// Disconnect stops all open subscribe requests, timers, heartbeats and unsubscribes from all channels
func (m *SubscriptionManager) Disconnect() {
	m.pubnub.Config.Log.Println("disconnect")

	if m.exitSubscriptionManager != nil {
		m.exitSubscriptionManager <- true
	}
	m.reconnectionManager.stopHeartbeatTimer()

	m.pubnub.heartbeatManager.stopHeartbeat(false, false)
	m.unsubscribeAll()
	m.stopSubscribeLoop()

}

func (m *SubscriptionManager) stopSubscribeLoop() {
	m.log("loop stop")

	if m.ctx != nil && m.subscribeCancel != nil {
		m.subscribeCancel()
		m.ctx = nil
		m.subscribeCancel = nil
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
		Channels:      m.stateManager.prepareChannelList(true),
		ChannelGroups: m.stateManager.prepareGroupList(true),
	})
}

func (m *SubscriptionManager) log(message string) {
	m.pubnub.Config.Log.Printf("pubnub: subscribe: %s: %s: %s/%s\n",
		message,
		m.pubnub.Config.UUID,
		m.stateManager.prepareChannelList(true),
		m.stateManager.prepareGroupList(true))
}
