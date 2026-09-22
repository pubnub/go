package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go/v10"
	"github.com/stretchr/testify/assert"
)

func TestDataSyncLiveChannelCreateGetGlobal(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("gchn")
	stored := h.createGlobalChannel(h.admin, id, globalChannelPayload(), "active")
	assert.Equal(t, id, stored.ID)
	assert.Equal(t, dsChannelClass, stored.EntityClass)
	assert.Equal(t, dsClassVersion, stored.EntityClassVersion)
	assert.Equal(t, pubnub.PNEntityClassLevelGlobal, stored.EntityClassLevel)
	assert.Equal(t, "active", stored.Status)
	assert.NotEmpty(t, stored.ETag)
	assert.Equal(t, "Lobby", stored.Payload["name"])

	res, status, err := h.client.DataSync.GetChannel().ID(id).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	if res == nil {
		t.Fatal("GetChannel returned a nil response")
	}
	assert.Equal(t, id, res.Data.ID)
	assert.Equal(t, dsChannelClass, res.Data.EntityClass)
	assert.Equal(t, "Lobby", res.Data.Payload["name"])
	assert.Equal(t, "public", res.Data.Payload["type"])
}

func TestDataSyncLiveChannelCreateGetCustom(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("chn")
	stored := h.createGoE2EChannel(h.admin, id, kitchenSinkGoE2EChannelPayload(), "active")
	assert.Equal(t, id, stored.ID)
	assert.Equal(t, dsGoE2EChannelClass, stored.EntityClass)
	assert.Equal(t, dsClassVersion, stored.EntityClassVersion)
	assert.Equal(t, pubnub.PNEntityClassLevelSubKey, stored.EntityClassLevel)
	assert.Equal(t, "active", stored.Status)
	assert.NotEmpty(t, stored.ETag)
	assert.Equal(t, "undeclared", stored.Payload["scratch"])

	res, status, err := h.client.DataSync.GetChannel().ID(id).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, id, res.Data.ID)
	assert.Equal(t, "Alpha", res.Data.Payload["topic"])
	assert.EqualValues(t, 10, res.Data.Payload["capacity"])
	assert.Equal(t, "Portland", nestedCity(res.Data.Payload))
	assert.Equal(t, "alpha@example.com", res.Data.Payload["email"])
	assert.Equal(t, "s3cret", res.Data.Payload["secret"])
	assert.False(t, payloadHas(res.Data.Payload, "scratch"))
}

func TestDataSyncLiveChannelCreateServerGeneratedID(t *testing.T) {
	h := newDSLive(t)
	h.grant(dsGrant{
		ChannelPatterns: map[string]pubnub.ChannelPermissions{
			".*":              dsChannelResourcePerms(),
			"^__admin__.*$":   {Read: true},
			"^__private__.*$": {Read: true},
		},
		ChannelProjPatterns: map[string]string{".*": "admin"},
		MembershipPatterns:  map[string]pubnub.DataSyncPermissions{".*": dsCRUD()},
	})

	created := h.createGoE2EChannel(h.client, "", adminGoE2EChannelPayload(), "")
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, dsGoE2EChannelClass, created.EntityClass)

	res, status, err := h.client.DataSync.GetChannel().ID(created.ID).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, created.ID, res.Data.ID)
}

func TestDataSyncLiveChannelCreateDuplicateID(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("chn")
	h.createGoE2EChannel(h.client, id, adminGoE2EChannelPayload(), "")

	_, status, err := h.client.DataSync.CreateChannel().
		ID(id).
		EntityClass(dsGoE2EChannelClass).
		EntityClassVersion(dsClassVersion).
		Payload(adminGoE2EChannelPayload()).
		Execute()
	assertServerCode(t, err, status, 409)
}

func TestDataSyncLiveChannelCreateMissingTopic(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	payload := adminGoE2EChannelPayload()
	delete(payload, "topic")
	_, status, err := h.client.DataSync.CreateChannel().
		ID(dsNewID("chn")).
		EntityClass(dsGoE2EChannelClass).
		EntityClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assertServerCode(t, err, status, 400)
}

func TestDataSyncLiveChannelCreateBadValueKind(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	payload := adminGoE2EChannelPayload()
	payload["capacity"] = "ten"
	_, status, err := h.client.DataSync.CreateChannel().
		ID(dsNewID("chn")).
		EntityClass(dsGoE2EChannelClass).
		EntityClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assertServerCode(t, err, status, 400)
}

