package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func init() {
	pnconfig = NewConfigWithUserId(UserId(GenerateUUID()))

	pnconfig.PublishKey = "pub_key"
	pnconfig.SubscribeKey = "sub_key"
	pnconfig.SecretKey = "secret_key"

	pubnub = NewPubNub(pnconfig)
}

func TestRemoveChannelRequestBasic(t *testing.T) {
	assert := assert.New(t)

	opts := newRemoveChannelOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroup = "cg"

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v1/channel-registration/sub-key/sub_key/channel-group/cg",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("remove", "ch1,ch2,ch3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestRemoveChannelRequestBasicQueryParam(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := newRemoveChannelOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroup = "cg"
	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("remove", "ch1,ch2,ch3")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewRemoveChannelFromChannelGroupBuilder(t *testing.T) {
	assert := assert.New(t)
	o := newRemoveChannelFromChannelGroupBuilder(pubnub)
	o.ChannelGroup("cg")
	o.Channels([]string{"ch1", "ch2", "ch3"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v1/channel-registration/sub-key/sub_key/channel-group/cg",
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("remove", "ch1,ch2,ch3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewRemoveChannelFromChannelGroupBuilderContext(t *testing.T) {
	assert := assert.New(t)
	o := newRemoveChannelFromChannelGroupBuilderWithContext(pubnub, pubnub.ctx)
	o.ChannelGroup("cg")
	o.Channels([]string{"ch1", "ch2", "ch3"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v1/channel-registration/sub-key/sub_key/channel-group/cg",
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("remove", "ch1,ch2,ch3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestRemChannelsFromCGValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newRemoveChannelOpts(pn, pn.ctx)

	assert.Equal("pubnub/validation: pubnub: Remove Channel From Channel Group: Missing Subscribe Key", opts.validate().Error())
}

// Additional Validation Tests

func TestRemoveChannelsFromChannelGroupValidateMissingChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{} // Empty channels
	opts.ChannelGroup = "cg"

	assert.Equal("pubnub/validation: pubnub: Remove Channel From Channel Group: Missing Channel", opts.validate().Error())
}

func TestRemoveChannelsFromChannelGroupValidateMissingChannelsNil(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = nil // Nil channels
	opts.ChannelGroup = "cg"

	assert.Equal("pubnub/validation: pubnub: Remove Channel From Channel Group: Missing Channel", opts.validate().Error())
}

func TestRemoveChannelsFromChannelGroupValidateMissingChannelGroup(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1"}
	opts.ChannelGroup = ""

	assert.Equal("pubnub/validation: pubnub: Remove Channel From Channel Group: Missing Channel Group", opts.validate().Error())
}

func TestRemoveChannelsFromChannelGroupValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1", "ch2"}
	opts.ChannelGroup = "test-group"

	assert.Nil(opts.validate())
}

// Builder Pattern Tests

func TestRemoveChannelsFromChannelGroupBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelFromChannelGroupBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestRemoveChannelsFromChannelGroupBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelFromChannelGroupBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestRemoveChannelsFromChannelGroupBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{"key": "value"}
	channels := []string{"ch1", "ch2"}

	builder := newRemoveChannelFromChannelGroupBuilder(pn)
	result := builder.Channels(channels).ChannelGroup("test-group").QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal("test-group", builder.opts.ChannelGroup)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestRemoveChannelsFromChannelGroupBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelFromChannelGroupBuilder(pn)

	// Test Channels setter
	channels := []string{"ch1", "ch2", "ch3"}
	builder.Channels(channels)
	assert.Equal(channels, builder.opts.Channels)

	// Test ChannelGroup setter
	builder.ChannelGroup("my-group")
	assert.Equal("my-group", builder.opts.ChannelGroup)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestRemoveChannelsFromChannelGroupBuilderChannelGroup(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelFromChannelGroupBuilder(pn)
	builder.Channels([]string{"ch1"})
	builder.ChannelGroup("test-group")

	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expected := "/v1/channel-registration/sub-key/demo/channel-group/test-group"
	assert.Equal(expected, path)
}

func TestRemoveChannelsFromChannelGroupBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelFromChannelGroupBuilder(pn)
	builder.Channels([]string{"ch1"})
	builder.ChannelGroup("test-group")

	queryParam := map[string]string{
		"custom": "param",
		"test":   "value",
	}
	builder.QueryParam(queryParam)

	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("param", query.Get("custom"))
	assert.Equal("value", query.Get("test"))
	assert.Equal("ch1", query.Get("remove"))
}

// URL/Path Building Tests

func TestRemoveChannelsFromChannelGroupBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1"}
	opts.ChannelGroup = "test-group"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/channel-registration/sub-key/demo/channel-group/test-group"
	assert.Equal(expected, path)
}

func TestRemoveChannelsFromChannelGroupBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1"}
	opts.ChannelGroup = "group-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	// Should URL encode special characters
	assert.Contains(path, "group-with-special%40chars%23and%24symbols")
}

func TestRemoveChannelsFromChannelGroupBuildQuerySingleChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Should have remove parameter with single channel
	assert.Equal("ch1", query.Get("remove"))

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestRemoveChannelsFromChannelGroupBuildQueryMultipleChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Should have remove parameter with comma-separated channels
	assert.Equal("ch1,ch2,ch3", query.Get("remove"))
}

func TestRemoveChannelsFromChannelGroupBuildQueryWithParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1"}
	opts.QueryParam = map[string]string{
		"custom": "value",
		"test":   "param",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Should have custom parameters
	assert.Equal("value", query.Get("custom"))
	assert.Equal("param", query.Get("test"))

	// Should have remove parameter
	assert.Equal("ch1", query.Get("remove"))

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestRemoveChannelsFromChannelGroupBuildQueryEmptyChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Should have empty remove parameter
	assert.Equal("", query.Get("remove"))
}

// HTTP Method and Operation Tests

func TestRemoveChannelsFromChannelGroupOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)

	assert.Equal(PNRemoveChannelFromChannelGroupOperation, opts.operationType())
}

func TestRemoveChannelsFromChannelGroupIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestRemoveChannelsFromChannelGroupTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Edge Case Tests

func TestRemoveChannelsFromChannelGroupWithManyChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a large list of channels
	channels := make([]string, 100)
	for i := 0; i < 100; i++ {
		channels[i] = fmt.Sprintf("channel_%d", i)
	}

	builder := newRemoveChannelFromChannelGroupBuilder(pn)
	builder.Channels(channels)
	builder.ChannelGroup("large-group")

	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(100, len(builder.opts.Channels))

	// Test query building with many channels
	query, err := builder.opts.buildQuery()
	assert.Nil(err)

	removeParam := query.Get("remove")
	assert.NotEmpty(removeParam)
	assert.Contains(removeParam, "channel_0")
	assert.Contains(removeParam, "channel_99")
	// Should contain comma-separated channels
	assert.Contains(removeParam, ",")
}

func TestRemoveChannelsFromChannelGroupWithUnicodeChannelGroup(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1"}
	opts.ChannelGroup = "测试群组-русский-ファイル"

	// Should pass validation
	assert.Nil(opts.validate())

	// Should build path with URL encoding
	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/channel-group/")
	// Unicode should be URL encoded
	assert.Contains(path, "%")
}

func TestRemoveChannelsFromChannelGroupWithUnicodeChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"频道1", "канал2", "チャンネル3"}
	opts.ChannelGroup = "test-group"

	// Should pass validation
	assert.Nil(opts.validate())

	// Should build query with Unicode channels
	query, err := opts.buildQuery()
	assert.Nil(err)
	removeParam := query.Get("remove")
	assert.Contains(removeParam, "频道1")
	assert.Contains(removeParam, "канал2")
	assert.Contains(removeParam, "チャンネル3")
	assert.Contains(removeParam, ",")
}

func TestRemoveChannelsFromChannelGroupWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialChannelNames := []string{
		"channel_with_underscores",
		"channel-with-hyphens",
		"channel.with.dots",
		"channel with spaces",
		"channel@with#special$chars",
		"channel%already%encoded",
	}

	specialGroupNames := []string{
		"group_with_underscores",
		"group-with-hyphens",
		"group.with.dots",
		"group with spaces",
		"group@with#special$chars",
		"group%already%encoded",
	}

	for _, channelName := range specialChannelNames {
		for _, groupName := range specialGroupNames {
			opts := newRemoveChannelOpts(pn, pn.ctx)
			opts.Channels = []string{channelName}
			opts.ChannelGroup = groupName

			// Should pass validation
			assert.Nil(opts.validate(), "Should validate channel: %s, group: %s", channelName, groupName)

			// Should build valid path
			path, err := opts.buildPath()
			assert.Nil(err, "Should build path for channel: %s, group: %s", channelName, groupName)
			assert.Contains(path, "/channel-group/", "Should contain correct path for: %s, %s", channelName, groupName)
		}
	}
}

