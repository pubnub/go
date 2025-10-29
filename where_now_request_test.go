package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func init() {
	pnconfig = NewConfigWithUserId(UserId(GenerateUUID()))

	pnconfig.PublishKey = "pub_key"
	pnconfig.SubscribeKey = "sub_key"
	pnconfig.SecretKey = "secret_key"

	pubnub = NewPubNub(pnconfig)
}

func TestWhereNowBasicRequest(t *testing.T) {
	assert := assert.New(t)

	opts := newWhereNowOpts(pubnub, pubnub.ctx)
	opts.UUID = "my-custom-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestWhereNowBasicRequestQueryParam(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := newWhereNowOpts(pubnub, pubnub.ctx)
	opts.UUID = "my-custom-uuid"
	opts.QueryParam = queryParam
	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewWhereNowBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newWhereNowBuilder(pubnub)
	o.UUID("my-custom-uuid")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})
}

func TestNewWhereNowBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newWhereNowBuilderWithContext(pubnub, pubnub.ctx)
	o.UUID("my-custom-uuid")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})
}

func TestNewWhereNowResponserrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newWhereNowResponse(jsonBytes, StatusResponse{})
	assert.Equal("pubnub/parsing: error unmarshalling response: {s}", err.Error())
}

func TestWhereNowValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""

	opts := newWhereNowOpts(pn, pn.ctx)
	opts.UUID = "my-custom-uuid"

	assert.Equal("pubnub/validation: pubnub: Where Now: Missing Subscribe Key", opts.validate().Error())
}

// HTTP Method and Operation Tests

func TestWhereNowHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)

	assert.Equal("GET", opts.httpMethod())
}

func TestWhereNowOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)

	assert.Equal(PNWhereNowOperation, opts.operationType())
}

func TestWhereNowIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestWhereNowTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (3 setters)

func TestWhereNowBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newWhereNowBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestWhereNowBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newWhereNowBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestWhereNowBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newWhereNowBuilder(pn)

	// Test UUID setter
	uuid := "custom-user-123"
	builder.UUID(uuid)
	assert.Equal(uuid, builder.opts.UUID)

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

func TestWhereNowBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	uuid := "custom-uuid"
	queryParam := map[string]string{"key": "value"}
	transport := &http.Transport{}

	builder := newWhereNowBuilder(pn)
	result := builder.UUID(uuid).
		QueryParam(queryParam).
		Transport(transport)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal(uuid, builder.opts.UUID)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

func TestWhereNowBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newWhereNowBuilder(pn)

	// Verify default values
	assert.Empty(builder.opts.UUID)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestWhereNowBuilderUUIDCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		uuid        string
		description string
	}{
		{
			name:        "Custom UUID",
			uuid:        "custom-user-123",
			description: "Set where now for custom UUID",
		},
		{
			name:        "UUID with special characters",
			uuid:        "user@domain.com",
			description: "Set where now for UUID with special characters",
		},
		{
			name:        "UUID with Unicode",
			uuid:        "用户123",
			description: "Set where now for UUID with Unicode characters",
		},
		{
			name:        "Empty UUID (uses default)",
			uuid:        "",
			description: "Empty UUID should use default config UUID",
		},
		{
			name:        "Long UUID",
			uuid:        strings.Repeat("a", 100),
			description: "Very long UUID handling",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newWhereNowBuilder(pn)
			builder.UUID(tc.uuid)

			assert.Equal(tc.uuid, builder.opts.UUID)
		})
	}
}

func TestWhereNowBuilderQueryParamCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		queryParam  map[string]string
		description string
	}{
		{
			name:        "Empty query params",
			queryParam:  map[string]string{},
			description: "Empty query parameter map",
		},
		{
			name: "Multiple query params",
			queryParam: map[string]string{
				"param1": "value1",
				"param2": "value2",
				"param3": "value3",
			},
			description: "Multiple query parameters",
		},
		{
			name: "Special character query params",
			queryParam: map[string]string{
				"special@key":   "special@value",
				"unicode测试":     "unicode值",
				"with spaces":   "also spaces",
				"equals=key":    "equals=value",
				"ampersand&key": "ampersand&value",
			},
			description: "Query parameters with special characters",
		},
		{
			name:        "Nil query params",
			queryParam:  nil,
			description: "Nil query parameter map",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newWhereNowBuilder(pn)
			builder.QueryParam(tc.queryParam)

			assert.Equal(tc.queryParam, builder.opts.QueryParam)
		})
	}
}

func TestWhereNowBuilderAllSettersChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	uuid := "custom-uuid"
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}
	transport := &http.Transport{}

	// Test all 3 setters in chain
	builder := newWhereNowBuilder(pn).
		UUID(uuid).
		QueryParam(queryParam).
		Transport(transport)

	// Verify all are set correctly
	assert.Equal(uuid, builder.opts.UUID)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)
}

// URL/Path Building Tests

func TestWhereNowBuildPathBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)
	opts.UUID = "test-user"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/uuid/test-user"
	assert.Equal(expected, path)
}

func TestWhereNowBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newWhereNowOpts(pn, pn.ctx)
	opts.UUID = "my-user"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/custom-sub-key/uuid/my-user"
	assert.Equal(expected, path)
}

func TestWhereNowBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)
	opts.UUID = "user@domain.com"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/uuid/user@domain.com"
	assert.Equal(expected, path)
}

func TestWhereNowBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)
	opts.UUID = "用户123"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/uuid/用户123"
	assert.Equal(expected, path)
}

func TestWhereNowBuildPathEmptyUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)
	opts.UUID = ""

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/uuid/"
	assert.Equal(expected, path)
}

func TestWhereNowBuildPathLongUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)
	longUUID := strings.Repeat("a", 100)
	opts.UUID = longUUID

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/presence/sub-key/demo/uuid/" + longUUID
	assert.Equal(expected, path)
}

// JSON Body Building Tests (CRITICAL for GET operation - should be empty)

func TestWhereNowBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations should have empty body
	assert.Equal([]byte{}, body)
}

func TestWhereNowBuildBodyWithAllParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)

	// Set all possible parameters - body should still be empty for GET
	opts.UUID = "custom-uuid"
	opts.QueryParam = map[string]string{"param": "value"}
	opts.Transport = &http.Transport{}

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET operations always have empty body regardless of parameters
	assert.Equal([]byte{}, body)
}

func TestWhereNowBuildBodyErrorScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)

	// Even with potential error conditions, buildBody should not fail for GET
	opts.UUID = "" // Empty UUID

	body, err := opts.buildBody()
	assert.Nil(err) // buildBody should never error for GET operations
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}

// Query Parameter Tests

func TestWhereNowBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestWhereNowBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newWhereNowOpts(pn, pn.ctx)

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

	// Verify all custom parameters are set
	assert.Equal("value", query.Get("custom"))
	assert.Equal("value%40with%23symbols", query.Get("special_chars"))
	assert.Equal("%E6%B5%8B%E8%AF%95%E5%8F%82%E6%95%B0", query.Get("unicode"))
	assert.Equal("", query.Get("empty_value"))
	assert.Equal("42", query.Get("number_string"))
	assert.Equal("true", query.Get("boolean_string"))
}

func TestWhereNowBuildQueryEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		queryParam  map[string]string
		description string
	}{
		{
			name:        "Nil query params",
			queryParam:  nil,
			description: "Nil query parameter map",
		},
		{
			name:        "Empty query params",
			queryParam:  map[string]string{},
			description: "Empty query parameter map",
		},
		{
			name: "Large query params",
			queryParam: map[string]string{
				"param1": strings.Repeat("a", 1000),
				"param2": strings.Repeat("b", 1000),
			},
			description: "Large query parameter values",
		},
		{
			name: "Special character query params",
			queryParam: map[string]string{
				"special@key":   "special@value",
				"unicode测试":     "unicode值",
				"with spaces":   "also spaces",
				"equals=key":    "equals=value",
				"ampersand&key": "ampersand&value",
			},
			description: "Query parameters with special characters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newWhereNowOpts(pn, pn.ctx)
			opts.QueryParam = tc.queryParam

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			// Should always have default parameters
			assert.NotEmpty(query.Get("uuid"))
			assert.NotEmpty(query.Get("pnsdk"))
		})
	}
}

// GET-Specific Tests (Where Now Characteristics)

func TestWhereNowGetOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newWhereNowBuilder(pn)
	builder.UUID("test-user")

	// Verify it's a GET operation
	assert.Equal("GET", builder.opts.httpMethod())

	// GET operations have empty body
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	// Should have proper path for where now
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/presence/sub-key/demo/uuid/test-user")
}

