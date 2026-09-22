package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go/v10"
	"github.com/stretchr/testify/assert"
)

func TestDataSyncLiveUserCreateGetGlobal(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("gusr")
	stored := h.createGlobalUser(h.admin, id, globalUserPayload(), "active")
	assert.Equal(t, id, stored.ID)
	assert.Equal(t, dsUserClass, stored.EntityClass)
	assert.Equal(t, dsClassVersion, stored.EntityClassVersion)
	assert.Equal(t, pubnub.PNEntityClassLevelGlobal, stored.EntityClassLevel)
	assert.Equal(t, "active", stored.Status)
	assert.NotEmpty(t, stored.ETag)
	assert.Equal(t, "Ada", stored.Payload["name"])

	res, status, err := h.client.DataSync.GetUser().ID(id).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	if res == nil {
		t.Fatal("GetUser returned a nil response")
	}
	assert.Equal(t, id, res.Data.ID)
	assert.Equal(t, dsUserClass, res.Data.EntityClass)
	assert.Equal(t, "Ada", res.Data.Payload["name"])
	assert.Equal(t, "member", res.Data.Payload["type"])
}

func TestDataSyncLiveUserCreateGetCustom(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("usr")
	stored := h.createGoE2EUser(h.admin, id, kitchenSinkGoE2EUserPayload(), "active")
	assert.Equal(t, id, stored.ID)
	assert.Equal(t, dsGoE2EUserClass, stored.EntityClass)
	assert.Equal(t, dsClassVersion, stored.EntityClassVersion)
	assert.Equal(t, pubnub.PNEntityClassLevelSubKey, stored.EntityClassLevel)
	assert.Equal(t, "active", stored.Status)
	assert.NotEmpty(t, stored.ETag)
	assert.Equal(t, "undeclared", stored.Payload["scratch"])

	res, status, err := h.client.DataSync.GetUser().ID(id).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	if res == nil {
		t.Fatal("GetUser returned a nil response")
	}
	assert.Equal(t, id, res.Data.ID)
	assert.Equal(t, "Alpha", res.Data.Payload["displayName"])
	assert.EqualValues(t, 10, res.Data.Payload["score"])
	assert.Equal(t, "Portland", nestedCity(res.Data.Payload))
	assert.Equal(t, "alpha@example.com", res.Data.Payload["email"])
	assert.Equal(t, "s3cret", res.Data.Payload["secret"])
	assert.False(t, payloadHas(res.Data.Payload, "scratch"))
}

func TestDataSyncLiveUserCreateServerGeneratedID(t *testing.T) {
	h := newDSLive(t)
	h.grant(dsGrant{
		UsersPatterns:       dsUserAnyPerms(),
		UserProjPatterns:    map[string]string{".*": "admin"},
		MembershipPatterns:  map[string]pubnub.DataSyncPermissions{".*": dsCRUD()},
		ChannelPatterns:     dsEventChannelPatterns(),
		ChannelProjPatterns: map[string]string{".*": "admin"},
	})

	created := h.createGoE2EUser(h.client, "", adminGoE2EUserPayload(), "")
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, dsGoE2EUserClass, created.EntityClass)

	res, status, err := h.client.DataSync.GetUser().ID(created.ID).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, created.ID, res.Data.ID)
}

func TestDataSyncLiveUserCreateDuplicateID(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("usr")
	h.createGoE2EUser(h.client, id, adminGoE2EUserPayload(), "")

	_, status, err := h.client.DataSync.CreateUser().
		ID(id).
		EntityClass(dsGoE2EUserClass).
		EntityClassVersion(dsClassVersion).
		Payload(adminGoE2EUserPayload()).
		Execute()
	assertServerCode(t, err, status, 409)
}

func TestDataSyncLiveUserCreateMissingDisplayName(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	payload := adminGoE2EUserPayload()
	delete(payload, "displayName")
	_, status, err := h.client.DataSync.CreateUser().
		ID(dsNewID("usr")).
		EntityClass(dsGoE2EUserClass).
		EntityClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assertServerCode(t, err, status, 400)
}

func TestDataSyncLiveUserCreateBadValueKind(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	payload := adminGoE2EUserPayload()
	payload["score"] = "ten"
	_, status, err := h.client.DataSync.CreateUser().
		ID(dsNewID("usr")).
		EntityClass(dsGoE2EUserClass).
		EntityClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assertServerCode(t, err, status, 400)
}

