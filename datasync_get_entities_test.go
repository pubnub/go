package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetEntitiesDefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetEntitiesBuilder(pn)
	assert.Equal(entitiesDefaultLimit, o.opts.Limit)
}

func TestGetEntitiesBuildPathQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetEntitiesBuilder(pn).
		EntityClass("vehicle").
		EntityClassVersion(1).
		EntityClassLevel(PNEntityClassLevelSubKey).
		Cursor("TjIw").
		Limit(10).
		Filter("status == 'active'").
		FilterAdvanced("status == 'active' AND year > 2020").
		Sort([]string{"-createdAt", "+id"}).
		QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/entities", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("vehicle", query.Get("entity_class"))
	assert.Equal("1", query.Get("entity_class_version"))
	assert.Equal("SubKey", query.Get("entity_class_level"))
	assert.Equal("TjIw", query.Get("cursor"))
	assert.Equal(strconv.Itoa(10), query.Get("limit"))
	assert.Equal("status == 'active'", query.Get("filter"))
	assert.Equal("-createdAt,+id", query.Get("sort"))
	assert.Equal("v1", query.Get("q1"))
	assert.NotEmpty(query.Get("filter_advanced"))
}

func TestGetEntitiesVersionOmittedWhenUnset(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetEntitiesBuilder(pn).EntityClass("vehicle")
	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("", query.Get("entity_class_version"))
	assert.Equal("", query.Get("entity_class_level"))
}

func TestGetEntitiesHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGetEntitiesBuilder(pn)
	assert.Equal("GET", o.opts.httpMethod())
	assert.Equal(PNGetEntitiesOperation, o.opts.operationType())
}

func TestGetEntitiesValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetEntitiesBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingEntityClass)

	pn.Config.SubscribeKey = ""
	o2 := newGetEntitiesBuilder(pn).EntityClass("vehicle")
	assert.Contains(o2.opts.validate().Error(), StrMissingSubKey)
}

func TestGetEntitiesResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"data":[{"id":"i789","status":"active","entityClass":"user","entityClassVersion":1,"entityClassLevel":"SubKey","eTag":"1","payload":{"custom":"fields"}}],"links":{"self":"/self","next":"/next","prev":null},"meta":{"has_next":true,"has_prev":false,"next_cursor":"TjIw","prev_cursor":null,"limit":20}}`)

	r, _, err := newPNEntitiesResponse(jsonBytes, newGetEntitiesOpts(pn, pn.ctx), StatusResponse{})
	assert.Nil(err)
	assert.Len(r.Data, 1)
	assert.Equal("i789", r.Data[0].ID)
	assert.Equal(PNEntityClassLevelSubKey, r.Data[0].EntityClassLevel)
	assert.NotNil(r.Meta)
	assert.True(r.Meta.HasNext)
	assert.Equal("TjIw", r.Meta.NextCursor)
	assert.Equal(20, r.Meta.Limit)
	assert.NotNil(r.Links)
	assert.Equal("/next", r.Links.Next)
}

func TestGetEntitiesExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.GetEntities().Execute()
	assert.Contains(err.Error(), StrMissingEntityClass)
}
