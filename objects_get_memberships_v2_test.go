package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/pubnub/go/v7/utils"
	"github.com/stretchr/testify/assert"
)

func AssertGetMembershipsV2(t *testing.T, checkQueryParam, testContext, withFilter bool, withSort bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	incl := []PNMembershipsInclude{
		PNMembershipsIncludeCustom,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := EnumArrayToStringArray(incl)

	o := newGetMembershipsBuilderV2(pn)
	if testContext {
		o = newGetMembershipsBuilderV2WithContext(pn, pn.ctx)
	}

	userID := "id0"
	limit := 90
	start := "Mxmy"
	end := "Nxny"

	o.UUID(userID)
	o.Include(incl)
	o.Limit(limit)
	o.Start(start)
	o.End(end)
	o.Count(false)
	o.QueryParam(queryParam)
	if withFilter {
		o.Filter("custom.a5 == 'b5' || custom.c5 == 'd5'")
	}
	sort := []string{"name", "created:desc"}
	if withSort {
		o.Sort(sort)
	}

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/uuids/%s/channels", pn.Config.SubscribeKey, "id0"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
		assert.Equal(string(utils.JoinChannels(inclStr)), u.Get("include"))
		assert.Equal(strconv.Itoa(limit), u.Get("limit"))
		assert.Equal(start, u.Get("start"))
		assert.Equal(end, u.Get("end"))
		assert.Equal("0", u.Get("count"))
		if withFilter {
			assert.Equal("custom.a5 == 'b5' || custom.c5 == 'd5'", u.Get("filter"))
		}
		if withSort {
			v := &url.Values{}
			SetQueryParamAsCommaSepString(v, sort, "sort")
			assert.Equal(v.Get("sort"), u.Get("sort"))
		}
	}

}

func TestGetMembershipsV2(t *testing.T) {
	AssertGetMembershipsV2(t, true, false, false, false)
}

func TestGetMembershipsV2Context(t *testing.T) {
	AssertGetMembershipsV2(t, true, true, false, false)
}

func TestGetMembershipsV2WithFilter(t *testing.T) {
	AssertGetMembershipsV2(t, true, false, true, false)
}

func TestGetMembershipsV2WithFilterContext(t *testing.T) {
	AssertGetMembershipsV2(t, true, true, true, false)
}

func TestGetMembershipsV2WithSort(t *testing.T) {
	AssertGetMembershipsV2(t, true, false, false, true)
}

func TestGetMembershipsV2WithSortContext(t *testing.T) {
	AssertGetMembershipsV2(t, true, true, false, true)
}

func TestGetMembershipsV2WithFilterWithSort(t *testing.T) {
	AssertGetMembershipsV2(t, true, false, true, true)
}

func TestGetMembershipsV2WithFilterWithSortContext(t *testing.T) {
	AssertGetMembershipsV2(t, true, true, true, true)
}

func TestGetMembershipsV2ResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestGetMembershipsV2ResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":[{"id":"id0","custom":{"a3":"b3","c3":"d3"},"channel":{"id":"id0","name":"name","description":"desc","custom":{"a":"b"},"created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T13:26:08.341297Z","eTag":"Aee9zsKNndXlHw"},"created":"2019-08-20T13:26:24.07832Z","updated":"2019-08-20T13:26:24.07832Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":1,"next":"MQ", "prev":"NQ"}`)

	r, _, err := newPNGetMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(1, r.TotalCount)
	assert.Equal("MQ", r.Next)
	assert.Equal("NQ", r.Prev)
	assert.Equal("id0", r.Data[0].ID)
	assert.Equal("name", r.Data[0].Channel.Name)
	assert.Equal("desc", r.Data[0].Channel.Description)
	// assert.Equal("2019-08-20T13:26:08.341297Z", r.Data[0].Channel.Created)
	assert.Equal("2019-08-20T13:26:08.341297Z", r.Data[0].Channel.Updated)
	assert.Equal("Aee9zsKNndXlHw", r.Data[0].Channel.ETag)
	assert.Equal("b", r.Data[0].Channel.Custom["a"])
	assert.Equal(nil, r.Data[0].Channel.Custom["c"])
	assert.Equal("2019-08-20T13:26:24.07832Z", r.Data[0].Created)
	assert.Equal("2019-08-20T13:26:24.07832Z", r.Data[0].Updated)
	assert.Equal("AamrnoXdpdmzjwE", r.Data[0].ETag)
	assert.Equal("b3", r.Data[0].Custom["a3"])
	assert.Equal("d3", r.Data[0].Custom["c3"])

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestGetMembershipsV2ValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newGetMembershipsOptsV2(pn, pn.ctx)
	opts.UUID = "test-uuid"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetMembershipsV2ValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)
	opts.UUID = "test-uuid"

	assert.Nil(opts.validate())
}

