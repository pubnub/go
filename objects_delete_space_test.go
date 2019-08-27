package pubnub

import (
	"fmt"
	"testing"

	h "github.com/sprucehealth/pubnub-go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertDeleteSpace(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newDeleteSpaceBuilder(pn)
	if testContext {
		o = newDeleteSpaceBuilderWithContext(pn, backgroundContext)
	}

	o.ID("id0")
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/objects/%s/spaces/%s", pn.Config.SubscribeKey, "id0"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

}

func TestDeleteSpace(t *testing.T) {
	AssertDeleteSpace(t, true, false)
}

func TestDeleteSpaceContext(t *testing.T) {
	AssertDeleteSpace(t, true, true)
}

func TestDeleteSpaceResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &deleteSpaceOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNDeleteSpaceResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

//{"status":200,"data":null}
func TestDeleteSpaceResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &deleteSpaceOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":null}`)

	r, _, err := newPNDeleteSpaceResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(200, r.Status)
	assert.Equal(nil, r.Data)

	assert.Nil(err)
}
