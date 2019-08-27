package e2e

import (
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	pubnub "github.com/sprucehealth/pubnub-go"
	"github.com/sprucehealth/pubnub-go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

import _ "net/http/pprof"
import "net/http"

var timeout = 3

func SubscribesLogsForQueryParams(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	assert := assert.New(t)
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.SecretKey = "sec-key"
	pn.Config.AuthKey = "myAuthKey"
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	pn.Subscribe().
		Channels([]string{"ch1", "ch2"}).
		QueryParam(queryParam).
		Execute()

	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-tic.C:
		tic.Stop()

	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	//fmt.Printf("Captured: %s", out)

	s := fmt.Sprintf("%s", out)
	expected2 := fmt.Sprintf("q1=v1")
	expected3 := fmt.Sprintf("q2=v2")

	assert.Contains(s, expected2)
	assert.Contains(s, expected3)

	//https://ps.pndsn.com/v1/auth/grant/sub-key/sub-c-e41d50d4-43ce-11e8-a433-9e6b275e7b64?m=1&auth=authkey1,authkey2&channel=ch1,ch2&timestamp=1535719219&pnsdk=PubNub-Go/4.1.3&uuid=pn-a83164fe-7ecf-42ab-ba14-d2d8e6eabd7a&r=1&w=1&signature=0SkyfvohAq8_0phVi0YhCL4c2ZRSPBVwCwQ9fANvPmM=

}

func TestRequestMesssageOverflow(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-message-overflow")

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.MessageQueueOverflowCount = 2
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	timestamp1 := GetTimetoken(pn)
	for i := 0; i < 3; i++ {
		message := fmt.Sprintf("message %d", i)
		pn.Publish().Channel(ch).Message(message).Execute()
	}

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					continue
				case pubnub.PNRequestMessageCountExceededCategory:
					doneSubscribe <- true
					break
				default:
					errChan <- fmt.Sprintf("error ===> %v", status)
					break
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				break
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				break
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().Channels([]string{ch}).Timetoken(timestamp1).Execute()
	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")
	}
}

/////////////////////////////
/////////////////////////////
// Structure
// - Channel Subscription
// - Groups Subscription
// - Misc
/////////////////////////////
/////////////////////////////

/////////////////////////////
// Channel Subscription
/////////////////////////////

func TestSubscribeUnsubscribe(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-u-ch")
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6063", nil))
	// }()

	pn := pubnub.NewPubNub(configCopy())
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	listener := pubnub.NewListener()

	go func() {
		for {
			//fmt.Println("listening...")
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					//fmt.Println("PNConnectedCategory...")
					doneSubscribe <- true
					break
				case pubnub.PNDisconnectedCategory:
					//fmt.Println("PNDisconnectedCategory...")
					doneUnsubscribe <- true
					break
				case pubnub.PNAcknowledgmentCategory:
					doneUnsubscribe <- true
					break
				case pubnub.PNCancelledCategory:
					continue
				default:
					//fmt.Println("default...", status)
					errChan <- fmt.Sprintf("error ===> %v", status)
					//break
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				//break
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				//break
			}
		}
		//fmt.Println("exit listening...")
	}()

	pn.AddListener(listener)

	pn.Subscribe().Channels([]string{ch}).Execute()

	tic1 := time.NewTicker(time.Duration(timeout) * time.Second * 3)
	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	case <-tic1.C:
		tic1.Stop()
		assert.Fail("timeout")
	}

	//fmt.Println("calling Unsubscribe...")
	time.Sleep(3 * time.Second)

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	tic := time.NewTicker(time.Duration(timeout) * time.Second * 3)
	select {
	case <-doneUnsubscribe:
		//fmt.Println("doneUnsubscribe...")
	case err := <-errChan:
		assert.Fail(err)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")
	}

	//fmt.Println("after select")
	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
}

func GenRandom() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func TestSubscribePublishUnsubscribeString(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := "hey"

	go SubscribePublishUnsubscribeMultiCommon(t, s, "", pubMessage, false, false)
	m := <-pubMessage
	msg := m.(string)
	assert.Equal(s, msg)
}

func TestSubscribePublishUnsubscribeStringEnc(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := "yay!"

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, false, false)
	m := <-pubMessage
	msg := m.(string)
	assert.Equal(s, msg)
}

func TestSubscribePublishUnsubscribeInt(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := 1

	go SubscribePublishUnsubscribeMultiCommon(t, s, "", pubMessage, false, false)
	m := <-pubMessage
	msg := m.(float64)
	assert.Equal(float64(1), msg)
}

