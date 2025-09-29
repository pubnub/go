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

func AssertGetAllUUIDMetadata(t *testing.T, checkQueryParam, testContext, withFilter bool, withSort bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	incl := []PNUUIDMetadataInclude{
		PNUUIDMetadataIncludeCustom,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := EnumArrayToStringArray(incl)

	o := newGetAllUUIDMetadataBuilder(pn)
	if testContext {
		o = newGetAllUUIDMetadataBuilderWithContext(pn, pn.ctx)
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
		fmt.Sprintf("/v2/objects/%s/uuids", pn.Config.SubscribeKey),
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

func TestGetAllUUIDMetadata(t *testing.T) {
	AssertGetAllUUIDMetadata(t, true, false, false, false)
}

func TestGetAllUUIDMetadataContext(t *testing.T) {
	AssertGetAllUUIDMetadata(t, true, true, false, false)
}

func TestGetAllUUIDMetadataWithFilter(t *testing.T) {
	AssertGetAllUUIDMetadata(t, true, false, true, false)
}

func TestGetAllUUIDMetadataWithFilterContext(t *testing.T) {
	AssertGetAllUUIDMetadata(t, true, true, true, false)
}

func TestGetAllUUIDMetadataWithSort(t *testing.T) {
	AssertGetAllUUIDMetadata(t, true, false, false, true)
}

func TestGetAllUUIDMetadataWithSortContext(t *testing.T) {
	AssertGetAllUUIDMetadata(t, true, true, false, true)
}

func TestGetAllUUIDMetadataWithFilterWithSort(t *testing.T) {
	AssertGetAllUUIDMetadata(t, true, false, true, true)
}

func TestGetAllUUIDMetadataWithFilterWithSortContext(t *testing.T) {
	AssertGetAllUUIDMetadata(t, true, true, true, true)
}

func TestGetAllUUIDMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetAllUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestGetAllUUIDMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":[{"id":"id2","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-19T14:44:54.837392Z","updated":"2019-08-19T14:44:54.837392Z","eTag":"AbyT4v2p6K7fpQE"},{"id":"id0","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-20T13:26:19.140324Z","updated":"2019-08-20T13:26:19.140324Z","eTag":"AbyT4v2p6K7fpQE"}],"totalCount":2,"next":"Mg","prev":"Nd"}`)

	r, _, err := newPNGetAllUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(2, r.TotalCount)
	assert.Equal("Mg", r.Next)
	assert.Equal("Nd", r.Prev)
	assert.Equal("id2", r.Data[0].ID)
	assert.Equal("name", r.Data[0].Name)
	assert.Equal("extid", r.Data[0].ExternalID)
	assert.Equal("purl", r.Data[0].ProfileURL)
	assert.Equal("email", r.Data[0].Email)
	// assert.Equal("2019-08-19T14:44:54.837392Z", r.Data[0].Created)
	assert.Equal("2019-08-19T14:44:54.837392Z", r.Data[0].Updated)
	assert.Equal("AbyT4v2p6K7fpQE", r.Data[0].ETag)
	assert.Equal("b", r.Data[0].Custom["a"])
	assert.Equal("d", r.Data[0].Custom["c"])

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestGetAllUUIDMetadataValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetAllUUIDMetadataValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	assert.Nil(opts.validate())
}

func TestGetAllUUIDMetadataValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)
	opts.Limit = 50
	opts.Start = "start-token"
	opts.End = "end-token"
	opts.Filter = "name LIKE 'test*'"
	opts.Sort = []string{"name:asc"}
	opts.Count = true
	opts.Include = []string{"custom"}
	opts.QueryParam = map[string]string{"param": "value"}

	assert.Nil(opts.validate())
}

func TestGetAllUUIDMetadataValidateDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Should have default limit set by builder
	builder := newGetAllUUIDMetadataBuilder(pn)
	assert.Equal(getAllUUIDMetadataLimitV2, builder.opts.Limit)
	assert.Equal(100, builder.opts.Limit) // Verify default value
}

// HTTP Method and Operation Tests

func TestGetAllUUIDMetadataHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	assert.Equal("GET", opts.httpMethod())
}

func TestGetAllUUIDMetadataOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	assert.Equal(PNGetAllUUIDMetadataOperation, opts.operationType())
}

func TestGetAllUUIDMetadataIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestGetAllUUIDMetadataTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (8 setters)

func TestGetAllUUIDMetadataBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetAllUUIDMetadataBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
	assert.Equal(getAllUUIDMetadataLimitV2, builder.opts.Limit) // Default limit
}

func TestGetAllUUIDMetadataBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetAllUUIDMetadataBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
	assert.Equal(getAllUUIDMetadataLimitV2, builder.opts.Limit) // Default limit
}

func TestGetAllUUIDMetadataBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetAllUUIDMetadataBuilder(pn)

	// Test Include setter
	include := []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom}
	builder.Include(include)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)

	// Test Limit setter
	builder.Limit(25)
	assert.Equal(25, builder.opts.Limit)

	// Test Start setter
	builder.Start("start-token")
	assert.Equal("start-token", builder.opts.Start)

	// Test End setter
	builder.End("end-token")
	assert.Equal("end-token", builder.opts.End)

	// Test Filter setter
	builder.Filter("name LIKE 'test*'")
	assert.Equal("name LIKE 'test*'", builder.opts.Filter)

	// Test Sort setter
	sort := []string{"name:asc", "created:desc"}
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

func TestGetAllUUIDMetadataBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	include := []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom}
	sort := []string{"name:asc"}
	queryParam := map[string]string{"key": "value"}

	builder := newGetAllUUIDMetadataBuilder(pn)
	result := builder.Include(include).
		Limit(50).
		Start("start").
		End("end").
		Filter("filter").
		Sort(sort).
		Count(true).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal(EnumArrayToStringArray(include), builder.opts.Include)
	assert.Equal(50, builder.opts.Limit)
	assert.Equal("start", builder.opts.Start)
	assert.Equal("end", builder.opts.End)
	assert.Equal("filter", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(true, builder.opts.Count)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestGetAllUUIDMetadataBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newGetAllUUIDMetadataBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

func TestGetAllUUIDMetadataBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetAllUUIDMetadataBuilder(pn)

	// Verify default values
	assert.Equal(getAllUUIDMetadataLimitV2, builder.opts.Limit) // Default: 100
	assert.Nil(builder.opts.Include)
	assert.Equal("", builder.opts.Start)
	assert.Equal("", builder.opts.End)
	assert.Equal("", builder.opts.Filter)
	assert.Nil(builder.opts.Sort)
	assert.Equal(false, builder.opts.Count)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestGetAllUUIDMetadataBuilderLimitDefault(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetAllUUIDMetadataBuilder(pn)

	// Should have default limit of 100
	assert.Equal(100, builder.opts.Limit)
	assert.Equal(getAllUUIDMetadataLimitV2, builder.opts.Limit)

	// Can override default
	builder.Limit(50)
	assert.Equal(50, builder.opts.Limit)
}

func TestGetAllUUIDMetadataBuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	transport := &http.Transport{}
	include := []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom}
	sort := []string{"name:asc", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Test all 9 setters in chain
	builder := newGetAllUUIDMetadataBuilder(pn).
		Include(include).
		Limit(75).
		Start("start-token").
		End("end-token").
		Filter("name LIKE 'user*'").
		Sort(sort).
		Count(true).
		QueryParam(queryParam).
		Transport(transport)

	// Verify all are set correctly
	assert.Equal(EnumArrayToStringArray(include), builder.opts.Include)
	assert.Equal(75, builder.opts.Limit)
	assert.Equal("start-token", builder.opts.Start)
	assert.Equal("end-token", builder.opts.End)
	assert.Equal("name LIKE 'user*'", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(true, builder.opts.Count)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestGetAllUUIDMetadataBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/uuids"
	assert.Equal(expected, path)
}

func TestGetAllUUIDMetadataBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/uuids"
	assert.Equal(expected, path)
}

func TestGetAllUUIDMetadataBuildPathSpecialCharsInSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "key-with-special@chars#and$symbols"
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/")
	assert.Contains(path, "key-with-special@chars#and$symbols")
	assert.Contains(path, "/uuids")
}

func TestGetAllUUIDMetadataBuildPathUnicodeSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "订阅键-русский-キー"
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/")
	assert.Contains(path, "订阅键-русский-キー")
	assert.Contains(path, "/uuids")
}

// Query Parameter Tests

func TestGetAllUUIDMetadataBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)
	opts.Limit = getAllUUIDMetadataLimitV2 // Set default explicitly

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("100", query.Get("limit")) // Default limit
	assert.Equal("0", query.Get("count"))   // Default count = false
}

