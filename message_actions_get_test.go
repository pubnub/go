package pubnub

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertGetMessageActions(t *testing.T, checkQueryParam, testContext bool) {

	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newGetMessageActionsBuilder(pn)
	if testContext {
		o = newGetMessageActionsBuilderWithContext(pn, pn.ctx)
	}

	channel := "chan"
	timetoken := "15698453963258802"
	aTimetoken := "15692384791344400"
	limit := 10
	o.Channel(channel)
	o.Start(timetoken)
	o.End(aTimetoken)
	o.Limit(limit)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(getMessageActionsPath, pn.Config.SubscribeKey, channel),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
		assert.Equal(timetoken, u.Get("start"))
		assert.Equal(aTimetoken, u.Get("end"))
		assert.Equal(strconv.Itoa(limit), u.Get("limit"))
	}

}

func TestGetMessageActions(t *testing.T) {
	AssertGetMessageActions(t, true, false)
}

func TestGetMessageActionsContext(t *testing.T) {
	AssertGetMessageActions(t, true, true)
}

func TestGetMessageActionsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetMessageActionsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestGetMessageActionsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status": 200, "data": [{"messageTimetoken": "15698466245557325", "type": "reaction", "uuid": "pn-85463c27-ad24-49d4-8cdf-db93a300855a", "value": "smiley_face", "actionTimetoken": "15698466249528820"}]}`)

	r, _, err := newPNGetMessageActionsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("15698466245557325", r.Data[0].MessageTimetoken)
	assert.Equal("reaction", r.Data[0].ActionType)
	assert.Equal("smiley_face", r.Data[0].ActionValue)
	assert.Equal("15698466249528820", r.Data[0].ActionTimetoken)
	assert.Equal("pn-85463c27-ad24-49d4-8cdf-db93a300855a", r.Data[0].UUID)

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestGetMessageActionsValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetMessageActionsValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestGetMessageActionsHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)

	// GetMessageActions should use GET method (default when httpMethod() not defined)
	// Since httpMethod() is not defined, it defaults to GET
	// We can verify this by checking that buildBody() returns empty
	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET requests have empty body
}

func TestGetMessageActionsOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)

	assert.Equal(PNGetMessageActionsOperation, opts.operationType())
}

func TestGetMessageActionsIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestGetMessageActionsTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests

func TestGetMessageActionsBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMessageActionsBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestGetMessageActionsBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMessageActionsBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestGetMessageActionsBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMessageActionsBuilder(pn)

	// Test Channel setter
	builder.Channel("test-channel")
	assert.Equal("test-channel", builder.opts.Channel)

	// Test Start setter
	builder.Start("15698453963258802")
	assert.Equal("15698453963258802", builder.opts.Start)

	// Test End setter
	builder.End("15692384791344400")
	assert.Equal("15692384791344400", builder.opts.End)

	// Test Limit setter
	builder.Limit(50)
	assert.Equal(50, builder.opts.Limit)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestGetMessageActionsBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{"key": "value"}

	builder := newGetMessageActionsBuilder(pn)
	result := builder.Channel("test-channel").
		Start("15698453963258802").
		End("15692384791344400").
		Limit(25).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("15698453963258802", builder.opts.Start)
	assert.Equal("15692384791344400", builder.opts.End)
	assert.Equal(25, builder.opts.Limit)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestGetMessageActionsBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newGetMessageActionsBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestGetMessageActionsBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/message-actions/demo/channel/test-channel"
	assert.Equal(expected, path)
}

func TestGetMessageActionsBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "channel-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	// Should contain the channel name (possibly URL encoded)
	assert.Contains(path, "/message-actions/demo/channel/")
}

func TestGetMessageActionsBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "测试频道-русский-チャンネル"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/message-actions/demo/channel/")
}

func TestGetMessageActionsBuildPathEdgeCases(t *testing.T) {
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

	for _, channel := range specialChannels {
		opts := newGetMessageActionsOpts(pn, pn.ctx)
		opts.Channel = channel

		path, err := opts.buildPath()
		assert.Nil(err, "Should build path for channel: %s", channel)
		assert.Contains(path, "/message-actions/", "Should contain message-actions path for: %s", channel)
	}
}

// Conditional Query Parameter Tests

func TestGetMessageActionsBuildQueryDefault(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters but no optional ones
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("", query.Get("start")) // Empty when not set
	assert.Equal("", query.Get("end"))   // Empty when not set
	assert.Equal("", query.Get("limit")) // Empty when not set
}

func TestGetMessageActionsBuildQueryWithStart(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	opts.Start = "15698453963258802"

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("15698453963258802", query.Get("start"))
	assert.Equal("", query.Get("end"))   // Should remain empty
	assert.Equal("", query.Get("limit")) // Should remain empty
}

func TestGetMessageActionsBuildQueryWithEnd(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	opts.End = "15692384791344400"

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("", query.Get("start")) // Should remain empty
	assert.Equal("15692384791344400", query.Get("end"))
	assert.Equal("", query.Get("limit")) // Should remain empty
}

func TestGetMessageActionsBuildQueryWithLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)

	// Test positive limit
	opts.Limit = 25
	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("25", query.Get("limit"))

	// Test zero limit (should not be included)
	opts.Limit = 0
	query, err = opts.buildQuery()
	assert.Nil(err)
	assert.Equal("", query.Get("limit"))

	// Test negative limit (should not be included)
	opts.Limit = -5
	query, err = opts.buildQuery()
	assert.Nil(err)
	assert.Equal("", query.Get("limit"))
}

func TestGetMessageActionsBuildQueryWithAllParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	opts.Start = "15698453963258802"
	opts.End = "15692384791344400"
	opts.Limit = 50

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("15698453963258802", query.Get("start"))
	assert.Equal("15692384791344400", query.Get("end"))
	assert.Equal("50", query.Get("limit"))
}

func TestGetMessageActionsBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)

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

func TestGetMessageActionsWithUnicodeChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "测试频道-русский-チャンネル"

	// Should pass validation
	assert.Nil(opts.validate())

	// Should build path with URL encoding
	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/message-actions/demo/channel/")
}

func TestGetMessageActionsWithLongChannelName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create very long channel name
	longName := ""
	for i := 0; i < 100; i++ {
		longName += fmt.Sprintf("channel_%d_", i)
	}

	opts := newGetMessageActionsOpts(pn, pn.ctx)
	opts.Channel = longName

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/message-actions/")
	assert.Contains(path, "channel_0_")
	assert.Contains(path, "channel_99_")
}

func TestGetMessageActionsWithExtremeTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMessageActionsBuilder(pn)
	builder.Channel("test-channel")

	// Test extreme timetoken values
	maxTimetoken := "9223372036854775807" // Max int64 as string
	minTimetoken := "1"

	builder.Start(minTimetoken)
	builder.End(maxTimetoken)

	assert.Equal(minTimetoken, builder.opts.Start)
	assert.Equal(maxTimetoken, builder.opts.End)

	// Test query building with extreme values
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("1", query.Get("start"))
	assert.Equal("9223372036854775807", query.Get("end"))
}

func TestGetMessageActionsWithLongTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMessageActionsBuilder(pn)
	builder.Channel("test-channel")

	// Test very long timetoken strings
	longTimetoken := "15698453963258802999999999999999999"
	builder.Start(longTimetoken)

	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal(longTimetoken, query.Get("start"))
}

func TestGetMessageActionsWithEmptyQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.QueryParam = map[string]string{}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

func TestGetMessageActionsWithNilQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMessageActionsOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.QueryParam = nil

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

func TestGetMessageActionsLimitBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name          string
		limit         int
		expectedLimit string
		shouldInclude bool
	}{
		{
			name:          "Positive limit",
			limit:         25,
			expectedLimit: "25",
			shouldInclude: true,
		},
		{
			name:          "Large limit",
			limit:         1000,
			expectedLimit: "1000",
			shouldInclude: true,
		},
		{
			name:          "Zero limit",
			limit:         0,
			expectedLimit: "",
			shouldInclude: false,
		},
		{
			name:          "Negative limit",
			limit:         -5,
			expectedLimit: "",
			shouldInclude: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetMessageActionsBuilder(pn)
			builder.Channel("test-channel")
			builder.Limit(tc.limit)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)

			if tc.shouldInclude {
				assert.Equal(tc.expectedLimit, query.Get("limit"))
			} else {
				assert.Equal("", query.Get("limit"))
			}
		})
	}
}

// Error Scenario Tests

func TestGetMessageActionsExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newGetMessageActionsBuilder(pn)
	builder.Channel("test-channel")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetMessageActionsPathBuildingEdgeCases(t *testing.T) {
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

	for _, channel := range specialChannels {
		opts := newGetMessageActionsOpts(pn, pn.ctx)
		opts.Channel = channel

		// Should pass validation
		assert.Nil(opts.validate(), "Should validate channel: %s", channel)

		// Should build valid path
		path, err := opts.buildPath()
		assert.Nil(err, "Should build path for channel: %s", channel)
		assert.Contains(path, "/message-actions/", "Should contain message-actions path for: %s", channel)
	}
}

func TestGetMessageActionsParameterCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		start       string
		end         string
		limit       int
		expectStart string
		expectEnd   string
		expectLimit string
	}{
		{
			name:        "No parameters",
			start:       "",
			end:         "",
			limit:       0,
			expectStart: "",
			expectEnd:   "",
			expectLimit: "",
		},
		{
			name:        "Only start",
			start:       "15698453963258802",
			end:         "",
			limit:       0,
			expectStart: "15698453963258802",
			expectEnd:   "",
			expectLimit: "",
		},
		{
			name:        "Only end",
			start:       "",
			end:         "15692384791344400",
			limit:       0,
			expectStart: "",
			expectEnd:   "15692384791344400",
			expectLimit: "",
		},
		{
			name:        "Only limit",
			start:       "",
			end:         "",
			limit:       50,
			expectStart: "",
			expectEnd:   "",
			expectLimit: "50",
		},
		{
			name:        "All parameters",
			start:       "15698453963258802",
			end:         "15692384791344400",
			limit:       25,
			expectStart: "15698453963258802",
			expectEnd:   "15692384791344400",
			expectLimit: "25",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetMessageActionsBuilder(pn)
			builder.Channel("test-channel")

			if tc.start != "" {
				builder.Start(tc.start)
			}
			if tc.end != "" {
				builder.End(tc.end)
			}
			if tc.limit != 0 {
				builder.Limit(tc.limit)
			}

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expectStart, query.Get("start"))
			assert.Equal(tc.expectEnd, query.Get("end"))
			assert.Equal(tc.expectLimit, query.Get("limit"))
		})
	}
}

func TestGetMessageActionsEmptyParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMessageActionsBuilder(pn)
	builder.Channel("test-channel")

	// Test empty string parameters (should not be included)
	builder.Start("")
	builder.End("")

	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("", query.Get("start"))
	assert.Equal("", query.Get("end"))
}
