package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pubnub/go/utils"
)

const ADD_CHANNEL_CHANNEL_GROUP_PATH = "/v1/channel-registration/sub-key/%s/channel-group/%s"

var emptyAddChannelChannelGroupResp *AddChannelChannelGroupResponse

type addChannelChannelGroupBuilder struct {
	opts *addChannelOpts
}

func newAddChannelChannelGroupBuilder(
	pubnub *PubNub) *addChannelChannelGroupBuilder {
	builder := addChannelChannelGroupBuilder{
		opts: &addChannelOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newAddChannelChannelGroupBuilderWithContext(
	pubnub *PubNub, context Context) *addChannelChannelGroupBuilder {
	builder := addChannelChannelGroupBuilder{
		opts: &addChannelOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *addChannelChannelGroupBuilder) Channels(
	ch []string) *addChannelChannelGroupBuilder {
	b.opts.Channels = ch

	return b
}

func (b *addChannelChannelGroupBuilder) Group(
	cg string) *addChannelChannelGroupBuilder {
	b.opts.Group = cg

	return b
}

func (b *addChannelChannelGroupBuilder) Transport(
	tr http.RoundTripper) *addChannelChannelGroupBuilder {
	b.opts.Transport = tr

	return b
}

func (b *addChannelChannelGroupBuilder) Execute() (
	*AddChannelChannelGroupResponse, StatusResponse, error) {
	rawJson, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyAddChannelChannelGroupResp, status, err
	}

	return newAddChannelChannelGroupsResponse(rawJson, status)
}

type addChannelOpts struct {
	pubnub *PubNub

	Channels []string

	Group string

	Transport http.RoundTripper

	ctx Context
}

func (o *addChannelOpts) config() Config {
	return *o.pubnub.Config
}

func (o *addChannelOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *addChannelOpts) context() Context {
	return o.ctx
}

func (o *addChannelOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return ErrMissingSubKey
	}

	if len(o.Channels) == 0 {
		return ErrMissingChannel
	}

	if o.Group == "" {
		return ErrMissingChannelGroup
	}

	return nil
}

func (o *addChannelOpts) buildPath() (string, error) {
	return fmt.Sprintf(ADD_CHANNEL_CHANNEL_GROUP_PATH,
		o.pubnub.Config.SubscribeKey,
		utils.UrlEncode(o.Group)), nil
}

func (o *addChannelOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	q.Set("add", strings.Join(o.Channels, ","))

	return q, nil
}

func (o *addChannelOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *addChannelOpts) httpMethod() string {
	return "GET"
}

func (o *addChannelOpts) isAuthRequired() bool {
	return true
}

func (o *addChannelOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *addChannelOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *addChannelOpts) operationType() OperationType {
	return PNAddChannelsToChannelGroupOperation
}

type AddChannelChannelGroupResponse struct {
}

func newAddChannelChannelGroupsResponse(jsonBytes []byte, status StatusResponse) (
	*AddChannelChannelGroupResponse, StatusResponse, error) {

	return emptyAddChannelChannelGroupResp, status, nil
}
