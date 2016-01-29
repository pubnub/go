// Package tests has the unit tests of package messaging.
// pubnubGroupSubscribe_test.go contains the tests related to the Group
// Subscribe requests on pubnub Api
package tests

import (
	"encoding/json"
	"fmt"
	"github.com/pubnub/go/messaging"
	"github.com/stretchr/testify/assert"
	"strings"
	"time"
	// "os"
	"testing"
)

// TestGroupSubscribeStart prints a message on the screen to mark the beginning of
// subscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestGroupSubscribeStart(t *testing.T) {
	PrintTestMessage("==========Group Subscribe tests start==========")
}

func TestGroupSubscriptionConnectedAndUnsubscribedSingle(t *testing.T) {
	//messaging.SetLogOutput(os.Stderr)
	assert := assert.New(t)
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	r := GenRandom()
	group := fmt.Sprintf("testChannelGroup_sub_%d", r.Intn(20))

	createChannelGroups(pubnubInstance, []string{group})
	defer removeChannelGroups(pubnubInstance, []string{group})

	time.Sleep(1 * time.Second)

	subscribeSuccessChannel := make(chan []byte)
	subscribeErrorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.ChannelGroupSubscribe(group,
		subscribeSuccessChannel, subscribeErrorChannel)
	select {
	case msg := <-subscribeSuccessChannel:
		val := string(msg)
		assert.Equal(val, fmt.Sprintf(
			"[1, \"Subscription to channel group '%s' connected\", \"%s\"]",
			group, group))
	case err := <-subscribeErrorChannel:
		assert.Fail(string(err))
	}

	go pubnubInstance.ChannelGroupUnsubscribe(group, successChannel, errorChannel)
	select {
	case msg := <-successChannel:
		val := string(msg)
		assert.Equal(val, fmt.Sprintf(
			"[1, \"Subscription to channel group '%s' unsubscribed\", \"%s\"]",
			group, group))
	case err := <-errorChannel:
		assert.Fail(string(err))
	}

	pubnubInstance.CloseExistingConnection()
}

func TestGroupSubscriptionConnectedAndUnsubscribedMultiple(t *testing.T) {
	//messaging.SetLogOutput(os.Stderr)
	assert := assert.New(t)
	pubnub := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	groupsString, _ := GenerateTwoRandomChannelStrings(3)
	groups := strings.Split(groupsString, ",")

	createChannelGroups(pubnub, groups)
	defer removeChannelGroups(pubnub, groups)

	time.Sleep(1 * time.Second)

	subscribeSuccessChannel := make(chan []byte)
	subscribeErrorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	await := make(chan bool)

	go pubnub.ChannelGroupSubscribe(groupsString,
		subscribeSuccessChannel, subscribeErrorChannel)

	go func() {
		var messages []string

		for {
			select {
			case message := <-subscribeSuccessChannel:
				var msg []interface{}

				err := json.Unmarshal(message, &msg)
				if err != nil {
					assert.Fail(err.Error())
				}

				assert.Contains(msg[1].(string), "Subscription to channel group")
				assert.Contains(msg[1].(string), "connected")
				assert.Len(msg, 3)

				messages = append(messages, string(msg[2].(string)))
			case err := <-subscribeErrorChannel:
				assert.Fail("Subscribe error", err)
			case <-timeouts(4):
				assert.Fail("For looop timed out")
				break
			}

			if len(messages) == 3 {
				break
			}
		}
		assert.True(AssertStringSliceElementsEqual(groups, messages))
		await <- true
	}()

	<-await

	go pubnub.ChannelGroupUnsubscribe(groupsString, successChannel, errorChannel)
	go func() {
		var messages []string

		for {
			select {
			case message := <-successChannel:
				var msg []interface{}

				err := json.Unmarshal(message, &msg)
				if err != nil {
					assert.Fail(err.Error())
				}

				assert.Contains(msg[1].(string), "Subscription to channel group")
				assert.Contains(msg[1].(string), "unsubscribed")
				assert.Len(msg, 3)

				messages = append(messages, string(msg[2].(string)))
			case err := <-errorChannel:
				assert.Fail("Subscribe error", err)
			case <-timeouts(4):
				assert.Fail("For looop timed out")
				break
			}

			if len(messages) == 3 {
				break
			}
		}

		assert.True(AssertStringSliceElementsEqual(groups, messages))
		await <- true
	}()

	<-await

	pubnub.CloseExistingConnection()
}

