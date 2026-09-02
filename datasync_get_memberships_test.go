package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetMembershipsDefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetMembershipsBuilder(pn)
	assert.Equal(entitiesDefaultLimit, o.opts.Limit)
}

func TestGetMembershipsBuildPathQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetMembershipsBuilder(pn).
		UserID("alice").
		ChannelID("general").
		RelationshipClassVersion(1).
		Cursor("TjIw").
		Limit(10).
		FilterFast("status == 'active'").
		Sort([]string{"-createdAt", "+id"}).
		QueryParam(map[string]string{"q1": "v1"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/datasync/subkeys/%s/memberships", pn.Config.SubscribeKey),
		path, []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("alice", query.Get("user_id"))
	assert.Equal("general", query.Get("channel_id"))
	assert.Equal("1", query.Get("relationship_class_version"))
	assert.Equal("TjIw", query.Get("cursor"))
	assert.Equal(strconv.Itoa(10), query.Get("limit"))
	assert.Equal("status == 'active'", query.Get("filter_fast"))
	assert.Equal("", query.Get("filter"))
	assert.Equal("-createdAt,+id", query.Get("sort"))
	assert.Equal("v1", query.Get("q1"))
}

func TestGetMembershipsFilterQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetMembershipsBuilder(pn).Filter("status == 'active' AND role == 'admin'")
	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("status == 'active' AND role == 'admin'", query.Get("filter"))
	assert.Equal("", query.Get("filter_fast"))
}

func TestGetMembershipsOptionalFiltersOmittedWhenUnset(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetMembershipsBuilder(pn)
	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("", query.Get("user_id"))
	assert.Equal("", query.Get("channel_id"))
	assert.Equal("", query.Get("relationship_class_version"))
	assert.Equal(strconv.Itoa(entitiesDefaultLimit), query.Get("limit"))
}

func TestGetMembershipsHTTPMethodAndOperation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGetMembershipsBuilder(pn)
	assert.Equal("GET", o.opts.httpMethod())
	assert.Equal(PNGetDataSyncMembershipsOperation, o.opts.operationType())
}

func TestGetMembershipsValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	o := newGetMembershipsBuilder(pn)
	assert.Contains(o.opts.validate().Error(), StrMissingSubKey)
}

func TestGetMembershipsValidateRejectsBothFilters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetMembershipsBuilder(pn).
		FilterFast("status == 'active'").
		Filter("status == 'active' AND role == 'admin'")
	assert.Contains(o.opts.validate().Error(), StrExclusiveDataSyncFilter)
}

func TestGetMembershipsResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	jsonBytes := []byte(`{"data":[{"id":"m-123","channelId":"general","userId":"alice","relationshipClass":"Membership","relationshipClassVersion":1,"status":"active","eTag":"1","payload":{"role":"moderator"}}],"links":{"self":"/self","next":"/next","prev":null},"meta":{"has_next":true,"has_prev":false,"next_cursor":"TjIw","prev_cursor":null,"limit":20}}`)

	r, _, err := newPNDataSyncMembershipsResponse(jsonBytes, newGetMembershipsOpts(pn, pn.ctx), StatusResponse{})
	assert.Nil(err)
	assert.Len(r.Data, 1)
	assert.Equal("m-123", r.Data[0].ID)
	assert.Equal("general", r.Data[0].ChannelID)
	assert.Equal("alice", r.Data[0].UserID)
	assert.NotNil(r.Meta)
	assert.True(r.Meta.HasNext)
	assert.Equal("TjIw", r.Meta.NextCursor)
	assert.Equal(20, r.Meta.Limit)
	assert.NotNil(r.Links)
	assert.Equal("/next", r.Links.Next)
}

func TestGetMembershipsExecuteValidationError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	_, _, err := pn.DataSync.GetMemberships().Execute()
	assert.Contains(err.Error(), StrMissingSubKey)
}

func TestGetMembershipsExecuteBothFiltersError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	_, _, err := pn.DataSync.GetMemberships().
		FilterFast("status == 'active'").
		Filter("status == 'active' AND role == 'admin'").
		Execute()
	assert.Contains(err.Error(), StrExclusiveDataSyncFilter)
}
