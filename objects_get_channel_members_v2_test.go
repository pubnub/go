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

func AssertGetChannelMembersV2(t *testing.T, checkQueryParam, testContext, withFilter bool, withSort bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	incl := []PNChannelMembersInclude{
		PNChannelMembersIncludeCustom,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := EnumArrayToStringArray(incl)

	o := newGetChannelMembersBuilderV2(pn)
	if testContext {
		o = newGetChannelMembersBuilderV2WithContext(pn, pn.ctx)
	}

	spaceID := "id0"
	limit := 90
	start := "Mxmy"
	end := "Nxny"

	o.Channel(spaceID)
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
		fmt.Sprintf("/v2/objects/%s/channels/%s/uuids", pn.Config.SubscribeKey, "id0"),
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

func TestGetChannelMembersV2(t *testing.T) {
	AssertGetChannelMembersV2(t, true, false, false, false)
}

func TestGetChannelMembersV2Context(t *testing.T) {
	AssertGetChannelMembersV2(t, true, true, false, false)
}

func TestGetChannelMembersV2WithFilter(t *testing.T) {
	AssertGetChannelMembersV2(t, true, false, true, false)
}

func TestGetChannelMembersV2WithFilterContext(t *testing.T) {
	AssertGetChannelMembersV2(t, true, true, true, false)
}

func TestGetChannelMembersV2WithSort(t *testing.T) {
	AssertGetChannelMembersV2(t, true, false, false, true)
}

func TestGetChannelMembersV2WithSortContext(t *testing.T) {
	AssertGetChannelMembersV2(t, true, true, false, true)
}

func TestGetChannelMembersV2WithFilterWithSort(t *testing.T) {
	AssertGetChannelMembersV2(t, true, false, true, true)
}

func TestGetChannelMembersV2WithFilterWithSortContext(t *testing.T) {
	AssertGetChannelMembersV2(t, true, true, true, true)
}

func TestGetChannelMembersV2ResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetChannelMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestGetChannelMembersV2ResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":[{"id":"id0","custom":{"a3":"b3","c3":"d3"},"uuid":{"id":"id0","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-20T13:26:19.140324Z","updated":"2019-08-20T13:26:19.140324Z","eTag":"AbyT4v2p6K7fpQE"},"created":"2019-08-20T13:26:24.07832Z","updated":"2019-08-20T13:26:24.07832Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":1,"next":"MQ","prev":"NQ"}`)

	r, _, err := newPNGetChannelMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(1, r.TotalCount)
	assert.Equal("MQ", r.Next)
	assert.Equal("NQ", r.Prev)
	assert.Equal("id0", r.Data[0].ID)
	assert.Equal("name", r.Data[0].UUID.Name)
	assert.Equal("extid", r.Data[0].UUID.ExternalID)
	assert.Equal("purl", r.Data[0].UUID.ProfileURL)
	assert.Equal("email", r.Data[0].UUID.Email)
	// assert.Equal("2019-08-20T13:26:19.140324Z", r.Data[0].UUID.Created)
	assert.Equal("2019-08-20T13:26:19.140324Z", r.Data[0].UUID.Updated)
	assert.Equal("AbyT4v2p6K7fpQE", r.Data[0].UUID.ETag)
	assert.Equal("b", r.Data[0].UUID.Custom["a"])
	assert.Equal("d", r.Data[0].UUID.Custom["c"])
	assert.Equal("2019-08-20T13:26:24.07832Z", r.Data[0].Created)
	assert.Equal("2019-08-20T13:26:24.07832Z", r.Data[0].Updated)
	assert.Equal("AamrnoXdpdmzjwE", r.Data[0].ETag)
	assert.Equal("b3", r.Data[0].Custom["a3"])
	assert.Equal("d3", r.Data[0].Custom["c3"])

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestGetChannelMembersV2ValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)
	opts.Channel = "test-channel"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetChannelMembersV2ValidateMissingChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)
	opts.Channel = ""

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestGetChannelMembersV2ValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)
	opts.Channel = "test-channel"

	assert.Nil(opts.validate())
}

func TestGetChannelMembersV2ValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.Include = []string{"custom", "uuid"}
	opts.QueryParam = map[string]string{"param": "value"}

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestGetChannelMembersV2HTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)

	assert.Equal("GET", opts.httpMethod())
}

