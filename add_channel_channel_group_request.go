package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pubnub/go/v8/utils"
)

const addChannelToChannelGroupPath = "/v1/channel-registration/sub-key/%s/channel-group/%s"

var emptyAddChannelToChannelGroupResp *AddChannelToChannelGroupResponse

// addChannelToChannelGroupBuilder provides a builder to add channel to a channel group
type addChannelToChannelGroupBuilder struct {
	opts *addChannelOpts
}

func newAddChannelToChannelGroupBuilder(pubnub *PubNub) *addChannelToChannelGroupBuilder {
	return newAddChannelToChannelGroupBuilderWithContext(pubnub, pubnub.ctx)
}

func newAddChannelToChannelGroupBuilderWithContext(pubnub *PubNub, context Context) *addChannelToChannelGroupBuilder {
	builder := addChannelToChannelGroupBuilder{
		opts: newAddChannelOpts(pubnub, context),
	}

	return &builder
}

// Channels sets the channels to add to the channel group
func (b *addChannelToChannelGroupBuilder) Channels(
	ch []string) *addChannelToChannelGroupBuilder {

	b.opts.Channels = ch

	return b
}

// ChannelGroup sets the channel group to add the channels
func (b *addChannelToChannelGroupBuilder) ChannelGroup(
	cg string) *addChannelToChannelGroupBuilder {
	b.opts.ChannelGroup = cg

	return b
}

// Transport sets the transport for the request
func (b *addChannelToChannelGroupBuilder) Transport(
	tr http.RoundTripper) *addChannelToChannelGroupBuilder {
	b.opts.Transport = tr

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *addChannelToChannelGroupBuilder) QueryParam(queryParam map[string]string) *addChannelToChannelGroupBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *addChannelOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"Channels":     o.Channels,
		"ChannelGroup": o.ChannelGroup,
	}
}

// Execute runs AddChannelToChannelGroup request
func (b *addChannelToChannelGroupBuilder) Execute() (
	*AddChannelToChannelGroupResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNAddChannelsToChannelGroupOperation, b.opts.GetLogParams(), true)
	
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyAddChannelToChannelGroupResp, status, err
	}

	return newAddChannelToChannelGroupsResponse(rawJSON, status)
}

func newAddChannelOpts(pubnub *PubNub, context Context) *addChannelOpts {
	return &addChannelOpts{
		endpointOpts: endpointOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}
}

type addChannelOpts struct {
	endpointOpts
	Channels     []string
	ChannelGroup string
	QueryParam   map[string]string
	Transport    http.RoundTripper
}

func (o *addChannelOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 {
		return newValidationError(o, StrMissingChannel)
	}

	if o.ChannelGroup == "" {
		return newValidationError(o, StrMissingChannelGroup)
	}

	return nil
}

func (o *addChannelOpts) buildPath() (string, error) {
	return fmt.Sprintf(addChannelToChannelGroupPath,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.ChannelGroup)), nil
}

func (o *addChannelOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	var channels []string

	for _, v := range o.Channels {
		channels = append(channels, v)
	}

	q.Set("add", strings.Join(channels, ","))
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *addChannelOpts) operationType() OperationType {
	return PNAddChannelsToChannelGroupOperation
}

// AddChannelToChannelGroupResponse is the struct returned when the Execute function of AddChannelToChannelGroup is called.
type AddChannelToChannelGroupResponse struct {
}

func newAddChannelToChannelGroupsResponse(jsonBytes []byte, status StatusResponse) (
	*AddChannelToChannelGroupResponse, StatusResponse, error) {

	return emptyAddChannelToChannelGroupResp, status, nil
}
