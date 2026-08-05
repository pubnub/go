package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCreateEntityBuildPathQueryBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newCreateEntityBuilder(pn)
	o.ID("entity-abc").
		EntityClass("vehicle").
		EntityClassVersion(1).
		Status("active").
		Payload(map[string]interface{}{"make": "Toyota", "year": 2025}).
		QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/entities", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))

	body, err := o.opts.buildBody()
	assert.Nil(err)

	var parsed createEntityBody
	assert.Nil(json.Unmarshal(body, &parsed))
	assert.Equal("entity-abc", parsed.Data.ID)
	assert.Equal("vehicle", parsed.Data.EntityClass)
	assert.Equal(1, parsed.Data.EntityClassVersion)
	assert.Equal("active", parsed.Data.Status)
	assert.Equal("Toyota", parsed.Data.Payload["make"])
}

func TestCreateEntityHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newCreateEntityBuilder(pn).EntityClass("vehicle").EntityClassVersion(1)

	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(entityContentType, headers["Content-Type"])
	_, hasIdempotencyKey := headers["Idempotency-Key"]
	assert.False(hasIdempotencyKey)
}

func TestCreateEntityHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newCreateEntityBuilder(pn)
	assert.Equal("POST", o.opts.httpMethod())
	assert.Equal(PNCreateEntityOperation, o.opts.operationType())
}

func TestCreateEntityValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Missing entity class
	o := newCreateEntityBuilder(pn).EntityClassVersion(1)
	assert.Contains(o.opts.validate().Error(), StrMissingEntityClass)

	// Invalid version
	o2 := newCreateEntityBuilder(pn).EntityClass("vehicle").EntityClassVersion(0)
	assert.Contains(o2.opts.validate().Error(), StrInvalidEntityClassVersion)

	// Missing sub key
	pn.Config.SubscribeKey = ""
	o3 := newCreateEntityBuilder(pn).EntityClass("vehicle").EntityClassVersion(1)
	assert.Contains(o3.opts.validate().Error(), StrMissingSubKey)
}

func TestCreateEntityResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"status":201,"data":{"id":"entity-abc","entityClass":"vehicle","entityClassVersion":1,"status":"active","payload":{"make":"Toyota"},"createdAt":"2026-03-20T12:00:00.000Z","updatedAt":"2026-03-20T12:00:00.000Z","eTag":"BfklQ...","expiresAt":null}}`)

	r, _, err := newPNEntityResponse(jsonBytes, StatusResponse{}, PNCreateEntityOperation, pn)
	assert.Nil(err)
	assert.Equal(201, r.Status)
	assert.Equal("entity-abc", r.Data.ID)
	assert.Equal("vehicle", r.Data.EntityClass)
	assert.Equal(1, r.Data.EntityClassVersion)
	assert.Equal("BfklQ...", r.Data.ETag)
	assert.Equal("Toyota", r.Data.Payload["make"])
}

func TestCreateEntityExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	_, _, err := pn.DataSync.CreateEntity().EntityClass("vehicle").EntityClassVersion(1).Execute()
	assert.Contains(err.Error(), StrMissingSubKey)
}
