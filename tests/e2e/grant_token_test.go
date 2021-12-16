package e2e

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v6"
	"github.com/stretchr/testify/assert"
)

func TestGrantToken(t *testing.T) {

	assert := assert.New(t)

	pcc := pamConfigCopy()
	pcc.SubscribeKey = "sub-c-78c27be6-001a-11ec-b0c0-62dfa3a98328"
	pcc.PublishKey = "pub-c-668e1ce9-5f2b-4c51-980f-26e16ff3698e"
	pcc.SecretKey = "sec-c-NWUyYmRhZjMtNTZlNC00ZDNiLTkyMmEtN2NmMmU2MjY3Y2Rm"
	pcc.Origin = "aws-hnd-1-ingress-tls10.pubnub.com"
	pn := pubnub.NewPubNub(pcc)
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	subscribed := make(chan bool)

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
		cc := configCopy()
		cc.SubscribeKey = "sub-c-78c27be6-001a-11ec-b0c0-62dfa3a98328"
		cc.PublishKey = "pub-c-668e1ce9-5f2b-4c51-980f-26e16ff3698e"
		cc.Origin = "aws-hnd-1-ingress-tls10.pubnub.com"

		cc.AuthKey = token
		pnClient := pubnub.NewPubNub(cc)
		if enableDebuggingInTests {
			pnClient.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
		}

		listener := pubnub.NewListener()
		exitListener := make(chan bool)

		go func() {
		ExitLabel:
			for {
				select {
				case status := <-listener.Status:
					switch status.Category {
					case pubnub.PNConnectedCategory:
						subscribed <- true
					}
				case <-exitListener:
					break ExitLabel

				}
			}
		}()

		pnClient.AddListener(listener)

		pnClient.Subscribe().Channels([]string{ch1}).Execute()

		<-subscribed

		_, _, errPub := pnClient.Publish().Channel(ch1).Message("expectedMsg").Execute()
		if errPub != nil {
			assert.Fail("Failed when publishing: " + errPub.Error())
		}

		_, _, errRev := pn.RevokeToken().Token(token).Execute()
		if errRev != nil {
			assert.Fail("Failed when revoking: " + errRev.Error())
		}

		time.Sleep(60 * time.Second)

		_, _, errPub2 := pnClient.Publish().Channel(ch1).Message("expectedMsg").Execute()
		if errPub2 == nil {
			assert.Fail("Didn't fail when publishing")
		} else {
			fmt.Println(errPub2)
		}

	}

}
