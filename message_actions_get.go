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
	"strconv"

	"github.com/pubnub/go/v5/pnerr"
)

var emptyPNGetMessageActionsResponse *PNGetMessageActionsResponse

const getMessageActionsPath = "/v1/message-actions/%s/channel/%s"

type getMessageActionsBuilder struct {
	opts *getMessageActionsOpts
}

func newGetMessageActionsBuilder(pubnub *PubNub) *getMessageActionsBuilder {
	builder := getMessageActionsBuilder{
		opts: &getMessageActionsOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newGetMessageActionsBuilderWithContext(pubnub *PubNub,
	context Context) *getMessageActionsBuilder {
	builder := getMessageActionsBuilder{
		opts: &getMessageActionsOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *getMessageActionsBuilder) Channel(channel string) *getMessageActionsBuilder {
	b.opts.Channel = channel

	return b
}

func (b *getMessageActionsBuilder) Start(timetoken string) *getMessageActionsBuilder {
	b.opts.Start = timetoken

	return b
}

func (b *getMessageActionsBuilder) End(timetoken string) *getMessageActionsBuilder {
	b.opts.End = timetoken

	return b
}

func (b *getMessageActionsBuilder) Limit(limit int) *getMessageActionsBuilder {
	b.opts.Limit = limit

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getMessageActionsBuilder) QueryParam(queryParam map[string]string) *getMessageActionsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getMessageActions request.
func (b *getMessageActionsBuilder) Transport(tr http.RoundTripper) *getMessageActionsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getMessageActions request.
func (b *getMessageActionsBuilder) Execute() (*PNGetMessageActionsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNGetMessageActionsResponse, status, err
	}

	return newPNGetMessageActionsResponse(rawJSON, b.opts, status)
}

type getMessageActionsOpts struct {
	pubnub *PubNub

	Channel    string
	Start      string
	End        string
	Limit      int
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getMessageActionsOpts) config() Config {
	return *o.pubnub.Config
}

func (o *getMessageActionsOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getMessageActionsOpts) context() Context {
	return o.ctx
}

func (o *getMessageActionsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getMessageActionsOpts) buildPath() (string, error) {
	return fmt.Sprintf(getMessageActionsPath,
		o.pubnub.Config.SubscribeKey, o.Channel), nil
}

func (o *getMessageActionsOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Start != "" {
		q.Set("start", o.Start)
	}

	if o.End != "" {
		q.Set("end", o.End)
	}

	if o.Limit > 0 {
		q.Set("limit", strconv.Itoa(o.Limit))
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getMessageActionsOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getMessageActionsOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getMessageActionsOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *getMessageActionsOpts) httpMethod() string {
	return "GET"
}

func (o *getMessageActionsOpts) isAuthRequired() bool {
	return true
}

func (o *getMessageActionsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getMessageActionsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getMessageActionsOpts) operationType() OperationType {
	return PNGetMessageActionsOperation
}

func (o *getMessageActionsOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNGetMessageActionsMore is the struct used when the PNGetMessageActionsResponse has more link
type PNGetMessageActionsMore struct {
	URL   string `json:"url"`
	Start string `json:"start"`
	End   string `json:"end"`
	Limit int    `json:"limit"`
}

// PNGetMessageActionsResponse is the GetMessageActions API Response
type PNGetMessageActionsResponse struct {
	status int                        `json:"status"`
	Data   []PNMessageActionsResponse `json:"data"`
	More   PNGetMessageActionsMore    `json:"more"`
}

func newPNGetMessageActionsResponse(jsonBytes []byte, o *getMessageActionsOpts,
	status StatusResponse) (*PNGetMessageActionsResponse, StatusResponse, error) {

	resp := &PNGetMessageActionsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGetMessageActionsResponse, status, e
	}

	return resp, status, nil
}
