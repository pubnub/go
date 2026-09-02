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
)

var emptyGetEntitiesResponse *PNEntitiesResponse

type getEntitiesBuilder struct {
	opts *getEntitiesOpts
}

func newGetEntitiesBuilder(pubnub *PubNub) *getEntitiesBuilder {
	return newGetEntitiesBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetEntitiesOpts(pubnub *PubNub, ctx Context) *getEntitiesOpts {
	return &getEntitiesOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newGetEntitiesBuilderWithContext(pubnub *PubNub, context Context) *getEntitiesBuilder {
	builder := &getEntitiesBuilder{opts: newGetEntitiesOpts(pubnub, context)}
	builder.opts.Limit = entitiesDefaultLimit
	return builder
}

// EntityClass sets the required entity class name to list entities for.
func (b *getEntitiesBuilder) EntityClass(entityClass string) *getEntitiesBuilder {
	b.opts.EntityClass = entityClass
	return b
}

// EntityClassVersion sets the optional schema version to read against. When
// unset, the server uses the latest version.
func (b *getEntitiesBuilder) EntityClassVersion(version int) *getEntitiesBuilder {
	b.opts.EntityClassVersion = version
	b.opts.setEntityClassVersion = true
	return b
}

// EntityClassLevel sets the optional class scope (Global, Account, or SubKey) used to
// disambiguate classes with the same name defined at different levels.
func (b *getEntitiesBuilder) EntityClassLevel(level PNEntityClassLevel) *getEntitiesBuilder {
	b.opts.EntityClassLevel = level
	return b
}

// Cursor sets the pagination token returned from a previous request.
func (b *getEntitiesBuilder) Cursor(cursor string) *getEntitiesBuilder {
	b.opts.Cursor = cursor
	return b
}

// Limit sets the maximum number of items to return per page (1-100, default 20).
func (b *getEntitiesBuilder) Limit(limit int) *getEntitiesBuilder {
	b.opts.Limit = limit
	return b
}

// FilterFast sets the strongly consistent filter expression (query param
// filter_fast). It always reflects the latest writes but accepts fewer
// conditions than Filter. FilterFast and Filter cannot both be set.
func (b *getEntitiesBuilder) FilterFast(filterFast string) *getEntitiesBuilder {
	b.opts.FilterFast = filterFast
	return b
}

// Filter sets the eventually consistent filter expression (query param
// filter). It accepts more conditions than FilterFast; recent writes may
// not yet be reflected. FilterFast and Filter cannot both be set.
func (b *getEntitiesBuilder) Filter(filter string) *getEntitiesBuilder {
	b.opts.Filter = filter
	return b
}

// Sort sets the comma-separated list of fields to sort by, each prefixed with
// + for ascending or - for descending. Example: "-createdAt,+id".
func (b *getEntitiesBuilder) Sort(sort []string) *getEntitiesBuilder {
	b.opts.Sort = sort
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getEntitiesBuilder) QueryParam(queryParam map[string]string) *getEntitiesBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the getEntities request.
func (b *getEntitiesBuilder) Transport(tr http.RoundTripper) *getEntitiesBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *getEntitiesOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"EntityClass": o.EntityClass,
		"Limit":       o.Limit,
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
	if o.FilterFast != "" {
		params["FilterFast"] = o.FilterFast
	}
	if o.Filter != "" {
		params["Filter"] = o.Filter
	}
	if len(o.Sort) > 0 {
		params["Sort"] = o.Sort
	}
	return params
}

// Execute runs the getEntities request.
func (b *getEntitiesBuilder) Execute() (*PNEntitiesResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNGetEntitiesOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetEntitiesResponse, status, err
	}

	return newPNEntitiesResponse(rawJSON, b.opts, status)
}

type getEntitiesOpts struct {
	endpointOpts

	EntityClass           string
	EntityClassVersion    int
	setEntityClassVersion bool
	EntityClassLevel      PNEntityClassLevel
	Cursor                string
	Limit                 int
	FilterFast            string
	Filter                string
	Sort                  []string
	QueryParam            map[string]string

	Transport http.RoundTripper
}

func (o *getEntitiesOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.EntityClass == "" {
		return newValidationError(o, StrMissingEntityClass)
	}
	if conflictingDataSyncFilters(o.FilterFast, o.Filter) {
		return newValidationError(o, StrExclusiveDataSyncFilter)
	}
	return nil
}

func (o *getEntitiesOpts) buildPath() (string, error) {
	return fmt.Sprintf(entitiesPath, o.pubnub.Config.SubscribeKey), nil
}

func (o *getEntitiesOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	q.Set("entity_class", o.EntityClass)

	if o.setEntityClassVersion {
		q.Set("entity_class_version", strconv.Itoa(o.EntityClassVersion))
	}

	if o.EntityClassLevel != "" {
		q.Set("entity_class_level", string(o.EntityClassLevel))
	}

	q.Set("limit", strconv.Itoa(o.Limit))

	if o.Cursor != "" {
		q.Set("cursor", o.Cursor)
	}
	// "filter", "filter_fast" and "sort" are URL-encoded centrally in buildURL.
	if o.FilterFast != "" {
		q.Set("filter_fast", o.FilterFast)
	}
	if o.Filter != "" {
		q.Set("filter", o.Filter)
	}
	if len(o.Sort) > 0 {
		SetQueryParamAsCommaSepString(q, o.Sort, "sort")
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getEntitiesOpts) httpMethod() string {
	return "GET"
}

func (o *getEntitiesOpts) operationType() OperationType {
	return PNGetEntitiesOperation
}

func newPNEntitiesResponse(jsonBytes []byte, o *getEntitiesOpts,
	status StatusResponse) (*PNEntitiesResponse, StatusResponse, error) {

	resp := &PNEntitiesResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		o.pubnub.loggerManager.LogError(e, "EntitiesResponseParsingFailed", PNGetEntitiesOperation, true)
		return emptyGetEntitiesResponse, status, e
	}

	return resp, status, nil
}
