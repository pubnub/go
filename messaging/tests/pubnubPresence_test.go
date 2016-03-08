// Package tests has the unit tests of package messaging.
// pubnubPresence_test.go contains the tests related to the presence requests on pubnub Api
package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pubnub/go/messaging"
	"github.com/stretchr/testify/assert"
)

// TestPresenceStart prints a message on the screen to mark the beginning of
// presence tests.
// PrintTestMessage is defined in the common.go file.
func TestPresenceStart(t *testing.T) {
	PrintTestMessage("==========Presence tests start==========")
}

const PresenceServerTimeoutHighter = 5
const PresenceServerTimeoutLower = 3

// TestCustomUuid subscribes to a pubnub channel using a custom uuid and then
// makes a call to the herenow method of the pubnub api. The custom id should
// be present in the response else the test fails.
func TestCustomUuid(t *testing.T) {
	assert := assert.New(t)

	stop, sleep := NewVCRBothWithSleep(
		"fixtures/presence/customUuid", []string{"uuid"}, 1)
	defer stop()

	messaging.SetSubscribeTimeout(10)

	uuid := "customuuid"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, uuid)
	channel := "customUuid"

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)
	successGet := make(chan []byte)
	errorGet := make(chan []byte)

	go pubnubInstance.Subscribe(channel, "", successChannel, false, errorChannel)
	ExpectConnectedEvent(t, channel, "", successChannel, errorChannel)

	sleep(PresenceServerTimeoutHighter)

	go pubnubInstance.HereNow(channel, true, true, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Contains(string(value), uuid)
		var occupants struct {
			Uuids     []map[string]string
			Occupancy int
		}

		err := json.Unmarshal(value, &occupants)
		if err != nil {
			assert.Fail(err.Error())
		}

		found := false
		for _, v := range occupants.Uuids {
			if v["uuid"] == uuid {
				found = true
			}
		}

		assert.True(found)
	case err := <-errorGet:
		assert.Fail("Failed to get online users", string(err))
	case <-messaging.Timeout():
		assert.Fail("HereNow timeout")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectUnsubscribedEvent(t, channel, "", unsubscribeSuccessChannel, unsubscribeErrorChannel)
}

// TestPresence subscribes to the presence notifications on a pubnub channel and
// then subscribes to a pubnub channel. The test waits till we get a response from
// the subscribe call. The method that parses the presence response sets the global
// variable _endPresenceTestAsSuccess to true if the presence contains a join info
// on the channel and _endPresenceTestAsFailure is otherwise.
func Test0Presence(t *testing.T) {
	assert := assert.New(t)

	// TODO: a random request urls prevents to cover test with VCR
	// stop := NewVCRBoth(
	// 	"fixtures/presence/zeroPresence", []string{"uuid"}, 6)
	// defer stop()

	customUuid := "UUID_zeroPresence"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, customUuid)
	channel := "Channel_ZeroPresence"

	successSubscribe := make(chan []byte)
	errorSubscribe := make(chan []byte)
	successPresence := make(chan []byte)
	errorPresence := make(chan []byte)
	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	await := make(chan bool)

	go pubnubInstance.Subscribe(channel, "", successPresence, true, errorPresence)
	ExpectConnectedEvent(t, channel, "", successPresence, errorPresence)

	go func() {
		select {
		case value := <-successPresence:
			data, _, returnedChannel, err := messaging.ParseJSON(value, "")
			if err != nil {
				assert.Fail(err.Error())
			}

			var occupants []struct {
				Action    string
				Uuid      string
				Timestamp float64
				Occupancy int
			}

			err = json.Unmarshal([]byte(data), &occupants)
			if err != nil {
				assert.Fail(err.Error())
			}

			channelSubRepsonseReceived := false
			for i := 0; i < len(occupants); i++ {
				if (occupants[i].Action == "join") && occupants[i].Uuid == customUuid {
					channelSubRepsonseReceived = true
					break
				}
			}

			assert.True(channelSubRepsonseReceived)
			assert.Equal(channel, returnedChannel)

			await <- true
		case err := <-errorPresence:
			await <- false
			assert.Fail("Failed to subscribe to presence", string(err))
		case <-timeouts(15):
			assert.Fail("Presence timeout")
			await <- false
		}
	}()

	go pubnubInstance.Subscribe(channel, "", successSubscribe, false, errorSubscribe)
	ExpectConnectedEvent(t, channel, "", successSubscribe, errorSubscribe)

	<-await

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectUnsubscribedEvent(t, channel, "", unsubscribeSuccessChannel, unsubscribeErrorChannel)

	pubnubInstance.CloseExistingConnection()
}

