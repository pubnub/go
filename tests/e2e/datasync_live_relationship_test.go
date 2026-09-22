package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go/v10"
	"github.com/stretchr/testify/assert"
)

func TestDataSyncLiveRelationshipCreateGet(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	aID := dsNewID("ent")
	bID := dsNewID("st")
	h.createEntity(h.admin, aID, kitchenSinkEntityPayload(), "active")
	h.createStatusEntity(h.admin, bID, "Alpha", "active")

	id := dsNewID("rel")
	created := h.createRelationship(h.admin, id, aID, bID, kitchenSinkRelPayload(), "active")
	assert.Equal(t, id, created.ID)
	assert.Equal(t, aID, created.EntityAID)
	assert.Equal(t, bID, created.EntityBID)
	assert.Equal(t, dsRelClass, created.RelationshipClass)
	assert.Equal(t, dsClassVersion, created.RelationshipClassVersion)
	assert.Equal(t, "active", created.Status)
	assert.NotEmpty(t, created.ETag)

	res, status, err := h.client.DataSync.GetRelationship().ID(id).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, id, res.Data.ID)
	assert.Equal(t, "owner", res.Data.Payload["role"])
	assert.Equal(t, "west", nestedRegion(res.Data.Payload))
	assert.Equal(t, "because", res.Data.Payload["reason"])
	assert.Equal(t, "rel-secret", res.Data.Payload["secret"])
	assert.Empty(t, res.Data.Status)
}

func TestDataSyncLiveRelationshipCreateServerGeneratedID(t *testing.T) {
	h := newDSLive(t)
	h.grant(dsGrant{
		EntityPatterns:       map[string]pubnub.DataSyncPermissions{".*": dsCRUD()},
		RelationshipPatterns: map[string]pubnub.DataSyncPermissions{".*": dsCRUD()},
		EntityProjPatterns:   map[string]string{".*": "admin"},
		RelProjPatterns:      map[string]string{".*": "admin"},
	})

	aID := dsNewID("ent")
	bID := dsNewID("st")
	h.createEntity(h.client, aID, adminEntityPayload(), "")
	h.createStatusEntity(h.client, bID, "Bravo", "active")

	created := h.createRelationship(h.client, "", aID, bID, adminRelPayload(), "")
	assert.NotEmpty(t, created.ID)

	res, status, err := h.client.DataSync.GetRelationship().ID(created.ID).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, created.ID, res.Data.ID)
}

func TestDataSyncLiveRelationshipCreateDuplicateID(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	aID := dsNewID("ent")
	bID := dsNewID("st")
	h.createEntity(h.client, aID, adminEntityPayload(), "")
	h.createStatusEntity(h.client, bID, "Dup", "active")
	id := dsNewID("rel")
	h.createRelationship(h.client, id, aID, bID, adminRelPayload(), "")

	a2 := dsNewID("ent")
	b2 := dsNewID("st")
	h.createEntity(h.client, a2, adminEntityPayload(), "")
	h.createStatusEntity(h.client, b2, "Dup2", "active")

	_, status, err := h.client.DataSync.CreateRelationship().
		ID(id).
		EntityAID(a2).
		EntityBID(b2).
		RelationshipClass(dsRelClass).
		RelationshipClassVersion(dsClassVersion).
		Payload(adminRelPayload()).
		Execute()
	assertServerCode(t, err, status, 409)
}

func TestDataSyncLiveRelationshipCreateMissingRole(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	aID := dsNewID("ent")
	bID := dsNewID("st")
	h.createEntity(h.client, aID, adminEntityPayload(), "")
	h.createStatusEntity(h.client, bID, "NoRole", "active")

	payload := adminRelPayload()
	delete(payload, "role")
	_, status, err := h.client.DataSync.CreateRelationship().
		ID(dsNewID("rel")).
		EntityAID(aID).
		EntityBID(bID).
		RelationshipClass(dsRelClass).
		RelationshipClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assertServerCode(t, err, status, 400)
}

