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
		o = newGetChannelMembersBuilderV2WithContext(pn, backgroundContext)
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
	opts := &getChannelMembersOptsV2{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetChannelMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestGetChannelMembersV2ResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getChannelMembersOptsV2{
		pubnub: pn,
	}
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
