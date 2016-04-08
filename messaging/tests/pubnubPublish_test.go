// Package tests has the unit tests of package messaging.
// pubnubPublish_test.go contains the tests related to the publish requests on pubnub Api
package tests

import (
	"encoding/json"
	"testing"

	"github.com/pubnub/go/messaging"
	"github.com/stretchr/testify/assert"
)

// TestPublishStart prints a message on the screen to mark the beginning of
// publish tests.
// PrintTestMessage is defined in the common.go file.
func TestPublishStart(t *testing.T) {
	PrintTestMessage("==========Publish tests start==========")
}

// TestNullMessage sends out a null message to a pubnub channel. The response should
// be an "Invalid Message".
func TestNullMessage(t *testing.T) {
	assert := assert.New(t)
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	channel := "nullMessage"
	var message interface{}
	message = nil

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.Publish(channel, message, successChannel, errorChannel)
	select {
	case msg := <-successChannel:
		assert.Fail("Response on success channel while expecting an error", string(msg))
	case err := <-errorChannel:
		assert.Contains(string(err), "Invalid Message")
	case <-timeout():
		assert.Fail("Publish timeout")
	}
}

// TestSuccessCodeAndInfo sends out a message to the pubnub channel
func TestSuccessCodeAndInfo(t *testing.T) {
	assert := assert.New(t)

	stop, _ := NewVCRNonSubscribe("fixtures/publish/successCodeAndInfo",
		[]string{"uuid"})
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	channel := "successCodeAndInfo"
	message := "Pubnub API Usage Example"

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.Publish(channel, message, successChannel, errorChannel)
	select {
	case msg := <-successChannel:
		assert.Contains(string(msg), "1,")
		assert.Contains(string(msg), "\"Sent\",")
	case err := <-errorChannel:
		assert.Fail(string(err))
	case <-timeout():
		assert.Fail("Publish timeout")
	}
}

// TestSuccessCodeAndInfoWithEncryption sends out an encrypted
// message to the pubnub channel
func TestSuccessCodeAndInfoWithEncryption(t *testing.T) {
	assert := assert.New(t)

	stop, _ := NewVCRNonSubscribe(
		"fixtures/publish/successCodeAndInfoWithEncryption", []string{"uuid"})
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "enigma", false, "")
	channel := "successCodeAndInfo"
	message := "Pubnub API Usage Example"

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.Publish(channel, message, successChannel, errorChannel)
	select {
	case msg := <-successChannel:
		assert.Contains(string(msg), "1,")
		assert.Contains(string(msg), "\"Sent\",")
	case err := <-errorChannel:
		assert.Fail(string(err))
	case <-timeout():
		assert.Fail("Publish timeout")
	}
}

// TestSuccessCodeAndInfoForComplexMessage sends out a complex message to the pubnub channel
func TestSuccessCodeAndInfoForComplexMessage(t *testing.T) {
	assert := assert.New(t)

	stop, _ := NewVCRNonSubscribe(
		"fixtures/publish/successCodeAndInfoForComplexMessage", []string{"uuid"})
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	channel := "successCodeAndInfoForComplexMessage"

	customStruct := CustomStruct{
		Foo: "hi!",
		Bar: []int{1, 2, 3, 4, 5},
	}

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.Publish(channel, customStruct, successChannel, errorChannel)
	select {
	case msg := <-successChannel:
		assert.Contains(string(msg), "1,")
		assert.Contains(string(msg), "\"Sent\",")
	case err := <-errorChannel:
		assert.Fail(string(err))
	case <-timeout():
		assert.Fail("Publish timeout")
	}
}

// TestSuccessCodeAndInfoForComplexMessage2 sends out a complex message to the pubnub channel
func TestSuccessCodeAndInfoForComplexMessage2(t *testing.T) {
	assert := assert.New(t)

	stop, _ := NewVCRNonSubscribe(
		"fixtures/publish/successCodeAndInfoForComplexMessage2", []string{"uuid"})
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	channel := "successCodeAndInfoForComplexMessage2"

	customComplexMessage := InitComplexMessage()

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.Publish(channel, customComplexMessage,
		successChannel, errorChannel)
	select {
	case msg := <-successChannel:
		assert.Contains(string(msg), "1,")
		assert.Contains(string(msg), "\"Sent\",")
	case err := <-errorChannel:
		assert.Fail(string(err))
	case <-timeout():
		assert.Fail("Publish timeout")
	}
}

// TestSuccessCodeAndInfoForComplexMessage2WithEncryption sends out an
// encypted complex message to the pubnub channel
func TestSuccessCodeAndInfoForComplexMessage2WithEncryption(t *testing.T) {
	assert := assert.New(t)

	stop, _ := NewVCRNonSubscribe(
		"fixtures/publish/successCodeAndInfoForComplexMessage2WithEncryption",
		[]string{"uuid"})
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "enigma", false, "")
	channel := "successCodeAndInfoForComplexMessage2WithEncryption"

	customComplexMessage := InitComplexMessage()

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.Publish(channel, customComplexMessage,
		successChannel, errorChannel)
	select {
	case msg := <-successChannel:
		assert.Contains(string(msg), "1,")
		assert.Contains(string(msg), "\"Sent\",")
	case err := <-errorChannel:
		assert.Fail(string(err))
	case <-timeout():
		assert.Fail("Publish timeout")
	}
}

