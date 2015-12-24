package messaging

import (
	"fmt"
	"strings"
	"sync"
)

const (
	connectionConnected connectionAction = 1 << iota
	connectionUnsubscribed
	connectionReconnected
)

type connectionAction int

type connectionEvent struct {
	Channel string
	Source  string
	Action  connectionAction
	Type    responseType
}

func newConnectionEventForChannel(channel string,
	action connectionAction) *connectionEvent {
	return &connectionEvent{
		Channel: channel,
		Action:  action,
		Type:    channelResponse,
	}
}

func newConnectionEventForChannelGroup(group string,
	action connectionAction) *connectionEvent {
	return &connectionEvent{
		Source: group,
		Action: action,
		Type:   channelGroupResponse,
	}
}

func (e connectionEvent) Bytes() []byte {
	switch e.Type {
	case channelResponse:
		fallthrough
	case wildcardResponse:
		return []byte(fmt.Sprintf(
			"[1, \"%s channel '%s' %sed\", \"%s\"]",
			stringPresenceOrSubscribe(e.Channel),
			removePnpres(e.Channel), stringConnectionAction(e.Action),
			removePnpres(e.Channel)))

	case channelGroupResponse:
		return []byte(fmt.Sprintf(
			"[1, \"%s channel group '%s' %sed\", \"%s\"]",
			stringPresenceOrSubscribe(e.Source),
			removePnpres(e.Source), stringConnectionAction(e.Action),
			strings.Replace(e.Source, presenceSuffix, "", -1)))

	default:
		panic(fmt.Sprintf("Undefined response type: %d", e.Type))
	}
}

type subscriptionItem struct {
	Name           string
	SuccessChannel chan<- []byte
	ErrorChannel   chan<- []byte
	Connected      bool
}

func (e *subscriptionItem) SetConnected() (changed bool) {
	if e.Connected == false {
		e.Connected = true
		return true
	} else {
		return false
	}
}

type subscriptionEntity struct {
	sync.RWMutex
	items         map[string]*subscriptionItem
	abortedMarker bool
}

func newSubscriptionEntity() *subscriptionEntity {
	e := new(subscriptionEntity)

	e.items = make(map[string]*subscriptionItem)

	return e
}

func (e *subscriptionEntity) Add(name string,
	successChannel chan<- []byte, errorChannel chan<- []byte) {
	e.add(name, false, successChannel, errorChannel)
}

func (e *subscriptionEntity) AddConnected(name string,
	successChannel chan<- []byte, errorChannel chan<- []byte) {
	e.add(name, true, successChannel, errorChannel)
}

func (e *subscriptionEntity) add(name string, connected bool,
	successChannel chan<- []byte, errorChannel chan<- []byte) {

	e.Lock()
	defer e.Unlock()

	item := &subscriptionItem{
		Name:           name,
		SuccessChannel: successChannel,
		ErrorChannel:   errorChannel,
		Connected:      connected,
	}

	e.items[name] = item
}

func (e *subscriptionEntity) Remove(name string) bool {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.items[name]; ok {
		delete(e.items, name)

		return true
	} else {
		return false
	}
}

func (e *subscriptionEntity) Length() int {
	e.RLock()
	defer e.RUnlock()

	return len(e.items)
}

func (e *subscriptionEntity) Exist(name string) bool {
	e.RLock()
	defer e.RUnlock()

	if _, ok := e.items[name]; ok {
		return true
	} else {
		return false
	}
}

func (e *subscriptionEntity) Empty() bool {
	e.RLock()
	defer e.RUnlock()

	return len(e.items) == 0
}

func (e *subscriptionEntity) Get(name string) (*subscriptionItem, bool) {
	e.RLock()
	defer e.RUnlock()

	if _, ok := e.items[name]; ok {
		return e.items[name], true
	} else {
		return nil, false
	}
}

func (e *subscriptionEntity) Names() []string {
	e.RLock()
	defer e.RUnlock()

	var names []string

	for k, _ := range e.items {
		names = append(names, k)
	}

	return names
}

func (e *subscriptionEntity) NamesString() string {
	names := e.Names()

	return strings.Join(names, ",")
}

func (e *subscriptionEntity) HasConnected() bool {
	e.RLock()
	defer e.RUnlock()

	for _, item := range e.items {
		if item.Connected {
			return true
		}
	}

	return false
}

func (e *subscriptionEntity) ConnectedNames() []string {
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

func (e *subscriptionEntity) ConnectedNamesString() string {
	names := e.ConnectedNames()

	return strings.Join(names, ",")
}

func (e *subscriptionEntity) Clear() {
	e.Lock()
	defer e.Unlock()

	e.items = make(map[string]*subscriptionItem)
}

func (e *subscriptionEntity) Abort() {
	e.Lock()
	defer e.Unlock()

	e.abortedMarker = true
}

func (e *subscriptionEntity) ApplyAbort() {
	e.Lock()
	defer e.Unlock()

	if e.abortedMarker == true {
		e.Clear()
	}
}

func (e *subscriptionEntity) ResetConnected() {
	e.Lock()
	defer e.Unlock()

	for _, item := range e.items {
		item.Connected = false
	}
}

func (e *subscriptionEntity) SetConnected() (changedItemNames []string) {
	e.Lock()
	defer e.Unlock()

	for name, item := range e.items {
		if item.SetConnected() == true {
			changedItemNames = append(changedItemNames, name)
		}
	}

	return changedItemNames
}

func createSubscriptionChannels() (chan []byte, chan []byte) {

	successResponse := make(chan []byte)
	errorResponse := make(chan []byte)

	return successResponse, errorResponse
}