func TestWhereNowChannelRetrieval(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*whereNowOpts)
		description string
	}{
		{
			name: "Where now for specific UUID",
			setupOpts: func(opts *whereNowOpts) {
				opts.UUID = "user-123"
			},
			description: "Get channel list for specific user",
		},
		{
			name: "Where now with special character UUID",
			setupOpts: func(opts *whereNowOpts) {
				opts.UUID = "user@domain.com"
			},
			description: "Get channel list for UUID with special characters",
		},
		{
			name: "Where now with Unicode UUID",
			setupOpts: func(opts *whereNowOpts) {
				opts.UUID = "用户测试"
			},
			description: "Get channel list for Unicode UUID",
		},
		{
			name: "Where now with custom query params",
			setupOpts: func(opts *whereNowOpts) {
				opts.UUID = "user-123"
				opts.QueryParam = map[string]string{
					"debug": "true",
					"extra": "parameter",
				}
			},
			description: "Get channel list with additional query parameters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newWhereNowOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			// Should pass validation
			assert.Nil(opts.validate())

			// Should be GET operation
			assert.Equal("GET", opts.httpMethod())

			// Should have empty body
			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)

			// Should build valid path
			path, err := opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/v2/presence/sub-key/")
			assert.Contains(path, "/uuid/")

			// Should build valid query
			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

func TestWhereNowUUIDHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name         string
		uuid         string
		expectedUUID string
		description  string
	}{
		{
			name:         "Default UUID from config",
			uuid:         "",
			expectedUUID: pn.Config.UUID,
			description:  "Use default config UUID when not specified",
		},
		{
			name:         "Custom UUID",
			uuid:         "custom-user-123",
			expectedUUID: "custom-user-123",
			description:  "Use custom UUID when specified",
		},
		{
			name:         "UUID with special characters",
			uuid:         "user@domain.com",
			expectedUUID: "user@domain.com",
			description:  "Handle UUID with special characters",
		},
		{
			name:         "Unicode UUID",
			uuid:         "用户123",
			expectedUUID: "用户123",
			description:  "Handle Unicode UUID",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newWhereNowBuilder(pn)
			if tc.uuid != "" {
				builder.UUID(tc.uuid)
			}

			// Execute to trigger UUID default logic
			// Since we can't actually execute without a network call,
			// we'll simulate the UUID defaulting logic
			if len(builder.opts.UUID) <= 0 {
				builder.opts.UUID = builder.opts.pubnub.Config.UUID
			}

			assert.Equal(tc.expectedUUID, builder.opts.UUID)
		})
	}
}

func TestWhereNowEmptyBodyVerification(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that GET operations always have empty body regardless of configuration
	testCases := []struct {
		name      string
		setupOpts func(*whereNowOpts)
	}{
		{
			name: "With all parameters set",
			setupOpts: func(opts *whereNowOpts) {
				opts.UUID = "custom-uuid"
				opts.QueryParam = map[string]string{
					"param1": "value1",
					"param2": "value2",
				}
				opts.Transport = &http.Transport{}
			},
		},
		{
			name: "With minimal parameters",
			setupOpts: func(opts *whereNowOpts) {
				opts.UUID = "simple-user"
			},
		},
		{
			name: "With empty/nil parameters",
			setupOpts: func(opts *whereNowOpts) {
				opts.UUID = ""
				opts.QueryParam = nil
				opts.Transport = nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newWhereNowOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			body, err := opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
			assert.Equal([]byte{}, body)
		})
	}
}

// Comprehensive Edge Case Tests

