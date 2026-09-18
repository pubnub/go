package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestRemoveEntityBuildPathHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newDeleteEntityBuilder(pn).ID("entity-abc")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/entities/entity-abc", pn.Config.SubscribeKey),
		path, []int{})

	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	_, hasIfMatch := headers["If-Match"]
	assert.False(hasIfMatch)

	o.IfMatchETag("Dpqr3...")
	headers, err = o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal("Dpqr3...", headers["If-Match"])
}

func TestRemoveEntityHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newDeleteEntityBuilder(pn)
	assert.Equal("DELETE", o.opts.httpMethod())
	assert.Equal(PNRemoveEntityOperation, o.opts.operationType())
}

func TestRemoveEntityValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newDeleteEntityBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingEntityID)

	pn.Config.SubscribeKey = ""
	o2 := newDeleteEntityBuilder(pn).ID("entity-abc")
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestRemoveEntityResponseParsingEmptyBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	r, _, err := newPNRemoveEntityResponse([]byte(""), newDeleteEntityOpts(pn, pn.ctx), StatusResponse{})
	assert.Nil(err)
	assert.NotNil(r)
}

func TestRemoveEntityExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.RemoveEntity().Execute()
	assert.Contains(err.Error(), StrMissingEntityID)
}
