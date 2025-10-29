package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/pubnub/go/v8/utils"
	"github.com/stretchr/testify/assert"
)

func AssertGetAllChannelMetadata(t *testing.T, checkQueryParam, testContext, withFilter bool, withSort bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	incl := []PNChannelMetadataInclude{
		PNChannelMetadataIncludeCustom,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := EnumArrayToStringArray(incl)

	o := newGetAllChannelMetadataBuilder(pn)
	if testContext {
		o = newGetAllChannelMetadataBuilderWithContext(pn, pn.ctx)
	}

	limit := 90
	start := "Mxmy"
	end := "Nxny"

	o.Include(incl)
	o.Limit(limit)
	o.Start(start)
	o.End(end)
	o.Count(false)
	o.QueryParam(queryParam)
	if withFilter {
		o.Filter("name like 'a*'")
	}
	sort := []string{"name", "created:desc"}
	if withSort {
		o.Sort(sort)
	}

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/channels", pn.Config.SubscribeKey),
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
			assert.Equal("name like 'a*'", u.Get("filter"))
		}
		if withSort {
			v := &url.Values{}
			SetQueryParamAsCommaSepString(v, sort, "sort")
			assert.Equal(v.Get("sort"), u.Get("sort"))
		}

	}

}

func TestGetAllChannelMetadata(t *testing.T) {
	AssertGetAllChannelMetadata(t, true, false, false, false)
}

func TestGetAllChannelMetadataContext(t *testing.T) {
	AssertGetAllChannelMetadata(t, true, true, false, false)
}

func TestGetAllChannelMetadataWithFilter(t *testing.T) {
	AssertGetAllChannelMetadata(t, true, false, true, false)
}

func TestGetAllChannelMetadataWithFilterContext(t *testing.T) {
	AssertGetAllChannelMetadata(t, true, true, true, false)
}

func TestGetAllChannelMetadataWithSort(t *testing.T) {
	AssertGetAllChannelMetadata(t, true, false, false, true)
}

func TestGetAllChannelMetadataWithSortContext(t *testing.T) {
	AssertGetAllChannelMetadata(t, true, true, false, true)
}

func TestGetAllChannelMetadataWithFilterWithSort(t *testing.T) {
	AssertGetAllChannelMetadata(t, true, false, true, true)
}

func TestGetAllChannelMetadataWithFilterWithSortContext(t *testing.T) {
	AssertGetAllChannelMetadata(t, true, true, true, true)
}

func TestGetAllChannelMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetAllChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestGetAllChannelMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":[{"id":"id0","name":"name","description":"desc","custom":{"a":"b"},"status":"active","type":"public","created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T13:26:08.341297Z","eTag":"Aee9zsKNndXlHw"},{"id":"id01","name":"name","description":"desc","custom":{"a":"b"},"status":"inactive","type":"private","created":"2019-08-20T14:44:52.799969Z","updated":"2019-08-20T14:44:52.799969Z","eTag":"Aee9zsKNndXlHw"}],"totalCount":2,"next":"Mg","prev":"Nd"}`)

	r, _, err := newPNGetAllChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(2, r.TotalCount)
	assert.Equal("Mg", r.Next)
	assert.Equal("Nd", r.Prev)
	assert.Equal("id0", r.Data[0].ID)
	assert.Equal("name", r.Data[0].Name)
	assert.Equal("desc", r.Data[0].Description)
	//assert.Equal("2019-08-20T13:26:08.341297Z", r.Data[0].Created)
	assert.Equal("2019-08-20T13:26:08.341297Z", r.Data[0].Updated)
	assert.Equal("Aee9zsKNndXlHw", r.Data[0].ETag)
	assert.Equal("b", r.Data[0].Custom["a"])
	assert.Equal("active", r.Data[0].Status)
	assert.Equal("public", r.Data[0].Type)

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestGetAllChannelMetadataValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetAllChannelMetadataValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)

	assert.Nil(opts.validate())
}

func TestGetAllChannelMetadataValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
	opts.Limit = 50
	opts.Include = []string{"custom"}
	opts.Start = "start_cursor"
	opts.End = "end_cursor"
	opts.Filter = "name LIKE 'test*'"
	opts.Sort = []string{"name", "created:desc"}
	opts.Count = true
	opts.QueryParam = map[string]string{"param": "value"}

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestGetAllChannelMetadataHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)

	assert.Equal("GET", opts.httpMethod())
}

func TestGetAllChannelMetadataOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)

	assert.Equal(PNGetAllChannelMetadataOperation, opts.operationType())
}

func TestGetAllChannelMetadataIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestGetAllChannelMetadataTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (9 setters)

func TestGetAllChannelMetadataBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetAllChannelMetadataBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
	// Should have default limit
	assert.Equal(getAllChannelMetadataLimitV2, builder.opts.Limit)
}

func TestGetAllChannelMetadataBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetAllChannelMetadataBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
	// Should have default limit
	assert.Equal(getAllChannelMetadataLimitV2, builder.opts.Limit)
}

func TestGetAllChannelMetadataBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetAllChannelMetadataBuilder(pn)

	// Test Include setter
	include := []PNChannelMetadataInclude{
		PNChannelMetadataIncludeCustom,
	}
	builder.Include(include)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)

	// Test Limit setter
	builder.Limit(50)
	assert.Equal(50, builder.opts.Limit)

	// Test Start setter
	builder.Start("start_cursor")
	assert.Equal("start_cursor", builder.opts.Start)

	// Test End setter
	builder.End("end_cursor")
	assert.Equal("end_cursor", builder.opts.End)

	// Test Filter setter
	builder.Filter("name LIKE 'test*'")
	assert.Equal("name LIKE 'test*'", builder.opts.Filter)

	// Test Sort setter
	sort := []string{"name", "created:desc"}
	builder.Sort(sort)
	assert.Equal(sort, builder.opts.Sort)

	// Test Count setter
	builder.Count(true)
	assert.Equal(true, builder.opts.Count)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestGetAllChannelMetadataBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	include := []PNChannelMetadataInclude{PNChannelMetadataIncludeCustom}
	sort := []string{"name:asc", "created:desc"}
	queryParam := map[string]string{"key": "value"}

	builder := newGetAllChannelMetadataBuilder(pn)
	result := builder.Include(include).
		Limit(25).
		Start("start").
		End("end").
		Filter("name LIKE 'a*'").
		Sort(sort).
		Count(true).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal(EnumArrayToStringArray(include), builder.opts.Include)
	assert.Equal(25, builder.opts.Limit)
	assert.Equal("start", builder.opts.Start)
	assert.Equal("end", builder.opts.End)
	assert.Equal("name LIKE 'a*'", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(true, builder.opts.Count)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestGetAllChannelMetadataBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newGetAllChannelMetadataBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

func TestGetAllChannelMetadataBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetAllChannelMetadataBuilder(pn)

	// Verify default values
	assert.Equal(getAllChannelMetadataLimitV2, builder.opts.Limit) // Default 100
	assert.Equal("", builder.opts.Start)
	assert.Equal("", builder.opts.End)
	assert.Equal("", builder.opts.Filter)
	assert.Nil(builder.opts.Sort)
	assert.Equal(false, builder.opts.Count)
	assert.Nil(builder.opts.Include)
	assert.Nil(builder.opts.QueryParam)
}

// URL/Path Building Tests

func TestGetAllChannelMetadataBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/channels"
	assert.Equal(expected, path)
}

func TestGetAllChannelMetadataBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/channels"
	assert.Equal(expected, path)
}

func TestGetAllChannelMetadataBuildPathWithSpecialCharsInSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "sub-key-with-special@chars#and$symbols"
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/")
	assert.Contains(path, "/channels")
	// SubscribeKey should be included in path
	assert.Contains(path, "sub-key-with-special@chars#and$symbols")
}

// Complex Query Parameter Tests

func TestGetAllChannelMetadataBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
	opts.Limit = getAllChannelMetadataLimitV2 // Default

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal(strconv.Itoa(getAllChannelMetadataLimitV2), query.Get("limit"))
	assert.Equal("0", query.Get("count")) // Default false
	assert.Equal("", query.Get("include"))
	assert.Equal("", query.Get("start"))
	assert.Equal("", query.Get("end"))
	assert.Equal("", query.Get("filter"))
	assert.Equal("", query.Get("sort"))
}

func TestGetAllChannelMetadataBuildQueryWithInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
	opts.Include = []string{"custom", "status", "type"}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("custom,status,type", query.Get("include"))
}

func TestGetAllChannelMetadataBuildQueryWithPagination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
	opts.Limit = 25
	opts.Start = "start_cursor_123"
	opts.End = "end_cursor_456"

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("25", query.Get("limit"))
	assert.Equal("start_cursor_123", query.Get("start"))
	assert.Equal("end_cursor_456", query.Get("end"))
}

