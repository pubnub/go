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
		o = newGetAllUUIDMetadataBuilderWithContext(pn, backgroundContext)
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
	opts := &getAllUUIDMetadataOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetAllUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestGetAllUUIDMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getAllUUIDMetadataOpts{
		pubnub: pn,
	}
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
