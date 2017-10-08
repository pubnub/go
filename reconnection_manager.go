package pubnub

import (
	"log"
	"math"
	"sync"
	"time"
)

const (
	RECONNECTION_INTERVAL              = 3
	RECONNECTION_MINEXPONENTIALBACKOFF = 1
	RECONNECTION_MAXEXPONENTIALBACKOFF = 32
)

type ReconnectionManager struct {
	sync.RWMutex

	timerMutex sync.RWMutex

	ExponentialMultiplier int
	FailedCalls           int
	Milliseconds          int

	OnReconnection              func()
	OnMaxReconnectionExhaustion func()

	DoneTimer chan bool

	Timer *time.Ticker

	pubnub *PubNub
}

func newReconnectionManager(pubnub *PubNub) *ReconnectionManager {
	manager := &ReconnectionManager{}

	manager.pubnub = pubnub

	manager.ExponentialMultiplier = 1
	manager.FailedCalls = 0
	manager.Milliseconds = 1000

	return manager
}

func (m *ReconnectionManager) HandleReconnection(handler func()) {
	m.Lock()
	m.OnReconnection = handler
	m.Unlock()
}

func (m *ReconnectionManager) HandleOnMaxReconnectionExhaustion(handler func()) {
	m.Lock()
	m.OnMaxReconnectionExhaustion = handler
	m.Unlock()
}

func (m *ReconnectionManager) startPolling() {
	if m.pubnub.Config.PNReconnectionPolicy == PNNonePolicy {
		log.Println("reconnection policy is disabled, please handle reconnection manually")
		return
	}

	m.Lock()
	m.ExponentialMultiplier = 1
	m.FailedCalls = 0
	m.Unlock()

	m.registerHeartbeatTimer()
}

func (m *ReconnectionManager) registerHeartbeatTimer() {
	m.stopHeartbeatTimer()

	if m.pubnub.Config.PNReconnectionPolicy == PNNonePolicy {
		log.Println("Reconnection policy is disabled, please handle reconnection manually.")
		return
	}

	maxRetries := m.pubnub.Config.MaximumReconnectionRetries

	m.RLock()
	failedCalls := m.FailedCalls
	m.RUnlock()

	if maxRetries != -1 && failedCalls >= maxRetries {
		go m.OnMaxReconnectionExhaustion()
		return
	}

	timerInterval := RECONNECTION_INTERVAL

	if m.pubnub.Config.PNReconnectionPolicy == PNExponentialPolicy {
		timerInterval = int(math.Pow(2, float64(m.ExponentialMultiplier)) - 1)
		if timerInterval > RECONNECTION_MAXEXPONENTIALBACKOFF {
			timerInterval = RECONNECTION_MINEXPONENTIALBACKOFF

			m.Lock()
			m.ExponentialMultiplier = 1
			m.Unlock()

			// TODO: add timestamp
			log.Printf("timerInterval > MAXEXPONENTIALBACKOFF at: \n")
		} else {
			timerInterval = RECONNECTION_MINEXPONENTIALBACKOFF
		}
		// log.debug("timerInterval = " + String.valueOf(timerInterval) + " at: " +
		// Calendar.getInstance().getTime().toString());
	}

	if m.pubnub.Config.PNReconnectionPolicy == PNLinearPolicy {
		timerInterval = RECONNECTION_INTERVAL
	}

	m.Lock()
	m.Timer = time.NewTicker(time.Duration(timerInterval) * time.Second)
	m.Unlock()

	go func() {
		// Lock??
		if m.Timer == nil {
			return
		}

		m.Lock()
		timer := m.Timer.C
		doneT := m.DoneTimer
		m.Unlock()

		for {
			select {
			case <-timer:
				go m.callTime()
			case <-doneT:
				return
			}
		}
	}()
}

func (m *ReconnectionManager) stopHeartbeatTimer() {
	m.timerMutex.Lock()

	if m.Timer != nil {
		m.Timer.Stop()
	}

	if m.DoneTimer != nil {
		m.DoneTimer <- true
	}

	m.timerMutex.Unlock()
}

func (m *ReconnectionManager) callTime() {
	_, status, err := m.pubnub.Time().Execute()
	if err != nil {
		return
	}

	if status.Error == nil {
		m.stopHeartbeatTimer()
		m.OnReconnection()
		return
	}

	m.Lock()
	m.ExponentialMultiplier++
	m.FailedCalls++
	m.Unlock()

	m.registerHeartbeatTimer()
}
