package pubnub

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPushProvisionsRequestValidate(t *testing.T) {
	assert := assert.New(t)

	opts := newListPushProvisionsRequestOpts(pubnub, pubnub.ctx)
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS

	err := opts.validate()
	assert.Nil(err)

	opts1 := newListPushProvisionsRequestOpts(pubnub, pubnub.ctx)
	opts1.DeviceIDForPush = "deviceId"
	opts1.PushType = PNPushTypeNone

	err1 := opts1.validate()
	assert.Contains(err1.Error(), "Missing Push Type")

	opts3 := newListPushProvisionsRequestOpts(pubnub, pubnub.ctx)
	opts3.PushType = PNPushTypeAPNS

	err3 := opts3.validate()

	assert.Contains(err3.Error(), "Missing Device ID")

}

func TestListPushProvisionsRequestBuildPath(t *testing.T) {
	assert := assert.New(t)

	opts := newListPushProvisionsRequestOpts(pubnub, pubnub.ctx)
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS

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

	o := newListPushProvisionsRequestBuilderWithContext(pubnub, pubnub.ctx)
	o.DeviceIDForPush("deviceId")
	o.PushType(PNPushTypeAPNS)
	str, err := o.opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId", str)
	assert.Nil(err)

}

func TestListPushProvisionsRequestBuildQuery(t *testing.T) {
	assert := assert.New(t)

	opts := newListPushProvisionsRequestOpts(pubnub, pubnub.ctx)
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS

	u, err := opts.buildQuery()
	assert.Equal("apns", u.Get("type"))
	assert.Nil(err)
}

func TestListPushProvisionsRequestBuildQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := newListPushProvisionsRequestOpts(pubnub, pubnub.ctx)
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS

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

	opts := newListPushProvisionsRequestOpts(pubnub, pubnub.ctx)
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS
	opts.Topic = "a"
	opts.Environment = PNPushEnvironmentDevelopment

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

	opts := newListPushProvisionsRequestOpts(pubnub, pubnub.ctx)
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS

	_, err := opts.buildBody()
	assert.Nil(err)

}

func TestListPushProvisionsNewListPushProvisionsRequestResponseErrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newListPushProvisionsRequestResponse(jsonBytes, StatusResponse{})
	assert.Equal("pubnub/parsing: error unmarshalling response: {s}", err.Error())
}

func TestListPushProvisionsValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "deviceId"
	opts.PushType = PNPushTypeAPNS

	assert.Equal("pubnub/validation: pubnub: List Push Enabled Channels: Missing Subscribe Key", opts.validate().Error())
}

// HTTP Method and Operation Tests

func TestListPushProvisionsHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)

	assert.Equal("GET", opts.httpMethod())
}

func TestListPushProvisionsOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)

	assert.Equal(PNPushNotificationsEnabledChannelsOperation, opts.operationType())
}

func TestListPushProvisionsIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestListPushProvisionsTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Comprehensive Validation Tests

func TestListPushProvisionsValidationComprehensive(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name          string
		setupOpts     func(*listPushProvisionsRequestOpts)
		expectedError string
	}{
		{
			name: "Missing subscribe key",
			setupOpts: func(opts *listPushProvisionsRequestOpts) {
				opts.pubnub.Config.SubscribeKey = ""
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeAPNS
			},
			expectedError: "Missing Subscribe Key",
		},
		{
			name: "Missing device ID",
			setupOpts: func(opts *listPushProvisionsRequestOpts) {
				opts.DeviceIDForPush = ""
				opts.PushType = PNPushTypeAPNS
			},
			expectedError: "Missing Device ID",
		},
		{
			name: "Missing push type",
			setupOpts: func(opts *listPushProvisionsRequestOpts) {
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeNone
			},
			expectedError: "Missing Push Type",
		},
		{
			name: "APNS2 missing topic",
			setupOpts: func(opts *listPushProvisionsRequestOpts) {
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeAPNS2
				opts.Topic = ""
			},
			expectedError: "Missing Push Topic",
		},
		{
			name: "Valid APNS configuration",
			setupOpts: func(opts *listPushProvisionsRequestOpts) {
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeAPNS
			},
			expectedError: "",
		},
		{
			name: "Valid APNS2 configuration",
			setupOpts: func(opts *listPushProvisionsRequestOpts) {
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeAPNS2
				opts.Topic = "com.example.app"
				opts.Environment = PNPushEnvironmentProduction
			},
			expectedError: "",
		},
		{
			name: "Valid GCM configuration",
			setupOpts: func(opts *listPushProvisionsRequestOpts) {
				opts.DeviceIDForPush = "device123"
				opts.PushType = PNPushTypeGCM
			},
			expectedError: "",
		},
		{
			name: "Valid MPNS configuration",
			setupOpts: func(opts *listPushProvisionsRequestOpts) {
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
			opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
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

func TestListPushProvisionsBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newListPushProvisionsRequestBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestListPushProvisionsBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newListPushProvisionsRequestBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestListPushProvisionsBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newListPushProvisionsRequestBuilder(pn)

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

func TestListPushProvisionsBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	pushType := PNPushTypeAPNS2
	deviceID := "device-123"
	topic := "com.example.app"
	environment := PNPushEnvironmentDevelopment
	queryParam := map[string]string{"key": "value"}
	transport := &http.Transport{}

	builder := newListPushProvisionsRequestBuilder(pn)
	result := builder.PushType(pushType).
		DeviceIDForPush(deviceID).
		Topic(topic).
		Environment(environment).
		QueryParam(queryParam).
		Transport(transport)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal(pushType, builder.opts.PushType)
	assert.Equal(deviceID, builder.opts.DeviceIDForPush)
	assert.Equal(topic, builder.opts.Topic)
	assert.Equal(environment, builder.opts.Environment)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

func TestListPushProvisionsBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newListPushProvisionsRequestBuilder(pn)

	// Verify default values (zero value of PNPushType is 0, not PNPushTypeNone which is 1)
	assert.Equal(PNPushType(0), builder.opts.PushType)
	assert.Empty(builder.opts.DeviceIDForPush)
	assert.Empty(builder.opts.Topic)
	assert.Empty(builder.opts.Environment) // Default environment is empty
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestListPushProvisionsBuilderPushTypeCombinations(t *testing.T) {
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
			builder := newListPushProvisionsRequestBuilder(pn)
			builder.PushType(tc.pushType)

			assert.Equal(tc.pushType, builder.opts.PushType)
		})
	}
}

