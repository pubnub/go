package pubnub

import "sync"

type Listener struct {
	Status   chan PNStatus
	Message  chan PNMessage
	Presence chan PNPresence
}

func NewListener() *Listener {
	return &Listener{
		Status:   make(chan PNStatus),
		Message:  make(chan PNMessage),
		Presence: make(chan PNPresence),
	}
}

type ListenerManager struct {
	sync.RWMutex
	listeners map[*Listener]bool
}

func newListenerManager() *ListenerManager {
	return &ListenerManager{
		listeners: make(map[*Listener]bool, 2),
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

func (m *ListenerManager) announceStatus(status *PNStatus) {
	m.RLock()
	defer m.RUnlock()

	for l, _ := range m.listeners {
		l.Status <- *status
	}
}

func (m *ListenerManager) announceMessage(message *PNMessage) {
	m.RLock()
	defer m.RUnlock()

	for l, _ := range m.listeners {
		l.Message <- *message
	}
}

func (m *ListenerManager) announcePresence(presence *PNPresence) {
	m.RLock()
	defer m.RUnlock()

	for l, _ := range m.listeners {
		l.Presence <- *presence
	}
}

type PNStatus int
type PNMessage int
type PNPresence int
