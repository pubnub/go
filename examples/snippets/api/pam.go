package main

import (
	"fmt"

	pubnub "github.com/pubnub/go/v5"
)

func main() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "sub-c-b9ab9508-43cf-11e8-9967-869954283fb4"
	config.PublishKey = "pub-c-1bd448ed-05ba-4dbc-81a5-7d6ff5c6e2bb"
	config.SecretKey = "sec-c-NjlmYzVkMjEtOWIxZi00YmJlLThjZDktMjI4NGQwZDUxZDQ0"

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.Grant().
		Channels([]string{"ch1"}).
		ChannelGroups([]string{"cg1"}).
		AuthKeys([]string{"key-1"}).
		Read(true).
		Write(true).
		Manage(true).
		Execute()

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(res, status)
}
