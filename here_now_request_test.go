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

func TestHereNowChannelsGroups(t *testing.T) {
	assert := assert.New(t)

	opts := newHereNowOpts(pubnub, pubnub.ctx)

	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroups = []string{"cg1", "cg2", "cg3"}
	opts.IncludeUUIDs = true
	opts.SetIncludeUUIDs = true

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
	expected.Set("disable-uuids", "0")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid", "limit"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHereNowNoChannel(t *testing.T) {
	assert := assert.New(t)

	opts := newHereNowOpts(pubnub, pubnub.ctx)

	opts.ChannelGroups = []string{"cg1", "cg2", "cg3"}

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v2/presence/sub_key/sub_key/channel/,", path)
}

func TestNewHereNowBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newHereNowBuilder(pubnub)
	o.ChannelGroups([]string{"cg1", "cg2", "cg3"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v2/presence/sub_key/sub_key/channel/,", path)
}

func TestNewHereNowBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newHereNowBuilderWithContext(pubnub, pubnub.ctx)
	o.ChannelGroups([]string{"cg1", "cg2", "cg3"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v2/presence/sub_key/sub_key/channel/,", path)
}

func TestHereNowMultipleWithOpts(t *testing.T) {
	assert := assert.New(t)

	opts := newHereNowOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroups = []string{"cg1", "cg2", "cg3"}
	opts.IncludeUUIDs = false
	opts.IncludeState = true
	opts.SetIncludeState = true
	opts.SetIncludeUUIDs = true

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
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid", "limit"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHereNowMultipleWithOptsQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := newHereNowOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroups = []string{"cg1", "cg2", "cg3"}
	opts.IncludeUUIDs = false
	opts.IncludeState = true
	opts.SetIncludeState = true
	opts.SetIncludeUUIDs = true

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	expected.Set("state", "1")
	expected.Set("disable-uuids", "1")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid", "limit"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHereNowGlobal(t *testing.T) {
	assert := assert.New(t)

	opts := newHereNowOpts(pubnub, pubnub.ctx)

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
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid", "limit"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHereNowValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newHereNowOpts(pn, pn.ctx)

	assert.Equal("pubnub/validation: pubnub: Here Now: Missing Subscribe Key", opts.validate().Error())
}

func TestHereNowBuildPath(t *testing.T) {
	assert := assert.New(t)
	opts := newHereNowOpts(pubnub, pubnub.ctx)

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v2/presence/sub_key/sub_key", path)

}

func TestHereNowBuildQuery(t *testing.T) {
	assert := assert.New(t)

	opts := newHereNowOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroups = []string{"cg1", "cg2", "cg3"}
	opts.IncludeUUIDs = false
	opts.IncludeState = true
	opts.SetIncludeState = true
	opts.SetIncludeUUIDs = false

	query, err := opts.buildQuery()
	assert.Nil(err)
	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	expected.Set("state", "1")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid", "limit"}, []string{})

}

func TestNewHereNowResponseErrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newHereNowResponse(jsonBytes, nil, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestNewHereNowResponseOneChannel(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte("{\"status\":200,\"message\":\"OK\",\"service\":\"Presence\",\"uuids\":[{\"uuid\":\"a3ffd012-a3b9-478c-8705-64089f24d71e\",\"state\":{\"age\":10}}],\"occupancy\":1}")

	_, _, err := newHereNowResponse(jsonBytes, []string{"a"}, StatusResponse{})
	assert.Nil(err)
}

func TestNewHereNowResponseOccupancyZero(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte("{\"status\":200,\"message\":\"OK\",\"service\":\"Presence\",\"occupancy\":0,\"total_channels\":1,\"total_occupancy\":1}")

	r, _, err := newHereNowResponse(jsonBytes, []string{"a"}, StatusResponse{})
	assert.Nil(err)
	assert.Equal(1, r.TotalChannels)
	assert.Equal(0, r.TotalOccupancy)

}

func TestNewHereNowResponseOccupancyZeroPayload(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte("{\"status\":200,\"message\":\"OK\",\"service\":\"Presence\",\"occupancy\":\"0\",\"total_channels\":1,\"total_occupancy\":1}")

	r, _, err := newHereNowResponse(jsonBytes, []string{"a"}, StatusResponse{})
	assert.Nil(err)
	assert.Equal(1, r.TotalChannels)
	assert.Equal(0, r.TotalOccupancy)
}

func TestNewHereNowResponseOccupancyZeroPayloadWithCh(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte("{\"status\":200,\"message\":\"OK\",\"payload\":{\"total_occupancy\":3,\"total_channels\":1,\"channels\":{\"ch1\":{\"occupancy\":1,\"uuids\":[{\"uuid\":\"user1\",\"state\":{\"age\":10}}]}}},\"service\":\"Presence\"}")

	r, _, err := newHereNowResponse(jsonBytes, []string{"a"}, StatusResponse{})
	assert.Nil(err)
	assert.Equal(1, r.TotalChannels)
	assert.Equal(3, r.TotalOccupancy)
}

func TestNewHereNowResponseOccupancyZeroPayloadWithoutCh(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte("{\"status\":200,\"message\":\"OK\",\"payload\":{\"total_occupancy\":3,\"total_channels\":2},\"service\":\"Presence\"}")

	r, _, err := newHereNowResponse(jsonBytes, []string{"a"}, StatusResponse{})
	assert.Nil(err)
	assert.Equal(1, r.TotalChannels)
	assert.Equal(0, r.TotalOccupancy)

}

// HTTP Method and Operation Tests

func TestHereNowHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	assert.Equal("GET", opts.httpMethod())
}

func TestHereNowOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	assert.Equal(PNHereNowOperation, opts.operationType())
}

func TestHereNowIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestHereNowTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (6 setters)

func TestHereNowBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHereNowBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestHereNowBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHereNowBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestHereNowBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHereNowBuilder(pn)

	// Test Channels setter
	channels := []string{"channel1", "channel2"}
	builder.Channels(channels)
	assert.Equal(channels, builder.opts.Channels)

	// Test ChannelGroups setter
	channelGroups := []string{"group1", "group2"}
	builder.ChannelGroups(channelGroups)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)

	// Test IncludeState setter
	builder.IncludeState(true)
	assert.True(builder.opts.IncludeState)
	assert.True(builder.opts.SetIncludeState)

	// Test IncludeUUIDs setter
	builder.IncludeUUIDs(false)
	assert.False(builder.opts.IncludeUUIDs)
	assert.True(builder.opts.SetIncludeUUIDs)

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

func TestHereNowBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	channels := []string{"channel1"}
	channelGroups := []string{"group1"}
	queryParam := map[string]string{"key": "value"}
	transport := &http.Transport{}

	builder := newHereNowBuilder(pn)
	result := builder.Channels(channels).
		ChannelGroups(channelGroups).
		IncludeState(true).
		IncludeUUIDs(false).
		QueryParam(queryParam).
		Transport(transport)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)
	assert.True(builder.opts.IncludeState)
	assert.True(builder.opts.SetIncludeState)
	assert.False(builder.opts.IncludeUUIDs)
	assert.True(builder.opts.SetIncludeUUIDs)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

func TestHereNowBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHereNowBuilder(pn)

	// Verify default values
	assert.Nil(builder.opts.Channels)
	assert.Nil(builder.opts.ChannelGroups)
	assert.False(builder.opts.IncludeState)
	assert.False(builder.opts.IncludeUUIDs)
	assert.False(builder.opts.SetIncludeState)
	assert.False(builder.opts.SetIncludeUUIDs)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestHereNowBuilderChannelCombinations(t *testing.T) {
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
			description: "Get presence for single channel",
		},
		{
			name:        "Multiple channels",
			channels:    []string{"channel1", "channel2", "channel3"},
			description: "Get presence for multiple channels",
		},
		{
			name:          "Single channel group",
			channelGroups: []string{"group1"},
			description:   "Get presence for single channel group",
		},
		{
			name:          "Multiple channel groups",
			channelGroups: []string{"group1", "group2", "group3"},
			description:   "Get presence for multiple channel groups",
		},
		{
			name:          "Channels and groups combination",
			channels:      []string{"channel1", "channel2"},
			channelGroups: []string{"group1", "group2"},
			description:   "Get presence for both channels and channel groups",
		},
		{
			name:        "Global presence (no channels/groups)",
			description: "Get global presence for all subscribed channels",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newHereNowBuilder(pn)
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

func TestHereNowBuilderIncludeStateCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name         string
		includeState bool
		description  string
	}{
		{
			name:         "Include state true",
			includeState: true,
			description:  "Include user state in presence response",
		},
		{
			name:         "Include state false",
			includeState: false,
			description:  "Exclude user state from presence response",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newHereNowBuilder(pn)
			builder.IncludeState(tc.includeState)

			assert.Equal(tc.includeState, builder.opts.IncludeState)
			assert.True(builder.opts.SetIncludeState)
		})
	}
}

func TestHereNowBuilderIncludeUUIDsCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name         string
		includeUUIDs bool
		description  string
	}{
		{
			name:         "Include UUIDs true",
			includeUUIDs: true,
			description:  "Include user UUIDs in presence response",
		},
		{
			name:         "Include UUIDs false",
			includeUUIDs: false,
			description:  "Exclude user UUIDs from presence response",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newHereNowBuilder(pn)
			builder.IncludeUUIDs(tc.includeUUIDs)

			assert.Equal(tc.includeUUIDs, builder.opts.IncludeUUIDs)
			assert.True(builder.opts.SetIncludeUUIDs)
		})
	}
}

func TestHereNowBuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	transport := &http.Transport{}
	channels := []string{"channel1", "channel2"}
	channelGroups := []string{"group1", "group2"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Test all 6 setters in chain
	builder := newHereNowBuilder(pn).
		Channels(channels).
		ChannelGroups(channelGroups).
		IncludeState(true).
		IncludeUUIDs(false).
		QueryParam(queryParam).
		Transport(transport)

	// Verify all are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)
	assert.True(builder.opts.IncludeState)
	assert.True(builder.opts.SetIncludeState)
	assert.False(builder.opts.IncludeUUIDs)
	assert.True(builder.opts.SetIncludeUUIDs)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestHereNowBuildPathGlobal(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub_key/demo"
	assert.Equal(expected, path)
}

func TestHereNowBuildPathSingleChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub_key/demo/channel/test-channel"
	assert.Equal(expected, path)
}

func TestHereNowBuildPathMultipleChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1", "channel2", "channel3"}

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub_key/demo/channel/channel1,channel2,channel3"
	assert.Equal(expected, path)
}

func TestHereNowBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newHereNowOpts(pn, pn.ctx)
	opts.Channels = []string{"my-channel"}

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub_key/custom-sub-key/channel/my-channel"
	assert.Equal(expected, path)
}

func TestHereNowBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)
	opts.Channels = []string{"channel@with#symbols", "channel-with-dashes"}

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub_key/demo/channel/channel@with#symbols,channel-with-dashes"
	assert.Equal(expected, path)
}

func TestHereNowBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)
	opts.Channels = []string{"频道中文", "канал-русский"}

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub_key/demo/channel/频道中文,канал-русский"
	assert.Equal(expected, path)
}

func TestHereNowBuildPathChannelGroupsOnly(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)
	opts.ChannelGroups = []string{"group1", "group2"}

	path, err := opts.buildPath()
	assert.Nil(err)
	// When only channel groups, path uses empty channel list
	expected := "/v2/presence/sub_key/demo/channel/,"
	assert.Equal(expected, path)
}

func TestHereNowBuildPathChannelsAndGroups(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1", "channel2"}
	opts.ChannelGroups = []string{"group1", "group2"}

	path, err := opts.buildPath()
	assert.Nil(err)
	// Channels take precedence in path, groups go in query
	expected := "/v2/presence/sub_key/demo/channel/channel1,channel2"
	assert.Equal(expected, path)
}

// JSON Body Building Tests (CRITICAL for GET operation - should be empty)

func TestHereNowBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations should have empty body
	assert.Equal([]byte{}, body)
}

func TestHereNowBuildBodyWithAllParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	// Set all possible parameters - body should still be empty for GET
	opts.Channels = []string{"channel1", "channel2"}
	opts.ChannelGroups = []string{"group1", "group2"}
	opts.IncludeState = true
	opts.IncludeUUIDs = false
	opts.SetIncludeState = true
	opts.SetIncludeUUIDs = true
	opts.QueryParam = map[string]string{"param": "value"}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations always have empty body regardless of parameters
	assert.Equal([]byte{}, body)
}

func TestHereNowBuildBodyErrorScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	// Even with potential error conditions, buildBody should not fail for GET
	opts.Channels = []string{}      // Empty channels
	opts.ChannelGroups = []string{} // Empty groups

	body, err := opts.buildBody()
	assert.Nil(err) // buildBody should never error for GET operations
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}

// Query Parameter Tests

func TestHereNowBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestHereNowBuildQueryWithChannelGroups(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	opts.ChannelGroups = []string{"group1", "group2"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	channelGroupValue := query.Get("channel-group")
	assert.Equal("group1,group2", channelGroupValue)
}

func TestHereNowBuildQueryIncludeState(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name          string
		includeState  bool
		setFlag       bool
		expectedValue string
	}{
		{
			name:          "Include state true",
			includeState:  true,
			setFlag:       true,
			expectedValue: "1",
		},
		{
			name:          "Include state false",
			includeState:  false,
			setFlag:       true,
			expectedValue: "0",
		},
		{
			name:          "Include state not set",
			includeState:  false,
			setFlag:       false,
			expectedValue: "", // Parameter should not be present
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newHereNowOpts(pn, pn.ctx)
			opts.IncludeState = tc.includeState
			opts.SetIncludeState = tc.setFlag

			query, err := opts.buildQuery()
			assert.Nil(err)

			stateValue := query.Get("state")
			assert.Equal(tc.expectedValue, stateValue)
		})
	}
}

func TestHereNowBuildQueryIncludeUUIDs(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name          string
		includeUUIDs  bool
		setFlag       bool
		expectedValue string
	}{
		{
			name:          "Include UUIDs true",
			includeUUIDs:  true,
			setFlag:       true,
			expectedValue: "0", // disable-uuids=0 means include UUIDs
		},
		{
			name:          "Include UUIDs false",
			includeUUIDs:  false,
			setFlag:       true,
			expectedValue: "1", // disable-uuids=1 means exclude UUIDs
		},
		{
			name:          "Include UUIDs not set",
			includeUUIDs:  false,
			setFlag:       false,
			expectedValue: "", // Parameter should not be present
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newHereNowOpts(pn, pn.ctx)
			opts.IncludeUUIDs = tc.includeUUIDs
			opts.SetIncludeUUIDs = tc.setFlag

			query, err := opts.buildQuery()
			assert.Nil(err)

			disableUuidsValue := query.Get("disable-uuids")
			assert.Equal(tc.expectedValue, disableUuidsValue)
		})
	}
}

func TestHereNowBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	// Set original unencoded values for the opts
	opts.QueryParam = map[string]string{
		"custom":         "value",
		"special_chars":  "value@with#symbols",
		"unicode":        "测试参数",
		"empty_value":    "",
		"number_string":  "42",
		"boolean_string": "true",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify URL-encoded values (buildQuery encodes the parameters)
	assert.Equal("value", query.Get("custom"))
	assert.Equal("value%40with%23symbols", query.Get("special_chars"))
	assert.Equal("%E6%B5%8B%E8%AF%95%E5%8F%82%E6%95%B0", query.Get("unicode"))
	assert.Equal("", query.Get("empty_value"))
	assert.Equal("42", query.Get("number_string"))
	assert.Equal("true", query.Get("boolean_string"))
}

func TestHereNowBuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	// Set all possible query parameters
	opts.ChannelGroups = []string{"group1", "group2"}
	opts.IncludeState = true
	opts.SetIncludeState = true
	opts.IncludeUUIDs = false
	opts.SetIncludeUUIDs = true
	opts.QueryParam = map[string]string{
		"extra": "parameter",
		"debug": "true",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all parameters are set correctly
	assert.Equal("group1,group2", query.Get("channel-group"))
	assert.Equal("1", query.Get("state"))
	assert.Equal("1", query.Get("disable-uuids"))
	assert.Equal("parameter", query.Get("extra"))
	assert.Equal("true", query.Get("debug"))

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestHereNowBuildQueryEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name          string
		channelGroups []string
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
			expectValue:   "group@with#symbols,group-with-dashes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newHereNowOpts(pn, pn.ctx)
			opts.ChannelGroups = tc.channelGroups

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expectValue, query.Get("channel-group"))
		})
	}
}