func TestSubscribePublishUnsubscribeIntEnc(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := 1

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, false, false)
	m := <-pubMessage
	msg := m.(float64)
	assert.Equal(float64(1), msg)
}

func TestSubscribePublishUnsubscribePNOtherDisable(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := map[string]interface{}{
		"id":        2,
		"not_other": "123456",
		"pn_other":  "yay!",
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, true, false)
	m := <-pubMessage
	msg := m.(map[string]interface{})
	assert.Equal("123456", msg["not_other"])
	assert.Equal("yay!", msg["pn_other"])
}

func TestSubscribePublishUnsubscribePNOther(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := map[string]interface{}{
		"id":        1,
		"not_other": "12345",
		"pn_other":  "yay!",
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, false, false)
	m := <-pubMessage
	msg := m.(map[string]interface{})
	assert.Equal("12345", msg["not_other"])
	assert.Equal("yay!", msg["pn_other"])

}

func TestSubscribePublishUnsubscribePNOtherComplex(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s1 := map[string]interface{}{
		"id":        1,
		"not_other": "1234567",
	}
	s := map[string]interface{}{
		"id":        1,
		"not_other": "12345",
		"pn_other":  s1,
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, false, false)
	m := <-pubMessage
	msg := m.(map[string]interface{})
	assert.Equal("12345", msg["not_other"])
	if msgOther, ok := msg["pn_other"].(map[string]interface{}); !ok {
		assert.Fail("!map[string]interface{}")
	} else {
		assert.Equal("1234567", msgOther["not_other"])
	}

}

func TestSubscribePublishUnsubscribeInterfaceWithoutPNOther(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := map[string]interface{}{
		"id":        3,
		"not_other": "1234",
		"ss_other":  "yay!",
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "", pubMessage, false, false)
	m := <-pubMessage
	msg := m.(map[string]interface{})
	assert.Equal("1234", msg["not_other"])
	assert.Equal("yay!", msg["ss_other"])

}

func TestSubscribePublishUnsubscribeInterfaceWithoutPNOtherEnc(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := map[string]interface{}{
		"id":        4,
		"not_other": "123",
		"ss_other":  "yay!",
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, false, false)
	m := <-pubMessage
	msg := m.(map[string]interface{})
	assert.Equal("123", msg["not_other"])
	assert.Equal("yay!", msg["ss_other"])
}

type customStruct struct {
	Foo string
	Bar []int
}

func TestSubscribePublishUnsubscribeInterfaceCustom(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := customStruct{
		Foo: "hi!",
		Bar: []int{1, 2, 3, 4, 5},
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "", pubMessage, false, false)
	m := <-pubMessage
	//s1 := reflect.ValueOf(m)
	//fmt.Println("s:::", s1, s1.Type())
	if msg, ok := m.(map[string]interface{}); !ok {
		//fmt.Println(msg)
		assert.Fail("not map")
	} else {
		//fmt.Println(msg)
		//byt := []byte(message.Message)
		//fmt.Println(message.Message.(string))
		//err := json.Unmarshal(byt, &msg)
		//assert.Nil(err)
		assert.Equal("hi!", msg["Foo"])
		//assert.Equal("1", msg["Bar"].(map[string]interface{})[0])
		//assert.Equal("\"yay!\"", msg["pn_other"])
	}
}

func TestSubscribePublishUnsubscribeInterfaceWithoutCustomEnc(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := customStruct{
		Foo: "hi!",
		Bar: []int{1, 2, 3, 4, 5},
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, false, false)
	m := <-pubMessage
	//s1 := reflect.ValueOf(m)
	//fmt.Println("s:::", s1, s1.Type())
	if msg, ok := m.(map[string]interface{}); !ok {
		//fmt.Println(msg)
		assert.Fail("not map")
	} else {
		//fmt.Println(msg)
		//byt := []byte(message.Message)
		//fmt.Println(message.Message.(string))
		//err := json.Unmarshal(byt, &msg)
		//assert.Nil(err)
		assert.Equal("hi!", msg["Foo"])
		//assert.Equal("1", msg["Bar"].(map[string]interface{})[0])
		//assert.Equal("\"yay!\"", msg["pn_other"])
	}
}

func TestSubscribePublishUnsubscribeStringPost(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := "hey"

	go SubscribePublishUnsubscribeMultiCommon(t, s, "", pubMessage, false, true)
	m := <-pubMessage
	msg := m.(string)
	assert.Equal(s, msg)
}

func TestSubscribePublishUnsubscribeStringEncPost(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := "hey"

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, false, true)
	m := <-pubMessage
	msg := m.(string)
	assert.Equal(s, msg)
}