func TestDataSyncLiveRelationshipSetAndIfMatch(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	aID := dsNewID("ent")
	bID := dsNewID("st")
	h.createEntity(h.client, aID, adminEntityPayload(), "")
	h.createStatusEntity(h.client, bID, "Set", "active")
	id := dsNewID("rel")
	created := h.createRelationship(h.client, id, aID, bID, adminRelPayload(), "")

	updated := adminRelPayload()
	updated["role"] = "editor"
	updated["weight"] = 9.5
	res, status, err := h.client.DataSync.SetRelationship().
		ID(id).
		RelationshipClassVersion(dsClassVersion).
		Payload(updated).
		IfMatchETag(created.ETag).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "editor", res.Data.Payload["role"])
	assert.EqualValues(t, 9.5, res.Data.Payload["weight"])
	assert.NotEqual(t, created.ETag, res.Data.ETag)

	_, staleStatus, staleErr := h.client.DataSync.SetRelationship().
		ID(id).
		RelationshipClassVersion(dsClassVersion).
		Payload(updated).
		IfMatchETag(created.ETag).
		Execute()
	assertServerCode(t, staleErr, staleStatus, 412)
}

func TestDataSyncLiveRelationshipJSONPatch(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	aID := dsNewID("ent")
	bID := dsNewID("st")
	h.createEntity(h.client, aID, adminEntityPayload(), "")
	h.createStatusEntity(h.client, bID, "Patch", "active")
	id := dsNewID("rel")
	created := h.createRelationship(h.client, id, aID, bID, adminRelPayload(), "")

	res, status, err := h.client.DataSync.UpdateRelationship().
		ID(id).
		IfMatchETag(created.ETag).
		Test("/payload/role", "owner").
		Replace("/payload/role", "reviewer").
		Replace("/payload/weight", 3.25).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "reviewer", res.Data.Payload["role"])
	assert.EqualValues(t, 3.25, res.Data.Payload["weight"])
	etag := res.Data.ETag

	res, status, err = h.client.DataSync.UpdateRelationship().
		ID(id).
		IfMatchETag(etag).
		Remove("/payload/comment").
		Execute()
	assert.Nil(t, err)
	assert.False(t, payloadHas(res.Data.Payload, "comment"))
	etag = res.Data.ETag

	res, status, err = h.client.DataSync.UpdateRelationship().
		ID(id).
		IfMatchETag(etag).
		Add("/payload/comment", "restored").
		Copy("/payload/role", "/payload/comment").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, "reviewer", res.Data.Payload["comment"])
	etag = res.Data.ETag

	res, status, err = h.client.DataSync.UpdateRelationship().
		ID(id).
		IfMatchETag(etag).
		Move("/payload/comment", "/payload/kind").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, "reviewer", res.Data.Payload["kind"])
	assert.False(t, payloadHas(res.Data.Payload, "comment"))

	_, staleStatus, staleErr := h.client.DataSync.UpdateRelationship().
		ID(id).
		IfMatchETag(created.ETag).
		Replace("/payload/role", "stale").
		Execute()
	assertServerCode(t, staleErr, staleStatus, 412)
}

func TestDataSyncLiveRelationshipRemove(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	aID := dsNewID("ent")
	bID := dsNewID("st")
	h.createEntity(h.client, aID, adminEntityPayload(), "")
	h.createStatusEntity(h.client, bID, "Del", "active")
	id := dsNewID("rel")
	created := h.createRelationship(h.client, id, aID, bID, adminRelPayload(), "")

	_, status, err := h.client.DataSync.RemoveRelationship().ID(id).IfMatchETag(created.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	_, getStatus, getErr := h.client.DataSync.GetRelationship().ID(id).Execute()
	assertServerCode(t, getErr, getStatus, 404)

	_, aStatus, aErr := h.clientGetEntity(aID)
	assert.Nil(t, aErr)
	assert.Equal(t, 200, aStatus.StatusCode)
}
