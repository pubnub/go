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

func AssertManageMembersV2(t *testing.T, checkQueryParam, testContext bool, withFilter bool, withSort bool) {
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

	o := newManageChannelMembersBuilderV2(pn)
	if testContext {
		o = newManageChannelMembersBuilderV2WithContext(pn, pn.ctx)
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

	id0 := "id0"
	if withFilter {
		o.Filter("name like 'a*'")
	}
	sort := []string{"name", "created:desc"}
	if withSort {
		o.Sort(sort)
	}

	custom := make(map[string]interface{})
	custom["a1"] = "b1"
	custom["c1"] = "d1"

	uuid := PNChannelMembersUUID{
		ID: id0,
	}

	in := PNChannelMembersSet{
		UUID:   uuid,
		Custom: custom,
		Status: "active",
		Type:   "member",
	}

	inArr := []PNChannelMembersSet{
		in,
	}

	re := PNChannelMembersRemove{
		UUID: uuid,
	}

	reArr := []PNChannelMembersRemove{
		re,
	}
	o.Set(inArr)
	o.Remove(reArr)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/channels/%s/uuids", pn.Config.SubscribeKey, spaceID),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	expectedBody := "{\"set\":[{\"uuid\":{\"id\":\"id0\"},\"custom\":{\"a1\":\"b1\",\"c1\":\"d1\"},\"status\":\"active\",\"type\":\"member\"}],\"delete\":[{\"uuid\":{\"id\":\"id0\"}}]}"

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

func TestManageMembersV2(t *testing.T) {
	AssertManageMembersV2(t, true, false, false, false)
}

func TestManageMembersV2Context(t *testing.T) {
	AssertManageMembersV2(t, true, true, false, false)
}

func TestManageMembersV2WithFilter(t *testing.T) {
	AssertManageMembersV2(t, true, false, true, false)
}

func TestManageMembersV2WithFilterContext(t *testing.T) {
	AssertManageMembersV2(t, true, true, true, false)
}

func TestManageMembersV2WithSort(t *testing.T) {
	AssertManageMembersV2(t, true, false, false, true)
}

func TestManageMembersV2WithSortContext(t *testing.T) {
	AssertManageMembersV2(t, true, true, false, true)
}

func TestManageMembersV2WithFilterWithSort(t *testing.T) {
	AssertManageMembersV2(t, true, false, true, true)
}

func TestManageMembersV2WithFilterWithSortContext(t *testing.T) {
	AssertManageMembersV2(t, true, true, true, true)
}

func TestManageMembersV2ResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNManageMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestManageMembersV2ResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":[{"id":"userid4","custom":{"a1":"b1","c1":"d1"},"status":"active","type":"member","uuid":{"id":"userid4","name":"userid4name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"status":"active","type":"user","created":"2019-08-23T10:36:27.083453Z","updated":"2019-08-23T10:36:27.083453Z","eTag":"AbuLvdnC9JnYEA"},"created":"2019-08-23T10:41:35.503214Z","updated":"2019-08-23T10:41:35.503214Z","eTag":"AZK3l4nQsrWG9gE"}],"totalCount":1,"next":"MQ"}`)

	r, _, err := newPNManageMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(1, r.TotalCount)
	assert.Equal("MQ", r.Next)
	assert.Equal("userid4", r.Data[0].ID)
	assert.Equal("2019-08-23T10:41:35.503214Z", r.Data[0].Created)
	assert.Equal("2019-08-23T10:41:35.503214Z", r.Data[0].Updated)
	assert.Equal("AZK3l4nQsrWG9gE", r.Data[0].ETag)
	assert.Equal("b1", r.Data[0].Custom["a1"])
	assert.Equal("d1", r.Data[0].Custom["c1"])
	assert.Equal("active", r.Data[0].Status)
	assert.Equal("member", r.Data[0].Type)
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
	assert.Equal("active", r.Data[0].UUID.Status)
	assert.Equal("user", r.Data[0].UUID.Type)

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestManageMembersV2ValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newManageMembersOptsV2(pn, pn.ctx)
	opts.Channel = "test-channel"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestManageMembersV2ValidateMissingChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)
	opts.Channel = ""

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestManageMembersV2ValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)
	opts.Channel = "test-channel"

	assert.Nil(opts.validate())
}

func TestManageMembersV2ValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)
	opts.Channel = "test-channel"
	opts.Include = []string{"custom", "uuid"}
	opts.QueryParam = map[string]string{"param": "value"}
	opts.MembersSet = []PNChannelMembersSet{{UUID: PNChannelMembersUUID{ID: "user1"}}}
	opts.MembersRemove = []PNChannelMembersRemove{{UUID: PNChannelMembersUUID{ID: "user2"}}}

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestManageMembersV2HTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	assert.Equal("PATCH", opts.httpMethod())
}

func TestManageMembersV2OperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	assert.Equal(PNManageMembersOperation, opts.operationType())
}

func TestManageMembersV2IsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestManageMembersV2Timeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (12 setters)

func TestManageMembersV2BuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newManageChannelMembersBuilderV2(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
	assert.Equal(manageMembersLimitV2, builder.opts.Limit) // Default limit (100)
}

func TestManageMembersV2BuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newManageChannelMembersBuilderV2WithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestManageMembersV2BuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newManageChannelMembersBuilderV2(pn)

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

	// Test Set setter (CRITICAL for PATCH operation)
	setMembers := []PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: map[string]interface{}{"role": "admin"}},
		{UUID: PNChannelMembersUUID{ID: "user2"}, Custom: map[string]interface{}{"role": "member"}},
	}
	builder.Set(setMembers)
	assert.Equal(setMembers, builder.opts.MembersSet)

	// Test Remove setter (CRITICAL for PATCH operation)
	removeMembers := []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user3"}},
		{UUID: PNChannelMembersUUID{ID: "user4"}},
	}
	builder.Remove(removeMembers)
	assert.Equal(removeMembers, builder.opts.MembersRemove)

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

func TestManageMembersV2BuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	sort := []string{"name"}
	queryParam := map[string]string{"key": "value"}
	setMembers := []PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: map[string]interface{}{"role": "admin"}},
	}
	removeMembers := []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user2"}},
	}
	transport := &http.Transport{}

	builder := newManageChannelMembersBuilderV2(pn)
	result := builder.Channel("test-channel").
		Include(include).
		Limit(75).
		Start("start").
		End("end").
		Count(true).
		Filter("filter").
		Sort(sort).
		Set(setMembers).
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
	assert.Equal(setMembers, builder.opts.MembersSet)
	assert.Equal(removeMembers, builder.opts.MembersRemove)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

func TestManageMembersV2BuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newManageChannelMembersBuilderV2(pn)

	// Verify default values
	assert.Equal("", builder.opts.Channel)
	assert.Nil(builder.opts.Include)
	assert.Equal(manageMembersLimitV2, builder.opts.Limit) // 100
	assert.Equal("", builder.opts.Start)
	assert.Equal("", builder.opts.End)
	assert.Equal(false, builder.opts.Count)
	assert.Equal("", builder.opts.Filter)
	assert.Nil(builder.opts.Sort)
	assert.Nil(builder.opts.MembersSet)
	assert.Nil(builder.opts.MembersRemove)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestManageMembersV2BuilderIncludeTypes(t *testing.T) {
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
			builder := newManageChannelMembersBuilderV2(pn)
			builder.Include(tc.includes)

			expectedInclude := EnumArrayToStringArray(tc.includes)
			assert.Equal(expectedInclude, builder.opts.Include)
		})
	}
}

func TestManageMembersV2BuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	transport := &http.Transport{}
	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	sort := []string{"name", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}
	setMembers := []PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: map[string]interface{}{"role": "admin"}},
		{UUID: PNChannelMembersUUID{ID: "user2"}, Custom: map[string]interface{}{"level": 5}},
	}
	removeMembers := []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user3"}},
		{UUID: PNChannelMembersUUID{ID: "user4"}},
	}

	// Test all 12 setters in chain
	builder := newManageChannelMembersBuilderV2(pn).
		Channel("test-channel").
		Include(include).
		Limit(50).
		Start("start-token").
		End("end-token").
		Count(true).
		Filter("name LIKE 'test*'").
		Sort(sort).
		Set(setMembers).
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
	assert.Equal(setMembers, builder.opts.MembersSet)
	assert.Equal(removeMembers, builder.opts.MembersRemove)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestManageMembersV2BuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)
	opts.Channel = "test-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/channels/test-channel/uuids"
	assert.Equal(expected, path)
}

func TestManageMembersV2BuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newManageMembersOptsV2(pn, pn.ctx)
	opts.Channel = "my-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/channels/my-channel/uuids"
	assert.Equal(expected, path)
}

func TestManageMembersV2BuildPathWithSpecialCharsInChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)
	opts.Channel = "channel-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")
	assert.Contains(path, "channel-with-special@chars#and$symbols")
	assert.Contains(path, "/uuids")
}

func TestManageMembersV2BuildPathWithUnicodeChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)
	opts.Channel = "频道名称-канал-チャンネル"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/")
	assert.Contains(path, "频道名称-канал-チャンネル")
	assert.Contains(path, "/uuids")
}

// JSON Body Building Tests (CRITICAL for dual PATCH operation)

func TestManageMembersV2BuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"set":null,"delete":null}`
	assert.Equal(expected, string(body))
}

