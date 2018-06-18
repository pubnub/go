package main

import (
	"fmt"

	pubnub "github.com/pubnub/go"
)

func main() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SecretKey = "demo"

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.History().
		Channel("my_channel").
		Count(100).
		Execute()

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(res, status)

	res, status, err = pn.History().
		Channel("my_channel").
		Count(100).
		Start(int64(-1)).
		End(int64(15093483374296431)).
		Reverse(true).
		Execute()

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(res, status)
}
