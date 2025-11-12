package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/pubnub/go/v8/pnerr"
)

const grantTokenPath = "/v3/pam/%s/grant"

var emptyPNGrantTokenResponse *PNGrantTokenResponse

type grantTokenBuilder struct {
	opts *grantTokenOpts
}

type SpaceId string

type grantTokenEntitiesBuilder grantTokenBuilder

type grantTokenObjectsBuilder grantTokenBuilder

func newGrantTokenBuilder(pubnub *PubNub) *grantTokenBuilder {
	return newGrantTokenBuilderWithContext(pubnub, pubnub.ctx)
}

func newGrantTokenObjectsBuilder(opts *grantTokenOpts) *grantTokenObjectsBuilder {
	return &grantTokenObjectsBuilder{opts}
}

func newGrantTokenEntitiesBuilder(opts *grantTokenOpts) *grantTokenEntitiesBuilder {
	return &grantTokenEntitiesBuilder{opts}
}

func newGrantTokenBuilderWithContext(pubnub *PubNub, context Context) *grantTokenBuilder {
	builder := grantTokenBuilder{
		opts: newGrantTokenOpts(pubnub, context)}
	return &builder
}

// TTL in minutes for which granted permissions are valid.
//
// Min: 1
// Max: 525600
// Default: 1440
//
// Setting value to 0 will apply the grant indefinitely (forever grant).
func (b *grantTokenBuilder) TTL(ttl int) *grantTokenBuilder {
	b.opts.TTL = ttl
	b.opts.setTTL = true

	return b
}

// Meta sets the Meta for the Grant request.
func (b *grantTokenBuilder) Meta(meta map[string]interface{}) *grantTokenBuilder {
	b.opts.Meta = meta

	return b
}

func (b *grantTokenBuilder) AuthorizedUUID(uuid string) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).AuthorizedUUID(uuid)
}

func (b *grantTokenBuilder) AuthorizedUserId(userId UserId) *grantTokenEntitiesBuilder {
	return newGrantTokenEntitiesBuilder(b.opts).AuthorizedUserId(userId)
}

// Channels sets the Channels for the Grant request.
func (b *grantTokenBuilder) Channels(channels map[string]ChannelPermissions) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).Channels(channels)
}

// ChannelGroups sets the ChannelGroups for the Grant request.
func (b *grantTokenBuilder) ChannelGroups(groups map[string]GroupPermissions) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).ChannelGroups(groups)
}

func (b *grantTokenBuilder) UUIDs(uuids map[string]UUIDPermissions) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).UUIDs(uuids)
}

// ChannelsPattern sets the ChannelPermissions for the Grant request.
func (b *grantTokenBuilder) ChannelsPattern(channels map[string]ChannelPermissions) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).ChannelsPattern(channels)
}

// ChannelGroupsPattern sets the GroupPermissions for the Grant request.
func (b *grantTokenBuilder) ChannelGroupsPattern(groups map[string]GroupPermissions) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).ChannelGroupsPattern(groups)
}

func (b *grantTokenBuilder) UUIDsPattern(uuids map[string]UUIDPermissions) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).UUIDsPattern(uuids)
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *grantTokenBuilder) QueryParam(queryParam map[string]string) *grantTokenBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *grantTokenOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{}
	if len(o.AuthKeys) > 0 {
		params["AuthKeys"] = o.AuthKeys
	}
	if len(o.Channels) > 0 {
		channelNames := make([]string, 0, len(o.Channels))
		for k := range o.Channels {
			channelNames = append(channelNames, k)
		}
		params["Channels"] = channelNames
	}
	if len(o.ChannelGroups) > 0 {
		groupNames := make([]string, 0, len(o.ChannelGroups))
		for k := range o.ChannelGroups {
			groupNames = append(groupNames, k)
		}
		params["ChannelGroups"] = groupNames
	}
	if len(o.UUIDs) > 0 {
		uuidNames := make([]string, 0, len(o.UUIDs))
		for k := range o.UUIDs {
			uuidNames = append(uuidNames, k)
		}
		params["UUIDs"] = uuidNames
	}
	if len(o.ChannelsPattern) > 0 {
		params["ChannelsPattern"] = fmt.Sprintf("(%d patterns)", len(o.ChannelsPattern))
	}
	if len(o.ChannelGroupsPattern) > 0 {
		params["ChannelGroupsPattern"] = fmt.Sprintf("(%d patterns)", len(o.ChannelGroupsPattern))
	}
	if len(o.UUIDsPattern) > 0 {
		params["UUIDsPattern"] = fmt.Sprintf("(%d patterns)", len(o.UUIDsPattern))
	}
	if o.AuthorizedUUID != "" {
		params["AuthorizedUUID"] = o.AuthorizedUUID
	}
	if o.setTTL {
		params["TTL"] = o.TTL
	}
	if o.Meta != nil {
		params["Meta"] = fmt.Sprintf("%v", o.Meta)
	}
	return params
}

