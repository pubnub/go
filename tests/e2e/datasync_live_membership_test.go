package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go/v9"
	"github.com/stretchr/testify/assert"
)

func TestDataSyncLiveMembershipCreateGet(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	userID := dsNewID("gusr")
	chID := dsNewID("gchn")
	h.createGlobalUser(h.admin, userID, globalUserPayload(), "active")
	h.createGlobalChannel(h.admin, chID, globalChannelPayload(), "active")

	id := dsNewID("mem")
	created := h.createMembership(h.admin, id, userID, chID, nil, "active")
	assert.Equal(t, id, created.ID)
	assert.Equal(t, userID, created.UserID)
	assert.Equal(t, chID, created.ChannelID)
	assert.Equal(t, dsMembershipClass, created.RelationshipClass)
	assert.Equal(t, dsClassVersion, created.RelationshipClassVersion)
	assert.Equal(t, "active", created.Status)
	assert.NotEmpty(t, created.ETag)

	res, status, err := h.client.DataSync.GetMembership().ID(id).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, id, res.Data.ID)
	assert.Equal(t, userID, res.Data.UserID)
	assert.Equal(t, chID, res.Data.ChannelID)
}

func TestDataSyncLiveMembershipCreateServerGeneratedID(t *testing.T) {
	h := newDSLive(t)
	h.grant(dsGrant{
		UsersPatterns:       dsUserAnyPerms(),
		UserProjPatterns:    dsUserProjPatterns("admin"),
		ChannelProjPatterns: dsChannelProjPatterns("admin"),
		MembershipPatterns:  map[string]pubnub.DataSyncPermissions{".*": dsCRUD()},
		MemProjPatterns:     map[string]string{".*": pubnub.PNDataSyncDefaultProjection},
		ChannelPatterns:     dsEventChannelPatterns(),
	})

	userID := dsNewID("gusr")
	chID := dsNewID("gchn")
	h.createGlobalUser(h.admin, userID, globalUserPayload(), "")
	h.createGlobalChannel(h.admin, chID, globalChannelPayload(), "")

	created := h.createMembership(h.client, "", userID, chID, nil, "")
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, dsMembershipClass, created.RelationshipClass)

	res, status, err := h.client.DataSync.GetMembership().ID(created.ID).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, created.ID, res.Data.ID)
}

func TestDataSyncLiveMembershipCreateDuplicateID(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	u1 := dsNewID("gusr")
	c1 := dsNewID("gchn")
	h.createGlobalUser(h.admin, u1, globalUserPayload(), "")
	h.createGlobalChannel(h.admin, c1, globalChannelPayload(), "")
	id := dsNewID("mem")
	h.createMembership(h.client, id, u1, c1, nil, "")

	u2 := dsNewID("gusr")
	c2 := dsNewID("gchn")
	h.createGlobalUser(h.admin, u2, globalUserPayload(), "")
	h.createGlobalChannel(h.admin, c2, globalChannelPayload(), "")

	_, status, err := h.client.DataSync.CreateMembership().
		ID(id).
		UserID(u2).
		ChannelID(c2).
		RelationshipClassVersion(dsClassVersion).
		Execute()
	assertServerCode(t, err, status, 409)
}

func TestDataSyncLiveMembershipAllowsUserChannelSubclasses(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	userID := dsNewID("usr")
	chID := dsNewID("chn")
	h.createGoE2EUser(h.admin, userID, kitchenSinkGoE2EUserPayload(), "")
	h.createGoE2EChannel(h.admin, chID, kitchenSinkGoE2EChannelPayload(), "")

	created := h.createMembership(h.client, dsNewID("mem"), userID, chID, nil, "")
	assert.Equal(t, userID, created.UserID)
	assert.Equal(t, chID, created.ChannelID)
	assert.Equal(t, dsMembershipClass, created.RelationshipClass)
}

func TestDataSyncLiveMembershipRejectsGenericEntities(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	aID := dsNewID("ent")
	bID := dsNewID("st")
	h.createEntity(h.admin, aID, kitchenSinkEntityPayload(), "")
	h.createStatusEntity(h.admin, bID, "Mem", "")

	_, status, err := h.client.DataSync.CreateMembership().
		ID(dsNewID("mem")).
		UserID(aID).
		ChannelID(bID).
		RelationshipClassVersion(dsClassVersion).
		Execute()
	assertServerCode(t, err, status, 400)
}

