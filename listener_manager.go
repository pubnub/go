package pubnub

import (
	"sync"
)

// Listener type has all the `types` of response events
type Listener struct {
	Status              chan *PNStatus
	Message             chan *PNMessage
	Presence            chan *PNPresence
	Signal              chan *PNMessage
	UUIDEvent           chan *PNUUIDEvent
	ChannelEvent        chan *PNChannelEvent
	MembershipEvent     chan *PNMembershipEvent
	MessageActionsEvent chan *PNMessageActionsEvent
	File                chan *PNFilesEvent
}

// NewListener initates the listener to facilitate the event handling
func NewListener() *Listener {
	return &Listener{
		Status:              make(chan *PNStatus),
		Message:             make(chan *PNMessage),
		Presence:            make(chan *PNPresence),
		Signal:              make(chan *PNMessage),
		UUIDEvent:           make(chan *PNUUIDEvent),
		ChannelEvent:        make(chan *PNChannelEvent),
		MembershipEvent:     make(chan *PNMembershipEvent),
		MessageActionsEvent: make(chan *PNMessageActionsEvent),
		File:                make(chan *PNFilesEvent),
	}
}

// ListenerManager is used in the internal handling of listeners.
type ListenerManager struct {
	sync.RWMutex
	ctx                  Context
	listeners            map[*Listener]bool
	exitListener         chan bool
	exitListenerAnnounce chan bool
	pubnub               *PubNub
}

func newListenerManager(ctx Context, pn *PubNub) *ListenerManager {
	return &ListenerManager{
		listeners:            make(map[*Listener]bool, 2),
		ctx:                  ctx,
		exitListener:         make(chan bool),
		exitListenerAnnounce: make(chan bool),
		pubnub:               pn,
	}
}

func (m *ListenerManager) addListener(listener *Listener) {
	m.Lock()

	m.listeners[listener] = true
	m.Unlock()
}

func (m *ListenerManager) removeListener(listener *Listener) {
	m.pubnub.Config.Log.Println("before removeListener")
	m.Lock()
	m.pubnub.Config.Log.Println("in removeListener lock")
	delete(m.listeners, listener)
	m.Unlock()
	m.pubnub.Config.Log.Println("after removeListener")
}

func (m *ListenerManager) removeAllListeners() {
	m.pubnub.Config.Log.Println("in removeAllListeners")
	m.Lock()
	lis := m.listeners
	for l := range lis {
		delete(m.listeners, l)
	}
	m.Unlock()
}

func (m *ListenerManager) copyListeners() map[*Listener]bool {
	m.Lock()
	lis := make(map[*Listener]bool)
	for k, v := range m.listeners {
		lis[k] = v
	}
	m.Unlock()
	return lis
}

func (m *ListenerManager) announceStatus(status *PNStatus) {
	go func() {
		lis := m.copyListeners()
	AnnounceStatusLabel:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.pubnub.Config.Log.Println("announceStatus exitListener")
				break AnnounceStatusLabel
			case l.Status <- status:
			}
		}
		m.pubnub.Config.Log.Println("announceStatus exit")
	}()
}

func (m *ListenerManager) announceMessage(message *PNMessage) {
	go func() {
		lis := m.copyListeners()
	AnnounceMessageLabel:
		for l := range lis {
			select {
			case <-m.exitListenerAnnounce:
				m.pubnub.Config.Log.Println("announceMessage exitListenerAnnounce")
				break AnnounceMessageLabel
			case l.Message <- message:
			}
		}

	}()
}

func (m *ListenerManager) announceSignal(message *PNMessage) {
	go func() {
		lis := m.copyListeners()

	AnnounceSignalLabel:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.pubnub.Config.Log.Println("announceSignal exitListener")
				break AnnounceSignalLabel

			case l.Signal <- message:
			}
		}
	}()
}

func (m *ListenerManager) announceUUIDEvent(message *PNUUIDEvent) {
	go func() {
		lis := m.copyListeners()

	AnnounceUUIDEventLabel:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.pubnub.Config.Log.Println("announceUUIDEvent exitListener")
				break AnnounceUUIDEventLabel

			case l.UUIDEvent <- message:
				m.pubnub.Config.Log.Println("l.UUIDEvent", message)
			}
		}
	}()
}

