package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/pubnub/go/v7/utils"
	"github.com/stretchr/testify/assert"
)

func AssertSetChannelMetadata(t *testing.T, checkQueryParam, testContext bool) {
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

	o := newSetChannelMetadataBuilder(pn)
	if testContext {
		o = newSetChannelMetadataBuilderWithContext(pn, pn.ctx)
	}

	o.Include(incl)
	o.Channel("id0")
	o.Name("name")
	o.Description("exturl")
	o.Custom(custom)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/channels/%s", pn.Config.SubscribeKey, "id0"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	expectedBody := "{\"name\":\"name\",\"description\":\"exturl\",\"custom\":{\"a\":\"b\",\"c\":\"d\"}}"

	assert.Equal(expectedBody, string(body))

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
		assert.Equal(string(utils.JoinChannels(inclStr)), u.Get("include"))
	}

}

func TestExcludeInChannelMetadataBodyNotSetFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newSetChannelMetadataBuilder(pn)

	o.Channel("id0")
	o.Name("name")

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/channels/%s", pn.Config.SubscribeKey, "id0"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	expectedBody := "{\"name\":\"name\"}"

	assert.Equal(expectedBody, string(body))
}

func TestSetChannelMetadata(t *testing.T) {
	AssertSetChannelMetadata(t, true, false)
}

func TestSetChannelMetadataContext(t *testing.T) {
	AssertSetChannelMetadata(t, true, true)
}

func TestSetChannelMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNSetChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestSetChannelMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":{"id":"id0","name":"name","description":"desc","custom":{"a":"b","c":"d"},"created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T14:48:11.675743Z","eTag":"AYKH2s7ZlYKoJA"}}`)

	r, _, err := newPNSetChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("id0", r.Data.ID)
	assert.Equal("name", r.Data.Name)
	assert.Equal("desc", r.Data.Description)
	// assert.Equal("2019-08-20T13:26:08.341297Z", r.Data.Created)
	assert.Equal("2019-08-20T14:48:11.675743Z", r.Data.Updated)
	assert.Equal("AYKH2s7ZlYKoJA", r.Data.ETag)
	assert.Equal("b", r.Data.Custom["a"])

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestSetChannelMetadataValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestSetChannelMetadataValidateMissingChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = ""

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestSetChannelMetadataValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	assert.Nil(opts.validate())
}

func TestSetChannelMetadataValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.Name = "Test Channel"
	opts.Description = "A test channel for validation"
	opts.Custom = map[string]interface{}{
		"type":     "test",
		"priority": 1,
	}

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestSetChannelMetadataHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)

	assert.Equal("PATCH", opts.httpMethod())
}

func TestSetChannelMetadataOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)

	assert.Equal(PNSetChannelMetadataOperation, opts.operationType())
}

func TestSetChannelMetadataIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestSetChannelMetadataTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (7 setters)

func TestSetChannelMetadataBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetChannelMetadataBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestSetChannelMetadataBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetChannelMetadataBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestSetChannelMetadataBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetChannelMetadataBuilder(pn)

	// Test Channel setter
	builder.Channel("test-channel")
	assert.Equal("test-channel", builder.opts.Channel)

	// Test Name setter
	builder.Name("Test Channel Name")
	assert.Equal("Test Channel Name", builder.opts.Name)

	// Test Description setter
	builder.Description("Test Channel Description")
	assert.Equal("Test Channel Description", builder.opts.Description)

	// Test Custom setter
	custom := map[string]interface{}{
		"category": "testing",
		"priority": 1,
		"tags":     []string{"test", "channel"},
	}
	builder.Custom(custom)
	assert.Equal(custom, builder.opts.Custom)

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

func TestSetChannelMetadataBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	custom := map[string]interface{}{"type": "test"}
	include := []PNChannelMetadataInclude{PNChannelMetadataIncludeCustom}
	queryParam := map[string]string{"key": "value"}

	builder := newSetChannelMetadataBuilder(pn)
	result := builder.Channel("test-channel").
		Name("Test Name").
		Description("Test Description").
		Custom(custom).
		Include(include).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("Test Name", builder.opts.Name)
	assert.Equal("Test Description", builder.opts.Description)
	assert.Equal(custom, builder.opts.Custom)
	assert.Equal(EnumArrayToStringArray(include), builder.opts.Include)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestSetChannelMetadataBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newSetChannelMetadataBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

// Complex JSON Body Building Tests

func TestSetChannelMetadataBuildBodyMinimal(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.Name = "Test Name"

	body, err := opts.buildBody()
	assert.Nil(err)

	expectedBody := `{"name":"Test Name"}`
	assert.Equal(expectedBody, string(body))
}

func TestSetChannelMetadataBuildBodyAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.Name = "Test Name"
	opts.Description = "Test Description"
	opts.Custom = map[string]interface{}{
		"category": "test",
		"priority": 1,
	}

	body, err := opts.buildBody()
	assert.Nil(err)

	// Parse JSON to verify structure (order may vary)
	var parsedBody map[string]interface{}
	err = json.Unmarshal(body, &parsedBody)
	assert.Nil(err)

	assert.Equal("Test Name", parsedBody["name"])
	assert.Equal("Test Description", parsedBody["description"])
	assert.NotNil(parsedBody["custom"])

	customMap := parsedBody["custom"].(map[string]interface{})
	assert.Equal("test", customMap["category"])
	assert.Equal(float64(1), customMap["priority"]) // JSON numbers are float64
}

func TestSetChannelMetadataBuildBodyEmptyFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	// Don't set any optional fields

	body, err := opts.buildBody()
	assert.Nil(err)

	expectedBody := `{}`
	assert.Equal(expectedBody, string(body))
}

func TestSetChannelMetadataBuildBodyComplexCustom(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.Custom = map[string]interface{}{
		"nested": map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": "deep value",
				"array":  []interface{}{1, "two", true},
			},
		},
		"unicode":       "测试数据-русский-チャンネル",
		"special_chars": "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		"empty_string":  "",
		"null_value":    nil,
		"boolean_true":  true,
		"boolean_false": false,
		"number_int":    42,
		"number_float":  3.14159,
	}

	body, err := opts.buildBody()
	assert.Nil(err)

	// Parse and verify complex structure
	var parsedBody map[string]interface{}
	err = json.Unmarshal(body, &parsedBody)
	assert.Nil(err)

	customMap := parsedBody["custom"].(map[string]interface{})
	assert.Equal("测试数据-русский-チャンネル", customMap["unicode"])
	assert.Equal("!@#$%^&*()_+-=[]{}|;':\",./<>?", customMap["special_chars"])
	assert.Equal(true, customMap["boolean_true"])
	assert.Equal(false, customMap["boolean_false"])
	assert.Equal(float64(42), customMap["number_int"])
	assert.Equal(3.14159, customMap["number_float"])
}

func TestSetChannelMetadataBuildBodyLargeCustom(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	// Create large custom data
	largeCustom := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		largeCustom[fmt.Sprintf("field_%d", i)] = fmt.Sprintf("value_%d", i)
	}
	opts.Custom = largeCustom

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.NotEmpty(body)

	// Verify it's valid JSON
	var parsedBody map[string]interface{}
	err = json.Unmarshal(body, &parsedBody)
	assert.Nil(err)

	customMap := parsedBody["custom"].(map[string]interface{})
	assert.Equal(100, len(customMap))
	assert.Equal("value_0", customMap["field_0"])
	assert.Equal("value_99", customMap["field_99"])
}

// URL/Path Building Tests

func TestSetChannelMetadataBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/channels/test-channel"
	assert.Equal(expected, path)
}

func TestSetChannelMetadataBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "channel-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	// Should contain the channel name (possibly URL encoded)
	assert.Contains(path, "/v2/objects/demo/channels/")
}

func TestSetChannelMetadataBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "测试频道-русский-チャンネル"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")
}

func TestSetChannelMetadataBuildPathEdgeCases(t *testing.T) {
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
		opts := newSetChannelMetadataOpts(pn, pn.ctx)
		opts.Channel = channel

		path, err := opts.buildPath()
		assert.Nil(err, "Should build path for channel: %s", channel)
		assert.Contains(path, "/v2/objects/demo/channels/", "Should contain base path for: %s", channel)
	}
}

// Include Parameter Tests (Comma-separated enum conversion)

func TestSetChannelMetadataBuildQueryWithoutInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters but no include
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("", query.Get("include"))
}

func TestSetChannelMetadataBuildQueryWithIncludeCustom(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Include = []string{"custom"}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("custom", query.Get("include"))
}

func TestSetChannelMetadataBuildQueryWithMultipleIncludes(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Include = []string{"custom", "type"}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("custom,type", query.Get("include"))
}

func TestSetChannelMetadataBuilderIncludeEnums(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		includes []PNChannelMetadataInclude
		expected []string
	}{
		{
			name:     "Single include",
			includes: []PNChannelMetadataInclude{PNChannelMetadataIncludeCustom},
			expected: []string{"custom"},
		},
		{
			name:     "Multiple includes",
			includes: []PNChannelMetadataInclude{PNChannelMetadataIncludeCustom},
			expected: []string{"custom"},
		},
		{
			name:     "Empty includes",
			includes: []PNChannelMetadataInclude{},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetChannelMetadataBuilder(pn)
			builder.Include(tc.includes)

			assert.Equal(tc.expected, builder.opts.Include)
		})
	}
}

func TestSetChannelMetadataBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)

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

func TestSetChannelMetadataWithUnicodeChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = "测试频道-русский-チャンネル"
	opts.Name = "Unicode Name: 测试名称"
	opts.Description = "Unicode Description: русское описание"

	// Should pass validation
	assert.Nil(opts.validate())

	// Should build path
	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")

	// Should build body with unicode content
	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), "Unicode Name")
	assert.Contains(string(body), "Unicode Description")
}

func TestSetChannelMetadataWithVeryLongChannelName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create very long channel name
	longName := ""
	for i := 0; i < 100; i++ {
		longName += fmt.Sprintf("segment_%d_", i)
	}

	opts := newSetChannelMetadataOpts(pn, pn.ctx)
	opts.Channel = longName

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/")
	assert.Contains(path, "segment_0_")
	assert.Contains(path, "segment_99_")
}

func TestSetChannelMetadataWithLargeMetadata(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetChannelMetadataBuilder(pn)
	builder.Channel("large-metadata-channel")

	// Very long name and description
	longName := ""
	longDescription := ""
	for i := 0; i < 100; i++ {
		longName += fmt.Sprintf("Name_Segment_%d_", i)
		longDescription += fmt.Sprintf("Description_Segment_%d_", i)
	}

	builder.Name(longName)
	builder.Description(longDescription)

	// Large custom data
	largeCustom := make(map[string]interface{})
	for i := 0; i < 50; i++ {
		largeCustom[fmt.Sprintf("field_%d", i)] = map[string]interface{}{
			"nested_field": fmt.Sprintf("nested_value_%d", i),
			"array":        []interface{}{i, fmt.Sprintf("item_%d", i), i%2 == 0},
		}
	}
	builder.Custom(largeCustom)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build valid body
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.NotEmpty(body)

	// Should be valid JSON
	var parsedBody map[string]interface{}
	err = json.Unmarshal(body, &parsedBody)
	assert.Nil(err)
}

