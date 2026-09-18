package e2e

import (
	"fmt"
	"strings"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v9"
	"github.com/stretchr/testify/assert"
)

var queryTitles = []string{"Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot", "Golf", "Hotel"}

type dsQueryCorpus struct {
	prefix string
	entity []pubnub.PNEntity
	status []pubnub.PNEntity
	rel    []pubnub.PNRelationship
}

func corpusTitleFilter(prefix string) string {
	parts := make([]string, len(queryTitles))
	for i, title := range queryTitles {
		parts[i] = fmt.Sprintf(`title == "%s-%s"`, prefix, title)
	}
	return "(" + strings.Join(parts, " || ") + ")"
}

func seedQueryCorpus(t *testing.T, h *dsLive) dsQueryCorpus {
	t.Helper()
	prefix := randomized("goe2e-q")
	titles := queryTitles
	categories := []string{"tools", "parts", "tools", "parts", "tools", "parts", "tools", "parts"}
	cities := []string{"Portland", "Seattle", "Portland", "Austin", "Seattle", "Portland", "Austin", "Seattle"}
	actives := []bool{true, true, false, true, false, true, false, true}

	c := dsQueryCorpus{prefix: prefix, entity: make([]pubnub.PNEntity, 0, 8)}
	for i, title := range titles {
		p := kitchenSinkEntityPayload()
		p["title"] = prefix + "-" + title
		p["category"] = categories[i]
		p["score"] = float64((i + 1) * 10)
		p["active"] = actives[i]
		p["address"] = map[string]interface{}{"city": cities[i]}
		p["email"] = fmt.Sprintf("%s-%s@example.com", prefix, title)
		p["updatedTs"] = fmt.Sprintf("2026-09-%02dT12:00:00Z", i+1)
		ent := h.createEntity(h.admin, fmt.Sprintf("%s-ent-%d", prefix, i), p, "active")
		c.entity = append(c.entity, *ent)
	}

	for i, label := range []string{"One", "Two", "Three"} {
		st := h.createStatusEntity(h.admin, fmt.Sprintf("%s-st-%d", prefix, i), prefix+"-"+label, "active")
		c.status = append(c.status, *st)
	}

	kinds := []string{"link", "bind", "link", "bind", "link"}
	roles := []string{"owner", "editor", "viewer", "owner", "editor"}
	for i := 0; i < 5; i++ {
		p := kitchenSinkRelPayload()
		p["role"] = prefix + "-" + roles[i]
		p["kind"] = prefix + "-" + kinds[i]
		p["weight"] = float64((i + 1) * 10)
		p["confirmed"] = i%2 == 0
		p["reason"] = fmt.Sprintf("reason-%d", i)
		rel := h.createRelationship(h.admin, fmt.Sprintf("%s-rel-%d", prefix, i),
			c.entity[i].ID, c.status[i%len(c.status)].ID, p, "active")
		c.rel = append(c.rel, *rel)
	}

	titleFilter := corpusTitleFilter(prefix)
	waitEntitiesFilterFast(t, h.client, titleFilter, 8)
	waitEntitiesFilter(t, h.client, fmt.Sprintf(`title == "%s-Hotel"`, prefix), 1)
	_ = eventually(t, dsFilterTimeout, 250*time.Millisecond,
		"relationship corpus visible by entityAId",
		func() (*pubnub.PNRelationshipsResponse, bool) {
			r, _, err := h.client.DataSync.GetRelationships().
				RelationshipClass(dsRelClass).
				EntityAID(c.entity[4].ID).
				Execute()
			return r, err == nil && r != nil && len(r.Data) == 1
		})
	return c
}

