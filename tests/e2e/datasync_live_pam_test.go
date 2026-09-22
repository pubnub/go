package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go/v10"
	"github.com/stretchr/testify/assert"
)

func TestDataSyncLivePAMMissingPermissions(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("ent")
	h.createEntity(h.admin, id, kitchenSinkEntityPayload(), "active")

	t.Run("GetDenied", func(t *testing.T) {
		h.grant(dsGrant{
			EntityPatterns: map[string]pubnub.DataSyncPermissions{
				dsIDPattern(): {Create: true, Update: true, Delete: true},
			},
			EntityProjPatterns: map[string]string{dsIDPattern(): "admin"},
		})
		_, status, err := h.clientGetEntity(id)
		assertServerCode(t, err, status, 403)
	})

	t.Run("CreateDenied", func(t *testing.T) {
		h.grant(dsGrant{
			EntityPatterns: map[string]pubnub.DataSyncPermissions{
				dsIDPattern(): {Get: true, Update: true, Delete: true},
			},
			EntityProjPatterns: map[string]string{dsIDPattern(): "admin"},
		})
		_, status, err := h.client.DataSync.CreateEntity().
			ID(dsNewID("ent")).
			EntityClass(dsEntityClass).
			EntityClassVersion(dsClassVersion).
			Payload(adminEntityPayload()).
			Execute()
		assertServerCode(t, err, status, 403)
	})

	t.Run("UpdateDenied", func(t *testing.T) {
		h.grant(dsGrant{
			EntityPatterns: map[string]pubnub.DataSyncPermissions{
				dsIDPattern(): {Get: true, Create: true, Delete: true},
			},
			EntityProjPatterns: map[string]string{dsIDPattern(): "admin"},
		})
		_, status, err := h.client.DataSync.UpdateEntity().
			ID(id).
			Replace("/payload/title", "Nope").
			Execute()
		assertServerCode(t, err, status, 403)
	})

	t.Run("DeleteDenied", func(t *testing.T) {
		h.grant(dsGrant{
			EntityPatterns: map[string]pubnub.DataSyncPermissions{
				dsIDPattern(): {Get: true, Create: true, Update: true},
			},
			EntityProjPatterns: map[string]string{dsIDPattern(): "admin"},
		})
		_, status, err := h.client.DataSync.RemoveEntity().ID(id).Execute()
		assertServerCode(t, err, status, 403)
	})
}

func TestDataSyncLivePAMExactIDGrant(t *testing.T) {
	h := newDSLive(t)
	allowed := dsNewID("ent")
	denied := dsNewID("ent")
	h.createEntity(h.admin, allowed, kitchenSinkEntityPayload(), "active")
	h.createEntity(h.admin, denied, kitchenSinkEntityPayload(), "active")

	h.grant(dsGrant{
		Entities: map[string]pubnub.DataSyncPermissions{
			allowed: {Get: true},
		},
		EntityProjections: map[string]string{allowed: "admin"},
	})

	res, status, err := h.clientGetEntity(allowed)
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, allowed, res.Data.ID)

	_, deniedStatus, deniedErr := h.clientGetEntity(denied)
	assertServerCode(t, deniedErr, deniedStatus, 403)
}

func TestDataSyncLivePAMRelationshipDenied(t *testing.T) {
	h := newDSLive(t)
	aID := dsNewID("ent")
	bID := dsNewID("st")
	h.createEntity(h.admin, aID, kitchenSinkEntityPayload(), "active")
	h.createStatusEntity(h.admin, bID, "PAM", "active")
	relID := dsNewID("rel")
	h.createRelationship(h.admin, relID, aID, bID, kitchenSinkRelPayload(), "active")

	h.grant(dsGrant{
		EntityPatterns: map[string]pubnub.DataSyncPermissions{dsIDPattern(): dsCRUD()},
		RelationshipPatterns: map[string]pubnub.DataSyncPermissions{
			dsIDPattern(): {Get: true, Create: true, Update: true},
		},
		EntityProjPatterns: map[string]string{dsIDPattern(): "admin"},
		RelProjPatterns:    map[string]string{dsIDPattern(): "admin"},
	})

	_, status, err := h.client.DataSync.RemoveRelationship().ID(relID).Execute()
	assertServerCode(t, err, status, 403)

	h.grant(dsGrant{
		EntityPatterns: map[string]pubnub.DataSyncPermissions{dsIDPattern(): dsCRUD()},
		RelationshipPatterns: map[string]pubnub.DataSyncPermissions{
			dsIDPattern(): {Create: true, Update: true, Delete: true},
		},
		EntityProjPatterns: map[string]string{dsIDPattern(): "admin"},
		RelProjPatterns:    map[string]string{dsIDPattern(): "admin"},
	})
	_, getStatus, getErr := h.client.DataSync.GetRelationship().ID(relID).Execute()
	assertServerCode(t, getErr, getStatus, 403)
}

func TestDataSyncLiveProjectionDefault(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("ent")
	h.createEntity(h.admin, id, kitchenSinkEntityPayload(), "active")
	h.grantProjection(pubnub.PNDataSyncDefaultProjection)

	res, status, err := h.clientGetEntity(id)
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "Alpha", res.Data.Payload["title"])
	assert.Equal(t, "tools", res.Data.Payload["category"])
	assert.Equal(t, "Portland", nestedCity(res.Data.Payload))
	assert.Equal(t, "undeclared", res.Data.Payload["scratch"])
	assert.False(t, payloadHas(res.Data.Payload, "email"))
	assert.False(t, payloadHas(res.Data.Payload, "secret"))
	assert.Equal(t, "active", res.Data.Status)

	payload := defaultEntityPayload()
	payload["email"] = "blocked@example.com"
	_, writeStatus, writeErr := h.client.DataSync.SetEntity().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assertServerCode(t, writeErr, writeStatus, 403)
}