func TestListPushProvisionsBuilderDeviceIDCombinations(t *testing.T) {
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
			builder := newListPushProvisionsRequestBuilder(pn)
			builder.DeviceIDForPush(tc.deviceID)

			assert.Equal(tc.deviceID, builder.opts.DeviceIDForPush)
		})
	}
}

func TestListPushProvisionsBuilderAPNS2Combinations(t *testing.T) {
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
			builder := newListPushProvisionsRequestBuilder(pn)
			builder.Topic(tc.topic).Environment(tc.environment)

			assert.Equal(tc.topic, builder.opts.Topic)
			assert.Equal(tc.environment, builder.opts.Environment)
		})
	}
}

func TestListPushProvisionsBuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

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
	builder := newListPushProvisionsRequestBuilder(pn).
		PushType(pushType).
		DeviceIDForPush(deviceID).
		Topic(topic).
		Environment(environment).
		QueryParam(queryParam).
		Transport(transport)

	// Verify all are set correctly
	assert.Equal(pushType, builder.opts.PushType)
	assert.Equal(deviceID, builder.opts.DeviceIDForPush)
	assert.Equal(topic, builder.opts.Topic)
	assert.Equal(environment, builder.opts.Environment)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests (CRITICAL: APNS1 vs APNS2 paths)

func TestListPushProvisionsBuildPathAPNS(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeAPNS

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/push/sub-key/demo/devices/device123"
	assert.Equal(expected, path)
}

func TestListPushProvisionsBuildPathAPNS2(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeAPNS2

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/push/sub-key/demo/devices-apns2/device123"
	assert.Equal(expected, path)
}

func TestListPushProvisionsBuildPathGCM(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeGCM

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/push/sub-key/demo/devices/device123"
	assert.Equal(expected, path)
}

func TestListPushProvisionsBuildPathMPNS(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeMPNS

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/push/sub-key/demo/devices/device123"
	assert.Equal(expected, path)
}

func TestListPushProvisionsBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "my-device"
	opts.PushType = PNPushTypeAPNS

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/push/sub-key/custom-sub-key/devices/my-device"
	assert.Equal(expected, path)
}

func TestListPushProvisionsBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device@domain.com"
	opts.PushType = PNPushTypeAPNS

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v1/push/sub-key/demo/devices/device%40domain.com"
	assert.Equal(expected, path)
}

func TestListPushProvisionsBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "设备123"
	opts.PushType = PNPushTypeAPNS2

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/push/sub-key/demo/devices-apns2/%E8%AE%BE%E5%A4%87123"
	assert.Equal(expected, path)
}

// JSON Body Building Tests (CRITICAL for GET operation - should be empty)

func TestListPushProvisionsBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations should have empty body
	assert.Equal([]byte{}, body)
}

func TestListPushProvisionsBuildBodyWithAllParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)

	// Set all possible parameters - body should still be empty for GET
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

func TestListPushProvisionsBuildBodyErrorScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)

	// Even with potential error conditions, buildBody should not fail for GET
	opts.DeviceIDForPush = ""      // Empty device ID
	opts.PushType = PNPushTypeNone // Invalid push type

	body, err := opts.buildBody()
	assert.Nil(err) // buildBody should never error for GET operations
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}

// Query Parameter Tests (Push Type Specific)

func TestListPushProvisionsBuildQueryAPNS(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeAPNS

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have push type and default parameters
	assert.Equal("apns", query.Get("type"))
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestListPushProvisionsBuildQueryAPNS2(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeAPNS2
	opts.Topic = "com.example.app"
	opts.Environment = PNPushEnvironmentProduction

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Should have push type, topic, environment and default parameters
	assert.Equal("apns2", query.Get("type"))
	assert.Equal("com.example.app", query.Get("topic"))
	assert.Equal("production", query.Get("environment"))
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestListPushProvisionsBuildQueryGCM(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeGCM

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("gcm", query.Get("type"))
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestListPushProvisionsBuildQueryMPNS(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
	opts.DeviceIDForPush = "device123"
	opts.PushType = PNPushTypeMPNS

	query, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("mpns", query.Get("type"))
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestListPushProvisionsBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
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

	// Verify push type and custom parameters (note: special chars and Unicode are URL-encoded in query strings)
	assert.Equal("apns", query.Get("type"))
	assert.Equal("value", query.Get("custom"))
	assert.Equal("value%40with%23symbols", query.Get("special_chars"))
	assert.Equal("%E6%B5%8B%E8%AF%95%E5%8F%82%E6%95%B0", query.Get("unicode"))
	assert.Equal("", query.Get("empty_value"))
	assert.Equal("42", query.Get("number_string"))
	assert.Equal("true", query.Get("boolean_string"))
}

func TestListPushProvisionsBuildQueryEnvironmentCombinations(t *testing.T) {
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
			opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
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

// Push Type Specific Tests

func TestListPushProvisionsAPNSPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newListPushProvisionsRequestBuilder(pn)
	builder.DeviceIDForPush("ios-device-token")
	builder.PushType(PNPushTypeAPNS)

	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v1/push/")
	assert.Contains(path, "/devices/")
	assert.NotContains(path, "devices-apns2")
}

func TestListPushProvisionsAPNS2Path(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newListPushProvisionsRequestBuilder(pn)
	builder.DeviceIDForPush("ios-device-token")
	builder.PushType(PNPushTypeAPNS2)
	builder.Topic("com.example.app")

	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/push/")
	assert.Contains(path, "/devices-apns2/")
	assert.NotContains(path, "/v1/push/")
}

func TestListPushProvisionsGCMConfiguration(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newListPushProvisionsRequestBuilder(pn)
	builder.DeviceIDForPush("gcm-registration-token")
	builder.PushType(PNPushTypeGCM)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should use v1 path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v1/push/")

	// Should have correct query type
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("gcm", query.Get("type"))
}

func TestListPushProvisionsMPNSConfiguration(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newListPushProvisionsRequestBuilder(pn)
	builder.DeviceIDForPush("mpns-device-token")
	builder.PushType(PNPushTypeMPNS)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should use v1 path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v1/push/")

	// Should have correct query type
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("mpns", query.Get("type"))
}

// Device ID Encoding Tests

func TestListPushProvisionsDeviceIDEncoding(t *testing.T) {
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
			opts := newListPushProvisionsRequestOpts(pn, pn.ctx)
			opts.DeviceIDForPush = tc.deviceID
			opts.PushType = PNPushTypeAPNS

			path, err := opts.buildPath()
			assert.Nil(err)
			assert.Equal(tc.expectedPath, path)
		})
	}
}

// Response Parsing Tests

func TestListPushProvisionsResponseParsing(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name             string
		jsonResponse     string
		expectError      bool
		expectedChannels []string
		description      string
	}{
		{
			name:             "Valid response with channels",
			jsonResponse:     `["channel1", "channel2", "channel3"]`,
			expectError:      false,
			expectedChannels: []string{"channel1", "channel2", "channel3"},
			description:      "Parse valid response with multiple channels",
		},
		{
			name:             "Valid response with single channel",
			jsonResponse:     `["single-channel"]`,
			expectError:      false,
			expectedChannels: []string{"single-channel"},
			description:      "Parse valid response with single channel",
		},
		{
			name:             "Valid response with no channels",
			jsonResponse:     `[]`,
			expectError:      false,
			expectedChannels: []string{},
			description:      "Parse valid response with empty channel list",
		},
		{
			name:             "Invalid JSON",
			jsonResponse:     `{invalid json}`,
			expectError:      true,
			expectedChannels: nil,
			description:      "Handle invalid JSON gracefully",
		},
		{
			name:             "Non-array JSON",
			jsonResponse:     `{"channels": ["channel1"]}`,
			expectError:      false,
			expectedChannels: nil,
			description:      "Handle non-array JSON (returns nil channels)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, _, err := newListPushProvisionsRequestResponse([]byte(tc.jsonResponse), StatusResponse{})

			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.NotNil(resp)
				assert.Equal(tc.expectedChannels, resp.Channels)
			}
		})
	}
}

