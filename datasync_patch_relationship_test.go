package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestPatchRelationshipBuildPathHeadersBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPatchRelationshipBuilder(pn).
		ID("r-123").
		Add("/payload/custom/role", "admin").
		Replace("/status", "inactive").
		IfMatchETag("1")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/relationships/r-123", pn.Config.SubscribeKey),
		path, []int{})

	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(entityPatchContentType, headers["Content-Type"])
	assert.Equal("1", headers["If-Match"])

	body, err := o.opts.buildBody()
	assert.Nil(err)

	var ops []PNJSONPatchOperation
	assert.Nil(json.Unmarshal(body, &ops))
	assert.Len(ops, 2)
	assert.Equal("add", ops[0].Op)
	assert.Equal("/payload/custom/role", ops[0].Path)
	assert.Equal("admin", ops[0].Value)
	assert.Equal("replace", ops[1].Op)
	assert.Equal("/status", ops[1].Path)
}

func TestPatchRelationshipHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newPatchRelationshipBuilder(pn)
	assert.Equal("PATCH", o.opts.httpMethod())
	assert.Equal(PNPatchRelationshipOperation, o.opts.operationType())
}

func TestPatchRelationshipValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPatchRelationshipBuilder(pn).Add("/status", "active")
	assert.Contains(o.opts.validate().Error(), StrMissingRelationshipID)

	o2 := newPatchRelationshipBuilder(pn).ID("r-123")
	assert.Contains(o2.opts.validate().Error(), StrMissingPatchOperations)

	o3 := newPatchRelationshipBuilder(pn).ID("r-123").Remove("/payload/temp")
	assert.Nil(o3.opts.validate())
}

func TestPatchRelationshipResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"status":200,"data":{"id":"r-123","entityAId":"u123","entityBId":"s456","relationshipClass":"ProductOwner","relationshipClassVersion":1,"status":"active","payload":{"custom":{"role":"admin"}},"eTag":"2"}}`)

	r, _, err := newPNRelationshipResponse(jsonBytes, StatusResponse{}, PNPatchRelationshipOperation, pn)
	assert.Nil(err)
	assert.Equal("2", r.Data.ETag)
	custom := r.Data.Payload["custom"].(map[string]interface{})
	assert.Equal("admin", custom["role"])
}

func TestPatchRelationshipExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.PatchRelationship().ID("r-123").Execute()
	assert.Contains(err.Error(), StrMissingPatchOperations)
}
