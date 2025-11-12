package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v8/pnerr"
	"github.com/pubnub/go/v8/utils"
)

const allChannelGroupPath = "/v1/channel-registration/sub-key/%s/channel-group/%s"

var emptyAllChannelGroupResponse *AllChannelGroupResponse

type allChannelGroupBuilder struct {
	opts *allChannelGroupOpts
}

func newAllChannelGroupBuilder(pubnub *PubNub) *allChannelGroupBuilder {
	return newAllChannelGroupBuilderWithContext(pubnub, pubnub.ctx)
}

func newAllChannelGroupOpts(pubnub *PubNub, ctx Context) *allChannelGroupOpts {
	return &allChannelGroupOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}
func newAllChannelGroupBuilderWithContext(pubnub *PubNub,
	context Context) *allChannelGroupBuilder {
	builder := allChannelGroupBuilder{
		opts: newAllChannelGroupOpts(pubnub, context)}
	return &builder
}

// ChannelGroup sets the channel group to list channels.
func (b *allChannelGroupBuilder) ChannelGroup(
	cg string) *allChannelGroupBuilder {
	b.opts.ChannelGroup = cg
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *allChannelGroupBuilder) QueryParam(queryParam map[string]string) *allChannelGroupBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *allChannelGroupOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ChannelGroup": o.ChannelGroup,
	}
}

// Execute runs the ListChannelsInChannelGroup request.
func (b *allChannelGroupBuilder) Execute() (
	*AllChannelGroupResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNChannelsForGroupOperation, b.opts.GetLogParams(), true)
	
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyAllChannelGroupResponse, status, err
	}

	return newAllChannelGroupResponse(rawJSON, status)
}

type allChannelGroupOpts struct {
	endpointOpts

	ChannelGroup string
	QueryParam   map[string]string
	Transport    http.RoundTripper
}

func (o *allChannelGroupOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.ChannelGroup == "" {
		return newValidationError(o, StrMissingChannelGroup)
	}

	return nil
}

func (o *allChannelGroupOpts) buildPath() (string, error) {
	return fmt.Sprintf(allChannelGroupPath,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.ChannelGroup)), nil
}

func (o *allChannelGroupOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *allChannelGroupOpts) isAuthRequired() bool {
	return true
}

func (o *allChannelGroupOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *allChannelGroupOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *allChannelGroupOpts) operationType() OperationType {
	return PNChannelsForGroupOperation
}

// AllChannelGroupResponse is the struct returned when the Execute function of List All Channel Groups is called.
type AllChannelGroupResponse struct {
	Channels     []string
	ChannelGroup string
}

func newAllChannelGroupResponse(jsonBytes []byte, status StatusResponse) (
	*AllChannelGroupResponse, StatusResponse, error) {
	resp := &AllChannelGroupResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyAllChannelGroupResponse, status, e
	}

	if parsedValue, ok := value.(map[string]interface{}); ok {
		if payload, ok := parsedValue["payload"].(map[string]interface{}); ok {
			if group, ok := payload["group"].(string); ok {
				resp.ChannelGroup = group
			}

			if channels, ok := payload["channels"].([]interface{}); ok {
				parsedChannels := []string{}

				for _, channel := range channels {
					if ch, ok := channel.(string); ok {
						parsedChannels = append(parsedChannels, ch)
					}
				}

				resp.Channels = parsedChannels
			}
		}
	}

	return resp, status, nil
}
