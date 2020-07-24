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

	"github.com/pubnub/go/pnerr"
)

const grantTokenPath = "/v3/pam/%s/grant"

var emptyPNGrantTokenResponse *PNGrantTokenResponse

type grantTokenBuilder struct {
	opts *grantTokenOpts
}

func newGrantTokenBuilder(pubnub *PubNub) *grantTokenBuilder {
	builder := grantTokenBuilder{
		opts: &grantTokenOpts{
			pubnub: pubnub,
		},
	}

	return &builder
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

// Uncomment when PAMv3 is fully functional.
// Channels sets the Channels for the Grant request.
// func (b *grantTokenBuilder) Channels(channels map[string]ChannelPermissions) *grantTokenBuilder {
// 	b.opts.Channels = channels

// 	return b
// }

// // ChannelGroups sets the ChannelGroups for the Grant request.
// func (b *grantTokenBuilder) ChannelGroups(groups map[string]GroupPermissions) *grantTokenBuilder {
// 	b.opts.ChannelGroups = groups

// 	return b
// }

// Users sets the Users for the Grant request.
func (b *grantTokenBuilder) Users(users map[string]UserSpacePermissions) *grantTokenBuilder {
	b.opts.Users = users

	return b
}

// Spaces sets the Spaces for the Grant request.
func (b *grantTokenBuilder) Spaces(spaces map[string]UserSpacePermissions) *grantTokenBuilder {
	b.opts.Spaces = spaces

	return b
}

// Uncomment when PAMv3 is fully functional.
// // Channels sets the Channels for the Grant request.
// func (b *grantTokenBuilder) ChannelsPattern(channels map[string]ChannelPermissions) *grantTokenBuilder {
// 	b.opts.ChannelsPattern = channels

// 	return b
// }

// // ChannelGroups sets the ChannelGroups for the Grant request.
// func (b *grantTokenBuilder) ChannelGroupsPattern(groups map[string]GroupPermissions) *grantTokenBuilder {
// 	b.opts.ChannelGroupsPattern = groups

// 	return b
// }

// Users sets the Users for the Grant request.
func (b *grantTokenBuilder) UsersPattern(users map[string]UserSpacePermissions) *grantTokenBuilder {
	b.opts.UsersPattern = users

	return b
}

// Spaces sets the Spaces for the Grant request.
func (b *grantTokenBuilder) SpacesPattern(spaces map[string]UserSpacePermissions) *grantTokenBuilder {
	b.opts.SpacesPattern = spaces

	return b
}

// Meta sets the Meta for the Grant request.
func (b *grantTokenBuilder) Meta(meta map[string]interface{}) *grantTokenBuilder {
	b.opts.Meta = meta

	return b
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

	return newGrantTokenResponse(b, rawJSON, status)
}

type grantTokenOpts struct {
	pubnub *PubNub
	ctx    Context

	AuthKeys             []string
	Channels             map[string]ChannelPermissions
	ChannelGroups        map[string]GroupPermissions
	Spaces               map[string]UserSpacePermissions
	Users                map[string]UserSpacePermissions
	ChannelsPattern      map[string]ChannelPermissions
	ChannelGroupsPattern map[string]GroupPermissions
	SpacesPattern        map[string]UserSpacePermissions
	UsersPattern         map[string]UserSpacePermissions
	QueryParam           map[string]string
	Meta                 map[string]interface{}

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

	default:
		//case PNUsers:
		//case PNSpaces:
		resourceWithPerms := resource.(map[string]UserSpacePermissions)
		resourceWithPermsLen := len(resourceWithPerms)
		if resourceWithPermsLen > 0 {
			r := make(map[string]int64, resourceWithPermsLen)
			for k, v := range resourceWithPerms {
				bmVal = int64(0)
				bmVal = o.setBitmask(v.Read, PNRead, bmVal)
				bmVal = o.setBitmask(v.Write, PNWrite, bmVal)
				bmVal = o.setBitmask(v.Manage, PNManage, bmVal)
				bmVal = o.setBitmask(v.Delete, PNDelete, bmVal)
				bmVal = o.setBitmask(v.Create, PNCreate, bmVal)
				o.pubnub.Config.Log.Println("bmVal UserSpacePermissions:", bmVal)
				r[k] = bmVal
			}
			return r
		}
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
			Users:    o.parseResourcePermissions(o.Users, PNUsers),
			Spaces:   o.parseResourcePermissions(o.Spaces, PNSpaces),
		},
		Patterns: GrantResources{
			Channels: o.parseResourcePermissions(o.ChannelsPattern, PNChannels),
			Groups:   o.parseResourcePermissions(o.ChannelGroupsPattern, PNGroups),
			Users:    o.parseResourcePermissions(o.UsersPattern, PNUsers),
			Spaces:   o.parseResourcePermissions(o.SpacesPattern, PNSpaces),
		},
		Meta: meta,
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

func newGrantTokenResponse(b *grantTokenBuilder, jsonBytes []byte, status StatusResponse) (*PNGrantTokenResponse, StatusResponse, error) {
	resp := &PNGrantTokenResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGrantTokenResponse, status, e
	}

	b.opts.pubnub.tokenManager.StoreToken(resp.Data.Token)

	return resp, status, nil
}
