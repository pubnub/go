package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go/v5"
	"github.com/stretchr/testify/assert"
)

func TestRemoveAllPushNotifications(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemoveAllPushNotifications().
		DeviceIDForPush("cg").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	assert.Nil(err)
}

func TestRemoveAllPushNotificationsContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemoveAllPushNotificationsWithContext(backgroundContext).
		DeviceIDForPush("cg").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	assert.Nil(err)
}
