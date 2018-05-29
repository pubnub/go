package pubnub

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	reconnectionInterval              = 10
	reconnectionMinExponentialBackoff = 1
	reconnectionMaxExponentialBackoff = 32
)

type ReconnectionManager struct {
	sync.RWMutex

	timerMutex sync.RWMutex

	ExponentialMultiplier       int
	FailedCalls                 int
	Milliseconds                int
	OnReconnection              func()
	OnMaxReconnectionExhaustion func()
	DoneTimer                   chan bool
	//Timer                       *time.Ticker
	hbRunning               bool
	pubnub                  *PubNub
	exitReconnectionManager chan bool
}

func newReconnectionManager(pubnub *PubNub) *ReconnectionManager {
	manager := &ReconnectionManager{}

	manager.pubnub = pubnub

	manager.ExponentialMultiplier = 1
	manager.FailedCalls = 0
	manager.Milliseconds = 1000
	manager.exitReconnectionManager = make(chan bool)
	manager.hbRunning = false

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
		m.pubnub.Config.Log.Println("Reconnection policy is disabled, please handle reconnection manually.")
		return
	}

	m.Lock()
	m.ExponentialMultiplier = 1
	m.FailedCalls = 0
	hbRunning := m.hbRunning
	m.Unlock()

	if !hbRunning {
		m.pubnub.Config.Log.Println(fmt.Sprintf("Reconnection policy: %d, retries: %d", m.pubnub.Config.PNReconnectionPolicy, m.pubnub.Config.MaximumReconnectionRetries))
		m.startHeartbeatTimer()
	}

}

func (m *ReconnectionManager) startHeartbeatTimer() {
	//m.stopHeartbeatTimer()

	/*if m.pubnub.Config.PNReconnectionPolicy == PNNonePolicy {
		return
	}*/

	timerInterval := reconnectionInterval

	//m.Lock()
	//m.Timer = time.NewTicker(time.Duration(timerInterval) * time.Second)
	//timeout := time.After(time.Duration(timeoutVal) * time.Second)

	//m.Unlock()

	for {

		m.Lock()
		m.hbRunning = true
		m.Unlock()
		_, status, err := m.pubnub.Time().Execute()
		if status.Error == nil {
			m.RLock()
			failedCalls := m.FailedCalls
			m.RUnlock()
			if failedCalls > 0 {
				timerInterval = reconnectionInterval
				m.Lock()
				m.FailedCalls = 0
				m.Unlock()
				m.pubnub.Config.Log.Println(fmt.Sprintf("Network reconnected"))
				m.OnReconnection()
			}
			//break
		} else {
			if m.pubnub.Config.PNReconnectionPolicy == PNExponentialPolicy {
				timerInterval = m.GetExponentialInterval()
			}
			m.Lock()
			m.FailedCalls++
			m.pubnub.Config.Log.Println(fmt.Sprintf("Network disconnected, reconnection try %d of %d\n %s %v", m.FailedCalls, m.pubnub.Config.MaximumReconnectionRetries, status, err))
			m.ExponentialMultiplier++

			failedCalls := m.FailedCalls
			retries := m.pubnub.Config.MaximumReconnectionRetries
			m.Unlock()
			if retries != -1 && failedCalls >= retries {
				m.pubnub.Config.Log.Printf(fmt.Sprintf("Network connection retry limit (%d) exceeded", retries))
				go m.OnMaxReconnectionExhaustion()
				m.Lock()
				m.hbRunning = false
				m.Unlock()
				return
			}
		}

		select {
		case <-time.After(time.Duration(timerInterval) * time.Second):
		case <-m.pubnub.ctx.Done():
			m.pubnub.Config.Log.Printf(fmt.Sprintf("==========> pubnub.ctx.Done\n"))
			m.Lock()
			m.hbRunning = false
			m.Unlock()
			return
		case <-m.exitReconnectionManager:
			m.pubnub.Config.Log.Printf(fmt.Sprintf("==========> exitReconnectionManager\n"))
			m.Lock()
			m.hbRunning = false
			m.Unlock()
			return
		}
		//m.registerHeartbeatTimer()
		//}
	}

	/*go func() {
		if m.Timer == nil {
			return
		}

		m.Lock()
		timer := m.Timer.C
		doneT := m.DoneTimer
		m.Unlock()

		for {
			select {
			case <-m.pubnub.ctx.Done():
				return
			case <-timer:
				go m.callTime()
			case <-doneT:
				return
			case <-m.exitReconnectionManager:
				return
			}
		}
	}()*/
}

func (m *ReconnectionManager) GetExponentialInterval() int {
	timerInterval := int(math.Pow(2, float64(m.ExponentialMultiplier)) - 1)
	if timerInterval > reconnectionMaxExponentialBackoff {
		timerInterval = reconnectionMinExponentialBackoff

		m.Lock()
		m.ExponentialMultiplier = 1
		m.pubnub.Config.Log.Printf(fmt.Sprintf("==========> timerInterval > MaxExponentialBackoff at: %d\n", m.ExponentialMultiplier))
		m.Unlock()

	} else if timerInterval < 1 {
		timerInterval = reconnectionMinExponentialBackoff
		m.Lock()
		m.ExponentialMultiplier = 1
		m.pubnub.Config.Log.Printf(fmt.Sprintf("==========> timerInterval < 1 at: %d\n", m.ExponentialMultiplier))
		m.Unlock()
	}
	return timerInterval
}

/*func (m *ReconnectionManager) stopHeartbeatTimer() {
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
	//_, status, err := m.pubnub.Time().Execute()

	//m.stopHeartbeatTimer()

}*/
