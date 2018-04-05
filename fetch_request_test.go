package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessFetchGet(t *testing.T, expectedString string, channels []string) {
	assert := assert.New(t)

	opts := &fetchOpts{
		Channels:         channels,
		Reverse:          false,
		IncludeTimetoken: true,
		pubnub:           pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/channel/%s", expectedString),
		path, []int{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}

func AssertSuccessFetchQuery(t *testing.T, expectedString string, channels []string) {
	opts := &fetchOpts{
		Channels: channels,
		Reverse:  false,
		pubnub:   pubnub,
	}

	query, _ := opts.buildQuery()

	assert.Equal(t, "25", query.Get("max"))
	assert.Equal(t, "false", query.Get("reverse"))

}

func TestFetchPath(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGet(t, "test1,test2", channels)
}

func TestFetchQuery(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchQuery(t, "%22test%22?max=25&reverse=false", channels)
}

/*func Fetch(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(config)

	channels := []string{"test1", "test2"}

	res, status, err := pn.Fetch().
		Channels(channels).
		Count(count).
		Reverse(reverse).
		Execute()

	assert.Nil(err)

	if status.Error == nil {
		for channel, messages := range res.Messages {
			fmt.Println("channel", channel)
			for _, messageInt := range messages {
				message := pubnub.FetchResponseItem(messageInt)
				fmt.Println(message.Message)
				fmt.Println(message.Timetoken)
			}
		}
	} else {
		fmt.Println("ParseFetch", err)
		fmt.Println("ParseFetch", status.StatusCode)
	}
}*/
