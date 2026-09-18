package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCreateMembershipBuildPathQueryBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newCreateMembershipBuilder(pn)
	o.ID("m-123").
		ChannelID("general").
		UserID("alice").
		RelationshipClassVersion(1).
		Status("active").
		Payload(map[string]interface{}{"role": "moderator"}).
		QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/memberships", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))

	body, err := o.opts.buildBody()
	assert.Nil(err)

	var parsed createMembershipBody
	assert.Nil(json.Unmarshal(body, &parsed))
	assert.Equal("m-123", parsed.Data.ID)
	assert.Equal("general", parsed.Data.ChannelID)
	assert.Equal("alice", parsed.Data.UserID)
	assert.Equal(1, parsed.Data.RelationshipClassVersion)
	assert.Equal("active", parsed.Data.Status)
	assert.Equal("moderator", parsed.Data.Payload["role"])
}

func TestCreateMembershipHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newCreateMembershipBuilder(pn).ChannelID("general").UserID("alice").RelationshipClassVersion(1)

	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(membershipContentType, headers["Content-Type"])
}

func TestCreateMembershipHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newCreateMembershipBuilder(pn)
	assert.Equal("POST", o.opts.httpMethod())
	assert.Equal(PNCreateDataSyncMembershipOperation, o.opts.operationType())
}

func TestCreateMembershipValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newCreateMembershipBuilder(pn).UserID("alice").RelationshipClassVersion(1)
	assert.Contains(o.opts.validate().Error(), StrMissingChannelID)

	o2 := newCreateMembershipBuilder(pn).ChannelID("general").RelationshipClassVersion(1)
	assert.Contains(o2.opts.validate().Error(), StrMissingUserID)

	o3 := newCreateMembershipBuilder(pn).ChannelID("general").UserID("alice").RelationshipClassVersion(0)
	assert.Contains(o3.opts.validate().Error(), StrInvalidRelationshipClassVersion)

	pn.Config.SubscribeKey = ""
	o4 := newCreateMembershipBuilder(pn).ChannelID("general").UserID("alice").RelationshipClassVersion(1)
	assert.Contains(o4.opts.validate().Error(), StrMissingSubKey)
}

func TestCreateMembershipResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"status":201,"data":{"id":"m-123","channelId":"general","userId":"alice","relationshipClass":"Membership","relationshipClassVersion":1,"status":"active","payload":{"role":"moderator"},"eTag":"1"}}`)

	r, _, err := newPNDataSyncMembershipResponse(jsonBytes, StatusResponse{}, PNCreateDataSyncMembershipOperation, pn)
	assert.Nil(err)
	assert.Equal(201, r.Status)
	assert.Equal("m-123", r.Data.ID)
	assert.Equal("general", r.Data.ChannelID)
	assert.Equal("alice", r.Data.UserID)
	assert.Equal("Membership", r.Data.RelationshipClass)
	assert.Equal("moderator", r.Data.Payload["role"])
}

func TestCreateMembershipExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	_, _, err := pn.DataSync.CreateMembership().ChannelID("general").UserID("alice").RelationshipClassVersion(1).Execute()
	assert.Contains(err.Error(), StrMissingSubKey)
}
