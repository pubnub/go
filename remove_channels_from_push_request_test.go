package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveChannelsFromPushRequestValidate(t *testing.T) {
	assert := assert.New(t)

	opts := newRemoveChannelsFromPushOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS

	err := opts.validate()
	assert.Nil(err)

	opts1 := newRemoveChannelsFromPushOpts(pubnub, pubnub.ctx)
	opts1.Channels = []string{"ch1", "ch2", "ch3"}
	opts1.DeviceIDForPush = "deviceId"
	opts1.PushType = PNPushTypeNone

	err1 := opts1.validate()
	assert.Contains(err1.Error(), "Missing Push Type")

	opts2 := newRemoveChannelsFromPushOpts(pubnub, pubnub.ctx)
	opts2.DeviceIDForPush = "deviceId"
	opts2.PushType = PNPushTypeAPNS

	err2 := opts2.validate()
	assert.Contains(err2.Error(), "Missing Channel")

	opts3 := newRemoveChannelsFromPushOpts(pubnub, pubnub.ctx)
	opts3.Channels = []string{"ch1", "ch2", "ch3"}
	opts3.PushType = PNPushTypeAPNS

	err3 := opts3.validate()
	assert.Contains(err3.Error(), "Missing Device ID")
}

func TestRemoveChannelsFromPushRequestBuildPath(t *testing.T) {
	assert := assert.New(t)

	opts := newRemoveChannelsFromPushOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS

	str, err := opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId", str)
	assert.Nil(err)

}

func TestRemoveChannelsFromPushRequestBuildQueryParam(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := newRemoveChannelsFromPushOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS
	opts.QueryParam = queryParam

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("remove"))
	assert.Equal("apns", u.Get("type"))
	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))

	assert.Nil(err)
}

func TestRemoveChannelsFromPushRequestBuildQueryParamTopicAndEnv(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := newRemoveChannelsFromPushOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS
	opts.QueryParam = queryParam
	opts.Topic = "a"
	opts.Environment = PNPushEnvironmentProduction

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("remove"))
	assert.Equal("apns", u.Get("type"))
	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))
	assert.Equal("production", u.Get("environment"))
	assert.Equal("a", u.Get("topic"))

	assert.Nil(err)
}

func TestRemoveChannelsFromPushRequestBuildQuery(t *testing.T) {
	assert := assert.New(t)

	opts := newRemoveChannelsFromPushOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("remove"))
	assert.Equal("apns", u.Get("type"))

	assert.Nil(err)
}

func TestRemoveChannelsFromPushRequestBuildBody(t *testing.T) {
	assert := assert.New(t)

	opts := newRemoveChannelsFromPushOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS

	_, err := opts.buildBody()
	assert.Nil(err)

}

func TestNewRemoveChannelsFromPushBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newRemoveChannelsFromPushBuilder(pubnub)
	o.Channels([]string{"ch1", "ch2", "ch3"})
	o.DeviceIDForPush("deviceId")
	o.PushType(PNPushTypeAPNS)
	u, err := o.opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("remove"))
	assert.Equal("apns", u.Get("type"))
	assert.Nil(err)
}

func TestNewRemoveChannelsFromPushBuilderWithContext(t *testing.T) {
	assert := assert.New(t)

	o := newRemoveChannelsFromPushBuilderWithContext(pubnub, backgroundContext)
	o.Channels([]string{"ch1", "ch2", "ch3"})
	o.DeviceIDForPush("deviceId")
	o.PushType(PNPushTypeAPNS)
	u, err := o.opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("remove"))
	assert.Equal("apns", u.Get("type"))
	assert.Nil(err)

}

func TestRemChannelsFromPushValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS

	assert.Equal("pubnub/validation: pubnub: Remove Channel Group: Missing Subscribe Key", opts.validate().Error())
}
