package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go/v5"
	"github.com/stretchr/testify/assert"
)

func TestRemovePushNotificationsFromChannels(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemovePushNotificationsFromChannels().
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	assert.Nil(err)
}

func TestRemovePushNotificationsFromChannelsContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemovePushNotificationsFromChannelsWithContext(backgroundContext).
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	assert.Nil(err)
}
func TestRemovePushNotificationsFromChannelsTopicAndEnv(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemovePushNotificationsFromChannels().
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(pubnub.PNPushTypeGCM).
		Topic("a").
		Environment(pubnub.PNPushEnvironmentProduction).
		Execute()
	assert.Nil(err)
}

func TestRemovePushNotificationsFromChannelsTopicAndEnvContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemovePushNotificationsFromChannelsWithContext(backgroundContext).
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(pubnub.PNPushTypeGCM).
		Topic("a").
		Environment(pubnub.PNPushEnvironmentProduction).
		Execute()
	assert.Nil(err)
}
