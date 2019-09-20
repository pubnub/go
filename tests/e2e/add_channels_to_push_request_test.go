package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestAddChannelToPushNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	assert.Nil(err)
}

func TestAddChannelToPushNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{"ch1"}).
		DeviceIDForPush("cg1").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	assert.Nil(err)
}

func TestAddChannelToPushNotStubbedContextWithQueryParam(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{"ch1"}).
		DeviceIDForPush("cg1").
		PushType(pubnub.PNPushTypeGCM).
		QueryParam(queryParam).
		Execute()
	assert.Nil(err)
}
