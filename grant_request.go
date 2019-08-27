package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sprucehealth/pubnub-go/pnerr"
)

// PNGrantType grant types
type PNGrantType int

const grantPath = "/v2/auth/grant/sub-key/%s"
const grantV3Path = "/v3/auth/grant/sub-key/%s"
const (
	// PNReadEnabled Read Enabled. Applies to Subscribe, History, Presence, Objects
	PNReadEnabled PNGrantType = 1 + iota
	// PNWriteEnabled Write Enabled. Applies to Publish, Objects
	PNWriteEnabled
	// PNManageEnabled Manage Enabled. Applies to Channel-Groups, Objects
	PNManageEnabled
	// PNDeleteEnabled Delete Enabled. Applies to History, Objects
	PNDeleteEnabled
	// PNCreateEnabled Create Enabled. Applies to Objects
	PNCreateEnabled
)

var emptyGrantResponse *GrantResponse

type grantBuilder struct {
	opts *grantOpts
}

type patternPermissions struct {
}
type patterns struct {
	ChannelsPattern      string
	ChannelGroupsPattern string
	SpacesPattern        string
	UsersPattern         string
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

func (b *grantBuilder) Read(read bool) *grantBuilder {
	b.opts.Read = read

	return b
}

func (b *grantBuilder) Write(write bool) *grantBuilder {
	b.opts.Write = write

	return b
}

func (b *grantBuilder) Manage(manage bool) *grantBuilder {
	b.opts.Manage = manage

	return b
}

func (b *grantBuilder) Delete(del bool) *grantBuilder {
	b.opts.Delete = del

	return b
}

func (b *grantBuilder) Create(create bool) *grantBuilder {
	b.opts.Create = create

	return b
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
func (b *grantBuilder) Channels(channels []string) *grantBuilder {
	b.opts.Channels = channels

	return b
}

// ChannelGroups sets the ChannelGroups for the Grant request.
func (b *grantBuilder) ChannelGroups(groups []string) *grantBuilder {
	b.opts.ChannelGroups = groups

	return b
}

// Users sets the Users for the Grant request.
func (b *grantBuilder) Users(users []string) *grantBuilder {
	b.opts.Users = users

	return b
}

// Patterns sets the Patterns for the Grant request.
func (b *grantBuilder) Patterns(pattern string, resourceTypes patterns) *grantBuilder {
	// b.opts.Patterns = patterns

	return b
}

// Spaces sets the Spaces for the Grant request.
func (b *grantBuilder) Spaces(spaces []string) *grantBuilder {
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
func (b *grantBuilder) Execute() (*GrantResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGrantResponse, status, err
	}

	return newGrantResponse(rawJSON, status)
}

type grantOpts struct {
	pubnub *PubNub
	ctx    Context

	AuthKeys      []string
	Channels      []string
	ChannelGroups []string
	QueryParam    map[string]string
	Meta          map[string]interface{}
	Spaces        []string
	Users         []string

	// Stringified permissions
	// Setting 'true' or 'false' will apply permissions to level
	Read   bool
	Write  bool
	Manage bool
	Delete bool
	Create bool
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

func (o *grantOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Read {
		q.Set("r", "1")
	} else {
		q.Set("r", "0")
	}

	if o.Write {
		q.Set("w", "1")
	} else {
		q.Set("w", "0")
	}

	if o.Manage {
		q.Set("m", "1")
	} else {
		q.Set("m", "0")
	}

	if o.Delete {
		q.Set("d", "1")
	} else {
		q.Set("d", "0")
	}

	if len(o.AuthKeys) > 0 {
		q.Set("auth", strings.Join(o.AuthKeys, ","))
	}

	if len(o.Channels) > 0 {
		q.Set("channel", strings.Join(o.Channels, ","))
	}

	if len(o.ChannelGroups) > 0 {
		q.Set("channel-group", strings.Join(o.ChannelGroups, ","))
	}

	if o.setTTL {
		if o.TTL >= -1 {
			q.Set("ttl", fmt.Sprintf("%d", o.TTL))
		}
	}

	timestamp := time.Now().Unix()
	q.Set("timestamp", strconv.Itoa(int(timestamp)))
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *grantOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *grantOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *grantOpts) httpMethod() string {
	return "GET"
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

// GrantResponse is the struct returned when the Execute function of Grant is called.
type GrantResponse struct {
	Level        string
	SubscribeKey string

	TTL int

	Channels      map[string]*PNPAMEntityData
	ChannelGroups map[string]*PNPAMEntityData

	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
	DeleteEnabled bool
}

// PNPAMEntityData is the struct containing the access details of the channels.
type PNPAMEntityData struct {
	Name          string
	AuthKeys      map[string]*PNAccessManagerKeyData
	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
	DeleteEnabled bool
	TTL           int
}

// PNAccessManagerKeyData is the struct containing the access details of the channel groups.
type PNAccessManagerKeyData struct {
	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
	DeleteEnabled bool
	TTL           int
}

func newGrantResponse(jsonBytes []byte, status StatusResponse) (
	*GrantResponse, StatusResponse, error) {
	resp := &GrantResponse{}
	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyGrantResponse, status, e
	}

	constructedChannels := make(map[string]*PNPAMEntityData)
	constructedGroups := make(map[string]*PNPAMEntityData)

	grantData, _ := value.(map[string]interface{})
	payload := grantData["payload"]
	parsedPayload := payload.(map[string]interface{})
	auths, _ := parsedPayload["auths"].(map[string]interface{})
	ttl, _ := parsedPayload["ttl"].(float64)

	if val, ok := parsedPayload["channel"]; ok {
		channelName := val.(string)
		auths := make(map[string]*PNAccessManagerKeyData)
		channelMap, _ := parsedPayload["auths"].(map[string]interface{})
		entityData := &PNPAMEntityData{
			Name: channelName,
		}

		for key, value := range channelMap {
			auths[key] = createPNAccessManagerKeyData(value, entityData, false)
		}

		entityData.AuthKeys = auths
		entityData.TTL = int(ttl)
		constructedChannels[channelName] = entityData
	}

	if val, ok := parsedPayload["channel-groups"]; ok {
		groupName, _ := val.(string)
		constructedAuthKey := make(map[string]*PNAccessManagerKeyData)
		entityData := &PNPAMEntityData{
			Name: groupName,
		}

		if _, ok := val.(string); ok {
			for authKeyName, value := range auths {
				constructedAuthKey[authKeyName] = createPNAccessManagerKeyData(value, entityData, false)
			}

			entityData.AuthKeys = constructedAuthKey
			constructedGroups[groupName] = entityData
		}

		if groupMap, ok := val.(map[string]interface{}); ok {
			groupName, _ := val.(string)
			constructedAuthKey := make(map[string]*PNAccessManagerKeyData)
			entityData := &PNPAMEntityData{
				Name: groupName,
			}

			for groupName, value := range groupMap {
				valueMap := value.(map[string]interface{})

				if keys, ok := valueMap["auths"]; ok {
					parsedKeys, _ := keys.(map[string]interface{})

					for keyName, value := range parsedKeys {
						constructedAuthKey[keyName] = createPNAccessManagerKeyData(value, entityData, false)
					}
				}

				createPNAccessManagerKeyData(valueMap, entityData, true)

				if val, ok := parsedPayload["ttl"]; ok {
					parsedVal, _ := val.(float64)
					entityData.TTL = int(parsedVal)
				}

				entityData.AuthKeys = constructedAuthKey
				constructedGroups[groupName] = entityData
			}
		}
	}

	if val, ok := parsedPayload["channels"]; ok {
		channelMap, _ := val.(map[string]interface{})

		for channelName, value := range channelMap {
			constructedChannels[channelName] = fetchChannel(channelName,
				value, parsedPayload)
		}
	}

	level, _ := parsedPayload["level"].(string)
	subKey, _ := parsedPayload["subscribe_key"].(string)

	resp.Level = level
	resp.SubscribeKey = subKey
	resp.Channels = constructedChannels
	resp.ChannelGroups = constructedGroups

	if r, ok := parsedPayload["r"]; ok {
		parsedValue, _ := r.(float64)
		if parsedValue == float64(1) {
			resp.ReadEnabled = true
		} else {
			resp.ReadEnabled = false
		}
	}

	if r, ok := parsedPayload["w"]; ok {
		parsedValue, _ := r.(float64)
		if parsedValue == float64(1) {
			resp.WriteEnabled = true
		} else {
			resp.WriteEnabled = false
		}
	}

	if r, ok := parsedPayload["m"]; ok {
		parsedValue, _ := r.(float64)
		if parsedValue == float64(1) {
			resp.ManageEnabled = true
		} else {
			resp.ManageEnabled = false
		}
	}

	if r, ok := parsedPayload["d"]; ok {
		parsedValue, _ := r.(float64)
		if parsedValue == float64(1) {
			resp.DeleteEnabled = true
		} else {
			resp.DeleteEnabled = false
		}
	}

	if r, ok := parsedPayload["ttl"]; ok {
		parsedValue, _ := r.(float64)
		resp.TTL = int(parsedValue)
	}

	return resp, status, nil
}

func fetchChannel(channelName string,
	value interface{}, parsedPayload map[string]interface{}) *PNPAMEntityData {

	auths := make(map[string]*PNAccessManagerKeyData)
	entityData := &PNPAMEntityData{
		Name: channelName,
	}

	valueMap, _ := value.(map[string]interface{})

	if val, ok := valueMap["auths"]; ok {
		parsedValue := val.(map[string]interface{})

		for key, value := range parsedValue {
			auths[key] = createPNAccessManagerKeyData(value, entityData, false)
		}
	}

	createPNAccessManagerKeyData(value, entityData, true)

	if val, ok := parsedPayload["ttl"]; ok {
		parsedVal, _ := val.(float64)
		entityData.TTL = int(parsedVal)
	}

	entityData.AuthKeys = auths

	return entityData
}

func readKeyData(val interface{}, keyData *PNAccessManagerKeyData, entityData *PNPAMEntityData, writeToEntityData bool, grantType PNGrantType) {
	parsedValue, _ := val.(float64)
	readValue := false
	if parsedValue == float64(1) {
		readValue = true
	}
	if writeToEntityData {
		switch grantType {
		case PNReadEnabled:
			entityData.ReadEnabled = readValue
		case PNWriteEnabled:
			entityData.WriteEnabled = readValue
		case PNManageEnabled:
			entityData.ManageEnabled = readValue
		case PNDeleteEnabled:
			entityData.DeleteEnabled = readValue
		}
	} else {
		switch grantType {
		case PNReadEnabled:
			keyData.ReadEnabled = readValue
		case PNWriteEnabled:
			keyData.WriteEnabled = readValue
		case PNManageEnabled:
			keyData.ManageEnabled = readValue
		case PNDeleteEnabled:
			keyData.DeleteEnabled = readValue
		}
	}
}

func createPNAccessManagerKeyData(value interface{}, entityData *PNPAMEntityData, writeToEntityData bool) *PNAccessManagerKeyData {
	valueMap := value.(map[string]interface{})
	keyData := &PNAccessManagerKeyData{}

	if val, ok := valueMap["r"]; ok {
		readKeyData(val, keyData, entityData, writeToEntityData, PNReadEnabled)
	}

	if val, ok := valueMap["w"]; ok {
		readKeyData(val, keyData, entityData, writeToEntityData, PNWriteEnabled)
	}

	if val, ok := valueMap["m"]; ok {
		readKeyData(val, keyData, entityData, writeToEntityData, PNManageEnabled)
	}

	if val, ok := valueMap["d"]; ok {
		readKeyData(val, keyData, entityData, writeToEntityData, PNDeleteEnabled)
	}

	if val, ok := valueMap["ttl"]; ok {
		parsedVal, _ := val.(int)
		entityData.TTL = parsedVal
	}
	return keyData
}
