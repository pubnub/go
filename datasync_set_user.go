package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptySetUserResponse *PNUserResponse

type setUserBuilder struct {
	opts *setUserOpts
}

func newSetUserBuilder(pubnub *PubNub) *setUserBuilder {
	return newSetUserBuilderWithContext(pubnub, pubnub.ctx)
}

func newSetUserOpts(pubnub *PubNub, ctx Context) *setUserOpts {
	return &setUserOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newSetUserBuilderWithContext(pubnub *PubNub, context Context) *setUserBuilder {
	return &setUserBuilder{opts: newSetUserOpts(pubnub, context)}
}

// setUserBody is the request envelope sent on PUT /users/{id}.
// entityClass is intentionally omitted because it is immutable.
type setUserBody struct {
	Data setUserBodyData `json:"data"`
}

type setUserBodyData struct {
	EntityClassVersion int                    `json:"entityClassVersion"`
	Status             string                 `json:"status,omitempty"`
	Payload            map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the required identifier of the user to replace.
func (b *setUserBuilder) ID(id string) *setUserBuilder {
	b.opts.ID = id
	return b
}

// EntityClassVersion sets the required schema version of the entity class (>= 1).
func (b *setUserBuilder) EntityClassVersion(version int) *setUserBuilder {
	b.opts.EntityClassVersion = version
	return b
}

// Status sets the optional user status (e.g. "active").
func (b *setUserBuilder) Status(status string) *setUserBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the optional custom profile fields. Under the default PAM
// projection this replaces the entire payload; omitted fields are removed.
func (b *setUserBuilder) Payload(payload map[string]interface{}) *setUserBuilder {
	b.opts.Payload = payload
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *setUserBuilder) IfMatchETag(eTag string) *setUserBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *setUserBuilder) QueryParam(queryParam map[string]string) *setUserBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the setUser request.
func (b *setUserBuilder) Transport(tr http.RoundTripper) *setUserBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *setUserOpts) GetLogParams() map[string]interface{} {
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

// Execute runs the setUser request.
func (b *setUserBuilder) Execute() (*PNUserResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNSetDataSyncUserOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySetUserResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNSetDataSyncUserOperation, b.opts.pubnub)
}

type setUserOpts struct {
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

func (o *setUserOpts) validate() error {
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

func (o *setUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(usersIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *setUserOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *setUserOpts) buildBody() ([]byte, error) {
	b := &setUserBody{
		Data: setUserBodyData{
			EntityClassVersion: o.EntityClassVersion,
			Status:             o.Status,
			Payload:            o.Payload,
		},
	}

	jsonEncBytes, errEnc := json.Marshal(b)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "SetUserSerializationFailed", PNSetDataSyncUserOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *setUserOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": userContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *setUserOpts) httpMethod() string {
	return "PUT"
}

func (o *setUserOpts) operationType() OperationType {
	return PNSetDataSyncUserOperation
}
