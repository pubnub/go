// Package pubnubMessaging has the unit tests of package pubnubMessaging.
// pubnubSubscribe_test.go contains the tests related to the Subscribe requests on pubnub Api
package pubnubTests

import (
    "testing"
    "github.com/pubnub/go/3.4/pubnubMessaging"
    "fmt"
    "strings"
    "encoding/json"
    "strconv"
    "net/url"
    "encoding/xml"
)

// TestSubscribeStart prints a message on the screen to mark the beginning of 
// subscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestSubscribeStart(t *testing.T){
    PrintTestMessage("==========Subscribe tests start==========")
}

// TestSubscriptionConnectStatus sends out a subscribe request to a pubnub channel
// and validates the response for the connect status.
func TestSubscriptionConnectStatus(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "", "SubscriptionConnectStatus", "")    
}

// TestSubscriptionAlreadySubscribed sends out a subscribe request to a pubnub channel
// and when connected sends out another subscribe request. The response for the second 
func TestSubscriptionAlreadySubscribed(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    returnSubscribeChannel2 := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel2, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "already subscribed", "SubscriptionAlreadySubscribed", "")    
}

// TestMultiSubscriptionConnectStatus send out a pubnub multi channel subscribe request and 
// parses the response for multiple connection status.
func TestMultiSubscriptionConnectStatus(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    testName := "TestMultiSubscriptionConnectStatus"
    channels := "testChannel1,testChannel2"

    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channels, returnSubscribeChannel, false)
    ParseSubscribeResponseForMultipleChannels(returnSubscribeChannel, channels, t, testName)    
}

// ParseSubscribeResponseForMultipleChannels parses the pubnub multi channel response 
// for the number or channels connected and matches them to the connected channels.
func ParseSubscribeResponseForMultipleChannels(returnChannel chan []byte, channels string, t *testing.T, testName string){
    noOfChannelsConnected := 0
    channelArray := strings.Split(channels, ",");
    loops := 0
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := fmt.Sprintf("%s", value)
            message := "' connected"
            messageReconn := "' reconnected"
            if((strings.Contains(response, message)) || (strings.Contains(response, messageReconn))){
                noOfChannelsConnected++
                if(noOfChannelsConnected >= len(channelArray)){
                    fmt.Println("Test '" + testName + "': passed.")
                    break
                }
            } 
        }
        loops++
        if(loops > 30){
            t.Error("Test '" + testName + "': failed.");
            break	
        }
    }
}

// TestSubscriptionForSimpleMessage first subscribes to a pubnub channel and then publishes 
// a message on the same pubnub channel. The subscribe response should receive this same message.  
func TestSubscriptionForSimpleMessage(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "", "SubscriptionConnectedForSimple", "")    
}

// TestSubscriptionForSimpleMessageWithCipher first subscribes to a pubnub channel and then publishes 
// an encrypted message on the same pubnub channel. The subscribe response should receive 
// the decrypted message.   
func TestSubscriptionForSimpleMessageWithCipher(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "enigma", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "", "SubscriptionConnectedForSimpleWithCipher", "enigma")    
}

// TestSubscriptionForComplexMessage first subscribes to a pubnub channel and then publishes 
// a complex message on the same pubnub channel. The subscribe response should receive 
// the same message.  
func TestSubscriptionForComplexMessage(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "", "SubscriptionConnectedForComplex", "")    
}

// TestSubscriptionForComplexMessageWithCipher first subscribes to a pubnub channel and then publishes 
// an encrypted complex message on the same pubnub channel. The subscribe response should receive 
// the decrypted message.   
func TestSubscriptionForComplexMessageWithCipher(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "enigma", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "", "SubscriptionConnectedForComplexWithCipher", "enigma")    
}

// PublishComplexMessage publises a complex message on a pubnub channel and 
// calls the parse method to validate the message subscription.
// CustomComplexMessage and InitComplexMessage are defined in the common.go file.
func PublishComplexMessage(pubnubInstance *pubnubMessaging.Pubnub, t *testing.T, channel string, testName string, cipherKey string){
    customComplexMessage := InitComplexMessage()
    
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, customComplexMessage, returnChannel)
    
    ParseSubscribeResponse(pubnubInstance, returnChannel, t, channel, "", testName, cipherKey)    
}

