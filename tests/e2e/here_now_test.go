package e2e

import (
	//"log"
	"fmt"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestHereNowNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.HereNow().
		Channels([]string{"ch"}).
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
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
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
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
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
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
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
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
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
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
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
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
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
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
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
