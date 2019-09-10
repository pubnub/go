package pubnub

import (
	"bytes"
	"encoding/base64"
	"fmt"
	//"io/ioutil"
	//"log"
	//"encoding/json"
	"math/bits"
	"strconv"
	"strings"

	cbor "github.com/brianolson/cbor_go"
)

type PNGrantBitMask int64

const (
	PNRead   PNGrantBitMask = 1
	PNWrite                 = 2
	PNManage                = 4
	PNDelete                = 8
	PNCreate                = 16
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
)

// PNResourceType grant types
type PNResourceType int

const (
	PNChannels PNResourceType = 1 + iota
	PNGroups
	PNUsers
	PNSpaces
)

type ChannelPermissions struct {
	Read   bool
	Write  bool
	Delete bool
}

type GroupPermissions struct {
	Read   bool
	Manage bool
}

type UserSpacePermissions struct {
	Read   bool
	Write  bool
	Manage bool
	Delete bool
	Create bool
}

type ResourcePermissions struct {
	Read   bool
	Write  bool
	Manage bool
	Delete bool
	Create bool
}

// type patterns struct {
// 	ChannelsPattern      string
// 	ChannelGroupsPattern string
// 	SpacesPattern        string
// 	UsersPattern         string
// }

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

// DecodeCBORToken
func GetPermissions(token string) (PNGrantTokenDecoded, error) {
	token = strings.Replace(token, "-", "+", -1)
	token = strings.Replace(token, "_", "/", -1)
	fmt.Println("\nStrings `-`, `_` replaced Token--->", token)
	if i := len(token) % 4; i != 0 {
		token += strings.Repeat("=", 4-i)
	}
	fmt.Println("Padded Token--->", token)

	var cborObject PNGrantTokenDecoded
	value, decodeErr := base64.StdEncoding.DecodeString(token)
	if decodeErr != nil {
		fmt.Println("\nDecoding Error--->", decodeErr)
	} else {
		fmt.Println("\nDecoded Token--->", string(value), value)

		c := cbor.NewDecoder(bytes.NewReader(value))
		// var res map[string]interface{}
		// c.Decode(&res)
		// //fmt.Println(res)
		// jsonSerialized, e := json.Marshal(res)
		// fmt.Println(e)
		// fmt.Println("jsonSerialized--->", string(jsonSerialized), res)

		err1 := c.Decode(&cborObject)
		if err1 != nil {
			fmt.Printf("\nCBOR decode Error---> %#v", err1)
			//log.Println("Write file:", ioutil.WriteFile("data.json", value, 0600))
			return cborObject, err1
		}
		//log.Println("Write file:", ioutil.WriteFile("data.json", value, 0600))
		//log.Println("Write file:", ioutil.WriteFile("data.json", []byte(string(cborObject)), 0600))

		return cborObject, nil
	}
	return cborObject, decodeErr
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
		fmt.Println(k, string(v))
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
	fmt.Println(r)
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

// func ParseGrantResources(res GrantResources) *GrantResourcesWithBoolPerms {
// 	channels := make(map[string]ResourcePermissions, len(res.Channels))

// 	for k, v := range res.Channels {
// 		fmt.Println("", k, v)
// 		channels[k] = ParseGrantPerms(v)
// 	}

// 	groups := make(map[string]ResourcePermissions, len(res.Groups))
// 	for k, v := range res.Groups {
// 		fmt.Println("", k, v)
// 		groups[k] = ParseGrantPerms(v)
// 	}

// 	spaces := make(map[string]ResourcePermissions, len(res.Spaces))
// 	for k, v := range res.Spaces {
// 		fmt.Println("", k, v)
// 		spaces[k] = ParseGrantPerms(v)
// 	}

// 	users := make(map[string]ResourcePermissions, len(res.Users))
// 	for k, v := range res.Users {
// 		fmt.Println("", k, v)
// 		users[k] = ParseGrantPerms(v)
// 	}

// 	g := GrantResourcesWithBoolPerms{
// 		Channels: channels,
// 		Users:    users,
// 		Groups:   groups,
// 		Spaces:   spaces,
// 	}
// 	return &g
// }

func ParseGrantResources(res GrantResources, token string, timetoken int64) *GrantResourcesWithPermissions {
	channels := make(map[string]ChannelPermissionsWithToken, len(res.Channels))

	for k, v := range res.Channels {
		fmt.Println("", k, v)
		channels[k] = ChannelPermissionsWithToken{
			Permissions:  parseGrantPerms(v, PNChannels).(ChannelPermissions),
			BitMaskPerms: v,
			Token:        token,
			Timestamp:    timetoken,
		}
	}

	groups := make(map[string]GroupPermissionsWithToken, len(res.Groups))
	for k, v := range res.Groups {
		fmt.Println("", k, v)
		groups[k] = GroupPermissionsWithToken{
			Permissions:  parseGrantPerms(v, PNGroups).(GroupPermissions),
			BitMaskPerms: v,
			Token:        token,
			Timestamp:    timetoken,
		}
	}

	spaces := make(map[string]UserSpacePermissionsWithToken, len(res.Spaces))
	for k, v := range res.Spaces {
		fmt.Println("", k, v)
		spaces[k] = UserSpacePermissionsWithToken{
			Permissions:  parseGrantPerms(v, PNSpaces).(UserSpacePermissions),
			BitMaskPerms: v,
			Token:        token,
			Timestamp:    timetoken,
		}
	}

	users := make(map[string]UserSpacePermissionsWithToken, len(res.Users))
	for k, v := range res.Users {
		fmt.Println("", k, v)
		users[k] = UserSpacePermissionsWithToken{
			Permissions:  parseGrantPerms(v, PNUsers).(UserSpacePermissions),
			BitMaskPerms: v,
			Token:        token,
			Timestamp:    timetoken,
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

// type GrantResourcesWithBoolPerms struct {
// 	Channels map[string]ResourcePermissions
// 	Groups   map[string]ResourcePermissions
// 	Users    map[string]ResourcePermissions
// 	Spaces   map[string]ResourcePermissions
// }

type ChannelPermissionsWithToken struct {
	Permissions  ChannelPermissions
	BitMaskPerms int64
	Token        string
	Timestamp    int64
}

type GroupPermissionsWithToken struct {
	Permissions  GroupPermissions
	BitMaskPerms int64
	Token        string
	Timestamp    int64
}

type UserSpacePermissionsWithToken struct {
	Permissions  UserSpacePermissions
	BitMaskPerms int64
	Token        string
	Timestamp    int64
}

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

type PermissionsBody struct {
	Resources GrantResources         `json:"resources"`
	Patterns  GrantResources         `json:"patterns"`
	Meta      map[string]interface{} `json:"meta"`
}

type GrantResources struct {
	Channels map[string]int64 `json:"channels" cbor:"chan"`
	Groups   map[string]int64 `json:"groups" cbor:"grp"`
	Users    map[string]int64 `json:"users" cbor:"usr"`
	Spaces   map[string]int64 `json:"spaces" cbor:"spc"`
}

type PNGrantTokenDecoded struct {
	Resources GrantResources         `cbor:"res"`
	Patterns  GrantResources         `cbor:"pat"`
	Meta      map[string]interface{} `cbor:"meta"`
	Signature []byte                 `cbor:"sig"`
	Version   int                    `cbor:"v"`
	Timestamp int64                  `cbor:"t"`
	TTL       int                    `cbor:"ttl"`
}
