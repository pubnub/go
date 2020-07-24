package pubnub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
)

const setStatePath = "/v2/presence/sub-key/%s/channel/%s/uuid/%s/data"

var emptySetStateResponse *SetStateResponse

type setStateBuilder struct {
	opts *setStateOpts
}

func newSetStateBuilder(pubnub *PubNub) *setStateBuilder {
	builder := setStateBuilder{
		opts: &setStateOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newSetStateBuilderWithContext(pubnub *PubNub, context Context) *setStateBuilder {
	builder := setStateBuilder{
		opts: &setStateOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

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

// Execute runs the the Set State request and returns the SetStateResponse
func (b *setStateBuilder) Execute() (*SetStateResponse, StatusResponse, error) {
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
	State         map[string]interface{}
	Channels      []string
	ChannelGroups []string
	UUID          string
	QueryParam    map[string]string
	pubnub        *PubNub
	stringState   string
	ctx           Context
}

func (o *setStateOpts) config() Config {
	return *o.pubnub.Config
}

func (o *setStateOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *setStateOpts) context() Context {
	return o.ctx
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

func (o *setStateOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *setStateOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *setStateOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *setStateOpts) httpMethod() string {
	return "GET"
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

func (o *setStateOpts) operationType() OperationType {
	return PNSetStateOperation
}

func (o *setStateOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

func newSetStateResponse(jsonBytes []byte, status StatusResponse) (
	*SetStateResponse, StatusResponse, error) {
	resp := &SetStateResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptySetStateResponse, status, e
	}

	v, ok := value.(map[string]interface{})
	if !ok {
		return emptySetStateResponse, status, errors.New("Response parsing error")
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
		if val, ok := v["payload"].(interface{}); ok {
			resp.State = val
		}
	}
	resp.Message = message

	return resp, status, nil
}

// SetStateResponse is the response returned when the Execute function of SetState is called.
type SetStateResponse struct {
	State   interface{}
	Message string
}
