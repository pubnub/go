package pubnub

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"

	h "github.com/pubnub/go/v5/tests/helpers"
	"github.com/pubnub/go/v5/utils"
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
		o = newRemoveMembershipsBuilderWithContext(pn, backgroundContext)
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
	opts := &removeMembershipsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNRemoveMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRemoveMembershipsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &removeMembershipsOpts{
		pubnub: pn,
	}
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
