package pubnub

import (
	"fmt"
	"net/http"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessPublishFileMessageGetAllParameters(t *testing.T, expectedString, messageText, fileID, fileName string, message interface{}, cipher string, genFromIDAndName bool) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = cipher
	pn.Config.UseRandomInitializationVector = false

	o := newPublishFileMessageBuilder(pn)
	m1 := PNPublishFileMessage{}
	if genFromIDAndName {
		if message == nil {
			m := &PNPublishMessage{
				Text: messageText,
			}

			file := &PNFileInfoForPublish{
				ID:   fileID,
				Name: fileName,
			}

			m1 = PNPublishFileMessage{
				PNFile:    file,
				PNMessage: m,
			}
		} else {
			m1 = message.(PNPublishFileMessage)
		}
		o.Message(m1)
	} else {
		o.MessageText(messageText)
		o.FileID(fileID)
		o.FileName(fileName)
	}

	channel := "ch"
	o.Channel(channel)

	o.opts.setTTL = true
	o.TTL(20)
	o.Meta("a")

	path, err := o.opts.buildPath()
	assert.Nil(err)

	query, _ := o.opts.buildQuery()
	for k, v := range *query {
		if k == "pnsdk" || k == "uuid" || k == "seqn" {
			continue
		}
		switch k {
		case "meta":
			assert.Equal("\"a\"", v[0])
		case "store":
			assert.Equal("0", v[0])
		case "norep":
			assert.Equal("true", v[0])
		}
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf(publishFileMessageGetPath, pn.Config.PublishKey, pn.Config.SubscribeKey, channel, "0", expectedString),
		fmt.Sprintf("%s", path),
		[]int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	c := o.opts.config()

	assert.Empty(body)
	assert.Equal(o.opts.Meta, "a")
	assert.Equal(o.opts.TTL, 20)
	assert.Equal(o.opts.UsePost, false)
	assert.Equal(c.UUID, pn.Config.UUID)
	assert.Equal(o.opts.httpMethod(), "GET")
}

func TestPublishFileMessageValidatePublishKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.PublishKey = ""
	opts := newPublishFileMessageOpts(pn, pn.ctx)
	assert.Equal("pubnub/validation: pubnub: Publish File: Missing Publish Key", opts.validate().Error())
}

func TestPublishFileMessageValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newPublishFileMessageOpts(pn, pn.ctx)
	assert.Equal("pubnub/validation: pubnub: Publish File: Missing Subscribe Key", opts.validate().Error())
}

func TestPublishFileMessageValidateFileID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newPublishFileMessageOpts(pn, pn.ctx)
	assert.Equal("pubnub/validation: pubnub: Publish File: Missing File ID", opts.validate().Error())
}

func TestPublishFileMessageValidateFileName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newPublishFileMessageOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.FileID = "sdd"
	assert.Equal("pubnub/validation: pubnub: Publish File: Missing File Name", opts.validate().Error())
}

func TestPublishFileMessageValidateFileMessageNilFileID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	m1 := PNPublishFileMessage{
		PNFile:    nil,
		PNMessage: nil,
	}
	opts := newPublishFileMessageOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.Message = m1
	assert.Equal("pubnub/validation: pubnub: Publish File: Missing File ID", opts.validate().Error())
}

func TestPublishFileMessageValidateFileMessageNilFileName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	file := &PNFileInfoForPublish{
		ID:   "a",
		Name: "",
	}
	m1 := PNPublishFileMessage{
		PNFile:    file,
		PNMessage: nil,
	}
	opts := newPublishFileMessageOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.Message = m1
	assert.Equal("pubnub/validation: pubnub: Publish File: Missing File Name", opts.validate().Error())
}

func TestPublishFileMessageValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newPublishFileMessageOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.Message = "a"
	assert.Equal("pubnub/validation: pubnub: Publish File: Missing Message", opts.validate().Error())
}

func TestPublishFileMessageGetAllParametersFromInterface(t *testing.T) {
	AssertSuccessPublishFileMessageGetAllParameters(t, "%7B%22message%22%3A%7B%22text%22%3A%22test%20message%22%7D%2C%22file%22%3A%7B%22name%22%3A%22test%20file.txt%22%2C%22id%22%3A%22asds%22%7D%7D", "test message", "asds", "test file.txt", nil, "", true)
}

