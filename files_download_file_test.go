package pubnub

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertDownloadFile(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newDownloadFileBuilder(pn)
	if testContext {
		o = newDownloadFileBuilderWithContext(pn, pn.ctx)
	}

	channel := "test-channel"
	fileID := "file123"
	fileName := "test.txt"
	cipherKey := "test-cipher"

	o.Channel(channel)
	o.ID(fileID)
	o.Name(fileName)
	o.CipherKey(cipherKey)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/files/%s/channels/%s/files/%s/%s", pn.Config.SubscribeKey, channel, fileID, fileName),
		path, []int{})

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestDownloadFile(t *testing.T) {
	AssertDownloadFile(t, true, false)
}

func TestDownloadFileContext(t *testing.T) {
	AssertDownloadFile(t, true, true)
}

func TestDownloadFileQueryParam(t *testing.T) {
	AssertDownloadFile(t, false, false)
}

func TestDownloadFileValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newDownloadFileOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.ID = "id"
	opts.Name = "name"

	assert.Equal("pubnub/validation: pubnub: Download File: Missing Subscribe Key", opts.validate().Error())
}

func TestDownloadFileValidateChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pn, pn.ctx)
	opts.ID = "id"
	opts.Name = "name"

	assert.Equal("pubnub/validation: pubnub: Download File: Missing Channel", opts.validate().Error())
}

func TestDownloadFileValidateFileID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.Name = "name"

	assert.Equal("pubnub/validation: pubnub: Download File: Missing File ID", opts.validate().Error())
}

func TestDownloadFileValidateFileName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.ID = "id"

	assert.Equal("pubnub/validation: pubnub: Download File: Missing File Name", opts.validate().Error())
}

func TestDownloadFileValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.ID = "id"
	opts.Name = "name"

	assert.Nil(opts.validate())
}

func TestNewDownloadFileBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newDownloadFileBuilder(pubnub)
	o.Channel("ch")
	o.ID("id")
	o.Name("name")

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/files/%s/channels/%s/files/%s/%s", pubnub.Config.SubscribeKey, "ch", "id", "name"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestDownloadFileBuilderSetters(t *testing.T) {
	assert := assert.New(t)

	o := newDownloadFileBuilder(pubnub)

	// Test fluent interface
	result := o.Channel("test-channel").
		ID("test-id").
		Name("test-name.txt").
		CipherKey("test-cipher")

	assert.Equal(o, result) // Should return same instance for chaining
	assert.Equal("test-channel", o.opts.Channel)
	assert.Equal("test-id", o.opts.ID)
	assert.Equal("test-name.txt", o.opts.Name)
	assert.Equal("test-cipher", o.opts.CipherKey)
}

// ===========================
// HTTP Method and Operation Tests
// ===========================

func TestDownloadFileHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)

	method := opts.httpMethod()

	assert.Equal("GET", method)
}

func TestDownloadFileOperationType(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)

	opType := opts.operationType()

	assert.Equal(PNDownloadFileOperation, opType)
}

func TestDownloadFileIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)

	authRequired := opts.isAuthRequired()

	assert.True(authRequired)
}

func TestDownloadFileRequestTimeout(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)

	timeout := opts.requestTimeout()

	assert.Equal(pubnub.Config.NonSubscribeRequestTimeout, timeout)
}

func TestDownloadFileBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Empty(body)
}

// ===========================
// Enhanced Builder Pattern Tests
// ===========================

func TestDownloadFileBuilderChannel(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDownloadFileBuilder(pubnub)

	result := builder.Channel("test-channel")

	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal(builder, result) // Fluent interface
}

func TestDownloadFileBuilderID(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDownloadFileBuilder(pubnub)

	result := builder.ID("test-id")

	assert.Equal("test-id", builder.opts.ID)
	assert.Equal(builder, result) // Fluent interface
}

func TestDownloadFileBuilderName(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDownloadFileBuilder(pubnub)

	result := builder.Name("test.txt")

	assert.Equal("test.txt", builder.opts.Name)
	assert.Equal(builder, result) // Fluent interface
}

func TestDownloadFileBuilderCipherKey(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDownloadFileBuilder(pubnub)

	result := builder.CipherKey("my-cipher-key")

	assert.Equal("my-cipher-key", builder.opts.CipherKey)
	assert.Equal(builder, result) // Fluent interface
}

func TestDownloadFileBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDownloadFileBuilder(pubnub)
	queryParam := map[string]string{"q1": "v1", "q2": "v2"}

	result := builder.QueryParam(queryParam)

	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(builder, result) // Fluent interface
}

func TestDownloadFileBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDownloadFileBuilder(pubnub)
	transport := &http.Transport{}

	result := builder.Transport(transport)

	assert.Equal(transport, builder.opts.Transport)
	assert.Equal(builder, result) // Fluent interface
}

func TestDownloadFileBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDownloadFileBuilder(pubnub)

	result := builder.
		Channel("test-channel").
		ID("test-id").
		Name("test.txt").
		CipherKey("my-cipher-key").
		QueryParam(map[string]string{"q1": "v1"}).
		Transport(&http.Transport{})

	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("test-id", builder.opts.ID)
	assert.Equal("test.txt", builder.opts.Name)
	assert.Equal("my-cipher-key", builder.opts.CipherKey)
	assert.Equal(map[string]string{"q1": "v1"}, builder.opts.QueryParam)
	assert.NotNil(builder.opts.Transport)
	assert.Equal(builder, result) // Fluent interface
}

// ===========================
// URL/Path Building Tests
// ===========================

func TestDownloadFileBuildPath(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel"
	opts.ID = "test-id"
	opts.Name = "test.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	expectedPath := fmt.Sprintf("/v1/files/%s/channels/test-channel/files/test-id/test.txt", pubnub.Config.SubscribeKey)
	assert.Equal(expectedPath, path)
}

func TestDownloadFileBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel@#$%"
	opts.ID = "test-id@#$%"
	opts.Name = "test@#$%.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	expectedPath := fmt.Sprintf("/v1/files/%s/channels/test-channel@#$%%/files/test-id@#$%%/test@#$%%.txt", pubnub.Config.SubscribeKey)
	assert.Equal(expectedPath, path)
}

func TestDownloadFileBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "频道"
	opts.ID = "文件ID"
	opts.Name = "文件.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	expectedPath := fmt.Sprintf("/v1/files/%s/channels/频道/files/文件ID/文件.txt", pubnub.Config.SubscribeKey)
	assert.Equal(expectedPath, path)
}

func TestDownloadFileBuildPathWithSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	pubnub.Config.SubscribeKey = "special-sub-key"
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel"
	opts.ID = "test-id"
	opts.Name = "test.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	expectedPath := "/v1/files/special-sub-key/channels/test-channel/files/test-id/test.txt"
	assert.Equal(expectedPath, path)
}

// ===========================
// Query Parameter Tests
// ===========================

func TestDownloadFileBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.NotNil(query)
	assert.Equal(pubnub.Config.UUID, query.Get("uuid"))
}

func TestDownloadFileBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	opts.QueryParam = map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.Equal("value1", query.Get("param1"))
	assert.Equal("value2", query.Get("param2"))
}

func TestDownloadFileBuildQueryWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	opts.QueryParam = map[string]string{
		"special_chars": "value@#$%&*",
	}

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.Equal("value%40%23%24%25%26%2A", query.Get("special_chars"))
}

func TestDownloadFileBuildQueryWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	opts.QueryParam = map[string]string{
		"unicode": "文件下载",
	}

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.Equal("%E6%96%87%E4%BB%B6%E4%B8%8B%E8%BD%BD", query.Get("unicode"))
}

// ===========================
// Encryption/Cipher Key Tests
// ===========================

func TestDownloadFileWithCipherKey(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDownloadFileBuilder(pubnub)

	builder.CipherKey("my-cipher-key")

	assert.Equal("my-cipher-key", builder.opts.CipherKey)
}

func TestDownloadFileWithEmptyCipherKey(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDownloadFileBuilder(pubnub)

	builder.CipherKey("")

	assert.Equal("", builder.opts.CipherKey)
}

func TestDownloadFileWithLongCipherKey(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDownloadFileBuilder(pubnub)
	longCipherKey := strings.Repeat("a", 256)

	builder.CipherKey(longCipherKey)

	assert.Equal(longCipherKey, builder.opts.CipherKey)
}

// ===========================
// Edge Case Tests
// ===========================

func TestDownloadFileWithLongValues(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	longChannel := strings.Repeat("a", 1000)
	longID := strings.Repeat("b", 1000)
	longName := strings.Repeat("c", 1000) + ".txt"

	opts.Channel = longChannel
	opts.ID = longID
	opts.Name = longName

	path, err := opts.buildPath()

	assert.Nil(err)
	assert.Contains(path, longChannel)
	assert.Contains(path, longID)
	assert.Contains(path, longName)
}

func TestDownloadFileWithUnicodeValues(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "测试频道"
	opts.ID = "文件标识符"
	opts.Name = "测试文件.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	assert.Contains(path, "测试频道")
	assert.Contains(path, "文件标识符")
	assert.Contains(path, "测试文件.txt")
}

func TestDownloadFileWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "channel@#$%"
	opts.ID = "id!@#$%^&*()"
	opts.Name = "file@#$%.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	assert.Contains(path, "channel@#$%")
	assert.Contains(path, "id!@#$%^&*()")
	assert.Contains(path, "file@#$%.txt")
}

func TestDownloadFileWithSpaces(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "my channel"
	opts.ID = "my file id"
	opts.Name = "my file.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	assert.Contains(path, "my channel")
	assert.Contains(path, "my file id")
	assert.Contains(path, "my file.txt")
}

