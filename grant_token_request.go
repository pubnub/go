package pubnub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v7/pnerr"
)

const grantTokenPath = "/v3/pam/%s/grant"

var emptyPNGrantTokenResponse *PNGrantTokenResponse

type grantTokenBuilder struct {
	opts *grantTokenOpts
}

type SpaceId string

type grantTokenSumBuilder grantTokenBuilder

type grantTokenObjectsBuilder grantTokenBuilder

func newGrantTokenBuilder(pubnub *PubNub) *grantTokenBuilder {
	builder := grantTokenBuilder{
		opts: &grantTokenOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newGrantTokenObjectsBuilder(opts *grantTokenOpts) *grantTokenObjectsBuilder {
	return &grantTokenObjectsBuilder{opts}
}

func newGrantTokenSumBuilder(opts *grantTokenOpts) *grantTokenSumBuilder {
	return &grantTokenSumBuilder{opts}
}

func newGrantTokenBuilderWithContext(pubnub *PubNub, context Context) *grantTokenBuilder {
	builder := grantTokenBuilder{
		opts: &grantTokenOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

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

// Deprecated: use AuthorizedUserId instead
func (b *grantTokenBuilder) AuthorizedUUID(uuid string) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).AuthorizedUUID(uuid)
}

func (b *grantTokenBuilder) AuthorizedUserId(userId UserId) *grantTokenSumBuilder {
	return newGrantTokenSumBuilder(b.opts).AuthorizedUserId(userId)
}

// Channels sets the Channels for the Grant request.
// Deprecated: Use SpacesPermissions instead
func (b *grantTokenBuilder) Channels(channels map[string]ChannelPermissions) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).Channels(channels)
}

// ChannelGroups sets the ChannelGroups for the Grant request.
// Deprecated
func (b *grantTokenBuilder) ChannelGroups(groups map[string]GroupPermissions) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).ChannelGroups(groups)
}

// Deprecated: Use UsersPermissions instead
func (b *grantTokenBuilder) UUIDs(uuids map[string]UUIDPermissions) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).UUIDs(uuids)
}

// ChannelsPattern sets the ChannelPermissions for the Grant request.
// Deprecated: Use SpacePatternsPermissions instead
func (b *grantTokenBuilder) ChannelsPattern(channels map[string]ChannelPermissions) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).ChannelsPattern(channels)
}

// ChannelGroupsPattern sets the GroupPermissions for the Grant request.
// Deprecated
func (b *grantTokenBuilder) ChannelGroupsPattern(groups map[string]GroupPermissions) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).ChannelGroupsPattern(groups)
}

