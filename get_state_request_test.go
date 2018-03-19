package pubnub

import (
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetStateBasicRequest(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.Uuid = "my-custom-uuid"

	opts := &getStateOpts{
		Channels:      []string{"ch"},
		ChannelGroups: []string{"cg"},
		pubnub:        pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/ch/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestGetStateMultipleChannelsChannelGroups(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.Uuid = "my-custom-uuid"

	opts := &getStateOpts{
		Channels:      []string{"ch1", "ch2", "ch3"},
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		pubnub:        pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}
