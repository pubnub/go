package pubnub

import (
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExponentialExhaustion(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/time/0",
		Query:              "",
		ResponseBody:       `[15078947309567840]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	config := NewConfig()

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

	r.startHeartbeatTimer()
	t2 := time.Now()
	diff := t2.Unix() - t1.Unix()
	assert.True((diff >= 11) && (diff <= 12))
	assert.True(reconnectionExhausted)
	r.stopHeartbeatTimer()
}

func TestLinearExhaustion(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/time/0",
		Query:              "",
		ResponseBody:       `[15078947309567840]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	config := NewConfig()

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

	r.startHeartbeatTimer()
	t2 := time.Now()
	diff := t2.Unix() - t1.Unix()
	assert.True((diff >= 10) && (diff <= 11))
	assert.True(reconnectionExhausted)
	r.stopHeartbeatTimer()
}

func TestReconnect(t *testing.T) {
	assert := assert.New(t)

	config := NewConfig()

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
