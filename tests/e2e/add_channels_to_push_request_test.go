package e2e

import (
	//"fmt"
	//"log"
	//"os"
	"testing"

	pubnub "github.com/sprucehealth/pubnub-go"
	"github.com/stretchr/testify/assert"
)

func TestAddChannelToPushNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	_, _, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	//fmt.Println(err.Error())
	assert.Nil(err)
}

func TestAddChannelToPushNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{"ch1"}).
		DeviceIDForPush("cg1").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	//fmt.Println(err.Error())
	assert.Nil(err)
}

func TestAddChannelToPushNotStubbedContextWithQueryParam(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{"ch1"}).
		DeviceIDForPush("cg1").
		PushType(pubnub.PNPushTypeGCM).
		QueryParam(queryParam).
		Execute()
	//fmt.Println(err.Error())
	assert.Nil(err)
}