func TestGroupSubscriptionReceiveSingleMessage(t *testing.T) {
	//messaging.SetLogOutput(os.Stderr)
	assert := assert.New(t)
	pubnub := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	r := GenRandom()
	group := fmt.Sprintf("testChannelGroup_sub_%d", r.Intn(20))
	channel := fmt.Sprintf("testChannel_sub_%d", r.Intn(20))

	populateChannelGroup(pubnub, group, channel)
	defer removeChannelGroups(pubnub, []string{group})

	time.Sleep(1 * time.Second)

	subscribeSuccessChannel := make(chan []byte)
	subscribeErrorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	msgReceived := make(chan bool)

	go pubnub.ChannelGroupSubscribe(group,
		subscribeSuccessChannel, subscribeErrorChannel)
	ExpectConnectedEvent(t, "", group, subscribeSuccessChannel,
		subscribeErrorChannel)

	go func() {
		select {
		case message := <-subscribeSuccessChannel:
			var msg []interface{}

			err := json.Unmarshal(message, &msg)
			if err != nil {
				assert.Fail(err.Error())
			}

			assert.Len(msg, 4)
			assert.Equal(msg[2], channel)
			assert.Equal(msg[3], group)
			msgReceived <- true
		case err := <-subscribeErrorChannel:
			assert.Fail(string(err))
		case <-timeouts(3):
			assert.Fail("Subscription timeout")
		}
	}()

	go pubnub.Publish(channel, "hey", successChannel, errorChannel)

	<-msgReceived

	go pubnub.ChannelGroupUnsubscribe(group, unsubscribeSuccessChannel,
		unsubscribeErrorChannel)
	ExpectUnsubscribedEvent(t, "", group, unsubscribeSuccessChannel,
		unsubscribeErrorChannel)

	pubnub.CloseExistingConnection()
}

func TestGroupSubscriptionPresence(t *testing.T) {
	//messaging.SetLogOutput(os.Stderr)
	assert := assert.New(t)
	pubnub := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	r := GenRandom()
	group := fmt.Sprintf("testChannelGroup_sub_%d", r.Intn(20))
	groupPresence := fmt.Sprintf("%s%s", group, presenceSuffix)

	createChannelGroups(pubnub, []string{group})
	defer removeChannelGroups(pubnub, []string{group})

	time.Sleep(1 * time.Second)

	presenceSuccessChannel := make(chan []byte)
	presenceErrorChannel := make(chan []byte)
	subscribeSuccessChannel := make(chan []byte)
	subscribeErrorChannel := make(chan []byte)
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go func() {
		select {
		case <-subscribeSuccessChannel:
		case <-subscribeErrorChannel:
		case <-successChannel:
		case <-errorChannel:
		}
	}()

	go pubnub.ChannelGroupSubscribe(groupPresence,
		presenceSuccessChannel, presenceErrorChannel)
	ExpectConnectedEvent(t, "", group, presenceSuccessChannel,
		presenceErrorChannel)

	go pubnub.ChannelGroupSubscribe(group,
		subscribeSuccessChannel, subscribeErrorChannel)
	select {
	case message := <-presenceSuccessChannel:
		var msg []interface{}

		// msgString := string(message)

		err := json.Unmarshal(message, &msg)
		if err != nil {
			assert.Fail(err.Error())
		}

		assert.Equal("adsf", msg[2].(string))
		assert.Equal(group, msg[3].(string))
		// Message instead join: [{"occupancy":22,"timestamp":1.454068276e+09}]
		// assert.Contains(msgString, "join")
	case err := <-presenceErrorChannel:
		assert.Fail(string(err))
	}

	go pubnub.ChannelGroupUnsubscribe(group, successChannel, errorChannel)
	select {
	case message := <-presenceSuccessChannel:
		var msg []interface{}

		// msgString := string(message)

		err := json.Unmarshal(message, &msg)
		if err != nil {
			assert.Fail(err.Error())
		}

		assert.Len(msg, 4)
		assert.Equal("adsf", msg[2].(string))
		assert.Equal(group, msg[3].(string))
		// Message instead join: [{"occupancy":22,"timestamp":1.454068276e+09}]
		// assert.Contains(msgString, "join")
	case err := <-presenceErrorChannel:
		assert.Fail(string(err))
	}

	pubnub.CloseExistingConnection()
}

