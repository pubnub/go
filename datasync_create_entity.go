package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v9/pnerr"
)

var emptyCreateEntityResponse *PNEntityResponse

type createEntityBuilder struct {
	opts *createEntityOpts
}

func newCreateEntityBuilder(pubnub *PubNub) *createEntityBuilder {
	return newCreateEntityBuilderWithContext(pubnub, pubnub.ctx)
}

func newCreateEntityOpts(pubnub *PubNub, ctx Context) *createEntityOpts {
	return &createEntityOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newCreateEntityBuilderWithContext(pubnub *PubNub, context Context) *createEntityBuilder {
	return &createEntityBuilder{opts: newCreateEntityOpts(pubnub, context)}
}

// createEntityBody is the request envelope sent on POST /entities.
type createEntityBody struct {
	Data createEntityBodyData `json:"data"`
}

type createEntityBodyData struct {
	ID                 string                 `json:"id,omitempty"`
	EntityClass        string                 `json:"entityClass"`
	EntityClassVersion int                    `json:"entityClassVersion"`
	EntityClassLevel   PNEntityClassLevel     `json:"entityClassLevel,omitempty"`
	Status             string                 `json:"status,omitempty"`
	Payload            map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the optional entity identifier. When omitted the server generates one.
func (b *createEntityBuilder) ID(id string) *createEntityBuilder {
	b.opts.ID = id
	return b
}

// EntityClass sets the required entity class identifier. Immutable after creation.
func (b *createEntityBuilder) EntityClass(entityClass string) *createEntityBuilder {
	b.opts.EntityClass = entityClass
	return b
}

// EntityClassVersion sets the required schema version of the entity class (>= 1).
func (b *createEntityBuilder) EntityClassVersion(version int) *createEntityBuilder {
	b.opts.EntityClassVersion = version
	return b
}

// EntityClassLevel sets the optional class scope (Global, Account, or SubKey) used to
// disambiguate classes with the same name defined at different levels.
func (b *createEntityBuilder) EntityClassLevel(level PNEntityClassLevel) *createEntityBuilder {
	b.opts.EntityClassLevel = level
	return b
}

// Status sets the optional entity status (e.g. "active").
func (b *createEntityBuilder) Status(status string) *createEntityBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the optional user-defined custom properties.
func (b *createEntityBuilder) Payload(payload map[string]interface{}) *createEntityBuilder {
	b.opts.Payload = payload
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *createEntityBuilder) QueryParam(queryParam map[string]string) *createEntityBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the createEntity request.
func (b *createEntityBuilder) Transport(tr http.RoundTripper) *createEntityBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *createEntityOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"EntityClass":        o.EntityClass,
		"EntityClassVersion": o.EntityClassVersion,
	}
	if o.ID != "" {
		params["ID"] = o.ID
	}
	if o.EntityClassLevel != "" {
		params["EntityClassLevel"] = o.EntityClassLevel
	}
	if o.Status != "" {
		params["Status"] = o.Status
	}
	if o.Payload != nil {
		params["Payload"] = fmt.Sprintf("%v", o.Payload)
	}
	return params
}

// Execute runs the createEntity request.
func (b *createEntityBuilder) Execute() (*PNEntityResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNCreateEntityOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyCreateEntityResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNCreateEntityOperation, b.opts.pubnub)
}

type createEntityOpts struct {
	endpointOpts

	ID                 string
	EntityClass        string
	EntityClassVersion int
	EntityClassLevel   PNEntityClassLevel
	Status             string
	Payload            map[string]interface{}
	QueryParam         map[string]string

	Transport http.RoundTripper
}

func (o *createEntityOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.EntityClass == "" {
		return newValidationError(o, StrMissingEntityClass)
	}
	if o.EntityClassVersion < 1 {
		return newValidationError(o, StrInvalidEntityClassVersion)
	}
	return nil
}

func (o *createEntityOpts) buildPath() (string, error) {
	return fmt.Sprintf(entitiesPath, o.pubnub.Config.SubscribeKey), nil
}

func (o *createEntityOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *createEntityOpts) buildBody() ([]byte, error) {
	b := &createEntityBody{
		Data: createEntityBodyData{
			ID:                 o.ID,
			EntityClass:        o.EntityClass,
			EntityClassVersion: o.EntityClassVersion,
			EntityClassLevel:   o.EntityClassLevel,
			Status:             o.Status,
			Payload:            o.Payload,
		},
	}

	jsonEncBytes, errEnc := json.Marshal(b)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "CreateEntitySerializationFailed", PNCreateEntityOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *createEntityOpts) buildHeaders() (map[string]string, error) {
	return map[string]string{
		"Content-Type": entityContentType,
	}, nil
}

func (o *createEntityOpts) httpMethod() string {
	return "POST"
}

func (o *createEntityOpts) operationType() OperationType {
	return PNCreateEntityOperation
}

// newPNEntityResponse unmarshals a single-entity envelope shared by the create,
// get, update and patch operations.
func newPNEntityResponse(jsonBytes []byte, status StatusResponse,
	operation OperationType, pn *PubNub) (*PNEntityResponse, StatusResponse, error) {

	resp := &PNEntityResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		pn.loggerManager.LogError(e, "EntityResponseParsingFailed", operation, true)
		return nil, status, e
	}

	return resp, status, nil
}