func TestSubscribePublishUnsubscribeIntPost(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := 1

	go SubscribePublishUnsubscribeMultiCommon(t, s, "", pubMessage, false, true)
	m := <-pubMessage
	msg := m.(float64)
	assert.Equal(float64(1), msg)
}

func TestSubscribePublishUnsubscribeIntEncPost(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := 1

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, false, true)
	m := <-pubMessage
	msg := m.(float64)
	assert.Equal(float64(1), msg)
}

func TestSubscribePublishUnsubscribePNOtherDisablePost(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := map[string]interface{}{
		"id":        2,
		"not_other": "123456",
		"pn_other":  "yay!",
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, true, true)
	m := <-pubMessage
	msg := m.(map[string]interface{})
	assert.Equal("123456", msg["not_other"])
	assert.Equal("yay!", msg["pn_other"])
}

func TestSubscribePublishUnsubscribePNOtherPost(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := map[string]interface{}{
		"id":        1,
		"not_other": "12345",
		"pn_other":  "yay!",
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, false, true)
	m := <-pubMessage
	msg := m.(map[string]interface{})
	assert.Equal("12345", msg["not_other"])
	assert.Equal("yay!", msg["pn_other"])

}

func TestSubscribePublishUnsubscribePNOtherComplexPost(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s1 := map[string]interface{}{
		"id":        1,
		"not_other": "1234567",
	}
	s := map[string]interface{}{
		"id":        1,
		"not_other": "12345",
		"pn_other":  s1,
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, false, true)
	m := <-pubMessage
	msg := m.(map[string]interface{})
	assert.Equal("12345", msg["not_other"])
	if msgOther, ok := msg["pn_other"].(map[string]interface{}); !ok {
		assert.Fail("!map[string]interface{}")
	} else {
		assert.Equal("1234567", msgOther["not_other"])
	}

}

func TestSubscribePublishUnsubscribeInterfaceWithoutPNOtherPost(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := map[string]interface{}{
		"id":        3,
		"not_other": "1234",
		"ss_other":  "yay!",
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "", pubMessage, false, true)
	m := <-pubMessage
	msg := m.(map[string]interface{})
	assert.Equal("1234", msg["not_other"])
	assert.Equal("yay!", msg["ss_other"])

}

func TestSubscribePublishUnsubscribeInterfaceWithoutPNOtherEncPost(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := map[string]interface{}{
		"id":        4,
		"not_other": "123",
		"ss_other":  "yay!",
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, false, true)
	m := <-pubMessage
	msg := m.(map[string]interface{})
	assert.Equal("123", msg["not_other"])
	assert.Equal("yay!", msg["ss_other"])
}

func TestSubscribePublishUnsubscribeInterfaceCustomPost(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := customStruct{
		Foo: "hi!",
		Bar: []int{1, 2, 3, 4, 5},
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "", pubMessage, false, true)
	m := <-pubMessage
	//s1 := reflect.ValueOf(m)
	//fmt.Println("s:::", s1, s1.Type())
	if msg, ok := m.(map[string]interface{}); !ok {
		//fmt.Println(msg)
		assert.Fail("not map")
	} else {
		//fmt.Println(msg)
		//byt := []byte(message.Message)
		//fmt.Println(message.Message.(string))
		//err := json.Unmarshal(byt, &msg)
		//assert.Nil(err)
		assert.Equal("hi!", msg["Foo"])
		//assert.Equal("1", msg["Bar"].(map[string]interface{})[0])
		//assert.Equal("\"yay!\"", msg["pn_other"])
	}
}

func TestSubscribePublishUnsubscribeInterfaceWithoutCustomEncPost(t *testing.T) {
	assert := assert.New(t)
	pubMessage := make(chan interface{})
	s := customStruct{
		Foo: "hi!",
		Bar: []int{1, 2, 3, 4, 5},
	}

	go SubscribePublishUnsubscribeMultiCommon(t, s, "enigma", pubMessage, false, true)
	m := <-pubMessage
	//s1 := reflect.ValueOf(m)
	//fmt.Println("s:::", s1, s1.Type())
	if msg, ok := m.(map[string]interface{}); !ok {
		//fmt.Println(msg)
		assert.Fail("not map")
	} else {
		//fmt.Println(msg)
		//byt := []byte(message.Message)
		//fmt.Println(message.Message.(string))
		//err := json.Unmarshal(byt, &msg)
		//assert.Nil(err)
		assert.Equal("hi!", msg["Foo"])
		//assert.Equal("1", msg["Bar"].(map[string]interface{})[0])
		//assert.Equal("\"yay!\"", msg["pn_other"])
	}
}

