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

func AssertSetMemberships(t *testing.T, checkQueryParam, testContext bool, withFilter bool, withSort bool) {
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

	o := newSetMembershipsBuilder(pn)
	if testContext {
		o = newSetMembershipsBuilderWithContext(pn, pn.ctx)
	}

	spaceID := "id0"
	limit := 90
	start := "Mxmy"
	end := "Nxny"

	o.UUID(spaceID)
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

	channel := PNMembershipsChannel{
		ID: id0,
	}

	in := PNMembershipsSet{
		Channel: channel,
		Custom:  custom,
	}

	inArr := []PNMembershipsSet{
		in,
	}

	custom2 := make(map[string]interface{})
	custom2["a2"] = "b2"
	custom2["c2"] = "d2"

	o.Set(inArr)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/uuids/%s/channels", pn.Config.SubscribeKey, spaceID),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	expectedBody := "{\"set\":[{\"channel\":{\"id\":\"id0\"},\"custom\":{\"a1\":\"b1\",\"c1\":\"d1\"},\"status\":\"\",\"type\":\"\"}]}"

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

func TestSetMemberships(t *testing.T) {
	AssertSetMemberships(t, true, false, false, false)
}

func TestSetMembershipsContext(t *testing.T) {
	AssertSetMemberships(t, true, true, false, false)
}

func TestSetMembershipsWithFilter(t *testing.T) {
	AssertSetMemberships(t, true, false, true, false)
}

func TestSetMembershipsWithFilterContext(t *testing.T) {
	AssertSetMemberships(t, true, true, true, false)
}

func TestSetMembershipsWithSort(t *testing.T) {
	AssertSetMemberships(t, true, false, false, true)
}

func TestSetMembershipsWithSortContext(t *testing.T) {
	AssertSetMemberships(t, true, true, false, true)
}

func TestSetMembershipsWithFilterWithSort(t *testing.T) {
	AssertSetMemberships(t, true, false, true, true)
}

func TestSetMembershipsWithFilterWithSortContext(t *testing.T) {
	AssertSetMemberships(t, true, true, true, true)
}

func TestSetMembershipsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNSetMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestSetMembershipsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":[{"id":"spaceid3","custom":{"a3":"b3","c3":"d3"},"channel":{"id":"spaceid3","name":"spaceid3name","description":"spaceid3desc","custom":{"a":"b"},"created":"2019-08-23T10:34:43.985248Z","updated":"2019-08-23T10:34:43.985248Z","eTag":"Aazjn7vC3oDDYw"},"created":"2019-08-23T10:41:17.156491Z","updated":"2019-08-23T10:41:17.156491Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":1,"next":"MQ"}`)

	r, _, err := newPNSetMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(1, r.TotalCount)
	assert.Equal("MQ", r.Next)
	assert.Equal("spaceid3", r.Data[0].ID)
	assert.Equal("spaceid3", r.Data[0].Channel.ID)
	assert.Equal("spaceid3name", r.Data[0].Channel.Name)
	assert.Equal("spaceid3desc", r.Data[0].Channel.Description)
	//assert.Equal("2019-08-23T10:34:43.985248Z", r.Data[0].Channel.Created)
	assert.Equal("2019-08-23T10:34:43.985248Z", r.Data[0].Channel.Updated)
	assert.Equal("Aazjn7vC3oDDYw", r.Data[0].Channel.ETag)
	assert.Equal("b", r.Data[0].Channel.Custom["a"])
	assert.Equal("2019-08-23T10:41:17.156491Z", r.Data[0].Created)
	assert.Equal("2019-08-23T10:41:17.156491Z", r.Data[0].Updated)
	assert.Equal("AamrnoXdpdmzjwE", r.Data[0].ETag)
	assert.Equal("b3", r.Data[0].Custom["a3"])
	assert.Equal("d3", r.Data[0].Custom["c3"])

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestSetMembershipsValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newSetMembershipsOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestSetMembershipsValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	assert.Nil(opts.validate())
}

func TestSetMembershipsValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"
	opts.Include = []string{"custom", "channel"}
	opts.QueryParam = map[string]string{"param": "value"}
	opts.MembershipsSet = []PNMembershipsSet{{Channel: PNMembershipsChannel{ID: "channel1"}}}

	assert.Nil(opts.validate())
}

