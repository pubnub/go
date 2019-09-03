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

type ResourcePermissions struct {
	Read   bool
	Write  bool
	Manage bool
	Delete bool
	Create bool
}

type patterns struct {
	ChannelsPattern      string
	ChannelGroupsPattern string
	SpacesPattern        string
	UsersPattern         string
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

func DecodeCBORToken(token string) (PNGrantTokenDecoded, error) {
	token = strings.Replace(token, "-", "+", -1)
	token = strings.Replace(token, "_", "/", -1)
	fmt.Println("\nStrings `-`, `_` replaced Token--->", token)
	// if i := len(s) % 4; i != 0 {
	// 	token += strings.Repeat("=", 4-i)
	// }
	// fmt.Println("Padded Token--->", token)
	var cborObject PNGrantTokenDecoded
	value, decodeErr := base64.StdEncoding.DecodeString(token)
	if decodeErr != nil {
		fmt.Println("\nDecoding Error--->", decodeErr)
	} else {
		fmt.Println("\nDecoded Token--->", string(value))
		//var resp interface{}
		c := cbor.NewDecoder(bytes.NewReader(value))

		err1 := c.Decode(&cborObject)
		if err1 != nil {
			fmt.Printf("\nCBOR decode Error---> %#v", err1)
			//log.Println("Write file:", ioutil.WriteFile("data.json", value, 0600))
			return cborObject, err1
		}
		return cborObject, nil
	}
	return cborObject, decodeErr
}

func ParseGrantPerms(i int64) ResourcePermissions {
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
	return r
}

func ParseGrantResources(res GrantResources) *GrantResourcesWithBoolPerms {
	channels := make(map[string]ResourcePermissions, len(res.Channels))

	for k, v := range res.Channels {
		fmt.Println("", k, v)
		channels[k] = ParseGrantPerms(v)
	}

	groups := make(map[string]ResourcePermissions, len(res.Groups))
	for k, v := range res.Groups {
		fmt.Println("", k, v)
		groups[k] = ParseGrantPerms(v)
	}

	spaces := make(map[string]ResourcePermissions, len(res.Spaces))
	for k, v := range res.Spaces {
		fmt.Println("", k, v)
		spaces[k] = ParseGrantPerms(v)
	}

	users := make(map[string]ResourcePermissions, len(res.Users))
	for k, v := range res.Users {
		fmt.Println("", k, v)
		users[k] = ParseGrantPerms(v)
	}

	g := GrantResourcesWithBoolPerms{
		Channels: channels,
		Users:    users,
		Groups:   groups,
		Spaces:   spaces,
	}
	return &g
}

type GrantResourcesWithBoolPerms struct {
	Channels map[string]ResourcePermissions
	Groups   map[string]ResourcePermissions
	Users    map[string]ResourcePermissions
	Spaces   map[string]ResourcePermissions
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
	Timetoken int64                  `cbor:"t"`
	TTL       int                    `cbor:"ttl"`
}