// PublishSimpleMessage publises a message on a pubnub channel and 
// calls the parse method to validate the message subscription.
func PublishSimpleMessage(pubnubInstance *pubnubMessaging.Pubnub, t *testing.T, channel string, testName string, cipherKey string){
    message := "Test message"
    
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnChannel)
    
    ParseSubscribeResponse(pubnubInstance, returnChannel, t, channel, "", testName, cipherKey)    
}

// ValidateComplexData takes an interafce as a parameter and iterates through it 
// It validates each field of the response interface against the initialized struct
// CustomComplexMessage, ReplaceEncodedChars and InitComplexMessage are defined in the common.go file.
func ValidateComplexData(m map[string]interface{}) (bool){
    //m := vv.(map[string]interface{})
    //if val,ok := m["VersionId"]; ok {
        //fmt.Println("VersionId",m["VersionId"])
    //}
    customComplexMessage := InitComplexMessage()            
    valid := false            
    for k, v := range m {
         //fmt.Println("k:", k, "v:", v)
         if (k == "OperationName") {
             if(m["OperationName"].(string) == customComplexMessage.OperationName){
                 valid = true
             } else {
                 return false
             }
        } else if (k == "VersionId") {
             if a, ok := v.(string); ok {
                verId, convErr := strconv.ParseFloat(a, 64)
                if (convErr != nil) { 
                    fmt.Println(convErr)
                    return false 
                }else{
                     if(float32(verId) == customComplexMessage.VersionId){
                         valid = true
                     } else {
                         return false
                     }
                 }    
             }
        } else if (k == "TimeToken") {
            i, convErr := strconv.ParseInt(v.(string), 10, 64)
            if (convErr != nil) { 
                fmt.Println(convErr)
                return false 
            }else{             
                 if(i == customComplexMessage.TimeToken){
                     valid = true
                 } else {
                     return false
                 }
             }    
        } else if (k == "DemoMessage") {
            b1 := v.(map[string]interface{})
            jsonData, _ := json.Marshal(customComplexMessage.DemoMessage.DefaultMessage)
            if val,ok := b1["DefaultMessage"]; ok {
                if(val.(string) != string(jsonData)){
                    return false 
                } else {
                    valid = true
                }
            }
        } else if (k == "SampleXml") {
            data := &Data{}
            s1, _ := url.QueryUnescape(m["SampleXml"].(string))
            
            reader := strings.NewReader(ReplaceEncodedChars(s1))
            err := xml.NewDecoder(reader).Decode(&data)
            
            if(err != nil){
                fmt.Println(err)
                return false 
            } else {
                 jsonData, _ := json.Marshal(customComplexMessage.SampleXml)
                 if(s1 == string(jsonData)){
                     valid = true
                 } else {
                     return false
                 }
             }    
        } else if (k == "Channels") {
             strSlice1, _ := json.Marshal(v) 
            strSlice2, _ := json.Marshal(customComplexMessage.Channels)
            s1, err := url.QueryUnescape(string(strSlice1))    
            if(err != nil){    
                fmt.Println(err)
                return false 
            } else {
                 if(s1 == string(strSlice2)){
                     valid = true
                 } else {
                     return false
                 }
             }    
        }
    }
    return valid
}

// CheckComplexData iterates through the json interafce and will read when 
// map type is encountered. 
// CustomComplexMessage and InitComplexMessage are defined in the common.go file.
func CheckComplexData(b interface{}) bool{
    valid := false
    switch vv := b.(type) {
        case string:
            //fmt.Println( "is string", vv)
        case int:
            //fmt.Println( "is int", vv)
        case []interface{}:
            //fmt.Println( "is an array:")
            //for i, u := range vv {
            for _, u := range vv {
                return CheckComplexData(u)
                //fmt.Println(i, u)
            }
        case map[string]interface{}:   
            m := vv 
            return ValidateComplexData(m)
        default:
        }
    return valid        
}

