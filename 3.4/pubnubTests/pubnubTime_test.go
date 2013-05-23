// Package pubnubMessaging has the unit tests of package pubnubMessaging.
// pubnubTime_test.go contains the tests related to the Time requests on pubnub Api
package pubnubTests

import (
    "testing"
    "github.com/pubnub/go/3.4/pubnubMessaging"
    "fmt"
)

// TestTimeStart prints a message on the screen to mark the beginning of 
// time tests.
// PrintTestMessage is defined in the common.go file.
func TestTimeStart(t *testing.T){
    PrintTestMessage("==========Time tests start==========")
}

// TestServerTime calls the GetTime method of the pubnubMessaging to test the time
func TestServerTime(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    returnTimeChannel := make(chan []byte)
    go pubnubInstance.GetTime(returnTimeChannel)
    ParseTimeResponse(returnTimeChannel, t)    
}

// ParseTimeResponse parses the time response from the pubnub api.
// On error the test fails.
func ParseTimeResponse(returnChannel chan []byte,t *testing.T){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := fmt.Sprintf("%s", value)
            if(response != ""){
                fmt.Println("Server time: passed.")
                break 
            } else {
                t.Error("Server time: failed.");
                break
            }
        }
    }
}

// TestTimeEnd prints a message on the screen to mark the end of 
// time tests.
// PrintTestMessage is defined in the common.go file.
func TestTimeEnd(t *testing.T){
    PrintTestMessage("==========Time tests end==========")
}   