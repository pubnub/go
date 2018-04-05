package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestHistoryDeleteNotStubbed(t *testing.T) {
	assert := assert.New(t)

	ch := randomized("h-ch")
	pn := pubnub.NewPubNub(pamConfigCopy())

	_, _, err := pn.DeleteMessages().
		Channel(ch).
		Execute()

	assert.Nil(err)
}

func TestHistoryDeleteMissingChannelError(t *testing.T) {
	assert := assert.New(t)

	/*config.PublishKey = "pub-c-4f1dbd79-ab94-487d-b779-5881927db87c"
	config.SubscribeKey = "sub-c-f2489488-2dbd-11e8-a27a-a2b5bab5b996"
	config.SecretKey = "sec-c-NjlmYzVkMjEtOWIxZi00YmJlLThjZDktMjI4NGQwZDUxZDQ0"*/

	config2 := pamConfigCopy()

	pn := pubnub.NewPubNub(config2)

	res, _, err := pn.DeleteMessages().
		Channel("").
		Execute()

	assert.Nil(res)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestHistoryDeleteSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()

	// Not allowed characters: /?#,
	validCharacters := "-._~:[]@!$&'()*+;=`|"

	config.Uuid = validCharacters
	config.AuthKey = validCharacters

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.DeleteMessages().
		Channel(validCharacters).
		Start(int64(123)).
		End(int64(456)).
		Execute()

	assert.Nil(err)
}
