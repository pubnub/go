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
			case <-listener.Message:
			case <-listener.Presence:
			}
		}
	}()

	pn.AddListener(listener)

	res, status, err := pn.SetState().
		Channels([]string{"ch"}).
		State(map[string]interface{}{
			"field_a": "cool",
			"field_b": 21,
		}).
		Execute()

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(res, status)

	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()

	<-done
}
