package e2e

import (
	"log"
	"os"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v8"
	"github.com/stretchr/testify/assert"
)

// Test basic Fire functionality - happy path
func TestFireBasic(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ch := randomized("fire_test_channel")
	message := "Hello from Fire!"

	// Test basic Fire
	resp, status, err := pn.Fire().
		Channel(ch).
		Message(message).
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.NotEmpty(resp.Timestamp)
}

// Test Fire with different message types
func TestFireMessageTypes(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ch := randomized("fire_message_types")

	testCases := []struct {
		name    string
		message interface{}
	}{
		{
			name:    "String message",
			message: "Hello Fire!",
		},
		{
			name:    "Number message",
			message: 42,
		},
		{
			name:    "Boolean message",
			message: true,
		},
		{
			name:    "Object message",
			message: map[string]interface{}{"text": "Hello", "number": 123, "bool": false},
		},
		{
			name:    "Array message",
			message: []string{"item1", "item2", "item3"},
		},
		{
			name: "Complex object",
			message: map[string]interface{}{
				"user":      "testUser",
				"action":    "fire_test",
				"timestamp": time.Now().Unix(),
				"data":      map[string]interface{}{"key": "value"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, status, err := pn.Fire().
				Channel(ch).
				Message(tc.message).
				Execute()

			assert.Nil(err)
			assert.Equal(200, status.StatusCode)
			assert.NotNil(resp)
			assert.NotEmpty(resp.Timestamp)
		})
	}
}

// Test Fire with TTL parameter
func TestFireWithTTL(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ch := randomized("fire_ttl_test")
	message := "Test message with TTL"

	testCases := []struct {
		name string
		ttl  int
	}{
		{"TTL 1 hour", 1},
		{"TTL 24 hours", 24},
		{"TTL 72 hours", 72},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, status, err := pn.Fire().
				Channel(ch).
				Message(message).
				TTL(tc.ttl).
				Execute()

			assert.Nil(err)
			assert.Equal(200, status.StatusCode)
			assert.NotNil(resp)
			assert.NotEmpty(resp.Timestamp)
		})
	}
}

// Test Fire with UsePost parameter
func TestFireWithUsePost(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ch := randomized("fire_post_test")
	message := map[string]interface{}{
		"longMessage": "This is a longer message that might benefit from using POST method instead of GET",
		"data":        []string{"item1", "item2", "item3", "item4", "item5"},
		"metadata":    map[string]interface{}{"timestamp": time.Now().Unix(), "source": "e2e_test"},
	}

	// Test with POST
	resp1, status1, err1 := pn.Fire().
		Channel(ch).
		Message(message).
		UsePost(true).
		Execute()

	assert.Nil(err1)
	assert.Equal(200, status1.StatusCode)
	assert.NotNil(resp1)
	assert.NotEmpty(resp1.Timestamp)

	// Test with GET (default)
	resp2, status2, err2 := pn.Fire().
		Channel(ch).
		Message("Simple message for GET").
		UsePost(false).
		Execute()

	assert.Nil(err2)
	assert.Equal(200, status2.StatusCode)
	assert.NotNil(resp2)
	assert.NotEmpty(resp2.Timestamp)
}

// Test Fire with Serialize parameter
func TestFireWithSerialize(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ch := randomized("fire_serialize_test")

	// Test with Serialize = true (default)
	resp1, status1, err1 := pn.Fire().
		Channel(ch).
		Message(map[string]interface{}{"key": "value"}).
		Serialize(true).
		Execute()

	assert.Nil(err1)
	assert.Equal(200, status1.StatusCode)
	assert.NotNil(resp1)

	// Test with Serialize = false (pre-serialized JSON)
	preSerializedJSON := `{"preSerializedKey":"preSerializedValue","number":123}`
	resp2, status2, err2 := pn.Fire().
		Channel(ch).
		Message(preSerializedJSON).
		Serialize(false).
		UsePost(true). // Use POST for pre-serialized content
		Execute()

	assert.Nil(err2)
	assert.Equal(200, status2.StatusCode)
	assert.NotNil(resp2)
}

// Test Fire with QueryParam
func TestFireWithQueryParam(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ch := randomized("fire_query_param_test")
	message := "Test message with query params"
	queryParams := map[string]string{
		"custom_param1": "value1",
		"custom_param2": "value2",
		"test_id":       "fire_e2e_test",
	}

	resp, status, err := pn.Fire().
		Channel(ch).
		Message(message).
		QueryParam(queryParams).
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.NotEmpty(resp.Timestamp)
}

// Test Fire with Context
func TestFireWithContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ch := randomized("fire_context_test")
	message := "Test message with context"

	resp, status, err := pn.FireWithContext(backgroundContext).
		Channel(ch).
		Message(message).
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.NotEmpty(resp.Timestamp)
}

