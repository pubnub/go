package pubnub

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

type patternPermissions struct {
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
