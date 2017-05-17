package pubnub

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const PUBLISH_GET_PATH = "/publish/%s/%s/0/%s/%s/%s"
const PUBLISH_POST_PATH = "/publish/%s/%s/0/%s/%s"

type Publish struct {
	pubnub *PubNub

	Channel string
	Message interface{}
	UsePost bool

	Transport *http.Client

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
	ctx, _ := context.WithCancel(context.Background())
	okCh := make(chan interface{})
	errCh := make(chan error)

	go executeRequest(ctx, e, okCh, errCh)

	select {
	case resp := <-okCh:
		if e.SuccessChannel != nil {
			e.SuccessChannel <- resp
		}
		return resp, nil
	case err := <-errCh:
		if e.ErrorChannel != nil {
			e.ErrorChannel <- err
		}
		return nil, err
	case <-ctx.Done():
		return nil, errors.New("pubnub: Context cancelled")
	}
}

func (e *Publish) ExecuteWithContext(ctx context.Context) (interface{}, error) {
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
