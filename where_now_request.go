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

	"github.com/pubnub/go/v5/pnerr"
)

var whereNowPath = "/v2/presence/sub-key/%s/uuid/%s"

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

// UUID sets the UUID to fetch the where now info.
func (b *whereNowBuilder) UUID(uuid string) *whereNowBuilder {
	b.opts.UUID = uuid

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *whereNowBuilder) QueryParam(queryParam map[string]string) *whereNowBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the WhereNow request.
func (b *whereNowBuilder) Execute() (*WhereNowResponse, StatusResponse, error) {
	if len(b.opts.UUID) <= 0 {
		b.opts.UUID = b.opts.pubnub.Config.UUID
	}

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyWhereNowResponse, status, err
	}

	return newWhereNowResponse(rawJSON, status)
}

type whereNowOpts struct {
	pubnub *PubNub

	UUID       string
	QueryParam map[string]string
	Transport  http.RoundTripper

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
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *whereNowOpts) buildPath() (string, error) {
	return fmt.Sprintf(whereNowPath,
		o.pubnub.Config.SubscribeKey,
		o.UUID), nil
}

func (o *whereNowOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *whereNowOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *whereNowOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *whereNowOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
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

func (o *whereNowOpts) operationType() OperationType {
	return PNWhereNowOperation
}

func (o *whereNowOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// WhereNowResponse is the response of the WhereNow request. Contains channels info.
type WhereNowResponse struct {
	Channels []string
}

func newWhereNowResponse(jsonBytes []byte, status StatusResponse) (
	*WhereNowResponse, StatusResponse, error) {
	resp := &WhereNowResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyWhereNowResponse, status, e
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

	return resp, status, nil
}
