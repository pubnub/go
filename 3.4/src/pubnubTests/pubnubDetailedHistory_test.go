package pubnubTests

import (
    "testing"
    "pubnubMessaging"
    "strings"
    "fmt"
    "strconv"
    "encoding/json"
)

// Start indicator
func TestDetailedHistoryStart(t *testing.T){
    PrintTestMessage("==========DetailedHistory tests start==========")
}
    
func TestDetailedHistory(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    message := "Test Message"
    returnPublishChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnPublishChannel)
    ParseResponse(returnPublishChannel, t, pubnubInstance, channel, message, "DetailedHistory", 1)
}

func TestEncryptedDetailedHistory(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "enigma", false, "")    
    
    channel := "testChannel"
    message := "Test Message"
    returnPublishChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, message, returnPublishChannel)
    ParseResponse(returnPublishChannel, t, pubnubInstance, channel, message, "EncryptedDetailedHistory", 1)
}

func TestDetailedHistoryFor10Messages(t *testing.T) {
    testName := "TestDetailedHistoryFor10Messages"
    DetailedHistoryFor10Messages(t, "", testName)
}

func TestDetailedHistoryFor10EncryptedMessages(t *testing.T) {
    testName := "TestDetailedHistoryFor10EncryptedMessages"
    DetailedHistoryFor10Messages(t, "enigma", testName)
}

func DetailedHistoryFor10Messages(t *testing.T, cipherKey string, testName string) {
    numberOfMessages := 10
    startMessagesFrom := 0
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", cipherKey, false, "")    
    
    message := "Test Message "
    channel := "testChannel"
    
    messagesSent := PublishMessages(pubnubInstance, channel, t, startMessagesFrom, numberOfMessages, message)
    
    if(messagesSent){
    returnHistoryChannel := make(chan []byte)
    go pubnubInstance.History(channel, numberOfMessages, 0, 0, false, returnHistoryChannel)
        ParseHistoryResponseForMultipleMessages(returnHistoryChannel, t, channel, message, testName, startMessagesFrom, numberOfMessages, cipherKey)
    }else{
        t.Error("Test '" + testName + "': failed.");
    }    
}

func TestDetailedHistoryParamsFor10MessagesWithSeretKey(t *testing.T) {
    testName := "TestDetailedHistoryFor10MessagesWithSeretKey"
    DetailedHistoryParamsFor10Messages(t, "", "secret", testName)
}

func TestDetailedHistoryParamsFor10EncryptedMessagesWithSeretKey(t *testing.T) {
    testName := "TestDetailedHistoryFor10EncryptedMessagesWithSeretKey"
    DetailedHistoryParamsFor10Messages(t, "enigma", "secret", testName)
}

func TestDetailedHistoryParamsFor10Messages(t *testing.T) {
    testName := "TestDetailedHistoryFor10Messages"
    DetailedHistoryParamsFor10Messages(t, "", "", testName)
}

func TestDetailedHistoryParamsFor10EncryptedMessages(t *testing.T) {
    testName := "TestDetailedHistoryFor10EncryptedMessages"
    DetailedHistoryParamsFor10Messages(t, "enigma", "", testName)
}

func DetailedHistoryParamsFor10Messages(t *testing.T, cipherKey string, secretKey string, testName string) {
    numberOfMessages := 5
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", secretKey, cipherKey, false, "")    
    
    message := "Test Message "
    channel := "testChannel"
    
    startTime := GetServerTime(pubnubInstance, t, testName) 
    startMessagesFrom := 0
    messagesSent := PublishMessages(pubnubInstance, channel, t, startMessagesFrom, numberOfMessages, message)
    
    midTime := GetServerTime(pubnubInstance, t, testName)
    startMessagesFrom = 5
    messagesSent2 := PublishMessages(pubnubInstance, channel, t, startMessagesFrom, numberOfMessages, message)
    endTime := GetServerTime(pubnubInstance, t, testName)
    
    startMessagesFrom = 0
    if(messagesSent){
    returnHistoryChannel := make(chan []byte)
    go pubnubInstance.History(channel, numberOfMessages, startTime, midTime, false, returnHistoryChannel)
        ParseHistoryResponseForMultipleMessages(returnHistoryChannel, t, channel, message, testName, startMessagesFrom, numberOfMessages, cipherKey)
    }else{
        t.Error("Test '" + testName + "': failed.");
    }
    
    startMessagesFrom = 5
    if(messagesSent2){
    returnHistoryChannel := make(chan []byte)
    go pubnubInstance.History(channel, numberOfMessages, midTime, endTime, false, returnHistoryChannel)
        ParseHistoryResponseForMultipleMessages(returnHistoryChannel, t, channel, message, testName, startMessagesFrom, numberOfMessages, cipherKey)
    }else{
        t.Error("Test '" + testName + "': failed.");
    }    
}

