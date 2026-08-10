package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyUpdateUserResponse *PNUserResponse

type updateUserBuilder struct {
	opts *updateUserOpts
}

func newUpdateUserBuilder(pubnub *PubNub) *updateUserBuilder {
	return newUpdateUserBuilderWithContext(pubnub, pubnub.ctx)
}

func newUpdateUserOpts(pubnub *PubNub, ctx Context) *updateUserOpts {
	return &updateUserOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newUpdateUserBuilderWithContext(pubnub *PubNub, context Context) *updateUserBuilder {
	return &updateUserBuilder{opts: newUpdateUserOpts(pubnub, context)}
}

// updateUserBody is the request envelope sent on PUT /users/{id}.
// entityClass is intentionally omitted because it is immutable.
type updateUserBody struct {
	Data updateUserBodyData `json:"data"`
}

type updateUserBodyData struct {
	EntityClassVersion int                    `json:"entityClassVersion"`
	Status             string                 `json:"status,omitempty"`
	Payload            map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the required identifier of the user to replace.
func (b *updateUserBuilder) ID(id string) *updateUserBuilder {
	b.opts.ID = id
	return b
}

// EntityClassVersion sets the required schema version of the entity class (>= 1).
func (b *updateUserBuilder) EntityClassVersion(version int) *updateUserBuilder {
	b.opts.EntityClassVersion = version
	return b
}

// Status sets the optional user status (e.g. "active").
func (b *updateUserBuilder) Status(status string) *updateUserBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the optional custom profile fields. Under the default PAM
// projection this replaces the entire payload; omitted fields are removed.
func (b *updateUserBuilder) Payload(payload map[string]interface{}) *updateUserBuilder {
	b.opts.Payload = payload
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *updateUserBuilder) IfMatchETag(eTag string) *updateUserBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
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

// GetLogParams returns the user-provided parameters for logging
func (o *updateUserOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"ID":                 o.ID,
		"EntityClassVersion": o.EntityClassVersion,
	}
	if o.Status != "" {
		params["Status"] = o.Status
	}
	if o.Payload != nil {
		params["Payload"] = fmt.Sprintf("%v", o.Payload)
	}
	return params
}

// Execute runs the updateUser request.
func (b *updateUserBuilder) Execute() (*PNUserResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNUpdateDataSyncUserOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyUpdateUserResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNUpdateDataSyncUserOperation, b.opts.pubnub)
}

type updateUserOpts struct {
	endpointOpts

	ID                 string
	EntityClassVersion int
	Status             string
	Payload            map[string]interface{}
	IfMatchETag        string
	setIfMatchETag     bool
	QueryParam         map[string]string

	Transport http.RoundTripper
}

func (o *updateUserOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingUserID)
	}
	if o.EntityClassVersion < 1 {
		return newValidationError(o, StrInvalidEntityClassVersion)
	}
	return nil
}

func (o *updateUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(usersIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *updateUserOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *updateUserOpts) buildBody() ([]byte, error) {
	b := &updateUserBody{
		Data: updateUserBodyData{
			EntityClassVersion: o.EntityClassVersion,
			Status:             o.Status,
			Payload:            o.Payload,
		},
	}

	jsonEncBytes, errEnc := json.Marshal(b)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "UpdateUserSerializationFailed", PNUpdateDataSyncUserOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *updateUserOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": userContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *updateUserOpts) httpMethod() string {
	return "PUT"
}

func (o *updateUserOpts) operationType() OperationType {
	return PNUpdateDataSyncUserOperation
}