func TestPublishFileMessageGetAllParameters(t *testing.T) {
	AssertSuccessPublishFileMessageGetAllParameters(t, "%7B%22message%22%3A%7B%22text%22%3A%22test%20message%22%7D%2C%22file%22%3A%7B%22name%22%3A%22test%20file.txt%22%2C%22id%22%3A%22asds%22%7D%7D", "test message", "asds", "test file.txt", nil, "", false)
}
func TestPublishFileMessageGetAllParametersFromInterfaceCipher(t *testing.T) {
	AssertSuccessPublishFileMessageGetAllParameters(t, "%22g31ercyjak2YG6ZCA4ii587rApOVOoDTCGCB06CudfJoZhrfRXVpWOAD5mbh44P9%2FdBeUCOEcJEjQRdRmsLm633IHTzPNlFD1AfIDut4f5k%3D%22", "test message", "asds", "test file.txt", nil, "enigma", true)
}

func TestPublishFileMessageGetAllParametersCipher(t *testing.T) {
	AssertSuccessPublishFileMessageGetAllParameters(t, "%22g31ercyjak2YG6ZCA4ii587rApOVOoDTCGCB06CudfJoZhrfRXVpWOAD5mbh44P9%2FdBeUCOEcJEjQRdRmsLm633IHTzPNlFD1AfIDut4f5k%3D%22", "test message", "asds", "test file.txt", nil, "enigma", false)
}

func TestPublishFileMessageGetAllParametersFromMessage(t *testing.T) {
	messageText := "asasdasd"
	fileID := "asasdasd"
	fileName := "asasdasd"
	m := &PNPublishMessage{
		Text: messageText,
	}

	file := &PNFileInfoForPublish{
		ID:   fileID,
		Name: fileName,
	}

	m1 := PNPublishFileMessage{
		PNFile:    file,
		PNMessage: m,
	}
	AssertSuccessPublishFileMessageGetAllParameters(t, "%7B%22message%22%3A%7B%22text%22%3A%22asasdasd%22%7D%2C%22file%22%3A%7B%22name%22%3A%22asasdasd%22%2C%22id%22%3A%22asasdasd%22%7D%7D", "test message", "asds", "test file.txt", m1, "", true)
}
func TestPublishFileMessageGetAllParametersFromMessageCipher(t *testing.T) {
	messageText := "asasdasd1"
	fileID := "asasdasd"
	fileName := "asasdasd"
	m := &PNPublishMessage{
		Text: messageText,
	}

	file := &PNFileInfoForPublish{
		ID:   fileID,
		Name: fileName,
	}

	m1 := PNPublishFileMessage{
		PNFile:    file,
		PNMessage: m,
	}
	AssertSuccessPublishFileMessageGetAllParameters(t, "%22g31ercyjak2YG6ZCA4ii59BezrtHgy%2BYy58G0fftdJbiWKqQKUlENxvOR5F5liVOx51PDn0jJ59adQVj9bWdcGI4s2Qb1sFlo4JHzWEX81M%3D%22", "test message", "asds", "test file.txt", m1, "enigma", true)
}

func AssertPublishFileMessage(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newPublishFileMessageBuilder(pn)
	if testContext {
		o = newPublishFileMessageBuilderWithContext(pn, backgroundContext)
	}

	channel := "chan"
	o.Channel(channel)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)
	messageText := "asasdasd"
	fileID := "asasdasd"
	fileName := "asasdasd"
	m := &PNPublishMessage{
		Text: messageText,
	}

	file := &PNFileInfoForPublish{
		ID:   fileID,
		Name: fileName,
	}

	m1 := PNPublishFileMessage{
		PNFile:    file,
		PNMessage: m,
	}
	o.Message(m1)
	h.AssertPathsEqual(t,
		fmt.Sprintf(publishFileMessageGetPath, pn.Config.SubscribeKey, pn.Config.PublishKey, channel,
			"0",
			"%7B%22message%22%3A%7B%22text%22%3A%22%22%7D%2C%22file%22%3A%7B%22name%22%3A%22%22%2C%22id%22%3A%22%22%7D%7D"),
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

func TestPublishFileMessage(t *testing.T) {
	AssertPublishFileMessage(t, true, false)
}

func TestPublishFileMessageContext(t *testing.T) {
	AssertPublishFileMessage(t, true, true)
}

func TestPublishFileMessageResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newPublishFileMessageOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPublishFileMessageResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestPublishFileMessageResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newPublishFileMessageOpts(pn, pn.ctx)
	jsonBytes := []byte(`[1, "Sent", "12142342544254"]`)

	r, _, err := newPublishFileMessageResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(int64(12142342544254), r.Timestamp)

	assert.Nil(err)
}

