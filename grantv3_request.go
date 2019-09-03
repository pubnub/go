package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pubnub/go/pnerr"
	"io/ioutil"
	"net/http"
	"net/url"
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
func (b *grantBuilder) Channels(channels map[string]ResourcePermissions) *grantBuilder {
	b.opts.Channels = channels

	return b
}

// ChannelGroups sets the ChannelGroups for the Grant request.
func (b *grantBuilder) ChannelGroups(groups map[string]ResourcePermissions) *grantBuilder {
	b.opts.ChannelGroups = groups

	return b
}

// Users sets the Users for the Grant request.
func (b *grantBuilder) Users(users map[string]ResourcePermissions) *grantBuilder {
	b.opts.Users = users

	return b
}

// Patterns sets the Patterns for the Grant request.
// func (b *grantBuilder) Patterns(pattern string, resourceTypes patterns) *grantBuilder {
// 	// b.opts.Patterns = patterns

// 	return b
// }

// Spaces sets the Spaces for the Grant request.
func (b *grantBuilder) Spaces(spaces map[string]ResourcePermissions) *grantBuilder {
	b.opts.Spaces = spaces

	return b
}

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

	return newGrantResponse(rawJSON, status)
}

type grantOpts struct {
	pubnub *PubNub
	ctx    Context

	AuthKeys      []string
	Channels      map[string]ResourcePermissions
	ChannelGroups map[string]ResourcePermissions
	QueryParam    map[string]string
	Meta          map[string]interface{}
	Spaces        map[string]ResourcePermissions
	Users         map[string]ResourcePermissions
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

type permissionsBody struct {
	Resources GrantResources         `json:"resources"`
	Patterns  GrantResources         `json:"patterns"`
	Meta      map[string]interface{} `json:"meta"`
}

type grantBody struct {
	TTL         int             `json:"ttl"`
	Permissions permissionsBody `json:"permissions"`
}

func parseResourcePermissions(resource map[string]ResourcePermissions) map[string]int64 {
	r := make(map[string]int64, len(resource))

	for k, v := range resource {
		bm := int64(0)
		if r := v.Read; r {
			bm |= 1
		}
		if w := v.Write; w {
			bm |= 2
		}
		if m := v.Manage; m {
			bm |= 4
		}
		if c := v.Delete; c {
			bm |= 8
		}
		if d := v.Create; d {
			bm |= 16
		}
		r[k] = bm
	}
	return r
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

	var channels map[string]int64
	var groups map[string]int64
	var users map[string]int64
	var spaces map[string]int64

	if len(o.Channels) > 0 {
		channels = parseResourcePermissions(o.Channels)
	} else {
		channels = make(map[string]int64)
	}

	if len(o.ChannelGroups) > 0 {
		groups = parseResourcePermissions(o.ChannelGroups)
	} else {
		groups = make(map[string]int64)
	}

	if len(o.Users) > 0 {
		users = parseResourcePermissions(o.Users)
	} else {
		users = make(map[string]int64)
	}

	if len(o.Spaces) > 0 {
		spaces = parseResourcePermissions(o.Spaces)
	} else {
		spaces = make(map[string]int64)
	}

	rb := GrantResources{
		Channels: channels,
		Users:    users,
		Groups:   groups,
		Spaces:   spaces,
	}

	meta := o.Meta

	if meta == nil {
		meta = make(map[string]interface{})
	}

	permissions := permissionsBody{
		Resources: rb,
		Patterns: GrantResources{
			Channels: make(map[string]int64),
			Users:    make(map[string]int64),
			Groups:   make(map[string]int64),
			Spaces:   make(map[string]int64),
		},
		Meta: meta,
	}

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

func newGrantResponse(jsonBytes []byte, status StatusResponse) (*PNGrantResponse, StatusResponse, error) {
	resp := &PNGrantResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGrantResponse, status, e
	}

	fmt.Println("resp.Data.Token--->", resp.Data.Token)

	return resp, status, nil
}
