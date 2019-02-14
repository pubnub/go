package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessHistoryWithMessagesGet(t *testing.T, expectedString string, channels []string, timetoken string, channelTimetokens []string) {
	assert := assert.New(t)

	opts := &historyWithMessagesOpts{
		Channels:          channels,
		Timetoken:         timetoken,
		ChannelTimetokens: channelTimetokens,
		pubnub:            pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/channels-with-messages/%s", expectedString),
		path, []int{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Empty(body)
}

func TestHistoryWithMessagesPath(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelTimetokens := []string{"15499825804610610", "15499925804610615"}
	AssertSuccessHistoryWithMessagesGet(t, "test1,test2", channels, "15499825804610610", channelTimetokens)
}

func TestHistoryWithMessagesQuery(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelTimetokens := []string{"15499825804610610", "15499925804610615"}
	AssertSuccessHistoryWithMessagesGetQuery(t, "15499825804610610", "15499825804610610,15499925804610615", channels, "15499825804610610", channelTimetokens)
}

func TestHistoryWithMessagesQuery2(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelTimetokens := []string{}
	AssertSuccessHistoryWithMessagesGetQuery(t, "15499825804610610", "", channels, "15499825804610610", channelTimetokens)
}

func TestHistoryWithMessagesQuery3(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelTimetokens := []string{"15499825804610610", "15499925804610615"}
	AssertSuccessHistoryWithMessagesGetQuery(t, "", "15499825804610610,15499925804610615", channels, "", channelTimetokens)
}

func AssertSuccessHistoryWithMessagesGetQuery(t *testing.T, expectedString1 string, expectedString2 string, channels []string, timetoken string, channelTimetokens []string) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := &historyWithMessagesOpts{
		Channels:          channels,
		Timetoken:         timetoken,
		ChannelTimetokens: channelTimetokens,
		pubnub:            pubnub,
		QueryParam:        queryParam,
	}

	u, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))

	assert.Equal(expectedString1, u.Get("timetoken"))
	assert.Equal(expectedString2, u.Get("channelTimetokens"))

}

func AssertNewHistoryWithMessagesBuilder(t *testing.T, testQueryParam bool, testContext bool, expectedString string, expectedString1 string, expectedString2 string, channels []string, timetoken string, channelTimetokens []string) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	o := newHistoryWithMessagesBuilder(pubnub)
	if testContext {
		o = newHistoryWithMessagesBuilderWithContext(pubnub, backgroundContext)
	}
	o.Channels(channels)
	o.Timetoken(timetoken)
	o.ChannelTimetokens(channelTimetokens)
	if testQueryParam {
		o.QueryParam(queryParam)
	}

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/channels-with-messages/%s", expectedString),
		path, []int{})

	u, _ := o.opts.buildQuery()

	if testQueryParam {
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

	assert.Equal(expectedString1, u.Get("timetoken"))
	assert.Equal(expectedString2, u.Get("channelTimetokens"))

}

func TestHistoryWithMessagesBuilder(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelTimetokens := []string{"15499825804610610", "15499925804610615"}
	AssertNewHistoryWithMessagesBuilder(t, false, false, "test1,test2", "15499825804610610", "15499825804610610,15499925804610615", channels, "15499825804610610", channelTimetokens)
}

func TestHistoryWithMessagesBuilderQP(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelTimetokens := []string{"15499825804610610", "15499925804610615"}
	AssertNewHistoryWithMessagesBuilder(t, true, false, "test1,test2", "15499825804610610", "15499825804610610,15499925804610615", channels, "15499825804610610", channelTimetokens)
}

func TestHistoryWithMessagesBuilderContext(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelTimetokens := []string{"15499825804610610", "15499925804610615"}
	AssertNewHistoryWithMessagesBuilder(t, false, true, "test1,test2", "15499825804610610", "15499825804610610,15499925804610615", channels, "15499825804610610", channelTimetokens)
}

func TestHistoryWithMessagesBuilderContextQP(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelTimetokens := []string{"15499825804610610", "15499925804610615"}
	AssertNewHistoryWithMessagesBuilder(t, true, true, "test1,test2", "15499825804610610", "15499825804610610,15499925804610615", channels, "15499825804610610", channelTimetokens)
}

func TestHistoryWithMessagesResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &historyWithMessagesOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newHistoryWithMessagesResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

//{"status": 200, "error": false, "error_message": "", "channels": {"my-channel1":1,"my-channel":2}}
func TestHistoryWithMessagesResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &historyWithMessagesOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status": 200, "error": false, "error_message": "", "channels": {"my-channel1":1,"my-channel":2}}`)

	res, _, err := newHistoryWithMessagesResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(2, res.Channels["my-channel"])
	assert.Equal(1, res.Channels["my-channel1"])
	assert.Nil(err)
}
