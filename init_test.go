package pubnub

var pnconfig *Config
var pubnub *PubNub

func init() {
	pnconfig = NewConfig()

	pnconfig.PublishKey = "pub_key"
	pnconfig.SubscribeKey = "sub_key"
	pnconfig.SecretKey = "secret_key"

	pubnub = NewPubNub(pnconfig)
}

func pubnubCopy() *PubNub {
	pn := new(PubNub)
	*pn = *pubnub
	return pn
}
