package pubnub

import (
	"fmt"
	"net/http"
	"strings"
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

	o := newRemoveChannelsFromPushBuilderWithContext(pubnub, pubnub.ctx)
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

	assert.Equal("pubnub/validation: pubnub: Remove Push From Channel: Missing Subscribe Key", opts.validate().Error())
}

// HTTP Method and Operation Tests

func TestRemoveChannelsFromPushHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)

	assert.Equal("GET", opts.httpMethod())
}

func TestRemoveChannelsFromPushOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)

	assert.Equal(PNRemovePushNotificationsFromChannelsOperation, opts.operationType())
}

func TestRemoveChannelsFromPushIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestRemoveChannelsFromPushTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Comprehensive Validation Tests

func TestRemoveChannelsFromPushValidationComprehensive(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name          string
		setupOpts     func(*removeChannelsFromPushOpts)
		expectedError string
	}{
		{
			name: "Missing subscribe key",
			setupOpts: func(opts *removeChannelsFromPushOpts) {
				opts.pubnub.Config.SubscribeKey = ""
				opts.Channels = []string{"channel1"}
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeAPNS
			},
			expectedError: "Missing Subscribe Key",
		},
		{
			name: "Missing channels",
			setupOpts: func(opts *removeChannelsFromPushOpts) {
				opts.Channels = []string{}
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeAPNS
			},
			expectedError: "Missing Channel",
		},
		{
			name: "Nil channels",
			setupOpts: func(opts *removeChannelsFromPushOpts) {
				opts.Channels = nil
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeAPNS
			},
			expectedError: "Missing Channel",
		},
		{
			name: "Missing device ID",
			setupOpts: func(opts *removeChannelsFromPushOpts) {
				opts.Channels = []string{"channel1"}
				opts.DeviceIDForPush = ""
				opts.PushType = PNPushTypeAPNS
			},
			expectedError: "Missing Device ID",
		},
		{
			name: "Missing push type",
			setupOpts: func(opts *removeChannelsFromPushOpts) {
				opts.Channels = []string{"channel1"}
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeNone
			},
			expectedError: "Missing Push Type",
		},
		{
			name: "APNS2 missing topic",
			setupOpts: func(opts *removeChannelsFromPushOpts) {
				opts.Channels = []string{"channel1"}
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeAPNS2
				opts.Topic = ""
			},
			expectedError: "Missing Push Topic",
		},
		{
			name: "Valid APNS configuration",
			setupOpts: func(opts *removeChannelsFromPushOpts) {
				opts.Channels = []string{"channel1"}
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeAPNS
			},
			expectedError: "",
		},
		{
			name: "Valid APNS2 configuration",
			setupOpts: func(opts *removeChannelsFromPushOpts) {
				opts.Channels = []string{"channel1"}
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeAPNS2
				opts.Topic = "com.example.app"
				opts.Environment = PNPushEnvironmentProduction
			},
			expectedError: "",
		},
		{
			name: "Valid GCM configuration",
			setupOpts: func(opts *removeChannelsFromPushOpts) {
				opts.Channels = []string{"channel1"}
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeGCM
			},
			expectedError: "",
		},
		{
			name: "Valid MPNS configuration",
			setupOpts: func(opts *removeChannelsFromPushOpts) {
				opts.Channels = []string{"channel1"}
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeMPNS
			},
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh PubNub instance for each test case to avoid shared state
			pn := NewPubNub(NewDemoConfig())
			opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			err := opts.validate()
			if tc.expectedError == "" {
				assert.Nil(err)
			} else {
				assert.NotNil(err)
				assert.Contains(err.Error(), tc.expectedError)
			}
		})
	}
}

// Systematic Builder Pattern Tests (6 setters)

func TestRemoveChannelsFromPushBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelsFromPushBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestRemoveChannelsFromPushBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelsFromPushBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestRemoveChannelsFromPushBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelsFromPushBuilder(pn)

	// Test Channels setter
	channels := []string{"channel1", "channel2"}
	builder.Channels(channels)
	assert.Equal(channels, builder.opts.Channels)

	// Test PushType setter
	pushType := PNPushTypeAPNS2
	builder.PushType(pushType)
	assert.Equal(pushType, builder.opts.PushType)

	// Test DeviceIDForPush setter
	deviceID := "device-123-abc"
	builder.DeviceIDForPush(deviceID)
	assert.Equal(deviceID, builder.opts.DeviceIDForPush)

	// Test Topic setter
	topic := "com.example.myapp"
	builder.Topic(topic)
	assert.Equal(topic, builder.opts.Topic)

	// Test Environment setter
	environment := PNPushEnvironmentProduction
	builder.Environment(environment)
	assert.Equal(environment, builder.opts.Environment)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Test Transport setter
	transport := &http.Transport{}
	builder.Transport(transport)
	assert.Equal(transport, builder.opts.Transport)
}

func TestRemoveChannelsFromPushBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	channels := []string{"channel1", "channel2"}
	pushType := PNPushTypeAPNS2
	deviceID := "device-123"
	topic := "com.example.app"
	environment := PNPushEnvironmentDevelopment
	queryParam := map[string]string{"key": "value"}
	transport := &http.Transport{}

	builder := newRemoveChannelsFromPushBuilder(pn)
	result := builder.Channels(channels).
		PushType(pushType).
		DeviceIDForPush(deviceID).
		Topic(topic).
		Environment(environment).
		QueryParam(queryParam).
		Transport(transport)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(pushType, builder.opts.PushType)
	assert.Equal(deviceID, builder.opts.DeviceIDForPush)
	assert.Equal(topic, builder.opts.Topic)
	assert.Equal(environment, builder.opts.Environment)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

func TestRemoveChannelsFromPushBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelsFromPushBuilder(pn)

	// Verify default values (zero value of PNPushType is 0, not PNPushTypeNone which is 1)
	assert.Equal(PNPushType(0), builder.opts.PushType)
	assert.Nil(builder.opts.Channels)
	assert.Empty(builder.opts.DeviceIDForPush)
	assert.Empty(builder.opts.Topic)
	assert.Empty(builder.opts.Environment) // Default environment is empty
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestRemoveChannelsFromPushBuilderPushTypeCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		pushType    PNPushType
		description string
	}{
		{
			name:        "APNS push type",
			pushType:    PNPushTypeAPNS,
			description: "iOS push notifications (legacy)",
		},
		{
			name:        "APNS2 push type",
			pushType:    PNPushTypeAPNS2,
			description: "iOS push notifications (HTTP/2)",
		},
		{
			name:        "GCM push type",
			pushType:    PNPushTypeGCM,
			description: "Google Cloud Messaging",
		},
		{
			name:        "MPNS push type",
			pushType:    PNPushTypeMPNS,
			description: "Microsoft Push Notification Service",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveChannelsFromPushBuilder(pn)
			builder.PushType(tc.pushType)

			assert.Equal(tc.pushType, builder.opts.PushType)
		})
	}
}

func TestRemoveChannelsFromPushBuilderChannelCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		channels    []string
		description string
	}{
		{
			name:        "Single channel",
			channels:    []string{"channel1"},
			description: "Basic single channel",
		},
		{
			name:        "Multiple channels",
			channels:    []string{"channel1", "channel2", "channel3"},
			description: "Multiple channels for batch removal",
		},
		{
			name:        "Channels with special characters",
			channels:    []string{"channel@domain.com", "channel#with$symbols"},
			description: "Channels with special characters",
		},
		{
			name:        "Channels with Unicode",
			channels:    []string{"频道中文", "канал-русский"},
			description: "Channels with Unicode characters",
		},
		{
			name:        "Empty channel list",
			channels:    []string{},
			description: "Empty channel list (should fail validation)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveChannelsFromPushBuilder(pn)
			builder.Channels(tc.channels)

			assert.Equal(tc.channels, builder.opts.Channels)
		})
	}
}

