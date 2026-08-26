package pubnub

import (
	"bytes"
	"encoding/base64"
	"strings"

	cbor "github.com/brianolson/cbor_go"
)

const (
	dataSyncEntitiesKey      = "datasync:entities"
	dataSyncUsersKey         = "datasync:users"
	dataSyncChannelsKey      = "datasync:channels"
	dataSyncRelationshipsKey = "datasync:relationships"
	dataSyncMembershipsKey   = "datasync:memberships"
	pnProjectionsMetaKey     = "pn-projections"

	// PNDataSyncDefaultProjection is the implicit projection when a resource
	// has no entry in meta.pn-projections.
	PNDataSyncDefaultProjection = "__default__"
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
	// PNDataSync for DataSync entities, relationships, and memberships
	PNDataSync
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
	Create bool
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
		Create: p.Create,
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
	Create bool
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
		Create: p.Create,
	}
}

// DataSyncPermissions contains the CRUD flags accepted for DataSync resource
// types (entities, relationships, memberships).
type DataSyncPermissions struct {
	Get    bool
	Create bool
	Update bool
	Delete bool
}

// PNDataSyncTokenScopes holds entity-level DataSync permissions keyed by
// resource id (resources) or pattern (patterns).
type PNDataSyncTokenScopes struct {
	Entities      map[string]DataSyncPermissions
	Relationships map[string]DataSyncPermissions
	Memberships   map[string]DataSyncPermissions
}

func (s PNDataSyncTokenScopes) empty() bool {
	return len(s.Entities) == 0 && len(s.Relationships) == 0 && len(s.Memberships) == 0
}

// PNDataSyncProjectionScope maps a DataSync resource id (or pattern) to the
// single projection name the principal is looking through.
// Users and Channels encode to datasync:users:{id} and datasync:channels:{id}
// in meta.pn-projections; their CRUD permissions still use UUIDs/Channels.
type PNDataSyncProjectionScope struct {
	Entities      map[string]string
	Users         map[string]string
	Channels      map[string]string
	Relationships map[string]string
	Memberships   map[string]string
}

func (s PNDataSyncProjectionScope) empty() bool {
	return len(s.Entities) == 0 && len(s.Users) == 0 && len(s.Channels) == 0 &&
		len(s.Relationships) == 0 && len(s.Memberships) == 0
}

// PNDataSyncProjections holds DataSync projection assignments for exact
// resources and patterns. Encoded into token meta under "pn-projections".
type PNDataSyncProjections struct {
	Resources PNDataSyncProjectionScope
	Patterns  PNDataSyncProjectionScope
}

func (p PNDataSyncProjections) empty() bool {
	return p.Resources.empty() && p.Patterns.empty()
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
	Projections    *PNDataSyncProjections
}

type PNTokenResources struct {
	Channels      map[string]ChannelPermissions
	ChannelGroups map[string]GroupPermissions
	UUIDs         map[string]UUIDPermissions
	DataSync      PNDataSyncTokenScopes
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
		Projections:    parseDataSyncProjections(permissions.Meta),
	}, nil
}

