package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
)

const GET_STATE_PATH = "/v2/presence/sub-key/%s/channel/%s/uuid/%s"

var emptyGetStateResp *GetStateResponse

type getStateBuilder struct {
	opts *getStateOpts
}

func newGetStateBuilder(pubnub *PubNub) *getStateBuilder {
	builder := getStateBuilder{
		opts: &getStateOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newGetStateBuilderWithContext(pubnub *PubNub,
	context Context) *getStateBuilder {
	builder := getStateBuilder{
		opts: &getStateOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *getStateBuilder) Channels(ch []string) *getStateBuilder {
	b.opts.Channels = ch

	return b
}

func (b *getStateBuilder) ChannelGroups(cg []string) *getStateBuilder {
	b.opts.ChannelGroups = cg

	return b
}

func (b *getStateBuilder) Transport(
	tr http.RoundTripper) *getStateBuilder {
	b.opts.Transport = tr

	return b
}

func (b *getStateBuilder) Execute() (
	*GetStateResponse, StatusResponse, error) {
	rawJson, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetStateResp, status, err
	}

	return newGetStateResponse(rawJson, status)
}

type getStateOpts struct {
	pubnub *PubNub

	Channels []string

	ChannelGroups []string

	Transport http.RoundTripper

	ctx Context
}

func (o *getStateOpts) config() Config {
	return *o.pubnub.Config
}

func (o *getStateOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getStateOpts) context() Context {
	return o.ctx
}

func (o *getStateOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 && len(o.ChannelGroups) == 0 {
		return newValidationError(o, "Missing Channel or Channel Group")
	}

	return nil
}

func (o *getStateOpts) buildPath() (string, error) {
	return fmt.Sprintf(GET_STATE_PATH,
		o.pubnub.Config.SubscribeKey,
		utils.PamEncode(strings.Join(o.Channels, ",")),
		utils.UrlEncode(o.pubnub.Config.Uuid)), nil
}

func (o *getStateOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	q.Set("channel-group", strings.Join(o.ChannelGroups, ","))

	return q, nil
}

func (o *getStateOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getStateOpts) httpMethod() string {
	return "GET"
}

func (o *getStateOpts) isAuthRequired() bool {
	return true
}

func (o *getStateOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getStateOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getStateOpts) operationType() OperationType {
	return PNGetStateOperation
}

type GetStateResponse struct {
	State map[string]interface{}
}

func newGetStateResponse(jsonBytes []byte, status StatusResponse) (
	*GetStateResponse, StatusResponse, error) {

	resp := &GetStateResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyGetStateResp, status, e
	}

	if parsedValue, ok := value.(map[string]interface{}); ok {
		if payload, ok := parsedValue["payload"].(map[string]interface{}); ok {
			resp.State = payload
		}
	}

	return resp, status, nil
}
