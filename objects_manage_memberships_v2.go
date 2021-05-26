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

var emptyManageMembershipsResponse *PNManageMembershipsResponse

const manageMembershipsPathV2 = "/v2/objects/%s/uuids/%s/channels"

const manageMembershipsLimitV2 = 100

type manageMembershipsBuilderV2 struct {
	opts *manageMembershipsOptsV2
}

func newManageMembershipsBuilderV2(pubnub *PubNub) *manageMembershipsBuilderV2 {
	builder := manageMembershipsBuilderV2{
		opts: &manageMembershipsOptsV2{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = manageMembershipsLimitV2

	return &builder
}

func newManageMembershipsBuilderV2WithContext(pubnub *PubNub,
	context Context) *manageMembershipsBuilderV2 {
	builder := manageMembershipsBuilderV2{
		opts: &manageMembershipsOptsV2{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *manageMembershipsBuilderV2) Include(include []PNMembershipsInclude) *manageMembershipsBuilderV2 {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *manageMembershipsBuilderV2) UUID(uuid string) *manageMembershipsBuilderV2 {
	b.opts.UUID = uuid

	return b
}

func (b *manageMembershipsBuilderV2) Limit(limit int) *manageMembershipsBuilderV2 {
	b.opts.Limit = limit

	return b
}

func (b *manageMembershipsBuilderV2) Start(start string) *manageMembershipsBuilderV2 {
	b.opts.Start = start

	return b
}

func (b *manageMembershipsBuilderV2) End(end string) *manageMembershipsBuilderV2 {
	b.opts.End = end

	return b
}

func (b *manageMembershipsBuilderV2) Count(count bool) *manageMembershipsBuilderV2 {
	b.opts.Count = count

	return b
}

func (b *manageMembershipsBuilderV2) Filter(filter string) *manageMembershipsBuilderV2 {
	b.opts.Filter = filter

	return b
}

func (b *manageMembershipsBuilderV2) Sort(sort []string) *manageMembershipsBuilderV2 {
	b.opts.Sort = sort

	return b
}

func (b *manageMembershipsBuilderV2) Set(membershipsSet []PNMembershipsSet) *manageMembershipsBuilderV2 {
	b.opts.MembershipsSet = membershipsSet

	return b
}

func (b *manageMembershipsBuilderV2) Remove(membershipsRemove []PNMembershipsRemove) *manageMembershipsBuilderV2 {
	b.opts.MembershipsRemove = membershipsRemove

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *manageMembershipsBuilderV2) QueryParam(queryParam map[string]string) *manageMembershipsBuilderV2 {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the manageMemberships request.
func (b *manageMembershipsBuilderV2) Transport(tr http.RoundTripper) *manageMembershipsBuilderV2 {
	b.opts.Transport = tr
	return b
}

// Execute runs the manageMemberships request.
func (b *manageMembershipsBuilderV2) Execute() (*PNManageMembershipsResponse, StatusResponse, error) {
	if len(b.opts.UUID) <= 0 {
		b.opts.UUID = b.opts.pubnub.Config.UUID
	}

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyManageMembershipsResponse, status, err
	}

	return newPNManageMembershipsResponse(rawJSON, b.opts, status)
}

type manageMembershipsOptsV2 struct {
	pubnub            *PubNub
	UUID              string
	Limit             int
	Include           []string
	Start             string
	End               string
	Filter            string
	Sort              []string
	Count             bool
	QueryParam        map[string]string
	MembershipsRemove []PNMembershipsRemove
	MembershipsSet    []PNMembershipsSet
	Transport         http.RoundTripper

	ctx Context
}

func (o *manageMembershipsOptsV2) config() Config {
	return *o.pubnub.Config
}

func (o *manageMembershipsOptsV2) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *manageMembershipsOptsV2) context() Context {
	return o.ctx
}

func (o *manageMembershipsOptsV2) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *manageMembershipsOptsV2) buildPath() (string, error) {
	return fmt.Sprintf(manageMembershipsPathV2,
		o.pubnub.Config.SubscribeKey, o.UUID), nil
}

func (o *manageMembershipsOptsV2) buildQuery() (*url.Values, error) {

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

func (o *manageMembershipsOptsV2) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *manageMembershipsOptsV2) buildBody() ([]byte, error) {
	b := &PNManageMembershipsBody{
		Set:    o.MembershipsSet,
		Remove: o.MembershipsRemove,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *manageMembershipsOptsV2) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *manageMembershipsOptsV2) httpMethod() string {
	return "PATCH"
}

func (o *manageMembershipsOptsV2) isAuthRequired() bool {
	return true
}

func (o *manageMembershipsOptsV2) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *manageMembershipsOptsV2) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *manageMembershipsOptsV2) operationType() OperationType {
	return PNManageMembershipsOperation
}

func (o *manageMembershipsOptsV2) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNManageMembershipsResponse is the Objects API Response for ManageMemberships
type PNManageMembershipsResponse struct {
	status     int             `json:"status"`
	Data       []PNMemberships `json:"data"`
	TotalCount int             `json:"totalCount"`
	Next       string          `json:"next"`
	Prev       string          `json:"prev"`
}

func newPNManageMembershipsResponse(jsonBytes []byte, o *manageMembershipsOptsV2,
	status StatusResponse) (*PNManageMembershipsResponse, StatusResponse, error) {

	resp := &PNManageMembershipsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyManageMembershipsResponse, status, e
	}

	return resp, status, nil
}
