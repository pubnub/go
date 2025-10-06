package pubnub

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSendFile(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newSendFileBuilder(pn)
	if testContext {
		o = newSendFileBuilderWithContext(pn, pn.ctx)
	}

	channel := "chan"
	o.Channel(channel)
	o.QueryParam(queryParam)
	o.CustomMessageType("custom")

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(sendFilePath, pn.Config.SubscribeKey, channel),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{123, 34, 110, 97, 109, 101, 34, 58, 34, 34, 125}, body)

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
		assert.Equal("custom", u.Get("custom_message_type"))
	}

}

func TestSendFile(t *testing.T) {
	AssertSendFile(t, true, false)
}

func TestSendFileContext(t *testing.T) {
	AssertSendFile(t, true, true)
}

func TestSendFileResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNSendFileResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: error unmarshalling response: {s}", err.Error())
}

func TestSendFileCustomMessageTypeValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pn, pn.ctx)
	opts.CustomMessageType = "custom-message_type"
	assert.True(opts.isCustomMessageTypeCorrect())
	opts.CustomMessageType = "a"
	assert.False(opts.isCustomMessageTypeCorrect())
	opts.CustomMessageType = "!@#$%^&*("
	assert.False(opts.isCustomMessageTypeCorrect())
}

// ===========================
// Validation Tests
// ===========================

func TestSendFileValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	pubnub.Config.SubscribeKey = ""
	opts := newSendFileOpts(pubnub, pubnub.ctx)

	err := opts.validate()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestSendFileValidateChannel(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Channel = ""

	err := opts.validate()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestSendFileValidateName(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel"
	opts.Name = ""

	err := opts.validate()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing File Name")
}

func TestSendFileValidateCustomMessageType(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel"
	opts.Name = "test.txt"
	opts.File = &os.File{} // Mock file
	opts.CustomMessageType = "!"

	err := opts.validate()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Invalid CustomMessageType")
}

// ===========================
// HTTP Method and Operation Tests
// ===========================

func TestSendFileHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)

	method := opts.httpMethod()

	assert.Equal("POST", method)
}

func TestSendFileOperationType(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)

	opType := opts.operationType()

	assert.Equal(PNSendFileOperation, opType)
}

func TestSendFileIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)

	authRequired := opts.isAuthRequired()

	assert.True(authRequired)
}

func TestSendFileRequestTimeout(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)

	timeout := opts.requestTimeout()

	assert.Equal(pubnub.Config.NonSubscribeRequestTimeout, timeout)
}

// ===========================
// Builder Pattern Tests
// ===========================

