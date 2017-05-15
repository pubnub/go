package pubnub

import (
	"context"
	"fmt"
	"net/url"
)

const PUBLISH_GET_PATH = "/publish/%s/%s/0/%s/%s/%s"
const PUBLISH_POST_PATH = "/publish/%s/%s/0/%s/%s"

type Publish struct {
	TransactionalEndpoint

	pubnub *PubNub

	Channel string
	Message interface{}
	UsePost bool

	SuccessChannel chan interface{}
	ErrorChannel   chan error
}

func NewPublish(pubnub *PubNub) *Publish {
	return &Publish{
		pubnub: pubnub,
	}
}
func (e *Publish) PubNub() *PubNub {
	return e.pubnub
}

func (e *Publish) buildPath() string {
	if e.UsePost == true {
		return fmt.Sprintf(PUBLISH_POST_PATH,
			e.pubnub.Config.PublishKey,
			e.pubnub.Config.SubscribeKey,
			e.Channel,
			"0")
	}

	return fmt.Sprintf(PUBLISH_GET_PATH,
		e.pubnub.Config.PublishKey,
		e.pubnub.Config.SubscribeKey,
		e.Channel,
		"0",
		e.Message)
}

func (e *Publish) Execute() (interface{}, error) {
	panic("not implemented")
}

func (e *Publish) ExecuteWithContext(ctx context.Context) (interface{}, error) {
	// TODO: execute with context
	return executeRequest(ctx, e, e.SuccessChannel, e.ErrorChannel)
}

func (e *Publish) buildQuery() *url.Values {
	q := defaultQuery()

	q.Set("blah", "hey")

	return q
}

func (e *Publish) buildBody() string {
	return ""
}
