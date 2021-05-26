package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v5/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertAddMessageActions(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newAddMessageActionsBuilder(pn)
	if testContext {
		o = newAddMessageActionsBuilderWithContext(pn, backgroundContext)
	}

	ma := MessageAction{
		ActionType:  "action",
		ActionValue: "smiley",
	}

	channel := "chan"
	timetoken := "15698453963258802"
	o.Channel(channel)
	o.MessageTimetoken(timetoken)
	o.Action(ma)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(addMessageActionsPath, pn.Config.SubscribeKey, channel, timetoken),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	expectedBody := "{\"type\":\"action\",\"value\":\"smiley\"}"

	assert.Equal(expectedBody, string(body))

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

}

func TestAddMessageActions(t *testing.T) {
	AssertAddMessageActions(t, true, false)
}

func TestAddMessageActionsContext(t *testing.T) {
	AssertAddMessageActions(t, true, true)
}

func TestAddMessageActionsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &addMessageActionsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNAddMessageActionsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestAddMessageActionsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &addMessageActionsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status": 200, "data": {"messageTimetoken": "15210190573608384", "type": "reaction", "uuid": "pn-871b8325-a11f-48cb-9c15-64984790703e", "value": "smiley_face", "actionTimetoken": "15692384791344400"}}`)

	r, _, err := newPNAddMessageActionsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("15210190573608384", r.Data.MessageTimetoken)
	assert.Equal("reaction", r.Data.ActionType)
	assert.Equal("smiley_face", r.Data.ActionValue)
	assert.Equal("15692384791344400", r.Data.ActionTimetoken)
	assert.Equal("pn-871b8325-a11f-48cb-9c15-64984790703e", r.Data.UUID)

	assert.Nil(err)
}
