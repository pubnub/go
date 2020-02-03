// +build go1.9

package e2e

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeParseLogsForAuthKey(t *testing.T) {

	assert := assert.New(t)
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.AuthKey = "myAuthKey"
	channel := "ch"
	if enableDebuggingInTests {

		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	pn.Subscribe().
		Channels([]string{channel}).
		Execute()

	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-tic.C:
		tic.Stop()
	}

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	s := fmt.Sprintf("%s", out)
	expected := fmt.Sprintf("https://%s/v2/subscribe/%s/%s/0?pnsdk=PubNub-Go/%s&uuid=%s&auth=%s",
		pn.Config.Origin,
		pn.Config.SubscribeKey,
		channel,
		pubnub.Version,
		pn.Config.UUID,
		pn.Config.AuthKey)

	//https://ps.pndsn.com/v2/subscribe/sub-c-e41d50d4-43ce-11e8-a433-9e6b275e7b64/ch/0?pnsdk=PubNub-Go/4.1.3&uuid=pn-ac860b0d-d078-46b1-b142-c5492101dc82&auth=myAuthKey
	assert.Contains(s, expected)
}
