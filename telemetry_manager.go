package pubnub

import (
	"fmt"
	"sync"
	"time"
)

const timestampDivider = 1000

const cleanUpInterval = 1
const cleanUpIntervalMultiplier = 1000

// LatencyEntry is the struct to store the timestamp and latency values.
type LatencyEntry struct {
	D int64
	L float64
}

// Operations is the struct to store the latency values of different operations.
type Operations struct {
	latencies []LatencyEntry
}

// TelemetryManager is the struct to store the Telemetry details.
type TelemetryManager struct {
	sync.RWMutex

	operations map[string][]LatencyEntry

	ctx Context

	cleanUpTimer *time.Ticker

	maxLatencyDataAge int
	IsRunning         bool
}

func newTelemetryManager(maxLatencyDataAge int, ctx Context) *TelemetryManager {
	manager := &TelemetryManager{
		maxLatencyDataAge: maxLatencyDataAge,
		operations:        make(map[string][]LatencyEntry),
		ctx:               ctx,
	}

	go manager.startCleanUpTimer()

	return manager
}

// OperationLatency returns a map of the stored latencies by operation.
func (m *TelemetryManager) OperationLatency() map[string]string {
	operationLatencies := make(map[string]string)

	//var ops map[string][]LatencyEntry
	m.RLock()

	for endpointName := range m.operations {
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

// StoreLatency stores the latency values of the different operations.
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

// CleanUpTelemetryData cleans up telemetry data of all operations.
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
	m.ctx.Done()
	m.Unlock()
}

func (m *TelemetryManager) startCleanUpTimer() {
	m.cleanUpTimer = time.NewTicker(
		time.Duration(
			cleanUpInterval*cleanUpIntervalMultiplier) * time.Millisecond)

	go func() {
	CleanUpTimerLabel:
		for {
			timerCh := m.cleanUpTimer.C

			select {
			case <-timerCh:
				m.CleanUpTelemetryData()
			case <-m.ctx.Done():
				m.cleanUpTimer.Stop()
				break CleanUpTimerLabel
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
	case PNMessageCountsOperation:
		endpoint = "mc"
		break
	case PNHistoryOperation:
		fallthrough
	case PNFetchMessagesOperation:
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
	case PNAccessManagerGrantToken:
		endpoint = "pamv3"
		break
	case PNSignalOperation:
		endpoint = "sig"
		break
	case PNGetMessageActionsOperation:
		fallthrough
	case PNAddMessageActionsOperation:
		fallthrough
	case PNRemoveMessageActionsOperation:
		endpoint = "msga"
		break
	case PNHistoryWithActionsOperation:
		endpoint = "hist"
		break
	case PNCreateUserOperation:
		fallthrough
	case PNGetUsersOperation:
		fallthrough
	case PNGetUserOperation:
		fallthrough
	case PNUpdateUserOperation:
		fallthrough
	case PNDeleteUserOperation:
		fallthrough
	case PNGetSpaceOperation:
		fallthrough
	case PNGetSpacesOperation:
		fallthrough
	case PNCreateSpaceOperation:
		fallthrough
	case PNDeleteSpaceOperation:
		fallthrough
	case PNUpdateSpaceOperation:
		fallthrough
	case PNGetMembershipsOperation:
		fallthrough
	case PNGetChannelMembersOperation:
		fallthrough
	case PNManageMembershipsOperation:
		fallthrough
	case PNManageMembersOperation:
		fallthrough
	case PNSetChannelMembersOperation:
		fallthrough
	case PNSetMembershipsOperation:
		fallthrough
	case PNRemoveChannelMetadataOperation:
		fallthrough
	case PNRemoveUUIDMetadataOperation:
		fallthrough
	case PNGetAllChannelMetadataOperation:
		fallthrough
	case PNGetAllUUIDMetadataOperation:
		fallthrough
	case PNGetUUIDMetadataOperation:
		fallthrough
	case PNRemoveMembershipsOperation:
		fallthrough
	case PNRemoveChannelMembersOperation:
		fallthrough
	case PNSetUUIDMetadataOperation:
		fallthrough
	case PNGetChannelMetadataOperation:
		fallthrough
	case PNSetChannelMetadataOperation:
		endpoint = "obj"
		break
	case PNDeleteFileOperation:
		fallthrough
	case PNDownloadFileOperation:
		fallthrough
	case PNGetFileURLOperation:
		fallthrough
	case PNListFilesOperation:
		fallthrough
	case PNSendFileOperation:
		fallthrough
	case PNSendFileToS3Operation:
		fallthrough
	case PNPublishFileMessageOperation:
		endpoint = "file"
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
