package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetEntityBuildPathQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetEntityBuilder(pn).ID("entity-abc").QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/entities/entity-abc", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", query.Get("q1"))

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}

func TestGetEntityHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGetEntityBuilder(pn)
	assert.Equal("GET", o.opts.httpMethod())
	assert.Equal(PNGetEntityOperation, o.opts.operationType())
}

func TestGetEntityValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetEntityBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingEntityID)

	pn.Config.SubscribeKey = ""
	o2 := newGetEntityBuilder(pn).ID("entity-abc")
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestGetEntityResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"status":200,"data":{"id":"entity-abc","entityClass":"vehicle","entityClassVersion":1,"status":"active","payload":{"make":"Toyota"},"eTag":"BfklQ..."}}`)

	r, _, err := newPNEntityResponse(jsonBytes, StatusResponse{}, PNGetEntityOperation, pn)
	assert.Nil(err)
	assert.Equal("entity-abc", r.Data.ID)
	assert.Equal("Toyota", r.Data.Payload["make"])
}

func TestGetEntityExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.GetEntity().Execute()
	assert.Contains(err.Error(), StrMissingEntityID)
}
