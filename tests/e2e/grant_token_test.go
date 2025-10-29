package e2e

import (
	//"fmt"

	"log"
	"os"
	"testing"

	pubnub "github.com/pubnub/go/v8"
	"github.com/stretchr/testify/assert"
)

func TestGrantToken(t *testing.T) {

	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	ch1 := randomized("channel1")
	ch := map[string]pubnub.ChannelPermissions{
		ch1: {
			Read:   true,
			Write:  true,
			Delete: false,
		},
	}

	cg1 := randomized("group1")
	cg2 := randomized("group2")
	cg := map[string]pubnub.GroupPermissions{
		cg1: {
			Read:   true,
			Manage: true,
		},
		cg2: {
			Read:   true,
			Manage: false,
		},
	}

	res, _, err := pn.GrantToken().TTL(10).
		Channels(ch).
		ChannelGroups(cg).
		Execute()

	assert.Nil(err)

	assert.NotNil(res)
	if res != nil {
		token := res.Data.Token
		cborObject, err := pubnub.GetPermissions(token)
		if err == nil {
			chResources := pubnub.ParseGrantResources(cborObject.Resources, token, cborObject.Timestamp, cborObject.TTL)

			assert.Equal(ch[ch1], chResources.Channels[ch1].Permissions)
			assert.Equal(cg[cg1], chResources.Groups[cg1].Permissions)
			assert.Equal(cg[cg2], chResources.Groups[cg2].Permissions)
		}

	}

}

func TestGrantTokenEntities(t *testing.T) {

	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	space := "space"
	user := "user"
	spacePattern := "space*"
	userPattern := "user*"

	expectedTokenResourcePermissions := pubnub.PNTokenResources{
		Channels: map[string]pubnub.ChannelPermissions{
			space: {
				Read:  true,
				Write: true,
			},
		},
		ChannelGroups: map[string]pubnub.GroupPermissions{},
		UUIDs: map[string]pubnub.UUIDPermissions{
			user: {
				Get:    true,
				Update: true,
			},
		},
	}

	expectedTokenPatternPermissions := pubnub.PNTokenResources{
		Channels: map[string]pubnub.ChannelPermissions{
			spacePattern: {
				Read:  true,
				Write: true,
			},
		},
		ChannelGroups: map[string]pubnub.GroupPermissions{},
		UUIDs: map[string]pubnub.UUIDPermissions{
			userPattern: {
				Get:    true,
				Update: true,
			},
		},
	}

	res, _, err := pn.GrantToken().TTL(10).
		SpacesPermissions(map[pubnub.SpaceId]pubnub.SpacePermissions{
			pubnub.SpaceId(space): {
				Read:  true,
				Write: true,
			},
		}).
		UsersPermissions(map[pubnub.UserId]pubnub.UserPermissions{
			pubnub.UserId(user): {
				Get:    true,
				Delete: false,
				Update: true,
			},
		}).
		SpacePatternsPermissions(map[string]pubnub.SpacePermissions{
			spacePattern: {
				Read:  true,
				Write: true,
			},
		}).
		UserPatternsPermissions(map[string]pubnub.UserPermissions{
			userPattern: {
				Get:    true,
				Delete: false,
				Update: true,
			},
		}).
		Execute()

	assert.Nil(err)

	assert.NotNil(res)
	parsedToken, err := pubnub.ParseToken(res.Data.Token)

	assert.Nil(err)

	assert.Equal(expectedTokenResourcePermissions, parsedToken.Resources)
	assert.Equal(expectedTokenPatternPermissions, parsedToken.Patterns)
}
