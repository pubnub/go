package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
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
		o = newDownloadFileBuilderWithContext(pn, backgroundContext)
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
