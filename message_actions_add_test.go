package pubnub

import (
	"fmt"
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
		o = newAddMessageActionsBuilderWithContext(pn, backgroundContext)
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
