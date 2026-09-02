package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetRelationshipsDefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetRelationshipsBuilder(pn)
	assert.Equal(entitiesDefaultLimit, o.opts.Limit)
}

func TestGetRelationshipsBuildPathQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetRelationshipsBuilder(pn).
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		EntityAID("u123").
		EntityBID("s456").
		Cursor("TjIw").
		Limit(10).
		FilterFast("status == 'active'").
		Sort([]string{"-createdAt", "+id"}).
		QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/relationships", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("ProductOwner", query.Get("relationship_class"))
	assert.Equal("1", query.Get("relationship_class_version"))
	assert.Equal("u123", query.Get("entity_a_id"))
	assert.Equal("s456", query.Get("entity_b_id"))
	assert.Equal("TjIw", query.Get("cursor"))
	assert.Equal(strconv.Itoa(10), query.Get("limit"))
	assert.Equal("status == 'active'", query.Get("filter_fast"))
	assert.Equal("", query.Get("filter"))
	assert.Equal("-createdAt,+id", query.Get("sort"))
	assert.Equal("v1", query.Get("q1"))
}

func TestGetRelationshipsFilterQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetRelationshipsBuilder(pn).
		RelationshipClass("ProductOwner").
		Filter("status == 'active' AND role == 'admin'")
	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("status == 'active' AND role == 'admin'", query.Get("filter"))
	assert.Equal("", query.Get("filter_fast"))
}

func TestGetRelationshipsVersionOmittedWhenUnset(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetRelationshipsBuilder(pn).RelationshipClass("ProductOwner")
	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("", query.Get("relationship_class_version"))
	assert.Equal("", query.Get("entity_a_id"))
	assert.Equal("", query.Get("entity_b_id"))
}

func TestGetRelationshipsHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGetRelationshipsBuilder(pn)
	assert.Equal("GET", o.opts.httpMethod())
	assert.Equal(PNGetRelationshipsOperation, o.opts.operationType())
}

func TestGetRelationshipsValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetRelationshipsBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingRelationshipClass)

	o3 := newGetRelationshipsBuilder(pn).RelationshipClass("ProductOwner").
		FilterFast("status == 'active'").
		Filter("status == 'active' AND role == 'admin'")
	assert.Contains(o3.opts.validate().Error(), StrExclusiveDataSyncFilter)

	pn.Config.SubscribeKey = ""
	o2 := newGetRelationshipsBuilder(pn).RelationshipClass("ProductOwner")
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestGetRelationshipsResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"data":[{"id":"r-123","entityAId":"u123","entityBId":"s456","relationshipClass":"ProductOwner","relationshipClassVersion":1,"status":"active","eTag":"1","payload":{"custom":"fields"}}],"links":{"self":"/self","next":"/next","prev":null},"meta":{"has_next":true,"has_prev":false,"next_cursor":"TjIw","prev_cursor":null,"limit":20}}`)

	r, _, err := newPNRelationshipsResponse(jsonBytes, newGetRelationshipsOpts(pn, pn.ctx), StatusResponse{})
	assert.Nil(err)
	assert.Len(r.Data, 1)
	assert.Equal("r-123", r.Data[0].ID)
	assert.Equal("u123", r.Data[0].EntityAID)
	assert.NotNil(r.Meta)
	assert.True(r.Meta.HasNext)
	assert.Equal("TjIw", r.Meta.NextCursor)
	assert.Equal(20, r.Meta.Limit)
	assert.NotNil(r.Links)
	assert.Equal("/next", r.Links.Next)
}

func TestGetRelationshipsExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.GetRelationships().Execute()
	assert.Contains(err.Error(), StrMissingRelationshipClass)

	_, _, err = pn.DataSync.GetRelationships().
		RelationshipClass("ProductOwner").
		FilterFast("status == 'active'").
		Filter("status == 'active' AND role == 'admin'").
		Execute()
	assert.Contains(err.Error(), StrExclusiveDataSyncFilter)
}