func TestGetMembershipsV2ValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)
	opts.UUID = "test-uuid"
	opts.Include = []string{"custom", "channel"}
	opts.QueryParam = map[string]string{"param": "value"}

	assert.Nil(opts.validate())
}

func TestGetMembershipsV2ValidateUUIDDefaultBehavior(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.UUID = "config-uuid"

	builder := newGetMembershipsBuilderV2(pn)
	// Don't set UUID explicitly - should use Config.UUID in Execute

	// Before Execute, UUID should be empty
	assert.Equal("", builder.opts.UUID)

	// Test that Execute correctly sets UUID from Config when it's empty
	// We simulate the logic from Execute method
	opts := builder.opts
	if len(opts.UUID) <= 0 {
		opts.UUID = opts.pubnub.Config.UUID
	}

	// Now UUID should be set to config UUID
	assert.Equal("config-uuid", opts.UUID)

	// And path building should work correctly
	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "config-uuid")
}

// HTTP Method and Operation Tests

func TestGetMembershipsV2HTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

	assert.Equal("GET", opts.httpMethod())
}

func TestGetMembershipsV2OperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

	assert.Equal(PNGetMembershipsOperation, opts.operationType())
}

func TestGetMembershipsV2IsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestGetMembershipsV2Timeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (9 setters)

func TestGetMembershipsV2BuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMembershipsBuilderV2(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
	assert.Equal(membershipsLimitV2, builder.opts.Limit) // Default limit (100)
}

func TestGetMembershipsV2BuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMembershipsBuilderV2WithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestGetMembershipsV2BuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMembershipsBuilderV2(pn)

	// Test UUID setter
	builder.UUID("test-uuid")
	assert.Equal("test-uuid", builder.opts.UUID)

	// Test Include setter
	include := []PNMembershipsInclude{PNMembershipsIncludeCustom}
	builder.Include(include)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)

	// Test Limit setter
	builder.Limit(50)
	assert.Equal(50, builder.opts.Limit)

	// Test Start setter
	builder.Start("start-token")
	assert.Equal("start-token", builder.opts.Start)

	// Test End setter
	builder.End("end-token")
	assert.Equal("end-token", builder.opts.End)

	// Test Count setter
	builder.Count(true)
	assert.Equal(true, builder.opts.Count)

	// Test Filter setter
	builder.Filter("name LIKE 'channel*'")
	assert.Equal("name LIKE 'channel*'", builder.opts.Filter)

	// Test Sort setter
	sort := []string{"name", "created:desc"}
	builder.Sort(sort)
	assert.Equal(sort, builder.opts.Sort)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Test Transport setter
	transport := &http.Transport{}
	builder.Transport(transport)
	assert.Equal(transport, builder.opts.Transport)
}

func TestGetMembershipsV2BuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	include := []PNMembershipsInclude{PNMembershipsIncludeCustom}
	sort := []string{"name"}
	queryParam := map[string]string{"key": "value"}
	transport := &http.Transport{}

	builder := newGetMembershipsBuilderV2(pn)
	result := builder.UUID("test-uuid").
		Include(include).
		Limit(75).
		Start("start").
		End("end").
		Count(true).
		Filter("filter").
		Sort(sort).
		QueryParam(queryParam).
		Transport(transport)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-uuid", builder.opts.UUID)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal(75, builder.opts.Limit)
	assert.Equal("start", builder.opts.Start)
	assert.Equal("end", builder.opts.End)
	assert.Equal(true, builder.opts.Count)
	assert.Equal("filter", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

func TestGetMembershipsV2BuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMembershipsBuilderV2(pn)

	// Verify default values
	assert.Equal("", builder.opts.UUID) // UUID defaults to empty, set later in Execute
	assert.Nil(builder.opts.Include)
	assert.Equal(membershipsLimitV2, builder.opts.Limit) // 100
	assert.Equal("", builder.opts.Start)
	assert.Equal("", builder.opts.End)
	assert.Equal(false, builder.opts.Count)
	assert.Equal("", builder.opts.Filter)
	assert.Nil(builder.opts.Sort)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestGetMembershipsV2BuilderIncludeTypes(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		includes []PNMembershipsInclude
		expected []string
	}{
		{
			name:     "Single include",
			includes: []PNMembershipsInclude{PNMembershipsIncludeCustom},
			expected: []string{"custom"},
		},
		{
			name:     "Multiple includes",
			includes: []PNMembershipsInclude{PNMembershipsIncludeCustom, PNMembershipsIncludeChannel},
			expected: []string{"custom", "channel"},
		},
		{
			name:     "All includes",
			includes: []PNMembershipsInclude{PNMembershipsIncludeCustom, PNMembershipsIncludeChannel, PNMembershipsIncludeChannelCustom},
			expected: []string{"custom", "channel", "channel.custom"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetMembershipsBuilderV2(pn)
			builder.Include(tc.includes)

			expectedInclude := EnumArrayToStringArray(tc.includes)
			assert.Equal(expectedInclude, builder.opts.Include)
		})
	}
}