func TestSendFileBuilderChannel(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	result := builder.Channel("test-channel")

	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileBuilderName(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	result := builder.Name("test.txt")

	assert.Equal("test.txt", builder.opts.Name)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileBuilderMessage(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	result := builder.Message("test message")

	assert.Equal("test message", builder.opts.Message)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileBuilderFile(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)
	tempFile, _ := os.CreateTemp("", "test")
	defer os.Remove(tempFile.Name())

	result := builder.File(tempFile)

	assert.Equal(tempFile, builder.opts.File)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileBuilderCipherKey(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	result := builder.CipherKey("my-cipher-key")

	assert.Equal("my-cipher-key", builder.opts.CipherKey)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileBuilderTTL(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	result := builder.TTL(24)

	assert.Equal(24, builder.opts.TTL)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileBuilderMeta(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)
	meta := map[string]interface{}{"key": "value"}

	result := builder.Meta(meta)

	assert.Equal(meta, builder.opts.Meta)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileBuilderShouldStore(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	result := builder.ShouldStore(true)

	assert.True(builder.opts.ShouldStore)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileBuilderCustomMessageType(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	result := builder.CustomMessageType("custom-type")

	assert.Equal("custom-type", builder.opts.CustomMessageType)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)
	queryParam := map[string]string{"q1": "v1", "q2": "v2"}

	result := builder.QueryParam(queryParam)

	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)
	transport := &http.Transport{}

	result := builder.Transport(transport)

	assert.Equal(transport, builder.opts.Transport)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)
	tempFile, _ := os.CreateTemp("", "test")
	defer os.Remove(tempFile.Name())

	result := builder.
		Channel("test-channel").
		Name("test.txt").
		Message("test message").
		File(tempFile).
		CipherKey("my-cipher-key").
		TTL(24).
		Meta(map[string]interface{}{"key": "value"}).
		ShouldStore(true).
		CustomMessageType("custom-type").
		QueryParam(map[string]string{"q1": "v1"}).
		Transport(&http.Transport{})

	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("test.txt", builder.opts.Name)
	assert.Equal("test message", builder.opts.Message)
	assert.Equal(tempFile, builder.opts.File)
	assert.Equal("my-cipher-key", builder.opts.CipherKey)
	assert.Equal(24, builder.opts.TTL)
	assert.Equal(map[string]interface{}{"key": "value"}, builder.opts.Meta)
	assert.True(builder.opts.ShouldStore)
	assert.Equal("custom-type", builder.opts.CustomMessageType)
	assert.Equal(map[string]string{"q1": "v1"}, builder.opts.QueryParam)
	assert.NotNil(builder.opts.Transport)
	assert.Equal(builder, result) // Fluent interface
}

// ===========================
// URL/Path Building Tests
// ===========================

func TestSendFileBuildPath(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel"

	path, err := opts.buildPath()

	assert.Nil(err)
	expectedPath := fmt.Sprintf("/v1/files/%s/channels/test-channel/generate-upload-url", pubnub.Config.SubscribeKey)
	assert.Equal(expectedPath, path)
}

func TestSendFileBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel@#$%"

	path, err := opts.buildPath()

	assert.Nil(err)
	expectedPath := fmt.Sprintf("/v1/files/%s/channels/test-channel@#$%%/generate-upload-url", pubnub.Config.SubscribeKey)
	assert.Equal(expectedPath, path)
}

func TestSendFileBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "频道"

	path, err := opts.buildPath()

	assert.Nil(err)
	expectedPath := fmt.Sprintf("/v1/files/%s/channels/频道/generate-upload-url", pubnub.Config.SubscribeKey)
	assert.Equal(expectedPath, path)
}

func TestSendFileBuildPathWithSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	pubnub.Config.SubscribeKey = "special-sub-key"
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel"

	path, err := opts.buildPath()

	assert.Nil(err)
	expectedPath := "/v1/files/special-sub-key/channels/test-channel/generate-upload-url"
	assert.Equal(expectedPath, path)
}

// ===========================
// JSON Body Building Tests
// ===========================

func TestSendFileBuildBodyBasic(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Name = "test.txt"

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Contains(string(body), `"name":"test.txt"`)
}

func TestSendFileBuildBodyWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Name = "test@#$%.txt"

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Contains(string(body), `"name":"test@#$%.txt"`)
}

func TestSendFileBuildBodyWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Name = "文件.txt"

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Contains(string(body), `"name":"文件.txt"`)
}

func TestSendFileBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Name = ""

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Contains(string(body), `"name":""`)
}

func TestSendFileBuildBodyWithQuotes(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Name = `test"file.txt`

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Contains(string(body), `"name":"test\"file.txt"`)
}

// ===========================
// Query Parameter Tests
// ===========================

func TestSendFileBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.NotNil(query)
	assert.Equal(pubnub.Config.UUID, query.Get("uuid"))
}

func TestSendFileBuildQueryWithCustomMessageType(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.CustomMessageType = "custom-type"

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.Equal("custom-type", query.Get("custom_message_type"))
}

func TestSendFileBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.QueryParam = map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.Equal("value1", query.Get("param1"))
	assert.Equal("value2", query.Get("param2"))
}

func TestSendFileBuildQueryWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.QueryParam = map[string]string{
		"special_chars": "value@#$%&*",
	}

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.Equal("value%40%23%24%25%26%2A", query.Get("special_chars"))
}

func TestSendFileBuildQueryWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.QueryParam = map[string]string{
		"unicode": "文件上传",
	}

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.Equal("%E6%96%87%E4%BB%B6%E4%B8%8A%E4%BC%A0", query.Get("unicode"))
}

func TestSendFileBuildQueryCombinations(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.CustomMessageType = "custom-type"
	opts.QueryParam = map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.Equal("custom-type", query.Get("custom_message_type"))
	assert.Equal("value1", query.Get("param1"))
	assert.Equal("value2", query.Get("param2"))
	assert.Equal(pubnub.Config.UUID, query.Get("uuid"))
}

// ===========================
// File Handling Tests
// ===========================

func TestSendFileWithFile(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)
	tempFile, _ := os.CreateTemp("", "test")
	defer os.Remove(tempFile.Name())

	builder.File(tempFile)

	assert.Equal(tempFile, builder.opts.File)
}

func TestSendFileWithCipherKey(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	builder.CipherKey("my-cipher-key")

	assert.Equal("my-cipher-key", builder.opts.CipherKey)
}

