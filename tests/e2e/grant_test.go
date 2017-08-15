package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestGrantSucccessNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	_, err := pn.Grant(&pubnub.GrantOpts{
		Channels: []string{"ch"},
		Write:    true,
		Read:     true,
	})

	assert.Nil(err)
}

func TestGrantMultipleMixed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	res, err := pn.Grant(&pubnub.GrantOpts{
		Channels: []string{"ch1", "ch2", "ch3"},
		Groups:   []string{"cg1", "cg2", "cg3"},
		Write:    true,
		Read:     true,
		Manage:   true,
		AuthKeys: []string{"my-auth-key-1", "my-auth-key-2"},
	})

	assert.Nil(err)
	assert.NotNil(res)
}

func TestGrantSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()

	config.Uuid = SPECIAL_CHARACTERS
	config.AuthKey = SPECIAL_CHARACTERS

	pn := pubnub.NewPubNub(config)

	_, err := pn.Grant(&pubnub.GrantOpts{
		Channels: []string{SPECIAL_CHANNEL},
		Groups:   []string{SPECIAL_CHANNEL},
		Write:    true,
		Read:     true,
		Manage:   true,
		AuthKeys: []string{SPECIAL_CHARACTERS},
	})

	assert.Nil(err)
}