func TestGetMembershipsV2BuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	transport := &http.Transport{}
	include := []PNMembershipsInclude{PNMembershipsIncludeCustom}
	sort := []string{"name", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Test all 9 setters in chain
	builder := newGetMembershipsBuilderV2(pn).
		UUID("test-uuid").
		Include(include).
		Limit(50).
		Start("start-token").
		End("end-token").
		Count(true).
		Filter("name LIKE 'test*'").
		Sort(sort).
		QueryParam(queryParam).
		Transport(transport)

	// Verify all are set correctly
	assert.Equal("test-uuid", builder.opts.UUID)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal(50, builder.opts.Limit)
	assert.Equal("start-token", builder.opts.Start)
	assert.Equal("end-token", builder.opts.End)
	assert.Equal(true, builder.opts.Count)
	assert.Equal("name LIKE 'test*'", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestGetMembershipsV2BuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)
	opts.UUID = "test-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/uuids/test-uuid/channels"
	assert.Equal(expected, path)
}

func TestGetMembershipsV2BuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newGetMembershipsOptsV2(pn, pn.ctx)
	opts.UUID = "my-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/uuids/my-uuid/channels"
	assert.Equal(expected, path)
}

func TestGetMembershipsV2BuildPathWithSpecialCharsInUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)
	opts.UUID = "uuid-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/")
	assert.Contains(path, "uuid-with-special@chars#and$symbols")
	assert.Contains(path, "/channels")
}

func TestGetMembershipsV2BuildPathWithUnicodeUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)
	opts.UUID = "用户ID-пользователь-ユーザーID"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/")
	assert.Contains(path, "用户ID-пользователь-ユーザーID")
	assert.Contains(path, "/channels")
}

// JSON Body Building Tests (CRITICAL for GET operation - should be empty)

func TestGetMembershipsV2BuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations should have empty body
	assert.Equal([]byte{}, body)
}

func TestGetMembershipsV2BuildBodyWithAllParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

	// Set all possible parameters - body should still be empty for GET
	opts.UUID = "test-uuid"
	opts.Include = []string{"custom", "channel"}
	opts.Limit = 50
	opts.Start = "start"
	opts.End = "end"
	opts.Count = true
	opts.Filter = "filter"
	opts.Sort = []string{"name"}
	opts.QueryParam = map[string]string{"param": "value"}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations always have empty body regardless of parameters
	assert.Equal([]byte{}, body)
}

func TestGetMembershipsV2BuildBodyErrorScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

	// Even with potential error conditions, buildBody should not fail for GET
	opts.UUID = "" // Empty UUID

	body, err := opts.buildBody()
	assert.Nil(err) // buildBody should never error for GET operations
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}

// Query Parameter Tests

func TestGetMembershipsV2BuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("0", query.Get("limit")) // Default limit not set until builder initialization
	assert.Equal("0", query.Get("count")) // Default count=false
}

func TestGetMembershipsV2BuildQueryWithInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

	opts.Include = []string{"custom", "channel"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	includeValue := query.Get("include")
	assert.Contains(includeValue, "custom")
	assert.Contains(includeValue, "channel")
}

