// Package tests has the unit tests of package messaging.
// pubnubDetailedHistory_test.go contains the tests related to the History requests on pubnub Api
package tests

import (
	"encoding/json"
	"fmt"
	"github.com/pubnub/go/messaging"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestDetailedHistoryStart prints a message on the screen to mark the beginning of
// detailed history tests.
// PrintTestMessage is defined in the common.go file.
func TestDetailedHistoryStart(t *testing.T) {
	PrintTestMessage("==========DetailedHistory tests start==========")
}

// TestDetailedHistoryFor10Messages publish's 10 unencrypted messages to a pubnub channel, and after that
// calls the history method of the messaging package to fetch last 10 messages. These received
// messages are compared to the messages sent and if all match test is successful.
func TestDetailedHistoryFor10Messages(t *testing.T) {
	testName := "TestDetailedHistoryFor10Messages"
	DetailedHistoryFor10Messages(t, "", testName)
}

// TestDetailedHistoryFor10EncryptedMessages publish's 10 encrypted messages to a pubnub channel, and after that
// calls the history method of the messaging package to fetch last 10 messages. These received
// messages are compared to the messages sent and if all match test is successful.
func TestDetailedHistoryFor10EncryptedMessages(t *testing.T) {
	testName := "TestDetailedHistoryFor10EncryptedMessages"
	DetailedHistoryFor10Messages(t, "enigma", testName)
}

// DetailedHistoryFor10Messages is a common method used by both TestDetailedHistoryFor10EncryptedMessages
// and TestDetailedHistoryFor10Messages to publish's 10 messages to a pubnub channel, and after that
// call the history method of the messaging package to fetch last 10 messages. These received
// messages are compared to the messages sent and if all match test is successful.
func DetailedHistoryFor10Messages(t *testing.T, cipherKey string, testName string) {
	numberOfMessages := 10
	startMessagesFrom := 0
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, cipherKey, false, "")

	message := "Test Message "
	channel := RandomChannel()

	messagesSent := PublishMessages(pubnubInstance, channel, t, startMessagesFrom, numberOfMessages, message)
	if messagesSent {
		returnHistoryChannel := make(chan []byte)
		errorChannel := make(chan []byte)
		responseChannel := make(chan string)
		waitChannel := make(chan string)

		go pubnubInstance.History(channel, numberOfMessages, 0, 0, false, returnHistoryChannel, errorChannel)
		go ParseHistoryResponseForMultipleMessages(returnHistoryChannel, channel, message, testName, startMessagesFrom, numberOfMessages, cipherKey, responseChannel)
		go ParseErrorResponse(errorChannel, responseChannel)
		go WaitForCompletion(responseChannel, waitChannel)
		ParseWaitResponse(waitChannel, t, testName)
	} else {
		t.Error("Test '" + testName + "': failed.")
	}
}

// TestDetailedHistoryParamsFor10MessagesWithSeretKey publish's 10 unencrypted secret keyed messages
// to a pubnub channel, and after that calls the history method of the messaging package to fetch
// last 10 messages with time parameters between which the messages were sent. These received
// messages are compared to the messages sent and if all match test is successful.
func TestDetailedHistoryParamsFor10MessagesWithSeretKey(t *testing.T) {
	testName := "TestDetailedHistoryFor10MessagesWithSeretKey"
	DetailedHistoryParamsFor10Messages(t, "", "secret", testName)
}

// TestDetailedHistoryParamsFor10EncryptedMessagesWithSeretKey publish's 10 encrypted secret keyed messages
// to a pubnub channel, and after that calls the history method of the messaging package to fetch
// last 10 messages with time parameters between which the messages were sent. These received
// messages are compared to the messages sent and if all match test is successful.
func TestDetailedHistoryParamsFor10EncryptedMessagesWithSeretKey(t *testing.T) {
	testName := "TestDetailedHistoryFor10EncryptedMessagesWithSeretKey"
	DetailedHistoryParamsFor10Messages(t, "enigma", "secret", testName)
}

// TestDetailedHistoryParamsFor10Messages publish's 10 unencrypted messages
// to a pubnub channel, and after that calls the history method of the messaging package to fetch
// last 10 messages with time parameters between which the messages were sent. These received
// messages are compared to the messages sent and if all match test is successful.
func TestDetailedHistoryParamsFor10Messages(t *testing.T) {
	testName := "TestDetailedHistoryFor10Messages"
	DetailedHistoryParamsFor10Messages(t, "", "", testName)
}

