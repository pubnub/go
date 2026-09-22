package e2e

import (
	"regexp"
	"testing"

	pubnub "github.com/pubnub/go/v10"
	"github.com/stretchr/testify/assert"
)

func TestDataSyncLivePAMUserMissingPermissions(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("usr")
	h.createGoE2EUser(h.admin, id, kitchenSinkGoE2EUserPayload(), "active")

	t.Run("GetDenied", func(t *testing.T) {
		h.grant(dsGrant{
			UsersPatterns:    dsUserPerms(pubnub.UUIDPermissions{Create: true, Update: true, Delete: true}),
			UserProjPatterns: map[string]string{dsIDPattern(): "admin"},
		})
		_, status, err := h.client.DataSync.GetUser().ID(id).Execute()
		assertServerCode(t, err, status, 403)
	})

	t.Run("CreateDenied", func(t *testing.T) {
		h.grant(dsGrant{
			UsersPatterns:    dsUserPerms(pubnub.UUIDPermissions{Get: true, Update: true, Delete: true}),
			UserProjPatterns: map[string]string{dsIDPattern(): "admin"},
		})
		_, status, err := h.client.DataSync.CreateUser().
			ID(dsNewID("usr")).
			EntityClass(dsGoE2EUserClass).
			EntityClassVersion(dsClassVersion).
			Payload(adminGoE2EUserPayload()).
			Execute()
		assertServerCode(t, err, status, 403)
	})

	t.Run("UpdateDenied", func(t *testing.T) {
		h.grant(dsGrant{
			UsersPatterns:    dsUserPerms(pubnub.UUIDPermissions{Get: true, Create: true, Delete: true}),
			UserProjPatterns: map[string]string{dsIDPattern(): "admin"},
		})
		_, status, err := h.client.DataSync.UpdateUser().
			ID(id).
			Replace("/payload/displayName", "Nope").
			Execute()
		assertServerCode(t, err, status, 403)
	})

	t.Run("DeleteDenied", func(t *testing.T) {
		h.grant(dsGrant{
			UsersPatterns:    dsUserPerms(pubnub.UUIDPermissions{Get: true, Create: true, Update: true}),
			UserProjPatterns: map[string]string{dsIDPattern(): "admin"},
		})
		_, status, err := h.client.DataSync.RemoveUser().ID(id).Execute()
		assertServerCode(t, err, status, 403)
	})
}

func TestDataSyncLivePAMUserExactIDGrant(t *testing.T) {
	h := newDSLive(t)
	allowed := dsNewID("usr")
	denied := dsNewID("usr")
	h.createGoE2EUser(h.admin, allowed, kitchenSinkGoE2EUserPayload(), "active")
	h.createGoE2EUser(h.admin, denied, kitchenSinkGoE2EUserPayload(), "active")

	h.grant(dsGrant{
		Users: map[string]pubnub.UUIDPermissions{allowed: dsUUIDCRUD()},
		UsersPatterns: map[string]pubnub.UUIDPermissions{
			allowed:                          dsUUIDCRUD(),
			allowed + "*":                    dsUUIDCRUD(),
			regexp.QuoteMeta(allowed) + ".*": dsUUIDCRUD(),
		},
		UserProjections:     map[string]string{allowed: "admin"},
		UserProjPatterns:    dsUserProjPatterns("admin"),
		ChannelPatterns:     dsEventChannelPatterns(),
		ChannelProjPatterns: dsChannelProjPatterns("admin"),
	})

	res, status, err := h.client.DataSync.GetUser().ID(allowed).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	if res == nil {
		t.Fatal("GetUser returned a nil response")
	}
	assert.Equal(t, allowed, res.Data.ID)

	_, deniedStatus, deniedErr := h.client.DataSync.GetUser().ID(denied).Execute()
	assertServerCode(t, deniedErr, deniedStatus, 403)
}

