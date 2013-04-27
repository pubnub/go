package pubnubTests

import (
    "testing"
    "pubnubMessaging"
    "strings"
    "fmt"
)

// Start indicator
func TestDetailedHistoryStart(t *testing.T){
    fmt.Println("")
    fmt.Println("==========DetailedHistory tests start==========")
    fmt.Println("")
}
    
func TestDetailedHistory(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    message := "Test Message"
    returnPublishChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnPublishChannel)
    ParseResponse(returnPublishChannel, t, pubnubInstance, channel, message)
}

func ParseHistoryResponse(returnChannel chan []byte, t *testing.T, channel string, message string){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := fmt.Sprintf("%s", value)
            fmt.Println("response", response)
            if(strings.Contains(response, message)){
                fmt.Println("Detailed history: passed.") 
                break
            } else {
                t.Error("Detailed history: failed.");
            }
        }
    }
}

func ParseResponse(returnChannel chan []byte,t *testing.T, pubnubInstance *pubnubMessaging.Pubnub, channel string, message string){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            returnHistoryChannel := make(chan []byte)
            go pubnubInstance.History(channel, 1, 0, 0, false, returnHistoryChannel)
            ParseHistoryResponse(returnHistoryChannel, t, channel, message)
            break
        }
    }
}

// End indicator
func TestDetailedHistoryEnd(t *testing.T){
    fmt.Println("")
    fmt.Println("==========DetailedHistory tests end==========")
    fmt.Println("")
}