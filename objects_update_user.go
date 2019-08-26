package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
)

var emptyPNUpdateUserResponse *PNUpdateUserResponse

const updateUserPath = "/v1/objects/%s/users/%s"

type updateUserBuilder struct {
	opts *updateUserOpts
}

func newUpdateUserBuilder(pubnub *PubNub) *updateUserBuilder {
	builder := updateUserBuilder{
		opts: &updateUserOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newUpdateUserBuilderWithContext(pubnub *PubNub,
	context Context) *updateUserBuilder {
	builder := updateUserBuilder{
		opts: &updateUserOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

type UpdateUserBody struct {
	Id         string                 `json:"id"`
	Name       string                 `json:"name"`
	ExternalId string                 `json:"externalId"`
	ProfileUrl string                 `json:"profileUrl"`
	Email      string                 `json:"email"`
	Custom     map[string]interface{} `json:"custom"`
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateUserBuilder) Include(include []PNUserSpaceInclude) *updateUserBuilder {
	b.opts.Include = utils.EnumArrayToStringArray(fmt.Sprint(include))

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateUserBuilder) Id(id string) *updateUserBuilder {
	b.opts.Id = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateUserBuilder) Name(name string) *updateUserBuilder {
	b.opts.Name = name

	return b
}

func (b *updateUserBuilder) ExternalId(externalId string) *updateUserBuilder {
	b.opts.ExternalId = externalId

	return b
}

func (b *updateUserBuilder) ProfileUrl(profileUrl string) *updateUserBuilder {
	b.opts.ProfileUrl = profileUrl

	return b
}

func (b *updateUserBuilder) Email(email string) *updateUserBuilder {
	b.opts.Email = email

	return b
}

func (b *updateUserBuilder) Custom(custom map[string]interface{}) *updateUserBuilder {
	b.opts.Custom = custom

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *updateUserBuilder) QueryParam(queryParam map[string]string) *updateUserBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the updateUser request.
func (b *updateUserBuilder) Transport(tr http.RoundTripper) *updateUserBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the updateUser request.
func (b *updateUserBuilder) Execute() (*PNUpdateUserResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNUpdateUserResponse, status, err
	}

	return newPNUpdateUserResponse(rawJSON, b.opts, status)
}

type updateUserOpts struct {
	pubnub     *PubNub
	Include    []string
	Id         string
	Name       string
	ExternalId string
	ProfileUrl string
	Email      string
	Custom     map[string]interface{}
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *updateUserOpts) config() Config {
	return *o.pubnub.Config
}

func (o *updateUserOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *updateUserOpts) context() Context {
	return o.ctx
}

func (o *updateUserOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *updateUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(updateUserPath,
		o.pubnub.Config.SubscribeKey, o.Id), nil
}

func (o *updateUserOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *updateUserOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *updateUserOpts) buildBody() ([]byte, error) {
	b := &UpdateUserBody{
		Id:         o.Id,
		Name:       o.Name,
		ExternalId: o.ExternalId,
		ProfileUrl: o.ProfileUrl,
		Email:      o.Email,
		Custom:     o.Custom,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

}

func (o *updateUserOpts) httpMethod() string {
	return "PATCH"
}

func (o *updateUserOpts) isAuthRequired() bool {
	return true
}

func (o *updateUserOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *updateUserOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *updateUserOpts) operationType() OperationType {
	return PNUpdateUserOperation
}

func (o *updateUserOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

type PNUpdateUserResponse struct {
	Status int    `json:"status"`
	Data   PNUser `json:"data"`
}

func newPNUpdateUserResponse(jsonBytes []byte, o *updateUserOpts,
	status StatusResponse) (*PNUpdateUserResponse, StatusResponse, error) {

	resp := &PNUpdateUserResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNUpdateUserResponse, status, e
	}

	return resp, status, nil
}
