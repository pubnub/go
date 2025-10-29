package pubnub

import (
	"testing"
	"time"

	"github.com/pubnub/go/v8/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestExponentialExhaustion(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()

	// Return error status to trigger reconnection logic
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/time/0",
		Query:              "",
		ResponseBody:       `{"error": "simulated network failure"}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 500,
	})

	config := NewConfigWithUserId(UserId(GenerateUUID()))

	pn := NewPubNub(config)
	pn.Config.MaximumReconnectionRetries = 2
	pn.Config.NonSubscribeRequestTimeout = 2
	pn.Config.ConnectTimeout = 2
	pn.Config.PNReconnectionPolicy = PNExponentialPolicy

	pn.SetClient(interceptor.GetClient())
	t1 := time.Now()
	r := newReconnectionManager(pn)
	reconnectionExhausted := false
	r.HandleOnMaxReconnectionExhaustion(func() {
		reconnectionExhausted = true
	})

	// Use timeout to prevent test hanging
	done := make(chan bool, 1)
	go func() {
		r.startHeartbeatTimer()
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(20 * time.Second):
		r.stopHeartbeatTimer()
		<-done
	}

	t2 := time.Now()
	diff := t2.Unix() - t1.Unix()

	assert.True((diff >= 1) && (diff <= 4), "Expected 1-4 seconds, got %d", diff)
	assert.True(reconnectionExhausted)
	r.stopHeartbeatTimer()
}

func TestLinearExhaustion(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()

	// Return error status to trigger reconnection logic
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/time/0",
		Query:              "",
		ResponseBody:       `{"error": "simulated network failure"}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 500,
	})

	config := NewConfigWithUserId(UserId(GenerateUUID()))

	pn := NewPubNub(config)
	pn.Config.MaximumReconnectionRetries = 1
	pn.Config.PNReconnectionPolicy = PNLinearPolicy
	pn.SetClient(interceptor.GetClient())
	t1 := time.Now()
	r := newReconnectionManager(pn)
	reconnectionExhausted := false
	r.HandleOnMaxReconnectionExhaustion(func() {
		reconnectionExhausted = true
	})

	// Use timeout to prevent test hanging
	done := make(chan bool, 1)
	go func() {
		r.startHeartbeatTimer()
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		r.stopHeartbeatTimer()
		<-done
	}

	t2 := time.Now()
	diff := t2.Unix() - t1.Unix()

	assert.True((diff >= 0) && (diff <= 3), "Expected 0-3 seconds, got %d", diff)
	assert.True(reconnectionExhausted)
	r.stopHeartbeatTimer()
}

func TestReconnect(t *testing.T) {
	assert := assert.New(t)

	config := NewConfigWithUserId(UserId(GenerateUUID()))
	pn := NewPubNub(config)
	pn.Config.MaximumReconnectionRetries = 1
	pn.Config.PNReconnectionPolicy = PNLinearPolicy
	r := newReconnectionManager(pn)
	r.FailedCalls = 1
	reconnected := false
	doneReconnected := make(chan bool)
	r.HandleReconnection(func() {
		reconnected = true
		doneReconnected <- true
	})
	go r.startHeartbeatTimer()
	<-doneReconnected
	assert.True(reconnected)
	r.stopHeartbeatTimer()
}
