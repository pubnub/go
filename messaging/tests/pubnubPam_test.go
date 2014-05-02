// Package tests has the unit tests of package messaging.
// pubnubEncryption_test.go contains the tests related to the Encryption/Decryption of messages
package tests

import (
	//"encoding/json"
	"fmt"
	"github.com/pubnub/go/messaging"
	"testing"
	"strings"
	"time"
	//"unicode/utf16"
)

func TestPamStart(t *testing.T) {
	PrintTestMessage("==========PAM tests start==========")
}

func TestSecretKeyRequired(t *testing.T){
	pubnubInstance := messaging.NewPubnub("demo-36", "demo-36", "", "", false, "")
	channel := "testChannel"

	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)
	
	go pubnubInstance.GrantSubscribe (channel, true, true, 12, returnPamChannel, errorChannel)
	go ParsePamErrorResponse(errorChannel, "SecretKeyRequired", "Secret key is required", responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SecretKeyRequired")
}

func ParsePamErrorResponse(channel chan []byte, testName string, message string, responseChannel chan string) {
	for {
		value, ok := <-channel
		if !ok {
			break
		}
		returnVal := string(value)
		//fmt.Println("returnVal:",returnVal);
		if returnVal != "[]" {
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

func ParsePamResponse(returnChannel chan []byte, pubnubInstance *messaging.Pubnub, message string, channel string, testName string, responseChannel chan string) {
	for {
		value, ok := <-returnChannel
		if !ok {
			break
		}
		if string(value) != "[]" {
			response := string(value)
			//fmt.Println("response:", response)
			if strings.Contains(response, message) {

				responseChannel <- "Test '" + testName + "': passed."
				break
			} else {
				responseChannel <- "Test '" + testName + "': failed."
			}
		}
	}
}

func TestSubscribeGrantPositive(t *testing.T){
	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	channel := "testChannelSubscribeGrantPositive"
	ttl := 1
	message := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s":{"r":1,"w":1}},"subscribe_key":"%s","ttl":%d,"level":"channel"}}`, channel, PamSubKey, ttl)
	message2 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s":{"r":0,"w":0}},"subscribe_key":"%s","ttl":%d,"level":"channel"}}`, channel, PamSubKey, ttl)
	
	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)
	
	go pubnubInstance.GrantSubscribe (channel, true, true, ttl, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message, channel, "SubscribeGrantPositiveGrant", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SubscribeGrantPositiveGrant")
	
	go pubnubInstance.GrantSubscribe (channel, false, false, -1, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message2, channel, "SubscribeGrantPositiveRevoke", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SubscribeGrantPositiveRevoke")
}

func TestSubscribeGrantNegative(t *testing.T){
	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	channel := "testChannelSubscribeGrantNegative"
	message := fmt.Sprintf(`{"status":403,"service":"Access Manager","error":true,"message":"Forbidden","payload":{"channels":["%s"]}}`, channel)
	
	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)
	
	go pubnubInstance.Subscribe (channel, "", returnPamChannel, false, errorChannel)
	go ParsePamErrorResponse(errorChannel, "SubscribeGrantNegative", message, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SubscribeGrantNegative")
}

func TestPresenceGrantPositive(t *testing.T){
	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	channel := "testChannelPresenceGrantPositive"
	ttl := 1
	message := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s-pnpres":{"r":1,"w":1}},"subscribe_key":"%s","ttl":%d,"level":"channel"}}`, channel, PamSubKey, ttl)
	message2 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s-pnpres":{"r":0,"w":0}},"subscribe_key":"%s","ttl":%d,"level":"channel"}}`, channel, PamSubKey, ttl)
	
	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)
	
	go pubnubInstance.GrantPresence (channel, true, true, ttl, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message, channel, "PresenceGrantPositiveGrant", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "PresenceGrantPositiveGrant")

	go pubnubInstance.GrantPresence (channel, false, false, -1, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message2, channel, "PresenceGrantPositiveRevoke", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "PresenceGrantPositiveRevoke")
}

func TestPresenceGrantNegative(t *testing.T){
	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	channel := "testChannelPresenceGrantNegative"
	message := fmt.Sprintf(`{"status":403,"service":"Access Manager","error":true,"message":"Forbidden","payload":{"channels":["%s-pnpres"]}}`, channel)
	
	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)
	
	go pubnubInstance.Subscribe (channel, "", returnPamChannel, true, errorChannel)
	go ParsePamErrorResponse(errorChannel, "PresenceGrantNegative", message, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "PresenceGrantNegative")
}

func TestSubscribeAudit(t *testing.T){
	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	channel := "testChannelSubscribeAudit"
	time.Sleep(time.Duration(5) * time.Second)
	ttl:=1
	//message1 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{},"subscribe_key":"%s","level":"subkey"}}`, PamSubKey)
	message1 := fmt.Sprintf(`"subscribe_key":"%s","level":"subkey"`, PamSubKey)
	//	{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{},"subscribe_key":"sub-c-a3d5a1c8-ae97-11e3-a952-02ee2ddab7fe","level":"channel"}}
	//message := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{},"subscribe_key":"%s","level":"channel"}}`, PamSubKey)
	message := fmt.Sprintf(`"subscribe_key":"%s","level":"channel"`, PamSubKey)
	message2 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s":{"r":1,"w":1}},"subscribe_key":"%s","ttl":%d,"level":"channel"}}`, channel, PamSubKey, ttl)
	//{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"testChannelSubscribeAudit":{"r":1,"w":1}},"subscribe_key":"sub-c-a3d5a1c8-ae97-11e3-a952-02ee2ddab7fe","ttl":1,"level":"channel"}}	
	message3 := fmt.Sprintf(`"%s":{"r":1,"w":1,"ttl":%d}`, channel, ttl)
	message4 := fmt.Sprintf(`"%s":{"r":1,"w":1,"ttl":%d}`, channel, ttl)
	message5 := fmt.Sprintf(`[1, "Subscription to channel '%s' connected", "%s"]`, channel, channel)
	
	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)
	
	//audit
	go pubnubInstance.AuditSubscribe(channel, returnPamChannel, errorChannel)
	//fmt.Println("message:", message)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message, channel, "SubscribeAuditChannel", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SubscribeAuditChannel")
	
	//audit	
	go pubnubInstance.AuditSubscribe("", returnPamChannel, errorChannel)
	//fmt.Println("message1:", message1)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message1, "", "SubscribeAuditSubKey", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SubscribeAuditSubKey")
		
	//grant
	go pubnubInstance.GrantSubscribe (channel, true, true, ttl, returnPamChannel, errorChannel)
	//fmt.Println("message2:", message2)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message2, channel, "SubscribeAuditGrantPositiveGrant", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SubscribeAuditGrantPositiveGrant")
	
	go pubnubInstance.Subscribe(channel, "", returnPamChannel, false, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message5, channel, "SubscribeAudit", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SubscribeAudit")
	
	time.Sleep(time.Duration(5) * time.Second)
	
	//audit
	go pubnubInstance.AuditSubscribe(channel, returnPamChannel, errorChannel)
	//fmt.Println("message3:", message3)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message3, channel, "SubscribeAuditChannel2", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SubscribeAuditChannel2")
	
	//audit	
	go pubnubInstance.AuditSubscribe("", returnPamChannel, errorChannel)
	//fmt.Println("message4:", message4)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message4, channel, "SubscribeAuditSubKey2", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SubscribeAuditSubKey2")

	go pubnubInstance.Unsubscribe(channel, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message5, channel, "SubscribeAuditUnsub", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	
	go pubnubInstance.GrantSubscribe (channel, false, false, -1, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message5, channel, "SubscribeAuditRevoke", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)

}

func TestPresenceAudit(t *testing.T){
	/*{"status":200,"service":"Access Manager","message":"Success","payload":{"ch
	annels":{"test":{"r":1,"w":1,"ttl":9},"testChannelSubscribeGrantPositive":{"r":0
	,"w":0,"ttl":1},"testChannelPresenceGrantPositive":{"r":0,"w":0,"ttl":1},"test-p
	npres":{"r":1,"w":1,"ttl":12}},"subscribe_key":"sub-c-a3d5a1c8-ae97-11e3-a952-02
	ee2ddab7fe","level":"subkey"}}*/
	
	/*{"status":200,"service":"Access Manager","message":"Success","payload":{"ch
	annels":{"test-pnpres":{"r":1,"w":1,"ttl":12}},"subscribe_key":"sub-c-a3d5a1c8-a
	e97-11e3-a952-02ee2ddab7fe","level":"channel"}}*/

	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	channel := "testChannelPresenceAudit"
	time.Sleep(time.Duration(10) * time.Second)
	ttl:=1
	//message1 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{},"subscribe_key":"%s","level":"subkey"}}`, PamSubKey)
	message1 := fmt.Sprintf(`"subscribe_key":"%s","level":"subkey"`, PamSubKey)
	//	{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{},"subscribe_key":"sub-c-a3d5a1c8-ae97-11e3-a952-02ee2ddab7fe","level":"channel"}}
	//message := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{},"subscribe_key":"%s","level":"channel"}}`, PamSubKey)
	message := fmt.Sprintf(`"subscribe_key":"%s","level":"channel"`, PamSubKey)
	message2 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s-pnpres":{"r":1,"w":1}},"subscribe_key":"%s","ttl":%d,"level":"channel"}}`, channel, PamSubKey, ttl)
	//{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"testChannelSubscribeAudit":{"r":1,"w":1}},"subscribe_key":"sub-c-a3d5a1c8-ae97-11e3-a952-02ee2ddab7fe","ttl":1,"level":"channel"}}	
	message3 := fmt.Sprintf(`"%s-pnpres":{"r":1,"w":1,"ttl":%d}`, channel, ttl)
	message4 := fmt.Sprintf(`"%s-pnpres":{"r":1,"w":1,"ttl":%d}`, channel, ttl)
	//"testChannelPresenceAudit-pnpres":{"r":1,"w":1,"ttl":1}
	//"testChannelPresenceAudit-pnpres":{"r":1,"w":1,"ttl":1}
	message5 := fmt.Sprintf(`"Presence notifications for channel '%s' connected", "%s"`, channel, channel)
	
	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)
	
	//audit
	go pubnubInstance.AuditPresence(channel, returnPamChannel, errorChannel)
	//fmt.Println("message:", message)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message, channel, "PresenceAuditChannel", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "PresenceAuditChannel")
	//audit	
	go pubnubInstance.AuditPresence("", returnPamChannel, errorChannel)
	//fmt.Println("message1:", message1)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message1, "", "PresenceAuditSubKey", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "PresenceAuditSubKey")
		
	//grant
	go pubnubInstance.GrantPresence (channel, true, true, ttl, returnPamChannel, errorChannel)
	//fmt.Println("message2:", message2)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message2, channel, "PresenceAuditGrant", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "PresenceAuditGrant")
	
	go pubnubInstance.Subscribe(channel, "", returnPamChannel, true, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message5, channel, "PresenceAudit", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "PresenceAudit")
	
	time.Sleep(time.Duration(5) * time.Second)
	
	//audit
	go pubnubInstance.AuditPresence(channel, returnPamChannel, errorChannel)
	//fmt.Println("message3:", message3)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message3, channel, "PresenceAuditChannel2", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "PresenceAuditChannel2")
	//audit	
	go pubnubInstance.AuditPresence("", returnPamChannel, errorChannel)
	//fmt.Println("message4:", message4)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message4, channel, "PresenceAuditSubKey2", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "PresenceAuditSubKey2")

	go pubnubInstance.PresenceUnsubscribe(channel, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message5, channel, "PresenceAuditUnsub", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	

	go pubnubInstance.GrantPresence (channel, false, false, -1, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message5, channel, "PresenceAuditRevoke", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	
}

func TestAuthSubscribe(t *testing.T){
	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	pubnubInstance.SetAuthenticationKey("authkey")
	channel := "testChannelSubscribeAuth"
	
	time.Sleep(time.Duration(10) * time.Second)
	ttl := 1
	message := fmt.Sprintf(`{"auths":{"authkey":{"r":1,"w":1}}`)
	message2 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s":{"r":0,"w":0}},"subscribe_key":"%s","ttl":%d,"level":"channel"}}`, channel, PamSubKey, ttl)
	message5 := fmt.Sprintf(`'%s' connected`, channel)
	message6 := fmt.Sprintf(`'%s' unsubscribed`, channel)
	
	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)
	
	//grant	
	go pubnubInstance.GrantSubscribe (channel, true, true, ttl, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message, channel, "SubscribeAuthGrant", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SubscribeAuthGrant")

	//subscribe
	go pubnubInstance.Subscribe(channel, "", returnPamChannel, false, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message5, channel, "SubscribeAuthSubscribe", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	//check connect 
	ParseWaitResponse(waitChannel, t, "SubscribeAuthSubscribe")

	go pubnubInstance.Unsubscribe(channel, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message6, channel, "SubscribeAuthUnsubscribe", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "SubscribeAuthUnsubscribe")
	
	//revoke
	go pubnubInstance.GrantSubscribe (channel, false, false, -1, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message2, channel, "SubscribeAuthRevoke", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
}

func TestAuthPresence(t *testing.T){
	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	pubnubInstance.SetAuthenticationKey("authkey")
	channel := "testChannelPresenceAuth"
	
	time.Sleep(time.Duration(10) * time.Second)
	ttl := 1
	message := fmt.Sprintf(`{"auths":{"authkey":{"r":1,"w":1}}`)
	message2 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s":{"r":0,"w":0}},"subscribe_key":"%s","ttl":%d,"level":"channel"}}`, channel, PamSubKey, ttl)
	message5 := fmt.Sprintf(`'%s' connected`, channel)
	message6 := fmt.Sprintf(`'%s' unsubscribed`, channel)
	
	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)
	
	//grant	
	go pubnubInstance.GrantPresence (channel, true, true, ttl, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message, channel, "PresenceAuthGrant", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "PresenceAuthGrant")

	//subscribe
	go pubnubInstance.Subscribe(channel, "", returnPamChannel, true, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message5, channel, "PresenceAuthSubscribe", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	//check connect 
	ParseWaitResponse(waitChannel, t, "PresenceAuthSubscribe")

	go pubnubInstance.PresenceUnsubscribe(channel, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message6, channel, "PresenceAuthUnsubscribe", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "PresenceAuthUnsubscribe")
	
	//revoke
	go pubnubInstance.GrantSubscribe (channel, false, false, -1, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message2, channel, "PresenceAuthRevoke", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
}

func TestAuthHereNow(t *testing.T){
	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	pubnubInstance.SetAuthenticationKey("authkey")
	channel := "testChannelHereNowAuth"
	
	time.Sleep(time.Duration(10) * time.Second)
	ttl := 1
	message := fmt.Sprintf(`{"auths":{"authkey":{"r":1,"w":1}}`)
	message2 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s":{"r":0,"w":0}},"subscribe_key":"%s","ttl":%d,"level":"channel"}}`, channel, PamSubKey, ttl)
	message5 := fmt.Sprintf(`'%s' connected`, channel)
	message4 := pubnubInstance.GetUUID()
	message6 := fmt.Sprintf(`'%s' unsubscribed`, channel)
	
	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)
	
	//grant	
	go pubnubInstance.GrantPresence (channel, true, true, ttl, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message, channel, "HereNowAuthGrant", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "HereNowAuthGrant")
	
	//grant	
	go pubnubInstance.GrantSubscribe (channel, true, true, ttl, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message, channel, "HereNowAuthSubscribe", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "HereNowAuthSubscribe")
	
	//subscribe
	go pubnubInstance.Subscribe(channel, "", returnPamChannel, false, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message5, channel, "HereNowAuthSubscribe", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	//check connect 
	ParseWaitResponse(waitChannel, t, "HereNowAuthSubscribe")

	time.Sleep(time.Duration(10) * time.Second)
	//herenow
	go pubnubInstance.HereNow(channel, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message4, channel, "HereNowAuthHereNow", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	//check connect 
	ParseWaitResponse(waitChannel, t, "HereNowAuthHereNow")
	
	go pubnubInstance.Unsubscribe(channel, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message6, channel, "HereNowAuthUnsubscribe", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "HereNowAuthUnsubscribe")

	//revoke
	go pubnubInstance.GrantSubscribe (channel, false, false, -1, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message2, channel, "HereNowAuthHereNow", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	
}	

func TestAuthHistory(t *testing.T){
	pubnubInstance := messaging.NewPubnub(PamPubKey, PamSubKey, PamSecKey, "", false, "")
	pubnubInstance.SetAuthenticationKey("authkey")
	channel := "testChannelHistoryAuth"
	
	time.Sleep(time.Duration(10) * time.Second)
	ttl := 1
	message := fmt.Sprintf(`{"auths":{"authkey":{"r":1,"w":1}}`)
	message2 := fmt.Sprintf(`{"status":200,"service":"Access Manager","message":"Success","payload":{"channels":{"%s":{"r":0,"w":0}},"subscribe_key":"%s","ttl":%d,"level":"channel"}}`, channel, PamSubKey, ttl)
	message5 := "Test Message"
	
	returnPamChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	responseChannel := make(chan string)
	waitChannel := make(chan string)
	
	//grant	
	go pubnubInstance.GrantPresence (channel, true, true, ttl, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message, channel, "HistoryAuthGrant", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "HistoryAuthGrant")

	//grant	
	go pubnubInstance.GrantSubscribe (channel, true, true, ttl, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message, channel, "HistoryAuthSubscribe", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "HistoryAuthSubscribe")
	
	//publish
	go pubnubInstance.Publish(channel, message5, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, "Sent", channel, "HistoryAuthPublish", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	ParseWaitResponse(waitChannel, t, "HistoryAuthPublish")

	//history
	go pubnubInstance.History(channel, 1, 0, 0, false, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message5, channel, "HistoryAuthHistory", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	//check connect 
	ParseWaitResponse(waitChannel, t, "HistoryAuthHistory")
	
	//revoke
	go pubnubInstance.GrantSubscribe (channel, false, false, -1, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message2, channel, "HistoryAuthHistory", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	
	//revoke
	go pubnubInstance.GrantSubscribe (channel, false, false, -1, returnPamChannel, errorChannel)
	go ParsePamResponse(returnPamChannel, pubnubInstance, message2, channel, "HistoryAuthHereNow", responseChannel)
	go ParseErrorResponse(errorChannel, responseChannel)
	go WaitForCompletion(responseChannel, waitChannel)
	
}

func TestPamEnd(t *testing.T) {
	PrintTestMessage("==========PAM tests End==========")
}