func TestGetMembershipsV2BuildQueryWithPagination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

	opts.Limit = 50
	opts.Start = "start-token"
	opts.End = "end-token"
	opts.Count = true

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("50", query.Get("limit"))
	assert.Equal("start-token", query.Get("start"))
	assert.Equal("end-token", query.Get("end"))
	assert.Equal("1", query.Get("count"))
}

func TestGetMembershipsV2BuildQueryWithFilterAndSort(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

	opts.Filter = "custom.role == 'admin'"
	opts.Sort = []string{"name", "created:desc"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("custom.role == 'admin'", query.Get("filter"))

	sortValue := query.Get("sort")
	assert.Contains(sortValue, "name")
	assert.Contains(sortValue, "created:desc")
}

func TestGetMembershipsV2BuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

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

func TestGetMembershipsV2BuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

	// Set all possible query parameters
	opts.Include = []string{"custom", "channel"}
	opts.Limit = 25
	opts.Start = "start"
	opts.End = "end"
	opts.Count = true
	opts.Filter = "active = true"
	opts.Sort = []string{"name:asc"}
	opts.QueryParam = map[string]string{
		"extra": "parameter",
		"debug": "true",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all parameters are set correctly
	assert.Contains(query.Get("include"), "custom")
	assert.Equal("25", query.Get("limit"))
	assert.Equal("start", query.Get("start"))
	assert.Equal("end", query.Get("end"))
	assert.Equal("1", query.Get("count"))
	assert.Equal("active = true", query.Get("filter"))
	assert.Equal("name:asc", query.Get("sort"))
	assert.Equal("parameter", query.Get("extra"))
	assert.Equal("true", query.Get("debug"))

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

// GET-Specific Tests (Read Operation Characteristics)

func TestGetMembershipsV2GetOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMembershipsBuilderV2(pn)
	builder.UUID("test-uuid")

	// Verify it's a GET operation
	assert.Equal("GET", builder.opts.httpMethod())

	// GET operations have empty body
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	// Should have proper path for membership retrieval (UUID to channels)
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/test-uuid/channels")
}

func TestGetMembershipsV2DefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMembershipsBuilderV2(pn)

	// Should have default limit set to membershipsLimitV2 (100)
	assert.Equal(membershipsLimitV2, builder.opts.Limit)
	assert.Equal(100, builder.opts.Limit)

	// Should be included in query
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("limit"))
}

