package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestNewGetStateResponse(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.SetUserId(UserId("my-custom-uuid"))

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

	pubnub.Config.SetUserId(UserId("my-custom-uuid"))

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

	pubnub.Config.SetUserId(UserId("my-custom-uuid"))

	jsonBytes := []byte(`{"status": 400, "error": 1, "message": "Invalid JSON specified.", "service": "Presence"}`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("Invalid JSON specified.", err.Error())
}

func TestGetStateBasicRequest(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.SetUserId(UserId("my-custom-uuid"))

	opts := newGetStateOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch"}
	opts.ChannelGroups = []string{"cg"}

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

	opts := newGetStateOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch"}
	opts.ChannelGroups = []string{"cg"}
	opts.UUID = uuid

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

	pubnub.Config.SetUserId(UserId("my-custom-uuid"))

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

	pubnub.Config.SetUserId(UserId("my-custom-uuid"))

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

	pubnub.Config.SetUserId(UserId("my-custom-uuid"))

	o := newGetStateBuilderWithContext(pubnub, pubnub.ctx)
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

	pubnub.Config.SetUserId(UserId("my-custom-uuid"))

	opts := newGetStateOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroups = []string{"cg1", "cg2", "cg3"}

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
	opts := newGetStateOpts(pn, pn.ctx)
	assert.Equal("pubnub/validation: pubnub: Get State: Missing Channel or Channel Group", opts.validate().Error())
}

func TestGetStateValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newGetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroups = []string{"cg1", "cg2", "cg3"}

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
	assert.Equal("response parsing error", err.Error())
}

func TestNewGetStateResponseParsingPayloadError(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"status": 200, "message": "OK", "payload": "error", "uuid": "my-custom-uuid", "service": "Presence"}`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("response parsing payload 2", err.Error())
}

func TestNewGetStateResponseParsingPayloadChannelsError(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"status": 200, "message": "OK", "payload": {"channels": "a"}, "uuid": "my-custom-uuid", "service": "Presence"}`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("response parsing channels", err.Error())
}

func TestNewGetStateResponseParsingPayloadChannelError(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"status": 200, "message": "OK", "payload": null, "uuid": "my-custom-uuid", "channel": "my-channel", "service": "Presence"}`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("response parsing channel", err.Error())
}

func TestNewGetStateResponseParsingChannelError(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"status": 200, "message": "OK", "uuid": "my-custom-uuid", "channel": "my-channel", "service": "Presence"}`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("response parsing channel", err.Error())
}

