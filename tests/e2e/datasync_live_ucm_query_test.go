package e2e

import (
	"fmt"
	"strings"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v10"
	"github.com/stretchr/testify/assert"
)

func corpusDisplayNameFilter(prefix string) string {
	parts := make([]string, len(queryTitles))
	for i, title := range queryTitles {
		parts[i] = fmt.Sprintf(`displayName == "%s-%s"`, prefix, title)
	}
	return "(" + strings.Join(parts, " || ") + ")"
}

func corpusTopicFilter(prefix string) string {
	parts := make([]string, len(queryTitles))
	for i, title := range queryTitles {
		parts[i] = fmt.Sprintf(`topic == "%s-%s"`, prefix, title)
	}
	return "(" + strings.Join(parts, " || ") + ")"
}

type dsUCMQueryCorpus struct {
	prefix     string
	users      []pubnub.PNEntity
	channels   []pubnub.PNEntity
	gUsers     []pubnub.PNEntity
	gChannels  []pubnub.PNEntity
	membership []pubnub.PNDataSyncMembership
}

func seedUCMQueryCorpus(t *testing.T, h *dsLive) dsUCMQueryCorpus {
	t.Helper()
	prefix := randomized("goe2e-q")
	categories := []string{"tools", "parts", "tools", "parts", "tools", "parts", "tools", "parts"}
	cities := []string{"Portland", "Seattle", "Portland", "Austin", "Seattle", "Portland", "Austin", "Seattle"}
	actives := []bool{true, true, false, true, false, true, false, true}

	c := dsUCMQueryCorpus{prefix: prefix}
	usrPrefix := randomized("goe2e-usr-q")
	chnPrefix := randomized("goe2e-chn-q")
	gusrPrefix := randomized("goe2e-gusr-q")
	gchnPrefix := randomized("goe2e-gchn-q")
	memPrefix := randomized("goe2e-mem-q")
	for i, title := range queryTitles {
		up := kitchenSinkGoE2EUserPayload()
		up["displayName"] = prefix + "-" + title
		up["category"] = categories[i]
		up["score"] = float64((i + 1) * 10)
		up["active"] = actives[i]
		up["address"] = map[string]interface{}{"city": cities[i]}
		up["email"] = fmt.Sprintf("%s-%s@example.com", prefix, title)
		up["updatedTs"] = fmt.Sprintf("2026-09-%02dT12:00:00Z", i+1)
		user := h.createGoE2EUser(h.admin, fmt.Sprintf("%s-%d", usrPrefix, i), up, "active")
		c.users = append(c.users, *user)

		cp := kitchenSinkGoE2EChannelPayload()
		cp["topic"] = prefix + "-" + title
		cp["category"] = categories[i]
		cp["capacity"] = float64((i + 1) * 10)
		cp["live"] = actives[i]
		cp["address"] = map[string]interface{}{"city": cities[i]}
		cp["email"] = fmt.Sprintf("%s-%s-ch@example.com", prefix, title)
		cp["lastActive"] = fmt.Sprintf("2026-09-%02dT12:00:00Z", i+1)
		ch := h.createGoE2EChannel(h.admin, fmt.Sprintf("%s-%d", chnPrefix, i), cp, "active")
		c.channels = append(c.channels, *ch)
	}

	for i := 0; i < 3; i++ {
		gp := globalUserPayload()
		gp["name"] = fmt.Sprintf("%s-GUser-%d", prefix, i)
		gu := h.createGlobalUser(h.admin, fmt.Sprintf("%s-%d", gusrPrefix, i), gp, "active")
		c.gUsers = append(c.gUsers, *gu)

		gcp := globalChannelPayload()
		gcp["name"] = fmt.Sprintf("%s-GChan-%d", prefix, i)
		gc := h.createGlobalChannel(h.admin, fmt.Sprintf("%s-%d", gchnPrefix, i), gcp, "active")
		c.gChannels = append(c.gChannels, *gc)
	}

	pairs := [][2]int{{0, 0}, {1, 0}, {2, 1}}
	for i, pair := range pairs {
		mem := h.createMembership(h.admin, fmt.Sprintf("%s-%d", memPrefix, i),
			c.gUsers[pair[0]].ID, c.gChannels[pair[1]].ID, nil, "active")
		c.membership = append(c.membership, *mem)
	}

	waitUsersFilterFast(t, h.client, dsGoE2EUserClass, pubnub.PNEntityClassLevelSubKey, corpusDisplayNameFilter(prefix), 8)
	waitUsersFilter(t, h.client, dsGoE2EUserClass, pubnub.PNEntityClassLevelSubKey, fmt.Sprintf(`displayName == "%s-Hotel"`, prefix), 1)
	waitChannelsFilterFast(t, h.client, dsGoE2EChannelClass, pubnub.PNEntityClassLevelSubKey, corpusTopicFilter(prefix), 8)
	waitChannelsFilter(t, h.client, dsGoE2EChannelClass, pubnub.PNEntityClassLevelSubKey, fmt.Sprintf(`topic == "%s-Hotel"`, prefix), 1)
	_ = eventually(t, dsFilterTimeout, 250*time.Millisecond,
		"membership corpus visible by userId",
		func() (*pubnub.PNDataSyncMembershipsResponse, bool) {
			r, _, err := h.client.DataSync.GetMemberships().
				UserID(c.gUsers[0].ID).
				Execute()
			return r, err == nil && r != nil && len(r.Data) == 1
		})
	return c
}