// Additional validation tests for PublishFileMessage
func TestPublishFileMessageValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test successful validation with Message object
	file := &PNFileInfoForPublish{
		ID:   "test_id",
		Name: "test_file.txt",
	}
	m := &PNPublishMessage{
		Text: "test message",
	}
	m1 := PNPublishFileMessage{
		PNFile:    file,
		PNMessage: m,
	}

	opts := newPublishFileMessageOpts(pn, pn.ctx)
	opts.Channel = "test_channel"
	opts.Message = m1

	assert.Nil(opts.validate())

	// Test successful validation with individual fields
	opts2 := newPublishFileMessageOpts(pn, pn.ctx)
	opts2.Channel = "test_channel"
	opts2.FileID = "test_id"
	opts2.FileName = "test_file.txt"
	opts2.MessageText = "test message"

	assert.Nil(opts2.validate())
}

func TestPublishFileMessageValidateNoChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// PublishFileMessage doesn't require Channel in validation (unlike other functions)
	opts := newPublishFileMessageOpts(pn, pn.ctx)
	opts.FileID = "test_id"
	opts.FileName = "test_file.txt"

	assert.Nil(opts.validate()) // Should pass even without Channel
}

// Builder pattern tests for PublishFileMessage
func TestPublishFileMessageBuilder(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPublishFileMessageBuilder(pn)
	o.Channel("test_channel")
	o.FileID("test_id")
	o.FileName("test_file.txt")
	o.MessageText("test message")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "test_channel")
	assert.Contains(path, "test_id")
	assert.Contains(path, "test_file.txt")

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}

func TestPublishFileMessageBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newPublishFileMessageBuilder(pn)

	// Test method chaining
	result := o.Channel("test_channel").FileID("test_id").FileName("test_file.txt").MessageText("test message").TTL(24).Meta("meta_data").ShouldStore(true)
	assert.Equal(o, result) // Should return same instance for chaining

	// Verify all values are set correctly
	assert.Equal("test_channel", o.opts.Channel)
	assert.Equal("test_id", o.opts.FileID)
	assert.Equal("test_file.txt", o.opts.FileName)
	assert.Equal("test message", o.opts.MessageText)
	assert.Equal(24, o.opts.TTL)
	assert.Equal("meta_data", o.opts.Meta)
	assert.Equal(true, o.opts.ShouldStore)
	assert.Equal(true, o.opts.setTTL)
	assert.Equal(true, o.opts.setShouldStore)
}

func TestPublishFileMessageBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParams := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	o := newPublishFileMessageBuilder(pn)
	result := o.QueryParam(queryParams)
	assert.Equal(o, result) // Should return same instance for chaining
	assert.Equal(queryParams, o.opts.QueryParam)
}

func TestPublishFileMessageBuilderMessageObject(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	file := &PNFileInfoForPublish{
		ID:   "test_id",
		Name: "test_file.txt",
	}
	m := &PNPublishMessage{
		Text: "test message",
	}
	m1 := PNPublishFileMessage{
		PNFile:    file,
		PNMessage: m,
	}

	o := newPublishFileMessageBuilder(pn)
	result := o.Message(m1)
	assert.Equal(o, result) // Should return same instance for chaining
	assert.Equal(m1, o.opts.Message)
}

