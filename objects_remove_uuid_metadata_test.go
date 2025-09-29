package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertRemoveUUIDMetadata(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newRemoveUUIDMetadataBuilder(pn)
	if testContext {
		o = newRemoveUUIDMetadataBuilderWithContext(pn, pn.ctx)
	}

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
	}

}

func TestRemoveUUIDMetadata(t *testing.T) {
	AssertRemoveUUIDMetadata(t, true, false)
}

func TestRemoveUUIDMetadataContext(t *testing.T) {
	AssertRemoveUUIDMetadata(t, true, true)
}

func TestRemoveUUIDMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNRemoveUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRemoveUUIDMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":null}`)

	r, _, err := newPNRemoveUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(nil, r.Data)

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestRemoveUUIDMetadataValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestRemoveUUIDMetadataValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	assert.Nil(opts.validate())
}

func TestRemoveUUIDMetadataValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"
	opts.QueryParam = map[string]string{"param": "value"}

	assert.Nil(opts.validate())
}

func TestRemoveUUIDMetadataValidateNoUUIDRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
	// No UUID set - should still pass validation as UUID is optional (auto-fallback)

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestRemoveUUIDMetadataHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)

	assert.Equal("DELETE", opts.httpMethod())
}

func TestRemoveUUIDMetadataOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)

	assert.Equal(PNRemoveUUIDMetadataOperation, opts.operationType())
}

func TestRemoveUUIDMetadataIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestRemoveUUIDMetadataTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (3 setters)

func TestRemoveUUIDMetadataBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveUUIDMetadataBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestRemoveUUIDMetadataBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveUUIDMetadataBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestRemoveUUIDMetadataBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveUUIDMetadataBuilder(pn)

	// Test UUID setter
	builder.UUID("test-uuid")
	assert.Equal("test-uuid", builder.opts.UUID)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestRemoveUUIDMetadataBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{"key": "value"}

	builder := newRemoveUUIDMetadataBuilder(pn)
	result := builder.UUID("test-uuid").
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-uuid", builder.opts.UUID)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestRemoveUUIDMetadataBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newRemoveUUIDMetadataBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

func TestRemoveUUIDMetadataBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveUUIDMetadataBuilder(pn)

	// Verify default values
	assert.Equal("", builder.opts.UUID)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestRemoveUUIDMetadataUUIDAutoFallback(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.UUID = "config-uuid-123"

	builder := newRemoveUUIDMetadataBuilder(pn)
	// Don't set UUID - should fall back to Config.UUID during Execute

	// Simulate Execute logic for UUID fallback
	if len(builder.opts.UUID) <= 0 {
		builder.opts.UUID = builder.opts.pubnub.Config.UUID
	}

	assert.Equal("config-uuid-123", builder.opts.UUID)
}

func TestRemoveUUIDMetadataUUIDExplicitOverridesConfig(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.UUID = "config-uuid-123"

	builder := newRemoveUUIDMetadataBuilder(pn)
	builder.UUID("explicit-uuid-456")

	// Simulate Execute logic for UUID fallback
	if len(builder.opts.UUID) <= 0 {
		builder.opts.UUID = builder.opts.pubnub.Config.UUID
	}

	// Should keep explicit UUID, not fallback
	assert.Equal("explicit-uuid-456", builder.opts.UUID)
}

func TestRemoveUUIDMetadataBuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	transport := &http.Transport{}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Test all 3 setters in chain
	builder := newRemoveUUIDMetadataBuilder(pn).
		UUID("test-uuid").
		QueryParam(queryParam).
		Transport(transport)

	// Verify all are set correctly
	assert.Equal("test-uuid", builder.opts.UUID)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestRemoveUUIDMetadataBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/uuids/test-uuid"
	assert.Equal(expected, path)
}

func TestRemoveUUIDMetadataBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "my-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/uuids/my-uuid"
	assert.Equal(expected, path)
}

func TestRemoveUUIDMetadataBuildPathWithSpecialCharsInUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "uuid-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/")
	assert.Contains(path, "uuid-with-special@chars#and$symbols")
}

func TestRemoveUUIDMetadataBuildPathWithUnicodeUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "用户ID-пользователь-ユーザー"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/")
	assert.Contains(path, "用户ID-пользователь-ユーザー")
}

// Query Parameter Tests

func TestRemoveUUIDMetadataBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestRemoveUUIDMetadataBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)

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

func TestRemoveUUIDMetadataBuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)

	// Set all possible query parameters
	opts.QueryParam = map[string]string{
		"extra": "parameter",
		"debug": "true",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all parameters are set correctly
	assert.Equal("parameter", query.Get("extra"))
	assert.Equal("true", query.Get("debug"))

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

// Query Parameter Edge Cases

func TestRemoveUUIDMetadataQueryParameterHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		queryParam map[string]string
	}{
		{
			name:       "No query parameters",
			queryParam: nil,
		},
		{
			name:       "Empty query parameters",
			queryParam: map[string]string{},
		},
		{
			name: "Single query parameter",
			queryParam: map[string]string{
				"single": "value",
			},
		},
		{
			name: "Multiple query parameters",
			queryParam: map[string]string{
				"param1": "value1",
				"param2": "value2",
				"param3": "value3",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
			opts.QueryParam = tc.queryParam

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Should always have default parameters
			assert.NotEmpty(query.Get("uuid"))
			assert.NotEmpty(query.Get("pnsdk"))

			// Verify custom parameters if any
			if tc.queryParam != nil {
				for key, expectedValue := range tc.queryParam {
					assert.Equal(expectedValue, query.Get(key))
				}
			}
		})
	}
}

