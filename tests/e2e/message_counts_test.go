package e2e

import (
	"fmt"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v5"
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
		if i < 6 {
			pn.Publish().Channel(ch2).Message(fmt.Sprintf("testch2 %d", i)).Execute()
			time.Sleep(1 * time.Second)
		}

	}

	timestamp3 := GetTimetoken(pn)

	_, _, err0 := pn.MessageCounts().
		Channels([]string{ch1, ch2}).
		ChannelsTimetoken([]int64{timestamp1, timestamp2, timestamp3}).
		Execute()
	assert.Contains(err0.Error(), pubnub.StrChannelsTimetokenLength)

	ret, _, err := pn.MessageCounts().
		Channels([]string{ch1, ch2}).
		ChannelsTimetoken([]int64{timestamp2, timestamp3}).
		Execute()

	assert.Nil(err)
	MatchMessageCounts(ret, 5, 0, ch1, ch2, assert)

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	ret3, _, err := pn.MessageCounts().
		Channels([]string{ch1, ch2}).
		Timetoken(timestamp1).
		QueryParam(queryParam).
		Execute()

	MatchMessageCounts(ret3, 10, 6, ch1, ch2, assert)
	assert.Nil(err)

	ret1, _, err1 := pn.MessageCountsWithContext(backgroundContext).
		Channels([]string{ch1, ch2}).
		Timetoken(timestamp2).
		QueryParam(queryParam).
		Execute()

	assert.Nil(err1)

	MatchMessageCounts(ret1, 5, 1, ch1, ch2, assert)

	ret4, _, err4 := pn.MessageCountsWithContext(backgroundContext).
		Channels([]string{ch1, ch2}).
		ChannelsTimetoken([]int64{timestamp2}).
		QueryParam(queryParam).
		Execute()

	assert.Nil(err4)

	MatchMessageCounts(ret4, 5, 1, ch1, ch2, assert)

}
