package pubnub

import (
	"fmt"
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

func ParsePerms(i int64) {
	b := fmt.Sprintf("%b", i)
	fmt.Println(b)
	for k, v := range b {
		fmt.Println(k, string(v))
	}
}

func ParseGrantResources(res GrantResources) {
	for k, v := range res.Channels {
		fmt.Println("", k, v)
		// b1 := make([]byte, 8)
		// binary.LittleEndian.PutUint64(b1, uint64(v))
		// b := fmtBits(b1) //fmt.Sprintf("%b", v) //strconv.FormatInt(int64(v), 2)
		//fmt.Println(v, b, b1, fmt.Sprintf("%b", v))
		ParsePerms(v)
		//bm := [8]byte(int64(v))
		// for i := 1; i <= 64; i++ {
		// 	if b.IsSet(i) {
		// 		fmt.Print(i, " ")
		// 	}
		// }
	}
	for k, v := range res.Groups {
		fmt.Println("", k, v)
		ParsePerms(v)
	}
	for k, v := range res.Spaces {
		fmt.Println("", k, v)
		ParsePerms(v)
	}
	for k, v := range res.Users {
		fmt.Println("", k, v)
		ParsePerms(v)
	}
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
