package pubnub

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v5/pnerr"
)

const timePath = "/time/0"

var emptyTimeResp *TimeResponse

type timeBuilder struct {
	opts *timeOpts
}

func newTimeBuilder(pubnub *PubNub) *timeBuilder {
	builder := timeBuilder{
		opts: &timeOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newTimeBuilderWithContext(pubnub *PubNub, context Context) *timeBuilder {
	builder := timeBuilder{
		opts: &timeOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Transport sets the Transport for the request.
func (b *timeBuilder) Transport(tr http.RoundTripper) *timeBuilder {
	b.opts.Transport = tr
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *timeBuilder) QueryParam(queryParam map[string]string) *timeBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Time request and fetches the time from the server.
func (b *timeBuilder) Execute() (*TimeResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyTimeResp, status, err
	}

	return newTimeResponse(rawJSON, status)
}

type timeOpts struct {
	pubnub     *PubNub
	QueryParam map[string]string
	Transport  http.RoundTripper

	ctx Context
}

func (o *timeOpts) config() Config {
	return *o.pubnub.Config
}

func (o *timeOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *timeOpts) context() Context {
	return o.ctx
}

func (o *timeOpts) validate() error {
	return nil
}

func (o *timeOpts) buildPath() (string, error) {
	return timePath, nil
}

func (o *timeOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *timeOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *timeOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *timeOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *timeOpts) httpMethod() string {
	return "GET"
}

func (o *timeOpts) isAuthRequired() bool {
	return false
}

func (o *timeOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *timeOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *timeOpts) operationType() OperationType {
	return PNTimeOperation
}

func (o *timeOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// TimeResponse is the response when Time call is executed.
type TimeResponse struct {
	Timetoken int64
}

func newTimeResponse(jsonBytes []byte, status StatusResponse) (*TimeResponse, StatusResponse, error) {
	resp := &TimeResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyTimeResp, status, e
	}

	if parsedValue, ok := value.([]interface{}); ok {
		if tt, ok := parsedValue[0].(float64); ok {
			resp.Timetoken = int64(tt)
		}
	}

	return resp, status, nil
}
