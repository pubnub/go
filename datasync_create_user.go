package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyCreateUserResponse *PNUserResponse

type createUserBuilder struct {
	opts *createUserOpts
}

func newCreateUserBuilder(pubnub *PubNub) *createUserBuilder {
	return newCreateUserBuilderWithContext(pubnub, pubnub.ctx)
}

func newCreateUserOpts(pubnub *PubNub, ctx Context) *createUserOpts {
	return &createUserOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newCreateUserBuilderWithContext(pubnub *PubNub, context Context) *createUserBuilder {
	return &createUserBuilder{opts: newCreateUserOpts(pubnub, context)}
}

// createUserBody is the request envelope sent on POST /users.
type createUserBody struct {
	Data createUserBodyData `json:"data"`
}

type createUserBodyData struct {
	ID                 string                 `json:"id,omitempty"`
	EntityClass        string                 `json:"entityClass,omitempty"`
	EntityClassVersion int                    `json:"entityClassVersion"`
	EntityClassLevel   PNEntityClassLevel     `json:"entityClassLevel,omitempty"`
	Status             string                 `json:"status,omitempty"`
	Payload            map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the optional user identifier. When omitted the server generates one.
func (b *createUserBuilder) ID(id string) *createUserBuilder {
	b.opts.ID = id
	return b
}

// EntityClass sets the optional entity class. Defaults to "User" on the server;
// when set it must be User or a subclass of User. Immutable after creation.
func (b *createUserBuilder) EntityClass(entityClass string) *createUserBuilder {
	b.opts.EntityClass = entityClass
	return b
}

// EntityClassVersion sets the required schema version of the entity class (>= 1).
func (b *createUserBuilder) EntityClassVersion(version int) *createUserBuilder {
	b.opts.EntityClassVersion = version
	return b
}

// EntityClassLevel sets the optional class scope (Global, Account, or SubKey) used to
// disambiguate classes with the same name defined at different levels.
func (b *createUserBuilder) EntityClassLevel(level PNEntityClassLevel) *createUserBuilder {
	b.opts.EntityClassLevel = level
	return b
}

// Status sets the optional user status (e.g. "active").
func (b *createUserBuilder) Status(status string) *createUserBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the optional custom profile fields.
func (b *createUserBuilder) Payload(payload map[string]interface{}) *createUserBuilder {
	b.opts.Payload = payload
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

// GetLogParams returns the user-provided parameters for logging
func (o *createUserOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"EntityClassVersion": o.EntityClassVersion,
	}
	if o.ID != "" {
		params["ID"] = o.ID
	}
	if o.EntityClass != "" {
		params["EntityClass"] = o.EntityClass
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

// Execute runs the createUser request.
func (b *createUserBuilder) Execute() (*PNUserResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNCreateDataSyncUserOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyCreateUserResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNCreateDataSyncUserOperation, b.opts.pubnub)
}

type createUserOpts struct {
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

func (o *createUserOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.EntityClassVersion < 1 {
		return newValidationError(o, StrInvalidEntityClassVersion)
	}
	return nil
}

func (o *createUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(usersPath, o.pubnub.Config.SubscribeKey), nil
}

func (o *createUserOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *createUserOpts) buildBody() ([]byte, error) {
	b := &createUserBody{
		Data: createUserBodyData{
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
		o.pubnub.loggerManager.LogError(errEnc, "CreateUserSerializationFailed", PNCreateDataSyncUserOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *createUserOpts) buildHeaders() (map[string]string, error) {
	return map[string]string{
		"Content-Type": userContentType,
	}, nil
}

func (o *createUserOpts) httpMethod() string {
	return "POST"
}

func (o *createUserOpts) operationType() OperationType {
	return PNCreateDataSyncUserOperation
}
