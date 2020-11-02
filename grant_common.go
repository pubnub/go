package pubnub

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/bits"
	"strconv"
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
	PNUsers
	// PNSpaces for spaces
	PNSpaces
)

// ChannelPermissions contains all the acceptable perms for channels
type ChannelPermissions struct {
	Read   bool
	Write  bool
	Delete bool
}

// GroupPermissions contains all the acceptable perms for groups
type GroupPermissions struct {
	Read   bool
	Manage bool
}

// UserSpacePermissions contains all the acceptable perms for Users and Spaces
type UserSpacePermissions struct {
	Read   bool
	Write  bool
	Manage bool
	Delete bool
	Create bool
}

// ResourcePermissions contains all the applicable perms for bitmask translations.
type ResourcePermissions struct {
	Read   bool
	Write  bool
	Manage bool
	Delete bool
	Create bool
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

func parseGrantPerms(i int64, resourceType PNResourceType) interface{} {
	b := fmt.Sprintf("%08b\n", bits.Reverse8(uint8(i)))
	r := ResourcePermissions{
		Read:   false,
		Write:  false,
		Manage: false,
		Delete: false,
		Create: false,
	}
	for k, v := range b {
		i, _ := strconv.Atoi(string(v))
		switch k {
		case 0:
			r.Read = (i == 1)
		case 1:
			r.Write = (i == 1)
		case 2:
			r.Manage = (i == 1)
		case 3:
			r.Delete = (i == 1)
		case 4:
			r.Create = (i == 1)
		}
	}

	switch resourceType {
	case PNChannels:
		return ChannelPermissions{
			Read:   r.Read,
			Write:  r.Write,
			Delete: r.Delete,
		}
	case PNGroups:
		return GroupPermissions{
			Read:   r.Read,
			Manage: r.Manage,
		}
	default:
		return UserSpacePermissions{
			Read:   r.Read,
			Write:  r.Write,
			Delete: r.Delete,
			Manage: r.Manage,
			Create: r.Create,
		}

	}
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

	spaces := make(map[string]UserSpacePermissionsWithToken, len(res.Spaces))
	for k, v := range res.Spaces {
		spaces[k] = UserSpacePermissionsWithToken{
			Permissions:  parseGrantPerms(v, PNSpaces).(UserSpacePermissions),
			BitMaskPerms: v,
			Token:        token,
			Timestamp:    timetoken,
			TTL:          ttl,
		}
	}

	users := make(map[string]UserSpacePermissionsWithToken, len(res.Users))
	for k, v := range res.Users {
		users[k] = UserSpacePermissionsWithToken{
			Permissions:  parseGrantPerms(v, PNUsers).(UserSpacePermissions),
			BitMaskPerms: v,
			Token:        token,
			Timestamp:    timetoken,
			TTL:          ttl,
		}
	}

	g := GrantResourcesWithPermissions{
		Channels: channels,
		Users:    users,
		Groups:   groups,
		Spaces:   spaces,
	}
	return &g
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

// UserSpacePermissionsWithToken is used for users/spaces resource type permissions
type UserSpacePermissionsWithToken struct {
	Permissions  UserSpacePermissions
	BitMaskPerms int64
	Token        string
	Timestamp    int64
	TTL          int
}

// GrantResourcesWithPermissions is used as a common struct to store all resource type permissions
type GrantResourcesWithPermissions struct {
	Channels        map[string]ChannelPermissionsWithToken
	Groups          map[string]GroupPermissionsWithToken
	Users           map[string]UserSpacePermissionsWithToken
	Spaces          map[string]UserSpacePermissionsWithToken
	ChannelsPattern map[string]ChannelPermissionsWithToken
	GroupsPattern   map[string]GroupPermissionsWithToken
	UsersPattern    map[string]UserSpacePermissionsWithToken
	SpacesPattern   map[string]UserSpacePermissionsWithToken
}

// PermissionsBody is the struct used to decode the server response
type PermissionsBody struct {
	Resources GrantResources         `json:"resources"`
	Patterns  GrantResources         `json:"patterns"`
	Meta      map[string]interface{} `json:"meta"`
}

// GrantResources is the struct used to decode the server response
type GrantResources struct {
	Channels map[string]int64 `json:"channels" cbor:"chan"`
	Groups   map[string]int64 `json:"groups" cbor:"grp"`
	Users    map[string]int64 `json:"users" cbor:"usr"`
	Spaces   map[string]int64 `json:"spaces" cbor:"spc"`
}

// PNGrantTokenDecoded is the struct used to decode the server response
type PNGrantTokenDecoded struct {
	Resources GrantResources         `cbor:"res"`
	Patterns  GrantResources         `cbor:"pat"`
	Meta      map[string]interface{} `cbor:"meta"`
	Signature []byte                 `cbor:"sig"`
	Version   int                    `cbor:"v"`
	Timestamp int64                  `cbor:"t"`
	TTL       int                    `cbor:"ttl"`
}
