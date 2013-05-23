// Package pubnubMessaging has the unit tests of package pubnubMessaging.
// pubnubUnsubscribe_test.go contains the tests related to the Unsubscribe requests on pubnub Api
package pubnubTests

import (
    "testing"
    "fmt"
    "strings"
    "github.com/pubnub/go/3.4/pubnubMessaging"
    "time"
)

// TestUnsubscribeStart prints a message on the screen to mark the beginning of 
// unsubscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestUnsubscribeStart(t *testing.T){
    PrintTestMessage("==========Unsubscribe tests start==========")
}

// TestUnsubscribeNotSubscribed will try to unsubscribe a non subscribed pubnub channel. 
// The response should contain 'not subscribed'
func TestUnsubscribeNotSubscribed(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    

    currentTime := time.Now()
    channel := "testChannel" + currentTime.Format("20060102150405")
    
    returnUnsubscribeChannel := make(chan []byte)
    go pubnubInstance.Unsubscribe(channel, returnUnsubscribeChannel)
    ParseUnsubscribeResponse(returnUnsubscribeChannel, t, channel, "not subscribed")    
}

// TestUnsubscribe will subscribe to a pubnub channel and then send an unsubscribe request
// The response should contain 'unsubscribed'
func TestUnsubscribe(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponseAndCallUnsubscribe(pubnubInstance, returnSubscribeChannel, t, channel, "connected")    
}

// ParseSubscribeResponseAndCallUnsubscribe will parse the response on the go channel.
// It will check the subscribe connection status and when connected
// it will initiate the unsubscribe request. 
func ParseSubscribeResponseAndCallUnsubscribe(pubnubInstance *pubnubMessaging.Pubnub, returnChannel chan []byte, t *testing.T, channel string, message string){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := fmt.Sprintf("%s", value)
            
            message = "'" + channel + "' " + message
            messageAbort := "'" + channel + "' aborted" 
            if(strings.Contains(response, message)){
                returnUnsubscribeChannel := make(chan []byte)
                go pubnubInstance.Unsubscribe(channel, returnUnsubscribeChannel)
                ParseUnsubscribeResponse(returnUnsubscribeChannel, t, channel, "unsubscribed")    
                break
            } else if (strings.Contains(response, messageAbort)){
                t.Error("Test unsubscribed: failed.");
                break
            } else {
                t.Error("Test unsubscribed: failed.");
                break
            }
        }
    }
}

// ParseUnsubscribeResponse will parse the unsubscribe response on the go channel. 
// If it contains unsubscribed the test will pass.
func ParseUnsubscribeResponse(returnChannel chan []byte, t *testing.T, channel string, message string){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := fmt.Sprintf("%s", value)
            if(strings.Contains(response, message)){
                fmt.Println("Test '" + message + "': passed.")
                break
            } else {
                t.Error("Test '" + message + "': failed.");
                break
            }
        }
    }
}

// TestUnsubscribeEnd prints a message on the screen to mark the end of 
// unsubscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestUnsubscribeEnd(t *testing.T){
    PrintTestMessage("==========Unsubscribe tests end==========")
}   