// GET-Specific Tests (Presence Retrieval Characteristics)

func TestHereNowGetOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHereNowBuilder(pn)
	builder.Channels([]string{"test-channel"})
	builder.IncludeState(true)
	builder.IncludeUUIDs(true)

	// Verify it's a GET operation
	assert.Equal("GET", builder.opts.httpMethod())

	// GET operations have empty body
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	// Should have proper path for presence retrieval
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/presence/sub_key/demo/channel/test-channel")
}

func TestHereNowPresenceRetrievalValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*hereNowOpts)
		description string
	}{
		{
			name: "Global presence retrieval",
			setupOpts: func(opts *hereNowOpts) {
				// No channels or groups - global presence
			},
			description: "Get presence for all subscribed channels",
		},
		{
			name: "Single channel presence",
			setupOpts: func(opts *hereNowOpts) {
				opts.Channels = []string{"channel1"}
				opts.IncludeState = true
				opts.SetIncludeState = true
			},
			description: "Get presence for specific channel with state",
		},
		{
			name: "Multiple channels presence",
			setupOpts: func(opts *hereNowOpts) {
				opts.Channels = []string{"channel1", "channel2", "channel3"}
				opts.IncludeUUIDs = true
				opts.SetIncludeUUIDs = true
			},
			description: "Get presence for multiple channels with UUIDs",
		},
		{
			name: "Channel groups presence",
			setupOpts: func(opts *hereNowOpts) {
				opts.ChannelGroups = []string{"group1", "group2"}
				opts.IncludeState = false
				opts.SetIncludeState = true
			},
			description: "Get presence for channel groups without state",
		},
		{
			name: "Mixed channels and groups presence",
			setupOpts: func(opts *hereNowOpts) {
				opts.Channels = []string{"channel1"}
				opts.ChannelGroups = []string{"group1"}
				opts.IncludeState = true
				opts.IncludeUUIDs = false
				opts.SetIncludeState = true
				opts.SetIncludeUUIDs = true
			},
			description: "Get presence for both channels and channel groups",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newHereNowOpts(pn, pn.ctx)
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
			assert.Contains(path, "/v2/presence/sub_key/")

			// Should build valid query
			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

func TestHereNowResponseStructureValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHereNowBuilder(pn)
	builder.Channels([]string{"test-channel"})
	builder.IncludeState(true)
	builder.IncludeUUIDs(true)

	// Response should contain presence data after GET operation
	// This is tested in the existing response parsing tests
	// but verify the operation is configured correctly
	opts := builder.opts

	// Verify operation is configured correctly
	assert.Equal("GET", opts.httpMethod())
	assert.Equal(PNHereNowOperation, opts.operationType())
	assert.True(opts.isAuthRequired())
}

func TestHereNowEmptyBodyVerification(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that GET operations always have empty body regardless of configuration
	testCases := []struct {
		name      string
		setupOpts func(*hereNowOpts)
	}{
		{
			name: "With all parameters set",
			setupOpts: func(opts *hereNowOpts) {
				opts.Channels = []string{"channel1", "channel2"}
				opts.ChannelGroups = []string{"group1", "group2"}
				opts.IncludeState = true
				opts.IncludeUUIDs = false
				opts.SetIncludeState = true
				opts.SetIncludeUUIDs = true
				opts.QueryParam = map[string]string{
					"param1": "value1",
					"param2": "value2",
				}
			},
		},
		{
			name: "With minimal parameters",
			setupOpts: func(opts *hereNowOpts) {
				opts.Channels = []string{"simple-channel"}
			},
		},
		{
			name: "With empty/nil parameters",
			setupOpts: func(opts *hereNowOpts) {
				opts.Channels = []string{}
				opts.ChannelGroups = nil
				opts.QueryParam = nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newHereNowOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
			assert.Equal([]byte{}, body)
		})
	}
}

func TestHereNowGlobalVsChannelSpecific(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name          string
		channels      []string
		channelGroups []string
		expectedPath  string
		description   string
	}{
		{
			name:         "Global presence",
			expectedPath: "/v2/presence/sub_key/demo",
			description:  "No channels or groups specified - global presence",
		},
		{
			name:          "Channel groups only",
			channelGroups: []string{"group1"},
			expectedPath:  "/v2/presence/sub_key/demo/channel/,",
			description:   "Only channel groups - uses empty channel path",
		},
		{
			name:         "Single channel",
			channels:     []string{"channel1"},
			expectedPath: "/v2/presence/sub_key/demo/channel/channel1",
			description:  "Single channel specified",
		},
		{
			name:         "Multiple channels",
			channels:     []string{"channel1", "channel2"},
			expectedPath: "/v2/presence/sub_key/demo/channel/channel1,channel2",
			description:  "Multiple channels specified",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newHereNowOpts(pn, pn.ctx)
			opts.Channels = tc.channels
			opts.ChannelGroups = tc.channelGroups

			path, err := opts.buildPath()
			assert.Nil(err)
			assert.Equal(tc.expectedPath, path)
		})
	}
}

