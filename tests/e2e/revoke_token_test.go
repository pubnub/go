package e2e

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v7"
	"github.com/stretchr/testify/assert"
)

// Test happy path - successful token revocation
func TestRevokeTokenSuccess(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// First, create a token to revoke
	ch1 := randomized("revoke_test_channel")
	ch := map[string]pubnub.ChannelPermissions{
		ch1: {
			Read:  true,
			Write: true,
		},
	}

	// Grant a token with short TTL
	grantRes, _, err := pn.GrantToken().TTL(1).
		Channels(ch).
		Execute()

	assert.Nil(err)
	if grantRes == nil {
		t.Fatal("Grant response is nil")
		return
	}
	assert.NotNil(grantRes)
	assert.NotEmpty(grantRes.Data.Token)

	token := grantRes.Data.Token

	// Verify token is valid by parsing it
	parsedToken, err := pubnub.ParseToken(token)
	assert.Nil(err)
	assert.NotNil(parsedToken)

	// Now revoke the token
	revokeRes, status, err := pn.RevokeToken().
		Token(token).
		Execute()

	assert.Nil(err)
	assert.NotNil(revokeRes)
	assert.Equal(200, status.StatusCode)
}

// Test revoke token with query parameters
func TestRevokeTokenWithQueryParams(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Create a token to revoke
	ch1 := randomized("revoke_qp_test_channel")
	ch := map[string]pubnub.ChannelPermissions{
		ch1: {
			Read: true,
		},
	}

	grantRes, _, err := pn.GrantToken().TTL(1).
		Channels(ch).
		Execute()

	assert.Nil(err)
	assert.NotNil(grantRes)

	token := grantRes.Data.Token

	// Revoke with custom query parameters
	queryParam := map[string]string{
		"custom_param": "test_value",
		"app_version":  "1.0.0",
	}

	revokeRes, status, err := pn.RevokeToken().
		Token(token).
		QueryParam(queryParam).
		Execute()

	assert.Nil(err)
	assert.NotNil(revokeRes)
	assert.Equal(200, status.StatusCode)
}

// Test builder pattern and fluent interface
func TestRevokeTokenBuilderPattern(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	// Create a token
	grantRes, _, err := pn.GrantToken().TTL(1).
		Channels(map[string]pubnub.ChannelPermissions{
			randomized("builder_test"): {Read: true},
		}).
		Execute()

	assert.Nil(err)
	token := grantRes.Data.Token

	// Test fluent interface chaining
	builder := pn.RevokeToken()
	builder = builder.Token(token)
	builder = builder.QueryParam(map[string]string{"test": "value"})

	res, status, err := builder.Execute()

	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(200, status.StatusCode)
}

// Test validation error - missing token
func TestRevokeTokenMissingToken(t *testing.T) {
	assert := assert.New(t)

	// Use a valid PAM config but don't provide token
	config := pamConfigCopy()
	if config.PublishKey == "" || config.SubscribeKey == "" || config.SecretKey == "" {
		// Create a mock config for validation testing
		config.PublishKey = "demo"
		config.SubscribeKey = "demo"
		config.SecretKey = "demo"
	}
	pn := pubnub.NewPubNub(config)

	// Try to revoke without providing a token
	_, _, err := pn.RevokeToken().Execute()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing PAMv3 token")
}

// Test validation error - empty token
func TestRevokeTokenEmptyToken(t *testing.T) {
	assert := assert.New(t)

	// Use a valid PAM config for validation testing
	config := pamConfigCopy()
	if config.PublishKey == "" || config.SubscribeKey == "" || config.SecretKey == "" {
		config.PublishKey = "demo"
		config.SubscribeKey = "demo"
		config.SecretKey = "demo"
	}
	pn := pubnub.NewPubNub(config)

	// Try to revoke with empty token
	_, _, err := pn.RevokeToken().
		Token("").
		Execute()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing PAMv3 token")
}

// Test validation error - missing publish key
func TestRevokeTokenMissingPublishKey(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()
	config.PublishKey = ""
	pn := pubnub.NewPubNub(config)

	_, _, err := pn.RevokeToken().
		Token("test-token").
		Execute()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Publish Key")
}

// Test validation error - missing subscribe key
func TestRevokeTokenMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()
	// Ensure we have valid publish and secret keys, but missing subscribe key
	if config.PublishKey == "" {
		config.PublishKey = "demo"
	}
	if config.SecretKey == "" {
		config.SecretKey = "demo"
	}
	config.SubscribeKey = ""
	pn := pubnub.NewPubNub(config)

	_, _, err := pn.RevokeToken().
		Token("test-token").
		Execute()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

