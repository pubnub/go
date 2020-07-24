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

func TestAssertSuccessFireGetContext(t *testing.T) {
	assert := assert.New(t)
	message := "test"
	pn := NewPubNub(NewDemoConfig())

	o := newFireBuilderWithContext(pn, backgroundContext)
	o.Channel("ch")
	o.Message(message)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/publish/demo/demo/0/ch/0/%s", "%22test%22"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Empty(body)
}

func AssertSuccessFirePostAllParameters(t *testing.T, expectedString string, message interface{}, cipher string) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = cipher

	o := newFireBuilder(pn)
	o.Channel("ch")
	o.Message(message)
	o.Serialize(false)
	o.UsePost(true)
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
		"/publish/demo/demo/0/ch/0",
		path,
		[]int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	c := o.opts.config()

	assert.Equal(expectedString, string(body))
	assert.Equal(o.opts.Meta, "a")
	assert.Equal(o.opts.TTL, 20)
	assert.Equal(o.opts.UsePost, true)
	assert.Equal(c.UUID, pn.Config.UUID)
	assert.Equal(o.opts.Serialize, false)
	assert.Equal(o.opts.httpMethod(), "POST")
}

func AssertSuccessFireGetAllParameters(t *testing.T, expectedString string, message interface{}, cipher string) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = cipher

	o := newFireBuilder(pn)
	o.Channel("ch")
	o.Message(message)
	o.Serialize(false)
	o.UsePost(false)
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
		fmt.Sprintf("/publish/demo/demo/0/ch/0/%s", expectedString),
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
	assert.Equal(o.opts.Serialize, false)
	assert.Equal(o.opts.httpMethod(), "GET")
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

	assert.Equal(opts.UsePost, true)
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

func TestFireDoNotSerializePost(t *testing.T) {
	assert := assert.New(t)

	message := "{\"one\":\"hey\"}"

	opts := &fireOpts{
		Channel:   "ch",
		Message:   message,
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
	assert.Nil(err)
	assert.NotEmpty(body)
}

func TestFireDoNotSerializeQueryParam(t *testing.T) {
	assert := assert.New(t)

	message := "{\"one\":\"hey\"}"
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := &fireOpts{
		Channel:    "ch",
		Message:    message,
		pubnub:     pubnub,
		QueryParam: queryParam,
	}

	b, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("v1", b.Get("q1"))
	assert.Equal("v2", b.Get("q2"))

}

func TestValidatePublishKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.PublishKey = ""
	opts := &fireOpts{
		pubnub: pn,
	}
	assert.Equal("pubnub/validation: pubnub: Fire: Missing Publish Key", opts.validate().Error())
}

func TestValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &fireOpts{
		pubnub: pn,
	}
	assert.Equal("pubnub/validation: pubnub: Fire: Missing Subscribe Key", opts.validate().Error())
}

func TestValidateChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &fireOpts{
		pubnub: pn,
	}
	assert.Equal("pubnub/validation: pubnub: Fire: Missing Channel", opts.validate().Error())
}

func TestValidateMessage(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &fireOpts{
		Channel: "ch",
		pubnub:  pn,
	}
	assert.Equal("pubnub/validation: pubnub: Fire: Missing Message", opts.validate().Error())
}

func TestValidate(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &fireOpts{
		Channel: "ch",
		Message: "a",
		pubnub:  pn,
	}
	err := opts.validate()
	assert.Nil(err)
}

func TestFirePath(t *testing.T) {
	message := "test"
	AssertSuccessFireGet(t, "%22test%22", message)
}

func TestFireQuery(t *testing.T) {
	message := "test"
	AssertSuccessFireQuery(t, "%22test%22?store=0&norep=true&", message)
}

func TestFireGetAllParameters(t *testing.T) {
	message := "test"
	AssertSuccessFireGetAllParameters(t, "%22test%22", message, "")
}
func TestFireGetAllParametersCipher(t *testing.T) {
	message := "test"
	AssertSuccessFireGetAllParameters(t, "%22c3dSanMrRnc4ZnNNT1BEaGFnZmd1QT09%22", message, "enigma")
}

func TestFirePostAllParameters(t *testing.T) {
	message := "test"
	AssertSuccessFirePostAllParameters(t, "\"+3AfkVAl8saHsXJdtOhRVQ==\"", message, "enigma")
}

func TestFirePathPost(t *testing.T) {

	AssertSuccessFirePost(t, "[1,2,3]", []int{1, 2, 3})
}