// Comprehensive Edge Case Tests

func TestHereNowWithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*hereNowBuilder)
	}{
		{
			name: "Many channels",
			setupFn: func(builder *hereNowBuilder) {
				var manyChannels []string
				for i := 0; i < 100; i++ {
					manyChannels = append(manyChannels, fmt.Sprintf("channel_%d", i))
				}
				builder.Channels(manyChannels)
			},
		},
		{
			name: "Many channel groups",
			setupFn: func(builder *hereNowBuilder) {
				var manyGroups []string
				for i := 0; i < 100; i++ {
					manyGroups = append(manyGroups, fmt.Sprintf("group_%d", i))
				}
				builder.ChannelGroups(manyGroups)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *hereNowBuilder) {
				largeQueryParam := make(map[string]string)
				for i := 0; i < 100; i++ {
					largeQueryParam[fmt.Sprintf("param_%d", i)] = strings.Repeat(fmt.Sprintf("value_%d_", i), 20)
				}
				builder.QueryParam(largeQueryParam)
				builder.Channels([]string{"test-channel"})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newHereNowBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation for all cases
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

func TestHereNowSpecialCharacterHandling(t *testing.T) {
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
			builder := newHereNowBuilder(pn)
			builder.Channels([]string{specialString})
			builder.ChannelGroups([]string{specialString})
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

func TestHereNowParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name          string
		channels      []string
		channelGroups []string
		description   string
	}{
		{
			name:        "Empty string channel",
			channels:    []string{""},
			description: "Channel with empty string",
		},
		{
			name:        "Single character channel",
			channels:    []string{"a"},
			description: "Channel with single character",
		},
		{
			name:        "Unicode-only channel",
			channels:    []string{"测试"},
			description: "Channel with Unicode characters",
		},
		{
			name:          "Empty channel groups",
			channelGroups: []string{""},
			description:   "Channel group with empty string",
		},
		{
			name:          "Empty arrays",
			channels:      []string{},
			channelGroups: []string{},
			description:   "Empty channels and groups - should use global presence",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newHereNowBuilder(pn)
			if tc.channels != nil {
				builder.Channels(tc.channels)
			}
			if tc.channelGroups != nil {
				builder.ChannelGroups(tc.channelGroups)
			}

			// Should pass validation (HereNow allows any configuration)
			assert.Nil(builder.opts.validate())

			// Should build valid components
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/v2/presence/sub_key/")

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body) // GET operation always has empty body
		})
	}
}

func TestHereNowComplexPresenceScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*hereNowBuilder)
		validateFn func(*testing.T, string, *url.Values)
	}{
		{
			name: "International channels with state",
			setupFn: func(builder *hereNowBuilder) {
				builder.Channels([]string{"频道中文123", "канал-русский"})
				builder.IncludeState(true)
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, "/v2/presence/sub_key/demo/channel/")
				assert.Equal("1", query.Get("state"))
			},
		},
		{
			name: "Professional presence monitoring",
			setupFn: func(builder *hereNowBuilder) {
				builder.Channels([]string{"company-channel-1", "company-channel-2"})
				builder.ChannelGroups([]string{"company-groups", "admin-groups"})
				builder.IncludeUUIDs(true)
				builder.IncludeState(false)
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, "company-channel-1,company-channel-2")
				assert.Equal("company-groups,admin-groups", query.Get("channel-group"))
				assert.Equal("0", query.Get("disable-uuids"))
				assert.Equal("0", query.Get("state"))
			},
		},
		{
			name: "Gaming presence with detailed monitoring",
			setupFn: func(builder *hereNowBuilder) {
				builder.Channels([]string{"game-lobby", "game-room-1"})
				builder.IncludeState(true)
				builder.IncludeUUIDs(true)
				builder.QueryParam(map[string]string{
					"game_mode": "battle_royale",
					"region":    "us-east",
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, "game-lobby,game-room-1")
				assert.Equal("1", query.Get("state"))
				assert.Equal("0", query.Get("disable-uuids"))
				assert.Equal("battle_royale", query.Get("game_mode"))
				assert.Equal("us-east", query.Get("region"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newHereNowBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path
			path, err := builder.opts.buildPath()
			assert.Nil(err)

			// Should build valid query
			query, err := builder.opts.buildQuery()
			assert.Nil(err)

			// Should always have empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)

			// Run custom validation
			tc.validateFn(t, path, query)
		})
	}
}

// Error Scenario Tests

func TestHereNowExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newHereNowBuilder(pn)
	builder.Channels([]string{"test-channel"})

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestHereNowPathBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCases := []struct {
		name         string
		subscribeKey string
		channels     []string
		expectError  bool
	}{
		{
			name:         "Empty SubscribeKey",
			subscribeKey: "",
			channels:     []string{"test-channel"},
			expectError:  false, // buildPath doesn't validate, only validate() does
		},
		{
			name:         "Empty channels",
			subscribeKey: "demo",
			channels:     []string{},
			expectError:  false, // buildPath handles empty channels (global presence)
		},
		{
			name:         "SubscribeKey with spaces",
			subscribeKey: "   sub key   ",
			channels:     []string{"test-channel"},
			expectError:  false,
		},
		{
			name:         "Very special characters",
			subscribeKey: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			channels:     []string{"!@#$%^&*()_+-=[]{}|;':\",./<>?"},
			expectError:  false,
		},
		{
			name:         "Unicode SubscribeKey and channels",
			subscribeKey: "测试订阅键-русский-キー",
			channels:     []string{"频道测试-канал-チャンネル"},
			expectError:  false,
		},
		{
			name:         "Very long values",
			subscribeKey: strings.Repeat("a", 1000),
			channels:     []string{strings.Repeat("b", 1000)},
			expectError:  false,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			pn.Config.SubscribeKey = tc.subscribeKey
			opts := newHereNowOpts(pn, pn.ctx)
			opts.Channels = tc.channels

			path, err := opts.buildPath()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Contains(path, "/v2/presence/sub_key/")
			}
		})
	}
}

func TestHereNowQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*hereNowOpts)
		expectError bool
	}{
		{
			name: "Nil query params",
			setupOpts: func(opts *hereNowOpts) {
				opts.QueryParam = nil
			},
			expectError: false,
		},
		{
			name: "Empty query params",
			setupOpts: func(opts *hereNowOpts) {
				opts.QueryParam = map[string]string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *hereNowOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *hereNowOpts) {
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
			opts := newHereNowOpts(pn, pn.ctx)
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

func TestHereNowBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newHereNowBuilder(pn)

	channels := []string{"channel1", "channel2"}
	channelGroups := []string{"group1", "group2"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.Channels(channels).
		ChannelGroups(channelGroups).
		IncludeState(true).
		IncludeUUIDs(false).
		QueryParam(queryParam)

	// Verify all values are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)
	assert.True(builder.opts.IncludeState)
	assert.True(builder.opts.SetIncludeState)
	assert.False(builder.opts.IncludeUUIDs)
	assert.True(builder.opts.SetIncludeUUIDs)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v2/presence/sub_key/demo/channel/channel1,channel2"
	assert.Equal(expectedPath, path)

	// Should build query with all params
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("group1,group2", query.Get("channel-group"))
	assert.Equal("1", query.Get("state"))
	assert.Equal("1", query.Get("disable-uuids"))
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))

	// Should always have empty body (GET operation)
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}

func TestHereNowValidationErrors(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name          string
		setupOpts     func(*hereNowOpts)
		expectedError string
	}{
		{
			name: "Missing subscribe key",
			setupOpts: func(opts *hereNowOpts) {
				opts.pubnub.Config.SubscribeKey = ""
			},
			expectedError: "Missing Subscribe Key",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh PubNub instance for each test case to avoid shared state
			pn := NewPubNub(NewDemoConfig())
			opts := newHereNowOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			err := opts.validate()
			assert.NotNil(err)
			assert.Contains(err.Error(), tc.expectedError)
		})
	}
}