func TestDataSyncLiveProjectionPrivate(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("ent")
	h.createEntity(h.admin, id, kitchenSinkEntityPayload(), "active")
	h.grantProjection("private")

	res, status, err := h.clientGetEntity(id)
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "Alpha", res.Data.Payload["title"])
	assert.Equal(t, "alpha@example.com", res.Data.Payload["email"])
	assert.False(t, payloadHas(res.Data.Payload, "secret"))
	assert.False(t, payloadHas(res.Data.Payload, "scratch"))
	assert.Empty(t, res.Data.Status)

	payload := adminEntityPayload()
	delete(payload, "secret")
	payload["secret"] = "nope"
	_, writeStatus, writeErr := h.client.DataSync.SetEntity().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assertServerCode(t, writeErr, writeStatus, 403)

	scratch := adminEntityPayload()
	delete(scratch, "secret")
	scratch["scratch"] = "nope"
	_, scratchStatus, scratchErr := h.client.DataSync.SetEntity().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(scratch).
		Execute()
	assertServerCode(t, scratchErr, scratchStatus, 403)

	_, patchStatus, patchErr := h.client.DataSync.UpdateEntity().
		ID(id).
		Replace("/payload/secret", "nope").
		Execute()
	assertServerCode(t, patchErr, patchStatus, 403)
}

func TestDataSyncLiveProjectionAdmin(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("ent")
	h.createEntity(h.admin, id, kitchenSinkEntityPayload(), "active")
	h.grantProjection("admin")

	res, status, err := h.clientGetEntity(id)
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "s3cret", res.Data.Payload["secret"])
	assert.Equal(t, "alpha@example.com", res.Data.Payload["email"])
	assert.False(t, payloadHas(res.Data.Payload, "scratch"))
	assert.Empty(t, res.Data.Status)
}

func TestDataSyncLiveProjectionReplaceDoesNotWipeHiddenFields(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("ent")
	// Do not store /status: a named projection that does not include it rejects
	// both setting and omitting a stored status (DS-0202).
	h.createEntity(h.admin, id, kitchenSinkEntityPayload(), "")
	h.grantProjection("private")

	payload := adminEntityPayload()
	delete(payload, "secret")
	payload["title"] = "PrivateEdit"
	payload["email"] = "private@example.com"

	res, status, err := h.client.DataSync.SetEntity().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Payload(payload).
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	if res == nil {
		t.Fatal("SetEntity returned a nil response")
	}
	assert.Equal(t, "PrivateEdit", res.Data.Payload["title"])
	assert.False(t, payloadHas(res.Data.Payload, "secret"))

	oracle := h.adminGetEntity(id)
	assert.Equal(t, "s3cret", oracle.Payload["secret"])
	assert.Equal(t, "PrivateEdit", oracle.Payload["title"])
}

func TestDataSyncLiveProjectionStatusEntity(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("st")
	h.createStatusEntity(h.admin, id, "Visible", "active")

	h.grantProjection(pubnub.PNDataSyncDefaultProjection)
	def, status, err := h.clientGetEntity(id)
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "active", def.Data.Status)
	assert.Equal(t, "Visible", def.Data.Payload["label"])

	h.grantProjection("admin")
	adm, status, err := h.clientGetEntity(id)
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "active", adm.Data.Status)

	h.grantProjection("private")
	priv, status, err := h.clientGetEntity(id)
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Empty(t, priv.Data.Status)
	assert.Equal(t, "Visible", priv.Data.Payload["label"])

	_, writeStatus, writeErr := h.client.DataSync.SetEntity().
		ID(id).
		EntityClassVersion(dsClassVersion).
		Status("inactive").
		Payload(map[string]interface{}{"label": "Visible"}).
		Execute()
	assertServerCode(t, writeErr, writeStatus, 403)
}

func TestDataSyncLiveProjectionRelationship(t *testing.T) {
	h := newDSLive(t)
	aID := dsNewID("ent")
	bID := dsNewID("st")
	h.createEntity(h.admin, aID, kitchenSinkEntityPayload(), "active")
	h.createStatusEntity(h.admin, bID, "Rel", "active")
	relID := dsNewID("rel")
	h.createRelationship(h.admin, relID, aID, bID, kitchenSinkRelPayload(), "active")

	h.grantProjection(pubnub.PNDataSyncDefaultProjection)
	def, status, err := h.client.DataSync.GetRelationship().ID(relID).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "owner", def.Data.Payload["role"])
	assert.False(t, payloadHas(def.Data.Payload, "reason"))
	assert.False(t, payloadHas(def.Data.Payload, "secret"))

	h.grantProjection("private")
	priv, status, err := h.client.DataSync.GetRelationship().ID(relID).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "because", priv.Data.Payload["reason"])
	assert.False(t, payloadHas(priv.Data.Payload, "secret"))

	_, writeStatus, writeErr := h.client.DataSync.UpdateRelationship().
		ID(relID).
		Replace("/payload/secret", "nope").
		Execute()
	assertServerCode(t, writeErr, writeStatus, 403)

	h.grantProjection("admin")
	adm, status, err := h.client.DataSync.GetRelationship().ID(relID).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "rel-secret", adm.Data.Payload["secret"])
	assert.Equal(t, "because", adm.Data.Payload["reason"])
}
