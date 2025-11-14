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

// ReconnectionManager is used to store the properties required in running the Reconnection Manager.
type ReconnectionManager struct {
	sync.RWMutex

	timerMutex sync.RWMutex

	ExponentialMultiplier       int
	FailedCalls                 int
	Milliseconds                int
	OnReconnection              func()
	OnMaxReconnectionExhaustion func()
	DoneTimer                   chan bool
	hbRunning                   bool
	pubnub                      *PubNub
	exitReconnectionManager     chan bool
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

// HandleReconnection sets the handler that will be called when the network reconnects after a disconnect.
func (m *ReconnectionManager) HandleReconnection(handler func()) {
	m.Lock()
	m.OnReconnection = handler
	m.Unlock()
}

// HandleOnMaxReconnectionExhaustion sets the handler that will be called when the max reconnection attempts are exhausted.
func (m *ReconnectionManager) HandleOnMaxReconnectionExhaustion(handler func()) {
	m.Lock()
	m.OnMaxReconnectionExhaustion = handler
	m.Unlock()
}

func (m *ReconnectionManager) startPolling() {

	if m.pubnub.Config.PNReconnectionPolicy == PNNonePolicy {
		m.pubnub.loggerManager.LogSimple(PNLogLevelInfo, "Reconnection policy is disabled", false)
		return
	}

	m.Lock()
	m.ExponentialMultiplier = 1
	m.FailedCalls = 0
	hbRunning := m.hbRunning
	m.Unlock()

	if !hbRunning {
		m.pubnub.loggerManager.LogSimple(PNLogLevelInfo, fmt.Sprintf("Starting reconnection manager: policy=%d, maxRetries=%d", m.pubnub.Config.PNReconnectionPolicy, m.pubnub.Config.MaximumReconnectionRetries), false)

		m.startHeartbeatTimer()
	} else {
		m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Reconnection manager already running", false)
	}

}

func (m *ReconnectionManager) startHeartbeatTimer() {

	timerInterval := reconnectionInterval

	for {

		m.Lock()
		m.hbRunning = true
		failedCalls := m.FailedCalls
		m.Unlock()
		_, status, _ := m.pubnub.Time().Execute()
		if status.Error == nil {
			if failedCalls > 0 {
				timerInterval = reconnectionInterval
				m.Lock()
				m.FailedCalls = 0
				m.Unlock()
				m.pubnub.loggerManager.LogSimple(PNLogLevelInfo, "Network reconnected", false)
				m.OnReconnection()
			}
		} else {
			if m.pubnub.Config.PNReconnectionPolicy == PNExponentialPolicy {
				timerInterval = m.getExponentialInterval()
			}
			m.Lock()
			m.FailedCalls++
			m.pubnub.loggerManager.LogSimple(PNLogLevelWarn, fmt.Sprintf("Network disconnected, reconnection attempt %d of %d", m.FailedCalls, m.pubnub.Config.MaximumReconnectionRetries), false)
			m.ExponentialMultiplier++

			failedCalls := m.FailedCalls
			retries := m.pubnub.Config.MaximumReconnectionRetries
			m.Unlock()
			if retries != -1 && failedCalls >= retries {
				m.pubnub.loggerManager.LogSimple(PNLogLevelError, fmt.Sprintf("Network connection retry limit (%d) exceeded", retries), false)
				m.Lock()
				m.hbRunning = false
				m.Unlock()
				m.OnMaxReconnectionExhaustion()
				return
			}
		}

		select {
		case <-time.After(time.Duration(timerInterval) * time.Second):
		case <-m.pubnub.ctx.Done():
			m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Reconnection manager stopping: context done", false)
			m.Lock()
			m.hbRunning = false
			m.Unlock()
			return
		case <-m.exitReconnectionManager:
			m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Reconnection manager stopping: exit signal received", false)
			return
		}
	}
}

func (m *ReconnectionManager) getExponentialInterval() int {
	timerInterval := int(math.Pow(2, float64(m.ExponentialMultiplier)) - 1)
	if timerInterval > reconnectionMaxExponentialBackoff {
		timerInterval = reconnectionMinExponentialBackoff

		m.Lock()
		m.ExponentialMultiplier = 1
		m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, fmt.Sprintf("Reconnection backoff exceeded maximum, resetting to minimum: multiplier=%d", m.ExponentialMultiplier), false)
		m.Unlock()

	} else if timerInterval < 1 {
		timerInterval = reconnectionMinExponentialBackoff
		m.Lock()
		m.ExponentialMultiplier = 1
		m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, fmt.Sprintf("Reconnection interval too small, resetting: multiplier=%d", m.ExponentialMultiplier), false)
		m.Unlock()
	}
	return timerInterval
}

func (m *ReconnectionManager) stopHeartbeatTimer() {
	m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Stopping reconnection heartbeat timer", false)
	m.Lock()
	if m.hbRunning {
		m.hbRunning = false
		// Use non-blocking send to prevent deadlock when the receiver
		// is blocked on network I/O (like Time().Execute())
		select {
		case m.exitReconnectionManager <- true:
			// Successfully sent exit signal - immediate shutdown
			m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Reconnection exit signal sent successfully", false)
		default:
			// Channel is full or no receiver ready - this is OK since we set hbRunning = false
			// The heartbeat timer will eventually check hbRunning and exit gracefully
			m.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Reconnection exit signal not sent, relying on hbRunning flag", false)
		}
	}
	m.Unlock()
}