func TestDataSyncLiveMembershipSetAndIfMatch(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	userID := dsNewID("gusr")
	chID := dsNewID("gchn")
	h.createGlobalUser(h.admin, userID, globalUserPayload(), "")
	h.createGlobalChannel(h.admin, chID, globalChannelPayload(), "")
	id := dsNewID("mem")
	created := h.createMembership(h.client, id, userID, chID, nil, "")

	res, status, err := h.client.DataSync.SetMembership().
		ID(id).
		RelationshipClassVersion(dsClassVersion).
		Status("inactive").
		IfMatchETag(created.ETag).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "inactive", res.Data.Status)
	assert.NotEqual(t, created.ETag, res.Data.ETag)

	_, staleStatus, staleErr := h.client.DataSync.SetMembership().
		ID(id).
		RelationshipClassVersion(dsClassVersion).
		Status("active").
		IfMatchETag(created.ETag).
		Execute()
	assertServerCode(t, staleErr, staleStatus, 412)
}

func TestDataSyncLiveMembershipJSONPatch(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	userID := dsNewID("gusr")
	chID := dsNewID("gchn")
	h.createGlobalUser(h.admin, userID, globalUserPayload(), "")
	h.createGlobalChannel(h.admin, chID, globalChannelPayload(), "")
	id := dsNewID("mem")
	created := h.createMembership(h.client, id, userID, chID, nil, "active")

	res, status, err := h.client.DataSync.UpdateMembership().
		ID(id).
		IfMatchETag(created.ETag).
		Test("/status", "active").
		Replace("/status", "inactive").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "inactive", res.Data.Status)

	_, staleStatus, staleErr := h.client.DataSync.UpdateMembership().
		ID(id).
		IfMatchETag(created.ETag).
		Replace("/status", "active").
		Execute()
	assertServerCode(t, staleErr, staleStatus, 412)
}

func TestDataSyncLiveMembershipRemove(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	userID := dsNewID("gusr")
	chID := dsNewID("gchn")
	h.createGlobalUser(h.admin, userID, globalUserPayload(), "")
	h.createGlobalChannel(h.admin, chID, globalChannelPayload(), "")
	id := dsNewID("mem")
	created := h.createMembership(h.client, id, userID, chID, nil, "")

	_, status, err := h.client.DataSync.RemoveMembership().ID(id).IfMatchETag(created.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	_, getStatus, getErr := h.client.DataSync.GetMembership().ID(id).Execute()
	assertServerCode(t, getErr, getStatus, 404)

	_, userStatus, userErr := h.client.DataSync.GetUser().ID(userID).Execute()
	assert.Nil(t, userErr)
	assert.Equal(t, 200, userStatus.StatusCode)
	_, chStatus, chErr := h.client.DataSync.GetChannel().ID(chID).Execute()
	assert.Nil(t, chErr)
	assert.Equal(t, 200, chStatus.StatusCode)
}

func TestDataSyncLiveMembershipCascadeOnUserRemove(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	userID := dsNewID("gusr")
	chID := dsNewID("gchn")
	user := h.createGlobalUser(h.admin, userID, globalUserPayload(), "")
	h.createGlobalChannel(h.admin, chID, globalChannelPayload(), "")
	mem := h.createMembership(h.client, dsNewID("mem"), userID, chID, nil, "")

	_, status, err := h.client.DataSync.RemoveUser().ID(userID).IfMatchETag(user.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	waitMembershipGone(t, h.client, mem.ID)
}

func TestDataSyncLiveMembershipCascadeOnChannelRemove(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	userID := dsNewID("gusr")
	chID := dsNewID("gchn")
	h.createGlobalUser(h.admin, userID, globalUserPayload(), "")
	ch := h.createGlobalChannel(h.admin, chID, globalChannelPayload(), "")
	mem := h.createMembership(h.client, dsNewID("mem"), userID, chID, nil, "")

	_, status, err := h.client.DataSync.RemoveChannel().ID(chID).IfMatchETag(ch.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	waitMembershipGone(t, h.client, mem.ID)
}