func SubscribePublishUnsubscribeMultiCommon(t *testing.T, s interface{}, cipher string, pubMessage chan interface{}, disablePNOtherProcessing bool, usePost bool) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)
	donePublish := make(chan bool)
	exit := make(chan bool)
	errChan := make(chan string)
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6062", nil))
	// }()

	//r := GenRandom()

	ch := "testChannel_sub_96112" //fmt.Sprintf("testChannel_sub_%d", r.Intn(99999))

	pn := pubnub.NewPubNub(configCopy())
	ips, err1 := net.LookupIP(pn.Config.Origin)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err1)
		os.Exit(1)
	}
	for _, ip := range ips {
		fmt.Printf("%s IN A %s\n", pn.Config.Origin, ip.String())
	}

	pn.Config.CipherKey = cipher
	pn.Config.DisablePNOtherProcessing = disablePNOtherProcessing
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	listener := pubnub.NewListener()

	tic := time.NewTicker(time.Duration(timeout) * time.Second)

	go func() {
	CloseLoop:
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
				case pubnub.PNAcknowledgmentCategory:
					doneUnsubscribe <- true
				default:
					fmt.Println("SubscribePublishUnsubscribeMultiCommon status", status)
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				donePublish <- true
				if pubMessage != nil {
					pubMessage <- message.Message
				} else {
					fmt.Println("pubMessage nil")
				}

			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			case <-tic.C:
				fmt.Println("SubscribePublishUnsubscribeMultiCommon timeout")
				assert.Fail("timeout")
				errChan <- "timeout"

				break CloseLoop
			case <-exit:
				tic.Stop()
				break CloseLoop
			}
		}
		//fmt.Println("SubscribePublishUnsubscribeMultiCommon exiting loop")
	}()

	pn.AddListener(listener)

	pn.Subscribe().Channels([]string{ch}).Execute()

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
		//return
	}

	pn.Publish().Channel(ch).Message(s).UsePost(usePost).Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
		//return
	}

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	fmt.Println("SubscribePublishUnsubscribeMultiCommon before doneUnsubscribe")
	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}
	fmt.Println("SubscribePublishUnsubscribeMultiCommon after doneUnsubscribe")
	exit <- true
	fmt.Println("SubscribePublishUnsubscribeMultiCommon after exit")

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
	fmt.Println("SubscribePublishUnsubscribeMultiCommon after zero")

}

/*func TestSubscribePublishUnsubscribePNOther(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)
	donePublish := make(chan bool)
	errChan := make(chan string)

	//r := GenRandom()

	ch := "testChannel_sub_96112" //fmt.Sprintf("testChannel_sub_%d", r.Intn(99999))

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.CipherKey = "enigma"
	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	s := map[string]interface{}{
		"id":        1,
		"not_other": "12345",
		"pn_other":  "yay!",
	}
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				msg := message.Message.(map[string]interface{})
				assert.Equal("12345", msg["not_other"])
				assert.Equal("\"yay!\"", msg["pn_other"])
				donePublish <- true
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().Channels([]string{ch}).Execute()

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Publish().Channel(ch).Message(s).Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
}*/

/*func TestSubscribePublishUnsubscribePNOtherDisable(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)
	donePublish := make(chan bool)
	errChan := make(chan string)

	//r := GenRandom()

	ch := "testChannel_sub_96112" //fmt.Sprintf("testChannel_sub_%d", r.Intn(99999))

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.CipherKey = "enigma"
	pn.Config.DisablePNOtherProcessing = true
	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	s := map[string]interface{}{
		"id":        2,
		"not_other": "1234",
		"pn_other":  "\"yay!\"",
	}
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				//var msg map[string]interface{}
				fmt.Println("reflect.TypeOf(data).Kind()", reflect.TypeOf(message.Message).Kind(), message.Message)
				if msg, ok := message.Message.(map[string]interface{}); !ok {
					fmt.Println(msg)
					assert.Fail("not map")
				} else {
					fmt.Println(msg)
					//byt := []byte(message.Message)
					//fmt.Println(message.Message.(string))
					//err := json.Unmarshal(byt, &msg)
					//assert.Nil(err)
					assert.Equal("1234", msg["not_other"])
					assert.Equal("\"yay!\"", msg["pn_other"])
				}
				donePublish <- true
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().Channels([]string{ch}).Execute()

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Publish().Channel(ch).Message(s).Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
}*/

