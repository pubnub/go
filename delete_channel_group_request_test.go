package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestDeleteChannelGroupRequestBasic(t *testing.T) {
	assert := assert.New(t)

	opts := newDeleteChannelGroupOpts(pubnub, pubnub.ctx, deleteChannelGroupOpts{
		ChannelGroup: "cg",
	})

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v1/channel-registration/sub-key/sub_key/channel-group/cg/remove",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestDeleteChannelGroupRequestBasicQueryParam(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := newDeleteChannelGroupOpts(pubnub, pubnub.ctx, deleteChannelGroupOpts{
		ChannelGroup: "cg",
		QueryParam:   queryParam,
	})

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v1/channel-registration/sub-key/sub_key/channel-group/cg/remove",
		u.EscapedPath(), []int{})

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

func TestNewDeleteChannelGroupBuilder(t *testing.T) {
	assert := assert.New(t)
	o := newDeleteChannelGroupBuilder(pubnub)
	o.ChannelGroup("cg")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v1/channel-registration/sub-key/sub_key/channel-group/cg/remove",
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestNewDeleteChannelGroupBuilderContext(t *testing.T) {
	assert := assert.New(t)
	o := newDeleteChannelGroupBuilderWithContext(pubnub, pubnub.ctx)
	o.ChannelGroup("cg")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v1/channel-registration/sub-key/sub_key/channel-group/cg/remove",
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}
func TestDeleteChannelGroupOptsValidateSub(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{
		ChannelGroup: "cg",
	})

	assert.Equal("pubnub/validation: pubnub: Remove Channel Group: Missing Subscribe Key", opts.validate().Error())
}

// Additional Validation Tests

func TestDeleteChannelGroupValidateMissingChannelGroup(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{
		ChannelGroup: "",
	})

	assert.Equal("pubnub/validation: pubnub: Remove Channel Group: Missing Channel Group", opts.validate().Error())
}

func TestDeleteChannelGroupValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{
		ChannelGroup: "test-group",
	})

	assert.Nil(opts.validate())
}

// Builder Pattern Tests

func TestDeleteChannelGroupBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newDeleteChannelGroupBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestDeleteChannelGroupBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newDeleteChannelGroupBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestDeleteChannelGroupBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{"key": "value"}

	builder := newDeleteChannelGroupBuilder(pn)
	result := builder.ChannelGroup("test-group").QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-group", builder.opts.ChannelGroup)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestDeleteChannelGroupBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newDeleteChannelGroupBuilder(pn)

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

func TestDeleteChannelGroupBuilderChannelGroup(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newDeleteChannelGroupBuilder(pn)
	builder.ChannelGroup("test-group")

	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expected := "/v1/channel-registration/sub-key/demo/channel-group/test-group/remove"
	assert.Equal(expected, path)
}

func TestDeleteChannelGroupBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newDeleteChannelGroupBuilder(pn)
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

func TestDeleteChannelGroupBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{
		ChannelGroup: "test-group",
	})

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/channel-registration/sub-key/demo/channel-group/test-group/remove"
	assert.Equal(expected, path)
}

func TestDeleteChannelGroupBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{
		ChannelGroup: "group-with-special@chars#and$symbols",
	})

	path, err := opts.buildPath()
	assert.Nil(err)
	// Should URL encode special characters
	assert.Contains(path, "group-with-special%40chars%23and%24symbols")
	assert.Contains(path, "/remove")
}

func TestDeleteChannelGroupBuildQueryEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Should have default parameters only
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestDeleteChannelGroupBuildQueryWithParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	queryParam := map[string]string{
		"custom": "value",
		"test":   "param",
	}
	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{
		QueryParam: queryParam,
	})

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

func TestDeleteChannelGroupOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{})

	assert.Equal(PNRemoveGroupOperation, opts.operationType())
}

// Edge Case Tests

func TestDeleteChannelGroupWithUnicodeChannelGroup(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{
		ChannelGroup: "测试群组-русский-ファイル",
	})

	// Should pass validation
	assert.Nil(opts.validate())

	// Should build path with URL encoding
	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/channel-group/")
	assert.Contains(path, "/remove")
	// Unicode should be URL encoded
	assert.Contains(path, "%")
}

func TestDeleteChannelGroupWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialGroupNames := []string{
		"group_with_underscores",
		"group-with-hyphens",
		"group.with.dots",
		"group with spaces",
		"group@with#special$chars",
		"group%already%encoded",
	}

	for _, groupName := range specialGroupNames {
		opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{
			ChannelGroup: groupName,
		})

		// Should pass validation
		assert.Nil(opts.validate(), "Should validate group name: %s", groupName)

		// Should build valid path
		path, err := opts.buildPath()
		assert.Nil(err, "Should build path for group name: %s", groupName)
		assert.Contains(path, "/channel-group/", "Should contain correct path for: %s", groupName)
		assert.Contains(path, "/remove", "Should contain /remove suffix for: %s", groupName)
	}
}

func TestDeleteChannelGroupWithVeryLongChannelGroupName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a very long channel group name
	longName := ""
	for i := 0; i < 100; i++ {
		longName += fmt.Sprintf("group_%d_", i)
	}

	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{
		ChannelGroup: longName,
	})

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/channel-group/")
	assert.Contains(path, "/remove")
}

func TestDeleteChannelGroupWithEmptyQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{
		QueryParam: map[string]string{},
	})

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

func TestDeleteChannelGroupWithNilQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{
		QueryParam: nil,
	})

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

func TestDeleteChannelGroupWithComplexQueryParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

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
	opts := newDeleteChannelGroupOpts(pn, pn.ctx, deleteChannelGroupOpts{
		QueryParam: complexParams,
	})

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

func TestDeleteChannelGroupResponseProcessing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newDeleteChannelGroupBuilder(pn)
	builder.ChannelGroup("test-group")

	// Test that Execute returns empty response for successful case
	// Note: This would normally require a mock, but we're testing the response structure
	resp := emptyDeleteChannelGroupResponse
	assert.Nil(resp) // Should be nil as defined in the implementation
}

func TestDeleteChannelGroupResponseStructure(t *testing.T) {
	assert := assert.New(t)

	// Test that DeleteChannelGroupResponse is a simple empty struct
	response := &DeleteChannelGroupResponse{}
	assert.NotNil(response)

	// Should be a zero-value struct
	assert.Equal(DeleteChannelGroupResponse{}, *response)
}

// Error Scenario Tests

func TestDeleteChannelGroupExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newDeleteChannelGroupBuilder(pn)
	builder.ChannelGroup("test-group")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestDeleteChannelGroupWithInvalidChannelGroup(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newDeleteChannelGroupBuilder(pn)
	// Don't set ChannelGroup, should fail validation

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel Group")
}