func TestManageMembersV2BuildBodySetOnly(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	opts.MembersSet = []PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: map[string]interface{}{"role": "admin"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"set":[{"uuid":{"id":"user1"},"custom":{"role":"admin"},"status":"","type":""}],"delete":null}`
	assert.Equal(expected, string(body))
}

func TestManageMembersV2BuildBodyRemoveOnly(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	opts.MembersRemove = []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user1"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"set":null,"delete":[{"uuid":{"id":"user1"}}]}`
	assert.Equal(expected, string(body))
}

func TestManageMembersV2BuildBodyCombinedOperations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	opts.MembersSet = []PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: map[string]interface{}{"role": "admin"}},
		{UUID: PNChannelMembersUUID{ID: "user2"}, Custom: map[string]interface{}{"level": 5}},
	}
	opts.MembersRemove = []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user3"}},
		{UUID: PNChannelMembersUUID{ID: "user4"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"set":[{"uuid":{"id":"user1"},"custom":{"role":"admin"},"status":"","type":""},{"uuid":{"id":"user2"},"custom":{"level":5},"status":"","type":""}],"delete":[{"uuid":{"id":"user3"}},{"uuid":{"id":"user4"}}]}`
	assert.Equal(expected, string(body))
}

func TestManageMembersV2BuildBodyWithUnicodeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	opts.MembersSet = []PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "用户123"}, Custom: map[string]interface{}{"名字": "张三", "角色": "管理员"}},
		{UUID: PNChannelMembersUUID{ID: "пользователь456"}, Custom: map[string]interface{}{"имя": "Иван", "роль": "участник"}},
	}
	opts.MembersRemove = []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "ユーザー789"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"id":"用户123"`)
	assert.Contains(string(body), `"名字":"张三"`)
	assert.Contains(string(body), `"id":"пользователь456"`)
	assert.Contains(string(body), `"id":"ユーザー789"`)
	assert.Contains(string(body), `"set":[`)
	assert.Contains(string(body), `"delete":[`)
}

