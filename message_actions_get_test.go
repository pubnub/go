package pubnub

import (
	"fmt"
	"strconv"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertGetMessageActions(t *testing.T, checkQueryParam, testContext bool) {

	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newGetMessageActionsBuilder(pn)
	if testContext {
		o = newGetMessageActionsBuilderWithContext(pn, backgroundContext)
	}

	channel := "chan"
	timetoken := "15698453963258802"
	aTimetoken := "15692384791344400"
	limit := 10
	o.Channel(channel)
	o.Start(timetoken)
	o.End(aTimetoken)
	o.Limit(limit)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(getMessageActionsPath, pn.Config.SubscribeKey, channel),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
		assert.Equal(timetoken, u.Get("start"))
		assert.Equal(aTimetoken, u.Get("end"))
		assert.Equal(strconv.Itoa(limit), u.Get("limit"))
	}

}

func TestGetMessageActions(t *testing.T) {
	AssertGetMessageActions(t, true, false)
}

func TestGetMessageActionsContext(t *testing.T) {
	AssertGetMessageActions(t, true, true)
}

func TestGetMessageActionsResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getMessageActionsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetMessageActionsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestGetMessageActionsResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getMessageActionsOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status": 200, "data": [{"messageTimetoken": "15698466245557325", "type": "reaction", "uuid": "pn-85463c27-ad24-49d4-8cdf-db93a300855a", "value": "smiley_face", "actionTimetoken": "15698466249528820"}]}`)

	r, _, err := newPNGetMessageActionsResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("15698466245557325", r.Data[0].MessageTimetoken)
	assert.Equal("reaction", r.Data[0].ActionType)
	assert.Equal("smiley_face", r.Data[0].ActionValue)
	assert.Equal("15698466249528820", r.Data[0].ActionTimetoken)
	assert.Equal("pn-85463c27-ad24-49d4-8cdf-db93a300855a", r.Data[0].UUID)

	assert.Nil(err)
}
