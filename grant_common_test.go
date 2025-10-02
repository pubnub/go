package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_toUUIDsPermissionsMap(t *testing.T) {
	tests := []struct {
		name string
		have map[UserId]UserPermissions
		want map[string]UUIDPermissions
	}{{
		name: "Get",
		have: map[UserId]UserPermissions{"a": {
			Get: true,
		}},
		want: map[string]UUIDPermissions{"a": {
			Get: true,
		}}}, {
		name: "Update",
		have: map[UserId]UserPermissions{"a": {
			Update: true,
		}},
		want: map[string]UUIDPermissions{"a": {
			Update: true,
		}}}, {
		name: "Delete",
		have: map[UserId]UserPermissions{"a": {
			Delete: true,
		}},
		want: map[string]UUIDPermissions{"a": {
			Delete: true,
		}}}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, toUUIDsPermissionsMap(tt.have), "toUUIDsPermissionsMap(%v)", tt.have)
		})
	}
}

func Test_toChannelsPermissionsMap(t *testing.T) {
	tests := []struct {
		name string
		have map[SpaceId]SpacePermissions
		want map[string]ChannelPermissions
	}{{
		name: "Get",
		have: map[SpaceId]SpacePermissions{"a": {
			Get: true,
		}},
		want: map[string]ChannelPermissions{"a": {
			Get: true,
		}}}, {
		name: "Update",
		have: map[SpaceId]SpacePermissions{"a": {
			Update: true,
		}},
		want: map[string]ChannelPermissions{"a": {
			Update: true,
		}}}, {
		name: "Write",
		have: map[SpaceId]SpacePermissions{"a": {
			Write: true,
		}},
		want: map[string]ChannelPermissions{"a": {
			Write: true,
		}}}, {
		name: "Delete",
		have: map[SpaceId]SpacePermissions{"a": {
			Delete: true,
		}},
		want: map[string]ChannelPermissions{"a": {
			Delete: true,
		}}}, {
		name: "Read",
		have: map[SpaceId]SpacePermissions{"a": {
			Read: true,
		}},
		want: map[string]ChannelPermissions{"a": {
			Read: true,
		}}}, {
		name: "Manage",
		have: map[SpaceId]SpacePermissions{"a": {
			Manage: true,
		}},
		want: map[string]ChannelPermissions{"a": {
			Manage: true,
		}}}, {
		name: "Join",
		have: map[SpaceId]SpacePermissions{"a": {
			Join: true,
		}},
		want: map[string]ChannelPermissions{"a": {
			Join: true,
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, toChannelsPermissionsMap(tt.have), "toChannelsPermissionsMap(%v)", tt.have)
		})
	}
}
