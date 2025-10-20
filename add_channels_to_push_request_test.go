package pubnub

import (
	"fmt"
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

	o := newAddPushNotificationsOnChannelsBuilderWithContext(pubnub, pubnub.ctx)
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

	assert.Equal("pubnub/validation: pubnub: Add Push From Channel: Missing Subscribe Key", opts.validate().Error())
}

// New comprehensive tests for missing coverage

func TestAddChannelsToPushValidateAPNS2Topic(t *testing.T) {
	assert := assert.New(t)

	// Test APNS2 without topic should fail
	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS2,
	})

	err := opts.validate()
	assert.Contains(err.Error(), "Missing Push Topic")

	// Test APNS2 with topic should pass
	opts2 := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS2,
		Topic:           "com.example.app",
	})

	err2 := opts2.validate()
	assert.Nil(err2)
}

func TestAddChannelsToPushAPNS2PathBuilding(t *testing.T) {
	assert := assert.New(t)

	// Test APNS2 uses different path
	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS2,
		Topic:           "com.example.app",
	})

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v2/push/sub-key/sub_key/devices-apns2/deviceId", path)

	// Test regular APNS uses v1 path
	opts2 := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeAPNS,
	})

	path2, err2 := opts2.buildPath()
	assert.Nil(err2)
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId", path2)
}

func TestAddChannelsToPushHttpMethod(t *testing.T) {
	assert := assert.New(t)

	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeGCM,
	})

	assert.Equal("GET", opts.httpMethod())
}

func TestAddChannelsToPushIsAuthRequired(t *testing.T) {
	assert := assert.New(t)

	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeGCM,
	})

	assert.True(opts.isAuthRequired())
}

func TestAddChannelsToPushTimeouts(t *testing.T) {
	assert := assert.New(t)

	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeGCM,
	})

	assert.Equal(pubnub.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pubnub.Config.ConnectTimeout, opts.connectTimeout())
}

func TestAddChannelsToPushOperationType(t *testing.T) {
	assert := assert.New(t)

	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1"},
		DeviceIDForPush: "deviceId",
		PushType:        PNPushTypeGCM,
	})

	assert.Equal(PNAddPushNotificationsOnChannelsOperation, opts.operationType())
}

func TestAddChannelsToPushBuilderSetters(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name       string
		setterFunc func(*addPushNotificationsOnChannelsBuilder)
		validateFn func(*addPushNotificationsOnChannelsBuilder, *testing.T)
	}{
		{
			name: "Channels setter",
			setterFunc: func(b *addPushNotificationsOnChannelsBuilder) {
				b.Channels([]string{"ch1", "ch2"})
			},
			validateFn: func(b *addPushNotificationsOnChannelsBuilder, t *testing.T) {
				assert.Equal([]string{"ch1", "ch2"}, b.opts.Channels)
			},
		},
		{
			name: "DeviceIDForPush setter",
			setterFunc: func(b *addPushNotificationsOnChannelsBuilder) {
				b.DeviceIDForPush("test-device-123")
			},
			validateFn: func(b *addPushNotificationsOnChannelsBuilder, t *testing.T) {
				assert.Equal("test-device-123", b.opts.DeviceIDForPush)
			},
		},
		{
			name: "Topic setter",
			setterFunc: func(b *addPushNotificationsOnChannelsBuilder) {
				b.Topic("com.example.topic")
			},
			validateFn: func(b *addPushNotificationsOnChannelsBuilder, t *testing.T) {
				assert.Equal("com.example.topic", b.opts.Topic)
			},
		},
		{
			name: "Environment setter",
			setterFunc: func(b *addPushNotificationsOnChannelsBuilder) {
				b.Environment(PNPushEnvironmentProduction)
			},
			validateFn: func(b *addPushNotificationsOnChannelsBuilder, t *testing.T) {
				assert.Equal(PNPushEnvironmentProduction, b.opts.Environment)
			},
		},
		{
			name: "QueryParam setter",
			setterFunc: func(b *addPushNotificationsOnChannelsBuilder) {
				b.QueryParam(map[string]string{"key": "value"})
			},
			validateFn: func(b *addPushNotificationsOnChannelsBuilder, t *testing.T) {
				assert.Equal(map[string]string{"key": "value"}, b.opts.QueryParam)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			builder := newAddPushNotificationsOnChannelsBuilder(pubnub)
			tc.setterFunc(builder)
			tc.validateFn(builder, t)
		})
	}
}

