package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestPatchUserBuildPathQueryBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPatchUserBuilder(pn).
		ID("alice").
		Replace("/payload/profileUrl", "https://example.com/alice").
		Add("/payload/phone", "+1-555-0100").
		Remove("/payload/isActive").
		IfMatchETag("1").
		QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/users/alice", pn.Config.SubscribeKey),
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
	assert.Equal("/payload/profileUrl", ops[0].Path)
	assert.Equal("add", ops[1].Op)
	assert.Equal("remove", ops[2].Op)
}

func TestPatchUserHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPatchUserBuilder(pn).ID("alice").Replace("/status", "inactive").IfMatchETag("1")
	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(entityPatchContentType, headers["Content-Type"])
	assert.Equal("1", headers["If-Match"])
}

func TestPatchUserHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newPatchUserBuilder(pn)
	assert.Equal("PATCH", o.opts.httpMethod())
	assert.Equal(PNPatchDataSyncUserOperation, o.opts.operationType())
}

func TestPatchUserValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPatchUserBuilder(pn).Replace("/status", "inactive")
	assert.Contains(o.opts.validate().Error(), StrMissingUserID)

	o2 := newPatchUserBuilder(pn).ID("alice")
	assert.Contains(o2.opts.validate().Error(), StrMissingPatchOperations)

	pn.Config.SubscribeKey = ""
	o3 := newPatchUserBuilder(pn).ID("alice").Replace("/status", "inactive")
	assert.Contains(o3.opts.validate().Error(), StrMissingSubKey)
}

func TestPatchUserExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.PatchUser().ID("alice").Execute()
	assert.Contains(err.Error(), StrMissingPatchOperations)
}
