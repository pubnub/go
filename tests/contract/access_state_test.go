package contract

import pubnub "github.com/pubnub/go/v7"

type accessStateKey struct{}

type accessState struct {
	CurrentPermissions             interface{}
	ChannelPermissions             map[string]*pubnub.ChannelPermissions
	ChannelPatternPermissions      map[string]*pubnub.ChannelPermissions
	ChannelGroupPermissions        map[string]*pubnub.GroupPermissions
	ChannelGroupPatternPermissions map[string]*pubnub.GroupPermissions
	UUIDPermissions                map[string]*pubnub.UUIDPermissions
	UUIDPatternPermissions         map[string]*pubnub.UUIDPermissions
	TTL                            int
	TokenString                    string
	AuthorizedUUID                 string
	GrantTokenResult               pubnub.PNGrantTokenResponse
	ParsedToken                    *pubnub.PNToken
	ResourcePermissions            interface{}
	RevokeTokenResult              pubnub.PNRevokeTokenResponse
}

func newAccessState(pn *pubnub.PubNub) *accessState {
	return &accessState{
		TTL:                            0,
		ChannelPermissions:             make(map[string]*pubnub.ChannelPermissions),
		ChannelPatternPermissions:      make(map[string]*pubnub.ChannelPermissions),
		ChannelGroupPermissions:        make(map[string]*pubnub.GroupPermissions),
		ChannelGroupPatternPermissions: make(map[string]*pubnub.GroupPermissions),
		UUIDPermissions:                make(map[string]*pubnub.UUIDPermissions),
		UUIDPatternPermissions:         make(map[string]*pubnub.UUIDPermissions),
		CurrentPermissions:             nil}
}
