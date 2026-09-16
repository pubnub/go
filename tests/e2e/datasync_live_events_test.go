package e2e

import (
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v9"
	"github.com/stretchr/testify/assert"
)

func TestDataSyncLiveEventsEntityLifecycle(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("ent")
	sub := h.subscribeDataSync([]string{id, "__admin__" + id})

	created := h.createEntity(h.client, id, adminEntityPayload(), "")
	createEv := waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventCreate &&
			ev.Type == pubnub.PNDataSyncEventTypeEntity &&
			ev.Entity != nil && ev.Entity.ID == id
	})
	assert.Equal(t, "data-sync", createEv.Source)
	assert.Equal(t, dsEntityClass, createEv.ClassName)
	assert.Equal(t, pubnub.PNEntityClassLevelSubKey, createEv.ClassLevel)
	assert.Equal(t, dsClassVersion, createEv.ClassVersion)
	assert.Equal(t, "Alpha", createEv.Entity.Payload["title"])

	_, status, err := h.client.DataSync.UpdateEntity().
		ID(id).
		IfMatchETag(created.ETag).
		Replace("/payload/title", "LiveUpdated").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	updateEv := waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventUpdate &&
			ev.Type == pubnub.PNDataSyncEventTypeEntity &&
			ev.Entity != nil && ev.Entity.ID == id &&
			ev.Entity.Payload["title"] == "LiveUpdated"
	})
	assert.Equal(t, dsEntityClass, updateEv.ClassName)

	current := h.adminGetEntity(id)
	_, status, err = h.client.DataSync.RemoveEntity().ID(id).IfMatchETag(current.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	delEv := waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventDelete &&
			ev.Type == pubnub.PNDataSyncEventTypeEntity &&
			ev.ID == id
	})
	assert.NotEmpty(t, delEv.DeletedAt)
	assert.Nil(t, delEv.Entity)
}

func TestDataSyncLiveEventsProjectionChannels(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("ent")
	adminCh := "__admin__" + id
	h.createEntity(h.admin, id, kitchenSinkEntityPayload(), "active")
	h.grantHappy()

	// One client / one listener: events from every subscribed channel are
	// delivered on the same DataSyncEvent stream, so match on Channel.
	sub := h.subscribeDataSync([]string{id, adminCh})
	time.Sleep(500 * time.Millisecond)

	current := h.adminGetEntity(id)
	_, status, err := h.admin.DataSync.UpdateEntity().
		ID(id).
		IfMatchETag(current.ETag).
		Replace("/payload/title", "Projected").
		Replace("/payload/email", "proj@example.com").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	// Both projection views arrive on the same listener; collect them
	// together so one is not discarded while waiting for the other.
	var defEv, adminEv *pubnub.PNDataSyncEventResult
	deadline := time.After(dsEventTimeout)
	for defEv == nil || adminEv == nil {
		select {
		case ev := <-sub.events:
			if ev.Event != pubnub.PNDataSyncEventUpdate ||
				ev.Type != pubnub.PNDataSyncEventTypeEntity ||
				ev.Entity == nil || ev.Entity.ID != id {
				continue
			}
			switch ev.Channel {
			case id:
				defEv = ev
			case adminCh:
				adminEv = ev
			}
		case <-deadline:
			t.Fatalf("timeout waiting for projection events default=%v admin=%v", defEv != nil, adminEv != nil)
		}
	}

	assert.Equal(t, "Projected", defEv.Entity.Payload["title"])
	assert.False(t, payloadHas(defEv.Entity.Payload, "email"))
	assert.False(t, payloadHas(defEv.Entity.Payload, "secret"))

	assert.Equal(t, "Projected", adminEv.Entity.Payload["title"])
	assert.Equal(t, "proj@example.com", adminEv.Entity.Payload["email"])
	assert.Equal(t, "s3cret", adminEv.Entity.Payload["secret"])
}

func TestDataSyncLiveEventsRelationshipFanout(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	aID := dsNewID("ent")
	bID := dsNewID("st")
	h.createEntity(h.client, aID, adminEntityPayload(), "")
	h.createStatusEntity(h.client, bID, "Fanout", "active")

	sub := h.subscribeDataSync([]string{aID, bID})
	relID := dsNewID("rel")
	created := h.createRelationship(h.client, relID, aID, bID, adminRelPayload(), "")

	sawA := false
	sawB := false
	deadline := time.After(dsEventTimeout)
	for !sawA || !sawB {
		select {
		case ev := <-sub.events:
			if ev.Event != pubnub.PNDataSyncEventCreate || ev.Type != pubnub.PNDataSyncEventTypeRelationship {
				continue
			}
			if ev.Relationship == nil || ev.Relationship.ID != relID {
				continue
			}
			assert.Equal(t, aID, ev.Relationship.EntityAID)
			assert.Equal(t, bID, ev.Relationship.EntityBID)
			assert.Equal(t, dsRelClass, ev.ClassName)
			switch ev.Channel {
			case aID:
				sawA = true
			case bID:
				sawB = true
			}
		case <-deadline:
			t.Fatalf("timeout waiting for relationship create fan-out a=%v b=%v", sawA, sawB)
		}
	}

	_, status, err := h.client.DataSync.UpdateRelationship().
		ID(relID).
		IfMatchETag(created.ETag).
		Replace("/payload/role", "editor").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventUpdate &&
			ev.Type == pubnub.PNDataSyncEventTypeRelationship &&
			ev.Relationship != nil && ev.Relationship.ID == relID &&
			ev.Relationship.Payload["role"] == "editor"
	})

	current, _, err := h.client.DataSync.GetRelationship().ID(relID).Execute()
	assert.Nil(t, err)
	_, status, err = h.client.DataSync.RemoveRelationship().ID(relID).IfMatchETag(current.Data.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	delEv := waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventDelete &&
			ev.Type == pubnub.PNDataSyncEventTypeRelationship &&
			ev.ID == relID
	})
	assert.NotEmpty(t, delEv.DeletedAt)
	assert.Nil(t, delEv.Relationship)
}

func TestDataSyncLiveEventsSubscribeDenied(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("ent")
	h.createEntity(h.admin, id, kitchenSinkEntityPayload(), "active")
	h.grant(dsGrant{
		EntityPatterns:     map[string]pubnub.DataSyncPermissions{dsIDPattern(): dsCRUD()},
		EntityProjPatterns: map[string]string{dsIDPattern(): "admin"},
	})

	listener := pubnub.NewListener()
	denied := make(chan struct{}, 1)
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-stop:
				return
			case st := <-listener.Status:
				if st.Category == pubnub.PNAccessDeniedCategory {
					select {
					case denied <- struct{}{}:
					default:
					}
				}
			case <-listener.DataSyncEvent:
			case <-listener.Message:
			}
		}
	}()
	h.client.AddListener(listener)
	h.client.Subscribe().Channels([]string{id}).Execute()
	t.Cleanup(func() {
		close(stop)
		h.client.Unsubscribe().Channels([]string{id}).Execute()
		h.client.RemoveListener(listener)
	})

	select {
	case <-denied:
	case <-time.After(10 * time.Second):
		t.Fatal("expected PNAccessDeniedCategory when subscribing without channel Read")
	}
}
