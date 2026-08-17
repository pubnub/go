package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestSetChannelBuildPathQueryBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newSetChannelBuilder(pn).
		ID("general").
		EntityClassVersion(1).
		Status("active").
		Payload(map[string]interface{}{"name": "General"}).
		IfMatchETag("1").
		QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/channels/general", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))

	body, err := o.opts.buildBody()
	assert.Nil(err)

	var parsed setChannelBody
	assert.Nil(json.Unmarshal(body, &parsed))
	assert.Equal(1, parsed.Data.EntityClassVersion)
	assert.Equal("active", parsed.Data.Status)
	assert.Equal("General", parsed.Data.Payload["name"])
}

func TestSetChannelHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newSetChannelBuilder(pn).ID("general").EntityClassVersion(1).IfMatchETag("1")
	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(channelContentType, headers["Content-Type"])
	assert.Equal("1", headers["If-Match"])
}

func TestSetChannelHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newSetChannelBuilder(pn)
	assert.Equal("PUT", o.opts.httpMethod())
	assert.Equal(PNSetDataSyncChannelOperation, o.opts.operationType())
}

func TestSetChannelValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newSetChannelBuilder(pn).EntityClassVersion(1)
	assert.Contains(o.opts.validate().Error(), StrMissingChannelID)

	o2 := newSetChannelBuilder(pn).ID("general").EntityClassVersion(0)
	assert.Contains(o2.opts.validate().Error(), StrInvalidEntityClassVersion)

	pn.Config.SubscribeKey = ""
	o3 := newSetChannelBuilder(pn).ID("general").EntityClassVersion(1)
	assert.Contains(o3.opts.validate().Error(), StrMissingSubKey)
}

func TestSetChannelExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.SetChannel().EntityClassVersion(1).Execute()
	assert.Contains(err.Error(), StrMissingChannelID)
}
