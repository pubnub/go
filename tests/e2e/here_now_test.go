package e2e

import (
	//"log"
	"fmt"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v7"
	"github.com/pubnub/go/v7/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestHereNowNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.HereNow().
		Channels([]string{randomized("ch")}).
		Execute()

	assert.Nil(err)
}

func TestHereNowMultipleChannelsWithState(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub_key/%s/channel/ch1,ch2", config.SubscribeKey),
		Query:              "state=1",
		ResponseBody:       "{\"status\":200,\"message\":\"OK\",\"payload\":{\"total_occupancy\":3,\"total_channels\":2,\"channels\":{\"ch1\":{\"occupancy\":1,\"uuids\":[{\"uuid\":\"user1\",\"state\":{\"age\":10}}]},\"ch2\":{\"occupancy\":2,\"uuids\":[{\"uuid\":\"user1\",\"state\":{\"age\":10}},{\"uuid\":\"user3\",\"state\":{\"age\":30}}]}}},\"service\":\"Presence\"}",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid", "limit"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.HereNow().
		Channels([]string{"ch1", "ch2"}).
		IncludeState(true).
		Execute()

	assert.Equal(2, res.TotalChannels)
	assert.Equal(3, res.TotalOccupancy)

	if res.Channels[0].ChannelName == "ch1" {
		assert.Equal("ch1", res.Channels[0].ChannelName)
		assert.Equal(1, res.Channels[0].Occupancy)
		assert.Equal("user1", res.Channels[0].Occupants[0].UUID)
		assert.Equal(map[string]interface{}{"age": float64(10)}, res.Channels[0].Occupants[0].State)

		assert.Equal("ch2", res.Channels[1].ChannelName)
		assert.Equal(2, res.Channels[1].Occupancy)
		assert.Equal("user1", res.Channels[1].Occupants[0].UUID)
		assert.Equal(map[string]interface{}{"age": float64(10)}, res.Channels[1].Occupants[0].State)
		assert.Equal("user3", res.Channels[1].Occupants[1].UUID)
		assert.Equal(map[string]interface{}{"age": float64(30)}, res.Channels[1].Occupants[1].State)
	} else if res.Channels[1].ChannelName == "ch2" {
		assert.Equal("ch1", res.Channels[1].ChannelName)
		assert.Equal(1, res.Channels[1].Occupancy)
		assert.Equal("user1", res.Channels[1].Occupants[0].UUID)
		assert.Equal(map[string]interface{}{"age": float64(10)}, res.Channels[1].Occupants[0].State)

		assert.Equal("ch2", res.Channels[0].ChannelName)
		assert.Equal(2, res.Channels[0].Occupancy)
		assert.Equal("user1", res.Channels[0].Occupants[0].UUID)
		assert.Equal(map[string]interface{}{"age": float64(10)}, res.Channels[0].Occupants[0].State)
		assert.Equal("user3", res.Channels[0].Occupants[1].UUID)
		assert.Equal(map[string]interface{}{"age": float64(30)}, res.Channels[0].Occupants[1].State)
	}

	assert.Nil(err)
}

func TestMultipleChannelWithoutStateSync(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub_key/%s/channel/game1,game2", config.SubscribeKey),
		Query:              "state=0",
		ResponseBody:       "{\"status\": 200, \"message\": \"OK\", \"payload\": {\"channels\": {\"game1\": {\"uuids\": [\"a3ffd012-a3b9-478c-8705-64089f24d71e\"], \"occupancy\": 1}}, \"total_channels\": 1, \"total_occupancy\": 1}, \"service\": \"Presence\"}",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid", "limit"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.HereNow().
		Channels([]string{"game1", "game2"}).
		IncludeState(false).
		Execute()

	assert.Equal(1, res.TotalChannels)
	assert.Equal(1, res.TotalOccupancy)

	assert.Equal("game1", res.Channels[0].ChannelName)
	assert.Equal(1, res.Channels[0].Occupancy)
	//log.Println(res.Channels[0])
	assert.Equal("a3ffd012-a3b9-478c-8705-64089f24d71e", res.Channels[0].Occupants[0].UUID)
	assert.Equal(map[string]interface{}{}, res.Channels[0].Occupants[0].State)

	assert.Nil(err)
}