func (m *ListenerManager) announceChannelEvent(message *PNChannelEvent) {
	go func() {
		lis := m.copyListeners()

	AnnounceChannelEventLabel:
		for l := range lis {
			m.pubnub.Config.Log.Println("l.ChannelEvent", l)
			select {
			case <-m.exitListener:
				m.pubnub.Config.Log.Println("announceChannelEvent exitListener")
				break AnnounceChannelEventLabel

			case l.ChannelEvent <- message:
				m.pubnub.Config.Log.Println("l.ChannelEvent", message)
			}
		}
	}()
}

func (m *ListenerManager) announceMembershipEvent(message *PNMembershipEvent) {
	go func() {
		lis := m.copyListeners()

	AnnounceMembershipEvent:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.pubnub.Config.Log.Println("announceMembershipEvent exitListener")
				break AnnounceMembershipEvent

			case l.MembershipEvent <- message:
				m.pubnub.Config.Log.Println("l.MembershipEvent", message)
			}
		}
	}()
}

func (m *ListenerManager) announceMessageActionsEvent(message *PNMessageActionsEvent) {
	go func() {
		lis := m.copyListeners()

	AnnounceMessageActionsEvent:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.pubnub.Config.Log.Println("announceMessageActionsEvent exitListener")
				break AnnounceMessageActionsEvent

			case l.MessageActionsEvent <- message:
				m.pubnub.Config.Log.Println("l.MessageActionsEvent", message)
			}
		}
	}()
}

func (m *ListenerManager) announcePresence(presence *PNPresence) {
	go func() {
		lis := m.copyListeners()

	AnnouncePresenceLabel:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.pubnub.Config.Log.Println("announcePresence exitListener")
				break AnnouncePresenceLabel

			case l.Presence <- presence:
			}
		}
	}()
}

func (m *ListenerManager) announceFile(file *PNFilesEvent) {
	go func() {
		lis := m.copyListeners()

	AnnounceFileLabel:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.pubnub.Config.Log.Println("announceFile exitListener")
				break AnnounceFileLabel

			case l.File <- file:
			}
		}
	}()
}

// PNStatus is the status struct
type PNStatus struct {
	Category              StatusCategory
	Operation             OperationType
	ErrorData             error
	Error                 bool
	TLSEnabled            bool
	StatusCode            int
	UUID                  string
	AuthKey               string
	Origin                string
	ClientRequest         interface{} // Should be same for non-google environment
	AffectedChannels      []string
	AffectedChannelGroups []string
}

// PNMessage is the Message Response for Subscribe
type PNMessage struct {
	Message           interface{}
	UserMetadata      interface{}
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
	Publisher         string
	Timetoken         int64
}

// PNPresence is the Message Response for Presence
type PNPresence struct {
	Event             string
	UUID              string
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
	Occupancy         int
	Timetoken         int64
	Timestamp         int64
	UserMetadata      map[string]interface{}
	State             interface{}
	Join              []string
	Leave             []string
	Timeout           []string
	HereNowRefresh    bool
}

// PNUUIDEvent is the Response for an User Event
type PNUUIDEvent struct {
	Event             PNObjectsEvent
	UUID              string
	Description       string
	Timestamp         string
	Name              string
	ExternalID        string
	ProfileURL        string
	Email             string
	Updated           string
	ETag              string
	Custom            map[string]interface{}
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
}

// PNChannelEvent is the Response for a Space Event
type PNChannelEvent struct {
	Event             PNObjectsEvent
	ChannelID         string
	Description       string
	Timestamp         string
	Name              string
	Updated           string
	ETag              string
	Custom            map[string]interface{}
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
}

// PNMembershipEvent is the Response for a Membership Event
type PNMembershipEvent struct {
	Event             PNObjectsEvent
	UUID              string
	ChannelID         string
	Description       string
	Timestamp         string
	Custom            map[string]interface{}
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
}

// PNMessageActionsEvent is the Response for a Message Actions Event
type PNMessageActionsEvent struct {
	Event             PNMessageActionsEventType
	Data              PNMessageActionsResponse
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
}

// PNFilesEvent is the Response for a Files Event
type PNFilesEvent struct {
	File              PNFileMessageAndDetails
	UserMetadata      interface{}
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
	Publisher         string
	Timetoken         int64
}
