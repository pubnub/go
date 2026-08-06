package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCreateRelationshipBuildPathQueryBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newCreateRelationshipBuilder(pn)
	o.ID("r-123").
		EntityAID("u123").
		EntityBID("s456").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Status("active").
		Payload(map[string]interface{}{"role": "admin"}).
		QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/relationships", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))

	body, err := o.opts.buildBody()
	assert.Nil(err)

	var parsed createRelationshipBody
	assert.Nil(json.Unmarshal(body, &parsed))
	assert.Equal("r-123", parsed.Data.ID)
	assert.Equal("u123", parsed.Data.EntityAID)
	assert.Equal("s456", parsed.Data.EntityBID)
	assert.Equal("ProductOwner", parsed.Data.RelationshipClass)
	assert.Equal(1, parsed.Data.RelationshipClassVersion)
	assert.Equal("active", parsed.Data.Status)
	assert.Equal("admin", parsed.Data.Payload["role"])
}

func TestCreateRelationshipHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newCreateRelationshipBuilder(pn).
		EntityAID("u123").
		EntityBID("s456").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1)

	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(relationshipContentType, headers["Content-Type"])
}

func TestCreateRelationshipHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newCreateRelationshipBuilder(pn)
	assert.Equal("POST", o.opts.httpMethod())
	assert.Equal(PNCreateRelationshipOperation, o.opts.operationType())
}

func TestCreateRelationshipValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newCreateRelationshipBuilder(pn).
		EntityBID("s456").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1)
	assert.Contains(o.opts.validate().Error(), StrMissingEntityAID)

	o2 := newCreateRelationshipBuilder(pn).
		EntityAID("u123").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1)
	assert.Contains(o2.opts.validate().Error(), StrMissingEntityBID)

	o3 := newCreateRelationshipBuilder(pn).
		EntityAID("u123").
		EntityBID("s456").
		RelationshipClassVersion(1)
	assert.Contains(o3.opts.validate().Error(), StrMissingRelationshipClass)

	o4 := newCreateRelationshipBuilder(pn).
		EntityAID("u123").
		EntityBID("s456").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(0)
	assert.Contains(o4.opts.validate().Error(), StrInvalidRelationshipClassVersion)

	pn.Config.SubscribeKey = ""
	o5 := newCreateRelationshipBuilder(pn).
		EntityAID("u123").
		EntityBID("s456").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1)
	assert.Contains(o5.opts.validate().Error(), StrMissingSubKey)
}

func TestCreateRelationshipResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"status":201,"data":{"id":"r-123","entityAId":"u123","entityBId":"s456","relationshipClass":"ProductOwner","relationshipClassVersion":1,"status":"active","payload":{"role":"admin"},"createdAt":"2021-01-01T00:00:00.000Z","updatedAt":"2021-01-01T00:00:00.000Z","eTag":"1"}}`)

	r, _, err := newPNRelationshipResponse(jsonBytes, StatusResponse{}, PNCreateRelationshipOperation, pn)
	assert.Nil(err)
	assert.Equal(201, r.Status)
	assert.Equal("r-123", r.Data.ID)
	assert.Equal("u123", r.Data.EntityAID)
	assert.Equal("s456", r.Data.EntityBID)
	assert.Equal("ProductOwner", r.Data.RelationshipClass)
	assert.Equal(1, r.Data.RelationshipClassVersion)
	assert.Equal("1", r.Data.ETag)
	assert.Equal("admin", r.Data.Payload["role"])
}

func TestCreateRelationshipExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	_, _, err := pn.DataSync.CreateRelationship().
		EntityAID("u123").
		EntityBID("s456").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Execute()
	assert.Contains(err.Error(), StrMissingSubKey)
}
