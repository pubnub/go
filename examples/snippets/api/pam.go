package main

import (
	"fmt"

	pubnub "github.com/pubnub/go"
)

func main() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "sub-c-90c51098-c040-11e5-a316-0619f8945a4f"
	config.PublishKey = "pub-c-1bd448ed-05ba-4dbc-81a5-7d6ff5c6e2bb"
	config.SecretKey = "sec-c-ZDA1ZTdlNzAtYzU4Zi00MmEwLTljZmItM2ZhMDExZTE2ZmQ5"

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
