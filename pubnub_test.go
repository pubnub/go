package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializer(t *testing.T) {
	assert := assert.New(t)

	pnconfig := NewConfig()
	pnconfig.PublishKey = "my_pub_key"
	pnconfig.SubscribeKey = "my_sub_key"
	pnconfig.SecretKey = "my_secret_key"

	pubnub := NewPubNub(pnconfig)

	assert.Equal("my_pub_key", pubnub.Config.PublishKey)
	assert.Equal("my_sub_key", pubnub.Config.SubscribeKey)
	assert.Equal("my_secret_key", pubnub.Config.SecretKey)
}

func TestDemoInitializer(t *testing.T) {
	demo := NewPubNubDemo()

	assert := assert.New(t)

	assert.Equal("demo", demo.Config.PublishKey)
	assert.Equal("demo", demo.Config.SubscribeKey)
	assert.Equal("demo", demo.Config.SecretKey)
}

func TestMultipleConcurrentInit(t *testing.T) {
	go NewPubNub(NewConfig())
	NewPubNub(NewConfig())
}