// Execute runs the Grant request.
func (b *grantTokenBuilder) Execute() (*PNGrantTokenResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNAccessManagerGrantToken, b.opts.GetLogParams(), true)
	
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNGrantTokenResponse, status, err
	}
	resp, status, e := newGrantTokenResponse(rawJSON, status)
	if e != nil {
		b.opts.pubnub.tokenManager.StoreToken(resp.Data.Token)
	}

	return resp, status, e
}

func (b *grantTokenBuilder) SpacesPermissions(spacesPermissions map[SpaceId]SpacePermissions) *grantTokenEntitiesBuilder {
	return newGrantTokenEntitiesBuilder(b.opts).SpacesPermissions(spacesPermissions)
}

func (b *grantTokenBuilder) UsersPermissions(usersPermissions map[UserId]UserPermissions) *grantTokenEntitiesBuilder {
	return newGrantTokenEntitiesBuilder(b.opts).UsersPermissions(usersPermissions)
}

func (b *grantTokenBuilder) SpacePatternsPermissions(spacePatternsPermissions map[string]SpacePermissions) *grantTokenEntitiesBuilder {
	return newGrantTokenEntitiesBuilder(b.opts).SpacePatternsPermissions(spacePatternsPermissions)
}

func (b *grantTokenBuilder) UserPatternsPermissions(userPatternsPermissions map[string]UserPermissions) *grantTokenEntitiesBuilder {
	return newGrantTokenEntitiesBuilder(b.opts).UserPatternsPermissions(userPatternsPermissions)
}

// TTL in minutes for which granted permissions are valid.
//
// Min: 1
// Max: 525600
// Default: 1440
//
// Setting value to 0 will apply the grant indefinitely (forever grant).
func (b *grantTokenObjectsBuilder) TTL(ttl int) *grantTokenObjectsBuilder {
	b.opts.TTL = ttl
	b.opts.setTTL = true

	return b
}

// Meta sets the Meta for the Grant request.
func (b *grantTokenObjectsBuilder) Meta(meta map[string]interface{}) *grantTokenObjectsBuilder {
	b.opts.Meta = meta

	return b
}

func (b *grantTokenObjectsBuilder) AuthorizedUUID(uuid string) *grantTokenObjectsBuilder {
	b.opts.AuthorizedUUID = uuid

	return b
}

// Channels sets the Channels for the Grant request.
func (b *grantTokenObjectsBuilder) Channels(channels map[string]ChannelPermissions) *grantTokenObjectsBuilder {
	b.opts.Channels = channels

	return b
}

// ChannelGroups sets the ChannelGroups for the Grant request.
func (b *grantTokenObjectsBuilder) ChannelGroups(groups map[string]GroupPermissions) *grantTokenObjectsBuilder {
	b.opts.ChannelGroups = groups

	return b
}

func (b *grantTokenObjectsBuilder) UUIDs(uuids map[string]UUIDPermissions) *grantTokenObjectsBuilder {
	b.opts.UUIDs = uuids

	return b
}

// Channels sets the Channels for the Grant request.
func (b *grantTokenObjectsBuilder) ChannelsPattern(channels map[string]ChannelPermissions) *grantTokenObjectsBuilder {
	b.opts.ChannelsPattern = channels

	return b
}

// ChannelGroups sets the ChannelGroups for the Grant request.
func (b *grantTokenObjectsBuilder) ChannelGroupsPattern(groups map[string]GroupPermissions) *grantTokenObjectsBuilder {
	b.opts.ChannelGroupsPattern = groups

	return b
}

func (b *grantTokenObjectsBuilder) UUIDsPattern(uuids map[string]UUIDPermissions) *grantTokenObjectsBuilder {
	b.opts.UUIDsPattern = uuids

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *grantTokenObjectsBuilder) QueryParam(queryParam map[string]string) *grantTokenObjectsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Grant request.
func (b *grantTokenObjectsBuilder) Execute() (*PNGrantTokenResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNAccessManagerGrantToken, b.opts.GetLogParams(), true)
	
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNGrantTokenResponse, status, err
	}
	resp, status, e := newGrantTokenResponse(rawJSON, status)
	if e != nil {
		b.opts.pubnub.tokenManager.StoreToken(resp.Data.Token)
	}

	return resp, status, e
}

// TTL in minutes for which granted permissions are valid.
//
// Min: 1
// Max: 525600
// Default: 1440
//
// Setting value to 0 will apply the grant indefinitely (forever grant).
func (b *grantTokenEntitiesBuilder) TTL(ttl int) *grantTokenEntitiesBuilder {
	b.opts.TTL = ttl
	b.opts.setTTL = true

	return b
}

