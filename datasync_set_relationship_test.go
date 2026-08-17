package pubnub

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestSetRelationshipBuildPathHeadersBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newSetRelationshipBuilder(pn).
		ID("r-123").
		RelationshipClassVersion(2).
		Status("active").
		Payload(map[string]interface{}{"custom": "fields"}).
		IfMatchETag("1")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/relationships/r-123", pn.Config.SubscribeKey),
		path, []int{})

	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(relationshipContentType, headers["Content-Type"])
	assert.Equal("1", headers["If-Match"])

	body, err := o.opts.buildBody()
	assert.Nil(err)

	var parsed setRelationshipBody
	assert.Nil(json.Unmarshal(body, &parsed))
	assert.Equal(2, parsed.Data.RelationshipClassVersion)
	assert.Equal("active", parsed.Data.Status)
	assert.Equal("fields", parsed.Data.Payload["custom"])

	// Immutable link fields must NOT be present in the body.
	assert.False(strings.Contains(string(body), "entityAId"))
	assert.False(strings.Contains(string(body), "entityBId"))
	assert.False(strings.Contains(string(body), "relationshipClass\":"))
}

func TestSetRelationshipHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newSetRelationshipBuilder(pn)
	assert.Equal("PUT", o.opts.httpMethod())
	assert.Equal(PNSetRelationshipOperation, o.opts.operationType())
}

func TestSetRelationshipValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newSetRelationshipBuilder(pn).RelationshipClassVersion(1)
	assert.Contains(o.opts.validate().Error(), StrMissingRelationshipID)

	o2 := newSetRelationshipBuilder(pn).ID("r-123").RelationshipClassVersion(0)
	assert.Contains(o2.opts.validate().Error(), StrInvalidRelationshipClassVersion)
}

func TestSetRelationshipResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"status":200,"data":{"id":"r-123","entityAId":"u123","entityBId":"s456","relationshipClass":"ProductOwner","relationshipClassVersion":2,"status":"active","eTag":"2"}}`)

	r, _, err := newPNRelationshipResponse(jsonBytes, StatusResponse{}, PNSetRelationshipOperation, pn)
	assert.Nil(err)
	assert.Equal(2, r.Data.RelationshipClassVersion)
	assert.Equal("2", r.Data.ETag)
}

func TestSetRelationshipExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.SetRelationship().RelationshipClassVersion(2).Execute()
	assert.Contains(err.Error(), StrMissingRelationshipID)
}
