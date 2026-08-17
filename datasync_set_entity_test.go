package pubnub

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestSetEntityBuildPathBody(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newSetEntityBuilder(pn).
		ID("entity-abc").
		EntityClassVersion(2).
		Status("active").
		Payload(map[string]interface{}{"color": "blue"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/entities/entity-abc", pn.Config.SubscribeKey),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	var parsed setEntityBody
	assert.Nil(json.Unmarshal(body, &parsed))
	assert.Equal(2, parsed.Data.EntityClassVersion)
	assert.Equal("active", parsed.Data.Status)
	assert.Equal("blue", parsed.Data.Payload["color"])
	// entityClass must NOT be present in the body (immutable)
	assert.False(strings.Contains(string(body), "entityClass\":"))
}

func TestSetEntityHeaders(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newSetEntityBuilder(pn).ID("entity-abc").EntityClassVersion(2)
	headers, err := o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal(entityContentType, headers["Content-Type"])
	_, hasIfMatch := headers["If-Match"]
	assert.False(hasIfMatch)

	o.IfMatchETag("BfklQ...")
	headers, err = o.opts.buildHeaders()
	assert.Nil(err)
	assert.Equal("BfklQ...", headers["If-Match"])
}

func TestSetEntityHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newSetEntityBuilder(pn)
	assert.Equal("PUT", o.opts.httpMethod())
	assert.Equal(PNSetEntityOperation, o.opts.operationType())
}

func TestSetEntityValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newSetEntityBuilder(pn).EntityClassVersion(2)
	assert.Contains(o.opts.validate().Error(), StrMissingEntityID)

	o2 := newSetEntityBuilder(pn).ID("entity-abc").EntityClassVersion(0)
	assert.Contains(o2.opts.validate().Error(), StrInvalidEntityClassVersion)

	pn.Config.SubscribeKey = ""
	o3 := newSetEntityBuilder(pn).ID("entity-abc").EntityClassVersion(2)
	assert.Contains(o3.opts.validate().Error(), StrMissingSubKey)
}

func TestSetEntityResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"status":200,"data":{"id":"entity-abc","entityClass":"vehicle","entityClassVersion":2,"status":"active","eTag":"Cxyz1..."}}`)

	r, _, err := newPNEntityResponse(jsonBytes, StatusResponse{}, PNSetEntityOperation, pn)
	assert.Nil(err)
	assert.Equal(2, r.Data.EntityClassVersion)
	assert.Equal("Cxyz1...", r.Data.ETag)
}

func TestSetEntityExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.SetEntity().EntityClassVersion(2).Execute()
	assert.Contains(err.Error(), StrMissingEntityID)
}
