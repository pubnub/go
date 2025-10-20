package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go/v7"

	"github.com/stretchr/testify/assert"
)

func TestListPushProvisionsNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	ch1 := randomized("testChannel_sub_")
	cg1 := randomized("testCG_sub_")

	_, _, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{ch1}).
		DeviceIDForPush(cg1).
		PushType(pubnub.PNPushTypeFCM).
		Execute()

	assert.Nil(err)

	resp, _, err := pn.ListPushProvisions().
		DeviceIDForPush(cg1).
		PushType(pubnub.PNPushTypeFCM).
		Execute()
	assert.Contains(resp.Channels, ch1)
	assert.Nil(err)
}

func TestListPushProvisionsNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	ch1 := randomized("testChannel_sub_")
	cg1 := randomized("testCG_sub_")

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{ch1}).
		DeviceIDForPush(cg1).
		PushType(pubnub.PNPushTypeFCM).
		Execute()

	assert.Nil(err)

	resp, _, err := pn.ListPushProvisionsWithContext(backgroundContext).
		DeviceIDForPush(cg1).
		PushType(pubnub.PNPushTypeFCM).
		Execute()
	assert.Contains(resp.Channels, ch1)
	assert.Nil(err)
}

func TestListPushProvisionsTopicAndEnvNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	ch1 := randomized("testChannel_sub_")
	cg1 := randomized("testCG_sub_")

	_, _, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{ch1}).
		DeviceIDForPush(cg1).
		PushType(pubnub.PNPushTypeFCM).
		Execute()

	assert.Nil(err)

	resp, _, err := pn.ListPushProvisions().
		DeviceIDForPush(cg1).
		PushType(pubnub.PNPushTypeFCM).
		Topic("a").
		Environment(pubnub.PNPushEnvironmentProduction).
		Execute()
	assert.Contains(resp.Channels, ch1)
	assert.Nil(err)
}

func TestListPushProvisionsTopicAndEnvNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	ch1 := randomized("testChannel_sub_")
	cg1 := randomized("testCG_sub_")

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{ch1}).
		DeviceIDForPush(cg1).
		PushType(pubnub.PNPushTypeFCM).
		Topic("a").
		Environment(pubnub.PNPushEnvironmentProduction).
		Execute()

	assert.Nil(err)

	resp, _, err := pn.ListPushProvisionsWithContext(backgroundContext).
		DeviceIDForPush(cg1).
		PushType(pubnub.PNPushTypeFCM).
		Topic("a").
		Environment(pubnub.PNPushEnvironmentProduction).
		Execute()
	assert.Contains(resp.Channels, ch1)
	assert.Nil(err)
}
