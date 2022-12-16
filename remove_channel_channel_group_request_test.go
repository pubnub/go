package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func init() {
	pnconfig = NewConfigWithUserId(UserId(GenerateUUID()))

	pnconfig.PublishKey = "pub_key"
	pnconfig.SubscribeKey = "sub_key"
	pnconfig.SecretKey = "secret_key"

	pubnub = NewPubNub(pnconfig)
}

func TestRemoveChannelRequestBasic(t *testing.T) {
	assert := assert.New(t)

	opts := newRemoveChannelOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroup = "cg"

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg"),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("remove", "ch1,ch2,ch3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestRemoveChannelRequestBasicQueryParam(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := newRemoveChannelOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch1", "ch2", "ch3"}
	opts.ChannelGroup = "cg"
	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("remove", "ch1,ch2,ch3")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewRemoveChannelFromChannelGroupBuilder(t *testing.T) {
	assert := assert.New(t)
	o := newRemoveChannelFromChannelGroupBuilder(pubnub)
	o.ChannelGroup("cg")
	o.Channels([]string{"ch1", "ch2", "ch3"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg"),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("remove", "ch1,ch2,ch3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewRemoveChannelFromChannelGroupBuilderContext(t *testing.T) {
	assert := assert.New(t)
	o := newRemoveChannelFromChannelGroupBuilderWithContext(pubnub, backgroundContext)
	o.ChannelGroup("cg")
	o.Channels([]string{"ch1", "ch2", "ch3"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg"),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("remove", "ch1,ch2,ch3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestRemChannelsFromCGValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newRemoveChannelOpts(pn, pn.ctx)

	assert.Equal("pubnub/validation: pubnub: Remove Channel From Channel Group: Missing Subscribe Key", opts.validate().Error())
}