func TestManageMembersV2BuildBodyWithComplexCustomData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	opts.MembersSet = []PNChannelMembersSet{
		{
			UUID: PNChannelMembersUUID{ID: "user1"},
			Custom: map[string]interface{}{
				"role":        "admin",
				"permissions": []string{"read", "write", "delete"},
				"profile": map[string]interface{}{
					"name":   "John Doe",
					"email":  "john@example.com",
					"active": true,
				},
				"level": 10,
			},
		},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"role":"admin"`)
	assert.Contains(string(body), `"permissions":["read","write","delete"]`)
	assert.Contains(string(body), `"profile":{`)
	assert.Contains(string(body), `"level":10`)
}

func TestManageMembersV2BuildBodyWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	opts.MembersSet = []PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "user@domain.com"}, Custom: map[string]interface{}{"email": "test@example.com"}},
		{UUID: PNChannelMembersUUID{ID: "user-with-dashes"}, Custom: map[string]interface{}{"name": "John O'Connor"}},
	}
	opts.MembersRemove = []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user_with_underscores"}},
		{UUID: PNChannelMembersUUID{ID: "user.with.dots"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"id":"user@domain.com"`)
	assert.Contains(string(body), `"id":"user-with-dashes"`)
	assert.Contains(string(body), `"id":"user_with_underscores"`)
	assert.Contains(string(body), `"id":"user.with.dots"`)
	assert.Contains(string(body), `"email":"test@example.com"`)
}

func TestManageMembersV2BuildBodyEmptyArrays(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	opts.MembersSet = []PNChannelMembersSet{}
	opts.MembersRemove = []PNChannelMembersRemove{}

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"set":[],"delete":[]}`
	assert.Equal(expected, string(body))
}

func TestManageMembersV2BuildBodySetWithNilCustom(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	opts.MembersSet = []PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: nil},
		{UUID: PNChannelMembersUUID{ID: "user2"}}, // Custom field not set
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"id":"user1"`)
	assert.Contains(string(body), `"id":"user2"`)
	assert.Contains(string(body), `"custom":null`)
}

// Query Parameter Tests

func TestManageMembersV2BuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("0", query.Get("limit")) // Default limit not set until builder initialization
	assert.Equal("0", query.Get("count")) // Default count=false
}

func TestManageMembersV2BuildQueryWithInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	opts.Include = []string{"custom", "uuid"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	includeValue := query.Get("include")
	assert.Contains(includeValue, "custom")
	assert.Contains(includeValue, "uuid")
}

