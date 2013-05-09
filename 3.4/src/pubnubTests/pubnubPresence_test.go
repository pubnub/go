package pubnubTests

import (
    "testing"
    "pubnubMessaging"
    "strings"
    "fmt"
    "time"
    "encoding/json"
)

var _endPresenceTestAsFailure = false
var _endPresenceTestAsSuccess = false

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
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, customUuid)  
    channel := "testForPresenceChannel"
    returnPresenceChannel := make(chan []byte)
    go SubscribeToPresence(channel, pubnubInstance, t, testName, customUuid, returnPresenceChannel)
    time.Sleep(500*time.Millisecond)
    go SubscribeRoutine(channel, pubnubInstance)

    _endPresenceTestAsSuccess = false
    _endPresenceTestAsFailure = false
    
    for {
        if(_endPresenceTestAsSuccess){
        	fmt.Println("Test '" + testName + "': passed.");
            break
        } else if (_endPresenceTestAsFailure) {
            t.Error("Test '" + testName + "': failed.");
        	break
        }
        time.Sleep(500*time.Millisecond) 
    }
}

func SubscribeRoutine(channel string, pubnubInstance *pubnubMessaging.Pubnub){
    var subscribeChannel = make(chan []byte)
    go pubnubInstance.Subscribe(channel, subscribeChannel, false)
    ParseSubscribeResponseForPresence(subscribeChannel, channel)   
}

func SubscribeToPresence(channel string, pubnubInstance *pubnubMessaging.Pubnub, t *testing.T, testName string, customUuid string, returnPresenceChannel chan []byte){
    go pubnubInstance.Subscribe(channel, returnPresenceChannel, true)
    ParsePresenceResponse(pubnubInstance, t, returnPresenceChannel, channel, testName, customUuid, false)    
}

func HereNow(t *testing.T, cipherKey string, customUuid string, testName string){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", cipherKey, false, customUuid)  
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponseForPresence(returnSubscribeChannel, channel)  
    time.Sleep(500*time.Millisecond)  
    returnChannel := make(chan []byte)
    go pubnubInstance.HereNow(channel, returnChannel)
    ParseHereNowResponse(returnChannel, t, channel, customUuid, testName)
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
            //fmt.Println("resp pres:", response)
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
                        Uuid string
                        Timestamp float64
                        Occupancy int
                    }
                    //fmt.Println("data '" + testName + "':",data)
                    if(err2 != nil){
                        fmt.Println("err2 '" + testName + "':",err2)
                        _endPresenceTestAsFailure = true
                        break
                    }
                    
                    err := json.Unmarshal([]byte(data), &occupants)
                    if(err != nil) { 
                        fmt.Println("err '" + testName + "':",err)
                        _endPresenceTestAsFailure = true
                        break
                    } else {
                        channelSubRepsonseReceived := false
                        for i:=0; i<len(occupants); i++ {
                            if((occupants[i].Action == "join") && occupants[i].Uuid == customUuid){
                                channelSubRepsonseReceived = true
                                break
                            }
                        }
                        if(!channelSubRepsonseReceived){
                            fmt.Println("Test '" + testName + "': failed. Err2")
                            _endPresenceTestAsFailure = true
                            break
                        }
                        if(channel == returnedChannel){
                            _endPresenceTestAsSuccess = true
                            return true    
                        } else {
                            fmt.Println("Test '" + testName + "': failed. Err3")
                            _endPresenceTestAsFailure = true
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