package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestRemoveRelationshipBuildPathHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newDeleteRelationshipBuilder(pn).ID("r-123").IfMatchETag("1")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/relationships/r-123", pn.Config.SubscribeKey),
		path, []int{})

	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal("1", headers["If-Match"])
}

func TestRemoveRelationshipHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newDeleteRelationshipBuilder(pn)
	assert.Equal("DELETE", o.opts.httpMethod())
	assert.Equal(PNRemoveRelationshipOperation, o.opts.operationType())
}

func TestRemoveRelationshipValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newDeleteRelationshipBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingRelationshipID)

	o2 := newDeleteRelationshipBuilder(pn).ID("r-123")
	assert.Nil(o2.opts.validate())
}

func TestRemoveRelationshipResponseParsingEmptyBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	r, _, err := newPNRemoveRelationshipResponse([]byte(""), newDeleteRelationshipOpts(pn, pn.ctx), StatusResponse{})
	assert.Nil(err)
	assert.NotNil(r)
}

func TestRemoveRelationshipExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.RemoveRelationship().Execute()
	assert.Contains(err.Error(), StrMissingRelationshipID)
}