func TestGetChannelMembersV2OperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)

	assert.Equal(PNGetChannelMembersOperation, opts.operationType())
}

func TestGetChannelMembersV2IsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestGetChannelMembersV2Timeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (9 setters)

func TestGetChannelMembersV2BuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetChannelMembersBuilderV2(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
	assert.Equal(membersLimitV2, builder.opts.Limit) // Default limit
}

func TestGetChannelMembersV2BuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetChannelMembersBuilderV2WithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestGetChannelMembersV2BuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetChannelMembersBuilderV2(pn)

	// Test Channel setter
	builder.Channel("test-channel")
	assert.Equal("test-channel", builder.opts.Channel)

	// Test Include setter
	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
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

	// Test Filter setter
	builder.Filter("name LIKE 'user*'")
	assert.Equal("name LIKE 'user*'", builder.opts.Filter)

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

func TestGetChannelMembersV2BuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	sort := []string{"name"}
	queryParam := map[string]string{"key": "value"}

	builder := newGetChannelMembersBuilderV2(pn)
	result := builder.Channel("test-channel").
		Include(include).
		Limit(75).
		Start("start").
		End("end").
		Filter("filter").
		Sort(sort).
		Count(true).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal(75, builder.opts.Limit)
	assert.Equal("start", builder.opts.Start)
	assert.Equal("end", builder.opts.End)
	assert.Equal("filter", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(true, builder.opts.Count)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestGetChannelMembersV2BuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newGetChannelMembersBuilderV2(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

func TestGetChannelMembersV2BuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetChannelMembersBuilderV2(pn)

	// Verify default values
	assert.Equal("", builder.opts.Channel)
	assert.Nil(builder.opts.Include)
	assert.Equal(membersLimitV2, builder.opts.Limit)
	assert.Equal("", builder.opts.Start)
	assert.Equal("", builder.opts.End)
	assert.Equal("", builder.opts.Filter)
	assert.Nil(builder.opts.Sort)
	assert.Equal(false, builder.opts.Count)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestGetChannelMembersV2BuilderIncludeTypes(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		includes []PNChannelMembersInclude
		expected []string
	}{
		{
			name:     "Single include",
			includes: []PNChannelMembersInclude{PNChannelMembersIncludeCustom},
			expected: []string{"custom"},
		},
		{
			name:     "Multiple includes",
			includes: []PNChannelMembersInclude{PNChannelMembersIncludeCustom, PNChannelMembersIncludeUUID},
			expected: []string{"custom", "uuid"},
		},
		{
			name:     "All includes",
			includes: []PNChannelMembersInclude{PNChannelMembersIncludeCustom, PNChannelMembersIncludeUUID, PNChannelMembersIncludeUUIDCustom},
			expected: []string{"custom", "uuid", "uuid.custom"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetChannelMembersBuilderV2(pn)
			builder.Include(tc.includes)

			expectedInclude := EnumArrayToStringArray(tc.includes)
			assert.Equal(expectedInclude, builder.opts.Include)
		})
	}
}

func TestGetChannelMembersV2BuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	transport := &http.Transport{}
	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	sort := []string{"name", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Test all 9 setters in chain
	builder := newGetChannelMembersBuilderV2(pn).
		Channel("test-channel").
		Include(include).
		Limit(50).
		Start("start-token").
		End("end-token").
		Filter("name LIKE 'test*'").
		Sort(sort).
		Count(true).
		QueryParam(queryParam).
		Transport(transport)

	// Verify all are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal(50, builder.opts.Limit)
	assert.Equal("start-token", builder.opts.Start)
	assert.Equal("end-token", builder.opts.End)
	assert.Equal("name LIKE 'test*'", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(true, builder.opts.Count)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestGetChannelMembersV2BuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)
	opts.Channel = "test-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/channels/test-channel/uuids"
	assert.Equal(expected, path)
}

func TestGetChannelMembersV2BuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)
	opts.Channel = "my-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/channels/my-channel/uuids"
	assert.Equal(expected, path)
}

func TestGetChannelMembersV2BuildPathWithSpecialCharsInChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)
	opts.Channel = "channel-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")
	assert.Contains(path, "channel-with-special@chars#and$symbols")
	assert.Contains(path, "/uuids")
}

func TestGetChannelMembersV2BuildPathWithUnicodeChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)
	opts.Channel = "频道名称-канал-チャンネル"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")
	assert.Contains(path, "频道名称-канал-チャンネル")
	assert.Contains(path, "/uuids")
}

