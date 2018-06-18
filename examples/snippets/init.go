package main

var pn *PubNub

func Init() {
	pnconfig = NewConfig()

	pnconfig.PublishKey = "demo"
	pnconfig.SubscribeKey = "demo"

	pn = NewPubNub(pnconfig)
}

func pubnubCopy() *PubNub {
	pn := new(PubNub)
	*pn = *pubnub
	return pn
}
