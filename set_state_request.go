package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
)

const SET_STATE_PATH = "/v2/presence/sub-key/%s/channel/%s/uuid/%s/data"

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

func (b *setStateBuilder) State(state map[string]interface{}) *setStateBuilder {
	b.opts.State = state
	return b
}

func (b *setStateBuilder) Channels(channels []string) *setStateBuilder {
	b.opts.Channels = channels
	return b
}

func (b *setStateBuilder) ChannelGroups(groups []string) *setStateBuilder {
	b.opts.ChannelGroups = groups
	return b
}

func (b *setStateBuilder) Execute() (*SetStateResponse, StatusResponse, error) {
	stateOperation := StateOperation{}
	stateOperation.channels = b.opts.Channels
	stateOperation.channelGroups = b.opts.ChannelGroups
	stateOperation.state = b.opts.State

	b.opts.pubnub.subscriptionManager.adaptState(stateOperation)

	rawJson, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySetStateResponse, status, err
	}

	return newSetStateResponse(rawJson, status)
}

type setStateOpts struct {
	State         map[string]interface{}
	Channels      []string
	ChannelGroups []string

	pubnub *PubNub
	ctx    Context
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
		return ErrMissingSubKey
	}

	if len(o.Channels) == 0 && len(o.ChannelGroups) == 0 {
		return pnerr.NewValidationError("Channel or channel group is missing")
	}

	if o.State == nil {
		return pnerr.NewValidationError("State missing")
	}

	return nil
}

func (o *setStateOpts) buildPath() (string, error) {
	channels := string(utils.JoinChannels(o.Channels))

	return fmt.Sprintf(SET_STATE_PATH,
		o.pubnub.Config.SubscribeKey,
		channels,
		utils.UrlEncode(o.pubnub.Config.Uuid),
	), nil
}

func (o *setStateOpts) buildQuery() (*url.Values, error) {
	var err error
	var state, groups []byte

	q := defaultQuery(o.pubnub.Config.Uuid)

	state, err = json.Marshal(o.State)
	if err != nil {
		return nil, err
	}

	groups = utils.JoinChannels(o.ChannelGroups)

	if o.State != nil {
		q.Set("state", string(state))
	}

	if len(o.ChannelGroups) > 0 {
		q.Set("channel-group", string(groups))
	}

	return q, nil
}

func (o *setStateOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
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

func (o *setStateOpts) operationType() PNOperationType {
	return PNSetStateOperation
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

	v, _ := value.(map[string]interface{})
	val, _ := v["payload"].([]interface{})

	resp.State = val

	return resp, status, nil
}

type SetStateResponse struct {
	State []interface{}
}
