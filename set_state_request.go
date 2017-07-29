package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
)

const SET_STATE_PATH = "/v2/presence/sub-key/%s/channel/%s/uuid/%s/data"

var emptySetStateResponse *SetStateResponse

func SetStateRequest(pn *PubNub, opts *LeaveOpts) (*SetStateResponse, error) {
	opts.pubnub = pn
	rawJson, err := executeRequest(opts)
	if err != nil {
		return emptySetStateResponse, err
	}

	return newSetStateResponse(rawJson)
}

func SetStateRequestWithContext(ctx Context, pn *PubNub,
	opts *HistoryOpts) (*SetStateResponse, error) {
	opts.pubnub = pn
	opts.ctx = ctx

	_, err := executeRequest(opts)
	if err != nil {
		return emptySetStateResponse, err
	}

	return emptySetStateResponse, nil
}

type SetStateOpts struct {
	State         interface{}
	Channels      []string
	ChannelGroups []string

	pubnub *PubNub
	ctx    Context
}

func (o *SetStateOpts) config() Config {
	return *o.pubnub.Config
}

func (o *SetStateOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *SetStateOpts) context() Context {
	return o.ctx
}

func (o *SetStateOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return ErrMissingSubKey
	}

	if len(o.Channels) == 0 || len(o.ChannelGroups) == 0 {
		return pnerr.NewValidationError("Channel or channel group is missing")
	}

	log.Println(o.State)

	if o.State == nil {
		return pnerr.NewValidationError("State missing")
	}

	return nil
}

func (o *SetStateOpts) buildPath() (string, error) {
	channels := utils.JoinChannels(o.Channels)

	return fmt.Sprintf(SET_STATE_PATH,
		o.pubnub.Config.SubscribeKey,
		channels,
		o.pubnub.Config.Uuid,
	), nil
}

func (o *SetStateOpts) buildQuery() (*url.Values, error) {
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

func (o *SetStateOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *SetStateOpts) httpMethod() string {
	return "GET"
}

func (o *SetStateOpts) isAuthRequired() bool {
	return true
}

func (o *SetStateOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *SetStateOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func newSetStateResponse(jsonBytes []byte) (*SetStateResponse, error) {
	resp := &SetStateResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptySetStateResponse, e
	}

	v, _ := value.(map[string]interface{})
	val, _ := v["payload"].([]interface{})

	resp.State = val

	return resp, nil
}

type SetStateResponse struct {
	State []interface{}
}
