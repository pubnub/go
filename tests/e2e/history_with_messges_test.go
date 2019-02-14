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

func MatchHistoryWithMessages(ret *pubnub.HistoryWithMessagesResponse, count1, count2 int, ch1, ch2 string, assert *a.Assertions) {
	for ch, v := range ret.Channels {
		if ch == ch1 {
			assert.Equal(count1, v)
		}
		if ch == ch2 {
			assert.Equal(count2, v)
		}

	}
}

func TestHistoryWithMessages(t *testing.T) {
	assert := a.New(t)

	config2 := pubnub.NewDemoConfig()
	config2.Secure = false
	config2.Origin = "balancer1g.bronze.aws-pdx-1.ps.pn"

	pn := pubnub.NewPubNub(config2)
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
		if i < 9 {
			pn.Publish().Channel(ch2).Message(fmt.Sprintf("testch2 %d", i)).Execute()
		}
	}

	timestamp3 := GetTimetoken(pn)
	fmt.Println("here", strconv.FormatInt(timestamp2, 10), strconv.FormatInt(timestamp3, 10))

	ret, s, err := pn.HistoryWithMessages().
		Channels([]string{ch1, ch2}).
		ChannelsTimetoken([]string{strconv.FormatInt(timestamp2, 10), strconv.FormatInt(timestamp3, 10)}).
		Execute()

	fmt.Println("s", s)
	fmt.Println("s.StatusCode", s.StatusCode)

	assert.Nil(err)
	MatchHistoryWithMessages(ret, 5, 0, ch1, ch2, assert)

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	ret3, _, err := pn.HistoryWithMessages().
		Channels([]string{ch1, ch2}).
		Timetoken(strconv.FormatInt(timestamp1, 10)).
		QueryParam(queryParam).
		Execute()

	MatchHistoryWithMessages(ret3, 10, 9, ch1, ch2, assert)
	assert.Nil(err)

	ret1, _, err1 := pn.HistoryWithMessagesWithContext(backgroundContext).
		Channels([]string{ch1, ch2}).
		Timetoken(strconv.FormatInt(timestamp2, 10)).
		QueryParam(queryParam).
		Execute()

	assert.Nil(err1)

	MatchHistoryWithMessages(ret1, 5, 4, ch1, ch2, assert)

}
