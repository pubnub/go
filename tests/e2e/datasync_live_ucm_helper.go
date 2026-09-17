package e2e

import (
	"fmt"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v9"
)

func globalUserPayload() map[string]interface{} {
	return map[string]interface{}{
		"name": "Ada",
		"type": "member",
	}
}

func kitchenSinkGoE2EUserPayload() map[string]interface{} {
	return map[string]interface{}{
		"name":        "Ada",
		"type":        "member",
		"displayName": "Alpha",
		"category":    "tools",
		"score":       float64(10),
		"active":      true,
		"bornOn":      "2020-01-15",
		"updatedTs":   "2026-09-01T12:00:00Z",
		"notes":       "not searchable",
		"address":     map[string]interface{}{"city": "Portland"},
		"email":       "alpha@example.com",
		"secret":      "s3cret",
		"scratch":     "undeclared",
	}
}

func adminGoE2EUserPayload() map[string]interface{} {
	p := kitchenSinkGoE2EUserPayload()
	delete(p, "scratch")
	return p
}

func defaultGoE2EUserPayload() map[string]interface{} {
	p := kitchenSinkGoE2EUserPayload()
	delete(p, "email")
	delete(p, "secret")
	return p
}

func globalChannelPayload() map[string]interface{} {
	return map[string]interface{}{
		"name": "Lobby",
		"type": "public",
	}
}

func kitchenSinkGoE2EChannelPayload() map[string]interface{} {
	return map[string]interface{}{
		"name":       "Lobby",
		"type":       "public",
		"topic":      "Alpha",
		"category":   "tools",
		"capacity":   float64(10),
		"live":       true,
		"launchedOn": "2020-01-15",
		"lastActive": "2026-09-01T12:00:00Z",
		"notes":      "not searchable",
		"address":    map[string]interface{}{"city": "Portland"},
		"email":      "alpha@example.com",
		"secret":     "s3cret",
		"scratch":    "undeclared",
	}
}

func adminGoE2EChannelPayload() map[string]interface{} {
	p := kitchenSinkGoE2EChannelPayload()
	delete(p, "scratch")
	return p
}

func defaultGoE2EChannelPayload() map[string]interface{} {
	p := kitchenSinkGoE2EChannelPayload()
	delete(p, "email")
	delete(p, "secret")
	return p
}

func (h *dsLive) createUser(pn *pubnub.PubNub, id, class string, level pubnub.PNEntityClassLevel, payload map[string]interface{}, status string) *pubnub.PNEntity {
	h.t.Helper()
	b := pn.DataSync.CreateUser().
		EntityClass(class).
		EntityClassVersion(dsClassVersion).
		EntityClassLevel(level).
		Payload(payload)
	if id != "" {
		b = b.ID(id)
	}
	if status != "" {
		b = b.Status(status)
	}
	res, st, err := b.Execute()
	if err != nil || res == nil {
		h.t.Fatalf("CreateUser id=%s class=%s: err=%v status=%d body=%s", id, class, err, st.StatusCode, st.OriginalResponse)
	}
	h.trackUser(res.Data.ID)
	return &res.Data
}

func (h *dsLive) createGlobalUser(pn *pubnub.PubNub, id string, payload map[string]interface{}, status string) *pubnub.PNEntity {
	h.t.Helper()
	return h.createUser(pn, id, dsUserClass, pubnub.PNEntityClassLevelGlobal, payload, status)
}

func (h *dsLive) createGoE2EUser(pn *pubnub.PubNub, id string, payload map[string]interface{}, status string) *pubnub.PNEntity {
	h.t.Helper()
	return h.createUser(pn, id, dsGoE2EUserClass, pubnub.PNEntityClassLevelSubKey, payload, status)
}

func (h *dsLive) createChannel(pn *pubnub.PubNub, id, class string, level pubnub.PNEntityClassLevel, payload map[string]interface{}, status string) *pubnub.PNEntity {
	h.t.Helper()
	b := pn.DataSync.CreateChannel().
		EntityClass(class).
		EntityClassVersion(dsClassVersion).
		EntityClassLevel(level).
		Payload(payload)
	if id != "" {
		b = b.ID(id)
	}
	if status != "" {
		b = b.Status(status)
	}
	res, st, err := b.Execute()
	if err != nil || res == nil {
		h.t.Fatalf("CreateChannel id=%s class=%s: err=%v status=%d body=%s", id, class, err, st.StatusCode, st.OriginalResponse)
	}
	h.trackChannel(res.Data.ID)
	return &res.Data
}

func (h *dsLive) createGlobalChannel(pn *pubnub.PubNub, id string, payload map[string]interface{}, status string) *pubnub.PNEntity {
	h.t.Helper()
	return h.createChannel(pn, id, dsChannelClass, pubnub.PNEntityClassLevelGlobal, payload, status)
}

func (h *dsLive) createGoE2EChannel(pn *pubnub.PubNub, id string, payload map[string]interface{}, status string) *pubnub.PNEntity {
	h.t.Helper()
	return h.createChannel(pn, id, dsGoE2EChannelClass, pubnub.PNEntityClassLevelSubKey, payload, status)
}

