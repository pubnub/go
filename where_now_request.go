package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
)

var WHERE_NOW_PATH = "/v2/presence/sub-key/%s/uuid/%s"

var emptyWhereNowResponse *WhereNowResponse

type whereNowBuilder struct {
	opts *whereNowOpts
}

func newWhereNowBuilder(pubnub *PubNub) *whereNowBuilder {
	builder := whereNowBuilder{
		opts: &whereNowOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newWhereNowBuilderWithContext(pubnub *PubNub,
	context Context) *whereNowBuilder {
	builder := whereNowBuilder{
		opts: &whereNowOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *whereNowBuilder) Uuid(uuid string) *whereNowBuilder {
	b.opts.Uuid = uuid

	return b
}

func (b *whereNowBuilder) Execute() (*WhereNowResponse, error) {
	rawJson, err := executeRequest(b.opts)
	if err != nil {
		return emptyWhereNowResponse, err
	}

	return newWhereNowResponse(rawJson)
}

type whereNowOpts struct {
	pubnub *PubNub

	Uuid string

	Transport http.RoundTripper

	ctx Context
}

func (o *whereNowOpts) config() Config {
	return *o.pubnub.Config
}

func (o *whereNowOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *whereNowOpts) context() Context {
	return o.ctx
}

func (o *whereNowOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return ErrMissingSubKey
	}

	if o.Uuid == "" {
		return ErrMissingUuid
	}

	return nil
}

func (o *whereNowOpts) buildPath() (string, error) {
	return fmt.Sprintf(WHERE_NOW_PATH,
		o.pubnub.Config.SubscribeKey,
		o.Uuid), nil
}

func (o *whereNowOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	return q, nil
}

func (o *whereNowOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *whereNowOpts) httpMethod() string {
	return "GET"
}

func (o *whereNowOpts) isAuthRequired() bool {
	return true
}

func (o *whereNowOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *whereNowOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *whereNowOpts) operationType() PNOperationType {
	return PNWhereNowOperation
}

type WhereNowResponse struct {
	Channels []string
}

func newWhereNowResponse(jsonBytes []byte) (*WhereNowResponse, error) {
	resp := &WhereNowResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyWhereNowResponse, e
	}

	if parsedValue, ok := value.(map[string]interface{}); ok {
		if payload, ok := parsedValue["payload"].(map[string]interface{}); ok {
			if channels, ok := payload["channels"].([]interface{}); ok {
				for _, ch := range channels {
					if channel, ok := ch.(string); ok {
						resp.Channels = append(resp.Channels, channel)
					}
				}
			}
		}
	}

	return resp, nil
}
