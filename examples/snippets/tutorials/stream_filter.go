package main

import pubnub "github.com/pubnub/go/v7"

func main() {
	config := pubnub.NewConfigWithUserId(UserId(pubnub.GenerateUUID()))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SetUserId(UserId("my_uuid"))

	pn := pubnub.NewPubNub(config)

	meta := map[string]interface{}{
		"my":   "meta",
		"name": "PubNub",
	}

	pn.Subscribe().
		Channels([]string{"ch1"}).
		Execute()

	pn.Publish().
		Meta(meta).
		Message("hello").
		Channel("ch1").
		Execute()
}
