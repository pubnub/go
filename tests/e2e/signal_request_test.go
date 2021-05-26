package e2e

import (
	"fmt"
	//"log"
	"net"
	"os"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v5"
	"github.com/stretchr/testify/assert"
)

func TestSignal(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(config)
	ips, err1 := net.LookupIP(pn.Config.Origin)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err1)
		os.Exit(1)
	}
	for _, ip := range ips {
		fmt.Printf("%s IN A %s\n", pn.Config.Origin, ip.String())
	}
	pn.Config.SubscribeKey = "demo"
	pn.Config.PublishKey = "demo"

	_, _, err := pn.Signal().
		Channel("ch").
		Message("hey").
		Execute()

	assert.Nil(err)

}

func TestSubscribeSignalUnsubscribeInt(t *testing.T) {
	assert := assert.New(t)
	sigMessage := make(chan interface{})
	s := 1

	go SubscribeSignalUnsubscribeMultiCommon(t, s, "", sigMessage, false, false)
	m := <-sigMessage
	msg := m.(float64)
	assert.Equal(float64(s), msg)
}

func SubscribeSignalUnsubscribeMultiCommon(t *testing.T, s interface{}, cipher string, sigMessage chan interface{}, disablePNOtherProcessing bool, usePost bool) {
	assert := assert.New(t)

	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)
	donePublish := make(chan bool)
	exit := make(chan bool)
	errChan := make(chan string)
	r := GenRandom()

	ch := fmt.Sprintf("testChannel_sub_%d", r.Intn(99999))

	pn := pubnub.NewPubNub(configCopy())

	pn.Config.SubscribeKey = "demo"
	pn.Config.PublishKey = "demo"

	pn.Config.CipherKey = cipher
	pn.Config.DisablePNOtherProcessing = disablePNOtherProcessing

	listener := pubnub.NewListener()

	tic := time.NewTicker(time.Duration(timeout) * time.Second)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
				case pubnub.PNAcknowledgmentCategory:
					doneUnsubscribe <- true
				default:
					doneUnsubscribe <- true
				}
			case message := <-listener.Signal:
				donePublish <- true
				if sigMessage != nil {
					sigMessage <- message.Message
				}

			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			case <-tic.C:
				assert.Fail("timeout")
				errChan <- "timeout"

				return
			case <-exit:
				tic.Stop()
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().Channels([]string{ch}).Execute()

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Signal().Channel(ch).Message(s).Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}
	exit <- true

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
}
