package pubnubTests

import (
    "testing"
    "pubnubMessaging"
    "strings"
    "fmt"
    //"strconv"
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
	
	pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, customUuid)  
	
	channel := "testChannel"
	
	returnPresenceChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnPresenceChannel, true)
    ParsePresenceResponse(t, returnPresenceChannel, channel, testName)    

	returnSubscribeChannel := make(chan []byte)
	go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
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
			if(strings.Contains(response, message)){
                return true
			}else {
				break
			} 
        }
    }
    return false
}

func ParsePresenceResponse(t *testing.T, returnChannel chan []byte, channel string, testName string) {
	connected := false
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := fmt.Sprintf("%s", value)
 			message := "'" + channel + "' connected"
			if ((!connected) && (strings.Contains(response, message))){
				connected = true
			} 
			if(connected){
				data, _, returnedChannel, err2 := pubnubMessaging.ParseJson(value, "")
				var occupants []struct {
					Action string
					Timestamp string
					Uuid string
					Occupancy string
	        	}
	        	if(err2 != nil){
	        		fmt.Println("err2 '" + testName + "':",err2)
            		t.Error("Test '" + testName + "': failed.");
            		break
	        	}
	        	fmt.Println("data '" + testName + "':",data)
	        	fmt.Println("returnedChannel '" + testName + "':",returnedChannel)
	        	err := json.Unmarshal([]byte(data), &occupants)
	        	if(err != nil) { 
	        		fmt.Println("err '" + testName + "':",err)
            		t.Error("Test '" + testName + "': failed.");
            		break
	            } else {
		        	if(channel == returnedChannel){
		        		fmt.Println("Test '" + testName + "': passed.");	
		        		break
		        	} else {
		        		t.Error("Test '" + testName + "': failed.");
		        		break
		        	}
		        }	        	
			}
        }
    }
}

// End indicator
func TestPresenceEnd(t *testing.T){
    PrintTestMessage("==========Presence tests end==========")
}