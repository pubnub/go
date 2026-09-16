package e2e

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v9"
	"github.com/pubnub/go/v9/pnerr"
	"github.com/stretchr/testify/assert"
)

const (
	dsEntityClass   = "GoE2EEntity"
	dsStatusClass   = "GoE2EStatus"
	dsRelClass      = "GoE2ERel"
	dsClassVersion  = 1
	dsGrantTTL      = 30
	dsFilterTimeout = 15 * time.Second
	dsEventTimeout  = 15 * time.Second
)

func dsCRUD() pubnub.DataSyncPermissions {
	return pubnub.DataSyncPermissions{Get: true, Create: true, Update: true, Delete: true}
}

func dsIDPattern() string {
	return "goe2e-.*"
}

func skipWithoutDSKeys(t *testing.T) {
	t.Helper()
	if os.Getenv("DS_PUBLISH_KEY") == "" ||
		os.Getenv("DS_SUBSCRIBE_KEY") == "" ||
		os.Getenv("DS_SECRET_KEY") == "" {
		t.Skip("DS_PUBLISH_KEY / DS_SUBSCRIBE_KEY / DS_SECRET_KEY not set")
	}
}

func dsAdminConfig() *pubnub.Config {
	cfg := pubnub.NewConfigWithUserId(pubnub.UserId(pubnub.GenerateUUID()))
	cfg.PublishKey = os.Getenv("DS_PUBLISH_KEY")
	cfg.SubscribeKey = os.Getenv("DS_SUBSCRIBE_KEY")
	cfg.SecretKey = os.Getenv("DS_SECRET_KEY")
	return cfg
}

func dsClientConfig() *pubnub.Config {
	cfg := pubnub.NewConfigWithUserId(pubnub.UserId(pubnub.GenerateUUID()))
	cfg.PublishKey = os.Getenv("DS_PUBLISH_KEY")
	cfg.SubscribeKey = os.Getenv("DS_SUBSCRIBE_KEY")
	return cfg
}

// dsGrant describes a GrantToken request for the client under test.
type dsGrant struct {
	Entities             map[string]pubnub.DataSyncPermissions
	EntityPatterns       map[string]pubnub.DataSyncPermissions
	Relationships        map[string]pubnub.DataSyncPermissions
	RelationshipPatterns map[string]pubnub.DataSyncPermissions
	EntityProjections    map[string]string
	EntityProjPatterns   map[string]string
	RelProjections       map[string]string
	RelProjPatterns      map[string]string
	Channels             map[string]pubnub.ChannelPermissions
	ChannelPatterns      map[string]pubnub.ChannelPermissions
}

type dsLive struct {
	t        *testing.T
	admin    *pubnub.PubNub
	client   *pubnub.PubNub
	entities []string
	rels     []string
}

func newDSLive(t *testing.T) *dsLive {
	t.Helper()
	skipWithoutDSKeys(t)

	h := &dsLive{
		t:      t,
		admin:  pubnub.NewPubNub(dsAdminConfig()),
		client: pubnub.NewPubNub(dsClientConfig()),
	}
	t.Cleanup(func() {
		for i := len(h.rels) - 1; i >= 0; i-- {
			_, _, _ = h.admin.DataSync.RemoveRelationship().ID(h.rels[i]).Execute()
		}
		for i := len(h.entities) - 1; i >= 0; i-- {
			_, _, _ = h.admin.DataSync.RemoveEntity().ID(h.entities[i]).Execute()
		}
		h.client.Destroy()
		h.admin.Destroy()
	})
	return h
}

func (h *dsLive) trackEntity(id string) {
	if id == "" {
		return
	}
	h.entities = append(h.entities, id)
}

func (h *dsLive) trackRel(id string) {
	if id == "" {
		return
	}
	h.rels = append(h.rels, id)
}

