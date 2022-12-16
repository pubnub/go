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

var emptySetMembershipsResponse *PNSetMembershipsResponse

const setMembershipsPath = "/v2/objects/%s/uuids/%s/channels"

const setMembershipsLimit = 100

type setMembershipsBuilder struct {
	opts *setMembershipsOpts
}

func newSetMembershipsBuilder(pubnub *PubNub) *setMembershipsBuilder {
	return newSetMembershipsBuilderWithContext(pubnub, pubnub.ctx)
}

func newSetMembershipsBuilderWithContext(pubnub *PubNub,
	context Context) *setMembershipsBuilder {
	return &setMembershipsBuilder{
		opts: newSetMembershipsOpts(
			pubnub,
			context,
		),
	}
}

func (b *setMembershipsBuilder) Include(include []PNMembershipsInclude) *setMembershipsBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *setMembershipsBuilder) UUID(uuid string) *setMembershipsBuilder {
	b.opts.UUID = uuid

	return b
}

func (b *setMembershipsBuilder) Limit(limit int) *setMembershipsBuilder {
	b.opts.Limit = limit

	return b
}

func (b *setMembershipsBuilder) Start(start string) *setMembershipsBuilder {
	b.opts.Start = start

	return b
}

func (b *setMembershipsBuilder) End(end string) *setMembershipsBuilder {
	b.opts.End = end

	return b
}

func (b *setMembershipsBuilder) Count(count bool) *setMembershipsBuilder {
	b.opts.Count = count

	return b
}

func (b *setMembershipsBuilder) Filter(filter string) *setMembershipsBuilder {
	b.opts.Filter = filter

	return b
}

func (b *setMembershipsBuilder) Sort(sort []string) *setMembershipsBuilder {
	b.opts.Sort = sort

	return b
}

func (b *setMembershipsBuilder) Set(membershipSet []PNMembershipsSet) *setMembershipsBuilder {
	b.opts.MembershipsSet = membershipSet

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *setMembershipsBuilder) QueryParam(queryParam map[string]string) *setMembershipsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the setMemberships request.
func (b *setMembershipsBuilder) Transport(tr http.RoundTripper) *setMembershipsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the setMemberships request.
func (b *setMembershipsBuilder) Execute() (*PNSetMembershipsResponse, StatusResponse, error) {
	if len(b.opts.UUID) <= 0 {
		b.opts.UUID = b.opts.pubnub.Config.UUID
	}

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySetMembershipsResponse, status, err
	}

	return newPNSetMembershipsResponse(rawJSON, b.opts, status)
}

func newSetMembershipsOpts(pubnub *PubNub, ctx Context) *setMembershipsOpts {
	return &setMembershipsOpts{
		endpointOpts: endpointOpts{
			pubnub: pubnub,
			ctx:    ctx,
		},
		Limit: setMembershipsLimit}
}

type setMembershipsOpts struct {
	endpointOpts
	UUID           string
	Limit          int
	Include        []string
	Start          string
	End            string
	Filter         string
	Sort           []string
	Count          bool
	QueryParam     map[string]string
	MembershipsSet []PNMembershipsSet
	Transport      http.RoundTripper
}

func (o *setMembershipsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *setMembershipsOpts) buildPath() (string, error) {
	return fmt.Sprintf(setMembershipsPath,
		o.pubnub.Config.SubscribeKey, o.UUID), nil
}

func (o *setMembershipsOpts) buildQuery() (*url.Values, error) {

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

// PNMembersAddChangeSet is the Objects API input to add, remove or update members
type PNMembersAddChangeSet struct {
	Set []PNMembershipsSet `json:"set"`
}

func (o *setMembershipsOpts) buildBody() ([]byte, error) {
	b := &PNMembersAddChangeSet{
		Set: o.MembershipsSet,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

}

func (o *setMembershipsOpts) httpMethod() string {
	return "PATCH"
}

func (o *setMembershipsOpts) isAuthRequired() bool {
	return true
}

func (o *setMembershipsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *setMembershipsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *setMembershipsOpts) operationType() OperationType {
	return PNSetMembershipsOperation
}

// PNSetMembershipsResponse is the Objects API Response for SetMemberships
type PNSetMembershipsResponse struct {
	status     int             `json:"status"`
	Data       []PNMemberships `json:"data"`
	TotalCount int             `json:"totalCount"`
	Next       string          `json:"next"`
	Prev       string          `json:"prev"`
}

func newPNSetMembershipsResponse(jsonBytes []byte, o *setMembershipsOpts,
	status StatusResponse) (*PNSetMembershipsResponse, StatusResponse, error) {

	resp := &PNSetMembershipsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptySetMembershipsResponse, status, e
	}

	return resp, status, nil
}