/*func TestSubscribePublishUnsubscribeInterfaceWithoutPNOther(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)
	donePublish := make(chan bool)
	errChan := make(chan string)

	//r := GenRandom()

	ch := "testChannel_sub_96112" //fmt.Sprintf("testChannel_sub_%d", r.Intn(99999))

	pn := pubnub.NewPubNub(configCopy())

	s := map[string]interface{}{
		"id":        3,
		"not_other": "1234",
		"ss_other":  "\"yay!\"",
	}
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				var msg map[string]interface{}
				fmt.Println("reflect.TypeOf(data).Kind()", reflect.TypeOf(message.Message).Kind(), message.Message)
				msg = message.Message.(map[string]interface{})
				fmt.Println(msg)
				assert.Equal("1234", msg["not_other"])
				assert.Equal("\"yay!\"", msg["ss_other"])
				donePublish <- true
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().Channels([]string{ch}).Execute()

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Publish().Channel(ch).Message(s).Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
}*/

/*func TestSubscribePublishUnsubscribeInterfaceWithoutPNOtherEnc(t *testing.T) {
assert := assert.New(t)
doneSubscribe := make(chan bool)
doneUnsubscribe := make(chan bool)
donePublish := make(chan bool)
errChan := make(chan string)

//r := GenRandom()

ch := "testChannel_sub_96112" //fmt.Sprintf("testChannel_sub_%d", r.Intn(99999))

pn := pubnub.NewPubNub(configCopy())
pn.Config.CipherKey = "enigma"
pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

/*s := map[string]interface{}{
	"not_other": "1234",
	"ss_other":  "\"yay!\"",
}*/
//s := 1.1
/*s := customStruct{
		Foo: "hi!",
		Bar: []int{1, 2, 3, 4, 5},
	}

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				fmt.Println("reflect.TypeOf(data).Kind()", reflect.TypeOf(message.Message).Kind(), message.Message)
				s := reflect.ValueOf(message.Message)
				fmt.Println("s:::", s, s.Type())
				if msg, ok := message.Message.(map[string]interface{}); !ok {
					fmt.Println(msg)
					assert.Fail("not map")
				} else {
					fmt.Println(msg)
					//byt := []byte(message.Message)
					//fmt.Println(message.Message.(string))
					//err := json.Unmarshal(byt, &msg)
					//assert.Nil(err)
					assert.Equal("hi!", msg["Foo"])
					//assert.Equal("1", msg["Bar"].(map[string]interface{})[0])
					//assert.Equal("\"yay!\"", msg["pn_other"])
				}
				donePublish <- true
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().Channels([]string{ch}).Execute()

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Publish().Channel(ch).Message(s).Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
}*/

/*func TestSubscribePublishUnsubscribe(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)
	donePublish := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-pu-ch")

	pn := pubnub.NewPubNub(configCopy())

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				assert.Equal(message.Message, "hey")
				donePublish <- true
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().Channels([]string{ch}).Execute()

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Publish().Channel(ch).Message("hey").Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
}*/

// Also tests:
// - test operations like publish/unsubscribe invoked inside another goroutine
// - test unsubscribe all
func TestSubscribePublishPartialUnsubscribe(t *testing.T) {
	assert := assert.New(t)
	doneUnsubscribe := make(chan bool)
	errChan := make(chan string)
	var once sync.Once

	ch1 := randomized("sub-partialu-ch1")
	ch2 := randomized("sub-partialu-ch2")
	heyPub := heyIterator(3)
	heySub := heyIterator(3)

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.UUID = randomized("sub-partialu-uuid")

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					once.Do(func() {
						pn.Publish().Channel(ch1).Message(<-heyPub).Execute()
					})
					continue
				}
				if len(status.AffectedChannels) == 1 && status.Operation == pubnub.PNUnsubscribeOperation {
					assert.Equal(status.AffectedChannels[0], ch2)
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				if message.Message == <-heySub {
					pn.Unsubscribe().
						Channels([]string{ch2}).
						Execute()
				} else {
					errChan <- fmt.Sprintf("Unexpected message: %s",
						message.Message)
				}
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().Channels([]string{ch1, ch2}).Execute()

	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")
	}

	pn.RemoveListener(listener)
	pn.UnsubscribeAll()

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
}

