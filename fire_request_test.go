package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessFireGet(t *testing.T, expectedString string, message interface{}) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())

	o := newFireBuilder(pn)
	o.Channel("ch")
	o.Message(message)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/publish/demo/demo/0/ch/0/%s", expectedString),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Empty(body)
}

func AssertSuccessFireGetAllParameters(t *testing.T, expectedString string, message interface{}) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())

	o := newFireBuilder(pn)
	o.Channel("ch")
	o.Message(message)
	o.Serialize(false)
	o.TTL(20)
	o.Meta("a")

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/publish/demo/demo/0/ch/0/%s", expectedString),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Empty(body)
	assert.Equal(o.opts.Meta, "a")
	assert.Equal(o.opts.TTL, 20)
	assert.Equal(o.opts.Serialize, false)
}

func AssertSuccessFirePost(t *testing.T, expectedBody string, message interface{}) {
	assert := assert.New(t)

	opts := &fireOpts{
		Channel:   "ch",
		Message:   message,
		pubnub:    pubnub,
		UsePost:   true,
		Serialize: true,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		"/publish/pub_key/sub_key/0/ch/0",
		u.EscapedPath(), []int{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal(expectedBody, string(body))
}

func AssertSuccessFireQuery(t *testing.T, expectedString string, message interface{}) {
	opts := &fireOpts{
		Channel: "ch",
		Message: "hey",
		pubnub:  pubnub,
	}

	query, _ := opts.buildQuery()

	expected := &url.Values{}
	expected.Set("store", "0")
	expected.Set("norep", "true")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"seqn", "pnsdk", "uuid", "store", "norep"}, []string{})

}

func TestFirePath(t *testing.T) {
	message := "test"
	AssertSuccessFireGet(t, "%22test%22", message)
}

func TestFireQuery(t *testing.T) {
	message := "test"
	AssertSuccessFireQuery(t, "%22test%22?store=0&norep=true&", message)
}

func TestFirePathPost(t *testing.T) {
	AssertSuccessFirePost(t, "[1,2,3]", []int{1, 2, 3})
}
