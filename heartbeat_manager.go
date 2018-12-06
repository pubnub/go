package pubnub

import (
	"fmt"
	"sync"
	"time"
)

type HeartbeatManager struct {
	sync.RWMutex

	heartbeatChannels map[string]*SubscriptionItem
	heartbeatGroups   map[string]*SubscriptionItem
	pubnub            *PubNub

	hbLoopMutex               sync.RWMutex
	hbTimer                   *time.Ticker
	hbDone                    chan bool
	ctx                       Context
	runIndependentOfSubscribe bool
	hbRunning                 bool
}

func newHeartbeatManager(pn *PubNub, context Context) *HeartbeatManager {
	return &HeartbeatManager{
		heartbeatChannels: make(map[string]*SubscriptionItem),
		heartbeatGroups:   make(map[string]*SubscriptionItem),
		ctx:               context,
		pubnub:            pn,
	}
}

func (m *HeartbeatManager) Destroy() {
	m.stopHeartbeat(true, true)
}

func (m *HeartbeatManager) startHeartbeatTimer(runIndependentOfSubscribe bool) {
	m.RLock()
	hbRunning := m.hbRunning
	m.RUnlock()
	if hbRunning && !runIndependentOfSubscribe {
		return
	}
	m.stopHeartbeat(runIndependentOfSubscribe, true)

	m.Lock()
	m.hbRunning = true
	m.Unlock()

	m.pubnub.Config.Log.Println("heartbeat: new timer", m.pubnub.Config.HeartbeatInterval)
	if m.pubnub.Config.PresenceTimeout <= 0 && m.pubnub.Config.HeartbeatInterval <= 0 {
		return
	}

	m.hbLoopMutex.Lock()
	m.Lock()
	m.hbDone = make(chan bool)
	m.hbTimer = time.NewTicker(time.Duration(m.pubnub.Config.HeartbeatInterval) * time.Second)
	m.Unlock()

	if runIndependentOfSubscribe {
		m.performHeartbeatLoop()
	}

	go func() {
		defer m.hbLoopMutex.Unlock()
		defer func() {
			m.Lock()
			m.hbDone = nil
			m.Unlock()
		}()

		for {
			m.RLock()
			timerCh := m.hbTimer.C
			doneCh := m.hbDone
			m.RUnlock()

			select {
			case <-timerCh:
				timeNow := time.Now().Unix()
				if runIndependentOfSubscribe {
					m.performHeartbeatLoop()
				} else {
					m.pubnub.subscriptionManager.hbDataMutex.RLock()
					reqSentAt := m.pubnub.subscriptionManager.requestSentAt
					m.pubnub.subscriptionManager.hbDataMutex.RUnlock()

					if reqSentAt > 0 {
						timediff := int64(m.pubnub.Config.HeartbeatInterval) - (timeNow - reqSentAt)
						m.pubnub.Config.Log.Println(fmt.Sprintf("heartbeat timediff: %d", timediff))
						m.pubnub.subscriptionManager.hbDataMutex.Lock()
						m.pubnub.subscriptionManager.requestSentAt = 0
						m.pubnub.subscriptionManager.hbDataMutex.Unlock()
						if timediff > 10 {
							m.Lock()
							m.hbTimer.Stop()
							m.Unlock()

							m.pubnub.Config.Log.Println(fmt.Sprintf("heartbeat sleeping timediff: %d", timediff))
							time.Sleep(time.Duration(timediff) * time.Second)
							m.pubnub.Config.Log.Println("heartbeat sleep end")
							m.Lock()
							m.hbTimer = time.NewTicker(time.Duration(m.pubnub.Config.HeartbeatInterval) * time.Second)
							m.Unlock()
						}
					}
					m.performHeartbeatLoop()
				}
			case <-doneCh:
				m.pubnub.Config.Log.Println("heartbeat: loop after stop")
				return
			}
		}
	}()
}

func (m *HeartbeatManager) stopHeartbeat(runIndependentOfSubscribe bool, skipRuncheck bool) {
	if !skipRuncheck {
		m.RLock()
		hbRunning := m.hbRunning
		m.RUnlock()

		if hbRunning && !runIndependentOfSubscribe {
			return
		}
	}
	m.pubnub.Config.Log.Println("heartbeat: loop: stopping...")

	m.Lock()
	if m.hbTimer != nil {
		m.hbTimer.Stop()
		m.pubnub.Config.Log.Println("heartbeat: loop: timer stopped")
	}

	if m.hbDone != nil {
		m.hbDone <- true
		m.pubnub.Config.Log.Println("heartbeat: loop: done channel stopped")
	}
	m.hbRunning = false
	m.Unlock()
	m.pubnub.subscriptionManager.hbDataMutex.Lock()
	m.pubnub.subscriptionManager.requestSentAt = 0
	m.pubnub.subscriptionManager.hbDataMutex.Unlock()
}

func (m *HeartbeatManager) prepareList(subItem map[string]*SubscriptionItem) []string {
	response := []string{}

	for _, v := range subItem {
		response = append(response, v.name)
	}
	return response
}

func (m *HeartbeatManager) performHeartbeatLoop() error {
	m.RLock()
	presenceChannels := m.prepareList(m.heartbeatChannels)
	presenceGroups := m.prepareList(m.heartbeatGroups)
	m.pubnub.Config.Log.Println("performHeartbeatLoop: count presenceChannels, presenceGroups", len(presenceChannels), len(presenceGroups))
	m.RUnlock()
	var stateStorage map[string]interface{}
	if (len(presenceChannels) == 0) && (len(presenceGroups) == 0) {
		m.pubnub.Config.Log.Println("performHeartbeatLoop: count presenceChannels, presenceGroups nil")
		presenceChannels = m.pubnub.subscriptionManager.stateManager.prepareChannelList(false)
		presenceGroups = m.pubnub.subscriptionManager.stateManager.prepareGroupList(false)
		stateStorage = m.pubnub.subscriptionManager.stateManager.createStatePayload()
		m.pubnub.Config.Log.Println("performHeartbeatLoop: count sub presenceChannels, presenceGroups", len(presenceChannels), len(presenceGroups))
	}

	if len(presenceChannels) <= 0 && len(presenceGroups) <= 0 {
		m.pubnub.Config.Log.Println("heartbeat: no channels left")
		go m.stopHeartbeat(true, true)
		return nil
	}

	_, status, err := newHeartbeatBuilder(m.pubnub).
		Channels(presenceChannels).
		ChannelGroups(presenceGroups).
		State(stateStorage).
		Execute()

	if err != nil {

		pnStatus := &PNStatus{
			Operation: PNHeartBeatOperation,
			Category:  PNBadRequestCategory,
			Error:     true,
			ErrorData: err,
		}
		m.pubnub.Config.Log.Println("performHeartbeatLoop: err", err, pnStatus)

		m.pubnub.subscriptionManager.listenerManager.announceStatus(pnStatus)

		return err
	}

	pnStatus := &PNStatus{
		Category:   PNUnknownCategory,
		Error:      false,
		Operation:  PNHeartBeatOperation,
		StatusCode: status.StatusCode,
	}
	m.pubnub.Config.Log.Println("performHeartbeatLoop: err", err, pnStatus)

	m.pubnub.subscriptionManager.listenerManager.announceStatus(pnStatus)

	return nil
}
