package e2e

import (
	//"log"
	//"os"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestListPushProvisionsNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	_, _, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{"ch1"}).
		DeviceIDForPush("cg1").
		PushType(pubnub.PNPushTypeGCM).
		Execute()

	assert.Nil(err)

	resp, _, err := pn.ListPushProvisions().
		DeviceIDForPush("cg1").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	assert.Equal("ch1", resp.Channels[0])
	assert.Nil(err)
}
