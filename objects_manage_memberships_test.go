package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/pubnub/go/utils"
	"github.com/stretchr/testify/assert"
)

func AssertManageMemberships(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	incl := []PNMembersInclude{
		PNMembersCustom,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := utils.EnumArrayToStringArray(fmt.Sprint(incl))

	o := newManageMembershipsBuilder(pn)
	if testContext {
		o = newManageMembershipsBuilderWithContext(pn, backgroundContext)
	}

	userId := "id0"
	limit := 90
	start := "Mxmy"
	end := "Nxny"

	o.UserId(userId)
	o.Include(incl)
	o.Limit(limit)
	o.Start(start)
	o.End(end)
	o.Count(false)
	o.QueryParam(queryParam)

	id0 := "id0"

	custom3 := make(map[string]interface{})
	custom3["a3"] = "b3"
	custom3["c3"] = "d3"

	in := PNMembershipsInput{
		Id:     id0,
		Custom: custom3,
	}

	inArr := []PNMembershipsInput{
		in,
	}

	custom4 := make(map[string]interface{})
	custom4["a4"] = "b4"
	custom4["c4"] = "d4"

	up := PNMembershipsInput{
		Id:     id0,
		Custom: custom4,
	}

	upArr := []PNMembershipsInput{
		up,
	}

	re := PNMembershipsRemove{
		Id: id0,
	}

	reArr := []PNMembershipsRemove{
		re,
	}

	o.Add(inArr)
	o.Update(upArr)
	o.Remove(reArr)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/objects/%s/users/%s/spaces", pn.Config.SubscribeKey, userId),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	expectedBody := "{\"add\":[{\"id\":\"id0\",\"custom\":{\"a3\":\"b3\",\"c3\":\"d3\"}}],\"update\":[{\"id\":\"id0\",\"custom\":{\"a4\":\"b4\",\"c4\":\"d4\"}}],\"remove\":[{\"id\":\"id0\"}]}"

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
	}

}

func TestManageMemberships(t *testing.T) {
	AssertManageMemberships(t, true, false)
}

func TestManageMembershipsContext(t *testing.T) {
	AssertManageMemberships(t, true, true)
}

func TestManageMembershipsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &manageMembershipsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNManageMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

// add{"status":200,"data":[{"id":"spaceid2","custom":{"a1":"b1","c1":"d1"},"created":"2019-08-21T11:43:35.889327Z","updated":"2019-08-21T11:43:35.889327Z","eTag":"AZK3l4nQsrWG9gE"},{"id":"spaceid0","custom":{"a3":"b3","c3":"d3"},"created":"2019-08-21T11:44:30.893128Z","updated":"2019-08-21T11:44:30.893128Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":2,"next":"Mg"}
// update: {"status":200,"data":[{"id":"spaceid0","custom":{"a4":"b4","c4":"d4"},"created":"2019-08-21T09:08:22.49193Z","updated":"2019-08-21T11:39:15.159336Z","eTag":"AZa25Pq3w6iHjwE"}],"totalCount":1,"next":"MQ"}
func TestManageMembershipsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &manageMembershipsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":[{"id":"spaceid2","custom":{"a1":"b1","c1":"d1"},"created":"2019-08-21T11:43:35.889327Z","updated":"2019-08-21T11:43:35.889327Z","eTag":"AZK3l4nQsrWG9gE"},{"id":"spaceid0","custom":{"a3":"b3","c3":"d3"},"created":"2019-08-21T11:44:30.893128Z","updated":"2019-08-21T11:44:30.893128Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":2,"next":"Mg"}`)

	r, _, err := newPNManageMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(200, r.Status)
	assert.Equal(2, r.TotalCount)
	assert.Equal("Mg", r.Next)
	assert.Equal("spaceid2", r.Data[0].Id)
	assert.Equal("2019-08-21T11:43:35.889327Z", r.Data[0].Created)
	assert.Equal("2019-08-21T11:43:35.889327Z", r.Data[0].Updated)
	assert.Equal("AZK3l4nQsrWG9gE", r.Data[0].ETag)
	assert.Equal("b1", r.Data[0].Custom["a1"])
	assert.Equal("d1", r.Data[0].Custom["c1"])

	assert.Nil(err)
}