func JoinLeaveChannel(t *testing.T) {
	assert := assert.New(t)

	// await both connected event on emitter and join presence event received
	var wg sync.WaitGroup
	wg.Add(2)

	donePresenceConnect := make(chan bool)
	doneJoin := make(chan bool)
	doneLeave := make(chan bool)
	errChan := make(chan string)
	ch := randomized("ch")

	configEmitter := configCopy()
	configPresenceListener := configCopy()

	configEmitter.UUID = randomized("sub-lj-emitter")
	configPresenceListener.UUID = randomized("sub-lj-listener")

	pn := pubnub.NewPubNub(configEmitter)
	pnPresenceListener := pubnub.NewPubNub(configPresenceListener)

	listenerEmitter := pubnub.NewListener()
	listenerPresenceListener := pubnub.NewListener()

	// emitter
	go func() {
		for {
			select {
			case status := <-listenerEmitter.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					wg.Done()
					return
				}
			case <-listenerEmitter.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listenerEmitter.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	// listener
	go func() {
		for {
			select {
			case status := <-listenerPresenceListener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					donePresenceConnect <- true
				}
			case message := <-listenerPresenceListener.Message:
				errChan <- fmt.Sprintf("Unexpected message: %s",
					message.Message)
			case presence := <-listenerPresenceListener.Presence:
				// ignore join event of presence listener
				if presence.UUID == configPresenceListener.UUID {
					continue
				}

				assert.Equal(ch, presence.Channel)

				if presence.Event == "leave" {
					assert.Equal(configEmitter.UUID, presence.UUID)
					doneLeave <- true
					return
				}
				assert.Equal("join", presence.Event)
				assert.Equal(configEmitter.UUID, presence.UUID)
				wg.Done()
			}
		}
	}()

	pn.AddListener(listenerEmitter)
	pnPresenceListener.AddListener(listenerPresenceListener)

	pnPresenceListener.Subscribe().
		Channels([]string{ch}).
		WithPresence(true).
		Execute()

	select {
	case <-donePresenceConnect:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Subscribe().
		Channels([]string{ch}).
		Execute()

	go func() {
		wg.Wait()
		doneJoin <- true
	}()

	select {
	case <-doneJoin:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-doneLeave:
	case err := <-errChan:
		assert.Fail(err)
		return
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")

	}
}

func SubscribeJoinLeaveGroup(t *testing.T) {
	assert := assert.New(t)

	// await both connected event on emitter and join presence event received
	var wg sync.WaitGroup
	wg.Add(2)

	donePresenceConnect := make(chan bool)
	doneJoinEvent := make(chan bool)
	doneLeaveEvent := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-jlg-ch")
	cg := randomized("sub-jlg-cg")

	configEmitter := configCopy()
	configPresenceListener := configCopy()

	configEmitter.UUID = randomized("emitter")
	configPresenceListener.UUID = randomized("listener")

	pn := pubnub.NewPubNub(configEmitter)
	pnPresenceListener := pubnub.NewPubNub(configPresenceListener)

	listenerEmitter := pubnub.NewListener()
	listenerPresenceListener := pubnub.NewListener()

	// emitter
	go func() {
		for {
			select {
			case status := <-listenerEmitter.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					wg.Done()
					return
				}
			case <-listenerEmitter.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listenerEmitter.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	// listener
	go func() {
		for {
			select {
			case status := <-listenerPresenceListener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					donePresenceConnect <- true
				}
			case <-listenerPresenceListener.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case presence := <-listenerPresenceListener.Presence:
				// ignore join event of presence listener
				if presence.UUID == configPresenceListener.UUID {
					continue
				}

				assert.Equal(presence.Channel, ch)

				if presence.Event == "leave" {
					assert.Equal(configEmitter.UUID, presence.UUID)
					doneLeaveEvent <- true
					return
				}
				assert.Equal("join", presence.Event)
				assert.Equal(configEmitter.UUID, presence.UUID)
				wg.Done()
			}
		}
	}()

	pn.AddListener(listenerEmitter)
	pnPresenceListener.AddListener(listenerPresenceListener)

	pnPresenceListener.AddChannelToChannelGroup().
		Channels([]string{ch}).
		ChannelGroup(cg).
		Execute()

	pnPresenceListener.Subscribe().
		ChannelGroups([]string{cg}).
		WithPresence(true).
		Execute()

	select {
	case <-donePresenceConnect:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Subscribe().
		ChannelGroups([]string{cg}).
		Execute()

	go func() {
		wg.Wait()
		doneJoinEvent <- true
	}()

	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-doneJoinEvent:
	case err := <-errChan:
		assert.Fail(err)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")

	}

	pn.Unsubscribe().
		ChannelGroups([]string{cg}).
		Execute()

	select {
	case <-doneLeaveEvent:
	case err := <-errChan:
		assert.Fail(err)
	}
}

/////////////////////////////
// Unsubscribe
/////////////////////////////

func TestUnsubscribeAll(t *testing.T) {
	assert := assert.New(t)
	pn := pubnub.NewPubNub(configCopy())
	channels := []string{
		randomized("sub-ua-ch1"),
		randomized("sub-ua-ch2"),
		randomized("sub-ua-ch3")}

	groups := []string{
		randomized("sub-ua-cg1"),
		randomized("sub-ua-cg2"),
		randomized("sub-ua-cg3")}

	pn.Subscribe().
		Channels(channels).
		ChannelGroups(groups).
		WithPresence(true).
		Execute()

	assert.Equal(len(pn.GetSubscribedChannels()), 3)
	assert.Equal(len(pn.GetSubscribedGroups()), 3)

	pn.UnsubscribeAll()

	assert.Equal(len(pn.GetSubscribedChannels()), 0)
	assert.Equal(len(pn.GetSubscribedGroups()), 0)
}

/////////////////////////////
// Misc
/////////////////////////////

func TestSubscribe403Error(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneAccessDenied := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-403-ch")

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	pamConfig := pamConfigCopy()
	pamConfig.SecretKey = ""
	pn2 := pubnub.NewPubNub(pamConfig)

	pn2.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
					break
				case pubnub.PNAccessDeniedCategory:
					doneAccessDenied <- true
					break
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				break
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				break
			}
		}
	}()

	pn2.AddListener(listener)

	pn.Grant().
		Read(false).
		Write(false).
		Manage(false).
		TTL(10).
		Execute()

	fmt.Println("sleeping")
	time.Sleep(5 * time.Second)
	fmt.Println("after sleeping")
	pn2.Subscribe().
		Channels([]string{ch}).
		Execute()

	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-doneSubscribe:
		assert.Fail("Access denied expected")
	case <-doneAccessDenied:
	case err := <-errChan:
		assert.Fail(err)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")

	}
}

