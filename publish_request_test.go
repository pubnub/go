package pubnub

import (
	"log"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

var pnconfig *Config
var pubnub *PubNub

func init() {
	pnconfig = NewConfig()

	pnconfig.PublishKey = "pub_key"
	pnconfig.SubscribeKey = "sub_key"
	pnconfig.SecretKey = "secret_key"

	pubnub = NewPubNub(pnconfig)
}

func TestSimplePublish(t *testing.T) {
	assert := assert.New(t)

	opts := &PublishOpts{
		Channel: "ch",
		Message: "hey",
		pubnub:  pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	h.AssertPathsEqual(t, "/publish/pub_key/sub_key/0/ch/0/\"hey\"", path, []int{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}

func TestPublishSliceMessage(t *testing.T) {
	assert := assert.New(t)

	message := []string{"hey1", "hey2", "hey3"}

	opts := &PublishOpts{
		Channel: "ch",
		Message: message,
		pubnub:  pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	log.Println(path)
	h.AssertPathsEqual(t,
		"/publish/pub_key/sub_key/0/ch/0/%5B%22hey1%22%2C%22hey2%22%2C%22hey3%22%5D",
		path, []int{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}
