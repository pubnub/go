package e2e

import (
	"fmt"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestGrantToken(t *testing.T) {

	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	// pn.Config.Origin = "ingress.bronze.aws-pdx-1.ps.pn"
	// pn.Config.Secure = false
	// pn.Config.PublishKey = "pub-c-03f156ea-a2e3-4c35-a733-9535824be897"
	// pn.Config.SubscribeKey = "sub-c-d7da9e58-c997-11e9-a139-dab2c75acd6f"
	// pn.Config.SecretKey = "sec-c-MmUxNTZjMmYtNzFkNS00ODkzLWE2YjctNmQ4YzE5NWNmZDA3"

	// ch1 := randomized("channel")
	// cg1 := randomized("cg")
	// cg2 := randomized("cg")
	u1 := randomized("u")
	s1 := randomized("s")

	// ch := map[string]pubnub.ChannelPermissions{
	// 	ch1: pubnub.ChannelPermissions{
	// 		Read:   true,
	// 		Write:  true,
	// 		Delete: false,
	// 	},
	// }

	u := map[string]pubnub.UserSpacePermissions{
		u1: pubnub.UserSpacePermissions{
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
			Create: false,
		},
	}

	s := map[string]pubnub.UserSpacePermissions{
		s1: pubnub.UserSpacePermissions{
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
			Create: true,
		},
	}

	// cg := map[string]pubnub.GroupPermissions{
	// 	cg1: pubnub.GroupPermissions{
	// 		Read:   true,
	// 		Manage: true,
	// 	},
	// 	cg2: pubnub.GroupPermissions{
	// 		Read:   true,
	// 		Manage: false,
	// 	},
	// }

	res, _, err := pn.GrantToken().TTL(10).
		//Channels(ch).
		//ChannelGroups(cg).
		Users(u).
		Spaces(s).
		Execute()

	assert.Nil(err)

	//fmt.Println(res)
	assert.NotNil(res)
	if res != nil {
		token := res.Data.Token
		//token = "p0F2AkF0Gl043rhDdHRsCkNyZXOkRGNoYW6hZnNlY3JldAFDZ3JwoEN1c3KgQ3NwY6BDcGF0pERjaGFuoENncnCgQ3VzcqBDc3BjoERtZXRhoENzaWdYIGOAeTyWGJI-blahPGD9TuKlaW1YQgiB4uR_edmfq-61"
		//map[pat:map[usr:map[] spc:map[] chan:map[] grp:map[]] meta:map[] sig:[205 161 131 38 100 38 57 220 2 234 208 130 204 167 117 48 224 91 132 70 12 192 211 34 47 43 64 188 207 118 55 110] v:2 t:1567502256 ttl:10 res:map[grp:map[cg-1623328:23 cg-6488712:19] usr:map[u-3244801:15] spc:map[s-8225817:31] chan:map[channel-7076766:7]]]
		cborObject, err := pubnub.GetPermissions(token)
		if err == nil {
			// fmt.Printf("\nCBOR decode Token---> %#v", cborObject)
			// fmt.Println("")
			// fmt.Println("Sig: ", string(cborObject.Signature))
			// fmt.Println("Version: ", cborObject.Version)
			// fmt.Println("Timestamp: ", cborObject.Timestamp)
			// fmt.Println("TTL: ", cborObject.TTL)
			// fmt.Println(fmt.Sprintf("Meta: %#v", cborObject.Meta))
			// fmt.Println("")
			// fmt.Println(" --- Resources")
			chResources := pubnub.ParseGrantResources(cborObject.Resources, token, cborObject.Timestamp, cborObject.TTL)

			fmt.Println(chResources)

			// assert.Equal(ch[ch1].Read, chResources.Channels[ch1].Permissions.Read)
			// assert.Equal(ch[ch1].Write, chResources.Channels[ch1].Permissions.Write)
			// //assert.Equal(ch[ch1].Manage, chResources.Channels[ch1].Permissions.Manage)
			// assert.Equal(ch[ch1].Delete, chResources.Channels[ch1].Permissions.Delete)
			// //assert.Equal(ch[ch1].Create, chResources.Channels[ch1].Permissions.Create)

			assert.Equal(u[u1].Read, chResources.Users[u1].Permissions.Read)
			assert.Equal(u[u1].Write, chResources.Users[u1].Permissions.Write)
			assert.Equal(u[u1].Manage, chResources.Users[u1].Permissions.Manage)
			assert.Equal(u[u1].Delete, chResources.Users[u1].Permissions.Delete)
			assert.Equal(u[u1].Create, chResources.Users[u1].Permissions.Create)

			assert.Equal(s[s1].Read, chResources.Spaces[s1].Permissions.Read)
			assert.Equal(s[s1].Write, chResources.Spaces[s1].Permissions.Write)
			assert.Equal(s[s1].Manage, chResources.Spaces[s1].Permissions.Manage)
			assert.Equal(s[s1].Delete, chResources.Spaces[s1].Permissions.Delete)
			assert.Equal(s[s1].Create, chResources.Spaces[s1].Permissions.Create)

			// fmt.Println(cg1, cg[cg1], chResources.Groups[cg1])
			// assert.Equal(cg[cg1].Read, chResources.Groups[cg1].Permissions.Read)
			// //assert.Equal(cg[cg1].Write, chResources.Groups[cg1].Permissions.Write)
			// assert.Equal(cg[cg1].Manage, chResources.Groups[cg1].Permissions.Manage)
			// //assert.Equal(cg[cg1].Delete, chResources.Groups[cg1].Permissions.Delete)
			// //assert.Equal(cg[cg1].Create, chResources.Groups[cg1].Permissions.Create)

			// fmt.Println(cg2, cg[cg2], chResources.Groups[cg2])
			// assert.Equal(cg[cg2].Read, chResources.Groups[cg2].Permissions.Read)
			// //assert.Equal(cg[cg2].Write, chResources.Groups[cg2].Permissions.Write)
			// assert.Equal(cg[cg2].Manage, chResources.Groups[cg2].Permissions.Manage)
			//assert.Equal(cg[cg2].Delete, chResources.Groups[cg2].Permissions.Delete)
			//assert.Equal(cg[cg2].Create, chResources.Groups[cg2].Permissions.Create)

			//fmt.Println(" --- Patterns")
			pubnub.ParseGrantResources(cborObject.Patterns, token, cborObject.Timestamp, cborObject.TTL)
		}

		// t2 := pn.TokenManager.GetTokens([]string{ch1}, nil, nil, nil)
		// fmt.Println(t2.Channels[ch1].Token)

		// err := json.Unmarshal(value, &resp)
		// if err != nil {
		// 	fmt.Printf("\nUnmarshal Error---> %#v", err)
		// 	//log.Println("Write file:", ioutil.WriteFile("data.json", value, 0600))
		// } else {
		// 	fmt.Println("\nUnmarshalled Token--->", resp)
		// }
		//}

	}

}
