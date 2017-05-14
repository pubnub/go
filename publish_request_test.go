package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var pubnub *PubNub

func init() {
	pnconfig := NewPNConfiguration()
	pnconfig.PublishKey = "my_pub_key"
	pnconfig.SubscribeKey = "my_sub_key"
	pnconfig.SecretKey = "my_secret_key"

	pubnub = NewPubNub(pnconfig)
}

func TestPublishInitializer(t *testing.T) {
	assert := assert.New(t)

	publish := pubnub.Publish()
	assert.Implements((*Endpoint)(nil), publish)
	assert.IsType((*Publish)(nil), publish)
}

func TestPublishBuilder(t *testing.T) {
	assert := assert.New(t)

	publish := pubnub.Publish()
	publish.Channel = "news"
	publish.Message = "hello!"

	assert.Implements((*Endpoint)(nil), publish)
	assert.IsType((*Publish)(nil), publish)

	assert.Equal("/publish/my_pub_key/my_sub_key/0/news/0/hello!",
		publish.buildPath())
}
