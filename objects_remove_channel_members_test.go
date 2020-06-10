package pubnub

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/pubnub/go/utils"
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
		o = newRemoveChannelMembersBuilderWithContext(pn, backgroundContext)
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
	opts := &removeChannelMembersOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNRemoveChannelMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRemoveChannelMembersResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &removeChannelMembersOpts{
		pubnub: pn,
	}
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
