package pubnub

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestNewSetStateBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newSetStateBuilder(pubnub)
	o.Channels([]string{"ch1", "ch2", "ch3"})

	path, err := o.opts.buildPath()
	assert.Nil(err)

	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/uuid/%s/data",
			o.opts.pubnub.Config.UUID),
		u.EscapedPath(), []int{})
}

func TestNewSetStateBuilderWithUUID(t *testing.T) {
	assert := assert.New(t)

	o := newSetStateBuilder(pubnub)
	o.Channels([]string{"ch1", "ch2", "ch3"})
	uuid := "customuuid"
	o.UUID(uuid)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/uuid/%s/data",
			uuid),
		u.EscapedPath(), []int{})
}

func TestNewSetStateBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newSetStateBuilderWithContext(pubnub, backgroundContext)
	o.Channels([]string{"ch1", "ch2", "ch3"})

	path, err := o.opts.buildPath()
	assert.Nil(err)

	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/uuid/%s/data",
			o.opts.pubnub.Config.UUID),
		u.EscapedPath(), []int{})
}

func TestNewSetStateResponse(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.SetUserId(UserId("my-custom-uuid"))

	jsonBytes := []byte(`{"status": 200, "message": "OK", "payload": {"k": "v"}, "uuid": "my-custom-uuid", "service": "Presence"}`)

	res, _, err := newSetStateResponse(jsonBytes, fakeResponseState)
	assert.Nil(err)
	if s, ok := res.State.(map[string]interface{}); ok {
		assert.Equal("v", s["k"])
	} else {
		assert.Fail("!map[string]interface{}")
	}
}

func TestSetStateRequestBasic(t *testing.T) {
	assert := assert.New(t)
	state := make(map[string]interface{})
	state["name"] = "Alex"
	state["count"] = 5

	opts := newSetStateOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch"}
	opts.ChannelGroups = []string{"cg"}
	opts.State = state

	err := opts.validate()
	assert.Nil(err)

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/%s/uuid/%s/data",
			opts.Channels[0], opts.pubnub.Config.UUID),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")
	expected.Set("state", `{"count":5,"name":"Alex"}`)
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)

}

func TestSetStateMultipleChannels(t *testing.T) {
	assert := assert.New(t)

	opts := newSetStateOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroups = []string{"cg"}

	path, err := opts.buildPath()
	assert.Nil(err)

	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/uuid/%s/data",
			opts.pubnub.Config.UUID),
		u.EscapedPath(), []int{})
}

func TestSetStateMultipleChannelGroups(t *testing.T) {
	assert := assert.New(t)

	opts := newSetStateOpts(pubnub, pubnub.ctx)
	opts.ChannelGroups = []string{"cg1", "cg2", "cg3"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}

func TestSetStateMultipleChannelGroupsQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := newSetStateOpts(pubnub, pubnub.ctx)
	opts.ChannelGroups = []string{"cg1", "cg2", "cg3"}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}

func TestSetStateValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newSetStateOpts(pn, pn.ctx)
	opts.ChannelGroups = []string{"cg1", "cg2", "cg3"}

	assert.Equal("pubnub/validation: pubnub: Set State: Missing Subscribe Key", opts.validate().Error())
}

func TestSetStateValidateCG(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)

	assert.Equal("pubnub/validation: pubnub: Set State: Missing Channel or Channel Group", opts.validate().Error())
}

func TestSetStateValidateState(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)
	opts.ChannelGroups = []string{"cg1", "cg2", "cg3"}

	assert.Equal("pubnub/validation: pubnub: Set State: Missing State", opts.validate().Error())
}

func TestNewSetStateResponseErrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newSetStateResponse(jsonBytes, StatusResponse{})
	assert.Equal("pubnub/parsing: error unmarshalling response: {s}", err.Error())
}

func TestNewSetStateResponseValueError(t *testing.T) {
	assert := assert.New(t)
	state := make(map[string]interface{})
	state["name"] = "Alex"
	state["error"] = 5
	b, err1 := json.Marshal(state)
	if err1 != nil {
		panic(err1)
	}

	_, _, err := newSetStateResponse([]byte(b), StatusResponse{})
	assert.Equal("", err.Error())
}

// HTTP Method and Operation Tests

func TestSetStateHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)

	assert.Equal("GET", opts.httpMethod())
}

func TestSetStateOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)

	assert.Equal(PNSetStateOperation, opts.operationType())
}

func TestSetStateIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestSetStateTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (5 setters)

func TestSetStateBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetStateBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestSetStateBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetStateBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestSetStateBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetStateBuilder(pn)

	// Test State setter
	state := map[string]interface{}{
		"name": "Alice",
		"age":  25,
	}
	builder.State(state)
	assert.Equal(state, builder.opts.State)

	// Test Channels setter
	channels := []string{"channel1", "channel2"}
	builder.Channels(channels)
	assert.Equal(channels, builder.opts.Channels)

	// Test ChannelGroups setter
	channelGroups := []string{"group1", "group2"}
	builder.ChannelGroups(channelGroups)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)

	// Test UUID setter
	uuid := "custom-uuid"
	builder.UUID(uuid)
	assert.Equal(uuid, builder.opts.UUID)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestSetStateBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	state := map[string]interface{}{"key": "value"}
	channels := []string{"channel1"}
	channelGroups := []string{"group1"}
	uuid := "custom-uuid"
	queryParam := map[string]string{"key": "value"}

	builder := newSetStateBuilder(pn)
	result := builder.State(state).
		Channels(channels).
		ChannelGroups(channelGroups).
		UUID(uuid).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal(state, builder.opts.State)
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)
	assert.Equal(uuid, builder.opts.UUID)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestSetStateBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetStateBuilder(pn)

	// Verify default values
	assert.Nil(builder.opts.State)
	assert.Nil(builder.opts.Channels)
	assert.Nil(builder.opts.ChannelGroups)
	assert.Empty(builder.opts.UUID)
	assert.Nil(builder.opts.QueryParam)
}

func TestSetStateBuilderStateCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		state       map[string]interface{}
		description string
	}{
		{
			name:        "Simple string state",
			state:       map[string]interface{}{"status": "online"},
			description: "Set simple string state",
		},
		{
			name: "Complex nested state",
			state: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"profile": map[string]interface{}{
						"age":     25,
						"premium": true,
					},
				},
				"settings": []string{"notifications", "theme"},
			},
			description: "Set complex nested state with multiple data types",
		},
		{
			name:        "Numeric state",
			state:       map[string]interface{}{"score": 100, "level": 5, "percentage": 85.5},
			description: "Set numeric state values",
		},
		{
			name:        "Boolean state",
			state:       map[string]interface{}{"active": true, "visible": false},
			description: "Set boolean state values",
		},
		{
			name:        "Mixed data types",
			state:       map[string]interface{}{"name": "Bob", "count": 42, "enabled": true, "tags": []string{"admin", "power"}},
			description: "Set state with mixed data types",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetStateBuilder(pn)
			builder.State(tc.state)

			assert.Equal(tc.state, builder.opts.State)
		})
	}
}

func TestSetStateBuilderChannelCombinations(t *testing.T) {
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
			description: "Set state for single channel",
		},
		{
			name:        "Multiple channels",
			channels:    []string{"channel1", "channel2", "channel3"},
			description: "Set state for multiple channels",
		},
		{
			name:          "Single channel group",
			channelGroups: []string{"group1"},
			description:   "Set state for single channel group",
		},
		{
			name:          "Multiple channel groups",
			channelGroups: []string{"group1", "group2", "group3"},
			description:   "Set state for multiple channel groups",
		},
		{
			name:          "Channels and groups combination",
			channels:      []string{"channel1", "channel2"},
			channelGroups: []string{"group1", "group2"},
			description:   "Set state for both channels and channel groups",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetStateBuilder(pn)
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

func TestSetStateBuilderUUIDCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		uuid        string
		description string
	}{
		{
			name:        "Custom UUID",
			uuid:        "custom-user-123",
			description: "Set state for custom UUID",
		},
		{
			name:        "UUID with special characters",
			uuid:        "user@domain.com",
			description: "Set state for UUID with special characters",
		},
		{
			name:        "UUID with Unicode",
			uuid:        "用户123",
			description: "Set state for UUID with Unicode characters",
		},
		{
			name:        "Empty UUID (uses default)",
			uuid:        "",
			description: "Empty UUID should use default config UUID",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetStateBuilder(pn)
			builder.UUID(tc.uuid)

			assert.Equal(tc.uuid, builder.opts.UUID)
		})
	}
}

func TestSetStateBuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	state := map[string]interface{}{
		"status": "online",
		"level":  10,
	}
	channels := []string{"channel1", "channel2"}
	channelGroups := []string{"group1", "group2"}
	uuid := "custom-uuid"
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Test all 5 setters in chain
	builder := newSetStateBuilder(pn).
		State(state).
		Channels(channels).
		ChannelGroups(channelGroups).
		UUID(uuid).
		QueryParam(queryParam)

	// Verify all are set correctly
	assert.Equal(state, builder.opts.State)
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)
	assert.Equal(uuid, builder.opts.UUID)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

// URL/Path Building Tests

func TestSetStateBuildPathBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/channel/test-channel/uuid/" + pn.Config.UUID + "/data"
	assert.Equal(expected, path)
}

func TestSetStateBuildPathWithCustomUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}
	opts.UUID = "custom-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/channel/test-channel/uuid/custom-uuid/data"
	assert.Equal(expected, path)
}

func TestSetStateBuildPathMultipleChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1", "channel2", "channel3"}

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/channel/channel1,channel2,channel3/uuid/" + pn.Config.UUID + "/data"
	assert.Equal(expected, path)
}

func TestSetStateBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newSetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"my-channel"}

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/custom-sub-key/channel/my-channel/uuid/" + pn.Config.UUID + "/data"
	assert.Equal(expected, path)
}

func TestSetStateBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"channel@with#symbols"}
	opts.UUID = "user@domain.com"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/channel/channel%40with%23symbols/uuid/user%40domain.com/data"
	assert.Equal(expected, path)
}

func TestSetStateBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"频道中文"}
	opts.UUID = "用户123"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/channel/%E9%A2%91%E9%81%93%E4%B8%AD%E6%96%87/uuid/%E7%94%A8%E6%88%B7123/data"
	assert.Equal(expected, path)
}

// JSON Body Building Tests (CRITICAL for GET operation - should be empty)

func TestSetStateBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations should have empty body
	assert.Equal([]byte{}, body)
}

func TestSetStateBuildBodyWithAllParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)

	// Set all possible parameters - body should still be empty for GET
	opts.State = map[string]interface{}{"key": "value"}
	opts.Channels = []string{"channel1", "channel2"}
	opts.ChannelGroups = []string{"group1", "group2"}
	opts.UUID = "custom-uuid"
	opts.QueryParam = map[string]string{"param": "value"}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations always have empty body regardless of parameters
	assert.Equal([]byte{}, body)
}

func TestSetStateBuildBodyErrorScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)

	// Even with potential error conditions, buildBody should not fail for GET
	opts.Channels = []string{}      // Empty channels
	opts.ChannelGroups = []string{} // Empty groups

	body, err := opts.buildBody()
	assert.Nil(err) // buildBody should never error for GET operations
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}

// Query Parameter Tests

func TestSetStateBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestSetStateBuildQueryWithState(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}
	opts.State = map[string]interface{}{
		"name":   "Alice",
		"status": "online",
		"count":  42,
	}

	// Validate to serialize state
	err := opts.validate()
	assert.Nil(err)

	query, err := opts.buildQuery()
	assert.Nil(err)

	stateValue := query.Get("state")
	assert.NotEmpty(stateValue)

	// Parse the JSON to verify it contains our state
	var parsedState map[string]interface{}
	err = json.Unmarshal([]byte(stateValue), &parsedState)
	assert.Nil(err)
	assert.Equal("Alice", parsedState["name"])
	assert.Equal("online", parsedState["status"])
	assert.Equal(float64(42), parsedState["count"]) // JSON unmarshals numbers as float64
}

func TestSetStateBuildQueryWithChannelGroups(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)

	opts.ChannelGroups = []string{"group1", "group2"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	channelGroupValue := query.Get("channel-group")
	assert.Equal("group1,group2", channelGroupValue)
}

func TestSetStateBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)

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

func TestSetStateBuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetStateOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}

	// Set all possible query parameters
	opts.State = map[string]interface{}{"status": "active"}
	opts.ChannelGroups = []string{"group1", "group2"}
	opts.QueryParam = map[string]string{
		"extra": "parameter",
		"debug": "true",
	}

	// Validate to serialize state
	err := opts.validate()
	assert.Nil(err)

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all parameters are set correctly
	assert.Equal("group1,group2", query.Get("channel-group"))
	assert.NotEmpty(query.Get("state"))
	assert.Equal("parameter", query.Get("extra"))
	assert.Equal("true", query.Get("debug"))

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestSetStateBuildQueryEdgeCases(t *testing.T) {
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
			expectValue:   "group%40with%23symbols,group-with-dashes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newSetStateOpts(pn, pn.ctx)
			opts.ChannelGroups = tc.channelGroups

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expectValue, query.Get("channel-group"))
		})
	}
}

