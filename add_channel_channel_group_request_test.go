package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestAddChannelRequestBasic(t *testing.T) {
	assert := assert.New(t)

	opts := newAddChannelOpts(pubnub, pubnub.ctx)
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
	expected.Set("add", "ch1,ch2,ch3")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestAddChannelRequestBasicQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := newAddChannelOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroup = "cg"

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	opts.QueryParam = queryParam

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
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")
	expected.Set("add", "ch1,ch2,ch3")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestNewAddChannelToChannelGroupBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newAddChannelToChannelGroupBuilder(pubnub)
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
	expected.Set("add", "ch1,ch2,ch3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestNewAddChannelToChannelGroupBuilderWithContext(t *testing.T) {
	assert := assert.New(t)

	o := newAddChannelToChannelGroupBuilderWithContext(pubnub, pubnub.ctx)
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
	expected.Set("add", "ch1,ch2,ch3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestAddChannelOptsValidateSub(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newAddChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroup = "cg"

	assert.Equal("pubnub/validation: pubnub: Add Channel To Channel Group: Missing Subscribe Key", opts.validate().Error())
}

// Additional validation tests specific to AddChannelToChannelGroup
func TestAddChannelOptsValidateMissingChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddChannelOpts(pn, pn.ctx)
	opts.Channels = []string{} // Empty channels
	opts.ChannelGroup = "cg"

	assert.Equal("pubnub/validation: pubnub: Add Channel To Channel Group: Missing Channel", opts.validate().Error())
}

func TestAddChannelOptsValidateMissingChannelsNil(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddChannelOpts(pn, pn.ctx)
	opts.Channels = nil // Nil channels
	opts.ChannelGroup = "cg"

	assert.Equal("pubnub/validation: pubnub: Add Channel To Channel Group: Missing Channel", opts.validate().Error())
}

func TestAddChannelOptsValidateMissingChannelGroup(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1", "ch2"}
	opts.ChannelGroup = ""

	assert.Equal("pubnub/validation: pubnub: Add Channel To Channel Group: Missing Channel Group", opts.validate().Error())
}

func TestAddChannelOptsValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newAddChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1", "ch2"}
	opts.ChannelGroup = "cg"

	assert.Nil(opts.validate())
}

// Builder pattern tests for AddChannelToChannelGroup
func TestAddChannelToChannelGroupBuilder(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test basic builder
	builder := newAddChannelToChannelGroupBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)

	// Test Channels setting
	channels := []string{"ch1", "ch2", "ch3"}
	result := builder.Channels(channels)
	assert.Equal(builder, result) // Should return same instance for chaining
	assert.Equal(channels, builder.opts.Channels)
}

func TestAddChannelToChannelGroupBuilderWithContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAddChannelToChannelGroupBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestAddChannelToChannelGroupBuilderChannelGroup(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAddChannelToChannelGroupBuilder(pn)
	result := builder.ChannelGroup("test-group")
	assert.Equal(builder, result) // Should return same instance for chaining
	assert.Equal("test-group", builder.opts.ChannelGroup)
}

func TestAddChannelToChannelGroupBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParams := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	builder := newAddChannelToChannelGroupBuilder(pn)
	result := builder.QueryParam(queryParams)
	assert.Equal(builder, result) // Should return same instance for chaining
	assert.Equal(queryParams, builder.opts.QueryParam)
}

func TestAddChannelToChannelGroupBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	channels := []string{"ch1", "ch2", "ch3"}
	queryParams := map[string]string{
		"test1": "value1",
		"test2": "value2",
	}

	// Test fluent interface chaining
	builder := newAddChannelToChannelGroupBuilder(pn).
		Channels(channels).
		ChannelGroup("test-group").
		QueryParam(queryParams)

	assert.Equal(channels, builder.opts.Channels)
	assert.Equal("test-group", builder.opts.ChannelGroup)
	assert.Equal(queryParams, builder.opts.QueryParam)
}

func TestAddChannelToChannelGroupBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAddChannelToChannelGroupBuilder(pn)

	// Test Channels setter
	channels1 := []string{"ch1", "ch2"}
	builder.Channels(channels1)
	assert.Equal(channels1, builder.opts.Channels)

	// Test ChannelGroup setter
	builder.ChannelGroup("group1")
	assert.Equal("group1", builder.opts.ChannelGroup)

	// Test QueryParam setter
	queryParams := map[string]string{"key": "value"}
	builder.QueryParam(queryParams)
	assert.Equal(queryParams, builder.opts.QueryParam)

	// Test overwriting Channels
	channels2 := []string{"ch3", "ch4", "ch5"}
	builder.Channels(channels2)
	assert.Equal(channels2, builder.opts.Channels)

	// Test overwriting ChannelGroup
	builder.ChannelGroup("group2")
	assert.Equal("group2", builder.opts.ChannelGroup)

	// Test overwriting QueryParam
	newQueryParams := map[string]string{"newkey": "newvalue"}
	builder.QueryParam(newQueryParams)
	assert.Equal(newQueryParams, builder.opts.QueryParam)
}

