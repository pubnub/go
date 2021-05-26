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

var emptyGetAllChannelMetadataResponse *PNGetAllChannelMetadataResponse

const getAllChannelMetadataPath = "/v2/objects/%s/channels"

const getAllChannelMetadataLimitV2 = 100

type getAllChannelMetadataBuilder struct {
	opts *getAllChannelMetadataOpts
}

func newGetAllChannelMetadataBuilder(pubnub *PubNub) *getAllChannelMetadataBuilder {
	builder := getAllChannelMetadataBuilder{
		opts: &getAllChannelMetadataOpts{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = getAllChannelMetadataLimitV2

	return &builder
}

func newGetAllChannelMetadataBuilderWithContext(pubnub *PubNub,
	context Context) *getAllChannelMetadataBuilder {
	builder := getAllChannelMetadataBuilder{
		opts: &getAllChannelMetadataOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *getAllChannelMetadataBuilder) Include(include []PNChannelMetadataInclude) *getAllChannelMetadataBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getAllChannelMetadataBuilder) Limit(limit int) *getAllChannelMetadataBuilder {
	b.opts.Limit = limit

	return b
}

func (b *getAllChannelMetadataBuilder) Start(start string) *getAllChannelMetadataBuilder {
	b.opts.Start = start

	return b
}

func (b *getAllChannelMetadataBuilder) End(end string) *getAllChannelMetadataBuilder {
	b.opts.End = end

	return b
}

func (b *getAllChannelMetadataBuilder) Filter(filter string) *getAllChannelMetadataBuilder {
	b.opts.Filter = filter

	return b
}

func (b *getAllChannelMetadataBuilder) Sort(sort []string) *getAllChannelMetadataBuilder {
	b.opts.Sort = sort

	return b
}

func (b *getAllChannelMetadataBuilder) Count(count bool) *getAllChannelMetadataBuilder {
	b.opts.Count = count

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getAllChannelMetadataBuilder) QueryParam(queryParam map[string]string) *getAllChannelMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getAllChannelMetadata request.
func (b *getAllChannelMetadataBuilder) Transport(tr http.RoundTripper) *getAllChannelMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getAllChannelMetadata request.
func (b *getAllChannelMetadataBuilder) Execute() (*PNGetAllChannelMetadataResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetAllChannelMetadataResponse, status, err
	}

	return newPNGetAllChannelMetadataResponse(rawJSON, b.opts, status)
}

type getAllChannelMetadataOpts struct {
	pubnub *PubNub

	Limit      int
	Include    []string
	Start      string
	End        string
	Filter     string
	Sort       []string
	Count      bool
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getAllChannelMetadataOpts) config() Config {
	return *o.pubnub.Config
}

func (o *getAllChannelMetadataOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getAllChannelMetadataOpts) context() Context {
	return o.ctx
}

func (o *getAllChannelMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getAllChannelMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(getAllChannelMetadataPath,
		o.pubnub.Config.SubscribeKey), nil
}

func (o *getAllChannelMetadataOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		SetQueryParamAsCommaSepString(q, o.Include, "include")
	}

	q.Set("limit", strconv.Itoa(o.Limit))

	if o.Start != "" {
		q.Set("start", o.Start)
	}

	if o.Count {
		q.Set("count", "1")
	} else {
		q.Set("count", "0")
	}

	if o.End != "" {
		q.Set("end", o.End)
	}
	if o.Filter != "" {
		q.Set("filter", o.Filter)
	}
	if o.Sort != nil {
		SetQueryParamAsCommaSepString(q, o.Sort, "sort")
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getAllChannelMetadataOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getAllChannelMetadataOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getAllChannelMetadataOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *getAllChannelMetadataOpts) httpMethod() string {
	return "GET"
}

func (o *getAllChannelMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *getAllChannelMetadataOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getAllChannelMetadataOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getAllChannelMetadataOpts) operationType() OperationType {
	return PNGetAllChannelMetadataOperation
}

func (o *getAllChannelMetadataOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNGetAllChannelMetadataResponse is the Objects API Response for Get Spaces
type PNGetAllChannelMetadataResponse struct {
	status     int         `json:"status"`
	Data       []PNChannel `json:"data"`
	TotalCount int         `json:"totalCount"`
	Next       string      `json:"next"`
	Prev       string      `json:"prev"`
}

func newPNGetAllChannelMetadataResponse(jsonBytes []byte, o *getAllChannelMetadataOpts,
	status StatusResponse) (*PNGetAllChannelMetadataResponse, StatusResponse, error) {

	resp := &PNGetAllChannelMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyGetAllChannelMetadataResponse, status, e
	}

	return resp, status, nil
}