// GET-Specific Tests (State Setting Characteristics)

func TestSetStateGetOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetStateBuilder(pn)
	builder.Channels([]string{"test-channel"}).
		State(map[string]interface{}{"status": "online"})

	// Verify it's a GET operation
	assert.Equal("GET", builder.opts.httpMethod())

	// GET operations have empty body
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	// Should have proper path for state setting
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/presence/sub-key/demo/channel/test-channel")
	assert.Contains(path, "/data")
}

func TestSetStateStateSetting(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*setStateOpts)
		description string
	}{
		{
			name: "Set state for single channel",
			setupOpts: func(opts *setStateOpts) {
				opts.Channels = []string{"channel1"}
				opts.State = map[string]interface{}{"status": "online"}
			},
			description: "Set presence state for specific channel",
		},
		{
			name: "Set state for multiple channels",
			setupOpts: func(opts *setStateOpts) {
				opts.Channels = []string{"channel1", "channel2", "channel3"}
				opts.State = map[string]interface{}{"level": 5, "score": 100}
			},
			description: "Set presence state for multiple channels",
		},
		{
			name: "Set state for channel groups",
			setupOpts: func(opts *setStateOpts) {
				opts.ChannelGroups = []string{"group1", "group2"}
				opts.State = map[string]interface{}{"mode": "active"}
			},
			description: "Set presence state for channel groups",
		},
		{
			name: "Set state with custom UUID",
			setupOpts: func(opts *setStateOpts) {
				opts.Channels = []string{"channel1"}
				opts.UUID = "custom-user-id"
				opts.State = map[string]interface{}{"name": "Alice", "role": "admin"}
			},
			description: "Set presence state with custom UUID",
		},
		{
			name: "Set complex nested state",
			setupOpts: func(opts *setStateOpts) {
				opts.Channels = []string{"channel1"}
				opts.State = map[string]interface{}{
					"user": map[string]interface{}{
						"name": "Bob",
						"preferences": map[string]interface{}{
							"theme":         "dark",
							"notifications": true,
						},
					},
					"session": map[string]interface{}{
						"start_time": "2023-01-01T00:00:00Z",
						"duration":   3600,
					},
				}
			},
			description: "Set complex nested presence state",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newSetStateOpts(pn, pn.ctx)
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
			assert.Contains(path, "/data")

			// Should build valid query with state
			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
			if opts.State != nil {
				assert.NotEmpty(query.Get("state"))
			}
		})
	}
}

func TestSetStateUUIDHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name         string
		uuid         string
		expectedUUID string
		description  string
	}{
		{
			name:         "Default UUID",
			uuid:         "",
			expectedUUID: pn.Config.UUID,
			description:  "Use default config UUID when not specified",
		},
		{
			name:         "Custom UUID",
			uuid:         "custom-user-123",
			expectedUUID: "custom-user-123",
			description:  "Use custom UUID when specified",
		},
		{
			name:         "UUID with special characters",
			uuid:         "user@domain.com",
			expectedUUID: "user%40domain.com", // URL encoded
			description:  "Properly encode UUID with special characters",
		},
		{
			name:         "Unicode UUID",
			uuid:         "用户123",
			expectedUUID: "%E7%94%A8%E6%88%B7123", // URL encoded
			description:  "Properly encode Unicode UUID",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newSetStateOpts(pn, pn.ctx)
			opts.Channels = []string{"test-channel"}
			opts.UUID = tc.uuid

			path, err := opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/uuid/"+tc.expectedUUID+"/data")
		})
	}
}

