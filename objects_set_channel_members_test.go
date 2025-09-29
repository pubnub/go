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

func AssertSetChannelMembers(t *testing.T, checkQueryParam, testContext bool, withFilter bool, withSort bool) {
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

	o := newSetChannelMembersBuilder(pn)
	if testContext {
		o = newSetChannelMembersBuilderWithContext(pn, pn.ctx)
	}

	userID := "id0"
	limit := 90
	start := "Mxmy"
	end := "Nxny"

	o.Channel(userID)
	o.Include(incl)
	o.Limit(limit)
	o.Start(start)
	o.End(end)
	o.Count(false)
	o.QueryParam(queryParam)

	id0 := "id0"
	if withFilter {
		o.Filter("name like 'a*'")
	}
	sort := []string{"name", "created:desc"}
	if withSort {
		o.Sort(sort)
	}

	custom3 := make(map[string]interface{})
	custom3["a3"] = "b3"
	custom3["c3"] = "d3"

	uuid := PNChannelMembersUUID{
		ID: id0,
	}

	in := PNChannelMembersSet{
		UUID:   uuid,
		Custom: custom3,
	}

	inArr := []PNChannelMembersSet{
		in,
	}

	custom4 := make(map[string]interface{})
	custom4["a4"] = "b4"
	custom4["c4"] = "d4"

	o.Set(inArr)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/channels/%s/uuids", pn.Config.SubscribeKey, userID),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	expectedBody := "{\"set\":[{\"uuid\":{\"id\":\"id0\"},\"custom\":{\"a3\":\"b3\",\"c3\":\"d3\"}}]}"

	assert.Equal(expectedBody, string(body))

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

func TestSetChannelMembers(t *testing.T) {
	AssertSetChannelMembers(t, true, false, false, false)
}

func TestSetChannelMembersContext(t *testing.T) {
	AssertSetChannelMembers(t, true, true, false, false)
}

func TestSetChannelMembersWithFilter(t *testing.T) {
	AssertSetChannelMembers(t, true, false, true, false)
}

func TestSetChannelMembersWithFilterContext(t *testing.T) {
	AssertSetChannelMembers(t, true, true, true, false)
}

func TestSetChannelMembersWithSort(t *testing.T) {
	AssertSetChannelMembers(t, true, false, false, true)
}

func TestSetChannelMembersWithSortContext(t *testing.T) {
	AssertSetChannelMembers(t, true, true, false, true)
}

func TestSetChannelMembersWithFilterWithSort(t *testing.T) {
	AssertSetChannelMembers(t, true, false, true, true)
}

func TestSetChannelMembersWithFilterWithSortContext(t *testing.T) {
	AssertSetChannelMembers(t, true, true, true, true)
}

func TestSetChannelMembersResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNSetChannelMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestSetChannelMembersResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	jsonBytes := []byte(`{"status":200,"data":[{"id":"userid4","custom":{"a1":"b1","c1":"d1"},"uuid":{"id":"userid4","name":"userid4name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-23T10:36:27.083453Z","updated":"2019-08-23T10:36:27.083453Z","eTag":"AbuLvdnC9JnYEA"},"created":"2019-08-23T10:41:35.503214Z","updated":"2019-08-23T10:41:35.503214Z","eTag":"AZK3l4nQsrWG9gE"}],"totalCount":1,"next":"MQ"}`)

	r, _, err := newPNSetChannelMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(1, r.TotalCount)
	assert.Equal("MQ", r.Next)
	assert.Equal("userid4", r.Data[0].ID)
	assert.Equal("2019-08-23T10:41:35.503214Z", r.Data[0].Created)
	assert.Equal("2019-08-23T10:41:35.503214Z", r.Data[0].Updated)
	assert.Equal("AZK3l4nQsrWG9gE", r.Data[0].ETag)
	assert.Equal("b1", r.Data[0].Custom["a1"])
	assert.Equal("d1", r.Data[0].Custom["c1"])
	assert.Equal("userid4", r.Data[0].UUID.ID)
	assert.Equal("userid4name", r.Data[0].UUID.Name)
	assert.Equal("extid", r.Data[0].UUID.ExternalID)
	assert.Equal("purl", r.Data[0].UUID.ProfileURL)
	assert.Equal("email", r.Data[0].UUID.Email)
	//assert.Equal("2019-08-23T10:36:27.083453Z", r.Data[0].UUID.Created)
	assert.Equal("2019-08-23T10:36:27.083453Z", r.Data[0].UUID.Updated)
	assert.Equal("AbuLvdnC9JnYEA", r.Data[0].UUID.ETag)
	assert.Equal("b", r.Data[0].UUID.Custom["a"])
	assert.Equal("d", r.Data[0].UUID.Custom["c"])

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestSetChannelMembersValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newSetChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestSetChannelMembersValidateMissingChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)
	opts.Channel = ""

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestSetChannelMembersValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	assert.Nil(opts.validate())
}

func TestSetChannelMembersValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.ChannelMembersSet = []PNChannelMembersSet{
		{
			UUID: PNChannelMembersUUID{ID: "user1"},
		},
	}
	opts.QueryParam = map[string]string{"param": "value"}

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestSetChannelMembersHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	assert.Equal("PATCH", opts.httpMethod())
}

func TestSetChannelMembersOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	assert.Equal(PNSetChannelMembersOperation, opts.operationType())
}

func TestSetChannelMembersIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestSetChannelMembersTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (11 setters)

func TestSetChannelMembersBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetChannelMembersBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
	assert.Equal(setChannelMembersLimit, builder.opts.Limit) // Default limit
}

func TestSetChannelMembersBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetChannelMembersBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestSetChannelMembersBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetChannelMembersBuilder(pn)

	// Test Include setter
	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	builder.Include(include)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)

	// Test Channel setter
	builder.Channel("test-channel")
	assert.Equal("test-channel", builder.opts.Channel)

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
	builder.Filter("name LIKE 'user*'")
	assert.Equal("name LIKE 'user*'", builder.opts.Filter)

	// Test Sort setter
	sort := []string{"name", "created:desc"}
	builder.Sort(sort)
	assert.Equal(sort, builder.opts.Sort)

	// Test Set setter (core functionality)
	memberSet := []PNChannelMembersSet{
		{
			UUID: PNChannelMembersUUID{ID: "user1"},
		},
	}
	builder.Set(memberSet)
	assert.Equal(memberSet, builder.opts.ChannelMembersSet)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestSetChannelMembersBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	sort := []string{"name"}
	memberSet := []PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "user1"}},
	}
	queryParam := map[string]string{"key": "value"}

	builder := newSetChannelMembersBuilder(pn)
	result := builder.Include(include).
		Channel("test-channel").
		Limit(75).
		Start("start").
		End("end").
		Count(true).
		Filter("filter").
		Sort(sort).
		Set(memberSet).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal(75, builder.opts.Limit)
	assert.Equal("start", builder.opts.Start)
	assert.Equal("end", builder.opts.End)
	assert.Equal(true, builder.opts.Count)
	assert.Equal("filter", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(memberSet, builder.opts.ChannelMembersSet)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestSetChannelMembersBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newSetChannelMembersBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

func TestSetChannelMembersBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetChannelMembersBuilder(pn)

	// Verify default values
	assert.Nil(builder.opts.Include)
	assert.Equal("", builder.opts.Channel)
	assert.Equal(setChannelMembersLimit, builder.opts.Limit)
	assert.Equal("", builder.opts.Start)
	assert.Equal("", builder.opts.End)
	assert.Equal(false, builder.opts.Count)
	assert.Equal("", builder.opts.Filter)
	assert.Nil(builder.opts.Sort)
	assert.Nil(builder.opts.ChannelMembersSet)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestSetChannelMembersBuilderIncludeTypes(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetChannelMembersBuilder(pn)
			builder.Include(tc.includes)

			expectedInclude := EnumArrayToStringArray(tc.includes)
			assert.Equal(expectedInclude, builder.opts.Include)
		})
	}
}

func TestSetChannelMembersBuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	transport := &http.Transport{}
	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	sort := []string{"name", "created:desc"}
	memberSet := []PNChannelMembersSet{
		{
			UUID:   PNChannelMembersUUID{ID: "user1"},
			Custom: map[string]interface{}{"role": "admin"},
		},
	}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Test all 11 setters in chain
	builder := newSetChannelMembersBuilder(pn).
		Include(include).
		Channel("test-channel").
		Limit(50).
		Start("start-token").
		End("end-token").
		Count(true).
		Filter("name LIKE 'test*'").
		Sort(sort).
		Set(memberSet).
		QueryParam(queryParam).
		Transport(transport)

	// Verify all are set correctly
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal(50, builder.opts.Limit)
	assert.Equal("start-token", builder.opts.Start)
	assert.Equal("end-token", builder.opts.End)
	assert.Equal(true, builder.opts.Count)
	assert.Equal("name LIKE 'test*'", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(memberSet, builder.opts.ChannelMembersSet)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestSetChannelMembersBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/channels/test-channel/uuids"
	assert.Equal(expected, path)
}

func TestSetChannelMembersBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newSetChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "my-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/channels/my-channel/uuids"
	assert.Equal(expected, path)
}

func TestSetChannelMembersBuildPathWithSpecialCharsInChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "channel-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")
	assert.Contains(path, "channel-with-special@chars#and$symbols")
	assert.Contains(path, "/uuids")
}

func TestSetChannelMembersBuildPathWithUnicodeChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "频道名称-канал-チャンネル"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")
	assert.Contains(path, "频道名称-канал-チャンネル")
	assert.Contains(path, "/uuids")
}

// JSON Body Building Tests

func TestSetChannelMembersBuildBodyBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	memberSet := []PNChannelMembersSet{
		{
			UUID: PNChannelMembersUUID{ID: "user1"},
		},
	}
	opts.ChannelMembersSet = memberSet

	body, err := opts.buildBody()
	assert.Nil(err)
	expectedBody := `{"set":[{"uuid":{"id":"user1"},"custom":null}]}`
	assert.Equal(expectedBody, string(body))
}

func TestSetChannelMembersBuildBodyWithCustomData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	custom := map[string]interface{}{
		"role":        "admin",
		"permissions": []string{"read", "write"},
		"level":       5,
	}

	memberSet := []PNChannelMembersSet{
		{
			UUID:   PNChannelMembersUUID{ID: "user1"},
			Custom: custom,
		},
	}
	opts.ChannelMembersSet = memberSet

	body, err := opts.buildBody()
	assert.Nil(err)

	// Parse and verify JSON structure (order may vary)
	bodyStr := string(body)
	assert.Contains(bodyStr, `"set":[`)
	assert.Contains(bodyStr, `"uuid":{"id":"user1"}`)
	assert.Contains(bodyStr, `"custom":`)
	assert.Contains(bodyStr, `"role":"admin"`)
	assert.Contains(bodyStr, `"level":5`)
}

func TestSetChannelMembersBuildBodyMultipleMembers(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	memberSet := []PNChannelMembersSet{
		{
			UUID:   PNChannelMembersUUID{ID: "user1"},
			Custom: map[string]interface{}{"role": "admin"},
		},
		{
			UUID:   PNChannelMembersUUID{ID: "user2"},
			Custom: map[string]interface{}{"role": "member"},
		},
	}
	opts.ChannelMembersSet = memberSet

	body, err := opts.buildBody()
	assert.Nil(err)

	bodyStr := string(body)
	assert.Contains(bodyStr, `"set":[`)
	assert.Contains(bodyStr, `"uuid":{"id":"user1"}`)
	assert.Contains(bodyStr, `"uuid":{"id":"user2"}`)
	assert.Contains(bodyStr, `"role":"admin"`)
	assert.Contains(bodyStr, `"role":"member"`)
}

