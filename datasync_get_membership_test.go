package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetMembershipBuildPathQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetMembershipBuilder(pn).ID("m-123").QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/memberships/m-123", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))
}

func TestGetMembershipHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGetMembershipBuilder(pn)
	assert.Equal("GET", o.opts.httpMethod())
	assert.Equal(PNGetDataSyncMembershipOperation, o.opts.operationType())
}

func TestGetMembershipValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetMembershipBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingMembershipID)

	pn.Config.SubscribeKey = ""
	o2 := newGetMembershipBuilder(pn).ID("m-123")
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestGetMembershipExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.GetMembership().Execute()
	assert.Contains(err.Error(), StrMissingMembershipID)
}
