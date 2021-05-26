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

func AssertGetMembershipsV2(t *testing.T, checkQueryParam, testContext, withFilter bool, withSort bool) {
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

	o := newGetMembershipsBuilderV2(pn)
	if testContext {
		o = newGetMembershipsBuilderV2WithContext(pn, backgroundContext)
	}

	userID := "id0"
	limit := 90
	start := "Mxmy"
	end := "Nxny"

	o.UUID(userID)
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
		fmt.Sprintf("/v2/objects/%s/uuids/%s/channels", pn.Config.SubscribeKey, "id0"),
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

func TestGetMembershipsV2(t *testing.T) {
	AssertGetMembershipsV2(t, true, false, false, false)
}

func TestGetMembershipsV2Context(t *testing.T) {
	AssertGetMembershipsV2(t, true, true, false, false)
}

func TestGetMembershipsV2WithFilter(t *testing.T) {
	AssertGetMembershipsV2(t, true, false, true, false)
}

func TestGetMembershipsV2WithFilterContext(t *testing.T) {
	AssertGetMembershipsV2(t, true, true, true, false)
}

func TestGetMembershipsV2WithSort(t *testing.T) {
	AssertGetMembershipsV2(t, true, false, false, true)
}

func TestGetMembershipsV2WithSortContext(t *testing.T) {
	AssertGetMembershipsV2(t, true, true, false, true)
}

func TestGetMembershipsV2WithFilterWithSort(t *testing.T) {
	AssertGetMembershipsV2(t, true, false, true, true)
}

func TestGetMembershipsV2WithFilterWithSortContext(t *testing.T) {
	AssertGetMembershipsV2(t, true, true, true, true)
}

func TestGetMembershipsV2ResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getMembershipsOptsV2{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestGetMembershipsV2ResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getMembershipsOptsV2{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":[{"id":"id0","custom":{"a3":"b3","c3":"d3"},"channel":{"id":"id0","name":"name","description":"desc","custom":{"a":"b"},"created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T13:26:08.341297Z","eTag":"Aee9zsKNndXlHw"},"created":"2019-08-20T13:26:24.07832Z","updated":"2019-08-20T13:26:24.07832Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":1,"next":"MQ", "prev":"NQ"}`)

	r, _, err := newPNGetMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(1, r.TotalCount)
	assert.Equal("MQ", r.Next)
	assert.Equal("NQ", r.Prev)
	assert.Equal("id0", r.Data[0].ID)
	assert.Equal("name", r.Data[0].Channel.Name)
	assert.Equal("desc", r.Data[0].Channel.Description)
	// assert.Equal("2019-08-20T13:26:08.341297Z", r.Data[0].Channel.Created)
	assert.Equal("2019-08-20T13:26:08.341297Z", r.Data[0].Channel.Updated)
	assert.Equal("Aee9zsKNndXlHw", r.Data[0].Channel.ETag)
	assert.Equal("b", r.Data[0].Channel.Custom["a"])
	assert.Equal(nil, r.Data[0].Channel.Custom["c"])
	assert.Equal("2019-08-20T13:26:24.07832Z", r.Data[0].Created)
	assert.Equal("2019-08-20T13:26:24.07832Z", r.Data[0].Updated)
	assert.Equal("AamrnoXdpdmzjwE", r.Data[0].ETag)
	assert.Equal("b3", r.Data[0].Custom["a3"])
	assert.Equal("d3", r.Data[0].Custom["c3"])

	assert.Nil(err)
}
