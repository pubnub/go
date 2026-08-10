package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyCreateChannelResponse *PNDataSyncChannelResponse

type createChannelBuilder struct {
	opts *createChannelOpts
}

func newCreateChannelBuilder(pubnub *PubNub) *createChannelBuilder {
	return newCreateChannelBuilderWithContext(pubnub, pubnub.ctx)
}

func newCreateChannelOpts(pubnub *PubNub, ctx Context) *createChannelOpts {
	return &createChannelOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newCreateChannelBuilderWithContext(pubnub *PubNub, context Context) *createChannelBuilder {
	return &createChannelBuilder{opts: newCreateChannelOpts(pubnub, context)}
}

// createChannelBody is the request envelope sent on POST /channels.
type createChannelBody struct {
	Data createChannelBodyData `json:"data"`
}

type createChannelBodyData struct {
	ID                 string                 `json:"id,omitempty"`
	EntityClass        string                 `json:"entityClass,omitempty"`
	EntityClassVersion int                    `json:"entityClassVersion"`
	EntityClassLevel   string                 `json:"entityClassLevel,omitempty"`
	Status             string                 `json:"status,omitempty"`
	Payload            map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the optional channel identifier. When omitted the server generates one.
func (b *createChannelBuilder) ID(id string) *createChannelBuilder {
	b.opts.ID = id
	return b
}

// EntityClass sets the optional entity class. Defaults to "Channel" on the server;
// when set it must be Channel or a subclass of Channel. Immutable after creation.
func (b *createChannelBuilder) EntityClass(entityClass string) *createChannelBuilder {
	b.opts.EntityClass = entityClass
	return b
}

// EntityClassVersion sets the required schema version of the entity class (>= 1).
func (b *createChannelBuilder) EntityClassVersion(version int) *createChannelBuilder {
	b.opts.EntityClassVersion = version
	return b
}

// EntityClassLevel sets the optional class scope (Global or SubKey) used to
// disambiguate classes with the same name defined at different levels.
func (b *createChannelBuilder) EntityClassLevel(level string) *createChannelBuilder {
	b.opts.EntityClassLevel = level
	return b
}

// Status sets the optional channel status (e.g. "active").
func (b *createChannelBuilder) Status(status string) *createChannelBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the optional custom channel fields.
func (b *createChannelBuilder) Payload(payload map[string]interface{}) *createChannelBuilder {
	b.opts.Payload = payload
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *createChannelBuilder) QueryParam(queryParam map[string]string) *createChannelBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the createChannel request.
func (b *createChannelBuilder) Transport(tr http.RoundTripper) *createChannelBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *createChannelOpts) GetLogParams() map[string]interface{} {
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

// Execute runs the createChannel request.
func (b *createChannelBuilder) Execute() (*PNDataSyncChannelResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNCreateDataSyncChannelOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyCreateChannelResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNCreateDataSyncChannelOperation, b.opts.pubnub)
}

type createChannelOpts struct {
	endpointOpts

	ID                 string
	EntityClass        string
	EntityClassVersion int
	EntityClassLevel   string
	Status             string
	Payload            map[string]interface{}
	QueryParam         map[string]string

	Transport http.RoundTripper
}

func (o *createChannelOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.EntityClassVersion < 1 {
		return newValidationError(o, StrInvalidEntityClassVersion)
	}
	return nil
}

func (o *createChannelOpts) buildPath() (string, error) {
	return fmt.Sprintf(channelsPath, o.pubnub.Config.SubscribeKey), nil
}

func (o *createChannelOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *createChannelOpts) buildBody() ([]byte, error) {
	b := &createChannelBody{
		Data: createChannelBodyData{
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
		o.pubnub.loggerManager.LogError(errEnc, "CreateChannelSerializationFailed", PNCreateDataSyncChannelOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *createChannelOpts) buildHeaders() (map[string]string, error) {
	return map[string]string{
		"Content-Type": channelContentType,
	}, nil
}

func (o *createChannelOpts) httpMethod() string {
	return "POST"
}

func (o *createChannelOpts) operationType() OperationType {
	return PNCreateDataSyncChannelOperation
}
