package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v5/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessMessageCountsGet(t *testing.T, expectedString string, channels []string, timetoken int64, channelsTimetoken []int64) {
	assert := assert.New(t)

	opts := &messageCountsOpts{
		Channels:          channels,
		Timetoken:         timetoken,
		ChannelsTimetoken: channelsTimetoken,
		pubnub:            pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/message-counts/%s", expectedString),
		path, []int{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Empty(body)
}

func TestMessageCountsPath(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}
	AssertSuccessMessageCountsGet(t, "test1,test2", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsQuery(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}
	AssertSuccessMessageCountsGetQuery(t, "", "15499825804610610,15499925804610615", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsQuery2(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{}
	AssertSuccessMessageCountsGetQuery(t, "", "", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsQuery3(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}
	AssertSuccessMessageCountsGetQuery(t, "", "15499825804610610,15499925804610615", channels, 0, channelsTimetoken)
}

func AssertSuccessMessageCountsGetQuery(t *testing.T, expectedString1 string, expectedString2 string, channels []string, timetoken int64, channelsTimetoken []int64) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := &messageCountsOpts{
		Channels:          channels,
		Timetoken:         timetoken,
		ChannelsTimetoken: channelsTimetoken,
		pubnub:            pubnub,
		QueryParam:        queryParam,
	}

	u, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))

	assert.Equal(expectedString1, u.Get("timetoken"))
	assert.Equal(expectedString2, u.Get("channelsTimetoken"))

}

func AssertNewMessageCountsBuilder(t *testing.T, testQueryParam bool, testContext bool, expectedString string, expectedString1 string, expectedString2 string, channels []string, timetoken int64, channelsTimetoken []int64) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	o := newMessageCountsBuilder(pubnub)
	if testContext {
		o = newMessageCountsBuilderWithContext(pubnub, backgroundContext)
	}
	o.Channels(channels)
	o.Timetoken(timetoken)
	o.ChannelsTimetoken(channelsTimetoken)
	if testQueryParam {
		o.QueryParam(queryParam)
	}

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/message-counts/%s", expectedString),
		path, []int{})

	u, _ := o.opts.buildQuery()

	if testQueryParam {
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

	assert.Equal(expectedString1, u.Get("timetoken"))
	assert.Equal(expectedString2, u.Get("channelsTimetoken"))

}

func TestMessageCountsBuilder(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}
	AssertNewMessageCountsBuilder(t, false, false, "test1,test2", "", "15499825804610610,15499925804610615", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsBuilderQP(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}
	AssertNewMessageCountsBuilder(t, true, false, "test1,test2", "", "15499825804610610,15499925804610615", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsBuilderContext(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}
	AssertNewMessageCountsBuilder(t, false, true, "test1,test2", "", "15499825804610610,15499925804610615", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsBuilderContextQP(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}

	AssertNewMessageCountsBuilder(t, true, true, "test1,test2", "", "15499825804610610,15499925804610615", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &messageCountsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newMessageCountsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

//{"status": 200, "error": false, "error_message": "", "channels": {"my-channel1":1,"my-channel":2}}
func TestMessageCountsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &messageCountsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status": 200, "error": false, "error_message": "", "channels": {"my-channel1":1,"my-channel":2}}`)

	res, _, err := newMessageCountsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(2, res.Channels["my-channel"])
	assert.Equal(1, res.Channels["my-channel1"])
	assert.Nil(err)
}
