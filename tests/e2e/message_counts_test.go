package e2e

import (
	"fmt"
	//"log"
	//"os"
	"strconv"
	"testing"

	pubnub "github.com/pubnub/go"
	a "github.com/stretchr/testify/assert"
)

func MatchMessageCounts(ret *pubnub.MessageCountsResponse, count1, count2 int, ch1, ch2 string, assert *a.Assertions) {
	for ch, v := range ret.Channels {
		if ch == ch1 {
			assert.Equal(count1, v)
		}
		if ch == ch2 {
			assert.Equal(count2, v)
		}

	}
}

func TestMessageCounts(t *testing.T) {
	assert := a.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

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
		if i < 6 {
			pn.Publish().Channel(ch2).Message(fmt.Sprintf("testch2 %d", i)).Execute()
		}
	}

	timestamp3 := GetTimetoken(pn)
	fmt.Println("here", strconv.FormatInt(timestamp2, 10), strconv.FormatInt(timestamp3, 10))

	ret, s, err := pn.MessageCounts().
		Channels([]string{ch1, ch2}).
		ChannelsTimetoken([]string{strconv.FormatInt(timestamp2, 10), strconv.FormatInt(timestamp3, 10)}).
		Execute()

	fmt.Println("s", s)
	fmt.Println("s.StatusCode", s.StatusCode)

	assert.Nil(err)
	MatchMessageCounts(ret, 5, 0, ch1, ch2, assert)

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	ret3, _, err := pn.MessageCounts().
		Channels([]string{ch1, ch2}).
		Timetoken(strconv.FormatInt(timestamp1, 10)).
		QueryParam(queryParam).
		Execute()

	MatchMessageCounts(ret3, 10, 6, ch1, ch2, assert)
	assert.Nil(err)

	ret1, _, err1 := pn.MessageCountsWithContext(backgroundContext).
		Channels([]string{ch1, ch2}).
		Timetoken(strconv.FormatInt(timestamp2, 10)).
		QueryParam(queryParam).
		Execute()

	assert.Nil(err1)

	MatchMessageCounts(ret1, 5, 1, ch1, ch2, assert)

}