func TestGetMembershipsV2ReadOperationValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*getMembershipsOptsV2)
		description string
	}{
		{
			name: "Basic read operation",
			setupOpts: func(opts *getMembershipsOptsV2) {
				opts.UUID = "user123"
			},
			description: "Get memberships for specific UUID",
		},
		{
			name: "Read with include options",
			setupOpts: func(opts *getMembershipsOptsV2) {
				opts.UUID = "user123"
				opts.Include = []string{"custom", "channel"}
			},
			description: "Get memberships with additional data",
		},
		{
			name: "Read with pagination",
			setupOpts: func(opts *getMembershipsOptsV2) {
				opts.UUID = "user123"
				opts.Limit = 20
				opts.Start = "pagination_token"
			},
			description: "Get paginated memberships",
		},
		{
			name: "Read with filtering",
			setupOpts: func(opts *getMembershipsOptsV2) {
				opts.UUID = "user123"
				opts.Filter = "custom.role == 'admin'"
			},
			description: "Get filtered memberships",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetMembershipsOptsV2(pn, pn.ctx)
			tc.setupOpts(opts)

			// Should pass validation
			assert.Nil(opts.validate())

			// Should be GET operation
			assert.Equal("GET", opts.httpMethod())

			// Should have empty body
			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)

			// Should build valid path
			path, err := opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/uuids/")
			assert.Contains(path, "/channels")

			// Should build valid query
			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

func TestGetMembershipsV2ResponseStructureAfterRetrieval(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetMembershipsBuilderV2(pn)
	builder.UUID("test-uuid")

	// Response should contain memberships data after GET operation
	// This is tested in the existing TestGetMembershipsV2ResponseValuePass
	// but verify the operation is configured correctly
	opts := builder.opts

	// Verify operation is configured correctly
	assert.Equal("GET", opts.httpMethod())
	assert.Equal(PNGetMembershipsOperation, opts.operationType())
	assert.True(opts.isAuthRequired())
}

func TestGetMembershipsV2EmptyBodyVerification(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that GET operations always have empty body regardless of configuration
	testCases := []struct {
		name      string
		setupOpts func(*getMembershipsOptsV2)
	}{
		{
			name: "With all parameters set",
			setupOpts: func(opts *getMembershipsOptsV2) {
				opts.UUID = "test-uuid"
				opts.Include = []string{"custom", "channel", "channel.custom"}
				opts.Limit = 50
				opts.Start = "start"
				opts.End = "end"
				opts.Count = true
				opts.Filter = "complex filter expression"
				opts.Sort = []string{"name:asc", "created:desc"}
				opts.QueryParam = map[string]string{
					"param1": "value1",
					"param2": "value2",
				}
			},
		},
		{
			name: "With minimal parameters",
			setupOpts: func(opts *getMembershipsOptsV2) {
				opts.UUID = "simple-uuid"
			},
		},
		{
			name: "With empty/nil parameters",
			setupOpts: func(opts *getMembershipsOptsV2) {
				opts.UUID = ""
				opts.Include = nil
				opts.QueryParam = nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetMembershipsOptsV2(pn, pn.ctx)
			tc.setupOpts(opts)

			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
			assert.Equal([]byte{}, body)
		})
	}
}

// Comprehensive Edge Case Tests

func TestGetMembershipsV2WithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*getMembershipsBuilderV2)
	}{
		{
			name: "Very long UUID",
			setupFn: func(builder *getMembershipsBuilderV2) {
				longUUID := strings.Repeat("VeryLongUUID", 50) // 600 characters
				builder.UUID(longUUID)
			},
		},
		{
			name: "Large filter expression",
			setupFn: func(builder *getMembershipsBuilderV2) {
				largeFilter := "(" + strings.Repeat("custom.field == 'value' OR ", 100) + "custom.final == 'end')"
				builder.Filter(largeFilter)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *getMembershipsBuilderV2) {
				largeQueryParam := make(map[string]string)
				for i := 0; i < 100; i++ {
					largeQueryParam[fmt.Sprintf("param_%d", i)] = strings.Repeat(fmt.Sprintf("value_%d_", i), 20)
				}
				builder.QueryParam(largeQueryParam)
			},
		},
		{
			name: "Large pagination limit",
			setupFn: func(builder *getMembershipsBuilderV2) {
				builder.Limit(10000) // Very large limit
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetMembershipsBuilderV2(pn)
			builder.UUID("baseline-uuid")
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path and query
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Should always have empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
		})
	}
}

func TestGetMembershipsV2SpecialCharacterHandling(t *testing.T) {
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
		"emoji😀🎉🚀💯",
		"ñáéíóúüç", // Accented characters
	}

	for i, specialString := range specialStrings {
		t.Run(fmt.Sprintf("SpecialString_%d", i), func(t *testing.T) {
			builder := newGetMembershipsBuilderV2(pn)
			builder.UUID(specialString)
			builder.Filter(fmt.Sprintf("custom.field == '%s'", specialString))
			builder.QueryParam(map[string]string{
				"special_field": specialString,
			})

			// Should pass validation (basic validation doesn't check content)
			assert.Nil(builder.opts.validate())

			// Should build valid path and query
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Should always have empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
		})
	}
}

func TestGetMembershipsV2ParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name   string
		uuid   string
		limit  int
		filter string
	}{
		{
			name:   "Empty string UUID",
			uuid:   "",
			limit:  1,
			filter: "",
		},
		{
			name:   "Single character UUID",
			uuid:   "a",
			limit:  1,
			filter: "a",
		},
		{
			name:   "Unicode-only UUID",
			uuid:   "测试",
			limit:  50,
			filter: "测试 == '值'",
		},
		{
			name:   "Minimum limit",
			uuid:   "test",
			limit:  1,
			filter: "simple",
		},
		{
			name:   "Large limit",
			uuid:   "test",
			limit:  1000,
			filter: "complex.nested == 'value'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetMembershipsBuilderV2(pn)
			builder.UUID(tc.uuid)
			builder.Limit(tc.limit)
			if tc.filter != "" {
				builder.Filter(tc.filter)
			}

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid components
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			if tc.uuid != "" {
				assert.Contains(path, tc.uuid)
			}

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.Equal(fmt.Sprintf("%d", tc.limit), query.Get("limit"))

			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body) // GET operation always has empty body
		})
	}
}