func TestRemoveChannelsFromChannelGroupWithVeryLongChannelNames(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create very long channel names
	longChannels := make([]string, 10)
	for i := 0; i < 10; i++ {
		longName := ""
		for j := 0; j < 50; j++ {
			longName += fmt.Sprintf("channel_%d_%d_", i, j)
		}
		longChannels[i] = longName
	}

	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = longChannels
	opts.ChannelGroup = "test-group"

	assert.Nil(opts.validate())

	query, err := opts.buildQuery()
	assert.Nil(err)
	removeParam := query.Get("remove")
	assert.NotEmpty(removeParam)
	assert.Contains(removeParam, ",")
}

func TestRemoveChannelsFromChannelGroupWithEmptyQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1"}
	opts.QueryParam = map[string]string{}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
	assert.Equal("ch1", query.Get("remove"))
}

func TestRemoveChannelsFromChannelGroupWithNilQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1"}
	opts.QueryParam = nil

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
	assert.Equal("ch1", query.Get("remove"))
}

func TestRemoveChannelsFromChannelGroupWithComplexQueryParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1", "ch2"}

	complexParams := map[string]string{
		"filter":         "status=active",
		"sort":           "name,created_at",
		"include":        "metadata,custom",
		"special_chars":  "value@with#symbols",
		"unicode":        "测试参数",
		"empty_value":    "",
		"number_string":  "42",
		"boolean_string": "true",
	}
	opts.QueryParam = complexParams

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all parameters are present
	for key, expectedValue := range complexParams {
		actualValue := query.Get(key)
		if key == "special_chars" {
			// Special characters should be URL encoded
			assert.Contains(actualValue, "%", "Query parameter %s should be URL encoded", key)
		} else if key == "unicode" {
			// Unicode should be URL encoded
			assert.Contains(actualValue, "%", "Query parameter %s should contain URL encoded Unicode", key)
		} else if key == "filter" {
			// Filter parameter contains = which gets URL encoded
			assert.Equal("status%3Dactive", actualValue, "Query parameter %s should be URL encoded", key)
		} else if key == "sort" {
			// Sort parameter contains , which gets URL encoded
			assert.Equal("name%2Ccreated_at", actualValue, "Query parameter %s should be URL encoded", key)
		} else if key == "include" {
			// Include parameter contains , which gets URL encoded
			assert.Equal("metadata%2Ccustom", actualValue, "Query parameter %s should be URL encoded", key)
		} else {
			assert.Equal(expectedValue, actualValue, "Query parameter %s should match", key)
		}
	}

	// Should still have remove parameter
	assert.Equal("ch1,ch2", query.Get("remove"))
}

// Response Processing Tests

func TestRemoveChannelsFromChannelGroupResponseProcessing(t *testing.T) {
	assert := assert.New(t)

	// Test empty response processing (normal case)
	jsonBytes := []byte(`{"status": 200, "message": "OK"}`)
	resp, status, err := newRemoveChannelFromChannelGroupResponse(jsonBytes, StatusResponse{})

	assert.Nil(err)
	assert.Nil(resp) // Response is expected to be nil for this endpoint
	assert.Equal(StatusResponse{}, status)
}

func TestRemoveChannelsFromChannelGroupResponseWithNull(t *testing.T) {
	assert := assert.New(t)

	// Test null response
	jsonBytes := []byte(`null`)
	resp, _, err := newRemoveChannelFromChannelGroupResponse(jsonBytes, StatusResponse{})

	assert.Nil(err)
	assert.Nil(resp) // Response is expected to be nil for this endpoint
}

func TestRemoveChannelsFromChannelGroupResponseWithEmptyJSON(t *testing.T) {
	assert := assert.New(t)

	// Test empty JSON
	jsonBytes := []byte(`{}`)
	resp, _, err := newRemoveChannelFromChannelGroupResponse(jsonBytes, StatusResponse{})

	assert.Nil(err)
	assert.Nil(resp) // Response is expected to be nil for this endpoint
}

// Error Scenario Tests

func TestRemoveChannelsFromChannelGroupExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newRemoveChannelFromChannelGroupBuilder(pn)
	builder.Channels([]string{"ch1"})
	builder.ChannelGroup("test-group")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}
