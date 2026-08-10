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

var emptyGetDataSyncMembershipsResponse *PNDataSyncMembershipsResponse

type getMembershipsBuilder struct {
	opts *getMembershipsOpts
}

func newGetMembershipsBuilder(pubnub *PubNub) *getMembershipsBuilder {
	return newGetMembershipsBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetMembershipsOpts(pubnub *PubNub, ctx Context) *getMembershipsOpts {
	return &getMembershipsOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newGetMembershipsBuilderWithContext(pubnub *PubNub, context Context) *getMembershipsBuilder {
	builder := &getMembershipsBuilder{opts: newGetMembershipsOpts(pubnub, context)}
	builder.opts.Limit = entitiesDefaultLimit
	return builder
}

// UserID filters memberships by the linked user.
func (b *getMembershipsBuilder) UserID(userID string) *getMembershipsBuilder {
	b.opts.UserID = userID
	return b
}

// ChannelID filters memberships by the linked channel.
func (b *getMembershipsBuilder) ChannelID(channelID string) *getMembershipsBuilder {
	b.opts.ChannelID = channelID
	return b
}

// RelationshipClassVersion sets the optional schema version to filter by. When
// unset, the server returns all versions (filter/sort use the latest schema).
func (b *getMembershipsBuilder) RelationshipClassVersion(version int) *getMembershipsBuilder {
	b.opts.RelationshipClassVersion = version
	b.opts.setRelationshipClassVersion = true
	return b
}

// Cursor sets the pagination token returned from a previous request.
func (b *getMembershipsBuilder) Cursor(cursor string) *getMembershipsBuilder {
	b.opts.Cursor = cursor
	return b
}

// Limit sets the maximum number of items to return per page (1-100, default 20).
func (b *getMembershipsBuilder) Limit(limit int) *getMembershipsBuilder {
	b.opts.Limit = limit
	return b
}

// Filter sets the AppContext Query Language filter expression.
func (b *getMembershipsBuilder) Filter(filter string) *getMembershipsBuilder {
	b.opts.Filter = filter
	return b
}

// FilterAdvanced sets the extended filter expression supporting logical
// operators and nested conditions.
func (b *getMembershipsBuilder) FilterAdvanced(filterAdvanced string) *getMembershipsBuilder {
	b.opts.FilterAdvanced = filterAdvanced
	return b
}

// Sort sets the comma-separated list of fields to sort by, each prefixed with
// + for ascending or - for descending. Example: "-createdAt,+id".
func (b *getMembershipsBuilder) Sort(sort []string) *getMembershipsBuilder {
	b.opts.Sort = sort
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getMembershipsBuilder) QueryParam(queryParam map[string]string) *getMembershipsBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the getMemberships request.
func (b *getMembershipsBuilder) Transport(tr http.RoundTripper) *getMembershipsBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *getMembershipsOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"Limit": o.Limit,
	}
	if o.UserID != "" {
		params["UserID"] = o.UserID
	}
	if o.ChannelID != "" {
		params["ChannelID"] = o.ChannelID
	}
	if o.setRelationshipClassVersion {
		params["RelationshipClassVersion"] = o.RelationshipClassVersion
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

// Execute runs the getMemberships request.
func (b *getMembershipsBuilder) Execute() (*PNDataSyncMembershipsResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNGetDataSyncMembershipsOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetDataSyncMembershipsResponse, status, err
	}

	return newPNDataSyncMembershipsResponse(rawJSON, b.opts, status)
}

type getMembershipsOpts struct {
	endpointOpts

	UserID                      string
	ChannelID                   string
	RelationshipClassVersion    int
	setRelationshipClassVersion bool
	Cursor                      string
	Limit                       int
	Filter                      string
	FilterAdvanced              string
	Sort                        []string
	QueryParam                  map[string]string

	Transport http.RoundTripper
}

func (o *getMembershipsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	return nil
}

func (o *getMembershipsOpts) buildPath() (string, error) {
	return fmt.Sprintf(membershipsPath, o.pubnub.Config.SubscribeKey), nil
}

func (o *getMembershipsOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.UserID != "" {
		q.Set("user_id", o.UserID)
	}
	if o.ChannelID != "" {
		q.Set("channel_id", o.ChannelID)
	}
	if o.setRelationshipClassVersion {
		q.Set("relationship_class_version", strconv.Itoa(o.RelationshipClassVersion))
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

func (o *getMembershipsOpts) httpMethod() string {
	return "GET"
}

func (o *getMembershipsOpts) operationType() OperationType {
	return PNGetDataSyncMembershipsOperation
}

func newPNDataSyncMembershipsResponse(jsonBytes []byte, o *getMembershipsOpts,
	status StatusResponse) (*PNDataSyncMembershipsResponse, StatusResponse, error) {

	resp := &PNDataSyncMembershipsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		o.pubnub.loggerManager.LogError(e, "MembershipsResponseParsingFailed", PNGetDataSyncMembershipsOperation, true)
		return emptyGetDataSyncMembershipsResponse, status, e
	}

	return resp, status, nil
}
