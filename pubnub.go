package pubnub

type PubNub struct {
	PNConfig *PNConfiguration
}

func (pn *PubNub) Publish() *Publish {
	return NewPublish(pn)
}

func NewPubNub(pnconf *PNConfiguration) *PubNub {
	return &PubNub{
		PNConfig: pnconf,
	}
}

func NewPubNubDemo() *PubNub {
	return &PubNub{
		PNConfig: NewPNConfigurationDemo(),
	}
}