// Deprecated: Use UserPatternsPermissions instead
func (b *grantTokenBuilder) UUIDsPattern(uuids map[string]UUIDPermissions) *grantTokenObjectsBuilder {
	return newGrantTokenObjectsBuilder(b.opts).UUIDsPattern(uuids)
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *grantTokenBuilder) QueryParam(queryParam map[string]string) *grantTokenBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Grant request.
func (b *grantTokenBuilder) Execute() (*PNGrantTokenResponse, StatusResponse, error) {
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

func (b *grantTokenBuilder) SpacesPermissions(spacesPermissions map[SpaceId]SpacePermissions) *grantTokenSumBuilder {
	return newGrantTokenSumBuilder(b.opts).SpacesPermissions(spacesPermissions)
}

func (b *grantTokenBuilder) UsersPermissions(usersPermissions map[UserId]UserPermissions) *grantTokenSumBuilder {
	return newGrantTokenSumBuilder(b.opts).UsersPermissions(usersPermissions)
}

func (b *grantTokenBuilder) SpacePatternsPermissions(spacePatternsPermissions map[SpaceId]SpacePermissions) *grantTokenSumBuilder {
	return newGrantTokenSumBuilder(b.opts).SpacePatternsPermissions(spacePatternsPermissions)
}

func (b *grantTokenBuilder) UserPatternsPermissions(userPatternsPermissions map[UserId]UserPermissions) *grantTokenSumBuilder {
	return newGrantTokenSumBuilder(b.opts).UserPatternsPermissions(userPatternsPermissions)
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

//Channels sets the Channels for the Grant request.
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
func (b *grantTokenSumBuilder) TTL(ttl int) *grantTokenSumBuilder {
	b.opts.TTL = ttl
	b.opts.setTTL = true

	return b
}

func (b *grantTokenSumBuilder) AuthorizedUserId(userId UserId) *grantTokenSumBuilder {
	b.opts.AuthorizedUUID = string(userId)

	return b
}

//SpacesPermissions sets the Spaces for the Grant request.
func (b *grantTokenSumBuilder) SpacesPermissions(spaces map[SpaceId]SpacePermissions) *grantTokenSumBuilder {
	b.opts.Channels = toChannelsPermissionsMap(spaces)

	return b
}

func (b *grantTokenSumBuilder) UsersPermissions(users map[UserId]UserPermissions) *grantTokenSumBuilder {
	b.opts.UUIDs = toUUIDsPermissionsMap(users)

	return b
}

// SpacePatternsPermissions sets the Channels for the Grant request.
func (b *grantTokenSumBuilder) SpacePatternsPermissions(spaces map[SpaceId]SpacePermissions) *grantTokenSumBuilder {
	b.opts.ChannelsPattern = toChannelsPermissionsMap(spaces)

	return b
}

func (b *grantTokenSumBuilder) UserPatternsPermissions(users map[UserId]UserPermissions) *grantTokenSumBuilder {
	b.opts.UUIDsPattern = toUUIDsPermissionsMap(users)

	return b
}

// Meta sets the Meta for the Grant request.
func (b *grantTokenSumBuilder) Meta(meta map[string]interface{}) *grantTokenSumBuilder {
	b.opts.Meta = meta

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *grantTokenSumBuilder) QueryParam(queryParam map[string]string) *grantTokenSumBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Grant request.
func (b *grantTokenSumBuilder) Execute() (*PNGrantTokenResponse, StatusResponse, error) {
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

type grantTokenOpts struct {
	pubnub *PubNub
	ctx    Context

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

func (o *grantTokenOpts) config() Config {
	return *o.pubnub.Config
}

func (o *grantTokenOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *grantTokenOpts) context() Context {
	return o.ctx
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
	o.pubnub.Config.Log.Println(fmt.Sprintf("bmVal: %t %d %d", value, bitmask, bm))
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
				o.pubnub.Config.Log.Println("bmVal ChannelPermissions:", bmVal)
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
				o.pubnub.Config.Log.Println("bmVal GroupPermissions:", bmVal)
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
				o.pubnub.Config.Log.Println("bmVal UUIDPermissions:", bmVal)
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

func (o *grantTokenOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
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

	o.pubnub.Config.Log.Println("permissions: ", permissions)

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
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *grantTokenOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *grantTokenOpts) httpMethod() string {
	return "POST"
}

func (o *grantTokenOpts) isAuthRequired() bool {
	return true
}

func (o *grantTokenOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *grantTokenOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *grantTokenOpts) operationType() OperationType {
	return PNAccessManagerGrantToken
}

func (o *grantTokenOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

func (o *grantTokenOpts) tokenManager() *TokenManager {
	return o.pubnub.tokenManager
}

// PNGrantTokenData is the struct used to decode the server response
type PNGrantTokenData struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

// PNGrantTokenResponse is the struct returned when the Execute function of Grant Token is called.
type PNGrantTokenResponse struct {
	status  int              `json:"status"`
	Data    PNGrantTokenData `json:"data"`
	service string           `json:"service"`
}

func newGrantTokenResponse(jsonBytes []byte, status StatusResponse) (*PNGrantTokenResponse, StatusResponse, error) {
	resp := &PNGrantTokenResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGrantTokenResponse, status, e
	}

	return resp, status, nil
}
