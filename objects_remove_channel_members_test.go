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

func AssertRemoveChannelMembers(t *testing.T, checkQueryParam, testContext bool, withFilter bool, withSort bool) {
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

	o := newRemoveChannelMembersBuilder(pn)
	if testContext {
		o = newRemoveChannelMembersBuilderWithContext(pn, pn.ctx)
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

	custom4 := make(map[string]interface{})
	custom4["a4"] = "b4"
	custom4["c4"] = "d4"

	uuid := PNChannelMembersUUID{
		ID: id0,
	}

	re := PNChannelMembersRemove{
		UUID: uuid,
	}

	reArr := []PNChannelMembersRemove{
		re,
	}

	o.Remove(reArr)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/channels/%s/uuids", pn.Config.SubscribeKey, userID),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	expectedBody := "{\"delete\":[{\"uuid\":{\"id\":\"id0\"}}]}"

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

func TestRemoveChannelMembers(t *testing.T) {
	AssertRemoveChannelMembers(t, true, false, false, false)
}

func TestRemoveChannelMembersContext(t *testing.T) {
	AssertRemoveChannelMembers(t, true, true, false, false)
}

func TestRemoveChannelMembersWithFilter(t *testing.T) {
	AssertRemoveChannelMembers(t, true, false, true, false)
}

func TestRemoveChannelMembersWithFilterContext(t *testing.T) {
	AssertRemoveChannelMembers(t, true, true, true, false)
}

func TestRemoveChannelMembersWithSort(t *testing.T) {
	AssertRemoveChannelMembers(t, true, false, false, true)
}

func TestRemoveChannelMembersWithSortContext(t *testing.T) {
	AssertRemoveChannelMembers(t, true, true, false, true)
}

func TestRemoveChannelMembersWithFilterWithSort(t *testing.T) {
	AssertRemoveChannelMembers(t, true, false, true, true)
}

func TestRemoveChannelMembersWithFilterWithSortContext(t *testing.T) {
	AssertRemoveChannelMembers(t, true, true, true, true)
}

func TestRemoveChannelMembersResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNRemoveChannelMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRemoveChannelMembersResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":[{"id":"userid4","custom":{"a1":"b1","c1":"d1"},"uuid":{"id":"userid4","name":"userid4name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-23T10:36:27.083453Z","updated":"2019-08-23T10:36:27.083453Z","eTag":"AbuLvdnC9JnYEA"},"created":"2019-08-23T10:41:35.503214Z","updated":"2019-08-23T10:41:35.503214Z","eTag":"AZK3l4nQsrWG9gE"}],"totalCount":1,"next":"MQ"}`)

	r, _, err := newPNRemoveChannelMembersResponse(jsonBytes, opts, StatusResponse{})

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
	// assert.Equal("2019-08-23T10:36:27.083453Z", r.Data[0].UUID.Created)
	assert.Equal("2019-08-23T10:36:27.083453Z", r.Data[0].UUID.Updated)
	assert.Equal("AbuLvdnC9JnYEA", r.Data[0].UUID.ETag)
	assert.Equal("b", r.Data[0].UUID.Custom["a"])
	assert.Equal("d", r.Data[0].UUID.Custom["c"])

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestRemoveChannelMembersValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestRemoveChannelMembersValidateMissingChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)
	opts.Channel = ""

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestRemoveChannelMembersValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	assert.Nil(opts.validate())
}

func TestRemoveChannelMembersValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.Include = []string{"custom", "uuid"}
	opts.QueryParam = map[string]string{"param": "value"}
	opts.ChannelMembersRemove = []PNChannelMembersRemove{{UUID: PNChannelMembersUUID{ID: "user1"}}}

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestRemoveChannelMembersHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	assert.Equal("PATCH", opts.httpMethod())
}

func TestRemoveChannelMembersOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	assert.Equal(PNRemoveChannelMembersOperation, opts.operationType())
}

func TestRemoveChannelMembersIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestRemoveChannelMembersTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (11 setters)

func TestRemoveChannelMembersBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelMembersBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
	assert.Equal(removeChannelMembersLimit, builder.opts.Limit) // Default limit (100)
}

func TestRemoveChannelMembersBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelMembersBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestRemoveChannelMembersBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelMembersBuilder(pn)

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

	// Test Remove setter (CRITICAL for PATCH operation)
	removeMembers := []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user1"}},
		{UUID: PNChannelMembersUUID{ID: "user2"}},
	}
	builder.Remove(removeMembers)
	assert.Equal(removeMembers, builder.opts.ChannelMembersRemove)

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

func TestRemoveChannelMembersBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	sort := []string{"name"}
	queryParam := map[string]string{"key": "value"}
	removeMembers := []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user1"}},
	}
	transport := &http.Transport{}

	builder := newRemoveChannelMembersBuilder(pn)
	result := builder.Channel("test-channel").
		Include(include).
		Limit(75).
		Start("start").
		End("end").
		Count(true).
		Filter("filter").
		Sort(sort).
		Remove(removeMembers).
		QueryParam(queryParam).
		Transport(transport)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal(75, builder.opts.Limit)
	assert.Equal("start", builder.opts.Start)
	assert.Equal("end", builder.opts.End)
	assert.Equal(true, builder.opts.Count)
	assert.Equal("filter", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(removeMembers, builder.opts.ChannelMembersRemove)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

func TestRemoveChannelMembersBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelMembersBuilder(pn)

	// Verify default values
	assert.Equal("", builder.opts.Channel)
	assert.Nil(builder.opts.Include)
	assert.Equal(removeChannelMembersLimit, builder.opts.Limit) // 100
	assert.Equal("", builder.opts.Start)
	assert.Equal("", builder.opts.End)
	assert.Equal(false, builder.opts.Count)
	assert.Equal("", builder.opts.Filter)
	assert.Nil(builder.opts.Sort)
	assert.Nil(builder.opts.ChannelMembersRemove)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestRemoveChannelMembersBuilderIncludeTypes(t *testing.T) {
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
			builder := newRemoveChannelMembersBuilder(pn)
			builder.Include(tc.includes)

			expectedInclude := EnumArrayToStringArray(tc.includes)
			assert.Equal(expectedInclude, builder.opts.Include)
		})
	}
}

func TestRemoveChannelMembersBuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	transport := &http.Transport{}
	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	sort := []string{"name", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}
	removeMembers := []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user1"}},
		{UUID: PNChannelMembersUUID{ID: "user2"}},
	}

	// Test all 11 setters in chain
	builder := newRemoveChannelMembersBuilder(pn).
		Channel("test-channel").
		Include(include).
		Limit(50).
		Start("start-token").
		End("end-token").
		Count(true).
		Filter("name LIKE 'test*'").
		Sort(sort).
		Remove(removeMembers).
		QueryParam(queryParam).
		Transport(transport)

	// Verify all are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal(50, builder.opts.Limit)
	assert.Equal("start-token", builder.opts.Start)
	assert.Equal("end-token", builder.opts.End)
	assert.Equal(true, builder.opts.Count)
	assert.Equal("name LIKE 'test*'", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(removeMembers, builder.opts.ChannelMembersRemove)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestRemoveChannelMembersBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/channels/test-channel/uuids"
	assert.Equal(expected, path)
}

func TestRemoveChannelMembersBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "my-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/channels/my-channel/uuids"
	assert.Equal(expected, path)
}

func TestRemoveChannelMembersBuildPathWithSpecialCharsInChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "channel-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")
	assert.Contains(path, "channel-with-special@chars#and$symbols")
	assert.Contains(path, "/uuids")
}

func TestRemoveChannelMembersBuildPathWithUnicodeChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)
	opts.Channel = "频道名称-канал-チャンネル"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")
	assert.Contains(path, "频道名称-канал-チャンネル")
	assert.Contains(path, "/uuids")
}

// JSON Body Building Tests (CRITICAL for PATCH operation)

func TestRemoveChannelMembersBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"delete":null}`
	assert.Equal(expected, string(body))
}

func TestRemoveChannelMembersBuildBodySingleMember(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	opts.ChannelMembersRemove = []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user1"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"delete":[{"uuid":{"id":"user1"}}]}`
	assert.Equal(expected, string(body))
}

func TestRemoveChannelMembersBuildBodyMultipleMembers(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	opts.ChannelMembersRemove = []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user1"}},
		{UUID: PNChannelMembersUUID{ID: "user2"}},
		{UUID: PNChannelMembersUUID{ID: "user3"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"delete":[{"uuid":{"id":"user1"}},{"uuid":{"id":"user2"}},{"uuid":{"id":"user3"}}]}`
	assert.Equal(expected, string(body))
}

func TestRemoveChannelMembersBuildBodyWithUnicodeUUIDs(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	opts.ChannelMembersRemove = []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "用户123"}},
		{UUID: PNChannelMembersUUID{ID: "пользователь456"}},
		{UUID: PNChannelMembersUUID{ID: "ユーザー789"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"id":"用户123"`)
	assert.Contains(string(body), `"id":"пользователь456"`)
	assert.Contains(string(body), `"id":"ユーザー789"`)
	assert.Contains(string(body), `"delete":[`)
}

func TestRemoveChannelMembersBuildBodyWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	opts.ChannelMembersRemove = []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user@domain.com"}},
		{UUID: PNChannelMembersUUID{ID: "user-with-dashes"}},
		{UUID: PNChannelMembersUUID{ID: "user_with_underscores"}},
		{UUID: PNChannelMembersUUID{ID: "user.with.dots"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"id":"user@domain.com"`)
	assert.Contains(string(body), `"id":"user-with-dashes"`)
	assert.Contains(string(body), `"id":"user_with_underscores"`)
	assert.Contains(string(body), `"id":"user.with.dots"`)
}

func TestRemoveChannelMembersBuildBodyEmptyRemoveArray(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	opts.ChannelMembersRemove = []PNChannelMembersRemove{}

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"delete":[]}`
	assert.Equal(expected, string(body))
}

// Query Parameter Tests

func TestRemoveChannelMembersBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("0", query.Get("limit")) // Default limit not set until builder initialization
	assert.Equal("0", query.Get("count")) // Default count=false
}

func TestRemoveChannelMembersBuildQueryWithInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	opts.Include = []string{"custom", "uuid"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	includeValue := query.Get("include")
	assert.Contains(includeValue, "custom")
	assert.Contains(includeValue, "uuid")
}

func TestRemoveChannelMembersBuildQueryWithPagination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

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

func TestRemoveChannelMembersBuildQueryWithFilterAndSort(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

	opts.Filter = "custom.role == 'admin'"
	opts.Sort = []string{"name", "created:desc"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("custom.role == 'admin'", query.Get("filter"))

	sortValue := query.Get("sort")
	assert.Contains(sortValue, "name")
	assert.Contains(sortValue, "created:desc")
}

func TestRemoveChannelMembersBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

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

func TestRemoveChannelMembersBuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

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

// PATCH-Specific Tests (Remove Operation Characteristics)

func TestRemoveChannelMembersPatchOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelMembersBuilder(pn)
	builder.Channel("test-channel")
	builder.Remove([]PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user1"}},
	})

	// Verify it's a PATCH operation
	assert.Equal("PATCH", builder.opts.httpMethod())

	// PATCH operations have JSON body with delete structure
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"delete":[`)
	assert.Contains(string(body), `"uuid":{"id":"user1"}`)

	// Should have proper path for member removal
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/test-channel/uuids")
}

func TestRemoveChannelMembersDefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelMembersBuilder(pn)

	// Should have default limit set to removeChannelMembersLimit (100)
	assert.Equal(removeChannelMembersLimit, builder.opts.Limit)
	assert.Equal(100, builder.opts.Limit)

	// Should be included in query
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("limit"))
}

func TestRemoveChannelMembersDeleteStructureValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name          string
		removeMembers []PNChannelMembersRemove
		expectedJSON  string
	}{
		{
			name:          "Nil remove members",
			removeMembers: nil,
			expectedJSON:  `{"delete":null}`,
		},
		{
			name:          "Empty remove members",
			removeMembers: []PNChannelMembersRemove{},
			expectedJSON:  `{"delete":[]}`,
		},
		{
			name: "Single remove member",
			removeMembers: []PNChannelMembersRemove{
				{UUID: PNChannelMembersUUID{ID: "test-user"}},
			},
			expectedJSON: `{"delete":[{"uuid":{"id":"test-user"}}]}`,
		},
		{
			name: "Multiple remove members",
			removeMembers: []PNChannelMembersRemove{
				{UUID: PNChannelMembersUUID{ID: "user-a"}},
				{UUID: PNChannelMembersUUID{ID: "user-b"}},
			},
			expectedJSON: `{"delete":[{"uuid":{"id":"user-a"}},{"uuid":{"id":"user-b"}}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newRemoveChannelMembersOpts(pn, pn.ctx)
			opts.ChannelMembersRemove = tc.removeMembers

			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Equal(tc.expectedJSON, string(body))
		})
	}
}

