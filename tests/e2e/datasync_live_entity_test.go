package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go/v9"
	"github.com/stretchr/testify/assert"
)

func TestDataSyncLiveEntityCreateGet(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("ent")
	stored := h.createEntity(h.admin, id, kitchenSinkEntityPayload(), "active")

	assert.Equal(t, id, stored.ID)
	assert.Equal(t, dsEntityClass, stored.EntityClass)
	assert.Equal(t, dsClassVersion, stored.EntityClassVersion)
	assert.Equal(t, pubnub.PNEntityClassLevelSubKey, stored.EntityClassLevel)
	assert.Equal(t, "active", stored.Status)
	assert.NotEmpty(t, stored.ETag)
	assert.NotEmpty(t, stored.CreatedAt)
	assert.NotEmpty(t, stored.UpdatedAt)
	assert.Equal(t, "undeclared", stored.Payload["scratch"])

	res, status, err := h.clientGetEntity(id)
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, id, res.Data.ID)
	assert.Equal(t, "Alpha", res.Data.Payload["title"])
	assert.EqualValues(t, 10, res.Data.Payload["score"])
	assert.Equal(t, "Portland", nestedCity(res.Data.Payload))
	assert.Equal(t, "alpha@example.com", res.Data.Payload["email"])
	assert.Equal(t, "s3cret", res.Data.Payload["secret"])
	assert.False(t, payloadHas(res.Data.Payload, "scratch"))
}

func TestDataSyncLiveEntityCreateServerGeneratedID(t *testing.T) {
	h := newDSLive(t)
	// Server-generated ids are not goe2e-* prefixed, so grant create on any id.
	h.grant(dsGrant{
		EntityPatterns:       map[string]pubnub.DataSyncPermissions{".*": dsCRUD()},
		RelationshipPatterns: map[string]pubnub.DataSyncPermissions{".*": dsCRUD()},
		EntityProjPatterns:   map[string]string{".*": "admin"},
		RelProjPatterns:      map[string]string{".*": "admin"},
	})

	created := h.createEntity(h.client, "", adminEntityPayload(), "")
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, dsEntityClass, created.EntityClass)

	res, status, err := h.clientGetEntity(created.ID)
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, created.ID, res.Data.ID)
}

func TestDataSyncLiveEntityCreateDuplicateID(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("ent")
	h.createEntity(h.client, id, adminEntityPayload(), "")

	_, status, err := h.client.DataSync.CreateEntity().
		ID(id).
		EntityClass(dsEntityClass).
		EntityClassVersion(dsClassVersion).
		Payload(adminEntityPayload()).
		Execute()
	assertServerCode(t, err, status, 409)
}

func TestDataSyncLiveEntityCreateMissingTitle(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	payload := adminEntityPayload()
	delete(payload, "title")
	_, status, err := h.client.DataSync.CreateEntity().
		ID(dsNewID("ent")).
		EntityClass(dsEntityClass).
		EntityClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assertServerCode(t, err, status, 400)
}

func TestDataSyncLiveEntityCreateBadValueKind(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	payload := adminEntityPayload()
	payload["score"] = "ten"
	_, status, err := h.client.DataSync.CreateEntity().
		ID(dsNewID("ent")).
		EntityClass(dsEntityClass).
		EntityClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assertServerCode(t, err, status, 400)
}

func TestDataSyncLiveEntitySetAndIfMatch(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("ent")
	created := h.createEntity(h.client, id, adminEntityPayload(), "")

	updatedPayload := adminEntityPayload()
	updatedPayload["title"] = "Bravo"
	updatedPayload["score"] = float64(20)
	updatedPayload["email"] = "bravo@example.com"

	res, status, err := h.client.DataSync.SetEntity().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(updatedPayload).
		IfMatchETag(created.ETag).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "Bravo", res.Data.Payload["title"])
	assert.EqualValues(t, 20, res.Data.Payload["score"])
	assert.NotEqual(t, created.ETag, res.Data.ETag)

	_, staleStatus, staleErr := h.client.DataSync.SetEntity().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(updatedPayload).
		IfMatchETag(created.ETag).
		Execute()
	assertServerCode(t, staleErr, staleStatus, 412)
}

func TestDataSyncLiveEntityJSONPatch(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("ent")
	created := h.createEntity(h.client, id, adminEntityPayload(), "")

	res, status, err := h.client.DataSync.UpdateEntity().
		ID(id).
		IfMatchETag(created.ETag).
		Test("/payload/title", "Alpha").
		Replace("/payload/title", "Patched").
		Replace("/payload/score", float64(99)).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "Patched", res.Data.Payload["title"])
	assert.EqualValues(t, 99, res.Data.Payload["score"])
	etag := res.Data.ETag

	res, status, err = h.client.DataSync.UpdateEntity().
		ID(id).
		IfMatchETag(etag).
		Remove("/payload/notes").
		Execute()
	assert.Nil(t, err)
	assert.False(t, payloadHas(res.Data.Payload, "notes"))
	etag = res.Data.ETag

	res, status, err = h.client.DataSync.UpdateEntity().
		ID(id).
		IfMatchETag(etag).
		Add("/payload/notes", "restored").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, "restored", res.Data.Payload["notes"])
	etag = res.Data.ETag

	res, status, err = h.client.DataSync.UpdateEntity().
		ID(id).
		IfMatchETag(etag).
		Copy("/payload/title", "/payload/notes").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, "Patched", res.Data.Payload["notes"])
	etag = res.Data.ETag

	res, status, err = h.client.DataSync.UpdateEntity().
		ID(id).
		IfMatchETag(etag).
		Move("/payload/notes", "/payload/category").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, "Patched", res.Data.Payload["category"])
	assert.False(t, payloadHas(res.Data.Payload, "notes"))

	_, staleStatus, staleErr := h.client.DataSync.UpdateEntity().
		ID(id).
		IfMatchETag(created.ETag).
		Replace("/payload/title", "stale").
		Execute()
	assertServerCode(t, staleErr, staleStatus, 412)
}

func TestDataSyncLiveEntityRemove(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("ent")
	created := h.createEntity(h.client, id, adminEntityPayload(), "")

	_, status, err := h.client.DataSync.RemoveEntity().ID(id).IfMatchETag(created.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	_, getStatus, getErr := h.clientGetEntity(id)
	assertServerCode(t, getErr, getStatus, 404)
}

func TestDataSyncLiveEntityRemoveCascadeRelationships(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	aID := dsNewID("ent")
	bID := dsNewID("st")
	h.createEntity(h.client, aID, adminEntityPayload(), "")
	h.createStatusEntity(h.client, bID, "Alpha", "active")
	rel := h.createRelationship(h.client, dsNewID("rel"), aID, bID, adminRelPayload(), "")

	entity := h.adminGetEntity(aID)
	_, status, err := h.client.DataSync.RemoveEntity().ID(aID).IfMatchETag(entity.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	_, relStatus, relErr := h.client.DataSync.GetRelationship().ID(rel.ID).Execute()
	assertServerCode(t, relErr, relStatus, 404)
}

func TestDataSyncLiveEntitySetStatusDefaultProjection(t *testing.T) {
	h := newDSLive(t)
	h.grantProjection(pubnub.PNDataSyncDefaultProjection)

	id := dsNewID("ent")
	created := h.createEntity(h.client, id, defaultEntityPayload(), "active")
	assert.Equal(t, "active", created.Status)

	res, status, err := h.client.DataSync.SetEntity().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Status("inactive").
		Payload(defaultEntityPayload()).
		IfMatchETag(created.ETag).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "inactive", res.Data.Status)
}