// Parameter-specific tests for PublishFileMessage
func TestPublishFileMessageWithTTL(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test normal TTL
	o1 := newPublishFileMessageBuilder(pn)
	o1.TTL(24)
	assert.Equal(24, o1.opts.TTL)
	assert.Equal(true, o1.opts.setTTL)

	// Test zero TTL
	o2 := newPublishFileMessageBuilder(pn)
	o2.TTL(0)
	assert.Equal(0, o2.opts.TTL)
	assert.Equal(true, o2.opts.setTTL)

	// Test negative TTL
	o3 := newPublishFileMessageBuilder(pn)
	o3.TTL(-1)
	assert.Equal(-1, o3.opts.TTL)
	assert.Equal(true, o3.opts.setTTL)

	// Test very large TTL
	o4 := newPublishFileMessageBuilder(pn)
	o4.TTL(999999)
	assert.Equal(999999, o4.opts.TTL)
	assert.Equal(true, o4.opts.setTTL)
}

func TestPublishFileMessageWithMeta(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test string meta
	o1 := newPublishFileMessageBuilder(pn)
	o1.Meta("string_meta")
	assert.Equal("string_meta", o1.opts.Meta)

	// Test map meta
	metaMap := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}
	o2 := newPublishFileMessageBuilder(pn)
	o2.Meta(metaMap)
	assert.Equal(metaMap, o2.opts.Meta)

	// Test nil meta
	o3 := newPublishFileMessageBuilder(pn)
	o3.Meta(nil)
	assert.Nil(o3.opts.Meta)

	// Test numeric meta
	o4 := newPublishFileMessageBuilder(pn)
	o4.Meta(12345)
	assert.Equal(12345, o4.opts.Meta)
}

func TestPublishFileMessageWithShouldStore(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test ShouldStore true
	o1 := newPublishFileMessageBuilder(pn)
	o1.ShouldStore(true)
	assert.Equal(true, o1.opts.ShouldStore)
	assert.Equal(true, o1.opts.setShouldStore)

	// Test ShouldStore false
	o2 := newPublishFileMessageBuilder(pn)
	o2.ShouldStore(false)
	assert.Equal(false, o2.opts.ShouldStore)
	assert.Equal(true, o2.opts.setShouldStore)
}

func TestPublishFileMessageWithUsePost(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test UsePost true (though not implemented)
	o1 := newPublishFileMessageBuilder(pn)
	o1.usePost(true)
	assert.Equal(true, o1.opts.UsePost)

	// Test UsePost false
	o2 := newPublishFileMessageBuilder(pn)
	o2.usePost(false)
	assert.Equal(false, o2.opts.UsePost)
}

// Edge case tests for PublishFileMessage
func TestPublishFileMessageWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with special characters in various fields
	channel := "test-channel_with@special#chars"
	fileID := "file-id_with@special#chars"
	fileName := "file-name_with@special#chars.txt"
	messageText := "message with @special #characters & symbols!"

	o := newPublishFileMessageBuilder(pn)
	o.Channel(channel)
	o.FileID(fileID)
	o.FileName(fileName)
	o.MessageText(messageText)

	assert.Nil(o.opts.validate())

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.NotEmpty(path)
	// URL encoding should handle special characters
	assert.Contains(path, "test-channel_with")
}

func TestPublishFileMessageWithUnicodeCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with Unicode characters
	channel := "测试频道"
	fileID := "файл-идентификатор"
	fileName := "ファイル名.txt"
	messageText := "Message with 中文, русский, and 日本語 text"

	o := newPublishFileMessageBuilder(pn)
	o.Channel(channel)
	o.FileID(fileID)
	o.FileName(fileName)
	o.MessageText(messageText)

	assert.Nil(o.opts.validate())

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.NotEmpty(path)
}

func TestPublishFileMessageWithLongValues(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with very long values
	longChannel := "very_long_channel_name_" + string(make([]byte, 200))
	longFileID := "very_long_file_id_" + string(make([]byte, 200))
	longFileName := "very_long_file_name_" + string(make([]byte, 200)) + ".txt"
	longMessageText := "very_long_message_text_" + string(make([]byte, 500))

	// Fill with valid characters
	for i := 23; i < len(longChannel); i++ {
		longChannel = longChannel[:i] + "a" + longChannel[i+1:]
	}
	for i := 17; i < len(longFileID); i++ {
		longFileID = longFileID[:i] + "b" + longFileID[i+1:]
	}
	for i := 18; i < len(longFileName)-4; i++ {
		longFileName = longFileName[:i] + "c" + longFileName[i+1:]
	}
	for i := 23; i < len(longMessageText); i++ {
		longMessageText = longMessageText[:i] + "d" + longMessageText[i+1:]
	}

	o := newPublishFileMessageBuilder(pn)
	o.Channel(longChannel)
	o.FileID(longFileID)
	o.FileName(longFileName)
	o.MessageText(longMessageText)

	assert.Nil(o.opts.validate())

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.NotEmpty(path)
}

