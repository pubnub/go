package pubnubTests

import (
    "testing"
    "fmt"
    "strings"
    "pubnubMessaging"
    "time"
)

// Start indicator
func TestUnsubscribeStart(t *testing.T){
	PrintTestMessage("==========Unsubscribe tests start==========")
}

func TestUnsubscribeNotSubscribed(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    

    currentTime := time.Now()
    channel := "testChannel" + currentTime.Format("20060102150405")
    
    returnUnsubscribeChannel := make(chan []byte)
    go pubnubInstance.Unsubscribe(channel, returnUnsubscribeChannel)
    ParseUnsubscribeResponse(returnUnsubscribeChannel, t, channel, "not subscribed")    
}

func TestUnsubscribe(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponseAndCallUnsubscribe(pubnubInstance, returnSubscribeChannel, t, channel, "connected")    
}

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

// End indicator
func TestUnsubscribeEnd(t *testing.T){
	PrintTestMessage("==========Unsubscribe tests end==========")
}   
