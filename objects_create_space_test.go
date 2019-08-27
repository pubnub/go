package pubnub

import (
	"fmt"
	"testing"

	h "github.com/sprucehealth/pubnub-go/tests/helpers"
	"github.com/sprucehealth/pubnub-go/utils"
	"github.com/stretchr/testify/assert"
)

func AssertCreateSpace(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	incl := []PNUserSpaceInclude{
		PNUserSpaceCustom,
	}
	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := utils.EnumArrayToStringArray(fmt.Sprint(incl))

	o := newCreateSpaceBuilder(pn)
	if testContext {
		o = newCreateSpaceBuilderWithContext(pn, backgroundContext)
	}

	o.Include(incl)
	o.ID("id0")
	o.Name("name")
	o.Description("exturl")
	o.Custom(custom)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/objects/%s/spaces", pn.Config.SubscribeKey),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	expectedBody := "{\"id\":\"id0\",\"name\":\"name\",\"description\":\"exturl\",\"custom\":{\"a\":\"b\",\"c\":\"d\"}}"

	assert.Equal(expectedBody, string(body))

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
		assert.Equal(string(utils.JoinChannels(inclStr)), u.Get("include"))
	}

}

func TestCreateSpace(t *testing.T) {
	AssertCreateSpace(t, true, false)
}

func TestCreateSpaceContext(t *testing.T) {
	AssertCreateSpace(t, true, true)
}

func TestCreateSpaceResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &createSpaceOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNCreateSpaceResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

//{"status":200,"data":{"id":"id2","name":"name","description":"desc","custom":{"a":"b"},"created":"2019-08-19T15:05:51.493894Z","updated":"2019-08-19T15:05:51.493894Z","eTag":"Aee9zsKNndXlHw"}}
func TestCreateSpaceResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &createSpaceOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":{"id":"id2","name":"name","description":"desc","custom":{"a":"b"},"created":"2019-08-19T15:05:51.493894Z","updated":"2019-08-19T15:05:51.493894Z","eTag":"Aee9zsKNndXlHw"}}`)

	r, _, err := newPNCreateSpaceResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(200, r.Status)
	assert.Equal("id2", r.Data.ID)
	assert.Equal("name", r.Data.Name)
	assert.Equal("desc", r.Data.Description)
	assert.Equal("2019-08-19T15:05:51.493894Z", r.Data.Created)
	assert.Equal("2019-08-19T15:05:51.493894Z", r.Data.Updated)
	assert.Equal("Aee9zsKNndXlHw", r.Data.ETag)
	assert.Equal("b", r.Data.Custom["a"])

	assert.Nil(err)
}