func GetServerTime(pubnubInstance *pubnubMessaging.Pubnub, t *testing.T, testName string) int64{
    returnTimeChannel := make(chan []byte)
    go pubnubInstance.GetTime(returnTimeChannel)
    return ParseServerTimeResponse(returnTimeChannel, t, testName)    
}

func ParseServerTimeResponse(returnChannel chan []byte,t *testing.T, testName string) int64 {
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := string(value)
            if(response != ""){
                var arr []int64
                err2 := json.Unmarshal(value, &arr)
                if(err2 != nil){
                    fmt.Println("err2", err2)
                    t.Error("Test '" + testName + "': failed.");
                    break
                }else {    
                return arr[0]
            }         
            } else {
                fmt.Println("response", response)
                t.Error("Test '" + testName + "': failed.");
                break
            }
        }
    }
    return 0
}

func PublishMessages(pubnubInstance *pubnubMessaging.Pubnub, channel string, t *testing.T, startMessagesFrom int, numberOfMessages int, message string) bool{
    messagesReceived := 0
    messageToSend := ""
    for i:=startMessagesFrom; i< startMessagesFrom+numberOfMessages; i++{
        messageToSend = message + strconv.Itoa(i)
    
        returnPublishChannel := make(chan []byte)
        go pubnubInstance.Publish(channel, messageToSend, returnPublishChannel)
        published := ParsePublishResponseForMultipleMessages(returnPublishChannel, t, channel, publishSuccessMessage, "PublishTenMessages")
        if (published) {
            messagesReceived++
        }
    }
    if(messagesReceived == numberOfMessages){
        return true
    }
    
    return false
}

func ParsePublishResponseForMultipleMessages(returnChannel chan []byte, t *testing.T, channel string, message string, testname string) bool{
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := fmt.Sprintf("%s", value)
            if(strings.Contains(response, message)){
                return true
            } else {
                return false
            }
        }
    }
    return false
}


func ParseHistoryResponseForMultipleMessages(returnChannel chan []byte, t *testing.T, channel string, message string, testName string, startMessagesFrom int, numberOfMessages int, cipherKey string){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            data, _, _, err := pubnubMessaging.ParseJson(value, cipherKey)
            if(err != nil) {
                t.Error("Test '" + testName + "': failed.");
            } else {
                var arr []string
                err2 := json.Unmarshal([]byte(data), &arr)
                if(err2 != nil){
                    t.Error("Test '" + testName + "': failed.");
                }else {    
                    messagesReceived := 0
                    
                    if(len(arr) != numberOfMessages){
                        t.Error("Test '" + testName + "': failed.");
                        break
                    }
                    for i:=0; i < numberOfMessages; i++ {
                        if(arr[i] == message + strconv.Itoa(startMessagesFrom+i)){
                            //fmt.Println("data:",arr[i])
                            messagesReceived++
                        }
                    }   
                    if(messagesReceived == numberOfMessages){
                        fmt.Println("Test '" + testName + "': passed.")
                    } else {
                        t.Error("Test '" + testName + "': failed.");
                    }
                break
                }
            }
        }
    }
}

func ParseHistoryResponse(returnChannel chan []byte, t *testing.T, channel string, message string, testName string){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := string(value)
            //fmt.Println("response", response)
            if(strings.Contains(response, message)){
                fmt.Println("Test '" + testName + "': passed.") 
                break
            } else {
                t.Error("Test '" + testName + "': failed.");
            }
        }
    }
}

func ParseResponse(returnChannel chan []byte,t *testing.T, pubnubInstance *pubnubMessaging.Pubnub, channel string, message string, testName string, numberOfMessages int){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            returnHistoryChannel := make(chan []byte)
            go pubnubInstance.History(channel, 1, 0, 0, false, returnHistoryChannel)
            ParseHistoryResponse(returnHistoryChannel, t, channel, message, testName)
            break
        }
    }
}

// End indicator
func TestDetailedHistoryEnd(t *testing.T){
    PrintTestMessage("==========DetailedHistory tests end==========")
}