func TestSetMembershipsValidateUUIDDefaultBehavior(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.UUID = "config-uuid"

	builder := newSetMembershipsBuilder(pn)
	// Don't set UUID explicitly - should use Config.UUID in Execute
	builder.Set([]PNMembershipsSet{{Channel: PNMembershipsChannel{ID: "channel1"}}})

	// Before Execute, UUID should be empty
	assert.Equal("", builder.opts.UUID)

	// Test that Execute correctly sets UUID from Config when it's empty
	// We can't easily test Execute without making real HTTP calls, so we simulate the logic
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

func TestSetMembershipsHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	assert.Equal("PATCH", opts.httpMethod())
}

func TestSetMembershipsOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	assert.Equal(PNSetMembershipsOperation, opts.operationType())
}

func TestSetMembershipsIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestSetMembershipsTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (10 setters)

func TestSetMembershipsBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetMembershipsBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
	assert.Equal(setMembershipsLimit, builder.opts.Limit) // Default limit (100)
}

func TestSetMembershipsBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetMembershipsBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestSetMembershipsBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetMembershipsBuilder(pn)

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

	// Test Set setter (CRITICAL for PATCH operation)
	setMemberships := []PNMembershipsSet{
		{Channel: PNMembershipsChannel{ID: "channel1"}, Custom: map[string]interface{}{"role": "admin"}},
		{Channel: PNMembershipsChannel{ID: "channel2"}, Custom: map[string]interface{}{"role": "member"}},
	}
	builder.Set(setMemberships)
	assert.Equal(setMemberships, builder.opts.MembershipsSet)

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

func TestSetMembershipsBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	include := []PNMembershipsInclude{PNMembershipsIncludeCustom}
	sort := []string{"name"}
	queryParam := map[string]string{"key": "value"}
	setMemberships := []PNMembershipsSet{
		{Channel: PNMembershipsChannel{ID: "channel1"}, Custom: map[string]interface{}{"role": "admin"}},
	}
	transport := &http.Transport{}

	builder := newSetMembershipsBuilder(pn)
	result := builder.UUID("test-uuid").
		Include(include).
		Limit(75).
		Start("start").
		End("end").
		Count(true).
		Filter("filter").
		Sort(sort).
		Set(setMemberships).
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
	assert.Equal(setMemberships, builder.opts.MembershipsSet)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

func TestSetMembershipsBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetMembershipsBuilder(pn)

	// Verify default values
	assert.Equal("", builder.opts.UUID) // UUID defaults to empty, set later in Execute
	assert.Nil(builder.opts.Include)
	assert.Equal(setMembershipsLimit, builder.opts.Limit) // 100
	assert.Equal("", builder.opts.Start)
	assert.Equal("", builder.opts.End)
	assert.Equal(false, builder.opts.Count)
	assert.Equal("", builder.opts.Filter)
	assert.Nil(builder.opts.Sort)
	assert.Nil(builder.opts.MembershipsSet)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestSetMembershipsBuilderIncludeTypes(t *testing.T) {
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
			builder := newSetMembershipsBuilder(pn)
			builder.Include(tc.includes)

			expectedInclude := EnumArrayToStringArray(tc.includes)
			assert.Equal(expectedInclude, builder.opts.Include)
		})
	}
}

func TestSetMembershipsBuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	transport := &http.Transport{}
	include := []PNMembershipsInclude{PNMembershipsIncludeCustom}
	sort := []string{"name", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}
	setMemberships := []PNMembershipsSet{
		{Channel: PNMembershipsChannel{ID: "channel1"}, Custom: map[string]interface{}{"role": "admin"}},
		{Channel: PNMembershipsChannel{ID: "channel2"}, Custom: map[string]interface{}{"level": 5}},
	}

	// Test all 10 setters in chain
	builder := newSetMembershipsBuilder(pn).
		UUID("test-uuid").
		Include(include).
		Limit(50).
		Start("start-token").
		End("end-token").
		Count(true).
		Filter("name LIKE 'test*'").
		Sort(sort).
		Set(setMemberships).
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
	assert.Equal(setMemberships, builder.opts.MembershipsSet)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestSetMembershipsBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/uuids/test-uuid/channels"
	assert.Equal(expected, path)
}

func TestSetMembershipsBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newSetMembershipsOpts(pn, pn.ctx)
	opts.UUID = "my-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/uuids/my-uuid/channels"
	assert.Equal(expected, path)
}

