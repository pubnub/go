package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/sprucehealth/pubnub-go/tests/helpers"
	"github.com/sprucehealth/pubnub-go/utils"
	"github.com/stretchr/testify/assert"
)

func AssertGetMemberships(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	incl := []PNMembershipsInclude{
		PNMembershipsCustom,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := utils.EnumArrayToStringArray(fmt.Sprint(incl))

	o := newGetMembershipsBuilder(pn)
	if testContext {
		o = newGetMembershipsBuilderWithContext(pn, backgroundContext)
	}

	userID := "id0"
	limit := 90
	start := "Mxmy"
	end := "Nxny"

	o.UserID(userID)
	o.Include(incl)
	o.Limit(limit)
	o.Start(start)
	o.End(end)
	o.Count(false)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/objects/%s/users/%s/spaces", pn.Config.SubscribeKey, "id0"),
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
	}

}

func TestGetMemberships(t *testing.T) {
	AssertGetMemberships(t, true, false)
}

func TestGetMembershipsContext(t *testing.T) {
	AssertGetMemberships(t, true, true)
}

func TestGetMembershipsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getMembershipsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

//{"status":200,"data":[{"id":"id0","custom":{"a3":"b3","c3":"d3"},"space":{"id":"id0","name":"name","description":"desc","custom":{"a":"b"},"created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T13:26:08.341297Z","eTag":"Aee9zsKNndXlHw"},"created":"2019-08-20T13:26:24.07832Z","updated":"2019-08-20T13:26:24.07832Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":1,"next":"MQ"}
func TestGetMembershipsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getMembershipsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":[{"id":"id0","custom":{"a3":"b3","c3":"d3"},"space":{"id":"id0","name":"name","description":"desc","custom":{"a":"b"},"created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T13:26:08.341297Z","eTag":"Aee9zsKNndXlHw"},"created":"2019-08-20T13:26:24.07832Z","updated":"2019-08-20T13:26:24.07832Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":1,"next":"MQ", "prev":"NQ"}`)

	r, _, err := newPNGetMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(200, r.Status)
	assert.Equal(1, r.TotalCount)
	assert.Equal("MQ", r.Next)
	assert.Equal("NQ", r.Prev)
	assert.Equal("id0", r.Data[0].ID)
	assert.Equal("name", r.Data[0].Space.Name)
	assert.Equal("desc", r.Data[0].Space.Description)
	assert.Equal("2019-08-20T13:26:08.341297Z", r.Data[0].Space.Created)
	assert.Equal("2019-08-20T13:26:08.341297Z", r.Data[0].Space.Updated)
	assert.Equal("Aee9zsKNndXlHw", r.Data[0].Space.ETag)
	assert.Equal("b", r.Data[0].Space.Custom["a"])
	assert.Equal(nil, r.Data[0].Space.Custom["c"])
	assert.Equal("2019-08-20T13:26:24.07832Z", r.Data[0].Created)
	assert.Equal("2019-08-20T13:26:24.07832Z", r.Data[0].Updated)
	assert.Equal("AamrnoXdpdmzjwE", r.Data[0].ETag)
	assert.Equal("b3", r.Data[0].Custom["a3"])
	assert.Equal("d3", r.Data[0].Custom["c3"])

	assert.Nil(err)
}