func TestPublishFileMessageWithEmptyValues(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with empty channel (should still validate since Channel not required in validation)
	o1 := newPublishFileMessageBuilder(pn)
	o1.Channel("")
	o1.FileID("test_id")
	o1.FileName("test_file.txt")
	o1.MessageText("test message")

	assert.Nil(o1.opts.validate())

	// Test with empty FileID (should fail validation)
	o2 := newPublishFileMessageBuilder(pn)
	o2.Channel("test_channel")
	o2.FileID("")
	o2.FileName("test_file.txt")
	o2.MessageText("test message")

	assert.NotNil(o2.opts.validate())

	// Test with empty FileName (should fail validation)
	o3 := newPublishFileMessageBuilder(pn)
	o3.Channel("test_channel")
	o3.FileID("test_id")
	o3.FileName("")
	o3.MessageText("test message")

	assert.NotNil(o3.opts.validate())

	// Test with empty MessageText (should still validate if FileID and FileName are present)
	o4 := newPublishFileMessageBuilder(pn)
	o4.Channel("test_channel")
	o4.FileID("test_id")
	o4.FileName("test_file.txt")
	o4.MessageText("")

	assert.Nil(o4.opts.validate())
}

func TestPublishFileMessageWithComplexMeta(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with complex nested meta structure
	complexMeta := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   123,
			"name": "Test User",
			"settings": map[string]interface{}{
				"theme":         "dark",
				"notifications": true,
			},
		},
		"file_info": map[string]interface{}{
			"uploaded_at": "2023-01-01T00:00:00Z",
			"size":        1024,
			"tags":        []string{"important", "document"},
		},
		"unicode_text": "测试 русский ファイル",
	}

	o := newPublishFileMessageBuilder(pn)
	o.Channel("test_channel")
	o.FileID("test_id")
	o.FileName("test_file.txt")
	o.MessageText("test message")
	o.Meta(complexMeta)

	assert.Nil(o.opts.validate())
	assert.Equal(complexMeta, o.opts.Meta)

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.NotEmpty(path)
}

// ===========================
// UseRawMessage Tests
// ===========================

func TestPublishFileMessageBuilderUseRawMessage(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newPublishFileMessageBuilder(pubnub)

	// Test UseRawMessage true
	result := builder.UseRawMessage(true)
	assert.True(builder.opts.UseRawMessage)
	assert.Equal(builder, result) // Fluent interface

	// Test UseRawMessage false
	result = builder.UseRawMessage(false)
	assert.False(builder.opts.UseRawMessage)
	assert.Equal(builder, result) // Fluent interface
}

func TestPublishFileMessageUseRawMessageDefaults(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newPublishFileMessageOpts(pubnub, pubnub.ctx)

	// Test default value
	assert.False(opts.UseRawMessage)
}

func TestPublishFileMessageUseRawMessageMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newPublishFileMessageBuilder(pubnub)

	// Test method chaining with UseRawMessage
	result := builder.
		Channel("test-channel").
		FileID("test-id").
		FileName("test.txt").
		MessageText("test message").
		UseRawMessage(true).
		ShouldStore(true).
		TTL(24)

	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("test-id", builder.opts.FileID)
	assert.Equal("test.txt", builder.opts.FileName)
	assert.Equal("test message", builder.opts.MessageText)
	assert.True(builder.opts.UseRawMessage)
	assert.True(builder.opts.ShouldStore)
	assert.Equal(24, builder.opts.TTL)
	assert.Equal(builder, result) // Fluent interface
}

