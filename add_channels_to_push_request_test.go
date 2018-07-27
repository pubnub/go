package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestAddChannelsToPushOptsValidate(t *testing.T) {
	assert := assert.New(t)

	opts := &addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	err := opts.validate()
	assert.Nil(err)

	opts1 := &addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeNone,
		pubnub:          pubnub,
	}

	err1 := opts1.validate()
	assert.Contains(err1.Error(), "Missing Push Type")

	opts2 := &addChannelsToPushOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	err2 := opts2.validate()
	assert.Contains(err2.Error(), "Missing Channel")

	opts3 := &addChannelsToPushOpts{
		Channels: []string{"ch1", "ch2", "ch3"},
		PushType: PNPushTypeAPNS,
		pubnub:   pubnub,
	}

	err3 := opts3.validate()
	assert.Contains(err3.Error(), "Missing Device ID")

}

func TestAddChannelsToPushOptsBuildPath(t *testing.T) {
	assert := assert.New(t)

	opts := &addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	str, err := opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId", str)
	assert.Nil(err)

}

func TestAddChannelsToPushOptsBuildQuery(t *testing.T) {
	assert := assert.New(t)

	opts := &addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("add"))
	assert.Equal("apns", u.Get("type"))
	assert.Nil(err)
}

func TestAddChannelsToPushOptsBuildBody(t *testing.T) {
	assert := assert.New(t)

	opts := &addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	_, err := opts.buildBody()

	assert.Nil(err)

}

func TestNewAddPushNotificationsOnChannelsBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newAddPushNotificationsOnChannelsBuilder(pubnub)
	o.Channels([]string{"ch1", "ch2", "ch3"})
	o.DeviceIDForPush("deviceID")
	o.PushType(PNPushTypeAPNS)

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceID", path)
}

func TestNewAddPushNotificationsOnChannelsBuilderWithContext(t *testing.T) {
	assert := assert.New(t)

	o := newAddPushNotificationsOnChannelsBuilderWithContext(pubnub, context.Background())
	o.Channels([]string{"ch1", "ch2", "ch3"})
	o.DeviceIDForPush("deviceID")
	o.PushType(PNPushTypeAPNS)

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceID", path)
}

func TestAddChannelsToPushValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &addChannelsToPushOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pn,
	}

	assert.Equal("pubnub/validation: pubnub: \x0e: Missing Subscribe Key", opts.validate().Error())
}
