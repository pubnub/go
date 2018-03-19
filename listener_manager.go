package pubnub

import (
	"sync"
)

type Listener struct {
	Status   chan *PNStatus
	Message  chan *PNMessage
	Presence chan *PNPresence
}

func NewListener() *Listener {
	return &Listener{
		Status:   make(chan *PNStatus),
		Message:  make(chan *PNMessage),
		Presence: make(chan *PNPresence),
	}
}

type ListenerManager struct {
	sync.RWMutex
	ctx       Context
	listeners map[*Listener]bool
}

func newListenerManager(ctx Context) *ListenerManager {
	return &ListenerManager{
		listeners: make(map[*Listener]bool, 2),
		ctx:       ctx,
	}
}

func (m *ListenerManager) addListener(listener *Listener) {
	m.Lock()
	defer m.Unlock()

	m.listeners[listener] = true
}

func (m *ListenerManager) removeListener(listener *Listener) {
	m.Lock()
	defer m.Unlock()

	delete(m.listeners, listener)
}

func (m *ListenerManager) removeAllListeners() {
	m.Lock()
	defer m.Unlock()

	for l, _ := range m.listeners {
		delete(m.listeners, l)
	}
}

func (m *ListenerManager) announceStatus(status *PNStatus) {
	m.RLock()
	listeners := m.listeners
	m.RUnlock()

	go func() {
		for l, _ := range listeners {
			select {
			case <-m.ctx.Done():
				return
			case l.Status <- status:
			}
		}
	}()
}

func (m *ListenerManager) announceMessage(message *PNMessage) {
	m.RLock()
	listeners := m.listeners
	m.RUnlock()

	go func() {
		for l, _ := range listeners {
			select {
			case <-m.ctx.Done():
				return
			case l.Message <- message:
			}
		}
	}()
}

func (m *ListenerManager) announcePresence(presence *PNPresence) {
	m.RLock()
	listeners := m.listeners
	m.RUnlock()

	for l, _ := range listeners {
		select {
		case <-m.ctx.Done():
			return
		case l.Presence <- presence:
		}
	}
}

type PNStatus struct {
	Category  StatusCategory
	Operation OperationType

	ErrorData  error
	Error      bool
	TlsEnabled bool
	StatusCode int
	Uuid       string
	AuthKey    string
	Origin     string
	// Should be same for non-google environment
	ClientRequest interface{}

	AffectedChannels      []string
	AffectedChannelGroups []string
}

type PNMessage struct {
	Message      interface{}
	UserMetadata interface{}

	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
	Publisher         string
	Timetoken         int64
}

type PNPresence struct {
	Event             string
	Uuid              string
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string

	Occupancy int
	Timetoken int64
	Timestamp int64

	UserMetadata map[string]interface{}
	State        interface{}

	Join    []string
	Leave   []string
	Timeout []string
}
