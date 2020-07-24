package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
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
	opts := &listFilesOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNListFilesResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestListFilesResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &listFilesOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":[{"name":"test_file_upload_name_42893.txt","id":"9ef0e123-1e4a-40b9-89d5-f4be0e8b1f2c","size":21904,"created":"2020-07-21T09:10:55Z"}],"next":null,"count":1}`)

	r, _, err := newPNListFilesResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(r.Count, 1)
	assert.Equal(r.Data[0].ID, "9ef0e123-1e4a-40b9-89d5-f4be0e8b1f2c")
	assert.Equal(r.Data[0].Name, "test_file_upload_name_42893.txt")
	assert.Equal(r.Data[0].Size, 21904)
	assert.Equal(r.Data[0].Created, "2020-07-21T09:10:55Z")

	assert.Nil(err)
}