func grantResourcesToPNTokenResources(grantResources GrantResources) PNTokenResources {
	tokenResources := PNTokenResources{
		Channels:      make(map[string]ChannelPermissions),
		ChannelGroups: make(map[string]GroupPermissions),
		UUIDs:         make(map[string]UUIDPermissions),
		DataSync: PNDataSyncTokenScopes{
			Entities:      make(map[string]DataSyncPermissions),
			Relationships: make(map[string]DataSyncPermissions),
			Memberships:   make(map[string]DataSyncPermissions),
		},
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
	for k, v := range grantResources.DataSyncEntities {
		tokenResources.DataSync.Entities[k] = parseGrantPerms(v, PNDataSync).(DataSyncPermissions)
	}
	for k, v := range grantResources.DataSyncRelationships {
		tokenResources.DataSync.Relationships[k] = parseGrantPerms(v, PNDataSync).(DataSyncPermissions)
	}
	for k, v := range grantResources.DataSyncMemberships {
		tokenResources.DataSync.Memberships[k] = parseGrantPerms(v, PNDataSync).(DataSyncPermissions)
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
	create := i&int64(PNCreate) != 0
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
			Create: create,
		}
	case PNGroups:
		return GroupPermissions{
			Read:   read,
			Manage: manage,
		}
	case PNDataSync:
		return DataSyncPermissions{
			Get:    get,
			Create: create,
			Update: update,
			Delete: delete,
		}
	default:
		return UUIDPermissions{
			Get:    get,
			Update: update,
			Delete: delete,
			Create: create,
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
	Channels              map[string]int64 `json:"channels" cbor:"chan"`
	Groups                map[string]int64 `json:"groups" cbor:"grp"`
	UUIDs                 map[string]int64 `json:"uuids" cbor:"uuid"`
	Users                 map[string]int64 `json:"users" cbor:"usr"`
	Spaces                map[string]int64 `json:"spaces" cbor:"spc"`
	DataSyncEntities      map[string]int64 `json:"datasync:entities,omitempty" cbor:"datasync:entities"`
	DataSyncRelationships map[string]int64 `json:"datasync:relationships,omitempty" cbor:"datasync:relationships"`
	DataSyncMemberships   map[string]int64 `json:"datasync:memberships,omitempty" cbor:"datasync:memberships"`
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

func encodeProjectionScope(scope PNDataSyncProjectionScope) map[string]interface{} {
	out := make(map[string]interface{})
	putProjectionEntries(out, dataSyncEntitiesKey, scope.Entities)
	putProjectionEntries(out, dataSyncUsersKey, scope.Users)
	putProjectionEntries(out, dataSyncChannelsKey, scope.Channels)
	putProjectionEntries(out, dataSyncRelationshipsKey, scope.Relationships)
	putProjectionEntries(out, dataSyncMembershipsKey, scope.Memberships)
	return out
}

func putProjectionEntries(out map[string]interface{}, namespace string, ids map[string]string) {
	for id, name := range ids {
		if id == "" {
			continue
		}
		out[namespace+":"+id] = name
	}
}

func applyDataSyncProjections(meta map[string]interface{}, projections PNDataSyncProjections) map[string]interface{} {
	res := encodeProjectionScope(projections.Resources)
	pat := encodeProjectionScope(projections.Patterns)
	if len(res) == 0 && len(pat) == 0 {
		return meta
	}

	out := cloneStringInterfaceMap(meta)
	existing := asInterfaceMap(out[pnProjectionsMetaKey])

	if len(res) > 0 {
		mergedRes := asInterfaceMap(existing["res"])
		for k, v := range res {
			mergedRes[k] = v
		}
		existing["res"] = mergedRes
	}
	if len(pat) > 0 {
		mergedPat := asInterfaceMap(existing["pat"])
		for k, v := range pat {
			mergedPat[k] = v
		}
		existing["pat"] = mergedPat
	}

	out[pnProjectionsMetaKey] = existing
	return out
}

func parseDataSyncProjections(meta map[string]interface{}) *PNDataSyncProjections {
	if meta == nil {
		return nil
	}
	raw, ok := meta[pnProjectionsMetaKey]
	if !ok || raw == nil {
		return nil
	}

	block := asInterfaceMap(raw)
	if len(block) == 0 {
		return nil
	}

	projections := &PNDataSyncProjections{
		Resources: parseProjectionScope(block["res"]),
		Patterns:  parseProjectionScope(block["pat"]),
	}
	if projections.empty() {
		return nil
	}
	return projections
}

func parseProjectionScope(raw interface{}) PNDataSyncProjectionScope {
	scope := PNDataSyncProjectionScope{
		Entities:      make(map[string]string),
		Users:         make(map[string]string),
		Channels:      make(map[string]string),
		Relationships: make(map[string]string),
		Memberships:   make(map[string]string),
	}
	for key, val := range asInterfaceMap(raw) {
		name, ok := val.(string)
		if !ok || name == "" {
			continue
		}
		switch {
		case strings.HasPrefix(key, dataSyncEntitiesKey+":"):
			scope.Entities[key[len(dataSyncEntitiesKey)+1:]] = name
		case strings.HasPrefix(key, dataSyncUsersKey+":"):
			scope.Users[key[len(dataSyncUsersKey)+1:]] = name
		case strings.HasPrefix(key, dataSyncChannelsKey+":"):
			scope.Channels[key[len(dataSyncChannelsKey)+1:]] = name
		case strings.HasPrefix(key, dataSyncRelationshipsKey+":"):
			scope.Relationships[key[len(dataSyncRelationshipsKey)+1:]] = name
		case strings.HasPrefix(key, dataSyncMembershipsKey+":"):
			scope.Memberships[key[len(dataSyncMembershipsKey)+1:]] = name
		}
	}
	return scope
}

func cloneStringInterfaceMap(meta map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(meta)+1)
	for k, v := range meta {
		out[k] = v
	}
	return out
}

func asInterfaceMap(v interface{}) map[string]interface{} {
	switch m := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(m))
		for k, val := range m {
			out[k] = val
		}
		return out
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(m))
		for k, val := range m {
			if ks, ok := k.(string); ok {
				out[ks] = val
			}
		}
		return out
	case map[string]string:
		out := make(map[string]interface{}, len(m))
		for k, val := range m {
			out[k] = val
		}
		return out
	default:
		return map[string]interface{}{}
	}
}