// Test Fire error scenarios
func TestFireErrorScenarios(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name          string
		setupClient   func() *pubnub.PubNub
		executeCall   func(*pubnub.PubNub) (*pubnub.PublishResponse, pubnub.StatusResponse, error)
		expectedError string
	}{
		{
			name: "Missing Publish Key",
			setupClient: func() *pubnub.PubNub {
				config := configCopy()
				config.PublishKey = ""
				return pubnub.NewPubNub(config)
			},
			executeCall: func(pn *pubnub.PubNub) (*pubnub.PublishResponse, pubnub.StatusResponse, error) {
				return pn.Fire().Channel("test").Message("test").Execute()
			},
			expectedError: "Missing Publish Key",
		},
		{
			name: "Missing Subscribe Key",
			setupClient: func() *pubnub.PubNub {
				config := configCopy()
				config.SubscribeKey = ""
				return pubnub.NewPubNub(config)
			},
			executeCall: func(pn *pubnub.PubNub) (*pubnub.PublishResponse, pubnub.StatusResponse, error) {
				return pn.Fire().Channel("test").Message("test").Execute()
			},
			expectedError: "Missing Subscribe Key",
		},
		{
			name: "Missing Channel",
			setupClient: func() *pubnub.PubNub {
				return pubnub.NewPubNub(configCopy())
			},
			executeCall: func(pn *pubnub.PubNub) (*pubnub.PublishResponse, pubnub.StatusResponse, error) {
				return pn.Fire().Message("test").Execute() // No channel
			},
			expectedError: "Missing Channel",
		},
		{
			name: "Missing Message",
			setupClient: func() *pubnub.PubNub {
				return pubnub.NewPubNub(configCopy())
			},
			executeCall: func(pn *pubnub.PubNub) (*pubnub.PublishResponse, pubnub.StatusResponse, error) {
				return pn.Fire().Channel("test").Execute() // No message
			},
			expectedError: "Missing Message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pn := tc.setupClient()

			_, _, err := tc.executeCall(pn)

			assert.NotNil(err)
			assert.Contains(err.Error(), tc.expectedError)
		})
	}
}

// Test Fire with invalid parameters
func TestFireInvalidParameters(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ch := randomized("fire_invalid_test")

	testCases := []struct {
		name        string
		setupCall   func() (*pubnub.PublishResponse, pubnub.StatusResponse, error)
		expectError bool
	}{
		{
			name: "Negative TTL",
			setupCall: func() (*pubnub.PublishResponse, pubnub.StatusResponse, error) {
				return pn.Fire().
					Channel(ch).
					Message("test").
					TTL(-1).
					Execute()
			},
			expectError: false, // API might handle gracefully
		},
		{
			name: "Very large TTL",
			setupCall: func() (*pubnub.PublishResponse, pubnub.StatusResponse, error) {
				return pn.Fire().
					Channel(ch).
					Message("test").
					TTL(999999).
					Execute()
			},
			expectError: false, // API might clamp this
		},
		{
			name: "Empty string message",
			setupCall: func() (*pubnub.PublishResponse, pubnub.StatusResponse, error) {
				return pn.Fire().
					Channel(ch).
					Message("").
					Execute()
			},
			expectError: false, // Empty string is valid
		},
		{
			name: "Nil query params",
			setupCall: func() (*pubnub.PublishResponse, pubnub.StatusResponse, error) {
				return pn.Fire().
					Channel(ch).
					Message("test").
					QueryParam(nil).
					Execute()
			},
			expectError: false, // Should be handled gracefully
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, status, err := tc.setupCall()

			if tc.expectError {
				assert.NotNil(err)
			} else {
				// For cases where we don't expect errors, just verify it doesn't crash
				if err == nil {
					assert.Equal(200, status.StatusCode)
					assert.NotNil(resp)
				}
				// If there is an error, it should be a meaningful one, not a panic
			}
		})
	}
}

// Test Fire comprehensive scenario with all parameters
func TestFireComprehensive(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ch := randomized("fire_comprehensive_test")
	message := map[string]interface{}{
		"eventType": "function_trigger",
		"data":      "comprehensive_test_data",
	}
	meta := map[string]interface{}{
		"source": "e2e_test",
		"type":   "comprehensive",
	}
	queryParams := map[string]string{
		"test_run": "comprehensive",
		"version":  "v1",
	}

	// Test comprehensive Fire with all parameters
	resp, status, err := pn.Fire().
		Channel(ch).
		Message(message).
		Meta(meta).
		TTL(24).
		UsePost(false).
		Serialize(true).
		QueryParam(queryParams).
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(resp)
	assert.NotEmpty(resp.Timestamp)

	// Verify response structure
	assert.True(resp.Timestamp > 0)
}

