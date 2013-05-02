package pubnubTests

import (
    "testing"
    "pubnubMessaging"
    "strings"
    "fmt"
)

// Start indicator
func TestPublishStart(t *testing.T){
    PrintTestMessage("==========Publish tests start==========")
}

func TestNullMessage(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    channel := "testChannel"
    var message interface{}
    message = nil
    returnChannel := make(chan []byte)
    
    go pubnubInstance.Publish(channel, message, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, "Invalid Message", "NullMessage")
}

func TestUniqueGuid(t *testing.T){
    guid, err := pubnubMessaging.GenUuid()
    if(err != nil){
        fmt.Println("err: ", err)
        t.Error("Test 'UniqueGuid': failed.");
    } else if (guid == ""){
        t.Error("Test 'UniqueGuid': failed.");
    } else {
        fmt.Println("Test 'UniqueGuid': passed.");
    }
}

func TestSuccessCodeAndInfo(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    channel := "testChannel"
    message := "Pubnub API Usage Example"
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfo")
}

func TestSuccessCodeAndInfoWithEncryption(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "enigma", false, "")    
    channel := "testChannel"
    message := "Pubnub API Usage Example"
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfoWithEncryption")
}

func TestSuccessCodeAndInfoWithSecretAndEncryption(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "secret", "enigma", false, "")    
    channel := "testChannel"
    message := "Pubnub API Usage Example"
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfoWithSecretAndEncryption")
}

func TestSuccessCodeAndInfoForComplexMessage(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    channel := "testChannel"
    
    customStruct := CustomStruct{
        Foo : "hi!",
        Bar : []int{1,2,3,4,5},
    }
    
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, customStruct, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfoForComplexMessage")
}

func TestSuccessCodeAndInfoForComplexMessage2(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    channel := "testChannel"
    
    customComplexMessage := InitComplexMessage()
    
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, customComplexMessage, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfoForComplexMessage2")
}

func TestSuccessCodeAndInfoForComplexMessage2WithSecretAndEncryption(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "secret", "enigma", false, "")    
    channel := "testChannel"
    
    customComplexMessage := InitComplexMessage()
    
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, customComplexMessage, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfoForComplexMessage2WithSecretAndEncryption")
}

func TestSuccessCodeAndInfoForComplexMessage2WithEncryption(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "enigma", false, "")    
    channel := "testChannel"
    
    customComplexMessage := InitComplexMessage()
    
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, customComplexMessage, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfoForComplexMessage2WithEncryption")
}


func ParsePublishResponse(returnChannel chan []byte, t *testing.T, channel string, message string, testname string){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := fmt.Sprintf("%s", value)
            //fmt.Println("Test '" + testname + "':" +response)
            if(strings.Contains(response, message)){
                fmt.Println("Test '" + testname + "': passed.")
                break
            } else {
                t.Error("Test '" + testname + "': failed.");
                break
            }
        }
    }
 }   
 
 // End indicator
func TestPublishEnd(t *testing.T){
    PrintTestMessage("==========Publish tests end==========")
}