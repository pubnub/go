package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pubnub/go/v9/pnerr"
)

const grantPath = "/v2/auth/grant/sub-key/%s"

var emptyGrantResponse *GrantResponse

type grantBuilder struct {
	opts *grantOpts
}

func newGrantBuilder(pubnub *PubNub) *grantBuilder {
	return newGrantBuilderWithContext(pubnub, pubnub.ctx)
}

func newGrantBuilderWithContext(pubnub *PubNub, context Context) *grantBuilder {
	builder := grantBuilder{
		opts: newGrantOpts(
			pubnub,
			context,
		),
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

func (b *grantBuilder) Get(get bool) *grantBuilder {
	b.opts.Get = get
	b.opts.isGetSet = true
	return b
}

func (b *grantBuilder) Update(update bool) *grantBuilder {
	b.opts.Update = update
	b.opts.isUpdateSet = true

	return b
}

func (b *grantBuilder) Join(join bool) *grantBuilder {
	b.opts.Join = join
	b.opts.isJoinSet = true

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

// ChannelGroups sets the ChannelGroups for the Grant request.
func (b *grantBuilder) UUIDs(targetUUIDs []string) *grantBuilder {
	b.opts.UUIDs = targetUUIDs

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

// GetLogParams returns the user-provided parameters for logging
func (o *grantOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"AuthKeys":      o.AuthKeys,
		"Channels":      o.Channels,
		"ChannelGroups": o.ChannelGroups,
		"UUIDs":         o.UUIDs,
		"Read":          o.Read,
		"Write":         o.Write,
		"Manage":        o.Manage,
		"Delete":        o.Delete,
		"Get":           o.Get,
		"Update":        o.Update,
		"Join":          o.Join,
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
func (b *grantBuilder) Execute() (*GrantResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNAccessManagerGrant, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGrantResponse, status, err
	}

	return newGrantResponse(rawJSON, status)
}

func newGrantOpts(pubnub *PubNub, ctx Context) *grantOpts {
	return &grantOpts{
		endpointOpts: endpointOpts{
			pubnub: pubnub,
			ctx:    ctx,
		},
	}
}

type grantOpts struct {
	endpointOpts

	AuthKeys      []string
	Channels      []string
	ChannelGroups []string
	UUIDs         []string
	QueryParam    map[string]string
	Meta          map[string]interface{}

	// Stringified permissions
	// Setting 'true' or 'false' will apply permissions to level
	Read   bool
	Write  bool
	Manage bool
	Delete bool
	Get    bool
	Update bool
	Join   bool
	// Max: 525600
	// Min: 1
	// Default: 1440
	// Setting 0 will apply the grant indefinitely
	TTL int

	// nil hacks
	setTTL      bool
	isGetSet    bool
	isUpdateSet bool
	isJoinSet   bool
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

	if o.isGetSet {
		if o.Get {
			q.Set("g", "1")
		} else {
			q.Set("g", "0")
		}
	}

	if o.isUpdateSet {
		if o.Update {
			q.Set("u", "1")
		} else {
			q.Set("u", "0")
		}
	}

	if o.isJoinSet {
		if o.Join {
			q.Set("j", "1")
		} else {
			q.Set("j", "0")
		}
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

	if len(o.UUIDs) > 0 {
		q.Set("target-uuid", strings.Join(o.UUIDs, ","))
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

func (o *grantOpts) operationType() OperationType {
	return PNAccessManagerGrant
}

// GrantResponse is the struct returned when the Execute function of Grant is called.
type GrantResponse struct {
	Level        string
	SubscribeKey string

	TTL int

	Channels      map[string]*PNPAMEntityData
	ChannelGroups map[string]*PNPAMEntityData
	UUIDs         map[string]*PNPAMEntityData

	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
	DeleteEnabled bool
	GetEnabled    bool
	UpdateEnabled bool
	JoinEnabled   bool
}

func newGrantResponse(jsonBytes []byte, status StatusResponse) (
	*GrantResponse, StatusResponse, error) {
	resp := &GrantResponse{}
	var value interface{}
	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyGrantResponse, status, e
	}

	constructedChannels := make(map[string]*PNPAMEntityData)
	constructedGroups := make(map[string]*PNPAMEntityData)
	constructedUUIDs := make(map[string]*PNPAMEntityData)

	grantData, ok := value.(map[string]interface{})
	if !ok {
		e := newGrantResponseParsingError(jsonBytes, "invalid JSON structure",
			fmt.Errorf("expected map[string]interface{}, got %T", value))
		return emptyGrantResponse, status, e
	}

	payload, ok := grantData["payload"]
	if !ok {
		e := newGrantResponseParsingError(jsonBytes, "missing payload field",
			fmt.Errorf("payload field not found in response"))
		return emptyGrantResponse, status, e
	}

	if payload == nil {
		e := newGrantResponseParsingError(jsonBytes, "null payload",
			fmt.Errorf("payload field is null"))
		return emptyGrantResponse, status, e
	}

	parsedPayload, ok := payload.(map[string]interface{})
	if !ok {
		e := newGrantResponseParsingError(jsonBytes, "invalid payload structure",
			fmt.Errorf("expected map[string]interface{} for payload, got %T", payload))
		return emptyGrantResponse, status, e
	}
	rootAuths, _, err := grantOptionalMap(parsedPayload, "auths")
	if err != nil {
		return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid auths field", err)
	}
	ttl, _, err := grantOptionalNumber(parsedPayload, "ttl")
	if err != nil {
		return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid ttl field", err)
	}

	if val, ok := parsedPayload["channel"]; ok {
		channelName, ok := val.(string)
		if !ok {
			return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid channel field",
				fmt.Errorf("expected string for channel, got %T", val))
		}
		constructedAuthKeys := make(map[string]*PNAccessManagerKeyData)
		entityData := &PNPAMEntityData{
			Name: channelName,
		}

		for key, value := range rootAuths {
			keyData, err := createPNAccessManagerKeyData(value, entityData, false)
			if err != nil {
				return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, fmt.Sprintf("invalid auth key %q", key), err)
			}
			constructedAuthKeys[key] = keyData
		}

		entityData.AuthKeys = constructedAuthKeys
		entityData.TTL = int(ttl)
		constructedChannels[channelName] = entityData
	}

	if val, ok := parsedPayload["channel-groups"]; ok {
		constructedAuthKey := make(map[string]*PNAccessManagerKeyData)
		entityData := &PNPAMEntityData{
			Name: "",
		}

		if groupName, ok := val.(string); ok {
			entityData.Name = groupName
			for authKeyName, value := range rootAuths {
				keyData, err := createPNAccessManagerKeyData(value, entityData, false)
				if err != nil {
					return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, fmt.Sprintf("invalid auth key %q", authKeyName), err)
				}
				constructedAuthKey[authKeyName] = keyData
			}

			entityData.AuthKeys = constructedAuthKey
			constructedGroups[groupName] = entityData
		} else if groupMap, ok := val.(map[string]interface{}); ok {
			for groupName, value := range groupMap {
				valueMap, ok := value.(map[string]interface{})
				if !ok {
					return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, fmt.Sprintf("invalid channel group %q", groupName),
						fmt.Errorf("expected map[string]interface{}, got %T", value))
				}

				if keys, ok := valueMap["auths"]; ok {
					parsedKeys, ok := keys.(map[string]interface{})
					if !ok {
						return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, fmt.Sprintf("invalid auths for channel group %q", groupName),
							fmt.Errorf("expected map[string]interface{}, got %T", keys))
					}

					for keyName, value := range parsedKeys {
						keyData, err := createPNAccessManagerKeyData(value, entityData, false)
						if err != nil {
							return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, fmt.Sprintf("invalid auth key %q", keyName), err)
						}
						constructedAuthKey[keyName] = keyData
					}
				}

				if _, err := createPNAccessManagerKeyData(valueMap, entityData, true); err != nil {
					return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, fmt.Sprintf("invalid channel group %q permissions", groupName), err)
				}

				if val, ok := parsedPayload["ttl"]; ok {
					parsedVal, ok := grantNumber(val)
					if !ok {
						return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid ttl field",
							fmt.Errorf("expected number for ttl, got %T", val))
					}
					entityData.TTL = int(parsedVal)
				}

				entityData.AuthKeys = constructedAuthKey
				constructedGroups[groupName] = entityData
			}
		} else {
			return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid channel-groups field",
				fmt.Errorf("expected string or map[string]interface{}, got %T", val))
		}
	}

	if val, ok := parsedPayload["channels"]; ok {
		channelMap, ok := val.(map[string]interface{})
		if !ok {
			return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid channels field",
				fmt.Errorf("expected map[string]interface{}, got %T", val))
		}

		for channelName, value := range channelMap {
			channelData, err := fetchChannel(channelName, value, parsedPayload)
			if err != nil {
				return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, fmt.Sprintf("invalid channel %q", channelName), err)
			}
			constructedChannels[channelName] = channelData
		}
	}

	if val, ok := parsedPayload["uuids"]; ok {
		uuids, ok := val.(map[string]interface{})
		if !ok {
			return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid uuids field",
				fmt.Errorf("expected map[string]interface{}, got %T", val))
		}

		for uuid, value := range uuids {
			uuidData, err := fetchChannel(uuid, value, parsedPayload)
			if err != nil {
				return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, fmt.Sprintf("invalid uuid %q", uuid), err)
			}
			constructedUUIDs[uuid] = uuidData
		}
	}

	level, _ := parsedPayload["level"].(string)
	subKey, _ := parsedPayload["subscribe_key"].(string)

	resp.Level = level
	resp.SubscribeKey = subKey
	resp.Channels = constructedChannels
	resp.ChannelGroups = constructedGroups
	resp.UUIDs = constructedUUIDs

	if resp.ReadEnabled, err = parsePerms(parsedPayload, "r"); err != nil {
		return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid r permission", err)
	}
	if resp.WriteEnabled, err = parsePerms(parsedPayload, "w"); err != nil {
		return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid w permission", err)
	}
	if resp.ManageEnabled, err = parsePerms(parsedPayload, "m"); err != nil {
		return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid m permission", err)
	}
	if resp.DeleteEnabled, err = parsePerms(parsedPayload, "d"); err != nil {
		return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid d permission", err)
	}
	if resp.GetEnabled, err = parsePerms(parsedPayload, "g"); err != nil {
		return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid g permission", err)
	}
	if resp.UpdateEnabled, err = parsePerms(parsedPayload, "u"); err != nil {
		return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid u permission", err)
	}
	if resp.JoinEnabled, err = parsePerms(parsedPayload, "j"); err != nil {
		return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid j permission", err)
	}

	if r, ok := parsedPayload["ttl"]; ok {
		parsedValue, ok := grantNumber(r)
		if !ok {
			return emptyGrantResponse, status, newGrantResponseParsingError(jsonBytes, "invalid ttl field",
				fmt.Errorf("expected number for ttl, got %T", r))
		}
		resp.TTL = int(parsedValue)
	}

	return resp, status, nil
}

