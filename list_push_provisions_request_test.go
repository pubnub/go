package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPushProvisionsRequestValidate(t *testing.T) {
	assert := assert.New(t)

	opts := &listPushProvisionsRequestOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	err := opts.validate()
	assert.Nil(err)

	opts1 := &listPushProvisionsRequestOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeNone,
		pubnub:          pubnub,
	}

	err1 := opts1.validate()
	assert.Contains(err1.Error(), "Missing Push Type")

	opts3 := &listPushProvisionsRequestOpts{
		PushType: PNPushTypeAPNS,
		pubnub:   pubnub,
	}

	err3 := opts3.validate()

	assert.Contains(err3.Error(), "Missing Device ID")

}

func TestListPushProvisionsRequestBuildPath(t *testing.T) {
	assert := assert.New(t)

	opts := &listPushProvisionsRequestOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	str, err := opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId", str)
	assert.Nil(err)

}

func TestNewListPushProvisionsRequestBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newListPushProvisionsRequestBuilder(pubnub)
	o.DeviceIDForPush("deviceId")
	o.PushType(PNPushTypeAPNS)
	str, err := o.opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId", str)
	assert.Nil(err)

}

func TestNewListPushProvisionsRequestBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newListPushProvisionsRequestBuilderWithContext(pubnub, backgroundContext)
	o.DeviceIDForPush("deviceId")
	o.PushType(PNPushTypeAPNS)
	str, err := o.opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId", str)
	assert.Nil(err)

}

func TestListPushProvisionsRequestBuildQuery(t *testing.T) {
	assert := assert.New(t)

	opts := &listPushProvisionsRequestOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	u, err := opts.buildQuery()
	assert.Equal("apns", u.Get("type"))
	assert.Nil(err)
}

func TestListPushProvisionsRequestBuildQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := &listPushProvisionsRequestOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	u, err := opts.buildQuery()
	assert.Equal("apns", u.Get("type"))
	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))
	assert.Nil(err)
}

func TestListPushProvisionsRequestBuildQueryParamTopicAndEnv(t *testing.T) {
	assert := assert.New(t)

	opts := &listPushProvisionsRequestOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
		Topic:           "a",
		Environment:     PNPushEnvironmentDevelopment,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	u, err := opts.buildQuery()
	assert.Equal("apns", u.Get("type"))
	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))
	assert.Equal("development", u.Get("environment"))
	assert.Equal("a", u.Get("topic"))

	assert.Nil(err)
}

func TestListPushProvisionsRequestBuildBody(t *testing.T) {
	assert := assert.New(t)

	opts := &listPushProvisionsRequestOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pubnub,
	}

	_, err := opts.buildBody()
	assert.Nil(err)

}

func TestListPushProvisionsNewListPushProvisionsRequestResponseErrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newListPushProvisionsRequestResponse(jsonBytes, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestListPushProvisionsValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &listPushProvisionsRequestOpts{
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
		pubnub:          pn,
	}

	assert.Equal("pubnub/validation: pubnub: Remove Channel Group: Missing Subscribe Key", opts.validate().Error())
}
