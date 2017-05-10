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
}
