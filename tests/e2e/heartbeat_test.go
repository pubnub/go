package e2e

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

//TODO ENABLE
func HeartbeatTimeoutEvent(t *testing.T) {
	assert := assert.New(t)
	ch := randomized("hb-te")
	emitterUUID := randomized("emitter")

	var wg sync.WaitGroup
	wg.Add(2)

	donePresenceConnect := make(chan bool)
	doneJoin := make(chan bool)
	doneTimeout := make(chan bool)
	errChan := make(chan string)

	configEmitter := configCopy()
	configEmitter.SetPresenceTimeout(6)

	configPresenceListener := configCopy()

	configEmitter.UUID = emitterUUID
	configPresenceListener.UUID = randomized("listener")

	pn := pubnub.NewPubNub(configEmitter)
	pnPresenceListener := pubnub.NewPubNub(configPresenceListener)

	listenerEmitter := pubnub.NewListener()
	listenerPresenceListener := pubnub.NewListener()

	exitListener := make(chan bool)

	// emitter
	go func() {
	ExitLabel:
		for {
			select {
			case status := <-listenerEmitter.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					wg.Done()
					return
				}
			case <-listenerEmitter.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listenerEmitter.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			case <-exitListener:
				break ExitLabel
			}
		}
	}()

	// listener
	go func() {
	ExitLabel:
		for {
			select {
			case status := <-listenerPresenceListener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					donePresenceConnect <- true
				}
			case message := <-listenerPresenceListener.Message:
				errChan <- fmt.Sprintf("Unexpected message: %s",
					message.Message)
			case presence := <-listenerPresenceListener.Presence:
				// ignore join event of presence listener
				if presence.UUID == configPresenceListener.UUID {
					continue
				}

				assert.Equal(ch, presence.Channel)

				if presence.Event == "timeout" {
					assert.Equal(configEmitter.UUID, presence.UUID)
					doneTimeout <- true
					return
				} else if presence.Event == "join" {
					assert.Equal("join", presence.Event)
					assert.Equal(configEmitter.UUID, presence.UUID)
					wg.Done()
				}
			case <-exitListener:
				break ExitLabel
			}
		}
	}()

	pn.AddListener(listenerEmitter)
	pnPresenceListener.AddListener(listenerPresenceListener)

	pnPresenceListener.Subscribe().
		Channels([]string{ch}).
		WithPresence(true).
		Execute()

	select {
	case <-donePresenceConnect:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Subscribe().
		Channels([]string{ch}).
		Execute()

	go func() {
		wg.Wait()
		doneJoin <- true
	}()

	select {
	case <-doneJoin:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	cl := pubnub.NewHTTP1Client(15, 15, 20)
	cl.Transport = fakeTransport{
		Status:     "200 OK",
		StatusCode: 200,
		// WARNING: can be read only once
		Body: ioutil.NopCloser(strings.NewReader("Hijacked response")),
	}
	pn.SetClient(cl)

	defer pn.UnsubscribeAll()

	select {
	case <-doneTimeout:
	case <-time.After(8 * time.Second):
		assert.Fail("No timeout event received in 8 seconds")
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	exitListener <- true
}

func TestHeartbeatStubbedRequest(t *testing.T) {
	ch := randomized("ch-hsr")
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub-key/%s/channel/", config.SubscribeKey) + ch + "/heartbeat",
		Query:              "heartbeat=6",
		ResponseBody:       "{\"status\": 200, \"message\": \"OK\", \"service\": \"Presence\"}",
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub-key/%s/channel/", config.SubscribeKey) + ch + "/leave",
		Query:              "",
		ResponseBody:       "{\"status\": 200, \"message\": \"OK\", \"service\": \"Presence\", \"action\": \"leave\"}",
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	doneConnect := make(chan bool)
	doneHeartbeat := make(chan bool)
	errChan := make(chan string)

	config := configCopy()
	config.SetPresenceTimeout(6)

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()
	pn.AddListener(listener)
	exitListener := make(chan bool)

	go func() {
	ExitLabel:
		for {
			select {
			case status := <-listener.Status:
				switch status.Operation {
				case pubnub.PNHeartBeatOperation:
					doneHeartbeat <- true
				}

				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneConnect <- true
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			case <-exitListener:
				break ExitLabel
			}
		}
	}()

	pn.Subscribe().
		Channels([]string{ch}).
		Execute()

	select {
	case <-doneConnect:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.SetClient(interceptor.GetClient())

	select {
	case <-doneHeartbeat:
	case <-time.After(10 * time.Second):
		assert.Fail("Heartbeat status was expected")
	}

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	exitListener <- true
}

// Test triggers BadRequestCategory in subscription.Status channel
// for failed HB call
func TestHeartbeatRequestWithError(t *testing.T) {
	ch := randomized("ch-hrwe")
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub-key/%s/channel/"+ch+"/heartbeat", config.SubscribeKey),
		Query:              "heartbeat=6",
		ResponseBody:       "{\"status\": 404, \"message\": \"Not Found\", \"error\": \"1\", \"service\": \"Presence\"}",
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 404,
	})

	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub-key/%s/channel/"+ch+"/leave", config.SubscribeKey),
		Query:              "",
		ResponseBody:       "{\"status\": 200, \"message\": \"OK\", \"service\": \"Presence\", \"action\": \"leave\"}",
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	doneConnect := make(chan bool)
	doneHeartbeat := make(chan bool)
	errChan := make(chan string)

	config := configCopy()
	config.SetPresenceTimeout(6)

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()
	pn.AddListener(listener)
	exitListener := make(chan bool)

	go func() {
	ExitLabel:
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneConnect <- true
				case pubnub.PNBadRequestCategory:
					doneHeartbeat <- true
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			case <-exitListener:
				break ExitLabel
			}
		}
	}()

	pn.Subscribe().
		Channels([]string{ch}).
		Execute()

	select {
	case <-doneConnect:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.SetClient(interceptor.GetClient())

	select {
	case <-doneHeartbeat:
	case <-time.After(10 * time.Second):
		assert.Fail("Heartbeat status was expected")
	}

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	exitListener <- true
}

