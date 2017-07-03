package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	utils "github.com/pubnub/go/utils"
	"github.com/stretchr/testify/assert"
)

var pnconfig *Config
var pubnub *PubNub

func init() {
	pnconfig = NewConfig()

	pnconfig.PublishKey = "pub_key"
	pnconfig.SubscribeKey = "sub_key"
	pnconfig.SecretKey = "secret_key"

	pubnub = NewPubNub(pnconfig)
}

func AssertSuccessPublishGet(t *testing.T, expectedString string, message interface{}) {
	assert := assert.New(t)

	opts := &PublishOpts{
		Channel: "ch",
		Message: message,
		pubnub:  pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/publish/pub_key/sub_key/0/ch/0/%s", expectedString),
		u.EscapedPath(), []int{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}

func AssertSuccessPublishPost(t *testing.T, expectedBody string, message interface{}) {
	assert := assert.New(t)

	opts := &PublishOpts{
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

	opts := &PublishOpts{
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

	opts := &PublishOpts{
		Channel: "ch",
		Message: msgMap,
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
	assert.Contains(err.Error(), "Type error, only string is expected")
	assert.Empty(body)
}

func TestPublishMeta(t *testing.T) {
	assert := assert.New(t)

	meta := make(map[string]string)

	meta["one"] = "hey1"
	meta["two"] = "hey2"
	meta["three"] = "hey3"

	opts := &PublishOpts{
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

	opts := &PublishOpts{
		Channel:    "ch",
		Message:    "hey",
		pubnub:     pubnub,
		DoNotStore: true,
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

func TestPublishEncryptDecryptValidation(t *testing.T) {
	assert := assert.New(t)

	encrypted := utils.EncryptString("enigma", "hey")

	msg, err := utils.DecryptString("enigma",
		encrypted)

	assert.Nil(err)
	assert.Equal("hey", msg)
}

func TestPublishSequenceCounter(t *testing.T) {
	assert := assert.New(t)

	meta := make(map[string]string)

	meta["one"] = "hey1"
	meta["two"] = "hey2"
	meta["three"] = "hey3"

	opts := &PublishOpts{
		Channel: "ch",
		Message: "hey",
		pubnub:  pubnub,
		Meta:    meta,
	}
	for i := 1; i <= MaxSequence; i++ {
		counter := <-opts.pubnub.publishSequence
		if counter+1 == MaxSequence {
			assert.Equal(<-opts.pubnub.publishSequence, 1)
			break
		}
	}
}
