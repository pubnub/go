package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
)

var emptyGetMembershipsResponse *PNGetMembershipsResponse

const getMembershipsPath = "/v1/objects/%s/users/%s/spaces"

const spaceMembershipLimit = 100

type getMembershipsBuilder struct {
	opts *getMembershipsOpts
}

func newGetMembershipsBuilder(pubnub *PubNub) *getMembershipsBuilder {
	builder := getMembershipsBuilder{
		opts: &getMembershipsOpts{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = spaceMembershipLimit

	return &builder
}

func newGetMembershipsBuilderWithContext(pubnub *PubNub,
	context Context) *getMembershipsBuilder {
	builder := getMembershipsBuilder{
		opts: &getMembershipsOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *getMembershipsBuilder) UserID(id string) *getMembershipsBuilder {
	b.opts.ID = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getMembershipsBuilder) Include(include []PNMembershipsInclude) *getMembershipsBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getMembershipsBuilder) Limit(limit int) *getMembershipsBuilder {
	b.opts.Limit = limit

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getMembershipsBuilder) Start(start string) *getMembershipsBuilder {
	b.opts.Start = start

	return b
}

func (b *getMembershipsBuilder) End(end string) *getMembershipsBuilder {
	b.opts.End = end

	return b
}

func (b *getMembershipsBuilder) Count(count bool) *getMembershipsBuilder {
	b.opts.Count = count

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getMembershipsBuilder) QueryParam(queryParam map[string]string) *getMembershipsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getMemberships request.
func (b *getMembershipsBuilder) Transport(tr http.RoundTripper) *getMembershipsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getMemberships request.
func (b *getMembershipsBuilder) Execute() (*PNGetMembershipsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetMembershipsResponse, status, err
	}

	return newPNGetMembershipsResponse(rawJSON, b.opts, status)
}

type getMembershipsOpts struct {
	pubnub     *PubNub
	ID         string
	Limit      int
	Include    []string
	Start      string
	End        string
	Count      bool
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getMembershipsOpts) config() Config {
	return *o.pubnub.Config
}

func (o *getMembershipsOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getMembershipsOpts) context() Context {
	return o.ctx
}

func (o *getMembershipsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getMembershipsOpts) buildPath() (string, error) {
	return fmt.Sprintf(getMembershipsPath,
		o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *getMembershipsOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
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
	o.pubnub.tokenManager.SetAuthParan(q, o.ID, PNUsers)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getMembershipsOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getMembershipsOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getMembershipsOpts) httpMethod() string {
	return "GET"
}

func (o *getMembershipsOpts) isAuthRequired() bool {
	return true
}

func (o *getMembershipsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getMembershipsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getMembershipsOpts) operationType() OperationType {
	return PNGetMembershipsOperation
}

func (o *getMembershipsOpts) telemetryManager() *TelemetryManager {
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

func newPNGetMembershipsResponse(jsonBytes []byte, o *getMembershipsOpts,
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
