package e2e

import (
	"fmt"
	"log"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeConnectEvent(t *testing.T) {
	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)

	pn := pubnub.NewPubNub(configCopy())

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.ConnectedCategory:
					doneSubscribe <- true
				case pubnub.DisconnectedCategory:
					doneUnsubscribe <- true
					return
				}
			case message := <-listener.Message:
				fmt.Println(message)
			case presence := <-listener.Presence:
				fmt.Println(presence)
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"blah"},
	})

	<-doneSubscribe

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		Channels: []string{"blah"},
	})

	<-doneUnsubscribe
}

func TestSubscribePublishUnsubscribeChannels(t *testing.T) {
	assert := assert.New(t)
	doneUnsubscribe := make(chan bool)

	pn := pubnub.NewPubNub(configCopy())

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				if len(status.AffectedChannels) == 1 &&
					status.Operation == pubnub.PNUnsubscribeOperation {
					assert.Equal(status.AffectedChannels[0], "ch1")
				}

				switch status.Category {
				case pubnub.ConnectedCategory:
					pn.Unsubscribe(&pubnub.UnsubscribeOperation{
						Channels: []string{"ch1"},
					})
				case pubnub.CancelledCategory:
					pn.Publish(&pubnub.PublishOpts{
						Channel: "ch2",
						Message: "hey",
					})
				case pubnub.DisconnectedCategory:
					doneUnsubscribe <- true
					return
				}
			case message := <-listener.Message:
				if message.Message == "hey" {
					pn.Unsubscribe(&pubnub.UnsubscribeOperation{
						Channels: []string{"ch2"},
					})
				}
			case presence := <-listener.Presence:
				fmt.Println(presence)
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch1", "ch2"},
	})

	<-doneUnsubscribe
}

func TestSubscribeSingleMessage(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)
	donePublish := make(chan bool)

	pn := pubnub.NewPubNub(configCopy())

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.ConnectedCategory:
					doneSubscribe <- true
				case pubnub.DisconnectedCategory:
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				assert.Equal(message.Message, "hey")
				donePublish <- true
			case presence := <-listener.Presence:
				fmt.Println(presence)
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch"},
	})

	<-doneSubscribe

	pn.Publish(&pubnub.PublishOpts{
		Channel: "ch",
		Message: "hey",
	})

	<-donePublish

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		Channels: []string{"ch"},
	})

	<-doneUnsubscribe
}

func TestSubscribePresenceSingleChannel(t *testing.T) {
	pn := pubnub.NewPubNub(configCopy())

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.ConnectedCategory:
					fmt.Println("connected")
				case pubnub.DisconnectedCategory:
					fmt.Println("disconnected")
				}
			case message := <-listener.Message:
				log.Println(message)
			case presence := <-listener.Presence:
				fmt.Println(presence)
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels:        []string{"ch"},
		PresenceEnabled: true,
	})
}
