package pubnub

import (
	"fmt"
	"net/url"

	"github.com/pubnub/go/v8/utils"
)

const leavePath = "/v2/presence/sub-key/%s/channel/%s/leave"

type leaveBuilder struct {
	opts *leaveOpts
}

func newLeaveBuilder(pubnub *PubNub) *leaveBuilder {
	return newLeaveBuilderWithContext(pubnub, pubnub.ctx)
}

func newLeaveOpts(pubnub *PubNub, ctx Context) *leaveOpts {
	return &leaveOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}
func newLeaveBuilderWithContext(pubnub *PubNub, context Context) *leaveBuilder {
	builder := leaveBuilder{
		opts: newLeaveOpts(pubnub, context)}
	return &builder
}

// Channels sets the channel names in the Unsubscribe request.
func (b *leaveBuilder) Channels(channels []string) *leaveBuilder {
	b.opts.Channels = channels
	return b
}

// ChannelGroups sets the channel group names in the Unsubscribe request.
func (b *leaveBuilder) ChannelGroups(groups []string) *leaveBuilder {
	b.opts.ChannelGroups = groups
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *leaveBuilder) QueryParam(queryParam map[string]string) *leaveBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *leaveOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"Channels":      o.Channels,
		"ChannelGroups": o.ChannelGroups,
	}
}

// Execute runs the Leave request.
func (b *leaveBuilder) Execute() (StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNUnsubscribeOperation, b.opts.GetLogParams(), true)
	
	_, status, err := executeRequest(b.opts)
	if err != nil {
		return status, err
	}

	return status, nil
}

type leaveOpts struct {
	endpointOpts
	Channels      []string
	ChannelGroups []string
	QueryParam    map[string]string
}

func (o *leaveOpts) buildPath() (string, error) {
	channels := utils.JoinChannels(o.Channels)

	if string(channels) == "" {
		channels = []byte(",")
	}

	return fmt.Sprintf(leavePath,
		o.pubnub.Config.SubscribeKey,
		channels), nil
}

func (o *leaveOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if len(o.ChannelGroups) > 0 {
		channelGroup := utils.JoinChannels(o.ChannelGroups)
		q.Set("channel-group", string(channelGroup))
	}
	SetQueryParam(q, o.QueryParam)
	return q, nil
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
