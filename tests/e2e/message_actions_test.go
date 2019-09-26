package e2e

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestMessageActionsListeners(t *testing.T) {
	MessageActionsListenersCommon(t, false)
}

func TestMessageActionsListenersEnc(t *testing.T) {
	MessageActionsListenersCommon(t, true)
}

func MessageActionsListenersCommon(t *testing.T, encrypted bool) {
	eventWaitTime := 2
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	r := GenRandom()

	ch := fmt.Sprintf("test_message_actions_%d", r.Intn(99999))
	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	pn.Config.PublishKey = "pub-c-03f156ea-a2e3-4c35-a733-9535824be897"
	pn.Config.SubscribeKey = "sub-c-d7da9e58-c997-11e9-a139-dab2c75acd6f"
	pn.Config.SecretKey = "sec-c-MmUxNTZjMmYtNzFkNS00ODkzLWE2YjctNmQ4YzE5NWNmZDA3"

	pn.Config.Origin = "ingress.bronze.aws-pdx-1.ps.pn"
	pn.Config.Secure = false
	if encrypted {
		pn.Config.CipherKey = "enigma"
	}

	// Subscribe,

	listener := pubnub.NewListener()
	var mut sync.RWMutex

	addEvent := false
	delEvent := false
	var recActionType, recActionTimetoken, recActionValue, recMessageTimetoken string
	doneConnected := make(chan bool)

	go func() {
		for {
			fmt.Println("Running =--->")
			select {

			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneConnected <- true
				default:
					fmt.Println(" --- status: ", status)
				}
			case messageActionsEvent := <-listener.MessageActionsEvent:
				fmt.Println(" --- MessageActionsEvent: ")
				fmt.Println(fmt.Sprintf("%s", messageActionsEvent))
				fmt.Println(fmt.Sprintf("messageActionsEvent.Channel: %s", messageActionsEvent.Channel))
				fmt.Println(fmt.Sprintf("messageActionsEvent.SubscribedChannel: %s", messageActionsEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("messageActionsEvent.Event: %s", messageActionsEvent.Event))
				fmt.Println(fmt.Sprintf("messageActionsEvent.Data.ActionType: %s", messageActionsEvent.Data.ActionType))
				fmt.Println(fmt.Sprintf("messageActionsEvent.Data.ActionValue: %s", messageActionsEvent.Data.ActionValue))
				fmt.Println(fmt.Sprintf("messageActionsEvent.Data.ActionTimetoken: %s", messageActionsEvent.Data.ActionTimetoken))
				fmt.Println(fmt.Sprintf("messageActionsEvent.Data.MessageTimetoken: %s", messageActionsEvent.Data.MessageTimetoken))

				if (messageActionsEvent.Event == pubnub.PNMessageActionsAdded) && (messageActionsEvent.Channel == ch) {
					mut.Lock()
					addEvent = true
					recActionTimetoken = messageActionsEvent.Data.ActionTimetoken
					recActionType = messageActionsEvent.Data.ActionType
					recActionValue = messageActionsEvent.Data.ActionValue
					recMessageTimetoken = messageActionsEvent.Data.MessageTimetoken
					mut.Unlock()
				}
				if (messageActionsEvent.Event == pubnub.PNMessageActionsRemoved) && (messageActionsEvent.Channel == ch) {
					mut.Lock()
					delEvent = true
					mut.Unlock()
				}
			}
			fmt.Println("=>>>>>>>>>>>>> restart")
		}

	}()

	pn.AddListener(listener)

	pn.Subscribe().Channels([]string{ch}).Execute()

	tic := time.NewTicker(time.Duration(eventWaitTime) * time.Second)
	select {
	case <-doneConnected:
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")
	}

	// Publish,

	resPub, _, _ := pn.Publish().
		Channel("ch").
		Message("hey").
		ShouldStore(false).
		Execute()

	var messageTimetoken string

	// read tt,
	if resPub != nil {
		messageTimetoken = strconv.FormatInt(resPub.Timestamp, 10)

		// add action,
		ma := pubnub.MessageAction{
			ActionType:  "reaction",
			ActionValue: "smiley_face",
		}

		resAddMA, _, errAddMA := pn.AddMessageAction().Channel(ch).MessageTimetoken(messageTimetoken).Action(ma).Execute()

		assert.Nil(errAddMA)

		// read add event,
		time.Sleep(1 * time.Second)
		mut.Lock()
		assert.True(addEvent)
		assert.Equal(messageTimetoken, recMessageTimetoken)
		assert.Equal(ma.ActionType, recActionType)
		assert.Equal(ma.ActionValue, recActionValue)

		if resAddMA != nil {
			assert.Equal(recActionTimetoken, resAddMA.Data.ActionTimetoken)
			assert.Equal(ma.ActionType, resAddMA.Data.ActionType)
			assert.Equal(ma.ActionValue, resAddMA.Data.ActionValue)
			assert.Equal(messageTimetoken, resAddMA.Data.MessageTimetoken)
		} else {
			assert.Fail("resAddMA nil")
		}
		mut.Unlock()

		// get action,
		limit := 1

		recActionTimetokenM1 := recActionTimetoken

		n, err := strconv.ParseInt(recActionTimetoken, 10, 64)
		if err == nil {
			n = n + 1
			recActionTimetokenM1 = strconv.FormatInt(n, 10)
		}

		resGetMA1, _, errGetMA1 := pn.GetMessageActions().Channel(ch).Execute()
		assert.Nil(errGetMA1)
		mut.Lock()
		MatchGetMA(1, assert, resGetMA1, recActionType, recActionTimetoken, recActionValue, recMessageTimetoken)
		mut.Unlock()

		resGetMA2, _, errGetMA2 := pn.GetMessageActions().Channel(ch).Start(recActionTimetokenM1).Execute()
		assert.Nil(errGetMA2)
		mut.Lock()
		MatchGetMA(2, assert, resGetMA2, recActionType, recActionTimetoken, recActionValue, recMessageTimetoken)
		mut.Unlock()

		resGetMA3, _, errGetMA3 := pn.GetMessageActions().Channel(ch).Start(recActionTimetokenM1).End(recActionTimetoken).Execute()
		assert.Nil(errGetMA3)
		mut.Lock()
		MatchGetMA(3, assert, resGetMA3, recActionType, recActionTimetoken, recActionValue, recMessageTimetoken)
		mut.Unlock()

		resGetMA4, _, errGetMA4 := pn.GetMessageActions().Channel(ch).Limit(limit).Execute()
		assert.Nil(errGetMA4)
		mut.Lock()
		MatchGetMA(4, assert, resGetMA4, recActionType, recActionTimetoken, recActionValue, recMessageTimetoken)
		mut.Unlock()

		// remove action,
		resDelMA, _, errDelMA := pn.RemoveMessageAction().Channel(ch).MessageTimetoken(messageTimetoken).ActionTimetoken(recActionTimetoken).Execute()
		assert.Nil(errDelMA)
		assert.NotNil(resDelMA)

		// read delete event
		time.Sleep(1 * time.Second)
		mut.Lock()
		assert.True(delEvent)
		mut.Unlock()
	} else {
		assert.Fail("resPub nil")
	}
}

func MatchGetMA(i int, assert *assert.Assertions, res *pubnub.PNGetMessageActionsResponse, recActionType, recActionTimetoken, recActionValue, recMessageTimetoken string) {
	if len(res.Data) > 0 {
		assert.Equal(recActionTimetoken, res.Data[0].ActionTimetoken)
		assert.Equal(recActionType, res.Data[0].ActionType)
		assert.Equal(recActionValue, res.Data[0].ActionValue)
		assert.Equal(recMessageTimetoken, res.Data[0].MessageTimetoken)
	} else {
		assert.Fail(fmt.Sprintf("res.Data < 0, %d", i))
	}
}
