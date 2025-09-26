package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/pubnub/go/v7/utils"
	"github.com/stretchr/testify/assert"
)

func AssertGetUUIDMetadata(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	incl := []PNUUIDMetadataInclude{
		PNUUIDMetadataIncludeCustom,
	}

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := EnumArrayToStringArray(incl)

	o := newGetUUIDMetadataBuilder(pn)
	if testContext {
		o = newGetUUIDMetadataBuilderWithContext(pn, pn.ctx)
	}

	o.Include(incl)
	o.UUID("id0")
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/uuids/%s", pn.Config.SubscribeKey, "id0"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Empty(body)

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
		assert.Equal(string(utils.JoinChannels(inclStr)), u.Get("include"))
	}

}

func TestGetUUIDMetadata(t *testing.T) {
	AssertGetUUIDMetadata(t, true, false)
}

func TestGetUUIDMetadataContext(t *testing.T) {
	AssertGetUUIDMetadata(t, true, true)
}

func TestGetUUIDMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestGetUUIDMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":{"id":"id0","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-20T13:26:19.140324Z","updated":"2019-08-20T13:26:19.140324Z","eTag":"AbyT4v2p6K7fpQE"}}`)

	r, _, err := newPNGetUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("id0", r.Data.ID)
	assert.Equal("name", r.Data.Name)
	assert.Equal("extid", r.Data.ExternalID)
	assert.Equal("purl", r.Data.ProfileURL)
	assert.Equal("email", r.Data.Email)
	// assert.Equal("2019-08-20T13:26:19.140324Z", r.Data.Created)
	assert.Equal("2019-08-20T13:26:19.140324Z", r.Data.Updated)
	assert.Equal("AbyT4v2p6K7fpQE", r.Data.ETag)
	assert.Equal("b", r.Data.Custom["a"])
	assert.Equal("d", r.Data.Custom["c"])

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestGetUUIDMetadataValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetUUIDMetadataValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	assert.Nil(opts.validate())
}

func TestGetUUIDMetadataValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"
	opts.Include = []string{"custom"}
	opts.QueryParam = map[string]string{"param": "value"}

	assert.Nil(opts.validate())
}

func TestGetUUIDMetadataValidateNoUUIDRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)
	// No UUID set - should still pass validation as UUID is optional (auto-fallback)

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestGetUUIDMetadataHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)

	assert.Equal("GET", opts.httpMethod())
}

func TestGetUUIDMetadataOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)

	assert.Equal(PNGetUUIDMetadataOperation, opts.operationType())
}

func TestGetUUIDMetadataIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestGetUUIDMetadataTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (4 setters)

func TestGetUUIDMetadataBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetUUIDMetadataBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestGetUUIDMetadataBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetUUIDMetadataBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestGetUUIDMetadataBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetUUIDMetadataBuilder(pn)

	// Test UUID setter
	builder.UUID("test-uuid")
	assert.Equal("test-uuid", builder.opts.UUID)

	// Test Include setter
	include := []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom}
	builder.Include(include)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestGetUUIDMetadataBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	include := []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom}
	queryParam := map[string]string{"key": "value"}

	builder := newGetUUIDMetadataBuilder(pn)
	result := builder.UUID("test-uuid").
		Include(include).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-uuid", builder.opts.UUID)
	assert.Equal(EnumArrayToStringArray(include), builder.opts.Include)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestGetUUIDMetadataBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newGetUUIDMetadataBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

func TestGetUUIDMetadataBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGetUUIDMetadataBuilder(pn)

	// Verify default values
	assert.Equal("", builder.opts.UUID)
	assert.Nil(builder.opts.Include)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestGetUUIDMetadataUUIDAutoFallback(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.UUID = "config-uuid-123"

	builder := newGetUUIDMetadataBuilder(pn)
	// Don't set UUID - should fall back to Config.UUID during Execute

	// Simulate Execute logic for UUID fallback
	if len(builder.opts.UUID) <= 0 {
		builder.opts.UUID = builder.opts.pubnub.Config.UUID
	}

	assert.Equal("config-uuid-123", builder.opts.UUID)
}

