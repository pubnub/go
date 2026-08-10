package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserBuildPathQueryBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newCreateUserBuilder(pn)
	o.ID("alice").
		EntityClass("User").
		EntityClassVersion(1).
		EntityClassLevel("SubKey").
		Status("active").
		Payload(map[string]interface{}{"name": "Alice", "email": "alice@example.com"}).
		QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/users", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))

	body, err := o.opts.buildBody()
	assert.Nil(err)

	var parsed createUserBody
	assert.Nil(json.Unmarshal(body, &parsed))
	assert.Equal("alice", parsed.Data.ID)
	assert.Equal("User", parsed.Data.EntityClass)
	assert.Equal(1, parsed.Data.EntityClassVersion)
	assert.Equal("SubKey", parsed.Data.EntityClassLevel)
	assert.Equal("active", parsed.Data.Status)
	assert.Equal("Alice", parsed.Data.Payload["name"])
}

func TestCreateUserHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newCreateUserBuilder(pn).EntityClassVersion(1)

	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(userContentType, headers["Content-Type"])
}

func TestCreateUserHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newCreateUserBuilder(pn)
	assert.Equal("POST", o.opts.httpMethod())
	assert.Equal(PNCreateDataSyncUserOperation, o.opts.operationType())
}

func TestCreateUserValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Invalid version
	o := newCreateUserBuilder(pn).EntityClassVersion(0)
	assert.Contains(o.opts.validate().Error(), StrInvalidEntityClassVersion)

	// Missing sub key
	pn.Config.SubscribeKey = ""
	o2 := newCreateUserBuilder(pn).EntityClassVersion(1)
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestCreateUserResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"status":201,"data":{"id":"alice","entityClass":"User","entityClassVersion":1,"entityClassLevel":"SubKey","status":"active","payload":{"name":"Alice"},"createdAt":"2026-03-20T12:00:00.000Z","updatedAt":"2026-03-20T12:00:00.000Z","eTag":"1"}}`)

	r, _, err := newPNEntityResponse(jsonBytes, StatusResponse{}, PNCreateDataSyncUserOperation, pn)
	assert.Nil(err)
	assert.Equal(201, r.Status)
	assert.Equal("alice", r.Data.ID)
	assert.Equal("User", r.Data.EntityClass)
	assert.Equal(1, r.Data.EntityClassVersion)
	assert.Equal("1", r.Data.ETag)
	assert.Equal("Alice", r.Data.Payload["name"])
}

func TestCreateUserExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	_, _, err := pn.DataSync.CreateUser().EntityClassVersion(1).Execute()
	assert.Contains(err.Error(), StrMissingSubKey)
}
