package pubnub

import (
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestHereNowChannelsGroups(t *testing.T) {
	assert := assert.New(t)

	opts := &hereNowOpts{
		Channels:      []string{"ch1", "ch2", "ch3"},
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		pubnub:        pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub_key/sub_key/channel/ch1,ch2,ch3",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHereNowMultipleWithOpts(t *testing.T) {
	assert := assert.New(t)

	opts := &hereNowOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		ChannelGroups:   []string{"cg1", "cg2", "cg3"},
		IncludeUUIDs:    false,
		IncludeState:    true,
		SetIncludeState: true,
		SetIncludeUUIDs: true,
		pubnub:          pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub_key/sub_key/channel/ch1,ch2,ch3",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	expected.Set("state", "1")
	expected.Set("disable-uuids", "1")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHereNowGlobal(t *testing.T) {
	assert := assert.New(t)

	opts := &hereNowOpts{
		pubnub: pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub_key/sub_key",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHereNowValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &hereNowOpts{
		pubnub: pn,
	}

	assert.Equal("pubnub/validation: pubnub: \b: Missing Subscribe Key", opts.validate().Error())
}

func TestHereNowBuildPath(t *testing.T) {
	assert := assert.New(t)
	opts := &hereNowOpts{
		pubnub: pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v2/presence/sub_key/sub_key", path)

}

func TestHereNowBuildQuery(t *testing.T) {
	assert := assert.New(t)
	opts := &hereNowOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		ChannelGroups:   []string{"cg1", "cg2", "cg3"},
		IncludeUUIDs:    false,
		IncludeState:    true,
		SetIncludeState: true,
		SetIncludeUUIDs: false,
		pubnub:          pubnub,
	}
	query, err := opts.buildQuery()
	assert.Nil(err)
	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	expected.Set("state", "1")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

}

func TestNewHereNowResponseErrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newHereNowResponse(jsonBytes, nil, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}
