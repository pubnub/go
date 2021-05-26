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

var emptyGetMembershipsResponse *PNGetMembershipsResponse

const getMembershipsPathV2 = "/v2/objects/%s/uuids/%s/channels"

const membershipsLimitV2 = 100

type getMembershipsBuilderV2 struct {
	opts *getMembershipsOptsV2
}

func newGetMembershipsBuilderV2(pubnub *PubNub) *getMembershipsBuilderV2 {
	builder := getMembershipsBuilderV2{
		opts: &getMembershipsOptsV2{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = membershipsLimitV2

	return &builder
}

func newGetMembershipsBuilderV2WithContext(pubnub *PubNub,
	context Context) *getMembershipsBuilderV2 {
	builder := getMembershipsBuilderV2{
		opts: &getMembershipsOptsV2{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *getMembershipsBuilderV2) UUID(uuid string) *getMembershipsBuilderV2 {
	b.opts.UUID = uuid

	return b
}

func (b *getMembershipsBuilderV2) Include(include []PNMembershipsInclude) *getMembershipsBuilderV2 {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getMembershipsBuilderV2) Limit(limit int) *getMembershipsBuilderV2 {
	b.opts.Limit = limit

	return b
}

func (b *getMembershipsBuilderV2) Start(start string) *getMembershipsBuilderV2 {
	b.opts.Start = start

	return b
}

func (b *getMembershipsBuilderV2) End(end string) *getMembershipsBuilderV2 {
	b.opts.End = end

	return b
}

func (b *getMembershipsBuilderV2) Filter(filter string) *getMembershipsBuilderV2 {
	b.opts.Filter = filter

	return b
}

func (b *getMembershipsBuilderV2) Sort(sort []string) *getMembershipsBuilderV2 {
	b.opts.Sort = sort

	return b
}

func (b *getMembershipsBuilderV2) Count(count bool) *getMembershipsBuilderV2 {
	b.opts.Count = count

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getMembershipsBuilderV2) QueryParam(queryParam map[string]string) *getMembershipsBuilderV2 {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getMemberships request.
func (b *getMembershipsBuilderV2) Transport(tr http.RoundTripper) *getMembershipsBuilderV2 {
	b.opts.Transport = tr
	return b
}

// Execute runs the getMemberships request.
func (b *getMembershipsBuilderV2) Execute() (*PNGetMembershipsResponse, StatusResponse, error) {
	if len(b.opts.UUID) <= 0 {
		b.opts.UUID = b.opts.pubnub.Config.UUID
	}

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetMembershipsResponse, status, err
	}

	return newPNGetMembershipsResponse(rawJSON, b.opts, status)
}

type getMembershipsOptsV2 struct {
	pubnub     *PubNub
	UUID       string
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

func (o *getMembershipsOptsV2) config() Config {
	return *o.pubnub.Config
}

func (o *getMembershipsOptsV2) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getMembershipsOptsV2) context() Context {
	return o.ctx
}

func (o *getMembershipsOptsV2) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getMembershipsOptsV2) buildPath() (string, error) {
	return fmt.Sprintf(getMembershipsPathV2,
		o.pubnub.Config.SubscribeKey, o.UUID), nil
}

func (o *getMembershipsOptsV2) buildQuery() (*url.Values, error) {

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

func (o *getMembershipsOptsV2) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getMembershipsOptsV2) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getMembershipsOptsV2) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *getMembershipsOptsV2) httpMethod() string {
	return "GET"
}

func (o *getMembershipsOptsV2) isAuthRequired() bool {
	return true
}

func (o *getMembershipsOptsV2) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getMembershipsOptsV2) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getMembershipsOptsV2) operationType() OperationType {
	return PNGetMembershipsOperation
}

func (o *getMembershipsOptsV2) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNGetMembershipsResponse is the Objects API Response for Get Memberships
type PNGetMembershipsResponse struct {
	status     int             `json:"status"`
	Data       []PNMemberships `json:"data"`
	TotalCount int             `json:"totalCount"`
	Next       string          `json:"next"`
	Prev       string          `json:"prev"`
}

func newPNGetMembershipsResponse(jsonBytes []byte, o *getMembershipsOptsV2,
	status StatusResponse) (*PNGetMembershipsResponse, StatusResponse, error) {

	resp := &PNGetMembershipsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyGetMembershipsResponse, status, e
	}

	return resp, status, nil
}
