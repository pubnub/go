package pubnub

import (
	"fmt"
)

const PUBLISH_GET_PATH = "/publish/%s/%s/0/%s/%s/%s"
const PUBLISH_POST_PATH = "/publish/%s/%s/0/%s/%s"

type Publish struct {
	TransactionalEndpoint

	pubnub *PubNub

	channel string
	message interface{}

	usePost bool
}

func NewPublish(pubnub *PubNub) *Publish {
	return &Publish{
		pubnub: pubnub,
	}
}

func (e *Publish) Channel(ch string) *Publish {
	e.channel = ch
	return e
}

func (e *Publish) Message(msg interface{}) *Publish {
	// TODO: serialize
	e.message = msg
	return e
}

func (e *Publish) BuildPath() string {
	if e.usePost == true {
		return fmt.Sprintf(PUBLISH_GET_PATH, e.pubnub.PNConfig.SubscribeKey)
	}

	return fmt.Sprintf(PUBLISH_POST_PATH)
}

func (e *Publish) BuildQuery() map[string]string {
	return make(map[string]string)
}

func (e *Publish) BuildBody() string {
	return ""
}