func TestDataSyncLivePAMChannelMissingPermissions(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("chn")
	h.createGoE2EChannel(h.admin, id, kitchenSinkGoE2EChannelPayload(), "active")

	t.Run("GetDenied", func(t *testing.T) {
		h.grant(dsGrant{
			ChannelPatterns: map[string]pubnub.ChannelPermissions{
				"^goe2e-.*$": {Read: true, Create: true, Update: true, Delete: true},
			},
			ChannelProjPatterns: map[string]string{dsIDPattern(): "admin"},
		})
		_, status, err := h.client.DataSync.GetChannel().ID(id).Execute()
		assertServerCode(t, err, status, 403)
	})

	t.Run("CreateDenied", func(t *testing.T) {
		h.grant(dsGrant{
			ChannelPatterns: map[string]pubnub.ChannelPermissions{
				"^goe2e-.*$": {Read: true, Get: true, Update: true, Delete: true},
			},
			ChannelProjPatterns: map[string]string{dsIDPattern(): "admin"},
		})
		_, status, err := h.client.DataSync.CreateChannel().
			ID(dsNewID("chn")).
			EntityClass(dsGoE2EChannelClass).
			EntityClassVersion(dsClassVersion).
			Payload(adminGoE2EChannelPayload()).
			Execute()
		assertServerCode(t, err, status, 403)
	})
}

func TestDataSyncLivePAMMembershipDenied(t *testing.T) {
	h := newDSLive(t)
	userID := dsNewID("gusr")
	chID := dsNewID("gchn")
	h.createGlobalUser(h.admin, userID, globalUserPayload(), "active")
	h.createGlobalChannel(h.admin, chID, globalChannelPayload(), "active")
	memID := dsNewID("mem")
	h.createMembership(h.admin, memID, userID, chID, nil, "active")

	h.grant(dsGrant{
		UsersPatterns:       dsUserPerms(dsUUIDCRUD()),
		UserProjPatterns:    dsUserProjPatterns("admin"),
		ChannelProjPatterns: dsChannelProjPatterns("admin"),
		ChannelPatterns:     dsEventChannelPatterns(),
		MembershipPatterns: map[string]pubnub.DataSyncPermissions{
			dsIDPattern(): {Get: true, Create: true, Update: true},
		},
		MemProjPatterns: map[string]string{dsIDPattern(): pubnub.PNDataSyncDefaultProjection},
	})
	_, status, err := h.client.DataSync.RemoveMembership().ID(memID).Execute()
	assertServerCode(t, err, status, 403)

	h.grant(dsGrant{
		UsersPatterns:       dsUserPerms(dsUUIDCRUD()),
		UserProjPatterns:    dsUserProjPatterns("admin"),
		ChannelProjPatterns: dsChannelProjPatterns("admin"),
		ChannelPatterns:     dsEventChannelPatterns(),
		MembershipPatterns: map[string]pubnub.DataSyncPermissions{
			dsIDPattern(): {Create: true, Update: true, Delete: true},
		},
		MemProjPatterns: map[string]string{dsIDPattern(): pubnub.PNDataSyncDefaultProjection},
	})
	_, getStatus, getErr := h.client.DataSync.GetMembership().ID(memID).Execute()
	assertServerCode(t, getErr, getStatus, 403)
}

func TestDataSyncLiveProjectionUserDefault(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("usr")
	h.createGoE2EUser(h.admin, id, kitchenSinkGoE2EUserPayload(), "active")
	h.grantProjection(pubnub.PNDataSyncDefaultProjection)

	res, status, err := h.client.DataSync.GetUser().ID(id).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	if res == nil {
		t.Fatal("GetUser returned a nil response")
	}
	assert.Equal(t, "Alpha", res.Data.Payload["displayName"])
	assert.Equal(t, "tools", res.Data.Payload["category"])
	assert.Equal(t, "Portland", nestedCity(res.Data.Payload))
	assert.Equal(t, "undeclared", res.Data.Payload["scratch"])
	assert.False(t, payloadHas(res.Data.Payload, "email"))
	assert.False(t, payloadHas(res.Data.Payload, "secret"))
	assert.Equal(t, "active", res.Data.Status)

	payload := defaultGoE2EUserPayload()
	payload["email"] = "blocked@example.com"
	_, writeStatus, writeErr := h.client.DataSync.SetUser().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assertServerCode(t, writeErr, writeStatus, 403)
}

func TestDataSyncLiveProjectionUserPrivate(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("usr")
	h.createGoE2EUser(h.admin, id, kitchenSinkGoE2EUserPayload(), "active")
	h.grantProjection("private")

	res, status, err := h.client.DataSync.GetUser().ID(id).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	if res == nil {
		t.Fatal("GetUser returned a nil response")
	}
	assert.Equal(t, "Alpha", res.Data.Payload["displayName"])
	assert.Equal(t, "alpha@example.com", res.Data.Payload["email"])
	assert.False(t, payloadHas(res.Data.Payload, "secret"))
	assert.False(t, payloadHas(res.Data.Payload, "scratch"))
	assert.Empty(t, res.Data.Status)

	payload := adminGoE2EUserPayload()
	payload["secret"] = "nope"
	_, writeStatus, writeErr := h.client.DataSync.SetUser().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assertServerCode(t, writeErr, writeStatus, 403)

	scratch := adminGoE2EUserPayload()
	delete(scratch, "secret")
	scratch["scratch"] = "nope"
	_, scratchStatus, scratchErr := h.client.DataSync.SetUser().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(scratch).
		Execute()
	assertServerCode(t, scratchErr, scratchStatus, 403)

	_, patchStatus, patchErr := h.client.DataSync.UpdateUser().
		ID(id).
		Replace("/payload/secret", "nope").
		Execute()
	assertServerCode(t, patchErr, patchStatus, 403)
}

