package main

import (
	"fmt"

	pubnub "github.com/pubnub/go/v5"
)

func main() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()
	done := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					done <- true
				}
			case message := <-listener.Message:
				//Channel
				fmt.Println(message.Channel)
				//Subscription
				fmt.Println(message.Subscription)
				//Payload
				fmt.Println(message.Message)
				//Publisher ID
				fmt.Println(message.Publisher)
				//Timetoken
				fmt.Println(message.Timetoken)
			case <-listener.Presence:
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()

	<-done
}
