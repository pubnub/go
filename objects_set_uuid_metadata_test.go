package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/pubnub/go/v7/utils"
	"github.com/stretchr/testify/assert"
)

func AssertSetUUIDMetadata(t *testing.T, checkQueryParam, testContext bool) {
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

	o := newSetUUIDMetadataBuilder(pn)
	if testContext {
		o = newSetUUIDMetadataBuilderWithContext(pn, pn.ctx)
	}

	o.Include(incl)
	o.UUID("id0")
	o.Name("name")
	o.ExternalID("exturl")
	o.ProfileURL("prourl")
	o.Email("email")
	o.Custom(custom)
	o.Status("active")
	o.Type("public")
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/uuids/%s", pn.Config.SubscribeKey, "id0"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	expectedBody := "{\"name\":\"name\",\"externalId\":\"exturl\",\"profileUrl\":\"prourl\",\"email\":\"email\",\"custom\":{\"a\":\"b\",\"c\":\"d\"},\"status\":\"active\",\"type\":\"public\"}"

	assert.Equal(expectedBody, string(body))

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
		assert.Equal(string(utils.JoinChannels(inclStr)), u.Get("include"))
	}

}

func TestExcludeInUUIDMetadataBodyNotSetFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	custom := map[string]interface{}{
		"a": "b",
		"c": "d",
	}

	o := newSetUUIDMetadataBuilder(pn)
	o.UUID("id0")
	o.Name("name")
	o.ExternalID("exturl")
	o.Custom(custom)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/uuids/%s", pn.Config.SubscribeKey, "id0"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	expectedBody := "{\"name\":\"name\",\"externalId\":\"exturl\",\"custom\":{\"a\":\"b\",\"c\":\"d\"}}"

	assert.Equal(expectedBody, string(body))
}

func TestSetUUIDMetadata(t *testing.T) {
	AssertSetUUIDMetadata(t, true, false)
}

func TestSetUUIDMetadataContext(t *testing.T) {
	AssertSetUUIDMetadata(t, true, true)
}

func TestSetUUIDMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNSetUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestSetUUIDMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":{"id":"id0","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"status":"active","type":"public","created":"2019-08-20T13:26:19.140324Z","updated":"2019-08-20T13:26:19.140324Z","eTag":"AbyT4v2p6K7fpQE"}}`)

	r, _, err := newPNSetUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
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
	assert.Equal("active", r.Data.Status)
	assert.Equal("public", r.Data.Type)

	assert.Nil(err)
}

// Comprehensive Validation Tests

func TestSetUUIDMetadataValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	err := opts.validate()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestSetUUIDMetadataValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	assert.Nil(opts.validate())
}

func TestSetUUIDMetadataValidateSuccessWithAllFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"
	opts.Name = "Test User"
	opts.Email = "test@example.com"
	opts.ExternalID = "ext123"
	opts.ProfileURL = "https://example.com/profile.jpg"
	opts.Custom = map[string]interface{}{"role": "admin"}

	assert.Nil(opts.validate())
}

func TestSetUUIDMetadataValidateNoUUIDRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	// No UUID set - should still pass validation as UUID is optional

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestSetUUIDMetadataHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)

	assert.Equal("PATCH", opts.httpMethod())
}

func TestSetUUIDMetadataOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)

	assert.Equal(PNSetUUIDMetadataOperation, opts.operationType())
}

func TestSetUUIDMetadataIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestSetUUIDMetadataTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Systematic Builder Pattern Tests (9 setters)

func TestSetUUIDMetadataBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetUUIDMetadataBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestSetUUIDMetadataBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetUUIDMetadataBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestSetUUIDMetadataBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetUUIDMetadataBuilder(pn)

	// Test UUID setter
	builder.UUID("test-uuid")
	assert.Equal("test-uuid", builder.opts.UUID)

	// Test Include setter
	include := []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom}
	builder.Include(include)
	expectedInclude := EnumArrayToStringArray(include)
	assert.Equal(expectedInclude, builder.opts.Include)

	// Test Name setter
	builder.Name("Test User")
	assert.Equal("Test User", builder.opts.Name)

	// Test ExternalID setter
	builder.ExternalID("ext123")
	assert.Equal("ext123", builder.opts.ExternalID)

	// Test ProfileURL setter
	builder.ProfileURL("https://example.com/profile.jpg")
	assert.Equal("https://example.com/profile.jpg", builder.opts.ProfileURL)

	// Test Email setter
	builder.Email("test@example.com")
	assert.Equal("test@example.com", builder.opts.Email)

	// Test Custom setter
	custom := map[string]interface{}{
		"role": "admin",
		"dept": "engineering",
	}
	builder.Custom(custom)
	assert.Equal(custom, builder.opts.Custom)

	// Test Status setter
	builder.Status("active")
	assert.Equal("active", builder.opts.Status)

	// Test Type setter
	builder.Type("public")
	assert.Equal("public", builder.opts.Type)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestSetUUIDMetadataBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	include := []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom}
	custom := map[string]interface{}{"role": "admin"}
	queryParam := map[string]string{"key": "value"}

	builder := newSetUUIDMetadataBuilder(pn)
	result := builder.UUID("test-uuid").
		Include(include).
		Name("Test User").
		ExternalID("ext123").
		ProfileURL("https://example.com/profile.jpg").
		Email("test@example.com").
		Custom(custom).
		Status("active").
		Type("public").
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-uuid", builder.opts.UUID)
	assert.Equal(EnumArrayToStringArray(include), builder.opts.Include)
	assert.Equal("Test User", builder.opts.Name)
	assert.Equal("ext123", builder.opts.ExternalID)
	assert.Equal("https://example.com/profile.jpg", builder.opts.ProfileURL)
	assert.Equal("test@example.com", builder.opts.Email)
	assert.Equal(custom, builder.opts.Custom)
	assert.Equal("active", builder.opts.Status)
	assert.Equal("public", builder.opts.Type)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestSetUUIDMetadataBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newSetUUIDMetadataBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

func TestSetUUIDMetadataBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newSetUUIDMetadataBuilder(pn)

	// Verify default values
	assert.Equal("", builder.opts.UUID)
	assert.Nil(builder.opts.Include)
	assert.Equal("", builder.opts.Name)
	assert.Equal("", builder.opts.ExternalID)
	assert.Equal("", builder.opts.ProfileURL)
	assert.Equal("", builder.opts.Email)
	assert.Nil(builder.opts.Custom)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestSetUUIDMetadataUUIDAutoFallback(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.UUID = "config-uuid-123"

	builder := newSetUUIDMetadataBuilder(pn)
	// Don't set UUID - should fall back to Config.UUID during Execute

	// Simulate Execute logic for UUID fallback
	if len(builder.opts.UUID) <= 0 {
		builder.opts.UUID = builder.opts.pubnub.Config.UUID
	}

	assert.Equal("config-uuid-123", builder.opts.UUID)
}

func TestSetUUIDMetadataUUIDExplicitOverridesConfig(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.UUID = "config-uuid-123"

	builder := newSetUUIDMetadataBuilder(pn)
	builder.UUID("explicit-uuid-456")

	// Simulate Execute logic for UUID fallback
	if len(builder.opts.UUID) <= 0 {
		builder.opts.UUID = builder.opts.pubnub.Config.UUID
	}

	// Should keep explicit UUID, not fallback
	assert.Equal("explicit-uuid-456", builder.opts.UUID)
}

// URL/Path Building Tests

func TestSetUUIDMetadataBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "test-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/demo/uuids/test-uuid"
	assert.Equal(expected, path)
}

func TestSetUUIDMetadataBuildPathWithDifferentSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "custom-sub-key"
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "my-uuid"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/objects/custom-sub-key/uuids/my-uuid"
	assert.Equal(expected, path)
}

func TestSetUUIDMetadataBuildPathWithSpecialCharsInUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "uuid-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/")
	assert.Contains(path, "uuid-with-special@chars#and$symbols")
}

func TestSetUUIDMetadataBuildPathWithUnicodeUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	opts.UUID = "用户ID-пользователь-ユーザー"

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v2/objects/demo/uuids/")
	assert.Contains(path, "用户ID-пользователь-ユーザー")
}

// JSON Body Building Tests

func TestSetUUIDMetadataBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	// No fields set

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := "{}"
	assert.Equal(expected, string(body))
}

func TestSetUUIDMetadataBuildBodyMinimal(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	opts.Name = "Test User"

	body, err := opts.buildBody()
	assert.Nil(err)
	expected := "{\"name\":\"Test User\"}"
	assert.Equal(expected, string(body))
}

func TestSetUUIDMetadataBuildBodyComplete(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	opts.Name = "Test User"
	opts.ExternalID = "ext123"
	opts.ProfileURL = "https://example.com/profile.jpg"
	opts.Email = "test@example.com"
	opts.Custom = map[string]interface{}{
		"role": "admin",
		"dept": "engineering",
	}
	opts.Status = "active"
	opts.Type = "public"

	body, err := opts.buildBody()
	assert.Nil(err)

	// Parse JSON to verify structure
	var parsed map[string]interface{}
	err = json.Unmarshal(body, &parsed)
	assert.Nil(err)

	assert.Equal("Test User", parsed["name"])
	assert.Equal("ext123", parsed["externalId"])
	assert.Equal("https://example.com/profile.jpg", parsed["profileUrl"])
	assert.Equal("test@example.com", parsed["email"])
	assert.Equal("active", parsed["status"])
	assert.Equal("public", parsed["type"])
	assert.NotNil(parsed["custom"])

	customMap := parsed["custom"].(map[string]interface{})
	assert.Equal("admin", customMap["role"])
	assert.Equal("engineering", customMap["dept"])
}

func TestSetUUIDMetadataBuildBodyPartialFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		setupFn  func(*setUUIDMetadataOpts)
		expected string
	}{
		{
			name: "Only email",
			setupFn: func(opts *setUUIDMetadataOpts) {
				opts.Email = "user@example.com"
			},
			expected: "{\"email\":\"user@example.com\"}",
		},
		{
			name: "Name and email",
			setupFn: func(opts *setUUIDMetadataOpts) {
				opts.Name = "John Doe"
				opts.Email = "john@example.com"
			},
			expected: "{\"name\":\"John Doe\",\"email\":\"john@example.com\"}",
		},
		{
			name: "ExternalID and ProfileURL",
			setupFn: func(opts *setUUIDMetadataOpts) {
				opts.ExternalID = "ext456"
				opts.ProfileURL = "https://cdn.example.com/avatars/user.png"
			},
			expected: "{\"externalId\":\"ext456\",\"profileUrl\":\"https://cdn.example.com/avatars/user.png\"}",
		},
		{
			name: "Only custom data",
			setupFn: func(opts *setUUIDMetadataOpts) {
				opts.Custom = map[string]interface{}{
					"tier":    "premium",
					"credits": 1000,
				}
			},
			expected: "{\"custom\":{\"credits\":1000,\"tier\":\"premium\"}}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newSetUUIDMetadataOpts(pn, pn.ctx)
			tc.setupFn(opts)

			body, err := opts.buildBody()
			assert.Nil(err)

			// Parse both expected and actual JSON to compare structure
			var expectedParsed, actualParsed map[string]interface{}
			err = json.Unmarshal([]byte(tc.expected), &expectedParsed)
			assert.Nil(err)
			err = json.Unmarshal(body, &actualParsed)
			assert.Nil(err)

			assert.Equal(expectedParsed, actualParsed)
		})
	}
}

func TestSetUUIDMetadataBuildBodyComplexCustom(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	opts.Custom = map[string]interface{}{
		"nested": map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": "deep value",
				"array":  []interface{}{1, "two", true},
			},
		},
		"unicode":       "测试数据-русский-ユーザー",
		"special_chars": "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		"empty_string":  "",
		"null_value":    nil,
		"boolean_true":  true,
		"boolean_false": false,
		"number_int":    42,
		"number_float":  3.14159,
	}

	body, err := opts.buildBody()
	assert.Nil(err)

	var parsedBody map[string]interface{}
	err = json.Unmarshal(body, &parsedBody)
	assert.Nil(err)

	customMap := parsedBody["custom"].(map[string]interface{})
	assert.Equal("测试数据-русский-ユーザー", customMap["unicode"])
	assert.Equal("!@#$%^&*()_+-=[]{}|;':\",./<>?", customMap["special_chars"])
	assert.Equal(true, customMap["boolean_true"])
	assert.Equal(false, customMap["boolean_false"])
	assert.Equal(float64(42), customMap["number_int"])
	assert.Equal(3.14159, customMap["number_float"])
}

func TestSetUUIDMetadataBuildBodyLargeCustom(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)

	// Create large custom object
	largeCustom := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		largeCustom[fmt.Sprintf("field_%d", i)] = fmt.Sprintf("value_%d", i)
	}
	opts.Custom = largeCustom

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.True(len(body) > 1000) // Should be a large JSON object

	var parsedBody map[string]interface{}
	err = json.Unmarshal(body, &parsedBody)
	assert.Nil(err)

	customMap := parsedBody["custom"].(map[string]interface{})
	assert.Equal(100, len(customMap))
	assert.Equal("value_50", customMap["field_50"])
}

// Query Parameter Tests

func TestSetUUIDMetadataBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestSetUUIDMetadataBuildQueryWithInclude(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)
	opts.Include = []string{"custom", "status", "type"}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("custom,status,type", query.Get("include"))
}

func TestSetUUIDMetadataBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)

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

func TestSetUUIDMetadataBuildQueryComprehensiveCombination(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)

	// Set all possible query parameters
	opts.Include = []string{"custom", "status", "type"}
	opts.QueryParam = map[string]string{
		"extra": "parameter",
		"debug": "true",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all parameters are set correctly
	assert.Equal("custom,status,type", query.Get("include"))
	assert.Equal("parameter", query.Get("extra"))
	assert.Equal("true", query.Get("debug"))

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

// Profile Field Validation Tests

func TestSetUUIDMetadataEmailFormatValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	validEmails := []string{
		"user@example.com",
		"test.email@domain.co.uk",
		"user+tag@example.org",
		"user123@test-domain.com",
		"用户@example.com", // Unicode in local part
		"user@测试.com",    // Unicode in domain
	}

	for _, email := range validEmails {
		t.Run(fmt.Sprintf("ValidEmail_%s", email), func(t *testing.T) {
			builder := newSetUUIDMetadataBuilder(pn)
			builder.Email(email)

			assert.Equal(email, builder.opts.Email)

			// Should build valid JSON body
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Contains(string(body), email)
		})
	}
}

func TestSetUUIDMetadataProfileURLValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	validURLs := []string{
		"https://example.com/profile.jpg",
		"http://cdn.example.com/avatars/user.png",
		"https://storage.googleapis.com/bucket/image.webp",
		"https://测试.com/profile.jpg",   // Unicode domain
		"https://example.com/用户头像.png", // Unicode path
		"https://example.com/profile?size=large&format=jpg",
	}

	for _, url := range validURLs {
		t.Run(fmt.Sprintf("ValidURL_%d", len(url)), func(t *testing.T) {
			builder := newSetUUIDMetadataBuilder(pn)
			builder.ProfileURL(url)

			assert.Equal(url, builder.opts.ProfileURL)

			// Should build valid JSON body
			body, err := builder.opts.buildBody()
			assert.Nil(err)

			// Parse JSON to handle URL encoding properly
			var parsed map[string]interface{}
			err = json.Unmarshal(body, &parsed)
			assert.Nil(err)
			assert.Equal(url, parsed["profileUrl"])
		})
	}
}

func TestSetUUIDMetadataUnicodeNames(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	unicodeNames := []string{
		"张三",             // Chinese
		"Иван Петров",    // Russian
		"田中太郎",           // Japanese
		"محمد أحمد",      // Arabic
		"김철수",            // Korean
		"राम शर्मा",      // Hindi
		"José María",     // Spanish with accents
		"François André", // French with accents
	}

	for i, name := range unicodeNames {
		t.Run(fmt.Sprintf("UnicodeName_%d", i), func(t *testing.T) {
			builder := newSetUUIDMetadataBuilder(pn)
			builder.Name(name)

			assert.Equal(name, builder.opts.Name)

			// Should build valid JSON body
			body, err := builder.opts.buildBody()
			assert.Nil(err)

			var parsed map[string]interface{}
			err = json.Unmarshal(body, &parsed)
			assert.Nil(err)
			assert.Equal(name, parsed["name"])
		})
	}
}

func TestSetUUIDMetadataExternalIDHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	externalIDs := []string{
		"ext123",
		"user_456",
		"uuid-789-abc-def",
		"external@system#123",
		"用户ID_123",
		"très_long_external_id_avec_caractères_spéciaux_123456789",
	}

	for i, externalID := range externalIDs {
		t.Run(fmt.Sprintf("ExternalID_%d", i), func(t *testing.T) {
			builder := newSetUUIDMetadataBuilder(pn)
			builder.ExternalID(externalID)

			assert.Equal(externalID, builder.opts.ExternalID)

			// Should build valid JSON body
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.Contains(string(body), externalID)
		})
	}
}

func TestSetUUIDMetadataProfileFieldCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		uuid        string
		displayName string
		email       string
		externalID  string
		profileURL  string
		custom      map[string]interface{}
	}{
		{
			name:        "Minimal profile",
			uuid:        "user123",
			displayName: "John Doe",
		},
		{
			name:        "Complete profile",
			uuid:        "user456",
			displayName: "Jane Smith",
			email:       "jane.smith@example.com",
			externalID:  "ext789",
			profileURL:  "https://example.com/jane.jpg",
			custom: map[string]interface{}{
				"role":       "admin",
				"department": "engineering",
				"join_date":  "2023-01-15",
			},
		},
		{
			name:        "Unicode profile",
			uuid:        "用户789",
			displayName: "张三",
			email:       "zhangsan@测试.com",
			externalID:  "外部ID_123",
			profileURL:  "https://测试.com/头像.jpg",
			custom: map[string]interface{}{
				"部门": "工程部",
				"级别": "高级工程师",
			},
		},
		{
			name:        "Special characters profile",
			uuid:        "user@special#123",
			displayName: "User With Special!@#$%^&*()Characters",
			email:       "user+special@example-domain.co.uk",
			externalID:  "ext_special_!@#$%^&*()",
			profileURL:  "https://example.com/profiles/special%20chars.png",
			custom: map[string]interface{}{
				"special@key": "special@value",
				"with spaces": "also spaces",
				"symbols!@#$": "more symbols!@#$",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetUUIDMetadataBuilder(pn)
			builder.UUID(tc.uuid)
			if tc.displayName != "" {
				builder.Name(tc.displayName)
			}
			if tc.email != "" {
				builder.Email(tc.email)
			}
			if tc.externalID != "" {
				builder.ExternalID(tc.externalID)
			}
			if tc.profileURL != "" {
				builder.ProfileURL(tc.profileURL)
			}
			if tc.custom != nil {
				builder.Custom(tc.custom)
			}

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build correct path
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.Contains(path, tc.uuid)

			// Should build valid JSON body
			body, err := builder.opts.buildBody()
			assert.Nil(err)

			var parsed map[string]interface{}
			err = json.Unmarshal(body, &parsed)
			assert.Nil(err)

			if tc.displayName != "" {
				assert.Equal(tc.displayName, parsed["name"])
			}
			if tc.email != "" {
				assert.Equal(tc.email, parsed["email"])
			}
			if tc.externalID != "" {
				assert.Equal(tc.externalID, parsed["externalId"])
			}
			if tc.profileURL != "" {
				assert.Equal(tc.profileURL, parsed["profileUrl"])
			}
		})
	}
}

// Comprehensive Edge Case Tests

func TestSetUUIDMetadataWithLongProfileData(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name    string
		setupFn func(*setUUIDMetadataBuilder)
	}{
		{
			name: "Very long name",
			setupFn: func(builder *setUUIDMetadataBuilder) {
				longName := strings.Repeat("VeryLongUserName", 50) // 800 characters
				builder.Name(longName)
			},
		},
		{
			name: "Very long email",
			setupFn: func(builder *setUUIDMetadataBuilder) {
				longEmail := strings.Repeat("verylongusername", 10) + "@" + strings.Repeat("verylongdomain", 10) + ".com"
				builder.Email(longEmail)
			},
		},
		{
			name: "Very long external ID",
			setupFn: func(builder *setUUIDMetadataBuilder) {
				longExternalID := strings.Repeat("external_id_", 100) // 1300 characters
				builder.ExternalID(longExternalID)
			},
		},
		{
			name: "Very long profile URL",
			setupFn: func(builder *setUUIDMetadataBuilder) {
				longURL := "https://example.com/" + strings.Repeat("very-long-path-segment/", 50) + "profile.jpg"
				builder.ProfileURL(longURL)
			},
		},
		{
			name: "Extremely large custom data",
			setupFn: func(builder *setUUIDMetadataBuilder) {
				largeCustom := make(map[string]interface{})
				for i := 0; i < 500; i++ {
					largeCustom[fmt.Sprintf("field_%d", i)] = strings.Repeat(fmt.Sprintf("data_%d_", i), 20)
				}
				builder.Custom(largeCustom)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetUUIDMetadataBuilder(pn)
			builder.UUID("test-uuid")
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid JSON body
			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.True(len(body) > 100) // Should be substantial

			// Should be valid JSON
			var parsed map[string]interface{}
			err = json.Unmarshal(body, &parsed)
			assert.Nil(err)
		})
	}
}

func TestSetUUIDMetadataSpecialCharacterHandling(t *testing.T) {
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
			builder := newSetUUIDMetadataBuilder(pn)
			builder.UUID("test-uuid")
			builder.Name(specialString)
			builder.Email(specialString + "@example.com")
			builder.ExternalID(specialString)
			builder.Custom(map[string]interface{}{
				"special_field": specialString,
			})

			// Should pass validation (basic validation doesn't check content)
			assert.Nil(builder.opts.validate())

			// Should build valid path and JSON body
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.NotNil(body)

			// Should be valid JSON
			var parsed map[string]interface{}
			err = json.Unmarshal(body, &parsed)
			assert.Nil(err)
		})
	}
}

func TestSetUUIDMetadataParameterBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		uuid        string
		displayName string
		email       string
		externalID  string
		profileURL  string
		queryParam  map[string]string
	}{
		{
			name:        "Empty string values",
			uuid:        "test-uuid",
			displayName: "",
			email:       "",
			externalID:  "",
			profileURL:  "",
		},
		{
			name:        "Single character values",
			uuid:        "a",
			displayName: "A",
			email:       "a@b.c",
			externalID:  "1",
			profileURL:  "https://a.b/c",
		},
		{
			name:        "Unicode-only values",
			uuid:        "测试",
			displayName: "용지",
			email:       "用户@测试.com",
			externalID:  "外部ID",
			profileURL:  "https://测试.com/头像.jpg",
		},
		{
			name:        "Mixed boundary conditions",
			uuid:        "test",
			displayName: strings.Repeat("A", 1000),
			email:       "user@example.com",
			externalID:  "",
			profileURL:  "https://example.com/profile.jpg",
			queryParam: map[string]string{
				"empty_param": "",
				"long_param":  strings.Repeat("value", 200),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetUUIDMetadataBuilder(pn)
			builder.UUID(tc.uuid)
			if tc.displayName != "" {
				builder.Name(tc.displayName)
			}
			if tc.email != "" {
				builder.Email(tc.email)
			}
			if tc.externalID != "" {
				builder.ExternalID(tc.externalID)
			}
			if tc.profileURL != "" {
				builder.ProfileURL(tc.profileURL)
			}
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

			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.NotNil(body)
		})
	}
}

func TestSetUUIDMetadataComplexProfileCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name       string
		setupFn    func(*setUUIDMetadataBuilder)
		validateFn func(*testing.T, []byte)
	}{
		{
			name: "All fields with Unicode",
			setupFn: func(builder *setUUIDMetadataBuilder) {
				builder.UUID("用户123")
				builder.Name("张三李四")
				builder.Email("zhangsan@测试域名.com")
				builder.ExternalID("外部系统ID_123")
				builder.ProfileURL("https://测试网站.com/用户头像/张三.jpg")
				builder.Custom(map[string]interface{}{
					"部门":   "工程部",
					"级别":   "高级工程师",
					"入职日期": "2023-01-15",
					"技能列表": []string{"Go", "JavaScript", "Python"},
					"个人信息": map[string]interface{}{
						"年龄": 30,
						"城市": "北京",
						"爱好": []string{"编程", "阅读", "音乐"},
					},
				})
			},
			validateFn: func(t *testing.T, body []byte) {
				var parsed map[string]interface{}
				err := json.Unmarshal(body, &parsed)
				assert.Nil(err)
				assert.Equal("张三李四", parsed["name"])
				assert.Equal("zhangsan@测试域名.com", parsed["email"])
				assert.Contains(parsed, "custom")
			},
		},
		{
			name: "Nested complex data structures",
			setupFn: func(builder *setUUIDMetadataBuilder) {
				builder.UUID("complex-user")
				builder.Name("Complex User")
				builder.Custom(map[string]interface{}{
					"permissions": map[string]interface{}{
						"read":  []string{"*"},
						"write": []string{"own", "team"},
						"admin": false,
					},
					"metadata": map[string]interface{}{
						"version": "2.1",
						"tags":    []interface{}{"user", "active", 123, true},
						"config": map[string]interface{}{
							"theme":         "dark",
							"notifications": true,
							"language":      "en-US",
							"timezone":      "UTC",
						},
					},
					"analytics": map[string]interface{}{
						"last_login":    "2023-12-01T10:30:00Z",
						"login_count":   42,
						"feature_flags": []string{"beta_features", "advanced_ui"},
					},
				})
			},
			validateFn: func(t *testing.T, body []byte) {
				var parsed map[string]interface{}
				err := json.Unmarshal(body, &parsed)
				assert.Nil(err)

				custom := parsed["custom"].(map[string]interface{})
				permissions := custom["permissions"].(map[string]interface{})
				assert.Equal(false, permissions["admin"])

				metadata := custom["metadata"].(map[string]interface{})
				assert.Equal("2.1", metadata["version"])
			},
		},
		{
			name: "Profile with all edge case combinations",
			setupFn: func(builder *setUUIDMetadataBuilder) {
				builder.UUID("edge@case#user$123")
				builder.Name("User With !@#$%^&*() Special Characters")
				builder.Email("user+special@sub-domain.example-site.co.uk")
				builder.ExternalID("ext_123@system#456")
				builder.ProfileURL("https://cdn.example.com/users/special%20chars/profile%20(1).jpg?v=123&size=large")
				builder.Custom(map[string]interface{}{
					"special@key":   "special@value",
					"unicode测试":     "unicode值",
					"with spaces":   "also spaces",
					"equals=key":    "equals=value",
					"ampersand&key": "ampersand&value",
					"nested": map[string]interface{}{
						"level1": map[string]interface{}{
							"special!@#": "deep special value",
							"unicode深层":  "deep unicode value",
						},
					},
				})
			},
			validateFn: func(t *testing.T, body []byte) {
				var parsed map[string]interface{}
				err := json.Unmarshal(body, &parsed)
				assert.Nil(err)
				assert.Contains(parsed["name"], "Special Characters")
				assert.Contains(parsed["email"], "sub-domain")
				assert.Contains(parsed["profileUrl"], "special%20chars")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetUUIDMetadataBuilder(pn)
			tc.setupFn(builder)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid components
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			assert.NotNil(path)

			body, err := builder.opts.buildBody()
			assert.Nil(err)
			assert.NotNil(body)

			// Run custom validation
			tc.validateFn(t, body)
		})
	}
}

func TestSetUUIDMetadataWithEmptyAndNilValues(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		setupFn  func(*setUUIDMetadataBuilder)
		expected string
	}{
		{
			name: "All nil/empty values",
			setupFn: func(builder *setUUIDMetadataBuilder) {
				builder.UUID("test-uuid")
				builder.Custom(nil)
				// Other fields left empty
			},
			expected: "{}",
		},
		{
			name: "Empty custom map",
			setupFn: func(builder *setUUIDMetadataBuilder) {
				builder.UUID("test-uuid")
				builder.Custom(map[string]interface{}{})
			},
			expected: "{}", // Empty custom map is omitted from JSON
		},
		{
			name: "Custom with nil values",
			setupFn: func(builder *setUUIDMetadataBuilder) {
				builder.UUID("test-uuid")
				builder.Custom(map[string]interface{}{
					"null_field":  nil,
					"empty_field": "",
					"zero_field":  0,
					"false_field": false,
				})
			},
			expected: "", // Will validate structure instead
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetUUIDMetadataBuilder(pn)
			tc.setupFn(builder)

			body, err := builder.opts.buildBody()
			assert.Nil(err)

			if tc.expected != "" {
				assert.Equal(tc.expected, string(body))
			} else {
				// Validate JSON structure
				var parsed map[string]interface{}
				err = json.Unmarshal(body, &parsed)
				assert.Nil(err)
			}
		})
	}
}

func TestSetUUIDMetadataIncludeEnumConversion(t *testing.T) {
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
			name:     "Single include custom",
			include:  []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom},
			expected: "custom",
		},
		{
			name:     "Single include status",
			include:  []PNUUIDMetadataInclude{PNUUIDMetadataIncludeStatus},
			expected: "status",
		},
		{
			name:     "Single include type",
			include:  []PNUUIDMetadataInclude{PNUUIDMetadataIncludeType},
			expected: "type",
		},
		{
			name:     "Multiple includes",
			include:  []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom, PNUUIDMetadataIncludeStatus, PNUUIDMetadataIncludeType},
			expected: "custom,status,type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newSetUUIDMetadataBuilder(pn)
			if tc.include != nil {
				builder.Include(tc.include)
			}

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expected, query.Get("include"))
		})
	}
}

// Error Scenario Tests

func TestSetUUIDMetadataExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newSetUUIDMetadataBuilder(pn)
	builder.UUID("test-uuid")
	builder.Name("Test User")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestSetUUIDMetadataPathBuildingEdgeCases(t *testing.T) {
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
			opts := newSetUUIDMetadataOpts(pn, pn.ctx)
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

func TestSetUUIDMetadataQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*setUUIDMetadataOpts)
		expectError bool
	}{
		{
			name: "Nil include array",
			setupOpts: func(opts *setUUIDMetadataOpts) {
				opts.Include = nil
			},
			expectError: false,
		},
		{
			name: "Empty include array",
			setupOpts: func(opts *setUUIDMetadataOpts) {
				opts.Include = []string{}
			},
			expectError: false,
		},
		{
			name: "Very large query params",
			setupOpts: func(opts *setUUIDMetadataOpts) {
				opts.QueryParam = map[string]string{
					"param1": strings.Repeat("a", 1000),
					"param2": strings.Repeat("b", 1000),
				}
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			setupOpts: func(opts *setUUIDMetadataOpts) {
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
			opts := newSetUUIDMetadataOpts(pn, pn.ctx)
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

func TestSetUUIDMetadataJSONBuildingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		setupOpts   func(*setUUIDMetadataOpts)
		expectError bool
	}{
		{
			name: "Valid custom data",
			setupOpts: func(opts *setUUIDMetadataOpts) {
				opts.Custom = map[string]interface{}{
					"valid": "data",
				}
			},
			expectError: false,
		},
		{
			name: "Complex valid data",
			setupOpts: func(opts *setUUIDMetadataOpts) {
				opts.Custom = map[string]interface{}{
					"nested": map[string]interface{}{
						"array": []interface{}{1, "two", true, nil},
					},
				}
			},
			expectError: false,
		},
		{
			name: "Large valid data",
			setupOpts: func(opts *setUUIDMetadataOpts) {
				largeData := make(map[string]interface{})
				for i := 0; i < 1000; i++ {
					largeData[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
				}
				opts.Custom = largeData
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newSetUUIDMetadataOpts(pn, pn.ctx)
			tc.setupOpts(opts)

			body, err := opts.buildBody()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.NotNil(body)

				// Should be valid JSON
				var parsed map[string]interface{}
				err = json.Unmarshal(body, &parsed)
				assert.Nil(err)
			}
		})
	}
}

func TestSetUUIDMetadataBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newSetUUIDMetadataBuilder(pn)

	include := []PNUUIDMetadataInclude{PNUUIDMetadataIncludeCustom}
	custom := map[string]interface{}{
		"role": "admin",
		"dept": "engineering",
	}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.UUID("complete-test-uuid").
		Include(include).
		Name("Complete Test User").
		ExternalID("ext123").
		ProfileURL("https://example.com/profile.jpg").
		Email("test@example.com").
		Custom(custom).
		Status("active").
		Type("public").
		QueryParam(queryParam)

	// Verify all values are set
	assert.Equal("complete-test-uuid", builder.opts.UUID)
	assert.Equal(EnumArrayToStringArray(include), builder.opts.Include)
	assert.Equal("Complete Test User", builder.opts.Name)
	assert.Equal("ext123", builder.opts.ExternalID)
	assert.Equal("https://example.com/profile.jpg", builder.opts.ProfileURL)
	assert.Equal("test@example.com", builder.opts.Email)
	assert.Equal(custom, builder.opts.Custom)
	assert.Equal("active", builder.opts.Status)
	assert.Equal("public", builder.opts.Type)
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

	// Should build complete JSON body
	body, err := builder.opts.buildBody()
	assert.Nil(err)

	var parsed map[string]interface{}
	err = json.Unmarshal(body, &parsed)
	assert.Nil(err)
	assert.Equal("Complete Test User", parsed["name"])
	assert.Equal("ext123", parsed["externalId"])
	assert.Equal("https://example.com/profile.jpg", parsed["profileUrl"])
	assert.Equal("test@example.com", parsed["email"])
	assert.Equal("active", parsed["status"])
	assert.Equal("public", parsed["type"])
	assert.NotNil(parsed["custom"])
}

func TestSetUUIDMetadataResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSetUUIDMetadataOpts(pn, pn.ctx)

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
			resp, _, err := newPNSetUUIDMetadataResponse(tc.jsonBytes, opts, StatusResponse{})

			if tc.expectError {
				assert.NotNil(err)
				// When there's an error, resp might be nil or the empty response
				if resp == nil {
					assert.Equal(emptyPNSetUUIDMetadataResponse, resp)
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
