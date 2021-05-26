// +build go1.8

package e2e

import (
	"fmt"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v5"
	"github.com/stretchr/testify/assert"
)

/////////////////////////////
// Channel Group Subscription
/////////////////////////////
func TestSubscribeUnsubscribeGroup(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-sug-ch")
	cg := randomized("sub-sug-cg")

	pn := pubnub.NewPubNub(configCopy())

	listener := pubnub.NewListener()
	exitListener := make(chan bool)

	go func() {
	ExitLabel:
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
					break
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
					break
				case pubnub.PNAcknowledgmentCategory:
					doneUnsubscribe <- true
					break
				case pubnub.PNCancelledCategory:
					continue
				default:
					errChan <- fmt.Sprintf("%v", status)
					//break
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				//break
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				//break
			case <-exitListener:
				break ExitLabel

			}
		}
	}()

	pn.AddListener(listener)

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{ch}).
		ChannelGroup(cg).
		Execute()

	assert.Nil(err)

	// await for adding channels
	time.Sleep(3 * time.Second)

	pn.Subscribe().
		ChannelGroups([]string{cg}).
		Execute()

	tic1 := time.NewTicker(time.Duration(timeout) * time.Second * 3)
	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	case <-tic1.C:
		tic1.Stop()
		assert.Fail("timeout")

	}

	time.Sleep(2 * time.Second)

	pn.Unsubscribe().
		ChannelGroups([]string{cg}).
		Execute()

	tic := time.NewTicker(time.Duration(timeout) * time.Second * 3)
	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")

	}

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))

	_, _, err = pn.RemoveChannelFromChannelGroup().
		Channels([]string{ch}).
		ChannelGroup(cg).
		Execute()
	exitListener <- true
}

func TestSubscribePublishUnsubscribeAllGroup(t *testing.T) {
	assert := assert.New(t)
	pn := pubnub.NewPubNub(configCopy())
	listener := pubnub.NewListener()
	doneSubscribe := make(chan bool)
	donePublish := make(chan bool)
	doneUnsubscribe := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-spuag-ch")
	cg1 := randomized("sub-spuag-cg1")
	cg2 := randomized("sub-spuag-cg2")

	pn.AddListener(listener)
	exitListener := make(chan bool)

	go func() {
	ExitLabel:
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				donePublish <- true
				assert.Equal("hey", message.Message)
				assert.Equal(ch, message.Channel)
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			case <-exitListener:
				break ExitLabel
			}
		}
	}()

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{ch}).
		ChannelGroup(cg1).
		Execute()

	assert.Nil(err)

	// await for adding channel
	time.Sleep(2 * time.Second)

	pn.Subscribe().
		ChannelGroups([]string{cg1, cg2}).
		Execute()

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Publish().Channel(ch).Message("hey").Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Unsubscribe().
		ChannelGroups([]string{cg2}).
		Execute()

	assert.Equal(len(pn.GetSubscribedGroups()), 1)

	pn.UnsubscribeAll()
	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")

	}

	assert.Equal(len(pn.GetSubscribedGroups()), 0)

	_, _, err = pn.RemoveChannelFromChannelGroup().
		Channels([]string{ch}).
		ChannelGroup(cg1).
		Execute()

	assert.Nil(err)
	exitListener <- true
}