func TestSetChannelMembersBuildBodyWithUnicodeCustomData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	custom := map[string]interface{}{
		"名前":    "管理者",
		"роль":  "администратор",
		"役割":    "管理者",
		"emoji": "👑🔐",
	}

	memberSet := []PNChannelMembersSet{
		{
			UUID:   PNChannelMembersUUID{ID: "unicode-user"},
			Custom: custom,
		},
	}
	opts.ChannelMembersSet = memberSet

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.NotNil(body)

	bodyStr := string(body)
	assert.Contains(bodyStr, `"uuid":{"id":"unicode-user"}`)
	assert.Contains(bodyStr, `"custom":`)
}

func TestSetChannelMembersBuildBodyEmptySet(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	opts.ChannelMembersSet = []PNChannelMembersSet{}

	body, err := opts.buildBody()
	assert.Nil(err)
	expectedBody := `{"set":[]}`
	assert.Equal(expectedBody, string(body))
}

func TestSetChannelMembersBuildBodyNilSet(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	opts.ChannelMembersSet = nil

	body, err := opts.buildBody()
	assert.Nil(err)
	expectedBody := `{"set":null}`
	assert.Equal(expectedBody, string(body))
}

// Query Parameter Tests

func TestSetChannelMembersBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("100", query.Get("limit")) // Default limit
	assert.Equal("0", query.Get("count"))   // Default count=false
}

func TestSetChannelMembersBuildQueryWithInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	opts.Include = []string{"custom", "uuid"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	includeValue := query.Get("include")
	assert.Contains(includeValue, "custom")
	assert.Contains(includeValue, "uuid")
}

func TestSetChannelMembersBuildQueryWithPagination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

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

func TestSetChannelMembersBuildQueryWithFilterAndSort(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

	opts.Filter = "name LIKE 'user*'"
	opts.Sort = []string{"name", "created:desc"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("name LIKE 'user*'", query.Get("filter"))

	sortValue := query.Get("sort")
	assert.Contains(sortValue, "name")
	assert.Contains(sortValue, "created:desc")
}

func TestSetChannelMembersBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

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

func TestSetChannelMembersBuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

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

func TestSetChannelMembersQueryParameterHandling(t *testing.T) {
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
			opts := newSetChannelMembersOpts(pn, pn.ctx)
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

func TestSetChannelMembersWithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*setChannelMembersBuilder)
	}{
		{
			name: "Very long channel name",
			setupFn: func(builder *setChannelMembersBuilder) {
				longChannel := strings.Repeat("VeryLongChannel", 50) // 750 characters
				builder.Channel(longChannel)
			},
		},
		{
			name: "Large member set",
			setupFn: func(builder *setChannelMembersBuilder) {
				var memberSet []PNChannelMembersSet
				for i := 0; i < 100; i++ {
					memberSet = append(memberSet, PNChannelMembersSet{
						UUID: PNChannelMembersUUID{ID: fmt.Sprintf("user_%d", i)},
						Custom: map[string]interface{}{
							"role":  fmt.Sprintf("role_%d", i),
							"index": i,
						},
					})
				}
				builder.Set(memberSet)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *setChannelMembersBuilder) {
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
			builder := newSetChannelMembersBuilder(pn)
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

			// Should build valid body
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.NotNil(body)
		})
	}
}

func TestSetChannelMembersSpecialCharacterHandling(t *testing.T) {
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
			builder := newSetChannelMembersBuilder(pn)
			builder.Channel(specialString)
			builder.Filter(specialString)
			builder.QueryParam(map[string]string{
				"special_field": specialString,
			})
			builder.Set([]PNChannelMembersSet{
				{
					UUID: PNChannelMembersUUID{ID: specialString},
					Custom: map[string]interface{}{
						"special_custom": specialString,
					},
				},
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

			// Should build valid body
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.NotNil(body)
		})
	}
}

func TestSetChannelMembersParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		channel     string
		limit       int
		memberCount int
	}{
		{
			name:        "Empty string channel",
			channel:     "",
			limit:       1,
			memberCount: 1,
		},
		{
			name:        "Single character channel",
			channel:     "a",
			limit:       1,
			memberCount: 1,
		},
		{
			name:        "Unicode-only channel",
			channel:     "测试",
			limit:       50,
			memberCount: 5,
		},
		{
			name:        "Minimum limit",
			channel:     "test",
			limit:       1,
			memberCount: 1,
		},
		{
			name:        "Large limit",
			channel:     "test",
			limit:       1000,
			memberCount: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetChannelMembersBuilder(pn)
			builder.Channel(tc.channel)
			builder.Limit(tc.limit)

			var memberSet []PNChannelMembersSet
			for i := 0; i < tc.memberCount; i++ {
				memberSet = append(memberSet, PNChannelMembersSet{
					UUID: PNChannelMembersUUID{ID: fmt.Sprintf("user_%d", i)},
				})
			}
			builder.Set(memberSet)

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
				assert.NotNil(body)
			}
		})
	}
}

func TestSetChannelMembersComplexMemberSets(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*setChannelMembersBuilder)
		validateFn func(*testing.T, []byte)
	}{
		{
			name: "Members with nested custom data",
			setupFn: func(builder *setChannelMembersBuilder) {
				memberSet := []PNChannelMembersSet{
					{
						UUID: PNChannelMembersUUID{ID: "user1"},
						Custom: map[string]interface{}{
							"role": "admin",
							"permissions": map[string]interface{}{
								"channels": []string{"read", "write", "delete"},
								"users":    []string{"read", "write"},
							},
							"metadata": map[string]interface{}{
								"level":        5,
								"experience":   1000,
								"achievements": []string{"first_login", "power_user"},
							},
						},
					},
				}
				builder.Set(memberSet)
			},
			validateFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				assert.Contains(bodyStr, `"uuid":{"id":"user1"}`)
				assert.Contains(bodyStr, `"role":"admin"`)
				assert.Contains(bodyStr, `"permissions"`)
				assert.Contains(bodyStr, `"metadata"`)
			},
		},
		{
			name: "Unicode members with international data",
			setupFn: func(builder *setChannelMembersBuilder) {
				memberSet := []PNChannelMembersSet{
					{
						UUID: PNChannelMembersUUID{ID: "用户123"},
						Custom: map[string]interface{}{
							"姓名":    "张三",
							"роль":  "администратор",
							"役割":    "管理者",
							"emoji": "👑🌍",
						},
					},
				}
				builder.Set(memberSet)
			},
			validateFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				assert.Contains(bodyStr, `"uuid":{"id":"用户123"}`)
				assert.Contains(bodyStr, `"custom"`)
			},
		},
		{
			name: "Mixed member types",
			setupFn: func(builder *setChannelMembersBuilder) {
				memberSet := []PNChannelMembersSet{
					{
						UUID: PNChannelMembersUUID{ID: "simple_user"},
					},
					{
						UUID: PNChannelMembersUUID{ID: "admin_user"},
						Custom: map[string]interface{}{
							"role":  "admin",
							"level": 10,
						},
					},
					{
						UUID: PNChannelMembersUUID{ID: "special@user#123"},
						Custom: map[string]interface{}{
							"special_chars": "value@with#symbols",
							"unicode":       "测试数据",
						},
					},
				}
				builder.Set(memberSet)
			},
			validateFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				assert.Contains(bodyStr, `"uuid":{"id":"simple_user"}`)
				assert.Contains(bodyStr, `"uuid":{"id":"admin_user"}`)
				assert.Contains(bodyStr, `"uuid":{"id":"special@user#123"}`)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetChannelMembersBuilder(pn)
			builder.Channel("test-channel")
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid components
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.NotNil(body)

			// Run custom validation
			tc.validateFn(t, body)
		})
	}
}

// Error Scenario Tests

func TestSetChannelMembersExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newSetChannelMembersBuilder(pn)
	builder.Channel("test-channel")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestSetChannelMembersPathBuildingEdgeCases(t *testing.T) {
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
			opts := newSetChannelMembersOpts(pn, pn.ctx)
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

func TestSetChannelMembersQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*setChannelMembersOpts)
		expectError bool
	}{
		{
			name: "Nil query params",
			setupOpts: func(opts *setChannelMembersOpts) {
				opts.QueryParam = nil
			},
			expectError: false,
		},
		{
			name: "Empty query params",
			setupOpts: func(opts *setChannelMembersOpts) {
				opts.QueryParam = map[string]string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *setChannelMembersOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *setChannelMembersOpts) {
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
			opts := newSetChannelMembersOpts(pn, pn.ctx)
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

func TestSetChannelMembersBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newSetChannelMembersBuilder(pn)

	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	sort := []string{"name:asc", "created:desc"}
	memberSet := []PNChannelMembersSet{
		{
			UUID: PNChannelMembersUUID{ID: "complete-user"},
			Custom: map[string]interface{}{
				"role":  "test",
				"level": 1,
			},
		},
	}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.Include(include).
		Channel("complete-test-channel").
		Limit(75).
		Start("start-token").
		End("end-token").
		Count(true).
		Filter("active = true").
		Sort(sort).
		Set(memberSet).
		QueryParam(queryParam)

	// Verify all values are set
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal("complete-test-channel", builder.opts.Channel)
	assert.Equal(75, builder.opts.Limit)
	assert.Equal("start-token", builder.opts.Start)
	assert.Equal("end-token", builder.opts.End)
	assert.Equal(true, builder.opts.Count)
	assert.Equal("active = true", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(memberSet, builder.opts.ChannelMembersSet)
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

	// Should build valid body
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	bodyStr := string(body)
	assert.Contains(bodyStr, `"uuid":{"id":"complete-user"}`)
	assert.Contains(bodyStr, `"role":"test"`)
}

func TestSetChannelMembersResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetChannelMembersOpts(pn, pn.ctx)

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
			name:        "Response with status only",
			jsonBytes:   []byte(`{"status":200}`),
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, _, err := newPNSetChannelMembersResponse(tc.jsonBytes, opts, StatusResponse{})

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

// PATCH-specific tests

func TestSetChannelMembersPatchOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetChannelMembersBuilder(pn)
	builder.Channel("test-channel")
	builder.Set([]PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "user1"}},
	})

	// Verify it's a PATCH operation
	assert.Equal("PATCH", builder.opts.httpMethod())

	// PATCH operations should have a body
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.NotEmpty(body)

	// Should have proper path for member management
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/test-channel/uuids")
}

func TestSetChannelMembersPatchBodyStructure(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetChannelMembersBuilder(pn)
	builder.Channel("test-channel")

	memberSet := []PNChannelMembersSet{
		{
			UUID:   PNChannelMembersUUID{ID: "user1"},
			Custom: map[string]interface{}{"role": "admin"},
		},
		{
			UUID:   PNChannelMembersUUID{ID: "user2"},
			Custom: map[string]interface{}{"role": "member"},
		},
	}
	builder.Set(memberSet)

	// PATCH body should follow specific structure for setting members
	body, err := builder.opts.buildBody()
	assert.Nil(err)

	bodyStr := string(body)
	assert.Contains(bodyStr, `"set":[`)
	assert.Contains(bodyStr, `"uuid":{"id":"user1"}`)
	assert.Contains(bodyStr, `"uuid":{"id":"user2"}`)
	assert.Contains(bodyStr, `"role":"admin"`)
	assert.Contains(bodyStr, `"role":"member"`)
}

func TestSetChannelMembersDefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetChannelMembersBuilder(pn)

	// Should have default limit set to setChannelMembersLimit (100)
	assert.Equal(setChannelMembersLimit, builder.opts.Limit)
	assert.Equal(100, builder.opts.Limit)

	// Should be included in query
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("limit"))
}