func (b *grantTokenEntitiesBuilder) AuthorizedUserId(userId UserId) *grantTokenEntitiesBuilder {
	b.opts.AuthorizedUUID = string(userId)

	return b
}

// SpacesPermissions sets the Spaces for the Grant request.
func (b *grantTokenEntitiesBuilder) SpacesPermissions(spaces map[SpaceId]SpacePermissions) *grantTokenEntitiesBuilder {
	b.opts.Channels = toChannelsPermissionsMap(spaces)

	return b
}

func (b *grantTokenEntitiesBuilder) UsersPermissions(users map[UserId]UserPermissions) *grantTokenEntitiesBuilder {
	b.opts.UUIDs = toUUIDsPermissionsMap(users)

	return b
}

// SpacePatternsPermissions sets the Channels for the Grant request.
func (b *grantTokenEntitiesBuilder) SpacePatternsPermissions(spaces map[string]SpacePermissions) *grantTokenEntitiesBuilder {
	b.opts.ChannelsPattern = toChannelPatternsPermissionsMap(spaces)

	return b
}

func (b *grantTokenEntitiesBuilder) UserPatternsPermissions(users map[string]UserPermissions) *grantTokenEntitiesBuilder {
	b.opts.UUIDsPattern = toUUIDPatternsPermissionsMap(users)

	return b
}

// Meta sets the Meta for the Grant request.
func (b *grantTokenEntitiesBuilder) Meta(meta map[string]interface{}) *grantTokenEntitiesBuilder {
	b.opts.Meta = meta

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *grantTokenEntitiesBuilder) QueryParam(queryParam map[string]string) *grantTokenEntitiesBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Grant request.
func (b *grantTokenEntitiesBuilder) Execute() (*PNGrantTokenResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNAccessManagerGrantToken, b.opts.GetLogParams(), true)
	
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNGrantTokenResponse, status, err
	}
	resp, status, e := newGrantTokenResponse(rawJSON, status)
	if e != nil {
		b.opts.pubnub.tokenManager.StoreToken(resp.Data.Token)
	}

	return resp, status, e
}

func newGrantTokenOpts(pubnub *PubNub, ctx Context) *grantTokenOpts {
	return &grantTokenOpts{
		endpointOpts: endpointOpts{
			pubnub: pubnub,
			ctx:    ctx,
		},
	}
}

type grantTokenOpts struct {
	endpointOpts
	AuthKeys             []string
	Channels             map[string]ChannelPermissions
	ChannelGroups        map[string]GroupPermissions
	UUIDs                map[string]UUIDPermissions
	ChannelsPattern      map[string]ChannelPermissions
	ChannelGroupsPattern map[string]GroupPermissions
	UUIDsPattern         map[string]UUIDPermissions
	QueryParam           map[string]string
	Meta                 map[string]interface{}
	AuthorizedUUID       string

	// Max: 525600
	// Min: 1
	// Default: 1440
	// Setting 0 will apply the grant indefinitely
	TTL int

	// nil hacks
	setTTL bool
}

func (o *grantTokenOpts) validate() error {
	if o.config().PublishKey == "" {
		return newValidationError(o, StrMissingPubKey)
	}

	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.config().SecretKey == "" {
		return newValidationError(o, StrMissingSecretKey)
	}

	if o.TTL <= 0 {
		return newValidationError(o, StrInvalidTTL)
	}

	return nil
}

func (o *grantTokenOpts) buildPath() (string, error) {
	return fmt.Sprintf(grantTokenPath, o.pubnub.Config.SubscribeKey), nil
}

type grantBody struct {
	TTL         int             `json:"ttl"`
	Permissions PermissionsBody `json:"permissions"`
}

func (o *grantTokenOpts) setBitmask(value bool, bitmask PNGrantBitMask, bm int64) int64 {
	if value {
		bm |= int64(bitmask)
	}
	o.pubnub.loggerManager.LogSimple(PNLogLevelTrace, fmt.Sprintf("Grant token: bitmask value=%t, mask=%d, result=%d", value, bitmask, bm), false)
	return bm
}

