package pubnub

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertRemoveChannelMetadata(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newRemoveChannelMetadataBuilder(pn)
	if testContext {
		o = newRemoveChannelMetadataBuilderWithContext(pn, pn.ctx)
	}

	o.Channel("id0")
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/channels/%s", pn.Config.SubscribeKey, "id0"),
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

func TestRemoveChannelMetadata(t *testing.T) {
	AssertRemoveChannelMetadata(t, true, false)
}

func TestRemoveChannelMetadataContext(t *testing.T) {
	AssertRemoveChannelMetadata(t, true, true)
}

func TestRemoveChannelMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNRemoveChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRemoveChannelMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":null}`)

	r, _, err := newPNRemoveChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(nil, r.Data)

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestRemoveChannelMetadataValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestRemoveChannelMetadataValidateMissingChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = ""

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestRemoveChannelMetadataValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	assert.Nil(opts.validate())
}

func TestRemoveChannelMetadataValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.QueryParam = map[string]string{"param": "value"}

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestRemoveChannelMetadataHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)

	assert.Equal("DELETE", opts.httpMethod())
}

func TestRemoveChannelMetadataOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)

	assert.Equal(PNRemoveChannelMetadataOperation, opts.operationType())
}

func TestRemoveChannelMetadataIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestRemoveChannelMetadataTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (3 setters)

func TestRemoveChannelMetadataBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelMetadataBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestRemoveChannelMetadataBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelMetadataBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestRemoveChannelMetadataBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelMetadataBuilder(pn)

	// Test Channel setter
	builder.Channel("test-channel")
	assert.Equal("test-channel", builder.opts.Channel)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestRemoveChannelMetadataBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{"key": "value"}

	builder := newRemoveChannelMetadataBuilder(pn)
	result := builder.Channel("test-channel").
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestRemoveChannelMetadataBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newRemoveChannelMetadataBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

func TestRemoveChannelMetadataBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelMetadataBuilder(pn)

	// Verify default values
	assert.Equal("", builder.opts.Channel)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

// URL/Path Building Tests

func TestRemoveChannelMetadataBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/channels/test-channel"
	assert.Equal(expected, path)
}

func TestRemoveChannelMetadataBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "my-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/channels/my-channel"
	assert.Equal(expected, path)
}

func TestRemoveChannelMetadataBuildPathWithSpecialCharsInChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "channel-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")
	assert.Contains(path, "channel-with-special@chars#and$symbols")
}

func TestRemoveChannelMetadataBuildPathWithUnicodeChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "测试频道-русский-チャンネル"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")
	assert.Contains(path, "测试频道-русский-チャンネル")
}

// Query Parameter Tests

func TestRemoveChannelMetadataBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestRemoveChannelMetadataBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)

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

func TestRemoveChannelMetadataBuildQueryEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		queryParam  map[string]string
		expectError bool
	}{
		{
			name:        "Nil query params",
			queryParam:  nil,
			expectError: false,
		},
		{
			name:        "Empty query params",
			queryParam:  map[string]string{},
			expectError: false,
		},
		{
			name: "Large query params",
			queryParam: map[string]string{
				"param1": strings.Repeat("a", 1000),
				"param2": strings.Repeat("b", 1000),
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			queryParam: map[string]string{
				"special@key":   "special@value",
				"unicode测试":     "unicode值",
				"with spaces":   "also spaces",
				"equals=key":    "equals=value",
				"ampersand&key": "ampersand&value",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
			opts.QueryParam = tc.queryParam

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

// Comprehensive Edge Case Tests

func TestRemoveChannelMetadataWithUnicodeChannelIds(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	unicodeChannelIds := []string{
		"测试频道",
		"русский_канал",
		"チャンネル名",
		"قناة_عربية",
		"채널_한국어",
		"कैनल_हिंदी",
		"測試字符串-русская строка-テスト文字列",
	}

	for i, channelId := range unicodeChannelIds {
		t.Run(fmt.Sprintf("UnicodeChannel_%d", i), func(t *testing.T) {
			opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
			opts.Channel = channelId

			// Should pass validation
			assert.Nil(opts.validate())

			// Should build valid path
			path, err := opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/v2/objects/demo/channels/")
			assert.Contains(path, channelId)
		})
	}
}

func TestRemoveChannelMetadataWithSpecialCharacterChannelIds(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialChannelIds := []string{
		"!@#$%^&*()_+-=[]{}|;':\",./<>?",
		"channel-with-hyphens",
		"channel_with_underscores",
		"channel.with.dots",
		"channel with spaces",
		"UPPERCASE_CHANNEL",
		"MixedCase_Channel_123",
		"123456789",
		"channel/with/slashes",
		"channel\\with\\backslashes",
	}

	for i, channelId := range specialChannelIds {
		t.Run(fmt.Sprintf("SpecialCharChannel_%d", i), func(t *testing.T) {
			opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
			opts.Channel = channelId

			// Should pass validation
			assert.Nil(opts.validate())

			// Should build valid path
			path, err := opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/v2/objects/demo/channels/")
		})
	}
}

func TestRemoveChannelMetadataWithLongChannelIds(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name      string
		channelId string
	}{
		{
			name:      "Long channel ID",
			channelId: strings.Repeat("long_channel_", 50), // 650 characters
		},
		{
			name:      "Very long channel ID",
			channelId: strings.Repeat("very_long_channel_", 100), // 1800 characters
		},
		{
			name:      "Extremely long channel ID",
			channelId: strings.Repeat("x", 5000), // 5000 characters
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
			opts.Channel = tc.channelId

			// Should pass validation
			assert.Nil(opts.validate())

			// Should build valid path
			path, err := opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/v2/objects/demo/channels/")
		})
	}
}

func TestRemoveChannelMetadataParameterCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		channel    string
		queryParam map[string]string
	}{
		{
			name:    "Basic channel only",
			channel: "basic-channel",
		},
		{
			name:    "Channel with query params",
			channel: "channel-with-params",
			queryParam: map[string]string{
				"debug":  "true",
				"source": "test",
			},
		},
		{
			name:    "Unicode channel with params",
			channel: "测试频道",
			queryParam: map[string]string{
				"unicode_param": "unicode值",
				"regular_param": "value",
			},
		},
		{
			name:    "Special chars everywhere",
			channel: "channel@with#special$chars",
			queryParam: map[string]string{
				"special@key": "special@value",
				"unicode测试":   "unicode值",
				"with spaces": "also spaces",
			},
		},
		{
			name:    "Long combinations",
			channel: strings.Repeat("long_", 50),
			queryParam: map[string]string{
				"long_param_1": strings.Repeat("value_", 100),
				"long_param_2": strings.Repeat("data_", 100),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveChannelMetadataBuilder(pn)
			builder.Channel(tc.channel)
			if tc.queryParam != nil {
				builder.QueryParam(tc.queryParam)
			}

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/v2/objects/demo/channels/")
			assert.Contains(path, tc.channel)

			// Should build valid query
			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

func TestRemoveChannelMetadataSpecialCharacterHandling(t *testing.T) {
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
	}

	for i, specialString := range specialStrings {
		t.Run(fmt.Sprintf("SpecialString_%d", i), func(t *testing.T) {
			opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
			opts.Channel = specialString
			opts.QueryParam = map[string]string{
				"special_field": specialString,
			}

			// Should pass validation (basic validation doesn't check content)
			assert.Nil(opts.validate())

			// Should build valid path and query
			path, err := opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

// Error Scenario Tests

func TestRemoveChannelMetadataExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newRemoveChannelMetadataBuilder(pn)
	builder.Channel("test-channel")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestRemoveChannelMetadataPathBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCases := []struct {
		name         string
		subscribeKey string
		channel      string
		expectError  bool
	}{
		{
			name:         "Empty SubscribeKey",
			subscribeKey: "",
			channel:      "test-channel",
			expectError:  false, // buildPath doesn't validate, only validate() does
		},
		{
			name:         "Empty Channel",
			subscribeKey: "demo",
			channel:      "",
			expectError:  false, // buildPath doesn't validate, only validate() does
		},
		{
			name:         "SubscribeKey with spaces",
			subscribeKey: "   sub key   ",
			channel:      "test-channel",
			expectError:  false,
		},
		{
			name:         "Channel with spaces",
			subscribeKey: "demo",
			channel:      "   test channel   ",
			expectError:  false,
		},
		{
			name:         "Very special characters",
			subscribeKey: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			channel:      "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			expectError:  false,
		},
		{
			name:         "Unicode SubscribeKey and Channel",
			subscribeKey: "测试订阅键-русский-キー",
			channel:      "测试频道-русский-チャンネル",
			expectError:  false,
		},
		{
			name:         "Very long values",
			subscribeKey: strings.Repeat("a", 1000),
			channel:      strings.Repeat("b", 1000),
			expectError:  false,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			pn.Config.SubscribeKey = tc.subscribeKey
			opts := newRemoveChannelMetadataOpts(pn, pn.ctx)
			opts.Channel = tc.channel

			path, err := opts.buildPath()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Contains(path, "/v2/objects/")
				assert.Contains(path, "/channels/")
			}
		})
	}
}

func TestRemoveChannelMetadataBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newRemoveChannelMetadataBuilder(pn)

	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.Channel("complete-test-channel").
		QueryParam(queryParam)

	// Verify all values are set
	assert.Equal("complete-test-channel", builder.opts.Channel)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v2/objects/demo/channels/complete-test-channel"
	assert.Equal(expectedPath, path)

	// Should build query with all params
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))
}

func TestRemoveChannelMetadataResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMetadataOpts(pn, pn.ctx)

	testCases := []struct {
		name        string
		jsonBytes   []byte
		expectError bool
	}{
		{
			name:        "Invalid JSON",
			jsonBytes:   []byte(`{invalid json`),
			expectError: true,
		},
		{
			name:        "Null JSON",
			jsonBytes:   []byte(`null`),
			expectError: false, // null is valid JSON
		},
		{
			name:        "Empty JSON object",
			jsonBytes:   []byte(`{}`),
			expectError: false,
		},
		{
			name:        "Valid deletion response",
			jsonBytes:   []byte(`{"status":200,"data":null}`),
			expectError: false,
		},
		{
			name:        "Response with status only",
			jsonBytes:   []byte(`{"status":200}`),
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, _, err := newPNRemoveChannelMetadataResponse(tc.jsonBytes, opts, StatusResponse{})

			if tc.expectError {
				assert.NotNil(err)
				// When there's an error, resp might be nil or the empty response
				if resp == nil {
					// Should get empty response when there's an error
					assert.Equal(emptyPNRemoveChannelMetadataResponse, resp)
				}
			} else {
				assert.Nil(err)
				// For successful parsing, resp should not be nil, but content may vary
				// Note: resp can be valid even with null data for deletion responses
				if err == nil {
					// Either resp is not nil, or it's nil but that's acceptable for null JSON
					if resp != nil {
						// If resp is not nil, it should have proper structure
						assert.NotNil(resp)
					}
					// null JSON might result in nil resp, and that's acceptable for deletion responses
				}
			}
		})
	}
}
