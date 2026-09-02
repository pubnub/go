package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetChannelsDefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetChannelsBuilder(pn)
	assert.Equal(entitiesDefaultLimit, o.opts.Limit)
}

func TestGetChannelsBuildPathQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetChannelsBuilder(pn).
		EntityClass("Channel").
		EntityClassVersion(1).
		EntityClassLevel(PNEntityClassLevelSubKey).
		Cursor("TjIw").
		Limit(10).
		FilterFast("status == 'active'").
		Sort([]string{"-createdAt", "+id"}).
		QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/channels", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("Channel", query.Get("entity_class"))
	assert.Equal("1", query.Get("entity_class_version"))
	assert.Equal("SubKey", query.Get("entity_class_level"))
	assert.Equal("TjIw", query.Get("cursor"))
	assert.Equal(strconv.Itoa(10), query.Get("limit"))
	assert.Equal("status == 'active'", query.Get("filter_fast"))
	assert.Equal("", query.Get("filter"))
	assert.Equal("-createdAt,+id", query.Get("sort"))
	assert.Equal("v1", query.Get("q1"))
}

func TestGetChannelsFilterQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetChannelsBuilder(pn).Filter("status == 'active' AND location == 'Pune'")
	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("status == 'active' AND location == 'Pune'", query.Get("filter"))
	assert.Equal("", query.Get("filter_fast"))
}

func TestGetChannelsOptionalFiltersOmittedWhenUnset(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetChannelsBuilder(pn)
	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("", query.Get("entity_class"))
	assert.Equal("", query.Get("entity_class_version"))
	assert.Equal("", query.Get("entity_class_level"))
	assert.Equal(strconv.Itoa(entitiesDefaultLimit), query.Get("limit"))
}

func TestGetChannelsHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGetChannelsBuilder(pn)
	assert.Equal("GET", o.opts.httpMethod())
	assert.Equal(PNGetDataSyncChannelsOperation, o.opts.operationType())
}

func TestGetChannelsValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	o := newGetChannelsBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingSubKey)
}

func TestGetChannelsValidateRejectsBothFilters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetChannelsBuilder(pn).
		FilterFast("status == 'active'").
		Filter("status == 'active' AND location == 'Pune'")
	assert.Contains(o.opts.validate().Error(), StrExclusiveDataSyncFilter)
}

func TestGetChannelsResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"data":[{"id":"general","status":"active","entityClass":"Channel","entityClassVersion":1,"entityClassLevel":"SubKey","eTag":"1","payload":{"name":"General"}}],"links":{"self":"/self","next":"/next","prev":null},"meta":{"has_next":true,"has_prev":false,"next_cursor":"TjIw","prev_cursor":null,"limit":20}}`)

	r, _, err := newPNDataSyncChannelsResponse(jsonBytes, newGetChannelsOpts(pn, pn.ctx), StatusResponse{})
	assert.Nil(err)
	assert.Len(r.Data, 1)
	assert.Equal("general", r.Data[0].ID)
	assert.Equal(PNEntityClassLevelSubKey, r.Data[0].EntityClassLevel)
	assert.NotNil(r.Meta)
	assert.True(r.Meta.HasNext)
	assert.Equal("TjIw", r.Meta.NextCursor)
	assert.Equal(20, r.Meta.Limit)
	assert.NotNil(r.Links)
	assert.Equal("/next", r.Links.Next)
}

func TestGetChannelsExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	_, _, err := pn.DataSync.GetChannels().Execute()
	assert.Contains(err.Error(), StrMissingSubKey)
}

func TestGetChannelsExecuteBothFiltersError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.GetChannels().
		FilterFast("status == 'active'").
		Filter("status == 'active' AND location == 'Pune'").
		Execute()
	assert.Contains(err.Error(), StrExclusiveDataSyncFilter)
}
