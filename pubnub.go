package pubnub

type PubNub struct {
	Config *Config
}

func (pn *PubNub) Publish() *Publish {
	return NewPublish(pn)
}

func NewPubNub(pnconf *Config) *PubNub {
	return &PubNub{
		Config: pnconf,
	}
}

func NewPubNubDemo() *PubNub {
	return &PubNub{
		Config: NewDemoConfig(),
	}
}
