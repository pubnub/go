package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrantTokenParseResourcePermissions(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGrantTokenBuilder(pn)

	m := map[string]ChannelPermissions{
		"channel": {
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
		},
	}

	r := o.opts.parseResourcePermissions(m, PNChannels)
	for _, v := range r {
		assert.Equal(int64(15), v)
	}
}

func TestGrantTokenParseResourcePermissions2(t *testing.T) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	o := newGrantTokenBuilder(pn)
	m := map[string]ChannelPermissions{
		"channel": {
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: false,
		},
	}

	r := o.opts.parseResourcePermissions(m, PNChannels)
	for _, v := range r {
		assert.Equal(int64(7), v)
	}
}

func TestGrantTokenParseResourcePermissions3(t *testing.T) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	o := newGrantTokenBuilder(pn)
	m := map[string]ChannelPermissions{
		"channel": {
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: false,
		},
		"channel2": {
			Read:   true,
			Write:  false,
			Manage: true,
			Delete: false,
		},
	}

	r := o.opts.parseResourcePermissions(m, PNChannels)
	assert.Equal(int64(7), r["channel"])
	assert.Equal(int64(5), r["channel2"])
}

func Test_GrantToken(t *testing.T) {
	pn := NewPubNub(NewDemoConfig())

	tests := []struct {
		name string
		have endpointOpts
		want string
	}{{
		name: "GrantToken objects v2",
		have: pn.GrantToken().
			TTL(100).
			Channels(map[string]ChannelPermissions{
				"channel": {
					Write:  false,
					Read:   true,
					Delete: false,
				},
			}).
			ChannelGroups(map[string]GroupPermissions{
				"cg": {
					Read:   true,
					Manage: true,
				},
				"cg2": {
					Read:   true,
					Manage: false,
				},
			}).opts,
		want: `{"ttl":100,"permissions":{"resources":{"channels":{"channel":1},"groups":{"cg":5,"cg2":1},"uuids":{},"users":{},"spaces":{}},"patterns":{"channels":{},"groups":{},"uuids":{},"users":{},"spaces":{}},"meta":{}}}`},
		{
			name: "GrantToken SUM",
			have: pn.GrantToken().
				TTL(100).
				SpacesPermissions(map[SpaceId]SpacePermissions{
					"channel": {
						Write:  false,
						Read:   true,
						Delete: false,
					},
				}).
				SpacePatternsPermissions(map[SpaceId]SpacePermissions{
					"channel": {
						Write:  true,
						Read:   true,
						Delete: false,
					},
				}).
				UsersPermissions(map[UserId]UserPermissions{
					"user": {
						Get:    true,
						Update: true,
						Delete: true,
					},
				}).
				UserPatternsPermissions(map[UserId]UserPermissions{
					"users*": {
						Get:    true,
						Update: false,
						Delete: true,
					},
				}).opts,
			want: `{"ttl":100,"permissions":{"resources":{"channels":{"channel":1},"groups":{},"uuids":{"user":104},"users":{},"spaces":{}},"patterns":{"channels":{"channel":3},"groups":{},"uuids":{"users*":40},"users":{},"spaces":{}},"meta":{}}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := tt.have.buildBody()
			assert.Nil(t, err)
			_, err = tt.have.buildPath()
			assert.Nil(t, err)
			assert.Equalf(t, tt.want, string(body), "GrantToken(%v)", tt.have)
		})
	}
}
