package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/pubnub/go/utils"
	"github.com/stretchr/testify/assert"
)

func AssertManageMembers(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	incl := []PNSpaceMembershipsIncude{
		PNSpaceMembershipsCustom,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := utils.EnumArrayToStringArray(fmt.Sprint(incl))

	o := newManageMembersBuilder(pn)
	if testContext {
		o = newManageMembersBuilderWithContext(pn, backgroundContext)
	}

	spaceId := "id0"
	limit := 90
	start := "Mxmy"
	end := "Nxny"

	o.SpaceId(spaceId)
	o.Include(incl)
	o.Limit(limit)
	o.Start(start)
	o.End(end)
	o.Count(false)
	o.QueryParam(queryParam)

	id0 := "id0"

	custom := make(map[string]interface{})
	custom["a1"] = "b1"
	custom["c1"] = "d1"

	in := PNSpaceMembershipInput{
		Id:     id0,
		Custom: custom,
	}

	inArr := []PNSpaceMembershipInput{
		in,
	}

	custom2 := make(map[string]interface{})
	custom2["a2"] = "b2"
	custom2["c2"] = "d2"

	up := PNSpaceMembershipInput{
		Id:     id0,
		Custom: custom2,
	}

	upArr := []PNSpaceMembershipInput{
		up,
	}

	re := PNSpaceMembershipRemove{
		Id: id0,
	}

	reArr := []PNSpaceMembershipRemove{
		re,
	}
	o.Add(inArr)
	o.Update(upArr)
	o.Remove(reArr)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/objects/%s/spaces/%s/users", pn.Config.SubscribeKey, spaceId),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	expectedBody := "{\"add\":[{\"id\":\"id0\",\"custom\":{\"a1\":\"b1\",\"c1\":\"d1\"}}],\"update\":[{\"id\":\"id0\",\"custom\":{\"a2\":\"b2\",\"c2\":\"d2\"}}],\"remove\":[{\"id\":\"id0\"}]}"

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

func TestManageMembers(t *testing.T) {
	AssertManageMembers(t, true, false)
}

func TestManageMembersContext(t *testing.T) {
	AssertManageMembers(t, true, true)
}

func TestManageMembersResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &manageMembersOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNManageMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

//add: {"status":200,"data":[{"id":"userid2","custom":{"a1":"b1","c1":"d1"},"created":"2019-08-21T11:43:35.889327Z","updated":"2019-08-21T11:43:35.889327Z","eTag":"AZK3l4nQsrWG9gE"}],"totalCount":1,"next":"MQ"}

//update: {"status":200,"data":[{"id":"userid0","custom":{"a2":"b2","c2":"d2"},"created":"2019-08-21T09:08:22.49193Z","updated":"2019-08-21T11:41:43.613345Z","eTag":"AdCFwIrDze6g/AE"}],"totalCount":1,"next":"MQ"}
func TestManageMembersResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &manageMembersOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":[{"id":"userid2","custom":{"a1":"b1","c1":"d1"},"created":"2019-08-21T11:43:35.889327Z","updated":"2019-08-21T11:43:35.889327Z","eTag":"AZK3l4nQsrWG9gE"}],"totalCount":1,"next":"MQ"}`)

	r, _, err := newPNManageMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(200, r.Status)
	assert.Equal(1, r.TotalCount)
	assert.Equal("MQ", r.Next)
	assert.Equal("userid2", r.Data[0].Id)
	assert.Equal("2019-08-21T11:43:35.889327Z", r.Data[0].Created)
	assert.Equal("2019-08-21T11:43:35.889327Z", r.Data[0].Updated)
	assert.Equal("AZK3l4nQsrWG9gE", r.Data[0].ETag)
	assert.Equal("b1", r.Data[0].Custom["a1"])
	assert.Equal("d1", r.Data[0].Custom["c1"])

	assert.Nil(err)
}
