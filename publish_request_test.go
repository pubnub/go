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
	o.TTL(10)
	o.ShouldStore(true)
	o.DoNotReplicate(true)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/publish/demo/demo/0/ch/0/%s", expectedString),
		path, []int{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Empty(body)
	assert.Equal(10, o.opts.TTL)
	assert.Equal(true, o.opts.ShouldStore)
	assert.Equal(true, o.opts.DoNotReplicate)
}

func AssertSuccessPublishGetContext(t *testing.T, expectedString string, message interface{}) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())

	o := newPublishBuilderWithContext(pn, backgroundContext)
	o.Channel("ch")
	o.Message(message)
	o.TTL(10)
	o.ShouldStore(true)
	o.DoNotReplicate(true)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/publish/demo/demo/0/ch/0/%s", expectedString),
		path, []int{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Empty(body)
	assert.Equal(10, o.opts.TTL)
	assert.Equal(true, o.opts.ShouldStore)
	assert.Equal(true, o.opts.DoNotReplicate)
}

func AssertSuccessPublishGet2(t *testing.T, expectedString string, message interface{}) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	pn.Config.AuthKey = "a"

	o := newPublishBuilder(pn)
	o.Channel("ch")
	o.Message(message)
	o.TTL(10)
	o.ShouldStore(false)
	o.DoNotReplicate(true)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/publish/demo/demo/0/ch/0/%s", expectedString),
		path, []int{})

	query, err := o.opts.buildQuery()
	//log.Println(query)

	assert.Nil(err)
	expected := &url.Values{}
	expected.Set("seqn", "1")
	expected.Set("uuid", pn.Config.UUID)
	expected.Set("ttl", "10")
	expected.Set("pnsdk", Version)
	expected.Set("norep", "true")
	expected.Set("store", "0")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"seqn", "pnsdk", "uuid", "store"}, []string{})

}

func AssertSuccessPublishGet3(t *testing.T, expectedString string, message interface{}) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	pn.Config.AuthKey = "a"

	o := newPublishBuilder(pn)
	o.Channel("ch")
	o.Message(message)
	o.TTL(10)
	o.ShouldStore(false)
	o.DoNotReplicate(true)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	o.QueryParam(queryParam)

	query, err := o.opts.buildQuery()
	//log.Println(query)

	assert.Nil(err)
	expected := &url.Values{}
	expected.Set("seqn", "1")
	expected.Set("uuid", pn.Config.UUID)
	expected.Set("ttl", "10")
	expected.Set("pnsdk", Version)
	expected.Set("norep", "true")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")
	expected.Set("store", "0")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"seqn", "pnsdk", "uuid", "store"}, []string{})

}

func AssertSuccessPublishGetAuth(t *testing.T, expectedString string, message interface{}) {

	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	pn.Config.AuthKey = "PubAuthKey"

	o := newPublishBuilder(pn)
	o.Channel("ch")
	o.Message(message)
	o.TTL(10)
	o.ShouldStore(true)
	o.DoNotReplicate(true)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/publish/demo/demo/0/ch/0/%s", expectedString),
		path, []int{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Empty(body)
	assert.Equal(10, o.opts.TTL)
	assert.Equal(true, o.opts.ShouldStore)
	assert.Equal(true, o.opts.DoNotReplicate)

}

func AssertSuccessPublishGetMeta(t *testing.T, expectedString string, message interface{}) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())

	o := newPublishBuilder(pn)
	o.Meta(nil)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/publish/demo/demo/0/ch/0/%s", expectedString),
		path, []int{})

	_, err1 := o.opts.buildQuery()

	assert.Nil(err1)
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
	AssertSuccessPublishGetAuth(t, "12", 12)
	AssertSuccessPublishGet(t, "%22hey%22", "hey")
	AssertSuccessPublishGet(t, "true", true)
	AssertSuccessPublishGet(t, "%5B%22hey1%22%2C%22hey2%22%2C%22hey3%22%5D",
		[]string{"hey1", "hey2", "hey3"})
	AssertSuccessPublishGet(t, "%5B1%2C2%2C3%5D", []int{1, 2, 3})
	AssertSuccessPublishGet(t,
		"%7B%22one%22%3A%22hey1%22%2C%22two%22%3A%22hey2%22%2C%22three%22%3A%22hey3%22%7D",
		msgStruct)
	AssertSuccessPublishGet(t,
		"%7B%22one%22%3A%22hey1%22%2C%22three%22%3A%22hey3%22%2C%22two%22%3A%22hey2%22%7D",
		msgMap)

	AssertSuccessPublishGetContext(t, "12", 12)
	AssertSuccessPublishGetContext(t, "%22hey%22", "hey")
	AssertSuccessPublishGetContext(t, "true", true)
	AssertSuccessPublishGetContext(t, "%5B%22hey1%22%2C%22hey2%22%2C%22hey3%22%5D",
		[]string{"hey1", "hey2", "hey3"})
	AssertSuccessPublishGetContext(t, "%5B1%2C2%2C3%5D", []int{1, 2, 3})
	AssertSuccessPublishGetContext(t,
		"%7B%22one%22%3A%22hey1%22%2C%22two%22%3A%22hey2%22%2C%22three%22%3A%22hey3%22%7D",
		msgStruct)
	AssertSuccessPublishGetContext(t,
		"%7B%22one%22%3A%22hey1%22%2C%22three%22%3A%22hey3%22%2C%22two%22%3A%22hey2%22%7D",
		msgMap)

	AssertSuccessPublishGet2(t, "12", 12)
	AssertSuccessPublishGet3(t, "12", 12)
	AssertSuccessPublishGet2(t, "%22hey%22", "hey")
	AssertSuccessPublishGet2(t, "true", true)
	AssertSuccessPublishGet2(t, "%5B%22hey1%22%2C%22hey2%22%2C%22hey3%22%5D",
		[]string{"hey1", "hey2", "hey3"})
	AssertSuccessPublishGet2(t, "%5B1%2C2%2C3%5D", []int{1, 2, 3})
	AssertSuccessPublishGet2(t,
		"%7B%22one%22%3A%22hey1%22%2C%22two%22%3A%22hey2%22%2C%22three%22%3A%22hey3%22%7D",
		msgStruct)
	AssertSuccessPublishGet2(t,
		"%7B%22one%22%3A%22hey1%22%2C%22three%22%3A%22hey3%22%2C%22two%22%3A%22hey2%22%7D",
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
		setShouldStore: true,
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

	pnconfig.CipherKey = ""
}

func TestPublishEncryptPNOther(t *testing.T) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())

	pn.Config.CipherKey = "enigma"
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
		"/publish/demo/demo/0/ch/0/%22bCC%2FkQbGdScQ0teYcawUsnASfRpUioutNKQfUAQNc46gWR%2FJnz8Ks5n%2FvfKnDkE6%22", path)
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

func TestNewPublishResponse(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newPublishResponse(jsonBytes, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestNewPublishResponseTimestamp(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`[1, Sent, "a"]`)

	_, _, err := newPublishResponse(jsonBytes, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {[1, Sent, \"a\"]}", err.Error())
}

func TestNewPublishResponseTimestamp2(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`[1, "Sent", "a"]`)

	_, _, err := newPublishResponse(jsonBytes, StatusResponse{})
	assert.Contains(err.Error(), "parsing \"a\": invalid syntax")
}

func TestPublishValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &publishOpts{
		pubnub: pn,
	}

	assert.Equal("pubnub/validation: pubnub: Publish: Missing Subscribe Key", opts.validate().Error())
}
