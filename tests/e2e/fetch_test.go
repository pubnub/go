package e2e

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v5"
	a "github.com/stretchr/testify/assert"
)

func GetTimetoken(pn *pubnub.PubNub) int64 {
	res, _, _ := pn.Time().Execute()
	return res.Timetoken
}

func MatchFetchMessages(ret *pubnub.FetchResponse, start int, ch1, ch2 string, assert *a.Assertions) {
	chMessages := ret.Messages[ch1]
	messages1 := make([]string, len(chMessages))
	expectedMessages1 := make([]string, len(chMessages))

	for i := start; i < len(chMessages); i++ {
		messages1[i] = chMessages[i].Message.(string)
		expectedMessages1[i] = "testch1 %d"
	}
	assert.ElementsMatch(expectedMessages1, messages1)

	ch2Messages := ret.Messages[ch2]
	messages2 := make([]string, len(ch2Messages))
	expectedMessages2 := make([]string, len(ch2Messages))
	for i := start; i < len(ch2Messages); i++ {
		messages2[i] = ch2Messages[i].Message.(string)
		expectedMessages2[i] = "testch2 %d"
	}
	assert.ElementsMatch(expectedMessages2, messages2)
}

func TestFetch(t *testing.T) {
	FetchCommon(t, false, false)
}

func TestFetchWithMeta(t *testing.T) {
	FetchCommon(t, true, false)
}

func TestFetchWithMessageActions(t *testing.T) {
	FetchCommon(t, false, true)
}

func TestFetchWithMetaAndMessageActions(t *testing.T) {
	FetchCommon(t, true, true)
}

func FetchCommon(t *testing.T, withMeta, withMessageActions bool) {
	assert := a.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	reverse := true

	ch1 := randomized("testChannel_sub_%d")
	ch2 := ch1 + "_2"

	timestamp1 := GetTimetoken(pn)
	time.Sleep(500 * time.Millisecond)
	timestamp2 := int64(0)

	for i := 0; i < 10; i++ {
		if i == 5 {
			time.Sleep(500 * time.Millisecond)
			timestamp2 = GetTimetoken(pn)
			time.Sleep(500 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond) //to reduce possiblity of quota exceeded
		pn.Publish().Channel(ch1).Message(fmt.Sprintf("testch1 %d", i)).Execute()
		pn.Publish().Channel(ch2).Message(fmt.Sprintf("testch2 %d", i)).Execute()
	}
	time.Sleep(500 * time.Millisecond)
	timestamp3 := GetTimetoken(pn)

	ret, _, err := pn.Fetch().
		Channels([]string{ch1, ch2}).
		Count(25).
		Reverse(reverse).
		Execute()

	assert.Nil(err)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	MatchFetchMessages(ret, 0, ch1, ch2, assert)

	ret1, _, err1 := pn.FetchWithContext(backgroundContext).
		Channels([]string{ch1, ch2}).
		Count(25).
		Reverse(reverse).
		Start(timestamp1).
		End(timestamp2).
		QueryParam(queryParam).
		Execute()

	assert.Nil(err1)

	MatchFetchMessages(ret1, 0, ch1, ch2, assert)

	ret2, _, err2 := pn.Fetch().
		Channels([]string{ch1, ch2}).
		Count(25).
		Reverse(reverse).
		Start(timestamp2).
		End(timestamp3).
		Execute()

	assert.Nil(err2)

	MatchFetchMessages(ret2, 5, ch1, ch2, assert)

	ret3, _, err3 := pn.Fetch().
		Channels([]string{ch1, ch2}).
		Count(25).
		Reverse(reverse).
		Start(timestamp1).
		Execute()

	assert.Nil(err3)

	MatchFetchMessages(ret3, 0, ch1, ch2, assert)

	ret4, _, err4 := pn.Fetch().
		Channels([]string{ch1, ch2}).
		Count(25).
		Reverse(reverse).
		Start(timestamp2).
		Execute()

	assert.Nil(err4)

	MatchFetchMessages(ret4, 5, ch1, ch2, assert)

}
