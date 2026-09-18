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

var emptyCreateRelationshipResponse *PNRelationshipResponse

type createRelationshipBuilder struct {
	opts *createRelationshipOpts
}

func newCreateRelationshipBuilder(pubnub *PubNub) *createRelationshipBuilder {
	return newCreateRelationshipBuilderWithContext(pubnub, pubnub.ctx)
}

func newCreateRelationshipOpts(pubnub *PubNub, ctx Context) *createRelationshipOpts {
	return &createRelationshipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newCreateRelationshipBuilderWithContext(pubnub *PubNub, context Context) *createRelationshipBuilder {
	return &createRelationshipBuilder{opts: newCreateRelationshipOpts(pubnub, context)}
}

// createRelationshipBody is the request envelope sent on POST /relationships.
type createRelationshipBody struct {
	Data createRelationshipBodyData `json:"data"`
}

type createRelationshipBodyData struct {
	ID                       string                 `json:"id,omitempty"`
	EntityAID                string                 `json:"entityAId"`
	EntityBID                string                 `json:"entityBId"`
	RelationshipClass        string                 `json:"relationshipClass"`
	RelationshipClassVersion int                    `json:"relationshipClassVersion"`
	Status                   string                 `json:"status,omitempty"`
	Payload                  map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the optional relationship identifier. When omitted the server generates one.
func (b *createRelationshipBuilder) ID(id string) *createRelationshipBuilder {
	b.opts.ID = id
	return b
}

// EntityAID sets the required first linked entity identifier. Immutable after creation.
func (b *createRelationshipBuilder) EntityAID(entityAID string) *createRelationshipBuilder {
	b.opts.EntityAID = entityAID
	return b
}

// EntityBID sets the required second linked entity identifier. Immutable after creation.
func (b *createRelationshipBuilder) EntityBID(entityBID string) *createRelationshipBuilder {
	b.opts.EntityBID = entityBID
	return b
}

// RelationshipClass sets the required relationship class identifier. Immutable after creation.
func (b *createRelationshipBuilder) RelationshipClass(relationshipClass string) *createRelationshipBuilder {
	b.opts.RelationshipClass = relationshipClass
	return b
}

// RelationshipClassVersion sets the required schema version of the relationship class (>= 1).
func (b *createRelationshipBuilder) RelationshipClassVersion(version int) *createRelationshipBuilder {
	b.opts.RelationshipClassVersion = version
	return b
}

// Status sets the optional relationship status (e.g. "active").
func (b *createRelationshipBuilder) Status(status string) *createRelationshipBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the optional user-defined custom properties.
func (b *createRelationshipBuilder) Payload(payload map[string]interface{}) *createRelationshipBuilder {
	b.opts.Payload = payload
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *createRelationshipBuilder) QueryParam(queryParam map[string]string) *createRelationshipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the createRelationship request.
func (b *createRelationshipBuilder) Transport(tr http.RoundTripper) *createRelationshipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *createRelationshipOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"EntityAID":                o.EntityAID,
		"EntityBID":                o.EntityBID,
		"RelationshipClass":        o.RelationshipClass,
		"RelationshipClassVersion": o.RelationshipClassVersion,
	}
	if o.ID != "" {
		params["ID"] = o.ID
	}
	if o.Status != "" {
		params["Status"] = o.Status
	}
	if o.Payload != nil {
		params["Payload"] = fmt.Sprintf("%v", o.Payload)
	}
	return params
}

// Execute runs the createRelationship request.
func (b *createRelationshipBuilder) Execute() (*PNRelationshipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNCreateRelationshipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyCreateRelationshipResponse, status, err
	}

	return newPNRelationshipResponse(rawJSON, status, PNCreateRelationshipOperation, b.opts.pubnub)
}

type createRelationshipOpts struct {
	endpointOpts

	ID                       string
	EntityAID                string
	EntityBID                string
	RelationshipClass        string
	RelationshipClassVersion int
	Status                   string
	Payload                  map[string]interface{}
	QueryParam               map[string]string

	Transport http.RoundTripper
}

func (o *createRelationshipOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.EntityAID == "" {
		return newValidationError(o, StrMissingEntityAID)
	}
	if o.EntityBID == "" {
		return newValidationError(o, StrMissingEntityBID)
	}
	if o.RelationshipClass == "" {
		return newValidationError(o, StrMissingRelationshipClass)
	}
	if o.RelationshipClassVersion < 1 {
		return newValidationError(o, StrInvalidRelationshipClassVersion)
	}
	return nil
}

func (o *createRelationshipOpts) buildPath() (string, error) {
	return fmt.Sprintf(relationshipsPath, o.pubnub.Config.SubscribeKey), nil
}

func (o *createRelationshipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *createRelationshipOpts) buildBody() ([]byte, error) {
	b := &createRelationshipBody{
		Data: createRelationshipBodyData{
			ID:                       o.ID,
			EntityAID:                o.EntityAID,
			EntityBID:                o.EntityBID,
			RelationshipClass:        o.RelationshipClass,
			RelationshipClassVersion: o.RelationshipClassVersion,
			Status:                   o.Status,
			Payload:                  o.Payload,
		},
	}

	jsonEncBytes, errEnc := json.Marshal(b)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "CreateRelationshipSerializationFailed", PNCreateRelationshipOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *createRelationshipOpts) buildHeaders() (map[string]string, error) {
	return map[string]string{
		"Content-Type": relationshipContentType,
	}, nil
}

func (o *createRelationshipOpts) httpMethod() string {
	return "POST"
}

func (o *createRelationshipOpts) operationType() OperationType {
	return PNCreateRelationshipOperation
}

// newPNRelationshipResponse unmarshals a single-relationship envelope shared by
// the create, get, update and patch operations.
func newPNRelationshipResponse(jsonBytes []byte, status StatusResponse,
	operation OperationType, pn *PubNub) (*PNRelationshipResponse, StatusResponse, error) {

	resp := &PNRelationshipResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		pn.loggerManager.LogError(e, "RelationshipResponseParsingFailed", operation, true)
		return nil, status, e
	}

	return resp, status, nil
}
