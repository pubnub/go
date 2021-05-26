package e2e

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v5"
	"github.com/stretchr/testify/assert"
)

func TestMessageActionsListeners(t *testing.T) {
	MessageActionsListenersCommon(t, false, false, false)
}

func TestMessageActionsListenersEnc(t *testing.T) {
	MessageActionsListenersCommon(t, true, false, false)
}

func TestMessageActionsListenersWithMeta(t *testing.T) {
	MessageActionsListenersCommon(t, false, true, false)
}

func TestMessageActionsListenersEncWithMeta(t *testing.T) {
	MessageActionsListenersCommon(t, true, true, false)
}

func TestMessageActionsListenersWithMA(t *testing.T) {
	MessageActionsListenersCommon(t, false, false, true)
}

func TestMessageActionsListenersEncWithMA(t *testing.T) {
	MessageActionsListenersCommon(t, true, false, true)
}

func TestMessageActionsListenersWithMetaMA(t *testing.T) {
	MessageActionsListenersCommon(t, false, true, true)
}

func TestMessageActionsListenersEncWithMetaMA(t *testing.T) {
	MessageActionsListenersCommon(t, true, true, true)
}

func MessageActionsListenersCommon(t *testing.T, encrypted, withMeta, withMessageActions bool) {
	eventWaitTime := 2
	assert := assert.New(t)

	pnMA := pubnub.NewPubNub(configCopy())

	r := GenRandom()

	chMA := fmt.Sprintf("test_message_actions_%d", r.Intn(99999))
	if enableDebuggingInTests {
		pnMA.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	if encrypted {
		pnMA.Config.CipherKey = "enigma"
	}

	// Subscribe,

	listener := pubnub.NewListener()
	var mut sync.RWMutex

	addEvent := false
	delEvent := false
	var recActionType, recActionTimetoken, recActionValue, recMessageTimetoken string
	doneConnected := make(chan bool)
	exitListener := make(chan bool)

	go func() {
	ExitLabel:
		for {
			select {

			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneConnected <- true
				default:
					//fmt.Println(" --- status: ", status)
				}
			case messageActionsEvent := <-listener.MessageActionsEvent:
				if enableDebuggingInTests {
					fmt.Println(" --- MessageActionsEvent: ")
					fmt.Println(fmt.Sprintf("%s", messageActionsEvent))
					fmt.Println(fmt.Sprintf("messageActionsEvent.Channel: %s", messageActionsEvent.Channel))
					fmt.Println(fmt.Sprintf("messageActionsEvent.SubscribedChannel: %s", messageActionsEvent.SubscribedChannel))
					fmt.Println(fmt.Sprintf("messageActionsEvent.Event: %s", messageActionsEvent.Event))
					fmt.Println(fmt.Sprintf("messageActionsEvent.Data.ActionType: %s", messageActionsEvent.Data.ActionType))
					fmt.Println(fmt.Sprintf("messageActionsEvent.Data.ActionValue: %s", messageActionsEvent.Data.ActionValue))
					fmt.Println(fmt.Sprintf("messageActionsEvent.Data.ActionTimetoken: %s", messageActionsEvent.Data.ActionTimetoken))
					fmt.Println(fmt.Sprintf("messageActionsEvent.Data.MessageTimetoken: %s", messageActionsEvent.Data.MessageTimetoken))
					fmt.Println(fmt.Sprintf("messageActionsEvent.Data.UUID: %s", messageActionsEvent.Data.UUID))
				}

				if (messageActionsEvent.Event == pubnub.PNMessageActionsAdded) && (messageActionsEvent.Channel == chMA) {
					mut.Lock()
					addEvent = true
					recActionTimetoken = messageActionsEvent.Data.ActionTimetoken
					recActionType = messageActionsEvent.Data.ActionType
					recActionValue = messageActionsEvent.Data.ActionValue
					recMessageTimetoken = messageActionsEvent.Data.MessageTimetoken
					mut.Unlock()
				}
				if (messageActionsEvent.Event == pubnub.PNMessageActionsRemoved) && (messageActionsEvent.Channel == chMA) {
					mut.Lock()
					delEvent = true
					mut.Unlock()
				}
			case <-exitListener:
				break ExitLabel

			}
		}

	}()

	pnMA.AddListener(listener)

	pnMA.Subscribe().Channels([]string{chMA}).Execute()

	tic := time.NewTicker(time.Duration(eventWaitTime) * time.Second)
	select {
	case <-doneConnected:
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")
	}

	meta := map[string]string{
		"m1": "n1",
		"m2": "n2",
	}
	if !withMeta {
		meta = nil
	}

	// Publish,
	message := fmt.Sprintf("hey_%d", r.Intn(99999))

	resPub, _, _ := pnMA.Publish().
		Channel(chMA).
		Message(message).
		ShouldStore(true).
		Meta(meta).
		Execute()

	var messageTimetoken string

	// read tt,
	if resPub != nil {
		messageTimetoken = strconv.FormatInt(resPub.Timestamp, 10)
		//fmt.Println("messageTimetoken", messageTimetoken)

		// add action,
		ma := pubnub.MessageAction{
			ActionType:  "reaction",
			ActionValue: "smiley_face",
		}

		resAddMA, _, errAddMA := pnMA.AddMessageAction().Channel(chMA).MessageTimetoken(messageTimetoken).Action(ma).Execute()

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
		//fmt.Println("recActionTimetoken", recActionTimetoken)

		n, err := strconv.ParseInt(recActionTimetoken, 10, 64)
		if err == nil {
			n = n + 1
			recActionTimetokenM1 = strconv.FormatInt(n, 10)
		}
		//fmt.Println("recActionTimetokenM1", recActionTimetokenM1, limit)

		resGetMA1, _, errGetMA1 := pnMA.GetMessageActions().Channel(chMA).Execute()
		assert.Nil(errGetMA1)
		mut.Lock()
		MatchGetMA(1, assert, resGetMA1, recActionType, recActionTimetoken, recActionValue, recMessageTimetoken)
		mut.Unlock()

		resGetMA2, _, errGetMA2 := pnMA.GetMessageActions().Channel(chMA).Start(recActionTimetokenM1).Execute()
		assert.Nil(errGetMA2)
		mut.Lock()
		MatchGetMA(2, assert, resGetMA2, recActionType, recActionTimetoken, recActionValue, recMessageTimetoken)
		mut.Unlock()

		resGetMA3, _, errGetMA3 := pnMA.GetMessageActions().Channel(chMA).Start(recActionTimetokenM1).End(recActionTimetoken).Execute()
		assert.Nil(errGetMA3)
		mut.Lock()
		MatchGetMA(3, assert, resGetMA3, recActionType, recActionTimetoken, recActionValue, recMessageTimetoken)
		mut.Unlock()

		resGetMA4, _, errGetMA4 := pnMA.GetMessageActions().Channel(chMA).Limit(limit).Execute()
		assert.Nil(errGetMA4)
		mut.Lock()
		MatchGetMA(4, assert, resGetMA4, recActionType, recActionTimetoken, recActionValue, recMessageTimetoken)
		mut.Unlock()

		var att int64
		tt, err := strconv.ParseInt(recActionTimetoken, 10, 64)
		if err == nil {
			att = int64(tt)
		}

		var mtt int64
		tt1, err := strconv.ParseInt(resGetMA1.Data[0].MessageTimetoken, 10, 64)
		if err == nil {
			mtt = int64(tt1)
		}

		if enableDebuggingInTests {
			fmt.Println("att", att)
			fmt.Println("mtt", mtt)
			fmt.Println("recActionTimetoken", recActionTimetoken)
			fmt.Println("messageTimetoken", messageTimetoken)
			fmt.Println("resPub", resPub)
		}

		// Test History with Meta
		resHist, _, errHist := pnMA.History().
			Channel(chMA).
			Count(10).
			Reverse(true).
			Start(att).
			End(mtt).
			IncludeMeta(withMeta).
			IncludeTimetoken(true).
			Execute()
		assert.Nil(errHist)
		mut.Lock()
		MatchHistoryMessageWithMAResp(assert, resHist, chMA, message, mtt, meta, withMeta)
		mut.Unlock()

		// Test Fetch with Meta
		retFM2, _, errFM2 := pnMA.Fetch().
			Channels([]string{chMA}).
			Count(10).
			Reverse(true).
			Start(att).
			End(mtt).
			IncludeMeta(withMeta).
			Execute()
		assert.Nil(errFM2)
		mut.Lock()
		MatchFetchMessageWithMAResp(assert, retFM2, chMA, message, mtt, att, pnMA.Config.UUID, ma, meta, withMeta, false)
		mut.Unlock()

		// Test Fetch with Meta and Actions
		retFM, _, errFM := pnMA.Fetch().
			Channels([]string{chMA}).
			Count(10).
			Reverse(true).
			Start(att).
			End(mtt).
			IncludeMeta(withMeta).
			IncludeMessageActions(withMessageActions).
			Execute()

		assert.Nil(errFM)
		mut.Lock()
		MatchFetchMessageWithMAResp(assert, retFM, chMA, message, mtt, att, pnMA.Config.UUID, ma, meta, withMeta, withMessageActions)
		mut.Unlock()

		// remove action,
		resDelMA, _, errDelMA := pnMA.RemoveMessageAction().Channel(chMA).MessageTimetoken(messageTimetoken).ActionTimetoken(recActionTimetoken).Execute()
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
	exitListener <- true
}

func MatchHistoryMessageWithMAResp(assert *assert.Assertions, resp *pubnub.HistoryResponse, chMA, message string, messageTimetoken int64, meta interface{}, withMeta bool) {
	if resp != nil {
		messages := resp.Messages
		//fmt.Println("====> history messages:", messages)
		if messages != nil {
			assert.Equal(message, messages[0].Message)
			assert.Equal(messageTimetoken, messages[0].Timetoken)
			if withMeta {
				if meta != nil {
					meta := messages[0].Meta.(map[string]interface{})
					assert.Equal("n1", meta["m1"])
					assert.Equal("n2", meta["m2"])
					//fmt.Println("meta:", meta)
				} else {
					assert.Fail("meta nil")
				}
			}
		} else {
			assert.Fail("messages nil")
		}
	} else {
		assert.Fail("res nil")
	}
}

func MatchFetchMessageWithMAResp(assert *assert.Assertions, resp *pubnub.FetchResponse, chMA, message string, messageTimetoken, recActionTimetokenM1 int64, UUID string, ma pubnub.MessageAction, meta interface{}, withMeta, withMessageActions bool) {
	if resp != nil {
		messages := resp.Messages
		//fmt.Println("messages:", messages)
		m0 := messages[chMA]
		if m0 != nil {
			assert.Equal(message, m0[0].Message)
			assert.Equal(strconv.FormatInt(messageTimetoken, 10), m0[0].Timetoken)
			if withMeta {
				if meta != nil {
					meta := m0[0].Meta.(map[string]interface{})
					assert.Equal("n1", meta["m1"])
					assert.Equal("n2", meta["m2"])
					//fmt.Println("meta:", meta)
				} else {
					assert.Fail("meta nil")
				}
			}
			if withMessageActions {
				actionTypes := m0[0].MessageActions

				if len(actionTypes) > 0 {
					a0 := actionTypes[ma.ActionType]
					r00 := a0.ActionsTypeValues[ma.ActionValue]
					if r00 != nil {
						assert.Equal(UUID, r00[0].UUID)
						assert.Equal(strconv.FormatInt(recActionTimetokenM1, 10), r00[0].ActionTimetoken)
						//fmt.Println("action val:", r00[0].UUID, r00[0].ActionTimetoken)
					} else {
						assert.Fail("r0 nil")
					}
				} else {
					assert.Fail("actionTypes nil")
				}
			}
		} else {
			assert.Fail("m0 nil")
		}

	} else {
		assert.Fail("res nil")
	}
}

func MatchGetMA(i int, assert *assert.Assertions, res *pubnub.PNGetMessageActionsResponse, recActionType, recActionTimetoken, recActionValue, recMessageTimetoken string) {
	if res != nil {
		if len(res.Data) > 0 {
			assert.Equal(recActionTimetoken, res.Data[0].ActionTimetoken)
			assert.Equal(recActionType, res.Data[0].ActionType)
			assert.Equal(recActionValue, res.Data[0].ActionValue)
			assert.Equal(recMessageTimetoken, res.Data[0].MessageTimetoken)
		} else {
			assert.Fail(fmt.Sprintf("res.Data < 0, %d", i))
		}
	} else {
		assert.Fail(fmt.Sprintf("res.Data nil"))
	}
}
