package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/sprucehealth/pubnub-go/tests/helpers"
	"github.com/sprucehealth/pubnub-go/utils"
	"github.com/stretchr/testify/assert"
)

func AssertGetUsers(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	incl := []PNUserSpaceInclude{
		PNUserSpaceCustom,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := utils.EnumArrayToStringArray(fmt.Sprint(incl))

	o := newGetUsersBuilder(pn)
	if testContext {
		o = newGetUsersBuilderWithContext(pn, backgroundContext)
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

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/objects/%s/users", pn.Config.SubscribeKey),
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

func TestGetUsers(t *testing.T) {
	AssertGetUsers(t, true, false)
}

func TestGetUsersContext(t *testing.T) {
	AssertGetUsers(t, true, true)
}

func TestGetUsersResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getUsersOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetUsersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

//{"status":200,"data":[{"id":"id2","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-19T14:44:54.837392Z","updated":"2019-08-19T14:44:54.837392Z","eTag":"AbyT4v2p6K7fpQE"},{"id":"id0","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-20T13:26:19.140324Z","updated":"2019-08-20T13:26:19.140324Z","eTag":"AbyT4v2p6K7fpQE"}],"totalCount":2,"next":"Mg"}
func TestGetUsersResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getUsersOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":[{"id":"id2","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-19T14:44:54.837392Z","updated":"2019-08-19T14:44:54.837392Z","eTag":"AbyT4v2p6K7fpQE"},{"id":"id0","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-20T13:26:19.140324Z","updated":"2019-08-20T13:26:19.140324Z","eTag":"AbyT4v2p6K7fpQE"}],"totalCount":2,"next":"Mg","prev":"Nd"}`)

	r, _, err := newPNGetUsersResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(200, r.Status)
	assert.Equal(2, r.TotalCount)
	assert.Equal("Mg", r.Next)
	assert.Equal("Nd", r.Prev)
	assert.Equal("id2", r.Data[0].ID)
	assert.Equal("name", r.Data[0].Name)
	assert.Equal("extid", r.Data[0].ExternalID)
	assert.Equal("purl", r.Data[0].ProfileURL)
	assert.Equal("email", r.Data[0].Email)
	assert.Equal("2019-08-19T14:44:54.837392Z", r.Data[0].Created)
	assert.Equal("2019-08-19T14:44:54.837392Z", r.Data[0].Updated)
	assert.Equal("AbyT4v2p6K7fpQE", r.Data[0].ETag)
	assert.Equal("b", r.Data[0].Custom["a"])
	assert.Equal("d", r.Data[0].Custom["c"])

	assert.Nil(err)
}
