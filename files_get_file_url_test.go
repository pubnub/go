package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertGetFileURL(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newGetFileURLBuilder(pn)
	if testContext {
		o = newGetFileURLBuilderWithContext(pn, backgroundContext)
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
		fmt.Sprintf(getFileURLPath, pn.Config.SubscribeKey, channel, id, name),
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

func TestGetFileURL(t *testing.T) {
	AssertGetFileURL(t, true, false)
}

func TestGetFileURLContext(t *testing.T) {
	AssertGetFileURL(t, true, true)
}

func TestGetFileURLResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	channel := "chan"
	id := "fileid"
	name := "filename"

	r, _, err := pn.GetFileURL().ID(id).Name(name).Channel(channel).Execute()
	assert.Contains(r.URL, fmt.Sprintf("%s/files/%s/%s", channel, id, name))

	assert.Nil(err)
}

// Validation tests for GetFileURL
func TestGetFileURLValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newGetFileURLOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.ID = "id"
	opts.Name = "name"

	assert.Equal("pubnub/validation: pubnub: Get File URL: Missing Subscribe Key", opts.validate().Error())
}

func TestGetFileURLValidateChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetFileURLOpts(pn, pn.ctx)
	opts.ID = "id"
	opts.Name = "name"

	assert.Equal("pubnub/validation: pubnub: Get File URL: Missing Channel", opts.validate().Error())
}

func TestGetFileURLValidateFileID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetFileURLOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.Name = "name"

	assert.Equal("pubnub/validation: pubnub: Get File URL: Missing File ID", opts.validate().Error())
}

func TestGetFileURLValidateFileName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetFileURLOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.ID = "id"

	assert.Equal("pubnub/validation: pubnub: Get File URL: Missing File Name", opts.validate().Error())
}

func TestGetFileURLValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGetFileURLOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.ID = "id"
	opts.Name = "name"

	assert.Nil(opts.validate())
}

// Builder pattern tests for GetFileURL
func TestGetFileURLBuilder(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetFileURLBuilder(pn)
	o.Channel("ch")
	o.ID("id")
	o.Name("name")

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(getFileURLPath, pn.Config.SubscribeKey, "ch", "id", "name"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}

func TestGetFileURLBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGetFileURLBuilder(pn)

	// Test method chaining
	result := o.Channel("test_channel").ID("test_id").Name("test_name.txt")
	assert.Equal(o, result) // Should return same instance for chaining

	// Verify all values are set correctly
	assert.Equal("test_channel", o.opts.Channel)
	assert.Equal("test_id", o.opts.ID)
	assert.Equal("test_name.txt", o.opts.Name)
}

func TestGetFileURLBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParams := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	o := newGetFileURLBuilder(pn)
	result := o.QueryParam(queryParams)
	assert.Equal(o, result) // Should return same instance for chaining
	assert.Equal(queryParams, o.opts.QueryParam)
}

// Edge case tests for GetFileURL
func TestGetFileURLWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with special characters in channel, ID, and name
	channel := "test-channel_with@special#chars"
	id := "file-id_with@special#chars"
	name := "file-name_with@special#chars.txt"

	opts := newGetFileURLOpts(pn, pn.ctx)
	opts.Channel = channel
	opts.ID = id
	opts.Name = name

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, channel)
	assert.Contains(path, id)
	assert.Contains(path, name)
}

func TestGetFileURLWithUnicodeCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with Unicode characters
	channel := "测试频道"
	id := "файл-идентификатор"
	name := "ファイル名.txt"

	opts := newGetFileURLOpts(pn, pn.ctx)
	opts.Channel = channel
	opts.ID = id
	opts.Name = name

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	// Path should be built without errors even with Unicode characters
	assert.NotEmpty(path)
}

func TestGetFileURLWithLongNames(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with very long names
	longChannel := "very_long_channel_name_" + string(make([]byte, 200))
	longID := "very_long_file_id_" + string(make([]byte, 200))
	longName := "very_long_file_name_" + string(make([]byte, 200)) + ".txt"

	// Fill with valid characters
	for i := 20; i < len(longChannel); i++ {
		longChannel = longChannel[:i] + "a" + longChannel[i+1:]
	}
	for i := 17; i < len(longID); i++ {
		longID = longID[:i] + "b" + longID[i+1:]
	}
	for i := 18; i < len(longName)-4; i++ {
		longName = longName[:i] + "c" + longName[i+1:]
	}

	opts := newGetFileURLOpts(pn, pn.ctx)
	opts.Channel = longChannel
	opts.ID = longID
	opts.Name = longName

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.NotEmpty(path)
}
