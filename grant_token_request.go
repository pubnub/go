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

// AuthKeys sets the AuthKeys for the Grant request.
func (b *grantTokenBuilder) AuthKeys(authKeys []string) *grantTokenBuilder {
	b.opts.AuthKeys = authKeys

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

//Patterns sets the Patterns for the Grant request.
// func (b *grantTokenBuilder) Patterns(pattern string, resourceTypes patterns) *grantTokenBuilder {
// 	// b.opts.Patterns = patterns

// 	return b
// }

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

	return nil
}

func (o *grantTokenOpts) buildPath() (string, error) {
	return fmt.Sprintf(grantTokenPath, o.pubnub.Config.SubscribeKey), nil
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

	//fmt.Println("permissions", permissions)

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

type PNGrantTokenData struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

// GrantResponse is the struct returned when the Execute function of Grant is called.
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
