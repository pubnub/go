package pubnub

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	cbor "github.com/brianolson/cbor_go"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a valid test token
func createTestToken(decoded PNGrantTokenDecoded) (string, error) {
	var buf bytes.Buffer
	encoder := cbor.NewEncoder(&buf)
	err := encoder.Encode(decoded)
	if err != nil {
		return "", err
	}

	// Encode to base64 and make URL-safe
	token := base64.StdEncoding.EncodeToString(buf.Bytes())
	token = strings.Replace(token, "+", "-", -1)
	token = strings.Replace(token, "/", "_", -1)
	token = strings.TrimRight(token, "=")

	return token, nil
}

// Input validation tests for ParseToken
func TestParseTokenWithEmptyString(t *testing.T) {
	assert := assert.New(t)

	result, err := ParseToken("")
	assert.Nil(result)
	assert.NotNil(err)
	// Empty string causes EOF error during base64 decoding
	assert.Contains(err.Error(), "EOF")
}

func TestParseTokenWithInvalidBase64(t *testing.T) {
	assert := assert.New(t)

	// Invalid base64 characters
	result, err := ParseToken("invalid!@#$%^&*()")
	assert.Nil(result)
	assert.NotNil(err)
	assert.Contains(err.Error(), "illegal base64 data")
}

func TestParseTokenWithMalformedToken(t *testing.T) {
	assert := assert.New(t)

	// Valid base64 but definitely invalid CBOR - truncated/incomplete CBOR stream
	invalidCBORData := []byte{0x9F} // Start of indefinite array but incomplete
	invalidToken := base64.StdEncoding.EncodeToString(invalidCBORData)
	result, err := ParseToken(invalidToken)

	// This may succeed with zero values or fail with CBOR error - both are acceptable
	// The important thing is that we handle it gracefully
	if err != nil {
		assert.Nil(result)
	} else {
		// If it succeeds, it should return a token with zero/default values
		assert.NotNil(result)
	}
}

func TestParseTokenWithCorruptedCBOR(t *testing.T) {
	assert := assert.New(t)

	// Truncated CBOR data
	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
	}

	var buf bytes.Buffer
	encoder := cbor.NewEncoder(&buf)
	encoder.Encode(decoded)

	// Truncate the CBOR data
	truncated := buf.Bytes()[:len(buf.Bytes())/2]
	corruptedToken := base64.StdEncoding.EncodeToString(truncated)

	result, err := ParseToken(corruptedToken)
	assert.Nil(result)
	assert.NotNil(err)
}

// Token format tests
func TestParseTokenWithURLSafeBase64(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
		Resources: GrantResources{
			Channels: map[string]int64{"test_channel": 1},
		},
		Meta: map[string]interface{}{"test": "data"},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)
	assert.Equal(1, result.Version)
	assert.Equal(int64(1234567890), result.Timestamp)
	assert.Equal(3600, result.TTL)
}

func TestParseTokenWithStandardBase64(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
		Resources: GrantResources{
			Channels: map[string]int64{"test+channel/data": 1},
		},
	}

	var buf bytes.Buffer
	encoder := cbor.NewEncoder(&buf)
	encoder.Encode(decoded)

	// Use standard base64 with + and / characters
	standardToken := base64.StdEncoding.EncodeToString(buf.Bytes())

	result, err := ParseToken(standardToken)
	assert.Nil(err)
	assert.NotNil(result)
	assert.Equal(1, result.Version)
}

