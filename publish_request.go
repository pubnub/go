package pubnub

import (
	"fmt"
	"strconv"

	"net/http"
	"net/url"
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
	opts.ctx = ctx

	resp, err := executeRequest(opts)
	if err != nil {
		return PublishResponse{}, err
	}

	return PublishResponse{
		Timestamp: 123,
	}, nil
}

type PublishOpts struct {
	pubnub *PubNub

	Ttl     int
	Channel string
	Message interface{}
	Meta    interface{}

	UsePost     bool
	ShouldStore bool
	Serialize   bool
	Replicate   bool

	ctx Context
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

func (o *PublishOpts) context() Context {
	return o.ctx
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

func (o *PublishOpts) customParams() map[string]string {
	params := make(map[string]string)

	if o.Meta != nil {
		params["meta"] = o.Meta.(string)
	}

	if o.ShouldStore {
		params["store"] = "1"
	} else {
		params["store"] = "0"
	}

	if o.Ttl != 0 {
		params["ttl"] = strconv.Itoa(o.Ttl)
	}

	params["seqn"] = strconv.Itoa(<-o.pubnub.publishSequence)

	if o.Replicate {
		params["norep"] = "true"
	}

	return params
}

func (o *PublishOpts) buildData() string {
	return ""
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
	params := o.customParams()

	for k, v := range params {
		q.Set(k, v)
	}

	return q
}

func (o *PublishOpts) buildBody() string {
	return ""
}

func (o *PublishOpts) httpMethod() string {
	if o.UsePost {
		return "POST"
	} else {
		return "GET"
	}
}

func (o *PublishOpts) isAuthRequired() bool {
	return true
}

func (o *PublishOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *PublishOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectionTimeout
}
