package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestNewGetStateResponse(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.UUID = "my-custom-uuid"

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

func TestGetStateBasicRequestWithUUID(t *testing.T) {
	assert := assert.New(t)

	uuid := "customuuid"

	opts := &getStateOpts{
		Channels:      []string{"ch"},
		ChannelGroups: []string{"cg"},
		UUID:          uuid,
		pubnub:        pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/ch/uuid/%s", uuid),
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

func TestNewGetStateBuilder(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.UUID = "my-custom-uuid"

	o := newGetStateBuilder(pubnub)
	o.Channels([]string{"ch"})
	o.ChannelGroups([]string{"cg"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/ch/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewGetStateBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.UUID = "my-custom-uuid"

	o := newGetStateBuilder(pubnub)
	o.Channels([]string{"ch"})
	o.ChannelGroups([]string{"cg"})
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	path, err := o.opts.buildPath()
	o.opts.QueryParam = queryParam

	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/ch/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Equal("v1", query.Get("q1"))
	assert.Equal("v2", query.Get("q2"))

	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewGetStateBuilderContext(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.UUID = "my-custom-uuid"

	o := newGetStateBuilderWithContext(pubnub, backgroundContext)
	o.Channels([]string{"ch"})
	o.ChannelGroups([]string{"cg"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/ch/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
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

func TestGetStateValidateChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getStateOpts{
		pubnub: pn,
	}
	assert.Equal("pubnub/validation: pubnub: Get State: Missing Channel or Channel Group", opts.validate().Error())
}

func TestGetStateValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &getStateOpts{
		Channels:      []string{"ch1", "ch2", "ch3"},
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		pubnub:        pn,
	}

	assert.Equal("pubnub/validation: pubnub: Get State: Missing Subscribe Key", opts.validate().Error())
}

func TestNewGetStateResponseErrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestNewGetStateResponseParsingError(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`"s"`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("Response parsing error", err.Error())
}

func TestNewGetStateResponseParsingPayloadError(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"status": 200, "message": "OK", "payload": "error", "uuid": "my-custom-uuid", "service": "Presence"}`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("Response parsing payload 2", err.Error())
}

func TestNewGetStateResponseParsingPayloadChannelsError(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"status": 200, "message": "OK", "payload": {"channels": "a"}, "uuid": "my-custom-uuid", "service": "Presence"}`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("Response parsing channels", err.Error())
}

func TestNewGetStateResponseParsingPayloadChannelError(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"status": 200, "message": "OK", "payload": null, "uuid": "my-custom-uuid", "channel": "my-channel", "service": "Presence"}`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("Response parsing channel", err.Error())
}

func TestNewGetStateResponseParsingChannelError(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"status": 200, "message": "OK", "uuid": "my-custom-uuid", "channel": "my-channel", "service": "Presence"}`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("Response parsing channel", err.Error())
}

func TestNewGetStateResponseParsingChannelNull(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"status": 200, "message": "OK", "uuid": "my-custom-uuid", "channel": {}, "service": "Presence"}`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("Response parsing channel 2", err.Error())
}