func TestSubscribeSignal(t *testing.T) {
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/subscribe/%s/ch/0", config.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"t":{"t":"14858178301085322","r":7},"m":[{"a":"4","e":1,"f":512,"i":"02a7b822-220c-49b0-90c4-d9cbecc0fd85","s":1,"p":{"t":"14858178301075219","r":7},"k":"demo-36","c":"chTest","d":"Signal"}]}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid", "tt"},
		ResponseStatusCode: 200,
	})

	assert := assert.New(t)
	doneMeta := make(chan bool)
	errChan := make(chan string)

	pn := pubnub.NewPubNub(configCopy())
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	pn.SetSubscribeClient(interceptor.GetClient())
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				// ignore status messages
				if status.Error {
					errChan <- fmt.Sprintf("Status Error: %s", status.Category)
					break
				} else {
					//fmt.Println("status", status)
					//doneMeta <- true
					break
				}
			case message := <-listener.Signal:
				meta, ok := message.Message.(string)
				if !ok {
					errChan <- "Invalid message type"
				}
				//fmt.Println("signal", message)
				assert.Equal(meta, "Signal")

				doneMeta <- true
				break
			case message := <-listener.Message:
				meta, ok := message.UserMetadata.(string)
				if !ok {
					errChan <- "Invalid message type"
				}
				fmt.Println("message", message)
				assert.Equal(meta, "my-data")
				doneMeta <- true
				break
			case <-listener.Presence:
				fmt.Println("Presence")
				errChan <- "Got presence while awaiting for a status event"
				break
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()

	select {
	case <-doneMeta:
	case err := <-errChan:
		assert.Fail(err)
	}
}

func TestSubscribeParseUserMeta(t *testing.T) {
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/subscribe/%s/ch/0", config.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"t":{"t":"14858178301085322","r":7},"m":[{"a":"4","f":512,"i":"02a7b822-220c-49b0-90c4-d9cbecc0fd85","s":1,"p":{"t":"14858178301075219","r":7},"k":"demo-36","c":"chTest","u":"my-data","d":{"City":"Goiania","Name":"Marcelo"}}]}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	assert := assert.New(t)
	doneMeta := make(chan bool)
	errChan := make(chan string)

	pn := pubnub.NewPubNub(configCopy())

	pn.SetSubscribeClient(interceptor.GetClient())
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				// ignore status messages
				if status.Error {
					errChan <- fmt.Sprintf("Status Error: %s", status.Category)
					break
				} else {
					fmt.Println(status)
					doneMeta <- true
					break
				}
			case message := <-listener.Message:
				meta, ok := message.UserMetadata.(string)
				if !ok {
					errChan <- "Invalid message type"
				}

				assert.Equal(meta, "my-data")
				doneMeta <- true
				break
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				break
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()

	select {
	case <-doneMeta:
	case err := <-errChan:
		assert.Fail(err)
	}
}