func TestSetStateEmptyBodyVerification(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that GET operations always have empty body regardless of configuration
	testCases := []struct {
		name      string
		setupOpts func(*setStateOpts)
	}{
		{
			name: "With all parameters set",
			setupOpts: func(opts *setStateOpts) {
				opts.State = map[string]interface{}{"complex": "state"}
				opts.Channels = []string{"channel1", "channel2"}
				opts.ChannelGroups = []string{"group1", "group2"}
				opts.UUID = "custom-uuid"
				opts.QueryParam = map[string]string{
					"param1": "value1",
					"param2": "value2",
				}
			},
		},
		{
			name: "With minimal parameters",
			setupOpts: func(opts *setStateOpts) {
				opts.Channels = []string{"simple-channel"}
				opts.State = map[string]interface{}{"simple": "state"}
			},
		},
		{
			name: "With empty/nil parameters",
			setupOpts: func(opts *setStateOpts) {
				opts.Channels = []string{}
				opts.ChannelGroups = nil
				opts.QueryParam = nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newSetStateOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
			assert.Equal([]byte{}, body)
		})
	}
}

// Comprehensive Edge Case Tests

func TestSetStateWithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*setStateBuilder)
	}{
		{
			name: "Many channels",
			setupFn: func(builder *setStateBuilder) {
				var manyChannels []string
				for i := 0; i < 100; i++ {
					manyChannels = append(manyChannels, fmt.Sprintf("channel_%d", i))
				}
				builder.Channels(manyChannels).
					State(map[string]interface{}{"status": "active"})
			},
		},
		{
			name: "Many channel groups",
			setupFn: func(builder *setStateBuilder) {
				var manyGroups []string
				for i := 0; i < 100; i++ {
					manyGroups = append(manyGroups, fmt.Sprintf("group_%d", i))
				}
				builder.ChannelGroups(manyGroups).
					State(map[string]interface{}{"mode": "bulk"})
			},
		},
		{
			name: "Large state object",
			setupFn: func(builder *setStateBuilder) {
				largeState := make(map[string]interface{})
				for i := 0; i < 100; i++ {
					largeState[fmt.Sprintf("field_%d", i)] = fmt.Sprintf("value_%d", i)
				}
				builder.Channels([]string{"test-channel"}).
					State(largeState)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *setStateBuilder) {
				largeQueryParam := make(map[string]string)
				for i := 0; i < 100; i++ {
					largeQueryParam[fmt.Sprintf("param_%d", i)] = strings.Repeat(fmt.Sprintf("value_%d_", i), 20)
				}
				builder.Channels([]string{"test-channel"}).
					State(map[string]interface{}{"status": "active"}).
					QueryParam(largeQueryParam)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetStateBuilder(pn)
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

func TestSetStateSpecialCharacterHandling(t *testing.T) {
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
			builder := newSetStateBuilder(pn)
			builder.Channels([]string{specialString})
			builder.ChannelGroups([]string{specialString})
			builder.UUID(specialString)
			builder.State(map[string]interface{}{
				"special_field": specialString,
			})
			builder.QueryParam(map[string]string{
				"special_param": specialString,
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

func TestSetStateParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name          string
		channels      []string
		channelGroups []string
		state         map[string]interface{}
		description   string
	}{
		{
			name:        "Empty string channel",
			channels:    []string{""},
			state:       map[string]interface{}{"key": "value"},
			description: "Channel with empty string",
		},
		{
			name:        "Single character channel",
			channels:    []string{"a"},
			state:       map[string]interface{}{"key": "value"},
			description: "Channel with single character",
		},
		{
			name:        "Unicode-only channel",
			channels:    []string{"测试"},
			state:       map[string]interface{}{"key": "value"},
			description: "Channel with Unicode characters",
		},
		{
			name:          "Empty channel groups",
			channelGroups: []string{""},
			state:         map[string]interface{}{"key": "value"},
			description:   "Channel group with empty string",
		},
		{
			name:        "Empty state values",
			channels:    []string{"test"},
			state:       map[string]interface{}{"empty": "", "null": nil},
			description: "State with empty and null values",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetStateBuilder(pn)
			if tc.channels != nil {
				builder.Channels(tc.channels)
			}
			if tc.channelGroups != nil {
				builder.ChannelGroups(tc.channelGroups)
			}
			if tc.state != nil {
				builder.State(tc.state)
			}

			// Should pass validation for most cases
			if len(tc.channels) > 0 || len(tc.channelGroups) > 0 {
				if tc.state != nil {
					assert.Nil(builder.opts.validate())
				}
			}

			// Should build valid components
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/v2/presence/sub-key/")

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body) // GET operation always has empty body
		})
	}
}

func TestSetStateComplexStateScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*setStateBuilder)
		validateFn func(*testing.T, string, *url.Values)
	}{
		{
			name: "User profile state",
			setupFn: func(builder *setStateBuilder) {
				builder.Channels([]string{"user-profile-123"})
				builder.State(map[string]interface{}{
					"user": map[string]interface{}{
						"name":     "Alice Johnson",
						"avatar":   "https://example.com/avatar.jpg",
						"status":   "online",
						"lastSeen": "2023-01-01T12:00:00Z",
					},
					"preferences": map[string]interface{}{
						"theme":         "dark",
						"notifications": true,
						"language":      "en-US",
					},
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, "user-profile-123")
				stateJSON := query.Get("state")
				assert.NotEmpty(stateJSON)

				var state map[string]interface{}
				err := json.Unmarshal([]byte(stateJSON), &state)
				assert.Nil(err)
				assert.Contains(state, "user")
				assert.Contains(state, "preferences")
			},
		},
		{
			name: "Gaming session state",
			setupFn: func(builder *setStateBuilder) {
				builder.Channels([]string{"game-lobby", "game-room-5"})
				builder.ChannelGroups([]string{"gaming-groups"})
				builder.UUID("player-abc123")
				builder.State(map[string]interface{}{
					"player": map[string]interface{}{
						"level":        42,
						"score":        98765,
						"rank":         "gold",
						"achievements": []string{"first_win", "speed_demon", "champion"},
					},
					"session": map[string]interface{}{
						"start_time": "2023-01-01T20:00:00Z",
						"game_mode":  "battle_royale",
						"team_id":    "team-alpha",
					},
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, "game-lobby,game-room-5")
				assert.Contains(path, "player-abc123")
				assert.Equal("gaming-groups", query.Get("channel-group"))

				stateJSON := query.Get("state")
				assert.NotEmpty(stateJSON)
			},
		},
		{
			name: "International content state",
			setupFn: func(builder *setStateBuilder) {
				builder.Channels([]string{"频道中文123", "канал-русский"})
				builder.UUID("用户测试")
				builder.State(map[string]interface{}{
					"locale": map[string]interface{}{
						"language": "zh-CN",
						"region":   "中国",
						"timezone": "Asia/Shanghai",
					},
					"content": map[string]interface{}{
						"title":       "测试标题",
						"description": "Описание на русском языке",
						"tags":        []string{"测试", "тест", "test"},
					},
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, "%E9%A2%91%E9%81%93%E4%B8%AD%E6%96%87123,%D0%BA%D0%B0%D0%BD%D0%B0%D0%BB-%D1%80%D1%83%D1%81%D1%81%D0%BA%D0%B8%D0%B9")
				assert.Contains(path, "%E7%94%A8%E6%88%B7%E6%B5%8B%E8%AF%95") // URL-encoded 用户测试

				stateJSON := query.Get("state")
				assert.NotEmpty(stateJSON)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetStateBuilder(pn)
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

func TestSetStateExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newSetStateBuilder(pn)
	builder.Channels([]string{"test-channel"}).
		State(map[string]interface{}{"key": "value"})

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestSetStatePathBuildingEdgeCases(t *testing.T) {
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
			expectError:  false, // buildPath handles empty channels
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
			name:         "Unicode SubscribeKey, channels, and UUID",
			subscribeKey: "测试订阅键-русский-キー",
			channels:     []string{"频道测试-канал-チャンネル"},
			uuid:         "用户测试-пользователь-ユーザー",
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
			opts := newSetStateOpts(pn, pn.ctx)
			opts.Channels = tc.channels
			opts.UUID = tc.uuid

			path, err := opts.buildPath()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Contains(path, "/v2/presence/sub-key/")
				assert.Contains(path, "/data")
			}
		})
	}
}

func TestSetStateQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*setStateOpts)
		expectError bool
	}{
		{
			name: "Nil query params",
			setupOpts: func(opts *setStateOpts) {
				opts.QueryParam = nil
			},
			expectError: false,
		},
		{
			name: "Empty query params",
			setupOpts: func(opts *setStateOpts) {
				opts.QueryParam = map[string]string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *setStateOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *setStateOpts) {
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
			opts := newSetStateOpts(pn, pn.ctx)
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

func TestSetStateBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newSetStateBuilder(pn)

	state := map[string]interface{}{
		"name":   "Alice",
		"level":  10,
		"active": true,
	}
	channels := []string{"channel1", "channel2"}
	channelGroups := []string{"group1", "group2"}
	uuid := "custom-uuid"
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.State(state).
		Channels(channels).
		ChannelGroups(channelGroups).
		UUID(uuid).
		QueryParam(queryParam)

	// Verify all values are set correctly
	assert.Equal(state, builder.opts.State)
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(channelGroups, builder.opts.ChannelGroups)
	assert.Equal(uuid, builder.opts.UUID)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v2/presence/sub-key/demo/channel/channel1,channel2/uuid/custom-uuid/data"
	assert.Equal(expectedPath, path)

	// Should build query with all params
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("group1,group2", query.Get("channel-group"))
	assert.NotEmpty(query.Get("state"))
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))

	// Should always have empty body (GET operation)
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}