func TestGetAllChannelMetadataBuildQueryWithCount(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test Count = true
	opts1 := newGetAllChannelMetadataOpts(pn, pn.ctx)
	opts1.Count = true

	query1, err1 := opts1.buildQuery()
	assert.Nil(err1)
	assert.Equal("1", query1.Get("count"))

	// Test Count = false (default)
	opts2 := newGetAllChannelMetadataOpts(pn, pn.ctx)
	opts2.Count = false

	query2, err2 := opts2.buildQuery()
	assert.Nil(err2)
	assert.Equal("0", query2.Get("count"))
}

func TestGetAllChannelMetadataBuildQueryWithFilter(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
	opts.Filter = "name LIKE 'test*' AND custom.type = 'important'"

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("name LIKE 'test*' AND custom.type = 'important'", query.Get("filter"))
}

func TestGetAllChannelMetadataBuildQueryWithSort(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
	opts.Sort = []string{"name:asc", "created:desc", "updated"}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("name:asc,created:desc,updated", query.Get("sort"))
}

func TestGetAllChannelMetadataBuildQueryLimitBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		limit    int
		expected string
	}{
		{
			name:     "Default limit",
			limit:    getAllChannelMetadataLimitV2,
			expected: "100",
		},
		{
			name:     "Small limit",
			limit:    1,
			expected: "1",
		},
		{
			name:     "Custom limit",
			limit:    50,
			expected: "50",
		},
		{
			name:     "Large limit",
			limit:    1000,
			expected: "1000",
		},
		{
			name:     "Zero limit",
			limit:    0,
			expected: "0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
			opts.Limit = tc.limit

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expected, query.Get("limit"))
		})
	}
}

func TestGetAllChannelMetadataBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)

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

func TestGetAllChannelMetadataBuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)

	// Set all possible query parameters
	opts.Include = []string{"custom", "status", "type"}
	opts.Limit = 75
	opts.Start = "start_123"
	opts.End = "end_456"
	opts.Filter = "name LIKE 'a*'"
	opts.Sort = []string{"name:asc", "created:desc"}
	opts.Count = true
	opts.QueryParam = map[string]string{
		"extra": "parameter",
		"debug": "true",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all parameters are set correctly
	assert.Equal("custom,status,type", query.Get("include"))
	assert.Equal("75", query.Get("limit"))
	assert.Equal("start_123", query.Get("start"))
	assert.Equal("end_456", query.Get("end"))
	assert.Equal("name LIKE 'a*'", query.Get("filter"))
	assert.Equal("name:asc,created:desc", query.Get("sort"))
	assert.Equal("1", query.Get("count"))
	assert.Equal("parameter", query.Get("extra"))
	assert.Equal("true", query.Get("debug"))

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

// Comprehensive Edge Case Tests

func TestGetAllChannelMetadataWithUnicodeFiltersAndSorts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
	opts.Filter = "name LIKE '测试*' OR description LIKE 'русский*'"
	opts.Sort = []string{"测试字段:asc", "русское_поле:desc"}

	// Should pass validation
	assert.Nil(opts.validate())

	// Should build valid query
	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Filter and sort should be included (possibly URL encoded)
	filterValue := query.Get("filter")
	sortValue := query.Get("sort")
	assert.NotEmpty(filterValue)
	assert.NotEmpty(sortValue)
}

func TestGetAllChannelMetadataWithComplexFilterSyntax(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	complexFilters := []string{
		"name LIKE 'test*'",
		"name LIKE 'test*' AND custom.type = 'important'",
		"(name LIKE 'a*' OR name LIKE 'b*') AND updated > '2023-01-01'",
		"custom.priority > 5 AND custom.category IN ('urgent', 'high')",
		"name LIKE '%special%' AND description NOT LIKE '%deprecated%'",
		"custom.tags CONTAINS 'production' AND updated >= '2023-01-01T00:00:00Z'",
	}

	for i, filter := range complexFilters {
		t.Run(fmt.Sprintf("ComplexFilter_%d", i), func(t *testing.T) {
			opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
			opts.Filter = filter

			// Should pass validation
			assert.Nil(opts.validate())

			// Should build valid query
			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(filter, query.Get("filter"))
		})
	}
}

