package pubnub

import (
	"github.com/pubnub/go/v7/crypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializer(t *testing.T) {
	assert := assert.New(t)

	pnconfig := NewConfigWithUserId(UserId(GenerateUUID()))
	pnconfig.PublishKey = "my_pub_key"
	pnconfig.SubscribeKey = "my_sub_key"
	pnconfig.SecretKey = "my_secret_key"
	pubnub := NewPubNub(pnconfig)

	assert.Equal("my_pub_key", pubnub.Config.PublishKey)
	assert.Equal("my_sub_key", pubnub.Config.SubscribeKey)
	assert.Equal("my_secret_key", pubnub.Config.SecretKey)
}

func TestCryptoModuleChangesWhenCipherKeyChanges(t *testing.T) {
	pubnub := NewPubNubDemo()

	cryptoModule, _ := crypto.NewAesCbcCryptoModule("my_cipher_key", true)
	pubnub.Config.CryptoModule = cryptoModule
	a := assert.New(t)

	a.Equal(cryptoModule, pubnub.getCryptoModule())

	pubnub.Config.CipherKey = "new_cipher_key"
	a.NotEqual(cryptoModule, pubnub.getCryptoModule())
}

func TestCryptoModuleChangesIfIvFlagChanges(t *testing.T) {
	pubnub := NewPubNubDemo()

	cryptoModule, _ := crypto.NewAesCbcCryptoModule("my_cipher_key", true)
	pubnub.Config.CryptoModule = cryptoModule
	a := assert.New(t)

	a.Equal(cryptoModule, pubnub.getCryptoModule())

	pubnub.Config.UseRandomInitializationVector = false
	a.NotEqual(cryptoModule, pubnub.getCryptoModule())
}

func TestDemoInitializer(t *testing.T) {
	demo := NewPubNubDemo()

	assert := assert.New(t)

	assert.Equal("demo", demo.Config.PublishKey)
	assert.Equal("demo", demo.Config.SubscribeKey)
	assert.Equal("demo", demo.Config.SecretKey)
}

func TestMultipleConcurrentInit(t *testing.T) {
	c1 := NewConfigWithUserId(UserId(GenerateUUID()))
	go NewPubNub(c1)
	c2 := NewConfigWithUserId(UserId(GenerateUUID()))
	NewPubNub(c2)
}
