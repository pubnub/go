// Package tests has the unit tests of package messaging.
// pubnubSubscribe_test.go contains the tests related to the Subscribe requests on pubnub Api
package tests

import (
	"fmt"
	"github.com/pubnub/go/messaging"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

// TestSubscribeStart prints a message on the screen to mark the beginning of
// subscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestSubscribeStart(t *testing.T) {
	PrintTestMessage("==========Subscribe tests start==========")
}

func TestChannelSubscriptionWithTimetoken(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	r := GenRandom()
	channel := fmt.Sprintf("testChannel_sub_%d", r.Intn(20))
	timestamp := ""

	successChannel, errorChannel, eventsChannel :=
		messaging.CreateSubscriptionChannels()

	timeSuccessChannel := make(chan []byte)
	timeErrorChannel := make(chan []byte)

	go pubnubInstance.GetTime(timeSuccessChannel, timeErrorChannel)

	select {
	case response := <-timeSuccessChannel:
		timestamp = strings.Trim(string(response), "[]")
	case err := <-timeErrorChannel:
		assert.Fail(t, "Error while getting server timestamp:", string(err))
	}

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	publishSuccessChannel := make(chan []byte)
	publishErrorChannel := make(chan []byte)

	go pubnubInstance.SubscribeWithTimetoken(channel, timestamp, successChannel, errorChannel, eventsChannel)
	go pubnubInstance.Publish(channel, 123, publishSuccessChannel, publishErrorChannel)
	go func() {
		select {
		case <-publishSuccessChannel:
		case <-publishErrorChannel:
		}
	}()

	<-timeouts(1)

	select {
	case <-eventsChannel:
	case <-errorChannel:
		assert.Fail(t, "Received Error first instead of message")
	case msg := <-successChannel:
		assert.Equal(t, "123", string(msg.Data))
	case <-time.After(time.Second * time.Duration(testTimeout)):
		assert.Fail(t, "Timeout occured")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, channel, "", eventsChannel)

	pubnubInstance.CloseExistingConnection()
}

// TestSubscriptionConnectStatus sends out a subscribe request to a pubnub
// channel and validates the response for the connect status.
func TestChannelSubscriptionConnectAndDisconnect(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	r := GenRandom()
	channel := fmt.Sprintf("testChannel_sub_%d", r.Intn(20))

	successChannel, errorChannel, eventsChannel :=
		messaging.CreateSubscriptionChannels()

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)

	select {
	case event := <-eventsChannel:
		assert.Equal(t, messaging.ConnectionConnected, event.Action)
		assert.Equal(t, channel, event.Channel)
		assert.Equal(t, messaging.ChannelResponse, event.Type)
	case <-errorChannel:
		assert.Fail(t, "Received Error first instead of 'connected' event")
	case <-successChannel:
		assert.Fail(t, "Received Message first instead of 'connected' event")
	case <-time.After(time.Second * time.Duration(testTimeout)):
		assert.Fail(t, "Timeout occured")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)

	select {
	case event := <-eventsChannel:
		assert.Equal(t, messaging.ConnectionDisconnected, event.Action)
		assert.Equal(t, channel, event.Channel)
		assert.Equal(t, messaging.ChannelResponse, event.Type)
	case <-errorChannel:
		assert.Fail(t, "Received Error first instead of 'connected' event")
	case <-successChannel:
		assert.Fail(t, "Received Message first instead of 'connected' event")
	case <-time.After(time.Second * time.Duration(testTimeout)):
		assert.Fail(t, "Timeout occured")
	}

	pubnubInstance.CloseExistingConnection()
}