func TestListPushProvisionsResponseParsingEdgeCases(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name         string
		jsonResponse string
		expectError  bool
		description  string
	}{
		{
			name:         "Very large response",
			jsonResponse: fmt.Sprintf(`[%s]`, strings.Repeat(`"channel", `, 999)+`"last-channel"`),
			expectError:  false,
			description:  "Handle very large channel list",
		},
		{
			name:         "Unicode channel names",
			jsonResponse: `["频道中文", "канал-русский", "チャンネル"]`,
			expectError:  false,
			description:  "Handle Unicode channel names",
		},
		{
			name:         "Channels with special characters",
			jsonResponse: `["channel@domain.com", "channel#with$symbols", "channel with spaces"]`,
			expectError:  false,
			description:  "Handle channels with special characters",
		},
		{
			name:         "Empty string channels",
			jsonResponse: `["", "valid-channel", ""]`,
			expectError:  false,
			description:  "Handle empty string channels",
		},
		{
			name:         "Null JSON",
			jsonResponse: `null`,
			expectError:  false,
			description:  "Handle null JSON",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, _, err := newListPushProvisionsRequestResponse([]byte(tc.jsonResponse), StatusResponse{})

			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.NotNil(resp)
			}
		})
	}
}

// Comprehensive Edge Case Tests

func TestListPushProvisionsWithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*listPushProvisionsRequestBuilder)
	}{
		{
			name: "Very long device ID",
			setupFn: func(builder *listPushProvisionsRequestBuilder) {
				longDeviceID := strings.Repeat("device", 250) // 1500 characters
				builder.DeviceIDForPush(longDeviceID)
				builder.PushType(PNPushTypeAPNS)
			},
		},
		{
			name: "Very long topic",
			setupFn: func(builder *listPushProvisionsRequestBuilder) {
				longTopic := "com.example." + strings.Repeat("app", 100)
				builder.DeviceIDForPush("device123")
				builder.PushType(PNPushTypeAPNS2)
				builder.Topic(longTopic)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *listPushProvisionsRequestBuilder) {
				largeQueryParam := make(map[string]string)
				for i := 0; i < 100; i++ {
					largeQueryParam[fmt.Sprintf("param_%d", i)] = strings.Repeat(fmt.Sprintf("value_%d_", i), 20)
				}
				builder.DeviceIDForPush("device123")
				builder.PushType(PNPushTypeGCM)
				builder.QueryParam(largeQueryParam)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newListPushProvisionsRequestBuilder(pn)
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

func TestListPushProvisionsSpecialCharacterHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialStrings := []string{
		"<script>alert('xss')</script>",
		"SELECT * FROM devices; DROP TABLE devices;",
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
			builder := newListPushProvisionsRequestBuilder(pn)
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

func TestListPushProvisionsParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		deviceID    string
		pushType    PNPushType
		description string
	}{
		{
			name:        "Empty device ID",
			deviceID:    "",
			pushType:    PNPushTypeAPNS,
			description: "Device ID with empty string",
		},
		{
			name:        "Single character device ID",
			deviceID:    "a",
			pushType:    PNPushTypeGCM,
			description: "Device ID with single character",
		},
		{
			name:        "Unicode-only device ID",
			deviceID:    "测试",
			pushType:    PNPushTypeMPNS,
			description: "Device ID with Unicode characters",
		},
		{
			name:        "Very long device ID",
			deviceID:    strings.Repeat("a", 1000),
			pushType:    PNPushTypeAPNS2,
			description: "Very long device ID string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newListPushProvisionsRequestBuilder(pn)
			builder.DeviceIDForPush(tc.deviceID)
			builder.PushType(tc.pushType)
			if tc.pushType == PNPushTypeAPNS2 {
				builder.Topic("com.example.app") // Required for APNS2
			}

			// Should build valid components
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

func TestListPushProvisionsExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newListPushProvisionsRequestBuilder(pn)
	builder.DeviceIDForPush("device123")
	builder.PushType(PNPushTypeAPNS)

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestListPushProvisionsBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newListPushProvisionsRequestBuilder(pn)

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
	builder.PushType(pushType).
		DeviceIDForPush(deviceID).
		Topic(topic).
		Environment(environment).
		QueryParam(queryParam).
		Transport(transport)

	// Verify all values are set correctly
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
