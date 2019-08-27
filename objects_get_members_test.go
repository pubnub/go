package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/sprucehealth/pubnub-go/tests/helpers"
	"github.com/sprucehealth/pubnub-go/utils"
	"github.com/stretchr/testify/assert"
)

func AssertGetMembers(t *testing.T, checkQueryParam, testContext bool) {
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

	o := newGetMembersBuilder(pn)
	if testContext {
		o = newGetMembersBuilderWithContext(pn, backgroundContext)
	}

	spaceID := "id0"
	limit := 90
	start := "Mxmy"
	end := "Nxny"

	o.SpaceID(spaceID)
	o.Include(incl)
	o.Limit(limit)
	o.Start(start)
	o.End(end)
	o.Count(false)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/objects/%s/spaces/%s/users", pn.Config.SubscribeKey, "id0"),
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

func TestGetMembers(t *testing.T) {
	AssertGetMembers(t, true, false)
}

func TestGetMembersContext(t *testing.T) {
	AssertGetMembers(t, true, true)
}

func TestGetMembersResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getMembersOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

//{"status":200,"data":[{"id":"id0","custom":{"a3":"b3","c3":"d3"},"user":{"id":"id0","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-20T13:26:19.140324Z","updated":"2019-08-20T13:26:19.140324Z","eTag":"AbyT4v2p6K7fpQE"},"created":"2019-08-20T13:26:24.07832Z","updated":"2019-08-20T13:26:24.07832Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":1,"next":"MQ","prev":"NQ"}
func TestGetMembersResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getMembersOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":[{"id":"id0","custom":{"a3":"b3","c3":"d3"},"user":{"id":"id0","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-20T13:26:19.140324Z","updated":"2019-08-20T13:26:19.140324Z","eTag":"AbyT4v2p6K7fpQE"},"created":"2019-08-20T13:26:24.07832Z","updated":"2019-08-20T13:26:24.07832Z","eTag":"AamrnoXdpdmzjwE"}],"totalCount":1,"next":"MQ","prev":"NQ"}`)

	r, _, err := newPNGetMembersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(200, r.Status)
	assert.Equal(1, r.TotalCount)
	assert.Equal("MQ", r.Next)
	assert.Equal("NQ", r.Prev)
	assert.Equal("id0", r.Data[0].ID)
	assert.Equal("name", r.Data[0].User.Name)
	assert.Equal("extid", r.Data[0].User.ExternalID)
	assert.Equal("purl", r.Data[0].User.ProfileURL)
	assert.Equal("email", r.Data[0].User.Email)
	assert.Equal("2019-08-20T13:26:19.140324Z", r.Data[0].User.Created)
	assert.Equal("2019-08-20T13:26:19.140324Z", r.Data[0].User.Updated)
	assert.Equal("AbyT4v2p6K7fpQE", r.Data[0].User.ETag)
	assert.Equal("b", r.Data[0].User.Custom["a"])
	assert.Equal("d", r.Data[0].User.Custom["c"])
	assert.Equal("2019-08-20T13:26:24.07832Z", r.Data[0].Created)
	assert.Equal("2019-08-20T13:26:24.07832Z", r.Data[0].Updated)
	assert.Equal("AamrnoXdpdmzjwE", r.Data[0].ETag)
	assert.Equal("b3", r.Data[0].Custom["a3"])
	assert.Equal("d3", r.Data[0].Custom["c3"])

	assert.Nil(err)
}
