package pubnub

import (
	"fmt"
	"net/http"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertRemoveMessageActions(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newRemoveMessageActionsBuilder(pn)
	if testContext {
		o = newRemoveMessageActionsBuilderWithContext(pn, pn.ctx)
	}

	channel := "chan"
	timetoken := "15698453963258802"
	aTimetoken := "15692384791344400"
	o.Channel(channel)
	o.MessageTimetoken(timetoken)
	o.ActionTimetoken(aTimetoken)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(removeMessageActionsPath, pn.Config.SubscribeKey, channel, timetoken, aTimetoken),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

}

func TestRemoveMessageActions(t *testing.T) {
	AssertRemoveMessageActions(t, true, false)
}

func TestRemoveMessageActionsContext(t *testing.T) {
	AssertRemoveMessageActions(t, true, true)
}

func TestRemoveMessageActionsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNRemoveMessageActionsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRemoveMessageActionsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status": 200, "data": {}}`)

	r, _, err := newPNRemoveMessageActionsResponse(jsonBytes, opts, StatusResponse{})
	assert.Empty(r.Data)

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestRemoveMessageActionsValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.MessageTimetoken = "15698453963258802"
	opts.ActionTimetoken = "15692384791344400"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestRemoveMessageActionsValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.MessageTimetoken = "15698453963258802"
	opts.ActionTimetoken = "15692384791344400"

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestRemoveMessageActionsHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)

	assert.Equal("DELETE", opts.httpMethod())
}

func TestRemoveMessageActionsOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)

	assert.Equal(PNRemoveMessageActionsOperation, opts.operationType())
}

func TestRemoveMessageActionsIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestRemoveMessageActionsTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

func TestRemoveMessageActionsHTTPBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)

	// DELETE requests should have empty body
	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}

// Systematic Builder Pattern Tests

func TestRemoveMessageActionsBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMessageActionsBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestRemoveMessageActionsBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMessageActionsBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestRemoveMessageActionsBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMessageActionsBuilder(pn)

	// Test Channel setter
	builder.Channel("test-channel")
	assert.Equal("test-channel", builder.opts.Channel)

	// Test MessageTimetoken setter
	builder.MessageTimetoken("15698453963258802")
	assert.Equal("15698453963258802", builder.opts.MessageTimetoken)

	// Test ActionTimetoken setter
	builder.ActionTimetoken("15692384791344400")
	assert.Equal("15692384791344400", builder.opts.ActionTimetoken)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestRemoveMessageActionsBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{"key": "value"}

	builder := newRemoveMessageActionsBuilder(pn)
	result := builder.Channel("test-channel").
		MessageTimetoken("15698453963258802").
		ActionTimetoken("15692384791344400").
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("15698453963258802", builder.opts.MessageTimetoken)
	assert.Equal("15692384791344400", builder.opts.ActionTimetoken)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestRemoveMessageActionsBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newRemoveMessageActionsBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

// Complex URL/Path Building Tests (4-Parameter Path)

func TestRemoveMessageActionsBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.MessageTimetoken = "15698453963258802"
	opts.ActionTimetoken = "15692384791344400"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/message-actions/demo/channel/test-channel/message/15698453963258802/action/15692384791344400"
	assert.Equal(expected, path)
}

func TestRemoveMessageActionsBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "channel-with-special@chars#and$symbols"
	opts.MessageTimetoken = "15698453963258802"
	opts.ActionTimetoken = "15692384791344400"

	path, err := opts.buildPath()
	assert.Nil(err)
	// Should contain the base path structure
	assert.Contains(path, "/v1/message-actions/demo/channel/")
	assert.Contains(path, "/message/15698453963258802/action/15692384791344400")
}

func TestRemoveMessageActionsBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "测试频道-русский-チャンネル"
	opts.MessageTimetoken = "15698453963258802"
	opts.ActionTimetoken = "15692384791344400"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v1/message-actions/demo/channel/")
	assert.Contains(path, "/message/15698453963258802/action/15692384791344400")
}

func TestRemoveMessageActionsBuildPathEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name             string
		channel          string
		messageTimetoken string
		actionTimetoken  string
	}{
		{
			name:             "Special chars in channel",
			channel:          "channel@with%encoded&chars",
			messageTimetoken: "15698453963258802",
			actionTimetoken:  "15692384791344400",
		},
		{
			name:             "Unicode in channel",
			channel:          "测试频道-русский-チャンネル",
			messageTimetoken: "15698453963258802",
			actionTimetoken:  "15692384791344400",
		},
		{
			name:             "Very long timetokens",
			channel:          "test-channel",
			messageTimetoken: "15698453963258802999999999999999999",
			actionTimetoken:  "15692384791344400888888888888888888",
		},
		{
			name:             "Min timetokens",
			channel:          "test-channel",
			messageTimetoken: "1",
			actionTimetoken:  "2",
		},
		{
			name:             "Max int64 timetokens",
			channel:          "test-channel",
			messageTimetoken: "9223372036854775807",
			actionTimetoken:  "9223372036854775806",
		},
		{
			name:             "Complex combination",
			channel:          "channel/with?query=params&special#chars",
			messageTimetoken: "15698453963258802999",
			actionTimetoken:  "15692384791344400888",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newRemoveMessageActionsOpts(pn, pn.ctx)
			opts.Channel = tc.channel
			opts.MessageTimetoken = tc.messageTimetoken
			opts.ActionTimetoken = tc.actionTimetoken

			path, err := opts.buildPath()
			assert.Nil(err, "Should build path for case: %s", tc.name)
			assert.Contains(path, "/v1/message-actions/", "Should contain base path for: %s", tc.name)
			assert.Contains(path, "/message/"+tc.messageTimetoken+"/action/"+tc.actionTimetoken, "Should contain timetokens for: %s", tc.name)
		})
	}
}

// Query Parameter Tests

func TestRemoveMessageActionsBuildQueryDefault(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestRemoveMessageActionsBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)

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

func TestRemoveMessageActionsBuildQueryWithEmptyParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)
	opts.QueryParam = map[string]string{}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

func TestRemoveMessageActionsBuildQueryWithNilParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)
	opts.QueryParam = nil

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

// Comprehensive Edge Case Tests

func TestRemoveMessageActionsWithUnicodeChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "测试频道-русский-チャンネル"
	opts.MessageTimetoken = "15698453963258802"
	opts.ActionTimetoken = "15692384791344400"

	// Should pass validation
	assert.Nil(opts.validate())

	// Should build path
	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v1/message-actions/demo/channel/")
}

func TestRemoveMessageActionsWithLongChannelName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create very long channel name
	longName := ""
	for i := 0; i < 100; i++ {
		longName += fmt.Sprintf("channel_%d_", i)
	}

	opts := newRemoveMessageActionsOpts(pn, pn.ctx)
	opts.Channel = longName
	opts.MessageTimetoken = "15698453963258802"
	opts.ActionTimetoken = "15692384791344400"

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v1/message-actions/")
	assert.Contains(path, "channel_0_")
	assert.Contains(path, "channel_99_")
}

func TestRemoveMessageActionsWithExtremeTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMessageActionsBuilder(pn)
	builder.Channel("test-channel")

	// Test extreme timetoken values
	maxTimetoken := "9223372036854775807" // Max int64 as string
	minTimetoken := "1"

	builder.MessageTimetoken(maxTimetoken)
	builder.ActionTimetoken(minTimetoken)

	assert.Equal(maxTimetoken, builder.opts.MessageTimetoken)
	assert.Equal(minTimetoken, builder.opts.ActionTimetoken)

	// Test path building with extreme values
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/message/9223372036854775807/action/1")
}

func TestRemoveMessageActionsWithLongTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMessageActionsBuilder(pn)
	builder.Channel("test-channel")

	// Test very long timetoken strings
	longMessageTimetoken := "15698453963258802999999999999999999"
	longActionTimetoken := "15692384791344400888888888888888888"
	builder.MessageTimetoken(longMessageTimetoken)
	builder.ActionTimetoken(longActionTimetoken)

	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/message/"+longMessageTimetoken+"/action/"+longActionTimetoken)
}

func TestRemoveMessageActionsWithSpecialCharacterTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMessageActionsBuilder(pn)
	builder.Channel("test-channel")

	// Timetokens with special characters (though unusual, should be handled)
	messageTimetoken := "15698453963258802_special"
	actionTimetoken := "15692384791344400-action"
	builder.MessageTimetoken(messageTimetoken)
	builder.ActionTimetoken(actionTimetoken)

	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/message/"+messageTimetoken+"/action/"+actionTimetoken)
}

func TestRemoveMessageActionsTimetokenCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name             string
		messageTimetoken string
		actionTimetoken  string
	}{
		{
			name:             "Normal timetokens",
			messageTimetoken: "15698453963258802",
			actionTimetoken:  "15692384791344400",
		},
		{
			name:             "Same timetokens",
			messageTimetoken: "15698453963258802",
			actionTimetoken:  "15698453963258802",
		},
		{
			name:             "Reversed order (action newer than message)",
			messageTimetoken: "15692384791344400",
			actionTimetoken:  "15698453963258802",
		},
		{
			name:             "Zero timetokens",
			messageTimetoken: "0",
			actionTimetoken:  "0",
		},
		{
			name:             "Single digit",
			messageTimetoken: "1",
			actionTimetoken:  "2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveMessageActionsBuilder(pn)
			builder.Channel("test-channel")
			builder.MessageTimetoken(tc.messageTimetoken)
			builder.ActionTimetoken(tc.actionTimetoken)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build correct path
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			expectedPath := fmt.Sprintf("/message/%s/action/%s", tc.messageTimetoken, tc.actionTimetoken)
			assert.Contains(path, expectedPath)
		})
	}
}

// Error Scenario Tests

func TestRemoveMessageActionsExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newRemoveMessageActionsBuilder(pn)
	builder.Channel("test-channel")
	builder.MessageTimetoken("15698453963258802")
	builder.ActionTimetoken("15692384791344400")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestRemoveMessageActionsPathBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCases := []struct {
		name             string
		channel          string
		messageTimetoken string
		actionTimetoken  string
	}{
		{
			name:             "Special characters in all components",
			channel:          "channel@with%encoded&chars",
			messageTimetoken: "15698453963258802#special",
			actionTimetoken:  "15692384791344400$action",
		},
		{
			name:             "Unicode in channel",
			channel:          "测试频道-русский-チャンネル-한국어",
			messageTimetoken: "15698453963258802",
			actionTimetoken:  "15692384791344400",
		},
		{
			name:             "Very long values",
			channel:          "very-long-channel-name-with-many-segments-and-characters",
			messageTimetoken: "15698453963258802999999999999999999999999999999",
			actionTimetoken:  "15692384791344400888888888888888888888888888888",
		},
		{
			name:             "Empty timetokens",
			channel:          "test-channel",
			messageTimetoken: "",
			actionTimetoken:  "",
		},
		{
			name:             "Slash characters",
			channel:          "channel/with/slashes",
			messageTimetoken: "15698453963258802",
			actionTimetoken:  "15692384791344400",
		},
		{
			name:             "Query characters",
			channel:          "channel?with=query&params",
			messageTimetoken: "15698453963258802",
			actionTimetoken:  "15692384791344400",
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newRemoveMessageActionsOpts(pn, pn.ctx)
			opts.Channel = tc.channel
			opts.MessageTimetoken = tc.messageTimetoken
			opts.ActionTimetoken = tc.actionTimetoken

			// Should pass validation
			assert.Nil(opts.validate(), "Should validate for case: %s", tc.name)

			// Should build valid path
			path, err := opts.buildPath()
			assert.Nil(err, "Should build path for case: %s", tc.name)
			assert.Contains(path, "/v1/message-actions/", "Should contain base path for: %s", tc.name)
		})
	}
}

func TestRemoveMessageActionsBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all required fields can be set
	builder := newRemoveMessageActionsBuilder(pn)

	// Set all required parameters
	builder.Channel("complete-test-channel")
	builder.MessageTimetoken("15698453963258802")
	builder.ActionTimetoken("15692384791344400")

	// Add optional parameters
	queryParam := map[string]string{
		"metadata": "test-removal",
		"source":   "automated-test",
	}
	builder.QueryParam(queryParam)

	// Verify all values are set
	assert.Equal("complete-test-channel", builder.opts.Channel)
	assert.Equal("15698453963258802", builder.opts.MessageTimetoken)
	assert.Equal("15692384791344400", builder.opts.ActionTimetoken)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v1/message-actions/demo/channel/complete-test-channel/message/15698453963258802/action/15692384791344400"
	assert.Equal(expectedPath, path)

	// Should build query with custom params
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("test-removal", query.Get("metadata"))
	assert.Equal("automated-test", query.Get("source"))
}

func TestRemoveMessageActionsEmptyParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMessageActionsBuilder(pn)

	// Test with empty string parameters
	builder.Channel("")
	builder.MessageTimetoken("")
	builder.ActionTimetoken("")

	// Should still pass validation (only SubscribeKey is required)
	assert.Nil(builder.opts.validate())

	// Should build path with empty components
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v1/message-actions/demo/channel//message//action/"
	assert.Equal(expectedPath, path)
}