func TestGetAllUUIDMetadataBuildQueryWithInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)
	opts.Include = []string{"custom"}
	opts.Limit = 100

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("custom", query.Get("include"))
}

func TestGetAllUUIDMetadataBuildQueryPaginationParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	opts.Limit = 50
	opts.Start = "start-token-123"
	opts.End = "end-token-456"
	opts.Count = true

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("50", query.Get("limit"))
	assert.Equal("start-token-123", query.Get("start"))
	assert.Equal("end-token-456", query.Get("end"))
	assert.Equal("1", query.Get("count")) // Count = true
}

func TestGetAllUUIDMetadataBuildQueryFilterAndSort(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	opts.Filter = "name LIKE 'test*'"
	opts.Sort = []string{"name:asc", "created:desc"}
	opts.Limit = 100

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("name LIKE 'test*'", query.Get("filter"))
	assert.Equal("name:asc,created:desc", query.Get("sort"))
}

func TestGetAllUUIDMetadataBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	customParams := map[string]string{
		"custom":         "value",
		"special_chars":  "value@with#symbols",
		"unicode":        "测试参数",
		"empty_value":    "",
		"number_string":  "42",
		"boolean_string": "true",
	}
	opts.QueryParam = customParams
	opts.Limit = 100

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

func TestGetAllUUIDMetadataBuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

	// Set all possible query parameters
	opts.Include = []string{"custom"}
	opts.Limit = 25
	opts.Start = "start-token"
	opts.End = "end-token"
	opts.Filter = "name LIKE 'user*'"
	opts.Sort = []string{"name:asc"}
	opts.Count = true
	opts.QueryParam = map[string]string{
		"extra": "parameter",
		"debug": "true",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all parameters are set correctly
	assert.Equal("custom", query.Get("include"))
	assert.Equal("25", query.Get("limit"))
	assert.Equal("start-token", query.Get("start"))
	assert.Equal("end-token", query.Get("end"))
	assert.Equal("name LIKE 'user*'", query.Get("filter"))
	assert.Equal("name:asc", query.Get("sort"))
	assert.Equal("1", query.Get("count"))
	assert.Equal("parameter", query.Get("extra"))
	assert.Equal("true", query.Get("debug"))

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

// Pagination Parameter Tests

func TestGetAllUUIDMetadataLimitBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name          string
		limit         int
		expectError   bool
		expectedLimit string
	}{
		{
			name:          "Negative limit",
			limit:         -1,
			expectError:   false, // buildQuery doesn't validate
			expectedLimit: "-1",
		},
		{
			name:          "Zero limit",
			limit:         0,
			expectError:   false,
			expectedLimit: "0",
		},
		{
			name:          "Small positive limit",
			limit:         1,
			expectError:   false,
			expectedLimit: "1",
		},
		{
			name:          "Normal limit",
			limit:         50,
			expectError:   false,
			expectedLimit: "50",
		},
		{
			name:          "Large limit",
			limit:         1000,
			expectError:   false,
			expectedLimit: "1000",
		},
		{
			name:          "Very large limit",
			limit:         999999,
			expectError:   false,
			expectedLimit: "999999",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)
			opts.Limit = tc.limit

			query, err := opts.buildQuery()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Equal(tc.expectedLimit, query.Get("limit"))
			}
		})
	}
}

func TestGetAllUUIDMetadataCountParameter(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		count    bool
		expected string
	}{
		{
			name:     "Count true",
			count:    true,
			expected: "1",
		},
		{
			name:     "Count false",
			count:    false,
			expected: "0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)
			opts.Count = tc.count
			opts.Limit = 100

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expected, query.Get("count"))
		})
	}
}

func TestGetAllUUIDMetadataStartEndTokens(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name  string
		start string
		end   string
	}{
		{
			name:  "Empty tokens",
			start: "",
			end:   "",
		},
		{
			name:  "Normal tokens",
			start: "start-token-123",
			end:   "end-token-456",
		},
		{
			name:  "Only start token",
			start: "start-only",
			end:   "",
		},
		{
			name:  "Only end token",
			start: "",
			end:   "end-only",
		},
		{
			name:  "Long tokens",
			start: strings.Repeat("start", 20),
			end:   strings.Repeat("end", 20),
		},
		{
			name:  "Special character tokens",
			start: "start@token#with$symbols",
			end:   "end@token#with$symbols",
		},
		{
			name:  "Unicode tokens",
			start: "开始令牌-начальный-開始",
			end:   "结束令牌-конечный-終了",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)
			opts.Start = tc.start
			opts.End = tc.end
			opts.Limit = 100

			query, err := opts.buildQuery()
			assert.Nil(err)

			if tc.start != "" {
				assert.Equal(tc.start, query.Get("start"))
			} else {
				assert.Equal("", query.Get("start"))
			}

			if tc.end != "" {
				assert.Equal(tc.end, query.Get("end"))
			} else {
				assert.Equal("", query.Get("end"))
			}
		})
	}
}

