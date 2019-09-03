package e2e

import (
	"bytes"
	"encoding/base64"
	//"encoding/hex"
	//"encoding/binary"
	//"strconv"
	//"encoding/json"
	"fmt"
	cbor "github.com/brianolson/cbor_go"
	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
	"testing"
)

func TestGrantV3(t *testing.T) {

	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	pn.Config.Origin = "ingress.bronze.aws-pdx-1.ps.pn"
	pn.Config.Secure = false
	pn.Config.PublishKey = "pub-c-03f156ea-a2e3-4c35-a733-9535824be897"
	pn.Config.SubscribeKey = "sub-c-d7da9e58-c997-11e9-a139-dab2c75acd6f"
	pn.Config.SecretKey = "sec-c-MmUxNTZjMmYtNzFkNS00ODkzLWE2YjctNmQ4YzE5NWNmZDA3"

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

	//fmt.Println(res)
	token := res.Data.Token
	//token = "p0F2AkF0Gl043rhDdHRsCkNyZXOkRGNoYW6hZnNlY3JldAFDZ3JwoEN1c3KgQ3NwY6BDcGF0pERjaGFuoENncnCgQ3VzcqBDc3BjoERtZXRhoENzaWdYIGOAeTyWGJI-blahPGD9TuKlaW1YQgiB4uR_edmfq-61"
	token = strings.Replace(token, "-", "+", -1)
	token = strings.Replace(token, "_", "/", -1)
	fmt.Println("\nStrings `-`, `_` replaced Token--->", token)
	// if i := len(s) % 4; i != 0 {
	// 	token += strings.Repeat("=", 4-i)
	// }
	// fmt.Println("Padded Token--->", token)
	value, decodeErr := base64.StdEncoding.DecodeString(token)
	if decodeErr != nil {
		fmt.Println("\nDecoding Error--->", decodeErr)
	} else {
		fmt.Println("\nDecoded Token--->", string(value))
		//var resp interface{}
		c := cbor.NewDecoder(bytes.NewReader(value))
		var cborObject pubnub.PNGrantTokenDecoded
		err1 := c.Decode(&cborObject)
		if err1 != nil {
			fmt.Printf("\nCBOR decode Error---> %#v", err1)
			//log.Println("Write file:", ioutil.WriteFile("data.json", value, 0600))
		} else {
			//map[pat:map[usr:map[] spc:map[] chan:map[] grp:map[]] meta:map[] sig:[205 161 131 38 100 38 57 220 2 234 208 130 204 167 117 48 224 91 132 70 12 192 211 34 47 43 64 188 207 118 55 110] v:2 t:1567502256 ttl:10 res:map[grp:map[cg-1623328:23 cg-6488712:19] usr:map[u-3244801:15] spc:map[s-8225817:31] chan:map[channel-7076766:7]]]
			fmt.Printf("\nCBOR decode Token---> %v", cborObject)
			fmt.Println("")
			fmt.Println("Sig: ", string(cborObject.Signature))
			fmt.Println("Version: ", cborObject.Version)
			fmt.Println("Timetoken: ", cborObject.Timetoken)
			fmt.Println("TTL: ", cborObject.TTL)
			fmt.Println(fmt.Sprintf("Meta: %#v", cborObject.Meta))
			fmt.Println("")
			fmt.Println(" --- Resources")
			pubnub.ParseGrantResources(cborObject.Resources)

			fmt.Println(" --- Patterns")
			pubnub.ParseGrantResources(cborObject.Patterns)
		}

		// err := json.Unmarshal(value, &resp)
		// if err != nil {
		// 	fmt.Printf("\nUnmarshal Error---> %#v", err)
		// 	//log.Println("Write file:", ioutil.WriteFile("data.json", value, 0600))
		// } else {
		// 	fmt.Println("\nUnmarshalled Token--->", resp)
		// }
	}

	assert.NotNil(res)

}
