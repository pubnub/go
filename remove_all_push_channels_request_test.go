package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveAllPushNotificationsValidate(t *testing.T) {
	assert := assert.New(t)

	opts := &removeAllPushChannelsForDeviceOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	err := opts.validate()
	assert.Nil(err)

	opts1 := &removeAllPushChannelsForDeviceOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeNone,
		pubnub:          pubnub,
	}

	err1 := opts1.validate()
	assert.Contains(err1.Error(), "Missing Push Type")

	opts3 := &removeAllPushChannelsForDeviceOpts{
		PushType: PNPushTypeAPNS,
		pubnub:   pubnub,
	}

	err3 := opts3.validate()

	assert.Contains(err3.Error(), "Missing Device ID")

}

func TestRemoveAllPushNotificationsBuildPath(t *testing.T) {
	assert := assert.New(t)

	opts := &removeAllPushChannelsForDeviceOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	str, err := opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId/remove", str)
	assert.Nil(err)

}

func TestRemoveAllPushNotificationsBuildQuery(t *testing.T) {
	assert := assert.New(t)

	opts := &removeAllPushChannelsForDeviceOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	u, err := opts.buildQuery()
	assert.Equal("apns", u.Get("type"))
	assert.Nil(err)
}

func TestRemoveAllPushNotificationsBuildBody(t *testing.T) {
	assert := assert.New(t)

	opts := &removeAllPushChannelsForDeviceOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	_, err := opts.buildBody()
	assert.Nil(err)

}

func TestListPushProvisionsRequestValidateSub(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &listPushProvisionsRequestOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pn,
	}

	assert.Equal("pubnub/validation: pubnub: \x0e: Missing Subscribe Key", opts.validate().Error())
}