func TestSetChannelMetadataSpecialCharacterHandling(t *testing.T) {
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
			opts := newSetChannelMetadataOpts(pn, pn.ctx)
			opts.Channel = fmt.Sprintf("test-channel-%d", i)
			opts.Name = specialString
			opts.Description = specialString
			opts.Custom = map[string]interface{}{
				"special_field": specialString,
				"nested": map[string]interface{}{
					"special_nested": specialString,
				},
			}

			// Should pass validation
			assert.Nil(opts.validate())

			// Should build valid JSON body
			body, err := opts.buildBody()
			assert.Nil(err)

			// Should be valid JSON
			var parsedBody map[string]interface{}
			err = json.Unmarshal(body, &parsedBody)
			assert.Nil(err)
		})
	}
}

func TestSetChannelMetadataParameterCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		channel     string
		nameVal     string
		description string
		custom      map[string]interface{}
		include     []string
	}{
		{
			name:    "Minimal - only channel",
			channel: "minimal-channel",
		},
		{
			name:    "Name only",
			channel: "name-channel",
			nameVal: "Channel Name",
		},
		{
			name:        "Description only",
			channel:     "desc-channel",
			description: "Channel Description",
		},
		{
			name:    "Custom only",
			channel: "custom-channel",
			custom:  map[string]interface{}{"type": "test"},
		},
		{
			name:        "All metadata fields",
			channel:     "full-channel",
			nameVal:     "Full Channel",
			description: "Complete channel with all fields",
			custom: map[string]interface{}{
				"category": "complete",
				"priority": 1,
				"features": []string{"feature1", "feature2"},
			},
		},
		{
			name:    "With include",
			channel: "include-channel",
			nameVal: "Include Channel",
			include: []string{"custom"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetChannelMetadataBuilder(pn)
			builder.Channel(tc.channel)

			if tc.nameVal != "" {
				builder.Name(tc.nameVal)
			}
			if tc.description != "" {
				builder.Description(tc.description)
			}
			if tc.custom != nil {
				builder.Custom(tc.custom)
			}
			if tc.include != nil {
				builder.opts.Include = tc.include
			}

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/v2/objects/demo/channels/"+tc.channel)

			// Should build valid query
			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Should build valid body
			body, err := builder.opts.buildBody()
			assert.Nil(err)

			// Should be valid JSON
			var parsedBody map[string]interface{}
			err = json.Unmarshal(body, &parsedBody)
			assert.Nil(err)
		})
	}
}

// Error Scenario Tests

func TestSetChannelMetadataExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newSetChannelMetadataBuilder(pn)
	builder.Channel("test-channel")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestSetChannelMetadataExecuteErrorMissingChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetChannelMetadataBuilder(pn)
	// Don't set channel

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestSetChannelMetadataPathBuildingEdgeCases(t *testing.T) {
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
			opts := newSetChannelMetadataOpts(pn, pn.ctx)
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

func TestSetChannelMetadataQueryBuildingEdgeCases(t *testing.T) {
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
			opts := newSetChannelMetadataOpts(pn, pn.ctx)
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

func TestSetChannelMetadataBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newSetChannelMetadataBuilder(pn)

	include := []PNChannelMetadataInclude{PNChannelMetadataIncludeCustom}
	custom := map[string]interface{}{
		"category": "complete",
		"priority": 1,
		"metadata": map[string]interface{}{
			"nested": "value",
		},
	}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.Channel("complete-test-channel").
		Name("Complete Channel Name").
		Description("Complete channel description with all features").
		Custom(custom).
		Include(include).
		QueryParam(queryParam)

	// Verify all values are set
	assert.Equal("complete-test-channel", builder.opts.Channel)
	assert.Equal("Complete Channel Name", builder.opts.Name)
	assert.Equal("Complete channel description with all features", builder.opts.Description)
	assert.Equal(custom, builder.opts.Custom)
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

	// Should build valid JSON body
	body, err := builder.opts.buildBody()
	assert.Nil(err)

	var parsedBody map[string]interface{}
	err = json.Unmarshal(body, &parsedBody)
	assert.Nil(err)
	assert.Equal("Complete Channel Name", parsedBody["name"])
	assert.Equal("Complete channel description with all features", parsedBody["description"])
	assert.NotNil(parsedBody["custom"])
}
