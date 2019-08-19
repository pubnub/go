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

var emptyPNCreateUserResponse *PNCreateUserResponse

const createUserPath = "/v1/objects/%s/users"

type createUserBuilder struct {
	opts *createUserOpts
}

func newCreateUserBuilder(pubnub *PubNub) *createUserBuilder {
	builder := createUserBuilder{
		opts: &createUserOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newCreateUserBuilderWithContext(pubnub *PubNub,
	context Context) *createUserBuilder {
	builder := createUserBuilder{
		opts: &createUserOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

type createUserBody struct {
	Id         string                 `json:"id"`
	Name       string                 `json:"name"`
	ExternalId string                 `json:"externalId"`
	ProfileUrl string                 `json:"profileUrl"`
	Email      string                 `json:"email"`
	Custom     map[string]interface{} `json:"custom"`
}

// Auth sets the Authorization key with permissions to perform the request.
// func (b *createUserBuilder) Auth(auth string) *createUserBuilder {
// 	b.opts.Auth = auth

// 	return b
// }

// Auth sets the Authorization key with permissions to perform the request.
func (b *createUserBuilder) Include(include []PNUserSpaceInclude) *createUserBuilder {
	b.opts.Include = utils.EnumArrayToStringArray(fmt.Sprint(include))

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *createUserBuilder) Id(id string) *createUserBuilder {
	b.opts.Id = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *createUserBuilder) Name(name string) *createUserBuilder {
	b.opts.Name = name

	return b
}

func (b *createUserBuilder) ExternalId(externalId string) *createUserBuilder {
	b.opts.ExternalId = externalId

	return b
}

func (b *createUserBuilder) ProfileUrl(profileUrl string) *createUserBuilder {
	b.opts.ProfileUrl = profileUrl

	return b
}

func (b *createUserBuilder) Email(email string) *createUserBuilder {
	b.opts.Email = email

	return b
}

func (b *createUserBuilder) Custom(custom map[string]interface{}) *createUserBuilder {
	b.opts.Custom = custom

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *createUserBuilder) QueryParam(queryParam map[string]string) *createUserBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the createUser request.
func (b *createUserBuilder) Transport(tr http.RoundTripper) *createUserBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the createUser request.
func (b *createUserBuilder) Execute() (*PNCreateUserResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNCreateUserResponse, status, err
	}

	return newPNCreateUserResponse(rawJSON, b.opts, status)
}

type createUserOpts struct {
	pubnub *PubNub

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

func (o *createUserOpts) config() Config {
	return *o.pubnub.Config
}

func (o *createUserOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *createUserOpts) context() Context {
	return o.ctx
}

func (o *createUserOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *createUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(createUserPath,
		o.pubnub.Config.SubscribeKey), nil
}

func (o *createUserOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}

	// if o.Auth != "" {
	// 	q.Set("auth", o.Auth)
	// }

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *createUserOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *createUserOpts) buildBody() ([]byte, error) {
	b := &createUserBody{
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
	fmt.Println(fmt.Sprintf("%v %s", b, string(jsonEncBytes)))
	return jsonEncBytes, nil

}

func (o *createUserOpts) httpMethod() string {
	return "POST"
}

func (o *createUserOpts) isAuthRequired() bool {
	return true
}

func (o *createUserOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *createUserOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *createUserOpts) operationType() OperationType {
	return PNCreateUserOperation
}

func (o *createUserOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

type PNCreateUserResponse struct {
	Status int    `json:"status"`
	Data   PNUser `json:"data"`
}

func newPNCreateUserResponse(jsonBytes []byte, o *createUserOpts,
	status StatusResponse) (*PNCreateUserResponse, StatusResponse, error) {

	resp := &PNCreateUserResponse{}

	fmt.Println(string(jsonBytes))

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		fmt.Println("error", err)
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNCreateUserResponse, status, e
	}

	return resp, status, nil
}
