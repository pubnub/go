package pubnub

var pubnub *PubNub

func init() {
	pnconfig := NewConfig()
	pnconfig.PublishKey = "my_pub_key"
	pnconfig.SubscribeKey = "my_sub_key"
	pnconfig.SecretKey = "my_secret_key"

	pubnub = NewPubNub(pnconfig)
}