func (h *dsLive) grant(g dsGrant) {
	h.t.Helper()
	builder := h.admin.GrantToken().
		TTL(dsGrantTTL).
		AuthorizedUUID(h.client.Config.UUID)

	if len(g.Channels) > 0 {
		builder = builder.Channels(g.Channels)
	}
	if len(g.ChannelPatterns) > 0 {
		builder = builder.ChannelsPattern(g.ChannelPatterns)
	}
	if len(g.Entities) > 0 || len(g.Relationships) > 0 {
		builder = builder.DataSync(pubnub.PNDataSyncTokenScopes{
			Entities:      g.Entities,
			Relationships: g.Relationships,
		})
	}
	if len(g.EntityPatterns) > 0 || len(g.RelationshipPatterns) > 0 {
		builder = builder.DataSyncPattern(pubnub.PNDataSyncTokenScopes{
			Entities:      g.EntityPatterns,
			Relationships: g.RelationshipPatterns,
		})
	}
	proj := pubnub.PNDataSyncProjections{
		Resources: pubnub.PNDataSyncProjectionScope{
			Entities:      g.EntityProjections,
			Relationships: g.RelProjections,
		},
		Patterns: pubnub.PNDataSyncProjectionScope{
			Entities:      g.EntityProjPatterns,
			Relationships: g.RelProjPatterns,
		},
	}
	if len(g.EntityProjections) > 0 || len(g.RelProjections) > 0 ||
		len(g.EntityProjPatterns) > 0 || len(g.RelProjPatterns) > 0 {
		builder = builder.DataSyncProjections(proj)
	}

	res, status, err := builder.Execute()
	if err != nil || res == nil || res.Data.Token == "" {
		h.t.Fatalf("GrantToken failed: err=%v status=%d body=%s", err, status.StatusCode, status.OriginalResponse)
	}
	h.client.SetToken(res.Data.Token)
}

func (h *dsLive) grantHappy() {
	h.t.Helper()
	// Named projections (admin/private) do not include GoE2EEntity / GoE2ERel
	// /status. Client writes under this grant must omit Status.
	h.grant(dsGrant{
		EntityPatterns:       map[string]pubnub.DataSyncPermissions{dsIDPattern(): dsCRUD()},
		RelationshipPatterns: map[string]pubnub.DataSyncPermissions{dsIDPattern(): dsCRUD()},
		EntityProjPatterns:   map[string]string{dsIDPattern(): "admin"},
		RelProjPatterns:      map[string]string{dsIDPattern(): "admin"},
		ChannelPatterns: map[string]pubnub.ChannelPermissions{
			"^goe2e-.*$":            {Read: true},
			"^__admin__goe2e-.*$":   {Read: true},
			"^__private__goe2e-.*$": {Read: true},
		},
	})
}

func (h *dsLive) grantProjection(name string) {
	h.t.Helper()
	h.grant(dsGrant{
		EntityPatterns:       map[string]pubnub.DataSyncPermissions{dsIDPattern(): dsCRUD()},
		RelationshipPatterns: map[string]pubnub.DataSyncPermissions{dsIDPattern(): dsCRUD()},
		EntityProjPatterns:   map[string]string{dsIDPattern(): name},
		RelProjPatterns:      map[string]string{dsIDPattern(): name},
		ChannelPatterns: map[string]pubnub.ChannelPermissions{
			"^goe2e-.*$":            {Read: true},
			"^__admin__goe2e-.*$":   {Read: true},
			"^__private__goe2e-.*$": {Read: true},
		},
	})
}

func dsNewID(kind string) string {
	return randomized("goe2e-" + kind)
}

// kitchenSinkEntityPayload is the GoE2EEntity instance template. Callers that
// operate under a named projection must drop undeclared fields (scratch) and
// any field outside that projection.
func kitchenSinkEntityPayload() map[string]interface{} {
	return map[string]interface{}{
		"title":     "Alpha",
		"category":  "tools",
		"score":     float64(10),
		"active":    true,
		"bornOn":    "2020-01-15",
		"updatedTs": "2026-09-01T12:00:00Z",
		"notes":     "not searchable",
		"address":   map[string]interface{}{"city": "Portland"},
		"email":     "alpha@example.com",
		"secret":    "s3cret",
		"scratch":   "undeclared",
	}
}

func adminEntityPayload() map[string]interface{} {
	p := kitchenSinkEntityPayload()
	delete(p, "scratch")
	return p
}

func defaultEntityPayload() map[string]interface{} {
	p := kitchenSinkEntityPayload()
	delete(p, "email")
	delete(p, "secret")
	return p
}