func TestHereNowMultipleChannelsWithoutUUIDs(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub_key/%s/channel/game1,game2", config.SubscribeKey),
		Query:              "state=0&disable-uuids=1",
		ResponseBody:       "{\"status\": 200, \"message\": \"OK\", \"payload\": {\"channels\": {\"game1\": {\"occupancy\": 1}}, \"total_channels\": 1, \"total_occupancy\": 1}, \"service\": \"Presence\"}",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid", "limit"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.HereNow().
		Channels([]string{"game1", "game2"}).
		IncludeState(false).
		IncludeUUIDs(false).
		Execute()

	assert.Equal(1, res.TotalChannels)
	assert.Equal(1, res.TotalOccupancy)

	assert.Equal("game1", res.Channels[0].ChannelName)
	assert.Equal(1, res.Channels[0].Occupancy)
	assert.Equal(0, len(res.Channels[0].Occupants))

	assert.Nil(err)
}

func TestHereNowSingleChannelWithState(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub_key/%s/channel/game1", config.SubscribeKey),
		Query:              "state=1",
		ResponseBody:       "{\"status\":200,\"message\":\"OK\",\"service\":\"Presence\",\"uuids\":[{\"uuid\":\"a3ffd012-a3b9-478c-8705-64089f24d71e\",\"state\":{\"age\":10}}],\"occupancy\":1}",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid", "limit"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.HereNow().
		Channels([]string{"game1"}).
		IncludeState(true).
		Execute()

	assert.Equal(1, res.TotalChannels)
	assert.Equal(1, res.TotalOccupancy)
	assert.Equal(1, len(res.Channels))

	assert.Equal("game1", res.Channels[0].ChannelName)
	assert.Equal(1, res.Channels[0].Occupancy)
	assert.Equal("a3ffd012-a3b9-478c-8705-64089f24d71e",
		res.Channels[0].Occupants[0].UUID)
	assert.Equal(map[string]interface{}{"age": float64(10)}, res.Channels[0].Occupants[0].State)

	assert.Nil(err)
}

func TestHereNowSingleChannelWithStateContext(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub_key/%s/channel/game1", config.SubscribeKey),
		Query:              "state=1",
		ResponseBody:       "{\"status\":200,\"message\":\"OK\",\"service\":\"Presence\",\"uuids\":[{\"uuid\":\"a3ffd012-a3b9-478c-8705-64089f24d71e\",\"state\":{\"age\":10}}],\"occupancy\":1}",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid", "limit"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.HereNowWithContext(backgroundContext).
		Channels([]string{"game1"}).
		IncludeState(true).
		Execute()

	assert.Equal(1, res.TotalChannels)
	assert.Equal(1, res.TotalOccupancy)
	assert.Equal(1, len(res.Channels))

	assert.Equal("game1", res.Channels[0].ChannelName)
	assert.Equal(1, res.Channels[0].Occupancy)
	assert.Equal("a3ffd012-a3b9-478c-8705-64089f24d71e",
		res.Channels[0].Occupants[0].UUID)
	assert.Equal(map[string]interface{}{"age": float64(10)}, res.Channels[0].Occupants[0].State)

	assert.Nil(err)
}

func TestHereNowSingleChannelWithoutState(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub_key/%s/channel/game1", config.SubscribeKey),
		Query:              "state=0",
		ResponseBody:       "{\"status\": 200, \"message\": \"OK\", \"service\": \"Presence\", \"uuids\": [\"a3ffd012-a3b9-478c-8705-64089f24d71e\"], \"occupancy\": 1}",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid", "limit"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.HereNow().
		Channels([]string{"game1"}).
		IncludeState(false).
		Execute()

	assert.Equal(1, res.TotalChannels)
	assert.Equal(1, res.TotalOccupancy)
	assert.Equal(1, len(res.Channels))

	assert.Equal("game1", res.Channels[0].ChannelName)
	assert.Equal(1, res.Channels[0].Occupancy)
	assert.Equal("a3ffd012-a3b9-478c-8705-64089f24d71e",
		res.Channels[0].Occupants[0].UUID)
	assert.Equal(map[string]interface{}{}, res.Channels[0].Occupants[0].State)

	assert.Nil(err)
}

