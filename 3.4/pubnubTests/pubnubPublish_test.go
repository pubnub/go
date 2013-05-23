// Package pubnubMessaging has the unit tests of package pubnubMessaging.
// pubnubPublish_test.go contains the tests related to the publish requests on pubnub Api
package pubnubTests

import (
    "testing"
    "github.com/pubnub/go/3.4/pubnubMessaging"
    "strings"
    "fmt"
)

// TestPublishStart prints a message on the screen to mark the beginning of 
// publish tests.
// PrintTestMessage is defined in the common.go file.
func TestPublishStart(t *testing.T){
    PrintTestMessage("==========Publish tests start==========")
}

// TestNullMessage sends out a null message to a pubnub channel. The response should
// be an "Invalid Message".
func TestNullMessage(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    channel := "testChannel"
    var message interface{}
    message = nil
    returnChannel := make(chan []byte)
    
    go pubnubInstance.Publish(channel, message, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, "Invalid Message", "NullMessage")
}

// TestUniqueGuid tests the generation of a unique GUID for the client.
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

// TestSuccessCodeAndInfo sends out a message to the pubnub channel
// The response is parsed and should match the 'sent' status.
// publishSuccessMessage is defined in the common.go file
func TestSuccessCodeAndInfo(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    channel := "testChannel"
    message := "Pubnub API Usage Example"
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfo")
}

// TestSuccessCodeAndInfoWithEncryption sends out an encrypted 
// message to the pubnub channel
// The response is parsed and should match the 'sent' status.
// publishSuccessMessage is defined in the common.go file
func TestSuccessCodeAndInfoWithEncryption(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "enigma", false, "")    
    channel := "testChannel"
    message := "Pubnub API Usage Example"
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfoWithEncryption")
}

// TestSuccessCodeAndInfoWithSecretAndEncryption sends out an encrypted 
// secret keyed message to the pubnub channel
// The response is parsed and should match the 'sent' status.
// publishSuccessMessage is defined in the common.go file
func TestSuccessCodeAndInfoWithSecretAndEncryption(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "secret", "enigma", false, "")    
    channel := "testChannel"
    message := "Pubnub API Usage Example"
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfoWithSecretAndEncryption")
}

// TestSuccessCodeAndInfoForComplexMessage sends out a complex message to the pubnub channel
// The response is parsed and should match the 'sent' status.
// publishSuccessMessage and customstruct is defined in the common.go file
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

// TestSuccessCodeAndInfoForComplexMessage2 sends out a complex message to the pubnub channel
// The response is parsed and should match the 'sent' status.
// publishSuccessMessage and InitComplexMessage is defined in the common.go file
func TestSuccessCodeAndInfoForComplexMessage2(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    channel := "testChannel"
    
    customComplexMessage := InitComplexMessage()
    
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, customComplexMessage, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfoForComplexMessage2")
}

// TestSuccessCodeAndInfoForComplexMessage2WithSecretAndEncryption sends out an 
// encypted and secret keyed complex message to the pubnub channel
// The response is parsed and should match the 'sent' status.
// publishSuccessMessage and InitComplexMessage is defined in the common.go file
func TestSuccessCodeAndInfoForComplexMessage2WithSecretAndEncryption(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "secret", "enigma", false, "")    
    channel := "testChannel"
    
    customComplexMessage := InitComplexMessage()
    
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, customComplexMessage, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfoForComplexMessage2WithSecretAndEncryption")
}

// TestSuccessCodeAndInfoForComplexMessage2WithEncryption sends out an 
// encypted complex message to the pubnub channel
// The response is parsed and should match the 'sent' status.
// publishSuccessMessage and InitComplexMessage is defined in the common.go file
func TestSuccessCodeAndInfoForComplexMessage2WithEncryption(t *testing.T){
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "enigma", false, "")    
    channel := "testChannel"
    
    customComplexMessage := InitComplexMessage()
    
    returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, customComplexMessage, returnChannel)
    
    ParsePublishResponse(returnChannel, t, channel, publishSuccessMessage, "SuccessCodeAndInfoForComplexMessage2WithEncryption")
}

// ParsePublishResponse parses the response from the pubnub api to validate the
// sent status. 
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
 
// TestPublishEnd prints a message on the screen to mark the end of 
// publish tests.
// PrintTestMessage is defined in the common.go file.
func TestPublishEnd(t *testing.T){
    PrintTestMessage("==========Publish tests end==========")
}