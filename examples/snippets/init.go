package main

import (
	pubnub "github.com/pubnub/go/v7"
)

func initPubNub() *pubnub.PubNub {
	pnconfig := pubnub.NewConfigWithUserId(pubnub.UserId(pubnub.GenerateUUID()))

	pnconfig.PublishKey = "demo"
	pnconfig.SubscribeKey = "demo"

	return pubnub.NewPubNub(pnconfig)
}
