package pubnub

import (
	"bytes"
	"encoding/base64"
	"strings"

	cbor "github.com/brianolson/cbor_go"
)

// PNGrantBitMask is the type for perms BitMask
type PNGrantBitMask int64

const (
	// PNRead Read Perms
	PNRead PNGrantBitMask = 1
	// PNWrite Write Perms
	PNWrite = 2
	// PNManage Manage Perms
	PNManage = 4
	// PNDelete Delete Perms
	PNDelete = 8
	// PNCreate Create Perms
	PNCreate = 16
	// PNGet Get Perms
	PNGet = 32
	// PNUpdate Update Perms
	PNUpdate = 64
	// PNJoin Join Perms
	PNJoin = 128
)

// PNGrantType grant types
type PNGrantType int

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
	// PNGetEnabled Get Enabled. Applies to Objects
	PNGetEnabled
	// PNUpdateEnabled Update Enabled. Applies to Objects
	PNUpdateEnabled
	// PNJoinEnabled Join Enabled. Applies to Objects
	PNJoinEnabled
)

// PNResourceType grant types
type PNResourceType int

const (
	// PNChannels for channels
	PNChannels PNResourceType = 1 + iota
	// PNGroups for groups
	PNGroups
	// PNUsers for users
	PNUUIDs
)

// ChannelPermissions contains all the acceptable perms for channels
type ChannelPermissions struct {
	Read   bool
	Write  bool
	Delete bool
	Get    bool
	Manage bool
	Update bool
	Join   bool
}

type SpacePermissions ChannelPermissions

func toChannelsPermissionsMap(spacesPermissions map[SpaceId]SpacePermissions) map[string]ChannelPermissions {
	var channelsPermissions = make(map[string]ChannelPermissions)

	for name, p := range spacesPermissions {
		channelsPermissions[string(name)] = p.toChannelPermissions()
	}

	return channelsPermissions
}

func toChannelPatternsPermissionsMap(spacesPermissions map[string]SpacePermissions) map[string]ChannelPermissions {
	var channelsPermissions = make(map[string]ChannelPermissions)

	for name, p := range spacesPermissions {
		channelsPermissions[name] = p.toChannelPermissions()
	}

	return channelsPermissions
}

func (p SpacePermissions) toChannelPermissions() ChannelPermissions {
	return ChannelPermissions{
		Read:   p.Read,
		Write:  p.Write,
		Delete: p.Delete,
		Get:    p.Get,
		Manage: p.Manage,
		Update: p.Update,
		Join:   p.Join,
	}
}

// GroupPermissions contains all the acceptable perms for groups
type GroupPermissions struct {
	Read   bool
	Manage bool
}

type UUIDPermissions struct {
	Get    bool
	Update bool
	Delete bool
}

type UserPermissions UUIDPermissions

func toUUIDsPermissionsMap(usersPermissions map[UserId]UserPermissions) map[string]UUIDPermissions {
	var channelsPermissions = make(map[string]UUIDPermissions)

	for name, p := range usersPermissions {
		channelsPermissions[string(name)] = p.toUUIDPermissions()
	}

	return channelsPermissions
}

func toUUIDPatternsPermissionsMap(usersPermissions map[string]UserPermissions) map[string]UUIDPermissions {
	var channelsPermissions = make(map[string]UUIDPermissions)

	for name, p := range usersPermissions {
		channelsPermissions[name] = p.toUUIDPermissions()
	}

	return channelsPermissions
}

func (p UserPermissions) toUUIDPermissions() UUIDPermissions {
	return UUIDPermissions{
		Delete: p.Delete,
		Get:    p.Get,
		Update: p.Update,
	}
}

// PNPAMEntityData is the struct containing the access details of the channels.
type PNPAMEntityData struct {
	Name          string
	AuthKeys      map[string]*PNAccessManagerKeyData
	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
	DeleteEnabled bool
	GetEnabled    bool
	UpdateEnabled bool
	JoinEnabled   bool
	TTL           int
}

// PNAccessManagerKeyData is the struct containing the access details of the channel groups.
type PNAccessManagerKeyData struct {
	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
	DeleteEnabled bool
	GetEnabled    bool
	UpdateEnabled bool
	JoinEnabled   bool
	TTL           int
}

// GetPermissions decodes the CBORToken
func GetPermissions(token string) (PNGrantTokenDecoded, error) {
	token = strings.Replace(token, "-", "+", -1)
	token = strings.Replace(token, "_", "/", -1)
	if i := len(token) % 4; i != 0 {
		token += strings.Repeat("=", 4-i)
	}

	var cborObject PNGrantTokenDecoded
	value, decodeErr := base64.StdEncoding.DecodeString(token)
	if decodeErr != nil {
		return cborObject, decodeErr
	}

	c := cbor.NewDecoder(bytes.NewReader(value))
	err1 := c.Decode(&cborObject)
	if err1 != nil {
		return cborObject, err1
	}

	return cborObject, nil
}

type PNToken struct {
	Version        int
	Timestamp      int64
	TTL            int
	AuthorizedUUID string
	Resources      PNTokenResources
	Patterns       PNTokenResources
	Meta           map[string]interface{}
}

type PNTokenResources struct {
	Channels      map[string]ChannelPermissions
	ChannelGroups map[string]GroupPermissions
	UUIDs         map[string]UUIDPermissions
}

