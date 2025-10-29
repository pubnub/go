package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestListAllChannelGroupRequestBasic(t *testing.T) {
	assert := assert.New(t)

	opts := newAllChannelGroupOpts(pubnub, pubnub.ctx)
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

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestListAllChannelGroupRequestBasicQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := newAllChannelGroupOpts(pubnub, pubnub.ctx)
	opts.ChannelGroup = "cg"
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewAllChannelGroupBuilder(t *testing.T) {
	assert := assert.New(t)
	o := newAllChannelGroupBuilder(pubnub)
	o.ChannelGroup("cg")

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

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewAllChannelGroupBuilderContext(t *testing.T) {
	assert := assert.New(t)
	o := newAllChannelGroupBuilderWithContext(pubnub, pubnub.ctx)
	o.ChannelGroup("cg")

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

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestListAllChannelsNewAllChannelGroupResponseErrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newAllChannelGroupResponse(jsonBytes, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestListAllChannelsValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newAllChannelGroupOpts(pn, pn.ctx)
	opts.ChannelGroup = "cg"

	assert.Equal("pubnub/validation: pubnub: List Channels In Channel Group: Missing Subscribe Key", opts.validate().Error())
}

func TestListAllChannelsValidateChannelGrp(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)

	assert.Equal("pubnub/validation: pubnub: List Channels In Channel Group: Missing Channel Group", opts.validate().Error())
}

// Additional Validation Tests

func TestListAllChannelsValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)
	opts.ChannelGroup = "test-group"

	assert.Nil(opts.validate())
}

// Builder Pattern Tests

func TestListChannelsInChannelGroupBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAllChannelGroupBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestListChannelsInChannelGroupBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAllChannelGroupBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestListChannelsInChannelGroupBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{"key": "value"}

	builder := newAllChannelGroupBuilder(pn)
	result := builder.ChannelGroup("test-group").QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-group", builder.opts.ChannelGroup)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestListChannelsInChannelGroupBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAllChannelGroupBuilder(pn)

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

func TestListChannelsInChannelGroupBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAllChannelGroupBuilder(pn)
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
}

// URL/Path Building Tests

func TestListChannelsInChannelGroupBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)
	opts.ChannelGroup = "test-group"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/channel-registration/sub-key/demo/channel-group/test-group"
	assert.Equal(expected, path)
}

func TestListChannelsInChannelGroupBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)
	opts.ChannelGroup = "group-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	// Should URL encode special characters
	assert.Contains(path, "group-with-special%40chars%23and%24symbols")
}

func TestListChannelsInChannelGroupBuildQueryEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestListChannelsInChannelGroupBuildQueryWithParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)
	opts.QueryParam = map[string]string{
		"custom": "value",
		"test":   "param",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Should have custom parameters
	assert.Equal("value", query.Get("custom"))
	assert.Equal("param", query.Get("test"))

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

// HTTP Method and Operation Tests

func TestListChannelsInChannelGroupOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)

	assert.Equal(PNChannelsForGroupOperation, opts.operationType())
}

func TestListChannelsInChannelGroupIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestListChannelsInChannelGroupTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Edge Case Tests

func TestListChannelsInChannelGroupWithUnicodeChannelGroup(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)
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

func TestListChannelsInChannelGroupWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialChars := []string{
		"group_with_underscores",
		"group-with-hyphens",
		"group.with.dots",
		"group with spaces",
		"group@with#special$chars",
		"group%already%encoded",
	}

	for _, groupName := range specialChars {
		opts := newAllChannelGroupOpts(pn, pn.ctx)
		opts.ChannelGroup = groupName

		// Should pass validation
		assert.Nil(opts.validate(), "Should validate group name: %s", groupName)

		// Should build valid path
		path, err := opts.buildPath()
		assert.Nil(err, "Should build path for group name: %s", groupName)
		assert.Contains(path, "/channel-group/", "Should contain correct path for: %s", groupName)
	}
}

