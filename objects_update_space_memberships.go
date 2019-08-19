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

var emptyUpdateSpaceMembershipsResponse *PNUpdateSpaceMembershipsResponse

const updateSpaceMembershipsPath = "/v1/objects/%s/spaces/%s/users"

const updateSpaceMembershipsLimit = 100

type updateSpaceMembershipsBuilder struct {
	opts *updateSpaceMembershipsOpts
}

func newUpdateSpaceMembershipsBuilder(pubnub *PubNub) *updateSpaceMembershipsBuilder {
	builder := updateSpaceMembershipsBuilder{
		opts: &updateSpaceMembershipsOpts{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = spaceLimit

	return &builder
}

func newUpdateSpaceMembershipsBuilderWithContext(pubnub *PubNub,
	context Context) *updateSpaceMembershipsBuilder {
	builder := updateSpaceMembershipsBuilder{
		opts: &updateSpaceMembershipsOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Auth sets the Authorization key with permissions to perform the request.
// func (b *updateSpaceMembershipsBuilder) Auth(auth string) *updateSpaceMembershipsBuilder {
// 	//b.opts.Auth = auth

// 	return b
// }

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateSpaceMembershipsBuilder) Include(include []PNSpaceMembershipsIncude) *updateSpaceMembershipsBuilder {
	b.opts.Include = utils.EnumArrayToStringArray(fmt.Sprint(include))

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateSpaceMembershipsBuilder) SpaceId(id string) *updateSpaceMembershipsBuilder {
	b.opts.SpaceId = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateSpaceMembershipsBuilder) Limit(limit int) *updateSpaceMembershipsBuilder {
	b.opts.Limit = limit

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateSpaceMembershipsBuilder) Start(start string) *updateSpaceMembershipsBuilder {
	b.opts.Start = start

	return b
}

func (b *updateSpaceMembershipsBuilder) End(end string) *updateSpaceMembershipsBuilder {
	b.opts.End = end

	return b
}

func (b *updateSpaceMembershipsBuilder) Count(count bool) *updateSpaceMembershipsBuilder {
	b.opts.Count = count

	return b
}

func (b *updateSpaceMembershipsBuilder) Add(spaceMembershipInput []PNSpaceMembershipInput) *updateSpaceMembershipsBuilder {
	b.opts.SpaceMembershipAdd = spaceMembershipInput

	return b
}

func (b *updateSpaceMembershipsBuilder) Update(spaceMembershipInput []PNSpaceMembershipInput) *updateSpaceMembershipsBuilder {
	b.opts.SpaceMembershipUpdate = spaceMembershipInput

	return b
}

func (b *updateSpaceMembershipsBuilder) Remove(spaceMembershipRemove []PNSpaceMembershipRemove) *updateSpaceMembershipsBuilder {
	b.opts.SpaceMembershipRemove = spaceMembershipRemove

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *updateSpaceMembershipsBuilder) QueryParam(queryParam map[string]string) *updateSpaceMembershipsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the updateSpaceMemberships request.
func (b *updateSpaceMembershipsBuilder) Transport(tr http.RoundTripper) *updateSpaceMembershipsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the updateSpaceMemberships request.
func (b *updateSpaceMembershipsBuilder) Execute() (*PNUpdateSpaceMembershipsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyUpdateSpaceMembershipsResponse, status, err
	}

	return newPNUpdateSpaceMembershipsResponse(rawJSON, b.opts, status)
}

type updateSpaceMembershipsOpts struct {
	pubnub                *PubNub
	SpaceId               string
	Limit                 int
	Include               []string
	Start                 string
	End                   string
	Count                 bool
	QueryParam            map[string]string
	SpaceMembershipRemove []PNSpaceMembershipRemove
	SpaceMembershipAdd    []PNSpaceMembershipInput
	SpaceMembershipUpdate []PNSpaceMembershipInput
	Transport             http.RoundTripper

	ctx Context
}

func (o *updateSpaceMembershipsOpts) config() Config {
	return *o.pubnub.Config
}

func (o *updateSpaceMembershipsOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *updateSpaceMembershipsOpts) context() Context {
	return o.ctx
}

func (o *updateSpaceMembershipsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *updateSpaceMembershipsOpts) buildPath() (string, error) {
	return fmt.Sprintf(updateSpaceMembershipsPath,
		o.pubnub.Config.SubscribeKey, o.SpaceId), nil
}

func (o *updateSpaceMembershipsOpts) buildQuery() (*url.Values, error) {

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

func (o *updateSpaceMembershipsOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

type PNSpaceMembershipInputChangeSet struct {
	Add    []PNSpaceMembershipInput  `json:"add"`
	Update []PNSpaceMembershipInput  `json:"update"`
	Remove []PNSpaceMembershipRemove `json:"remove"`
}

func (o *updateSpaceMembershipsOpts) buildBody() ([]byte, error) {
	b := &PNSpaceMembershipInputChangeSet{
		Add:    o.SpaceMembershipAdd,
		Update: o.SpaceMembershipUpdate,
		Remove: o.SpaceMembershipRemove,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	fmt.Println(fmt.Sprintf("buildBody %v %s", b, string(jsonEncBytes)))
	return jsonEncBytes, nil

}

func (o *updateSpaceMembershipsOpts) httpMethod() string {
	return "PATCH"
}

func (o *updateSpaceMembershipsOpts) isAuthRequired() bool {
	return true
}

func (o *updateSpaceMembershipsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *updateSpaceMembershipsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *updateSpaceMembershipsOpts) operationType() OperationType {
	return PNUpdateSpaceMembershipsOperation
}

func (o *updateSpaceMembershipsOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

type PNUpdateSpaceMembershipsResponse struct {
	Status     int       `json:"status"`
	Data       []PNSpace `json:"data"`
	TotalCount int       `json:"totalCount"`
	Next       string    `json:"next"`
	Prev       string    `json:"prev"`
}

func newPNUpdateSpaceMembershipsResponse(jsonBytes []byte, o *updateSpaceMembershipsOpts,
	status StatusResponse) (*PNUpdateSpaceMembershipsResponse, StatusResponse, error) {

	resp := &PNUpdateSpaceMembershipsResponse{}

	fmt.Println(string(jsonBytes))

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		fmt.Println("error", err)
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyUpdateSpaceMembershipsResponse, status, e
	}

	return resp, status, nil
}
