package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertRevokeToken(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newRevokeTokenBuilder(pn)
	if testContext {
		o = newRevokeTokenBuilderWithContext(pn, pn.ctx)
	}

	token := "token"
	o.QueryParam(queryParam)
	o.Token(token)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(revokeTokenPath, pn.Config.SubscribeKey, token),
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

func TestRevokeToken(t *testing.T) {
	AssertRevokeToken(t, true, false)
}

func TestRevokeTokenContext(t *testing.T) {
	AssertRevokeToken(t, true, true)
}

func TestRevokeTokenResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRevokeTokenOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNRevokeTokenResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRevokeTokenResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRevokeTokenOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200}`)

	_, s, err := newPNRevokeTokenResponse(jsonBytes, opts, StatusResponse{StatusCode: 200})
	assert.Equal(200, s.StatusCode)

	assert.Nil(err)
}

// Additional validation tests specific to RevokeToken
func TestRevokeTokenValidateMissingPublishKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.PublishKey = ""
	opts := newRevokeTokenOpts(pn, pn.ctx)
	opts.Token = "test-token"

	assert.Equal("pubnub/validation: pubnub: No Category Matched: Missing Publish Key", opts.validate().Error())
}

func TestRevokeTokenValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newRevokeTokenOpts(pn, pn.ctx)
	opts.Token = "test-token"

	assert.Equal("pubnub/validation: pubnub: No Category Matched: Missing Subscribe Key", opts.validate().Error())
}

func TestRevokeTokenValidateMissingSecretKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SecretKey = ""
	opts := newRevokeTokenOpts(pn, pn.ctx)
	opts.Token = "test-token"

	assert.Equal("pubnub/validation: pubnub: No Category Matched: Missing Secret Key", opts.validate().Error())
}

func TestRevokeTokenValidateMissingToken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRevokeTokenOpts(pn, pn.ctx)
	opts.Token = ""

	assert.Equal("pubnub/validation: pubnub: No Category Matched: Missing PAMv3 token", opts.validate().Error())
}

func TestRevokeTokenValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRevokeTokenOpts(pn, pn.ctx)
	opts.Token = "valid-token"

	assert.Nil(opts.validate())
}

// Builder pattern tests for RevokeToken
func TestRevokeTokenBuilder(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test basic builder
	builder := newRevokeTokenBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)

	// Test Token setting
	result := builder.Token("test-token-123")
	assert.Equal(builder, result) // Should return same instance for chaining
	assert.Equal("test-token-123", builder.opts.Token)
}

func TestRevokeTokenBuilderWithContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRevokeTokenBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestRevokeTokenBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParams := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	builder := newRevokeTokenBuilder(pn)
	result := builder.QueryParam(queryParams)
	assert.Equal(builder, result) // Should return same instance for chaining
	assert.Equal(queryParams, builder.opts.QueryParam)
}

func TestRevokeTokenBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParams := map[string]string{
		"test1": "value1",
		"test2": "value2",
	}

	// Test fluent interface chaining
	builder := newRevokeTokenBuilder(pn).
		Token("chained-token-abc123").
		QueryParam(queryParams)

	assert.Equal("chained-token-abc123", builder.opts.Token)
	assert.Equal(queryParams, builder.opts.QueryParam)
}

func TestRevokeTokenBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRevokeTokenBuilder(pn)

	// Test Token setter
	builder.Token("setter-test-token")
	assert.Equal("setter-test-token", builder.opts.Token)

	// Test QueryParam setter
	queryParams := map[string]string{"key": "value"}
	builder.QueryParam(queryParams)
	assert.Equal(queryParams, builder.opts.QueryParam)

	// Test overwriting Token
	builder.Token("new-token")
	assert.Equal("new-token", builder.opts.Token)

	// Test overwriting QueryParam
	newQueryParams := map[string]string{"newkey": "newvalue"}
	builder.QueryParam(newQueryParams)
	assert.Equal(newQueryParams, builder.opts.QueryParam)
}

// URL path building tests
func TestRevokeTokenBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newRevokeTokenOpts(pn, pn.ctx)
	opts.Token = "test-token"

	path, err := opts.buildPath()
	assert.Nil(err)
	expectedPath := fmt.Sprintf("/v3/pam/%s/grant/test-token", pn.Config.SubscribeKey)
	assert.Equal(expectedPath, path)
}

func TestRevokeTokenBuildPathWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newRevokeTokenOpts(pn, pn.ctx)
	// Token with characters that need URL encoding
	opts.Token = "token+with/special=chars&more"

	path, err := opts.buildPath()
	assert.Nil(err)

	// Should contain URL-encoded token
	assert.Contains(path, "/v3/pam/")
	assert.Contains(path, pn.Config.SubscribeKey)
	assert.Contains(path, "/grant/")
	// The token should be URL encoded in the path
	assert.NotContains(path, "+") // + should be encoded
	assert.NotContains(path, "=") // = should be encoded
}

func TestRevokeTokenBuildQuery(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newRevokeTokenOpts(pn, pn.ctx)
	opts.QueryParam = map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should include custom query parameters
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))

	// Should include default PubNub parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestRevokeTokenBuildQueryEmpty(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newRevokeTokenOpts(pn, pn.ctx)
	opts.QueryParam = nil

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should still include default PubNub parameters even with no custom params
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

// HTTP method and operation type tests
func TestRevokeTokenHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newRevokeTokenOpts(pn, pn.ctx)
	method := opts.httpMethod()
	assert.Equal("DELETE", method)
}

func TestRevokeTokenOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newRevokeTokenOpts(pn, pn.ctx)
	opType := opts.operationType()
	assert.Equal(PNAccessManagerRevokeToken, opType)
}

// Edge case tests for RevokeToken
func TestRevokeTokenWithVeryLongToken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a very long token string
	longToken := ""
	for i := 0; i < 500; i++ {
		longToken += "a"
	}

	builder := newRevokeTokenBuilder(pn)
	builder.Token(longToken)

	assert.Equal(longToken, builder.opts.Token)
	assert.Equal(500, len(builder.opts.Token))

	// Test path building with long token
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v3/pam/")
	assert.Contains(path, "/grant/")
}

func TestRevokeTokenWithUnicodeCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Token with Unicode characters
	unicodeToken := "token-测试-русский-ファイル-🔒"

	builder := newRevokeTokenBuilder(pn)
	builder.Token(unicodeToken)

	assert.Equal(unicodeToken, builder.opts.Token)

	// Test path building with Unicode token
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v3/pam/")
	assert.Contains(path, "/grant/")
	// Unicode should be properly URL encoded
	assert.NotContains(path, "测试")      // Should be encoded
	assert.NotContains(path, "русский") // Should be encoded
	assert.NotContains(path, "🔒")       // Should be encoded
}

func TestRevokeTokenWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Token with various special characters
	specialToken := "token-with-dashes_and_underscores.and.dots:and:colons@and@symbols#and#hashes$and$dollars%and%percents"

	builder := newRevokeTokenBuilder(pn)
	builder.Token(specialToken)

	assert.Equal(specialToken, builder.opts.Token)

	// Test validation passes
	assert.Nil(builder.opts.validate())

	// Test path building
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v3/pam/")
	assert.Contains(path, "/grant/")
}

func TestRevokeTokenWithEmptyQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRevokeTokenBuilder(pn)
	builder.Token("test-token")
	builder.QueryParam(map[string]string{}) // Empty map

	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestRevokeTokenWithNilQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newRevokeTokenBuilder(pn)
	builder.Token("test-token")
	builder.QueryParam(nil) // Nil map

	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

func TestRevokeTokenWithComplexQueryParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	complexParams := map[string]string{
		"app_version":    "1.2.3",
		"user_id":        "user-123-abc",
		"session_id":     "session-xyz-789",
		"special_chars":  "value-with-special@chars#and$symbols",
		"unicode_value":  "测试值-русский-ファイル",
		"empty_value":    "",
		"number_string":  "42",
		"boolean_string": "true",
	}

	builder := newRevokeTokenBuilder(pn)
	builder.Token("test-token")
	builder.QueryParam(complexParams)

	query, err := builder.opts.buildQuery()
	assert.Nil(err)

	// Verify all custom parameters are present
	for key, expectedValue := range complexParams {
		actualValue := query.Get(key)
		if key == "special_chars" {
			// Special characters are URL encoded
			assert.Equal("value-with-special%40chars%23and%24symbols", actualValue, "Query parameter %s should be URL encoded", key)
		} else if key == "unicode_value" {
			// Unicode characters are URL encoded
			assert.Contains(actualValue, "%E6%B5%8B", "Query parameter %s should contain URL encoded Unicode", key)
		} else {
			assert.Equal(expectedValue, actualValue, "Query parameter %s should match", key)
		}
	}

	// Should still have default parameters
	assert.NotEmpty(query.Get("uuid"))
	assert.NotEmpty(query.Get("pnsdk"))
}

// Error scenario tests
func TestRevokeTokenResponseParsingErrors(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRevokeTokenOpts(pn, pn.ctx)

	// Test completely invalid JSON
	invalidJSON := []byte(`{invalid json}`)
	_, _, err := newPNRevokeTokenResponse(invalidJSON, opts, StatusResponse{})
	assert.NotNil(err)
	assert.Contains(err.Error(), "parsing")

	// Test malformed but valid JSON
	malformedJSON := []byte(`{"unexpected": "structure", "not": "expected"}`)
	resp, status, err := newPNRevokeTokenResponse(malformedJSON, opts, StatusResponse{StatusCode: 200})
	assert.Nil(err) // Should not error on unexpected fields
	assert.NotNil(resp)
	assert.Equal(200, status.StatusCode)

	// Test null response
	nullJSON := []byte(`null`)
	resp, _, err = newPNRevokeTokenResponse(nullJSON, opts, StatusResponse{})
	if err != nil {
		assert.Contains(err.Error(), "parsing")
	}
	// For null JSON, the response parsing might succeed but return a valid response

	// Test empty response
	emptyJSON := []byte(``)
	_, _, err = newPNRevokeTokenResponse(emptyJSON, opts, StatusResponse{})
	assert.NotNil(err)
	if err != nil {
		assert.Contains(err.Error(), "parsing")
	}

	// Test response with wrong data types
	wrongTypeJSON := []byte(`{"status": "not-a-number"}`)
	resp, _, err = newPNRevokeTokenResponse(wrongTypeJSON, opts, StatusResponse{StatusCode: 200})
	assert.NotNil(err) // Should error when status field has wrong type
	assert.Contains(err.Error(), "parsing")
}

func TestRevokeTokenResponseWithDifferentStatusCodes(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRevokeTokenOpts(pn, pn.ctx)

	testCases := []struct {
		name       string
		jsonBytes  []byte
		statusCode int
	}{
		{
			name:       "Success response",
			jsonBytes:  []byte(`{"status": 200}`),
			statusCode: 200,
		},
		{
			name:       "Success with no body",
			jsonBytes:  []byte(`{}`),
			statusCode: 200,
		},
		{
			name:       "Success with additional fields",
			jsonBytes:  []byte(`{"status": 200, "message": "Token revoked", "data": {}}`),
			statusCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, status, err := newPNRevokeTokenResponse(tc.jsonBytes, opts, StatusResponse{StatusCode: tc.statusCode})
			assert.Nil(err, "Should not error for test case: %s", tc.name)
			assert.NotNil(resp, "Response should not be nil for test case: %s", tc.name)
			assert.Equal(tc.statusCode, status.StatusCode, "Status code should match for test case: %s", tc.name)
		})
	}
}

func TestRevokeTokenBuilderExecuteErrorHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with invalid configuration (missing secret key)
	pn.Config.SecretKey = ""

	builder := newRevokeTokenBuilder(pn)
	builder.Token("test-token")

	// Execute should fail with validation error
	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Secret Key")
}

func TestRevokeTokenEdgeCaseTokenFormats(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCaseTokens := []struct {
		name  string
		token string
	}{
		{"Base64-like token", "dGVzdC10b2tlbi1leGFtcGxl"},
		{"URL-safe base64", "dGVzdC10b2tlbi1leGFtcGxl_-"},
		{"JWT-like format", "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIn0.abc123"},
		{"Very short token", "a"},
		{"Numeric token", "123456789"},
		{"Mixed case", "AbCdEfGhIjKlMnOpQrStUvWxYz"},
		{"With padding", "token==="},
		{"With dashes", "token-with-many-dashes-in-between"},
		{"With underscores", "token_with_many_underscores_in_between"},
	}

	for _, tc := range edgeCaseTokens {
		t.Run(tc.name, func(t *testing.T) {
			builder := newRevokeTokenBuilder(pn)
			builder.Token(tc.token)

			assert.Equal(tc.token, builder.opts.Token)
			assert.Nil(builder.opts.validate(), "Token %s should pass validation", tc.token)

			// Test path building
			path, err := builder.opts.buildPath()
			assert.Nil(err, "Path building should succeed for token: %s", tc.token)
			assert.Contains(path, "/v3/pam/")
			assert.Contains(path, "/grant/")
		})
	}
}