func TestHereNowSingleChannelAndGroup(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub_key/%s/channel/game1", config.SubscribeKey),
		Query:              "state=1&channel-group=cg",
		ResponseBody:       "{\"status\":200,\"message\":\"OK\",\"payload\":{\"channels\":{}, \"total_channels\":0, \"total_occupancy\":0},\"service\":\"Presence\"}",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid", "limit"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.HereNow().
		Channels([]string{"game1"}).
		ChannelGroups([]string{"cg"}).
		IncludeState(true).
		Execute()

	assert.Equal(0, res.TotalOccupancy)

	assert.Nil(err)
}

func TestHereNowPaginationBasic(t *testing.T) {
	assert := assert.New(t)

	channelName := randomized("pagination-test")

	// Create 3 PubNub instances with unique UUIDs
	config1 := configCopy()
	config1.SetUserId(pubnub.UserId("user-1-" + randomized("uuid")))
	pn1 := pubnub.NewPubNub(config1)

	config2 := configCopy()
	config2.SetUserId(pubnub.UserId("user-2-" + randomized("uuid")))
	pn2 := pubnub.NewPubNub(config2)

	config3 := configCopy()
	config3.SetUserId(pubnub.UserId("user-3-" + randomized("uuid")))
	pn3 := pubnub.NewPubNub(config3)

	// Subscribe all 3 instances to the same channel
	pn1.Subscribe().Channels([]string{channelName}).Execute()
	pn2.Subscribe().Channels([]string{channelName}).Execute()
	pn3.Subscribe().Channels([]string{channelName}).Execute()

	// Cleanup
	defer pn1.Unsubscribe().Channels([]string{channelName}).Execute()
	defer pn2.Unsubscribe().Channels([]string{channelName}).Execute()
	defer pn3.Unsubscribe().Channels([]string{channelName}).Execute()

	// Wait for presence to register
	time.Sleep(3 * time.Second)

	// Test 1: Get first 2 users (Limit=2, Offset=0)
	res1, _, err1 := pn1.HereNow().
		Channels([]string{channelName}).
		Limit(2).
		Offset(0).
		IncludeUUIDs(true).
		Execute()

	assert.Nil(err1)
	assert.Equal(3, res1.TotalOccupancy)
	assert.Equal(1, len(res1.Channels))
	assert.Equal(2, len(res1.Channels[0].Occupants))

	// Collect UUIDs from first page
	firstPageUUIDs := make(map[string]bool)
	for _, occupant := range res1.Channels[0].Occupants {
		firstPageUUIDs[occupant.UUID] = true
	}

	// Test 2: Get remaining user (Limit=2, Offset=2)
	res2, _, err2 := pn1.HereNow().
		Channels([]string{channelName}).
		Limit(2).
		Offset(2).
		IncludeUUIDs(true).
		Execute()

	assert.Nil(err2)
	assert.Equal(3, res2.TotalOccupancy)
	assert.Equal(1, len(res2.Channels))
	assert.Equal(1, len(res2.Channels[0].Occupants))

	// Verify no duplicate UUIDs between pages
	thirdUserUUID := res2.Channels[0].Occupants[0].UUID
	assert.False(firstPageUUIDs[thirdUserUUID], "UUID should not appear in both pages")

	// Verify we have all 3 unique users when combining pages
	assert.Equal(3, len(firstPageUUIDs)+1)
}

func TestHereNowPaginationFull(t *testing.T) {
	assert := assert.New(t)

	channelName := randomized("pagination-full")

	// Create 3 PubNub instances with unique UUIDs
	config1 := configCopy()
	config1.SetUserId(pubnub.UserId("user-1-" + randomized("uuid")))
	pn1 := pubnub.NewPubNub(config1)

	config2 := configCopy()
	config2.SetUserId(pubnub.UserId("user-2-" + randomized("uuid")))
	pn2 := pubnub.NewPubNub(config2)

	config3 := configCopy()
	config3.SetUserId(pubnub.UserId("user-3-" + randomized("uuid")))
	pn3 := pubnub.NewPubNub(config3)

	// Subscribe all 3 instances
	pn1.Subscribe().Channels([]string{channelName}).Execute()
	pn2.Subscribe().Channels([]string{channelName}).Execute()
	pn3.Subscribe().Channels([]string{channelName}).Execute()

	// Cleanup
	defer pn1.Unsubscribe().Channels([]string{channelName}).Execute()
	defer pn2.Unsubscribe().Channels([]string{channelName}).Execute()
	defer pn3.Unsubscribe().Channels([]string{channelName}).Execute()

	// Wait for presence to register
	time.Sleep(3 * time.Second)

	// Paginate through all users one by one (Limit=1)
	allUUIDs := make(map[string]bool)

	for offset := 0; offset < 3; offset++ {
		res, _, err := pn1.HereNow().
			Channels([]string{channelName}).
			Limit(1).
			Offset(offset).
			IncludeUUIDs(true).
			Execute()

		assert.Nil(err)
		assert.Equal(3, res.TotalOccupancy)
		assert.Equal(1, len(res.Channels))
		assert.Equal(1, len(res.Channels[0].Occupants))

		// Collect UUID
		uuid := res.Channels[0].Occupants[0].UUID
		assert.False(allUUIDs[uuid], "UUID should not be duplicated across pages")
		allUUIDs[uuid] = true
	}

	// Verify we collected all 3 unique users
	assert.Equal(3, len(allUUIDs))
}