// GET-Specific Tests (Empty Body)

func TestGetChannelMembersV2BuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)
	opts.Channel = "test-channel"

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations have no body
}

func TestGetChannelMembersV2BuildBodyEmptyWithAllParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)

	// Set all possible parameters
	opts.Channel = "test-channel"
	opts.Include = []string{"custom", "uuid"}
	opts.Limit = 50
	opts.Start = "start"
	opts.End = "end"
	opts.Filter = "active = true"
	opts.Sort = []string{"name:asc"}
	opts.Count = true
	opts.QueryParam = map[string]string{"extra": "param"}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations always have empty body regardless of parameters
}

func TestGetChannelMembersV2GetOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetChannelMembersBuilderV2(pn)
	builder.Channel("test-channel")

	// Verify it's a GET operation
	assert.Equal("GET", builder.opts.httpMethod())

	// GET operations have no body
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	// Should have proper path for member retrieval
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/test-channel/uuids")
}

func TestGetChannelMembersV2DefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetChannelMembersBuilderV2(pn)

	// Should have default limit set to membersLimitV2 (100)
	assert.Equal(membersLimitV2, builder.opts.Limit)
	assert.Equal(100, builder.opts.Limit)

	// Should be included in query
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("limit"))
}

// Query Parameter Tests

func TestGetChannelMembersV2BuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("0", query.Get("limit")) // Default limit not set until builder initialization
	assert.Equal("0", query.Get("count")) // Default count=false
}

func TestGetChannelMembersV2BuildQueryWithInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)

	opts.Include = []string{"custom", "uuid"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	includeValue := query.Get("include")
	assert.Contains(includeValue, "custom")
	assert.Contains(includeValue, "uuid")
}

func TestGetChannelMembersV2BuildQueryWithPagination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)

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

func TestGetChannelMembersV2BuildQueryWithFilterAndSort(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)

	opts.Filter = "custom.role == 'admin'"
	opts.Sort = []string{"name", "created:desc"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("custom.role == 'admin'", query.Get("filter"))

	sortValue := query.Get("sort")
	assert.Contains(sortValue, "name")
	assert.Contains(sortValue, "created:desc")
}

func TestGetChannelMembersV2BuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)

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

func TestGetChannelMembersV2BuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)

	// Set all possible query parameters
	opts.Include = []string{"custom", "uuid"}
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

// Query Parameter Edge Cases

func TestGetChannelMembersV2QueryParameterHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		queryParam map[string]string
	}{
		{
			name:       "No query parameters",
			queryParam: nil,
		},
		{
			name:       "Empty query parameters",
			queryParam: map[string]string{},
		},
		{
			name: "Single query parameter",
			queryParam: map[string]string{
				"single": "value",
			},
		},
		{
			name: "Multiple query parameters",
			queryParam: map[string]string{
				"param1": "value1",
				"param2": "value2",
				"param3": "value3",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetChannelMembersOptsV2(pn, pn.ctx)
			opts.QueryParam = tc.queryParam

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Should always have default parameters
			assert.NotEmpty(query.Get("uuid"))
			assert.NotEmpty(query.Get("pnsdk"))

			// Verify custom parameters if any
			if tc.queryParam != nil {
				for key, expectedValue := range tc.queryParam {
					assert.Equal(expectedValue, query.Get(key))
				}
			}
		})
	}
}

// Comprehensive Edge Case Tests

func TestGetChannelMembersV2WithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*getChannelMembersBuilderV2)
	}{
		{
			name: "Very long channel name",
			setupFn: func(builder *getChannelMembersBuilderV2) {
				longChannel := strings.Repeat("VeryLongChannel", 50) // 750 characters
				builder.Channel(longChannel)
			},
		},
		{
			name: "Large filter expression",
			setupFn: func(builder *getChannelMembersBuilderV2) {
				largeFilter := "(" + strings.Repeat("custom.field == 'value' OR ", 100) + "custom.final == 'end')"
				builder.Filter(largeFilter)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *getChannelMembersBuilderV2) {
				largeQueryParam := make(map[string]string)
				for i := 0; i < 100; i++ {
					largeQueryParam[fmt.Sprintf("param_%d", i)] = strings.Repeat(fmt.Sprintf("value_%d_", i), 20)
				}
				builder.QueryParam(largeQueryParam)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetChannelMembersBuilderV2(pn)
			builder.Channel("test-channel")
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

			// Should build empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
		})
	}
}

