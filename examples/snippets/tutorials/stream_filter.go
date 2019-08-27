package main

import pubnub "github.com/sprucehealth/pubnub-go"

func main() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.UUID = "my_uuid"

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