func TestSendFileWithEmptyCipherKey(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	builder.CipherKey("")

	assert.Equal("", builder.opts.CipherKey)
}

// ===========================
// Message and Metadata Tests
// ===========================

func TestSendFileWithMessage(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	builder.Message("test message")

	assert.Equal("test message", builder.opts.Message)
}

func TestSendFileWithTTL(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	builder.TTL(24)

	assert.Equal(24, builder.opts.TTL)
}

func TestSendFileWithMeta(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)
	meta := map[string]interface{}{"key": "value"}

	builder.Meta(meta)

	assert.Equal(meta, builder.opts.Meta)
}

func TestSendFileWithShouldStore(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	builder.ShouldStore(true)

	assert.True(builder.opts.ShouldStore)
}

// ===========================
// Edge Case Tests
// ===========================

func TestSendFileWithLongFileName(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	longName := strings.Repeat("a", 1000) + ".txt"
	opts.Name = longName

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Contains(string(body), longName)
}

func TestSendFileWithLongMessage(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)
	longMessage := strings.Repeat("This is a very long message. ", 1000)

	builder.Message(longMessage)

	assert.Equal(longMessage, builder.opts.Message)
}

func TestSendFileWithLongChannel(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	longChannel := strings.Repeat("a", 1000)
	opts.Channel = longChannel

	path, err := opts.buildPath()

	assert.Nil(err)
	assert.Contains(path, longChannel)
}

func TestSendFileWithComplexMeta(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)
	complexMeta := map[string]interface{}{
		"string":  "value",
		"number":  123,
		"boolean": true,
		"array":   []string{"a", "b", "c"},
		"object":  map[string]interface{}{"nested": "value"},
	}

	builder.Meta(complexMeta)

	assert.Equal(complexMeta, builder.opts.Meta)
}

func TestSendFileWithUnicodeFileName(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Name = "文件名.txt"

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Contains(string(body), "文件名.txt")
}

func TestSendFileWithUnicodeMessage(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	builder.Message("这是一个测试消息")

	assert.Equal("这是一个测试消息", builder.opts.Message)
}

func TestSendFileWithUnicodeChannel(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "频道"

	path, err := opts.buildPath()

	assert.Nil(err)
	assert.Contains(path, "频道")
}

func TestSendFileWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)
	opts.Name = "file@#$%.txt"
	opts.Channel = "channel@#$%"

	body, err := opts.buildBody()
	path, pathErr := opts.buildPath()

	assert.Nil(err)
	assert.Nil(pathErr)
	assert.Contains(string(body), "file@#$%.txt")
	assert.Contains(path, "channel@#$%")
}

func TestSendFileWithMaxTTL(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	builder.TTL(2147483647) // Max int32

	assert.Equal(2147483647, builder.opts.TTL)
}

func TestSendFileWithZeroTTL(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	builder.TTL(0)

	assert.Equal(0, builder.opts.TTL)
}

// ===========================
// Error Scenario Tests
// ===========================

func TestSendFileValidationErrors(t *testing.T) {
	testCases := []struct {
		name              string
		subscribeKey      string
		channel           string
		fileName          string
		customMessageType string
		expectedError     string
	}{
		{
			name:          "Missing subscribe key",
			subscribeKey:  "",
			channel:       "test-channel",
			fileName:      "test.txt",
			expectedError: "Missing Subscribe Key",
		},
		{
			name:          "Missing channel",
			subscribeKey:  "demo",
			channel:       "",
			fileName:      "test.txt",
			expectedError: "Missing Channel",
		},
		{
			name:          "Missing file name",
			subscribeKey:  "demo",
			channel:       "test-channel",
			fileName:      "",
			expectedError: "Missing File Name",
		},
		{
			name:              "Invalid custom message type",
			subscribeKey:      "demo",
			channel:           "test-channel",
			fileName:          "test.txt",
			customMessageType: "!@#",
			expectedError:     "Invalid CustomMessageType",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pubnub := NewPubNub(NewDemoConfig())
			pubnub.Config.SubscribeKey = tc.subscribeKey
			opts := newSendFileOpts(pubnub, pubnub.ctx)
			opts.Channel = tc.channel
			opts.Name = tc.fileName
			opts.CustomMessageType = tc.customMessageType
			opts.File = &os.File{} // Mock file for validation

			err := opts.validate()

			assert.NotNil(err)
			assert.Contains(err.Error(), tc.expectedError)
		})
	}
}

func TestSendFileBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	assert.Equal("", builder.opts.Channel)
	assert.Equal("", builder.opts.Name)
	assert.Equal("", builder.opts.Message)
	assert.Nil(builder.opts.File)
	assert.Equal("", builder.opts.CipherKey)
	assert.Equal(0, builder.opts.TTL)
	assert.Nil(builder.opts.Meta)
	assert.False(builder.opts.ShouldStore)
	assert.Equal("", builder.opts.CustomMessageType)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestSendFileBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	// Check that all builder methods exist and return the builder
	assert.NotNil(builder.Channel("test"))
	assert.NotNil(builder.Name("test"))
	assert.NotNil(builder.Message("test"))
	assert.NotNil(builder.File(nil))
	assert.NotNil(builder.CipherKey("test"))
	assert.NotNil(builder.TTL(1))
	assert.NotNil(builder.Meta(nil))
	assert.NotNil(builder.ShouldStore(true))
	assert.NotNil(builder.CustomMessageType("test"))
	assert.NotNil(builder.QueryParam(nil))
	assert.NotNil(builder.Transport(nil))
}

func TestSendFileNewBuilderWithContext(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())

	builder := newSendFileBuilderWithContext(pubnub, pubnub.ctx)

	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pubnub, builder.opts.pubnub)
	assert.Equal(pubnub.ctx, builder.opts.ctx)
}

func TestSendFileIsCustomMessageTypeCorrect(t *testing.T) {
	testCases := []struct {
		name     string
		msgType  string
		expected bool
	}{
		{"Valid message type", "custom-message_type", true},
		{"Empty message type", "", true},
		{"Too short", "a", false},
		{"Special characters", "!@#$%", false},
		{"Invalid with numbers", "type123", false},
		{"Valid with dashes", "custom-type", true},
		{"Valid with underscores", "custom_type", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pubnub := NewPubNub(NewDemoConfig())
			opts := newSendFileOpts(pubnub, pubnub.ctx)
			opts.CustomMessageType = tc.msgType

			result := opts.isCustomMessageTypeCorrect()

			assert.Equal(tc.expected, result)
		})
	}
}

// ===========================
// UseRawText Tests
// ===========================

func TestSendFileBuilderUseRawText(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)

	// Test UseRawText true
	result := builder.UseRawText(true)
	assert.True(builder.opts.UseRawText)
	assert.Equal(builder, result) // Fluent interface

	// Test UseRawText false
	result = builder.UseRawText(false)
	assert.False(builder.opts.UseRawText)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileUseRawTextDefaults(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newSendFileOpts(pubnub, pubnub.ctx)

	// Test default value
	assert.False(opts.UseRawText)
}

func TestSendFileUseRawTextMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)
	tempFile, _ := os.CreateTemp("", "test")
	defer os.Remove(tempFile.Name())

	// Test method chaining with UseRawText
	result := builder.
		Channel("test-channel").
		Name("test.txt").
		Message("test message").
		File(tempFile).
		UseRawText(true).
		ShouldStore(true).
		TTL(24)

	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("test.txt", builder.opts.Name)
	assert.Equal("test message", builder.opts.Message)
	assert.Equal(tempFile, builder.opts.File)
	assert.True(builder.opts.UseRawText)
	assert.True(builder.opts.ShouldStore)
	assert.Equal(24, builder.opts.TTL)
	assert.Equal(builder, result) // Fluent interface
}

func TestSendFileUseRawTextWithAllParameters(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newSendFileBuilder(pubnub)
	tempFile, _ := os.CreateTemp("", "test")
	defer os.Remove(tempFile.Name())

	// Test UseRawText with all other parameters
	queryParam := map[string]string{"param1": "value1"}
	meta := map[string]interface{}{"key": "value"}

	result := builder.
		Channel("test-channel").
		Name("test.txt").
		Message("test message").
		File(tempFile).
		CipherKey("my-cipher-key").
		TTL(24).
		Meta(meta).
		ShouldStore(true).
		CustomMessageType("custom-type").
		QueryParam(queryParam).
		UseRawText(true).
		Transport(&http.Transport{})

	// Verify all parameters are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("test.txt", builder.opts.Name)
	assert.Equal("test message", builder.opts.Message)
	assert.Equal(tempFile, builder.opts.File)
	assert.Equal("my-cipher-key", builder.opts.CipherKey)
	assert.Equal(24, builder.opts.TTL)
	assert.Equal(meta, builder.opts.Meta)
	assert.True(builder.opts.ShouldStore)
	assert.Equal("custom-type", builder.opts.CustomMessageType)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.True(builder.opts.UseRawText)
	assert.NotNil(builder.opts.Transport)
	assert.Equal(builder, result) // Fluent interface
}
