package pubnub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/pubnub/go/v8/pnerr"
	"github.com/pubnub/go/v8/utils"
)

const setStatePath = "/v2/presence/sub-key/%s/channel/%s/uuid/%s/data"

var emptySetStateResponse *SetStateResponse

type setStateBuilder struct {
	opts *setStateOpts
}

func newSetStateBuilder(pubnub *PubNub) *setStateBuilder {
	return newSetStateBuilderWithContext(pubnub, pubnub.ctx)
}

func newSetStateOpts(pubnub *PubNub, ctx Context) *setStateOpts {
	return &setStateOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}
func newSetStateBuilderWithContext(pubnub *PubNub, context Context) *setStateBuilder {
	builder := setStateBuilder{
		opts: newSetStateOpts(pubnub, context)}
	return &builder
}

// State sets the State for the Set State request.
func (b *setStateBuilder) State(state map[string]interface{}) *setStateBuilder {
	b.opts.State = state
	return b
}

// Channels sets the Channels for the Set State request.
func (b *setStateBuilder) Channels(channels []string) *setStateBuilder {
	b.opts.Channels = channels
	return b
}

// ChannelGroups sets the ChannelGroups for the Set State request.
func (b *setStateBuilder) ChannelGroups(groups []string) *setStateBuilder {
	b.opts.ChannelGroups = groups
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *setStateBuilder) QueryParam(queryParam map[string]string) *setStateBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// UUID sets the UUID for the Set State request.
func (b *setStateBuilder) UUID(uuid string) *setStateBuilder {
	b.opts.UUID = uuid

	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *setStateOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"Channels":      o.Channels,
		"ChannelGroups": o.ChannelGroups,
	}
	if o.UUID != "" {
		params["UUID"] = o.UUID
	}
	if o.State != nil {
		params["State"] = fmt.Sprintf("%v", o.State)
	}
	return params
}

// Execute runs the the Set State request and returns the SetStateResponse
func (b *setStateBuilder) Execute() (*SetStateResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNSetStateOperation, b.opts.GetLogParams(), true)
	
	stateOperation := StateOperation{}
	stateOperation.channels = b.opts.Channels
	stateOperation.channelGroups = b.opts.ChannelGroups
	stateOperation.state = b.opts.State

	b.opts.pubnub.subscriptionManager.adaptState(stateOperation)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySetStateResponse, status, err
	}

	return newSetStateResponse(rawJSON, status)
}

type setStateOpts struct {
	endpointOpts
	State         map[string]interface{}
	Channels      []string
	ChannelGroups []string
	UUID          string
	QueryParam    map[string]string
	stringState   string
}

func (o *setStateOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 && len(o.ChannelGroups) == 0 {
		return newValidationError(o, "Missing Channel or Channel Group")
	}

	if o.State == nil {
		return newValidationError(o, "Missing State")
	}
	state, err := json.Marshal(o.State)
	if err != nil {
		return newValidationError(o, err.Error())
	}

	o.stringState = string(state)

	return nil
}

func (o *setStateOpts) buildPath() (string, error) {
	channels := string(utils.JoinChannels(o.Channels))
	uuid := o.UUID
	if uuid == "" {
		uuid = o.pubnub.Config.UUID
	}

	return fmt.Sprintf(setStatePath,
		o.pubnub.Config.SubscribeKey,
		channels,
		utils.URLEncode(uuid),
	), nil
}

func (o *setStateOpts) buildQuery() (*url.Values, error) {
	var groups []byte

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	groups = utils.JoinChannels(o.ChannelGroups)

	if o.stringState != "" {
		q.Set("state", o.stringState)
	}

	if len(o.ChannelGroups) > 0 {
		q.Set("channel-group", string(groups))
	}
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *setStateOpts) isAuthRequired() bool {
	return true
}

func (o *setStateOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *setStateOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *setStateOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *setStateOpts) httpMethod() string {
	return "GET"
}

func (o *setStateOpts) operationType() OperationType {
	return PNSetStateOperation
}

func newSetStateResponse(jsonBytes []byte, status StatusResponse) (
	*SetStateResponse, StatusResponse, error) {
	resp := &SetStateResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptySetStateResponse, status, e
	}

	v, ok := value.(map[string]interface{})
	if !ok {
		return emptySetStateResponse, status, errors.New("response parsing error")
	}
	message := ""
	if v["message"] != nil {
		if msg, ok := v["message"].(string); ok {
			message = msg
		}
	}

	if v["error"] != nil {
		return emptySetStateResponse, status, errors.New(message)
	}

	if v["payload"] != nil {
		resp.State = v["payload"]
	}
	resp.Message = message

	return resp, status, nil
}

// SetStateResponse is the response returned when the Execute function of SetState is called.
type SetStateResponse struct {
	State   interface{}
	Message string
}
