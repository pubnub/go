package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddChannelsToPushOptsValidate(t *testing.T) {
	assert := assert.New(t)

	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
	})

	err := opts.validate()
	assert.Nil(err)

	opts1 := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeNone,
	})

	err1 := opts1.validate()
	assert.Contains(err1.Error(), "Missing Push Type")

	opts2 := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
	})

	err2 := opts2.validate()
	assert.Contains(err2.Error(), "Missing Channel")

	opts3 := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels: []string{"ch1", "ch2", "ch3"},
		PushType: PNPushTypeAPNS,
	})

	err3 := opts3.validate()
	assert.Contains(err3.Error(), "Missing Device ID")

}

func TestAddChannelsToPushOptsBuildPath(t *testing.T) {
	assert := assert.New(t)

	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
	})

	str, err := opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId", str)
	assert.Nil(err)

}

func TestAddChannelsToPushOptsBuildQuery(t *testing.T) {
	assert := assert.New(t)

	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
	})

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("add"))
	assert.Equal("apns", u.Get("type"))
	assert.Nil(err)
}

func TestAddChannelsToPushOptsBuildQueryParams(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		QueryParam:      queryParam,
	})

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("add"))
	assert.Equal("apns", u.Get("type"))
	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))
	assert.Nil(err)
}

func TestAddChannelsToPushOptsBuildQueryParamsTopicAndEnv(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		QueryParam:      queryParam,
		Topic:           "a",
		Environment:     PNPushEnvironmentProduction,
	})

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("add"))
	assert.Equal("apns", u.Get("type"))
	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))
	assert.Equal("production", u.Get("environment"))
	assert.Equal("a", u.Get("topic"))
	assert.Nil(err)
}

func TestAddChannelsToPushOptsBuildBody(t *testing.T) {
	assert := assert.New(t)

	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
	})

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

	o := newAddPushNotificationsOnChannelsBuilderWithContext(pubnub, backgroundContext)
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
	opts := newAddChannelsToPushOpts(pn, pn.ctx, addChannelsToPushOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
	})

	assert.Equal("pubnub/validation: pubnub: Remove Channel Group: Missing Subscribe Key", opts.validate().Error())
}