func TestGroupSubscriptionAlreadySubscribed(t *testing.T) {
	//messaging.SetLogOutput(os.Stderr)
	assert := assert.New(t)
	pubnub := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	r := GenRandom()
	group := fmt.Sprintf("testChannelGroup_sub_%d", r.Intn(20))

	createChannelGroups(pubnub, []string{group})
	defer removeChannelGroups(pubnub, []string{group})

	time.Sleep(1 * time.Second)

	subscribeSuccessChannel := make(chan []byte)
	subscribeErrorChannel := make(chan []byte)
	subscribeSuccessChannel2 := make(chan []byte)
	subscribeErrorChannel2 := make(chan []byte)
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnub.ChannelGroupSubscribe(group,
		subscribeSuccessChannel, subscribeErrorChannel)
	ExpectConnectedEvent(t, "", group, subscribeSuccessChannel, subscribeErrorChannel)

	go pubnub.ChannelGroupSubscribe(group,
		subscribeSuccessChannel2, subscribeErrorChannel2)
	select {
	case <-subscribeSuccessChannel2:
		assert.Fail("Received success message while expecting error")
	case err := <-subscribeErrorChannel2:
		assert.Contains(string(err), "Subscription to channel group")
		assert.Contains(string(err), "already subscribed")
	}

	go pubnub.ChannelGroupUnsubscribe(group, successChannel, errorChannel)
	ExpectUnsubscribedEvent(t, "", group, successChannel, errorChannel)

	pubnub.CloseExistingConnection()
}

func TestGroupSubscriptionNotSubscribed(t *testing.T) {
	assert := assert.New(t)
	pubnub := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	r := GenRandom()
	group := fmt.Sprintf("testChannelGroup_sub_%d", r.Intn(20))

	createChannelGroups(pubnub, []string{group})
	defer removeChannelGroups(pubnub, []string{group})

	time.Sleep(1 * time.Second)

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnub.ChannelGroupUnsubscribe(group, successChannel, errorChannel)
	select {
	case response := <-successChannel:
		assert.Fail("Received success message while expecting error", string(response))
	case err := <-errorChannel:
		assert.Contains(string(err), "Subscription to channel group")
		assert.Contains(string(err), "not subscribed")
	}

	pubnub.CloseExistingConnection()
}

func TestGroupSubscriptionToNotExistingChannelGroup(t *testing.T) {
	assert := assert.New(t)
	pubnub := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	r := GenRandom()
	group := fmt.Sprintf("testChannelGroup_sub_%d", r.Intn(20))

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	removeChannelGroups(pubnub, []string{group})

	time.Sleep(1 * time.Second)

	go pubnub.ChannelGroupSubscribe(group, successChannel, errorChannel)
	select {
	case response := <-successChannel:
		assert.Fail("Received success message while expecting error", string(response))
	case err := <-errorChannel:
		assert.Contains(string(err), "Channel group or groups result in empty subscription set")
		assert.Contains(string(err), group)
	}

	pubnub.CloseExistingConnection()
}

func createChannelGroups(pubnub *messaging.Pubnub, groups []string) {
	successChannel := make(chan []byte, 1)
	errorChannel := make(chan []byte, 1)

	for _, group := range groups {
		// fmt.Println("Creating group", group)

		pubnub.ChannelGroupAddChannel(group, "adsf", successChannel, errorChannel)

		select {
		case <-successChannel:
			// fmt.Println("Group created")
		case <-errorChannel:
			fmt.Println("Channel group creation error")
		case <-timeout():
			fmt.Println("Channel group creation timeout")
		}
	}
}

func populateChannelGroup(pubnub *messaging.Pubnub, group, channels string) {

	successChannel := make(chan []byte, 1)
	errorChannel := make(chan []byte, 1)

	pubnub.ChannelGroupAddChannel(group, channels, successChannel, errorChannel)

	select {
	case <-successChannel:
		// fmt.Println("Group created")
	case <-errorChannel:
		fmt.Println("Channel group creation error")
	case <-timeout():
		fmt.Println("Channel group creation timeout")
	}
}

func removeChannelGroups(pubnub *messaging.Pubnub, groups []string) {
	successChannel := make(chan []byte, 1)
	errorChannel := make(chan []byte, 1)

	for _, group := range groups {
		// fmt.Println("Removing group", group)

		pubnub.ChannelGroupRemoveGroup(group, successChannel, errorChannel)

		select {
		case <-successChannel:
			// fmt.Println("Group removed")
		case <-errorChannel:
			fmt.Println("Channel group removal error")
		case <-timeout():
			fmt.Println("Channel group removal timeout")
		}
	}
}

// TestGroupSubscribeEnd prints a message on the screen to mark the end of
// subscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestGroupSubscribeEnd(t *testing.T) {
	PrintTestMessage("==========Group Subscribe tests end==========")
}