// TestWhereNow subscribes to a pubnub channel and then
// makes a call to the herenow method of the pubnub api. The occupancy should
// be greater than one.
func TestWhereNow(t *testing.T) {
	assert := assert.New(t)

	stop, sleep := NewVCRBothWithSleep(
		"fixtures/presence/whereNow", []string{"uuid"}, 1)
	defer stop()

	uuid := "UUID_WhereNow"
	channel := "Channel_WhereNow"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, uuid)

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)
	successGet := make(chan []byte)
	errorGet := make(chan []byte)

	go pubnubInstance.Subscribe(channel, "", successChannel, false, errorChannel)
	ExpectConnectedEvent(t, channel, "", successChannel, errorChannel)

	sleep(PresenceServerTimeoutLower)

	go pubnubInstance.WhereNow(uuid, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Contains(string(value), channel)
	case err := <-errorGet:
		assert.Fail("Failed to get channel list", string(err))
	case <-messaging.Timeout():
		assert.Fail("WhereNow timeout")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectUnsubscribedEvent(t, channel, "", unsubscribeSuccessChannel, unsubscribeErrorChannel)

	pubnubInstance.CloseExistingConnection()
}

// TestGlobalHereNow subscribes to a pubnub channel and then
// makes a call to the herenow method of the pubnub api. The occupancy should
// be greater than one.
func TestGlobalHereNow(t *testing.T) {
	assert := assert.New(t)

	stop, sleep := NewVCRBothWithSleep(
		"fixtures/presence/globalWhereNow", []string{"uuid"}, 1)
	defer stop()

	uuid := "UUID_GlobalWhereNow"
	channel := "Channel_GlobalWhereNow"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, uuid)

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)
	successGet := make(chan []byte)
	errorGet := make(chan []byte)

	go pubnubInstance.Subscribe(channel, "", successChannel, false, errorChannel)
	ExpectConnectedEvent(t, channel, "", successChannel, errorChannel)

	sleep(PresenceServerTimeoutLower)

	go pubnubInstance.GlobalHereNow(true, false, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Contains(string(value), channel)
	case err := <-errorGet:
		assert.Fail("Failed to get online users", string(err))
	case <-messaging.Timeout():
		assert.Fail("GlobalHereNow timeout")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectUnsubscribedEvent(t, channel, "", unsubscribeSuccessChannel, unsubscribeErrorChannel)

	pubnubInstance.CloseExistingConnection()
}

// TestSetGetUserState subscribes to a pubnub channel and then
// makes a call to the herenow method of the pubnub api. The occupancy should
// be greater than one.
func TestSetGetUserState(t *testing.T) {
	assert := assert.New(t)

	stop, sleep := NewVCRBothWithSleep(
		"fixtures/presence/setGetUserState", []string{"uuid"}, 1)
	defer stop()

	uuid := "UUID_SetGetUserState"
	channel := "Channel_SetGetUserState"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, uuid)

	key := "testkey"
	val := "testval"

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	successSet := make(chan []byte)
	errorSet := make(chan []byte)
	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channel, "", successChannel, false, errorChannel)
	ExpectConnectedEvent(t, channel, "", successChannel, errorChannel)

	go pubnubInstance.SetUserStateKeyVal(channel, key, val, successSet, errorSet)
	select {
	case value := <-successSet:
		actual := string(value)
		expectedSubstring := fmt.Sprintf("{\"%s\": \"%s\"}", key, val)

		assert.Contains(actual, expectedSubstring)
	case err := <-errorSet:
		assert.Fail("Failed to set state", string(err))
	case <-messaging.Timeout():
		assert.Fail("Set state timeout")
	}

	sleep(PresenceServerTimeoutLower)

	go pubnubInstance.GetUserState(channel, successGet, errorGet)
	select {
	case value := <-successGet:
		actual := string(value)
		expectedSubstring := fmt.Sprintf("{\"%s\": \"%s\"}", key, val)

		assert.Contains(actual, expectedSubstring)
	case err := <-errorSet:
		assert.Fail("Failed to get state", string(err))
	case <-messaging.Timeout():
		assert.Fail("Get state timeout")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectUnsubscribedEvent(t, channel, "", unsubscribeSuccessChannel, unsubscribeErrorChannel)

	pubnubInstance.CloseExistingConnection()
}

