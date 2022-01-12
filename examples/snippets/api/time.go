package main

import (
	"fmt"

	pubnub "github.com/pubnub/go/v7"
)

func main() {
	config := pubnub.NewConfig(pubnub.GenerateUUID())
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.Time().Execute()

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(res, status)
}