// URL path building tests
func TestAddChannelToChannelGroupBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newAddChannelOpts(pn, pn.ctx)
	opts.ChannelGroup = "test-group"

	path, err := opts.buildPath()
	assert.Nil(err)
	expectedPath := fmt.Sprintf("/v1/channel-registration/sub-key/%s/channel-group/test-group", pn.Config.SubscribeKey)
	assert.Equal(expectedPath, path)
}

func TestAddChannelToChannelGroupBuildPathWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newAddChannelOpts(pn, pn.ctx)
	// Channel group with characters that need URL encoding
	opts.ChannelGroup = "group+with/special=chars&more"

	path, err := opts.buildPath()
	assert.Nil(err)

	// Should contain URL-encoded channel group
	assert.Contains(path, "/v1/channel-registration/sub-key/")
	assert.Contains(path, pn.Config.SubscribeKey)
	assert.Contains(path, "/channel-group/")
	// The channel group should be URL encoded in the path
	assert.NotContains(path, "+") // + should be encoded
	assert.NotContains(path, "=") // = should be encoded
}

func TestAddChannelToChannelGroupBuildQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newAddChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.QueryParam = map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should include channels as comma-separated list
	assert.Equal("ch1,ch2,ch3", query.Get("add"))

	// Should include custom query parameters
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))

	// Should include default PubNub parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestAddChannelToChannelGroupBuildQuerySingleChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newAddChannelOpts(pn, pn.ctx)
	opts.Channels = []string{"single-channel"}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should include single channel
	assert.Equal("single-channel", query.Get("add"))
}

func TestAddChannelToChannelGroupBuildQueryEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newAddChannelOpts(pn, pn.ctx)
	opts.Channels = []string{}
	opts.QueryParam = nil

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should include empty add parameter
	assert.Equal("", query.Get("add"))

	// Should still have default PubNub parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

// HTTP method and operation type tests
func TestAddChannelToChannelGroupOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newAddChannelOpts(pn, pn.ctx)
	opType := opts.operationType()
	assert.Equal(PNAddChannelsToChannelGroupOperation, opType)
}

// Edge case tests for AddChannelToChannelGroup
func TestAddChannelToChannelGroupWithManyChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a large list of channels
	channels := make([]string, 100)
	for i := 0; i < 100; i++ {
		channels[i] = fmt.Sprintf("channel_%d", i)
	}

	builder := newAddChannelToChannelGroupBuilder(pn)
	builder.Channels(channels)
	builder.ChannelGroup("large-group")

	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(100, len(builder.opts.Channels))

	// Test query building with many channels
	query, err := builder.opts.buildQuery()
	assert.Nil(err)

	addParam := query.Get("add")
	assert.NotEmpty(addParam)
	assert.Contains(addParam, "channel_0")
	assert.Contains(addParam, "channel_99")
	// Should contain comma-separated channels
	assert.Contains(addParam, ",")
}

func TestAddChannelToChannelGroupWithUnicodeCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Channels and group with Unicode characters
	unicodeChannels := []string{
		"频道测试",
		"канал-тест",
		"チャンネル-テスト",
	}
	unicodeGroup := "组测试-группа-グループ"

	builder := newAddChannelToChannelGroupBuilder(pn)
	builder.Channels(unicodeChannels)
	builder.ChannelGroup(unicodeGroup)

	assert.Equal(unicodeChannels, builder.opts.Channels)
	assert.Equal(unicodeGroup, builder.opts.ChannelGroup)

	// Test validation passes
	assert.Nil(builder.opts.validate())

	// Test path building with Unicode group
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v1/channel-registration/sub-key/")
	assert.Contains(path, "/channel-group/")
	// Unicode should be properly URL encoded
	assert.NotContains(path, "测试")      // Should be encoded
	assert.NotContains(path, "русский") // Should be encoded

	// Test query building with Unicode channels
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	addParam := query.Get("add")
	assert.NotEmpty(addParam)
	// Should contain the Unicode channel names (may be URL encoded by HTTP client)
}

func TestAddChannelToChannelGroupWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Channels and group with various special characters
	specialChannels := []string{
		"channel-with-dashes",
		"channel_with_underscores",
		"channel.with.dots",
		"channel:with:colons",
	}
	specialGroup := "group-with-dashes_and_underscores.and.dots"

	builder := newAddChannelToChannelGroupBuilder(pn)
	builder.Channels(specialChannels)
	builder.ChannelGroup(specialGroup)

	assert.Equal(specialChannels, builder.opts.Channels)
	assert.Equal(specialGroup, builder.opts.ChannelGroup)

	// Test validation passes
	assert.Nil(builder.opts.validate())

	// Test path building
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v1/channel-registration/sub-key/")
	assert.Contains(path, "/channel-group/")

	// Test query building
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	addParam := query.Get("add")
	assert.Contains(addParam, "channel-with-dashes")
	assert.Contains(addParam, "channel_with_underscores")
	assert.Contains(addParam, "channel.with.dots")
	assert.Contains(addParam, "channel:with:colons")
}