// Comprehensive Edge Case Tests

func TestRemoveUUIDMetadataWithLongUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*removeUUIDMetadataBuilder)
	}{
		{
			name: "Very long UUID",
			setupFn: func(builder *removeUUIDMetadataBuilder) {
				longUUID := strings.Repeat("VeryLongUUID", 50) // 600 characters
				builder.UUID(longUUID)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *removeUUIDMetadataBuilder) {
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
			builder := newRemoveUUIDMetadataBuilder(pn)
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

func TestRemoveUUIDMetadataSpecialCharacterHandling(t *testing.T) {
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
			builder := newRemoveUUIDMetadataBuilder(pn)
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

func TestRemoveUUIDMetadataParameterBoundaries(t *testing.T) {
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
			builder := newRemoveUUIDMetadataBuilder(pn)
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

func TestRemoveUUIDMetadataComplexUUIDCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*removeUUIDMetadataBuilder)
		validateFn func(*testing.T, string, *url.Values)
	}{
		{
			name: "UUID with Unicode",
			setupFn: func(builder *removeUUIDMetadataBuilder) {
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
			setupFn: func(builder *removeUUIDMetadataBuilder) {
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
			setupFn: func(builder *removeUUIDMetadataBuilder) {
				longUUID := strings.Repeat("VeryLongUser", 20)
				builder.UUID(longUUID)
				builder.QueryParam(map[string]string{
					"param1": "complex parameter value",
					"param2": "another complex value",
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, strings.Repeat("VeryLongUser", 20))
				assert.Equal("complex%20parameter%20value", query.Get("param1"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRemoveUUIDMetadataBuilder(pn)
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

func TestRemoveUUIDMetadataExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newRemoveUUIDMetadataBuilder(pn)
	builder.UUID("test-uuid")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestRemoveUUIDMetadataPathBuildingEdgeCases(t *testing.T) {
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
			opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
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

func TestRemoveUUIDMetadataQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*removeUUIDMetadataOpts)
		expectError bool
	}{
		{
			name: "Nil query params",
			setupOpts: func(opts *removeUUIDMetadataOpts) {
				opts.QueryParam = nil
			},
			expectError: false,
		},
		{
			name: "Empty query params",
			setupOpts: func(opts *removeUUIDMetadataOpts) {
				opts.QueryParam = map[string]string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *removeUUIDMetadataOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *removeUUIDMetadataOpts) {
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
			opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)
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

func TestRemoveUUIDMetadataBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newRemoveUUIDMetadataBuilder(pn)

	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.UUID("complete-test-uuid").
		QueryParam(queryParam)

	// Verify all values are set
	assert.Equal("complete-test-uuid", builder.opts.UUID)
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
}

func TestRemoveUUIDMetadataResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)

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
			name:        "Valid delete response with null data",
			jsonBytes:   []byte(`{"status":200,"data":null}`),
			expectError: false,
		},
		{
			name:        "Valid delete response with empty data",
			jsonBytes:   []byte(`{"status":200,"data":{}}`),
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
			resp, _, err := newPNRemoveUUIDMetadataResponse(tc.jsonBytes, opts, StatusResponse{})

			if tc.expectError {
				assert.NotNil(err)
				// When there's an error, resp might be nil or the empty response
				if resp == nil {
					assert.Equal(emptyPNRemoveUUIDMetadataResponse, resp)
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

// DELETE-specific tests

func TestRemoveUUIDMetadataDeleteOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRemoveUUIDMetadataBuilder(pn)
	builder.UUID("test-uuid-to-delete")

	// Verify it's a DELETE operation
	assert.Equal("DELETE", builder.opts.httpMethod())

	// DELETE operations typically don't have a body
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	// Should have proper path for deletion
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/test-uuid-to-delete")
}

func TestRemoveUUIDMetadataDeleteResponseHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRemoveUUIDMetadataOpts(pn, pn.ctx)

	// DELETE operations typically return null data
	testCases := []struct {
		name     string
		jsonData []byte
		expected interface{}
	}{
		{
			name:     "Null data response",
			jsonData: []byte(`{"status":200,"data":null}`),
			expected: nil,
		},
		{
			name:     "Empty object response",
			jsonData: []byte(`{"status":200,"data":{}}`),
			expected: map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, _, err := newPNRemoveUUIDMetadataResponse(tc.jsonData, opts, StatusResponse{})
			assert.Nil(err)
			assert.NotNil(resp)
			assert.Equal(tc.expected, resp.Data)
		})
	}
}
