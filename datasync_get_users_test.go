package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersDefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetUsersBuilder(pn)
	assert.Equal(entitiesDefaultLimit, o.opts.Limit)
}

func TestGetUsersBuildPathQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetUsersBuilder(pn).
		EntityClass("User").
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
		fmt.Sprintf("/v1/datasync/subkeys/%s/users", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("User", query.Get("entity_class"))
	assert.Equal("1", query.Get("entity_class_version"))
	assert.Equal("SubKey", query.Get("entity_class_level"))
	assert.Equal("TjIw", query.Get("cursor"))
	assert.Equal(strconv.Itoa(10), query.Get("limit"))
	assert.Equal("status == 'active'", query.Get("filter_fast"))
	assert.Equal("", query.Get("filter"))
	assert.Equal("-createdAt,+id", query.Get("sort"))
	assert.Equal("v1", query.Get("q1"))
}

func TestGetUsersFilterQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetUsersBuilder(pn).Filter("status == 'active' AND location == 'Pune'")
	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("status == 'active' AND location == 'Pune'", query.Get("filter"))
	assert.Equal("", query.Get("filter_fast"))
}

func TestGetUsersOptionalFiltersOmittedWhenUnset(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetUsersBuilder(pn)
	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("", query.Get("entity_class"))
	assert.Equal("", query.Get("entity_class_version"))
	assert.Equal("", query.Get("entity_class_level"))
	assert.Equal(strconv.Itoa(entitiesDefaultLimit), query.Get("limit"))
}

func TestGetUsersHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGetUsersBuilder(pn)
	assert.Equal("GET", o.opts.httpMethod())
	assert.Equal(PNGetDataSyncUsersOperation, o.opts.operationType())
}

func TestGetUsersValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	o := newGetUsersBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingSubKey)
}

func TestGetUsersValidateRejectsBothFilters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetUsersBuilder(pn).
		FilterFast("status == 'active'").
		Filter("status == 'active' AND location == 'Pune'")
	assert.Contains(o.opts.validate().Error(), StrExclusiveDataSyncFilter)
}

func TestGetUsersResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"data":[{"id":"alice","status":"active","entityClass":"User","entityClassVersion":1,"entityClassLevel":"SubKey","eTag":"1","payload":{"name":"Alice"}}],"links":{"self":"/self","next":"/next","prev":null},"meta":{"has_next":true,"has_prev":false,"next_cursor":"TjIw","prev_cursor":null,"limit":20}}`)

	r, _, err := newPNUsersResponse(jsonBytes, newGetUsersOpts(pn, pn.ctx), StatusResponse{})
	assert.Nil(err)
	assert.Len(r.Data, 1)
	assert.Equal("alice", r.Data[0].ID)
	assert.Equal(PNEntityClassLevelSubKey, r.Data[0].EntityClassLevel)
	assert.NotNil(r.Meta)
	assert.True(r.Meta.HasNext)
	assert.Equal("TjIw", r.Meta.NextCursor)
	assert.Equal(20, r.Meta.Limit)
	assert.NotNil(r.Links)
	assert.Equal("/next", r.Links.Next)
}

func TestGetUsersExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	_, _, err := pn.DataSync.GetUsers().Execute()
	assert.Contains(err.Error(), StrMissingSubKey)
}

func TestGetUsersExecuteBothFiltersError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.GetUsers().
		FilterFast("status == 'active'").
		Filter("status == 'active' AND location == 'Pune'").
		Execute()
	assert.Contains(err.Error(), StrExclusiveDataSyncFilter)
}
