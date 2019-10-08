package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertRemoveMessageActions(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newRemoveMessageActionsBuilder(pn)
	if testContext {
		o = newRemoveMessageActionsBuilderWithContext(pn, backgroundContext)
	}

	channel := "chan"
	timetoken := "15698453963258802"
	aTimetoken := "15692384791344400"
	o.Channel(channel)
	o.MessageTimetoken(timetoken)
	o.ActionTimetoken(aTimetoken)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(removeMessageActionsPath, pn.Config.SubscribeKey, channel, timetoken, aTimetoken),
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

func TestRemoveMessageActions(t *testing.T) {
	AssertRemoveMessageActions(t, true, false)
}

func TestRemoveMessageActionsContext(t *testing.T) {
	AssertRemoveMessageActions(t, true, true)
}

func TestRemoveMessageActionsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &removeMessageActionsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNRemoveMessageActionsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRemoveMessageActionsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &removeMessageActionsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status": 200, "data": {}}`)

	r, _, err := newPNRemoveMessageActionsResponse(jsonBytes, opts, StatusResponse{})
	assert.Empty(r.Data)

	assert.Nil(err)
}
