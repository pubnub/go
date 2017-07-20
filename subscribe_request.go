package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pubnub/go/utils"
)

const SUBSCRIBE_PATH = "/v2/subscribe/%s/%s/0"

func newSubscribeRequest(ctx Context) *SubscribeResponse {
	return &SubscribeResponse{}
}

type SubscribeResponse struct {
}

type SubscribeOpts struct {
	pubnub *PubNub

	Channels []string
	Groups   []string

	Region           string
	FilterExpression string

	Timetoken    int64
	WithPresence bool

	ctx Context
}

func (o *SubscribeOpts) config() Config {
	return *o.pubnub.Config
}

func (o *SubscribeOpts) client() *http.Client {
	return o.pubnub.GetSubscribeClient()
}

func (o *SubscribeOpts) context() Context {
	return o.ctx
}

func (o *SubscribeOpts) validate() error {
	if o.config().PublishKey == "" {
		return ErrMissingPubKey
	}

	if o.config().SubscribeKey == "" {
		return ErrMissingSubKey
	}

	if len(o.Channels) == 0 && len(o.Groups) == 0 {
		return ErrMissingChannel
	}

	return nil
}

func (o *SubscribeOpts) buildPath() (string, error) {
	channels, err := utils.ChannelsAsString(o.Channels)

	if err != nil {
		return "", err
	}

	if string(channels) == "" {
		channels = []byte(",")
	}

	return fmt.Sprintf(SUBSCRIBE_PATH,
		o.pubnub.Config.SubscribeKey,
		channels,
	), nil
}

func (o *SubscribeOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	if len(o.Groups) > 0 {
		channelGroup, _ := utils.ChannelsAsString(o.Groups)
		q.Set("channel-group", string(channelGroup))
	}

	if o.FilterExpression != "" {
		q.Set("filter-expr", o.FilterExpression)
	}

	if o.Timetoken != 0 {
		q.Set("tt", strconv.FormatInt(o.Timetoken, 10))
	}

	if o.Region != "" {
		q.Set("tr", o.Region)
	}

	return q, nil
}

func (o *SubscribeOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *SubscribeOpts) httpMethod() string {
	return "GET"
}

func (o *SubscribeOpts) isAuthRequired() bool {
	return true
}

func (o *SubscribeOpts) requestTimeout() int {
	return o.pubnub.Config.SubscribeRequestTimeout
}

func (o *SubscribeOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}
