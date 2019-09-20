package main

import (
	"fmt"
	pubnub "github.com/pubnub/go"
)

func main() {

	// Config
	config := pubnub.NewConfig()
	config.PublishKey = "demo"
	config.SubscribeKey = "demo"
	// End Config

	// Init
	pn := pubnub.NewPubNub(config)
	// End Init

	// Init Listener
	listener := pubnub.NewListener()

	waitForConnect := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					// Connect event. You can do stuff like publish, and know you'll get it.
					// Or just use the connected event to confirm you are subscribed for
					// UI / internal notifications, etc
					waitForConnect <- true
				}
			case msg := <-listener.Message:
				fmt.Println(" --- MESSAGE: ")
				fmt.Println(fmt.Sprintf("msg.Channel: %s", msg.Channel))
				fmt.Println(fmt.Sprintf("msg.Message: %s", msg.Message))
				fmt.Println(fmt.Sprintf("msg.SubscribedChannel: %s", msg.SubscribedChannel))
				fmt.Println(fmt.Sprintf("msg.Timetoken: %d", msg.Timetoken))
			case presence := <-listener.Presence:
				fmt.Println(" --- PRESENCE: ")
				fmt.Println(fmt.Sprintf("%s", presence))
			}
		}
	}()
	// End Init Listener

	// Add Listener
	pn.AddListener(listener)

	// Subscribe
	pn.Subscribe().
		Channels([]string{"pubnub_onboarding_channel"}).
		WithPresence(true).
		Execute()

	<-waitForConnect

	// Publish

	msg := map[string]interface{}{
		"sender":  pn.Config.UUID,
		"content": "Hello From Go SDK",
	}

	resPub, statusPub, errPub := pn.Publish().
		Channel("pubnub_onboarding_channel").
		Message(msg).
		Execute()
	fmt.Println(resPub, statusPub, errPub)

	// End Publish

	// History

	res, status, err := pn.History().
		Channel("pubnub_onboarding_channel").
		Count(10).
		IncludeTimetoken(true).
		Execute()

	if res != nil {
		fmt.Println(" --- HISTORY: ")
		if res.Messages != nil {
			for _, v := range res.Messages {
				fmt.Println(fmt.Sprintf("Message: %s, Timetoken: %d", v.Message, v.Timetoken))
			}
		} else {
			fmt.Println(fmt.Sprintf("Messages null"))
		}
		fmt.Println(fmt.Sprintf("EndTimetoken: %d", res.EndTimetoken))
		fmt.Println(fmt.Sprintf("StartTimetoken: %d", res.StartTimetoken))
		fmt.Println("")
	} else {
		fmt.Println(fmt.Sprintf("StatusResponse: %s %e", status.Error, err))
	}
	// End History
}