func TestAddChannelsToPushBuilderChaining(t *testing.T) {
	assert := assert.New(t)

	builder := newAddPushNotificationsOnChannelsBuilder(pubnub).
		Channels([]string{"ch1", "ch2"}).
		DeviceIDForPush("device123").
		PushType(PNPushTypeAPNS2).
		Topic("com.test.app").
		Environment(PNPushEnvironmentDevelopment).
		QueryParam(map[string]string{"test": "value"})

	assert.Equal([]string{"ch1", "ch2"}, builder.opts.Channels)
	assert.Equal("device123", builder.opts.DeviceIDForPush)
	assert.Equal(PNPushTypeAPNS2, builder.opts.PushType)
	assert.Equal("com.test.app", builder.opts.Topic)
	assert.Equal(PNPushEnvironmentDevelopment, builder.opts.Environment)
	assert.Equal(map[string]string{"test": "value"}, builder.opts.QueryParam)
}

func TestAddChannelsToPushValidationErrors(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		opts     addChannelsToPushOpts
		expected string
	}{
		{
			name: "Missing subscribe key",
			opts: addChannelsToPushOpts{
				Channels:        []string{"ch1"},
				DeviceIDForPush: "device123",
				PushType:        PNPushTypeGCM,
			},
			expected: "Missing Subscribe Key",
		},
		{
			name: "Missing device ID",
			opts: addChannelsToPushOpts{
				Channels: []string{"ch1"},
				PushType: PNPushTypeGCM,
			},
			expected: "Missing Device ID",
		},
		{
			name: "Missing channels",
			opts: addChannelsToPushOpts{
				DeviceIDForPush: "device123",
				PushType:        PNPushTypeGCM,
			},
			expected: "Missing Channel",
		},
		{
			name: "Missing push type",
			opts: addChannelsToPushOpts{
				Channels:        []string{"ch1"},
				DeviceIDForPush: "device123",
				PushType:        PNPushTypeNone,
			},
			expected: "Missing Push Type",
		},
		{
			name: "Missing topic for APNS2",
			opts: addChannelsToPushOpts{
				Channels:        []string{"ch1"},
				DeviceIDForPush: "device123",
				PushType:        PNPushTypeAPNS2,
			},
			expected: "Missing Push Topic",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pn := NewPubNub(NewDemoConfig())
			if tc.name == "Missing subscribe key" {
				pn.Config.SubscribeKey = ""
			}
			opts := newAddChannelsToPushOpts(pn, pn.ctx, tc.opts)
			err := opts.validate()
			assert.Contains(err.Error(), tc.expected)
		})
	}
}

func TestAddChannelsToPushAllPushTypes(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		pushType PNPushType
		expected string
	}{
		{
			name:     "GCM push type",
			pushType: PNPushTypeGCM,
			expected: "gcm",
		},
		{
			name:     "APNS push type",
			pushType: PNPushTypeAPNS,
			expected: "apns",
		},
		{
			name:     "APNS2 push type",
			pushType: PNPushTypeAPNS2,
			expected: "apns2",
		},
		{
			name:     "FCM push type",
			pushType: PNPushTypeFCM,
			expected: "fcm",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
				Channels:        []string{"ch1"},
				DeviceIDForPush: "device123",
				PushType:        tc.pushType,
				Topic:           "com.test.app", // Required for APNS2
			})

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expected, query.Get("type"))
		})
	}
}

