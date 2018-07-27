package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestRemoveChannelsFromPushRequestValidate(t *testing.T) {
	assert := assert.New(t)

	opts := &removeChannelsFromPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	err := opts.validate()
	assert.Nil(err)

	opts1 := &removeChannelsFromPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeNone,
		pubnub:          pubnub,
	}

	err1 := opts1.validate()
	assert.Contains(err1.Error(), "Missing Push Type")

	opts2 := &removeChannelsFromPushOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	err2 := opts2.validate()
	assert.Contains(err2.Error(), "Missing Channel")

	opts3 := &removeChannelsFromPushOpts{
		Channels: []string{"ch1", "ch2", "ch3"},
		PushType: PNPushTypeAPNS,
		pubnub:   pubnub,
	}

	err3 := opts3.validate()
	assert.Contains(err3.Error(), "Missing Device ID")
}

func TestRemoveChannelsFromPushRequestBuildPath(t *testing.T) {
	assert := assert.New(t)

	opts := &removeChannelsFromPushOpts{
		DeviceIDForPush: "deviceId",
		Channels:        []string{"ch1", "ch2", "ch3"},
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	str, err := opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId", str)
	assert.Nil(err)

}

func TestRemoveChannelsFromPushRequestBuildQuery(t *testing.T) {
	assert := assert.New(t)

	opts := &removeChannelsFromPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("remove"))
	assert.Equal("apns", u.Get("type"))
	assert.Nil(err)
}

func TestRemoveChannelsFromPushRequestBuildBody(t *testing.T) {
	assert := assert.New(t)

	opts := &removeChannelsFromPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

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

	o := newRemoveChannelsFromPushBuilderWithContext(pubnub, context.Background())
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
	opts := &removeChannelsFromPushOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pn,
	}

	assert.Equal("pubnub/validation: pubnub: \x0e: Missing Subscribe Key", opts.validate().Error())
}
