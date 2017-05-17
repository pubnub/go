package pubnub

type PubNub struct {
	Config          *Config
	publishSequence chan int
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

func runPublishSequenceManager(maxSequence int, ch chan int) {
	for i := 1; ; i++ {
		if i == maxSequence {
			i = 1
		}

		ch <- i
	}
}
