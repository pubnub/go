// Package tests has the unit tests of package messaging.
// pubnubUnsubscribe_test.go contains the tests related to the Unsubscribe requests on pubnub Api
package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pubnub/go/messaging"
	"github.com/stretchr/testify/assert"
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
	assert := assert.New(t)
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")

	currentTime := time.Now()
	channel := "testChannel" + currentTime.Format("20060102150405")

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.Unsubscribe(channel, successChannel, errorChannel)
	select {
	case <-successChannel:
		assert.Fail("Success unsubscribe response while expecting an error")
	case err := <-errorChannel:
		assert.Contains(string(err), "not subscribed")
		assert.Contains(string(err), channel)
	case <-timeout():
		assert.Fail("Unsubscribe request timeout")
	}
}

// TestUnsubscribe will subscribe to a pubnub channel and then send an unsubscribe request
// The response should contain 'unsubscribed'
func TestUnsubscribeChannel(t *testing.T) {
	assert := assert.New(t)

	stop, sleep := NewVCRBoth(
		"fixtures/unsubscribe/channel", []string{"uuid"})
	defer stop()

	channel := "Channel_UnsubscribeChannel"
	uuid := "UUID_UnsubscribeChannel"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, uuid)

	subscribeSuccessChannel := make(chan []byte)
	subscribeErrorChannel := make(chan []byte)
	unSubscribeSuccessChannel := make(chan []byte)
	unSubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channel, "", subscribeSuccessChannel,
		false, subscribeErrorChannel)
	select {
	case msg := <-subscribeSuccessChannel:
		val := string(msg)
		assert.Equal(val, fmt.Sprintf(
			"[1, \"Subscription to channel '%s' connected\", \"%s\"]",
			channel, channel))
	case err := <-subscribeErrorChannel:
		assert.Fail(string(err))
	}

	sleep(2)

	go pubnubInstance.Unsubscribe(channel, unSubscribeSuccessChannel,
		unSubscribeErrorChannel)
	select {
	case msg := <-unSubscribeSuccessChannel:
		val := string(msg)
		assert.Equal(val, fmt.Sprintf(
			"[1, \"Subscription to channel '%s' unsubscribed\", \"%s\"]",
			channel, channel))
	case err := <-unSubscribeErrorChannel:
		assert.Fail(string(err))
	}

	select {
	case ev := <-unSubscribeSuccessChannel:
		var event messaging.PresenceResonse

		err := json.Unmarshal(ev, &event)
		if err != nil {
			assert.Fail(err.Error())
		}

		assert.Equal("leave", event.Action)
		assert.Equal(200, event.Status)
	case err := <-unSubscribeErrorChannel:
		assert.Fail(string(err))
	}
}

func TestUnsubscribeNetworkError(t *testing.T) {
	assert := assert.New(t)

	stop, _ := NewVCRBoth(
		"fixtures/unsubscribe/networkError", []string{"uuid"})
	defer stop()
	messaging.SetNonSubscribeTransport(abortedTransport)

	channel := "Channel_UnsubscribeNetError"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	subscribeSuccess := make(chan []byte)
	subscribeError := make(chan []byte)

	go pubnubInstance.Subscribe(channel, "", subscribeSuccess, false, subscribeError)
	ExpectConnectedEvent(t, channel, "", subscribeSuccess, subscribeError)

	successGet := make(chan []byte)
	errorGet := make(chan []byte)

	go pubnubInstance.Unsubscribe(channel, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Contains(string(value), "unsubscribed")
		assert.Contains(string(value), channel)
	case err := <-errorGet:
		assert.Fail("Error response while expecting success", string(err))
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), abortedTransport.PnMessage)
		assert.Contains(string(err), channel)
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}

	messaging.SetNonSubscribeTransport(nil)
}

func TestGroupUnsubscribeNetworkError(t *testing.T) {
	assert := assert.New(t)

	stop, sleep := NewVCRBoth(
		"fixtures/unsubscribe/groupNetworkError", []string{"uuid"})
	defer stop()

	group := "Channel_GroupUnsubscribeNetError"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	createChannelGroups(pubnubInstance, []string{group})
	defer removeChannelGroups(pubnubInstance, []string{group})

	sleep(2)

	subscribeSuccess := make(chan []byte)
	subscribeError := make(chan []byte)

	go pubnubInstance.ChannelGroupSubscribe(group, subscribeSuccess, subscribeError)
	ExpectConnectedEvent(t, "", group, subscribeSuccess, subscribeError)

	successGet := make(chan []byte)
	errorGet := make(chan []byte)

	messaging.SetNonSubscribeTransport(abortedTransport)

	go pubnubInstance.ChannelGroupUnsubscribe(group, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Contains(string(value), "unsubscribed")
		assert.Contains(string(value), group)
	case err := <-errorGet:
		assert.Fail("Error response while expecting success", string(err))
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), abortedTransport.PnMessage)
		assert.Contains(string(err), group)
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}

	messaging.SetNonSubscribeTransport(nil)
}

// TestUnsubscribeEnd prints a message on the screen to mark the end of
// unsubscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestUnsubscribeEnd(t *testing.T) {
	PrintTestMessage("==========Unsubscribe tests end==========")
}
