package pubnub

import (
	"fmt"
	"net/http"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertAddMessageActions(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newAddMessageActionsBuilder(pn)
	if testContext {
		o = newAddMessageActionsBuilderWithContext(pn, pn.ctx)
	}

	ma := MessageAction{
		ActionType:  "action",
		ActionValue: "smiley",
	}

	channel := "chan"
	timetoken := "15698453963258802"
	o.Channel(channel)
	o.MessageTimetoken(timetoken)
	o.Action(ma)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(addMessageActionsPath, pn.Config.SubscribeKey, channel, timetoken),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	expectedBody := "{\"type\":\"action\",\"value\":\"smiley\"}"

	assert.Equal(expectedBody, string(body))

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

}

func TestAddMessageActions(t *testing.T) {
	AssertAddMessageActions(t, true, false)
}

func TestAddMessageActionsContext(t *testing.T) {
	AssertAddMessageActions(t, true, true)
}

func TestAddMessageActionsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNAddMessageActionsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestAddMessageActionsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status": 200, "data": {"messageTimetoken": "15210190573608384", "type": "reaction", "uuid": "pn-871b8325-a11f-48cb-9c15-64984790703e", "value": "smiley_face", "actionTimetoken": "15692384791344400"}}`)

	r, _, err := newPNAddMessageActionsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("15210190573608384", r.Data.MessageTimetoken)
	assert.Equal("reaction", r.Data.ActionType)
	assert.Equal("smiley_face", r.Data.ActionValue)
	assert.Equal("15692384791344400", r.Data.ActionTimetoken)
	assert.Equal("pn-871b8325-a11f-48cb-9c15-64984790703e", r.Data.UUID)

	assert.Nil(err)
}

// Enhanced Message Actions Tests for better coverage
func TestAddMessageActionsValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test missing subscribe key validation (what the validate method actually checks)
	pn.Config.SubscribeKey = ""
	opts := newAddMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.MessageTimetoken = "15698453963258802"
	opts.Action = MessageAction{ActionType: "reaction", ActionValue: "smile"}

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Subscribe Key")

	// Test valid configuration
	pn.Config.SubscribeKey = "demo"
	err = opts.validate()
	assert.Nil(err)

	// Test path building with missing channel (should still work but create bad path)
	opts.Channel = ""
	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/channel//message/") // Empty channel in path

	// Test path building with missing timetoken
	opts.Channel = "test-channel"
	opts.MessageTimetoken = ""
	path, err = opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/message/") // Empty timetoken at end
}

func TestAddMessageActionsInvalidTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newAddMessageActionsBuilder(pn)

	// Test with various invalid timetoken formats
	invalidTimetokens := []string{
		"invalid-timetoken",
		"",
		"abc123",
		"-1",
		"0",
		"not_a_number",
		"15698453963258802999999999999999999", // too long
	}

	for _, tt := range invalidTimetokens {
		o.Channel("test-channel")
		o.MessageTimetoken(tt)
		o.Action(MessageAction{ActionType: "reaction", ActionValue: "smile"})

		path, err := o.opts.buildPath()
		// Should still build path but with invalid timetoken
		assert.Nil(err)
		assert.Contains(path, tt)
	}
}

func TestAddMessageActionsActionValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newAddMessageActionsBuilder(pn)
	o.Channel("test-channel")
	o.MessageTimetoken("15698453963258802")

	// Test with empty action type
	o.Action(MessageAction{ActionType: "", ActionValue: "smile"})
	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"type":""`)

	// Test with empty action value
	o.Action(MessageAction{ActionType: "reaction", ActionValue: ""})
	body, err = o.opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"value":""`)

	// Test with special characters in action
	o.Action(MessageAction{ActionType: "custom-type_123", ActionValue: "special!@#$%^&*()"})
	body, err = o.opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), "custom-type_123")
	// JSON marshaling escapes & as \u0026
	assert.Contains(string(body), "special!@#$%^") // Check most characters work
	assert.Contains(string(body), "\\u0026")       // Check & is properly escaped

	// Test with very long action values
	longValue := "very_long_action_value_"
	for i := 0; i < 100; i++ {
		longValue += "0123456789"
	}
	o.Action(MessageAction{ActionType: "long_test", ActionValue: longValue})
	body, err = o.opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), "long_test")
}

func TestAddMessageActionsSpecialCharactersInChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newAddMessageActionsBuilder(pn)

	// Test channels with special characters
	specialChannels := []string{
		"ch-with-dash",
		"ch_with_underscore",
		"ch.with.dot",
		"ch:with:colon",
		"ch with space",
		"ch/with/slash",
		"ch%20encoded",
		"unicode-チャンネル",
	}

	for _, channel := range specialChannels {
		o.Channel(channel)
		o.MessageTimetoken("15698453963258802")
		o.Action(MessageAction{ActionType: "test", ActionValue: "value"})

		path, err := o.opts.buildPath()
		assert.Nil(err)
		// Path should contain the channel name (possibly URL encoded)
		assert.NotEmpty(path)
	}
}

func TestAddMessageActionsResponseEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)

	// Test response with missing fields
	incompleteJSON := []byte(`{"status": 200, "data": {"messageTimetoken": "15210190573608384"}}`)
	r, _, err := newPNAddMessageActionsResponse(incompleteJSON, opts, StatusResponse{})
	assert.Nil(err)
	assert.Equal("15210190573608384", r.Data.MessageTimetoken)
	assert.Empty(r.Data.ActionType)
	assert.Empty(r.Data.ActionValue)

	// Test response with null data
	nullDataJSON := []byte(`{"status": 200, "data": null}`)
	_, _, err = newPNAddMessageActionsResponse(nullDataJSON, opts, StatusResponse{})
	// Should handle null data gracefully
	assert.Nil(err)

	// Test response with extra fields
	extraFieldsJSON := []byte(`{"status": 200, "data": {"messageTimetoken": "15210190573608384", "type": "reaction", "value": "smile", "actionTimetoken": "15692384791344400", "uuid": "test-uuid", "extraField": "ignored", "anotherExtra": 123}}`)
	r, _, err = newPNAddMessageActionsResponse(extraFieldsJSON, opts, StatusResponse{})
	assert.Nil(err)
	assert.Equal("15210190573608384", r.Data.MessageTimetoken)
	assert.Equal("reaction", r.Data.ActionType)
	assert.Equal("smile", r.Data.ActionValue)
}

func TestAddMessageActionsBuilderChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test fluent interface chaining
	result := newAddMessageActionsBuilder(pn).
		Channel("test-channel").
		MessageTimetoken("15698453963258802").
		Action(MessageAction{ActionType: "reaction", ActionValue: "smile"}).
		QueryParam(map[string]string{"custom": "param"})

	assert.NotNil(result)
	assert.Equal("test-channel", result.opts.Channel)
	assert.Equal("15698453963258802", result.opts.MessageTimetoken)
	assert.Equal("reaction", result.opts.Action.ActionType)
	assert.Equal("smile", result.opts.Action.ActionValue)
	assert.NotNil(result.opts.QueryParam)
	assert.Equal("param", result.opts.QueryParam["custom"])
}

func TestAddMessageActionsEmptyQueryParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newAddMessageActionsBuilder(pn)
	o.Channel("test-channel")
	o.MessageTimetoken("15698453963258802")
	o.Action(MessageAction{ActionType: "reaction", ActionValue: "smile"})

	// Test with nil query params
	o.QueryParam(nil)
	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Test with empty query params
	o.QueryParam(map[string]string{})
	query, err = o.opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

// HTTP Method and Operation Tests

func TestAddMessageActionsHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)

	assert.Equal("POST", opts.httpMethod())
}

func TestAddMessageActionsOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)

	assert.Equal(PNAddMessageActionsOperation, opts.operationType())
}

func TestAddMessageActionsIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestAddMessageActionsTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Comprehensive URL/Path Building Tests

func TestAddMessageActionsBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.MessageTimetoken = "15698453963258802"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/message-actions/demo/channel/test-channel/message/15698453963258802"
	assert.Equal(expected, path)
}

func TestAddMessageActionsBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "channel-with-special@chars#and$symbols"
	opts.MessageTimetoken = "15698453963258802"

	path, err := opts.buildPath()
	assert.Nil(err)
	// Should contain the channel name (possibly URL encoded)
	assert.Contains(path, "/message-actions/demo/channel/")
	assert.Contains(path, "/message/15698453963258802")
}

func TestAddMessageActionsBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "测试频道-русский-チャンネル"
	opts.MessageTimetoken = "15698453963258802"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/message-actions/demo/channel/")
	assert.Contains(path, "/message/15698453963258802")
}

func TestAddMessageActionsBuildPathEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialChannels := []string{
		"channel@with%encoded",
		"channel/with/slashes",
		"channel?with=query&chars",
		"channel#with#hashes",
		"channel with spaces and símböls",
		"测试频道-русский-チャンネル-한국어",
	}

	specialTimetokens := []string{
		"15698453963258802",
		"1234567890123456789", // Very long
		"1",                   // Very short
		"0",                   // Zero
	}

	for _, channel := range specialChannels {
		for _, timetoken := range specialTimetokens {
			opts := newAddMessageActionsOpts(pn, pn.ctx)
			opts.Channel = channel
			opts.MessageTimetoken = timetoken

			path, err := opts.buildPath()
			assert.Nil(err, "Should build path for channel: %s, timetoken: %s", channel, timetoken)
			assert.Contains(path, "/message-actions/", "Should contain message-actions path for: %s", channel)
			assert.Contains(path, "/message/"+timetoken, "Should contain timetoken for: %s", timetoken)
		}
	}
}

func TestAddMessageActionsBuildQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestAddMessageActionsBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)

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

// Systematic Builder Pattern Tests

func TestAddMessageActionsBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAddMessageActionsBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestAddMessageActionsBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAddMessageActionsBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestAddMessageActionsBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAddMessageActionsBuilder(pn)

	// Test Channel setter
	builder.Channel("test-channel")
	assert.Equal("test-channel", builder.opts.Channel)

	// Test MessageTimetoken setter
	builder.MessageTimetoken("15698453963258802")
	assert.Equal("15698453963258802", builder.opts.MessageTimetoken)

	// Test Action setter
	action := MessageAction{ActionType: "reaction", ActionValue: "smile"}
	builder.Action(action)
	assert.Equal(action, builder.opts.Action)
	assert.Equal("reaction", builder.opts.Action.ActionType)
	assert.Equal("smile", builder.opts.Action.ActionValue)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestAddMessageActionsBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	action := MessageAction{ActionType: "reaction", ActionValue: "thumbs_up"}
	queryParam := map[string]string{"key": "value"}

	builder := newAddMessageActionsBuilder(pn)
	result := builder.Channel("test-channel").
		MessageTimetoken("15698453963258802").
		Action(action).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("15698453963258802", builder.opts.MessageTimetoken)
	assert.Equal(action, builder.opts.Action)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestAddMessageActionsBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newAddMessageActionsBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

// Error Scenario Tests

func TestAddMessageActionsExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newAddMessageActionsBuilder(pn)
	builder.Channel("test-channel")
	builder.MessageTimetoken("15698453963258802")
	builder.Action(MessageAction{ActionType: "reaction", ActionValue: "smile"})

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestAddMessageActionsBuildBodyJSONError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)

	// Test with a structure that would cause JSON marshaling to succeed
	// but test that valid JSON is produced
	opts.Action = MessageAction{
		ActionType:  "test-type",
		ActionValue: "test-value",
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.NotNil(body)

	// Verify it's valid JSON
	assert.Contains(string(body), `"type":"test-type"`)
	assert.Contains(string(body), `"value":"test-value"`)
}

func TestAddMessageActionsBuildBodyWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)

	// Test with special characters that need JSON escaping
	opts.Action = MessageAction{
		ActionType:  "custom\"type\\with\tspecial\nchars",
		ActionValue: "value\"with\\quotes\nand\ttabs",
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.NotNil(body)

	// Should be valid JSON with proper escaping
	bodyStr := string(body)
	assert.Contains(bodyStr, `\"`) // Escaped quotes
	assert.Contains(bodyStr, `\\`) // Escaped backslashes
	assert.Contains(bodyStr, `\n`) // Escaped newlines
	assert.Contains(bodyStr, `\t`) // Escaped tabs
}

func TestAddMessageActionsValidationSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.MessageTimetoken = "15698453963258802"
	opts.Action = MessageAction{ActionType: "reaction", ActionValue: "smile"}

	assert.Nil(opts.validate())
}

func TestAddMessageActionsValidationMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newAddMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.MessageTimetoken = "15698453963258802"
	opts.Action = MessageAction{ActionType: "reaction", ActionValue: "smile"}

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

// Edge Cases for Message Actions

func TestAddMessageActionsWithEmptyAction(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAddMessageActionsBuilder(pn)
	builder.Channel("test-channel")
	builder.MessageTimetoken("15698453963258802")
	builder.Action(MessageAction{}) // Empty action

	// Should still pass validation (no action validation in validate())
	assert.Nil(builder.opts.validate())

	// Should build body with empty values
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"type":""`)
	assert.Contains(string(body), `"value":""`)
}

func TestAddMessageActionsWithLongValues(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create very long action type and value
	longType := ""
	longValue := ""
	for i := 0; i < 1000; i++ {
		longType += "type"
		longValue += "value"
	}

	builder := newAddMessageActionsBuilder(pn)
	builder.Channel("test-channel")
	builder.MessageTimetoken("15698453963258802")
	builder.Action(MessageAction{ActionType: longType, ActionValue: longValue})

	// Should handle long values
	assert.Nil(builder.opts.validate())

	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), "type")
	assert.Contains(string(body), "value")
	assert.True(len(body) > 8000) // Should be a large JSON body
}

func TestAddMessageActionsWithUnicodeAction(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAddMessageActionsBuilder(pn)
	builder.Channel("test-channel")
	builder.MessageTimetoken("15698453963258802")
	builder.Action(MessageAction{
		ActionType:  "反应-тип-タイプ",
		ActionValue: "笑脸-улыбка-スマイル",
	})

	// Should handle Unicode
	assert.Nil(builder.opts.validate())

	body, err := builder.opts.buildBody()
	assert.Nil(err)

	// Should contain Unicode characters properly encoded in JSON
	bodyStr := string(body)
	assert.Contains(bodyStr, "反应")
	assert.Contains(bodyStr, "笑脸")
}