func TestSetUserStateHereNow(t *testing.T) {
	assert := assert.New(t)

	stop, sleep := NewVCRBothWithSleep(
		"fixtures/presence/setUserStateHereNow", []string{"uuid"}, 1)
	defer stop()

	channel := "Channel_SetUserStateHereNow"
	uuid := "UUID_SetUserStateHereNow"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, uuid)

	key := "testkey"
	val := "testval"

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	successSet := make(chan []byte)
	errorSet := make(chan []byte)
	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channel, "", successChannel, false, errorChannel)
	ExpectConnectedEvent(t, channel, "", successChannel, errorChannel)

	go pubnubInstance.SetUserStateKeyVal(channel, key, val, successSet, errorSet)
	select {
	case value := <-successSet:
		actual := string(value)
		expectedSubstring := fmt.Sprintf("{\"%s\": \"%s\"}", key, val)

		assert.Contains(actual, expectedSubstring)
	case err := <-errorSet:
		assert.Fail("Failed to set state", string(err))
	case <-messaging.Timeout():
		assert.Fail("Set state timeout")
	}

	sleep(PresenceServerTimeoutLower)

	go pubnubInstance.HereNow(channel, true, true, successGet, errorGet)
	select {
	case value := <-successGet:
		actual := string(value)
		expectedSubstring := fmt.Sprintf("{\"%s\": \"%s\"}", key, val)

		assert.Contains(actual, expectedSubstring)
		assert.Contains(actual, pubnubInstance.GetUUID())
	case err := <-errorSet:
		assert.Fail("Failed to get state", string(err))
	case <-messaging.Timeout():
		assert.Fail("Get state timeout")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectUnsubscribedEvent(t, channel, "", unsubscribeSuccessChannel, unsubscribeErrorChannel)

	pubnubInstance.CloseExistingConnection()
}

func TestSetUserStateGlobalHereNow(t *testing.T) {
	assert := assert.New(t)

	stop, sleep := NewVCRBothWithSleep(
		"fixtures/presence/setUserStateGlobalHereNow", []string{"uuid"}, 1)
	defer stop()

	channel := "Channel_SetUserStateGlobalHereNow"
	uuid := "UUID_SetUserStateGlobalHereNow"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, uuid)

	key := "testkey"
	val := "testval"

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	successSet := make(chan []byte)
	errorSet := make(chan []byte)
	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	go pubnubInstance.Subscribe(channel, "", successChannel, false, errorChannel)
	ExpectConnectedEvent(t, channel, "", successChannel, errorChannel)

	go pubnubInstance.SetUserStateKeyVal(channel, key, val, successSet, errorSet)
	select {
	case value := <-successSet:
		actual := string(value)
		expectedSubstring := fmt.Sprintf("{\"%s\": \"%s\"}", key, val)

		assert.Contains(actual, expectedSubstring)
	case err := <-errorSet:
		assert.Fail("Failed to set state", string(err))
	case <-messaging.Timeout():
		assert.Fail("Set state timeout")
	}

	sleep(PresenceServerTimeoutLower)

	go pubnubInstance.GlobalHereNow(true, true, successGet, errorGet)
	select {
	case value := <-successGet:
		actual := string(value)
		expectedSubstring := fmt.Sprintf("{\"%s\": \"%s\"}", key, val)

		assert.Contains(actual, expectedSubstring)
		assert.Contains(actual, pubnubInstance.GetUUID())
	case err := <-errorSet:
		assert.Fail("Failed to get state", string(err))
	case <-messaging.Timeout():
		assert.Fail("Get state timeout")
	}

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectUnsubscribedEvent(t, channel, "", unsubscribeSuccessChannel, unsubscribeErrorChannel)

	pubnubInstance.CloseExistingConnection()
}

func TestSetUserStateJSON(t *testing.T) {
	assert := assert.New(t)

	stop, sleep := NewVCRNonSubscribeWithSleep(
		"fixtures/presence/setUserStateJSON", []string{"uuid"})
	defer stop()

	channel := "Channel_SetUserStateJSON"
	uuid := "UUID_SetUserStateJSON"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, uuid)

	key1 := "testkey"
	val1 := "testval"
	key2 := "testkey2"
	val2 := "testval2"

	successSet := make(chan []byte)
	errorSet := make(chan []byte)

	jsonString := fmt.Sprintf("{\"%s\": \"%s\",\"%s\": \"%s\"}", key1, val1, key2, val2)

	go pubnubInstance.SetUserStateJSON(channel, jsonString, successSet, errorSet)
	select {
	case value := <-successSet:
		actual := string(value)
		expectedSubstring := fmt.Sprintf("{\"%s\": \"%s\", \"%s\": \"%s\"}", key2, val2, key1, val1)
		assert.Contains(actual, expectedSubstring)
	case err := <-errorSet:
		assert.Fail("Failed to set state", string(err))
	case <-messaging.Timeout():
		assert.Fail("Set state timeout")
	}

	sleep(PresenceServerTimeoutLower)

	go pubnubInstance.SetUserStateKeyVal(channel, key2, "", successSet, errorSet)
	select {
	case value := <-successSet:
		actual := string(value)
		expectedSubstring := fmt.Sprintf("{\"%s\": \"%s\"}", key1, val1)
		assert.Contains(actual, expectedSubstring)
	case err := <-errorSet:
		assert.Fail("Failed to set state", string(err))
	case <-messaging.Timeout():
		assert.Fail("Set state timeout")
	}
}

// TestPresenceEnd prints a message on the screen to mark the end of
// presence tests.
// PrintTestMessage is defined in the common.go file.
func TestPresenceEnd(t *testing.T) {
	PrintTestMessage("==========Presence tests end==========")
}
