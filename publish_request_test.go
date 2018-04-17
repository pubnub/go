package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessPublishGet(t *testing.T, expectedString string, message interface{}) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())

	o := newPublishBuilder(pn)
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

func AssertSuccessPublishPost(t *testing.T, expectedBody string, message interface{}) {
	assert := assert.New(t)

	opts := &publishOpts{
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

func TestPublishMixedGet(t *testing.T) {
	type msg struct {
		One   string `json:"one"`
		Two   string `json:"two"`
		Three string `json:"three"`
	}
	msgStruct := msg{One: "hey1", Two: "hey2", Three: "hey3"}
	msgMap := make(map[string]string)

	msgMap["one"] = "hey1"
	msgMap["two"] = "hey2"
	msgMap["three"] = "hey3"

	AssertSuccessPublishGet(t, "12", 12)
	AssertSuccessPublishGet(t, "%22hey%22", "hey")
	AssertSuccessPublishGet(t, "true", true)
	AssertSuccessPublishGet(t, "%5B%22hey1%22,%22hey2%22,%22hey3%22%5D",
		[]string{"hey1", "hey2", "hey3"})
	AssertSuccessPublishGet(t, "%5B1,2,3%5D", []int{1, 2, 3})
	AssertSuccessPublishGet(t,
		"%7B%22one%22:%22hey1%22,%22two%22:%22hey2%22,%22three%22:%22hey3%22%7D",
		msgStruct)
	AssertSuccessPublishGet(t,
		"%7B%22one%22:%22hey1%22,%22three%22:%22hey3%22,%22two%22:%22hey2%22%7D",
		msgMap)
}

func TestPublishMixedPost(t *testing.T) {
	type msg struct {
		One   string `json:"one"`
		Two   string `json:"two"`
		Three string `json:"three"`
	}
	msgStruct := msg{One: "hey1", Two: "hey2", Three: "hey3"}
	msgMap := make(map[string]string)

	msgMap["one"] = "hey1"
	msgMap["two"] = "hey2"
	msgMap["three"] = "hey3"

	AssertSuccessPublishPost(t, "12", 12)
	AssertSuccessPublishPost(t, "\"hey\"", "hey")
	AssertSuccessPublishPost(t, "true", true)
	AssertSuccessPublishPost(t, "[\"hey1\",\"hey2\",\"hey3\"]",
		[]string{"hey1", "hey2", "hey3"})
	AssertSuccessPublishPost(t, "[1,2,3]", []int{1, 2, 3})
	AssertSuccessPublishPost(t,
		"{\"one\":\"hey1\",\"two\":\"hey2\",\"three\":\"hey3\"}",
		msgStruct)
	AssertSuccessPublishPost(t,
		"{\"one\":\"hey1\",\"three\":\"hey3\",\"two\":\"hey2\"}",
		msgMap)
}

func TestPublishDoNotSerializePost(t *testing.T) {
	assert := assert.New(t)

	message := "{\"one\":\"hey\"}"

	opts := &publishOpts{
		Channel: "ch",
		Message: message,
		pubnub:  pubnub,
		UsePost: true,
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
	assert.NotEmpty(body)
}

func TestPublishDoNotSerializeInvalidPost(t *testing.T) {
	assert := assert.New(t)

	msgMap := make(map[string]string)

	msgMap["one"] = "hey1"
	msgMap["two"] = "hey2"
	msgMap["three"] = "hey3"

	opts := &publishOpts{
		Channel:   "ch",
		Message:   msgMap,
		pubnub:    pubnub,
		UsePost:   true,
		Serialize: false,
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
	assert.Contains(err.Error(), "Message is not JSON serialized.")
	assert.Empty(body)
}

func TestPublishMeta(t *testing.T) {
	assert := assert.New(t)

	meta := make(map[string]string)

	meta["one"] = "hey1"
	meta["two"] = "hey2"
	meta["three"] = "hey3"

	opts := &publishOpts{
		Channel: "ch",
		Message: "hey",
		pubnub:  pubnub,
		Meta:    meta,
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("meta",
		"{\"one\":\"hey1\",\"three\":\"hey3\",\"two\":\"hey2\"}")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"seqn", "pnsdk", "uuid", "store"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}

func TestPublishStore(t *testing.T) {
	assert := assert.New(t)

	opts := &publishOpts{
		Channel:        "ch",
		Message:        "hey",
		pubnub:         pubnub,
		ShouldStore:    true,
		SetShouldStore: true,
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("store", "1")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"seqn", "pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}

func TestPublishEncrypt(t *testing.T) {
	assert := assert.New(t)

	pnconfig.CipherKey = "testCipher"

	opts := &publishOpts{
		Channel: "ch",
		Message: "hey",
		pubnub:  pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)

	assert.Equal(
		"/publish/pub_key/sub_key/0/ch/0/%22%2Bc52pEK3TCTpuEjEFzukRw%3D%3D%22", path)

	//"MnwzPGdVgz2osQCIQJviGg=="
	pnconfig.CipherKey = ""
}

func TestPublishEncryptPNOther(t *testing.T) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())

	pn.Config.CipherKey = "enigma"
	s := map[string]interface{}{
		"not_other": "1234",
		"pn_other":  "\"yay!\"",
	}

	opts := &publishOpts{
		Channel: "ch",
		Message: s,
		pubnub:  pn,
	}

	path, err := opts.buildPath()
	assert.Nil(err)

	assert.Equal(
		"/publish/demo/demo/0/ch/0/%7B%22not_other%22%3A%221234%22%2C%22pn_other%22%3A%22Wi24KS4pcTzvyuGOHubiXg%3D%3D%22%7D", path)

	pn.Config.CipherKey = ""
}

func TestPublishEncryptPNOtherDisable(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	pn.Config.CipherKey = "enigma"
	pn.Config.DisablePNOtherProcessing = true

	s := map[string]interface{}{
		"not_other": "1234",
		"pn_other":  "yay!",
	}

	opts := &publishOpts{
		Channel:   "ch",
		Message:   s,
		pubnub:    pn,
		Serialize: true,
	}

	path, err := opts.buildPath()
	assert.Nil(err)

	assert.Equal(
		"/publish/demo/demo/0/ch/0/%22nG41KRyzJtBR9WWShepyXb3hNh7JJnOoTrQ0SNRcAwRyBYDG2dDhL99svymHR89n%22", path)
	pn.Config.CipherKey = ""
}

func TestPublishSequenceCounter(t *testing.T) {
	assert := assert.New(t)

	meta := make(map[string]string)

	meta["one"] = "hey1"
	meta["two"] = "hey2"
	meta["three"] = "hey3"

	opts := &publishOpts{
		Channel: "ch",
		Message: "hey",
		pubnub:  pubnub,
		Meta:    meta,
	}
	for i := 1; i <= MaxSequence; i++ {
		counter := opts.pubnub.getPublishSequence()
		if counter == MaxSequence {
			assert.Equal(1, opts.pubnub.getPublishSequence())
			break
		}
	}
}
