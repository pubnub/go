package pubnub

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestNewSetStateBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newSetStateBuilder(pubnub)
	o.Channels([]string{"ch1", "ch2", "ch3"})

	path, err := o.opts.buildPath()
	assert.Nil(err)

	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/uuid/%s/data",
			o.opts.pubnub.Config.UUID),
		u.EscapedPath(), []int{})
}

func TestNewSetStateBuilderWithUUID(t *testing.T) {
	assert := assert.New(t)

	o := newSetStateBuilder(pubnub)
	o.Channels([]string{"ch1", "ch2", "ch3"})
	uuid := "customuuid"
	o.UUID(uuid)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/uuid/%s/data",
			uuid),
		u.EscapedPath(), []int{})
}

func TestNewSetStateBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newSetStateBuilderWithContext(pubnub, backgroundContext)
	o.Channels([]string{"ch1", "ch2", "ch3"})

	path, err := o.opts.buildPath()
	assert.Nil(err)

	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/uuid/%s/data",
			o.opts.pubnub.Config.UUID),
		u.EscapedPath(), []int{})
}

func TestNewSetStateResponse(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.UUID = "my-custom-uuid"

	jsonBytes := []byte(`{"status": 200, "message": "OK", "payload": {"k": "v"}, "uuid": "my-custom-uuid", "service": "Presence"}`)

	res, _, err := newSetStateResponse(jsonBytes, fakeResponseState)
	assert.Nil(err)
	if s, ok := res.State.(map[string]interface{}); ok {
		assert.Equal("v", s["k"])
	} else {
		assert.Fail("!map[string]interface{}")
	}
}

func TestSetStateRequestBasic(t *testing.T) {
	assert := assert.New(t)
	state := make(map[string]interface{})
	state["name"] = "Alex"
	state["count"] = 5

	opts := &setStateOpts{
		Channels:      []string{"ch"},
		ChannelGroups: []string{"cg"},
		State:         state,
		pubnub:        pubnub,
	}

	err := opts.validate()
	assert.Nil(err)

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/%s/uuid/%s/data",
			opts.Channels[0], opts.pubnub.Config.UUID),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")
	expected.Set("state", `{"count":5,"name":"Alex"}`)
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)

}

func TestSetStateMultipleChannels(t *testing.T) {
	assert := assert.New(t)

	opts := &setStateOpts{
		Channels: []string{"ch1", "ch2", "ch3"},
		pubnub:   pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)

	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/uuid/%s/data",
			opts.pubnub.Config.UUID),
		u.EscapedPath(), []int{})
}

func TestSetStateMultipleChannelGroups(t *testing.T) {
	assert := assert.New(t)

	opts := &setStateOpts{
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		pubnub:        pubnub,
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}

func TestSetStateMultipleChannelGroupsQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := &setStateOpts{
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		pubnub:        pubnub,
	}
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}

func TestSetStateValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &setStateOpts{
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		pubnub:        pn,
	}

	assert.Equal("pubnub/validation: pubnub: Set State: Missing Subscribe Key", opts.validate().Error())
}

func TestSetStateValidateCG(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &setStateOpts{
		pubnub: pn,
	}

	assert.Equal("pubnub/validation: pubnub: Set State: Missing Channel or Channel Group", opts.validate().Error())
}

func TestSetStateValidateState(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &setStateOpts{
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		pubnub:        pn,
	}

	assert.Equal("pubnub/validation: pubnub: Set State: Missing State", opts.validate().Error())
}

func TestNewSetStateResponseErrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newSetStateResponse(jsonBytes, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestNewSetStateResponseValueError(t *testing.T) {
	assert := assert.New(t)
	state := make(map[string]interface{})
	state["name"] = "Alex"
	state["error"] = 5
	b, err1 := json.Marshal(state)
	if err1 != nil {
		panic(err1)
	}

	_, _, err := newSetStateResponse([]byte(b), StatusResponse{})
	assert.Equal("", err.Error())
}