func kitchenSinkRelPayload() map[string]interface{} {
	return map[string]interface{}{
		"role":      "owner",
		"kind":      "link",
		"weight":    1.5,
		"confirmed": true,
		"startedOn": "2021-03-01",
		"linkedAt":  "2026-09-01T12:00:00Z",
		"comment":   "not searchable",
		"meta":      map[string]interface{}{"region": "west"},
		"reason":    "because",
		"secret":    "rel-secret",
	}
}

func adminRelPayload() map[string]interface{} {
	return kitchenSinkRelPayload()
}

func (h *dsLive) createEntity(pn *pubnub.PubNub, id string, payload map[string]interface{}, status string) *pubnub.PNEntity {
	h.t.Helper()
	b := pn.DataSync.CreateEntity().
		EntityClass(dsEntityClass).
		EntityClassVersion(dsClassVersion).
		EntityClassLevel(pubnub.PNEntityClassLevelSubKey).
		Payload(payload)
	if id != "" {
		b = b.ID(id)
	}
	if status != "" {
		b = b.Status(status)
	}
	res, st, err := b.Execute()
	if err != nil || res == nil {
		h.t.Fatalf("CreateEntity id=%s: err=%v status=%d body=%s", id, err, st.StatusCode, st.OriginalResponse)
	}
	h.trackEntity(res.Data.ID)
	return &res.Data
}

func (h *dsLive) createStatusEntity(pn *pubnub.PubNub, id, label, status string) *pubnub.PNEntity {
	h.t.Helper()
	b := pn.DataSync.CreateEntity().
		EntityClass(dsStatusClass).
		EntityClassVersion(dsClassVersion).
		EntityClassLevel(pubnub.PNEntityClassLevelSubKey).
		Payload(map[string]interface{}{"label": label})
	if id != "" {
		b = b.ID(id)
	}
	if status != "" {
		b = b.Status(status)
	}
	res, st, err := b.Execute()
	if err != nil || res == nil {
		h.t.Fatalf("CreateEntity status class id=%s: err=%v status=%d body=%s", id, err, st.StatusCode, st.OriginalResponse)
	}
	h.trackEntity(res.Data.ID)
	return &res.Data
}

func (h *dsLive) createRelationship(pn *pubnub.PubNub, id, aID, bID string, payload map[string]interface{}, status string) *pubnub.PNRelationship {
	h.t.Helper()
	b := pn.DataSync.CreateRelationship().
		EntityAID(aID).
		EntityBID(bID).
		RelationshipClass(dsRelClass).
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
		h.t.Fatalf("CreateRelationship id=%s: err=%v status=%d body=%s", id, err, st.StatusCode, st.OriginalResponse)
	}
	h.trackRel(res.Data.ID)
	return &res.Data
}

func (h *dsLive) adminGetEntity(id string) *pubnub.PNEntity {
	h.t.Helper()
	res, st, err := h.admin.DataSync.GetEntity().ID(id).Execute()
	if err != nil || res == nil {
		h.t.Fatalf("admin GetEntity %s: err=%v status=%d body=%s", id, err, st.StatusCode, st.OriginalResponse)
	}
	return &res.Data
}

func (h *dsLive) clientGetEntity(id string) (*pubnub.PNEntityResponse, pubnub.StatusResponse, error) {
	h.t.Helper()
	return h.client.DataSync.GetEntity().ID(id).Execute()
}

func assertServerCode(t *testing.T, err error, status pubnub.StatusResponse, want int) {
	t.Helper()
	assert.NotNil(t, err)
	assert.Equal(t, want, status.StatusCode)
	var se *pnerr.ServerError
	if errors.As(err, &se) {
		assert.Equal(t, want, se.StatusCode)
	}
}

func payloadHas(p map[string]interface{}, key string) bool {
	if p == nil {
		return false
	}
	_, ok := p[key]
	return ok
}

func nestedCity(p map[string]interface{}) string {
	if p == nil {
		return ""
	}
	addr, _ := p["address"].(map[string]interface{})
	if addr == nil {
		return ""
	}
	s, _ := addr["city"].(string)
	return s
}

func nestedRegion(p map[string]interface{}) string {
	if p == nil {
		return ""
	}
	meta, _ := p["meta"].(map[string]interface{})
	if meta == nil {
		return ""
	}
	s, _ := meta["region"].(string)
	return s
}

