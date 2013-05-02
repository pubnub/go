package pubnubTests

import (
    "testing"
    "pubnubMessaging"
    "fmt"
)

// Start indicator
func TestTimeStart(t *testing.T){
    PrintTestMessage("==========Time tests start==========")
}

func TestServerTime(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    returnTimeChannel := make(chan []byte)
    go pubnubInstance.GetTime(returnTimeChannel)
    ParseTimeResponse(returnTimeChannel, t)    
}

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

// End indicator
func TestTimeEnd(t *testing.T){
    PrintTestMessage("==========Time tests end==========")
}   