func TestManageMembersV2BuildQueryWithPagination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

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

func TestManageMembersV2BuildQueryWithFilterAndSort(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

	opts.Filter = "custom.role == 'admin'"
	opts.Sort = []string{"name", "created:desc"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("custom.role == 'admin'", query.Get("filter"))

	sortValue := query.Get("sort")
	assert.Contains(sortValue, "name")
	assert.Contains(sortValue, "created:desc")
}

func TestManageMembersV2BuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

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

func TestManageMembersV2BuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

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

// PATCH-Specific Tests (Dual Operation Characteristics)

func TestManageMembersV2PatchOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newManageChannelMembersBuilderV2(pn)
	builder.Channel("test-channel")
	builder.Set([]PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: map[string]interface{}{"role": "admin"}},
	})
	builder.Remove([]PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user2"}},
	})

	// Verify it's a PATCH operation
	assert.Equal("PATCH", builder.opts.httpMethod())

	// PATCH operations have JSON body with both set and delete structures
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"set":[`)
	assert.Contains(string(body), `"delete":[`)
	assert.Contains(string(body), `"uuid":{"id":"user1"}`)
	assert.Contains(string(body), `"uuid":{"id":"user2"}`)
	assert.Contains(string(body), `"role":"admin"`)

	// Should have proper path for member management
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/channels/test-channel/uuids")
}

func TestManageMembersV2DefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newManageChannelMembersBuilderV2(pn)

	// Should have default limit set to manageMembersLimitV2 (100)
	assert.Equal(manageMembersLimitV2, builder.opts.Limit)
	assert.Equal(100, builder.opts.Limit)

	// Should be included in query
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("limit"))
}

func TestManageMembersV2DualOperationValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name          string
		setMembers    []PNChannelMembersSet
		removeMembers []PNChannelMembersRemove
		expectedJSON  string
	}{
		{
			name:          "Both nil",
			setMembers:    nil,
			removeMembers: nil,
			expectedJSON:  `{"set":null,"delete":null}`,
		},
		{
			name:          "Both empty",
			setMembers:    []PNChannelMembersSet{},
			removeMembers: []PNChannelMembersRemove{},
			expectedJSON:  `{"set":[],"delete":[]}`,
		},
		{
			name: "Set only",
			setMembers: []PNChannelMembersSet{
				{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: map[string]interface{}{"role": "admin"}},
			},
			removeMembers: nil,
			expectedJSON:  `{"set":[{"uuid":{"id":"user1"},"custom":{"role":"admin"},"status":"","type":""}],"delete":null}`,
		},
		{
			name:       "Remove only",
			setMembers: nil,
			removeMembers: []PNChannelMembersRemove{
				{UUID: PNChannelMembersUUID{ID: "user1"}},
			},
			expectedJSON: `{"set":null,"delete":[{"uuid":{"id":"user1"}}]}`,
		},
		{
			name: "Both operations",
			setMembers: []PNChannelMembersSet{
				{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: map[string]interface{}{"role": "admin"}},
			},
			removeMembers: []PNChannelMembersRemove{
				{UUID: PNChannelMembersUUID{ID: "user2"}},
			},
			expectedJSON: `{"set":[{"uuid":{"id":"user1"},"custom":{"role":"admin"},"status":"","type":""}],"delete":[{"uuid":{"id":"user2"}}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newManageMembersOptsV2(pn, pn.ctx)
			opts.MembersSet = tc.setMembers
			opts.MembersRemove = tc.removeMembers

			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Equal(tc.expectedJSON, string(body))
		})
	}
}

func TestManageMembersV2ResponseStructureAfterManagement(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newManageChannelMembersBuilderV2(pn)
	builder.Channel("test-channel")
	builder.Set([]PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: map[string]interface{}{"role": "admin"}},
	})
	builder.Remove([]PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user2"}},
	})

	// Response should contain current members after management operations
	// This is tested in the existing TestManageMembersV2ResponseValuePass
	// but verify the operation is configured correctly
	opts := builder.opts

	// Verify operation is configured correctly
	assert.Equal("PATCH", opts.httpMethod())
	assert.Equal(PNManageMembersOperation, opts.operationType())
	assert.NotNil(opts.MembersSet)
	assert.NotNil(opts.MembersRemove)
	assert.Equal(1, len(opts.MembersSet))
	assert.Equal(1, len(opts.MembersRemove))
}