func TestPublishFileMessageUseRawMessageWithAllParameters(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newPublishFileMessageBuilder(pubnub)

	// Test UseRawMessage with all other parameters
	queryParam := map[string]string{"param1": "value1"}
	meta := map[string]interface{}{"key": "value"}

	result := builder.
		Channel("test-channel").
		FileID("test-id").
		FileName("test.txt").
		MessageText("test message").
		TTL(24).
		Meta(meta).
		ShouldStore(true).
		QueryParam(queryParam).
		UseRawMessage(true).
		Transport(&http.Transport{})

	// Verify all parameters are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("test-id", builder.opts.FileID)
	assert.Equal("test.txt", builder.opts.FileName)
	assert.Equal("test message", builder.opts.MessageText)
	assert.Equal(24, builder.opts.TTL)
	assert.Equal(meta, builder.opts.Meta)
	assert.True(builder.opts.ShouldStore)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.True(builder.opts.UseRawMessage)
	assert.NotNil(builder.opts.Transport)
	assert.Equal(builder, result) // Fluent interface
}

func TestPublishFileMessageBuildRawMessage(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newPublishFileMessageOpts(pubnub, pubnub.ctx)

	// Test with PNPublishFileMessage
	file := &PNFileInfoForPublish{
		ID:   "test-id",
		Name: "test.txt",
	}
	message := &PNPublishMessage{
		Text: "Hello World",
	}
	opts.Message = PNPublishFileMessage{
		PNFile:    file,
		PNMessage: message,
	}

	result := opts.buildRawMessage()
	rawMessage, ok := result.(map[string]interface{})
	assert.True(ok)
	assert.Equal("Hello World", rawMessage["message"])

	fileInfo, ok := rawMessage["file"].(map[string]interface{})
	assert.True(ok)
	assert.Equal("test-id", fileInfo["id"])
	assert.Equal("test.txt", fileInfo["name"])

	// Test with PNPublishFileMessageRaw
	fileRaw := &PNFileInfoForPublish{
		ID:   "raw-id",
		Name: "raw.txt",
	}
	messageRaw := &PNPublishMessageRaw{
		Text: "Raw Message",
	}
	opts.Message = PNPublishFileMessageRaw{
		PNFile:    fileRaw,
		PNMessage: messageRaw,
	}

	result = opts.buildRawMessage()
	rawMessage, ok = result.(map[string]interface{})
	assert.True(ok)
	assert.Equal("Raw Message", rawMessage["message"])

	fileInfo, ok = rawMessage["file"].(map[string]interface{})
	assert.True(ok)
	assert.Equal("raw-id", fileInfo["id"])
	assert.Equal("raw.txt", fileInfo["name"])

	// Test with MessageText fallback
	opts.Message = nil
	opts.MessageText = "Fallback message"
	opts.FileID = "fallback-id"
	opts.FileName = "fallback.txt"
	result = opts.buildRawMessage()
	rawMessage, ok = result.(map[string]interface{})
	assert.True(ok, "Fallback should return map[string]interface{}")
	assert.Equal("Fallback message", rawMessage["message"])

	fileInfo, ok = rawMessage["file"].(map[string]interface{})
	assert.True(ok, "Fallback should have file info")
	assert.Equal("fallback-id", fileInfo["id"])
	assert.Equal("fallback.txt", fileInfo["name"])
}

// TestPublishFileMessageBuildRawMessageWithIndividualFields tests the complete message structure
// when using individual field setters (MessageText, FileID, FileName) with UseRawMessage
func TestPublishFileMessageBuildRawMessageWithIndividualFields(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newPublishFileMessageBuilder(pubnub)

	// Test using individual field setters with UseRawMessage(true)
	builder.
		Channel("test-channel").
		MessageText("Test message").
		FileID("test-id").
		FileName("test.txt").
		UseRawMessage(true)

	// Build the path which internally calls buildRawMessage
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.NotEmpty(path)

	// Verify the message structure
	result := builder.opts.buildRawMessage()
	rawMessage, ok := result.(map[string]interface{})
	assert.True(ok, "Should return map[string]interface{} when using individual fields")
	assert.Equal("Test message", rawMessage["message"], "Message text should be present")

	fileInfo, ok := rawMessage["file"].(map[string]interface{})
	assert.True(ok, "Should have file info when using individual fields")
	assert.Equal("test-id", fileInfo["id"], "File ID should be present")
	assert.Equal("test.txt", fileInfo["name"], "File name should be present")
}