func TestDataSyncLiveProjectionUserAdmin(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("usr")
	h.createGoE2EUser(h.admin, id, kitchenSinkGoE2EUserPayload(), "active")
	h.grantProjection("admin")

	res, status, err := h.client.DataSync.GetUser().ID(id).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	if res == nil {
		t.Fatal("GetUser returned a nil response")
	}
	assert.Equal(t, "s3cret", res.Data.Payload["secret"])
	assert.Equal(t, "alpha@example.com", res.Data.Payload["email"])
	assert.False(t, payloadHas(res.Data.Payload, "scratch"))
	assert.Empty(t, res.Data.Status)
}

func TestDataSyncLiveProjectionUserReplaceDoesNotWipeHiddenFields(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("usr")
	h.createGoE2EUser(h.admin, id, kitchenSinkGoE2EUserPayload(), "")
	h.grantProjection("private")

	payload := adminGoE2EUserPayload()
	delete(payload, "secret")
	payload["displayName"] = "PrivateEdit"
	payload["email"] = "private@example.com"

	res, status, err := h.client.DataSync.SetUser().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	if res == nil {
		t.Fatal("SetUser returned a nil response")
	}
	assert.Equal(t, "PrivateEdit", res.Data.Payload["displayName"])
	assert.False(t, payloadHas(res.Data.Payload, "secret"))

	oracle := h.adminGetUser(id)
	assert.Equal(t, "s3cret", oracle.Payload["secret"])
	assert.Equal(t, "PrivateEdit", oracle.Payload["displayName"])
}

func TestDataSyncLiveProjectionChannelDefault(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("chn")
	h.createGoE2EChannel(h.admin, id, kitchenSinkGoE2EChannelPayload(), "active")
	h.grantProjection(pubnub.PNDataSyncDefaultProjection)

	res, status, err := h.client.DataSync.GetChannel().ID(id).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	if res == nil {
		t.Fatal("GetChannel returned a nil response")
	}
	assert.Equal(t, "Alpha", res.Data.Payload["topic"])
	assert.Equal(t, "undeclared", res.Data.Payload["scratch"])
	assert.False(t, payloadHas(res.Data.Payload, "email"))
	assert.False(t, payloadHas(res.Data.Payload, "secret"))
	assert.Equal(t, "active", res.Data.Status)
}

func TestDataSyncLiveProjectionChannelPrivate(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("chn")
	h.createGoE2EChannel(h.admin, id, kitchenSinkGoE2EChannelPayload(), "active")
	h.grantProjection("private")

	res, status, err := h.client.DataSync.GetChannel().ID(id).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	if res == nil {
		t.Fatal("GetChannel returned a nil response")
	}
	assert.Equal(t, "alpha@example.com", res.Data.Payload["email"])
	assert.False(t, payloadHas(res.Data.Payload, "secret"))
	assert.False(t, payloadHas(res.Data.Payload, "scratch"))
	assert.Empty(t, res.Data.Status)

	_, patchStatus, patchErr := h.client.DataSync.UpdateChannel().
		ID(id).
		Replace("/payload/secret", "nope").
		Execute()
	assertServerCode(t, patchErr, patchStatus, 403)
}

func TestDataSyncLiveProjectionChannelAdmin(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("chn")
	h.createGoE2EChannel(h.admin, id, kitchenSinkGoE2EChannelPayload(), "active")
	h.grantProjection("admin")

	res, status, err := h.client.DataSync.GetChannel().ID(id).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	if res == nil {
		t.Fatal("GetChannel returned a nil response")
	}
	assert.Equal(t, "s3cret", res.Data.Payload["secret"])
	assert.Equal(t, "alpha@example.com", res.Data.Payload["email"])
	assert.False(t, payloadHas(res.Data.Payload, "scratch"))
}
