// Package pubnubMessaging has the unit tests of package pubnubMessaging.
// pubnubPresence_test.go contains the tests related to the presence requests on pubnub Api
package pubnubTests

import (
    "testing"
    "github.com/pubnub/go/3.4.1/pubnubMessaging"
    "strings"
    "fmt"
    "encoding/json"
    "time"
)

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

// HereNow is a common method used by the tests TestHereNow, HereNowWithCipher, CustomUuid
// It subscribes to a pubnub channel and then 
// makes a call to the herenow method of the pubnub api.
func HereNow(t *testing.T, cipherKey string, customUuid string, testName string){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", cipherKey, false, customUuid)  
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    errorChannel := make(chan []byte)
    responseChannel := make(chan string)
    waitChannel := make(chan string)
    
    go pubnubInstance.Subscribe(channel, "", returnSubscribeChannel, false, errorChannel)
    go ParseSubscribeResponseForPresence(pubnubInstance, customUuid, returnSubscribeChannel, channel, testName, responseChannel)
    go ParseErrorResponse(errorChannel, responseChannel)  
    go WaitForCompletion(responseChannel, waitChannel)
    ParseWaitResponse(waitChannel, t, testName)    
}

// ParseHereNowResponse parses the herenow response on the go channel.
// In case of customuuid it looks for the custom uuid in the response.
// And in other cases checks for the occupancy.
func ParseHereNowResponse(returnChannel chan []byte, channel string, message string, testName string, responseChannel chan string){
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
                    responseChannel <- "Test '" + testName + "': passed."
                    break
                } else {
                    responseChannel <- "Test '" + testName + "': failed."
                    break
                }
            } else {
                var occupants struct {
                    Uuids []string
                    Occupancy int
                }
                
                err := json.Unmarshal(value, &occupants)
                if(err != nil) { 
                    //fmt.Println("Test '" + testName + "':",err)
                    responseChannel <- "Test '" + testName + "': failed. Message: " + err.Error()
                    break
                } else {
                    i := occupants.Occupancy
                    if(i <= 0){    
                        responseChannel <- "Test '" + testName + "': failed. Occupancy mismatch";
                        break
                    } else {
                        responseChannel <- "Test '" + testName + "': passed."
                    }
                }
            }        
        }
    }
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
    errorChannel := make(chan []byte)
    responseChannel := make(chan string)
    waitChannel := make(chan string)
    
    go pubnubInstance.Subscribe(channel, "", returnPresenceChannel, true, errorChannel)
    go ParseSubscribeResponseForPresence(pubnubInstance, customUuid, returnPresenceChannel, channel, testName, responseChannel)
    go ParseResponseDummy(errorChannel)  
    go WaitForCompletion(responseChannel, waitChannel)
    ParseWaitResponse(waitChannel, t, testName)
}

// ParseSubscribeResponseForPresence will look for the connection status in the response 
// received on the go channel. 
func ParseSubscribeResponseForPresence(pubnubInstance *pubnubMessaging.Pubnub, customUuid string, returnChannel chan []byte, channel string, testName string, responseChannel chan string) {
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        //response := fmt.Sprintf("%s", value)
        
        if string(value) != "[]"{
            if ((testName == "CustomUuid") || (testName == "HereNow") || (testName == "HereNowWithCipher")){
                response := fmt.Sprintf("%s", value)
                message := "'" + channel + "' connected"
                messageReconn := "'" + channel + "' reconnected"
                if((strings.Contains(response, message)) || (strings.Contains(response, messageReconn))){
                    errorChannel := make(chan []byte)
                    returnChannel := make(chan []byte)
                    time.Sleep(3 * time.Second)
                    go pubnubInstance.HereNow(channel, returnChannel, errorChannel)
                    go ParseHereNowResponse(returnChannel, channel, customUuid, testName, responseChannel)
                    go ParseErrorResponse(errorChannel, responseChannel)
                    break
                }    
            } else {
                response := fmt.Sprintf("%s", value)
                message := "'" + channel + "' connected"
                messageReconn := "'" + channel + "' reconnected"
                //fmt.Println("Test3 '" + testName + "':" +response)
                if((strings.Contains(response, message)) || (strings.Contains(response, messageReconn))){
                    
                    errorChannel2 := make(chan []byte)
                    returnSubscribeChannel := make(chan []byte)
                    time.Sleep(1 * time.Second)
                    go pubnubInstance.Subscribe(channel, "", returnSubscribeChannel, false, errorChannel2)
                    go ParseResponseDummy(returnSubscribeChannel)
                    go ParseResponseDummy(errorChannel2)
                }else {
                    if(testName == "Presence") {
                        data, _, returnedChannel, err2 := pubnubMessaging.ParseJson(value, "")
                        
                        var occupants []struct {
                            Action string
                            Uuid string
                            Timestamp float64
                            Occupancy int
                        }
                        
                        if(err2 != nil){
                            responseChannel <- "Test '" + testName + "': failed. Message: 1 :" + err2.Error()
                            break
                        }
                        //fmt.Println("Test3 '" + testName + "':" +data)
                        err := json.Unmarshal([]byte(data), &occupants)
                        if(err != nil) { 
                            //fmt.Println("err '" + testName + "':",err)
                            responseChannel <- "Test '" + testName + "': failed. Message: 2 :" + err.Error()
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
                                responseChannel <- "Test '" + testName + "': failed. Message: err3"
                                break
                            }
                            if(channel == returnedChannel){
                                responseChannel <- "Test '" + testName + "': passed."
                                break
                            } else {
                                responseChannel <- "Test '" + testName + "': failed. Message: err4"
                                break
                            }
                        }                         
                    }
                } 
            }
        }
    }
}

// TestPresenceEnd prints a message on the screen to mark the end of 
// presence tests.
// PrintTestMessage is defined in the common.go file.
func TestPresenceEnd(t *testing.T){
    PrintTestMessage("==========Presence tests end==========")
}