package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushString(t *testing.T) {
	assert := assert.New(t)

	pushAPNS := PNPushTypeAPNS
	pushMPNS := PNPushTypeMPNS
	pushGCM := PNPushTypeGCM
	pushNONE := PNPushTypeNone

	assert.Equal("apns", pushAPNS.String())
	assert.Equal("mpns", pushMPNS.String())
	assert.Equal("gcm", pushGCM.String())
	assert.Equal("none", pushNONE.String())
}