func (o *grantTokenOpts) parseResourcePermissions(resource interface{}, resourceType PNResourceType) map[string]int64 {
	bmVal := int64(0)
	switch resourceType {
	case PNChannels:
		resourceWithPerms := resource.(map[string]ChannelPermissions)
		resourceWithPermsLen := len(resourceWithPerms)
		if resourceWithPermsLen > 0 {
			r := make(map[string]int64, resourceWithPermsLen)
			for k, v := range resourceWithPerms {
				bmVal = int64(0)
				bmVal = o.setBitmask(v.Read, PNRead, bmVal)
				bmVal = o.setBitmask(v.Write, PNWrite, bmVal)
				bmVal = o.setBitmask(v.Delete, PNDelete, bmVal)
				bmVal = o.setBitmask(v.Join, PNJoin, bmVal)
				bmVal = o.setBitmask(v.Update, PNUpdate, bmVal)
				bmVal = o.setBitmask(v.Manage, PNManage, bmVal)
				bmVal = o.setBitmask(v.Get, PNGet, bmVal)
				o.pubnub.loggerManager.LogSimple(PNLogLevelTrace, fmt.Sprintf("Grant token: channel permissions bitmask=%d", bmVal), false)
				r[k] = bmVal
			}
			return r
		}
		return make(map[string]int64)

	case PNGroups:
		resourceWithPerms := resource.(map[string]GroupPermissions)
		resourceWithPermsLen := len(resourceWithPerms)
		if resourceWithPermsLen > 0 {
			r := make(map[string]int64, resourceWithPermsLen)
			for k, v := range resourceWithPerms {
				bmVal = int64(0)
				bmVal = o.setBitmask(v.Read, PNRead, bmVal)
				bmVal = o.setBitmask(v.Manage, PNManage, bmVal)
				o.pubnub.loggerManager.LogSimple(PNLogLevelTrace, fmt.Sprintf("Grant token: group permissions bitmask=%d", bmVal), false)
				r[k] = bmVal
			}
			return r
		}
		return make(map[string]int64)

	case PNUUIDs:
		resourceWithPerms := resource.(map[string]UUIDPermissions)
		resourceWithPermsLen := len(resourceWithPerms)
		if resourceWithPermsLen > 0 {
			r := make(map[string]int64, resourceWithPermsLen)
			for k, v := range resourceWithPerms {
				bmVal = int64(0)
				bmVal = o.setBitmask(v.Get, PNGet, bmVal)
				bmVal = o.setBitmask(v.Update, PNUpdate, bmVal)
				bmVal = o.setBitmask(v.Delete, PNDelete, bmVal)
				o.pubnub.loggerManager.LogSimple(PNLogLevelTrace, fmt.Sprintf("Grant token: UUID permissions bitmask=%d", bmVal), false)
				r[k] = bmVal
			}
			return r
		}
		return make(map[string]int64)
	default:
		return make(map[string]int64)
	}

}

func (o *grantTokenOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *grantTokenOpts) buildBody() ([]byte, error) {

	meta := o.Meta

	if meta == nil {
		meta = make(map[string]interface{})
	}

	permissions := PermissionsBody{
		Resources: GrantResources{
			Channels: o.parseResourcePermissions(o.Channels, PNChannels),
			Groups:   o.parseResourcePermissions(o.ChannelGroups, PNGroups),
			UUIDs:    o.parseResourcePermissions(o.UUIDs, PNUUIDs),
			Users:    make(map[string]int64),
			Spaces:   make(map[string]int64),
		},
		Patterns: GrantResources{
			Channels: o.parseResourcePermissions(o.ChannelsPattern, PNChannels),
			Groups:   o.parseResourcePermissions(o.ChannelGroupsPattern, PNGroups),
			UUIDs:    o.parseResourcePermissions(o.UUIDsPattern, PNUUIDs),
			Users:    make(map[string]int64),
			Spaces:   make(map[string]int64),
		},
		Meta:           meta,
		AuthorizedUUID: o.AuthorizedUUID,
	}

	o.pubnub.loggerManager.LogSimple(PNLogLevelTrace, fmt.Sprintf("Grant token: permissions=%+v", permissions), false)

	ttl := -1
	if o.setTTL {
		if o.TTL >= -1 {
			ttl = o.TTL
		}
	}

	b := &grantBody{
		TTL:         ttl,
		Permissions: permissions,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "GrantTokenSerializationFailed", PNAccessManagerGrantToken, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *grantTokenOpts) httpMethod() string {
	return "POST"
}

func (o *grantTokenOpts) operationType() OperationType {
	return PNAccessManagerGrantToken
}

// PNGrantTokenData is the struct used to decode the server response
type PNGrantTokenData struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

// PNGrantTokenResponse is the struct returned when the Execute function of Grant Token is called.
type PNGrantTokenResponse struct {
	Status  int              `json:"status"`
	Data    PNGrantTokenData `json:"data"`
	Service string           `json:"service"`
}

func newGrantTokenResponse(jsonBytes []byte, status StatusResponse) (*PNGrantTokenResponse, StatusResponse, error) {
	resp := &PNGrantTokenResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGrantTokenResponse, status, e
	}

	return resp, status, nil
}