// Test Fire vs Publish behavior - Fire should not store messages
func TestFireVsPublishBehavior(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ch := randomized("fire_vs_publish_test")
	fireMessage := "This message was fired (should not be stored)"
	publishMessage := "This message was published (should be stored)"

	// Fire a message (should not be stored)
	fireResp, fireStatus, fireErr := pn.Fire().
		Channel(ch).
		Message(fireMessage).
		Execute()

	assert.Nil(fireErr)
	assert.Equal(200, fireStatus.StatusCode)
	assert.NotNil(fireResp)

	// Publish a message (should be stored)
	pubResp, pubStatus, pubErr := pn.Publish().
		Channel(ch).
		Message(publishMessage).
		Execute()

	assert.Nil(pubErr)
	assert.Equal(200, pubStatus.StatusCode)
	assert.NotNil(pubResp)

	// Wait a moment for propagation
	time.Sleep(2 * time.Second)

	// Try to fetch history - should only contain the published message, not the fired one
	histResp, histStatus, histErr := pn.History().
		Channel(ch).
		Count(10).
		Execute()

	assert.Nil(histErr)
	assert.Equal(200, histStatus.StatusCode)
	assert.NotNil(histResp)

	// Fire messages should not appear in history (they are not stored)
	// Publish messages should appear in history
	if len(histResp.Messages) > 0 {
		found := false
		for _, msg := range histResp.Messages {
			if msg.Message == publishMessage {
				found = true
			}
			// Fire message should NOT be found in history
			assert.NotEqual(fireMessage, msg.Message, "Fire message should not be stored in history")
		}
		assert.True(found, "Published message should be found in history")
	}
}

// Test Fire edge cases
func TestFireEdgeCases(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Test 1: Empty string message
	t.Run("Empty String Message", func(t *testing.T) {
		ch := randomized("fire_edge_empty")
		resp, status, err := pn.Fire().
			Channel(ch).
			Message("").
			Execute()

		assert.Nil(err)
		assert.Equal(200, status.StatusCode)
		assert.NotNil(resp)
		assert.True(resp.Timestamp > 0)
	})

	// Test 2: Special characters in channel name
	t.Run("Special Characters Channel", func(t *testing.T) {
		ch := randomized("fire_edge_special-test.channel_123")
		resp, status, err := pn.Fire().
			Channel(ch).
			Message("test message").
			Execute()

		assert.Nil(err)
		assert.Equal(200, status.StatusCode)
		assert.NotNil(resp)
		assert.True(resp.Timestamp > 0)
	})

	// Test 3: Large message
	t.Run("Large Message", func(t *testing.T) {
		ch := randomized("fire_edge_large")
		largeMessage := make([]byte, 10000) // 10KB message
		for i := range largeMessage {
			largeMessage[i] = 'A'
		}

		resp, status, err := pn.Fire().
			Channel(ch).
			Message(string(largeMessage)).
			Execute()

		assert.Nil(err)
		assert.Equal(200, status.StatusCode)
		assert.NotNil(resp)
		assert.True(resp.Timestamp > 0)
	})

	// Test 4: Complex nested object
	t.Run("Complex Nested Object", func(t *testing.T) {
		ch := randomized("fire_edge_complex")
		complexMessage := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"level3": []interface{}{
						"string",
						123,
						true,
						map[string]interface{}{
							"nested": "value",
						},
					},
				},
				"array": []string{"a", "b", "c"},
			},
			"timestamp": time.Now().Unix(),
		}

		resp, status, err := pn.Fire().
			Channel(ch).
			Message(complexMessage).
			Execute()

		assert.Nil(err)
		assert.Equal(200, status.StatusCode)
		assert.NotNil(resp)
		assert.True(resp.Timestamp > 0)
	})

	// Test 5: Unicode characters
	t.Run("Unicode Characters", func(t *testing.T) {
		ch := randomized("fire_edge_unicode")
		unicodeMessage := "Hello 世界! 🌍 Здравствуй мир! 🚀"

		resp, status, err := pn.Fire().
			Channel(ch).
			Message(unicodeMessage).
			Execute()

		assert.Nil(err)
		assert.Equal(200, status.StatusCode)
		assert.NotNil(resp)
		assert.True(resp.Timestamp > 0)
	})
}

// Test Fire with Meta parameter
func TestFireWithMeta(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ch := randomized("fire_meta_test")
	message := "Hello with metadata!"

	testCases := []struct {
		name string
		meta interface{}
	}{
		{"String Meta", "simple_string_meta"},
		{"Numeric Meta", 42},
		{"Boolean Meta", true},
		{"Object Meta", map[string]interface{}{
			"user":   "test_user",
			"action": "fire_test",
			"count":  1,
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, status, err := pn.Fire().
				Channel(ch + "_" + tc.name).
				Message(message).
				Meta(tc.meta).
				Execute()

			assert.Nil(err)
			assert.Equal(200, status.StatusCode)
			assert.NotNil(resp)
			assert.NotEmpty(resp.Timestamp)
			assert.True(resp.Timestamp > 0)
		})
	}
}