func TestChannelSubscriptionMessage(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	r := GenRandom()
	channel := fmt.Sprintf("testChannel_sub_%d", r.Intn(20))

	await := make(chan bool)
	await2 := make(chan bool)

	successChannel, errorChannel, eventsChannel :=
		messaging.CreateSubscriptionChannels()

	publishSuccessChannel := make(chan []byte)
	publishErrorChannel := make(chan []byte)

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go func() {
		var messages []string

		for {
			select {
			case msg := <-successChannel:
				fmt.Println("got:", string(msg.Data))
				messages = append(messages, string(msg.Data))
			case err := <-errorChannel:
				assert.Fail(t, "Subscribe error", err.Error())
			case <-eventsChannel:
				await <- true
			case <-timeout():
				assert.Fail(t, "For looop timed out")
				break
			}

			if len(messages) == 2 {
				break
			}
		}

		assert.Equal(t, []string{"\"hello\"", "\"blah\""}, messages)
		await2 <- true
	}()

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)

	<-await
	fmt.Println("passed await")

	go pubnubInstance.Publish(channel, "hello", publishSuccessChannel, publishErrorChannel)
	select {
	case <-publishSuccessChannel:
	case err := <-publishErrorChannel:
		assert.Fail(t, "Publish error", string(err))
	case <-timeout():
		assert.Fail(t, "Publish#2 timed out")
	}

	go pubnubInstance.Publish(channel, "blah", publishSuccessChannel, publishErrorChannel)
	select {
	case <-publishSuccessChannel:
	case err := <-publishErrorChannel:
		assert.Fail(t, "Publish error", string(err))
	case <-timeout():
		assert.Fail(t, "Publish#2 timed out")
	}

	fmt.Println("messages published")
	<-await2

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, channel, "", eventsChannel)

	pubnubInstance.CloseExistingConnection()
}

// TestSubscriptionAlreadySubscribed sends out a subscribe request to a pubnub channel
// and when connected sends out another subscribe request
func TestChannelSubscriptionAlreadySubscribed(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")

	r := GenRandom()
	channel := fmt.Sprintf("testChannel_sub_%d", r.Intn(20))

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()
	successChannel2, errorChannel2, eventsChannel2 := messaging.CreateSubscriptionChannels()

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)
	ExpectConnectedEvent(t, channel, "", eventsChannel)

	go pubnubInstance.Subscribe(channel, successChannel2, errorChannel2, eventsChannel2)

	select {
	case <-eventsChannel2:
		assert.Fail(t, "Received Event first instead of error")
	case err := <-errorChannel2:
		assert.Contains(t, err.Error(), "already subscribed to")
		assert.Contains(t, err.Error(), channel)
	case <-successChannel2:
		assert.Fail(t, "Received Message first instead of error")
	case <-time.After(time.Second * time.Duration(testTimeout)):
		assert.Fail(t, "Timeout occured")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, channel, "", eventsChannel2)

	pubnubInstance.CloseExistingConnection()
}

func TestChannelSubscriptionConnectedReconnectedAndDisconnectMultipleChannelsInMultipleCalls(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")

	channels, channels2 := GenerateTwoRandomChannelStrings(3)

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()
	successChannel2, errorChannel2, eventsChannel2 := messaging.CreateSubscriptionChannels()

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channels, successChannel, errorChannel, eventsChannel)
	ExpectConnectedEvent(t, channels, "", eventsChannel)

	go pubnubInstance.Subscribe(channels2, successChannel2, errorChannel2, eventsChannel2)
	ExpectConnectedEvent(t, channels2, "", eventsChannel2)

	// TODO: test reconnect using CloseExistingConnection()

	go pubnubInstance.Unsubscribe(channels, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, channels, "", eventsChannel)

	go pubnubInstance.Unsubscribe(channels2, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, channels2, "", eventsChannel2)

	pubnubInstance.CloseExistingConnection()
}