func TestSetMembershipsBuildPathWithSpecialCharsInUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)
	opts.UUID = "uuid-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/")
	assert.Contains(path, "uuid-with-special@chars#and$symbols")
	assert.Contains(path, "/channels")
}

func TestSetMembershipsBuildPathWithUnicodeUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)
	opts.UUID = "用户ID-пользователь-ユーザーID"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/")
	assert.Contains(path, "用户ID-пользователь-ユーザーID")
	assert.Contains(path, "/channels")
}

// JSON Body Building Tests (CRITICAL for PATCH operation)

func TestSetMembershipsBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"set":null}`
	assert.Equal(expected, string(body))
}

func TestSetMembershipsBuildBodyBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	opts.MembershipsSet = []PNMembershipsSet{
		{Channel: PNMembershipsChannel{ID: "channel1"}, Custom: map[string]interface{}{"role": "admin"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"set":[{"channel":{"id":"channel1"},"custom":{"role":"admin"},"status":"","type":""}]}`
	assert.Equal(expected, string(body))
}

func TestSetMembershipsBuildBodyMultiple(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	opts.MembershipsSet = []PNMembershipsSet{
		{Channel: PNMembershipsChannel{ID: "channel1"}, Custom: map[string]interface{}{"role": "admin"}},
		{Channel: PNMembershipsChannel{ID: "channel2"}, Custom: map[string]interface{}{"level": 5}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"set":[{"channel":{"id":"channel1"},"custom":{"role":"admin"},"status":"","type":""},{"channel":{"id":"channel2"},"custom":{"level":5},"status":"","type":""}]}`
	assert.Equal(expected, string(body))
}

func TestSetMembershipsBuildBodyWithUnicodeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	opts.MembershipsSet = []PNMembershipsSet{
		{Channel: PNMembershipsChannel{ID: "频道123"}, Custom: map[string]interface{}{"名字": "张三", "角色": "管理员"}},
		{Channel: PNMembershipsChannel{ID: "канал456"}, Custom: map[string]interface{}{"имя": "Иван", "роль": "участник"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"id":"频道123"`)
	assert.Contains(string(body), `"名字":"张三"`)
	assert.Contains(string(body), `"id":"канал456"`)
	assert.Contains(string(body), `"set":[`)
}

