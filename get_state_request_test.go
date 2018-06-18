package pubnub

import (
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestNewGetStateResponse(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.UUID = "my-custom-uuid"

	//https://ssp.pubnub.com/v2/presence/sub-key/s/channel/my-channel/uuid/pn-696b6ccf-b473-4b4e-b86e-02ce7eca68cb?pnsdk=PubNub-Go/4.0.0-beta.7&uuid=pn-696b6ccf-b473-4b4e-b86e-02ce7eca68cb

	jsonBytes := []byte(`{"status": 200, "message": "OK", "payload": {"k": "v"}, "uuid": "my-custom-uuid", "channel": "my-channel", "service": "Presence"}`)

	res, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Nil(err)

	if s, ok := res.State["my-channel"].(map[string]interface{}); ok {
		assert.Equal("v", s["k"])
	} else {
		assert.Fail("!map[string]interface{}")
	}
}

func TestNewGetStateResponse2(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.UUID = "my-custom-uuid"
	//https://ps.pubnub.com/v2/presence/sub-key/s/channel/my-channel3,my-channel2,my-channel/uuid/5fef96e6-a64b-4808-8712-3623af768c3b?pnsdk=PubNub-Go/4.0.0-beta.7&uuid=5fef96e6-a64b-4808-8712-3623af768c3b

	jsonBytes := []byte(`{"status": 200, "message": "OK", "payload": {"channels": {"my-channel3": {"k": "v4"}, "my-channel2": {"k": "v3"}, "my-channel": {"k": "v3"}}}, "uuid": "my-custom-uuid", "service": "Presence"}`)

	res, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Nil(err)
	if s, ok := res.State["my-channel"].(map[string]interface{}); ok {
		assert.Equal("v3", s["k"])
	} else {
		assert.Fail("!map[string]interface{}")
	}
	if s, ok := res.State["my-channel3"].(map[string]interface{}); ok {
		assert.Equal("v4", s["k"])
	} else {
		assert.Fail("!map[string]interface{}")
	}
	if s, ok := res.State["my-channel2"].(map[string]interface{}); ok {
		assert.Equal("v3", s["k"])
	} else {
		assert.Fail("!map[string]interface{}")
	}
}

func TestNewGetStateResponseErr(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.UUID = "my-custom-uuid"

	jsonBytes := []byte(`{"status": 400, "error": 1, "message": "Invalid JSON specified.", "service": "Presence"}`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("Invalid JSON specified.", err.Error())
}

func TestGetStateBasicRequest(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.UUID = "my-custom-uuid"

	opts := &getStateOpts{
		Channels:      []string{"ch"},
		ChannelGroups: []string{"cg"},
		pubnub:        pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/ch/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestGetStateMultipleChannelsChannelGroups(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.UUID = "my-custom-uuid"

	opts := &getStateOpts{
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
		"/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/uuid/my-custom-uuid",
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
