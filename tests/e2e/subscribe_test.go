package e2e

import (
	"fmt"
	"testing"

	pubnub "github.com/pubnub/go"
)

func TestSubscribeConnectEvent(t *testing.T) {
	// assert := assert.New(t)
	done := make(chan bool)

	pn := pubnub.NewPubNub(configCopy())

	listener := pubnub.NewListener()

	go func() {
		select {
		case status := <-listener.Status:
			fmt.Println(status)
			done <- true
		case message := <-listener.Message:
			fmt.Println(message)
		case presence := <-listener.Presence:
			fmt.Println(presence)
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"blah"},
	})

	<-done
}
