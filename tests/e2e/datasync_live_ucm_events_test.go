package e2e

import (
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v10"
	"github.com/stretchr/testify/assert"
)

func TestDataSyncLiveEventsUserLifecycle(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("usr")
	sub := h.subscribeDataSync([]string{id, "__admin__" + id})

	created := h.createGoE2EUser(h.admin, id, kitchenSinkGoE2EUserPayload(), "")
	createEv := waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventCreate &&
			ev.Type == pubnub.PNDataSyncEventTypeUser &&
			ev.Entity != nil && ev.Entity.ID == id
	})
	assert.Equal(t, "data-sync", createEv.Source)
	assert.Equal(t, dsGoE2EUserClass, createEv.ClassName)
	assert.Equal(t, pubnub.PNEntityClassLevelSubKey, createEv.ClassLevel)
	assert.Equal(t, dsClassVersion, createEv.ClassVersion)
	assert.Equal(t, "Alpha", createEv.Entity.Payload["displayName"])

	_, status, err := h.client.DataSync.UpdateUser().
		ID(id).
		IfMatchETag(created.ETag).
		Replace("/payload/displayName", "LiveUpdated").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	updateEv := waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventUpdate &&
			ev.Type == pubnub.PNDataSyncEventTypeUser &&
			ev.Entity != nil && ev.Entity.ID == id &&
			ev.Entity.Payload["displayName"] == "LiveUpdated"
	})
	assert.Equal(t, dsGoE2EUserClass, updateEv.ClassName)

	current := h.adminGetUser(id)
	_, status, err = h.client.DataSync.RemoveUser().ID(id).IfMatchETag(current.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	delEv := waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventDelete &&
			ev.Type == pubnub.PNDataSyncEventTypeUser &&
			ev.ID == id
	})
	assert.NotEmpty(t, delEv.DeletedAt)
	assert.Nil(t, delEv.Entity)
}

func TestDataSyncLiveEventsUserProjectionChannels(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("usr")
	adminCh := "__admin__" + id
	h.createGoE2EUser(h.admin, id, kitchenSinkGoE2EUserPayload(), "active")
	h.grantHappy()

	sub := h.subscribeDataSync([]string{id, adminCh})
	time.Sleep(500 * time.Millisecond)

	current := h.adminGetUser(id)
	_, status, err := h.admin.DataSync.UpdateUser().
		ID(id).
		IfMatchETag(current.ETag).
		Replace("/payload/displayName", "Projected").
		Replace("/payload/email", "proj@example.com").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	var defEv, adminEv *pubnub.PNDataSyncEventResult
	deadline := time.After(dsEventTimeout)
	for adminEv == nil {
		select {
		case ev := <-sub.events:
			if ev.Event != pubnub.PNDataSyncEventUpdate ||
				(ev.Type != pubnub.PNDataSyncEventTypeUser && ev.Type != pubnub.PNDataSyncEventTypeEntity) ||
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
			t.Fatalf("timeout waiting for user admin projection event admin=%v default=%v", adminEv != nil, defEv != nil)
		}
	}

	assert.Equal(t, "Projected", adminEv.Entity.Payload["displayName"])
	assert.Equal(t, "proj@example.com", adminEv.Entity.Payload["email"])
	assert.Equal(t, "s3cret", adminEv.Entity.Payload["secret"])
	if defEv != nil {
		assert.Equal(t, "Projected", defEv.Entity.Payload["displayName"])
		assert.False(t, payloadHas(defEv.Entity.Payload, "email"))
		assert.False(t, payloadHas(defEv.Entity.Payload, "secret"))
	}
}

func TestDataSyncLiveEventsChannelLifecycle(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("chn")
	sub := h.subscribeDataSync([]string{id})

	created := h.createGoE2EChannel(h.client, id, adminGoE2EChannelPayload(), "")
	createEv := waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventCreate &&
			ev.Type == pubnub.PNDataSyncEventTypeChannel &&
			ev.Entity != nil && ev.Entity.ID == id
	})
	assert.Equal(t, dsGoE2EChannelClass, createEv.ClassName)
	assert.Equal(t, pubnub.PNEntityClassLevelSubKey, createEv.ClassLevel)
	assert.Equal(t, "Alpha", createEv.Entity.Payload["topic"])

	_, status, err := h.client.DataSync.UpdateChannel().
		ID(id).
		IfMatchETag(created.ETag).
		Replace("/payload/topic", "LiveUpdated").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventUpdate &&
			ev.Type == pubnub.PNDataSyncEventTypeChannel &&
			ev.Entity != nil && ev.Entity.ID == id &&
			ev.Entity.Payload["topic"] == "LiveUpdated"
	})

	current := h.adminGetChannel(id)
	_, status, err = h.client.DataSync.RemoveChannel().ID(id).IfMatchETag(current.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	delEv := waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventDelete &&
			ev.Type == pubnub.PNDataSyncEventTypeChannel &&
			ev.ID == id
	})
	assert.NotEmpty(t, delEv.DeletedAt)
	assert.Nil(t, delEv.Entity)
}

