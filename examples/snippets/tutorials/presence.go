package main

import (
	"fmt"

	pubnub "github.com/pubnub/go"
)

func subscribe() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.Uuid = "Stephen"

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()
	doneConnect := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneConnect <- true
					return
				case pubnub.PNReconnectedCategory:
				}
			case <-listener.Message:
			case <-listener.Presence:
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"my-channel"},
	})

	<-doneConnect
}

func hereNow() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.Uuid = "Stephen"

	pn := pubnub.NewPubNub(config)

	res, _, _ := pn.HereNow().
		Channels([]string{"ch1", "ch2"}).
		IncludeUuids(true).
		Execute()

	for _, v := range res.Channels {
		fmt.Println("Channel: ", v.ChannelName)
		fmt.Println("Occupancy: ", v.Occupancy)
		fmt.Println("Occupants")

		for _, v := range v.Occupants {
			fmt.Println("Uuid: ", v.Uuid, ", state: ", v.State)
		}
	}
}

func globalHereNow() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.Uuid = "Stephen"

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.HereNow().
		IncludeState(true).
		IncludeUuids(true).
		Execute()

	fmt.Println(res, status, err)
}

func whereNow() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.Uuid = "Stephen"

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.WhereNow().
		Uuid("person-uuid").
		Execute()

	fmt.Println(res, status, err)
}

func getState() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.Uuid = "Stephen"

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.GetState().
		Channels([]string{"my-unique-ch"}).
		Execute()

	fmt.Println(res, status, err)
}

func main() {
	subscribe()
	hereNow()
	globalHereNow()
	getState()
}