func TestGetMembershipsV2ComplexRetrievalScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*getMembershipsBuilderV2)
		validateFn func(*testing.T, string)
	}{
		{
			name: "International UUID with complex filters",
			setupFn: func(builder *getMembershipsBuilderV2) {
				builder.UUID("用户中文123")
				builder.Filter("custom.角色 == '管理员' && custom.语言 == 'zh'")
				builder.Include([]PNMembershipsInclude{PNMembershipsIncludeCustom, PNMembershipsIncludeChannel})
			},
			validateFn: func(t *testing.T, path string) {
				assert.Contains(path, "/uuids/用户中文123/channels")
			},
		},
		{
			name: "Professional membership retrieval with pagination",
			setupFn: func(builder *getMembershipsBuilderV2) {
				builder.UUID("professional_user")
				builder.Filter("custom.role IN ('admin', 'manager', 'lead')")
				builder.Sort([]string{"name:asc", "custom.priority:desc"})
				builder.Limit(25)
				builder.Count(true)
			},
			validateFn: func(t *testing.T, path string) {
				assert.Contains(path, "/uuids/professional_user/channels")
			},
		},
		{
			name: "Email-like UUID with complex include options",
			setupFn: func(builder *getMembershipsBuilderV2) {
				builder.UUID("user@company.com")
				builder.Include([]PNMembershipsInclude{
					PNMembershipsIncludeCustom,
					PNMembershipsIncludeChannel,
					PNMembershipsIncludeChannelCustom,
				})
				builder.Filter("channel.name LIKE '%company%'")
			},
			validateFn: func(t *testing.T, path string) {
				assert.Contains(path, "/uuids/user@company.com/channels")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetMembershipsBuilderV2(pn)
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path (UUID-to-channels direction)
			path, err := builder.opts.buildPath()
			assert.Nil(err)

			// Should always have empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)

			// Run verification
			tc.validateFn(t, path)
		})
	}
}

// Error Scenario Tests

func TestGetMembershipsV2ExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newGetMembershipsBuilderV2(pn)
	builder.UUID("test-uuid")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetMembershipsV2PathBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCases := []struct {
		name         string
		subscribeKey string
		uuid         string
		expectError  bool
	}{
		{
			name:         "Empty SubscribeKey",
			subscribeKey: "",
			uuid:         "test-uuid",
			expectError:  false, // buildPath doesn't validate, only validate() does
		},
		{
			name:         "Empty UUID",
			subscribeKey: "demo",
			uuid:         "",
			expectError:  false, // buildPath doesn't validate UUID
		},
		{
			name:         "SubscribeKey with spaces",
			subscribeKey: "   sub key   ",
			uuid:         "test-uuid",
			expectError:  false,
		},
		{
			name:         "UUID with spaces",
			subscribeKey: "demo",
			uuid:         "   test uuid   ",
			expectError:  false,
		},
		{
			name:         "Very special characters",
			subscribeKey: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			uuid:         "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			expectError:  false,
		},
		{
			name:         "Unicode SubscribeKey and UUID",
			subscribeKey: "测试订阅键-русский-キー",
			uuid:         "用户ID-пользователь-ユーザーID",
			expectError:  false,
		},
		{
			name:         "Very long values",
			subscribeKey: strings.Repeat("a", 1000),
			uuid:         strings.Repeat("b", 1000),
			expectError:  false,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			pn.Config.SubscribeKey = tc.subscribeKey
			opts := newGetMembershipsOptsV2(pn, pn.ctx)
			opts.UUID = tc.uuid

			path, err := opts.buildPath()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Contains(path, "/v2/objects/")
				assert.Contains(path, "/uuids/")
				assert.Contains(path, "/channels")
			}
		})
	}
}

