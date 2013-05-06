package pubnubTests

import (
    "testing"
    "pubnubMessaging"
    "fmt"
    "strings"
    "encoding/json"
    "strconv"
    "net/url"
    "encoding/xml"
    //"bytes"
    //"html"
    //"encoding/hex"
)

// Start indicator
func TestSubscribeStart(t *testing.T){
    PrintTestMessage("==========Subscribe tests start==========")
}

func TestSubscriptionConnectStatus(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "", "SubscriptionConnectStatus", "")    
}

func TestSubscriptionAlreadySubscribed(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    returnSubscribeChannel2 := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel2, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "already subscribed", "SubscriptionAlreadySubscribed", "")    
}

func TestMultiSubscriptionConnectStatus(t *testing.T) {
    /*pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    channel2 := "testChannel2"
    var buff bytes.Buffer
    buff.WriteString(channel)
    buff.WriteString(",")
    buff.WriteString(channel2)

    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(buff.String(), returnSubscribeChannel, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, buff.String(), "", "MultiSubscriptionConnectStatus")*/    
}

func TestSubscriptionForSimpleMessage(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "", "SubscriptionConnectedForSimple", "")    
}

func TestSubscriptionForSimpleMessageWithCipher(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "enigma", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "", "SubscriptionConnectedForSimpleWithCipher", "enigma")    
}


func TestSubscriptionForComplexMessage(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "", "SubscriptionConnectedForComplex", "")    
}

func TestSubscriptionForComplexMessageWithCipher(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "enigma", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "", "SubscriptionConnectedForComplexWithCipher", "enigma")    
}

func PublishComplexMessage(pubnubInstance *pubnubMessaging.Pubnub, t *testing.T, channel string, testName string, cipherKey string){
    customComplexMessage := InitComplexMessage()
    
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, customComplexMessage, returnChannel)
    
    ParseSubscribeResponse(pubnubInstance, returnChannel, t, channel, "", testName, cipherKey)    
}

func PublishSimpleMessage(pubnubInstance *pubnubMessaging.Pubnub, t *testing.T, channel string, testName string, cipherKey string){
    message := "Test message"
    
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnChannel)
    
    ParseSubscribeResponse(pubnubInstance, returnChannel, t, channel, "", testName, cipherKey)    
}

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

//TODO: merge with ParseSubscribeResponse
func ParseMultiSubscribeResponse(pubnubInstance *pubnubMessaging.Pubnub, returnChannel chan []byte, t *testing.T, channels string, message string, testName string){
    channelArray := strings.Split(channels, ",");
    responsesReceived := 0
    loops := 0
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if (string(value) != "[]"){
            loops++
            response := fmt.Sprintf("%s", value)
            if (testName == "MultiSubscriptionConnectStatus"){
                for i:=0; i < len(channelArray); i++ {
                    ch := strings.TrimSpace(channelArray[i])
                    message = "'" + ch + "' connected"
                    
                    if(strings.Contains(response, message)){
                        responsesReceived++
                    }                
                }
            }
        }
        if(responsesReceived >= len(channelArray)){
            fmt.Println("Test '" + testName + "': passed.")
            break
        }
        if(loops>10){
            t.Error("Test '" + testName + "': failed.");
            break
        }
    }
 }
 
 // End indicator
func TestSubscribeEnd(t *testing.T){
    PrintTestMessage("==========Subscribe tests end==========")
}   