func TestRemoveChannelsFromPushBuilderDeviceIDCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		deviceID    string
		description string
	}{
		{
			name:        "Simple device ID",
			deviceID:    "device123",
			description: "Basic alphanumeric device ID",
		},
		{
			name:        "UUID device ID",
			deviceID:    "550e8400-e29b-41d4-a716-446655440000",
			description: "UUID format device ID",
		},
		{
			name:        "Device ID with special characters",
			deviceID:    "device@domain.com",
			description: "Device ID with special characters",
		},
		{
			name:        "Device ID with Unicode",
			deviceID:    "设备123",
			description: "Device ID with Unicode characters",
		},
		{
			name:        "Long device ID",
			deviceID:    strings.Repeat("a", 100),
			description: "Very long device ID",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveChannelsFromPushBuilder(pn)
			builder.DeviceIDForPush(tc.deviceID)

			assert.Equal(tc.deviceID, builder.opts.DeviceIDForPush)
		})
	}
}

func TestRemoveChannelsFromPushBuilderAPNS2Combinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		topic       string
		environment PNPushEnvironment
		description string
	}{
		{
			name:        "Production APNS2",
			topic:       "com.example.myapp",
			environment: PNPushEnvironmentProduction,
			description: "Production environment APNS2 configuration",
		},
		{
			name:        "Development APNS2",
			topic:       "com.example.myapp.dev",
			environment: PNPushEnvironmentDevelopment,
			description: "Development environment APNS2 configuration",
		},
		{
			name:        "Complex topic",
			topic:       "com.company.mobile.application.notifications",
			environment: PNPushEnvironmentProduction,
			description: "Complex topic identifier",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveChannelsFromPushBuilder(pn)
			builder.Topic(tc.topic).Environment(tc.environment)

			assert.Equal(tc.topic, builder.opts.Topic)
			assert.Equal(tc.environment, builder.opts.Environment)
		})
	}
}

func TestRemoveChannelsFromPushBuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	channels := []string{"channel1", "channel2"}
	pushType := PNPushTypeAPNS2
	deviceID := "device-abc123"
	topic := "com.example.app"
	environment := PNPushEnvironmentProduction
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}
	transport := &http.Transport{}

	// Test all 6 setters in chain
	builder := newRemoveChannelsFromPushBuilder(pn).
		Channels(channels).
		PushType(pushType).
		DeviceIDForPush(deviceID).
		Topic(topic).
		Environment(environment).
		QueryParam(queryParam).
		Transport(transport)

	// Verify all are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(pushType, builder.opts.PushType)
	assert.Equal(deviceID, builder.opts.DeviceIDForPush)
	assert.Equal(topic, builder.opts.Topic)
	assert.Equal(environment, builder.opts.Environment)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests (CRITICAL: APNS1 vs APNS2 paths)

func TestRemoveChannelsFromPushBuildPathAPNS(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeAPNS

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/push/sub-key/demo/devices/device123"
	assert.Equal(expected, path)
}

func TestRemoveChannelsFromPushBuildPathAPNS2(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeAPNS2

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/push/sub-key/demo/devices-apns2/device123"
	assert.Equal(expected, path)
}

func TestRemoveChannelsFromPushBuildPathGCM(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeGCM

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/push/sub-key/demo/devices/device123"
	assert.Equal(expected, path)
}

func TestRemoveChannelsFromPushBuildPathMPNS(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeMPNS

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/push/sub-key/demo/devices/device123"
	assert.Equal(expected, path)
}

func TestRemoveChannelsFromPushBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "my-device"
	opts.PushType = PNPushTypeAPNS

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/push/sub-key/custom-sub-key/devices/my-device"
	assert.Equal(expected, path)
}

func TestRemoveChannelsFromPushBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device@domain.com"
	opts.PushType = PNPushTypeAPNS

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/push/sub-key/demo/devices/device%40domain.com"
	assert.Equal(expected, path)
}

func TestRemoveChannelsFromPushBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "设备123"
	opts.PushType = PNPushTypeAPNS2

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/push/sub-key/demo/devices-apns2/%E8%AE%BE%E5%A4%87123"
	assert.Equal(expected, path)
}

// JSON Body Building Tests (CRITICAL for GET operation - should be empty)

func TestRemoveChannelsFromPushBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations should have empty body
	assert.Equal([]byte{}, body)
}

func TestRemoveChannelsFromPushBuildBodyWithAllParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)

	// Set all possible parameters - body should still be empty for GET
	opts.Channels = []string{"channel1", "channel2"}
	opts.PushType = PNPushTypeAPNS2
	opts.DeviceIDForPush = "device123"
	opts.Topic = "com.example.app"
	opts.Environment = PNPushEnvironmentProduction
	opts.QueryParam = map[string]string{"param": "value"}
	opts.Transport = &http.Transport{}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations always have empty body regardless of parameters
	assert.Equal([]byte{}, body)
}

func TestRemoveChannelsFromPushBuildBodyErrorScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)

	// Even with potential error conditions, buildBody should not fail for GET
	opts.Channels = []string{}     // Empty channels
	opts.DeviceIDForPush = ""      // Empty device ID
	opts.PushType = PNPushTypeNone // Invalid push type

	body, err := opts.buildBody()
	assert.Nil(err) // buildBody should never error for GET operations
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}

// Query Parameter Tests (Push Type Specific and Channel Removal)

func TestRemoveChannelsFromPushBuildQueryAPNS(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1", "channel2"}
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeAPNS

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have push type, channels and default parameters
	assert.Equal("apns", query.Get("type"))
	assert.Equal("channel1,channel2", query.Get("remove"))
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestRemoveChannelsFromPushBuildQueryAPNS2(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1"}
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeAPNS2
	opts.Topic = "com.example.app"
	opts.Environment = PNPushEnvironmentProduction

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Should have push type, channels, topic, environment and default parameters
	assert.Equal("apns2", query.Get("type"))
	assert.Equal("channel1", query.Get("remove"))
	assert.Equal("com.example.app", query.Get("topic"))
	assert.Equal("production", query.Get("environment"))
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestRemoveChannelsFromPushBuildQueryGCM(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1", "channel2", "channel3"}
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeGCM

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("gcm", query.Get("type"))
	assert.Equal("channel1,channel2,channel3", query.Get("remove"))
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestRemoveChannelsFromPushBuildQueryMPNS(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1"}
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeMPNS

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("mpns", query.Get("type"))
	assert.Equal("channel1", query.Get("remove"))
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestRemoveChannelsFromPushBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
	opts.Channels = []string{"channel1"}
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeAPNS

	opts.QueryParam = map[string]string{
		"custom":         "value",
		"special_chars":  "value@with#symbols",
		"unicode":        "测试参数",
		"empty_value":    "",
		"number_string":  "42",
		"boolean_string": "true",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify push type, channels, and custom parameters (note: special chars and Unicode are URL-encoded in query strings)
	assert.Equal("apns", query.Get("type"))
	assert.Equal("channel1", query.Get("remove"))
	assert.Equal("value", query.Get("custom"))
	assert.Equal("value%40with%23symbols", query.Get("special_chars"))
	assert.Equal("%E6%B5%8B%E8%AF%95%E5%8F%82%E6%95%B0", query.Get("unicode"))
	assert.Equal("", query.Get("empty_value"))
	assert.Equal("42", query.Get("number_string"))
	assert.Equal("true", query.Get("boolean_string"))
}

func TestRemoveChannelsFromPushBuildQueryEnvironmentCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		environment PNPushEnvironment
		expected    string
	}{
		{
			name:        "Development environment",
			environment: PNPushEnvironmentDevelopment,
			expected:    "development",
		},
		{
			name:        "Production environment",
			environment: PNPushEnvironmentProduction,
			expected:    "production",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
			opts.Channels = []string{"channel1"}
			opts.DeviceIDForPush = "device123"
			opts.PushType = PNPushTypeAPNS2
			opts.Topic = "com.example.app"
			opts.Environment = tc.environment

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expected, query.Get("environment"))
		})
	}
}

