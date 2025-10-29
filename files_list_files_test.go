package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertListFiles(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newListFilesBuilder(pn)
	if testContext {
		o = newListFilesBuilderWithContext(pn, backgroundContext)
	}

	channel := "chan"
	o.Channel(channel)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(listFilesPath, pn.Config.SubscribeKey, channel),
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

func TestListFiles(t *testing.T) {
	AssertListFiles(t, true, false)
}

func TestListFilesContext(t *testing.T) {
	AssertListFiles(t, true, true)
}

func TestListFilesResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListFilesOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNListFilesResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestListFilesResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListFilesOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200,"data":[{"name":"test_file_upload_name_42893.txt","id":"9ef0e123-1e4a-40b9-89d5-f4be0e8b1f2c","size":21904,"created":"2020-07-21T09:10:55Z"}],"next":null,"count":1}`)

	r, _, err := newPNListFilesResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(r.Count, 1)
	assert.Equal(r.Data[0].ID, "9ef0e123-1e4a-40b9-89d5-f4be0e8b1f2c")
	assert.Equal(r.Data[0].Name, "test_file_upload_name_42893.txt")
	assert.Equal(r.Data[0].Size, 21904)
	assert.Equal(r.Data[0].Created, "2020-07-21T09:10:55Z")

	assert.Nil(err)
}

// Validation tests for ListFiles
func TestListFilesValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newListFilesOpts(pn, pn.ctx)

	assert.Equal("pubnub/validation: pubnub: List Files: Missing Subscribe Key", opts.validate().Error())
}

func TestListFilesValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newListFilesOpts(pn, pn.ctx)

	assert.Nil(opts.validate())
}

// Builder pattern tests for ListFiles
func TestListFilesBuilder(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newListFilesBuilder(pn)
	o.Channel("ch")

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(listFilesPath, pn.Config.SubscribeKey, "ch"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}

func TestListFilesBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newListFilesBuilder(pn)

	// Test method chaining
	result := o.Channel("test_channel").Limit(50).Next("next_token")
	assert.Equal(o, result) // Should return same instance for chaining

	// Verify all values are set correctly
	assert.Equal("test_channel", o.opts.Channel)
	assert.Equal(50, o.opts.Limit)
	assert.Equal("next_token", o.opts.Next)
}

func TestListFilesBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParams := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	o := newListFilesBuilder(pn)
	result := o.QueryParam(queryParams)
	assert.Equal(o, result) // Should return same instance for chaining
	assert.Equal(queryParams, o.opts.QueryParam)
}

// Parameter-specific tests for ListFiles
func TestListFilesWithDefaultLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newListFilesOpts(pn, pn.ctx)
	assert.Equal(listFilesLimit, opts.Limit) // Should be 100 by default

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("limit"))
}

func TestListFilesWithCustomLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newListFilesBuilder(pn)
	o.Limit(25)

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("25", query.Get("limit"))
}

func TestListFilesWithNext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newListFilesBuilder(pn)
	o.Next("some_pagination_token")

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("some_pagination_token", query.Get("next"))
}

func TestListFilesWithoutNext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newListFilesBuilder(pn)
	// Don't set Next

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("", query.Get("next")) // Should be empty when not set
}

func TestListFilesWithoutChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := newListFilesOpts(pn, pn.ctx)
	// Don't set Channel - should still be valid

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	// Path should be built even without channel (empty channel)
	assert.Contains(path, fmt.Sprintf("/v1/files/%s/channels/", pn.Config.SubscribeKey))
}

// Edge case tests for ListFiles
func TestListFilesWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with special characters in channel name
	channel := "test-channel_with@special#chars"

	opts := newListFilesOpts(pn, pn.ctx)
	opts.Channel = channel

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, channel)
}

func TestListFilesWithUnicodeCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with Unicode characters
	channel := "测试频道"

	opts := newListFilesOpts(pn, pn.ctx)
	opts.Channel = channel

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	// Path should be built without errors even with Unicode characters
	assert.NotEmpty(path)
}

func TestListFilesWithLongChannelName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with very long channel name
	longChannel := "very_long_channel_name_" + string(make([]byte, 200))

	// Fill with valid characters
	for i := 23; i < len(longChannel); i++ {
		longChannel = longChannel[:i] + "a" + longChannel[i+1:]
	}

	opts := newListFilesOpts(pn, pn.ctx)
	opts.Channel = longChannel

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.NotEmpty(path)
}

func TestListFilesWithExtremeLimit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with zero limit
	o1 := newListFilesBuilder(pn)
	o1.Limit(0)
	query1, err1 := o1.opts.buildQuery()
	assert.Nil(err1)
	assert.Equal("0", query1.Get("limit"))

	// Test with negative limit
	o2 := newListFilesBuilder(pn)
	o2.Limit(-1)
	query2, err2 := o2.opts.buildQuery()
	assert.Nil(err2)
	assert.Equal("-1", query2.Get("limit"))

	// Test with very large limit
	o3 := newListFilesBuilder(pn)
	o3.Limit(999999)
	query3, err3 := o3.opts.buildQuery()
	assert.Nil(err3)
	assert.Equal("999999", query3.Get("limit"))
}

func TestListFilesWithLongNextToken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with very long next token
	longNext := "very_long_next_token_" + string(make([]byte, 500))

	// Fill with valid characters
	for i := 20; i < len(longNext); i++ {
		longNext = longNext[:i] + "x" + longNext[i+1:]
	}

	o := newListFilesBuilder(pn)
	o.Next(longNext)

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal(longNext, query.Get("next"))
}
