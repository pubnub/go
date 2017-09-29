package pubnub

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pubnub/go/utils"
)

const LEAVE_PATH = "/v2/presence/sub-key/%s/channel/%s/leave"

type leaveBuilder struct {
	opts *leaveOpts
}

func newLeaveBuilder(pubnub *PubNub) *leaveBuilder {
	builder := leaveBuilder{
		opts: &leaveOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newLeaveBuilderWithContext(pubnub *PubNub, context Context) *leaveBuilder {
	builder := leaveBuilder{
		opts: &leaveOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *leaveBuilder) Channels(channels []string) *leaveBuilder {
	b.opts.Channels = channels
	return b
}

func (b *leaveBuilder) ChannelGroups(groups []string) *leaveBuilder {
	b.opts.ChannelGroups = groups
	return b
}

func (b *leaveBuilder) Execute() (StatusResponse, error) {
	_, status, err := executeRequest(b.opts)
	if err != nil {
		return status, err
	}

	return status, nil
}

type leaveOpts struct {
	Channels      []string
	ChannelGroups []string

	pubnub *PubNub
	ctx    Context
}

func (o *leaveOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *leaveOpts) httpMethod() string {
	return "GET"
}

func (o *leaveOpts) buildPath() (string, error) {
	channels := utils.JoinChannels(o.Channels)

	if string(channels) == "" {
		channels = []byte(",")
	}

	return fmt.Sprintf(LEAVE_PATH,
		o.pubnub.Config.SubscribeKey,
		channels), nil
}

func (o *leaveOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	if len(o.ChannelGroups) > 0 {
		channelGroup := utils.JoinChannels(o.ChannelGroups)
		q.Set("channel-group", string(channelGroup))
	}

	return q, nil
}

func (o *leaveOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *leaveOpts) config() Config {
	return *o.pubnub.Config
}

func (o *leaveOpts) context() Context {
	return o.ctx
}

func (o *leaveOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 && len(o.ChannelGroups) == 0 {
		return newValidationError(o, "Missing Channel or Channel Group")
	}

	return nil
}

func (o *leaveOpts) operationType() OperationType {
	return PNUnsubscribeOperation
}