// ParseSubscribeData is used by multiple test cases and acts according to the testcase names.
// In case of complex message calls a sub method and in case of a simle message parses
// the response.
func ParseSubscribeData (t *testing.T, response []byte, testName string, cipherKey string){
    if(response != nil){
        var b interface{}
        err := json.Unmarshal(response, &b)

        isValid := false    
        if((testName == "SubscriptionConnectedForComplex") || (testName == "SubscriptionConnectedForComplexWithCipher")){            
            isValid = CheckComplexData(b)
        } else if((testName == "SubscriptionConnectedForSimple") || (testName == "SubscriptionConnectedForSimpleWithCipher")){
            data, _, _, err2 := pubnubMessaging.ParseJson(response, cipherKey)
            if(err2 != nil){
                fmt.Println("err2:", err2)
            } else {    
                var arr []string
                err3 := json.Unmarshal([]byte(data), &arr)
                if(err3 != nil){
                    fmt.Println("err3:", err2)
                } else {
                    if(len(arr)>0){
                        if(arr[0] == "Test message"){
                            isValid = true
                        }
                    }
                }    
            }                    
        }
        if(err != nil){
            fmt.Println("err2:", err)
            fmt.Println("Test '" + testName + "': failed.")
            t.Error("Test '" + testName + "': failed.");
        } else if (!isValid){
            fmt.Println("Test '" + testName + "': failed.")
            t.Error("Test '" + testName + "': failed.");
        } else {
            fmt.Println("Test '" + testName + "': passed.")
        }
    }
}

// ParseSubscribeResponse reads the response from the go channel and unmarshal's it.
// It is used by multiple test cases and acts according to the testcase names.
// The idea is to parse each message in the response based on the type of message
// and test against the sent message. If both match the test case is successful.
// publishSuccessMessage is defined in the common.go file. 
func ParseSubscribeResponse(pubnubInstance *pubnubMessaging.Pubnub, returnChannel chan []byte, t *testing.T, channel string, message string, testName string, cipherKey string){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if (string(value) != "[]"){
            response := fmt.Sprintf("%s", value)
            //fmt.Println("Response:", response)
            if((testName == "SubscriptionConnectedForComplex") || (testName == "SubscriptionConnectedForComplexWithCipher")){
                message = "'" + channel + "' connected"
                if(strings.Contains(response, message)){
                    PublishComplexMessage(pubnubInstance, t, channel, publishSuccessMessage, cipherKey)
                } else {
                    //fmt.Println("resp:", response)
                    ParseSubscribeData(t, value, testName, cipherKey)
                    break
                }
            } else if((testName == "SubscriptionConnectedForSimple") || (testName == "SubscriptionConnectedForSimpleWithCipher")){
                message = "'" + channel + "' connected"
                if(strings.Contains(response, message)){
                    PublishSimpleMessage(pubnubInstance, t, channel, publishSuccessMessage, cipherKey)
                } else {
                    //fmt.Println(response)
                    ParseSubscribeData(t, value, testName, cipherKey)
                    break
                }
            } else if (testName == "SubscriptionForComplexMessage"){
                message = "'" + channel + "' connected"
                if(strings.Contains(response, message)){
                    
                }    
            } else if (testName == "SubscriptionAlreadySubscribed"){
                message = "'" + channel + "' connected"
                if(strings.Contains(response, message)){
                    fmt.Println("Test '" + testName + "': passed.")
                } else {
                    fmt.Println("Test '" + testName + "': failed.")
                    t.Error("Test '" + testName + "': failed.");
                }   
                break                 
            } else if (testName == "SubscriptionConnectStatus"){
                message = "'" + channel + "' connected"
                if(strings.Contains(response, message)){
                    fmt.Println("Test '" + testName + "': passed.")
                } else {
                    fmt.Println("Test '" + testName + "': failed.")
                    t.Error("Test '" + testName + "': failed.");
                }   
                break                 
            }
        }
    }
 }

// TestSubscribeEnd prints a message on the screen to mark the end of 
// subscribe tests.
// PrintTestMessage is defined in the common.go file.
func TestSubscribeEnd(t *testing.T){
    PrintTestMessage("==========Subscribe tests end==========")
}   