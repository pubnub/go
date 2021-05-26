package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v5/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertRemoveUUIDMetadata(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newRemoveUUIDMetadataBuilder(pn)
	if testContext {
		o = newRemoveUUIDMetadataBuilderWithContext(pn, backgroundContext)
	}

	o.UUID("id0")
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/uuids/%s", pn.Config.SubscribeKey, "id0"),
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

func TestRemoveUUIDMetadata(t *testing.T) {
	AssertRemoveUUIDMetadata(t, true, false)
}

func TestRemoveUUIDMetadataContext(t *testing.T) {
	AssertRemoveUUIDMetadata(t, true, true)
}

func TestRemoveUUIDMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &removeUUIDMetadataOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNRemoveUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRemoveUUIDMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &removeUUIDMetadataOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":null}`)

	r, _, err := newPNRemoveUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(nil, r.Data)

	assert.Nil(err)
}
