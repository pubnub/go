package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const GRANT_PATH = "/v1/auth/grant/sub-key/%s"

var emptyGrantResponse *GrantResponse

func GrantRequest(pn *PubNub, opts *GrantOpts) (*GrantResponse, error) {
	opts.pubnub = pn
	rawJson, err := executeRequest(opts)
	if err != nil {
		return emptyGrantResponse, err
	}

	return newGrantResponse(rawJson)
}

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

type GrantOpts struct {
	pubnub *PubNub
	ctx    Context

	AuthKeys []string
	Channels []string
	Groups   []string

	Read   bool
	Write  bool
	Manage bool

	// Stringified TTL
	// Max: 525600
	// Min: 1
	// Default: 1440
	// Setting 0 will apply the grant indefinitely
	Ttl string
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

	if len(o.AuthKeys) > 0 {
		q.Set("auth", strings.Join(o.AuthKeys, ","))
	}

	if len(o.Channels) > 0 {
		q.Set("channel", strings.Join(o.Channels, ","))
	}

	if len(o.Groups) > 0 {
		q.Set("channel-group", strings.Join(o.Groups, ","))
	}

	if o.Ttl != "" {
		ttl, err := strconv.ParseInt(o.Ttl, 10, 64)
		if err != nil {
			return &url.Values{}, err
		}

		if ttl >= -1 {
			q.Set("ttl", o.Ttl)
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

	Channels      map[string]map[string]*PNAccessManagerKeyData
	ChannelGroups map[string]map[string]*PNAccessManagerKeyData
}

type PNAccessManagerKeyData struct {
	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
}

type PNPAMEntityData struct {
	Name          string
	AuthKeys      []string
	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
	Ttl           int
}

func newGrantResponse(jsonBytes []byte) (*GrantResponse, error) {
	resp := &GrantResponse{}
	//
	// var value interface{}
	//
	// err := json.Unmarshal(jsonBytes, &value)
	// if err != nil {
	// 	e := pnerr.NewResponseParsingError("Error unmarshalling response",
	// 		ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)
	//
	// 	return emptyGrantResponse, e
	// }
	//
	// constructedChannels := make(map[string]map[string]*PNPAMEntityData)
	// constructedGroups := make(map[string]map[string]*PNPAMEntityData)
	// grantResp := &GrantResponse{}
	//
	// grantData, _ := value.(map[string]interface{})
	// payload := grantData["payload"]
	// parsedPayload := payload.(map[string]interface{})
	// auths, _ := parsedPayload["auths"].(map[string]interface{})
	//
	// if val, ok := parsedPayload["channel"]; ok {
	// 	var ttl int
	// 	channelName := parsedPayload["channel"]
	// 	constructedAuthKey := make(map[string]*PNAccessManagerKeyData)
	// 	entityData := PNPAMEntityData{
	// 		Name: channelName,
	// 	}
	//
	// 	for authKeyName, value := range auths {
	// 		auth, _ := value.(map[string]interface{})
	//
	// 		managerKeyData := PNAccessManagerKeyData{}
	//
	// 		if val, ok := auth["r"]; ok {
	// 			if val == "1" {
	// 				managerKeyData.ReadEnabled = true
	// 			} else {
	// 				managerKeyData.ReadEnabled = false
	// 			}
	// 		}
	//
	// 		if val, ok := auth["w"]; ok {
	// 			if val == "1" {
	// 				managerKeyData.WriteEnabled = true
	// 			} else {
	// 				managerKeyData.WriteEnabled = false
	// 			}
	// 		}
	//
	// 		if val, ok := auth["m"]; ok {
	// 			if val == "1" {
	// 				managerKeyData.ManageEnabled = true
	// 			} else {
	// 				managerKeyData.ManageEnabled = false
	// 			}
	// 		}
	//
	// 		if val, ok := auth["ttl"]; ok {
	// 			parsedVal, _ := val.(int)
	// 			entityData.Ttl = parsedVal
	// 			ttl = parsedVal
	// 		}
	//
	// 		constructedAuthKey[authKeyName] = &resp
	// 	}
	//
	// 	entityData.AuthKeys = constructedAuthKey
	// 	entityData.Ttl = ttl
	// 	constructedChannels[channelName] = entityData
	// }
	//
	// if val, ok := parsedPayload["channel-group"]; ok {
	// 	var ttl int
	// 	groupName := val
	// 	constructedAuthKey := make(map[string]*PNAccessManagerKeyData)
	// 	entityData := PNPAMEntityData{
	// 		Name: groupName,
	// 	}
	//
	// 	if groupString, ok := val.(string); ok {
	// 		for authKeyName, value := range auths {
	// 			auth, _ := value.(map[string]interface{})
	//
	// 			resp := PNPAMEntityData{}
	// 			resp.Name = groupName
	//
	// 			if val, ok := auth["r"]; ok {
	// 				if val == "1" {
	// 					resp.ReadEnabled = true
	// 				} else {
	// 					resp.ReadEnabled = false
	// 				}
	// 			}
	//
	// 			if val, ok := auth["w"]; ok {
	// 				if val == "1" {
	// 					resp.WriteEnabled = true
	// 				} else {
	// 					resp.WriteEnabled = false
	// 				}
	// 			}
	//
	// 			if val, ok := auth["m"]; ok {
	// 				if val == "1" {
	// 					resp.ManageEnabled = true
	// 				} else {
	// 					resp.ManageEnabled = false
	// 				}
	// 			}
	//
	// 			if val, ok := auth["ttl"]; ok {
	// 				parsedVal, _ := val.(int)
	// 				grantResp.Ttl = parsedVal
	// 			}
	//
	// 			constructedAuthKey[authKeyName] = &resp
	// 		}
	//
	// 		entityData.AuthKeys = constructedAuthKey
	// 		entityData.Ttl = ttl
	// 		constructedGroups[groupName] = entityData
	// 	}
	//
	// 	if groupSlice, ok := val.(map[string]interface{}); ok {
	// 		var ttl int
	// 		groupName := val
	// 		constructedAuthKey := make(map[string]*PNAccessManagerKeyData)
	// 		entityData := PNPAMEntityData{
	// 			Name: groupName,
	// 		}
	// 	}
	// }

	return resp, nil
}
