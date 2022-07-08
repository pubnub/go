package main

import (
	"fmt"

	pubnub "github.com/pubnub/go/v7"
)

func main() {
	config := pubnub.NewConfigWithUserId(UserId(pubnub.GenerateUUID()))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.HereNow().
		Channels([]string{"my_channel", "demo"}).
		IncludeUUIDs(true).
		Execute()

	if err != nil {
		fmt.Println("Error :", err)
	}

	fmt.Println(status)

	for _, v := range res.Channels {
		fmt.Println("---")
		fmt.Println("channel: ", v.ChannelName)
		fmt.Println("occupancy: ", v.Occupancy)

		for _, occupant := range v.Occupants {
			fmt.Printf("UUID: %s, state: %s\n", occupant.UUID, occupant.State)
		}
	}
}
