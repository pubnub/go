// Package tests has the unit tests of package messaging.
// pubnubEncryption_test.go contains the tests related to the Encryption/Decryption of messages
package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pubnub/go/messaging"
	"github.com/stretchr/testify/assert"
)

func TestPamStart(t *testing.T) {
	PrintTestMessage("==========PAM tests start==========")
}

func TestSecretKeyRequired(t *testing.T) {
	assert := assert.New(t)

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, "", "", false, "", CreateLoggerForTests())
	channel := "testSecretKeyRequired"

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.GrantSubscribe(channel, true, true, 12, "",
		successChannel, errorChannel)
	select {
	case <-successChannel:
		assert.Fail("Response on success channel while expecting error")
	case err := <-errorChannel:
		assert.Contains(string(err), "Secret key is required")
	case <-timeout():
		assert.Fail("GrantSubscribe Timeout")
	}
}

func TestGrantAndRevokeSubKeyLevelSubscribe(t *testing.T) {
	assert := assert.New(t)

	stop, sleep := NewVCRNonSubscribe(
		"fixtures/pam/grantAndRevokeSubKeyLevelSubscribe",
		[]string{"uuid", "signature", "timestamp"})
	defer stop()

	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "", CreateLoggerForTests())
	ttl := 4
	message := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"subscribe_key":"%s","ttl":%d,"r":1,"w":1,"m":0,"d":0,"level":"subkey"}}`, PamSubKey, ttl)
	message2 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"level":"subkey","subscribe_key":"%s","ttl":%d,"r":0,"w":0,"m":0,"d":0}}`, PamSubKey, 1)

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.GrantSubscribe("", true, true, ttl, "", successChannel, errorChannel)
	select {
	case resp := <-successChannel:
		response := string(resp)
		assert.JSONEq(message, response)
	case err := <-errorChannel:
		assert.Fail(string(err))
	case <-timeout():
		assert.Fail("GrantSubscribe Timeout")
	}

	successChannel2 := make(chan []byte)
	errorChannel2 := make(chan []byte)

	sleep(5)

	go pubnubInstance.GrantSubscribe("", false, false, -1, "", successChannel2, errorChannel2)
	select {
	case resp := <-successChannel2:
		response := string(resp)
		assert.JSONEq(message2, response)
	case err := <-errorChannel2:
		assert.Fail(string(err))
	case <-timeout():
		assert.Fail("GrantSubscribe Timeout")
	}
}

func TestGrantAndRevokeChannelLevelSubscribe(t *testing.T) {
	assert := assert.New(t)

	stop, sleep := NewVCRNonSubscribe(
		"fixtures/pam/grantAndRevokeChannelLevelSubscribe",
		[]string{"uuid", "signature", "timestamp"})
	defer stop()

	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "", CreateLoggerForTests())
	channel := "testChannelGrantAndRevokeChannelLevelSubscribe"
	ttl := 8

	message := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s":{"r":1,"w":1,"m":0,"d":0}},"level":"channel","subscribe_key":"%s","ttl":%d}}`, channel, PamSubKey, ttl)
	message2 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"ttl":%d,"channels":{"%s":{"w":0,"m":0,"d":0,"r":0}},"level":"channel","subscribe_key":"%s"}}`, 1, channel, PamSubKey)

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.GrantSubscribe(channel, true, true, ttl, "", successChannel, errorChannel)
	select {
	case resp := <-successChannel:
		response := string(resp)
		assert.JSONEq(message, response)
	case err := <-errorChannel:
		assert.Fail(string(err))
	case <-timeout():
		assert.Fail("GrantSubscribe Timeout")
	}

	successChannel2 := make(chan []byte)
	errorChannel2 := make(chan []byte)

	sleep(5)

	go pubnubInstance.GrantSubscribe(channel, false, false, -1, "", successChannel2, errorChannel2)
	select {
	case resp := <-successChannel2:
		response := string(resp)
		assert.JSONEq(message2, response)
	case err := <-errorChannel2:
		assert.Fail(string(err))
	case <-timeout():
		assert.Fail("GrantSubscribe Timeout")
	}
}

func TestGrantChannelLevelSubscribeWithAuth(t *testing.T) {
	assert := assert.New(t)

	stop, _ := NewVCRNonSubscribe(
		"fixtures/pam/grantChannelLevelSubscribeWithAuth",
		[]string{"uuid", "signature", "timestamp"})
	defer stop()

	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "", CreateLoggerForTests())
	channel := "testGrantChannelLevelSubscribeWithAuth"
	authKey := "myAuthKey"

	ttl := 1
	expected := fmt.Sprintf(`{
		"auths":{"%s":{"d":0,"m":0,"r":1,"w":1}},
		"channel":"%s",
		"level":"user",
		"subscribe_key":"%s",
		"ttl":%d
	}`, authKey, channel, PamSubKey, ttl)

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	go pubnubInstance.GrantSubscribe(channel, true, true, ttl, authKey, successChannel, errorChannel)
	select {
	case resp := <-successChannel:
		var response PamResponse
		err := json.Unmarshal(resp, &response)
		if err != nil {
			assert.Fail(err.Error())
		}

		payload, err := json.Marshal(response.Payload)
		if err != nil {
			assert.Fail(err.Error())
		}

		assert.JSONEq(expected, string(payload))
	case err := <-errorChannel:
		assert.Fail(string(err))
	case <-timeout():
		assert.Fail("GrantSubscribe Timeout")
	}
}

func TestPamEnd(t *testing.T) {
	PrintTestMessage("==========PAM tests End==========")
}
