package pubnub

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/pubnub/go/v7/utils"
	"github.com/stretchr/testify/assert"
)

func AssertGetChannelMetadata(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	incl := []PNChannelMetadataInclude{
		PNChannelMetadataIncludeCustom,
	}
	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := EnumArrayToStringArray(incl)

	o := newGetChannelMetadataBuilder(pn)
	if testContext {
		o = newGetChannelMetadataBuilderWithContext(pn, pn.ctx)
	}

	o.Include(incl)
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
		assert.Equal(string(utils.JoinChannels(inclStr)), u.Get("include"))
	}

}

func TestGetChannelMetadata(t *testing.T) {
	AssertGetChannelMetadata(t, true, false)
}

func TestGetChannelMetadataContext(t *testing.T) {
	AssertGetChannelMetadata(t, true, true)
}

func TestGetChannelMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestGetChannelMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":{"id":"id0","name":"name","description":"desc","custom":{"a":"b"},"status":"active","type":"public","created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T13:26:08.341297Z","eTag":"Aee9zsKNndXlHw"}}`)

	r, _, err := newPNGetChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("id0", r.Data.ID)
	assert.Equal("name", r.Data.Name)
	assert.Equal("desc", r.Data.Description)
	//assert.Equal("2019-08-20T13:26:08.341297Z", r.Data.Created)
	assert.Equal("2019-08-20T13:26:08.341297Z", r.Data.Updated)
	assert.Equal("Aee9zsKNndXlHw", r.Data.ETag)
	assert.Equal("b", r.Data.Custom["a"])
	assert.Equal("active", r.Data.Status)
	assert.Equal("public", r.Data.Type)

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestGetChannelMetadataValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetChannelMetadataValidateMissingChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = ""

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestGetChannelMetadataValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	assert.Nil(opts.validate())
}

func TestGetChannelMetadataValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.Include = []string{"custom"}
	opts.QueryParam = map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestGetChannelMetadataHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)

	assert.Equal("GET", opts.httpMethod())
}

func TestGetChannelMetadataOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)

	assert.Equal(PNGetChannelMetadataOperation, opts.operationType())
}

func TestGetChannelMetadataIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestGetChannelMetadataTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (4 setters)

func TestGetChannelMetadataBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetChannelMetadataBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestGetChannelMetadataBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetChannelMetadataBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestGetChannelMetadataBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetChannelMetadataBuilder(pn)

	// Test Channel setter
	builder.Channel("test-channel")
	assert.Equal("test-channel", builder.opts.Channel)

	// Test Include setter
	include := []PNChannelMetadataInclude{
		PNChannelMetadataIncludeCustom,
	}
	builder.Include(include)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestGetChannelMetadataBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	include := []PNChannelMetadataInclude{PNChannelMetadataIncludeCustom}
	queryParam := map[string]string{"key": "value"}

	builder := newGetChannelMetadataBuilder(pn)
	result := builder.Channel("test-channel").
		Include(include).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal(EnumArrayToStringArray(include), builder.opts.Include)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestGetChannelMetadataBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newGetChannelMetadataBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestGetChannelMetadataBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/channels/test-channel"
	assert.Equal(expected, path)
}

func TestGetChannelMetadataBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "channel-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	// Should contain the channel name (possibly URL encoded)
	assert.Contains(path, "/v2/objects/demo/channels/")
}

func TestGetChannelMetadataBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "测试频道-русский-チャンネル"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")
}

func TestGetChannelMetadataBuildPathEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialChannels := []string{
		"channel@with%encoded",
		"channel/with/slashes",
		"channel?with=query&chars",
		"channel#with#hashes",
		"channel with spaces and símböls",
		"测试频道-русский-チャンネル-한국어",
		"channel_with_underscores",
		"channel-with-dashes",
		"channel.with.dots",
		"UPPERCASE_CHANNEL",
	}

	for _, channel := range specialChannels {
		opts := newGetChannelMetadataOpts(pn, pn.ctx)
		opts.Channel = channel

		path, err := opts.buildPath()
		assert.Nil(err, "Should build path for channel: %s", channel)
		assert.Contains(path, "/v2/objects/demo/channels/", "Should contain base path for: %s", channel)
	}
}

// Include Parameter Tests (Comma-separated enum conversion)

func TestGetChannelMetadataBuildQueryWithoutInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters but no include
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("", query.Get("include"))
}

func TestGetChannelMetadataBuildQueryWithIncludeCustom(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	opts.Include = []string{"custom"}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("custom", query.Get("include"))
}

func TestGetChannelMetadataBuildQueryWithMultipleIncludes(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	opts.Include = []string{"custom", "status", "type"}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("custom,status,type", query.Get("include"))
}

func TestGetChannelMetadataBuilderIncludeEnums(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		includes []PNChannelMetadataInclude
		expected []string
	}{
		{
			name:     "Single include custom",
			includes: []PNChannelMetadataInclude{PNChannelMetadataIncludeCustom},
			expected: []string{"custom"},
		},
		{
			name:     "Single include status",
			includes: []PNChannelMetadataInclude{PNChannelMetadataIncludeStatus},
			expected: []string{"status"},
		},
		{
			name:     "Single include type",
			includes: []PNChannelMetadataInclude{PNChannelMetadataIncludeType},
			expected: []string{"type"},
		},
		{
			name:     "Multiple includes",
			includes: []PNChannelMetadataInclude{PNChannelMetadataIncludeCustom, PNChannelMetadataIncludeStatus, PNChannelMetadataIncludeType},
			expected: []string{"custom", "status", "type"},
		},
		{
			name:     "Empty includes",
			includes: []PNChannelMetadataInclude{},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetChannelMetadataBuilder(pn)
			builder.Include(tc.includes)

			assert.Equal(tc.expected, builder.opts.Include)
		})
	}
}

func TestGetChannelMetadataBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)

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

// Comprehensive Edge Case Tests

func TestGetChannelMetadataWithUnicodeChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "测试频道-русский-チャンネル"

	// Should pass validation
	assert.Nil(opts.validate())

	// Should build path
	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")

	// Should build query
	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

func TestGetChannelMetadataWithVeryLongChannelName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create very long channel name
	longName := ""
	for i := 0; i < 100; i++ {
		longName += fmt.Sprintf("segment_%d_", i)
	}

	opts := newGetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = longName

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/")
	assert.Contains(path, "segment_0_")
	assert.Contains(path, "segment_99_")
}

func TestGetChannelMetadataSpecialCharacterHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialStrings := []string{
		"!@#$%^&*()_+-=[]{}|;':\",./<>?",
		"測試字符串-русская строка-テスト文字列",
		"<script>alert('xss')</script>",
		"SELECT * FROM users; DROP TABLE users;",
		"newline\ncharacter\ttab\rcarriage",
		"",
		"   ",
		"\u0000\u0001\u0002",
	}

	for i, specialString := range specialStrings {
		t.Run(fmt.Sprintf("SpecialString_%d", i), func(t *testing.T) {
			opts := newGetChannelMetadataOpts(pn, pn.ctx)

			// For empty channel, validation should fail
			if specialString == "" {
				opts.Channel = specialString
				err := opts.validate()
				assert.NotNil(err)
				return
			}

			opts.Channel = fmt.Sprintf("test-channel-%d", i)
			opts.QueryParam = map[string]string{
				"special_field": specialString,
			}

			// Should pass validation for non-empty channels
			assert.Nil(opts.validate())

			// Should build valid query
			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

func TestGetChannelMetadataParameterCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		channel    string
		include    []string
		queryParam map[string]string
	}{
		{
			name:    "Minimal - only channel",
			channel: "minimal-channel",
		},
		{
			name:    "With include",
			channel: "include-channel",
			include: []string{"custom"},
		},
		{
			name:    "With query params",
			channel: "query-channel",
			queryParam: map[string]string{
				"filter": "active",
				"sort":   "name",
			},
		},
		{
			name:    "Complete - all parameters",
			channel: "complete-channel",
			include: []string{"custom", "status", "type"},
			queryParam: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
		},
		{
			name:    "Unicode channel with parameters",
			channel: "测试频道-русский",
			include: []string{"custom"},
			queryParam: map[string]string{
				"language": "unicode",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetChannelMetadataBuilder(pn)
			builder.Channel(tc.channel)

			if tc.include != nil {
				builder.opts.Include = tc.include
			}
			if tc.queryParam != nil {
				builder.QueryParam(tc.queryParam)
			}

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/v2/objects/demo/channels/")

			// Should build valid query
			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

// Error Scenario Tests

func TestGetChannelMetadataExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newGetChannelMetadataBuilder(pn)
	builder.Channel("test-channel")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetChannelMetadataExecuteErrorMissingChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetChannelMetadataBuilder(pn)
	// Don't set channel

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestGetChannelMetadataPathBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCases := []struct {
		name    string
		channel string
	}{
		{
			name:    "Empty channel name",
			channel: "",
		},
		{
			name:    "Channel with only spaces",
			channel: "   ",
		},
		{
			name:    "Very special characters",
			channel: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name:    "Unicode channel",
			channel: "测试频道-русский-チャンネル-한국어",
		},
		{
			name:    "URL-like channel",
			channel: "https://example.com/channel?param=value",
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetChannelMetadataOpts(pn, pn.ctx)
			opts.Channel = tc.channel

			// For empty channel, validation should fail
			if tc.channel == "" {
				err := opts.validate()
				assert.NotNil(err)
				return
			}

			// Should pass validation for non-empty channels
			assert.Nil(opts.validate(), "Should validate for case: %s", tc.name)

			// Should build valid path
			path, err := opts.buildPath()
			assert.Nil(err, "Should build path for case: %s", tc.name)
			assert.Contains(path, "/v2/objects/", "Should contain base path for: %s", tc.name)
		})
	}
}

func TestGetChannelMetadataQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		queryParam  map[string]string
		expectError bool
	}{
		{
			name:        "Empty query params",
			queryParam:  map[string]string{},
			expectError: false,
		},
		{
			name:        "Nil query params",
			queryParam:  nil,
			expectError: false,
		},
		{
			name: "Very large query params",
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
			opts := newGetChannelMetadataOpts(pn, pn.ctx)
			opts.Channel = "test-channel"
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

func TestGetChannelMetadataBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newGetChannelMetadataBuilder(pn)

	include := []PNChannelMetadataInclude{PNChannelMetadataIncludeCustom}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.Channel("complete-test-channel").
		Include(include).
		QueryParam(queryParam)

	// Verify all values are set
	assert.Equal("complete-test-channel", builder.opts.Channel)
	assert.Equal(EnumArrayToStringArray(include), builder.opts.Include)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v2/objects/demo/channels/complete-test-channel"
	assert.Equal(expectedPath, path)

	// Should build query with custom params and include
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))
	assert.Equal("custom", query.Get("include"))
}
