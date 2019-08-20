package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
)

var emptyPNDeleteUserResponse *PNDeleteUserResponse

const deleteUserPath = "/v1/objects/%s/users/%s"

type deleteUserBuilder struct {
	opts *deleteUserOpts
}

func newDeleteUserBuilder(pubnub *PubNub) *deleteUserBuilder {
	builder := deleteUserBuilder{
		opts: &deleteUserOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newDeleteUserBuilderWithContext(pubnub *PubNub,
	context Context) *deleteUserBuilder {
	builder := deleteUserBuilder{
		opts: &deleteUserOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Auth sets the Authorization key with permissions to perform the request.
// func (b *deleteUserBuilder) Auth(auth string) *deleteUserBuilder {
// 	b.opts.Auth = auth

// 	return b
// }

func (b *deleteUserBuilder) Id(id string) *deleteUserBuilder {
	b.opts.Id = id

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *deleteUserBuilder) QueryParam(queryParam map[string]string) *deleteUserBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the deleteUser request.
func (b *deleteUserBuilder) Transport(tr http.RoundTripper) *deleteUserBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the deleteUser request.
func (b *deleteUserBuilder) Execute() (*PNDeleteUserResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNDeleteUserResponse, status, err
	}

	return newPNDeleteUserResponse(rawJSON, b.opts, status)
}

type deleteUserOpts struct {
	pubnub     *PubNub
	Id         string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *deleteUserOpts) config() Config {
	return *o.pubnub.Config
}

func (o *deleteUserOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *deleteUserOpts) context() Context {
	return o.ctx
}

func (o *deleteUserOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *deleteUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(deleteUserPath,
		o.pubnub.Config.SubscribeKey, o.Id), nil
}

func (o *deleteUserOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	// if o.Auth != "" {
	// 	q.Set("auth", o.Auth)
	// }

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *deleteUserOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *deleteUserOpts) buildBody() ([]byte, error) {
	return []byte{}, nil

}

func (o *deleteUserOpts) httpMethod() string {
	return "DELETE"
}

func (o *deleteUserOpts) isAuthRequired() bool {
	return true
}

func (o *deleteUserOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *deleteUserOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *deleteUserOpts) operationType() OperationType {
	return PNDeleteUserOperation
}

func (o *deleteUserOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

type PNDeleteUserResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func newPNDeleteUserResponse(jsonBytes []byte, o *deleteUserOpts,
	status StatusResponse) (*PNDeleteUserResponse, StatusResponse, error) {

	resp := &PNDeleteUserResponse{}

	fmt.Println(string(jsonBytes))

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		fmt.Println("error", err)
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNDeleteUserResponse, status, e
	}

	return resp, status, nil
}
