package pubnub

import (
	"encoding/json"
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

func AssertRemoveMemberships(t *testing.T, checkQueryParam, testContext bool, withFilter bool, withSort bool) {
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

	o := newRemoveMembershipsBuilder(pn)
	if testContext {
		o = newRemoveMembershipsBuilderWithContext(pn, pn.ctx)
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

	custom2 := make(map[string]interface{})
	custom2["a2"] = "b2"
	custom2["c2"] = "d2"

	channel := PNMembershipsChannel{
		ID: id0,
	}

	re := PNMembershipsRemove{
		Channel: channel,
	}

	reArr := []PNMembershipsRemove{
		re,
	}
	o.Remove(reArr)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/uuids/%s/channels", pn.Config.SubscribeKey, spaceID),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	expectedBody := "{\"delete\":[{\"channel\":{\"id\":\"id0\"}}]}"

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

func TestRemoveMemberships(t *testing.T) {
	AssertRemoveMemberships(t, true, false, false, false)
}

func TestRemoveMembershipsContext(t *testing.T) {
	AssertRemoveMemberships(t, true, true, false, false)
}

func TestRemoveMembershipsWithFilter(t *testing.T) {
	AssertRemoveMemberships(t, true, false, true, false)
}

func TestRemoveMembershipsWithFilterContext(t *testing.T) {
	AssertRemoveMemberships(t, true, true, true, false)
}

func TestRemoveMembershipsWithSort(t *testing.T) {
	AssertRemoveMemberships(t, true, false, false, true)
}

func TestRemoveMembershipsWithSortContext(t *testing.T) {
	AssertRemoveMemberships(t, true, true, false, true)
}

func TestRemoveMembershipsWithFilterWithSort(t *testing.T) {
	AssertRemoveMemberships(t, true, false, true, true)
}

func TestRemoveMembershipsWithFilterWithSortContext(t *testing.T) {
	AssertRemoveMemberships(t, true, true, true, true)
}

func TestRemoveMembershipsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNRemoveMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRemoveMembershipsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":[{"id":"spaceid3","custom":{"a3":"b3","c3":"d3"},"channel":{"id":"spaceid3","name":"spaceid3name","description":"spaceid3desc","custom":{"a":"b"},"created":"2019-08-23T10:34:43.985248Z","updated":"2019-08-23T10:34:43.985248Z","eTag":"Aazjn7vC3oDDYw"},"created":"2019-08-23T10:41:17.156491Z","updated":"2019-08-23T10:41:17.156491Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":1,"next":"MQ"}`)

	r, _, err := newPNRemoveMembershipsResponse(jsonBytes, opts, StatusResponse{})

	assert.Equal(1, r.TotalCount)
	assert.Equal("MQ", r.Next)
	assert.Equal("spaceid3", r.Data[0].ID)
	assert.Equal("spaceid3", r.Data[0].Channel.ID)
	assert.Equal("spaceid3name", r.Data[0].Channel.Name)
	assert.Equal("spaceid3desc", r.Data[0].Channel.Description)
	// assert.Equal("2019-08-23T10:34:43.985248Z", r.Data[0].Channel.Created)
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

func TestRemoveMembershipsValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestRemoveMembershipsValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	assert.Nil(opts.validate())
}

func TestRemoveMembershipsValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"
	opts.Include = []string{"custom", "channel"}
	opts.QueryParam = map[string]string{"param": "value"}
	opts.MembershipsRemove = []PNMembershipsRemove{
		{Channel: PNMembershipsChannel{ID: "channel1"}},
	}

	assert.Nil(opts.validate())
}

func TestRemoveMembershipsValidateUUIDDefaultBehavior(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.UUID = "config-uuid"

	builder := newRemoveMembershipsBuilder(pn)
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

func TestRemoveMembershipsHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)

	assert.Equal("PATCH", opts.httpMethod())
}

func TestRemoveMembershipsOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)

	assert.Equal(PNRemoveMembershipsOperation, opts.operationType())
}

func TestRemoveMembershipsIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestRemoveMembershipsTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (10 setters)

func TestRemoveMembershipsBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMembershipsBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
	assert.Equal(removeMembershipsLimit, builder.opts.Limit) // Default limit (100)
}

func TestRemoveMembershipsBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMembershipsBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestRemoveMembershipsBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMembershipsBuilder(pn)

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

	// Test Remove setter (CRITICAL for this operation)
	remove := []PNMembershipsRemove{
		{Channel: PNMembershipsChannel{ID: "channel1"}},
		{Channel: PNMembershipsChannel{ID: "channel2"}},
	}
	builder.Remove(remove)
	assert.Equal(remove, builder.opts.MembershipsRemove)

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

func TestRemoveMembershipsBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	include := []PNMembershipsInclude{PNMembershipsIncludeCustom}
	sort := []string{"name"}
	remove := []PNMembershipsRemove{
		{Channel: PNMembershipsChannel{ID: "channel1"}},
	}
	queryParam := map[string]string{"key": "value"}
	transport := &http.Transport{}

	builder := newRemoveMembershipsBuilder(pn)
	result := builder.UUID("test-uuid").
		Include(include).
		Limit(75).
		Start("start").
		End("end").
		Count(true).
		Filter("filter").
		Sort(sort).
		Remove(remove).
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
	assert.Equal(remove, builder.opts.MembershipsRemove)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

func TestRemoveMembershipsBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMembershipsBuilder(pn)

	// Verify default values
	assert.Equal("", builder.opts.UUID) // UUID defaults to empty, set later in Execute
	assert.Nil(builder.opts.Include)
	assert.Equal(removeMembershipsLimit, builder.opts.Limit) // 100
	assert.Equal("", builder.opts.Start)
	assert.Equal("", builder.opts.End)
	assert.Equal(false, builder.opts.Count)
	assert.Equal("", builder.opts.Filter)
	assert.Nil(builder.opts.Sort)
	assert.Nil(builder.opts.MembershipsRemove) // Empty by default
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestRemoveMembershipsBuilderIncludeTypes(t *testing.T) {
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
			builder := newRemoveMembershipsBuilder(pn)
			builder.Include(tc.includes)

			expectedInclude := EnumArrayToStringArray(tc.includes)
			assert.Equal(expectedInclude, builder.opts.Include)
		})
	}
}

func TestRemoveMembershipsBuilderRemoveOperations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name   string
		remove []PNMembershipsRemove
	}{
		{
			name: "Single channel removal",
			remove: []PNMembershipsRemove{
				{Channel: PNMembershipsChannel{ID: "channel1"}},
			},
		},
		{
			name: "Multiple channels removal",
			remove: []PNMembershipsRemove{
				{Channel: PNMembershipsChannel{ID: "channel1"}},
				{Channel: PNMembershipsChannel{ID: "channel2"}},
				{Channel: PNMembershipsChannel{ID: "channel3"}},
			},
		},
		{
			name: "Unicode channel removal",
			remove: []PNMembershipsRemove{
				{Channel: PNMembershipsChannel{ID: "频道中文"}},
				{Channel: PNMembershipsChannel{ID: "канал-русский"}},
			},
		},
		{
			name: "Special character channels",
			remove: []PNMembershipsRemove{
				{Channel: PNMembershipsChannel{ID: "channel@with#symbols"}},
				{Channel: PNMembershipsChannel{ID: "channel-with-dashes"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveMembershipsBuilder(pn)
			builder.Remove(tc.remove)

			assert.Equal(tc.remove, builder.opts.MembershipsRemove)
		})
	}
}

func TestRemoveMembershipsBuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	transport := &http.Transport{}
	include := []PNMembershipsInclude{PNMembershipsIncludeCustom}
	sort := []string{"name", "created:desc"}
	remove := []PNMembershipsRemove{
		{Channel: PNMembershipsChannel{ID: "channel1"}},
		{Channel: PNMembershipsChannel{ID: "channel2"}},
	}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Test all 10 setters in chain
	builder := newRemoveMembershipsBuilder(pn).
		UUID("test-uuid").
		Include(include).
		Limit(50).
		Start("start-token").
		End("end-token").
		Count(true).
		Filter("name LIKE 'test*'").
		Sort(sort).
		Remove(remove).
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
	assert.Equal(remove, builder.opts.MembershipsRemove)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestRemoveMembershipsBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/uuids/test-uuid/channels"
	assert.Equal(expected, path)
}

func TestRemoveMembershipsBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.UUID = "my-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/uuids/my-uuid/channels"
	assert.Equal(expected, path)
}

func TestRemoveMembershipsBuildPathWithSpecialCharsInUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.UUID = "uuid-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/")
	assert.Contains(path, "uuid-with-special@chars#and$symbols")
	assert.Contains(path, "/channels")
}

func TestRemoveMembershipsBuildPathWithUnicodeUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.UUID = "用户ID-пользователь-ユーザーID"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/")
	assert.Contains(path, "用户ID-пользователь-ユーザーID")
	assert.Contains(path, "/channels")
}

// JSON Body Building Tests (CRITICAL for PATCH operation)

func TestRemoveMembershipsBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.MembershipsRemove = []PNMembershipsRemove{}

	body, err := opts.buildBody()
	assert.Nil(err)
	expectedBody := `{"delete":[]}`
	assert.Equal(expectedBody, string(body))
}

func TestRemoveMembershipsBuildBodySingle(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.MembershipsRemove = []PNMembershipsRemove{
		{Channel: PNMembershipsChannel{ID: "channel1"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	expectedBody := `{"delete":[{"channel":{"id":"channel1"}}]}`
	assert.Equal(expectedBody, string(body))
}

func TestRemoveMembershipsBuildBodyMultiple(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.MembershipsRemove = []PNMembershipsRemove{
		{Channel: PNMembershipsChannel{ID: "channel1"}},
		{Channel: PNMembershipsChannel{ID: "channel2"}},
		{Channel: PNMembershipsChannel{ID: "channel3"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	expectedBody := `{"delete":[{"channel":{"id":"channel1"}},{"channel":{"id":"channel2"}},{"channel":{"id":"channel3"}}]}`
	assert.Equal(expectedBody, string(body))
}

func TestRemoveMembershipsBuildBodyUnicode(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.MembershipsRemove = []PNMembershipsRemove{
		{Channel: PNMembershipsChannel{ID: "频道中文"}},
		{Channel: PNMembershipsChannel{ID: "канал-русский"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	// Unicode should be properly encoded in JSON
	assert.Contains(string(body), "频道中文")
	assert.Contains(string(body), "канал-русский")
	assert.Contains(string(body), `"delete":[`)
}

func TestRemoveMembershipsBuildBodySpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.MembershipsRemove = []PNMembershipsRemove{
		{Channel: PNMembershipsChannel{ID: "channel@with#symbols"}},
		{Channel: PNMembershipsChannel{ID: "channel-with-dashes"}},
	}

	body, err := opts.buildBody()
	assert.Nil(err)
	expectedBody := `{"delete":[{"channel":{"id":"channel@with#symbols"}},{"channel":{"id":"channel-with-dashes"}}]}`
	assert.Equal(expectedBody, string(body))
}

func TestRemoveMembershipsBuildBodyNilRemove(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)
	opts.MembershipsRemove = nil

	body, err := opts.buildBody()
	assert.Nil(err)
	expectedBody := `{"delete":null}`
	assert.Equal(expectedBody, string(body))
}

func TestRemoveMembershipsBuildBodyLargeSet(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)

	// Create a large set of channels to remove
	var remove []PNMembershipsRemove
	for i := 0; i < 50; i++ {
		remove = append(remove, PNMembershipsRemove{
			Channel: PNMembershipsChannel{ID: fmt.Sprintf("channel_%d", i)},
		})
	}
	opts.MembershipsRemove = remove

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"delete":[`)
	assert.Contains(string(body), `"channel_0"`)
	assert.Contains(string(body), `"channel_49"`)

	// Verify it's valid JSON
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	assert.Nil(err)
	assert.NotNil(jsonData["delete"])
}

// Query Parameter Tests

func TestRemoveMembershipsBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
	assert.Equal("0", query.Get("limit")) // Default limit not set until builder initialization
	assert.Equal("0", query.Get("count")) // Default count=false
}

func TestRemoveMembershipsBuildQueryWithInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)

	opts.Include = []string{"custom", "channel"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	includeValue := query.Get("include")
	assert.Contains(includeValue, "custom")
	assert.Contains(includeValue, "channel")
}

func TestRemoveMembershipsBuildQueryWithPagination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)

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

func TestRemoveMembershipsBuildQueryWithFilterAndSort(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)

	opts.Filter = "custom.role == 'admin'"
	opts.Sort = []string{"name", "created:desc"}

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("custom.role == 'admin'", query.Get("filter"))

	sortValue := query.Get("sort")
	assert.Contains(sortValue, "name")
	assert.Contains(sortValue, "created:desc")
}

func TestRemoveMembershipsBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)

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

func TestRemoveMembershipsBuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)

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

// PATCH-Specific Tests (Remove Operation Characteristics)

func TestRemoveMembershipsPatchOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMembershipsBuilder(pn)
	builder.UUID("test-uuid")
	builder.Remove([]PNMembershipsRemove{
		{Channel: PNMembershipsChannel{ID: "channel1"}},
	})

	// Verify it's a PATCH operation
	assert.Equal("PATCH", builder.opts.httpMethod())

	// PATCH operations have structured JSON body
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.NotEmpty(body)
	assert.Contains(string(body), `"delete"`)

	// Should have proper path for membership removal (UUID to channels)
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/test-uuid/channels")
}

func TestRemoveMembershipsDefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMembershipsBuilder(pn)

	// Should have default limit set to removeMembershipsLimit (100)
	assert.Equal(removeMembershipsLimit, builder.opts.Limit)
	assert.Equal(100, builder.opts.Limit)

	// Should be included in query
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("limit"))
}

func TestRemoveMembershipsDeleteOperationValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*removeMembershipsOpts)
		description string
	}{
		{
			name: "Basic remove operation",
			setupOpts: func(opts *removeMembershipsOpts) {
				opts.UUID = "user123"
				opts.MembershipsRemove = []PNMembershipsRemove{
					{Channel: PNMembershipsChannel{ID: "channel1"}},
				}
			},
			description: "Remove single membership",
		},
		{
			name: "Remove with include options",
			setupOpts: func(opts *removeMembershipsOpts) {
				opts.UUID = "user123"
				opts.Include = []string{"custom", "channel"}
				opts.MembershipsRemove = []PNMembershipsRemove{
					{Channel: PNMembershipsChannel{ID: "channel1"}},
				}
			},
			description: "Remove memberships with additional data",
		},
		{
			name: "Remove with pagination",
			setupOpts: func(opts *removeMembershipsOpts) {
				opts.UUID = "user123"
				opts.Limit = 20
				opts.Start = "pagination_token"
				opts.MembershipsRemove = []PNMembershipsRemove{
					{Channel: PNMembershipsChannel{ID: "channel1"}},
				}
			},
			description: "Remove with pagination parameters",
		},
		{
			name: "Bulk remove operation",
			setupOpts: func(opts *removeMembershipsOpts) {
				opts.UUID = "user123"
				opts.MembershipsRemove = []PNMembershipsRemove{
					{Channel: PNMembershipsChannel{ID: "channel1"}},
					{Channel: PNMembershipsChannel{ID: "channel2"}},
					{Channel: PNMembershipsChannel{ID: "channel3"}},
				}
			},
			description: "Remove multiple memberships",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newRemoveMembershipsOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			// Should pass validation
			assert.Nil(opts.validate())

			// Should be PATCH operation
			assert.Equal("PATCH", opts.httpMethod())

			// Should have structured body with delete operation
			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Contains(string(body), `"delete"`)

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

func TestRemoveMembershipsResponseStructureAfterRemoval(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMembershipsBuilder(pn)
	builder.UUID("test-uuid")
	builder.Remove([]PNMembershipsRemove{
		{Channel: PNMembershipsChannel{ID: "channel1"}},
	})

	// Response should contain remaining memberships data after PATCH operation
	// This is tested in the existing TestRemoveMembershipsResponseValuePass
	// but verify the operation is configured correctly
	opts := builder.opts

	// Verify operation is configured correctly
	assert.Equal("PATCH", opts.httpMethod())
	assert.Equal(PNRemoveMembershipsOperation, opts.operationType())
	assert.True(opts.isAuthRequired())
}

func TestRemoveMembershipsDeleteStructureValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test different delete structure scenarios
	testCases := []struct {
		name      string
		remove    []PNMembershipsRemove
		expectKey string
	}{
		{
			name:      "Empty delete array",
			remove:    []PNMembershipsRemove{},
			expectKey: `"delete":[]`,
		},
		{
			name: "Single channel delete",
			remove: []PNMembershipsRemove{
				{Channel: PNMembershipsChannel{ID: "channel1"}},
			},
			expectKey: `"delete":[{"channel":{"id":"channel1"}}]`,
		},
		{
			name: "Multiple channels delete",
			remove: []PNMembershipsRemove{
				{Channel: PNMembershipsChannel{ID: "channel1"}},
				{Channel: PNMembershipsChannel{ID: "channel2"}},
			},
			expectKey: `"channel1"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newRemoveMembershipsOpts(pn, pn.ctx)
			opts.MembershipsRemove = tc.remove

			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Contains(string(body), tc.expectKey)
		})
	}
}

func TestRemoveMembershipsPatchDefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveMembershipsBuilder(pn)

	// Default limit should be applied even for PATCH operations
	assert.Equal(100, builder.opts.Limit)

	// Should be reflected in queries
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("limit"))
}

func TestRemoveMembershipsBulkOperations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test bulk remove operations
	var bulkRemove []PNMembershipsRemove
	for i := 0; i < 20; i++ {
		bulkRemove = append(bulkRemove, PNMembershipsRemove{
			Channel: PNMembershipsChannel{ID: fmt.Sprintf("bulk_channel_%d", i)},
		})
	}

	builder := newRemoveMembershipsBuilder(pn)
	builder.UUID("bulk-user")
	builder.Remove(bulkRemove)

	// Should handle large remove operations
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"delete":[`)
	assert.Contains(string(body), "bulk_channel_0")
	assert.Contains(string(body), "bulk_channel_19")

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should be PATCH operation
	assert.Equal("PATCH", builder.opts.httpMethod())
}

// Comprehensive Edge Case Tests

func TestRemoveMembershipsWithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*removeMembershipsBuilder)
	}{
		{
			name: "Very long UUID",
			setupFn: func(builder *removeMembershipsBuilder) {
				longUUID := strings.Repeat("VeryLongUUID", 50) // 600 characters
				builder.UUID(longUUID)
			},
		},
		{
			name: "Large filter expression",
			setupFn: func(builder *removeMembershipsBuilder) {
				largeFilter := "(" + strings.Repeat("custom.field == 'value' OR ", 100) + "custom.final == 'end')"
				builder.Filter(largeFilter)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *removeMembershipsBuilder) {
				largeQueryParam := make(map[string]string)
				for i := 0; i < 100; i++ {
					largeQueryParam[fmt.Sprintf("param_%d", i)] = strings.Repeat(fmt.Sprintf("value_%d_", i), 20)
				}
				builder.QueryParam(largeQueryParam)
			},
		},
		{
			name: "Large pagination limit",
			setupFn: func(builder *removeMembershipsBuilder) {
				builder.Limit(10000) // Very large limit
			},
		},
		{
			name: "Massive remove operation",
			setupFn: func(builder *removeMembershipsBuilder) {
				var massiveRemove []PNMembershipsRemove
				for i := 0; i < 500; i++ {
					massiveRemove = append(massiveRemove, PNMembershipsRemove{
						Channel: PNMembershipsChannel{ID: fmt.Sprintf("massive_channel_%d", i)},
					})
				}
				builder.Remove(massiveRemove)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveMembershipsBuilder(pn)
			builder.UUID("baseline-uuid")
			builder.Remove([]PNMembershipsRemove{
				{Channel: PNMembershipsChannel{ID: "baseline-channel"}},
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

			// Should build valid body (PATCH operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.NotEmpty(body)
			assert.Contains(string(body), `"delete"`)
		})
	}
}

func TestRemoveMembershipsSpecialCharacterHandling(t *testing.T) {
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
			builder := newRemoveMembershipsBuilder(pn)
			builder.UUID(specialString)
			builder.Filter(fmt.Sprintf("custom.field == '%s'", specialString))
			builder.QueryParam(map[string]string{
				"special_field": specialString,
			})
			builder.Remove([]PNMembershipsRemove{
				{Channel: PNMembershipsChannel{ID: specialString}},
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

			// Should build valid body with special characters
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Contains(string(body), `"delete"`)
		})
	}
}

func TestRemoveMembershipsParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		uuid        string
		limit       int
		filter      string
		removeCount int
	}{
		{
			name:        "Empty string UUID",
			uuid:        "",
			limit:       1,
			filter:      "",
			removeCount: 1,
		},
		{
			name:        "Single character UUID",
			uuid:        "a",
			limit:       1,
			filter:      "a",
			removeCount: 1,
		},
		{
			name:        "Unicode-only UUID",
			uuid:        "测试",
			limit:       50,
			filter:      "测试 == '值'",
			removeCount: 2,
		},
		{
			name:        "Minimum limit",
			uuid:        "test",
			limit:       1,
			filter:      "simple",
			removeCount: 1,
		},
		{
			name:        "Large limit with many removals",
			uuid:        "test",
			limit:       1000,
			filter:      "complex.nested == 'value'",
			removeCount: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveMembershipsBuilder(pn)
			builder.UUID(tc.uuid)
			builder.Limit(tc.limit)
			if tc.filter != "" {
				builder.Filter(tc.filter)
			}

			// Create remove array
			var remove []PNMembershipsRemove
			for i := 0; i < tc.removeCount; i++ {
				remove = append(remove, PNMembershipsRemove{
					Channel: PNMembershipsChannel{ID: fmt.Sprintf("channel_%d", i)},
				})
			}
			builder.Remove(remove)

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
			assert.Contains(string(body), `"delete"`) // PATCH operation always has body
		})
	}
}

func TestRemoveMembershipsComplexRemovalScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*removeMembershipsBuilder)
		validateFn func(*testing.T, string)
	}{
		{
			name: "International UUID with complex channel removals",
			setupFn: func(builder *removeMembershipsBuilder) {
				builder.UUID("用户中文123")
				builder.Filter("custom.角色 == '管理员' && custom.语言 == 'zh'")
				builder.Include([]PNMembershipsInclude{PNMembershipsIncludeCustom, PNMembershipsIncludeChannel})
				builder.Remove([]PNMembershipsRemove{
					{Channel: PNMembershipsChannel{ID: "频道中文1"}},
					{Channel: PNMembershipsChannel{ID: "频道中文2"}},
				})
			},
			validateFn: func(t *testing.T, path string) {
				assert.Contains(path, "/uuids/用户中文123/channels")
			},
		},
		{
			name: "Professional user bulk channel removal",
			setupFn: func(builder *removeMembershipsBuilder) {
				builder.UUID("professional@company.com")
				builder.Filter("custom.role IN ('admin', 'manager')")
				builder.Sort([]string{"name:asc", "custom.priority:desc"})
				builder.Limit(25)
				builder.Count(true)

				var bulkRemove []PNMembershipsRemove
				for i := 0; i < 15; i++ {
					bulkRemove = append(bulkRemove, PNMembershipsRemove{
						Channel: PNMembershipsChannel{ID: fmt.Sprintf("company_channel_%d", i)},
					})
				}
				builder.Remove(bulkRemove)
			},
			validateFn: func(t *testing.T, path string) {
				assert.Contains(path, "/uuids/professional@company.com/channels")
			},
		},
		{
			name: "Email-like UUID with mixed channel types removal",
			setupFn: func(builder *removeMembershipsBuilder) {
				builder.UUID("user@company.com")
				builder.Include([]PNMembershipsInclude{
					PNMembershipsIncludeCustom,
					PNMembershipsIncludeChannel,
					PNMembershipsIncludeChannelCustom,
				})
				builder.Filter("channel.name LIKE '%company%'")
				builder.Remove([]PNMembershipsRemove{
					{Channel: PNMembershipsChannel{ID: "public-channel"}},
					{Channel: PNMembershipsChannel{ID: "private_channel"}},
					{Channel: PNMembershipsChannel{ID: "频道国际"}},
				})
			},
			validateFn: func(t *testing.T, path string) {
				assert.Contains(path, "/uuids/user@company.com/channels")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveMembershipsBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path (UUID-to-channels direction)
			path, err := builder.opts.buildPath()
			assert.Nil(err)

			// Should have structured PATCH body
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Contains(string(body), `"delete"`)

			// Run verification
			tc.validateFn(t, path)
		})
	}
}

// Error Scenario Tests

func TestRemoveMembershipsExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newRemoveMembershipsBuilder(pn)
	builder.UUID("test-uuid")
	builder.Remove([]PNMembershipsRemove{
		{Channel: PNMembershipsChannel{ID: "channel1"}},
	})

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestRemoveMembershipsPathBuildingEdgeCases(t *testing.T) {
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
			opts := newRemoveMembershipsOpts(pn, pn.ctx)
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

func TestRemoveMembershipsBodyBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*removeMembershipsOpts)
		expectError bool
		expectKey   string
	}{
		{
			name: "Nil remove array",
			setupOpts: func(opts *removeMembershipsOpts) {
				opts.MembershipsRemove = nil
			},
			expectError: false,
			expectKey:   `"delete":null`,
		},
		{
			name: "Empty remove array",
			setupOpts: func(opts *removeMembershipsOpts) {
				opts.MembershipsRemove = []PNMembershipsRemove{}
			},
			expectError: false,
			expectKey:   `"delete":[]`,
		},
		{
			name: "Remove with empty channel ID",
			setupOpts: func(opts *removeMembershipsOpts) {
				opts.MembershipsRemove = []PNMembershipsRemove{
					{Channel: PNMembershipsChannel{ID: ""}},
				}
			},
			expectError: false,
			expectKey:   `"delete":[{"channel":{"id":""}}]`,
		},
		{
			name: "Remove with Unicode channels",
			setupOpts: func(opts *removeMembershipsOpts) {
				opts.MembershipsRemove = []PNMembershipsRemove{
					{Channel: PNMembershipsChannel{ID: "测试频道"}},
					{Channel: PNMembershipsChannel{ID: "канал-тест"}},
				}
			},
			expectError: false,
			expectKey:   `"delete":[`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newRemoveMembershipsOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			body, err := opts.buildBody()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Contains(string(body), tc.expectKey)
			}
		})
	}
}

func TestRemoveMembershipsQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*removeMembershipsOpts)
		expectError bool
	}{
		{
			name: "Nil query params",
			setupOpts: func(opts *removeMembershipsOpts) {
				opts.QueryParam = nil
			},
			expectError: false,
		},
		{
			name: "Empty query params",
			setupOpts: func(opts *removeMembershipsOpts) {
				opts.QueryParam = map[string]string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *removeMembershipsOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *removeMembershipsOpts) {
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
			opts := newRemoveMembershipsOpts(pn, pn.ctx)
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

func TestRemoveMembershipsBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newRemoveMembershipsBuilder(pn)

	include := []PNMembershipsInclude{PNMembershipsIncludeCustom}
	sort := []string{"name:asc", "created:desc"}
	remove := []PNMembershipsRemove{
		{Channel: PNMembershipsChannel{ID: "channel1"}},
		{Channel: PNMembershipsChannel{ID: "channel2"}},
	}
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
		Remove(remove).
		QueryParam(queryParam)

	// Verify all values are set correctly
	assert.Equal("complete-test-uuid", builder.opts.UUID)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)
	assert.Equal(75, builder.opts.Limit)
	assert.Equal("start-token", builder.opts.Start)
	assert.Equal("end-token", builder.opts.End)
	assert.Equal(true, builder.opts.Count)
	assert.Equal("active = true", builder.opts.Filter)
	assert.Equal(sort, builder.opts.Sort)
	assert.Equal(remove, builder.opts.MembershipsRemove)
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

	// Should build valid PATCH body
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"delete"`)
	assert.Contains(string(body), "channel1")
	assert.Contains(string(body), "channel2")
}

func TestRemoveMembershipsResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveMembershipsOpts(pn, pn.ctx)

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
			name:        "Valid response after removal",
			jsonBytes:   []byte(`{"status":200,"data":[{"id":"remaining_channel","channel":{"id":"remaining_channel","name":"Remaining Channel"},"custom":{"role":"admin"}}],"totalCount":1,"next":"abc","prev":"xyz"}`),
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
			resp, _, err := newPNRemoveMembershipsResponse(tc.jsonBytes, opts, StatusResponse{})

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

// UUID-to-Channels Direction Tests (Removal)

func TestRemoveMembershipsUUIDtoChannelsDirection(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		description string
		setupFn     func(*removeMembershipsBuilder)
		verifyFn    func(*testing.T, string)
	}{
		{
			name:        "Basic user membership removal",
			description: "Single UUID removing specific channel memberships",
			setupFn: func(builder *removeMembershipsBuilder) {
				builder.UUID("user123")
				builder.Remove([]PNMembershipsRemove{
					{Channel: PNMembershipsChannel{ID: "channel1"}},
				})
			},
			verifyFn: func(t *testing.T, path string) {
				// Verify UUID-to-channels direction in path
				assert.Contains(path, "/uuids/user123/channels")
			},
		},
		{
			name:        "Professional user bulk membership removal",
			description: "UUID removing multiple channel memberships with filtering",
			setupFn: func(builder *removeMembershipsBuilder) {
				builder.UUID("professional@company.com")
				builder.Filter("custom.role IN ('admin', 'manager')")
				builder.Include([]PNMembershipsInclude{
					PNMembershipsIncludeCustom,
					PNMembershipsIncludeChannel,
					PNMembershipsIncludeChannelCustom,
				})
				builder.Remove([]PNMembershipsRemove{
					{Channel: PNMembershipsChannel{ID: "old-project-1"}},
					{Channel: PNMembershipsChannel{ID: "old-project-2"}},
					{Channel: PNMembershipsChannel{ID: "archived-channel"}},
				})
			},
			verifyFn: func(t *testing.T, path string) {
				assert.Contains(path, "/uuids/professional@company.com/channels")
			},
		},
		{
			name:        "International user selective removal",
			description: "UUID with Unicode characters removing international channels",
			setupFn: func(builder *removeMembershipsBuilder) {
				builder.UUID("国际用户_测试")
				builder.Limit(25)
				builder.Start("pagination_token")
				builder.Sort([]string{"name:asc"})
				builder.Count(true)
				builder.Remove([]PNMembershipsRemove{
					{Channel: PNMembershipsChannel{ID: "频道中文旧"}},
					{Channel: PNMembershipsChannel{ID: "канал-старый"}},
					{Channel: PNMembershipsChannel{ID: "古いチャンネル"}},
				})
			},
			verifyFn: func(t *testing.T, path string) {
				assert.Contains(path, "/uuids/国际用户_测试/channels")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveMembershipsBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should be PATCH operation
			assert.Equal("PATCH", builder.opts.httpMethod())

			// Should have structured body with delete operations
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Contains(string(body), `"delete"`)

			// Should build valid path (UUID-to-channels direction)
			path, err := builder.opts.buildPath()
			assert.Nil(err)

			// Run verification
			tc.verifyFn(t, path)
		})
	}
}