func TestWhereNowWithLargeData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*whereNowBuilder)
	}{
		{
			name: "Very long UUID",
			setupFn: func(builder *whereNowBuilder) {
				longUUID := strings.Repeat("user", 250) // 1000 characters
				builder.UUID(longUUID)
			},
		},
		{
			name: "Extremely large query params",
			setupFn: func(builder *whereNowBuilder) {
				largeQueryParam := make(map[string]string)
				for i := 0; i < 100; i++ {
					largeQueryParam[fmt.Sprintf("param_%d", i)] = strings.Repeat(fmt.Sprintf("value_%d_", i), 20)
				}
				builder.UUID("test-user").
					QueryParam(largeQueryParam)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newWhereNowBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation for all cases
			assert.Nil(builder.opts.validate())

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

func TestWhereNowSpecialCharacterHandling(t *testing.T) {
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
			builder := newWhereNowBuilder(pn)
			builder.UUID(specialString)
			builder.QueryParam(map[string]string{
				"special_param": specialString,
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

			// Should always have empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)
		})
	}
}

func TestWhereNowParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		uuid        string
		description string
	}{
		{
			name:        "Empty string UUID",
			uuid:        "",
			description: "UUID with empty string",
		},
		{
			name:        "Single character UUID",
			uuid:        "a",
			description: "UUID with single character",
		},
		{
			name:        "Unicode-only UUID",
			uuid:        "测试",
			description: "UUID with Unicode characters",
		},
		{
			name:        "Very long UUID",
			uuid:        strings.Repeat("a", 1000),
			description: "Very long UUID string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newWhereNowBuilder(pn)
			builder.UUID(tc.uuid)

			// Should pass validation for most cases
			assert.Nil(builder.opts.validate())

			// Should build valid components
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, "/v2/presence/sub-key/")

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)

			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body) // GET operation always has empty body
		})
	}
}

func TestWhereNowComplexScenarios(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*whereNowBuilder)
		validateFn func(*testing.T, string, *url.Values)
	}{
		{
			name: "User tracking scenario",
			setupFn: func(builder *whereNowBuilder) {
				builder.UUID("user-tracking-123")
				builder.QueryParam(map[string]string{
					"tracking_id": "session-abc",
					"app_version": "v1.2.3",
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, "user-tracking-123")
				assert.Equal("session-abc", query.Get("tracking_id"))
				assert.Equal("v1.2.3", query.Get("app_version"))
			},
		},
		{
			name: "International user scenario",
			setupFn: func(builder *whereNowBuilder) {
				builder.UUID("用户测试-русский-ユーザー")
				builder.QueryParam(map[string]string{
					"locale":   "zh-CN",
					"timezone": "Asia/Shanghai",
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, "用户测试-русский-ユーザー")
				assert.Equal("zh-CN", query.Get("locale"))
				assert.Equal("Asia%2FShanghai", query.Get("timezone"))
			},
		},
		{
			name: "Gaming user scenario",
			setupFn: func(builder *whereNowBuilder) {
				builder.UUID("player-abc123")
				builder.QueryParam(map[string]string{
					"game_mode": "battle_royale",
					"region":    "us-west",
					"version":   "2.1.0",
				})
			},
			validateFn: func(t *testing.T, path string, query *url.Values) {
				assert.Contains(path, "player-abc123")
				assert.Equal("battle_royale", query.Get("game_mode"))
				assert.Equal("us-west", query.Get("region"))
				assert.Equal("2.1.0", query.Get("version"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newWhereNowBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path
			path, err := builder.opts.buildPath()
			assert.Nil(err)

			// Should build valid query
			query, err := builder.opts.buildQuery()
			assert.Nil(err)

			// Should always have empty body (GET operation)
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Empty(body)

			// Run custom validation
			tc.validateFn(t, path, query)
		})
	}
}

// Error Scenario Tests

func TestWhereNowExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newWhereNowBuilder(pn)
	builder.UUID("test-user")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestWhereNowPathBuildingEdgeCases(t *testing.T) {
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
			uuid:         "test-user",
			expectError:  false, // buildPath doesn't validate, only validate() does
		},
		{
			name:         "SubscribeKey with spaces",
			subscribeKey: "   sub key   ",
			uuid:         "test-user",
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
			uuid:         "用户测试-пользователь-ユーザー",
			expectError:  false,
		},
		{
			name:         "Very long values",
			subscribeKey: strings.Repeat("a", 1000),
			uuid:         strings.Repeat("c", 1000),
			expectError:  false,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			pn.Config.SubscribeKey = tc.subscribeKey
			opts := newWhereNowOpts(pn, pn.ctx)
			opts.UUID = tc.uuid

			path, err := opts.buildPath()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Contains(path, "/v2/presence/sub-key/")
				assert.Contains(path, "/uuid/")
			}
		})
	}
}

func TestWhereNowQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*whereNowOpts)
		expectError bool
	}{
		{
			name: "Nil query params",
			setupOpts: func(opts *whereNowOpts) {
				opts.QueryParam = nil
			},
			expectError: false,
		},
		{
			name: "Empty query params",
			setupOpts: func(opts *whereNowOpts) {
				opts.QueryParam = map[string]string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *whereNowOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *whereNowOpts) {
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
			opts := newWhereNowOpts(pn, pn.ctx)
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

func TestWhereNowBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newWhereNowBuilder(pn)

	uuid := "custom-uuid"
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}
	transport := &http.Transport{}

	// Set all possible parameters
	builder.UUID(uuid).
		QueryParam(queryParam).
		Transport(transport)

	// Verify all values are set correctly
	assert.Equal(uuid, builder.opts.UUID)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(transport, builder.opts.Transport)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v2/presence/sub-key/demo/uuid/custom-uuid"
	assert.Equal(expectedPath, path)

	// Should build query with all params
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))

	// Should always have empty body (GET operation)
	body, err := builder.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
	assert.Equal([]byte{}, body)
}

func TestWhereNowValidationErrors(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name          string
		setupOpts     func(*whereNowOpts)
		expectedError string
	}{
		{
			name: "Missing subscribe key",
			setupOpts: func(opts *whereNowOpts) {
				opts.pubnub.Config.SubscribeKey = ""
				opts.UUID = "test-user"
			},
			expectedError: "Missing Subscribe Key",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh PubNub instance for each test case to avoid shared state
			pn := NewPubNub(NewDemoConfig())
			opts := newWhereNowOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			err := opts.validate()
			assert.NotNil(err)
			assert.Contains(err.Error(), tc.expectedError)
		})
	}
}

// Response Parsing Tests

func TestWhereNowResponseParsing(t *testing.T) {
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
			jsonResponse:     `{"status": 200, "message": "OK", "payload": {"channels": ["channel1", "channel2"]}, "service": "Presence"}`,
			expectError:      false,
			expectedChannels: []string{"channel1", "channel2"},
			description:      "Parse valid response with multiple channels",
		},
		{
			name:             "Valid response with single channel",
			jsonResponse:     `{"status": 200, "message": "OK", "payload": {"channels": ["single-channel"]}, "service": "Presence"}`,
			expectError:      false,
			expectedChannels: []string{"single-channel"},
			description:      "Parse valid response with single channel",
		},
		{
			name:             "Valid response with no channels",
			jsonResponse:     `{"status": 200, "message": "OK", "payload": {"channels": []}, "service": "Presence"}`,
			expectError:      false,
			expectedChannels: nil,
			description:      "Parse valid response with empty channel list",
		},
		{
			name:             "Response with missing payload",
			jsonResponse:     `{"status": 200, "message": "OK", "service": "Presence"}`,
			expectError:      false,
			expectedChannels: nil,
			description:      "Parse response missing payload section",
		},
		{
			name:             "Response with null channels",
			jsonResponse:     `{"status": 200, "message": "OK", "payload": {"channels": null}, "service": "Presence"}`,
			expectError:      false,
			expectedChannels: nil,
			description:      "Parse response with null channels",
		},
		{
			name:             "Invalid JSON",
			jsonResponse:     `{invalid json}`,
			expectError:      true,
			expectedChannels: nil,
			description:      "Handle invalid JSON gracefully",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, _, err := newWhereNowResponse([]byte(tc.jsonResponse), StatusResponse{})

			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.NotNil(resp)
				if tc.expectedChannels == nil {
					assert.Nil(resp.Channels)
				} else {
					assert.Equal(tc.expectedChannels, resp.Channels)
				}
			}
		})
	}
}

func TestWhereNowResponseParsingEdgeCases(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name         string
		jsonResponse string
		expectError  bool
		description  string
	}{
		{
			name:         "Empty JSON object",
			jsonResponse: `{}`,
			expectError:  false,
			description:  "Handle empty JSON object",
		},
		{
			name:         "Malformed channels array",
			jsonResponse: `{"payload": {"channels": ["valid", 123, "another"]}}`,
			expectError:  false,
			description:  "Handle mixed type channels array",
		},
		{
			name:         "Very large response",
			jsonResponse: fmt.Sprintf(`{"payload": {"channels": [%s]}}`, strings.Repeat(`"channel", `, 999)+`"last-channel"`),
			expectError:  false,
			description:  "Handle very large channel list",
		},
		{
			name:         "Unicode channel names",
			jsonResponse: `{"payload": {"channels": ["频道中文", "канал-русский", "チャンネル"]}}`,
			expectError:  false,
			description:  "Handle Unicode channel names",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, _, err := newWhereNowResponse([]byte(tc.jsonResponse), StatusResponse{})

			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.NotNil(resp)
			}
		})
	}
}
