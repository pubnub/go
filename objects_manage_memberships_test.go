package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/sprucehealth/pubnub-go/tests/helpers"
	"github.com/sprucehealth/pubnub-go/utils"
	"github.com/stretchr/testify/assert"
)

func AssertManageMemberships(t *testing.T, checkQueryParam, testContext bool) {
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

	o := newManageMembershipsBuilder(pn)
	if testContext {
		o = newManageMembershipsBuilderWithContext(pn, backgroundContext)
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

	id0 := "id0"

	custom3 := make(map[string]interface{})
	custom3["a3"] = "b3"
	custom3["c3"] = "d3"

	in := PNMembershipsInput{
		ID:     id0,
		Custom: custom3,
	}

	inArr := []PNMembershipsInput{
		in,
	}

	custom4 := make(map[string]interface{})
	custom4["a4"] = "b4"
	custom4["c4"] = "d4"

	up := PNMembershipsInput{
		ID:     id0,
		Custom: custom4,
	}

	upArr := []PNMembershipsInput{
		up,
	}

	re := PNMembershipsRemove{
		ID: id0,
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
		fmt.Sprintf("/v1/objects/%s/users/%s/spaces", pn.Config.SubscribeKey, userID),
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

// {"status":200,"data":[{"id":"spaceid3","custom":{"a3":"b3","c3":"d3"},"space":{"id":"spaceid3","name":"spaceid3name","description":"spaceid3desc","custom":{"a":"b"},"created":"2019-08-23T10:34:43.985248Z","updated":"2019-08-23T10:34:43.985248Z","eTag":"Aazjn7vC3oDDYw"},"created":"2019-08-23T10:41:17.156491Z","updated":"2019-08-23T10:41:17.156491Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":1,"next":"MQ"}
func TestManageMembershipsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &manageMembershipsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":[{"id":"spaceid3","custom":{"a3":"b3","c3":"d3"},"space":{"id":"spaceid3","name":"spaceid3name","description":"spaceid3desc","custom":{"a":"b"},"created":"2019-08-23T10:34:43.985248Z","updated":"2019-08-23T10:34:43.985248Z","eTag":"Aazjn7vC3oDDYw"},"created":"2019-08-23T10:41:17.156491Z","updated":"2019-08-23T10:41:17.156491Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":1,"next":"MQ"}`)

	r, _, err := newPNManageMembershipsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(200, r.Status)
	assert.Equal(1, r.TotalCount)
	assert.Equal("MQ", r.Next)
	assert.Equal("spaceid3", r.Data[0].ID)
	assert.Equal("spaceid3", r.Data[0].Space.ID)
	assert.Equal("spaceid3name", r.Data[0].Space.Name)
	assert.Equal("spaceid3desc", r.Data[0].Space.Description)
	assert.Equal("2019-08-23T10:34:43.985248Z", r.Data[0].Space.Created)
	assert.Equal("2019-08-23T10:34:43.985248Z", r.Data[0].Space.Updated)
	assert.Equal("Aazjn7vC3oDDYw", r.Data[0].Space.ETag)
	assert.Equal("b", r.Data[0].Space.Custom["a"])
	assert.Equal("2019-08-23T10:41:17.156491Z", r.Data[0].Created)
	assert.Equal("2019-08-23T10:41:17.156491Z", r.Data[0].Updated)
	assert.Equal("AamrnoXdpdmzjwE", r.Data[0].ETag)
	assert.Equal("b3", r.Data[0].Custom["a3"])
	assert.Equal("d3", r.Data[0].Custom["c3"])

	assert.Nil(err)
}
