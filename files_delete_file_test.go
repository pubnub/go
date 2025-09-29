package pubnub

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertDeleteFile(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newDeleteFileBuilder(pn)
	if testContext {
		o = newDeleteFileBuilderWithContext(pn, pn.ctx)
	}

	channel := "chan"
	id := "fileid"
	name := "filename"
	o.Channel(channel)
	o.QueryParam(queryParam)
	o.ID(id)
	o.Name(name)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(deleteFilePath, pn.Config.SubscribeKey, channel, id, name),
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

func TestDeleteFile(t *testing.T) {
	AssertDeleteFile(t, true, false)
}

func TestDeleteFileContext(t *testing.T) {
	AssertDeleteFile(t, true, true)
}

func TestDeleteFileResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNDeleteFileResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: error unmarshalling response: {s}", err.Error())
}

func TestDeleteFileResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200}`)

	_, s, err := newPNDeleteFileResponse(jsonBytes, opts, StatusResponse{StatusCode: 200})
	assert.Equal(200, s.StatusCode)

	assert.Nil(err)
}

// ===========================
// Validation Tests
// ===========================

func TestDeleteFileValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	pubnub.Config.SubscribeKey = ""
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)

	err := opts.validate()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestDeleteFileValidateChannel(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.Channel = ""

	err := opts.validate()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestDeleteFileValidateName(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel"
	opts.Name = ""

	err := opts.validate()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing File Name")
}

func TestDeleteFileValidateID(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel"
	opts.Name = "test.txt"
	opts.ID = ""

	err := opts.validate()

	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing File ID")
}

// ===========================
// HTTP Method and Operation Tests
// ===========================

func TestDeleteFileHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)

	method := opts.httpMethod()

	assert.Equal("DELETE", method)
}

func TestDeleteFileOperationType(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)

	opType := opts.operationType()

	assert.Equal(PNDeleteFileOperation, opType)
}

func TestDeleteFileIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)

	authRequired := opts.isAuthRequired()

	assert.True(authRequired)
}

func TestDeleteFileRequestTimeout(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)

	timeout := opts.requestTimeout()

	assert.Equal(pubnub.Config.NonSubscribeRequestTimeout, timeout)
}

func TestDeleteFileBuildBodyEmpty(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Empty(body)
}

// ===========================
// Builder Pattern Tests
// ===========================

func TestDeleteFileBuilderChannel(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDeleteFileBuilder(pubnub)

	result := builder.Channel("test-channel")

	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal(builder, result) // Fluent interface
}

func TestDeleteFileBuilderID(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDeleteFileBuilder(pubnub)

	result := builder.ID("test-id")

	assert.Equal("test-id", builder.opts.ID)
	assert.Equal(builder, result) // Fluent interface
}

func TestDeleteFileBuilderName(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDeleteFileBuilder(pubnub)

	result := builder.Name("test.txt")

	assert.Equal("test.txt", builder.opts.Name)
	assert.Equal(builder, result) // Fluent interface
}

func TestDeleteFileBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDeleteFileBuilder(pubnub)
	queryParam := map[string]string{"q1": "v1", "q2": "v2"}

	result := builder.QueryParam(queryParam)

	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(builder, result) // Fluent interface
}

func TestDeleteFileBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDeleteFileBuilder(pubnub)
	transport := &http.Transport{}

	result := builder.Transport(transport)

	assert.Equal(transport, builder.opts.Transport)
	assert.Equal(builder, result) // Fluent interface
}

func TestDeleteFileBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDeleteFileBuilder(pubnub)

	result := builder.
		Channel("test-channel").
		ID("test-id").
		Name("test.txt").
		QueryParam(map[string]string{"q1": "v1"}).
		Transport(&http.Transport{})

	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("test-id", builder.opts.ID)
	assert.Equal("test.txt", builder.opts.Name)
	assert.Equal(map[string]string{"q1": "v1"}, builder.opts.QueryParam)
	assert.NotNil(builder.opts.Transport)
	assert.Equal(builder, result) // Fluent interface
}

// ===========================
// URL/Path Building Tests
// ===========================

func TestDeleteFileBuildPath(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel"
	opts.ID = "test-id"
	opts.Name = "test.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	expectedPath := fmt.Sprintf("/v1/files/%s/channels/test-channel/files/test-id/test.txt", pubnub.Config.SubscribeKey)
	assert.Equal(expectedPath, path)
}

func TestDeleteFileBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel@#$%"
	opts.ID = "test-id@#$%"
	opts.Name = "test@#$%.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	expectedPath := fmt.Sprintf("/v1/files/%s/channels/test-channel@#$%%/files/test-id@#$%%/test@#$%%.txt", pubnub.Config.SubscribeKey)
	assert.Equal(expectedPath, path)
}

func TestDeleteFileBuildPathWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "频道"
	opts.ID = "文件ID"
	opts.Name = "文件.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	expectedPath := fmt.Sprintf("/v1/files/%s/channels/频道/files/文件ID/文件.txt", pubnub.Config.SubscribeKey)
	assert.Equal(expectedPath, path)
}

func TestDeleteFileBuildPathWithSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	pubnub.Config.SubscribeKey = "special-sub-key"
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
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

func TestDeleteFileBuildQueryBasic(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.NotNil(query)
	assert.Equal(pubnub.Config.UUID, query.Get("uuid"))
}

func TestDeleteFileBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.QueryParam = map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.Equal("value1", query.Get("param1"))
	assert.Equal("value2", query.Get("param2"))
}

func TestDeleteFileBuildQueryWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.QueryParam = map[string]string{
		"special_chars": "value@#$%&*",
	}

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.Equal("value%40%23%24%25%26%2A", query.Get("special_chars"))
}

func TestDeleteFileBuildQueryWithUnicode(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.QueryParam = map[string]string{
		"unicode": "文件删除",
	}

	query, err := opts.buildQuery()

	assert.Nil(err)
	assert.Equal("%E6%96%87%E4%BB%B6%E5%88%A0%E9%99%A4", query.Get("unicode"))
}

// ===========================
// DELETE-Specific Tests
// ===========================

func TestDeleteFileMethodCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)

	// DELETE requests should have empty body
	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	// DELETE method should be properly set
	method := opts.httpMethod()
	assert.Equal("DELETE", method)

	// Should require authentication
	authRequired := opts.isAuthRequired()
	assert.True(authRequired)
}

func TestDeleteFileOperationCharacteristics(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)

	// Operation type should be correct
	opType := opts.operationType()
	assert.Equal(PNDeleteFileOperation, opType)

	// Timeout should be non-subscribe timeout
	timeout := opts.requestTimeout()
	assert.Equal(pubnub.Config.NonSubscribeRequestTimeout, timeout)

	// Connect timeout should be set
	connectTimeout := opts.connectTimeout()
	assert.Equal(pubnub.Config.ConnectTimeout, connectTimeout)
}

// ===========================
// Edge Case Tests
// ===========================

func TestDeleteFileWithLongValues(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
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

func TestDeleteFileWithUnicodeValues(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "测试频道"
	opts.ID = "文件标识符"
	opts.Name = "测试文件.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	assert.Contains(path, "测试频道")
	assert.Contains(path, "文件标识符")
	assert.Contains(path, "测试文件.txt")
}

func TestDeleteFileWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "channel@#$%"
	opts.ID = "id!@#$%^&*()"
	opts.Name = "file@#$%.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	assert.Contains(path, "channel@#$%")
	assert.Contains(path, "id!@#$%^&*()")
	assert.Contains(path, "file@#$%.txt")
}

func TestDeleteFileWithEmptyValues(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.Channel = ""
	opts.ID = ""
	opts.Name = ""

	// Should fail validation
	err := opts.validate()
	assert.NotNil(err)
}

func TestDeleteFileWithSpaces(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "my channel"
	opts.ID = "my file id"
	opts.Name = "my file.txt"

	path, err := opts.buildPath()

	assert.Nil(err)
	assert.Contains(path, "my channel")
	assert.Contains(path, "my file id")
	assert.Contains(path, "my file.txt")
}

func TestDeleteFileWithDots(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
	opts.Channel = "channel.name"
	opts.ID = "file.id"
	opts.Name = "file.name.ext"

	path, err := opts.buildPath()

	assert.Nil(err)
	assert.Contains(path, "channel.name")
	assert.Contains(path, "file.id")
	assert.Contains(path, "file.name.ext")
}

func TestDeleteFileWithQueryParamCombinations(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newDeleteFileOpts(pubnub, pubnub.ctx)
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

func TestDeleteFileValidationErrors(t *testing.T) {
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
			opts := newDeleteFileOpts(pubnub, pubnub.ctx)
			opts.Channel = tc.channel
			opts.ID = tc.id
			opts.Name = tc.fileName

			err := opts.validate()

			assert.NotNil(err)
			assert.Contains(err.Error(), tc.expectedError)
		})
	}
}

func TestDeleteFileBuilderDefaults(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDeleteFileBuilder(pubnub)

	assert.Equal("", builder.opts.Channel)
	assert.Equal("", builder.opts.ID)
	assert.Equal("", builder.opts.Name)
	assert.Nil(builder.opts.QueryParam)
	assert.Nil(builder.opts.Transport)
}

func TestDeleteFileBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newDeleteFileBuilder(pubnub)

	// Check that all builder methods exist and return the builder
	assert.NotNil(builder.Channel("test"))
	assert.NotNil(builder.ID("test"))
	assert.NotNil(builder.Name("test"))
	assert.NotNil(builder.QueryParam(nil))
	assert.NotNil(builder.Transport(nil))
}

func TestDeleteFileNewBuilderWithContext(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())

	builder := newDeleteFileBuilderWithContext(pubnub, pubnub.ctx)

	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pubnub, builder.opts.pubnub)
	assert.Equal(pubnub.ctx, builder.opts.ctx)
}

func TestDeleteFileResponseParsingEdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		jsonBytes   []byte
		expectError bool
	}{
		{
			name:        "Valid response",
			jsonBytes:   []byte(`{"status":200}`),
			expectError: false,
		},
		{
			name:        "Empty response",
			jsonBytes:   []byte(`{}`),
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			jsonBytes:   []byte(`{invalid}`),
			expectError: true,
		},
		{
			name:        "Empty bytes",
			jsonBytes:   []byte(``),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pubnub := NewPubNub(NewDemoConfig())
			opts := newDeleteFileOpts(pubnub, pubnub.ctx)

			resp, _, err := newPNDeleteFileResponse(tc.jsonBytes, opts, StatusResponse{})

			if tc.expectError {
				assert.NotNil(err)
				assert.Equal(emptyDeleteFileResponse, resp)
			} else {
				assert.Nil(err)
				assert.NotNil(resp)
			}
		})
	}
}