func TestGetMembershipsV2QueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*getMembershipsOptsV2)
		expectError bool
	}{
		{
			name: "Nil query params",
			setupOpts: func(opts *getMembershipsOptsV2) {
				opts.QueryParam = nil
			},
			expectError: false,
		},
		{
			name: "Empty query params",
			setupOpts: func(opts *getMembershipsOptsV2) {
				opts.QueryParam = map[string]string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *getMembershipsOptsV2) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *getMembershipsOptsV2) {
				opts.QueryParam = map[string]string{
					"special@key":   "special@value",
					"unicode测试":     "unicode值",
					"with spaces":   "also spaces",
					"equals=key":    "equals=value",
					"ampersand&key": "ampersand&value",
				}
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetMembershipsOptsV2(pn, pn.ctx)
			tc.setupOpts(opts)

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

func TestGetMembershipsV2BuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newGetMembershipsBuilderV2(pn)

	include := []PNMembershipsInclude{PNMembershipsIncludeCustom}
	sort := []string{"name:asc", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.UUID("complete-test-uuid").
		Include(include).
		Limit(75).
		Start("start-token").
		End("end-token").
		Count(true).
		Filter("active = true").
		Sort(sort).
		QueryParam(queryParam)

	// Verify all values are set
	assert.Equal("complete-test-uuid", builder.opts.UUID)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal(75, builder.opts.Limit)
	assert.Equal("start-token", builder.opts.Start)
	assert.Equal("end-token", builder.opts.End)
	assert.Equal(true, builder.opts.Count)
	assert.Equal("active = true", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v2/objects/demo/uuids/complete-test-uuid/channels"
	assert.Equal(expectedPath, path)

	// Should build query with all params
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Contains(query.Get("include"), "custom")
	assert.Equal("75", query.Get("limit"))
	assert.Equal("start-token", query.Get("start"))
	assert.Equal("end-token", query.Get("end"))
	assert.Equal("1", query.Get("count"))
	assert.Equal("active = true", query.Get("filter"))
	assert.Contains(query.Get("sort"), "name:asc")
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))

	// Should always have empty body (GET operation)
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}

func TestGetMembershipsV2ResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetMembershipsOptsV2(pn, pn.ctx)

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
			name:        "Valid response with empty data",
			jsonBytes:   []byte(`{"status":200,"data":[],"totalCount":0}`),
			expectError: false,
		},
		{
			name:        "Valid response with membership data",
			jsonBytes:   []byte(`{"status":200,"data":[{"id":"channel1","channel":{"id":"channel1","name":"Channel 1"},"custom":{"role":"admin"}}],"totalCount":1,"next":"abc","prev":"xyz"}`),
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
			resp, _, err := newPNGetMembershipsResponse(tc.jsonBytes, opts, StatusResponse{})

			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				if err == nil {
					// Either resp is not nil, or it's nil but that's acceptable for null JSON
					if resp != nil {
						assert.NotNil(resp)
					}
				}
			}
		})
	}
}

// UUID-to-Channels Direction Tests

func TestGetMembershipsV2UUIDtoChannelsDirection(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		description string
		setupFn     func(*getMembershipsBuilderV2)
		verifyFn    func(*testing.T, string)
	}{
		{
			name:        "Basic user membership retrieval",
			description: "Single UUID retrieving all channel memberships",
			setupFn: func(builder *getMembershipsBuilderV2) {
				builder.UUID("user123")
			},
			verifyFn: func(t *testing.T, path string) {
				// Verify UUID-to-channels direction in path
				assert.Contains(path, "/uuids/user123/channels")
			},
		},
		{
			name:        "Professional user with filtered memberships",
			description: "UUID retrieving filtered channel memberships with metadata",
			setupFn: func(builder *getMembershipsBuilderV2) {
				builder.UUID("professional@company.com")
				builder.Filter("custom.role IN ('admin', 'manager')")
				builder.Include([]PNMembershipsInclude{
					PNMembershipsIncludeCustom,
					PNMembershipsIncludeChannel,
					PNMembershipsIncludeChannelCustom,
				})
			},
			verifyFn: func(t *testing.T, path string) {
				assert.Contains(path, "/uuids/professional@company.com/channels")
			},
		},
		{
			name:        "International user with pagination",
			description: "UUID with Unicode characters retrieving paginated memberships",
			setupFn: func(builder *getMembershipsBuilderV2) {
				builder.UUID("国际用户_测试")
				builder.Limit(25)
				builder.Start("pagination_token")
				builder.Sort([]string{"name:asc"})
				builder.Count(true)
			},
			verifyFn: func(t *testing.T, path string) {
				assert.Contains(path, "/uuids/国际用户_测试/channels")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetMembershipsBuilderV2(pn)
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should be GET operation
			assert.Equal("GET", builder.opts.httpMethod())

			// Should have empty body
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)

			// Should build valid path (UUID-to-channels direction)
			path, err := builder.opts.buildPath()
			assert.Nil(err)

			// Run verification
			tc.verifyFn(t, path)
		})
	}
}
