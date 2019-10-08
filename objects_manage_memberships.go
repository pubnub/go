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

var emptyManageMembershipsResponse *PNManageMembershipsResponse

const manageMembershipsPath = "/v1/objects/%s/users/%s/spaces"

const userSpaceMembershipsLimit = 100

type manageMembershipsBuilder struct {
	opts *manageMembershipsOpts
}

func newManageMembershipsBuilder(pubnub *PubNub) *manageMembershipsBuilder {
	builder := manageMembershipsBuilder{
		opts: &manageMembershipsOpts{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = spaceLimit

	return &builder
}

func newManageMembershipsBuilderWithContext(pubnub *PubNub,
	context Context) *manageMembershipsBuilder {
	builder := manageMembershipsBuilder{
		opts: &manageMembershipsOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *manageMembershipsBuilder) Include(include []PNMembershipsInclude) *manageMembershipsBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *manageMembershipsBuilder) UserID(id string) *manageMembershipsBuilder {
	b.opts.UserID = id

	return b
}

func (b *manageMembershipsBuilder) Limit(limit int) *manageMembershipsBuilder {
	b.opts.Limit = limit

	return b
}

func (b *manageMembershipsBuilder) Start(start string) *manageMembershipsBuilder {
	b.opts.Start = start

	return b
}

func (b *manageMembershipsBuilder) End(end string) *manageMembershipsBuilder {
	b.opts.End = end

	return b
}

func (b *manageMembershipsBuilder) Count(count bool) *manageMembershipsBuilder {
	b.opts.Count = count

	return b
}

func (b *manageMembershipsBuilder) Add(userMembershipInput []PNMembershipsInput) *manageMembershipsBuilder {
	b.opts.MembershipsAdd = userMembershipInput

	return b
}

func (b *manageMembershipsBuilder) Update(userMembershipInput []PNMembershipsInput) *manageMembershipsBuilder {
	b.opts.MembershipsUpdate = userMembershipInput

	return b
}

func (b *manageMembershipsBuilder) Remove(userMembershipRemove []PNMembershipsRemove) *manageMembershipsBuilder {
	b.opts.MembershipsRemove = userMembershipRemove

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *manageMembershipsBuilder) QueryParam(queryParam map[string]string) *manageMembershipsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the manageMemberships request.
func (b *manageMembershipsBuilder) Transport(tr http.RoundTripper) *manageMembershipsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the manageMemberships request.
func (b *manageMembershipsBuilder) Execute() (*PNManageMembershipsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyManageMembershipsResponse, status, err
	}

	return newPNManageMembershipsResponse(rawJSON, b.opts, status)
}

type manageMembershipsOpts struct {
	pubnub            *PubNub
	UserID            string
	Limit             int
	Include           []string
	Start             string
	End               string
	Count             bool
	QueryParam        map[string]string
	MembershipsRemove []PNMembershipsRemove
	MembershipsAdd    []PNMembershipsInput
	MembershipsUpdate []PNMembershipsInput
	Transport         http.RoundTripper

	ctx Context
}

func (o *manageMembershipsOpts) config() Config {
	return *o.pubnub.Config
}

func (o *manageMembershipsOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *manageMembershipsOpts) context() Context {
	return o.ctx
}

func (o *manageMembershipsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *manageMembershipsOpts) buildPath() (string, error) {
	return fmt.Sprintf(manageMembershipsPath,
		o.pubnub.Config.SubscribeKey, o.UserID), nil
}

func (o *manageMembershipsOpts) buildQuery() (*url.Values, error) {

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
	o.pubnub.tokenManager.SetAuthParan(q, o.UserID, PNUsers)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *manageMembershipsOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

// PNMembershipsInputChangeSet is the Objects API input to add, remove or update membership
type PNMembershipsInputChangeSet struct {
	Add    []PNMembershipsInput  `json:"add"`
	Update []PNMembershipsInput  `json:"update"`
	Remove []PNMembershipsRemove `json:"remove"`
}

func (o *manageMembershipsOpts) buildBody() ([]byte, error) {
	b := &PNMembershipsInputChangeSet{
		Add:    o.MembershipsAdd,
		Update: o.MembershipsUpdate,
		Remove: o.MembershipsRemove,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *manageMembershipsOpts) httpMethod() string {
	return "PATCH"
}

func (o *manageMembershipsOpts) isAuthRequired() bool {
	return true
}

func (o *manageMembershipsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *manageMembershipsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *manageMembershipsOpts) operationType() OperationType {
	return PNManageMembershipsOperation
}

func (o *manageMembershipsOpts) telemetryManager() *TelemetryManager {
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

func newPNManageMembershipsResponse(jsonBytes []byte, o *manageMembershipsOpts,
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