func TestDownloadFileWithDots(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "channel.name"
	opts.ID = "file.id"
	opts.Name = "file.name.ext"

	path, err := opts.buildPath()

	assert.Nil(err)
	assert.Contains(path, "channel.name")
	assert.Contains(path, "file.id")
	assert.Contains(path, "file.name.ext")
}

func TestDownloadFileWithQueryParamCombinations(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)
	opts.QueryParam = map[string]string{
		"param1":  "value1",
		"param2":  "value2",
		"special": "value@#$%",
		"unicode": "测试",
		"empty":   "",
		"spaces":  "value with spaces",
	}

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.Equal("value1", query.Get("param1"))
	assert.Equal("value2", query.Get("param2"))
	assert.Equal("value%40%23%24%25", query.Get("special"))
	assert.Equal("%E6%B5%8B%E8%AF%95", query.Get("unicode"))
	assert.Equal("", query.Get("empty"))
	assert.Equal("value%20with%20spaces", query.Get("spaces"))
}

// ===========================
// Error Scenario Tests
// ===========================

func TestDownloadFileValidationErrors(t *testing.T) {
	testCases := []struct {
		name          string
		subscribeKey  string
		channel       string
		id            string
		fileName      string
		expectedError string
	}{
		{
			name:          "Missing subscribe key",
			subscribeKey:  "",
			channel:       "test-channel",
			id:            "test-id",
			fileName:      "test.txt",
			expectedError: "Missing Subscribe Key",
		},
		{
			name:          "Missing channel",
			subscribeKey:  "demo",
			channel:       "",
			id:            "test-id",
			fileName:      "test.txt",
			expectedError: "Missing Channel",
		},
		{
			name:          "Missing file ID",
			subscribeKey:  "demo",
			channel:       "test-channel",
			id:            "",
			fileName:      "test.txt",
			expectedError: "Missing File ID",
		},
		{
			name:          "Missing file name",
			subscribeKey:  "demo",
			channel:       "test-channel",
			id:            "test-id",
			fileName:      "",
			expectedError: "Missing File Name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pubnub := NewPubNub(NewDemoConfig())
			pubnub.Config.SubscribeKey = tc.subscribeKey
			opts := newDownloadFileOpts(pubnub, pubnub.ctx)
			opts.Channel = tc.channel
			opts.ID = tc.id
			opts.Name = tc.fileName

			err := opts.validate()

			assert.NotNil(err)
			assert.Contains(err.Error(), tc.expectedError)
		})
	}
}

func TestDownloadFileBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDownloadFileBuilder(pubnub)

	assert.Equal("", builder.opts.Channel)
	assert.Equal("", builder.opts.ID)
	assert.Equal("", builder.opts.Name)
	assert.Equal("", builder.opts.CipherKey)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestDownloadFileBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDownloadFileBuilder(pubnub)

	// Check that all builder methods exist and return the builder
	assert.NotNil(builder.Channel("test"))
	assert.NotNil(builder.ID("test"))
	assert.NotNil(builder.Name("test"))
	assert.NotNil(builder.CipherKey("test"))
	assert.NotNil(builder.QueryParam(nil))
	assert.NotNil(builder.Transport(nil))
}

func TestDownloadFileNewBuilderWithContext(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())

	builder := newDownloadFileBuilderWithContext(pubnub, pubnub.ctx)

	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pubnub, builder.opts.pubnub)
	assert.Equal(pubnub.ctx, builder.opts.ctx)
}

// ===========================
// GET-Specific Tests
// ===========================

func TestDownloadFileMethodCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)

	// GET requests should have empty body
	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	// GET method should be properly set
	method := opts.httpMethod()
	assert.Equal("GET", method)

	// Should require authentication
	authRequired := opts.isAuthRequired()
	assert.True(authRequired)
}

func TestDownloadFileOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)

	// Operation type should be correct
	opType := opts.operationType()
	assert.Equal(PNDownloadFileOperation, opType)

	// Timeout should be non-subscribe timeout
	timeout := opts.requestTimeout()
	assert.Equal(pubnub.Config.NonSubscribeRequestTimeout, timeout)

	// Connect timeout should be set
	connectTimeout := opts.connectTimeout()
	assert.Equal(pubnub.Config.ConnectTimeout, connectTimeout)
}

// ===========================
// Response Type Tests
// ===========================

func TestDownloadFileResponseStructure(t *testing.T) {
	assert := assert.New(t)

	// Test response structure
	resp := &PNDownloadFileResponse{
		Status: 200,
		File:   nil,
	}

	assert.Equal(200, resp.Status)
	assert.Nil(resp.File)
}

func TestDownloadFileNewResponseFunction(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDownloadFileOpts(pubnub, pubnub.ctx)

	resp, status, err := newPNDownloadFileResponse([]byte(`{}`), opts, StatusResponse{StatusCode: 200})

	assert.Nil(err)
	assert.NotNil(resp)
	assert.Equal(200, status.StatusCode)
}