func ParseToken(token string) (*PNToken, error) {
	permissions, err := GetPermissions(token)

	if err != nil {
		return nil, err
	}

	resources := grantResourcesToPNTokenResources(permissions.Resources)
	patterns := grantResourcesToPNTokenResources(permissions.Patterns)

	return &PNToken{
		Version:        permissions.Version,
		Meta:           permissions.Meta,
		TTL:            permissions.TTL,
		Timestamp:      permissions.Timestamp,
		AuthorizedUUID: permissions.AuthorizedUUID,
		Resources:      resources,
		Patterns:       patterns,
	}, nil
}

func grantResourcesToPNTokenResources(grantResources GrantResources) PNTokenResources {
	tokenResources := PNTokenResources{
		Channels:      make(map[string]ChannelPermissions),
		ChannelGroups: make(map[string]GroupPermissions),
		UUIDs:         make(map[string]UUIDPermissions),
	}
	for k, v := range grantResources.Channels {
		tokenResources.Channels[k] = parseGrantPerms(v, PNChannels).(ChannelPermissions)
	}
	for k, v := range grantResources.Groups {
		tokenResources.ChannelGroups[k] = parseGrantPerms(v, PNGroups).(GroupPermissions)
	}
	for k, v := range grantResources.UUIDs {
		tokenResources.UUIDs[k] = parseGrantPerms(v, PNUUIDs).(UUIDPermissions)
	}
	return tokenResources
}

// ParseGrantResources parses the token for permissions and adds them along the other values to the GrantResourcesWithPermissions struct
func ParseGrantResources(res GrantResources, token string, timetoken int64, ttl int) *GrantResourcesWithPermissions {
	channels := make(map[string]ChannelPermissionsWithToken, len(res.Channels))

	for k, v := range res.Channels {
		channels[k] = ChannelPermissionsWithToken{
			Permissions:  parseGrantPerms(v, PNChannels).(ChannelPermissions),
			BitMaskPerms: v,
			Token:        token,
			Timestamp:    timetoken,
			TTL:          ttl,
		}
	}

	groups := make(map[string]GroupPermissionsWithToken, len(res.Groups))
	for k, v := range res.Groups {
		groups[k] = GroupPermissionsWithToken{
			Permissions:  parseGrantPerms(v, PNGroups).(GroupPermissions),
			BitMaskPerms: v,
			Token:        token,
			Timestamp:    timetoken,
			TTL:          ttl,
		}
	}

	g := GrantResourcesWithPermissions{
		Channels: channels,
		Groups:   groups,
	}
	return &g
}

func parseGrantPerms(i int64, resourceType PNResourceType) interface{} {
	read := i&int64(PNRead) != 0
	write := i&int64(PNWrite) != 0
	manage := i&int64(PNManage) != 0
	delete := i&int64(PNDelete) != 0
	get := i&int64(PNGet) != 0
	update := i&int64(PNUpdate) != 0
	join := i&int64(PNJoin) != 0

	switch resourceType {
	case PNChannels:
		return ChannelPermissions{
			Read:   read,
			Write:  write,
			Delete: delete,
			Update: update,
			Get:    get,
			Join:   join,
			Manage: manage,
		}
	case PNGroups:
		return GroupPermissions{
			Read:   read,
			Manage: manage,
		}
	default:
		return UUIDPermissions{
			Get:    get,
			Update: update,
			Delete: delete,
		}
	}
}

// ChannelPermissionsWithToken is used for channels resource type permissions
type ChannelPermissionsWithToken struct {
	Permissions  ChannelPermissions
	BitMaskPerms int64
	Token        string
	Timestamp    int64
	TTL          int
}

// GroupPermissionsWithToken is used for groups resource type permissions
type GroupPermissionsWithToken struct {
	Permissions  GroupPermissions
	BitMaskPerms int64
	Token        string
	Timestamp    int64
	TTL          int
}

// GrantResourcesWithPermissions is used as a common struct to store all resource type permissions
type GrantResourcesWithPermissions struct {
	Channels        map[string]ChannelPermissionsWithToken
	Groups          map[string]GroupPermissionsWithToken
	ChannelsPattern map[string]ChannelPermissionsWithToken
	GroupsPattern   map[string]GroupPermissionsWithToken
}

// PermissionsBody is the struct used to decode the server response
type PermissionsBody struct {
	Resources      GrantResources         `json:"resources"`
	Patterns       GrantResources         `json:"patterns"`
	Meta           map[string]interface{} `json:"meta"`
	AuthorizedUUID string                 `json:"uuid,omitempty"`
}

// GrantResources is the struct used to decode the server response
type GrantResources struct {
	Channels map[string]int64 `json:"channels" cbor:"chan"`
	Groups   map[string]int64 `json:"groups" cbor:"grp"`
	UUIDs    map[string]int64 `json:"uuids" cbor:"uuid"`
	Users    map[string]int64 `json:"users" cbor:"usr"`
	Spaces   map[string]int64 `json:"spaces" cbor:"spc"`
}

// PNGrantTokenDecoded is the struct used to decode the server response
type PNGrantTokenDecoded struct {
	Resources      GrantResources         `cbor:"res"`
	Patterns       GrantResources         `cbor:"pat"`
	Meta           map[string]interface{} `cbor:"meta"`
	Signature      []byte                 `cbor:"sig"`
	Version        int                    `cbor:"v"`
	Timestamp      int64                  `cbor:"t"`
	TTL            int                    `cbor:"ttl"`
	AuthorizedUUID string                 `cbor:"uuid"`
}