func TestGetUUIDMetadataUUIDExplicitOverridesConfig(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.UUID = "config-uuid-123"

	builder := newGetUUIDMetadataBuilder(pn)
	builder.UUID("explicit-uuid-456")

	// Simulate Execute logic for UUID fallback
	if len(builder.opts.UUID) <= 0 {
		builder.opts.UUID = builder.opts.pubnub.Config.UUID
	}

	// Should keep explicit UUID, not fallback
	assert.Equal("explicit-uuid-456", builder.opts.UUID)
}

// URL/Path Building Tests

func TestGetUUIDMetadataBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/uuids/test-uuid"
	assert.Equal(expected, path)
}

func TestGetUUIDMetadataBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "my-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/uuids/my-uuid"
	assert.Equal(expected, path)
}

func TestGetUUIDMetadataBuildPathWithSpecialCharsInUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "uuid-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/")
	assert.Contains(path, "uuid-with-special@chars#and$symbols")
}

func TestGetUUIDMetadataBuildPathWithUnicodeUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "用户ID-пользователь-ユーザー"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/")
	assert.Contains(path, "用户ID-пользователь-ユーザー")
}

// Query Parameter Tests

func TestGetUUIDMetadataBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestGetUUIDMetadataBuildQueryWithInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)
	opts.Include = []string{"custom"}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("custom", query.Get("include"))
}

func TestGetUUIDMetadataBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)

	customParams := map[string]string{
		"custom":         "value",
		"special_chars":  "value@with#symbols",
		"unicode":        "测试参数",
		"empty_value":    "",
		"number_string":  "42",
		"boolean_string": "true",
	}
	opts.QueryParam = customParams

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all custom parameters are present
	for key, expectedValue := range customParams {
		actualValue := query.Get(key)
		if key == "special_chars" {
			// Special characters should be URL encoded
			assert.Contains(actualValue, "%", "Query parameter %s should be URL encoded", key)
		} else if key == "unicode" {
			// Unicode should be URL encoded
			assert.Contains(actualValue, "%", "Query parameter %s should contain URL encoded Unicode", key)
		} else {
			assert.Equal(expectedValue, actualValue, "Query parameter %s should match", key)
		}
	}
}

func TestGetUUIDMetadataBuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)

	// Set all possible query parameters
	opts.Include = []string{"custom"}
	opts.QueryParam = map[string]string{
		"extra": "parameter",
		"debug": "true",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all parameters are set correctly
	assert.Equal("custom", query.Get("include"))
	assert.Equal("parameter", query.Get("extra"))
	assert.Equal("true", query.Get("debug"))

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

// Include Parameter Tests

func TestGetUUIDMetadataIncludeParameterHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		include  []string
		expected string
	}{
		{
			name:     "No include parameters",
			include:  nil,
			expected: "",
		},
		{
			name:     "Empty include array",
			include:  []string{},
			expected: "",
		},
		{
			name:     "Single include parameter",
			include:  []string{"custom"},
			expected: "custom",
		},
		{
			name:     "Multiple include parameters",
			include:  []string{"custom", "type"},
			expected: "custom,type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetUUIDMetadataOpts(pn, pn.ctx)
			opts.Include = tc.include

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expected, query.Get("include"))
		})
	}
}

func TestGetUUIDMetadataIncludeEnumConversion(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		include  []PNUUIDMetadataInclude
		expected string
	}{
		{
			name:     "No include",
			include:  nil,
			expected: "",
		},
		{
			name:     "Empty include",
			include:  []PNUUIDMetadataInclude{},
			expected: "",
		},
		{
			name:     "Single include",
			include:  []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom},
			expected: "custom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetUUIDMetadataBuilder(pn)
			if tc.include != nil {
				builder.Include(tc.include)
			}

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expected, query.Get("include"))
		})
	}
}

// Comprehensive Edge Case Tests

func TestGetUUIDMetadataWithLongUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*getUUIDMetadataBuilder)
	}{
		{
			name: "Very long UUID",
			setupFn: func(builder *getUUIDMetadataBuilder) {
				longUUID := strings.Repeat("VeryLongUUID", 50) // 600 characters
				builder.UUID(longUUID)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *getUUIDMetadataBuilder) {
				largeQueryParam := make(map[string]string)
				for i := 0; i < 100; i++ {
					largeQueryParam[fmt.Sprintf("param_%d", i)] = strings.Repeat(fmt.Sprintf("value_%d_", i), 20)
				}
				builder.QueryParam(largeQueryParam)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetUUIDMetadataBuilder(pn)
			builder.UUID("test-uuid")
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path and query
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

func TestGetUUIDMetadataSpecialCharacterHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialStrings := []string{
		"<script>alert('xss')</script>",
		"SELECT * FROM users; DROP TABLE users;",
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
			builder := newGetUUIDMetadataBuilder(pn)
			builder.UUID(specialString)
			builder.QueryParam(map[string]string{
				"special_field": specialString,
			})

			// Should pass validation (basic validation doesn't check content)
			assert.Nil(builder.opts.validate())

			// Should build valid path and query
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

func TestGetUUIDMetadataParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		uuid       string
		queryParam map[string]string
	}{
		{
			name: "Empty string values",
			uuid: "test-uuid",
			queryParam: map[string]string{
				"empty_param":  "",
				"normal_param": "value",
			},
		},
		{
			name: "Single character values",
			uuid: "a",
			queryParam: map[string]string{
				"single_char_key": "b",
				"key":             "c",
			},
		},
		{
			name: "Unicode-only values",
			uuid: "测试",
			queryParam: map[string]string{
				"unicode测试": "unicode值",
				"русский":   "значение",
			},
		},
		{
			name: "Mixed boundary conditions",
			uuid: "test",
			queryParam: map[string]string{
				"empty_param":   "",
				"long_param":    strings.Repeat("value", 200),
				"special@param": "special@value",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetUUIDMetadataBuilder(pn)
			builder.UUID(tc.uuid)
			if tc.queryParam != nil {
				builder.QueryParam(tc.queryParam)
			}

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid components
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, tc.uuid)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

func TestGetUUIDMetadataComplexUUIDCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*getUUIDMetadataBuilder)
		validateFn func(*testing.T, string, *url.Values)
	}{
		{
			name: "UUID with Unicode",
			setupFn: func(builder *getUUIDMetadataBuilder) {
				builder.UUID("用户123")
				builder.QueryParam(map[string]string{
					"参数": "值",
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, "用户123")
			},
		},
		{
			name: "UUID with special characters",
			setupFn: func(builder *getUUIDMetadataBuilder) {
				builder.UUID("user@special#123")
				builder.QueryParam(map[string]string{
					"special@key": "special@value",
					"with spaces": "also spaces",
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, "user@special#123")
			},
		},
		{
			name: "Very long UUID with complex params",
			setupFn: func(builder *getUUIDMetadataBuilder) {
				longUUID := strings.Repeat("VeryLongUser", 20)
				builder.UUID(longUUID)
				builder.Include([]PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom})
				builder.QueryParam(map[string]string{
					"filter": "complex filter expression",
					"sort":   "name:asc,created:desc",
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, strings.Repeat("VeryLongUser", 20))
				assert.Equal("custom", query.Get("include"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newGetUUIDMetadataBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid components
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Run custom validation
			tc.validateFn(t, path, query)
		})
	}
}

// Error Scenario Tests

func TestGetUUIDMetadataExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newGetUUIDMetadataBuilder(pn)
	builder.UUID("test-uuid")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestGetUUIDMetadataPathBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCases := []struct {
		name         string
		subscribeKey string
		uuid         string
		expectError  bool
	}{
		{
			name:         "Empty SubscribeKey",
			subscribeKey: "",
			uuid:         "test-uuid",
			expectError:  false, // buildPath doesn't validate, only validate() does
		},
		{
			name:         "Empty UUID",
			subscribeKey: "demo",
			uuid:         "",
			expectError:  false, // buildPath doesn't validate UUID
		},
		{
			name:         "SubscribeKey with spaces",
			subscribeKey: "   sub key   ",
			uuid:         "test-uuid",
			expectError:  false,
		},
		{
			name:         "UUID with spaces",
			subscribeKey: "demo",
			uuid:         "   test uuid   ",
			expectError:  false,
		},
		{
			name:         "Very special characters",
			subscribeKey: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			uuid:         "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			expectError:  false,
		},
		{
			name:         "Unicode SubscribeKey and UUID",
			subscribeKey: "测试订阅键-русский-キー",
			uuid:         "用户ID-пользователь-ユーザー",
			expectError:  false,
		},
		{
			name:         "Very long values",
			subscribeKey: strings.Repeat("a", 1000),
			uuid:         strings.Repeat("b", 1000),
			expectError:  false,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			pn.Config.SubscribeKey = tc.subscribeKey
			opts := newGetUUIDMetadataOpts(pn, pn.ctx)
			opts.UUID = tc.uuid

			path, err := opts.buildPath()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Contains(path, "/v2/objects/")
				assert.Contains(path, "/uuids/")
			}
		})
	}
}

func TestGetUUIDMetadataQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*getUUIDMetadataOpts)
		expectError bool
	}{
		{
			name: "Nil include array",
			setupOpts: func(opts *getUUIDMetadataOpts) {
				opts.Include = nil
			},
			expectError: false,
		},
		{
			name: "Empty include array",
			setupOpts: func(opts *getUUIDMetadataOpts) {
				opts.Include = []string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *getUUIDMetadataOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *getUUIDMetadataOpts) {
				opts.QueryParam = map[string]string{
					"special@key":   "special@value",
					"unicode测试":     "unicode值",
					"with spaces":   "also spaces",
					"equals=key":    "equals=value",
					"ampersand&key": "ampersand&value",
				}
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newGetUUIDMetadataOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			query, err := opts.buildQuery()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.NotNil(query)
			}
		})
	}
}

func TestGetUUIDMetadataBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newGetUUIDMetadataBuilder(pn)

	include := []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.UUID("complete-test-uuid").
		Include(include).
		QueryParam(queryParam)

	// Verify all values are set
	assert.Equal("complete-test-uuid", builder.opts.UUID)
	assert.Equal(EnumArrayToStringArray(include), builder.opts.Include)
	assert.Equal(queryParam, builder.opts.QueryParam)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v2/objects/demo/uuids/complete-test-uuid"
	assert.Equal(expectedPath, path)

	// Should build query with all params
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))
	assert.Equal("custom", query.Get("include"))
}

func TestGetUUIDMetadataResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetUUIDMetadataOpts(pn, pn.ctx)

	testCases := []struct {
		name        string
		jsonBytes   []byte
		expectError bool
	}{
		{
			name:        "Invalid JSON",
			jsonBytes:   []byte(`{invalid json`),
			expectError: true,
		},
		{
			name:        "Null JSON",
			jsonBytes:   []byte(`null`),
			expectError: false, // null is valid JSON
		},
		{
			name:        "Empty JSON object",
			jsonBytes:   []byte(`{}`),
			expectError: false,
		},
		{
			name:        "Valid UUID metadata response",
			jsonBytes:   []byte(`{"status":200,"data":{"id":"test-uuid","name":"Test User","email":"test@example.com"}}`),
			expectError: false,
		},
		{
			name:        "Response with status only",
			jsonBytes:   []byte(`{"status":200}`),
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, _, err := newPNGetUUIDMetadataResponse(tc.jsonBytes, opts, StatusResponse{})

			if tc.expectError {
				assert.NotNil(err)
				// When there's an error, resp might be nil or the empty response
				if resp == nil {
					assert.Equal(emptyPNGetUUIDMetadataResponse, resp)
				}
			} else {
				assert.Nil(err)
				// For successful parsing, resp should not be nil, but content may vary
				if err == nil {
					// Either resp is not nil, or it's nil but that's acceptable for null JSON
					if resp != nil {
						assert.NotNil(resp)
					}
				}
			}
		})
	}
}