func TestNewGetStateResponseParsingChannelNull(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"status": 200, "message": "OK", "uuid": "my-custom-uuid", "channel": {}, "service": "Presence"}`)

	_, _, err := newGetStateResponse(jsonBytes, fakeResponseState)
	assert.Equal("response parsing channel 2", err.Error())
}

// HTTP Method and Operation Tests

func TestGetStateHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)

	assert.Equal("GET", opts.httpMethod())
}

func TestGetStateOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)

	assert.Equal(PNGetStateOperation, opts.operationType())
}

func TestGetStateIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestGetStateTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (5 setters)

func TestGetStateBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetStateBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestGetStateBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetStateBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestGetStateBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetStateBuilder(pn)

	// Test Channels setter
	channels := []string{"channel1", "channel2"}
	builder.Channels(channels)
	assert.Equal(channels, builder.opts.Channels)

	// Test ChannelGroups setter
	channelGroups := []string{"group1", "group2"}
	builder.ChannelGroups(channelGroups)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)

	// Test UUID setter
	builder.UUID("test-uuid")
	assert.Equal("test-uuid", builder.opts.UUID)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Test Transport setter
	transport := &http.Transport{}
	builder.Transport(transport)
	assert.Equal(transport, builder.opts.Transport)
}

func TestGetStateBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	channels := []string{"channel1"}
	channelGroups := []string{"group1"}
	queryParam := map[string]string{"key": "value"}
	transport := &http.Transport{}

	builder := newGetStateBuilder(pn)
	result := builder.Channels(channels).
		ChannelGroups(channelGroups).
		UUID("test-uuid").
		QueryParam(queryParam).
		Transport(transport)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)
	assert.Equal("test-uuid", builder.opts.UUID)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

func TestGetStateBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetStateBuilder(pn)

	// Verify default values
	assert.Nil(builder.opts.Channels)
	assert.Nil(builder.opts.ChannelGroups)
	assert.Equal("", builder.opts.UUID) // UUID defaults to empty, uses Config.UUID in buildPath
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestGetStateBuilderChannelCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name          string
		channels      []string
		channelGroups []string
		description   string
	}{
		{
			name:        "Single channel",
			channels:    []string{"channel1"},
			description: "Get state for single channel",
		},
		{
			name:        "Multiple channels",
			channels:    []string{"channel1", "channel2", "channel3"},
			description: "Get state for multiple channels",
		},
		{
			name:          "Single channel group",
			channelGroups: []string{"group1"},
			description:   "Get state for single channel group",
		},
		{
			name:          "Multiple channel groups",
			channelGroups: []string{"group1", "group2", "group3"},
			description:   "Get state for multiple channel groups",
		},
		{
			name:          "Channels and groups combination",
			channels:      []string{"channel1", "channel2"},
			channelGroups: []string{"group1", "group2"},
			description:   "Get state for both channels and channel groups",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetStateBuilder(pn)
			if tc.channels != nil {
				builder.Channels(tc.channels)
			}
			if tc.channelGroups != nil {
				builder.ChannelGroups(tc.channelGroups)
			}

			assert.Equal(tc.channels, builder.opts.Channels)
			assert.Equal(tc.channelGroups, builder.opts.ChannelGroups)
		})
	}
}

func TestGetStateBuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	transport := &http.Transport{}
	channels := []string{"channel1", "channel2"}
	channelGroups := []string{"group1", "group2"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Test all 5 setters in chain
	builder := newGetStateBuilder(pn).
		Channels(channels).
		ChannelGroups(channelGroups).
		UUID("test-uuid").
		QueryParam(queryParam).
		Transport(transport)

	// Verify all are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)
	assert.Equal("test-uuid", builder.opts.UUID)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestGetStateBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}
	opts.UUID = "test-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/channel/test-channel/uuid/test-uuid"
	assert.Equal(expected, path)
}

func TestGetStateBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newGetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"my-channel"}
	opts.UUID = "my-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/custom-sub-key/channel/my-channel/uuid/my-uuid"
	assert.Equal(expected, path)
}

func TestGetStateBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"channel@with#symbols"}
	opts.UUID = "uuid-with-special@chars"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/presence/sub-key/demo/channel/")
	assert.Contains(path, "/uuid/")
	// Should be URL encoded
	assert.Contains(path, "channel%40with%23symbols")
	assert.Contains(path, "uuid-with-special%40chars")
}

func TestGetStateBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"频道中文"}
	opts.UUID = "用户ID-пользователь"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/presence/sub-key/demo/channel/")
	assert.Contains(path, "/uuid/")
	// Should be URL encoded
	assert.Contains(path, "%E9%A2%91%E9%81%93%E4%B8%AD%E6%96%87") // URL encoded Chinese
}

func TestGetStateBuildPathMultipleChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1", "channel2", "channel3"}
	opts.UUID = "test-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/channel/channel1,channel2,channel3/uuid/test-uuid"
	assert.Equal(expected, path)
}

func TestGetStateBuildPathUUIDDefaultBehavior(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.UUID = "config-uuid"
	opts := newGetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}
	// Don't set UUID explicitly - should use Config.UUID

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/channel/test-channel/uuid/config-uuid"
	assert.Equal(expected, path)
}

// JSON Body Building Tests (CRITICAL for GET operation - should be empty)

func TestGetStateBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations should have empty body
	assert.Equal([]byte{}, body)
}

func TestGetStateBuildBodyWithAllParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)

	// Set all possible parameters - body should still be empty for GET
	opts.Channels = []string{"channel1", "channel2"}
	opts.ChannelGroups = []string{"group1", "group2"}
	opts.UUID = "test-uuid"
	opts.QueryParam = map[string]string{"param": "value"}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations always have empty body regardless of parameters
	assert.Equal([]byte{}, body)
}

func TestGetStateBuildBodyErrorScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)

	// Even with potential error conditions, buildBody should not fail for GET
	opts.Channels = []string{} // Empty channels
	opts.UUID = ""             // Empty UUID

	body, err := opts.buildBody()
	assert.Nil(err) // buildBody should never error for GET operations
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}

// Query Parameter Tests

func TestGetStateBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestGetStateBuildQueryWithChannelGroups(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)

	opts.ChannelGroups = []string{"group1", "group2"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	channelGroupValue := query.Get("channel-group")
	assert.Equal("group1,group2", channelGroupValue)
}

func TestGetStateBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)

	customParams := map[string]string{
		"custom":         "value",
		"special_chars":  "value@with#symbols",
		"unicode":        "测试参数",
		"empty_value":    "",
		"number_string":  "42",
		"boolean_string": "true",
	}
	opts.QueryParam = customParams

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all custom parameters are present
	for key, expectedValue := range customParams {
		actualValue := query.Get(key)
		if key == "special_chars" {
			// Special characters should be URL encoded
			assert.Contains(actualValue, "%", "Query parameter %s should be URL encoded", key)
		} else if key == "unicode" {
			// Unicode should be URL encoded
			assert.Contains(actualValue, "%", "Query parameter %s should contain URL encoded Unicode", key)
		} else {
			assert.Equal(expectedValue, actualValue, "Query parameter %s should match", key)
		}
	}
}

func TestGetStateBuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetStateOpts(pn, pn.ctx)

	// Set all possible query parameters
	opts.ChannelGroups = []string{"group1", "group2"}
	opts.QueryParam = map[string]string{
		"extra": "parameter",
		"debug": "true",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all parameters are set correctly
	assert.Equal("group1,group2", query.Get("channel-group"))
	assert.Equal("parameter", query.Get("extra"))
	assert.Equal("true", query.Get("debug"))

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestGetStateBuildQueryEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name          string
		channelGroups []string
		queryParam    map[string]string
		expectValue   string
	}{
		{
			name:          "Empty channel groups",
			channelGroups: []string{},
			expectValue:   "",
		},
		{
			name:          "Nil channel groups",
			channelGroups: nil,
			expectValue:   "",
		},
		{
			name:          "Single channel group",
			channelGroups: []string{"single-group"},
			expectValue:   "single-group",
		},
		{
			name:          "Channel groups with special chars",
			channelGroups: []string{"group@with#symbols", "group-with-dashes"},
			expectValue:   "group%40with%23symbols,group-with-dashes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetStateOpts(pn, pn.ctx)
			opts.ChannelGroups = tc.channelGroups
			if tc.queryParam != nil {
				opts.QueryParam = tc.queryParam
			}

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expectValue, query.Get("channel-group"))
		})
	}
}

// GET-Specific Tests (State Retrieval Characteristics)

func TestGetStateGetOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetStateBuilder(pn)
	builder.Channels([]string{"test-channel"})
	builder.UUID("test-uuid")

	// Verify it's a GET operation
	assert.Equal("GET", builder.opts.httpMethod())

	// GET operations have empty body
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	// Should have proper path for state retrieval
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/presence/sub-key/demo/channel/test-channel/uuid/test-uuid")
}

func TestGetStateStateRetrievalValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*getStateOpts)
		description string
	}{
		{
			name: "Basic state retrieval",
			setupOpts: func(opts *getStateOpts) {
				opts.Channels = []string{"channel1"}
				opts.UUID = "user123"
			},
			description: "Get state for specific channel and UUID",
		},
		{
			name: "State retrieval from multiple channels",
			setupOpts: func(opts *getStateOpts) {
				opts.Channels = []string{"channel1", "channel2", "channel3"}
				opts.UUID = "user123"
			},
			description: "Get state from multiple channels",
		},
		{
			name: "State retrieval from channel groups",
			setupOpts: func(opts *getStateOpts) {
				opts.ChannelGroups = []string{"group1", "group2"}
				opts.UUID = "user123"
			},
			description: "Get state from channel groups",
		},
		{
			name: "State retrieval with mixed channels and groups",
			setupOpts: func(opts *getStateOpts) {
				opts.Channels = []string{"channel1"}
				opts.ChannelGroups = []string{"group1"}
				opts.UUID = "user123"
			},
			description: "Get state from both channels and channel groups",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetStateOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			// Should pass validation
			assert.Nil(opts.validate())

			// Should be GET operation
			assert.Equal("GET", opts.httpMethod())

			// Should have empty body
			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)

			// Should build valid path
			path, err := opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/v2/presence/sub-key/")
			assert.Contains(path, "/channel/")
			assert.Contains(path, "/uuid/")

			// Should build valid query
			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

func TestGetStateResponseStructureValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetStateBuilder(pn)
	builder.Channels([]string{"test-channel"})
	builder.UUID("test-uuid")

	// Response should contain state data after GET operation
	// This is tested in the existing response parsing tests
	// but verify the operation is configured correctly
	opts := builder.opts

	// Verify operation is configured correctly
	assert.Equal("GET", opts.httpMethod())
	assert.Equal(PNGetStateOperation, opts.operationType())
	assert.True(opts.isAuthRequired())
}

func TestGetStateEmptyBodyVerification(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that GET operations always have empty body regardless of configuration
	testCases := []struct {
		name      string
		setupOpts func(*getStateOpts)
	}{
		{
			name: "With all parameters set",
			setupOpts: func(opts *getStateOpts) {
				opts.Channels = []string{"channel1", "channel2"}
				opts.ChannelGroups = []string{"group1", "group2"}
				opts.UUID = "test-uuid"
				opts.QueryParam = map[string]string{
					"param1": "value1",
					"param2": "value2",
				}
			},
		},
		{
			name: "With minimal parameters",
			setupOpts: func(opts *getStateOpts) {
				opts.Channels = []string{"simple-channel"}
				opts.UUID = "simple-uuid"
			},
		},
		{
			name: "With empty/nil parameters",
			setupOpts: func(opts *getStateOpts) {
				opts.Channels = []string{}
				opts.ChannelGroups = nil
				opts.UUID = ""
				opts.QueryParam = nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetStateOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
			assert.Equal([]byte{}, body)
		})
	}
}

// Comprehensive Edge Case Tests

func TestGetStateWithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*getStateBuilder)
	}{
		{
			name: "Very long UUID",
			setupFn: func(builder *getStateBuilder) {
				longUUID := strings.Repeat("VeryLongUUID", 50) // 600 characters
				builder.UUID(longUUID)
				builder.Channels([]string{"test-channel"})
			},
		},
		{
			name: "Many channels",
			setupFn: func(builder *getStateBuilder) {
				var manyChannels []string
				for i := 0; i < 100; i++ {
					manyChannels = append(manyChannels, fmt.Sprintf("channel_%d", i))
				}
				builder.Channels(manyChannels)
				builder.UUID("test-uuid")
			},
		},
		{
			name: "Many channel groups",
			setupFn: func(builder *getStateBuilder) {
				var manyGroups []string
				for i := 0; i < 100; i++ {
					manyGroups = append(manyGroups, fmt.Sprintf("group_%d", i))
				}
				builder.ChannelGroups(manyGroups)
				builder.UUID("test-uuid")
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *getStateBuilder) {
				largeQueryParam := make(map[string]string)
				for i := 0; i < 100; i++ {
					largeQueryParam[fmt.Sprintf("param_%d", i)] = strings.Repeat(fmt.Sprintf("value_%d_", i), 20)
				}
				builder.QueryParam(largeQueryParam)
				builder.Channels([]string{"test-channel"})
				builder.UUID("test-uuid")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetStateBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation for most cases (except when no channels/groups)
			if len(builder.opts.Channels) > 0 || len(builder.opts.ChannelGroups) > 0 {
				assert.Nil(builder.opts.validate())
			}

			// Should build valid path and query
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Should always have empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
		})
	}
}

func TestGetStateSpecialCharacterHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialStrings := []string{
		"<script>alert('xss')</script>",
		"SELECT * FROM channels; DROP TABLE channels;",
		"newline\ncharacter\ttab\rcarriage",
		"   ",                // Only spaces
		"\u0000\u0001\u0002", // Control characters
		"\"quoted_string\"",
		"'single_quoted'",
		"back`tick`string",
		"emoji😀🎉🚀💯",
		"ñáéíóúüç", // Accented characters
	}

	for i, specialString := range specialStrings {
		t.Run(fmt.Sprintf("SpecialString_%d", i), func(t *testing.T) {
			builder := newGetStateBuilder(pn)
			builder.Channels([]string{specialString})
			builder.ChannelGroups([]string{specialString})
			builder.UUID(specialString)
			builder.QueryParam(map[string]string{
				"special_field": specialString,
			})

			// Should pass validation (basic validation doesn't check content)
			assert.Nil(builder.opts.validate())

			// Should build valid path and query
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Should always have empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
		})
	}
}

func TestGetStateParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name           string
		channels       []string
		channelGroups  []string
		uuid           string
		shouldValidate bool
	}{
		{
			name:           "Empty string channel",
			channels:       []string{""},
			uuid:           "test",
			shouldValidate: true,
		},
		{
			name:           "Single character channel",
			channels:       []string{"a"},
			uuid:           "a",
			shouldValidate: true,
		},
		{
			name:           "Unicode-only channel",
			channels:       []string{"测试"},
			uuid:           "测试",
			shouldValidate: true,
		},
		{
			name:           "Empty channel groups",
			channelGroups:  []string{""},
			uuid:           "test",
			shouldValidate: true,
		},
		{
			name:           "No channels or groups",
			channels:       []string{},
			uuid:           "test",
			shouldValidate: false, // Should fail validation
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetStateBuilder(pn)
			if tc.channels != nil {
				builder.Channels(tc.channels)
			}
			if tc.channelGroups != nil {
				builder.ChannelGroups(tc.channelGroups)
			}
			builder.UUID(tc.uuid)

			// Validation should match expectation
			err := builder.opts.validate()
			if tc.shouldValidate {
				assert.Nil(err)
			} else {
				assert.NotNil(err)
			}

			// Should build valid components when validation passes
			if tc.shouldValidate {
				path, err := builder.opts.buildPath()
				assert.Nil(err)
				if len(tc.channels) > 0 {
					assert.Contains(path, "/channel/")
				}

				query, err := builder.opts.buildQuery()
				assert.Nil(err)
				assert.NotNil(query)

				body, err := builder.opts.buildBody()
				assert.Nil(err)
				assert.Empty(body) // GET operation always has empty body
			}
		})
	}
}

func TestGetStateComplexStateScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*getStateBuilder)
		validateFn func(*testing.T, string)
	}{
		{
			name: "International channels with complex UUID",
			setupFn: func(builder *getStateBuilder) {
				builder.Channels([]string{"频道中文123", "канал-русский"})
				builder.UUID("用户中文123")
			},
			validateFn: func(t *testing.T, path string) {
				assert.Contains(path, "/v2/presence/sub-key/demo/channel/")
				assert.Contains(path, "/uuid/")
			},
		},
		{
			name: "Professional state retrieval with groups",
			setupFn: func(builder *getStateBuilder) {
				builder.Channels([]string{"company-channel-1", "company-channel-2"})
				builder.ChannelGroups([]string{"company-groups", "admin-groups"})
				builder.UUID("professional@company.com")
			},
			validateFn: func(t *testing.T, path string) {
				assert.Contains(path, "/uuid/professional%40company.com")
			},
		},
		{
			name: "Email-like UUID with mixed channel types",
			setupFn: func(builder *getStateBuilder) {
				builder.Channels([]string{"public-channel", "private_channel"})
				builder.ChannelGroups([]string{"user-groups", "system-groups"})
				builder.UUID("user@company.com")
				builder.QueryParam(map[string]string{
					"include_metadata": "true",
					"format":           "detailed",
				})
			},
			validateFn: func(t *testing.T, path string) {
				assert.Contains(path, "/uuid/user%40company.com")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetStateBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path
			path, err := builder.opts.buildPath()
			assert.Nil(err)

			// Should always have empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)

			// Run verification
			tc.validateFn(t, path)
		})
	}
}

// Error Scenario Tests

func TestGetStateExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newGetStateBuilder(pn)
	builder.Channels([]string{"test-channel"})
	builder.UUID("test-uuid")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetStatePathBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCases := []struct {
		name         string
		subscribeKey string
		channels     []string
		uuid         string
		expectError  bool
	}{
		{
			name:         "Empty SubscribeKey",
			subscribeKey: "",
			channels:     []string{"test-channel"},
			uuid:         "test-uuid",
			expectError:  false, // buildPath doesn't validate, only validate() does
		},
		{
			name:         "Empty channels",
			subscribeKey: "demo",
			channels:     []string{},
			uuid:         "test-uuid",
			expectError:  false, // buildPath doesn't validate channels
		},
		{
			name:         "Empty UUID",
			subscribeKey: "demo",
			channels:     []string{"test-channel"},
			uuid:         "",
			expectError:  false, // buildPath uses Config.UUID if empty
		},
		{
			name:         "SubscribeKey with spaces",
			subscribeKey: "   sub key   ",
			channels:     []string{"test-channel"},
			uuid:         "test-uuid",
			expectError:  false,
		},
		{
			name:         "Very special characters",
			subscribeKey: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			channels:     []string{"!@#$%^&*()_+-=[]{}|;':\",./<>?"},
			uuid:         "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			expectError:  false,
		},
		{
			name:         "Unicode SubscribeKey, channels and UUID",
			subscribeKey: "测试订阅键-русский-キー",
			channels:     []string{"频道测试-канал-チャンネル"},
			uuid:         "用户ID-пользователь-ユーザーID",
			expectError:  false,
		},
		{
			name:         "Very long values",
			subscribeKey: strings.Repeat("a", 1000),
			channels:     []string{strings.Repeat("b", 1000)},
			uuid:         strings.Repeat("c", 1000),
			expectError:  false,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			pn.Config.SubscribeKey = tc.subscribeKey
			opts := newGetStateOpts(pn, pn.ctx)
			opts.Channels = tc.channels
			opts.UUID = tc.uuid

			path, err := opts.buildPath()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Contains(path, "/v2/presence/sub-key/")
				assert.Contains(path, "/channel/")
				assert.Contains(path, "/uuid/")
			}
		})
	}
}

func TestGetStateQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*getStateOpts)
		expectError bool
	}{
		{
			name: "Nil query params",
			setupOpts: func(opts *getStateOpts) {
				opts.QueryParam = nil
			},
			expectError: false,
		},
		{
			name: "Empty query params",
			setupOpts: func(opts *getStateOpts) {
				opts.QueryParam = map[string]string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *getStateOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *getStateOpts) {
				opts.QueryParam = map[string]string{
					"special@key":   "special@value",
					"unicode测试":     "unicode值",
					"with spaces":   "also spaces",
					"equals=key":    "equals=value",
					"ampersand&key": "ampersand&value",
				}
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetStateOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			query, err := opts.buildQuery()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.NotNil(query)
			}
		})
	}
}

func TestGetStateBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newGetStateBuilder(pn)

	channels := []string{"channel1", "channel2"}
	channelGroups := []string{"group1", "group2"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.Channels(channels).
		ChannelGroups(channelGroups).
		UUID("complete-test-uuid").
		QueryParam(queryParam)

	// Verify all values are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)
	assert.Equal("complete-test-uuid", builder.opts.UUID)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v2/presence/sub-key/demo/channel/channel1,channel2/uuid/complete-test-uuid"
	assert.Equal(expectedPath, path)

	// Should build query with all params
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("group1,group2", query.Get("channel-group"))
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))

	// Should always have empty body (GET operation)
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}

func TestGetStateValidationErrors(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name          string
		setupOpts     func(*getStateOpts)
		expectedError string
	}{
		{
			name: "Missing subscribe key",
			setupOpts: func(opts *getStateOpts) {
				opts.pubnub.Config.SubscribeKey = ""
				opts.Channels = []string{"channel1"}
			},
			expectedError: "Missing Subscribe Key",
		},
		{
			name: "Missing channels and channel groups",
			setupOpts: func(opts *getStateOpts) {
				// Keep valid SubscribeKey so we test the channel/group validation
				opts.Channels = []string{}
				opts.ChannelGroups = []string{}
			},
			expectedError: "Missing Channel or Channel Group",
		},
		{
			name: "Nil channels and channel groups",
			setupOpts: func(opts *getStateOpts) {
				// Keep valid SubscribeKey so we test the channel/group validation
				opts.Channels = nil
				opts.ChannelGroups = nil
			},
			expectedError: "Missing Channel or Channel Group",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh PubNub instance for each test case to avoid shared state
			pn := NewPubNub(NewDemoConfig())
			opts := newGetStateOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			err := opts.validate()
			assert.NotNil(err)
			assert.Contains(err.Error(), tc.expectedError)
		})
	}
}

// Extended Response Parsing Tests

func TestGetStateResponseParsingEdgeCases(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name        string
		jsonBytes   []byte
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Empty payload with channel",
			jsonBytes:   []byte(`{"status": 200, "message": "OK", "payload": {}, "uuid": "my-uuid", "channel": "my-channel", "service": "Presence"}`),
			expectError: false,
		},
		{
			name:        "Null values",
			jsonBytes:   []byte(`{"status": 200, "message": "OK", "payload": null, "uuid": null, "channel": null, "service": "Presence"}`),
			expectError: false, // This actually doesn't error - it just creates empty state
		},
		{
			name:        "Complex nested payload",
			jsonBytes:   []byte(`{"status": 200, "message": "OK", "payload": {"channels": {"ch1": {"nested": {"deep": "value"}}, "ch2": {"simple": "data"}}}, "uuid": "test-uuid", "service": "Presence"}`),
			expectError: false,
		},
		{
			name:        "Number values in payload",
			jsonBytes:   []byte(`{"status": 200, "message": "OK", "payload": {"age": 25, "score": 98.5}, "uuid": "test-uuid", "channel": "test-channel", "service": "Presence"}`),
			expectError: false,
		},
		{
			name:        "Boolean values in payload",
			jsonBytes:   []byte(`{"status": 200, "message": "OK", "payload": {"active": true, "verified": false}, "uuid": "test-uuid", "channel": "test-channel", "service": "Presence"}`),
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, _, err := newGetStateResponse(tc.jsonBytes, fakeResponseState)

			if tc.expectError {
				assert.NotNil(err)
				if tc.errorMsg != "" {
					assert.Contains(err.Error(), tc.errorMsg)
				}
				// When an error occurs, response might be nil
			} else {
				assert.Nil(err)
				assert.NotNil(resp)
				assert.NotNil(resp.State)
			}
		})
	}
}