func TestHereNowLimitLargerThanCount(t *testing.T) {
	assert := assert.New(t)

	channelName := randomized("limit-larger")

	// Create 3 PubNub instances
	config1 := configCopy()
	config1.SetUserId(pubnub.UserId("user-1-" + randomized("uuid")))
	pn1 := pubnub.NewPubNub(config1)

	config2 := configCopy()
	config2.SetUserId(pubnub.UserId("user-2-" + randomized("uuid")))
	pn2 := pubnub.NewPubNub(config2)

	config3 := configCopy()
	config3.SetUserId(pubnub.UserId("user-3-" + randomized("uuid")))
	pn3 := pubnub.NewPubNub(config3)

	// Subscribe all 3 instances
	pn1.Subscribe().Channels([]string{channelName}).Execute()
	pn2.Subscribe().Channels([]string{channelName}).Execute()
	pn3.Subscribe().Channels([]string{channelName}).Execute()

	// Cleanup
	defer pn1.Unsubscribe().Channels([]string{channelName}).Execute()
	defer pn2.Unsubscribe().Channels([]string{channelName}).Execute()
	defer pn3.Unsubscribe().Channels([]string{channelName}).Execute()

	// Wait for presence to register
	time.Sleep(3 * time.Second)

	// Set limit to 10 (larger than the 3 users present)
	res, _, err := pn1.HereNow().
		Channels([]string{channelName}).
		Limit(10).
		Offset(0).
		IncludeUUIDs(true).
		Execute()

	assert.Nil(err)
	assert.Equal(3, res.TotalOccupancy)
	assert.Equal(1, len(res.Channels))
	// Should return all 3 users even though limit is 10
	assert.Equal(3, len(res.Channels[0].Occupants))
}

func TestHereNowOffsetBeyondCount(t *testing.T) {
	assert := assert.New(t)

	channelName := randomized("offset-beyond")

	// Create 3 PubNub instances
	config1 := configCopy()
	config1.SetUserId(pubnub.UserId("user-1-" + randomized("uuid")))
	pn1 := pubnub.NewPubNub(config1)

	config2 := configCopy()
	config2.SetUserId(pubnub.UserId("user-2-" + randomized("uuid")))
	pn2 := pubnub.NewPubNub(config2)

	config3 := configCopy()
	config3.SetUserId(pubnub.UserId("user-3-" + randomized("uuid")))
	pn3 := pubnub.NewPubNub(config3)

	// Subscribe all 3 instances
	pn1.Subscribe().Channels([]string{channelName}).Execute()
	pn2.Subscribe().Channels([]string{channelName}).Execute()
	pn3.Subscribe().Channels([]string{channelName}).Execute()

	// Cleanup
	defer pn1.Unsubscribe().Channels([]string{channelName}).Execute()
	defer pn2.Unsubscribe().Channels([]string{channelName}).Execute()
	defer pn3.Unsubscribe().Channels([]string{channelName}).Execute()

	// Wait for presence to register
	time.Sleep(3 * time.Second)

	// Set offset to 5 (beyond the 3 users present)
	res, _, err := pn1.HereNow().
		Channels([]string{channelName}).
		Limit(10).
		Offset(5).
		IncludeUUIDs(true).
		Execute()

	assert.Nil(err)
	assert.Equal(3, res.TotalOccupancy)
	assert.Equal(1, len(res.Channels))
	// Should return 0 occupants since offset skips all users
	assert.Equal(0, len(res.Channels[0].Occupants))
}

