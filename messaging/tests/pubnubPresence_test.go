// Package tests has the unit tests of package messaging.
// pubnubPresence_test.go contains the tests related to the presence requests on pubnub Api
package tests

import (
	"encoding/json"
	"fmt"
	"github.com/pubnub/go/messaging"
	"github.com/stretchr/testify/assert"
	// "os"
	"testing"
	"time"
)

// TestPresenceStart prints a message on the screen to mark the beginning of
// presence tests.
// PrintTestMessage is defined in the common.go file.
func TestPresenceStart(t *testing.T) {
	PrintTestMessage("==========Presence tests start==========")
	// messaging.SetLogOutput(os.Stderr)
	// messaging.LoggingEnabled(true)
}

func TestPresence(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	channel := fmt.Sprintf("presence_hb")

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()
	successChannelP, errorChannelP, eventsChannelP := messaging.CreateSubscriptionChannels()

	go pubnubInstance.Presence(channel, successChannelP, errorChannelP, eventsChannelP)
	ExpectConnectedEvent(t, channel+presenceSuffix, "", eventsChannelP)

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)
	go ExpectConnectedEvent(t, fmt.Sprintf("%s%s", channel, presenceSuffix), "", eventsChannel)

	select {
	case msg := <-successChannelP:
		var event messaging.PresenceEvent
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			assert.Fail(t, err.Error())
		}

		assert.Equal(t, event.Action, "join")
		assert.Equal(t, event.Uuid, pubnubInstance.GetUUID())
	case <-errorChannelP:
		assert.Fail(t, "Received Error first instead of presence event")
	case <-timeout():
		assert.Fail(t, "Timeout occured")
	}

	pubnubInstance.CloseExistingConnection()
}

func TestPresenceHeartbeat(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")
	pubnubInstance.SetPresenceHeartbeat(6)
	pubnubInstance.SetPresenceHeartbeatInterval(2)

	channel := fmt.Sprintf("presence_hb")

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()
	successChannelP, errorChannelP, eventsChannelP := messaging.CreateSubscriptionChannels()

	go pubnubInstance.Presence(channel, successChannelP, errorChannelP, eventsChannelP)
	go ExpectConnectedEvent(t, fmt.Sprintf("%s%s", channel, presenceSuffix), "", eventsChannelP)

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)
	go ExpectConnectedEvent(t, channel, "", eventsChannel)

	select {
	case msg := <-successChannelP:
		var event messaging.PresenceEvent
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			assert.Fail(t, err.Error())
		}

		if event.Action == "leave" && event.Uuid == pubnubInstance.GetUUID() {
			assert.Equal(t, event.Uuid, pubnubInstance.GetUUID())
		}
	case <-errorChannelP:
		assert.Fail(t, "Received Error first instead of presence event")
	case <-timeout():
		assert.Fail(t, "Timeout occured")
	}

	<-timeouts(7)

	pubnubInstance.CloseExistingConnection()
}

type HereNowOccupants struct {
	Uuids     []map[string]string
	Occupancy int
}

// TestHereNowStart prints a message on the screen to mark the beginning of
// presence tests.
// PrintTestMessage is defined in the common.go file.
func TestHereNowStart(t *testing.T) {
	PrintTestMessage("==========Here Now tests start==========")
}

// TestCustomUuid subscribes to a pubnub channel using a custom uuid and then
// makes a call to the herenow method of the pubnub api. The custom id should
// be present in the response else the test fails.
func TestCustomUuid(t *testing.T) {
	cipherKey := ""
	customUuid := "customuuid"
	HereNow(t, cipherKey, customUuid)
}

// TestHereNow subscribes to a pubnub channel and then
// makes a call to the herenow method of the pubnub api. The occupancy should
// be greater than one.
func TestHereNow(t *testing.T) {
	cipherKey := ""
	customUuid := "customuuid"
	HereNow(t, cipherKey, customUuid)
}

// TestHereNowWithCipher subscribes to a pubnub channel and then
// makes a call to the herenow method of the pubnub api. The occupancy should
// be greater than one.
func TestHereNowWithCipher(t *testing.T) {
	cipherKey := ""
	customUuid := "mycustomuuid"
	HereNow(t, cipherKey, customUuid)
}

// HereNow is a common method used by the tests TestHereNow, HereNowWithCipher, CustomUuid
// It subscribes to a pubnub channel and then
// makes a call to the herenow method of the pubnub api.
func HereNow(t *testing.T, cipherKey string, customUuid string) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, cipherKey, false, customUuid)

	//defer time.Sleep(1 * time.Second)
	defer pubnubInstance.CloseExistingConnection()

	var occupants HereNowOccupants

	r := GenRandom()
	channel := fmt.Sprintf("testChannel_hn_%d", r.Intn(100))

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()

	hereNowSuccessChannel := make(chan []byte)
	hereNowErrorChannel := make(chan []byte)

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)
	go LogErrors(errorChannel)

	ExpectConnectedEvent(t, channel, "", eventsChannel)

	time.Sleep(5 * time.Second)

	go pubnubInstance.HereNow(channel, true, true, hereNowSuccessChannel, hereNowErrorChannel)

	select {
	case val := <-hereNowSuccessChannel:
		found := false
		err := json.Unmarshal(val, &occupants)
		if err != nil {
			assert.Fail(t, err.Error())
		}

		for _, occupant := range occupants.Uuids {
			if occupant["uuid"] == customUuid {
				found = true
			}
		}

		assert.True(t, found, "UUID is not found")

	case err := <-hereNowErrorChannel:
		assert.Fail(t, "Received Error instead of HereNow response", err)
	case <-timeout():
		assert.Fail(t, "Timeout occured")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, channel, "", eventsChannel)
}

