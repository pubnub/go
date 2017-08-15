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

	"github.com/pubnub/go/pnerr"
)

const GRANT_PATH = "/v1/auth/grant/sub-key/%s"

var emptyGrantResponse *GrantResponse

func GrantRequestWithContext(ctx Context, pn *PubNub, opts *GrantOpts) (
	*GrantResponse, error) {
	opts.pubnub = pn
	opts.ctx = ctx

	_, err := executeRequest(opts)
	if err != nil {
		return emptyGrantResponse, err
	}

	return emptyGrantResponse, nil
}

type grantBuilder struct {
	opts *GrantOpts
}

func newGrantBuilder(pubnub *PubNub) *grantBuilder {
	builder := grantBuilder{
		opts: &GrantOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func (b *grantBuilder) Read(read bool) *grantBuilder {
	b.opts.Read = read
	b.opts.SetRead = true

	return b
}

func (b *grantBuilder) Write(write bool) *grantBuilder {
	b.opts.Write = write
	b.opts.SetWrite = true

	return b
}

func (b *grantBuilder) Manage(manage bool) *grantBuilder {
	b.opts.Manage = manage
	b.opts.SetManage = true

	return b
}

func (b *grantBuilder) Ttl(ttl int) *grantBuilder {
	b.opts.Ttl = ttl
	b.opts.SetTtl = true

	return b
}

func (b *grantBuilder) AuthKeys(authKeys []string) *grantBuilder {
	b.opts.AuthKeys = authKeys

	return b
}

func (b *grantBuilder) Channels(channels []string) *grantBuilder {
	b.opts.Channels = channels

	return b
}

func (b *grantBuilder) Groups(groups []string) *grantBuilder {
	b.opts.Groups = groups

	return b
}

func (b *grantBuilder) Execute() (*GrantResponse, error) {
	rawJson, err := executeRequest(b.opts)
	if err != nil {
		return emptyGrantResponse, err
	}

	return newGrantResponse(rawJson)
}

// TODO: make private
type GrantOpts struct {
	pubnub *PubNub
	ctx    Context

	AuthKeys []string
	Channels []string
	Groups   []string

	// Stringified permissions
	// Setting 'true' or 'false' will apply permissions to level
	Read   bool
	Write  bool
	Manage bool

	// Max: 525600
	// Min: 1
	// Default: 1440
	// Setting 0 will apply the grant indefinitely
	Ttl int

	// nil hacks
	SetRead   bool
	SetWrite  bool
	SetManage bool
	SetTtl    bool
}

func (o *GrantOpts) config() Config {
	return *o.pubnub.Config
}

func (o *GrantOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *GrantOpts) context() Context {
	return o.ctx
}

func (o *GrantOpts) validate() error {
	if o.config().PublishKey == "" {
		return ErrMissingPubKey
	}

	if o.config().SubscribeKey == "" {
		return ErrMissingSubKey
	}

	if o.config().SecretKey == "" {
		return ErrMissingSecretKey
	}

	return nil
}

func (o *GrantOpts) buildPath() (string, error) {
	return fmt.Sprintf(GRANT_PATH, o.pubnub.Config.SubscribeKey), nil
}

func (o *GrantOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	if o.SetRead {
		if o.Read {
			q.Set("r", "1")
		} else {
			q.Set("r", "0")
		}
	}

	if o.SetWrite {
		if o.Write {
			q.Set("w", "1")
		} else {
			q.Set("w", "0")
		}
	}

	if o.SetManage {
		if o.Manage {
			q.Set("m", "1")
		} else {
			q.Set("m", "0")
		}
	}

	if len(o.AuthKeys) > 0 {
		q.Set("auth", strings.Join(o.AuthKeys, ","))
	}

	if len(o.Channels) > 0 {
		q.Set("channel", strings.Join(o.Channels, ","))
	}

	if len(o.Groups) > 0 {
		q.Set("channel-group", strings.Join(o.Groups, ","))
	}

	if o.SetTtl {
		if o.Ttl >= -1 {
			q.Set("ttl", fmt.Sprintf("%d", o.Ttl))
		}
	}

	timestamp := time.Now().Unix()
	q.Set("timestamp", strconv.Itoa(int(timestamp)))

	return q, nil
}

func (o *GrantOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *GrantOpts) httpMethod() string {
	return "GET"
}

func (o *GrantOpts) isAuthRequired() bool {
	return true
}

func (o *GrantOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *GrantOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *GrantOpts) operationType() PNOperationType {
	return PNAccessManagerGrant
}

type GrantResponse struct {
	Level        string
	SubscribeKey string

	Ttl int

	Channels      map[string]*PNPAMEntityData
	ChannelGroups map[string]*PNPAMEntityData

	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
}

type PNPAMEntityData struct {
	Name          string
	AuthKeys      map[string]*PNAccessManagerKeyData
	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
	Ttl           int
}

type PNAccessManagerKeyData struct {
	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
	Ttl           int
}

func newGrantResponse(jsonBytes []byte) (*GrantResponse, error) {
	resp := &GrantResponse{}
	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyGrantResponse, e
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
			valueMap := value.(map[string]interface{})
			keyData := &PNAccessManagerKeyData{}

			if val, ok := valueMap["r"]; ok {
				parsedValue, _ := val.(float64)
				if parsedValue == float64(1) {
					keyData.ReadEnabled = true
				} else {
					keyData.ReadEnabled = false
				}
			}

			if val, ok := valueMap["w"]; ok {
				parsedValue, _ := val.(float64)
				if parsedValue == float64(1) {
					keyData.WriteEnabled = true
				} else {
					keyData.WriteEnabled = false
				}
			}

			if val, ok := valueMap["m"]; ok {
				parsedValue, _ := val.(float64)
				if parsedValue == float64(1) {
					keyData.ManageEnabled = true
				} else {
					keyData.ManageEnabled = false
				}
			}

			auths[key] = keyData
		}

		entityData.AuthKeys = auths
		entityData.Ttl = int(ttl)
		constructedChannels[channelName] = entityData
	}

	if val, ok := parsedPayload["channel-groups"]; ok {
		groupName, _ := val.(string)
		constructedAuthKey := make(map[string]*PNAccessManagerKeyData)
		entityData := PNPAMEntityData{
			Name: groupName,
		}

		if _, ok := val.(string); ok {
			for authKeyName, value := range auths {
				auth, _ := value.(map[string]interface{})

				managerKeyData := &PNAccessManagerKeyData{}

				if val, ok := auth["r"]; ok {
					parsedValue, _ := val.(float64)
					if parsedValue == float64(1) {
						managerKeyData.ReadEnabled = true
					} else {
						managerKeyData.ReadEnabled = false
					}
				}

				if val, ok := auth["w"]; ok {
					parsedValue, _ := val.(float64)
					if parsedValue == float64(1) {
						managerKeyData.WriteEnabled = true
					} else {
						managerKeyData.WriteEnabled = false
					}
				}

				if val, ok := auth["m"]; ok {
					parsedValue, _ := val.(float64)
					if parsedValue == float64(1) {
						managerKeyData.ManageEnabled = true
					} else {
						managerKeyData.ManageEnabled = false
					}
				}

				if val, ok := auth["ttl"]; ok {
					parsedVal, _ := val.(int)
					entityData.Ttl = parsedVal
				}

				constructedAuthKey[authKeyName] = managerKeyData
			}

			entityData.AuthKeys = constructedAuthKey
			constructedGroups[groupName] = &entityData
		}

		if groupMap, ok := val.(map[string]interface{}); ok {
			groupName, _ := val.(string)
			constructedAuthKey := make(map[string]*PNAccessManagerKeyData)
			entityData := PNPAMEntityData{
				Name: groupName,
			}

			for groupName, value := range groupMap {
				valueMap := value.(map[string]interface{})

				if keys, ok := valueMap["auths"]; ok {
					parsedKeys, _ := keys.(map[string]interface{})
					keyData := &PNAccessManagerKeyData{}

					for keyName, value := range parsedKeys {
						valueMap, _ := value.(map[string]interface{})

						if val, ok := valueMap["r"]; ok {
							parsedValue, _ := val.(float64)
							if parsedValue == float64(1) {
								keyData.ReadEnabled = true
							} else {
								keyData.ReadEnabled = false
							}
						}

						if val, ok := valueMap["w"]; ok {
							parsedValue, _ := val.(float64)
							if parsedValue == float64(1) {
								keyData.WriteEnabled = true
							} else {
								keyData.WriteEnabled = false
							}
						}

						if val, ok := valueMap["m"]; ok {
							parsedValue, _ := val.(float64)
							if parsedValue == float64(1) {
								keyData.ManageEnabled = true
							} else {
								keyData.ManageEnabled = false
							}
						}

						constructedAuthKey[keyName] = keyData
					}
				}

				if val, ok := valueMap["r"]; ok {
					parsedValue, _ := val.(float64)
					if parsedValue == float64(1) {
						entityData.ReadEnabled = true
					} else {
						entityData.ReadEnabled = false
					}
				}

				if val, ok := valueMap["w"]; ok {
					parsedValue, _ := val.(float64)
					if parsedValue == float64(1) {
						entityData.WriteEnabled = true
					} else {
						entityData.WriteEnabled = false
					}
				}

				if val, ok := valueMap["m"]; ok {
					parsedValue, _ := val.(float64)
					if parsedValue == float64(1) {
						entityData.ManageEnabled = true
					} else {
						entityData.ManageEnabled = false
					}
				}

				if val, ok := parsedPayload["ttl"]; ok {
					parsedVal, _ := val.(float64)
					entityData.Ttl = int(parsedVal)
				}

				entityData.AuthKeys = constructedAuthKey
				constructedGroups[groupName] = &entityData
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

	if r, ok := parsedPayload["ttl"]; ok {
		parsedValue, _ := r.(float64)
		resp.Ttl = int(parsedValue)
	}

	return resp, nil
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
			valueMap := value.(map[string]interface{})
			keyData := &PNAccessManagerKeyData{}

			if val, ok := valueMap["r"]; ok {
				parsedValue, _ := val.(float64)
				if parsedValue == float64(1) {
					keyData.ReadEnabled = true
				} else {
					keyData.ReadEnabled = false
				}
			}

			if val, ok := valueMap["w"]; ok {
				parsedValue, _ := val.(float64)
				if parsedValue == float64(1) {
					keyData.WriteEnabled = true
				} else {
					keyData.WriteEnabled = false
				}
			}

			if val, ok := valueMap["m"]; ok {
				parsedValue, _ := val.(float64)
				if parsedValue == float64(1) {
					keyData.ManageEnabled = true
				} else {
					keyData.ManageEnabled = false
				}
			}

			auths[key] = keyData
		}
	}

	if val, ok := valueMap["r"]; ok {
		parsedValue, _ := val.(float64)
		if parsedValue == float64(1) {
			entityData.ReadEnabled = true
		} else {
			entityData.ReadEnabled = false
		}
	}

	if val, ok := valueMap["w"]; ok {
		parsedValue, _ := val.(float64)
		if parsedValue == float64(1) {
			entityData.WriteEnabled = true
		} else {
			entityData.WriteEnabled = false
		}
	}

	if val, ok := valueMap["m"]; ok {
		parsedValue, _ := val.(float64)
		if parsedValue == float64(1) {
			entityData.ManageEnabled = true
		} else {
			entityData.ManageEnabled = false
		}
	}

	if val, ok := parsedPayload["ttl"]; ok {
		parsedVal, _ := val.(float64)
		entityData.Ttl = int(parsedVal)
	}

	entityData.AuthKeys = auths

	return entityData
}
