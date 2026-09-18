package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetRelationshipBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetRelationshipBuilder(pn).ID("r-123").QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/relationships/r-123", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))
}

func TestGetRelationshipHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGetRelationshipBuilder(pn)
	assert.Equal("GET", o.opts.httpMethod())
	assert.Equal(PNGetRelationshipOperation, o.opts.operationType())
}

func TestGetRelationshipValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetRelationshipBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingRelationshipID)

	o2 := newGetRelationshipBuilder(pn).ID("r-123")
	assert.Nil(o2.opts.validate())
}

func TestGetRelationshipResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"status":200,"data":{"id":"r-123","entityAId":"u123","entityBId":"s456","relationshipClass":"ProductOwner","relationshipClassVersion":1,"status":"active","payload":{"custom":"fields"},"eTag":"1"}}`)

	r, _, err := newPNRelationshipResponse(jsonBytes, StatusResponse{}, PNGetRelationshipOperation, pn)
	assert.Nil(err)
	assert.Equal("r-123", r.Data.ID)
	assert.Equal("u123", r.Data.EntityAID)
	assert.Equal("fields", r.Data.Payload["custom"])
}

func TestGetRelationshipExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.GetRelationship().Execute()
	assert.Contains(err.Error(), StrMissingRelationshipID)
}
