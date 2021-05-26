package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v5/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertRemoveChannelMetadata(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newRemoveChannelMetadataBuilder(pn)
	if testContext {
		o = newRemoveChannelMetadataBuilderWithContext(pn, backgroundContext)
	}

	o.Channel("id0")
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/channels/%s", pn.Config.SubscribeKey, "id0"),
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

func TestRemoveChannelMetadata(t *testing.T) {
	AssertRemoveChannelMetadata(t, true, false)
}

func TestRemoveChannelMetadataContext(t *testing.T) {
	AssertRemoveChannelMetadata(t, true, true)
}

func TestRemoveChannelMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &removeChannelMetadataOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNRemoveChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRemoveChannelMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &removeChannelMetadataOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":null}`)

	r, _, err := newPNRemoveChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(nil, r.Data)

	assert.Nil(err)
}
