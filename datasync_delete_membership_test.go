package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestRemoveMembershipBuildPathQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newDeleteMembershipBuilder(pn).ID("m-123").IfMatchETag("1").QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/memberships/m-123", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))

	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal("1", headers["If-Match"])
}

func TestRemoveMembershipHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newDeleteMembershipBuilder(pn)
	assert.Equal("DELETE", o.opts.httpMethod())
	assert.Equal(PNRemoveDataSyncMembershipOperation, o.opts.operationType())
}

func TestRemoveMembershipValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newDeleteMembershipBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingMembershipID)

	pn.Config.SubscribeKey = ""
	o2 := newDeleteMembershipBuilder(pn).ID("m-123")
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestRemoveMembershipEmptyBodyResponse(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	r, _, err := newPNRemoveDataSyncMembershipResponse([]byte{}, newDeleteMembershipOpts(pn, pn.ctx), StatusResponse{})
	assert.Nil(err)
	assert.NotNil(r)
}

func TestRemoveMembershipExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.RemoveMembership().Execute()
	assert.Contains(err.Error(), StrMissingMembershipID)
}