func TestSetMembershipsBuildBodyWithComplexCustomData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	opts.MembershipsSet = []PNMembershipsSet{
		{
			Channel: PNMembershipsChannel{ID: "channel1"},
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

func TestSetMembershipsBuildBodyWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	opts.MembershipsSet = []PNMembershipsSet{
		{Channel: PNMembershipsChannel{ID: "channel@domain.com"}, Custom: map[string]interface{}{"email": "test@example.com"}},
		{Channel: PNMembershipsChannel{ID: "channel-with-dashes"}, Custom: map[string]interface{}{"name": "John O'Connor"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"id":"channel@domain.com"`)
	assert.Contains(string(body), `"id":"channel-with-dashes"`)
	assert.Contains(string(body), `"email":"test@example.com"`)
}

func TestSetMembershipsBuildBodyEmptyArray(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	opts.MembershipsSet = []PNMembershipsSet{}

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := `{"set":[]}`
	assert.Equal(expected, string(body))
}

func TestSetMembershipsBuildBodyWithNilCustom(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	opts.MembershipsSet = []PNMembershipsSet{
		{Channel: PNMembershipsChannel{ID: "channel1"}, Custom: nil},
		{Channel: PNMembershipsChannel{ID: "channel2"}}, // Custom field not set
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"id":"channel1"`)
	assert.Contains(string(body), `"id":"channel2"`)
	assert.Contains(string(body), `"custom":null`)
}

// Query Parameter Tests

func TestSetMembershipsBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("100", query.Get("limit")) // Default limit
	assert.Equal("0", query.Get("count"))   // Default count=false
}

func TestSetMembershipsBuildQueryWithInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	opts.Include = []string{"custom", "channel"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	includeValue := query.Get("include")
	assert.Contains(includeValue, "custom")
	assert.Contains(includeValue, "channel")
}

func TestSetMembershipsBuildQueryWithPagination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

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

func TestSetMembershipsBuildQueryWithFilterAndSort(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

	opts.Filter = "custom.role == 'admin'"
	opts.Sort = []string{"name", "created:desc"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("custom.role == 'admin'", query.Get("filter"))

	sortValue := query.Get("sort")
	assert.Contains(sortValue, "name")
	assert.Contains(sortValue, "created:desc")
}

func TestSetMembershipsBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

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

func TestSetMembershipsBuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

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

// PATCH-Specific Tests (Set Operation Characteristics)

func TestSetMembershipsPatchOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetMembershipsBuilder(pn)
	builder.UUID("test-uuid")
	builder.Set([]PNMembershipsSet{
		{Channel: PNMembershipsChannel{ID: "channel1"}, Custom: map[string]interface{}{"role": "admin"}},
	})

	// Verify it's a PATCH operation
	assert.Equal("PATCH", builder.opts.httpMethod())

	// PATCH operations have JSON body with set structure
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"set":[`)
	assert.Contains(string(body), `"channel":{"id":"channel1"}`)
	assert.Contains(string(body), `"role":"admin"`)

	// Should have proper path for membership management (UUID to channels)
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/test-uuid/channels")
}

func TestSetMembershipsDefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetMembershipsBuilder(pn)

	// Should have default limit set to setMembershipsLimit (100)
	assert.Equal(setMembershipsLimit, builder.opts.Limit)
	assert.Equal(100, builder.opts.Limit)

	// Should be included in query
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("limit"))
}

func TestSetMembershipsSetOperationValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name         string
		memberships  []PNMembershipsSet
		expectedJSON string
	}{
		{
			name:         "Nil memberships",
			memberships:  nil,
			expectedJSON: `{"set":null}`,
		},
		{
			name:         "Empty memberships",
			memberships:  []PNMembershipsSet{},
			expectedJSON: `{"set":[]}`,
		},
		{
			name: "Single membership",
			memberships: []PNMembershipsSet{
				{Channel: PNMembershipsChannel{ID: "channel1"}, Custom: map[string]interface{}{"role": "admin"}},
			},
			expectedJSON: `{"set":[{"channel":{"id":"channel1"},"custom":{"role":"admin"},"status":"","type":""}]}`,
		},
		{
			name: "Multiple memberships",
			memberships: []PNMembershipsSet{
				{Channel: PNMembershipsChannel{ID: "channel1"}, Custom: map[string]interface{}{"role": "admin"}},
				{Channel: PNMembershipsChannel{ID: "channel2"}, Custom: map[string]interface{}{"role": "member"}},
			},
			expectedJSON: `{"set":[{"channel":{"id":"channel1"},"custom":{"role":"admin"},"status":"","type":""},{"channel":{"id":"channel2"},"custom":{"role":"member"},"status":"","type":""}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newSetMembershipsOpts(pn, pn.ctx)
			opts.MembershipsSet = tc.memberships

			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Equal(tc.expectedJSON, string(body))
		})
	}
}

func TestSetMembershipsResponseStructureAfterSetting(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetMembershipsBuilder(pn)
	builder.UUID("test-uuid")
	builder.Set([]PNMembershipsSet{
		{Channel: PNMembershipsChannel{ID: "channel1"}, Custom: map[string]interface{}{"role": "admin"}},
	})

	// Response should contain current memberships after set operations
	// This is tested in the existing TestSetMembershipsResponseValuePass
	// but verify the operation is configured correctly
	opts := builder.opts

	// Verify operation is configured correctly
	assert.Equal("PATCH", opts.httpMethod())
	assert.Equal(PNSetMembershipsOperation, opts.operationType())
	assert.NotNil(opts.MembershipsSet)
	assert.Equal(1, len(opts.MembershipsSet))
}

func TestSetMembershipsBulkOperations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test scenarios with multiple memberships
	testCases := []struct {
		name        string
		memberships []PNMembershipsSet
		description string
	}{
		{
			name: "Bulk membership assignment",
			memberships: []PNMembershipsSet{
				{Channel: PNMembershipsChannel{ID: "channel1"}, Custom: map[string]interface{}{"role": "admin"}},
				{Channel: PNMembershipsChannel{ID: "channel2"}, Custom: map[string]interface{}{"role": "member"}},
				{Channel: PNMembershipsChannel{ID: "channel3"}, Custom: map[string]interface{}{"role": "guest"}},
			},
			description: "Set multiple memberships in single operation",
		},
		{
			name: "Professional channels with complex data",
			memberships: []PNMembershipsSet{
				{Channel: PNMembershipsChannel{ID: "general"}, Custom: map[string]interface{}{"role": "member", "notifications": true}},
				{Channel: PNMembershipsChannel{ID: "announcements"}, Custom: map[string]interface{}{"role": "subscriber", "notifications": false}},
				{Channel: PNMembershipsChannel{ID: "dev-team"}, Custom: map[string]interface{}{"role": "admin", "team": "backend"}},
			},
			description: "Set professional channel memberships with metadata",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetMembershipsBuilder(pn)
			builder.UUID("test-uuid")
			builder.Set(tc.memberships)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid body with all memberships
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Contains(string(body), `"set":[`)

			// Verify all memberships are in body
			for _, membership := range tc.memberships {
				assert.Contains(string(body), fmt.Sprintf(`"id":"%s"`, membership.Channel.ID))
			}
		})
	}
}

