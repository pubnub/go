package pubnub

import (
	"fmt"
	"net/http"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessMessageCountsGet(t *testing.T, expectedString string, channels []string, timetoken int64, channelsTimetoken []int64) {
	assert := assert.New(t)

	opts := newMessageCountsOpts(pubnub, pubnub.ctx)
	opts.Channels = channels
	opts.Timetoken = timetoken
	opts.ChannelsTimetoken = channelsTimetoken

	path, err := opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/message-counts/%s", expectedString),
		path, []int{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Empty(body)
}

func TestMessageCountsPath(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}
	AssertSuccessMessageCountsGet(t, "test1,test2", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsQuery(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}
	AssertSuccessMessageCountsGetQuery(t, "", "15499825804610610,15499925804610615", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsQuery2(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{}
	AssertSuccessMessageCountsGetQuery(t, "", "", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsQuery3(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}
	AssertSuccessMessageCountsGetQuery(t, "", "15499825804610610,15499925804610615", channels, 0, channelsTimetoken)
}

func AssertSuccessMessageCountsGetQuery(t *testing.T, expectedString1 string, expectedString2 string, channels []string, timetoken int64, channelsTimetoken []int64) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := newMessageCountsOpts(pubnub, pubnub.ctx)
	opts.Channels = channels
	opts.Timetoken = timetoken
	opts.ChannelsTimetoken = channelsTimetoken
	opts.QueryParam = queryParam

	u, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))

	assert.Equal(expectedString1, u.Get("timetoken"))
	assert.Equal(expectedString2, u.Get("channelsTimetoken"))

}

func AssertNewMessageCountsBuilder(t *testing.T, testQueryParam bool, testContext bool, expectedString string, expectedString1 string, expectedString2 string, channels []string, timetoken int64, channelsTimetoken []int64) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	o := newMessageCountsBuilder(pubnub)
	if testContext {
		o = newMessageCountsBuilderWithContext(pubnub, pubnub.ctx)
	}
	o.Channels(channels)
	o.Timetoken(timetoken)
	o.ChannelsTimetoken(channelsTimetoken)
	if testQueryParam {
		o.QueryParam(queryParam)
	}

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/message-counts/%s", expectedString),
		path, []int{})

	u, _ := o.opts.buildQuery()

	if testQueryParam {
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

	assert.Equal(expectedString1, u.Get("timetoken"))
	assert.Equal(expectedString2, u.Get("channelsTimetoken"))

}

func TestMessageCountsBuilder(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}
	AssertNewMessageCountsBuilder(t, false, false, "test1,test2", "", "15499825804610610,15499925804610615", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsBuilderQP(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}
	AssertNewMessageCountsBuilder(t, true, false, "test1,test2", "", "15499825804610610,15499925804610615", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsBuilderContext(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}
	AssertNewMessageCountsBuilder(t, false, true, "test1,test2", "", "15499825804610610,15499925804610615", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsBuilderContextQP(t *testing.T) {
	channels := []string{"test1", "test2"}
	channelsTimetoken := []int64{15499825804610610, 15499925804610615}

	AssertNewMessageCountsBuilder(t, true, true, "test1,test2", "", "15499825804610610,15499925804610615", channels, 15499825804610610, channelsTimetoken)
}

func TestMessageCountsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newMessageCountsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

// {"status": 200, "error": false, "error_message": "", "channels": {"my-channel1":1,"my-channel":2}}
func TestMessageCountsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status": 200, "error": false, "error_message": "", "channels": {"my-channel1":1,"my-channel":2}}`)

	res, _, err := newMessageCountsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(2, res.Channels["my-channel"])
	assert.Equal(1, res.Channels["my-channel1"])
	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestMessageCountsValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}
	opts.ChannelsTimetoken = []int64{123456789}

	assert.Equal("pubnub/validation: pubnub: No Category Matched: Missing Subscribe Key", opts.validate().Error())
}

func TestMessageCountsValidateMissingChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{} // Empty channels
	opts.ChannelsTimetoken = []int64{123456789}

	assert.Equal("pubnub/validation: pubnub: No Category Matched: Missing Channel", opts.validate().Error())
}

func TestMessageCountsValidateNilChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = nil // Nil channels
	opts.ChannelsTimetoken = []int64{123456789}

	assert.Equal("pubnub/validation: pubnub: No Category Matched: Missing Channel", opts.validate().Error())
}

func TestMessageCountsValidateMissingTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}
	opts.ChannelsTimetoken = []int64{} // Empty timetoken array
	opts.Timetoken = 0                 // Zero timetoken

	assert.Equal("pubnub/validation: pubnub: No Category Matched: Missing Channels Timetoken", opts.validate().Error())
}

func TestMessageCountsValidateChannelsTimetokenLengthMismatch(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1", "channel2", "channel3"} // 3 channels
	opts.ChannelsTimetoken = []int64{123456789, 987654321}       // 2 timetokens - mismatch!

	assert.Equal("pubnub/validation: pubnub: No Category Matched: Length of Channels Timetoken and Channels do not match", opts.validate().Error())
}

func TestMessageCountsValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}
	opts.ChannelsTimetoken = []int64{123456789}

	assert.Nil(opts.validate())
}

func TestMessageCountsValidateSuccessMultipleChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1", "channel2", "channel3"}
	opts.ChannelsTimetoken = []int64{123456789, 987654321, 555666777} // Matching lengths

	assert.Nil(opts.validate())
}

func TestMessageCountsValidateSuccessLegacyTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}
	opts.ChannelsTimetoken = nil // Nil array to trigger legacy path
	opts.Timetoken = 123456789   // Use legacy timetoken

	assert.Nil(opts.validate())
}

// Systematic Builder Pattern Tests

func TestMessageCountsBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newMessageCountsBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestMessageCountsBuilderContextNew(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newMessageCountsBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestMessageCountsBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	channels := []string{"ch1", "ch2"}
	channelsTimetoken := []int64{123456789, 987654321}
	queryParam := map[string]string{"key": "value"}

	builder := newMessageCountsBuilder(pn)
	result := builder.Channels(channels).
		ChannelsTimetoken(channelsTimetoken).
		Timetoken(555666777).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(channelsTimetoken, builder.opts.ChannelsTimetoken)
	assert.Equal(int64(555666777), builder.opts.Timetoken)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestMessageCountsBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newMessageCountsBuilder(pn)

	// Test Channels setter
	channels := []string{"channel1", "channel2", "channel3"}
	builder.Channels(channels)
	assert.Equal(channels, builder.opts.Channels)

	// Test ChannelsTimetoken setter
	channelsTimetoken := []int64{123456789, 987654321, 555666777}
	builder.ChannelsTimetoken(channelsTimetoken)
	assert.Equal(channelsTimetoken, builder.opts.ChannelsTimetoken)

	// Test Timetoken setter (deprecated)
	builder.Timetoken(999888777)
	assert.Equal(int64(999888777), builder.opts.Timetoken)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestMessageCountsBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newMessageCountsBuilder(pn)
	builder.Channels([]string{"test-channel"})
	builder.ChannelsTimetoken([]int64{123456789})

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

// Complex Query Building Tests

func TestMessageCountsQuerySingleChannelTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1"}
	opts.ChannelsTimetoken = []int64{123456789} // Single timetoken -> uses "timetoken" parameter

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("123456789", query.Get("timetoken"))
	assert.Equal("", query.Get("channelsTimetoken")) // Should be empty
}

func TestMessageCountsQueryMultipleChannelsTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1", "channel2", "channel3"}
	opts.ChannelsTimetoken = []int64{123456789, 987654321, 555666777} // Multiple -> uses "channelsTimetoken" parameter

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("", query.Get("timetoken")) // Should be empty
	assert.Equal("123456789,987654321,555666777", query.Get("channelsTimetoken"))
}

func TestMessageCountsQueryLegacyTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1"}
	opts.ChannelsTimetoken = nil // Nil array to trigger legacy path
	opts.Timetoken = 999888777   // Legacy timetoken

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("999888777", query.Get("timetoken"))
	assert.Equal("", query.Get("channelsTimetoken")) // Should be empty
}

func TestMessageCountsQueryComplexParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1", "channel2"}
	opts.ChannelsTimetoken = []int64{123456789, 987654321}

	complexParams := map[string]string{
		"filter":         "status=active",
		"sort":           "time,desc",
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
			assert.Equal("time%2Cdesc", actualValue, "Query parameter %s should be URL encoded", key)
		} else {
			assert.Equal(expectedValue, actualValue, "Query parameter %s should match", key)
		}
	}

	// Verify channelsTimetoken parameter
	assert.Equal("123456789,987654321", query.Get("channelsTimetoken"))
}

// URL/Path Building Tests

func TestMessageCountsBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v3/history/sub-key/demo/message-counts/test-channel"
	assert.Equal(expected, path)
}

func TestMessageCountsBuildPathMultipleChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1", "channel2", "channel3"}

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v3/history/sub-key/demo/message-counts/channel1,channel2,channel3"
	assert.Equal(expected, path)
}

func TestMessageCountsBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"channel-with-special@chars#and$symbols", "channel2"}

	path, err := opts.buildPath()
	assert.Nil(err)
	// Should URL encode special characters in channel names
	assert.Contains(path, "/message-counts/")
	assert.Contains(path, "channel-with-special%40chars%23and%24symbols")
}

// HTTP Method and Operation Tests

func TestMessageCountsOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)

	assert.Equal(PNMessageCountsOperation, opts.operationType())
}

func TestMessageCountsIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestMessageCountsTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Comprehensive Edge Case Tests

func TestMessageCountsWithUnicodeChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"测试频道-русский-チャンネル", "channel2"}
	opts.ChannelsTimetoken = []int64{123456789, 987654321}

	// Should pass validation
	assert.Nil(opts.validate())

	// Should build path with URL encoding
	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/message-counts/")
	// Unicode should be URL encoded
	assert.Contains(path, "%")
}

func TestMessageCountsWithManyChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a large list of channels
	channels := make([]string, 50)
	channelsTimetoken := make([]int64, 50)
	for i := 0; i < 50; i++ {
		channels[i] = fmt.Sprintf("channel_%d", i)
		channelsTimetoken[i] = int64(123456789 + i)
	}

	builder := newMessageCountsBuilder(pn)
	builder.Channels(channels)
	builder.ChannelsTimetoken(channelsTimetoken)

	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(channelsTimetoken, builder.opts.ChannelsTimetoken)
	assert.Equal(50, len(builder.opts.Channels))
	assert.Equal(50, len(builder.opts.ChannelsTimetoken))

	// Should pass validation (matching lengths)
	assert.Nil(builder.opts.validate())

	// Test path building with many channels
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/message-counts/")
	assert.Contains(path, "channel_0")
	assert.Contains(path, "channel_49")
	// Should contain comma-separated channels
	assert.Contains(path, ",")
}

func TestMessageCountsWithVeryLongChannelNames(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create very long channel names
	longName1 := ""
	longName2 := ""
	for i := 0; i < 100; i++ {
		longName1 += fmt.Sprintf("ch1_%d_", i)
		longName2 += fmt.Sprintf("ch2_%d_", i)
	}

	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{longName1, longName2}
	opts.ChannelsTimetoken = []int64{123456789, 987654321}

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/message-counts/")
	assert.Contains(path, "ch1_0_")
	assert.Contains(path, "ch2_99_")
}

func TestMessageCountsWithExtremeTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newMessageCountsBuilder(pn)
	builder.Channels([]string{"test-channel"})

	// Test extreme timetoken values
	maxTimetoken := int64(9223372036854775807) // Max int64
	minTimetoken := int64(1)

	builder.ChannelsTimetoken([]int64{maxTimetoken})
	assert.Equal([]int64{maxTimetoken}, builder.opts.ChannelsTimetoken)

	// Test query building with extreme values
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("9223372036854775807", query.Get("timetoken"))

	// Test with minimum timetoken
	builder.ChannelsTimetoken([]int64{minTimetoken})
	query, err = builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("1", query.Get("timetoken"))
}

func TestMessageCountsWithEmptyQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}
	opts.ChannelsTimetoken = []int64{123456789}
	opts.QueryParam = map[string]string{}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

func TestMessageCountsWithNilQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newMessageCountsOpts(pn, pn.ctx)
	opts.Channels = []string{"test-channel"}
	opts.ChannelsTimetoken = []int64{123456789}
	opts.QueryParam = nil

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

func TestMessageCountsChannelsTimetokenCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name              string
		channels          []string
		channelsTimetoken []int64
		timetoken         int64
		expectedTimetoken string
		expectedChannels  string
		shouldValidate    bool
	}{
		{
			name:              "Single channel with single timetoken",
			channels:          []string{"ch1"},
			channelsTimetoken: []int64{123456789},
			timetoken:         0,
			expectedTimetoken: "123456789",
			expectedChannels:  "",
			shouldValidate:    true,
		},
		{
			name:              "Multiple channels with matching timetokens",
			channels:          []string{"ch1", "ch2", "ch3"},
			channelsTimetoken: []int64{123456789, 987654321, 555666777},
			timetoken:         0,
			expectedTimetoken: "",
			expectedChannels:  "123456789,987654321,555666777",
			shouldValidate:    true,
		},
		{
			name:              "Legacy single timetoken",
			channels:          []string{"ch1"},
			channelsTimetoken: nil,
			timetoken:         999888777,
			expectedTimetoken: "999888777",
			expectedChannels:  "",
			shouldValidate:    true,
		},
		{
			name:              "Length mismatch should fail validation",
			channels:          []string{"ch1", "ch2", "ch3"},
			channelsTimetoken: []int64{123456789, 987654321}, // 2 timetokens for 3 channels
			timetoken:         0,
			expectedTimetoken: "",
			expectedChannels:  "",
			shouldValidate:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newMessageCountsBuilder(pn)
			builder.Channels(tc.channels)
			builder.ChannelsTimetoken(tc.channelsTimetoken)
			builder.Timetoken(tc.timetoken)

			if tc.shouldValidate {
				assert.Nil(builder.opts.validate(), "Should pass validation for: %s", tc.name)

				query, err := builder.opts.buildQuery()
				assert.Nil(err)
				assert.Equal(tc.expectedTimetoken, query.Get("timetoken"))
				assert.Equal(tc.expectedChannels, query.Get("channelsTimetoken"))
			} else {
				assert.NotNil(builder.opts.validate(), "Should fail validation for: %s", tc.name)
			}
		})
	}
}

// Error Scenario Tests

func TestMessageCountsExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newMessageCountsBuilder(pn)
	builder.Channels([]string{"test-channel"})
	builder.ChannelsTimetoken([]int64{123456789})

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestMessageCountsExecuteErrorMissingChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newMessageCountsBuilder(pn)
	// Don't set Channels, should fail validation
	builder.ChannelsTimetoken([]int64{123456789})

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestMessageCountsExecuteErrorMissingTimetoken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newMessageCountsBuilder(pn)
	builder.Channels([]string{"test-channel"})
	// Don't set any timetoken, should fail validation

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channels Timetoken")
}

func TestMessageCountsExecuteErrorLengthMismatch(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newMessageCountsBuilder(pn)
	builder.Channels([]string{"ch1", "ch2", "ch3"})          // 3 channels
	builder.ChannelsTimetoken([]int64{123456789, 987654321}) // 2 timetokens - mismatch!

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Length of Channels Timetoken and Channels do not match")
}

func TestMessageCountsPathBuildingEdgeCases(t *testing.T) {
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
		opts := newMessageCountsOpts(pn, pn.ctx)
		opts.Channels = []string{channel}
		opts.ChannelsTimetoken = []int64{123456789}

		// Should pass validation
		assert.Nil(opts.validate(), "Should validate channel: %s", channel)

		// Should build valid path
		path, err := opts.buildPath()
		assert.Nil(err, "Should build path for channel: %s", channel)
		assert.Contains(path, "/message-counts/", "Should contain message-counts path for: %s", channel)
	}
}

func TestMessageCountsTransportSetter(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newMessageCountsBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}