func TestGetAllChannelMetadataWithComplexSortCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	sortCombinations := [][]string{
		{"name"},
		{"name:asc"},
		{"name:desc"},
		{"name:asc", "created:desc"},
		{"name", "created", "updated"},
		{"custom.priority:desc", "name:asc"},
		{"created:desc", "updated:desc", "name:asc"},
		{"custom.type", "custom.priority:desc", "name"},
	}

	for i, sort := range sortCombinations {
		t.Run(fmt.Sprintf("SortCombination_%d", i), func(t *testing.T) {
			opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
			opts.Sort = sort

			// Should pass validation
			assert.Nil(opts.validate())

			// Should build valid query
			query, err := opts.buildQuery()
			assert.Nil(err)

			expectedSort := ""
			for j, field := range sort {
				if j > 0 {
					expectedSort += ","
				}
				expectedSort += field
			}
			assert.Equal(expectedSort, query.Get("sort"))
		})
	}
}

func TestGetAllChannelMetadataWithLargePaginationValues(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test very large limit
	builder := newGetAllChannelMetadataBuilder(pn)
	builder.Limit(999999)

	assert.Equal(999999, builder.opts.Limit)

	// Test very long cursors
	longCursor := ""
	for i := 0; i < 1000; i++ {
		longCursor += fmt.Sprintf("segment_%d_", i)
	}

	builder.Start(longCursor)
	builder.End(longCursor + "_end")

	assert.Equal(longCursor, builder.opts.Start)
	assert.Equal(longCursor+"_end", builder.opts.End)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build valid query
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("999999", query.Get("limit"))
	assert.Equal(longCursor, query.Get("start"))
	assert.Equal(longCursor+"_end", query.Get("end"))
}