func TestHereNowDefaultBehavior(t *testing.T) {
	assert := assert.New(t)

	channelName := randomized("default-behavior")

	// Create 3 PubNub instances
	config1 := configCopy()
	userId1 := "user-1-" + randomized("uuid")
	config1.SetUserId(pubnub.UserId(userId1))
	pn1 := pubnub.NewPubNub(config1)

	config2 := configCopy()
	userId2 := "user-2-" + randomized("uuid")
	config2.SetUserId(pubnub.UserId(userId2))
	pn2 := pubnub.NewPubNub(config2)

	config3 := configCopy()
	userId3 := "user-3-" + randomized("uuid")
	config3.SetUserId(pubnub.UserId(userId3))
	pn3 := pubnub.NewPubNub(config3)

	// Subscribe all 3 instances
	pn1.Subscribe().Channels([]string{channelName}).Execute()
	pn2.Subscribe().Channels([]string{channelName}).Execute()
	pn3.Subscribe().Channels([]string{channelName}).Execute()

	// Cleanup
	defer pn1.Unsubscribe().Channels([]string{channelName}).Execute()
	defer pn2.Unsubscribe().Channels([]string{channelName}).Execute()
	defer pn3.Unsubscribe().Channels([]string{channelName}).Execute()

	// Wait for presence to register
	time.Sleep(3 * time.Second)

	// Call HereNow without setting Limit or Offset (should use defaults: limit=1000, offset=0)
	res, _, err := pn1.HereNow().
		Channels([]string{channelName}).
		IncludeUUIDs(true).
		Execute()

	assert.Nil(err)
	assert.Equal(3, res.TotalOccupancy)
	assert.Equal(1, len(res.Channels))
	// Should return all 3 users with default settings
	assert.Equal(3, len(res.Channels[0].Occupants))
}

func TestHereNowOutOfRangeParameters(t *testing.T) {
	assert := assert.New(t)

	channelName := randomized("out-of-range")

	// Create 3 PubNub instances
	config1 := configCopy()
	userId1 := "user-1-" + randomized("uuid")
	config1.SetUserId(pubnub.UserId(userId1))
	pn1 := pubnub.NewPubNub(config1)

	config2 := configCopy()
	userId2 := "user-2-" + randomized("uuid")
	config2.SetUserId(pubnub.UserId(userId2))
	pn2 := pubnub.NewPubNub(config2)

	config3 := configCopy()
	userId3 := "user-3-" + randomized("uuid")
	config3.SetUserId(pubnub.UserId(userId3))
	pn3 := pubnub.NewPubNub(config3)

	// Subscribe all 3 instances
	pn1.Subscribe().Channels([]string{channelName}).Execute()
	pn2.Subscribe().Channels([]string{channelName}).Execute()
	pn3.Subscribe().Channels([]string{channelName}).Execute()

	// Cleanup
	defer pn1.Unsubscribe().Channels([]string{channelName}).Execute()
	defer pn2.Unsubscribe().Channels([]string{channelName}).Execute()
	defer pn3.Unsubscribe().Channels([]string{channelName}).Execute()

	// Wait for presence to register
	time.Sleep(3 * time.Second)

	// Test with out-of-range limit (above maximum)
	_, _, err1 := pn1.HereNow().
		Channels([]string{channelName}).
		Limit(2000). // Out of range (max is 1000)
		IncludeUUIDs(true).
		Execute()

	// Server should return an error for out-of-range limit
	assert.NotNil(err1, "Server should return error for limit > 1000")

	// Test with negative limit
	_, _, err2 := pn1.HereNow().
		Channels([]string{channelName}).
		Limit(-10). // Negative limit
		IncludeUUIDs(true).
		Execute()

	// Server should return an error for negative limit
	assert.NotNil(err2, "Server should return error for negative limit")

	// Test with negative offset
	_, _, err3 := pn1.HereNow().
		Channels([]string{channelName}).
		Offset(-50). // Negative offset
		IncludeUUIDs(true).
		Execute()

	// Server should return an error for negative offset
	assert.NotNil(err3, "Server should return error for negative offset")
}