func newGrantResponseParsingError(jsonBytes []byte, message string, err error) error {
	return pnerr.NewResponseParsingError("Error parsing response: "+message,
		io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)
}

func grantNumber(value interface{}) (float64, bool) {
	switch typedValue := value.(type) {
	case float64:
		return typedValue, true
	case int:
		return float64(typedValue), true
	default:
		return 0, false
	}
}

func grantOptionalNumber(parsedPayload map[string]interface{}, name string) (float64, bool, error) {
	value, ok := parsedPayload[name]
	if !ok {
		return 0, false, nil
	}

	parsedValue, ok := grantNumber(value)
	if !ok {
		return 0, true, fmt.Errorf("expected number for %s, got %T", name, value)
	}

	return parsedValue, true, nil
}

func grantOptionalMap(parsedPayload map[string]interface{}, name string) (map[string]interface{}, bool, error) {
	value, ok := parsedPayload[name]
	if !ok {
		return nil, false, nil
	}

	parsedValue, ok := value.(map[string]interface{})
	if !ok {
		return nil, true, fmt.Errorf("expected map[string]interface{} for %s, got %T", name, value)
	}

	return parsedValue, true, nil
}

func parsePerms(parsedPayload map[string]interface{}, name string) (bool, error) {
	if r, ok := parsedPayload[name]; ok {
		parsedValue, ok := grantNumber(r)
		if !ok {
			return false, fmt.Errorf("expected number for %s, got %T", name, r)
		}
		if parsedValue == float64(1) {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, nil
}

func fetchChannel(channelName string, value interface{}, parsedPayload map[string]interface{}) (*PNPAMEntityData, error) {

	auths := make(map[string]*PNAccessManagerKeyData)
	entityData := &PNPAMEntityData{
		Name: channelName,
	}

	valueMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map[string]interface{}, got %T", value)
	}

	if val, ok := valueMap["auths"]; ok {
		parsedValue, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected map[string]interface{} for auths, got %T", val)
		}

		for key, value := range parsedValue {
			keyData, err := createPNAccessManagerKeyData(value, entityData, false)
			if err != nil {
				return nil, fmt.Errorf("invalid auth key %q: %w", key, err)
			}
			auths[key] = keyData
		}
	}

	if _, err := createPNAccessManagerKeyData(value, entityData, true); err != nil {
		return nil, err
	}

	if val, ok := parsedPayload["ttl"]; ok {
		parsedVal, ok := grantNumber(val)
		if !ok {
			return nil, fmt.Errorf("expected number for ttl, got %T", val)
		}
		entityData.TTL = int(parsedVal)
	}

	entityData.AuthKeys = auths

	return entityData, nil
}

