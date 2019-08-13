package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessSignalGet(t *testing.T, channel string, checkQueryParam bool) {
	assert := assert.New(t)
	msgMap := make(map[string]string)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	msgMap["one"] = "hey1"

	opts := &signalOpts{
		Channel:    channel,
		Message:    msgMap,
		pubnub:     pubnub,
		QueryParam: queryParam,
	}

	path, err := opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/signal/%s/%s/0/%s/0/%s", pubnub.Config.PublishKey, pubnub.Config.SubscribeKey, channel, "%7B%22one%22%3A%22hey1%22%7D"),
		path, []int{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Empty(body)

	if checkQueryParam {
		u, _ := opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

}

func TestSignalPath(t *testing.T) {
	channels := "test1"
	AssertSuccessSignalGet(t, channels, false)
}

func TestSignalPathQP(t *testing.T) {
	channels := "test1"
	AssertSuccessSignalGet(t, channels, true)
}

func AssertNewSignalBuilder(t *testing.T, checkQueryParam bool, testContext bool, channel string) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	msgMap := make(map[string]string)

	if !checkQueryParam {
		queryParam = nil
	}

	msgMap["one"] = "hey1"
	expectedBody := "{\"one\":\"hey1\"}"

	o := newSignalBuilder(pubnub)
	if testContext {
		o = newSignalBuilderWithContext(pubnub, backgroundContext)
	}
	o.Channel(channel)
	o.Message(msgMap)
	o.usePost(true)
	if checkQueryParam {
		o.QueryParam(queryParam)
	}

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/signal/%s/%s/0/%s/0", pubnub.Config.PublishKey, pubnub.Config.SubscribeKey, channel),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Equal(expectedBody, string(body))

	u, _ := o.opts.buildQuery()

	if checkQueryParam {
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

}

func TestSignalBuilder(t *testing.T) {
	channels := "test1"
	AssertNewSignalBuilder(t, false, false, channels)
}

func TestSignalBuilderQP(t *testing.T) {
	channels := "test1"
	AssertNewSignalBuilder(t, true, false, channels)
}

func TestSignalBuilderContext(t *testing.T) {
	channels := "test1"
	AssertNewSignalBuilder(t, true, true, channels)
}

func TestSignalBuilderContextQP(t *testing.T) {
	channels := "test1"
	AssertNewSignalBuilder(t, true, true, channels)
}

func TestSignalResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &signalOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newSignalResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

//[1, "Sent", "1232423423423"]
func TestSignalResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &signalOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`[1, "Sent", "1232423423423"]`)

	_, _, err := newSignalResponse(jsonBytes, opts, StatusResponse{})
	assert.Nil(err)
}