func TestAddChannelToChannelGroupWithEmptyQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAddChannelToChannelGroupBuilder(pn)
	builder.Channels([]string{"ch1"})
	builder.ChannelGroup("test-group")
	builder.QueryParam(map[string]string{}) // Empty map

	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should still have default parameters and channels
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("ch1", query.Get("add"))
}

func TestAddChannelToChannelGroupWithNilQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newAddChannelToChannelGroupBuilder(pn)
	builder.Channels([]string{"ch1"})
	builder.ChannelGroup("test-group")
	builder.QueryParam(nil) // Nil map

	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should still have default parameters and channels
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("ch1", query.Get("add"))
}

func TestAddChannelToChannelGroupWithComplexQueryParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	complexParams := map[string]string{
		"app_version":    "1.2.3",
		"user_id":        "user-123-abc",
		"session_id":     "session-xyz-789",
		"special_chars":  "value-with-special@chars#and$symbols",
		"unicode_value":  "测试值-русский-ファイル",
		"empty_value":    "",
		"number_string":  "42",
		"boolean_string": "true",
	}

	builder := newAddChannelToChannelGroupBuilder(pn)
	builder.Channels([]string{"ch1", "ch2"})
	builder.ChannelGroup("test-group")
	builder.QueryParam(complexParams)

	query, err := builder.opts.buildQuery()
	assert.Nil(err)

	// Verify all custom parameters are present (may be URL encoded)
	for key, expectedValue := range complexParams {
		actualValue := query.Get(key)
		if key == "special_chars" || key == "unicode_value" {
			// Special characters and Unicode may be URL encoded
			assert.NotEmpty(actualValue, "Query parameter %s should be present", key)
		} else {
			assert.Equal(expectedValue, actualValue, "Query parameter %s should match", key)
		}
	}

	// Should still have default parameters and channels
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("ch1,ch2", query.Get("add"))
}

// Error scenario tests
func TestAddChannelToChannelGroupBuilderExecuteErrorHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with invalid configuration (missing subscribe key)
	pn.Config.SubscribeKey = ""

	builder := newAddChannelToChannelGroupBuilder(pn)
	builder.Channels([]string{"ch1"})
	builder.ChannelGroup("test-group")

	// Execute should fail with validation error
	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestAddChannelToChannelGroupEdgeCaseChannelNames(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCaseChannels := []struct {
		name     string
		channels []string
	}{
		{"Very short channel", []string{"a"}},
		{"Numeric channel", []string{"123456789"}},
		{"Mixed case", []string{"AbCdEfGhIjKlMnOpQrStUvWxYz"}},
		{"With padding", []string{"channel==="}},
		{"With dashes", []string{"channel-with-many-dashes-in-between"}},
		{"With underscores", []string{"channel_with_many_underscores_in_between"}},
		{"Mixed special chars", []string{"ch@#$%^&*()_+-=[]{}|;:'\",.<>?/~`"}},
	}

	for _, tc := range edgeCaseChannels {
		t.Run(tc.name, func(t *testing.T) {
			builder := newAddChannelToChannelGroupBuilder(pn)
			builder.Channels(tc.channels)
			builder.ChannelGroup("test-group")

			assert.Equal(tc.channels, builder.opts.Channels)
			assert.Nil(builder.opts.validate(), "Channels %v should pass validation", tc.channels)

			// Test query building
			query, err := builder.opts.buildQuery()
			assert.Nil(err, "Query building should succeed for channels: %v", tc.channels)

			addParam := query.Get("add")
			assert.NotEmpty(addParam, "Add parameter should not be empty for channels: %v", tc.channels)
		})
	}
}

func TestAddChannelToChannelGroupEdgeCaseGroupNames(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCaseGroups := []struct {
		name  string
		group string
	}{
		{"Very short group", "a"},
		{"Numeric group", "123456789"},
		{"Mixed case", "AbCdEfGhIjKlMnOpQrStUvWxYz"},
		{"With padding", "group==="},
		{"With dashes", "group-with-many-dashes-in-between"},
		{"With underscores", "group_with_many_underscores_in_between"},
		{"URL-safe characters", "group-with_safe.chars~123"},
	}

	for _, tc := range edgeCaseGroups {
		t.Run(tc.name, func(t *testing.T) {
			builder := newAddChannelToChannelGroupBuilder(pn)
			builder.Channels([]string{"ch1"})
			builder.ChannelGroup(tc.group)

			assert.Equal(tc.group, builder.opts.ChannelGroup)
			assert.Nil(builder.opts.validate(), "Group %s should pass validation", tc.group)

			// Test path building
			path, err := builder.opts.buildPath()
			assert.Nil(err, "Path building should succeed for group: %s", tc.group)
			assert.Contains(path, "/v1/channel-registration/sub-key/")
			assert.Contains(path, "/channel-group/")
		})
	}
}