// Extended Response Parsing Tests (Note: Basic response parsing is already tested in existing tests)

// Limit and Offset Tests

func TestHereNowBuilderLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHereNowBuilder(pn)
	builder.Limit(500)

	assert.Equal(500, builder.opts.Limit)
}

func TestHereNowBuilderOffset(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHereNowBuilder(pn)
	builder.Offset(100)

	assert.Equal(100, builder.opts.Offset)
	assert.True(builder.opts.SetOffset)
}

func TestHereNowBuildQueryWithLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test when limit is set
	opts := newHereNowOpts(pn, pn.ctx)
	opts.Limit = 500

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("500", query.Get("limit"))

	// Test default limit (1000)
	opts2 := newHereNowOpts(pn, pn.ctx)
	query2, err2 := opts2.buildQuery()
	assert.Nil(err2)
	assert.Equal("1000", query2.Get("limit"))
}

func TestHereNowBuildQueryWithOffset(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test when offset is set to non-zero value
	opts := newHereNowOpts(pn, pn.ctx)
	opts.Offset = 100
	opts.SetOffset = true

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("offset"))

	// Test when offset is set to 0 - should not appear in query
	opts2 := newHereNowOpts(pn, pn.ctx)
	opts2.Offset = 0
	opts2.SetOffset = true

	query2, err2 := opts2.buildQuery()
	assert.Nil(err2)
	assert.Equal("", query2.Get("offset"))

	// Test when offset is not set
	opts3 := newHereNowOpts(pn, pn.ctx)
	query3, err3 := opts3.buildQuery()
	assert.Nil(err3)
	assert.Equal("", query3.Get("offset"))
}

func TestHereNowBuildQueryWithLimitAndOffset(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHereNowOpts(pn, pn.ctx)

	// Set both limit and offset
	opts.Limit = 250
	opts.Offset = 500
	opts.SetOffset = true

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("250", query.Get("limit"))
	assert.Equal("500", query.Get("offset"))
}

func TestHereNowBuilderLimitOffset(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHereNowBuilder(pn)
	builder.Channels([]string{"test-channel"}).
		Limit(100).
		Offset(200)

	assert.Equal(100, builder.opts.Limit)
	assert.Equal(200, builder.opts.Offset)

	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("limit"))
	assert.Equal("200", query.Get("offset"))
}

func TestHereNowBuilderAllParametersIncludingLimitOffset(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	channels := []string{"channel1", "channel2"}
	channelGroups := []string{"group1", "group2"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Test all setters including new Limit and Offset
	builder := newHereNowBuilder(pn).
		Channels(channels).
		ChannelGroups(channelGroups).
		IncludeState(true).
		IncludeUUIDs(false).
		Limit(100).
		Offset(50).
		QueryParam(queryParam)

	// Verify all are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)
	assert.True(builder.opts.IncludeState)
	assert.True(builder.opts.SetIncludeState)
	assert.False(builder.opts.IncludeUUIDs)
	assert.True(builder.opts.SetIncludeUUIDs)
	assert.Equal(100, builder.opts.Limit)
	assert.Equal(50, builder.opts.Offset)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Verify query contains all parameters
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("group1,group2", query.Get("channel-group"))
	assert.Equal("1", query.Get("state"))
	assert.Equal("1", query.Get("disable-uuids"))
	assert.Equal("100", query.Get("limit"))
	assert.Equal("50", query.Get("offset"))
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))
}
