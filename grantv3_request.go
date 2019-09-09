package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pubnub/go/pnerr"
	"io/ioutil"
	"net/http"
	"net/url"
	//	"regexp"
)

const grantPath = "/v3/pam/%s/grant"

var emptyPNGrantResponse *PNGrantResponse

type grantBuilder struct {
	opts *grantOpts
}

func newGrantBuilder(pubnub *PubNub) *grantBuilder {
	builder := grantBuilder{
		opts: &grantOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newGrantBuilderWithContext(pubnub *PubNub, context Context) *grantBuilder {
	builder := grantBuilder{
		opts: &grantOpts{
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
func (b *grantBuilder) TTL(ttl int) *grantBuilder {
	b.opts.TTL = ttl
	b.opts.setTTL = true

	return b
}

// AuthKeys sets the AuthKeys for the Grant request.
func (b *grantBuilder) AuthKeys(authKeys []string) *grantBuilder {
	b.opts.AuthKeys = authKeys

	return b
}

// Channels sets the Channels for the Grant request.
func (b *grantBuilder) Channels(channels map[string]ChannelPermissions) *grantBuilder {
	b.opts.Channels = channels

	return b
}

// ChannelGroups sets the ChannelGroups for the Grant request.
func (b *grantBuilder) ChannelGroups(groups map[string]GroupPermissions) *grantBuilder {
	b.opts.ChannelGroups = groups

	return b
}

// Users sets the Users for the Grant request.
func (b *grantBuilder) Users(users map[string]UserSpacePermissions) *grantBuilder {
	b.opts.Users = users

	return b
}

// Spaces sets the Spaces for the Grant request.
func (b *grantBuilder) Spaces(spaces map[string]UserSpacePermissions) *grantBuilder {
	b.opts.Spaces = spaces

	return b
}

// Channels sets the Channels for the Grant request.
func (b *grantBuilder) ChannelsPattern(channels map[string]ChannelPermissions) *grantBuilder {
	b.opts.ChannelsPattern = channels

	return b
}

// ChannelGroups sets the ChannelGroups for the Grant request.
func (b *grantBuilder) ChannelGroupsPattern(groups map[string]GroupPermissions) *grantBuilder {
	b.opts.ChannelGroupsPattern = groups

	return b
}

// Users sets the Users for the Grant request.
func (b *grantBuilder) UsersPattern(users map[string]UserSpacePermissions) *grantBuilder {
	b.opts.UsersPattern = users

	return b
}

// Spaces sets the Spaces for the Grant request.
func (b *grantBuilder) SpacesPattern(spaces map[string]UserSpacePermissions) *grantBuilder {
	b.opts.SpacesPattern = spaces

	return b
}

//Patterns sets the Patterns for the Grant request.
// func (b *grantBuilder) Patterns(pattern string, resourceTypes patterns) *grantBuilder {
// 	// b.opts.Patterns = patterns

// 	return b
// }

// Meta sets the Meta for the Grant request.
func (b *grantBuilder) Meta(meta map[string]interface{}) *grantBuilder {
	b.opts.Meta = meta

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *grantBuilder) QueryParam(queryParam map[string]string) *grantBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Grant request.
func (b *grantBuilder) Execute() (*PNGrantResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNGrantResponse, status, err
	}

	return newGrantResponse(b, rawJSON, status)
}

type grantOpts struct {
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

func (o *grantOpts) config() Config {
	return *o.pubnub.Config
}

func (o *grantOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *grantOpts) context() Context {
	return o.ctx
}

func (o *grantOpts) validate() error {
	if o.config().PublishKey == "" {
		return newValidationError(o, StrMissingPubKey)
	}

	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.config().SecretKey == "" {
		return newValidationError(o, StrMissingSecretKey)
	}

	return nil
}

func (o *grantOpts) buildPath() (string, error) {
	return fmt.Sprintf(grantPath, o.pubnub.Config.SubscribeKey), nil
}

type grantBody struct {
	TTL         int             `json:"ttl"`
	Permissions PermissionsBody `json:"permissions"`
}

func setBitmask(value bool, bitmask PNGrantBitMask, bm int64) int64 {
	if value {
		bm |= int64(bitmask)
	}
	//fmt.Println("====>", bm)
	return bm
}

func parseResourcePermissions(resource interface{}, resourceType PNResourceType) map[string]int64 {
	bmVal := int64(0)
	switch resourceType {
	case PNChannels:
		resourceWithPerms := resource.(map[string]ChannelPermissions)
		resourceWithPermsLen := len(resourceWithPerms)
		if resourceWithPermsLen > 0 {
			//fmt.Println(c)
			r := make(map[string]int64, resourceWithPermsLen)
			for k, v := range resourceWithPerms {
				// _, err := regexp.Compile(k)
				// if err != nil {
				// 	fmt.Println(err.Error())
				// } else {
				// 	fmt.Println("Regex compiled", k)
				// }
				bmVal = int64(0)
				bmVal = setBitmask(v.Read, PNRead, bmVal)
				bmVal = setBitmask(v.Write, PNWrite, bmVal)
				bmVal = setBitmask(v.Delete, PNDelete, bmVal)
				//fmt.Println("bmVal====>", bmVal)
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
				bmVal = setBitmask(v.Read, PNRead, bmVal)
				bmVal = setBitmask(v.Manage, PNManage, bmVal)
				//fmt.Println("bmVal====>", bmVal)
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
			//fmt.Println(c)
			for k, v := range resourceWithPerms {
				bmVal = int64(0)
				bmVal = setBitmask(v.Read, PNRead, bmVal)
				bmVal = setBitmask(v.Write, PNWrite, bmVal)
				bmVal = setBitmask(v.Manage, PNManage, bmVal)
				bmVal = setBitmask(v.Delete, PNDelete, bmVal)
				bmVal = setBitmask(v.Create, PNCreate, bmVal)
				//fmt.Println("bmVal====>", bmVal)
				r[k] = bmVal
			}
			return r
		}
		return make(map[string]int64)
	}

}

func (o *grantOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *grantOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *grantOpts) buildBody() ([]byte, error) {

	meta := o.Meta

	if meta == nil {
		meta = make(map[string]interface{})
	}

	permissions := PermissionsBody{
		Resources: GrantResources{
			Channels: parseResourcePermissions(o.Channels, PNChannels),
			Groups:   parseResourcePermissions(o.ChannelGroups, PNGroups),
			Users:    parseResourcePermissions(o.Users, PNUsers),
			Spaces:   parseResourcePermissions(o.Spaces, PNSpaces),
		},
		Patterns: GrantResources{
			Channels: parseResourcePermissions(o.ChannelsPattern, PNChannels),
			Groups:   parseResourcePermissions(o.ChannelGroupsPattern, PNGroups),
			Users:    parseResourcePermissions(o.UsersPattern, PNUsers),
			Spaces:   parseResourcePermissions(o.SpacesPattern, PNSpaces),
		},
		Meta: meta,
	}

	fmt.Println("permissions", permissions)

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

func (o *grantOpts) httpMethod() string {
	return "POST"
}

func (o *grantOpts) isAuthRequired() bool {
	return true
}

func (o *grantOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *grantOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *grantOpts) operationType() OperationType {
	return PNAccessManagerGrant
}

func (o *grantOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

type PNGrantData struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

// GrantResponse is the struct returned when the Execute function of Grant is called.
type PNGrantResponse struct {
	status  int         `json:"status"`
	Data    PNGrantData `json:"data"`
	service string      `json:"service"`
}

func newGrantResponse(b *grantBuilder, jsonBytes []byte, status StatusResponse) (*PNGrantResponse, StatusResponse, error) {
	resp := &PNGrantResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGrantResponse, status, e
	}

	b.opts.pubnub.tokenManager.StoreToken(resp.Data.Token)

	return resp, status, nil
}