func TestDataSyncLiveEventsGlobalUser(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	id := dsNewID("gusr")
	sub := h.subscribeDataSync([]string{id, "__admin__" + id})

	created := h.createGlobalUser(h.admin, id, globalUserPayload(), "")
	createEv := waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventCreate &&
			ev.Type == pubnub.PNDataSyncEventTypeUser &&
			((ev.Entity != nil && ev.Entity.ID == id) || ev.ID == id)
	})
	assert.Equal(t, dsUserClass, createEv.ClassName)
	assert.Equal(t, pubnub.PNEntityClassLevelGlobal, createEv.ClassLevel)

	_, status, err := h.client.DataSync.RemoveUser().ID(id).IfMatchETag(created.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventDelete &&
			ev.Type == pubnub.PNDataSyncEventTypeUser &&
			ev.ID == id
	})
}

func TestDataSyncLiveEventsMembershipFanout(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()

	userID := dsNewID("gusr")
	chID := dsNewID("gchn")
	h.createGlobalUser(h.admin, userID, globalUserPayload(), "")
	h.createGlobalChannel(h.admin, chID, globalChannelPayload(), "")

	sub := h.subscribeDataSync([]string{userID, chID})
	memID := dsNewID("mem")
	created := h.createMembership(h.client, memID, userID, chID, nil, "")

	sawUser := false
	sawCh := false
	deadline := time.After(dsEventTimeout)
	for !sawUser || !sawCh {
		select {
		case ev := <-sub.events:
			if ev.Event != pubnub.PNDataSyncEventCreate || ev.Type != pubnub.PNDataSyncEventTypeMembership {
				continue
			}
			if ev.Membership == nil || ev.Membership.ID != memID {
				continue
			}
			assert.Equal(t, userID, ev.Membership.UserID)
			assert.Equal(t, chID, ev.Membership.ChannelID)
			assert.Equal(t, dsMembershipClass, ev.ClassName)
			assert.Equal(t, pubnub.PNEntityClassLevelGlobal, ev.ClassLevel)
			switch ev.Channel {
			case userID:
				sawUser = true
			case chID:
				sawCh = true
			}
		case <-deadline:
			t.Fatalf("timeout waiting for membership create fan-out user=%v channel=%v", sawUser, sawCh)
		}
	}

	_, status, err := h.client.DataSync.UpdateMembership().
		ID(memID).
		IfMatchETag(created.ETag).
		Replace("/status", "inactive").
		Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventUpdate &&
			ev.Type == pubnub.PNDataSyncEventTypeMembership &&
			ev.Membership != nil && ev.Membership.ID == memID
	})

	current, _, err := h.client.DataSync.GetMembership().ID(memID).Execute()
	assert.Nil(t, err)
	_, status, err = h.client.DataSync.RemoveMembership().ID(memID).IfMatchETag(current.Data.ETag).Execute()
	assert.Nil(t, err)
	assert.Equal(t, 200, status.StatusCode)

	delEv := waitDSEvent(t, sub.events, func(ev *pubnub.PNDataSyncEventResult) bool {
		return ev.Event == pubnub.PNDataSyncEventDelete &&
			ev.Type == pubnub.PNDataSyncEventTypeMembership &&
			ev.ID == memID
	})
	assert.NotEmpty(t, delEv.DeletedAt)
	assert.Nil(t, delEv.Membership)
}

func TestDataSyncLiveEventsUserSubscribeDenied(t *testing.T) {
	h := newDSLive(t)
	id := dsNewID("usr")
	h.createGoE2EUser(h.admin, id, kitchenSinkGoE2EUserPayload(), "active")
	h.grant(dsGrant{
		UsersPatterns:    dsUserPerms(dsUUIDCRUD()),
		UserProjPatterns: map[string]string{dsIDPattern(): "admin"},
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