func TestListChannelsInChannelGroupWithVeryLongChannelGroupName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a very long channel group name
	longName := ""
	for i := 0; i < 100; i++ {
		longName += fmt.Sprintf("group_%d_", i)
	}

	opts := newAllChannelGroupOpts(pn, pn.ctx)
	opts.ChannelGroup = longName

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/channel-group/")
}

func TestListChannelsInChannelGroupWithEmptyQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)
	opts.QueryParam = map[string]string{}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

func TestListChannelsInChannelGroupWithNilQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)
	opts.QueryParam = nil

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

func TestListChannelsInChannelGroupWithComplexQueryParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAllChannelGroupOpts(pn, pn.ctx)

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
}

// Response Processing Tests

func TestListChannelsInChannelGroupResponseWithEmptyChannels(t *testing.T) {
	assert := assert.New(t)

	jsonBytes := []byte(`{
		"status": 200,
		"payload": {
			"channels": [],
			"group": "empty-group"
		}
	}`)

	resp, status, err := newAllChannelGroupResponse(jsonBytes, StatusResponse{})
	assert.Nil(err)
	assert.NotNil(resp)
	assert.Equal("empty-group", resp.ChannelGroup)
	assert.Equal(0, len(resp.Channels))
	assert.Equal(StatusResponse{}, status)
}

func TestListChannelsInChannelGroupResponseWithMultipleChannels(t *testing.T) {
	assert := assert.New(t)

	jsonBytes := []byte(`{
		"status": 200,
		"payload": {
			"channels": ["channel1", "channel2", "channel3"],
			"group": "multi-channel-group"
		}
	}`)

	resp, _, err := newAllChannelGroupResponse(jsonBytes, StatusResponse{})
	assert.Nil(err)
	assert.NotNil(resp)
	assert.Equal("multi-channel-group", resp.ChannelGroup)
	assert.Equal(3, len(resp.Channels))
	assert.Contains(resp.Channels, "channel1")
	assert.Contains(resp.Channels, "channel2")
	assert.Contains(resp.Channels, "channel3")
}

func TestListChannelsInChannelGroupResponseMissingPayload(t *testing.T) {
	assert := assert.New(t)

	jsonBytes := []byte(`{
		"status": 200,
		"message": "OK"
	}`)

	resp, _, err := newAllChannelGroupResponse(jsonBytes, StatusResponse{})
	assert.Nil(err)
	assert.NotNil(resp)
	assert.Equal("", resp.ChannelGroup)
	assert.Equal(0, len(resp.Channels))
}

func TestListChannelsInChannelGroupResponseMalformedChannels(t *testing.T) {
	assert := assert.New(t)

	// Channels field contains non-string values
	jsonBytes := []byte(`{
		"status": 200,
		"payload": {
			"channels": ["channel1", 123, "channel2", null, "channel3"],
			"group": "test-group"
		}
	}`)

	resp, _, err := newAllChannelGroupResponse(jsonBytes, StatusResponse{})
	assert.Nil(err)
	assert.NotNil(resp)
	assert.Equal("test-group", resp.ChannelGroup)
	// Should only include valid string channels
	assert.Equal(3, len(resp.Channels))
	assert.Contains(resp.Channels, "channel1")
	assert.Contains(resp.Channels, "channel2")
	assert.Contains(resp.Channels, "channel3")
}

func TestListChannelsInChannelGroupResponseWithUnicodeChannels(t *testing.T) {
	assert := assert.New(t)

	jsonBytes := []byte(`{
		"status": 200,
		"payload": {
			"channels": ["频道1", "канал2", "チャンネル3"],
			"group": "unicode-group"
		}
	}`)

	resp, _, err := newAllChannelGroupResponse(jsonBytes, StatusResponse{})
	assert.Nil(err)
	assert.NotNil(resp)
	assert.Equal("unicode-group", resp.ChannelGroup)
	assert.Equal(3, len(resp.Channels))
	assert.Contains(resp.Channels, "频道1")
	assert.Contains(resp.Channels, "канал2")
	assert.Contains(resp.Channels, "チャンネル3")
}

// Error Scenario Tests

func TestListChannelsInChannelGroupExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newAllChannelGroupBuilder(pn)
	builder.ChannelGroup("test-group")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}