func TestSetStateValidationErrors(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name          string
		setupOpts     func(*setStateOpts)
		expectedError string
	}{
		{
			name: "Missing subscribe key",
			setupOpts: func(opts *setStateOpts) {
				opts.pubnub.Config.SubscribeKey = ""
				opts.Channels = []string{"test"}
				opts.State = map[string]interface{}{"key": "value"}
			},
			expectedError: "Missing Subscribe Key",
		},
		{
			name: "Missing channels and channel groups",
			setupOpts: func(opts *setStateOpts) {
				opts.State = map[string]interface{}{"key": "value"}
			},
			expectedError: "Missing Channel or Channel Group",
		},
		{
			name: "Missing state",
			setupOpts: func(opts *setStateOpts) {
				opts.Channels = []string{"test"}
			},
			expectedError: "Missing State",
		},
		{
			name: "Invalid state (channel exists)",
			setupOpts: func(opts *setStateOpts) {
				opts.Channels = []string{"test"}
				opts.State = make(map[string]interface{})
				// Create a circular reference to cause JSON marshal error
				opts.State["self"] = opts.State
			},
			expectedError: "unsupported value: encountered a cycle",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh PubNub instance for each test case to avoid shared state
			pn := NewPubNub(NewDemoConfig())
			opts := newSetStateOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			err := opts.validate()
			assert.NotNil(err)
			assert.Contains(err.Error(), tc.expectedError)
		})
	}
}

