package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pubnub/go/v9/pnerr"
	"github.com/pubnub/go/v9/utils"
)

var emptyGetChannelsResponse *PNDataSyncChannelsResponse

type getChannelsBuilder struct {
	opts *getChannelsOpts
}

func newGetChannelsBuilder(pubnub *PubNub) *getChannelsBuilder {
	return newGetChannelsBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetChannelsOpts(pubnub *PubNub, ctx Context) *getChannelsOpts {
	return &getChannelsOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newGetChannelsBuilderWithContext(pubnub *PubNub, context Context) *getChannelsBuilder {
	builder := &getChannelsBuilder{opts: newGetChannelsOpts(pubnub, context)}
	builder.opts.Limit = entitiesDefaultLimit
	return builder
}

// EntityClass sets the optional entity class filter. Must be Channel or a
// subclass of Channel when provided.
func (b *getChannelsBuilder) EntityClass(entityClass string) *getChannelsBuilder {
	b.opts.EntityClass = entityClass
	return b
}

// EntityClassVersion sets the optional schema version to read against. When
// unset, the server uses the latest version.
func (b *getChannelsBuilder) EntityClassVersion(version int) *getChannelsBuilder {
	b.opts.EntityClassVersion = version
	b.opts.setEntityClassVersion = true
	return b
}

// EntityClassLevel sets the optional class scope (Global or SubKey) used to
// disambiguate classes with the same name defined at different levels.
func (b *getChannelsBuilder) EntityClassLevel(level string) *getChannelsBuilder {
	b.opts.EntityClassLevel = level
	return b
}

// Cursor sets the pagination token returned from a previous request.
func (b *getChannelsBuilder) Cursor(cursor string) *getChannelsBuilder {
	b.opts.Cursor = cursor
	return b
}

// Limit sets the maximum number of items to return per page (1-100, default 20).
func (b *getChannelsBuilder) Limit(limit int) *getChannelsBuilder {
	b.opts.Limit = limit
	return b
}

// Filter sets the AppContext Query Language filter expression.
func (b *getChannelsBuilder) Filter(filter string) *getChannelsBuilder {
	b.opts.Filter = filter
	return b
}

// FilterAdvanced sets the extended filter expression supporting logical
// operators and nested conditions.
func (b *getChannelsBuilder) FilterAdvanced(filterAdvanced string) *getChannelsBuilder {
	b.opts.FilterAdvanced = filterAdvanced
	return b
}

// Sort sets the comma-separated list of fields to sort by, each prefixed with
// + for ascending or - for descending. Example: "-createdAt,+id".
func (b *getChannelsBuilder) Sort(sort []string) *getChannelsBuilder {
	b.opts.Sort = sort
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getChannelsBuilder) QueryParam(queryParam map[string]string) *getChannelsBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the getChannels request.
func (b *getChannelsBuilder) Transport(tr http.RoundTripper) *getChannelsBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *getChannelsOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"Limit": o.Limit,
	}
	if o.EntityClass != "" {
		params["EntityClass"] = o.EntityClass
	}
	if o.setEntityClassVersion {
		params["EntityClassVersion"] = o.EntityClassVersion
	}
	if o.EntityClassLevel != "" {
		params["EntityClassLevel"] = o.EntityClassLevel
	}
	if o.Cursor != "" {
		params["Cursor"] = o.Cursor
	}
	if o.Filter != "" {
		params["Filter"] = o.Filter
	}
	if o.FilterAdvanced != "" {
		params["FilterAdvanced"] = o.FilterAdvanced
	}
	if len(o.Sort) > 0 {
		params["Sort"] = o.Sort
	}
	return params
}

// Execute runs the getChannels request.
func (b *getChannelsBuilder) Execute() (*PNDataSyncChannelsResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNGetDataSyncChannelsOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetChannelsResponse, status, err
	}

	return newPNDataSyncChannelsResponse(rawJSON, b.opts, status)
}

type getChannelsOpts struct {
	endpointOpts

	EntityClass           string
	EntityClassVersion    int
	setEntityClassVersion bool
	EntityClassLevel      string
	Cursor                string
	Limit                 int
	Filter                string
	FilterAdvanced        string
	Sort                  []string
	QueryParam            map[string]string

	Transport http.RoundTripper
}

func (o *getChannelsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	return nil
}

func (o *getChannelsOpts) buildPath() (string, error) {
	return fmt.Sprintf(channelsPath, o.pubnub.Config.SubscribeKey), nil
}

func (o *getChannelsOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.EntityClass != "" {
		q.Set("entity_class", o.EntityClass)
	}

	if o.setEntityClassVersion {
		q.Set("entity_class_version", strconv.Itoa(o.EntityClassVersion))
	}

	if o.EntityClassLevel != "" {
		q.Set("entity_class_level", o.EntityClassLevel)
	}

	q.Set("limit", strconv.Itoa(o.Limit))

	if o.Cursor != "" {
		q.Set("cursor", o.Cursor)
	}
	// "filter" and "sort" are URL-encoded centrally in buildURL; "filter_advanced"
	// is not, so it is encoded here.
	if o.Filter != "" {
		q.Set("filter", o.Filter)
	}
	if o.FilterAdvanced != "" {
		q.Set("filter_advanced", utils.URLEncode(o.FilterAdvanced))
	}
	if len(o.Sort) > 0 {
		SetQueryParamAsCommaSepString(q, o.Sort, "sort")
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getChannelsOpts) httpMethod() string {
	return "GET"
}

func (o *getChannelsOpts) operationType() OperationType {
	return PNGetDataSyncChannelsOperation
}

func newPNDataSyncChannelsResponse(jsonBytes []byte, o *getChannelsOpts,
	status StatusResponse) (*PNDataSyncChannelsResponse, StatusResponse, error) {

	resp := &PNDataSyncChannelsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		o.pubnub.loggerManager.LogError(e, "ChannelsResponseParsingFailed", PNGetDataSyncChannelsOperation, true)
		return emptyGetChannelsResponse, status, e
	}

	return resp, status, nil
}