// Channel Encoding Tests (CRITICAL: comma-separated channels in query param)

func TestRemoveChannelsFromPushChannelEncoding(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name           string
		channels       []string
		expectedRemove string
		description    string
	}{
		{
			name:           "Single channel",
			channels:       []string{"channel1"},
			expectedRemove: "channel1",
			description:    "Single channel should not have commas",
		},
		{
			name:           "Multiple channels",
			channels:       []string{"channel1", "channel2", "channel3"},
			expectedRemove: "channel1,channel2,channel3",
			description:    "Multiple channels should be comma-separated",
		},
		{
			name:           "Channels with special characters",
			channels:       []string{"channel@domain.com", "channel#with$symbols"},
			expectedRemove: "channel@domain.com,channel#with$symbols",
			description:    "Special characters should not affect comma separation",
		},
		{
			name:           "Channels with Unicode",
			channels:       []string{"频道中文", "канал-русский", "チャンネル"},
			expectedRemove: "频道中文,канал-русский,チャンネル",
			description:    "Unicode should not affect comma separation",
		},
		{
			name:           "Channels with spaces",
			channels:       []string{"channel with spaces", "another channel"},
			expectedRemove: "channel with spaces,another channel",
			description:    "Spaces should not affect comma separation",
		},
		{
			name:           "Many channels",
			channels:       []string{"ch1", "ch2", "ch3", "ch4", "ch5", "ch6", "ch7", "ch8", "ch9", "ch10"},
			expectedRemove: "ch1,ch2,ch3,ch4,ch5,ch6,ch7,ch8,ch9,ch10",
			description:    "Many channels should be properly comma-separated",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
			opts.Channels = tc.channels
			opts.DeviceIDForPush = "device123"
			opts.PushType = PNPushTypeAPNS

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expectedRemove, query.Get("remove"))
		})
	}
}

// Push Type Specific Tests

func TestRemoveChannelsFromPushAPNSPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelsFromPushBuilder(pn)
	builder.Channels([]string{"channel1"})
	builder.DeviceIDForPush("ios-device-token")
	builder.PushType(PNPushTypeAPNS)

	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v1/push/")
	assert.Contains(path, "/devices/")
	assert.NotContains(path, "devices-apns2")
}

func TestRemoveChannelsFromPushAPNS2Path(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelsFromPushBuilder(pn)
	builder.Channels([]string{"channel1"})
	builder.DeviceIDForPush("ios-device-token")
	builder.PushType(PNPushTypeAPNS2)
	builder.Topic("com.example.app")

	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/push/")
	assert.Contains(path, "/devices-apns2/")
	assert.NotContains(path, "/v1/push/")
}

func TestRemoveChannelsFromPushGCMConfiguration(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelsFromPushBuilder(pn)
	builder.Channels([]string{"channel1", "channel2"})
	builder.DeviceIDForPush("gcm-registration-token")
	builder.PushType(PNPushTypeGCM)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should use v1 path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v1/push/")

	// Should have correct query type and channels
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("gcm", query.Get("type"))
	assert.Equal("channel1,channel2", query.Get("remove"))
}