func TestPublishStringWithSerialization(t *testing.T) {
	assert := assert.New(t)

	stop, _ := NewVCRBoth(
		"fixtures/publish/publishStringWithSerialization",
		[]string{"uuid"})
	defer stop()

	channel := "Channel_PublishStringWithSerialization"
	uuid := "UUID_PublishStringWithSerialization"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, uuid)
	messageToPost := "{\"name\": \"Alex\", \"age\": \"123\"}"

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	subscribeSuccessChannel := make(chan []byte)
	subscribeErrorChannel := make(chan []byte)

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	await := make(chan bool)

	go pubnubInstance.Subscribe(channel, "", subscribeSuccessChannel, false,
		subscribeErrorChannel)
	ExpectConnectedEvent(t, channel, "", subscribeSuccessChannel,
		subscribeErrorChannel)

	go func() {
		select {
		case message := <-subscribeSuccessChannel:
			var response []interface{}
			var msgs []interface{}
			var err error

			err = json.Unmarshal(message, &response)
			if err != nil {
				assert.Fail(err.Error())
			}

			switch t := response[0].(type) {
			case []interface{}:
				var messageToPostMap map[string]interface{}

				msgs = response[0].([]interface{})
				err := json.Unmarshal([]byte(messageToPost), &messageToPostMap)
				if err != nil {
					assert.Fail(err.Error())
				}

				assert.Equal(messageToPost, msgs[0])
			default:
				assert.Fail("Unexpected response type%s: ", t)
			}

			await <- true
		case err := <-subscribeErrorChannel:
			assert.Fail(string(err))
			await <- false
		case <-timeouts(10):
			assert.Fail("Timeout")
			await <- false
		}
	}()

	go pubnubInstance.Publish(channel, messageToPost, successChannel, errorChannel)
	select {
	case <-successChannel:
	case err := <-errorChannel:
		assert.Fail(string(err))
	case <-timeout():
		assert.Fail("Publish timeout")
	}

	<-await

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectUnsubscribedEvent(t, channel, "", unsubscribeSuccessChannel, unsubscribeErrorChannel)
}

func TestPublishStringWithoutSerialization(t *testing.T) {
	assert := assert.New(t)

	stop, _ := NewVCRBoth(
		"fixtures/publish/publishStringWithoutSerialization",
		[]string{"uuid"})
	defer stop()

	channel := "Channel_PublishStringWithoutSerialization"
	uuid := "UUID_PublishStringWithoutSerialization"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, uuid)
	messageToPost := "{\"name\": \"Alex\", \"age\": \"123\"}"

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	subscribeSuccessChannel := make(chan []byte)
	subscribeErrorChannel := make(chan []byte)

	unsubscribeSuccessChannel := make(chan []byte)
	unsubscribeErrorChannel := make(chan []byte)

	await := make(chan bool)

	go pubnubInstance.Subscribe(channel, "", subscribeSuccessChannel, false,
		subscribeErrorChannel)
	ExpectConnectedEvent(t, channel, "", subscribeSuccessChannel,
		subscribeErrorChannel)

	go func() {
		select {
		case message := <-subscribeSuccessChannel:
			var response []interface{}
			var msgs []interface{}
			var err error

			err = json.Unmarshal(message, &response)
			if err != nil {
				assert.Fail(err.Error())
			}

			switch t := response[0].(type) {
			case []interface{}:
				var messageToPostMap map[string]interface{}

				msgs = response[0].([]interface{})
				err := json.Unmarshal([]byte(messageToPost), &messageToPostMap)
				if err != nil {
					assert.Fail(err.Error())
				}

				assert.Equal(messageToPostMap, msgs[0])
			default:
				assert.Fail("Unexpected response type%s: ", t)
			}

			await <- true
		case err := <-subscribeErrorChannel:
			assert.Fail(string(err))
			await <- false
		case <-timeouts(10):
			assert.Fail("Timeout")
			await <- false
		}
	}()

	go pubnubInstance.PublishExtended(channel, messageToPost, false, true,
		successChannel, errorChannel)
	select {
	case <-successChannel:
	case err := <-errorChannel:
		assert.Fail(string(err))
	case <-timeout():
		assert.Fail("Publish timeout")
	}

	<-await

	go pubnubInstance.Unsubscribe(channel, unsubscribeSuccessChannel, unsubscribeErrorChannel)
	ExpectUnsubscribedEvent(t, channel, "", unsubscribeSuccessChannel, unsubscribeErrorChannel)
}

// TestPublishEnd prints a message on the screen to mark the end of
// publish tests.
// PrintTestMessage is defined in the common.go file.
func TestPublishEnd(t *testing.T) {
	PrintTestMessage("==========Publish tests end==========")
}