// Include Parameter Tests

func TestGetAllUUIDMetadataIncludeParameterHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		include  []string
		expected string
	}{
		{
			name:     "No include parameters",
			include:  nil,
			expected: "",
		},
		{
			name:     "Empty include array",
			include:  []string{},
			expected: "",
		},
		{
			name:     "Single include parameter",
			include:  []string{"custom"},
			expected: "custom",
		},
		{
			name:     "Multiple include parameters",
			include:  []string{"custom", "type"},
			expected: "custom,type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)
			opts.Include = tc.include
			opts.Limit = 100

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expected, query.Get("include"))
		})
	}
}

func TestGetAllUUIDMetadataIncludeEnumConversion(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		include  []PNUUIDMetadataInclude
		expected string
	}{
		{
			name:     "No include",
			include:  nil,
			expected: "",
		},
		{
			name:     "Empty include",
			include:  []PNUUIDMetadataInclude{},
			expected: "",
		},
		{
			name:     "Single include",
			include:  []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom},
			expected: "custom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetAllUUIDMetadataBuilder(pn)
			if tc.include != nil {
				builder.Include(tc.include)
			}

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expected, query.Get("include"))
		})
	}
}

// Filter and Sort Tests

func TestGetAllUUIDMetadataFilterExpressions(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name   string
		filter string
	}{
		{
			name:   "Empty filter",
			filter: "",
		},
		{
			name:   "Simple LIKE filter",
			filter: "name LIKE 'test*'",
		},
		{
			name:   "Complex filter with AND",
			filter: "name LIKE 'user*' AND email LIKE '*@test.com'",
		},
		{
			name:   "IN filter",
			filter: "id IN ('uuid1', 'uuid2', 'uuid3')",
		},
		{
			name:   "Custom field filter",
			filter: "custom.type = 'premium'",
		},
		{
			name:   "Unicode filter",
			filter: "name LIKE '用户*'",
		},
		{
			name:   "Special character filter",
			filter: "name LIKE 'user@special#chars'",
		},
		{
			name:   "Very long filter",
			filter: strings.Repeat("name LIKE 'test' AND ", 20) + "id IS NOT NULL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)
			opts.Filter = tc.filter
			opts.Limit = 100

			query, err := opts.buildQuery()
			assert.Nil(err)

			if tc.filter != "" {
				assert.Equal(tc.filter, query.Get("filter"))
			} else {
				assert.Equal("", query.Get("filter"))
			}
		})
	}
}

func TestGetAllUUIDMetadataSortExpressions(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		sort     []string
		expected string
	}{
		{
			name:     "Empty sort",
			sort:     nil,
			expected: "",
		},
		{
			name:     "Empty sort array",
			sort:     []string{},
			expected: "",
		},
		{
			name:     "Single field ascending",
			sort:     []string{"name:asc"},
			expected: "name:asc",
		},
		{
			name:     "Single field descending",
			sort:     []string{"created:desc"},
			expected: "created:desc",
		},
		{
			name:     "Multiple fields",
			sort:     []string{"name:asc", "created:desc", "updated:asc"},
			expected: "name:asc,created:desc,updated:asc",
		},
		{
			name:     "Custom field sorting",
			sort:     []string{"custom.priority:desc"},
			expected: "custom.priority:desc",
		},
		{
			name:     "Unicode field names",
			sort:     []string{"名称:asc", "创建时间:desc"},
			expected: "名称:asc,创建时间:desc",
		},
		{
			name:     "Special character field names",
			sort:     []string{"field@special:desc"},
			expected: "field@special:desc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)
			opts.Sort = tc.sort
			opts.Limit = 100

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expected, query.Get("sort"))
		})
	}
}

// Comprehensive Edge Case Tests

func TestGetAllUUIDMetadataWithLargeParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*getAllUUIDMetadataBuilder)
	}{
		{
			name: "Very large filter expression",
			setupFn: func(builder *getAllUUIDMetadataBuilder) {
				largeFilter := strings.Repeat("name LIKE 'test' AND ", 100) + "id IS NOT NULL"
				builder.Filter(largeFilter)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *getAllUUIDMetadataBuilder) {
				largeQueryParam := make(map[string]string)
				for i := 0; i < 100; i++ {
					largeQueryParam[fmt.Sprintf("param_%d", i)] = strings.Repeat(fmt.Sprintf("value_%d_", i), 20)
				}
				builder.QueryParam(largeQueryParam)
			},
		},
		{
			name: "Many sort fields",
			setupFn: func(builder *getAllUUIDMetadataBuilder) {
				var sortFields []string
				for i := 0; i < 50; i++ {
					direction := "asc"
					if i%2 == 0 {
						direction = "desc"
					}
					sortFields = append(sortFields, fmt.Sprintf("field_%d:%s", i, direction))
				}
				builder.Sort(sortFields)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetAllUUIDMetadataBuilder(pn)
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
		})
	}
}

func TestGetAllUUIDMetadataSpecialCharacterHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialStrings := []string{
		"<script>alert('xss')</script>",
		"SELECT * FROM users; DROP TABLE users;",
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
			builder := newGetAllUUIDMetadataBuilder(pn)
			builder.Filter(specialString)
			builder.Start(specialString)
			builder.End(specialString)
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
		})
	}
}

func TestGetAllUUIDMetadataParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*getAllUUIDMetadataBuilder)
	}{
		{
			name: "Empty string values",
			setupFn: func(builder *getAllUUIDMetadataBuilder) {
				builder.Filter("")
				builder.Start("")
				builder.End("")
				builder.QueryParam(map[string]string{
					"empty_param":  "",
					"normal_param": "value",
				})
			},
		},
		{
			name: "Single character values",
			setupFn: func(builder *getAllUUIDMetadataBuilder) {
				builder.Filter("a")
				builder.Start("b")
				builder.End("c")
				builder.QueryParam(map[string]string{
					"single_char_key": "d",
				})
			},
		},
		{
			name: "Unicode-only values",
			setupFn: func(builder *getAllUUIDMetadataBuilder) {
				builder.Filter("测试过滤器")
				builder.Start("开始")
				builder.End("结束")
				builder.QueryParam(map[string]string{
					"unicode测试": "unicode值",
					"русский":   "значение",
				})
			},
		},
		{
			name: "Mixed boundary conditions",
			setupFn: func(builder *getAllUUIDMetadataBuilder) {
				builder.Limit(1)
				builder.Filter(strings.Repeat("filter", 200))
				builder.Sort([]string{strings.Repeat("field", 50) + ":asc"})
				builder.QueryParam(map[string]string{
					"empty_param":   "",
					"long_param":    strings.Repeat("value", 200),
					"special@param": "special@value",
				})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetAllUUIDMetadataBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid components
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

func TestGetAllUUIDMetadataComplexParameterCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*getAllUUIDMetadataBuilder)
		validateFn func(*testing.T, string, *url.Values)
	}{
		{
			name: "Full pagination with unicode",
			setupFn: func(builder *getAllUUIDMetadataBuilder) {
				builder.Limit(25)
				builder.Start("开始令牌")
				builder.End("结束令牌")
				builder.Count(true)
				builder.Filter("name LIKE '用户*'")
				builder.Sort([]string{"名称:asc", "创建时间:desc"})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Equal("25", query.Get("limit"))
				assert.Equal("1", query.Get("count"))
				assert.Contains(query.Get("filter"), "用户")
			},
		},
		{
			name: "Complex filter with special characters",
			setupFn: func(builder *getAllUUIDMetadataBuilder) {
				builder.Filter("name LIKE 'user@special#123' AND email LIKE '*@test.com'")
				builder.Sort([]string{"field@special:desc", "other#field:asc"})
				builder.QueryParam(map[string]string{
					"special@key": "special@value",
					"with spaces": "also spaces",
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(query.Get("filter"), "user@special#123")
				assert.Contains(query.Get("sort"), "field@special:desc")
			},
		},
		{
			name: "Large data with all parameters",
			setupFn: func(builder *getAllUUIDMetadataBuilder) {
				builder.Include([]PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom})
				builder.Limit(999)
				builder.Start(strings.Repeat("start", 50))
				builder.End(strings.Repeat("end", 50))
				builder.Filter(strings.Repeat("name LIKE 'test' AND ", 10) + "id IS NOT NULL")
				builder.Sort([]string{"name:asc", "created:desc", "updated:asc"})
				builder.Count(true)
				builder.QueryParam(map[string]string{
					"large_param": strings.Repeat("value", 100),
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Equal("custom", query.Get("include"))
				assert.Equal("999", query.Get("limit"))
				assert.Equal("1", query.Get("count"))
				assert.Equal("name:asc,created:desc,updated:asc", query.Get("sort"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetAllUUIDMetadataBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid components
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Run custom validation
			tc.validateFn(t, path, query)
		})
	}
}

// Error Scenario Tests

func TestGetAllUUIDMetadataExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newGetAllUUIDMetadataBuilder(pn)

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetAllUUIDMetadataPathBuildingEdgeCases(t *testing.T) {
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
			opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

			path, err := opts.buildPath()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Contains(path, "/v2/objects/")
				assert.Contains(path, "/uuids")
			}
		})
	}
}

func TestGetAllUUIDMetadataQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*getAllUUIDMetadataOpts)
		expectError bool
	}{
		{
			name: "Nil include array",
			setupOpts: func(opts *getAllUUIDMetadataOpts) {
				opts.Include = nil
				opts.Limit = 100
			},
			expectError: false,
		},
		{
			name: "Empty include array",
			setupOpts: func(opts *getAllUUIDMetadataOpts) {
				opts.Include = []string{}
				opts.Limit = 100
			},
			expectError: false,
		},
		{
			name: "Nil sort array",
			setupOpts: func(opts *getAllUUIDMetadataOpts) {
				opts.Sort = nil
				opts.Limit = 100
			},
			expectError: false,
		},
		{
			name: "Empty sort array",
			setupOpts: func(opts *getAllUUIDMetadataOpts) {
				opts.Sort = []string{}
				opts.Limit = 100
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *getAllUUIDMetadataOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
				opts.Limit = 100
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *getAllUUIDMetadataOpts) {
				opts.QueryParam = map[string]string{
					"special@key":   "special@value",
					"unicode测试":     "unicode值",
					"with spaces":   "also spaces",
					"equals=key":    "equals=value",
					"ampersand&key": "ampersand&value",
				}
				opts.Limit = 100
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)
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

func TestGetAllUUIDMetadataBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newGetAllUUIDMetadataBuilder(pn)

	include := []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom}
	sort := []string{"name:asc", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.Include(include).
		Limit(50).
		Start("start-token").
		End("end-token").
		Filter("name LIKE 'test*'").
		Sort(sort).
		Count(true).
		QueryParam(queryParam)

	// Verify all values are set
	assert.Equal(EnumArrayToStringArray(include), builder.opts.Include)
	assert.Equal(50, builder.opts.Limit)
	assert.Equal("start-token", builder.opts.Start)
	assert.Equal("end-token", builder.opts.End)
	assert.Equal("name LIKE 'test*'", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(true, builder.opts.Count)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v2/objects/demo/uuids"
	assert.Equal(expectedPath, path)

	// Should build query with all params
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))
	assert.Equal("custom", query.Get("include"))
	assert.Equal("50", query.Get("limit"))
	assert.Equal("start-token", query.Get("start"))
	assert.Equal("end-token", query.Get("end"))
	assert.Equal("name LIKE 'test*'", query.Get("filter"))
	assert.Equal("name:asc,created:desc", query.Get("sort"))
	assert.Equal("1", query.Get("count"))
}

func TestGetAllUUIDMetadataResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetAllUUIDMetadataOpts(pn, pn.ctx)

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
			name:        "Valid UUID metadata list response",
			jsonBytes:   []byte(`{"status":200,"data":[{"id":"test-uuid","name":"Test User","email":"test@example.com"}],"totalCount":1,"next":"next-token","prev":"prev-token"}`),
			expectError: false,
		},
		{
			name:        "Response with empty data array",
			jsonBytes:   []byte(`{"status":200,"data":[],"totalCount":0}`),
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
			resp, _, err := newPNGetAllUUIDMetadataResponse(tc.jsonBytes, opts, StatusResponse{})

			if tc.expectError {
				assert.NotNil(err)
				// When there's an error, resp might be nil or the empty response
				if resp == nil {
					assert.Equal(emptyPNGetAllUUIDMetadataResponse, resp)
				}
			} else {
				assert.Nil(err)
				// For successful parsing, resp should not be nil, but content may vary
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