// Extended State Serialization Tests

func TestSetStateStateSerialization(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name         string
		state        map[string]interface{}
		expectedJSON string
		description  string
	}{
		{
			name:         "Simple string state",
			state:        map[string]interface{}{"status": "online"},
			expectedJSON: `{"status":"online"}`,
			description:  "Simple string value serialization",
		},
		{
			name:         "Numeric state",
			state:        map[string]interface{}{"count": 42, "score": 98.5},
			expectedJSON: `{"count":42,"score":98.5}`,
			description:  "Numeric value serialization",
		},
		{
			name:         "Boolean state",
			state:        map[string]interface{}{"active": true, "visible": false},
			expectedJSON: `{"active":true,"visible":false}`,
			description:  "Boolean value serialization",
		},
		{
			name: "Array state",
			state: map[string]interface{}{
				"tags":   []string{"admin", "premium"},
				"scores": []int{100, 95, 88},
				"flags":  []bool{true, false, true},
			},
			description: "Array value serialization",
		},
		{
			name: "Nested object state",
			state: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"profile": map[string]interface{}{
						"age":     25,
						"premium": true,
					},
				},
			},
			description: "Nested object serialization",
		},
		{
			name: "Mixed data types",
			state: map[string]interface{}{
				"name":   "Bob",
				"age":    30,
				"active": true,
				"tags":   []string{"user", "premium"},
				"metadata": map[string]interface{}{
					"created_at": "2023-01-01T00:00:00Z",
					"version":    1,
				},
				"empty_value": nil,
			},
			description: "Mixed data types serialization",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newSetStateOpts(pn, pn.ctx)
			opts.Channels = []string{"test-channel"}
			opts.State = tc.state

			err := opts.validate()
			assert.Nil(err)

			// Verify state was serialized
			assert.NotEmpty(opts.stringState)

			// Verify it's valid JSON
			var parsed map[string]interface{}
			err = json.Unmarshal([]byte(opts.stringState), &parsed)
			assert.Nil(err)

			// Verify content matches original state
			for key := range tc.state {
				assert.Contains(parsed, key)
				// Note: JSON unmarshaling may change types (e.g., all numbers become float64)
			}
		})
	}
}
