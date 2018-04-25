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

	assert.Equal("APNS", pushAPNS.String())
	assert.Equal("MPNS", pushMPNS.String())
	assert.Equal("GCM", pushGCM.String())
	assert.Equal("NONE", pushNONE.String())
}