func TestDataSyncLiveUserSetAndIfMatch(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("usr")
	created := h.createGoE2EUser(h.client, id, adminGoE2EUserPayload(), "")

	updated := adminGoE2EUserPayload()
	updated["displayName"] = "Bravo"
	updated["score"] = float64(20)
	updated["email"] = "bravo@example.com"

	res, status, err := h.client.DataSync.SetUser().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(updated).
		IfMatchETag(created.ETag).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "Bravo", res.Data.Payload["displayName"])
	assert.EqualValues(t, 20, res.Data.Payload["score"])
	assert.NotEqual(t, created.ETag, res.Data.ETag)

	_, staleStatus, staleErr := h.client.DataSync.SetUser().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(updated).
		IfMatchETag(created.ETag).
		Execute()
	assertServerCode(t, staleErr, staleStatus, 412)
}

func TestDataSyncLiveUserJSONPatch(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("usr")
	created := h.createGoE2EUser(h.client, id, adminGoE2EUserPayload(), "")

	res, status, err := h.client.DataSync.UpdateUser().
		ID(id).
		IfMatchETag(created.ETag).
		Test("/payload/displayName", "Alpha").
		Replace("/payload/displayName", "Patched").
		Replace("/payload/score", float64(99)).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "Patched", res.Data.Payload["displayName"])
	assert.EqualValues(t, 99, res.Data.Payload["score"])
	etag := res.Data.ETag

	res, status, err = h.client.DataSync.UpdateUser().
		ID(id).
		IfMatchETag(etag).
		Remove("/payload/notes").
		Execute()
	assert.Nil(t, err)
	assert.False(t, payloadHas(res.Data.Payload, "notes"))
	etag = res.Data.ETag

	res, status, err = h.client.DataSync.UpdateUser().
		ID(id).
		IfMatchETag(etag).
		Add("/payload/notes", "restored").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, "restored", res.Data.Payload["notes"])
	etag = res.Data.ETag

	res, status, err = h.client.DataSync.UpdateUser().
		ID(id).
		IfMatchETag(etag).
		Copy("/payload/displayName", "/payload/notes").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, "Patched", res.Data.Payload["notes"])
	etag = res.Data.ETag

	res, status, err = h.client.DataSync.UpdateUser().
		ID(id).
		IfMatchETag(etag).
		Move("/payload/notes", "/payload/category").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, "Patched", res.Data.Payload["category"])
	assert.False(t, payloadHas(res.Data.Payload, "notes"))

	_, staleStatus, staleErr := h.client.DataSync.UpdateUser().
		ID(id).
		IfMatchETag(created.ETag).
		Replace("/payload/displayName", "stale").
		Execute()
	assertServerCode(t, staleErr, staleStatus, 412)
}

func TestDataSyncLiveUserRemove(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("usr")
	created := h.createGoE2EUser(h.client, id, adminGoE2EUserPayload(), "")

	_, status, err := h.client.DataSync.RemoveUser().ID(id).IfMatchETag(created.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	_, getStatus, getErr := h.client.DataSync.GetUser().ID(id).Execute()
	assertServerCode(t, getErr, getStatus, 404)
}

func TestDataSyncLiveUserSetStatusDefaultProjection(t *testing.T) {
	h := newDSLive(t)
	h.grantProjection(pubnub.PNDataSyncDefaultProjection)

	id := dsNewID("usr")
	created := h.createGoE2EUser(h.client, id, defaultGoE2EUserPayload(), "active")
	assert.Equal(t, "active", created.Status)

	res, status, err := h.client.DataSync.SetUser().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Status("inactive").
		Payload(defaultGoE2EUserPayload()).
		IfMatchETag(created.ETag).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "inactive", res.Data.Status)
}

func TestDataSyncLiveUserGlobalSetAndRemove(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("gusr")
	created := h.createGlobalUser(h.client, id, globalUserPayload(), "")

	updated := globalUserPayload()
	updated["name"] = "Grace"
	res, status, err := h.client.DataSync.SetUser().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(updated).
		IfMatchETag(created.ETag).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "Grace", res.Data.Payload["name"])

	_, status, err = h.client.DataSync.RemoveUser().ID(id).IfMatchETag(res.Data.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
}