func TestGlobalHereNow(t *testing.T) {
	cipherKey := ""
	customUuid := "customuuid"

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, cipherKey, false, customUuid)

	//defer time.Sleep(1 * time.Second)
	defer pubnubInstance.CloseExistingConnection()

	r := GenRandom()
	channel := fmt.Sprintf("testChannel_hn_%d", r.Intn(100))

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()

	globalHereNowSuccessChannel := make(chan []byte)
	globalHereNowErrorChannel := make(chan []byte)

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)
	go LogErrors(errorChannel)

	ExpectConnectedEvent(t, channel, "", eventsChannel)

	time.Sleep(5 * time.Second)

	go pubnubInstance.GlobalHereNow(true, false, globalHereNowSuccessChannel,
		globalHereNowErrorChannel)

	select {
	case val := <-globalHereNowSuccessChannel:
		assert.Contains(t, string(val), customUuid)
		assert.Contains(t, string(val), channel)

	case err := <-globalHereNowErrorChannel:
		assert.Fail(t, "Received Error instead of HereNow response", err)
	case <-timeout():
		assert.Fail(t, "Timeout occured")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, channel, "", eventsChannel)
}

func TestWereNow(t *testing.T) {
	cipherKey := ""
	customUuid := "mycustomuuid"

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, cipherKey, false, customUuid)

	defer pubnubInstance.CloseExistingConnection()

	r := GenRandom()
	channel := fmt.Sprintf("testChannel_wn_%d", r.Intn(100))

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()

	whereNowSuccessChannel := make(chan []byte)
	whereNowErrorChannel := make(chan []byte)

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)
	go LogErrors(errorChannel)

	ExpectConnectedEvent(t, channel, "", eventsChannel)

	time.Sleep(5 * time.Second)

	go pubnubInstance.WhereNow(customUuid, whereNowSuccessChannel, whereNowErrorChannel)

	select {
	case val := <-whereNowSuccessChannel:
		assert.Contains(t, string(val), channel)
	case err := <-whereNowErrorChannel:
		assert.Fail(t, "Received Error instead of HereNow response", err)
	case <-timeout():
		assert.Fail(t, "Timeout occured")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, channel, "", eventsChannel)
}

// TestSetGetUserState subscribes to a pubnub channel and then
// makes a call to the herenow method of the pubnub api. The occupancy should
// be greater than one.
func TestSetGetUserState(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	r := GenRandom()
	channel := fmt.Sprintf("testChannel_us_%d", r.Intn(100))
	key := "testkey"
	val := "testval"

	assert := assert.New(t)
	setStateSuccessChannel := make(chan []byte)
	setStateErrorChannel := make(chan []byte)
	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)
	ExpectConnectedEvent(t, channel, "", eventsChannel)

	go pubnubInstance.SetUserStateKeyVal(channel, key, val, setStateSuccessChannel, setStateErrorChannel)
	select {
	case msg := <-setStateSuccessChannel:
		var state struct {
			Payload map[string]interface{}
		}

		err := json.Unmarshal(msg, &state)
		if err != nil {
			assert.Fail(err.Error())
		}
		assert.Equal(val, state.Payload[key])
	case err := <-setStateErrorChannel:
		assert.Fail(string(err))
	}

	time.Sleep(2 * time.Second)

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, channel, "", eventsChannel)

	pubnubInstance.CloseExistingConnection()
}

func TestSetUserStateJSON(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	r := GenRandom()
	channel := fmt.Sprintf("testChannel_us_%d", r.Intn(100))
	key1 := "testkey"
	val1 := "testval"
	key2 := "testkey2"
	val2 := "testval2"

	assert := assert.New(t)

	jsonString := fmt.Sprintf("{\"%s\": \"%s\",\"%s\": \"%s\"}", key1, val1, key2, val2)
	setStateSuccessChannel := make(chan []byte)
	setStateErrorChannel := make(chan []byte)
	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()

	go pubnubInstance.Subscribe(channel, successChannel, errorChannel, eventsChannel)
	ExpectConnectedEvent(t, channel, "", eventsChannel)

	go pubnubInstance.SetUserStateJSON(channel, jsonString, setStateSuccessChannel, setStateErrorChannel)
	select {
	case msg := <-setStateSuccessChannel:
		var state struct {
			Payload map[string]interface{}
		}

		err := json.Unmarshal(msg, &state)
		if err != nil {
			assert.Fail(err.Error())
		}
		assert.Equal(val1, state.Payload[key1])
		assert.Equal(val2, state.Payload[key2])
	case err := <-setStateErrorChannel:
		assert.Fail(string(err))
	}

	time.Sleep(2 * time.Second)

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectDisconnectedEvent(t, channel, "", eventsChannel)

	pubnubInstance.CloseExistingConnection()
}

// TestPresenceEnd prints a message on the screen to mark the end of
// presence tests.
// PrintTestMessage is defined in the common.go file.
func TestPresenceEnd(t *testing.T) {
	PrintTestMessage("==========Presence tests end==========")
}
