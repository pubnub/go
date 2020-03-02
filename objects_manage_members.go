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

var emptyManageMembersResponse *PNManageMembersResponse

const manageMembersPath = "/v1/objects/%s/spaces/%s/users"

const manageMembersLimit = 100

type manageMembersBuilder struct {
	opts *manageMembersOpts
}

func newManageMembersBuilder(pubnub *PubNub) *manageMembersBuilder {
	builder := manageMembersBuilder{
		opts: &manageMembersOpts{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = spaceLimit

	return &builder
}

func newManageMembersBuilderWithContext(pubnub *PubNub,
	context Context) *manageMembersBuilder {
	builder := manageMembersBuilder{
		opts: &manageMembersOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *manageMembersBuilder) Include(include []PNMembersInclude) *manageMembersBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *manageMembersBuilder) SpaceID(id string) *manageMembersBuilder {
	b.opts.SpaceID = id

	return b
}

func (b *manageMembersBuilder) Limit(limit int) *manageMembersBuilder {
	b.opts.Limit = limit

	return b
}

func (b *manageMembersBuilder) Start(start string) *manageMembersBuilder {
	b.opts.Start = start

	return b
}

func (b *manageMembersBuilder) End(end string) *manageMembersBuilder {
	b.opts.End = end

	return b
}

func (b *manageMembersBuilder) Count(count bool) *manageMembersBuilder {
	b.opts.Count = count

	return b
}

func (b *manageMembersBuilder) Filter(filter string) *manageMembersBuilder {
	b.opts.Filter = filter

	return b
}

func (b *manageMembersBuilder) Sort(sort []string) *manageMembersBuilder {
	b.opts.Sort = sort

	return b
}

func (b *manageMembersBuilder) Add(membershipInput []PNMembersInput) *manageMembersBuilder {
	b.opts.MembershipAdd = membershipInput

	return b
}

func (b *manageMembersBuilder) Update(membershipInput []PNMembersInput) *manageMembersBuilder {
	b.opts.MembershipUpdate = membershipInput

	return b
}

func (b *manageMembersBuilder) Remove(membershipRemove []PNMembersRemove) *manageMembersBuilder {
	b.opts.MembershipRemove = membershipRemove

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *manageMembersBuilder) QueryParam(queryParam map[string]string) *manageMembersBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the manageMembers request.
func (b *manageMembersBuilder) Transport(tr http.RoundTripper) *manageMembersBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the manageMembers request.
func (b *manageMembersBuilder) Execute() (*PNManageMembersResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyManageMembersResponse, status, err
	}

	return newPNManageMembersResponse(rawJSON, b.opts, status)
}

type manageMembersOpts struct {
	pubnub           *PubNub
	SpaceID          string
	Limit            int
	Include          []string
	Start            string
	End              string
	Filter           string
	Sort             []string
	Count            bool
	QueryParam       map[string]string
	MembershipRemove []PNMembersRemove
	MembershipAdd    []PNMembersInput
	MembershipUpdate []PNMembersInput
	Transport        http.RoundTripper

	ctx Context
}

func (o *manageMembersOpts) config() Config {
	return *o.pubnub.Config
}

func (o *manageMembersOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *manageMembersOpts) context() Context {
	return o.ctx
}

func (o *manageMembersOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *manageMembersOpts) buildPath() (string, error) {
	return fmt.Sprintf(manageMembersPath,
		o.pubnub.Config.SubscribeKey, o.SpaceID), nil
}

func (o *manageMembersOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		SetArrayTypeQueryParam(q, o.Include, "include")
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
	if o.Sort != nil {
		SetArrayTypeQueryParam(q, o.Sort, "sort")
	}

	o.pubnub.tokenManager.SetAuthParan(q, o.SpaceID, PNSpaces)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *manageMembersOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

// PNMembersInputChangeSet is the Objects API input to add, remove or update members
type PNMembersInputChangeSet struct {
	Add    []PNMembersInput  `json:"add"`
	Update []PNMembersInput  `json:"update"`
	Remove []PNMembersRemove `json:"remove"`
}

func (o *manageMembersOpts) buildBody() ([]byte, error) {
	b := &PNMembersInputChangeSet{
		Add:    o.MembershipAdd,
		Update: o.MembershipUpdate,
		Remove: o.MembershipRemove,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

}

func (o *manageMembersOpts) httpMethod() string {
	return "PATCH"
}

func (o *manageMembersOpts) isAuthRequired() bool {
	return true
}

func (o *manageMembersOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *manageMembersOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *manageMembersOpts) operationType() OperationType {
	return PNManageMembersOperation
}

func (o *manageMembersOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNManageMembersResponse is the Objects API Response for ManageMembers
type PNManageMembersResponse struct {
	status     int         `json:"status"`
	Data       []PNMembers `json:"data"`
	TotalCount int         `json:"totalCount"`
	Next       string      `json:"next"`
	Prev       string      `json:"prev"`
}

func newPNManageMembersResponse(jsonBytes []byte, o *manageMembersOpts,
	status StatusResponse) (*PNManageMembersResponse, StatusResponse, error) {

	resp := &PNManageMembersResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyManageMembersResponse, status, e
	}

	return resp, status, nil
}