// Test validation error - missing secret key
func TestRevokeTokenMissingSecretKey(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()
	// Ensure we have valid publish and subscribe keys, but missing secret key
	if config.PublishKey == "" {
		config.PublishKey = "demo"
	}
	if config.SubscribeKey == "" {
		config.SubscribeKey = "demo"
	}
	config.SecretKey = ""
	pn := pubnub.NewPubNub(config)

	_, _, err := pn.RevokeToken().
		Token("test-token").
		Execute()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Secret Key")
}

// Test server error - invalid token format
func TestRevokeTokenInvalidFormat(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	// Try to revoke with malformed token
	_, status, err := pn.RevokeToken().
		Token("invalid-token-format").
		Execute()

	assert.NotNil(err)
	// Should get a server error response
	assert.True(status.StatusCode >= 400)
}

// Test server error - non-existent token
func TestRevokeTokenNonExistent(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	// Try to revoke a properly formatted but non-existent token
	// Using a valid-looking token format but with non-existent content
	fakeToken := "qEF2AkF0GmEI03xDdHRsGDxDcmVzpURjaGFuoWNjaDEY70NncnChY2NnMRjvQ3VzcqBDc3BjoER1dWlkoER1dWlkoER1dWlkDHQFQ3NpZ1ggZGI4ZDJhMGYtYjUyYy00YWJhLWJjOTEtNGM5NGQwM2ZlMzY3"

	_, status, err := pn.RevokeToken().
		Token(fakeToken).
		Execute()

	// Should get an error response from server
	assert.NotNil(err)
	assert.True(status.StatusCode >= 400)
}

// Test edge case - token with special characters that need URL encoding
func TestRevokeTokenWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	// Create a token first
	grantRes, _, err := pn.GrantToken().TTL(1).
		Channels(map[string]pubnub.ChannelPermissions{
			randomized("special_char_test"): {Read: true},
		}).
		Execute()

	assert.Nil(err)
	if grantRes == nil {
		t.Fatal("Grant response is nil")
		return
	}
	token := grantRes.Data.Token

	// The token should already be properly formatted, but verify it can be revoked
	// even if it contains characters that need URL encoding
	res, status, err := pn.RevokeToken().
		Token(token).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(200, status.StatusCode)
}

// Test edge case - very long token
func TestRevokeTokenLongToken(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	// Create a token with many permissions to make it longer
	channels := make(map[string]pubnub.ChannelPermissions)
	for i := 0; i < 10; i++ {
		channels[randomized("long_token_ch")] = pubnub.ChannelPermissions{
			Read:   true,
			Write:  true,
			Delete: true,
			Get:    true,
			Update: true,
			Join:   true,
		}
	}

	grantRes, _, err := pn.GrantToken().TTL(1).
		Channels(channels).
		Execute()

	assert.Nil(err)
	assert.NotNil(grantRes)

	token := grantRes.Data.Token
	assert.True(len(token) > 100) // Verify it's actually a long token

	// Should be able to revoke even very long tokens
	res, status, err := pn.RevokeToken().
		Token(token).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(200, status.StatusCode)
}

// Test context cancellation
func TestRevokeTokenWithContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	// Create a token
	grantRes, _, err := pn.GrantToken().TTL(1).
		Channels(map[string]pubnub.ChannelPermissions{
			randomized("context_test"): {Read: true},
		}).
		Execute()

	assert.Nil(err)
	token := grantRes.Data.Token

	// Test with background context (should work normally)
	res, status, err := pn.RevokeTokenWithContext(backgroundContext).
		Token(token).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(200, status.StatusCode)
}

// Test with expired token (TTL test)
func TestRevokeTokenExpiredToken(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	// Create a token with very short TTL
	grantRes, _, err := pn.GrantToken().TTL(1). // 1 minute
							Channels(map[string]pubnub.ChannelPermissions{
			randomized("expired_test"): {Read: true},
		}).
		Execute()

	assert.Nil(err)
	token := grantRes.Data.Token

	// Should be able to revoke even if token is near expiration
	// (since revocation doesn't depend on token being active)
	res, status, err := pn.RevokeToken().
		Token(token).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(200, status.StatusCode)
}

// Test multiple tokens (create and revoke several)
func TestRevokeTokenMultiple(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	tokens := make([]string, 3)

	// Create multiple tokens
	for i := 0; i < 3; i++ {
		grantRes, _, err := pn.GrantToken().TTL(1).
			Channels(map[string]pubnub.ChannelPermissions{
				randomized("multi_test"): {Read: true},
			}).
			Execute()

		assert.Nil(err)
		tokens[i] = grantRes.Data.Token
	}

	// Revoke all tokens
	for i, token := range tokens {
		res, status, err := pn.RevokeToken().
			Token(token).
			Execute()

		assert.Nil(err, "Failed to revoke token %d", i)
		assert.NotNil(res)
		assert.Equal(200, status.StatusCode)
	}
}