func TestManageMembersV2AtomicOperations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test scenarios where the same user might be in both set and remove
	// This tests the atomic nature of the operation
	testCases := []struct {
		name          string
		setMembers    []PNChannelMembersSet
		removeMembers []PNChannelMembersRemove
		description   string
	}{
		{
			name: "Same user in set and remove (update scenario)",
			setMembers: []PNChannelMembersSet{
				{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: map[string]interface{}{"role": "admin", "updated": true}},
			},
			removeMembers: []PNChannelMembersRemove{
				{UUID: PNChannelMembersUUID{ID: "user1"}}, // Remove old version
			},
			description: "User1 gets updated with new data atomically",
		},
		{
			name: "Bulk management with mixed operations",
			setMembers: []PNChannelMembersSet{
				{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: map[string]interface{}{"role": "admin"}},
				{UUID: PNChannelMembersUUID{ID: "user2"}, Custom: map[string]interface{}{"role": "member"}},
				{UUID: PNChannelMembersUUID{ID: "user3"}, Custom: map[string]interface{}{"role": "guest"}},
			},
			removeMembers: []PNChannelMembersRemove{
				{UUID: PNChannelMembersUUID{ID: "user4"}},
				{UUID: PNChannelMembersUUID{ID: "user5"}},
			},
			description: "Add 3 users and remove 2 users in single atomic operation",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newManageChannelMembersBuilderV2(pn)
			builder.Channel("test-channel")
			builder.Set(tc.setMembers)
			builder.Remove(tc.removeMembers)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid body with both operations
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Contains(string(body), `"set":[`)
			assert.Contains(string(body), `"delete":[`)

			// Verify all set members are in body
			for _, setMember := range tc.setMembers {
				assert.Contains(string(body), fmt.Sprintf(`"id":"%s"`, setMember.UUID.ID))
			}

			// Verify all remove members are in body
			for _, removeMember := range tc.removeMembers {
				assert.Contains(string(body), fmt.Sprintf(`"id":"%s"`, removeMember.UUID.ID))
			}
		})
	}
}

// Comprehensive Edge Case Tests

func TestManageMembersV2WithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*manageChannelMembersBuilderV2)
	}{
		{
			name: "Very long channel name",
			setupFn: func(builder *manageChannelMembersBuilderV2) {
				longChannel := strings.Repeat("VeryLongChannel", 50) // 750 characters
				builder.Channel(longChannel)
			},
		},
		{
			name: "Large number of members to set",
			setupFn: func(builder *manageChannelMembersBuilderV2) {
				var largeSetList []PNChannelMembersSet
				for i := 0; i < 50; i++ {
					largeSetList = append(largeSetList, PNChannelMembersSet{
						UUID: PNChannelMembersUUID{ID: fmt.Sprintf("set_user_%d", i)},
						Custom: map[string]interface{}{
							"role":  fmt.Sprintf("role_%d", i),
							"level": i,
						},
					})
				}
				builder.Set(largeSetList)
			},
		},
		{
			name: "Large number of members to remove",
			setupFn: func(builder *manageChannelMembersBuilderV2) {
				var largeRemoveList []PNChannelMembersRemove
				for i := 0; i < 50; i++ {
					largeRemoveList = append(largeRemoveList, PNChannelMembersRemove{
						UUID: PNChannelMembersUUID{ID: fmt.Sprintf("remove_user_%d", i)},
					})
				}
				builder.Remove(largeRemoveList)
			},
		},
		{
			name: "Large filter expression",
			setupFn: func(builder *manageChannelMembersBuilderV2) {
				largeFilter := "(" + strings.Repeat("custom.field == 'value' OR ", 100) + "custom.final == 'end')"
				builder.Filter(largeFilter)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *manageChannelMembersBuilderV2) {
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
			builder := newManageChannelMembersBuilderV2(pn)
			builder.Channel("test-channel")
			builder.Set([]PNChannelMembersSet{
				{UUID: PNChannelMembersUUID{ID: "baseline-set-user"}, Custom: map[string]interface{}{"role": "admin"}},
			})
			builder.Remove([]PNChannelMembersRemove{
				{UUID: PNChannelMembersUUID{ID: "baseline-remove-user"}},
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
			assert.Contains(string(body), `"set":`)
			assert.Contains(string(body), `"delete":`)
		})
	}
}

func TestManageMembersV2SpecialCharacterHandling(t *testing.T) {
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
			builder := newManageChannelMembersBuilderV2(pn)
			builder.Channel(specialString)
			builder.Filter(fmt.Sprintf("custom.field == '%s'", specialString))
			builder.QueryParam(map[string]string{
				"special_field": specialString,
			})
			builder.Set([]PNChannelMembersSet{
				{UUID: PNChannelMembersUUID{ID: specialString}, Custom: map[string]interface{}{"name": specialString}},
			})
			builder.Remove([]PNChannelMembersRemove{
				{UUID: PNChannelMembersUUID{ID: fmt.Sprintf("remove_%s", specialString)}},
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
			assert.Contains(string(body), `"set":`)
			assert.Contains(string(body), `"delete":`)
		})
	}
}

