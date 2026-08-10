package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestPatchChannelBuildPathQueryBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPatchChannelBuilder(pn).
		ID("general").
		Replace("/payload/topic", "https://example.com/general").
		Add("/payload/phone", "+1-555-0100").
		Remove("/payload/isActive").
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

	var ops []PNJSONPatchOperation
	assert.Nil(json.Unmarshal(body, &ops))
	assert.Len(ops, 3)
	assert.Equal("replace", ops[0].Op)
	assert.Equal("/payload/topic", ops[0].Path)
	assert.Equal("add", ops[1].Op)
	assert.Equal("remove", ops[2].Op)
}

func TestPatchChannelHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPatchChannelBuilder(pn).ID("general").Replace("/status", "inactive").IfMatchETag("1")
	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(entityPatchContentType, headers["Content-Type"])
	assert.Equal("1", headers["If-Match"])
}

func TestPatchChannelHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newPatchChannelBuilder(pn)
	assert.Equal("PATCH", o.opts.httpMethod())
	assert.Equal(PNPatchDataSyncChannelOperation, o.opts.operationType())
}

func TestPatchChannelValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPatchChannelBuilder(pn).Replace("/status", "inactive")
	assert.Contains(o.opts.validate().Error(), StrMissingChannelID)

	o2 := newPatchChannelBuilder(pn).ID("general")
	assert.Contains(o2.opts.validate().Error(), StrMissingPatchOperations)

	pn.Config.SubscribeKey = ""
	o3 := newPatchChannelBuilder(pn).ID("general").Replace("/status", "inactive")
	assert.Contains(o3.opts.validate().Error(), StrMissingSubKey)
}

func TestPatchChannelExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.PatchChannel().ID("general").Execute()
	assert.Contains(err.Error(), StrMissingPatchOperations)
}
