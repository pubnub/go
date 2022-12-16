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

var emptyRemoveMembershipsResponse *PNRemoveMembershipsResponse

const removeMembershipsPath = "/v2/objects/%s/uuids/%s/channels"

const removeMembershipsLimit = 100

type removeMembershipsBuilder struct {
	opts *removeMembershipsOpts
}

func newRemoveMembershipsBuilder(pubnub *PubNub) *removeMembershipsBuilder {
	return newRemoveMembershipsBuilderWithContext(pubnub, pubnub.ctx)
}

func newRemoveMembershipsBuilderWithContext(pubnub *PubNub,
	context Context) *removeMembershipsBuilder {
	builder := removeMembershipsBuilder{
		opts: newRemoveMembershipsOpts(pubnub, context),
	}
	builder.opts.Limit = removeMembershipsLimit

	return &builder
}

func (b *removeMembershipsBuilder) Include(include []PNMembershipsInclude) *removeMembershipsBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *removeMembershipsBuilder) UUID(uuid string) *removeMembershipsBuilder {
	b.opts.UUID = uuid

	return b
}

func (b *removeMembershipsBuilder) Limit(limit int) *removeMembershipsBuilder {
	b.opts.Limit = limit

	return b
}

func (b *removeMembershipsBuilder) Start(start string) *removeMembershipsBuilder {
	b.opts.Start = start

	return b
}

func (b *removeMembershipsBuilder) End(end string) *removeMembershipsBuilder {
	b.opts.End = end

	return b
}

func (b *removeMembershipsBuilder) Count(count bool) *removeMembershipsBuilder {
	b.opts.Count = count

	return b
}

func (b *removeMembershipsBuilder) Filter(filter string) *removeMembershipsBuilder {
	b.opts.Filter = filter

	return b
}

func (b *removeMembershipsBuilder) Sort(sort []string) *removeMembershipsBuilder {
	b.opts.Sort = sort

	return b
}

func (b *removeMembershipsBuilder) Remove(membershipsRemove []PNMembershipsRemove) *removeMembershipsBuilder {
	b.opts.MembershipsRemove = membershipsRemove

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeMembershipsBuilder) QueryParam(queryParam map[string]string) *removeMembershipsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the removeMemberships request.
func (b *removeMembershipsBuilder) Transport(tr http.RoundTripper) *removeMembershipsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the removeMemberships request.
func (b *removeMembershipsBuilder) Execute() (*PNRemoveMembershipsResponse, StatusResponse, error) {
	if len(b.opts.UUID) <= 0 {
		b.opts.UUID = b.opts.pubnub.Config.UUID
	}

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveMembershipsResponse, status, err
	}

	return newPNRemoveMembershipsResponse(rawJSON, b.opts, status)
}

func newRemoveMembershipsOpts(pubnub *PubNub, ctx Context) *removeMembershipsOpts {
	return &removeMembershipsOpts{
		endpointOpts: endpointOpts{
			pubnub: pubnub,
			ctx:    ctx,
		},
	}
}

type removeMembershipsOpts struct {
	endpointOpts
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
	Transport         http.RoundTripper
}

func (o *removeMembershipsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *removeMembershipsOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeMembershipsPath,
		o.pubnub.Config.SubscribeKey, o.UUID), nil
}

func (o *removeMembershipsOpts) buildQuery() (*url.Values, error) {

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

// PNMembershipsRemoveChangeSet is the Objects API input to add, remove or update members
type PNMembershipsRemoveChangeSet struct {
	Remove []PNMembershipsRemove `json:"delete"`
}

func (o *removeMembershipsOpts) buildBody() ([]byte, error) {
	b := &PNMembershipsRemoveChangeSet{
		Remove: o.MembershipsRemove,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

}

func (o *removeMembershipsOpts) httpMethod() string {
	return "PATCH"
}

func (o *removeMembershipsOpts) isAuthRequired() bool {
	return true
}

func (o *removeMembershipsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *removeMembershipsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *removeMembershipsOpts) operationType() OperationType {
	return PNRemoveMembershipsOperation
}

// PNRemoveMembershipsResponse is the Objects API Response for RemoveMemberships
type PNRemoveMembershipsResponse struct {
	status     int             `json:"status"`
	Data       []PNMemberships `json:"data"`
	TotalCount int             `json:"totalCount"`
	Next       string          `json:"next"`
	Prev       string          `json:"prev"`
}

func newPNRemoveMembershipsResponse(jsonBytes []byte, o *removeMembershipsOpts,
	status StatusResponse) (*PNRemoveMembershipsResponse, StatusResponse, error) {

	resp := &PNRemoveMembershipsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyRemoveMembershipsResponse, status, e
	}

	return resp, status, nil
}
