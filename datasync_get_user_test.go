package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetUserBuildPathQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetUserBuilder(pn).ID("alice").QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/users/alice", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))
}

func TestGetUserHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGetUserBuilder(pn)
	assert.Equal("GET", o.opts.httpMethod())
	assert.Equal(PNGetDataSyncUserOperation, o.opts.operationType())
}

func TestGetUserValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetUserBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingUserID)

	pn.Config.SubscribeKey = ""
	o2 := newGetUserBuilder(pn).ID("alice")
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestGetUserExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.GetUser().Execute()
	assert.Contains(err.Error(), StrMissingUserID)
}