func TestGetAllChannelMetadataSpecialCharacterHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialStrings := []string{
		"!@#$%^&*()_+-=[]{}|;':\",./<>?",
		"測試字符串-русская строка-テスト文字列",
		"<script>alert('xss')</script>",
		"SELECT * FROM channels; DROP TABLE channels;",
		"newline\ncharacter\ttab\rcarriage",
		"   ",                // Only spaces
		"\u0000\u0001\u0002", // Control characters
	}

	for i, specialString := range specialStrings {
		t.Run(fmt.Sprintf("SpecialString_%d", i), func(t *testing.T) {
			opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
			opts.Filter = "name LIKE '" + specialString + "'"
			opts.Start = specialString
			opts.End = specialString + "_end"
			opts.QueryParam = map[string]string{
				"special_field": specialString,
			}

			// Should pass validation
			assert.Nil(opts.validate())

			// Should build valid query
			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

func TestGetAllChannelMetadataParameterCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		limit      int
		include    []string
		start      string
		end        string
		filter     string
		sort       []string
		count      bool
		queryParam map[string]string
	}{
		{
			name:  "Minimal - defaults only",
			limit: getAllChannelMetadataLimitV2,
		},
		{
			name:  "Pagination only",
			limit: 25,
			start: "cursor_start",
			end:   "cursor_end",
		},
		{
			name:   "Filtering only",
			limit:  getAllChannelMetadataLimitV2,
			filter: "name LIKE 'test*'",
		},
		{
			name:  "Sorting only",
			limit: getAllChannelMetadataLimitV2,
			sort:  []string{"name:asc", "created:desc"},
		},
		{
			name:    "Include with count",
			limit:   getAllChannelMetadataLimitV2,
			include: []string{"custom", "status", "type"},
			count:   true,
		},
		{
			name:    "Complete - all parameters",
			limit:   50,
			include: []string{"custom"},
			start:   "start_cursor",
			end:     "end_cursor",
			filter:  "name LIKE 'a*' AND custom.type = 'important'",
			sort:    []string{"name:asc", "created:desc"},
			count:   true,
			queryParam: map[string]string{
				"debug": "true",
				"extra": "param",
			},
		},
		{
			name:    "Unicode everything",
			limit:   30,
			include: []string{"custom"},
			start:   "测试开始",
			end:     "测试结束",
			filter:  "name LIKE '测试*'",
			sort:    []string{"测试字段:asc"},
			count:   false,
			queryParam: map[string]string{
				"unicode_param": "unicode值",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetAllChannelMetadataBuilder(pn)
			builder.Limit(tc.limit)

			if tc.include != nil {
				builder.opts.Include = tc.include
			}
			if tc.start != "" {
				builder.Start(tc.start)
			}
			if tc.end != "" {
				builder.End(tc.end)
			}
			if tc.filter != "" {
				builder.Filter(tc.filter)
			}
			if tc.sort != nil {
				builder.Sort(tc.sort)
			}
			builder.Count(tc.count)
			if tc.queryParam != nil {
				builder.QueryParam(tc.queryParam)
			}

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.Equal("/v2/objects/demo/channels", path)

			// Should build valid query
			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

// Error Scenario Tests

func TestGetAllChannelMetadataExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newGetAllChannelMetadataBuilder(pn)

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetAllChannelMetadataPathBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCases := []struct {
		name         string
		subscribeKey string
		expectError  bool
	}{
		{
			name:         "Empty SubscribeKey",
			subscribeKey: "",
			expectError:  false, // buildPath doesn't validate, only validate() does
		},
		{
			name:         "SubscribeKey with spaces",
			subscribeKey: "   sub key   ",
			expectError:  false,
		},
		{
			name:         "Very special characters",
			subscribeKey: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			expectError:  false,
		},
		{
			name:         "Unicode SubscribeKey",
			subscribeKey: "测试订阅键-русский-キー",
			expectError:  false,
		},
		{
			name:         "Very long SubscribeKey",
			subscribeKey: strings.Repeat("a", 1000),
			expectError:  false,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			pn.Config.SubscribeKey = tc.subscribeKey
			opts := newGetAllChannelMetadataOpts(pn, pn.ctx)

			path, err := opts.buildPath()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Contains(path, "/v2/objects/")
				assert.Contains(path, "/channels")
			}
		})
	}
}

func TestGetAllChannelMetadataQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*getAllChannelMetadataOpts)
		expectError bool
	}{
		{
			name: "Nil include array",
			setupOpts: func(opts *getAllChannelMetadataOpts) {
				opts.Include = nil
			},
			expectError: false,
		},
		{
			name: "Empty include array",
			setupOpts: func(opts *getAllChannelMetadataOpts) {
				opts.Include = []string{}
			},
			expectError: false,
		},
		{
			name: "Nil sort array",
			setupOpts: func(opts *getAllChannelMetadataOpts) {
				opts.Sort = nil
			},
			expectError: false,
		},
		{
			name: "Empty sort array",
			setupOpts: func(opts *getAllChannelMetadataOpts) {
				opts.Sort = []string{}
			},
			expectError: false,
		},
		{
			name: "Negative limit",
			setupOpts: func(opts *getAllChannelMetadataOpts) {
				opts.Limit = -10
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *getAllChannelMetadataOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *getAllChannelMetadataOpts) {
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
			opts := newGetAllChannelMetadataOpts(pn, pn.ctx)
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

func TestGetAllChannelMetadataBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newGetAllChannelMetadataBuilder(pn)

	include := []PNChannelMetadataInclude{PNChannelMetadataIncludeCustom}
	sort := []string{"name:asc", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.Include(include).
		Limit(75).
		Start("start_cursor_123").
		End("end_cursor_456").
		Filter("name LIKE 'complete*'").
		Sort(sort).
		Count(true).
		QueryParam(queryParam)

	// Verify all values are set
	assert.Equal(EnumArrayToStringArray(include), builder.opts.Include)
	assert.Equal(75, builder.opts.Limit)
	assert.Equal("start_cursor_123", builder.opts.Start)
	assert.Equal("end_cursor_456", builder.opts.End)
	assert.Equal("name LIKE 'complete*'", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(true, builder.opts.Count)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v2/objects/demo/channels"
	assert.Equal(expectedPath, path)

	// Should build query with all params
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))
	assert.Equal("custom", query.Get("include"))
	assert.Equal("75", query.Get("limit"))
	assert.Equal("start_cursor_123", query.Get("start"))
	assert.Equal("end_cursor_456", query.Get("end"))
	assert.Equal("name LIKE 'complete*'", query.Get("filter"))
	assert.Equal("name:asc,created:desc", query.Get("sort"))
	assert.Equal("1", query.Get("count"))
}

func TestGetAllChannelMetadataPaginationEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name  string
		start string
		end   string
		limit int
	}{
		{
			name:  "Only start cursor",
			start: "start_only",
			end:   "",
			limit: 50,
		},
		{
			name:  "Only end cursor",
			start: "",
			end:   "end_only",
			limit: 50,
		},
		{
			name:  "Both cursors empty",
			start: "",
			end:   "",
			limit: 50,
		},
		{
			name:  "Same start and end",
			start: "same_cursor",
			end:   "same_cursor",
			limit: 50,
		},
		{
			name:  "Very long cursors",
			start: strings.Repeat("start_", 100),
			end:   strings.Repeat("end_", 100),
			limit: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetAllChannelMetadataBuilder(pn)
			builder.Limit(tc.limit).Start(tc.start).End(tc.end)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid query
			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.Equal(strconv.Itoa(tc.limit), query.Get("limit"))
			assert.Equal(tc.start, query.Get("start"))
			assert.Equal(tc.end, query.Get("end"))
		})
	}
}