func TestRemoveChannelsFromPushMPNSConfiguration(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveChannelsFromPushBuilder(pn)
	builder.Channels([]string{"channel1"})
	builder.DeviceIDForPush("mpns-device-token")
	builder.PushType(PNPushTypeMPNS)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should use v1 path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v1/push/")

	// Should have correct query type and channels
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("mpns", query.Get("type"))
	assert.Equal("channel1", query.Get("remove"))
}

// Device ID Encoding Tests

func TestRemoveChannelsFromPushDeviceIDEncoding(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name         string
		deviceID     string
		expectedPath string
		description  string
	}{
		{
			name:         "Simple device ID",
			deviceID:     "abc123",
			expectedPath: "/v1/push/sub-key/demo/devices/abc123",
			description:  "No encoding needed",
		},
		{
			name:         "Device ID with @ symbol",
			deviceID:     "device@domain.com",
			expectedPath: "/v1/push/sub-key/demo/devices/device%40domain.com",
			description:  "@ symbol should be URL encoded",
		},
		{
			name:         "Device ID with spaces",
			deviceID:     "device 123",
			expectedPath: "/v1/push/sub-key/demo/devices/device%20123",
			description:  "Spaces should be URL encoded",
		},
		{
			name:         "Device ID with special characters",
			deviceID:     "device!@#$%^&*()",
			expectedPath: "/v1/push/sub-key/demo/devices/device%21%40%23%24%25%5E%26%2A%28%29",
			description:  "Special characters should be URL encoded",
		},
		{
			name:         "Unicode device ID",
			deviceID:     "设备123",
			expectedPath: "/v1/push/sub-key/demo/devices/%E8%AE%BE%E5%A4%87123",
			description:  "Unicode should be URL encoded",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newRemoveChannelsFromPushOpts(pn, pn.ctx)
			opts.Channels = []string{"channel1"}
			opts.DeviceIDForPush = tc.deviceID
			opts.PushType = PNPushTypeAPNS

			path, err := opts.buildPath()
			assert.Nil(err)
			assert.Equal(tc.expectedPath, path)
		})
	}
}

// Comprehensive Edge Case Tests

func TestRemoveChannelsFromPushWithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*removeChannelsFromPushBuilder)
	}{
		{
			name: "Very many channels",
			setupFn: func(builder *removeChannelsFromPushBuilder) {
				manyChannels := make([]string, 100)
				for i := 0; i < 100; i++ {
					manyChannels[i] = fmt.Sprintf("channel_%d", i)
				}
				builder.Channels(manyChannels)
				builder.DeviceIDForPush("device123")
				builder.PushType(PNPushTypeAPNS)
			},
		},
		{
			name: "Very long device ID",
			setupFn: func(builder *removeChannelsFromPushBuilder) {
				longDeviceID := strings.Repeat("device", 250) // 1500 characters
				builder.Channels([]string{"channel1"})
				builder.DeviceIDForPush(longDeviceID)
				builder.PushType(PNPushTypeAPNS)
			},
		},
		{
			name: "Very long topic",
			setupFn: func(builder *removeChannelsFromPushBuilder) {
				longTopic := "com.example." + strings.Repeat("app", 100)
				builder.Channels([]string{"channel1"})
				builder.DeviceIDForPush("device123")
				builder.PushType(PNPushTypeAPNS2)
				builder.Topic(longTopic)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *removeChannelsFromPushBuilder) {
				largeQueryParam := make(map[string]string)
				for i := 0; i < 100; i++ {
					largeQueryParam[fmt.Sprintf("param_%d", i)] = strings.Repeat(fmt.Sprintf("value_%d_", i), 20)
				}
				builder.Channels([]string{"channel1"})
				builder.DeviceIDForPush("device123")
				builder.PushType(PNPushTypeGCM)
				builder.QueryParam(largeQueryParam)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveChannelsFromPushBuilder(pn)
			tc.setupFn(builder)

			// Should build valid path and query (though validation might fail for incomplete configs)
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Should always have empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
		})
	}
}

