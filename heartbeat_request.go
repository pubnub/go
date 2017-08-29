package pubnub

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pubnub/go/utils"
)

const HEARTBEAT_PATH = "/v2/presence/sub-key/%s/channel/%s/hearbeat"

type heartbeatBuilder struct {
	opts *heartbeatOpts
}

func newHeartbeatBuilder(pubnub *PubNub) *heartbeatBuilder {
	builder := heartbeatBuilder{
		opts: &heartbeatOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newHeartbeatBuilderWithContext(pubnub *PubNub,
	context Context) *heartbeatBuilder {
	builder := heartbeatBuilder{
		opts: &heartbeatOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *heartbeatBuilder) State(state interface{}) *heartbeatBuilder {
	b.opts.State = state

	return b
}

func (b *heartbeatBuilder) Channels(ch []string) *heartbeatBuilder {
	b.opts.Channels = ch

	return b
}

func (b *heartbeatBuilder) ChannelGroups(cg []string) *heartbeatBuilder {
	b.opts.ChannelGroups = cg

	return b
}

type heartbeatOpts struct {
	pubnub *PubNub

	State interface{}

	Channels      []string
	ChannelGroups []string

	ctx Context
}

func (o *heartbeatOpts) config() Config {
	return *o.pubnub.Config
}

func (o *heartbeatOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *heartbeatOpts) context() Context {
	return o.ctx
}

func (o *heartbeatOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return ErrMissingSubKey
	}

	if len(o.Channels) == 0 && len(o.ChannelGroups) == 0 {
		return ErrMissingChannel
	}

	return nil
}

func (o *heartbeatOpts) buildPath() (string, error) {
	return fmt.Sprintf(HEARTBEAT_PATH,
		o.pubnub.Config.SubscribeKey,
		utils.UrlEncode(strings.Join(o.Channels, ","))), nil
}

func (o *heartbeatOpts) buildQuery() (string, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	q.Set("heartbeat", strconv.Itoa(o.pubnub.Config.PresenceTimeout))

	if len(o.ChannelGroups) > 0 {
		q.Set("channel-group", strings.Join(o.ChannelGroups, ","))
	}

	if o.State != nil {
		state, _ := utils.ValueAsString(o.State)
		// TODO: handle error
		q.Set("state", string(state))
	}

	return "", nil
}

func (o *heartbeatOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *heartbeatOpts) httpMethod() string {
	return "GET"
}

func (o *heartbeatOpts) isAuthRequired() bool {
	return true
}

func (o *heartbeatOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *heartbeatOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *heartbeatOpts) operationType() PNOperationType {
	return PNHeartBeatOperation
}
