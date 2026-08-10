package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestRemoveChannelBuildPathQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newDeleteChannelBuilder(pn).ID("general").IfMatchETag("1").QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/channels/general", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))

	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal("1", headers["If-Match"])
}

func TestRemoveChannelHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newDeleteChannelBuilder(pn)
	assert.Equal("DELETE", o.opts.httpMethod())
	assert.Equal(PNRemoveDataSyncChannelOperation, o.opts.operationType())
}

func TestRemoveChannelValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newDeleteChannelBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingChannelID)

	pn.Config.SubscribeKey = ""
	o2 := newDeleteChannelBuilder(pn).ID("general")
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestRemoveChannelEmptyBodyResponse(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	r, _, err := newPNRemoveChannelResponse([]byte{}, newDeleteChannelOpts(pn, pn.ctx), StatusResponse{})
	assert.Nil(err)
	assert.NotNil(r)
}

func TestRemoveChannelExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.RemoveChannel().Execute()
	assert.Contains(err.Error(), StrMissingChannelID)
}
