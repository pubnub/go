package messaging

var pubnub = &Pubnub{
	SubscribeKey: "demo",
}

var pubnubWithAuth = &Pubnub{
	SubscribeKey:      "demo",
	AuthenticationKey: "blah",
}
