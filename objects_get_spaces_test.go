package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/sprucehealth/pubnub-go/tests/helpers"
	"github.com/sprucehealth/pubnub-go/utils"
	"github.com/stretchr/testify/assert"
)

func AssertGetSpaces(t *testing.T, checkQueryParam, testContext bool) {
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

	o := newGetSpacesBuilder(pn)
	if testContext {
		o = newGetSpacesBuilderWithContext(pn, backgroundContext)
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
		fmt.Sprintf("/v1/objects/%s/spaces", pn.Config.SubscribeKey),
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

func TestGetSpaces(t *testing.T) {
	AssertGetSpaces(t, true, false)
}

func TestGetSpacesContext(t *testing.T) {
	AssertGetSpaces(t, true, true)
}

func TestGetSpacesResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getSpacesOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetSpacesResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

//{"status":200,"data":[{"id":"id0","name":"name","description":"desc","custom":{"a":"b"},"created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T13:26:08.341297Z","eTag":"Aee9zsKNndXlHw"},{"id":"id01","name":"name","description":"desc","custom":{"a":"b"},"created":"2019-08-20T14:44:52.799969Z","updated":"2019-08-20T14:44:52.799969Z","eTag":"Aee9zsKNndXlHw"}],"totalCount":2,"next":"Mg"}
func TestGetSpacesResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getSpacesOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":[{"id":"id0","name":"name","description":"desc","custom":{"a":"b"},"created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T13:26:08.341297Z","eTag":"Aee9zsKNndXlHw"},{"id":"id01","name":"name","description":"desc","custom":{"a":"b"},"created":"2019-08-20T14:44:52.799969Z","updated":"2019-08-20T14:44:52.799969Z","eTag":"Aee9zsKNndXlHw"}],"totalCount":2,"next":"Mg","prev":"Nd"}`)

	r, _, err := newPNGetSpacesResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(200, r.Status)
	assert.Equal(2, r.TotalCount)
	assert.Equal("Mg", r.Next)
	assert.Equal("Nd", r.Prev)
	assert.Equal("id0", r.Data[0].ID)
	assert.Equal("name", r.Data[0].Name)
	assert.Equal("desc", r.Data[0].Description)
	assert.Equal("2019-08-20T13:26:08.341297Z", r.Data[0].Created)
	assert.Equal("2019-08-20T13:26:08.341297Z", r.Data[0].Updated)
	assert.Equal("Aee9zsKNndXlHw", r.Data[0].ETag)
	assert.Equal("b", r.Data[0].Custom["a"])

	assert.Nil(err)
}
