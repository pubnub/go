package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v5/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSendFile(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newSendFileBuilder(pn)
	if testContext {
		o = newSendFileBuilderWithContext(pn, backgroundContext)
	}

	channel := "chan"
	o.Channel(channel)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(sendFilePath, pn.Config.SubscribeKey, channel),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{123, 34, 110, 97, 109, 101, 34, 58, 34, 34, 125}, body)

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

}

func TestSendFile(t *testing.T) {
	AssertSendFile(t, true, false)
}

func TestSendFileContext(t *testing.T) {
	AssertSendFile(t, true, true)
}

func TestSendFileResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &sendFileOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNSendFileResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}