// NOTICE: snippet for manual hb testing
// - subscribe 'first'
// - unsubscribeAll
// - subscribe 'first'
// - subscribe 'second'
// - unsubscribe 'first', 'second'
func xTestHeartbeatRandomizedBehaviour(t *testing.T) {
	assert := assert.New(t)
	first := "first"
	second := "second"
	emitterUUID := randomized("emitter")

	var wg sync.WaitGroup
	wg.Add(2)

	doneJoin := make(chan bool)
	doneTimeout := make(chan bool)
	errChan := make(chan string)

	configEmitter := configCopy()
	configEmitter.SetPresenceTimeout(6)

	configEmitter.UUID = emitterUUID

	pn := pubnub.NewPubNub(configEmitter)

	listenerEmitter := pubnub.NewListener()
	exitListener := make(chan bool)

	// emitter
	go func() {
	ExitLabel:
		for {
			select {
			case status := <-listenerEmitter.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneJoin <- true
					return
				}
			case <-listenerEmitter.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listenerEmitter.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			case <-exitListener:
				break ExitLabel

			}
		}
	}()

	pn.AddListener(listenerEmitter)

	pn.Subscribe().
		Channels([]string{first}).
		Execute()

	go func() {
		doneJoin <- true
	}()

	select {
	case <-doneJoin:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	//log.Println("Sleeping 8s")
	time.Sleep(8 * time.Second)

	pn.UnsubscribeAll()

	//log.Println("Unsubscribed")
	//log.Println("Sleeping 8s")
	time.Sleep(8 * time.Second)

	pn.Subscribe().
		Channels([]string{first}).
		Execute()

	//log.Println("Subscribed again")
	//log.Println("Sleeping 8s")
	time.Sleep(8 * time.Second)

	pn.Subscribe().
		Channels([]string{second}).
		Execute()

	//log.Println("Subsccribed to another channel again")
	//log.Println("Sleeping 8s")
	time.Sleep(8 * time.Second)

	pn.Unsubscribe().
		Channels([]string{first, second}).
		Execute()

	//log.Println("Unsubscribed")
	//log.Println("Sleeping 8s")
	time.Sleep(8 * time.Second)

	select {
	case <-doneTimeout:
	case <-time.After(8 * time.Second):
		assert.Fail("No timeout event received in 8 seconds")
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	exitListener <- true
}
