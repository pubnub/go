package contract

import (
	"sync"
)

type subscribeStateKey struct{}

type subscribeState struct {
	sync.RWMutex
	allSubscribeMessages chan interface{}
	all                  []interface{}
}

func (s *subscribeState) addSubscribeMessage(msg interface{}) {
	s.Lock()
	defer s.Unlock()
	s.all = append(s.all, msg)
}

func (s *subscribeState) readAllSubscribeMessages() []interface{} {
	s.RLock()
	defer s.RUnlock()
	r := make([]interface{}, len(s.all))
	copy(r, s.all)
	return r
}

func newSubscribeState() *subscribeState {
	return &subscribeState{
		allSubscribeMessages: make(chan interface{}),
	}
}