func TestRemoveChannelsFromPushSpecialCharacterHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialStrings := []string{
		"<script>alert('xss')</script>",
		"SELECT * FROM channels; DROP TABLE channels;",
		"newline\ncharacter\ttab\rcarriage",
		"   ",                // Only spaces
		"\u0000\u0001\u0002", // Control characters
		"\"quoted_string\"",
		"'single_quoted'",
		"back`tick`string",
		"emoji😀🎉🚀💯",
		"ñáéíóúüç", // Accented characters
	}

	for i, specialString := range specialStrings {
		t.Run(fmt.Sprintf("SpecialString_%d", i), func(t *testing.T) {
			builder := newRemoveChannelsFromPushBuilder(pn)
			builder.Channels([]string{specialString, "normal_channel"})
			builder.DeviceIDForPush(specialString)
			builder.PushType(PNPushTypeAPNS)
			builder.Topic(specialString)
			builder.QueryParam(map[string]string{
				"special_param": specialString,
			})

			// Should build valid path and query
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Should always have empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
		})
	}
}

func TestRemoveChannelsFromPushParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		channels    []string
		deviceID    string
		pushType    PNPushType
		description string
	}{
		{
			name:        "Empty channel list",
			channels:    []string{},
			deviceID:    "device123",
			pushType:    PNPushTypeAPNS,
			description: "Empty channel list should fail validation",
		},
		{
			name:        "Single character channel",
			channels:    []string{"a"},
			deviceID:    "device123",
			pushType:    PNPushTypeGCM,
			description: "Single character channel",
		},
		{
			name:        "Unicode-only channels",
			channels:    []string{"测试", "канал"},
			deviceID:    "device123",
			pushType:    PNPushTypeMPNS,
			description: "Channels with Unicode characters",
		},
		{
			name:        "Very long channel names",
			channels:    []string{strings.Repeat("channel", 100)},
			deviceID:    "device123",
			pushType:    PNPushTypeAPNS2,
			description: "Very long channel name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveChannelsFromPushBuilder(pn)
			builder.Channels(tc.channels)
			builder.DeviceIDForPush(tc.deviceID)
			builder.PushType(tc.pushType)
			if tc.pushType == PNPushTypeAPNS2 {
				builder.Topic("com.example.app") // Required for APNS2
			}

			// Should build valid components (validation may fail for edge cases)
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/push/sub-key/")

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body) // GET operation always has empty body
		})
	}
}

// Error Scenario Tests

func TestRemoveChannelsFromPushExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newRemoveChannelsFromPushBuilder(pn)
	builder.Channels([]string{"channel1"})
	builder.DeviceIDForPush("device123")
	builder.PushType(PNPushTypeAPNS)

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestRemoveChannelsFromPushBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newRemoveChannelsFromPushBuilder(pn)

	channels := []string{"channel1", "channel2"}
	pushType := PNPushTypeAPNS2
	deviceID := "device-abc123"
	topic := "com.example.app"
	environment := PNPushEnvironmentProduction
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}
	transport := &http.Transport{}

	// Set all possible parameters
	builder.Channels(channels).
		PushType(pushType).
		DeviceIDForPush(deviceID).
		Topic(topic).
		Environment(environment).
		QueryParam(queryParam).
		Transport(transport)

	// Verify all values are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(pushType, builder.opts.PushType)
	assert.Equal(deviceID, builder.opts.DeviceIDForPush)
	assert.Equal(topic, builder.opts.Topic)
	assert.Equal(environment, builder.opts.Environment)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path (APNS2)
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v2/push/sub-key/demo/devices-apns2/device-abc123"
	assert.Equal(expectedPath, path)

	// Should build query with all params
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("apns2", query.Get("type"))
	assert.Equal("channel1,channel2", query.Get("remove"))
	assert.Equal("com.example.app", query.Get("topic"))
	assert.Equal("production", query.Get("environment"))
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))

	// Should always have empty body (GET operation)
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}