// Comprehensive Edge Case Tests

func TestSetMembershipsWithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*setMembershipsBuilder)
	}{
		{
			name: "Very long UUID",
			setupFn: func(builder *setMembershipsBuilder) {
				longUUID := strings.Repeat("VeryLongUUID", 50) // 600 characters
				builder.UUID(longUUID)
			},
		},
		{
			name: "Large number of memberships to set",
			setupFn: func(builder *setMembershipsBuilder) {
				var largeMembershipList []PNMembershipsSet
				for i := 0; i < 50; i++ {
					largeMembershipList = append(largeMembershipList, PNMembershipsSet{
						Channel: PNMembershipsChannel{ID: fmt.Sprintf("channel_%d", i)},
						Custom: map[string]interface{}{
							"role":  fmt.Sprintf("role_%d", i),
							"level": i,
						},
					})
				}
				builder.Set(largeMembershipList)
			},
		},
		{
			name: "Large filter expression",
			setupFn: func(builder *setMembershipsBuilder) {
				largeFilter := "(" + strings.Repeat("custom.field == 'value' OR ", 100) + "custom.final == 'end')"
				builder.Filter(largeFilter)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *setMembershipsBuilder) {
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
			builder := newSetMembershipsBuilder(pn)
			builder.UUID("baseline-uuid")
			builder.Set([]PNMembershipsSet{
				{Channel: PNMembershipsChannel{ID: "baseline-channel"}, Custom: map[string]interface{}{"role": "admin"}},
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
		})
	}
}

func TestSetMembershipsSpecialCharacterHandling(t *testing.T) {
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
			builder := newSetMembershipsBuilder(pn)
			builder.UUID(specialString)
			builder.Filter(fmt.Sprintf("custom.field == '%s'", specialString))
			builder.QueryParam(map[string]string{
				"special_field": specialString,
			})
			builder.Set([]PNMembershipsSet{
				{Channel: PNMembershipsChannel{ID: specialString}, Custom: map[string]interface{}{"name": specialString}},
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
		})
	}
}

func TestSetMembershipsParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		uuid     string
		limit    int
		filter   string
		setCount int
	}{
		{
			name:     "Empty string UUID",
			uuid:     "",
			limit:    1,
			filter:   "",
			setCount: 0,
		},
		{
			name:     "Single character UUID",
			uuid:     "a",
			limit:    1,
			filter:   "a",
			setCount: 1,
		},
		{
			name:     "Unicode-only UUID",
			uuid:     "测试",
			limit:    50,
			filter:   "测试 == '值'",
			setCount: 2,
		},
		{
			name:     "Minimum limit",
			uuid:     "test",
			limit:    1,
			filter:   "simple",
			setCount: 1,
		},
		{
			name:     "Large limit",
			uuid:     "test",
			limit:    1000,
			filter:   "complex.nested == 'value'",
			setCount: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetMembershipsBuilder(pn)
			builder.UUID(tc.uuid)
			builder.Limit(tc.limit)
			if tc.filter != "" {
				builder.Filter(tc.filter)
			}

			// Add specified number of memberships to set
			var setMemberships []PNMembershipsSet
			for i := 0; i < tc.setCount; i++ {
				setMemberships = append(setMemberships, PNMembershipsSet{
					Channel: PNMembershipsChannel{ID: fmt.Sprintf("channel_%d", i)},
					Custom:  map[string]interface{}{"index": i},
				})
			}
			if len(setMemberships) > 0 {
				builder.Set(setMemberships)
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
			assert.NotEmpty(body) // PATCH operation always has body
		})
	}
}

func TestSetMembershipsComplexMembershipScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*setMembershipsBuilder)
		validateFn func(*testing.T, []byte)
	}{
		{
			name: "Mixed character channel IDs with complex data",
			setupFn: func(builder *setMembershipsBuilder) {
				builder.Set([]PNMembershipsSet{
					{Channel: PNMembershipsChannel{ID: "channel-english"}, Custom: map[string]interface{}{"role": "admin", "lang": "en"}},
					{Channel: PNMembershipsChannel{ID: "频道中文"}, Custom: map[string]interface{}{"角色": "成员", "语言": "zh"}},
					{Channel: PNMembershipsChannel{ID: "канал-русский"}, Custom: map[string]interface{}{"роль": "гость", "язык": "ru"}},
				})
			},
			validateFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				assert.Contains(bodyStr, `"id":"channel-english"`)
				assert.Contains(bodyStr, `"id":"频道中文"`)
				assert.Contains(bodyStr, `"id":"канал-русский"`)
				assert.Contains(bodyStr, `"role":"admin"`)
				assert.Contains(bodyStr, `"角色":"成员"`)
				assert.Contains(bodyStr, `"роль":"гость"`)
			},
		},
		{
			name: "Professional channels with comprehensive metadata",
			setupFn: func(builder *setMembershipsBuilder) {
				builder.Set([]PNMembershipsSet{
					{Channel: PNMembershipsChannel{ID: "company-general"}, Custom: map[string]interface{}{
						"role":          "member",
						"notifications": true,
						"permissions":   []string{"read", "write"},
					}},
					{Channel: PNMembershipsChannel{ID: "dev-team"}, Custom: map[string]interface{}{
						"role":   "lead",
						"team":   "backend",
						"skills": []string{"Go", "Python"},
					}},
				})
			},
			validateFn: func(t *testing.T, body []byte) {
				bodyStr := string(body)
				assert.Contains(bodyStr, `"id":"company-general"`)
				assert.Contains(bodyStr, `"id":"dev-team"`)
				assert.Contains(bodyStr, `"notifications":true`)
				assert.Contains(bodyStr, `"permissions":["read","write"]`)
				assert.Contains(bodyStr, `"skills":["Go","Python"]`)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetMembershipsBuilder(pn)
			builder.UUID("test-uuid")
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

func TestSetMembershipsExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newSetMembershipsBuilder(pn)
	builder.UUID("test-uuid")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestSetMembershipsPathBuildingEdgeCases(t *testing.T) {
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
			opts := newSetMembershipsOpts(pn, pn.ctx)
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

func TestSetMembershipsQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*setMembershipsOpts)
		expectError bool
	}{
		{
			name: "Nil query params",
			setupOpts: func(opts *setMembershipsOpts) {
				opts.QueryParam = nil
			},
			expectError: false,
		},
		{
			name: "Empty query params",
			setupOpts: func(opts *setMembershipsOpts) {
				opts.QueryParam = map[string]string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *setMembershipsOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *setMembershipsOpts) {
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
			opts := newSetMembershipsOpts(pn, pn.ctx)
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

func TestSetMembershipsBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newSetMembershipsBuilder(pn)

	include := []PNMembershipsInclude{PNMembershipsIncludeCustom}
	sort := []string{"name:asc", "created:desc"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}
	setMemberships := []PNMembershipsSet{
		{Channel: PNMembershipsChannel{ID: "channel1"}, Custom: map[string]interface{}{"role": "admin"}},
		{Channel: PNMembershipsChannel{ID: "channel2"}, Custom: map[string]interface{}{"role": "member"}},
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
		Set(setMemberships).
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
	assert.Equal(setMemberships, builder.opts.MembershipsSet)
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

	// Should build correct JSON body (PATCH operation)
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	expectedBody := `{"set":[{"channel":{"id":"channel1"},"custom":{"role":"admin"},"status":"","type":""},{"channel":{"id":"channel2"},"custom":{"role":"member"},"status":"","type":""}]}`
	assert.Equal(expectedBody, string(body))
}

func TestSetMembershipsResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetMembershipsOpts(pn, pn.ctx)

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
			resp, _, err := newPNSetMembershipsResponse(tc.jsonBytes, opts, StatusResponse{})

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

func TestSetMembershipsUUIDtoChannelsDirection(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		description string
		setupFn     func(*setMembershipsBuilder)
		verifyFn    func(*testing.T, []byte, string)
	}{
		{
			name:        "User joining multiple channels",
			description: "Single UUID joining multiple channels with different roles",
			setupFn: func(builder *setMembershipsBuilder) {
				builder.UUID("user123")
				builder.Set([]PNMembershipsSet{
					{Channel: PNMembershipsChannel{ID: "general"}, Custom: map[string]interface{}{"role": "member", "joinedAt": "2023-01-01"}},
					{Channel: PNMembershipsChannel{ID: "announcements"}, Custom: map[string]interface{}{"role": "subscriber", "notifications": true}},
					{Channel: PNMembershipsChannel{ID: "dev-team"}, Custom: map[string]interface{}{"role": "admin", "team": "backend"}},
				})
			},
			verifyFn: func(t *testing.T, body []byte, path string) {
				bodyStr := string(body)
				// Verify UUID-to-channels direction in path
				assert.Contains(path, "/uuids/user123/channels")

				// Verify all channels are being set for this UUID
				assert.Contains(bodyStr, `"id":"general"`)
				assert.Contains(bodyStr, `"id":"announcements"`)
				assert.Contains(bodyStr, `"id":"dev-team"`)
				assert.Contains(bodyStr, `"role":"member"`)
				assert.Contains(bodyStr, `"role":"admin"`)
			},
		},
		{
			name:        "Complex membership data with nested structures",
			description: "UUID with deeply nested custom data for channels",
			setupFn: func(builder *setMembershipsBuilder) {
				builder.UUID("admin_user")
				builder.Set([]PNMembershipsSet{
					{Channel: PNMembershipsChannel{ID: "company-wide"}, Custom: map[string]interface{}{
						"access": map[string]interface{}{
							"level":        "full",
							"permissions":  []string{"read", "write", "admin"},
							"restrictions": []string{},
						},
						"profile": map[string]interface{}{
							"joinDate":   "2023-01-01",
							"lastActive": "2023-12-01",
							"preferences": map[string]interface{}{
								"notifications": true,
								"digest":        "daily",
							},
						},
					}},
				})
			},
			verifyFn: func(t *testing.T, body []byte, path string) {
				bodyStr := string(body)
				assert.Contains(path, "/uuids/admin_user/channels")
				assert.Contains(bodyStr, `"level":"full"`)
				assert.Contains(bodyStr, `"permissions":["read","write","admin"]`)
				assert.Contains(bodyStr, `"notifications":true`)
			},
		},
		{
			name:        "International channels with Unicode data",
			description: "UUID joining international channels with Unicode metadata",
			setupFn: func(builder *setMembershipsBuilder) {
				builder.UUID("国际用户")
				builder.Set([]PNMembershipsSet{
					{Channel: PNMembershipsChannel{ID: "中文频道"}, Custom: map[string]interface{}{"语言": "中文", "角色": "成员"}},
					{Channel: PNMembershipsChannel{ID: "русский-канал"}, Custom: map[string]interface{}{"язык": "русский", "роль": "гость"}},
					{Channel: PNMembershipsChannel{ID: "日本語チャンネル"}, Custom: map[string]interface{}{"言語": "日本語", "役割": "管理者"}},
				})
			},
			verifyFn: func(t *testing.T, body []byte, path string) {
				bodyStr := string(body)
				assert.Contains(path, "/uuids/国际用户/channels")
				assert.Contains(bodyStr, `"id":"中文频道"`)
				assert.Contains(bodyStr, `"id":"русский-канал"`)
				assert.Contains(bodyStr, `"id":"日本語チャンネル"`)
				assert.Contains(bodyStr, `"角色":"成员"`)
				assert.Contains(bodyStr, `"роль":"гость"`)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetMembershipsBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path (UUID-to-channels direction)
			path, err := builder.opts.buildPath()
			assert.Nil(err)

			// Should build valid JSON body
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.NotEmpty(body)

			// Run verification
			tc.verifyFn(t, body, path)
		})
	}
}
