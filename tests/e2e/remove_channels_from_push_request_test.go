package e2e

import (
	//"log"
	//"os"
	"testing"

	pubnub "github.com/sprucehealth/pubnub-go"
	"github.com/stretchr/testify/assert"
)

func TestRemovePushNotificationsFromChannels(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

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
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	_, _, err := pn.RemovePushNotificationsFromChannelsWithContext(backgroundContext).
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	assert.Nil(err)
}
