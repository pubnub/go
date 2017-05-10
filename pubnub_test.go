package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializer(t *testing.T) {
	assert := assert.New(t)

	pnconfig := NewPNConfiguration()
	pnconfig.PublishKey = "my_pub_key"
	pnconfig.SubscribeKey = "my_sub_key"
	pnconfig.SecretKey = "my_secret_key"

	pubnub := NewPubNub(pnconfig)

	assert.Equal("my_pub_key", pubnub.PNConfig.PublishKey)
	assert.Equal("my_sub_key", pubnub.PNConfig.SubscribeKey)
	assert.Equal("my_secret_key", pubnub.PNConfig.SecretKey)
}

func TestDemoInitializer(t *testing.T) {
	demo := NewPubNubDemo()

	assert := assert.New(t)

	assert.Equal("demo", demo.PNConfig.PublishKey)
	assert.Equal("demo", demo.PNConfig.SubscribeKey)
	assert.Equal("demo", demo.PNConfig.SecretKey)
}