// TestPublishFileMessageUseRawMessageWithIndividualFieldsViaPath tests the full path generation
// when using individual field setters with UseRawMessage to ensure the message is properly encoded
func TestPublishFileMessageUseRawMessageWithIndividualFieldsViaPath(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newPublishFileMessageBuilder(pubnub)

	// Test the complete flow with individual field setters
	builder.
		Channel("test-channel").
		MessageText("Hello World").
		FileID("file-123").
		FileName("document.pdf").
		UseRawMessage(true)

	path, err := builder.opts.buildPath()
	assert.Nil(err)
	assert.NotEmpty(path)

	// The path should contain the URL-encoded JSON with raw message format
	// Expected structure: {"message":"Hello World","file":{"id":"file-123","name":"document.pdf"}}
	assert.Contains(path, "test-channel", "Path should contain channel name")
	// The message should be in the path as URL-encoded JSON
	assert.Contains(path, "%22message%22", "Path should contain encoded message field")
	assert.Contains(path, "%22file%22", "Path should contain encoded file field")
}

// ===========================
// CustomMessageType Tests
// ===========================

func TestPublishFileMessageBuilderCustomMessageType(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newPublishFileMessageBuilder(pubnub)

	// Test CustomMessageType
	result := builder.CustomMessageType("file-message")
	assert.Equal("file-message", builder.opts.CustomMessageType)
	assert.Equal(builder, result) // Fluent interface
}

func TestPublishFileMessageCustomMessageTypeDefaults(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	opts := newPublishFileMessageOpts(pubnub, pubnub.ctx)

	// Test default value
	assert.Empty(opts.CustomMessageType)
}

func TestPublishFileMessageCustomMessageTypeMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newPublishFileMessageBuilder(pubnub)

	// Test method chaining with CustomMessageType
	result := builder.
		Channel("test-channel").
		FileID("test-id").
		FileName("test.txt").
		MessageText("test message").
		CustomMessageType("file-message").
		ShouldStore(true).
		TTL(24)

	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("test-id", builder.opts.FileID)
	assert.Equal("test.txt", builder.opts.FileName)
	assert.Equal("test message", builder.opts.MessageText)
	assert.Equal("file-message", builder.opts.CustomMessageType)
	assert.True(builder.opts.ShouldStore)
	assert.Equal(24, builder.opts.TTL)
	assert.Equal(builder, result) // Fluent interface
}

func TestPublishFileMessageCustomMessageTypeValidation(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())

	// Test valid CustomMessageType
	file := &PNFileInfoForPublish{
		ID:   "test-id",
		Name: "test.txt",
	}
	m := &PNPublishMessage{
		Text: "test message",
	}
	m1 := PNPublishFileMessage{
		PNFile:    file,
		PNMessage: m,
	}

	opts := newPublishFileMessageOpts(pubnub, pubnub.ctx)
	opts.Channel = "test-channel"
	opts.Message = m1
	opts.CustomMessageType = "valid-type"
	assert.Nil(opts.validate())

	// Test invalid CustomMessageType - too short
	opts2 := newPublishFileMessageOpts(pubnub, pubnub.ctx)
	opts2.Channel = "test-channel"
	opts2.Message = m1
	opts2.CustomMessageType = "ab"
	assert.NotNil(opts2.validate())
	assert.Contains(opts2.validate().Error(), "Invalid CustomMessageType")

	// Test invalid CustomMessageType - too long
	opts3 := newPublishFileMessageOpts(pubnub, pubnub.ctx)
	opts3.Channel = "test-channel"
	opts3.Message = m1
	opts3.CustomMessageType = "this-is-a-very-long-custom-message-type-that-exceeds-the-limit"
	assert.NotNil(opts3.validate())
	assert.Contains(opts3.validate().Error(), "Invalid CustomMessageType")

	// Test invalid CustomMessageType - invalid characters
	opts4 := newPublishFileMessageOpts(pubnub, pubnub.ctx)
	opts4.Channel = "test-channel"
	opts4.Message = m1
	opts4.CustomMessageType = "invalid@type"
	assert.NotNil(opts4.validate())
	assert.Contains(opts4.validate().Error(), "Invalid CustomMessageType")

	// Test empty CustomMessageType (should be valid)
	opts5 := newPublishFileMessageOpts(pubnub, pubnub.ctx)
	opts5.Channel = "test-channel"
	opts5.Message = m1
	opts5.CustomMessageType = ""
	assert.Nil(opts5.validate())
}

