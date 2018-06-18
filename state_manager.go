package pubnub

import (
	"fmt"
	"strings"
	"sync"
)

// StateManager is used to store the subscriptions types
type StateManager struct {
	sync.RWMutex

	channels         map[string]*SubscriptionItem
	groups           map[string]*SubscriptionItem
	presenceChannels map[string]*SubscriptionItem
	presenceGroups   map[string]*SubscriptionItem
}

// SubscriptionItem is used to store the subscription item's properties.
type SubscriptionItem struct {
	name  string
	state map[string]interface{}
}

func newStateManager() *StateManager {
	return &StateManager{
		channels:         make(map[string]*SubscriptionItem),
		presenceChannels: make(map[string]*SubscriptionItem),
		groups:           make(map[string]*SubscriptionItem),
		presenceGroups:   make(map[string]*SubscriptionItem),
	}
}

func newSubscriptionItem(name string) *SubscriptionItem {
	return &SubscriptionItem{
		name:  name,
		state: make(map[string]interface{}),
	}
}

func newSubscriptionItemWithState(name string, state map[string]interface{}) *SubscriptionItem {
	return &SubscriptionItem{
		name:  name,
		state: state,
	}
}

func (m *StateManager) prepareChannelList(includePresence bool) []string {
	m.RLock()
	channels := m.channels
	presenceChannels := m.presenceChannels
	m.RUnlock()

	return m.prepareMembershipList(channels, presenceChannels, includePresence)
}

func (m *StateManager) prepareGroupList(includePresence bool) []string {
	m.RLock()
	groups := m.groups
	presenceGroups := m.presenceGroups
	m.RUnlock()

	return m.prepareMembershipList(groups, presenceGroups, includePresence)
}

func (m *StateManager) adaptSubscribeOperation(
	subscribeOperation *SubscribeOperation) {
	m.Lock()

	for _, ch := range subscribeOperation.Channels {
		if strings.Contains(ch, "-pnpres") {
			if len(subscribeOperation.State) > 0 {
				m.presenceChannels[ch] = newSubscriptionItemWithState(ch, subscribeOperation.State)
			} else {
				m.presenceChannels[ch] = newSubscriptionItem(ch)
			}
		} else {
			if len(subscribeOperation.State) > 0 {
				m.channels[ch] = newSubscriptionItemWithState(ch, subscribeOperation.State)
			} else {
				m.channels[ch] = newSubscriptionItem(ch)
			}

			if subscribeOperation.PresenceEnabled {
				m.presenceChannels[ch] = newSubscriptionItem(ch)
			}
		}

	}

	for _, cg := range subscribeOperation.ChannelGroups {
		if strings.Contains(cg, "-pnpres") {
			if len(subscribeOperation.State) > 0 {
				m.presenceGroups[cg] = newSubscriptionItemWithState(cg, subscribeOperation.State)
			} else {
				m.presenceGroups[cg] = newSubscriptionItem(cg)
			}
		} else {
			if len(subscribeOperation.State) > 0 {
				m.groups[cg] = newSubscriptionItemWithState(cg, subscribeOperation.State)
			} else {
				m.groups[cg] = newSubscriptionItem(cg)
			}

			if subscribeOperation.PresenceEnabled {
				m.presenceGroups[cg] = newSubscriptionItem(cg)
			}
		}
	}
	m.Unlock()
}

func (m *StateManager) adaptStateOperation(stateOperation StateOperation) {
	m.Lock()

	for _, ch := range stateOperation.channels {
		if _, ok := m.channels[ch]; ok {
			subscribedChannel := m.channels[ch]

			if subscribedChannel.name != "" {
				subscribedChannel.state = stateOperation.state
			}
		}
	}

	for _, cg := range stateOperation.channelGroups {
		if _, ok := m.groups[cg]; ok {
			subscribedChannelGroup := m.groups[cg]

			if subscribedChannelGroup.name != "" {
				subscribedChannelGroup.state = stateOperation.state
			}
		}
	}
	m.Unlock()
}

func (m *StateManager) adaptUnsubscribeOperation(unsubscribeOperation *UnsubscribeOperation) {
	m.Lock()

	for _, ch := range unsubscribeOperation.Channels {
		if strings.Contains(ch, "-pnpres") {
			delete(m.presenceChannels, strings.Replace(ch, "-pnpres", "", -1))
		} else {
			delete(m.channels, ch)
		}
	}

	for _, cg := range unsubscribeOperation.ChannelGroups {
		if strings.Contains(cg, "-pnpres") {
			delete(m.presenceGroups, strings.Replace(cg, "-pnpres", "", -1))
		} else {
			delete(m.groups, cg)
		}

	}
	m.Unlock()
}

func (m *StateManager) createStatePayload() map[string]interface{} {
	m.RLock()
	defer m.RUnlock()

	stateResponse := make(map[string]interface{})

	for _, ch := range m.channels {
		if len(ch.state) != 0 {
			stateResponse[ch.name] = ch.state
		}
	}

	for _, gr := range m.groups {
		if len(gr.state) != 0 {
			stateResponse[gr.name] = gr.state
		}
	}

	return stateResponse
}

func (m *StateManager) isEmpty() bool {
	m.RLock()
	defer m.RUnlock()

	return len(m.channels) != 0 && len(m.presenceChannels) != 0 &&
		len(m.groups) != 0 && len(m.presenceGroups) != 0
}

func (m *StateManager) hasNonPresenceChannels() bool {
	m.RLock()
	defer m.RUnlock()

	return len(m.channels) <= 0 && len(m.groups) <= 0
}

func (m *StateManager) prepareMembershipList(dataStorage map[string]*SubscriptionItem,
	presenceStorage map[string]*SubscriptionItem, includePresence bool) []string {

	response := []string{}

	m.Lock()
	for _, v := range dataStorage {
		response = append(response, v.name)
	}

	if includePresence {
		for _, v := range presenceStorage {
			if !strings.Contains(v.name, "-pnpres") {
				response = append(response, fmt.Sprintf("%s-pnpres", v.name))
			} else {
				response = append(response, v.name)
			}
		}
	}
	m.Unlock()

	return response
}
