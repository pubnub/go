package pubnubTests

import (
    "testing"
    "pubnubMessaging"
    "fmt"
    "strings"
    "encoding/json"
)

// Start indicator
func TestSubscribeStart(t *testing.T){
	PrintTestMessage("==========Subscribe tests start==========")
}

func TestSubscriptionForComplexMessage(t *testing.T) {
    pubnubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", "", false, "")    
    
    channel := "testChannel"
    
    returnSubscribeChannel := make(chan []byte)
    go pubnubInstance.Subscribe(channel, returnSubscribeChannel, false)
    ParseSubscribeResponse(pubnubInstance, returnSubscribeChannel, t, channel, "", "SubscriptionConnectedForComplex")    
}

func PublishComplexMessage(pubnubInstance *pubnubMessaging.Pubnub, t *testing.T, channel string, testName string){
	customComplexMessage := InitComplexMessage()
	
	returnChannel := make(chan []byte)
    go pubnubInstance.Publish(channel, customComplexMessage, returnChannel)
    
    ParseSubscribeResponse(pubnubInstance, returnChannel, t, channel, "", testName)    
}

func ParseSubscribeData (t *testing.T, response []byte, testName string){
	if(response != nil){
		var customComplexMessage CustomComplexMessage
		data, _, _, err := pubnubMessaging.ParseJson(response)
		if(err != nil){
				fmt.Println("err1:", err)
				fmt.Println("Test '" + testName + "': failed.")
		} else {
			fmt.Println("data:", data)
			
			//var objmap map[string]*json.RawMessage
			//err := json.Unmarshal(data, &objmap)
			//err := json.Unmarshal([]byte(data), &objmap)
			err = json.Unmarshal([]byte(data), &customComplexMessage)
			
			/*for k, v := range m {
		        switch vv := v.(type) {
		        case string:
		            fmt.Println(k, "is string", vv)
		        case int:
		            fmt.Println(k, "is int", vv)
		        case []interface{}:
		            fmt.Println(k, "is an array:")
		            for i, u := range vv {
		                fmt.Println(i, u)
		            }
		        default:
		            fmt.Println(k, "is of a type I don't know how to handle")
		        }
			}*/			
			if(err != nil){
				fmt.Println("err2:", err)
				fmt.Println("Test '" + testName + "': failed.")
			} else {
				fmt.Println("Test '" + testName + "': passed.")
			}
		}
	}
}

func ParseSubscribeResponse(pubnubInstance *pubnubMessaging.Pubnub, returnChannel chan []byte, t *testing.T, channel string, message string, testName string){
    for {
        value, ok := <-returnChannel
        if !ok {
            break
        }
        if string(value) != "[]"{
            response := fmt.Sprintf("%s", value)
            //fmt.Println(response)
            if(testName == "SubscriptionConnectedForComplex"){
            	message = "'" + channel + "' connected"
            	if(strings.Contains(response, message)){
            		PublishComplexMessage(pubnubInstance, t, channel, publishSuccessMessage)
            	} else {
            		//fmt.Println(response)
            		ParseSubscribeData(t, value, testName)
            		break
            	}
            } else if (testName == "SubscriptionForComplexMessage"){
            	if(strings.Contains(response, message)){
            		
            	}	
            }    
        }
    }
 }

 
 // End indicator
func TestSubscribeEnd(t *testing.T){
	PrintTestMessage("==========Subscribe tests end==========")
}   