func TestSubscribeWithCustomTimetoken(t *testing.T) {
	ch := "ch"
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/subscribe/%s/ch/0", config.SubscribeKey),
		ResponseBody:       `{"t":{"t":"15069659902324693","r":12},"m":[]}`,
		Query:              "heartbeat=300",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/subscribe/%s", config.SubscribeKey),
		ResponseBody:       fmt.Sprintf(`{"t":{"t":"14607577960932487","r":1},"m":[{"a":"4","f":0,"i":"Client-g5d4g","p":{"t":"14607577960925503","r":1},"k":"%s","c":"ch","d":{"text":"Enter Message Here"},"b":"ch"}]}`, config.SubscribeKey),
		Query:              "heartbeat=300&tt=1337",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
		Hang:               true,
	})

	assert := assert.New(t)
	doneConnected := make(chan bool)
	errChan := make(chan string)

	pn := pubnub.NewPubNub(configCopy())
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	//pn.SetSubscribeClient(interceptor.GetClient())
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				if status.Category == pubnub.PNConnectedCategory {
					doneConnected <- true
					break
				} else {
					errChan <- fmt.Sprintf("Got status while awaiting for a message: %s",
						status.Category)
					break
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a message"
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a message"
				break
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{ch}).
		Timetoken(int64(1337)).
		Execute()

	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-doneConnected:
	case err := <-errChan:
		assert.Fail(err)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")

	}

	pn.UnsubscribeAll()
}

func TestSubscribeWithFilter(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	donePublish := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-wf-ch")

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.FilterExpression = "language!=spanish"
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				}
			case message := <-listener.Message:
				if message.Message == "Hello!" {
					donePublish <- true
				}
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{ch}).
		Execute()

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	pnPublish := pubnub.NewPubNub(configCopy())

	meta := make(map[string]string)
	meta["language"] = "spanish"

	pnPublish.Publish().
		Channel("ch").
		Meta(meta).
		Message("Hola!").
		Execute()

	anotherMeta := make(map[string]string)
	anotherMeta["language"] = "english"

	pnPublish.Publish().
		Channel(ch).
		Meta(anotherMeta).
		Message("Hello!").
		Execute()

	<-donePublish
}

func TestSubscribePublishUnsubscribeWithEncrypt(t *testing.T) {
	assert := assert.New(t)
	doneConnect := make(chan bool)
	donePublish := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-puwe-ch")

	config := configCopy()
	config.CipherKey = "my-key"
	pn := pubnub.NewPubNub(config)
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneConnect <- true
				}
			case message := <-listener.Message:
				assert.Equal("hey", message.Message)
				donePublish <- true
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{ch}).
		Execute()

	select {
	case <-doneConnect:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Publish().
		UsePost(true).
		Channel(ch).
		Message("hey").
		Execute()

	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")

	}
}

func TestSubscribeSuperCall(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	errChan := make(chan string)
	config := pamConfigCopy()
	// Not allowed characters:
	// .,:*
	validCharacters := "-_~?#[]@!$&'()+;=`|"
	config.UUID = validCharacters
	//config.AuthKey = validCharacters
	//config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	pn := pubnub.NewPubNub(config)
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				default:
					errChan <- fmt.Sprintf("Not connected: %v", status)
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	pn.AddListener(listener)

	// Not allowed characters:
	// ?#[]@!$&'()+;=`|
	groupCharacters := "-_~"

	pn.Subscribe().
		Channels([]string{validCharacters + "channel"}).
		ChannelGroups([]string{groupCharacters + "cg"}).
		Timetoken(int64(1337)).
		Execute()

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}
}

func ReconnectionExhaustion(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	errChan := make(chan string)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/subscribe/%s/ch/0", config.SubscribeKey),
		ResponseBody:       "",
		Query:              "heartbeat=300",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 400,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub-key/%s/channel/ch/leave", config.SubscribeKey),
		ResponseBody:       `{"status": 200, "message": "OK", "action": "leave", "service": "Presence"}`,
		Query:              "",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})
	config.MaximumReconnectionRetries = 1
	config.PNReconnectionPolicy = pubnub.PNLinearPolicy

	pn := pubnub.NewPubNub(config)
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	pn.Config.AuthKey = "myAuthKey"
	pn.SetSubscribeClient(interceptor.GetClient())
	listener := pubnub.NewListener()
	count := 0

	go func() {
		for {
			select {
			case status := <-listener.Status:

				switch status.Category {
				case pubnub.PNReconnectionAttemptsExhausted:
					doneSubscribe <- true
				default:
					//if count > 1 {
					//errChan <- fmt.Sprintf("Non PNReconnectedCategory event, %s", status)
					//fmt.Println(status)
					//}
				}
				count++
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()
	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-doneSubscribe:
		fmt.Println("doneSubscribe")
	case err := <-errChan:
		assert.Fail(err)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")

	}
}