func (h *dsLive) createMembership(pn *pubnub.PubNub, id, userID, channelID string, payload map[string]interface{}, status string) *pubnub.PNDataSyncMembership {
	h.t.Helper()
	b := pn.DataSync.CreateMembership().
		UserID(userID).
		ChannelID(channelID).
		RelationshipClassVersion(dsClassVersion).
		Payload(payload)
	if id != "" {
		b = b.ID(id)
	}
	if status != "" {
		b = b.Status(status)
	}
	res, st, err := b.Execute()
	if err != nil || res == nil {
		h.t.Fatalf("CreateMembership id=%s user=%s channel=%s: err=%v status=%d body=%s", id, userID, channelID, err, st.StatusCode, st.OriginalResponse)
	}
	h.trackMembership(res.Data.ID)
	return &res.Data
}

func (h *dsLive) adminGetUser(id string) *pubnub.PNEntity {
	h.t.Helper()
	res, st, err := h.admin.DataSync.GetUser().ID(id).Execute()
	if err != nil || res == nil {
		h.t.Fatalf("admin GetUser %s: err=%v status=%d body=%s", id, err, st.StatusCode, st.OriginalResponse)
	}
	return &res.Data
}

func (h *dsLive) adminGetChannel(id string) *pubnub.PNEntity {
	h.t.Helper()
	res, st, err := h.admin.DataSync.GetChannel().ID(id).Execute()
	if err != nil || res == nil {
		h.t.Fatalf("admin GetChannel %s: err=%v status=%d body=%s", id, err, st.StatusCode, st.OriginalResponse)
	}
	return &res.Data
}

func (h *dsLive) adminGetMembership(id string) *pubnub.PNDataSyncMembership {
	h.t.Helper()
	res, st, err := h.admin.DataSync.GetMembership().ID(id).Execute()
	if err != nil || res == nil {
		h.t.Fatalf("admin GetMembership %s: err=%v status=%d body=%s", id, err, st.StatusCode, st.OriginalResponse)
	}
	return &res.Data
}

func waitUsersFilterFast(t *testing.T, pn *pubnub.PubNub, class string, level pubnub.PNEntityClassLevel, filter string, want int) *pubnub.PNUsersResponse {
	t.Helper()
	return eventually(t, dsFilterTimeout, 250*time.Millisecond,
		fmt.Sprintf("GetUsers class=%s FilterFast %q returns %d", class, filter, want),
		func() (*pubnub.PNUsersResponse, bool) {
			r, _, err := pn.DataSync.GetUsers().
				EntityClass(class).
				EntityClassLevel(level).
				Limit(100).
				FilterFast(filter).
				Execute()
			return r, err == nil && r != nil && len(r.Data) == want
		})
}

func waitUsersFilter(t *testing.T, pn *pubnub.PubNub, class string, level pubnub.PNEntityClassLevel, filter string, want int) *pubnub.PNUsersResponse {
	t.Helper()
	return eventually(t, dsFilterTimeout, 250*time.Millisecond,
		fmt.Sprintf("GetUsers class=%s Filter %q returns %d", class, filter, want),
		func() (*pubnub.PNUsersResponse, bool) {
			r, _, err := pn.DataSync.GetUsers().
				EntityClass(class).
				EntityClassLevel(level).
				Limit(100).
				Filter(filter).
				Execute()
			return r, err == nil && r != nil && len(r.Data) == want
		})
}

func waitChannelsFilterFast(t *testing.T, pn *pubnub.PubNub, class string, level pubnub.PNEntityClassLevel, filter string, want int) *pubnub.PNDataSyncChannelsResponse {
	t.Helper()
	return eventually(t, dsFilterTimeout, 250*time.Millisecond,
		fmt.Sprintf("GetChannels class=%s FilterFast %q returns %d", class, filter, want),
		func() (*pubnub.PNDataSyncChannelsResponse, bool) {
			r, _, err := pn.DataSync.GetChannels().
				EntityClass(class).
				EntityClassLevel(level).
				Limit(100).
				FilterFast(filter).
				Execute()
			return r, err == nil && r != nil && len(r.Data) == want
		})
}

func waitChannelsFilter(t *testing.T, pn *pubnub.PubNub, class string, level pubnub.PNEntityClassLevel, filter string, want int) *pubnub.PNDataSyncChannelsResponse {
	t.Helper()
	return eventually(t, dsFilterTimeout, 250*time.Millisecond,
		fmt.Sprintf("GetChannels class=%s Filter %q returns %d", class, filter, want),
		func() (*pubnub.PNDataSyncChannelsResponse, bool) {
			r, _, err := pn.DataSync.GetChannels().
				EntityClass(class).
				EntityClassLevel(level).
				Limit(100).
				Filter(filter).
				Execute()
			return r, err == nil && r != nil && len(r.Data) == want
		})
}
