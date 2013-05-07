package pubnubTests

import (
    "testing"
    "pubnubMessaging"
    "strings"
    "fmt"
    //"time"
    "encoding/json"
)

// Start indicator
func TestPresenceStart(t *testing.T){
    PrintTestMessage("==========Presence tests start==========")
}

func TestCustomUuid(t *testing.T) {
    cipherKey := ""
    testName := "CustomUuid"
    customUuid := "customuuid"
    HereNow(t, cipherKey, customUuid, testName)
}

func TestHereNow(t *testing.T) {
    cipherKey := ""
    testName := "HereNow"
    customUuid := "customuuid"
    HereNow(t, cipherKey, customUuid, testName)
}

func TestHereNowWithCipher(t *testing.T) {
    cipherKey := ""
    testName := "HereNowWithCipher"
    customUuid := "customuuid"
    HereNow(t, cipherKey, customUuid, testName)

}

func TestPresence(t *testing.T) {
    customUuid := "customuuid"
    testName := "Presence"
    t.Parallel()
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, customUuid)  
    channel := "testChannel"
    returnPresenceChannel := make(chan []byte)
    go SubscribeToPresence(channel, pubnubInstance, t, testName, customUuid, returnPresenceChannel)
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    subscribed := ParseSubscribeResponseForPresence(returnSubscribeChannel, channel)

    if(!subscribed) {
        t.Error("Test '" + testName + "': failed.")
    }
    //time.Sleep(10 * time.Second)
}

func SubscribeToPresence(channel string, pubnubInstance *pubnubMessaging.Pubnub, t *testing.T, testName string, customUuid string, returnPresenceChannel chan []byte){
    go pubnubInstance.Subscribe(channel, returnPresenceChannel, true)
    ParsePresenceResponse(pubnubInstance, t, returnPresenceChannel, channel, testName, customUuid, false)    
}

func ParseResponsePresence(channel chan []byte){
    for {
        value, ok := <-channel
        if !ok {  
            break
        }
        if string(value) != "[]"{
            fmt.Println(fmt.Sprintf("Presence: %s ", value))
            //fmt.Println(fmt.Sprintf("%s", value))
            fmt.Println("");
        }
    }
}

func HereNow(t *testing.T, cipherKey string, customUuid string, testName string){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", cipherKey, false, customUuid)  
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    subscribed := ParseSubscribeResponseForPresence(returnSubscribeChannel, channel)    
    if(subscribed){
        returnChannel := make(chan []byte)
        go pubnubInstance.HereNow(channel, returnChannel)
        ParseHereNowResponse(returnChannel, t, channel, customUuid, testName)
    } else {
        t.Error("Test '" + testName + "': failed.");
    }    
}

func ParseHereNowResponse(returnChannel chan []byte, t *testing.T, channel string, message string, testName string){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := fmt.Sprintf("%s", value)
            //fmt.Println("Test '" + testName + "':" +response)
            if(testName == "CustomUuid"){ 
                if(strings.Contains(response, message)){
                    fmt.Println("Test '" + testName + "': passed.")
                    break
                } else {
                    t.Error("Test '" + testName + "': failed.");
                    break
                }
            } else {
                var occupants struct {
                    Uuids []string
                    Occupancy int
                }
                
                err := json.Unmarshal(value, &occupants)
                if(err != nil) { 
                    fmt.Println("Test '" + testName + "':",err)
                    t.Error("Test '" + testName + "': failed.");
                    break
                } else {
                    i := occupants.Occupancy
                    if(i <= 0){    
                        t.Error("Test '" + testName + "': failed.");
                        break
                    } else {
                        fmt.Println("Test '" + testName + "': passed.");
                    }
                }
            }        
        }
    }
}   

func ParseSubscribeResponseForPresence(returnChannel chan []byte, channel string) bool{
    for {
        value, ok := <-returnChannel
        
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := fmt.Sprintf("%s", value)
            message := "'" + channel + "' connected"
            messageReconn := "'" + channel + "' reconnected"
            if((strings.Contains(response, message)) || (strings.Contains(response, messageReconn))){
                return true
            }else {
                break
            } 
        }
    }
    return false
}

func ParsePresenceResponse(pubnubInstance *pubnubMessaging.Pubnub, t *testing.T, returnChannel chan []byte, channel string, testName string, customUuid string, testConnected bool) bool {
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := fmt.Sprintf("%s", value)
            //fmt.Println("resp:", response)
            messagePresence := "Presence notifications for channel '" + channel + "' connected"
            messagePresenceReconn := "Presence notifications for channel '" + channel + "' reconnected"
            
            if (testConnected && ((strings.Contains(response, messagePresence)) || (strings.Contains(response, messagePresenceReconn)))){
                return true
            } else if(!testConnected) {
                
                message := "'" + channel + "' disconnected"
                messageConn := "'" + channel + "' connected"
                messageReconn := "'" + channel + "' reconnected"
                if(!strings.Contains(response, message) && (!strings.Contains(response, messageConn)) && (!strings.Contains(response, messageReconn))){
                    data, _, returnedChannel, err2 := pubnubMessaging.ParseJson(value, "")
                    var occupants []struct {
                        Action string
                        Timestamp float64 `json:",string"`
                        Uuid string
                        Occupancy int64 `json:",string"`
                    }
                    if(err2 != nil){
                        fmt.Println("err2 '" + testName + "':",err2)
                        t.Error("Test '" + testName + "': failed.");
                        break
                    }
                    fmt.Println("data '" + testName + "':",data)
                    //fmt.Println("ts '" + testName + "':",ts)
                    //fmt.Println("returnedChannel '" + testName + "':",returnedChannel)
                    err := json.Unmarshal([]byte(data), &occupants)
                    if(err != nil) { 
                        fmt.Println("err '" + testName + "':",err)
                        fmt.Println("Test '" + testName + "': failed.")    
                        t.Error("Test '" + testName + "': failed.");
                        break
                    } else {
                        channelSubRepsonseReceived := false
                        for i:=0; i<len(occupants); i++ {
                            if((occupants[i].Action == "join") && occupants[i].Uuid == customUuid){
                                channelSubRepsonseReceived = true
                                fmt.Println("Test '" + testName + "': failed. Err1")    
                                t.Error("Test '" + testName + "': failed.");
                                break
                            }
                        }
                        if(!channelSubRepsonseReceived){
                            fmt.Println("Test '" + testName + "': failed. Err2")
                            t.Error("Test '" + testName + "': failed.");
                            break
                        }
                        if(channel == returnedChannel){
                            fmt.Println("Test '" + testName + "': passed.")
                            return true    
                        } else {
                            fmt.Println("Test '" + testName + "': failed. Err3")
                            t.Error("Test '" + testName + "': failed.");
                            break
                        }
                    } 
                }               
            }
        }
    }
    return false
}

// End indicator
func TestPresenceEnd(t *testing.T){
    PrintTestMessage("==========Presence tests end==========")
}