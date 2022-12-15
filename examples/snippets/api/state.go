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

	res, status, err := pn.SetState().
		Channels([]string{"ch1"}).
		State(map[string]interface{}{
			"age": 20,
		}).
		Execute()

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(res, status)
}