func TestGetChannelMembersV2SpecialCharacterHandling(t *testing.T) {
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
			builder := newGetChannelMembersBuilderV2(pn)
			builder.Channel(specialString)
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

			// Should build empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
		})
	}
}

func TestGetChannelMembersV2ParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		channel string
		limit   int
		filter  string
	}{
		{
			name:    "Empty string channel",
			channel: "",
			limit:   1,
			filter:  "",
		},
		{
			name:    "Single character channel",
			channel: "a",
			limit:   1,
			filter:  "a",
		},
		{
			name:    "Unicode-only channel",
			channel: "测试",
			limit:   50,
			filter:  "测试 == '值'",
		},
		{
			name:    "Minimum limit",
			channel: "test",
			limit:   1,
			filter:  "simple",
		},
		{
			name:    "Large limit",
			channel: "test",
			limit:   1000,
			filter:  "complex.nested == 'value'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetChannelMembersBuilderV2(pn)
			builder.Channel(tc.channel)
			builder.Limit(tc.limit)
			if tc.filter != "" {
				builder.Filter(tc.filter)
			}

			// Should pass validation or fail gracefully
			err := builder.opts.validate()
			if tc.channel == "" {
				assert.NotNil(err) // Empty channel should fail validation
			} else {
				assert.Nil(err)

				// Should build valid components
				path, err := builder.opts.buildPath()
				assert.Nil(err)
				if tc.channel != "" {
					assert.Contains(path, tc.channel)
				}

				query, err := builder.opts.buildQuery()
				assert.Nil(err)
				assert.Equal(fmt.Sprintf("%d", tc.limit), query.Get("limit"))

				body, err := builder.opts.buildBody()
				assert.Nil(err)
				assert.Empty(body) // GET operation
			}
		})
	}
}

func TestGetChannelMembersV2ComplexFilterExpressions(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*getChannelMembersBuilderV2)
		validateFn func(*testing.T, *url.Values)
	}{
		{
			name: "Simple equality filter",
			setupFn: func(builder *getChannelMembersBuilderV2) {
				builder.Filter("custom.role == 'admin'")
			},
			validateFn: func(t *testing.T, query *url.Values) {
				assert.Equal("custom.role == 'admin'", query.Get("filter"))
			},
		},
		{
			name: "Complex OR filter",
			setupFn: func(builder *getChannelMembersBuilderV2) {
				builder.Filter("custom.role == 'admin' OR custom.level > 5")
			},
			validateFn: func(t *testing.T, query *url.Values) {
				assert.Contains(query.Get("filter"), "custom.role == 'admin'")
				assert.Contains(query.Get("filter"), "OR")
				assert.Contains(query.Get("filter"), "custom.level > 5")
			},
		},
		{
			name: "Nested conditions with parentheses",
			setupFn: func(builder *getChannelMembersBuilderV2) {
				builder.Filter("(custom.role == 'admin' OR custom.role == 'moderator') AND custom.active == true")
			},
			validateFn: func(t *testing.T, query *url.Values) {
				filter := query.Get("filter")
				assert.Contains(filter, "(")
				assert.Contains(filter, ")")
				assert.Contains(filter, "AND")
			},
		},
		{
			name: "Unicode filter with international characters",
			setupFn: func(builder *getChannelMembersBuilderV2) {
				builder.Filter("custom.姓名 == '张三' OR custom.роль == 'администратор'")
			},
			validateFn: func(t *testing.T, query *url.Values) {
				filter := query.Get("filter")
				assert.NotEmpty(filter)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetChannelMembersBuilderV2(pn)
			builder.Channel("test-channel")
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid components
			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Run custom validation
			tc.validateFn(t, query)
		})
	}
}

// Error Scenario Tests

