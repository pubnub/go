// Package tests has the unit tests of package messaging.
// pubnubUnsubscribe_test.go contains the tests related to the Unsubscribe requests on pubnub Api
package tests

import (
	"encoding/json"
	"github.com/pubnub/go/messaging"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestUnsubscribeStart prints a message on the screen to mark the beginning of
// unsubscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestUnsubscribeStart(t *testing.T) {
	PrintTestMessage("==========Unsubscribe tests start==========")
}

// TestUnsubscribeNotSubscribed will try to unsubscribe a non subscribed pubnub channel.
// The response should contain 'not subscribed'
func TestUnsubscribeNotSubscribed(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")

	currentTime := time.Now()
	channel := "testChannel" + currentTime.Format("20060102150405")

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)

	select {
	case <-unsubscribeSuccessChannel:
		assert.Fail(t, "Received message on success callback while error expected")
	case err := <-unsubscribeErrorChannel:
		assert.Contains(t, string(err), "not subscribed")
		assert.Contains(t, string(err), "channel")
	case <-timeout():
		assert.Fail(t, "Timed out")
	}
}

// TestUnsubscribe will subscribe to a pubnub channel and then send an unsubscribe request
// The response should contain 'unsubscribed'
func TestUnsubscribe(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")

	channel := "testChannel"

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)
	ExpectConnectedEvent(t, channel, "", eventsChannel)
	go LogErrors(errorChannel)

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	go ExpectDisconnectedEvent(t, channel, "", eventsChannel)

	select {
	case msg := <-unsubscribeSuccessChannel:
		var response map[string]interface{}
		err := json.Unmarshal(msg, &response)
		if err != nil {
			assert.Fail(t, err.Error())
		}
		assert.Equal(t, "leave", response["action"])
		assert.Equal(t, "Presence", response["service"])
	case err := <-unsubscribeErrorChannel:
		assert.Fail(t, string(err))
	case <-timeout():
		assert.Fail(t, "Timed out")
	}
}

// TestUnsubscribeEnd prints a message on the screen to mark the end of
// unsubscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestUnsubscribeEnd(t *testing.T) {
	PrintTestMessage("==========Unsubscribe tests end==========")
}
