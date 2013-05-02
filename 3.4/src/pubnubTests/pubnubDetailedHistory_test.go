package pubnubTests

import (
    "testing"
    "pubnubMessaging"
    "strings"
    "fmt"
)

// Start indicator
func TestDetailedHistoryStart(t *testing.T){
    PrintTestMessage("==========DetailedHistory tests start==========")
}
    
func TestDetailedHistory(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    message := "Test Message"
    returnPublishChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnPublishChannel)
    ParseResponse(returnPublishChannel, t, pubnubInstance, channel, message, "DetailedHistory", 1)
}

func TestEncryptedDetailedHistory(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "enigma", false, "")    
    
    channel := "testChannel"
    message := "Test Message"
    returnPublishChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnPublishChannel)
    ParseResponse(returnPublishChannel, t, pubnubInstance, channel, message, "EncryptedDetailedHistory", 1)
}


func TestDetailedHistoryFor10Messages(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    message := "Test Message"
    returnPublishChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnPublishChannel)
    ParseResponse(returnPublishChannel, t, pubnubInstance, channel, message, "DetailedHistoryFor10Messages", 1)
}

func ParseHistoryResponse(returnChannel chan []byte, t *testing.T, channel string, message string, testName string){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := string(value)
            //fmt.Println("response", response)
            if(strings.Contains(response, message)){
                fmt.Println("Test '" + testName + "': passed.") 
                break
            } else {
                t.Error("Test '" + testName + "': failed.");
            }
        }
    }
}

func ParseResponse(returnChannel chan []byte,t *testing.T, pubnubInstance *pubnubMessaging.Pubnub, channel string, message string, testName string, numberOfMessages int){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            returnHistoryChannel := make(chan []byte)
            go pubnubInstance.History(channel, 1, 0, 0, false, returnHistoryChannel)
            ParseHistoryResponse(returnHistoryChannel, t, channel, message, testName)
            break
        }
    }
}

// End indicator
func TestDetailedHistoryEnd(t *testing.T){
    PrintTestMessage("==========DetailedHistory tests end==========")
}