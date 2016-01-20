// Package tests has the unit tests of package messaging.
// pubnubWildcardSubscribe.go contains the tests related to the Group
// Subscribe requests on pubnub Api
package tests

import (
	"encoding/json"
	"fmt"
	"github.com/pubnub/go/messaging"
	"github.com/stretchr/testify/assert"
	// "strings"
	// "os"
	"testing"
)

// TestWildcardSubscribeEnd prints a message on the screen to mark the beginning of
// subscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestWildcardSubscribeStart(t *testing.T) {
	PrintTestMessage("==========Wildcard Subscribe tests start==========")
}

func TestWildcardSubscriptionConnectedAndUnsubscribedSingle(t *testing.T) {
	//messaging.SetLogOutput(os.Stderr)
	assert := assert.New(t)
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	r := GenRandom()
	major := fmt.Sprintf("testChannel_sub_%d", r.Intn(20))
	wildcard := fmt.Sprintf("%s.*", major)

	subscribeSuccessChannel := make(chan []byte)
	subscribeErrorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(wildcard, "",
		subscribeSuccessChannel, false, subscribeErrorChannel)
	select {
	case msg := <-subscribeSuccessChannel:
		val := string(msg)
		assert.Equal(fmt.Sprintf(
			"[1, \"Subscription to channel '%s' connected\", \"%s\"]",
			wildcard, wildcard), val)
	case err := <-subscribeErrorChannel:
		assert.Fail(string(err))
	}

	go pubnubInstance.Unsubscribe(wildcard, successChannel, errorChannel)
	select {
	case msg := <-successChannel:
		val := string(msg)
		assert.Equal(fmt.Sprintf(
			"[1, \"Subscription to channel '%s' unsubscribed\", \"%s\"]",
			wildcard, wildcard), val)
	case err := <-errorChannel:
		assert.Fail(string(err))
	}

	pubnubInstance.CloseExistingConnection()
}

func TestWildcardSubscriptionMessage(t *testing.T) {
	//messaging.SetLogOutput(os.Stderr)
	assert := assert.New(t)
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	r := GenRandom()
	major := fmt.Sprintf("testChannel_sub_%d", r.Intn(20))
	minor := fmt.Sprintf("testChannel_sub_%d", r.Intn(20))
	channel := fmt.Sprintf("%s.%s", major, minor)
	wildcard := fmt.Sprintf("%s.*", major)

	subscribeSuccessChannel := make(chan []byte)
	subscribeErrorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	await := make(chan bool)

	go pubnubInstance.Subscribe(wildcard, "",
		subscribeSuccessChannel, false, subscribeErrorChannel)
	ExpectConnectedEvent(t, wildcard, "", subscribeSuccessChannel)
	go func() {
		select {
		case message := <-subscribeSuccessChannel:
			var msg []interface{}

			err := json.Unmarshal(message, &msg)
			if err != nil {
				assert.Fail(err.Error())
			}

			assert.Contains(string(message), "hey")
			assert.Equal(channel, msg[2].(string))
			assert.Equal(wildcard, msg[3].(string))
			await <- true
		case err := <-subscribeErrorChannel:
			assert.Fail(string(err))
		case <-messaging.SubscribeTimeout():
			assert.Fail("Subscribe timeout")
		}
	}()

	go pubnubInstance.Publish(channel, "hey", successChannel, errorChannel)
	select {
	case <-successChannel:
	case err := <-errorChannel:
		assert.Fail(string(err))
	}

	<-await

	go pubnubInstance.Unsubscribe(wildcard, successChannel, errorChannel)
	ExpectUnsubscribedEvent(t, wildcard, "", successChannel)

	pubnubInstance.CloseExistingConnection()
}

// TODO test presence

// TestWildcardSubscribeEnd prints a message on the screen to mark the end of
// subscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestWildcardSubscribeEnd(t *testing.T) {
	PrintTestMessage("==========Wildcard Subscribe tests end==========")
}
