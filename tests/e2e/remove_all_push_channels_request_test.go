package e2e

import (
	"log"
	"os"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestRemoveAllPushNotifications(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	_, _, err := pn.RemoveAllPushNotifications().
		DeviceIDForPush("cg").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	assert.Nil(err)
}
