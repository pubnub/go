package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pubnub/go/utils"
)

const heartbeatPath = "/v2/presence/sub-key/%s/channel/%s/heartbeat"

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

func (b *heartbeatBuilder) Execute() (interface{}, StatusResponse, error) {
	rawJson, status, err := executeRequest(b.opts)
	if err != nil {
		return "", status, err
	}

	var value interface{}

	err = json.Unmarshal(rawJson, &value)
	if err != nil {
		return nil, status, err
	}

	return value, status, nil
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
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 && len(o.ChannelGroups) == 0 {
		return newValidationError(o, "Missing Channel or Channel Group")
	}

	return nil
}

func (o *heartbeatOpts) buildPath() (string, error) {
	channels := string(utils.JoinChannels(o.Channels))

	return fmt.Sprintf(heartbeatPath,
		o.pubnub.Config.SubscribeKey,
		channels), nil
}

func (o *heartbeatOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid, o.pubnub.telemetryManager)

	q.Set("heartbeat", strconv.Itoa(o.pubnub.Config.PresenceTimeout))

	if len(o.ChannelGroups) > 0 {
		q.Set("channel-group", strings.Join(o.ChannelGroups, ","))
	}

	if o.State != nil {
		state, err := utils.ValueAsString(o.State)
		if err != nil {
			return &url.Values{}, err
		}

		if string(state) != "{}" {
			q.Set("state", string(state))
		}
	}

	return q, nil
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

func (o *heartbeatOpts) operationType() OperationType {
	return PNHeartBeatOperation
}

func (o *heartbeatOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}