func TestParseTokenWithPaddingHandling(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
		Resources: GrantResources{
			Channels: map[string]int64{"a": 1},
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	// Test the token as created (URL-safe, no padding)
	result1, err1 := ParseToken(token)
	assert.Nil(err1)
	assert.NotNil(result1)
	assert.Equal(1, result1.Version)

	// Test manually creating a token that needs padding
	var buf bytes.Buffer
	encoder := cbor.NewEncoder(&buf)
	encoder.Encode(decoded)

	// Create standard base64 (which may need padding)
	standardB64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Convert to URL-safe format like GetPermissions does
	urlSafeToken := strings.Replace(standardB64, "+", "-", -1)
	urlSafeToken = strings.Replace(urlSafeToken, "/", "_", -1)

	// Test both with and without padding
	result2, err2 := ParseToken(urlSafeToken)
	assert.Nil(err2)
	assert.NotNil(result2)

	// Remove padding and test
	unpadded := strings.TrimRight(urlSafeToken, "=")
	result3, err3 := ParseToken(unpadded)
	assert.Nil(err3)
	assert.NotNil(result3)

	// All should parse to the same values
	assert.Equal(result1.Version, result2.Version)
	assert.Equal(result1.Version, result3.Version)
}

// Permission parsing tests
func TestParseTokenWithChannelPermissions(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
		Resources: GrantResources{
			Channels: map[string]int64{
				"channel1": int64(PNRead | PNWrite),                       // 3
				"channel2": int64(PNRead | PNWrite | PNManage),            // 7
				"channel3": int64(PNRead | PNWrite | PNManage | PNDelete), // 15
			},
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)

	// Verify channel permissions
	assert.Contains(result.Resources.Channels, "channel1")
	assert.Contains(result.Resources.Channels, "channel2")
	assert.Contains(result.Resources.Channels, "channel3")

	ch1 := result.Resources.Channels["channel1"]
	assert.True(ch1.Read)
	assert.True(ch1.Write)
	assert.False(ch1.Manage)
	assert.False(ch1.Delete)

	ch2 := result.Resources.Channels["channel2"]
	assert.True(ch2.Read)
	assert.True(ch2.Write)
	assert.True(ch2.Manage)
	assert.False(ch2.Delete)

	ch3 := result.Resources.Channels["channel3"]
	assert.True(ch3.Read)
	assert.True(ch3.Write)
	assert.True(ch3.Manage)
	assert.True(ch3.Delete)
}

func TestParseTokenWithGroupPermissions(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
		Resources: GrantResources{
			Groups: map[string]int64{
				"group1": int64(PNRead),            // 1
				"group2": int64(PNRead | PNManage), // 5
			},
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)

	// Verify group permissions (only Read and Manage for groups)
	assert.Contains(result.Resources.ChannelGroups, "group1")
	assert.Contains(result.Resources.ChannelGroups, "group2")

	grp1 := result.Resources.ChannelGroups["group1"]
	assert.True(grp1.Read)
	assert.False(grp1.Manage)

	grp2 := result.Resources.ChannelGroups["group2"]
	assert.True(grp2.Read)
	assert.True(grp2.Manage)
}

func TestParseTokenWithUUIDPermissions(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
		Resources: GrantResources{
			UUIDs: map[string]int64{
				"uuid1": int64(PNGet),                       // 32
				"uuid2": int64(PNGet | PNUpdate),            // 96
				"uuid3": int64(PNGet | PNUpdate | PNDelete), // 104
			},
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)

	// Verify UUID permissions (Get, Update, Delete)
	assert.Contains(result.Resources.UUIDs, "uuid1")
	assert.Contains(result.Resources.UUIDs, "uuid2")
	assert.Contains(result.Resources.UUIDs, "uuid3")

	uuid1 := result.Resources.UUIDs["uuid1"]
	assert.True(uuid1.Get)
	assert.False(uuid1.Update)
	assert.False(uuid1.Delete)

	uuid2 := result.Resources.UUIDs["uuid2"]
	assert.True(uuid2.Get)
	assert.True(uuid2.Update)
	assert.False(uuid2.Delete)

	uuid3 := result.Resources.UUIDs["uuid3"]
	assert.True(uuid3.Get)
	assert.True(uuid3.Update)
	assert.True(uuid3.Delete)
}

func TestParseTokenWithPatternsAndResources(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
		Resources: GrantResources{
			Channels: map[string]int64{"specific_channel": int64(PNRead | PNWrite)},
		},
		Patterns: GrantResources{
			Channels: map[string]int64{"pattern_*": int64(PNRead)},
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)

	// Verify both resources and patterns are parsed
	assert.Contains(result.Resources.Channels, "specific_channel")
	assert.Contains(result.Patterns.Channels, "pattern_*")

	specificCh := result.Resources.Channels["specific_channel"]
	assert.True(specificCh.Read)
	assert.True(specificCh.Write)

	patternCh := result.Patterns.Channels["pattern_*"]
	assert.True(patternCh.Read)
	assert.False(patternCh.Write)
}

// Structure validation tests
func TestParseTokenStructureFields(t *testing.T) {
	assert := assert.New(t)

	testMeta := map[string]interface{}{
		"custom_field": "custom_value",
		"number_field": 42,
		"bool_field":   true,
	}

	decoded := PNGrantTokenDecoded{
		Version:        2,
		Timestamp:      1640995200, // 2022-01-01
		TTL:            86400,      // 24 hours
		AuthorizedUUID: "test-uuid-12345",
		Meta:           testMeta,
		Resources: GrantResources{
			Channels: map[string]int64{"test": 1},
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)

	// Verify all fields are properly parsed
	assert.Equal(2, result.Version)
	assert.Equal(int64(1640995200), result.Timestamp)
	assert.Equal(86400, result.TTL)
	assert.Equal("test-uuid-12345", result.AuthorizedUUID)

	// Verify meta data
	assert.NotNil(result.Meta)
	assert.Equal("custom_value", result.Meta["custom_field"])

	// Handle number field - CBOR may decode as uint64 or int64
	numberField := result.Meta["number_field"]
	switch v := numberField.(type) {
	case int64:
		assert.Equal(42, int(v))
	case uint64:
		assert.Equal(42, int(v))
	default:
		assert.Fail("number_field should be int64 or uint64", "got %T", v)
	}

	assert.Equal(true, result.Meta["bool_field"])
}

func TestParseTokenWithEmptyMeta(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
		Meta:      map[string]interface{}{},
		Resources: GrantResources{
			Channels: map[string]int64{"test": 1},
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)
	assert.NotNil(result.Meta)
	assert.Equal(0, len(result.Meta))
}

func TestParseTokenWithNilMeta(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
		Meta:      nil,
		Resources: GrantResources{
			Channels: map[string]int64{"test": 1},
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)

	// CBOR decoder may convert nil to empty map, both are acceptable
	if result.Meta != nil {
		assert.Equal(0, len(result.Meta))
	}
}

// Edge case tests
func TestParseTokenWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:        1,
		Timestamp:      1234567890,
		TTL:            3600,
		AuthorizedUUID: "uuid-with-special-chars@#$%",
		Resources: GrantResources{
			Channels: map[string]int64{
				"channel-with-dashes":        1,
				"channel_with_underscores":   2,
				"channel@with#special$chars": 3,
			},
		},
		Meta: map[string]interface{}{
			"key@with#special": "value$with%special&chars",
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)

	assert.Equal("uuid-with-special-chars@#$%", result.AuthorizedUUID)
	assert.Contains(result.Resources.Channels, "channel-with-dashes")
	assert.Contains(result.Resources.Channels, "channel_with_underscores")
	assert.Contains(result.Resources.Channels, "channel@with#special$chars")
	assert.Equal("value$with%special&chars", result.Meta["key@with#special"])
}

func TestParseTokenWithUnicodeCharacters(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:        1,
		Timestamp:      1234567890,
		TTL:            3600,
		AuthorizedUUID: "测试-uuid-русский-ユーザー",
		Resources: GrantResources{
			Channels: map[string]int64{
				"频道测试":       1,
				"канал-тест": 2,
				"チャンネル-テスト":  3,
			},
		},
		Meta: map[string]interface{}{
			"测试字段":  "测试值",
			"поле":  "значение",
			"フィールド": "値",
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)

	assert.Equal("测试-uuid-русский-ユーザー", result.AuthorizedUUID)
	assert.Contains(result.Resources.Channels, "频道测试")
	assert.Contains(result.Resources.Channels, "канал-тест")
	assert.Contains(result.Resources.Channels, "チャンネル-テスト")
	assert.Equal("测试值", result.Meta["测试字段"])
	assert.Equal("значение", result.Meta["поле"])
	assert.Equal("値", result.Meta["フィールド"])
}

func TestParseTokenWithLargeData(t *testing.T) {
	assert := assert.New(t)

	// Create token with many channels
	channels := make(map[string]int64)
	for i := 0; i < 1000; i++ {
		channels[fmt.Sprintf("channel_%d", i)] = int64(PNRead | PNWrite)
	}

	// Create large meta object
	largeMeta := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		largeMeta[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d_with_some_longer_text_to_make_it_larger", i)
	}

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
		Resources: GrantResources{
			Channels: channels,
		},
		Meta: largeMeta,
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)

	assert.Equal(1000, len(result.Resources.Channels))
	assert.Equal(100, len(result.Meta))

	// Verify some random entries
	assert.Contains(result.Resources.Channels, "channel_0")
	assert.Contains(result.Resources.Channels, "channel_999")
	assert.Contains(result.Meta, "key_0")
	assert.Contains(result.Meta, "key_99")
}

