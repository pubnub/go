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

var emptyGetRelationshipsResponse *PNRelationshipsResponse

type getRelationshipsBuilder struct {
	opts *getRelationshipsOpts
}

func newGetRelationshipsBuilder(pubnub *PubNub) *getRelationshipsBuilder {
	return newGetRelationshipsBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetRelationshipsOpts(pubnub *PubNub, ctx Context) *getRelationshipsOpts {
	return &getRelationshipsOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newGetRelationshipsBuilderWithContext(pubnub *PubNub, context Context) *getRelationshipsBuilder {
	builder := &getRelationshipsBuilder{opts: newGetRelationshipsOpts(pubnub, context)}
	builder.opts.Limit = entitiesDefaultLimit
	return builder
}

// RelationshipClass sets the required relationship class name to list relationships for.
func (b *getRelationshipsBuilder) RelationshipClass(relationshipClass string) *getRelationshipsBuilder {
	b.opts.RelationshipClass = relationshipClass
	return b
}

// RelationshipClassVersion sets the optional schema version to filter by. When
// unset, the server returns all versions (filter/sort use the latest schema).
func (b *getRelationshipsBuilder) RelationshipClassVersion(version int) *getRelationshipsBuilder {
	b.opts.RelationshipClassVersion = version
	b.opts.setRelationshipClassVersion = true
	return b
}

// EntityAID filters relationships by the first linked entity.
func (b *getRelationshipsBuilder) EntityAID(entityAID string) *getRelationshipsBuilder {
	b.opts.EntityAID = entityAID
	return b
}

// EntityBID filters relationships by the second linked entity.
func (b *getRelationshipsBuilder) EntityBID(entityBID string) *getRelationshipsBuilder {
	b.opts.EntityBID = entityBID
	return b
}

// Cursor sets the pagination token returned from a previous request.
func (b *getRelationshipsBuilder) Cursor(cursor string) *getRelationshipsBuilder {
	b.opts.Cursor = cursor
	return b
}

// Limit sets the maximum number of items to return per page (1-100, default 20).
func (b *getRelationshipsBuilder) Limit(limit int) *getRelationshipsBuilder {
	b.opts.Limit = limit
	return b
}

// FilterFast sets the strongly consistent filter expression (query param
// filter_fast). It always reflects the latest writes but accepts fewer
// conditions than Filter. FilterFast and Filter cannot both be set.
func (b *getRelationshipsBuilder) FilterFast(filterFast string) *getRelationshipsBuilder {
	b.opts.FilterFast = filterFast
	return b
}

// Filter sets the eventually consistent filter expression (query param
// filter). It accepts more conditions than FilterFast; recent writes may
// not yet be reflected. FilterFast and Filter cannot both be set.
func (b *getRelationshipsBuilder) Filter(filter string) *getRelationshipsBuilder {
	b.opts.Filter = filter
	return b
}

// Sort sets the comma-separated list of fields to sort by, each prefixed with
// + for ascending or - for descending. Example: "-createdAt,+id".
func (b *getRelationshipsBuilder) Sort(sort []string) *getRelationshipsBuilder {
	b.opts.Sort = sort
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getRelationshipsBuilder) QueryParam(queryParam map[string]string) *getRelationshipsBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the getRelationships request.
func (b *getRelationshipsBuilder) Transport(tr http.RoundTripper) *getRelationshipsBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *getRelationshipsOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"RelationshipClass": o.RelationshipClass,
		"Limit":             o.Limit,
	}
	if o.setRelationshipClassVersion {
		params["RelationshipClassVersion"] = o.RelationshipClassVersion
	}
	if o.EntityAID != "" {
		params["EntityAID"] = o.EntityAID
	}
	if o.EntityBID != "" {
		params["EntityBID"] = o.EntityBID
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

// Execute runs the getRelationships request.
func (b *getRelationshipsBuilder) Execute() (*PNRelationshipsResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNGetRelationshipsOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetRelationshipsResponse, status, err
	}

	return newPNRelationshipsResponse(rawJSON, b.opts, status)
}

type getRelationshipsOpts struct {
	endpointOpts

	RelationshipClass           string
	RelationshipClassVersion    int
	setRelationshipClassVersion bool
	EntityAID                   string
	EntityBID                   string
	Cursor                      string
	Limit                       int
	FilterFast                  string
	Filter                      string
	Sort                        []string
	QueryParam                  map[string]string

	Transport http.RoundTripper
}

func (o *getRelationshipsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.RelationshipClass == "" {
		return newValidationError(o, StrMissingRelationshipClass)
	}
	if conflictingDataSyncFilters(o.FilterFast, o.Filter) {
		return newValidationError(o, StrExclusiveDataSyncFilter)
	}
	return nil
}

func (o *getRelationshipsOpts) buildPath() (string, error) {
	return fmt.Sprintf(relationshipsPath, o.pubnub.Config.SubscribeKey), nil
}

func (o *getRelationshipsOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	q.Set("relationship_class", o.RelationshipClass)

	if o.setRelationshipClassVersion {
		q.Set("relationship_class_version", strconv.Itoa(o.RelationshipClassVersion))
	}
	if o.EntityAID != "" {
		q.Set("entity_a_id", o.EntityAID)
	}
	if o.EntityBID != "" {
		q.Set("entity_b_id", o.EntityBID)
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

func (o *getRelationshipsOpts) httpMethod() string {
	return "GET"
}

func (o *getRelationshipsOpts) operationType() OperationType {
	return PNGetRelationshipsOperation
}

func newPNRelationshipsResponse(jsonBytes []byte, o *getRelationshipsOpts,
	status StatusResponse) (*PNRelationshipsResponse, StatusResponse, error) {

	resp := &PNRelationshipsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		o.pubnub.loggerManager.LogError(e, "RelationshipsResponseParsingFailed", PNGetRelationshipsOperation, true)
		return emptyGetRelationshipsResponse, status, e
	}

	return resp, status, nil
}
