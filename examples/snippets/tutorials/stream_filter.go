package main

import pubnub "github.com/pubnub/go"

func main() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.Uuid = "my_uuid"

	pn := pubnub.NewPubNub(config)

	meta := map[string]interface{}{
		"my":   "meta",
		"name": "PubNub",
	}

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch1"},
	})

	pn.Publish().
		Meta(meta).
		Message("hello").
		Channel("ch1").
		Execute()
}