func TestPublishFileMessageCustomMessageTypeInQuery(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())

	file := &PNFileInfoForPublish{
		ID:   "test-id",
		Name: "test.txt",
	}
	m := &PNPublishMessage{
		Text: "test message",
	}
	m1 := PNPublishFileMessage{
		PNFile:    file,
		PNMessage: m,
	}

	builder := newPublishFileMessageBuilder(pubnub)
	builder.Channel("test-channel")
	builder.Message(m1)
	builder.CustomMessageType("file-message")

	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("file-message", query.Get("custom_message_type"))
}

func TestPublishFileMessageCustomMessageTypeWithAllParameters(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())
	builder := newPublishFileMessageBuilder(pubnub)

	// Test CustomMessageType with all other parameters
	queryParam := map[string]string{"param1": "value1"}
	meta := map[string]interface{}{"key": "value"}

	result := builder.
		Channel("test-channel").
		FileID("test-id").
		FileName("test.txt").
		MessageText("test message").
		TTL(24).
		Meta(meta).
		ShouldStore(true).
		QueryParam(queryParam).
		CustomMessageType("file-message").
		UseRawMessage(true).
		Transport(&http.Transport{})

	// Verify all parameters are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal("test-id", builder.opts.FileID)
	assert.Equal("test.txt", builder.opts.FileName)
	assert.Equal("test message", builder.opts.MessageText)
	assert.Equal(24, builder.opts.TTL)
	assert.Equal(meta, builder.opts.Meta)
	assert.True(builder.opts.ShouldStore)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal("file-message", builder.opts.CustomMessageType)
	assert.True(builder.opts.UseRawMessage)
	assert.NotNil(builder.opts.Transport)
	assert.Equal(builder, result) // Fluent interface
}

func TestPublishFileMessageCustomMessageTypeValidCharacters(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())

	file := &PNFileInfoForPublish{
		ID:   "test-id",
		Name: "test.txt",
	}
	m := &PNPublishMessage{
		Text: "test message",
	}
	m1 := PNPublishFileMessage{
		PNFile:    file,
		PNMessage: m,
	}

	// Test valid combinations
	validTypes := []string{
		"abc",
		"ABC",
		"test-type",
		"test_type",
		"Test-Message_Type",
		"file-message",
		"custom_message_type",
	}

	for _, validType := range validTypes {
		opts := newPublishFileMessageOpts(pubnub, pubnub.ctx)
		opts.Channel = "test-channel"
		opts.Message = m1
		opts.CustomMessageType = validType
		assert.Nil(opts.validate(), "Type '%s' should be valid", validType)
	}
}

func TestPublishFileMessageCustomMessageTypeInvalidCharacters(t *testing.T) {
	assert := assert.New(t)
	pubnub := NewPubNub(NewDemoConfig())

	file := &PNFileInfoForPublish{
		ID:   "test-id",
		Name: "test.txt",
	}
	m := &PNPublishMessage{
		Text: "test message",
	}
	m1 := PNPublishFileMessage{
		PNFile:    file,
		PNMessage: m,
	}

	// Test invalid combinations
	invalidTypes := []string{
		"ab",         // too short
		"a@b",        // invalid character @
		"test type",  // space not allowed
		"test!type",  // ! not allowed
		"test#type",  // # not allowed
		"test$type",  // $ not allowed
		"test%type",  // % not allowed
		"test&type",  // & not allowed
		"abc123",     // digits not allowed
		"test123",    // digits not allowed
		"test.type",  // dot not allowed
		"test/type",  // slash not allowed
		"test\\type", // backslash not allowed
		"test:type",  // colon not allowed
	}

	for _, invalidType := range invalidTypes {
		opts := newPublishFileMessageOpts(pubnub, pubnub.ctx)
		opts.Channel = "test-channel"
		opts.Message = m1
		opts.CustomMessageType = invalidType
		assert.NotNil(opts.validate(), "Type '%s' should be invalid", invalidType)
	}
}
