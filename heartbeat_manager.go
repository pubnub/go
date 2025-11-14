package pubnub

import (
	"fmt"
	"sync"
	"time"
)

// HeartbeatManager is a struct that assists in running of the heartbeat.
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
	queryParam                map[string]string
	state                     map[string]interface{}
}

func newHeartbeatManager(pn *PubNub, context Context) *HeartbeatManager {
	return &HeartbeatManager{
		heartbeatChannels: make(map[string]*SubscriptionItem),
		heartbeatGroups:   make(map[string]*SubscriptionItem),
		ctx:               context,
		pubnub:            pn,
	}
}

// Destroy stops the running heartbeat.
func (m *HeartbeatManager) Destroy() {
	m.stopHeartbeat(true, true)
}

func (m *HeartbeatManager) nonIndependentHeartbeatLoop() {
	timeNow := time.Now().Unix()

	m.pubnub.subscriptionManager.hbDataMutex.RLock()
	reqSentAt := m.pubnub.subscriptionManager.requestSentAt
	m.pubnub.subscriptionManager.hbDataMutex.RUnlock()

	if reqSentAt > 0 {
		timediff := int64(m.pubnub.Config.HeartbeatInterval) - (timeNow - reqSentAt)
		m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, fmt.Sprintf("Heartbeat timediff: %d seconds", timediff), false)
		m.pubnub.subscriptionManager.hbDataMutex.Lock()
		m.pubnub.subscriptionManager.requestSentAt = 0
		m.pubnub.subscriptionManager.hbDataMutex.Unlock()
		if timediff > 10 {
			m.Lock()
			m.hbTimer.Stop()
			m.Unlock()

			m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, fmt.Sprintf("Heartbeat sleeping for %d seconds", timediff), false)
			waitTimer := time.NewTicker(time.Duration(timediff) * time.Second)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				waitTimerCh := waitTimer.C
				for {
					select {
					case <-m.hbDone:
						wg.Done()
						m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Heartbeat loop interrupted", false)
						return
					case <-waitTimerCh:
						m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Heartbeat sleep completed", false)
						wg.Done()
						return
					}
				}
			}()
			wg.Wait()

			m.Lock()
			m.hbTimer = time.NewTicker(time.Duration(m.pubnub.Config.HeartbeatInterval) * time.Second)
			m.Unlock()
		}
	}
	m.performHeartbeatLoop()
}

func (m *HeartbeatManager) readHeartBeatTimer(runIndependentOfSubscribe bool) {
	go func() {

		defer m.hbLoopMutex.Unlock()
		defer func() {
			m.Lock()
			m.hbDone = nil
			m.Unlock()
		}()
	HeartbeatLabel:
		for {
			m.RLock()
			timerCh := m.hbTimer.C
			m.RUnlock()

			select {
			case <-timerCh:
				if runIndependentOfSubscribe {
					m.performHeartbeatLoop()
				} else {
					m.nonIndependentHeartbeatLoop()
				}
			case <-m.hbDone:
				m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Heartbeat timer stopped", false)
				break HeartbeatLabel
			}
		}
	}()
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

	m.pubnub.loggerManager.LogSimple(PNLogLevelInfo, fmt.Sprintf("Starting heartbeat timer: interval=%d seconds", m.pubnub.Config.HeartbeatInterval), false)
	m.pubnub.Config.Lock()
	presenceTimeout := m.pubnub.Config.PresenceTimeout
	heartbeatInterval := m.pubnub.Config.HeartbeatInterval
	m.pubnub.Config.Unlock()
	if presenceTimeout <= 0 && heartbeatInterval <= 0 {
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

	m.readHeartBeatTimer(runIndependentOfSubscribe)

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
	m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Stopping heartbeat", false)

	m.Lock()
	if m.hbTimer != nil {
		m.hbTimer.Stop()
		m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Heartbeat timer stopped", false)
	}

	if m.hbDone != nil {
		m.hbDone <- true
		m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Heartbeat done channel stopped", false)
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
	var stateStorage map[string]interface{}

	m.RLock()
	presenceChannels := m.prepareList(m.heartbeatChannels)
	presenceGroups := m.prepareList(m.heartbeatGroups)
	stateStorage = m.state
	queryParam := m.queryParam
	m.RUnlock()

	if (len(presenceChannels) == 0) && (len(presenceGroups) == 0) {
		presenceChannels = m.pubnub.subscriptionManager.stateManager.prepareChannelList(false)
		presenceGroups = m.pubnub.subscriptionManager.stateManager.prepareGroupList(false)
		stateStorage = m.pubnub.subscriptionManager.stateManager.createStatePayload()
		queryParam = nil
	}

	if len(presenceChannels) <= 0 && len(presenceGroups) <= 0 {
		m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Heartbeat: no channels left", false)
		go m.stopHeartbeat(true, true)
		return nil
	}

	_, status, err := newHeartbeatBuilder(m.pubnub).
		Channels(presenceChannels).
		ChannelGroups(presenceGroups).
		State(stateStorage).
		QueryParam(queryParam).
		Execute()

	if err != nil {

		pnStatus := &PNStatus{
			Operation: PNHeartBeatOperation,
			Category:  PNBadRequestCategory,
			Error:     true,
			ErrorData: err,
		}
		m.pubnub.loggerManager.LogError(err, "HeartbeatFailed", PNHeartBeatOperation, true)

		m.pubnub.subscriptionManager.listenerManager.announceStatus(pnStatus)

		return err
	}

	pnStatus := &PNStatus{
		Category:   PNUnknownCategory,
		Error:      false,
		Operation:  PNHeartBeatOperation,
		StatusCode: status.StatusCode,
	}
	m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, fmt.Sprintf("Heartbeat sent successfully: channels=%d, groups=%d", len(presenceChannels), len(presenceGroups)), false)

	m.pubnub.subscriptionManager.listenerManager.announceStatus(pnStatus)

	return nil
}