func TestDataSyncLiveQuery(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()
	c := seedQueryCorpus(t, h)
	titleLike := corpusTitleFilter(c.prefix)

	t.Run("EntityClassLevelAndVersion", func(t *testing.T) {
		res, status, err := h.client.DataSync.GetEntities().
			EntityClass(dsEntityClass).
			EntityClassVersion(dsClassVersion).
			EntityClassLevel(pubnub.PNEntityClassLevelSubKey).
			FilterFast(titleLike).
			Limit(100).
			Execute()
		assert.Nil(t, err)
		assert.Equal(t, 200, status.StatusCode)
		assert.Len(t, res.Data, 8)
	})

	t.Run("EntityPagination", func(t *testing.T) {
		seen := map[string]bool{}
		cursor := ""
		pages := 0
		for {
			b := h.client.DataSync.GetEntities().
				EntityClass(dsEntityClass).
				EntityClassLevel(pubnub.PNEntityClassLevelSubKey).
				FilterFast(titleLike).
				Sort([]string{"score"}).
				Limit(2)
			if cursor != "" {
				b = b.Cursor(cursor)
			}
			res, status, err := b.Execute()
			assert.Nil(t, err)
			assert.Equal(t, 200, status.StatusCode)
			assert.LessOrEqual(t, len(res.Data), 2)
			for _, e := range res.Data {
				seen[e.ID] = true
			}
			pages++
			if res.Meta == nil || !res.Meta.HasNext {
				break
			}
			assert.NotEmpty(t, res.Meta.NextCursor)
			cursor = res.Meta.NextCursor
			if pages > 8 {
				t.Fatal("pagination did not terminate")
			}
		}
		assert.GreaterOrEqual(t, pages, 4)
		assert.Len(t, seen, 8)
	})

	t.Run("FilterFastSimpleAndFull", func(t *testing.T) {
		tools := waitEntitiesFilterFast(t, h.client, fmt.Sprintf(`category == "tools" && %s`, titleLike), 4)
		assert.Len(t, tools.Data, 4)
		for _, e := range tools.Data {
			assert.Equal(t, "tools", e.Payload["category"])
		}

		portland := waitEntitiesFilterFast(t, h.client, fmt.Sprintf(`city == "Portland" && %s`, titleLike), 3)
		assert.Len(t, portland.Data, 3)

		active := waitEntitiesFilterFast(t, h.client, fmt.Sprintf(`active == true && %s`, titleLike), 5)
		assert.Len(t, active.Data, 5)

		high := waitEntitiesFilterFast(t, h.client, fmt.Sprintf(`score > 35 && %s`, titleLike), 5)
		assert.Len(t, high.Data, 5)
	})

	t.Run("FilterFullOnly", func(t *testing.T) {
		alpha := waitEntitiesFilter(t, h.client, fmt.Sprintf(`title == "%s-Alpha"`, c.prefix), 1)
		assert.Equal(t, c.entity[0].ID, alpha.Data[0].ID)

		high := waitEntitiesFilter(t, h.client, fmt.Sprintf(`score >= 50 && %s`, titleLike), 4)
		assert.Len(t, high.Data, 4)

		email := waitEntitiesFilter(t, h.client, fmt.Sprintf(`email == "%s-Alpha@example.com"`, c.prefix), 1)
		assert.Equal(t, c.entity[0].ID, email.Data[0].ID)
	})

	t.Run("FilterRejectedFields", func(t *testing.T) {
		_, status, err := h.client.DataSync.GetEntities().
			EntityClass(dsEntityClass).
			Filter(fmt.Sprintf(`notes == "not searchable" && %s`, titleLike)).
			Execute()
		assertServerCode(t, err, status, 400)

		_, status, err = h.client.DataSync.GetEntities().
			EntityClass(dsEntityClass).
			Filter(fmt.Sprintf(`secret == "s3cret" && %s`, titleLike)).
			Execute()
		assertServerCode(t, err, status, 400)

		_, status, err = h.client.DataSync.GetEntities().
			EntityClass(dsEntityClass).
			Filter(fmt.Sprintf(`category == "tools" && %s`, titleLike)).
			Execute()
		assertServerCode(t, err, status, 400)
	})

	t.Run("Sort", func(t *testing.T) {
		res, status, err := h.client.DataSync.GetEntities().
			EntityClass(dsEntityClass).
			FilterFast(titleLike).
			Sort([]string{"score:desc"}).
			Limit(100).
			Execute()
		assert.Nil(t, err)
		assert.Equal(t, 200, status.StatusCode)
		assert.GreaterOrEqual(t, len(res.Data), 2)
		first, _ := res.Data[0].Payload["score"].(float64)
		last, _ := res.Data[len(res.Data)-1].Payload["score"].(float64)
		assert.GreaterOrEqual(t, first, last)

		_, badStatus, badErr := h.client.DataSync.GetEntities().
			EntityClass(dsEntityClass).
			Filter(titleLike).
			Sort([]string{"category"}).
			Execute()
		assertServerCode(t, badErr, badStatus, 400)
	})

	t.Run("RelationshipListFilters", func(t *testing.T) {
		byA, status, err := h.client.DataSync.GetRelationships().
			RelationshipClass(dsRelClass).
			EntityAID(c.entity[0].ID).
			Execute()
		assert.Nil(t, err)
		assert.Equal(t, 200, status.StatusCode)
		assert.Len(t, byA.Data, 1)
		assert.Equal(t, c.rel[0].ID, byA.Data[0].ID)

		byB, status, err := h.client.DataSync.GetRelationships().
			RelationshipClass(dsRelClass).
			EntityBID(c.status[0].ID).
			Execute()
		assert.Nil(t, err)
		assert.Equal(t, 200, status.StatusCode)
		assert.GreaterOrEqual(t, len(byB.Data), 1)

		owners := waitRelationshipsFilterFast(t, h.client, fmt.Sprintf(`role == "%s-owner"`, c.prefix), 2)
		assert.Len(t, owners.Data, 2)

		ownerFilter := eventually(t, dsFilterTimeout, 250*time.Millisecond,
			"GetRelationships Filter by unique owner role",
			func() (*pubnub.PNRelationshipsResponse, bool) {
				r, _, e := h.client.DataSync.GetRelationships().
					RelationshipClass(dsRelClass).
					RelationshipClassVersion(dsClassVersion).
					Filter(fmt.Sprintf(`role == "%s-owner"`, c.prefix)).
					Limit(100).
					Execute()
				return r, e == nil && r != nil && len(r.Data) == 2
			})
		assert.Len(t, ownerFilter.Data, 2)

		links := waitRelationshipsFilterFast(t, h.client, fmt.Sprintf(`kind == "%s-link"`, c.prefix), 3)
		assert.Len(t, links.Data, 3)

		_, noneStatus, noneErr := h.client.DataSync.GetRelationships().
			RelationshipClass(dsRelClass).
			Filter(`comment == "not searchable"`).
			Execute()
		assertServerCode(t, noneErr, noneStatus, 400)

		_, simpleStatus, simpleErr := h.client.DataSync.GetRelationships().
			RelationshipClass(dsRelClass).
			Filter(fmt.Sprintf(`kind == "%s-link"`, c.prefix)).
			Execute()
		assertServerCode(t, simpleErr, simpleStatus, 400)

		heavy, status, err := h.client.DataSync.GetRelationships().
			RelationshipClass(dsRelClass).
			FilterFast(fmt.Sprintf(`weight >= 30 && (role == "%s-owner" || role == "%s-editor" || role == "%s-viewer")`, c.prefix, c.prefix, c.prefix)).
			Sort([]string{"weight:desc"}).
			Limit(100).
			Execute()
		assert.Nil(t, err)
		assert.Equal(t, 200, status.StatusCode)
		assert.GreaterOrEqual(t, len(heavy.Data), 3)

		cursor := ""
		seen := 0
		for i := 0; i < 5; i++ {
			b := h.client.DataSync.GetRelationships().
				RelationshipClass(dsRelClass).
				EntityAID(c.entity[0].ID).
				Limit(1)
			if cursor != "" {
				b = b.Cursor(cursor)
			}
			res, _, err := b.Execute()
			assert.Nil(t, err)
			seen += len(res.Data)
			if res.Meta == nil || !res.Meta.HasNext {
				break
			}
			cursor = res.Meta.NextCursor
		}
		assert.Equal(t, 1, seen)
	})
}
