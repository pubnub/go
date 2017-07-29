package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func init() {
	pnconfig = NewConfig()

	pnconfig.PublishKey = "pub_key"
	pnconfig.SubscribeKey = "sub_key"
	pnconfig.SecretKey = "secret_key"

	pubnub = NewPubNub(pnconfig)
}

func TestSetStateRequestBasic(t *testing.T) {
	assert := assert.New(t)
	state := make(map[string]interface{})
	state["name"] = "Alex"
	state["count"] = 5

	opts := &SetStateOpts{
		Channels:      []string{"ch"},
		ChannelGroups: []string{"cg"},
		State:         state,
		pubnub:        pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/%s/uuid/%s/data",
			opts.Channels[0], opts.pubnub.Config.Uuid),
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

	opts := &SetStateOpts{
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
			opts.pubnub.Config.Uuid),
		u.EscapedPath(), []int{})
}

func TestSetStateMultipleChannelGroups(t *testing.T) {
	assert := assert.New(t)

	opts := &SetStateOpts{
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		pubnub:        pubnub,
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}