func TestDataSyncLiveUCMQuery(t *testing.T) {
	h := newDSLive(t)
	h.grantHappy()
	c := seedUCMQueryCorpus(t, h)
	nameLike := corpusDisplayNameFilter(c.prefix)
	topicLike := corpusTopicFilter(c.prefix)

	t.Run("UserClassLevelAndVersion", func(t *testing.T) {
		res, status, err := h.client.DataSync.GetUsers().
			EntityClass(dsGoE2EUserClass).
			EntityClassVersion(dsClassVersion).
			EntityClassLevel(pubnub.PNEntityClassLevelSubKey).
			FilterFast(nameLike).
			Limit(100).
			Execute()
		assert.Nil(t, err)
		assert.Equal(t, 200, status.StatusCode)
		assert.Len(t, res.Data, 8)
		for _, u := range res.Data {
			assert.Equal(t, dsGoE2EUserClass, u.EntityClass)
		}
	})

	t.Run("UserClassIsolation", func(t *testing.T) {
		custom, status, err := h.client.DataSync.GetUsers().
			EntityClass(dsGoE2EUserClass).
			EntityClassLevel(pubnub.PNEntityClassLevelSubKey).
			FilterFast(nameLike).
			Limit(100).
			Execute()
		assert.Nil(t, err)
		assert.Equal(t, 200, status.StatusCode)
		assert.Len(t, custom.Data, 8)

		global := waitUsersFilterFast(t, h.client, dsUserClass, pubnub.PNEntityClassLevelGlobal,
			fmt.Sprintf(`name == "%s-GUser-0"`, c.prefix), 1)
		assert.Equal(t, dsUserClass, global.Data[0].EntityClass)
		assert.Equal(t, c.gUsers[0].ID, global.Data[0].ID)
	})

	t.Run("UserPagination", func(t *testing.T) {
		seen := map[string]bool{}
		cursor := ""
		pages := 0
		for {
			b := h.client.DataSync.GetUsers().
				EntityClass(dsGoE2EUserClass).
				EntityClassLevel(pubnub.PNEntityClassLevelSubKey).
				FilterFast(nameLike).
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

	t.Run("UserFilterFastSimpleAndFull", func(t *testing.T) {
		tools := waitUsersFilterFast(t, h.client, dsGoE2EUserClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`category == "tools" && %s`, nameLike), 4)
		assert.Len(t, tools.Data, 4)

		portland := waitUsersFilterFast(t, h.client, dsGoE2EUserClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`city == "Portland" && %s`, nameLike), 3)
		assert.Len(t, portland.Data, 3)

		active := waitUsersFilterFast(t, h.client, dsGoE2EUserClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`active == true && %s`, nameLike), 5)
		assert.Len(t, active.Data, 5)

		high := waitUsersFilterFast(t, h.client, dsGoE2EUserClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`score > 35 && %s`, nameLike), 5)
		assert.Len(t, high.Data, 5)
	})

	t.Run("UserFilterFullOnly", func(t *testing.T) {
		alpha := waitUsersFilter(t, h.client, dsGoE2EUserClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`displayName == "%s-Alpha"`, c.prefix), 1)
		assert.Equal(t, c.users[0].ID, alpha.Data[0].ID)

		high := waitUsersFilter(t, h.client, dsGoE2EUserClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`score >= 50 && %s`, nameLike), 4)
		assert.Len(t, high.Data, 4)

		email := waitUsersFilter(t, h.client, dsGoE2EUserClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`email == "%s-Alpha@example.com"`, c.prefix), 1)
		assert.Equal(t, c.users[0].ID, email.Data[0].ID)
	})

	t.Run("UserFilterRejectedFields", func(t *testing.T) {
		_, status, err := h.client.DataSync.GetUsers().
			EntityClass(dsGoE2EUserClass).
			Filter(fmt.Sprintf(`notes == "not searchable" && %s`, nameLike)).
			Execute()
		assertServerCode(t, err, status, 400)

		_, status, err = h.client.DataSync.GetUsers().
			EntityClass(dsGoE2EUserClass).
			Filter(fmt.Sprintf(`secret == "s3cret" && %s`, nameLike)).
			Execute()
		assertServerCode(t, err, status, 400)

		_, status, err = h.client.DataSync.GetUsers().
			EntityClass(dsGoE2EUserClass).
			Filter(fmt.Sprintf(`category == "tools" && %s`, nameLike)).
			Execute()
		assertServerCode(t, err, status, 400)
	})

	t.Run("UserSort", func(t *testing.T) {
		res, status, err := h.client.DataSync.GetUsers().
			EntityClass(dsGoE2EUserClass).
			FilterFast(nameLike).
			Sort([]string{"score:desc"}).
			Limit(100).
			Execute()
		assert.Nil(t, err)
		assert.Equal(t, 200, status.StatusCode)
		assert.GreaterOrEqual(t, len(res.Data), 2)
		first, _ := res.Data[0].Payload["score"].(float64)
		last, _ := res.Data[len(res.Data)-1].Payload["score"].(float64)
		assert.GreaterOrEqual(t, first, last)

		_, badStatus, badErr := h.client.DataSync.GetUsers().
			EntityClass(dsGoE2EUserClass).
			Filter(nameLike).
			Sort([]string{"category"}).
			Execute()
		assertServerCode(t, badErr, badStatus, 400)
	})

	t.Run("ChannelClassLevelAndVersion", func(t *testing.T) {
		res, status, err := h.client.DataSync.GetChannels().
			EntityClass(dsGoE2EChannelClass).
			EntityClassVersion(dsClassVersion).
			EntityClassLevel(pubnub.PNEntityClassLevelSubKey).
			FilterFast(topicLike).
			Limit(100).
			Execute()
		assert.Nil(t, err)
		assert.Equal(t, 200, status.StatusCode)
		assert.Len(t, res.Data, 8)
	})

	t.Run("ChannelPagination", func(t *testing.T) {
		seen := map[string]bool{}
		cursor := ""
		pages := 0
		for {
			b := h.client.DataSync.GetChannels().
				EntityClass(dsGoE2EChannelClass).
				EntityClassLevel(pubnub.PNEntityClassLevelSubKey).
				FilterFast(topicLike).
				Sort([]string{"capacity"}).
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
			cursor = res.Meta.NextCursor
			if pages > 8 {
				t.Fatal("pagination did not terminate")
			}
		}
		assert.GreaterOrEqual(t, pages, 4)
		assert.Len(t, seen, 8)
	})

	t.Run("ChannelFilterFastSimpleAndFull", func(t *testing.T) {
		tools := waitChannelsFilterFast(t, h.client, dsGoE2EChannelClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`category == "tools" && %s`, topicLike), 4)
		assert.Len(t, tools.Data, 4)

		portland := waitChannelsFilterFast(t, h.client, dsGoE2EChannelClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`city == "Portland" && %s`, topicLike), 3)
		assert.Len(t, portland.Data, 3)

		live := waitChannelsFilterFast(t, h.client, dsGoE2EChannelClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`live == true && %s`, topicLike), 5)
		assert.Len(t, live.Data, 5)

		high := waitChannelsFilterFast(t, h.client, dsGoE2EChannelClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`capacity > 35 && %s`, topicLike), 5)
		assert.Len(t, high.Data, 5)
	})

	t.Run("ChannelFilterFullOnly", func(t *testing.T) {
		alpha := waitChannelsFilter(t, h.client, dsGoE2EChannelClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`topic == "%s-Alpha"`, c.prefix), 1)
		assert.Equal(t, c.channels[0].ID, alpha.Data[0].ID)

		high := waitChannelsFilter(t, h.client, dsGoE2EChannelClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`capacity >= 50 && %s`, topicLike), 4)
		assert.Len(t, high.Data, 4)

		email := waitChannelsFilter(t, h.client, dsGoE2EChannelClass, pubnub.PNEntityClassLevelSubKey,
			fmt.Sprintf(`email == "%s-Alpha-ch@example.com"`, c.prefix), 1)
		assert.Equal(t, c.channels[0].ID, email.Data[0].ID)
	})

	t.Run("ChannelFilterRejectedFields", func(t *testing.T) {
		_, status, err := h.client.DataSync.GetChannels().
			EntityClass(dsGoE2EChannelClass).
			Filter(fmt.Sprintf(`notes == "not searchable" && %s`, topicLike)).
			Execute()
		assertServerCode(t, err, status, 400)

		_, status, err = h.client.DataSync.GetChannels().
			EntityClass(dsGoE2EChannelClass).
			Filter(fmt.Sprintf(`secret == "s3cret" && %s`, topicLike)).
			Execute()
		assertServerCode(t, err, status, 400)

		_, status, err = h.client.DataSync.GetChannels().
			EntityClass(dsGoE2EChannelClass).
			Filter(fmt.Sprintf(`category == "tools" && %s`, topicLike)).
			Execute()
		assertServerCode(t, err, status, 400)
	})

	t.Run("MembershipListFilters", func(t *testing.T) {
		byUser, status, err := h.client.DataSync.GetMemberships().
			UserID(c.gUsers[0].ID).
			Execute()
		assert.Nil(t, err)
		assert.Equal(t, 200, status.StatusCode)
		assert.Len(t, byUser.Data, 1)
		assert.Equal(t, c.membership[0].ID, byUser.Data[0].ID)

		byCh, status, err := h.client.DataSync.GetMemberships().
			ChannelID(c.gChannels[0].ID).
			Execute()
		assert.Nil(t, err)
		assert.Equal(t, 200, status.StatusCode)
		assert.Len(t, byCh.Data, 2)

		cursor := ""
		seen := 0
		for i := 0; i < 5; i++ {
			b := h.client.DataSync.GetMemberships().
				ChannelID(c.gChannels[0].ID).
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
		assert.Equal(t, 2, seen)
	})
}