func TestDataSyncLiveChannelSetAndIfMatch(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("chn")
	created := h.createGoE2EChannel(h.client, id, adminGoE2EChannelPayload(), "")

	updated := adminGoE2EChannelPayload()
	updated["topic"] = "Bravo"
	updated["capacity"] = float64(20)
	updated["email"] = "bravo@example.com"

	res, status, err := h.client.DataSync.SetChannel().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(updated).
		IfMatchETag(created.ETag).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "Bravo", res.Data.Payload["topic"])
	assert.EqualValues(t, 20, res.Data.Payload["capacity"])
	assert.NotEqual(t, created.ETag, res.Data.ETag)

	_, staleStatus, staleErr := h.client.DataSync.SetChannel().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(updated).
		IfMatchETag(created.ETag).
		Execute()
	assertServerCode(t, staleErr, staleStatus, 412)
}

func TestDataSyncLiveChannelJSONPatch(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("chn")
	created := h.createGoE2EChannel(h.client, id, adminGoE2EChannelPayload(), "")

	res, status, err := h.client.DataSync.UpdateChannel().
		ID(id).
		IfMatchETag(created.ETag).
		Test("/payload/topic", "Alpha").
		Replace("/payload/topic", "Patched").
		Replace("/payload/capacity", float64(99)).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "Patched", res.Data.Payload["topic"])
	assert.EqualValues(t, 99, res.Data.Payload["capacity"])
	etag := res.Data.ETag

	res, status, err = h.client.DataSync.UpdateChannel().
		ID(id).
		IfMatchETag(etag).
		Remove("/payload/notes").
		Execute()
	assert.Nil(t, err)
	assert.False(t, payloadHas(res.Data.Payload, "notes"))
	etag = res.Data.ETag

	res, status, err = h.client.DataSync.UpdateChannel().
		ID(id).
		IfMatchETag(etag).
		Add("/payload/notes", "restored").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, "restored", res.Data.Payload["notes"])
	etag = res.Data.ETag

	res, status, err = h.client.DataSync.UpdateChannel().
		ID(id).
		IfMatchETag(etag).
		Copy("/payload/topic", "/payload/notes").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, "Patched", res.Data.Payload["notes"])
	etag = res.Data.ETag

	res, status, err = h.client.DataSync.UpdateChannel().
		ID(id).
		IfMatchETag(etag).
		Move("/payload/notes", "/payload/category").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, "Patched", res.Data.Payload["category"])
	assert.False(t, payloadHas(res.Data.Payload, "notes"))

	_, staleStatus, staleErr := h.client.DataSync.UpdateChannel().
		ID(id).
		IfMatchETag(created.ETag).
		Replace("/payload/topic", "stale").
		Execute()
	assertServerCode(t, staleErr, staleStatus, 412)
}

func TestDataSyncLiveChannelRemove(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("chn")
	created := h.createGoE2EChannel(h.client, id, adminGoE2EChannelPayload(), "")

	_, status, err := h.client.DataSync.RemoveChannel().ID(id).IfMatchETag(created.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	_, getStatus, getErr := h.client.DataSync.GetChannel().ID(id).Execute()
	assertServerCode(t, getErr, getStatus, 404)
}

func TestDataSyncLiveChannelSetStatusDefaultProjection(t *testing.T) {
	h := newDSLive(t)
	h.grantProjection(pubnub.PNDataSyncDefaultProjection)

	id := dsNewID("chn")
	created := h.createGoE2EChannel(h.client, id, defaultGoE2EChannelPayload(), "active")
	assert.Equal(t, "active", created.Status)

	res, status, err := h.client.DataSync.SetChannel().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Status("inactive").
		Payload(defaultGoE2EChannelPayload()).
		IfMatchETag(created.ETag).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "inactive", res.Data.Status)
}

func TestDataSyncLiveChannelGlobalSetAndRemove(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("gchn")
	created := h.createGlobalChannel(h.client, id, globalChannelPayload(), "")

	updated := globalChannelPayload()
	updated["name"] = "Ops"
	res, status, err := h.client.DataSync.SetChannel().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(updated).
		IfMatchETag(created.ETag).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "Ops", res.Data.Payload["name"])

	_, status, err = h.client.DataSync.RemoveChannel().ID(id).IfMatchETag(res.Data.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
}
