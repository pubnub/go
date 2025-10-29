package main

import pubnub "github.com/pubnub/go/v8"

func main() {
	// Uncomment to run any of the tutorial examples:
	// mainAccessManager()
	// mainPresence()
	// mainPublishSubscribe()
	// mainStorageAndPlayback()
	// mainStreamController()
	mainStreamFilter()
}

func mainStreamFilter() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId(pubnub.GenerateUUID()))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SetUserId(pubnub.UserId("my_uuid"))

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
