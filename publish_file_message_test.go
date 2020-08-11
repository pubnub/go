package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessPublishFileMessageGetAllParameters(t *testing.T, expectedString, messageText, fileID, fileName string, message interface{}, cipher string, genFromIDAndName bool) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = cipher
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
	opts := &publishFileMessageOpts{
		pubnub: pn,
	}
	assert.Equal("pubnub/validation: pubnub: Publish File: Missing Publish Key", opts.validate().Error())
}

func TestPublishFileMessageValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &publishFileMessageOpts{
		pubnub: pn,
	}
	assert.Equal("pubnub/validation: pubnub: Publish File: Missing Subscribe Key", opts.validate().Error())
}

func TestPublishFileMessageValidateFileID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &publishFileMessageOpts{
		pubnub: pn,
	}
	assert.Equal("pubnub/validation: pubnub: Publish File: Missing File ID", opts.validate().Error())
}

func TestPublishFileMessageValidateFileName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &publishFileMessageOpts{
		Channel: "ch",
		pubnub:  pn,
		FileID:  "sdd",
	}
	assert.Equal("pubnub/validation: pubnub: Publish File: Missing File Name", opts.validate().Error())
}

func TestPublishFileMessageValidateFileMessageNilFileID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	m1 := PNPublishFileMessage{
		PNFile:    nil,
		PNMessage: nil,
	}
	opts := &publishFileMessageOpts{
		Channel: "ch",
		pubnub:  pn,
		Message: m1,
	}
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
	opts := &publishFileMessageOpts{
		Channel: "ch",
		pubnub:  pn,
		Message: m1,
	}
	assert.Equal("pubnub/validation: pubnub: Publish File: Missing File Name", opts.validate().Error())
}

func TestPublishFileMessageValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &publishFileMessageOpts{
		Channel: "ch",
		Message: "a",
		pubnub:  pn,
	}
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
	opts := &publishFileMessageOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPublishFileMessageResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestPublishFileMessageResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &publishFileMessageOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`[1, "Sent", "12142342544254"]`)

	r, _, err := newPublishFileMessageResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(int64(12142342544254), r.Timestamp)

	assert.Nil(err)
}
