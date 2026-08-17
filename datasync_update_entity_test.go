package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestUpdateEntityBuildPathBodyHelpers(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newUpdateEntityBuilder(pn).
		ID("entity-abc").
		Replace("/status", "inactive").
		Add("/payload/mileage", 42000).
		Remove("/payload/owner/license")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/entities/entity-abc", pn.Config.SubscribeKey),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	var ops []PNJSONPatchOperation
	assert.Nil(json.Unmarshal(body, &ops))
	assert.Len(ops, 3)
	assert.Equal("replace", ops[0].Op)
	assert.Equal("/status", ops[0].Path)
	assert.Equal("inactive", ops[0].Value)
	assert.Equal("add", ops[1].Op)
	assert.Equal("remove", ops[2].Op)
	assert.Equal("/payload/owner/license", ops[2].Path)
}

func TestUpdateEntityMoveCopyTest(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newUpdateEntityBuilder(pn).
		ID("entity-abc").
		Move("/payload/a", "/payload/b").
		Copy("/payload/c", "/payload/d").
		Test("/status", "active")

	body, err := o.opts.buildBody()
	assert.Nil(err)

	var ops []PNJSONPatchOperation
	assert.Nil(json.Unmarshal(body, &ops))
	assert.Len(ops, 3)
	assert.Equal("move", ops[0].Op)
	assert.Equal("/payload/a", ops[0].From)
	assert.Equal("/payload/b", ops[0].Path)
	assert.Equal("copy", ops[1].Op)
	assert.Equal("test", ops[2].Op)
}

func TestUpdateEntityHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newUpdateEntityBuilder(pn).ID("entity-abc").Replace("/status", "inactive")
	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(entityPatchContentType, headers["Content-Type"])
	_, hasIdempotencyKey := headers["Idempotency-Key"]
	assert.False(hasIdempotencyKey)
	_, hasIfMatch := headers["If-Match"]
	assert.False(hasIfMatch)

	o.IfMatchETag("BfklQ...")
	headers, err = o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal("BfklQ...", headers["If-Match"])
}

func TestUpdateEntityHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newUpdateEntityBuilder(pn)
	assert.Equal("PATCH", o.opts.httpMethod())
	assert.Equal(PNUpdateEntityOperation, o.opts.operationType())
}

func TestUpdateEntityValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newUpdateEntityBuilder(pn).Replace("/status", "inactive")
	assert.Contains(o.opts.validate().Error(), StrMissingEntityID)

	o2 := newUpdateEntityBuilder(pn).ID("entity-abc")
	assert.Contains(o2.opts.validate().Error(), StrMissingPatchOperations)

	pn.Config.SubscribeKey = ""
	o3 := newUpdateEntityBuilder(pn).ID("entity-abc").Replace("/status", "x")
	assert.Contains(o3.opts.validate().Error(), StrMissingSubKey)
}

func TestUpdateEntityResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"status":200,"data":{"id":"entity-abc","entityClass":"vehicle","entityClassVersion":1,"status":"inactive","payload":{"mileage":42000},"eTag":"Dpqr3..."}}`)

	r, _, err := newPNEntityResponse(jsonBytes, StatusResponse{}, PNUpdateEntityOperation, pn)
	assert.Nil(err)
	assert.Equal("inactive", r.Data.Status)
	assert.Equal(float64(42000), r.Data.Payload["mileage"])
}

func TestUpdateEntityExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.UpdateEntity().ID("entity-abc").Execute()
	assert.Contains(err.Error(), StrMissingPatchOperations)
}
