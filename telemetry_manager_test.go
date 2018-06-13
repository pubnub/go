package pubnub

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAverageLatency(t *testing.T) {
	assert := assert.New(t)

	endpointLatencies := []LatencyEntry{
		LatencyEntry{D: int64(100), L: float64(10)},
		LatencyEntry{D: int64(100), L: float64(20)},
		LatencyEntry{D: int64(100), L: float64(30)},
		LatencyEntry{D: int64(100), L: float64(40)},
		LatencyEntry{D: int64(100), L: float64(50)}}

	averageLatency := averageLatencyFromData(endpointLatencies)
	assert.Equal(float64(30), averageLatency)
}

func TestCleanUp(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := contextWithCancel(backgroundContext)
	manager := newTelemetryManager(1, ctx)

	for i := 0; i < 10; i++ {
		manager.StoreLatency(float64(i), PNPublishOperation)
	}

	// await for store timestamp expired
	time.Sleep(2 * time.Second)

	manager.CleanUpTelemetryData()

	assert.Equal(0, len(manager.OperationLatency()))
}

func TestValidQueries(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := contextWithCancel(backgroundContext)
	manager := newTelemetryManager(60, ctx)

	manager.StoreLatency(float64(1), PNPublishOperation)
	manager.StoreLatency(float64(2), PNPublishOperation)
	manager.StoreLatency(float64(3), PNPublishOperation)

	manager.StoreLatency(float64(4), PNHistoryOperation)
	manager.StoreLatency(float64(5), PNHistoryOperation)
	manager.StoreLatency(float64(6), PNHistoryOperation)

	manager.StoreLatency(float64(7), PNRemoveGroupOperation)
	manager.StoreLatency(float64(8), PNRemoveGroupOperation)
	manager.StoreLatency(float64(9), PNRemoveGroupOperation)

	queries := manager.OperationLatency()

	assert.Equal("2", queries["l_pub"])
	assert.Equal("5", queries["l_hist"])
	assert.Equal("8", queries["l_cg"])
}
