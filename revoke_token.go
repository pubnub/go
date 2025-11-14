package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/pubnub/go/v8/pnerr"
	"github.com/pubnub/go/v8/utils"
)

const revokeTokenPath = "/v3/pam/%s/grant/%s"

var emptyPNRevokeTokenResponse *PNRevokeTokenResponse

type revokeTokenBuilder struct {
	opts *revokeTokenOpts
}

func newRevokeTokenBuilder(pubnub *PubNub) *revokeTokenBuilder {
	return newRevokeTokenBuilderWithContext(pubnub, pubnub.ctx)
}

func newRevokeTokenOpts(pubnub *PubNub, ctx Context) *revokeTokenOpts {
	return &revokeTokenOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}
func newRevokeTokenBuilderWithContext(pubnub *PubNub, context Context) *revokeTokenBuilder {
	builder := revokeTokenBuilder{
		opts: newRevokeTokenOpts(pubnub, context)}
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

// GetLogParams returns the user-provided parameters for logging
func (o *revokeTokenOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{}
	if o.Token != "" {
		// Mask token for security, show only first 8 chars
		if len(o.Token) > 8 {
			params["Token"] = o.Token[:8] + "..."
		} else {
			params["Token"] = "***"
		}
	}
	return params
}

// Execute runs the Grant request.
func (b *revokeTokenBuilder) Execute() (*PNRevokeTokenResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNAccessManagerRevokeToken, b.opts.GetLogParams(), true)
	
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNRevokeTokenResponse, status, err
	}

	return newPNRevokeTokenResponse(rawJSON, b.opts, status)
}

type revokeTokenOpts struct {
	endpointOpts

	QueryParam map[string]string
	Token      string
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
	return fmt.Sprintf(revokeTokenPath, o.pubnub.Config.SubscribeKey, utils.URLEncode(o.Token)), nil
}

func (o *revokeTokenOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *revokeTokenOpts) httpMethod() string {
	return "DELETE"
}

func (o *revokeTokenOpts) operationType() OperationType {
	return PNAccessManagerRevokeToken
}

// PNRevokeTokenResponse is the struct returned when the Execute function of Grant Token is called.
type PNRevokeTokenResponse struct {
	Status int `json:"status"`
}

func newPNRevokeTokenResponse(jsonBytes []byte, o *revokeTokenOpts, status StatusResponse) (*PNRevokeTokenResponse, StatusResponse, error) {
	resp := &PNRevokeTokenResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNRevokeTokenResponse, status, e
	}

	return resp, status, nil
}
