package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMembershipBuildPathQueryBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newUpdateMembershipBuilder(pn).
		ID("m-123").
		RelationshipClassVersion(1).
		Status("active").
		Payload(map[string]interface{}{"role": "admin"}).
		IfMatchETag("1").
		QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/memberships/m-123", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))

	body, err := o.opts.buildBody()
	assert.Nil(err)

	var parsed updateMembershipBody
	assert.Nil(json.Unmarshal(body, &parsed))
	assert.Equal(1, parsed.Data.RelationshipClassVersion)
	assert.Equal("active", parsed.Data.Status)
	assert.Equal("admin", parsed.Data.Payload["role"])
}

func TestUpdateMembershipHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newUpdateMembershipBuilder(pn).ID("m-123").RelationshipClassVersion(1).IfMatchETag("1")
	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(membershipContentType, headers["Content-Type"])
	assert.Equal("1", headers["If-Match"])
}

func TestUpdateMembershipHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newUpdateMembershipBuilder(pn)
	assert.Equal("PUT", o.opts.httpMethod())
	assert.Equal(PNUpdateDataSyncMembershipOperation, o.opts.operationType())
}

func TestUpdateMembershipValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newUpdateMembershipBuilder(pn).RelationshipClassVersion(1)
	assert.Contains(o.opts.validate().Error(), StrMissingMembershipID)

	o2 := newUpdateMembershipBuilder(pn).ID("m-123").RelationshipClassVersion(0)
	assert.Contains(o2.opts.validate().Error(), StrInvalidRelationshipClassVersion)

	pn.Config.SubscribeKey = ""
	o3 := newUpdateMembershipBuilder(pn).ID("m-123").RelationshipClassVersion(1)
	assert.Contains(o3.opts.validate().Error(), StrMissingSubKey)
}

func TestUpdateMembershipExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.UpdateMembership().RelationshipClassVersion(1).Execute()
	assert.Contains(err.Error(), StrMissingMembershipID)
}