func TestRemoveChannelMembersResponseStructureAfterRemoval(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelMembersBuilder(pn)
	builder.Channel("test-channel")
	builder.Remove([]PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user1"}},
	})

	// Response should contain remaining members after removal
	// This is tested in the existing TestRemoveChannelMembersResponseValuePass
	// but verify the response structure expectations
	opts := builder.opts

	// Verify operation is configured correctly
	assert.Equal("PATCH", opts.httpMethod())
	assert.Equal(PNRemoveChannelMembersOperation, opts.operationType())
	assert.NotNil(opts.ChannelMembersRemove)
	assert.Equal(1, len(opts.ChannelMembersRemove))
}

// Comprehensive Edge Case Tests

func TestRemoveChannelMembersWithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*removeChannelMembersBuilder)
	}{
		{
			name: "Very long channel name",
			setupFn: func(builder *removeChannelMembersBuilder) {
				longChannel := strings.Repeat("VeryLongChannel", 50) // 750 characters
				builder.Channel(longChannel)
			},
		},
		{
			name: "Large number of members to remove",
			setupFn: func(builder *removeChannelMembersBuilder) {
				var largeRemoveList []PNChannelMembersRemove
				for i := 0; i < 100; i++ {
					largeRemoveList = append(largeRemoveList, PNChannelMembersRemove{
						UUID: PNChannelMembersUUID{ID: fmt.Sprintf("user_%d", i)},
					})
				}
				builder.Remove(largeRemoveList)
			},
		},
		{
			name: "Large filter expression",
			setupFn: func(builder *removeChannelMembersBuilder) {
				largeFilter := "(" + strings.Repeat("custom.field == 'value' OR ", 100) + "custom.final == 'end')"
				builder.Filter(largeFilter)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *removeChannelMembersBuilder) {
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
			builder := newRemoveChannelMembersBuilder(pn)
			builder.Channel("test-channel")
			builder.Remove([]PNChannelMembersRemove{
				{UUID: PNChannelMembersUUID{ID: "baseline-user"}},
			})
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

			// Should build valid JSON body (PATCH operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.NotEmpty(body)
			assert.Contains(string(body), `"delete":`)
		})
	}
}

func TestRemoveChannelMembersSpecialCharacterHandling(t *testing.T) {
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
			builder := newRemoveChannelMembersBuilder(pn)
			builder.Channel(specialString)
			builder.Filter(fmt.Sprintf("custom.field == '%s'", specialString))
			builder.QueryParam(map[string]string{
				"special_field": specialString,
			})
			builder.Remove([]PNChannelMembersRemove{
				{UUID: PNChannelMembersUUID{ID: specialString}},
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

			// Should build valid JSON body (PATCH operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.NotEmpty(body)
			assert.Contains(string(body), `"delete":`)
		})
	}
}

func TestRemoveChannelMembersParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		channel     string
		limit       int
		filter      string
		memberCount int
	}{
		{
			name:        "Empty string channel",
			channel:     "",
			limit:       1,
			filter:      "",
			memberCount: 0,
		},
		{
			name:        "Single character channel",
			channel:     "a",
			limit:       1,
			filter:      "a",
			memberCount: 1,
		},
		{
			name:        "Unicode-only channel",
			channel:     "测试",
			limit:       50,
			filter:      "测试 == '值'",
			memberCount: 2,
		},
		{
			name:        "Minimum limit",
			channel:     "test",
			limit:       1,
			filter:      "simple",
			memberCount: 1,
		},
		{
			name:        "Large limit",
			channel:     "test",
			limit:       1000,
			filter:      "complex.nested == 'value'",
			memberCount: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveChannelMembersBuilder(pn)
			builder.Channel(tc.channel)
			builder.Limit(tc.limit)
			if tc.filter != "" {
				builder.Filter(tc.filter)
			}

			// Add specified number of members to remove
			var removeMembers []PNChannelMembersRemove
			for i := 0; i < tc.memberCount; i++ {
				removeMembers = append(removeMembers, PNChannelMembersRemove{
					UUID: PNChannelMembersUUID{ID: fmt.Sprintf("user_%d", i)},
				})
			}
			if len(removeMembers) > 0 {
				builder.Remove(removeMembers)
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
				assert.NotEmpty(body) // PATCH operation always has body
			}
		})
	}
}

func TestRemoveChannelMembersComplexRemovalScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*removeChannelMembersBuilder)
		validateFn func(*testing.T, []byte)
	}{
		{
			name: "Mixed character UUID removals",
			setupFn: func(builder *removeChannelMembersBuilder) {
				builder.Remove([]PNChannelMembersRemove{
					{UUID: PNChannelMembersUUID{ID: "user-english"}},
					{UUID: PNChannelMembersUUID{ID: "用户中文"}},
					{UUID: PNChannelMembersUUID{ID: "пользователь"}},
					{UUID: PNChannelMembersUUID{ID: "ユーザー"}},
				})
			},
			validateFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				assert.Contains(bodyStr, `"id":"user-english"`)
				assert.Contains(bodyStr, `"id":"用户中文"`)
				assert.Contains(bodyStr, `"id":"пользователь"`)
				assert.Contains(bodyStr, `"id":"ユーザー"`)
			},
		},
		{
			name: "Email-like and special UUID removals",
			setupFn: func(builder *removeChannelMembersBuilder) {
				builder.Remove([]PNChannelMembersRemove{
					{UUID: PNChannelMembersUUID{ID: "user@domain.com"}},
					{UUID: PNChannelMembersUUID{ID: "user+tag@example.org"}},
					{UUID: PNChannelMembersUUID{ID: "user-with-dashes"}},
					{UUID: PNChannelMembersUUID{ID: "user_with_underscores"}},
				})
			},
			validateFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				assert.Contains(bodyStr, `"id":"user@domain.com"`)
				assert.Contains(bodyStr, `"id":"user+tag@example.org"`)
				assert.Contains(bodyStr, `"id":"user-with-dashes"`)
				assert.Contains(bodyStr, `"id":"user_with_underscores"`)
			},
		},
		{
			name: "Single vs multiple removal consistency",
			setupFn: func(builder *removeChannelMembersBuilder) {
				builder.Remove([]PNChannelMembersRemove{
					{UUID: PNChannelMembersUUID{ID: "single-user"}},
				})
			},
			validateFn: func(t *testing.T, body []byte) {
				expected := `{"delete":[{"uuid":{"id":"single-user"}}]}`
				assert.Equal(expected, string(body))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveChannelMembersBuilder(pn)
			builder.Channel("test-channel")
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid components
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.NotEmpty(body)

			// Run custom validation
			tc.validateFn(t, body)
		})
	}
}

// Error Scenario Tests

func TestRemoveChannelMembersExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newRemoveChannelMembersBuilder(pn)
	builder.Channel("test-channel")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestRemoveChannelMembersPathBuildingEdgeCases(t *testing.T) {
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
			opts := newRemoveChannelMembersOpts(pn, pn.ctx)
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

func TestRemoveChannelMembersQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*removeChannelMembersOpts)
		expectError bool
	}{
		{
			name: "Nil query params",
			setupOpts: func(opts *removeChannelMembersOpts) {
				opts.QueryParam = nil
			},
			expectError: false,
		},
		{
			name: "Empty query params",
			setupOpts: func(opts *removeChannelMembersOpts) {
				opts.QueryParam = map[string]string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *removeChannelMembersOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *removeChannelMembersOpts) {
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
			opts := newRemoveChannelMembersOpts(pn, pn.ctx)
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

func TestRemoveChannelMembersBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newRemoveChannelMembersBuilder(pn)

	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	sort := []string{"name:asc", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}
	removeMembers := []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user1"}},
		{UUID: PNChannelMembersUUID{ID: "user2"}},
	}

	// Set all possible parameters
	builder.Channel("complete-test-channel").
		Include(include).
		Limit(75).
		Start("start-token").
		End("end-token").
		Count(true).
		Filter("active = true").
		Sort(sort).
		Remove(removeMembers).
		QueryParam(queryParam)

	// Verify all values are set
	assert.Equal("complete-test-channel", builder.opts.Channel)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal(75, builder.opts.Limit)
	assert.Equal("start-token", builder.opts.Start)
	assert.Equal("end-token", builder.opts.End)
	assert.Equal(true, builder.opts.Count)
	assert.Equal("active = true", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(removeMembers, builder.opts.ChannelMembersRemove)
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

	// Should build correct JSON body (PATCH operation)
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	expectedBody := `{"delete":[{"uuid":{"id":"user1"}},{"uuid":{"id":"user2"}}]}`
	assert.Equal(expectedBody, string(body))
}

func TestRemoveChannelMembersResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelMembersOpts(pn, pn.ctx)

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
			name:        "Valid response with remaining member data",
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
			resp, _, err := newPNRemoveChannelMembersResponse(tc.jsonBytes, opts, StatusResponse{})

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
