package pubnub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v6/pnerr"
)

const revokeTokenPath = "/v3/pam/%s/grant/%s"

var emptyPNRevokeTokenResponse *PNRevokeTokenResponse

type revokeTokenBuilder struct {
	opts *revokeTokenOpts
}

func newRevokeTokenBuilder(pubnub *PubNub) *revokeTokenBuilder {
	builder := revokeTokenBuilder{
		opts: &revokeTokenOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newRevokeTokenBuilderWithContext(pubnub *PubNub, context Context) *revokeTokenBuilder {
	builder := revokeTokenBuilder{
		opts: &revokeTokenOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *revokeTokenBuilder) Token(token string) *revokeTokenBuilder {
	b.opts.Token = token

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *revokeTokenBuilder) QueryParam(queryParam map[string]string) *revokeTokenBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Grant request.
func (b *revokeTokenBuilder) Execute() (*PNRevokeTokenResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNRevokeTokenResponse, status, err
	}

	return newPNRevokeTokenResponse(rawJSON, b.opts, status)
}

type revokeTokenOpts struct {
	pubnub *PubNub
	ctx    Context

	QueryParam map[string]string
	Token      string
}

func (o *revokeTokenOpts) config() Config {
	return *o.pubnub.Config
}

func (o *revokeTokenOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *revokeTokenOpts) context() Context {
	return o.ctx
}

func (o *revokeTokenOpts) validate() error {
	if o.config().PublishKey == "" {
		return newValidationError(o, StrMissingPubKey)
	}

	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.config().SecretKey == "" {
		return newValidationError(o, StrMissingSecretKey)
	}

	if o.Token == "" {
		return newValidationError(o, StrMissingToken)
	}
	return nil
}

func (o *revokeTokenOpts) buildPath() (string, error) {
	return fmt.Sprintf(revokeTokenPath, o.pubnub.Config.SubscribeKey, o.Token), nil
}

func (o *revokeTokenOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *revokeTokenOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *revokeTokenOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *revokeTokenOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *revokeTokenOpts) httpMethod() string {
	return "DELETE"
}

func (o *revokeTokenOpts) isAuthRequired() bool {
	return true
}

func (o *revokeTokenOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *revokeTokenOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *revokeTokenOpts) operationType() OperationType {
	return PNAccessManagerRevokeToken
}

func (o *revokeTokenOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

func (o *revokeTokenOpts) tokenManager() *TokenManager {
	return o.pubnub.tokenManager
}

// PNRevokeTokenResponse is the struct returned when the Execute function of Grant Token is called.
type PNRevokeTokenResponse struct {
	status int `json:"status"`
}

func newPNRevokeTokenResponse(jsonBytes []byte, o *revokeTokenOpts, status StatusResponse) (*PNRevokeTokenResponse, StatusResponse, error) {
	resp := &PNRevokeTokenResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNRevokeTokenResponse, status, e
	}

	return resp, status, nil
}
