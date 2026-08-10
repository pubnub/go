package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetChannelBuildPathQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetChannelBuilder(pn).ID("general").QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/channels/general", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))
}

func TestGetChannelHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGetChannelBuilder(pn)
	assert.Equal("GET", o.opts.httpMethod())
	assert.Equal(PNGetDataSyncChannelOperation, o.opts.operationType())
}

func TestGetChannelValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetChannelBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingChannelID)

	pn.Config.SubscribeKey = ""
	o2 := newGetChannelBuilder(pn).ID("general")
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestGetChannelExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.GetChannel().Execute()
	assert.Contains(err.Error(), StrMissingChannelID)
}
