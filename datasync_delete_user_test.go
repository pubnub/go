package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestRemoveUserBuildPathQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newDeleteUserBuilder(pn).ID("alice").IfMatchETag("1").QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/users/alice", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))

	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal("1", headers["If-Match"])
}

func TestRemoveUserHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newDeleteUserBuilder(pn)
	assert.Equal("DELETE", o.opts.httpMethod())
	assert.Equal(PNRemoveDataSyncUserOperation, o.opts.operationType())
}

func TestRemoveUserValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newDeleteUserBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingUserID)

	pn.Config.SubscribeKey = ""
	o2 := newDeleteUserBuilder(pn).ID("alice")
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestRemoveUserEmptyBodyResponse(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	r, _, err := newPNRemoveUserResponse([]byte{}, newDeleteUserOpts(pn, pn.ctx), StatusResponse{})
	assert.Nil(err)
	assert.NotNil(r)
}

func TestRemoveUserExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.RemoveUser().Execute()
	assert.Contains(err.Error(), StrMissingUserID)
}
