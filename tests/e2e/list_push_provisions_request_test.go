package e2e

import (
	//"log"
	//"os"
	"fmt"
	pubnub "github.com/sprucehealth/pubnub-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListPushProvisionsNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	r := GenRandom()
	ch1 := fmt.Sprintf("testChannel_sub_%d", r.Intn(99999))
	cg1 := fmt.Sprintf("testCG_sub_%d", r.Intn(99999))

	_, _, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{ch1}).
		DeviceIDForPush(cg1).
		PushType(pubnub.PNPushTypeGCM).
		Execute()

	assert.Nil(err)

	resp, _, err := pn.ListPushProvisions().
		DeviceIDForPush(cg1).
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	assert.Equal(ch1, resp.Channels[0])
	assert.Nil(err)
}

func TestListPushProvisionsNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{"ch2"}).
		DeviceIDForPush("cg2").
		PushType(pubnub.PNPushTypeGCM).
		Execute()

	assert.Nil(err)

	resp, _, err := pn.ListPushProvisionsWithContext(backgroundContext).
		DeviceIDForPush("cg2").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	assert.Equal("ch2", resp.Channels[0])
	assert.Nil(err)
}