func TestChannelGroupSubscriptionConnectAndDisconnect(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")

	r := GenRandom()
	group := fmt.Sprintf("testChannelGroup_sub_%d", r.Intn(20))

	createChannelGroups(pubnubInstance, []string{group})
	defer removeChannelGroups(pubnubInstance, []string{group})

	successChannel, errorChannel, eventsChannel :=
		messaging.CreateSubscriptionChannels()

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.ChannelGroupSubscribe(group, successChannel, errorChannel, eventsChannel)

	select {
	case event := <-eventsChannel:
		assert.Equal(t, messaging.ConnectionConnected, event.Action)
		assert.Equal(t, group, event.Group)
		assert.Equal(t, messaging.ChannelGroupResponse, event.Type)
	case err := <-errorChannel:
		assert.Fail(t, "Received Error first instead of 'connected' event ", err.Error())
	case <-successChannel:
		assert.Fail(t, "Received Message first instead of 'connected' event")
	case <-time.After(time.Second * time.Duration(testTimeout)):
		assert.Fail(t, "Timeout occured")
	}

	go pubnubInstance.ChannelGroupUnsubscribe(group, unsubscribeSuccessChannel, unsubscribeErrorChannel)

	select {
	case event := <-eventsChannel:
		assert.Equal(t, messaging.ConnectionDisconnected, event.Action)
		assert.Equal(t, group, event.Group)
		assert.Equal(t, messaging.ChannelGroupResponse, event.Type)
	case <-errorChannel:
		assert.Fail(t, "Received Error first instead of 'connected' event")
	case <-successChannel:
		assert.Fail(t, "Received Message first instead of 'connected' event")
	case <-time.After(time.Second * time.Duration(testTimeout)):
		assert.Fail(t, "Timeout occured")
	}

	pubnubInstance.CloseExistingConnection()
}

func TestChannelGroupSubscriptionAlreadySubscribed(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")

	r := GenRandom()
	group := fmt.Sprintf("testChannelGroup_sub_%d", r.Intn(20))

	createChannelGroups(pubnubInstance, []string{group})
	defer removeChannelGroups(pubnubInstance, []string{group})

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()
	successChannel2, errorChannel2, eventsChannel2 := messaging.CreateSubscriptionChannels()

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.ChannelGroupSubscribe(group, successChannel, errorChannel, eventsChannel)
	ExpectConnectedEvent(t, "", group, eventsChannel)

	go pubnubInstance.ChannelGroupSubscribe(group, successChannel2, errorChannel2, eventsChannel2)

	select {
	case <-eventsChannel2:
		assert.Fail(t, "Received Event first instead of error")
	case err := <-errorChannel2:
		assert.Contains(t, err.Error(), "already subscribed to")
		assert.Contains(t, err.Error(), group)
	case <-successChannel2:
		assert.Fail(t, "Received Message first instead of error")
	case <-time.After(time.Second * time.Duration(testTimeout)):
		assert.Fail(t, "Timeout occured")
	}

	go pubnubInstance.ChannelGroupUnsubscribe(group, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, "", group, eventsChannel2)

	pubnubInstance.CloseExistingConnection()
}

func TestChannelGroupSubscriptionConnectedReconnectedAndDisconnectMultipleChannelsInMultipleCalls(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")

	groups, groups2 := GenerateTwoRandomChannelStrings(3)

	createChannelGroups(pubnubInstance, strings.Split(groups, ","))
	createChannelGroups(pubnubInstance, strings.Split(groups2, ","))

	defer removeChannelGroups(pubnubInstance, strings.Split(groups, ","))
	defer removeChannelGroups(pubnubInstance, strings.Split(groups2, ","))

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()
	successChannel2, errorChannel2, eventsChannel2 := messaging.CreateSubscriptionChannels()

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go LogErrors(errorChannel)

	go pubnubInstance.ChannelGroupSubscribe(groups, successChannel, errorChannel, eventsChannel)
	ExpectConnectedEvent(t, "", groups, eventsChannel)

	go pubnubInstance.ChannelGroupSubscribe(groups2, successChannel2, errorChannel2, eventsChannel2)
	ExpectConnectedEvent(t, "", groups2, eventsChannel2)

	// TODO: test reconnect

	go pubnubInstance.ChannelGroupUnsubscribe(groups, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, "", groups, eventsChannel)

	go pubnubInstance.ChannelGroupUnsubscribe(groups2, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, "", groups2, eventsChannel2)

	pubnubInstance.CloseExistingConnection()
}

