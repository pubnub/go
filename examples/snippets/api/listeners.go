package main

import pubnub "github.com/pubnub/go"

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

	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()

	<-done
}
