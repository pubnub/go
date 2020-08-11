package e2e

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	a "github.com/stretchr/testify/assert"
)

func GetTimetoken(pn *pubnub.PubNub) int64 {
	res, _, _ := pn.Time().Execute()
	return res.Timetoken
}

func MatchFetchMessages(ret *pubnub.FetchResponse, start int, ch1, ch2 string, assert *a.Assertions) {
	chMessages := ret.Messages[ch1]
	for i := start; i < len(chMessages); i++ {
		assert.Equal(fmt.Sprintf("testch1 %d", i), chMessages[i].Message)
	}
	ch2Messages := ret.Messages[ch2]
	for i := start; i < len(ch2Messages); i++ {
		assert.Equal(fmt.Sprintf("testch2 %d", i), ch2Messages[i].Message)
	}

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

	r := GenRandom()
	ch1 := fmt.Sprintf("testChannel_sub_%d", r.Intn(99999))
	ch2 := fmt.Sprintf("testChannel_sub_%d", r.Intn(99999))

	timestamp1 := GetTimetoken(pn)
	timestamp2 := int64(0)

	for i := 0; i < 10; i++ {
		if i == 5 {
			timestamp2 = GetTimetoken(pn)
		}
		pn.Publish().Channel(ch1).Message(fmt.Sprintf("testch1 %d", i)).Execute()
		time.Sleep(1 * time.Second)
		pn.Publish().Channel(ch2).Message(fmt.Sprintf("testch2 %d", i)).Execute()
		time.Sleep(1 * time.Second)
	}

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
