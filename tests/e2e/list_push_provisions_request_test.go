package e2e

import (
	//"log"
	//"os"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestListPushProvisionsNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	resp, _, err := pn.ListPushProvisions().
		DeviceIDForPush("cg").
		PushType(pubnub.PNPushTypeGCM).
		Execute()
	assert.Equal("ch", resp.Channels[0])
	assert.Nil(err)
}