func TestAddChannelsToPushSpecialCharacters(t *testing.T) {
	assert := assert.New(t)

	// Test device ID with special characters gets URL encoded
	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1"},
		DeviceIDForPush: "device@123#test",
		PushType:        PNPushTypeGCM,
	})

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "device%40123%23test")
}

func TestAddChannelsToPushUnicodeChannels(t *testing.T) {
	assert := assert.New(t)

	// Test Unicode channel names
	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"频道1", "канал2", "🚀channel"},
		DeviceIDForPush: "device123",
		PushType:        PNPushTypeGCM,
	})

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("频道1,канал2,🚀channel", query.Get("add"))
}

func TestAddChannelsToPushLargeChannelList(t *testing.T) {
	assert := assert.New(t)

	// Test large number of channels
	channels := make([]string, 100)
	for i := 0; i < 100; i++ {
		channels[i] = fmt.Sprintf("channel_%d", i)
	}

	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        channels,
		DeviceIDForPush: "device123",
		PushType:        PNPushTypeGCM,
	})

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Contains(query.Get("add"), "channel_0,channel_1")
	assert.Contains(query.Get("add"), "channel_99")
}

func TestAddChannelsToPushQueryParameterConflicts(t *testing.T) {
	assert := assert.New(t)

	// Test custom query params DO override system params (SetQueryParam is called last)
	opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
		Channels:        []string{"ch1"},
		DeviceIDForPush: "device123",
		PushType:        PNPushTypeGCM,
		QueryParam: map[string]string{
			"add":  "custom-override",
			"type": "custom-type",
		},
	})

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("custom-override", query.Get("add")) // Custom value should override
	assert.Equal("custom-type", query.Get("type"))    // Custom value should override
}

func TestAddChannelsToPushEnvironmentParameters(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
				Channels:        []string{"ch1"},
				DeviceIDForPush: "device123",
				PushType:        PNPushTypeAPNS2,
				Topic:           "com.test.app",
				Environment:     tc.environment,
			})

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expected, query.Get("environment"))
		})
	}
}

func TestAddChannelsToPushBuilderDefaults(t *testing.T) {
	assert := assert.New(t)

	builder := newAddPushNotificationsOnChannelsBuilder(pubnub)

	// Test defaults
	assert.Empty(builder.opts.Channels)
	assert.Empty(builder.opts.DeviceIDForPush)
	assert.Equal(PNPushType(0), builder.opts.PushType) // Zero value
	assert.Empty(builder.opts.Topic)
	assert.Equal(PNPushEnvironment(""), builder.opts.Environment) // Zero value
	assert.Nil(builder.opts.QueryParam)
}

func TestAddChannelsToPushEmptyAndNilValues(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name       string
		channels   []string
		deviceID   string
		shouldFail bool
		errorMsg   string
	}{
		{
			name:       "Empty channels slice",
			channels:   []string{},
			deviceID:   "device123",
			shouldFail: true,
			errorMsg:   "Missing Channel",
		},
		{
			name:       "Nil channels slice",
			channels:   nil,
			deviceID:   "device123",
			shouldFail: true,
			errorMsg:   "Missing Channel",
		},
		{
			name:       "Empty device ID",
			channels:   []string{"ch1"},
			deviceID:   "",
			shouldFail: true,
			errorMsg:   "Missing Device ID",
		},
		{
			name:       "Valid values",
			channels:   []string{"ch1"},
			deviceID:   "device123",
			shouldFail: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts := newAddChannelsToPushOpts(pubnub, pubnub.ctx, addChannelsToPushOpts{
				Channels:        tc.channels,
				DeviceIDForPush: tc.deviceID,
				PushType:        PNPushTypeGCM,
			})

			err := opts.validate()
			if tc.shouldFail {
				assert.NotNil(err)
				assert.Contains(err.Error(), tc.errorMsg)
			} else {
				assert.Nil(err)
			}
		})
	}
}
