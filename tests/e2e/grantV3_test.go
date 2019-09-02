package e2e

import (
	"fmt"
	"log"
	"os"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestGrantV3(t *testing.T) {

	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	pn.Config.Origin = "ingress.bronze.aws-pdx-1.ps.pn"
	pn.Config.Secure = false

	ch1 := randomized("channel")
	cg1 := randomized("cg")
	cg2 := randomized("cg")
	u1 := randomized("u")
	s1 := randomized("s")

	ch := map[string]pubnub.ResourcePermissions{
		ch1: pubnub.ResourcePermissions{
			Create: false,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: false,
		},
	}

	u := map[string]pubnub.ResourcePermissions{
		u1: pubnub.ResourcePermissions{
			Create: false,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
		},
	}

	s := map[string]pubnub.ResourcePermissions{
		s1: pubnub.ResourcePermissions{
			Create: true,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
		},
	}

	cg := map[string]pubnub.ResourcePermissions{
		cg1: pubnub.ResourcePermissions{
			Create: true,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: false,
		},
		cg2: pubnub.ResourcePermissions{
			Create: true,
			Read:   true,
			Write:  true,
			Manage: false,
			Delete: false,
		},
	}

	res, _, err := pn.Grant().TTL(10).
		Channels(ch).
		ChannelGroups(cg).
		Users(u).
		Spaces(s).
		Execute()

	assert.Nil(err)
	fmt.Println(res)
	assert.NotNil(res)

}
