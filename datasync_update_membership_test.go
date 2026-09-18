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
		Replace("/payload/role", "admin").
		Add("/payload/pinned", true).
		Remove("/payload/notificationsEnabled").
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

	var ops []PNJSONPatchOperation
	assert.Nil(json.Unmarshal(body, &ops))
	assert.Len(ops, 3)
	assert.Equal("replace", ops[0].Op)
	assert.Equal("/payload/role", ops[0].Path)
	assert.Equal("add", ops[1].Op)
	assert.Equal("remove", ops[2].Op)
}

func TestUpdateMembershipHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newUpdateMembershipBuilder(pn).ID("m-123").Replace("/status", "inactive").IfMatchETag("1")
	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(entityPatchContentType, headers["Content-Type"])
	assert.Equal("1", headers["If-Match"])
}

func TestUpdateMembershipHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newUpdateMembershipBuilder(pn)
	assert.Equal("PATCH", o.opts.httpMethod())
	assert.Equal(PNUpdateDataSyncMembershipOperation, o.opts.operationType())
}

func TestUpdateMembershipValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newUpdateMembershipBuilder(pn).Replace("/status", "inactive")
	assert.Contains(o.opts.validate().Error(), StrMissingMembershipID)

	o2 := newUpdateMembershipBuilder(pn).ID("m-123")
	assert.Contains(o2.opts.validate().Error(), StrMissingPatchOperations)

	pn.Config.SubscribeKey = ""
	o3 := newUpdateMembershipBuilder(pn).ID("m-123").Replace("/status", "inactive")
	assert.Contains(o3.opts.validate().Error(), StrMissingSubKey)
}

func TestUpdateMembershipExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.UpdateMembership().ID("m-123").Execute()
	assert.Contains(err.Error(), StrMissingPatchOperations)
}