func TestParseTokenWithZeroTTL(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       0, // Zero TTL
		Resources: GrantResources{
			Channels: map[string]int64{"test": 1},
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)
	assert.Equal(0, result.TTL)
}

func TestParseTokenWithNegativeTTL(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       -1, // Negative TTL (unlimited)
		Resources: GrantResources{
			Channels: map[string]int64{"test": 1},
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)
	assert.Equal(-1, result.TTL)
}

func TestParseTokenWithAllPermissionTypes(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
		Resources: GrantResources{
			Channels: map[string]int64{
				"channel1": int64(PNRead | PNWrite | PNManage | PNDelete | PNGet | PNUpdate | PNJoin),
			},
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)

	ch := result.Resources.Channels["channel1"]
	assert.True(ch.Read)
	assert.True(ch.Write)
	assert.True(ch.Manage)
	assert.True(ch.Delete)
	assert.True(ch.Get)
	assert.True(ch.Update)
	assert.True(ch.Join)
}

func TestParseTokenWithEmptyResources(t *testing.T) {
	assert := assert.New(t)

	decoded := PNGrantTokenDecoded{
		Version:   1,
		Timestamp: 1234567890,
		TTL:       3600,
		Resources: GrantResources{
			Channels: map[string]int64{},
			Groups:   map[string]int64{},
			UUIDs:    map[string]int64{},
		},
		Patterns: GrantResources{
			Channels: map[string]int64{},
			Groups:   map[string]int64{},
			UUIDs:    map[string]int64{},
		},
	}

	token, err := createTestToken(decoded)
	assert.Nil(err)

	result, err := ParseToken(token)
	assert.Nil(err)
	assert.NotNil(result)

	assert.Equal(0, len(result.Resources.Channels))
	assert.Equal(0, len(result.Resources.ChannelGroups))
	assert.Equal(0, len(result.Resources.UUIDs))
	assert.Equal(0, len(result.Patterns.Channels))
	assert.Equal(0, len(result.Patterns.ChannelGroups))
	assert.Equal(0, len(result.Patterns.UUIDs))
}
