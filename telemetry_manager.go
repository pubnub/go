package pubnub

import (
	"fmt"
	"sync"
	"time"
)

const timestampDivider = 1000

const cleanUpInterval = 1
const cleanUpIntervalMultiplier = 1000

type LatencyEntry struct {
	D int64
	L float64
}

type Operations struct {
	latencies []LatencyEntry
}

type TelemetryManager struct {
	sync.RWMutex

	operations map[string][]LatencyEntry

	ctx Context

	cleanUpTimer *time.Ticker

	maxLatencyDataAge    int
	ExitTelemetryManager chan bool
}

func newTelemetryManager(maxLatencyDataAge int, ctx Context) *TelemetryManager {
	manager := &TelemetryManager{
		maxLatencyDataAge:    maxLatencyDataAge,
		operations:           make(map[string][]LatencyEntry),
		ctx:                  ctx,
		ExitTelemetryManager: make(chan bool),
	}

	go manager.startCleanUpTimer()

	return manager
}

func (m *TelemetryManager) OperationLatency() map[string]string {
	operationLatencies := make(map[string]string)

	//var ops map[string][]LatencyEntry
	m.RLock()

	for endpointName, _ := range m.operations {
		queryKey := fmt.Sprintf("l_%s", endpointName)

		endpointAverageLatency := averageLatencyFromData(
			m.operations[endpointName])

		if endpointAverageLatency > 0 {
			operationLatencies[queryKey] = fmt.Sprint(endpointAverageLatency)
		}
	}
	m.RUnlock()

	return operationLatencies
}

func (m *TelemetryManager) StoreLatency(latency float64, t OperationType) {
	if latency > float64(0) && t != PNSubscribeOperation {
		endpointName := telemetryEndpointNameForOperation(t)

		storeTimestamp := time.Now().Unix()

		m.Lock()
		m.operations[endpointName] = append(m.operations[endpointName], LatencyEntry{
			D: storeTimestamp,
			L: latency,
		})
		m.Unlock()
	}
}

func (m *TelemetryManager) CleanUpTelemetryData() {
	currentTimestamp := time.Now().Unix()

	m.Lock()
	for endpoint, latencies := range m.operations {
		index := 0

		for _, latency := range latencies {
			if currentTimestamp-latency.D > int64(m.maxLatencyDataAge) {
				m.operations[endpoint] = append(m.operations[endpoint][:index],
					m.operations[endpoint][index+1:]...)
				continue
			}
			index++
		}

		if len(m.operations[endpoint]) == 0 {
			delete(m.operations, endpoint)
		}
	}
	m.Unlock()
}

func (m *TelemetryManager) startCleanUpTimer() {
	m.cleanUpTimer = time.NewTicker(
		time.Duration(
			cleanUpInterval*cleanUpIntervalMultiplier) * time.Millisecond)

	go func() {
		for {
			m.RLock()
			timerCh := m.cleanUpTimer.C
			m.RUnlock()
			select {
			case <-timerCh:
				m.CleanUpTelemetryData()
			case <-m.ctx.Done():
				m.cleanUpTimer.Stop()
				return
			case <-m.ExitTelemetryManager:
				m.cleanUpTimer.Stop()
				return
			}
		}
	}()
}

func telemetryEndpointNameForOperation(t OperationType) string {
	var endpoint string

	switch t {
	case PNPublishOperation:
		endpoint = "pub"
		break
	case PNHistoryOperation:
		fallthrough
	case PNDeleteMessagesOperation:
		endpoint = "hist"
		break
	case PNUnsubscribeOperation:
		fallthrough
	case PNWhereNowOperation:
		fallthrough
	case PNHereNowOperation:
		fallthrough
	case PNHeartBeatOperation:
		fallthrough
	case PNSetStateOperation:
		fallthrough
	case PNGetStateOperation:
		endpoint = "pres"
		break
	case PNAddChannelsToChannelGroupOperation:
		fallthrough
	case PNRemoveChannelFromChannelGroupOperation:
		fallthrough
	case PNChannelsForGroupOperation:
		fallthrough
	case PNRemoveGroupOperation:
		endpoint = "cg"
		break
	case PNAccessManagerRevoke:
		fallthrough
	case PNAccessManagerGrant:
		endpoint = "pam"
		break
	default:
		endpoint = "time"
		break
	}

	return endpoint
}

func averageLatencyFromData(endpointLatencies []LatencyEntry) float64 {
	var totalLatency float64

	for _, latency := range endpointLatencies {
		totalLatency += latency.L
	}

	return totalLatency / float64(len(endpointLatencies))
}
