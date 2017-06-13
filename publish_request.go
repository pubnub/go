package pubnub

import (

	// "errors"

	"fmt"

	// "net/http"

	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
)

const PUBLISH_GET_PATH = "/publish/%s/%s/0/%s/%s/%s"
const PUBLISH_POST_PATH = "/publish/%s/%s/0/%s/%s"

func PublishRequest(pn *PubNub, opts *PublishOpts) (PublishResponse, error) {
	opts.pubnub = pn
	resp, err := executeRequest(opts)
	if err != nil {
		return PublishResponse{}, err
	}

	var vals []interface{}

	// TODO: cast to slice
	// fmt.Println("%#v", resp)
	switch v := resp.(type) {
	case []interface{}:
		fmt.Println("slice", v)
		vals = v
	default:
		fmt.Println("default")
	}

	timestamp := vals[1].(int)

	return PublishResponse{
		Timestamp: timestamp,
	}, nil
}

func PublishRequestWithContext(ctx Context,
	pn *PubNub, opts *PublishOpts) (PublishResponse, error) {
	opts.pubnub = pn

	resp, err := executeRequestWithContext(ctx, opts)
	if err != nil {
		return PublishResponse{}, err
	}

	fmt.Println(resp)
	return PublishResponse{
		Timestamp: 123,
	}, nil
}

type PublishOpts struct {
	pubnub *PubNub

	Channel string
	Message interface{}
	UsePost bool
}

type PublishResponse struct {
	Timestamp int
}

func (o *PublishOpts) config() Config {
	return *o.pubnub.Config
}

func (o *PublishOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *PublishOpts) validate() error {
	if o.config().PublishKey == "" {
		return ErrMissingPubKey
	}

	if o.config().SubscribeKey == "" {
		return ErrMissingSubKey
	}

	if o.Channel == "" {
		return ErrMissingChannel
	}

	if o.Message == nil {
		return ErrMissingMessage
	}

	return nil
}

func (o *PublishOpts) buildPath() string {
	if o.UsePost == true {
		return fmt.Sprintf(PUBLISH_POST_PATH,
			o.pubnub.Config.PublishKey,
			o.pubnub.Config.SubscribeKey,
			o.Channel,
			"0")
	}

	return fmt.Sprintf(PUBLISH_GET_PATH,
		o.pubnub.Config.PublishKey,
		o.pubnub.Config.SubscribeKey,
		o.Channel,
		"0",
		o.Message)
}

func (o *PublishOpts) buildQuery() *url.Values {
	q := defaultQuery()

	q.Set("blah", "hey")

	return q
}

func (o *PublishOpts) buildBody() string {
	return ""
}

func (o *PublishOpts) validateParams() error {
	if o.pubnub.Config.PublishKey == "" {
		return pnerr.NewValidationError(
			fmt.Sprintf("Publish key was expected but got: %s", o.pubnub.Config.PublishKey))
	}

	return nil
}