func TestGetChannelMembersV2ExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newGetChannelMembersBuilderV2(pn)
	builder.Channel("test-channel")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetChannelMembersV2PathBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCases := []struct {
		name         string
		subscribeKey string
		channel      string
		expectError  bool
	}{
		{
			name:         "Empty SubscribeKey",
			subscribeKey: "",
			channel:      "test-channel",
			expectError:  false, // buildPath doesn't validate, only validate() does
		},
		{
			name:         "Empty Channel",
			subscribeKey: "demo",
			channel:      "",
			expectError:  false, // buildPath doesn't validate channel
		},
		{
			name:         "SubscribeKey with spaces",
			subscribeKey: "   sub key   ",
			channel:      "test-channel",
			expectError:  false,
		},
		{
			name:         "Channel with spaces",
			subscribeKey: "demo",
			channel:      "   test channel   ",
			expectError:  false,
		},
		{
			name:         "Very special characters",
			subscribeKey: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			channel:      "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			expectError:  false,
		},
		{
			name:         "Unicode SubscribeKey and Channel",
			subscribeKey: "测试订阅键-русский-キー",
			channel:      "频道名称-канал-チャンネル",
			expectError:  false,
		},
		{
			name:         "Very long values",
			subscribeKey: strings.Repeat("a", 1000),
			channel:      strings.Repeat("b", 1000),
			expectError:  false,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			pn.Config.SubscribeKey = tc.subscribeKey
			opts := newGetChannelMembersOptsV2(pn, pn.ctx)
			opts.Channel = tc.channel

			path, err := opts.buildPath()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Contains(path, "/v2/objects/")
				assert.Contains(path, "/channels/")
				assert.Contains(path, "/uuids")
			}
		})
	}
}

func TestGetChannelMembersV2QueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*getChannelMembersOptsV2)
		expectError bool
	}{
		{
			name: "Nil query params",
			setupOpts: func(opts *getChannelMembersOptsV2) {
				opts.QueryParam = nil
			},
			expectError: false,
		},
		{
			name: "Empty query params",
			setupOpts: func(opts *getChannelMembersOptsV2) {
				opts.QueryParam = map[string]string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *getChannelMembersOptsV2) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *getChannelMembersOptsV2) {
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
			opts := newGetChannelMembersOptsV2(pn, pn.ctx)
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

func TestGetChannelMembersV2BuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newGetChannelMembersBuilderV2(pn)

	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	sort := []string{"name:asc", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.Channel("complete-test-channel").
		Include(include).
		Limit(75).
		Start("start-token").
		End("end-token").
		Filter("active = true").
		Sort(sort).
		Count(true).
		QueryParam(queryParam)

	// Verify all values are set
	assert.Equal("complete-test-channel", builder.opts.Channel)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal(75, builder.opts.Limit)
	assert.Equal("start-token", builder.opts.Start)
	assert.Equal("end-token", builder.opts.End)
	assert.Equal("active = true", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(true, builder.opts.Count)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v2/objects/demo/channels/complete-test-channel/uuids"
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

	// Should build empty body (GET operation)
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}

func TestGetChannelMembersV2ResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetChannelMembersOptsV2(pn, pn.ctx)

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
			name:        "Valid response with member data",
			jsonBytes:   []byte(`{"status":200,"data":[{"id":"user1","uuid":{"id":"user1","name":"User 1"}}],"totalCount":1,"next":"abc","prev":"xyz"}`),
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
			resp, _, err := newPNGetChannelMembersResponse(tc.jsonBytes, opts, StatusResponse{})

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

// Include Parameter Effects on Response

func TestGetChannelMembersV2IncludeParameterEffects(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		includes []PNChannelMembersInclude
		expected []string
	}{
		{
			name:     "Include custom only",
			includes: []PNChannelMembersInclude{PNChannelMembersIncludeCustom},
			expected: []string{"custom"},
		},
		{
			name:     "Include UUID only",
			includes: []PNChannelMembersInclude{PNChannelMembersIncludeUUID},
			expected: []string{"uuid"},
		},
		{
			name:     "Include UUID custom only",
			includes: []PNChannelMembersInclude{PNChannelMembersIncludeUUIDCustom},
			expected: []string{"uuid.custom"},
		},
		{
			name:     "Include all",
			includes: []PNChannelMembersInclude{PNChannelMembersIncludeCustom, PNChannelMembersIncludeUUID, PNChannelMembersIncludeUUIDCustom},
			expected: []string{"custom", "uuid", "uuid.custom"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetChannelMembersBuilderV2(pn)
			builder.Channel("test-channel")
			builder.Include(tc.includes)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)

			includeValue := query.Get("include")
			for _, expected := range tc.expected {
				assert.Contains(includeValue, expected)
			}
		})
	}
}
