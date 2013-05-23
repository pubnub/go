// Package pubnubMessaging has the unit tests of package pubnubMessaging.
// pubnubPresence_test.go contains the tests related to the presence requests on pubnub Api
package pubnubTests

import (
    "testing"
    "github.com/pubnub/go/3.4/pubnubMessaging"
    "strings"
    "fmt"
    "time"
    "encoding/json"
)

// used in TestPresence
var _endPresenceTestAsFailure = false
// used in TestPresence
var _endPresenceTestAsSuccess = false

// TestPresenceStart prints a message on the screen to mark the beginning of 
// presence tests.
// PrintTestMessage is defined in the common.go file.
func TestPresenceStart(t *testing.T){
    PrintTestMessage("==========Presence tests start==========")
}

// TestCustomUuid subscribes to a pubnub channel using a custom uuid and then 
// makes a call to the herenow method of the pubnub api. The custom id should
// be present in the response else the test fails.
func TestCustomUuid(t *testing.T) {
    cipherKey := ""
    testName := "CustomUuid"
    customUuid := "customuuid"
    HereNow(t, cipherKey, customUuid, testName)
}

// TestHereNow subscribes to a pubnub channel and then 
// makes a call to the herenow method of the pubnub api. The occupancy should
// be greater than one.
func TestHereNow(t *testing.T) {
    cipherKey := ""
    testName := "HereNow"
    customUuid := ""
    HereNow(t, cipherKey, customUuid, testName)
}

// TestHereNowWithCipher subscribes to a pubnub channel and then 
// makes a call to the herenow method of the pubnub api. The occupancy should
// be greater than one.
func TestHereNowWithCipher(t *testing.T) {
    cipherKey := ""
    testName := "HereNowWithCipher"
    customUuid := "customuuid"
    HereNow(t, cipherKey, customUuid, testName)
}

// TestPresence subscribes to the presence notifications on a pubnub channel and 
// then subscribes to a pubnub channel. The test waits till we get a response from 
// the subscribe call. The method that parses the presence response sets the global 
// variable _endPresenceTestAsSuccess to true if the presence contains a join info
// on the channel and _endPresenceTestAsFailure is otherwise.
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

// SubscribeRoutine subscribes to a pubnub channel and waits for the response.
// Used as a go routine.
func SubscribeRoutine(channel string, pubnubInstance *pubnubMessaging.Pubnub){
    var subscribeChannel = make(chan []byte)
    go pubnubInstance.Subscribe(channel, subscribeChannel, false)
    ParseSubscribeResponseForPresence(subscribeChannel, channel)   
}

// SubscribeRoutine presence notifications to a pubnub channel and waits for the response.
// Used as a go routine.
func SubscribeToPresence(channel string, pubnubInstance *pubnubMessaging.Pubnub, t *testing.T, testName string, customUuid string, returnPresenceChannel chan []byte){
    go pubnubInstance.Subscribe(channel, returnPresenceChannel, true)
    ParsePresenceResponse(pubnubInstance, t, returnPresenceChannel, channel, testName, customUuid, false)    
}

// HereNow is a common method used by the tests TestHereNow, HereNowWithCipher, CustomUuid
// It subscribes to a pubnub channel and then 
// makes a call to the herenow method of the pubnub api.
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

// ParseHereNowResponse parses the herenow response on the go channel.
// In case of customuuid it looks for the custom uuid in the response.
// And in other cases checks for the occupancy.
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

// ParseSubscribeResponseForPresence will look for the connection status in the response 
// received on the go channel. 
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

// The method that parses the presence response sets the global 
// variable _endPresenceTestAsSuccess to true if the presence contains a join info
// on the channel and _endPresenceTestAsFailure is otherwise.
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

// TestPresenceEnd prints a message on the screen to mark the end of 
// presence tests.
// PrintTestMessage is defined in the common.go file.
func TestPresenceEnd(t *testing.T){
    PrintTestMessage("==========Presence tests end==========")
}