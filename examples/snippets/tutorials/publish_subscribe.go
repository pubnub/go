package main

import (
	"fmt"

	pubnub "github.com/pubnub/go/v5"
)

var pn *pubnub.PubNub

func init() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	pn = pubnub.NewPubNub(config)
}

func main() {
	listener := pubnub.NewListener()
	doneSubscribe := make(chan bool)
	data := make(map[string]interface{})

	data["awesome"] = "data"

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					res, status, err := pn.Publish().
						Channel("awesome-channel").
						Message(data).
						Execute()

					fmt.Printf(res, status, err)

					doneSubscribe <- true
					return
				}
			case <-listener.Message:
			case <-listener.Presence:
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"awesome-channel"}).
		Execute()

	<-doneSubscribe

	pn.Unsubscribe().
		Channels([]string{"awesome-channel"}).
		Execute()
}
