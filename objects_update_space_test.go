package pubnub

import (
	"fmt"
	"testing"

	h "github.com/sprucehealth/pubnub-go/tests/helpers"
	"github.com/sprucehealth/pubnub-go/utils"
	"github.com/stretchr/testify/assert"
)

func AssertUpdateSpace(t *testing.T, checkQueryParam, testContext bool) {
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

	o := newUpdateSpaceBuilder(pn)
	if testContext {
		o = newUpdateSpaceBuilderWithContext(pn, backgroundContext)
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
		fmt.Sprintf("/v1/objects/%s/spaces/%s", pn.Config.SubscribeKey, "id0"),
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

func TestUpdateSpace(t *testing.T) {
	AssertUpdateSpace(t, true, false)
}

func TestUpdateSpaceContext(t *testing.T) {
	AssertUpdateSpace(t, true, true)
}

func TestUpdateSpaceResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &updateSpaceOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNUpdateSpaceResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

//{"status":200,"data":{"id":"id0","name":"name","description":"desc","custom":{"a":"b","c":"d"},"created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T14:48:11.675743Z","eTag":"AYKH2s7ZlYKoJA"}}
func TestUpdateSpaceResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &updateSpaceOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":{"id":"id0","name":"name","description":"desc","custom":{"a":"b","c":"d"},"created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T14:48:11.675743Z","eTag":"AYKH2s7ZlYKoJA"}}`)

	r, _, err := newPNUpdateSpaceResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(200, r.Status)
	assert.Equal("id0", r.Data.ID)
	assert.Equal("name", r.Data.Name)
	assert.Equal("desc", r.Data.Description)
	assert.Equal("2019-08-20T13:26:08.341297Z", r.Data.Created)
	assert.Equal("2019-08-20T14:48:11.675743Z", r.Data.Updated)
	assert.Equal("AYKH2s7ZlYKoJA", r.Data.ETag)
	assert.Equal("b", r.Data.Custom["a"])

	assert.Nil(err)
}
