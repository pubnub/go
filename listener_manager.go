package pubnub

import (
	"sync"
)

//
type Listener struct {
	Status   chan *PNStatus
	Message  chan *PNMessage
	Presence chan *PNPresence
	Signal   chan *PNMessage
}

func NewListener() *Listener {
	return &Listener{
		Status:   make(chan *PNStatus),
		Message:  make(chan *PNMessage),
		Presence: make(chan *PNPresence),
		Signal:   make(chan *PNMessage),
	}
}

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
	delete(m.listeners, listener)
}

func (m *ListenerManager) removeAllListeners() {
	m.Lock()
	m.pubnub.Config.Log.Println("in removeAllListeners")
	for l := range m.listeners {
		delete(m.listeners, l)
	}
	m.Unlock()
}

func (m *ListenerManager) announceStatus(status *PNStatus) {
	go func() {
		m.RLock()
		m.pubnub.Config.Log.Println("announceStatus lock")
	AnnounceStatusLabel:
		for l := range m.listeners {
			select {
			case <-m.exitListener:
				m.pubnub.Config.Log.Println("announceStatus exitListener")
				break AnnounceStatusLabel
			case l.Status <- status:
			}
		}
		m.pubnub.Config.Log.Println("announceStatus unlock")
		m.RUnlock()
		m.pubnub.Config.Log.Println("announceStatus exit")
	}()
}

func (m *ListenerManager) announceMessage(message *PNMessage) {
	go func() {
		m.RLock()
	AnnounceMessageLabel:
		for l := range m.listeners {
			select {
			case <-m.exitListenerAnnounce:
				m.pubnub.Config.Log.Println("announceMessage exitListenerAnnounce")
				break AnnounceMessageLabel
			case l.Message <- message:
			}
		}
		m.RUnlock()
	}()
}

func (m *ListenerManager) announceSignal(message *PNMessage) {
	go func() {
		m.RLock()
	AnnounceSignalLabel:
		for l := range m.listeners {
			select {
			case <-m.exitListener:
				m.pubnub.Config.Log.Println("announceSignal exitListener")
				break AnnounceSignalLabel

			case l.Signal <- message:
			}
		}
		m.RUnlock()
	}()
}

func (m *ListenerManager) announcePresence(presence *PNPresence) {
	go func() {
		m.RLock()
	AnnouncePresenceLabel:
		for l := range m.listeners {
			select {
			case <-m.exitListener:
				m.pubnub.Config.Log.Println("announcePresence exitListener")
				break AnnouncePresenceLabel

			case l.Presence <- presence:
			}
		}
		m.RUnlock()
	}()
}

//
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
