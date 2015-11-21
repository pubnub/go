// Package tests has the unit tests of package messaging.
// pubnubEncryption_test.go contains the tests related to the Encryption/Decryption of messages
package tests

import (
	"encoding/json"
	"fmt"
	"github.com/pubnub/go/messaging"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestPamStart(t *testing.T) {
	PrintTestMessage("==========PAM tests start==========")

	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")

	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.GrantSubscribe("", false, false, -1, "", returnPamChannel, errorChannel)

	<-returnPamChannel

	time.Sleep(time.Duration(5) * time.Second)
}

func ParsePamErrorResponse(channel chan []byte, testName string, message string, responseChannel chan string) {
	for {
		value, ok := <-channel
		if !ok {
			break
		}
		returnVal := string(value)
		fmt.Println("returnValErr:", returnVal)
		fmt.Println("messageErr:", message)
		if returnVal != "[]" {
			if strings.Contains(returnVal, "aborted") || strings.Contains(returnVal, "reset") {
				continue
			}
			if strings.Contains(returnVal, message) {
				responseChannel <- "Test '" + testName + "': passed."
				break
			} else {
				responseChannel <- "Test '" + testName + "': failed."
			}

			break
		}
	}
}

func ParsePamResponseEqual(returnChannel chan []byte, pubnubInstance *messaging.Pubnub, message string, channel string, testName string, responseChannel chan string, t *testing.T) {
	for {
		value, ok := <-returnChannel
		if !ok {
			break
		}

		if string(value) != "[]" {
			response := string(value)
			fmt.Println("returnValErr:", response)
			fmt.Println("messageErr:", message)

			if assert.JSONEq(t, message, response) {
				responseChannel <- "Test '" + testName + "': passed."
				break
			} else {
				responseChannel <- "Test '" + testName + "': failed."
			}
		}
	}
}

func TestSecretKeyRequired(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "")
	channel := "testChannel"

	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)

	go pubnubInstance.GrantSubscribe(channel, true, true, 12, "", returnPamChannel, errorChannel)
	go ParsePamErrorResponse(errorChannel, "SecretKeyRequired", "Secret key is required", responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SecretKeyRequired")

}

func TestGrantAndRevokeSubKeyLevelSubscribe(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	ttl := 4
	message := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"r":1,"m":0,"w":1,"subscribe_key":"%s","ttl":%d,"level":"subkey"}}`, PamSubKey, ttl)
	message2 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"r":0,"m":0,"w":0,"subscribe_key":"%s","ttl":%d,"level":"subkey"}}`, PamSubKey, 1)

	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)

	go pubnubInstance.GrantSubscribe("", true, true, ttl, "", returnPamChannel, errorChannel)
	go ParsePamResponseEqual(returnPamChannel, pubnubInstance, message, "", "GrantSubKeyLevelSubscribe", responseChannel, t)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "GrantSubKeyLevelSubscribe")

	returnPamChannel2 := make(chan []byte)
	errorChannel2 := make(chan []byte)
	responseChannel2 := make(chan string)
	waitChannel2 := make(chan string)

	time.Sleep(time.Duration(5) * time.Second)

	go pubnubInstance.GrantSubscribe("", false, false, -1, "", returnPamChannel2, errorChannel2)
	go ParsePamResponseEqual(returnPamChannel2, pubnubInstance, message2, "", "RevokeSubKeyLevelSubscribe", responseChannel2, t)
	go ParseErrorResponse(errorChannel2, responseChannel2)
	go WaitForCompletion(responseChannel2, waitChannel2)
	ParseWaitResponse(waitChannel2, t, "RevokeSubKeyLevelSubscribe")
}

func TestGrantAndRevokeChannelLevelSubscribe(t *testing.T) {
	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	channel := "testChannelGrantAndRevokeChannelLevelSubscribe"
	ttl := 8
	message := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s":{"r":1,"m":0,"w":1}},"subscribe_key":"%s","ttl":%d,"level":"channel"}}`, channel, PamSubKey, ttl)
	message2 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s":{"r":0,"m":0,"w":0}},"subscribe_key":"%s","ttl":%d,"level":"channel"}}`, channel, PamSubKey, 1)

	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)

	go pubnubInstance.GrantSubscribe(channel, true, true, ttl, "", returnPamChannel, errorChannel)
	go ParsePamResponseEqual(returnPamChannel, pubnubInstance, message, channel, "GrantChannelLevelSubscribe", responseChannel, t)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "GrantChannelLevelSubscribe")

	returnPamChannel2 := make(chan []byte)
	errorChannel2 := make(chan []byte)
	responseChannel2 := make(chan string)
	waitChannel2 := make(chan string)

	time.Sleep(time.Duration(5) * time.Second)

	go pubnubInstance.GrantSubscribe(channel, false, false, -1, "", returnPamChannel2, errorChannel2)
	go ParsePamResponseEqual(returnPamChannel2, pubnubInstance, message2, channel, "RevokeChannelLevelSubscribe", responseChannel2, t)
	go ParseErrorResponse(errorChannel2, responseChannel2)
	go WaitForCompletion(responseChannel2, waitChannel2)
	ParseWaitResponse(waitChannel2, t, "RevokeChannelLevelSubscribe")
}

func TestGrantChannelLevelSubscribeWithAuth(t *testing.T) {
	var response, sendAsReturn []byte
	var pamResponse PamResponse

	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	channel := "testGrantChannelLevelSubscribeWithAuth"
	authKey := "myAuthKey"

	ttl := 1
	expected := fmt.Sprintf(`{
		"auths":{"%s":{"r":1,"m":0,"w":1}},
		"channel":"%s",
		"level":"user",
		"ttl":%d,
		"subscribe_key":"%s"
	}`, authKey, channel, ttl, PamSubKey)

	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)

	go pubnubInstance.GrantSubscribe(channel, true, true, ttl, authKey, returnPamChannel, errorChannel)

	response = <-returnPamChannel
	json.Unmarshal(response, &pamResponse)

	go ParsePamResponseEqual(returnPamChannel, pubnubInstance, expected, channel, "GrantChannelLevelSubscribeWithAuth", responseChannel, t)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)

	sendAsReturn, _ = json.Marshal(pamResponse.Payload)
	returnPamChannel <- []byte(sendAsReturn)

	ParseWaitResponse(waitChannel, t, "GrantChannelLevelSubscribeWithAuth")
}

func TestPamEnd(t *testing.T) {
	PrintTestMessage("==========PAM tests End==========")
}
