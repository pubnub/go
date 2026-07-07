package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestDeleteEntityBuildPathHeaders(t *testing.T) {
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

	o.IfMatch("Dpqr3...")
	headers, err = o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal("Dpqr3...", headers["If-Match"])
}

func TestDeleteEntityHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newDeleteEntityBuilder(pn)
	assert.Equal("DELETE", o.opts.httpMethod())
	assert.Equal(PNDeleteEntityOperation, o.opts.operationType())
}

func TestDeleteEntityValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newDeleteEntityBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingEntityID)

	pn.Config.SubscribeKey = ""
	o2 := newDeleteEntityBuilder(pn).ID("entity-abc")
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestDeleteEntityResponseParsingEmptyBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	r, _, err := newPNDeleteEntityResponse([]byte(""), newDeleteEntityOpts(pn, pn.ctx), StatusResponse{})
	assert.Nil(err)
	assert.NotNil(r)
}

func TestDeleteEntityExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DeleteEntity().Execute()
	assert.Contains(err.Error(), StrMissingEntityID)
}