func readKeyData(val interface{}, keyData *PNAccessManagerKeyData, entityData *PNPAMEntityData, writeToEntityData bool, grantType PNGrantType) error {
	parsedValue, ok := grantNumber(val)
	if !ok {
		return fmt.Errorf("expected number for permission, got %T", val)
	}
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
		case PNGetEnabled:
			entityData.GetEnabled = readValue
		case PNUpdateEnabled:
			entityData.UpdateEnabled = readValue
		case PNJoinEnabled:
			entityData.JoinEnabled = readValue
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
		case PNGetEnabled:
			keyData.GetEnabled = readValue
		case PNUpdateEnabled:
			keyData.UpdateEnabled = readValue
		case PNJoinEnabled:
			keyData.JoinEnabled = readValue
		}
	}
	return nil
}

func createPNAccessManagerKeyData(value interface{}, entityData *PNPAMEntityData, writeToEntityData bool) (*PNAccessManagerKeyData, error) {
	valueMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map[string]interface{}, got %T", value)
	}
	keyData := &PNAccessManagerKeyData{}

	if val, ok := valueMap["r"]; ok {
		if err := readKeyData(val, keyData, entityData, writeToEntityData, PNReadEnabled); err != nil {
			return nil, fmt.Errorf("invalid r permission: %w", err)
		}
	}

	if val, ok := valueMap["w"]; ok {
		if err := readKeyData(val, keyData, entityData, writeToEntityData, PNWriteEnabled); err != nil {
			return nil, fmt.Errorf("invalid w permission: %w", err)
		}
	}

	if val, ok := valueMap["m"]; ok {
		if err := readKeyData(val, keyData, entityData, writeToEntityData, PNManageEnabled); err != nil {
			return nil, fmt.Errorf("invalid m permission: %w", err)
		}
	}

	if val, ok := valueMap["d"]; ok {
		if err := readKeyData(val, keyData, entityData, writeToEntityData, PNDeleteEnabled); err != nil {
			return nil, fmt.Errorf("invalid d permission: %w", err)
		}
	}

	if val, ok := valueMap["g"]; ok {
		if err := readKeyData(val, keyData, entityData, writeToEntityData, PNGetEnabled); err != nil {
			return nil, fmt.Errorf("invalid g permission: %w", err)
		}
	}

	if val, ok := valueMap["u"]; ok {
		if err := readKeyData(val, keyData, entityData, writeToEntityData, PNUpdateEnabled); err != nil {
			return nil, fmt.Errorf("invalid u permission: %w", err)
		}
	}

	if val, ok := valueMap["j"]; ok {
		if err := readKeyData(val, keyData, entityData, writeToEntityData, PNJoinEnabled); err != nil {
			return nil, fmt.Errorf("invalid j permission: %w", err)
		}
	}

	if val, ok := valueMap["ttl"]; ok {
		parsedVal, ok := grantNumber(val)
		if !ok {
			return nil, fmt.Errorf("expected number for ttl, got %T", val)
		}
		entityData.TTL = int(parsedVal)
	}
	return keyData, nil
}