func TestManageMembersV2ParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		channel     string
		limit       int
		filter      string
		setCount    int
		removeCount int
	}{
		{
			name:        "Empty string channel",
			channel:     "",
			limit:       1,
			filter:      "",
			setCount:    0,
			removeCount: 0,
		},
		{
			name:        "Single character channel",
			channel:     "a",
			limit:       1,
			filter:      "a",
			setCount:    1,
			removeCount: 1,
		},
		{
			name:        "Unicode-only channel",
			channel:     "测试",
			limit:       50,
			filter:      "测试 == '值'",
			setCount:    2,
			removeCount: 1,
		},
		{
			name:        "Minimum limit",
			channel:     "test",
			limit:       1,
			filter:      "simple",
			setCount:    1,
			removeCount: 1,
		},
		{
			name:        "Large limit",
			channel:     "test",
			limit:       1000,
			filter:      "complex.nested == 'value'",
			setCount:    10,
			removeCount: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newManageChannelMembersBuilderV2(pn)
			builder.Channel(tc.channel)
			builder.Limit(tc.limit)
			if tc.filter != "" {
				builder.Filter(tc.filter)
			}

			// Add specified number of members to set
			var setMembers []PNChannelMembersSet
			for i := 0; i < tc.setCount; i++ {
				setMembers = append(setMembers, PNChannelMembersSet{
					UUID:   PNChannelMembersUUID{ID: fmt.Sprintf("set_user_%d", i)},
					Custom: map[string]interface{}{"index": i},
				})
			}
			if len(setMembers) > 0 {
				builder.Set(setMembers)
			}

			// Add specified number of members to remove
			var removeMembers []PNChannelMembersRemove
			for i := 0; i < tc.removeCount; i++ {
				removeMembers = append(removeMembers, PNChannelMembersRemove{
					UUID: PNChannelMembersUUID{ID: fmt.Sprintf("remove_user_%d", i)},
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

func TestManageMembersV2ComplexManagementScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*manageChannelMembersBuilderV2)
		validateFn func(*testing.T, []byte)
	}{
		{
			name: "Mixed character UUIDs with complex data",
			setupFn: func(builder *manageChannelMembersBuilderV2) {
				builder.Set([]PNChannelMembersSet{
					{UUID: PNChannelMembersUUID{ID: "user-english"}, Custom: map[string]interface{}{"role": "admin", "lang": "en"}},
					{UUID: PNChannelMembersUUID{ID: "用户中文"}, Custom: map[string]interface{}{"角色": "成员", "语言": "zh"}},
					{UUID: PNChannelMembersUUID{ID: "пользователь"}, Custom: map[string]interface{}{"роль": "гость", "язык": "ru"}},
				})
				builder.Remove([]PNChannelMembersRemove{
					{UUID: PNChannelMembersUUID{ID: "ユーザー"}},
				})
			},
			validateFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				assert.Contains(bodyStr, `"id":"user-english"`)
				assert.Contains(bodyStr, `"id":"用户中文"`)
				assert.Contains(bodyStr, `"id":"пользователь"`)
				assert.Contains(bodyStr, `"id":"ユーザー"`)
				assert.Contains(bodyStr, `"role":"admin"`)
				assert.Contains(bodyStr, `"角色":"成员"`)
				assert.Contains(bodyStr, `"роль":"гость"`)
			},
		},
		{
			name: "Email-like UUIDs with professional data",
			setupFn: func(builder *manageChannelMembersBuilderV2) {
				builder.Set([]PNChannelMembersSet{
					{UUID: PNChannelMembersUUID{ID: "admin@company.com"}, Custom: map[string]interface{}{
						"role":        "admin",
						"department":  "IT",
						"permissions": []string{"read", "write", "delete"},
					}},
					{UUID: PNChannelMembersUUID{ID: "user+dev@example.org"}, Custom: map[string]interface{}{
						"role": "developer",
						"team": "backend",
					}},
				})
				builder.Remove([]PNChannelMembersRemove{
					{UUID: PNChannelMembersUUID{ID: "temp.user@domain.co.uk"}},
					{UUID: PNChannelMembersUUID{ID: "guest_user@test-domain.com"}},
				})
			},
			validateFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				assert.Contains(bodyStr, `"id":"admin@company.com"`)
				assert.Contains(bodyStr, `"id":"user+dev@example.org"`)
				assert.Contains(bodyStr, `"id":"temp.user@domain.co.uk"`)
				assert.Contains(bodyStr, `"id":"guest_user@test-domain.com"`)
				assert.Contains(bodyStr, `"department":"IT"`)
				assert.Contains(bodyStr, `"permissions":["read","write","delete"]`)
			},
		},
		{
			name: "Atomic update scenario (same user in set and remove)",
			setupFn: func(builder *manageChannelMembersBuilderV2) {
				builder.Set([]PNChannelMembersSet{
					{UUID: PNChannelMembersUUID{ID: "user123"}, Custom: map[string]interface{}{
						"role":    "admin",
						"updated": true,
						"version": 2,
					}},
				})
				builder.Remove([]PNChannelMembersRemove{
					{UUID: PNChannelMembersUUID{ID: "user123"}}, // Remove old version
				})
			},
			validateFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				// Should contain user123 in both set and delete sections
				setSection := strings.Split(bodyStr, `"delete"`)[0]
				deleteSection := strings.Split(bodyStr, `"delete"`)[1]
				assert.Contains(setSection, `"id":"user123"`)
				assert.Contains(deleteSection, `"id":"user123"`)
				assert.Contains(bodyStr, `"updated":true`)
				assert.Contains(bodyStr, `"version":2`)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newManageChannelMembersBuilderV2(pn)
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

func TestManageMembersV2ExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newManageChannelMembersBuilderV2(pn)
	builder.Channel("test-channel")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestManageMembersV2PathBuildingEdgeCases(t *testing.T) {
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
			opts := newManageMembersOptsV2(pn, pn.ctx)
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

func TestManageMembersV2QueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*manageMembersOptsV2)
		expectError bool
	}{
		{
			name: "Nil query params",
			setupOpts: func(opts *manageMembersOptsV2) {
				opts.QueryParam = nil
			},
			expectError: false,
		},
		{
			name: "Empty query params",
			setupOpts: func(opts *manageMembersOptsV2) {
				opts.QueryParam = map[string]string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *manageMembersOptsV2) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *manageMembersOptsV2) {
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
			opts := newManageMembersOptsV2(pn, pn.ctx)
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

func TestManageMembersV2BuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newManageChannelMembersBuilderV2(pn)

	include := []PNChannelMembersInclude{PNChannelMembersIncludeCustom}
	sort := []string{"name:asc", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}
	setMembers := []PNChannelMembersSet{
		{UUID: PNChannelMembersUUID{ID: "user1"}, Custom: map[string]interface{}{"role": "admin"}},
		{UUID: PNChannelMembersUUID{ID: "user2"}, Custom: map[string]interface{}{"role": "member"}},
	}
	removeMembers := []PNChannelMembersRemove{
		{UUID: PNChannelMembersUUID{ID: "user3"}},
		{UUID: PNChannelMembersUUID{ID: "user4"}},
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
		Set(setMembers).
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
	assert.Equal(setMembers, builder.opts.MembersSet)
	assert.Equal(removeMembers, builder.opts.MembersRemove)
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
	expectedBody := `{"set":[{"uuid":{"id":"user1"},"custom":{"role":"admin"},"status":"","type":""},{"uuid":{"id":"user2"},"custom":{"role":"member"},"status":"","type":""}],"delete":[{"uuid":{"id":"user3"}},{"uuid":{"id":"user4"}}]}`
	assert.Equal(expectedBody, string(body))
}

func TestManageMembersV2ResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newManageMembersOptsV2(pn, pn.ctx)

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
			name:        "Valid response with managed member data",
			jsonBytes:   []byte(`{"status":200,"data":[{"id":"user1","uuid":{"id":"user1","name":"User 1"},"custom":{"role":"admin"}}],"totalCount":1,"next":"abc","prev":"xyz"}`),
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
			resp, _, err := newPNManageMembersResponse(tc.jsonBytes, opts, StatusResponse{})

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

// Dual Operation Complexity Tests

func TestManageMembersV2DualOperationComplexity(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		description string
		setupFn     func(*manageChannelMembersBuilderV2)
		verifyFn    func(*testing.T, []byte)
	}{
		{
			name:        "Complex nested custom data",
			description: "Set members with deeply nested custom data structures",
			setupFn: func(builder *manageChannelMembersBuilderV2) {
				builder.Set([]PNChannelMembersSet{
					{UUID: PNChannelMembersUUID{ID: "admin_user"}, Custom: map[string]interface{}{
						"profile": map[string]interface{}{
							"personal": map[string]interface{}{
								"name":  "John Doe",
								"email": "john@example.com",
								"preferences": map[string]interface{}{
									"notifications": true,
									"theme":         "dark",
									"languages":     []string{"en", "es", "fr"},
								},
							},
							"professional": map[string]interface{}{
								"title":      "Senior Developer",
								"department": "Engineering",
								"skills":     []string{"Go", "JavaScript", "Python"},
								"certifications": []map[string]interface{}{
									{"name": "AWS Certified", "year": 2023},
									{"name": "Kubernetes Expert", "year": 2022},
								},
							},
						},
						"permissions": map[string]interface{}{
							"read":   true,
							"write":  true,
							"admin":  true,
							"scopes": []string{"users", "channels", "analytics"},
						},
					}},
				})
			},
			verifyFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				assert.Contains(bodyStr, `"name":"John Doe"`)
				assert.Contains(bodyStr, `"skills":["Go","JavaScript","Python"]`)
				assert.Contains(bodyStr, `"certifications":[`)
				assert.Contains(bodyStr, `"scopes":["users","channels","analytics"]`)
			},
		},
		{
			name:        "Large scale atomic operation",
			description: "Perform large scale set and remove operations atomically",
			setupFn: func(builder *manageChannelMembersBuilderV2) {
				var setMembers []PNChannelMembersSet
				var removeMembers []PNChannelMembersRemove

				// Add 20 new members with various roles
				for i := 0; i < 20; i++ {
					setMembers = append(setMembers, PNChannelMembersSet{
						UUID: PNChannelMembersUUID{ID: fmt.Sprintf("new_user_%d", i)},
						Custom: map[string]interface{}{
							"role":     []string{"admin", "member", "guest"}[i%3],
							"joinDate": fmt.Sprintf("2023-%02d-01", (i%12)+1),
							"level":    i + 1,
							"active":   i%2 == 0,
						},
					})
				}

				// Remove 15 old members
				for i := 0; i < 15; i++ {
					removeMembers = append(removeMembers, PNChannelMembersRemove{
						UUID: PNChannelMembersUUID{ID: fmt.Sprintf("old_user_%d", i)},
					})
				}

				builder.Set(setMembers)
				builder.Remove(removeMembers)
			},
			verifyFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				// Verify presence of set and delete sections
				assert.Contains(bodyStr, `"set":[`)
				assert.Contains(bodyStr, `"delete":[`)

				// Verify some specific entries
				assert.Contains(bodyStr, `"id":"new_user_0"`)
				assert.Contains(bodyStr, `"id":"new_user_19"`)
				assert.Contains(bodyStr, `"id":"old_user_0"`)
				assert.Contains(bodyStr, `"id":"old_user_14"`)

				// Verify role distribution
				assert.Contains(bodyStr, `"role":"admin"`)
				assert.Contains(bodyStr, `"role":"member"`)
				assert.Contains(bodyStr, `"role":"guest"`)
			},
		},
		{
			name:        "Mixed data types and encodings",
			description: "Handle various data types and character encodings",
			setupFn: func(builder *manageChannelMembersBuilderV2) {
				builder.Set([]PNChannelMembersSet{
					{UUID: PNChannelMembersUUID{ID: "unicode_user_测试"}, Custom: map[string]interface{}{
						"数字":  42,
						"布尔值": true,
						"字符串": "测试字符串",
						"数组":  []interface{}{"项目1", "项目2", 123, false},
						"对象": map[string]interface{}{
							"嵌套字段": "嵌套值",
							"数字字段": 3.14,
						},
					}},
					{UUID: PNChannelMembersUUID{ID: "russian_пользователь"}, Custom: map[string]interface{}{
						"число":  99,
						"булево": false,
						"строка": "русская строка",
						"массив": []interface{}{"элемент1", "элемент2"},
					}},
				})
				builder.Remove([]PNChannelMembersRemove{
					{UUID: PNChannelMembersUUID{ID: "japanese_ユーザー"}},
					{UUID: PNChannelMembersUUID{ID: "arabic_مستخدم"}},
				})
			},
			verifyFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				assert.Contains(bodyStr, `"id":"unicode_user_测试"`)
				assert.Contains(bodyStr, `"数字":42`)
				assert.Contains(bodyStr, `"字符串":"测试字符串"`)
				assert.Contains(bodyStr, `"id":"russian_пользователь"`)
				assert.Contains(bodyStr, `"строка":"русская строка"`)
				assert.Contains(bodyStr, `"id":"japanese_ユーザー"`)
				assert.Contains(bodyStr, `"id":"arabic_مستخدم"`)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newManageChannelMembersBuilderV2(pn)
			builder.Channel("test-channel")
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid JSON body
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.NotEmpty(body)

			// Run verification
			tc.verifyFn(t, body)
		})
	}
}