func waitEntitiesFilterFast(t *testing.T, pn *pubnub.PubNub, filter string, want int) *pubnub.PNEntitiesResponse {
	t.Helper()
	return eventually(t, dsFilterTimeout, 250*time.Millisecond,
		fmt.Sprintf("GetEntities FilterFast %q returns %d", filter, want),
		func() (*pubnub.PNEntitiesResponse, bool) {
			r, _, err := pn.DataSync.GetEntities().
				EntityClass(dsEntityClass).
				EntityClassLevel(pubnub.PNEntityClassLevelSubKey).
				Limit(100).
				FilterFast(filter).
				Execute()
			return r, err == nil && r != nil && len(r.Data) == want
		})
}

func waitEntitiesFilter(t *testing.T, pn *pubnub.PubNub, filter string, want int) *pubnub.PNEntitiesResponse {
	t.Helper()
	return eventually(t, dsFilterTimeout, 250*time.Millisecond,
		fmt.Sprintf("GetEntities Filter %q returns %d", filter, want),
		func() (*pubnub.PNEntitiesResponse, bool) {
			r, _, err := pn.DataSync.GetEntities().
				EntityClass(dsEntityClass).
				EntityClassLevel(pubnub.PNEntityClassLevelSubKey).
				Limit(100).
				Filter(filter).
				Execute()
			return r, err == nil && r != nil && len(r.Data) == want
		})
}

func waitRelationshipsFilterFast(t *testing.T, pn *pubnub.PubNub, filter string, want int) *pubnub.PNRelationshipsResponse {
	t.Helper()
	return eventually(t, dsFilterTimeout, 250*time.Millisecond,
		fmt.Sprintf("GetRelationships FilterFast %q returns %d", filter, want),
		func() (*pubnub.PNRelationshipsResponse, bool) {
			r, _, err := pn.DataSync.GetRelationships().
				RelationshipClass(dsRelClass).
				Limit(100).
				FilterFast(filter).
				Execute()
			return r, err == nil && r != nil && len(r.Data) == want
		})
}

type dsEventSub struct {
	events chan *pubnub.PNDataSyncEventResult
	stop   chan struct{}
	pn     *pubnub.PubNub
	l      *pubnub.Listener
	chs    []string
}

func (h *dsLive) subscribeDataSync(channels []string) *dsEventSub {
	h.t.Helper()
	listener := pubnub.NewListener()
	connected := make(chan struct{}, 1)
	denied := make(chan *pubnub.PNStatus, 1)
	events := make(chan *pubnub.PNDataSyncEventResult, 32)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				return
			case st := <-listener.Status:
				switch st.Category {
				case pubnub.PNConnectedCategory:
					select {
					case connected <- struct{}{}:
					default:
					}
				case pubnub.PNAccessDeniedCategory:
					select {
					case denied <- st:
					default:
					}
				}
			case ev := <-listener.DataSyncEvent:
				select {
				case events <- ev:
				case <-stop:
					return
				}
			case <-listener.Message:
			case <-listener.Presence:
			}
		}
	}()

	h.client.AddListener(listener)
	h.client.Subscribe().Channels(channels).Execute()

	select {
	case <-connected:
	case st := <-denied:
		h.t.Fatalf("subscribe access denied: %+v", st)
	case <-time.After(10 * time.Second):
		h.t.Fatalf("timeout waiting for subscribe connected on %v", channels)
	}

	sub := &dsEventSub{events: events, stop: stop, pn: h.client, l: listener, chs: channels}
	h.t.Cleanup(sub.close)
	return sub
}

func (s *dsEventSub) close() {
	select {
	case <-s.stop:
	default:
		close(s.stop)
	}
	s.pn.Unsubscribe().Channels(s.chs).Execute()
	s.pn.RemoveListener(s.l)
}

func waitDSEvent(t *testing.T, events <-chan *pubnub.PNDataSyncEventResult, match func(*pubnub.PNDataSyncEventResult) bool) *pubnub.PNDataSyncEventResult {
	t.Helper()
	deadline := time.After(dsEventTimeout)
	var last *pubnub.PNDataSyncEventResult
	for {
		select {
		case ev := <-events:
			last = ev
			if match(ev) {
				return ev
			}
		case <-deadline:
			t.Fatalf("timeout waiting for DataSync event (last=%#v). If subscribe connected but no event arrived, enable DataSync event rules for GoE2EEntity/GoE2ERel v1 on the keyset", last)
			return last
		}
	}
}
