package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestSetUserBuildPathQueryBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newSetUserBuilder(pn).
		ID("alice").
		EntityClassVersion(1).
		Status("active").
		Payload(map[string]interface{}{"name": "Alice"}).
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

	var parsed setUserBody
	assert.Nil(json.Unmarshal(body, &parsed))
	assert.Equal(1, parsed.Data.EntityClassVersion)
	assert.Equal("active", parsed.Data.Status)
	assert.Equal("Alice", parsed.Data.Payload["name"])
}

func TestSetUserHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newSetUserBuilder(pn).ID("alice").EntityClassVersion(1).IfMatchETag("1")
	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(userContentType, headers["Content-Type"])
	assert.Equal("1", headers["If-Match"])
}

func TestSetUserHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newSetUserBuilder(pn)
	assert.Equal("PUT", o.opts.httpMethod())
	assert.Equal(PNSetDataSyncUserOperation, o.opts.operationType())
}

func TestSetUserValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newSetUserBuilder(pn).EntityClassVersion(1)
	assert.Contains(o.opts.validate().Error(), StrMissingUserID)

	o2 := newSetUserBuilder(pn).ID("alice").EntityClassVersion(0)
	assert.Contains(o2.opts.validate().Error(), StrInvalidEntityClassVersion)

	pn.Config.SubscribeKey = ""
	o3 := newSetUserBuilder(pn).ID("alice").EntityClassVersion(1)
	assert.Contains(o3.opts.validate().Error(), StrMissingSubKey)
}

func TestSetUserExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.SetUser().EntityClassVersion(1).Execute()
	assert.Contains(err.Error(), StrMissingUserID)
}
