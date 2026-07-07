package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestPatchEntityBuildPathBodyHelpers(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPatchEntityBuilder(pn).
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

func TestPatchEntityMoveCopyTest(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPatchEntityBuilder(pn).
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

func TestPatchEntityHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPatchEntityBuilder(pn).ID("entity-abc").Replace("/status", "inactive")
	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(entityPatchContentType, headers["Content-Type"])
	assert.NotEmpty(headers["Idempotency-Key"])

	o.IfMatch("BfklQ...").IdempotencyKey("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
	headers, err = o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal("BfklQ...", headers["If-Match"])
	assert.Equal("a1b2c3d4-e5f6-7890-abcd-ef1234567890", headers["Idempotency-Key"])
}

func TestPatchEntityHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newPatchEntityBuilder(pn)
	assert.Equal("PATCH", o.opts.httpMethod())
	assert.Equal(PNPatchEntityOperation, o.opts.operationType())
}

func TestPatchEntityValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPatchEntityBuilder(pn).Replace("/status", "inactive")
	assert.Contains(o.opts.validate().Error(), StrMissingEntityID)

	o2 := newPatchEntityBuilder(pn).ID("entity-abc")
	assert.Contains(o2.opts.validate().Error(), StrMissingPatchOperations)

	pn.Config.SubscribeKey = ""
	o3 := newPatchEntityBuilder(pn).ID("entity-abc").Replace("/status", "x")
	assert.Contains(o3.opts.validate().Error(), StrMissingSubKey)
}

func TestPatchEntityResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"status":200,"data":{"id":"entity-abc","entityClass":"vehicle","entityClassVersion":1,"status":"inactive","payload":{"mileage":42000},"eTag":"Dpqr3..."}}`)

	r, _, err := newPNEntityResponse(jsonBytes, StatusResponse{}, PNPatchEntityOperation, pn)
	assert.Nil(err)
	assert.Equal("inactive", r.Data.Status)
	assert.Equal(float64(42000), r.Data.Payload["mileage"])
}

func TestPatchEntityExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.PatchEntity().ID("entity-abc").Execute()
	assert.Contains(err.Error(), StrMissingPatchOperations)
}
