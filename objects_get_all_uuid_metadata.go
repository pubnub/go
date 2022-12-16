package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pubnub/go/v7/pnerr"
)

var emptyPNGetAllUUIDMetadataResponse *PNGetAllUUIDMetadataResponse

const getAllUUIDMetadataPath = "/v2/objects/%s/uuids"

const getAllUUIDMetadataLimitV2 = 100

type getAllUUIDMetadataBuilder struct {
	opts *getAllUUIDMetadataOpts
}

func newGetAllUUIDMetadataBuilder(pubnub *PubNub) *getAllUUIDMetadataBuilder {
	return newGetAllUUIDMetadataBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetAllUUIDMetadataBuilderWithContext(pubnub *PubNub,
	context Context) *getAllUUIDMetadataBuilder {
	builder := getAllUUIDMetadataBuilder{
		opts: newGetAllUUIDMetadataOpts(
			pubnub,
			context,
		),
	}
	builder.opts.Limit = getAllUUIDMetadataLimitV2

	return &builder
}

func (b *getAllUUIDMetadataBuilder) Include(include []PNUUIDMetadataInclude) *getAllUUIDMetadataBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getAllUUIDMetadataBuilder) Limit(limit int) *getAllUUIDMetadataBuilder {
	b.opts.Limit = limit

	return b
}

func (b *getAllUUIDMetadataBuilder) Start(start string) *getAllUUIDMetadataBuilder {
	b.opts.Start = start

	return b
}

func (b *getAllUUIDMetadataBuilder) End(end string) *getAllUUIDMetadataBuilder {
	b.opts.End = end

	return b
}

func (b *getAllUUIDMetadataBuilder) Filter(filter string) *getAllUUIDMetadataBuilder {
	b.opts.Filter = filter

	return b
}

func (b *getAllUUIDMetadataBuilder) Sort(sort []string) *getAllUUIDMetadataBuilder {
	b.opts.Sort = sort

	return b
}

func (b *getAllUUIDMetadataBuilder) Count(count bool) *getAllUUIDMetadataBuilder {
	b.opts.Count = count

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getAllUUIDMetadataBuilder) QueryParam(queryParam map[string]string) *getAllUUIDMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getAllUUIDMetadata request.
func (b *getAllUUIDMetadataBuilder) Transport(tr http.RoundTripper) *getAllUUIDMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getAllUUIDMetadata request.
func (b *getAllUUIDMetadataBuilder) Execute() (*PNGetAllUUIDMetadataResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNGetAllUUIDMetadataResponse, status, err
	}

	return newPNGetAllUUIDMetadataResponse(rawJSON, b.opts, status)
}

func newGetAllUUIDMetadataOpts(pubnub *PubNub, ctx Context) *getAllUUIDMetadataOpts {
	return &getAllUUIDMetadataOpts{
		endpointOpts: endpointOpts{
			pubnub: pubnub,
			ctx:    ctx,
		},
	}
}

type getAllUUIDMetadataOpts struct {
	endpointOpts

	Limit      int
	Include    []string
	Start      string
	End        string
	Filter     string
	Sort       []string
	Count      bool
	QueryParam map[string]string

	Transport http.RoundTripper
}

func (o *getAllUUIDMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getAllUUIDMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(getAllUUIDMetadataPath,
		o.pubnub.Config.SubscribeKey), nil
}

func (o *getAllUUIDMetadataOpts) buildQuery() (*url.Values, error) {

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

func (o *getAllUUIDMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *getAllUUIDMetadataOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getAllUUIDMetadataOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getAllUUIDMetadataOpts) operationType() OperationType {
	return PNGetAllUUIDMetadataOperation
}

// PNGetAllUUIDMetadataResponse is the Objects API Response for Get Users
type PNGetAllUUIDMetadataResponse struct {
	status     int      `json:"status"`
	Data       []PNUUID `json:"data"`
	TotalCount int      `json:"totalCount"`
	Next       string   `json:"next"`
	Prev       string   `json:"prev"`
}

func newPNGetAllUUIDMetadataResponse(jsonBytes []byte, o *getAllUUIDMetadataOpts,
	status StatusResponse) (*PNGetAllUUIDMetadataResponse, StatusResponse, error) {

	resp := &PNGetAllUUIDMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGetAllUUIDMetadataResponse, status, e
	}

	return resp, status, nil
}
