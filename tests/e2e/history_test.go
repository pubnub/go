package e2e

import (
	"log"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestHistorySuccessNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	// for i := 0; i <= 10; i++ {
	// 	_, err := pn.Publish(&pubnub.PublishOpts{
	// 		Channel:    "ch",
	// 		Message:    fmt.Sprintf("hey%d", i),
	// 		DoNotStore: true,
	// 	})
	//
	// 	assert.Nil(err)
	// }

	res, err := pn.History(&pubnub.HistoryOpts{
		Channel: "ch",
	})
	assert.Nil(err)

	log.Println(res)
}