// TestDetailedHistoryParamsFor10EncryptedMessages publish's 10 encrypted messages
// to a pubnub channel, and after that calls the history method of the messaging package to fetch
// last 10 messages with time parameters between which the messages were sent. These received
// messages are compared to the messages sent and if all match test is successful.
func TestDetailedHistoryParamsFor10EncryptedMessages(t *testing.T) {
	testName := "TestDetailedHistoryParamsFor10EncryptedMessages"
	DetailedHistoryParamsFor10Messages(t, "enigma", "", testName)
}

// DetailedHistoryFor10Messages is a common method used by both TestDetailedHistoryFor10EncryptedMessages
// and TestDetailedHistoryFor10Messages to publish's 10 messages to a pubnub channel, and after that
// call the history method of the messaging package to fetch last 10 messages with time parameters
// between which the messages were sent. These received message is compared to the messages sent and
// if all match test is successful.
func DetailedHistoryParamsFor10Messages(t *testing.T, cipherKey string, secretKey string, testName string) {
	time.Sleep(5 * time.Second)
	assert := assert.New(t)
	numberOfMessages := 5
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, secretKey, cipherKey, false, "")

	message := "Test Message "
	channel := RandomChannel()

	startTime := GetServerTime(pubnubInstance, t, testName)
	startMessagesFrom := 0
	messagesSent := PublishMessages(pubnubInstance, channel, t, startMessagesFrom, numberOfMessages, message)
	midTime := GetServerTime(pubnubInstance, t, testName)
	startMessagesFrom = 5
	messagesSent2 := PublishMessages(pubnubInstance, channel, t, startMessagesFrom, numberOfMessages, message)
	endTime := GetServerTime(pubnubInstance, t, testName)
	startMessagesFrom = 0

	assert.True(messagesSent, "Error while sending a first bunch of messages")
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.History(channel, numberOfMessages, startTime, midTime,
		false, successChannel, errorChannel)
	select {
	case value := <-successChannel:
		data, _, _, err := messaging.ParseJSON(value, cipherKey)
		if err != nil {
			assert.Fail(err.Error())
		}

		var arr []string
		err = json.Unmarshal([]byte(data), &arr)
		if err != nil {
			assert.Fail(err.Error())
		}

		messagesReceived := 0

		assert.Len(arr, numberOfMessages)

		for i := 0; i < len(arr); i++ {
			if arr[i] == message+strconv.Itoa(startMessagesFrom+i) {
				messagesReceived++
			}
		}

		assert.Equal(numberOfMessages, messagesReceived)
	case err := <-errorChannel:
		assert.Fail(string(err))
	}

	startMessagesFrom = 5

	assert.True(messagesSent2, "Error while sending a second bunch of messages")

	go pubnubInstance.History(channel, numberOfMessages, midTime, endTime, false,
		successChannel, errorChannel)
	select {
	case value := <-successChannel:
		data, _, _, err := messaging.ParseJSON(value, cipherKey)
		if err != nil {
			assert.Fail(err.Error())
		}

		var arr []string
		err = json.Unmarshal([]byte(data), &arr)
		if err != nil {
			assert.Fail(err.Error())
		}

		messagesReceived := 0

		assert.Len(arr, numberOfMessages)

		for i := 0; i < len(arr); i++ {
			if arr[i] == message+strconv.Itoa(startMessagesFrom+i) {
				messagesReceived++
			}
		}

		assert.Equal(numberOfMessages, messagesReceived)
	case err := <-errorChannel:
		assert.Fail(string(err))
	}
}

// GetServerTime calls the GetTime method of the messaging, parses the response to get the
// value and return it.
func GetServerTime(pubnubInstance *messaging.Pubnub, t *testing.T, testName string) int64 {
	assert := assert.New(t)
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.GetTime(successChannel, errorChannel)
	select {
	case value := <-successChannel:
		response := string(value)
		timestamp, err := strconv.Atoi(strings.Trim(response, "[]"))
		if err != nil {
			assert.Fail(err.Error())
		}

		return int64(timestamp)
	case err := <-errorChannel:
		assert.Fail(string(err))
		return 0
	case <-timeouts(10):
		assert.Fail("Getting server timestamp timeout")
		return 0
	}
}