func TestChannelAndChannelGroupSubscriptionConnectAndDisconnect(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")

	r := GenRandom()
	channel := fmt.Sprintf("testChannel_sub_%d", r.Intn(20))
	group := fmt.Sprintf("testChannelGroup_sub_%d", r.Intn(20))

	createChannelGroups(pubnubInstance, []string{group})
	defer removeChannelGroups(pubnubInstance, []string{group})

	successChannel, errorChannel, eventsChannel :=
		messaging.CreateSubscriptionChannels()

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)
	go pubnubInstance.ChannelGroupSubscribe(group, successChannel, errorChannel, eventsChannel)

	ExpectConnectedEvent(t, channel, group, eventsChannel)

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	go pubnubInstance.ChannelGroupUnsubscribe(group, unsubscribeSuccessChannel, unsubscribeErrorChannel)

	ExpectDisconnectedEvent(t, channel, group, eventsChannel)

	pubnubInstance.CloseExistingConnection()
}

func TestChannelAndChannelGroupSubscriptionAlreadySubscribed(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")

	defer pubnubInstance.CloseExistingConnection()

	r := GenRandom()
	channel := fmt.Sprintf("testChannel_sub_%d", r.Intn(20))
	group := fmt.Sprintf("testChannelGroup_sub_%d", r.Intn(20))

	createChannelGroups(pubnubInstance, []string{group})

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()
	successChannel2, errorChannel2, eventsChannel2 := messaging.CreateSubscriptionChannels()

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)
	go pubnubInstance.ChannelGroupSubscribe(group, successChannel, errorChannel, eventsChannel)

	ExpectConnectedEvent(t, channel, group, eventsChannel)

	go pubnubInstance.Subscribe(channel, successChannel2, errorChannel2, eventsChannel2)

	select {
	case msg := <-eventsChannel2:
		assert.Fail(t, "Received Event first instead of error", msg)
	case err := <-errorChannel2:
		assert.Contains(t, err.Error(), "already subscribed to")
		assert.Contains(t, err.Error(), channel)
	case <-successChannel2:
		assert.Fail(t, "Received Message first instead of error")
	case <-time.After(time.Second * time.Duration(testTimeout)):
		assert.Fail(t, "Timeout occured")
	}

	go pubnubInstance.ChannelGroupSubscribe(group, successChannel2, errorChannel2, eventsChannel2)

	select {
	case msg := <-eventsChannel2:
		assert.Fail(t, "Received Event first instead of error", msg)
	case err := <-errorChannel2:
		assert.Contains(t, err.Error(), "already subscribed to")
		assert.Contains(t, err.Error(), group)
	case <-successChannel2:
		assert.Fail(t, "Received Message first instead of error")
	case <-time.After(time.Second * time.Duration(testTimeout)):
		assert.Fail(t, "Timeout occured")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	go pubnubInstance.ChannelGroupUnsubscribe(group, unsubscribeSuccessChannel, unsubscribeErrorChannel)

	ExpectDisconnectedEvent(t, channel, group, eventsChannel2)
	removeChannelGroups(pubnubInstance, []string{group})
}

//*********
// HELPERS
//*********

func createChannelGroups(pubnub *messaging.Pubnub, groups []string) {
	successChannel := make(chan []byte, 1)
	errorChannel := make(chan []byte, 1)

	for _, group := range groups {
		fmt.Println("Creating group", group)

		pubnub.ChannelGroupAddChannel(group, "adsf", successChannel, errorChannel)

		select {
		case <-successChannel:
			fmt.Println("Group created")
		case <-errorChannel:
			fmt.Println("Channel group creation error")
		case <-timeout():
			fmt.Println("Channel group creation timeout")
		}
	}
}

func removeChannelGroups(pubnub *messaging.Pubnub, groups []string) {
	successChannel := make(chan []byte, 1)
	errorChannel := make(chan []byte, 1)

	for _, group := range groups {
		fmt.Println("Removing group", group)

		pubnub.ChannelGroupRemoveGroup(group, successChannel, errorChannel)

		select {
		case <-successChannel:
			fmt.Println("Group removed")
		case <-errorChannel:
			fmt.Println("Channel group removal error")
		case <-timeout():
			fmt.Println("Channel group removal timeout")
		}
	}
}

// TestSubscribeEnd prints a message on the screen to mark the end of
// subscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestSubscribeEnd(t *testing.T) {
	PrintTestMessage("==========Subscribe tests end==========")
}
