package main

import (
	"fmt"
	"strconv"

	pubnub "github.com/pubnub/go"
)

func getAllMessages(startTT int64) {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SecretKey = "my-secret"
	config.AuthKey = "my-auth"

	pn := pubnub.NewPubNub(config)

	res, _, _ := pn.History().
		Channel("history_channel").
		Count(2).
		Execute()

	msgs := res.Messages
	start := res.StartTimetoken
	end := res.EndTimetoken

	if len(msgs) > 0 {
		fmt.Println(len(msgs))
		fmt.Println("start " + strconv.Itoa(int(start)))
		fmt.Println("end " + strconv.Itoa(int(end)))
	}

	if len(msgs) == 100 {
		getAllMessages(start)
	}
}

func main() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SecretKey = "my-secret"
	config.AuthKey = "my-auth"

	pn := pubnub.NewPubNub(config)

	for i := 0; i < 500; i++ {
		pn.Publish().
			Message("message #" + strconv.Itoa(i)).
			Channel("history_channel").
			ShouldStore(true).
			Execute()
	}

	res, status, err := pn.History().
		Channel("history_channel").
		Count(2).
		Execute()

	fmt.Println(res, status, err)

	res, status, err = pn.History().
		Channel("history_channel").
		Count(100).
		Start(13847168620721752).
		End(15090358935871532).
		Execute()

	fmt.Println(res, status, err)

	res, status, err = pn.History().
		Channel("history_channel").
		Count(100).
		IncludeTimetoken(true).
		Execute()

	fmt.Println(res, status, err)

	getAllMessages(int64(15090358935871532))
}
