package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
	//"reflect"

	"net/http"
	"net/url"
)

var emptyObjectAPICreateUserResp *ObjectAPICreateUserResponse

const objectAPICreateUserPath = "/v1/objects/%s/users"

type objectAPICreateUserBuilder struct {
	opts *objectAPICreateUserOpts
}

func newObjectAPICreateUserBuilder(pubnub *PubNub) *objectAPICreateUserBuilder {
	builder := objectAPICreateUserBuilder{
		opts: &objectAPICreateUserOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newObjectAPICreateUserBuilderWithContext(pubnub *PubNub,
	context Context) *objectAPICreateUserBuilder {
	builder := objectAPICreateUserBuilder{
		opts: &objectAPICreateUserOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

type createUserBody struct {
	id         string
	name       string
	externalId string
	profileUrl string
	email      string
	custom     map[string]interface{}
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *objectAPICreateUserBuilder) Auth(auth string) *objectAPICreateUserBuilder {
	b.opts.Auth = auth

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *objectAPICreateUserBuilder) Include(include []string) *objectAPICreateUserBuilder {
	b.opts.Include = include

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *objectAPICreateUserBuilder) Id(id string) *objectAPICreateUserBuilder {
	b.opts.Id = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *objectAPICreateUserBuilder) Name(name string) *objectAPICreateUserBuilder {
	b.opts.Name = name

	return b
}

func (b *objectAPICreateUserBuilder) ExternalId(externalId string) *objectAPICreateUserBuilder {
	b.opts.ExternalId = externalId

	return b
}

func (b *objectAPICreateUserBuilder) ProfileUrl(profileUrl string) *objectAPICreateUserBuilder {
	b.opts.ProfileUrl = profileUrl

	return b
}

func (b *objectAPICreateUserBuilder) Email(email string) *objectAPICreateUserBuilder {
	b.opts.Email = email

	return b
}

func (b *objectAPICreateUserBuilder) Custom(custom map[string]interface{}) *objectAPICreateUserBuilder {
	b.opts.Custom = custom

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *objectAPICreateUserBuilder) QueryParam(queryParam map[string]string) *objectAPICreateUserBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the objectAPICreateUser request.
func (b *objectAPICreateUserBuilder) Transport(tr http.RoundTripper) *objectAPICreateUserBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the objectAPICreateUser request.
func (b *objectAPICreateUserBuilder) Execute() (*ObjectAPICreateUserResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyObjectAPICreateUserResp, status, err
	}

	return newObjectAPICreateUserResponse(rawJSON, b.opts, status)
}

type objectAPICreateUserOpts struct {
	pubnub *PubNub

	Auth       string
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

func (o *objectAPICreateUserOpts) config() Config {
	return *o.pubnub.Config
}

func (o *objectAPICreateUserOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *objectAPICreateUserOpts) context() Context {
	return o.ctx
}

func (o *objectAPICreateUserOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *objectAPICreateUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(objectAPICreateUserPath,
		o.pubnub.Config.SubscribeKey), nil
}

func (o *objectAPICreateUserOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}

	if o.Auth != "" {
		q.Set("auth", o.Auth)
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *objectAPICreateUserOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *objectAPICreateUserOpts) buildBody() ([]byte, error) {
	b := &createUserBody{
		id:         o.Id,
		name:       o.Name,
		externalId: o.ExternalId,
		profileUrl: o.ProfileUrl,
		email:      o.Email,
		custom:     o.Custom,
	}
	jsonEncBytes, errEnc := json.Marshal(b)
	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

}

func (o *objectAPICreateUserOpts) httpMethod() string {
	return "POST"
}

func (o *objectAPICreateUserOpts) isAuthRequired() bool {
	return true
}

func (o *objectAPICreateUserOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *objectAPICreateUserOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *objectAPICreateUserOpts) operationType() OperationType {
	return PNCreateUserOperation
}

func (o *objectAPICreateUserOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// ObjectAPICreateUserResponse is the response to objectAPICreateUser request. It contains a map of type objectAPICreateUserResponseItem
type ObjectAPICreateUserResponse struct {
	Id         string                 `json:"id"`
	Name       string                 `json:"name"`
	ExternalId string                 `json:"externalId"`
	ProfileUrl string                 `json:"profileUrl"`
	Email      string                 `json:"email"`
	Custom     map[string]interface{} `json:"custom"`
	Created    string                 `json:"created"`
	Updated    string                 `json:"updated"`
	ETag       string                 `json:"eTag"`
}

type ObjectAPICreateUserWithData struct {
	Status string                      `json:"status"`
	Data   ObjectAPICreateUserResponse `json:"data"`
}

func newObjectAPICreateUserResponse(jsonBytes []byte, o *objectAPICreateUserOpts,
	status StatusResponse) (*ObjectAPICreateUserResponse, StatusResponse, error) {

	resp := &ObjectAPICreateUserWithData{}

	fmt.Println(string(jsonBytes))

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyObjectAPICreateUserResp, status, e
	}

	return &resp.Data, status, nil
}
