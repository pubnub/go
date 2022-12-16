package pubnub

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v7/utils"
)

const deleteChannelGroupPath = "/v1/channel-registration/sub-key/%s/channel-group/%s/remove"

var emptyDeleteChannelGroupResponse *DeleteChannelGroupResponse

type deleteChannelGroupBuilder struct {
	opts *deleteChannelGroupOpts
}

func newDeleteChannelGroupBuilder(pubnub *PubNub) *deleteChannelGroupBuilder {
	return newDeleteChannelGroupBuilderWithContext(pubnub, pubnub.ctx)
}

func newDeleteChannelGroupBuilderWithContext(
	pubnub *PubNub, context Context) *deleteChannelGroupBuilder {
	builder := deleteChannelGroupBuilder{
		opts: newDeleteChannelGroupOpts(
			pubnub,
			context,
			deleteChannelGroupOpts{},
		),
	}

	return &builder
}

// ChannelGroup sets the channel group to delete.
func (b *deleteChannelGroupBuilder) ChannelGroup(
	cg string) *deleteChannelGroupBuilder {
	b.opts.ChannelGroup = cg
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *deleteChannelGroupBuilder) QueryParam(queryParam map[string]string) *deleteChannelGroupBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the DeleteChannelGroup request.
func (b *deleteChannelGroupBuilder) Execute() (
	*DeleteChannelGroupResponse, StatusResponse, error) {
	_, status, err := executeRequest(b.opts)

	if err != nil {
		return emptyDeleteChannelGroupResponse, status, err
	}

	return emptyDeleteChannelGroupResponse, status, nil
}

func newDeleteChannelGroupOpts(pubnub *PubNub, ctx Context, opts deleteChannelGroupOpts) *deleteChannelGroupOpts {
	opts.endpointOpts = endpointOpts{
		pubnub: pubnub,
		ctx:    ctx,
	}
	return &opts
}

type deleteChannelGroupOpts struct {
	endpointOpts
	ChannelGroup string
	Transport    http.RoundTripper
	QueryParam   map[string]string
}

func (o *deleteChannelGroupOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.ChannelGroup == "" {
		return newValidationError(o, StrMissingChannelGroup)
	}

	return nil
}

// DeleteChannelGroupResponse is response structure for Delete Channel Group function
type DeleteChannelGroupResponse struct{}

func (o *deleteChannelGroupOpts) buildPath() (string, error) {
	return fmt.Sprintf(deleteChannelGroupPath,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.ChannelGroup)), nil
}

func (o *deleteChannelGroupOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *deleteChannelGroupOpts) httpMethod() string {
	return "GET"
}

func (o *deleteChannelGroupOpts) operationType() OperationType {
	return PNRemoveGroupOperation
}
