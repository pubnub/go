package messaging

import (
	"strings"
	"sync"
)

type ConnectionAction int

type ConnectionEvent struct {
	Channel string
	Group   string
	Action  ConnectionAction
	Type    ResponseType
}

const (
	ConnectionConnected ConnectionAction = 1 << iota
	ConnectionDisconnected
	ConnectionReconnected
)

type SubscriptionItem struct {
	Name           string
	SuccessChannel chan<- SuccessResponse
	ErrorChannel   chan<- ErrorResponse
	EventsChannel  chan<- ConnectionEvent
	Connected      bool
}

func (e *SubscriptionItem) SetConnected() (changed bool) {
	if e.Connected == false {
		e.Connected = true
		return true
	} else {
		return false
	}
}

type SubscriptionEntity struct {
	sync.RWMutex
	items map[string]*SubscriptionItem
}

func NewSubscriptionEntity() *SubscriptionEntity {
	e := new(SubscriptionEntity)

	e.items = make(map[string]*SubscriptionItem)

	return e
}

func (e *SubscriptionEntity) Add(name string,
	successChannel chan<- SuccessResponse, errorChannel chan<- ErrorResponse,
	eventsChannel chan<- ConnectionEvent) {

	e.Lock()
	defer e.Unlock()

	item := &SubscriptionItem{
		Name:           name,
		SuccessChannel: successChannel,
		ErrorChannel:   errorChannel,
		EventsChannel:  eventsChannel,
		Connected:      false,
	}

	e.items[name] = item
}

func (e *SubscriptionEntity) Remove(name string) bool {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.items[name]; ok {
		delete(e.items, name)

		return true
	} else {
		return false
	}
}

func (e *SubscriptionEntity) Length() int {
	return len(e.items)
}

func (e *SubscriptionEntity) Exist(name string) bool {
	if _, ok := e.items[name]; ok {
		return true
	} else {
		return false
	}
}

func (e *SubscriptionEntity) Empty() bool {
	return len(e.items) == 0
}

func (e *SubscriptionEntity) Get(name string) (*SubscriptionItem, bool) {
	if _, ok := e.items[name]; ok {
		return e.items[name], true
	} else {
		return nil, false
	}
}

func (e *SubscriptionEntity) Names() []string {
	e.RLock()
	defer e.RUnlock()

	var names []string

	for k, _ := range e.items {
		names = append(names, k)
	}

	return names
}

func (e *SubscriptionEntity) NamesString() string {
	names := e.Names()

	return strings.Join(names, ",")
}

func (e *SubscriptionEntity) ConnectedNames() []string {
	e.RLock()
	defer e.RUnlock()

	var names []string

	for k, item := range e.items {
		if item.Connected {
			names = append(names, k)
		}
	}

	return names
}

func (e *SubscriptionEntity) ConnectedNamesString() string {
	names := e.ConnectedNames()

	return strings.Join(names, ",")
}

func (e *SubscriptionEntity) Clear() {
	e.Lock()
	defer e.Unlock()

	e.items = make(map[string]*SubscriptionItem)
}

func (e *SubscriptionEntity) ResetConnected() {
	e.Lock()
	defer e.Unlock()

	for _, item := range e.items {
		item.Connected = false
	}
}

func (e *SubscriptionEntity) SetConnected() (changedItemNames []string) {
	e.Lock()
	defer e.Unlock()

	for name, item := range e.items {
		if item.SetConnected() == true {
			changedItemNames = append(changedItemNames, name)
		}
	}

	return changedItemNames
}

func CreateSubscriptionChannels() (chan SuccessResponse, chan ErrorResponse,
	chan ConnectionEvent) {

	successResponse := make(chan SuccessResponse)
	errorResponse := make(chan ErrorResponse)
	connectionEvent := make(chan ConnectionEvent)

	return successResponse, errorResponse, connectionEvent
}

func StringConnectionAction(status ConnectionAction) string {
	switch status {
	case ConnectionConnected:
		return "connect"
	case ConnectionDisconnected:
		return "disconnect"
	case ConnectionReconnected:
		return "reconnect"
	default:
		return ""
	}
}
