package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCreateChannelBuildPathQueryBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newCreateChannelBuilder(pn)
	o.ID("general").
		EntityClass("Channel").
		EntityClassVersion(1).
		EntityClassLevel(PNEntityClassLevelSubKey).
		Status("active").
		Payload(map[string]interface{}{"name": "General", "email": "team-chat"}).
		QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/channels", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))

	body, err := o.opts.buildBody()
	assert.Nil(err)

	var parsed createChannelBody
	assert.Nil(json.Unmarshal(body, &parsed))
	assert.Equal("general", parsed.Data.ID)
	assert.Equal("Channel", parsed.Data.EntityClass)
	assert.Equal(1, parsed.Data.EntityClassVersion)
	assert.Equal(PNEntityClassLevelSubKey, parsed.Data.EntityClassLevel)
	assert.Equal("active", parsed.Data.Status)
	assert.Equal("General", parsed.Data.Payload["name"])
}

func TestCreateChannelHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newCreateChannelBuilder(pn).EntityClassVersion(1)

	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(channelContentType, headers["Content-Type"])
}

func TestCreateChannelHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newCreateChannelBuilder(pn)
	assert.Equal("POST", o.opts.httpMethod())
	assert.Equal(PNCreateDataSyncChannelOperation, o.opts.operationType())
}

func TestCreateChannelValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Invalid version
	o := newCreateChannelBuilder(pn).EntityClassVersion(0)
	assert.Contains(o.opts.validate().Error(), StrInvalidEntityClassVersion)

	// Missing sub key
	pn.Config.SubscribeKey = ""
	o2 := newCreateChannelBuilder(pn).EntityClassVersion(1)
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestCreateChannelResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"status":201,"data":{"id":"general","entityClass":"Channel","entityClassVersion":1,"entityClassLevel":"SubKey","status":"active","payload":{"name":"General"},"createdAt":"2026-03-20T12:00:00.000Z","updatedAt":"2026-03-20T12:00:00.000Z","eTag":"1"}}`)

	r, _, err := newPNEntityResponse(jsonBytes, StatusResponse{}, PNCreateDataSyncChannelOperation, pn)
	assert.Nil(err)
	assert.Equal(201, r.Status)
	assert.Equal("general", r.Data.ID)
	assert.Equal("Channel", r.Data.EntityClass)
	assert.Equal(1, r.Data.EntityClassVersion)
	assert.Equal("1", r.Data.ETag)
	assert.Equal("General", r.Data.Payload["name"])
}

func TestCreateChannelExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	_, _, err := pn.DataSync.CreateChannel().EntityClassVersion(1).Execute()
	assert.Contains(err.Error(), StrMissingSubKey)
}