// PublishMessages calls the publish method of messaging package numberOfMessages times
// and appends the count with the message to distinguish from the others.
//
// Parameters:
// pubnubInstance: a reference of *messaging.Pubnub,
// channel: the pubnub channel to publish the messages,
// t: a reference to *testing.T,
// startMessagesFrom: the message identifer,
// numberOfMessages: number of messages to send,
// message: message to send.
//
// returns a bool if the publish of all messages is successful.
func PublishMessages(pubnubInstance *messaging.Pubnub, channel string, t *testing.T, startMessagesFrom int, numberOfMessages int, message string) bool {
	assert := assert.New(t)
	messagesReceived := 0
	messageToSend := ""
	tOut := messaging.GetNonSubscribeTimeout()
	messaging.SetNonSubscribeTimeout(30)

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	for i := startMessagesFrom; i < startMessagesFrom+numberOfMessages; i++ {
		messageToSend = message + strconv.Itoa(i)

		go pubnubInstance.Publish(channel, messageToSend, successChannel, errorChannel)
		select {
		case <-successChannel:
			messagesReceived++
		case err := <-errorChannel:
			assert.Fail("Failed to get channel list", string(err))
		case <-messaging.Timeout():
			assert.Fail("WhereNow timeout")
		}
	}

	time.Sleep(3 * time.Second)
	if messagesReceived == numberOfMessages {
		return true
	}

	messaging.SetNonSubscribeTimeout(tOut)

	return false
}

// ParseHistoryResponseForMultipleMessages unmarshalls the response of the history call to the
// pubnub api and compares the received messages to the sent messages. If the response match the
// test is successful.
//
// Parameters:
// returnChannel: channel to read the response from,
// t: a reference to *testing.T,
// channel: the pubnub channel to publish the messages,
// message: message to compare,
// testname: the test name form where this method is called,
// startMessagesFrom: the message identifer,
// numberOfMessages: number of messages to send,
// cipherKey: the cipher key if used. Can be empty.
func ParseHistoryResponseForMultipleMessages(returnChannel chan []byte, channel string, message string, testName string, startMessagesFrom int, numberOfMessages int, cipherKey string, responseChannel chan string) {
	for {
		value, ok := <-returnChannel
		if !ok {
			break
		}
		if string(value) != "[]" {
			data, _, _, err := messaging.ParseJSON(value, cipherKey)
			if err != nil {
				//t.Error("Test '" + testName + "': failed.")
				responseChannel <- "Test '" + testName + "': failed. Message: " + err.Error()
			} else {
				var arr []string
				err2 := json.Unmarshal([]byte(data), &arr)
				if err2 != nil {
					//t.Error("Test '" + testName + "': failed.");
					responseChannel <- "Test '" + testName + "': failed. Message: " + err2.Error()
				} else {
					messagesReceived := 0

					if len(arr) != numberOfMessages {
						responseChannel <- "Test '" + testName + "': failed."
						//t.Error("Test '" + testName + "': failed.");
						break
					}
					for i := 0; i < numberOfMessages; i++ {
						if arr[i] == message+strconv.Itoa(startMessagesFrom+i) {
							//fmt.Println("data:",arr[i])
							messagesReceived++
						}
					}
					if messagesReceived == numberOfMessages {
						fmt.Println("Test '" + testName + "': passed.")
						responseChannel <- "Test '" + testName + "': passed."
					} else {
						responseChannel <- "Test '" + testName + "': failed. Returned message mismatch"
						//t.Error("Test '" + testName + "': failed.");
					}
					break
				}
			}
		}
	}
}

// ParseHistoryResponse parses the history response from the pubnub api on the returnChannel
// and checks if the response contains the message. If true then the test is successful.
func ParseHistoryResponse(returnChannel chan []byte, channel string, message string, testName string, responseChannel chan string) {
	for {
		value, ok := <-returnChannel
		if !ok {
			break
		}
		if string(value) != "[]" {
			response := string(value)
			//fmt.Println("response", response)

			if strings.Contains(response, message) {
				//fmt.Println("Test '" + testName + "': passed.")
				responseChannel <- "Test '" + testName + "': passed."
				break
			} else {
				responseChannel <- "Test '" + testName + "': failed."
			}
		}
	}
}

// TestDetailedHistoryEnd prints a message on the screen to mark the end of
// detailed history tests.
// PrintTestMessage is defined in the common.go file.
func TestDetailedHistoryEnd(t *testing.T) {
	PrintTestMessage("==========DetailedHistory tests end==========")
}
