package pubnub

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/pubnub/go/v8/utils"
)

const heartbeatPath = "/v2/presence/sub-key/%s/channel/%s/heartbeat"

type heartbeatBuilder struct {
	opts *heartbeatOpts
}

func newHeartbeatBuilder(pubnub *PubNub) *heartbeatBuilder {
	return newHeartbeatBuilderWithContext(pubnub, pubnub.ctx)
}

func newHeartbeatBuilderWithContext(pubnub *PubNub,
	context Context) *heartbeatBuilder {
	builder := heartbeatBuilder{
		opts: newHeartbeatOpts(
			pubnub,
			context,
		),
	}

	return &builder
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *heartbeatBuilder) QueryParam(queryParam map[string]string) *heartbeatBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// State sets the state for the Heartbeat request.
func (b *heartbeatBuilder) State(state interface{}) *heartbeatBuilder {
	b.opts.State = state

	return b
}

// Channels sets the Channels for the Heartbeat request.
func (b *heartbeatBuilder) Channels(ch []string) *heartbeatBuilder {
	b.opts.Channels = ch

	return b
}

// ChannelGroups sets the ChannelGroups for the Heartbeat request.
func (b *heartbeatBuilder) ChannelGroups(cg []string) *heartbeatBuilder {
	b.opts.ChannelGroups = cg

	return b
}

// Execute runs the Heartbeat request
func (b *heartbeatBuilder) Execute() (interface{}, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return "", status, err
	}

	var value interface{}

	err = json.Unmarshal(rawJSON, &value)
	if err != nil {
		return nil, status, err
	}

	return value, status, nil
}

func newHeartbeatOpts(pubnub *PubNub, ctx Context) *heartbeatOpts {
	return &heartbeatOpts{
		endpointOpts: endpointOpts{
			pubnub: pubnub,
			ctx:    ctx,
		},
	}
}

type heartbeatOpts struct {
	endpointOpts

	State interface{}

	Channels      []string
	ChannelGroups []string
	QueryParam    map[string]string
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
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

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
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *heartbeatOpts) operationType() OperationType {
	return PNHeartBeatOperation
}
