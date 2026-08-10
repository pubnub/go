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

var emptyGetUsersResponse *PNUsersResponse

type getUsersBuilder struct {
	opts *getUsersOpts
}

func newGetUsersBuilder(pubnub *PubNub) *getUsersBuilder {
	return newGetUsersBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetUsersOpts(pubnub *PubNub, ctx Context) *getUsersOpts {
	return &getUsersOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newGetUsersBuilderWithContext(pubnub *PubNub, context Context) *getUsersBuilder {
	builder := &getUsersBuilder{opts: newGetUsersOpts(pubnub, context)}
	builder.opts.Limit = entitiesDefaultLimit
	return builder
}

// EntityClass sets the optional entity class filter. Must be User or a subclass
// of User when provided.
func (b *getUsersBuilder) EntityClass(entityClass string) *getUsersBuilder {
	b.opts.EntityClass = entityClass
	return b
}

// EntityClassVersion sets the optional schema version to read against. When
// unset, the server uses the latest version.
func (b *getUsersBuilder) EntityClassVersion(version int) *getUsersBuilder {
	b.opts.EntityClassVersion = version
	b.opts.setEntityClassVersion = true
	return b
}

// EntityClassLevel sets the optional class scope (Global or SubKey) used to
// disambiguate classes with the same name defined at different levels.
func (b *getUsersBuilder) EntityClassLevel(level string) *getUsersBuilder {
	b.opts.EntityClassLevel = level
	return b
}

// Cursor sets the pagination token returned from a previous request.
func (b *getUsersBuilder) Cursor(cursor string) *getUsersBuilder {
	b.opts.Cursor = cursor
	return b
}

// Limit sets the maximum number of items to return per page (1-100, default 20).
func (b *getUsersBuilder) Limit(limit int) *getUsersBuilder {
	b.opts.Limit = limit
	return b
}

// Filter sets the AppContext Query Language filter expression.
func (b *getUsersBuilder) Filter(filter string) *getUsersBuilder {
	b.opts.Filter = filter
	return b
}

// FilterAdvanced sets the extended filter expression supporting logical
// operators and nested conditions.
func (b *getUsersBuilder) FilterAdvanced(filterAdvanced string) *getUsersBuilder {
	b.opts.FilterAdvanced = filterAdvanced
	return b
}

// Sort sets the comma-separated list of fields to sort by, each prefixed with
// + for ascending or - for descending. Example: "-createdAt,+id".
func (b *getUsersBuilder) Sort(sort []string) *getUsersBuilder {
	b.opts.Sort = sort
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getUsersBuilder) QueryParam(queryParam map[string]string) *getUsersBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the getUsers request.
func (b *getUsersBuilder) Transport(tr http.RoundTripper) *getUsersBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *getUsersOpts) GetLogParams() map[string]interface{} {
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

// Execute runs the getUsers request.
func (b *getUsersBuilder) Execute() (*PNUsersResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNGetDataSyncUsersOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetUsersResponse, status, err
	}

	return newPNUsersResponse(rawJSON, b.opts, status)
}

type getUsersOpts struct {
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

func (o *getUsersOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	return nil
}

func (o *getUsersOpts) buildPath() (string, error) {
	return fmt.Sprintf(usersPath, o.pubnub.Config.SubscribeKey), nil
}

func (o *getUsersOpts) buildQuery() (*url.Values, error) {
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

func (o *getUsersOpts) httpMethod() string {
	return "GET"
}

func (o *getUsersOpts) operationType() OperationType {
	return PNGetDataSyncUsersOperation
}

func newPNUsersResponse(jsonBytes []byte, o *getUsersOpts,
	status StatusResponse) (*PNUsersResponse, StatusResponse, error) {

	resp := &PNUsersResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		o.pubnub.loggerManager.LogError(e, "UsersResponseParsingFailed", PNGetDataSyncUsersOperation, true)
		return emptyGetUsersResponse, status, e
	}

	return resp, status, nil
}
