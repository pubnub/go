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
	//"reflect"
)

var emptyUpdateUserSpaceMembershipsResponse *PNUpdateUserSpaceMembershipsResponse

const updateUserSpaceMembershipsPath = "/v1/objects/%s/users/%s/spaces"

const userSpaceMembershipsLimit = 100

type updateUserSpaceMembershipsBuilder struct {
	opts *updateUserSpaceMembershipsOpts
}

func newUpdateUserSpaceMembershipsBuilder(pubnub *PubNub) *updateUserSpaceMembershipsBuilder {
	builder := updateUserSpaceMembershipsBuilder{
		opts: &updateUserSpaceMembershipsOpts{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = spaceLimit

	return &builder
}

func newUpdateUserSpaceMembershipsBuilderWithContext(pubnub *PubNub,
	context Context) *updateUserSpaceMembershipsBuilder {
	builder := updateUserSpaceMembershipsBuilder{
		opts: &updateUserSpaceMembershipsOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Auth sets the Authorization key with permissions to perform the request.
// func (b *updateUserSpaceMembershipsBuilder) Auth(auth string) *updateUserSpaceMembershipsBuilder {
// 	//b.opts.Auth = auth

// 	return b
// }

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateUserSpaceMembershipsBuilder) Include(include []PNMembersInclude) *updateUserSpaceMembershipsBuilder {
	b.opts.Include = utils.EnumArrayToStringArray(fmt.Sprint(include))

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateUserSpaceMembershipsBuilder) UserId(id string) *updateUserSpaceMembershipsBuilder {
	b.opts.UserId = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateUserSpaceMembershipsBuilder) Limit(limit int) *updateUserSpaceMembershipsBuilder {
	b.opts.Limit = limit

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateUserSpaceMembershipsBuilder) Start(start string) *updateUserSpaceMembershipsBuilder {
	b.opts.Start = start

	return b
}

func (b *updateUserSpaceMembershipsBuilder) End(end string) *updateUserSpaceMembershipsBuilder {
	b.opts.End = end

	return b
}

func (b *updateUserSpaceMembershipsBuilder) Count(count bool) *updateUserSpaceMembershipsBuilder {
	b.opts.Count = count

	return b
}

func (b *updateUserSpaceMembershipsBuilder) Add(userMembershipInput []PNUserMembershipInput) *updateUserSpaceMembershipsBuilder {
	b.opts.UserMembershipAdd = userMembershipInput

	return b
}

func (b *updateUserSpaceMembershipsBuilder) Update(userMembershipInput []PNUserMembershipInput) *updateUserSpaceMembershipsBuilder {
	b.opts.UserMembershipUpdate = userMembershipInput

	return b
}

func (b *updateUserSpaceMembershipsBuilder) Remove(userMembershipRemove []PNUserMembershipRemove) *updateUserSpaceMembershipsBuilder {
	b.opts.UserMembershipRemove = userMembershipRemove

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *updateUserSpaceMembershipsBuilder) QueryParam(queryParam map[string]string) *updateUserSpaceMembershipsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the updateUserSpaceMemberships request.
func (b *updateUserSpaceMembershipsBuilder) Transport(tr http.RoundTripper) *updateUserSpaceMembershipsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the updateUserSpaceMemberships request.
func (b *updateUserSpaceMembershipsBuilder) Execute() (*PNUpdateUserSpaceMembershipsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyUpdateUserSpaceMembershipsResponse, status, err
	}

	return newPNUpdateUserSpaceMembershipsResponse(rawJSON, b.opts, status)
}

type updateUserSpaceMembershipsOpts struct {
	pubnub               *PubNub
	UserId               string
	Limit                int
	Include              []string
	Start                string
	End                  string
	Count                bool
	QueryParam           map[string]string
	UserMembershipRemove []PNUserMembershipRemove
	UserMembershipAdd    []PNUserMembershipInput
	UserMembershipUpdate []PNUserMembershipInput
	Transport            http.RoundTripper

	ctx Context
}

func (o *updateUserSpaceMembershipsOpts) config() Config {
	return *o.pubnub.Config
}

func (o *updateUserSpaceMembershipsOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *updateUserSpaceMembershipsOpts) context() Context {
	return o.ctx
}

func (o *updateUserSpaceMembershipsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *updateUserSpaceMembershipsOpts) buildPath() (string, error) {
	return fmt.Sprintf(updateUserSpaceMembershipsPath,
		o.pubnub.Config.SubscribeKey, o.UserId), nil
}

func (o *updateUserSpaceMembershipsOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}

	// if o.Auth != "" {
	// 	q.Set("auth", o.Auth)
	// }

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

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *updateUserSpaceMembershipsOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

type PNUserMembershipInputChangeSet struct {
	Add    []PNUserMembershipInput  `json:"add"`
	Update []PNUserMembershipInput  `json:"update"`
	Remove []PNUserMembershipRemove `json:"remove"`
}

func (o *updateUserSpaceMembershipsOpts) buildBody() ([]byte, error) {
	b := &PNUserMembershipInputChangeSet{
		Add:    o.UserMembershipAdd,
		Update: o.UserMembershipUpdate,
		Remove: o.UserMembershipRemove,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	fmt.Println(fmt.Sprintf("%v %s", b, string(jsonEncBytes)))
	return jsonEncBytes, nil
}

func (o *updateUserSpaceMembershipsOpts) httpMethod() string {
	return "PATCH"
}

func (o *updateUserSpaceMembershipsOpts) isAuthRequired() bool {
	return true
}

func (o *updateUserSpaceMembershipsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *updateUserSpaceMembershipsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *updateUserSpaceMembershipsOpts) operationType() OperationType {
	return PNUpdateMembersOperation
}

func (o *updateUserSpaceMembershipsOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

type PNUpdateUserSpaceMembershipsResponse struct {
	Status     int                `json:"status"`
	Data       []PNUserMembership `json:"data"`
	TotalCount int                `json:"totalCount"`
	Next       string             `json:"next"`
	Prev       string             `json:"prev"`
}

func newPNUpdateUserSpaceMembershipsResponse(jsonBytes []byte, o *updateUserSpaceMembershipsOpts,
	status StatusResponse) (*PNUpdateUserSpaceMembershipsResponse, StatusResponse, error) {

	resp := &PNUpdateUserSpaceMembershipsResponse{}

	fmt.Println(string(jsonBytes))

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		fmt.Println("error", err)
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyUpdateUserSpaceMembershipsResponse, status, e
	}

	return resp, status, nil
}
