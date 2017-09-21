package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pubnub/go/utils"
)

const REMOVE_CHANNEL_CHANNEL_GROUP = "/v1/channel-registration/sub-key/%s/channel-group/%s"

var emptyRemoveChannelChannelGroupResponse *RemoveChannelChannelGroupResponse

type removeChannelChannelGroupBuilder struct {
	opts *removeChannelOpts
}

func newRemoveChannelChannelGroupBuilder(
	pubnub *PubNub) *removeChannelChannelGroupBuilder {
	builder := removeChannelChannelGroupBuilder{
		opts: &removeChannelOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newRemoveChannelChannelGroupBuilderWithContext(
	pubnub *PubNub, context Context) *removeChannelChannelGroupBuilder {
	builder := removeChannelChannelGroupBuilder{
		opts: &removeChannelOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *removeChannelChannelGroupBuilder) Channels(
	ch []string) *removeChannelChannelGroupBuilder {
	b.opts.Channels = ch
	return b
}

func (b *removeChannelChannelGroupBuilder) Group(
	cg string) *removeChannelChannelGroupBuilder {
	b.opts.Group = cg
	return b
}

func (b *removeChannelChannelGroupBuilder) Execute() (
	*RemoveChannelChannelGroupResponse, StatusResponse, error) {
	rawJson, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveChannelChannelGroupResponse, status, err
	}

	return newRemoveChannelChannelGroupResponse(rawJson, status)
}

type removeChannelOpts struct {
	pubnub *PubNub

	Channels []string

	Group string

	Transport http.RoundTripper

	ctx Context
}

func (o *removeChannelOpts) config() Config {
	return *o.pubnub.Config
}

func (o *removeChannelOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *removeChannelOpts) context() Context {
	return o.ctx
}

func (o *removeChannelOpts) validate() error {
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

func (o *removeChannelOpts) buildPath() (string, error) {
	return fmt.Sprintf(REMOVE_CHANNEL_CHANNEL_GROUP,
		o.pubnub.Config.SubscribeKey,
		utils.UrlEncode(o.Group)), nil
}

func (o *removeChannelOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	q.Set("remove", strings.Join(o.Channels, ","))

	return q, nil
}

func (o *removeChannelOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *removeChannelOpts) httpMethod() string {
	return "GET"
}

func (o *removeChannelOpts) isAuthRequired() bool {
	return true
}

func (o *removeChannelOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *removeChannelOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *removeChannelOpts) operationType() PNOperationType {
	return PNRemoveChannelFromChannelGroupOperation
}

type RemoveChannelChannelGroupResponse struct {
}

func newRemoveChannelChannelGroupResponse(jsonBytes []byte,
	status StatusResponse) (*RemoveChannelChannelGroupResponse,
	StatusResponse, error) {
	return emptyRemoveChannelChannelGroupResponse, status, nil
}
