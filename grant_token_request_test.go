package pubnub

import (
	"fmt"
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
		have endpoint
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
			name: "GrantToken Entities",
			have: pn.GrantToken().
				TTL(100).
				SpacesPermissions(map[SpaceId]SpacePermissions{
					"channel": {
						Write:  false,
						Read:   true,
						Delete: false,
					},
				}).
				SpacePatternsPermissions(map[string]SpacePermissions{
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
				UserPatternsPermissions(map[string]UserPermissions{
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

// Additional validation tests specific to GrantToken
func TestGrantTokenValidateTTL(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGrantTokenOpts(pn, pn.ctx)
	opts.TTL = 0 // Invalid TTL

	assert.Equal("pubnub/validation: pubnub: Grant Token: Invalid TTL", opts.validate().Error())
}

func TestGrantTokenValidateTTLNegative(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGrantTokenOpts(pn, pn.ctx)
	opts.TTL = -1 // Invalid TTL

	assert.Equal("pubnub/validation: pubnub: Grant Token: Invalid TTL", opts.validate().Error())
}

func TestGrantTokenValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newGrantTokenOpts(pn, pn.ctx)
	opts.TTL = 1440 // Valid TTL

	assert.Nil(opts.validate())
}

func TestGrantTokenValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newGrantTokenOpts(pn, pn.ctx)
	opts.TTL = 1440

	assert.Equal("pubnub/validation: pubnub: Grant Token: Missing Subscribe Key", opts.validate().Error())
}

func TestGrantTokenValidateMissingPublishKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.PublishKey = ""
	opts := newGrantTokenOpts(pn, pn.ctx)
	opts.TTL = 1440

	assert.Equal("pubnub/validation: pubnub: Grant Token: Missing Publish Key", opts.validate().Error())
}

func TestGrantTokenValidateMissingSecretKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SecretKey = ""
	opts := newGrantTokenOpts(pn, pn.ctx)
	opts.TTL = 1440

	assert.Equal("pubnub/validation: pubnub: Grant Token: Missing Secret Key", opts.validate().Error())
}

// Builder pattern tests for GrantToken
func TestGrantTokenBuilder(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test basic builder
	builder := newGrantTokenBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)

	// Test TTL setting
	result := builder.TTL(1440)
	assert.Equal(builder, result) // Should return same instance for chaining
	assert.Equal(1440, builder.opts.TTL)
	assert.True(builder.opts.setTTL)
}

func TestGrantTokenBuilderWithContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGrantTokenBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestGrantTokenBuilderMeta(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testMeta := map[string]interface{}{
		"custom_field": "custom_value",
		"number_field": 42,
		"bool_field":   true,
	}

	builder := newGrantTokenBuilder(pn)
	result := builder.Meta(testMeta)
	assert.Equal(builder, result) // Should return same instance for chaining
	assert.Equal(testMeta, builder.opts.Meta)
}

func TestGrantTokenBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParams := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	builder := newGrantTokenBuilder(pn)
	result := builder.QueryParam(queryParams)
	assert.Equal(builder, result) // Should return same instance for chaining
	assert.Equal(queryParams, builder.opts.QueryParam)
}

func TestGrantTokenBuilderToObjectsBuilder(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGrantTokenBuilder(pn)
	builder.TTL(1440)

	// Test transition to objects builder
	objectsBuilder := builder.Channels(map[string]ChannelPermissions{
		"test_channel": {Read: true, Write: true},
	})
	assert.NotNil(objectsBuilder)
	assert.Equal(1440, objectsBuilder.opts.TTL) // Should preserve state
}

func TestGrantTokenBuilderToEntitiesBuilder(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGrantTokenBuilder(pn)
	builder.TTL(1440)

	// Test transition to entities builder
	entitiesBuilder := builder.SpacesPermissions(map[SpaceId]SpacePermissions{
		"test_space": {Read: true, Write: true},
	})
	assert.NotNil(entitiesBuilder)
	assert.Equal(1440, entitiesBuilder.opts.TTL) // Should preserve state
}

func TestGrantTokenObjectsBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	channels := map[string]ChannelPermissions{
		"channel1": {Read: true, Write: true, Manage: true},
	}
	groups := map[string]GroupPermissions{
		"group1": {Read: true, Manage: true},
	}
	uuids := map[string]UUIDPermissions{
		"uuid1": {Get: true, Update: true},
	}

	builder := newGrantTokenBuilder(pn)
	objectsBuilder := builder.Channels(channels)

	// Test method chaining
	result := objectsBuilder.ChannelGroups(groups).UUIDs(uuids)
	assert.Equal(objectsBuilder, result)

	// Verify all values are set correctly
	assert.Equal(channels, objectsBuilder.opts.Channels)
	assert.Equal(groups, objectsBuilder.opts.ChannelGroups)
	assert.Equal(uuids, objectsBuilder.opts.UUIDs)
}

func TestGrantTokenObjectsBuilderPatterns(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	channelPatterns := map[string]ChannelPermissions{
		"channel.*": {Read: true, Write: true},
	}
	groupPatterns := map[string]GroupPermissions{
		"group.*": {Read: true, Manage: true},
	}
	uuidPatterns := map[string]UUIDPermissions{
		"uuid.*": {Get: true, Update: true},
	}

	builder := newGrantTokenBuilder(pn)
	objectsBuilder := builder.ChannelsPattern(channelPatterns)

	// Test pattern methods
	result := objectsBuilder.ChannelGroupsPattern(groupPatterns).UUIDsPattern(uuidPatterns)
	assert.Equal(objectsBuilder, result)

	// Verify all pattern values are set correctly
	assert.Equal(channelPatterns, objectsBuilder.opts.ChannelsPattern)
	assert.Equal(groupPatterns, objectsBuilder.opts.ChannelGroupsPattern)
	assert.Equal(uuidPatterns, objectsBuilder.opts.UUIDsPattern)
}

func TestGrantTokenObjectsBuilderAuthorizedUUID(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGrantTokenBuilder(pn)
	objectsBuilder := builder.AuthorizedUUID("test-uuid-12345")

	assert.Equal("test-uuid-12345", objectsBuilder.opts.AuthorizedUUID)
}

func TestGrantTokenEntitiesBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	spaces := map[SpaceId]SpacePermissions{
		"space1": {Read: true, Write: true, Manage: true},
	}
	users := map[UserId]UserPermissions{
		"user1": {Get: true, Update: true, Delete: true},
	}

	builder := newGrantTokenBuilder(pn)
	entitiesBuilder := builder.SpacesPermissions(spaces)

	// Test method chaining
	result := entitiesBuilder.UsersPermissions(users)
	assert.Equal(entitiesBuilder, result)

	// Verify values are converted and set correctly
	assert.NotNil(entitiesBuilder.opts.Channels)
	assert.NotNil(entitiesBuilder.opts.UUIDs)
}

func TestGrantTokenEntitiesBuilderPatterns(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	spacePatterns := map[string]SpacePermissions{
		"space.*": {Read: true, Write: true},
	}
	userPatterns := map[string]UserPermissions{
		"user.*": {Get: true, Update: true},
	}

	builder := newGrantTokenBuilder(pn)
	entitiesBuilder := builder.SpacePatternsPermissions(spacePatterns)

	// Test pattern methods
	result := entitiesBuilder.UserPatternsPermissions(userPatterns)
	assert.Equal(entitiesBuilder, result)

	// Verify pattern values are converted and set correctly
	assert.NotNil(entitiesBuilder.opts.ChannelsPattern)
	assert.NotNil(entitiesBuilder.opts.UUIDsPattern)
}

// Parameter-specific tests for GrantToken
func TestGrantTokenTTLBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test minimum valid TTL
	builder1 := newGrantTokenBuilder(pn)
	builder1.TTL(1) // Minimum valid
	assert.Equal(1, builder1.opts.TTL)
	assert.Nil(builder1.opts.validate())

	// Test default TTL
	builder2 := newGrantTokenBuilder(pn)
	builder2.TTL(1440) // Default (24 hours)
	assert.Equal(1440, builder2.opts.TTL)
	assert.Nil(builder2.opts.validate())

	// Test maximum reasonable TTL
	builder3 := newGrantTokenBuilder(pn)
	builder3.TTL(525600) // Maximum (1 year)
	assert.Equal(525600, builder3.opts.TTL)
	assert.Nil(builder3.opts.validate())

	// Test invalid TTL values
	builder4 := newGrantTokenBuilder(pn)
	builder4.TTL(0) // Invalid
	assert.NotNil(builder4.opts.validate())

	builder5 := newGrantTokenBuilder(pn)
	builder5.TTL(-1) // Invalid
	assert.NotNil(builder5.opts.validate())
}

func TestGrantTokenWithComplexMeta(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with complex nested meta structure
	complexMeta := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   123,
			"name": "Test User",
			"settings": map[string]interface{}{
				"theme":         "dark",
				"notifications": true,
			},
		},
		"permissions": map[string]interface{}{
			"granted_at": "2023-01-01T00:00:00Z",
			"level":      "premium",
			"features":   []string{"read", "write", "manage"},
		},
		"unicode_text": "测试 русский ファイル",
	}

	builder := newGrantTokenBuilder(pn)
	builder.TTL(1440)
	builder.Meta(complexMeta)

	assert.Equal(complexMeta, builder.opts.Meta)
	assert.Nil(builder.opts.validate())
}

func TestGrantTokenWithEmptyMeta(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with empty meta
	emptyMeta := map[string]interface{}{}

	builder := newGrantTokenBuilder(pn)
	builder.TTL(1440)
	builder.Meta(emptyMeta)

	assert.Equal(emptyMeta, builder.opts.Meta)
	assert.Nil(builder.opts.validate())
}

func TestGrantTokenWithNilMeta(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGrantTokenBuilder(pn)
	builder.TTL(1440)
	builder.Meta(nil)

	assert.Nil(builder.opts.Meta)
	assert.Nil(builder.opts.validate())
}

// Edge case tests for GrantToken
func TestGrantTokenWithSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with special characters in various fields
	channels := map[string]ChannelPermissions{
		"channel-with-dashes":        {Read: true, Write: true},
		"channel_with_underscores":   {Read: true, Write: true},
		"channel@with#special$chars": {Read: true, Write: true},
		"channel.with.dots":          {Read: true, Write: true},
		"channel:with:colons":        {Read: true, Write: true},
	}

	groups := map[string]GroupPermissions{
		"group-with-dashes":        {Read: true, Manage: true},
		"group_with_underscores":   {Read: true, Manage: true},
		"group@with#special$chars": {Read: true, Manage: true},
	}

	uuids := map[string]UUIDPermissions{
		"uuid-with-dashes":        {Get: true, Update: true},
		"uuid_with_underscores":   {Get: true, Update: true},
		"uuid@with#special$chars": {Get: true, Update: true},
	}

	builder := newGrantTokenBuilder(pn)
	objectsBuilder := builder.TTL(1440).
		Channels(channels).
		ChannelGroups(groups).
		UUIDs(uuids).
		AuthorizedUUID("uuid-with-special@chars#123")

	assert.Nil(objectsBuilder.opts.validate())

	// Test path building
	path, err := objectsBuilder.opts.buildPath()
	assert.Nil(err)
	assert.NotEmpty(path)
}

func TestGrantTokenWithUnicodeCharacters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with Unicode characters
	channels := map[string]ChannelPermissions{
		"频道测试":       {Read: true, Write: true},
		"канал-тест": {Read: true, Write: true},
		"チャンネル-テスト":  {Read: true, Write: true},
	}

	groups := map[string]GroupPermissions{
		"组测试":         {Read: true, Manage: true},
		"группа-тест": {Read: true, Manage: true},
		"グループ-テスト":    {Read: true, Manage: true},
	}

	meta := map[string]interface{}{
		"测试字段":  "测试值",
		"поле":  "значение",
		"フィールド": "値",
	}

	builder := newGrantTokenBuilder(pn)
	objectsBuilder := builder.TTL(1440).
		Channels(channels).
		ChannelGroups(groups).
		Meta(meta).
		AuthorizedUUID("测试-uuid-русский-ユーザー")

	assert.Nil(objectsBuilder.opts.validate())

	// Test body building
	body, err := objectsBuilder.opts.buildBody()
	assert.Nil(err)
	assert.NotEmpty(body)
}

func TestGrantTokenWithLargeScale(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with many channels and permissions
	channels := make(map[string]ChannelPermissions)
	for i := 0; i < 100; i++ {
		channels[fmt.Sprintf("channel_%d", i)] = ChannelPermissions{
			Read:   true,
			Write:  i%2 == 0,
			Manage: i%3 == 0,
			Delete: i%4 == 0,
			Get:    i%5 == 0,
			Update: i%6 == 0,
			Join:   i%7 == 0,
		}
	}

	groups := make(map[string]GroupPermissions)
	for i := 0; i < 50; i++ {
		groups[fmt.Sprintf("group_%d", i)] = GroupPermissions{
			Read:   true,
			Manage: i%2 == 0,
		}
	}

	uuids := make(map[string]UUIDPermissions)
	for i := 0; i < 50; i++ {
		uuids[fmt.Sprintf("uuid_%d", i)] = UUIDPermissions{
			Get:    true,
			Update: i%2 == 0,
			Delete: i%3 == 0,
		}
	}

	// Large meta object
	largeMeta := make(map[string]interface{})
	for i := 0; i < 50; i++ {
		largeMeta[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d_with_some_longer_text", i)
	}

	builder := newGrantTokenBuilder(pn)
	objectsBuilder := builder.TTL(1440).
		Channels(channels).
		ChannelGroups(groups).
		UUIDs(uuids).
		Meta(largeMeta)

	assert.Nil(objectsBuilder.opts.validate())

	// Test body building with large data
	body, err := objectsBuilder.opts.buildBody()
	assert.Nil(err)
	assert.NotEmpty(body)

	// Verify some data is present in the body
	bodyStr := string(body)
	assert.Contains(bodyStr, "channel_0")
	assert.Contains(bodyStr, "group_0")
	assert.Contains(bodyStr, "uuid_0")
	assert.Contains(bodyStr, "key_0")
}

func TestGrantTokenWithEmptyPermissions(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with empty permission maps
	channels := map[string]ChannelPermissions{}
	groups := map[string]GroupPermissions{}
	uuids := map[string]UUIDPermissions{}

	builder := newGrantTokenBuilder(pn)
	objectsBuilder := builder.TTL(1440).
		Channels(channels).
		ChannelGroups(groups).
		UUIDs(uuids)

	assert.Nil(objectsBuilder.opts.validate())

	// Test body building
	body, err := objectsBuilder.opts.buildBody()
	assert.Nil(err)
	assert.NotEmpty(body)

	// Should contain empty resource structures
	bodyStr := string(body)
	assert.Contains(bodyStr, `"channels":{}`)
	assert.Contains(bodyStr, `"groups":{}`)
	assert.Contains(bodyStr, `"uuids":{}`)
}

func TestGrantTokenWithAllPermissionCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test all possible permission combinations for channels
	channels := map[string]ChannelPermissions{
		"all_permissions": {
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
			Get:    true,
			Update: true,
			Join:   true,
		},
		"no_permissions": {
			Read:   false,
			Write:  false,
			Manage: false,
			Delete: false,
			Get:    false,
			Update: false,
			Join:   false,
		},
		"read_only": {
			Read:   true,
			Write:  false,
			Manage: false,
			Delete: false,
			Get:    false,
			Update: false,
			Join:   false,
		},
		"read_write": {
			Read:   true,
			Write:  true,
			Manage: false,
			Delete: false,
			Get:    false,
			Update: false,
			Join:   false,
		},
	}

	builder := newGrantTokenBuilder(pn)
	objectsBuilder := builder.TTL(1440).Channels(channels)

	assert.Nil(objectsBuilder.opts.validate())

	// Test permission parsing
	parsed := objectsBuilder.opts.parseResourcePermissions(channels, PNChannels)

	// All permissions: Read(1) + Write(2) + Manage(4) + Delete(8) + Get(32) + Update(64) + Join(128) = 239
	assert.Equal(int64(239), parsed["all_permissions"])

	// No permissions
	assert.Equal(int64(0), parsed["no_permissions"])

	// Read only
	assert.Equal(int64(1), parsed["read_only"])

	// Read + Write
	assert.Equal(int64(3), parsed["read_write"])
}

func TestGrantTokenPatternVsResourceDistinction(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test distinction between resources and patterns
	specificChannels := map[string]ChannelPermissions{
		"specific_channel": {Read: true, Write: true},
	}

	channelPatterns := map[string]ChannelPermissions{
		"pattern_*": {Read: true},
		"*.pattern": {Write: true},
	}

	builder := newGrantTokenBuilder(pn)
	objectsBuilder := builder.TTL(1440).
		Channels(specificChannels).
		ChannelsPattern(channelPatterns)

	assert.Nil(objectsBuilder.opts.validate())

	// Test body building
	body, err := objectsBuilder.opts.buildBody()
	assert.Nil(err)

	bodyStr := string(body)
	// Should have both resources and patterns sections
	assert.Contains(bodyStr, `"resources"`)
	assert.Contains(bodyStr, `"patterns"`)
	assert.Contains(bodyStr, "specific_channel")
	assert.Contains(bodyStr, "pattern_*")
}

func TestGrantTokenPathBuilding(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newGrantTokenBuilder(pn)
	opts := builder.opts

	path, err := opts.buildPath()
	assert.Nil(err)
	expectedPath := fmt.Sprintf("/v3/pam/%s/grant", pn.Config.SubscribeKey)
	assert.Equal(expectedPath, path)
}

func TestGrantTokenExecuteErrorHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test with invalid configuration (missing secret key)
	pn.Config.SecretKey = ""

	builder := newGrantTokenBuilder(pn)
	objectsBuilder := builder.TTL(1440).Channels(map[string]ChannelPermissions{
		"test": {Read: true},
	})

	// Execute should fail with validation error
	_, _, err := objectsBuilder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Secret Key")
}

func TestGrantTokenDataSyncBitmasks(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGrantTokenBuilder(pn)

	r := o.opts.parseResourcePermissions(map[string]DataSyncPermissions{
		"get":        {Get: true},
		"get_update": {Get: true, Update: true},
		"create":     {Create: true},
		"all":        {Get: true, Create: true, Update: true, Delete: true},
		"none":       {},
	}, PNDataSync)

	assert.Equal(int64(32), r["get"])
	assert.Equal(int64(96), r["get_update"])
	assert.Equal(int64(16), r["create"])
	assert.Equal(int64(120), r["all"])
	assert.Equal(int64(0), r["none"])
}

func TestGrantTokenUUIDAndChannelCreateBit(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGrantTokenBuilder(pn)

	uuids := o.opts.parseResourcePermissions(map[string]UUIDPermissions{
		"user-123": {Get: true, Create: true, Update: true, Delete: true},
	}, PNUUIDs)
	assert.Equal(int64(120), uuids["user-123"])

	channels := o.opts.parseResourcePermissions(map[string]ChannelPermissions{
		"channel-X": {Create: true},
	}, PNChannels)
	assert.Equal(int64(16), channels["channel-X"])
}

func Test_GrantTokenDataSync(t *testing.T) {
	pn := NewPubNub(NewDemoConfig())

	tests := []struct {
		name string
		have endpoint
		want string
	}{{
		name: "mixed channels datasync projections and user meta",
		have: pn.GrantToken().
			TTL(1440).
			AuthorizedUUID("user-123").
			Channels(map[string]ChannelPermissions{
				"channel-X": {Read: true, Get: true, Update: true, Join: true},
			}).
			DataSync(PNDataSyncTokenScopes{
				Entities: map[string]DataSyncPermissions{
					"order-456": {Get: true, Update: true},
				},
				Relationships: map[string]DataSyncPermissions{
					"user.A:channel.X": {Get: true},
				},
				Memberships: map[string]DataSyncPermissions{
					"user-123:channel-X": {Get: true},
				},
			}).
			DataSyncPattern(PNDataSyncTokenScopes{
				Entities: map[string]DataSyncPermissions{
					"order-*": {Get: true},
				},
			}).
			DataSyncProjections(PNDataSyncProjections{
				Resources: PNDataSyncProjectionScope{
					Entities: map[string]string{"user.A": "admin"},
					Relationships: map[string]string{
						"user.A:channel.X": "admin",
					},
				},
			}).
			Meta(map[string]interface{}{
				"app_version": "2.0",
			}).opts,
		want: `{"ttl":1440,"permissions":{"resources":{"channels":{"channel-X":225},"groups":{},"uuids":{},"users":{},"spaces":{},"datasync:entities":{"order-456":96},"datasync:relationships":{"user.A:channel.X":32},"datasync:memberships":{"user-123:channel-X":32}},"patterns":{"channels":{},"groups":{},"uuids":{},"users":{},"spaces":{},"datasync:entities":{"order-*":32}},"meta":{"app_version":"2.0","pn-projections":{"res":{"datasync:entities:user.A":"admin","datasync:relationships:user.A:channel.X":"admin"}}},"uuid":"user-123"}}`,
	}, {
		name: "projections only",
		have: pn.GrantToken().
			TTL(30).
			DataSyncProjections(PNDataSyncProjections{
				Resources: PNDataSyncProjectionScope{
					Entities: map[string]string{
						"school-greenwood-001": PNDataSyncDefaultProjection,
					},
				},
			}).opts,
		want: `{"ttl":30,"permissions":{"resources":{"channels":{},"groups":{},"uuids":{},"users":{},"spaces":{}},"patterns":{"channels":{},"groups":{},"uuids":{},"users":{},"spaces":{}},"meta":{"pn-projections":{"res":{"datasync:entities:school-greenwood-001":"__default__"}}}}}`,
	}, {
		name: "user and channel projections",
		have: pn.GrantToken().
			TTL(60).
			DataSyncProjections(PNDataSyncProjections{
				Resources: PNDataSyncProjectionScope{
					Users:    map[string]string{"user.A": "admin"},
					Channels: map[string]string{"channel-X": "moderator"},
				},
				Patterns: PNDataSyncProjectionScope{
					Users: map[string]string{"user.*": PNDataSyncDefaultProjection},
				},
			}).opts,
		want: `{"ttl":60,"permissions":{"resources":{"channels":{},"groups":{},"uuids":{},"users":{},"spaces":{}},"patterns":{"channels":{},"groups":{},"uuids":{},"users":{},"spaces":{}},"meta":{"pn-projections":{"pat":{"datasync:users:user.*":"__default__"},"res":{"datasync:channels:channel-X":"moderator","datasync:users:user.A":"admin"}}}}}`,
	}, {
		name: "empty datasync omitted",
		have: pn.GrantToken().
			TTL(100).
			DataSync(PNDataSyncTokenScopes{}).
			DataSyncPattern(PNDataSyncTokenScopes{}).opts,
		want: `{"ttl":100,"permissions":{"resources":{"channels":{},"groups":{},"uuids":{},"users":{},"spaces":{}},"patterns":{"channels":{},"groups":{},"uuids":{},"users":{},"spaces":{}},"meta":{}}}`,
	}, {
		name: "entities builder datasync",
		have: pn.GrantToken().
			TTL(100).
			UsersPermissions(map[UserId]UserPermissions{
				"user": {Get: true, Create: true},
			}).
			DataSync(PNDataSyncTokenScopes{
				Memberships: map[string]DataSyncPermissions{
					"user:channel-X": {Get: true, Update: true},
				},
			}).opts,
		want: `{"ttl":100,"permissions":{"resources":{"channels":{},"groups":{},"uuids":{"user":48},"users":{},"spaces":{},"datasync:memberships":{"user:channel-X":96}},"patterns":{"channels":{},"groups":{},"uuids":{},"users":{},"spaces":{}},"meta":{}}}`,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := tt.have.buildBody()
			assert.Nil(t, err)
			assert.Equalf(t, tt.want, string(body), "GrantToken(%v)", tt.have)
		})
	}
}

func TestGrantTokenDataSyncDoesNotMutateCallerMeta(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	meta := map[string]interface{}{
		"app_version": "2.0",
	}
	opts := pn.GrantToken().
		TTL(1440).
		Meta(meta).
		DataSyncProjections(PNDataSyncProjections{
			Resources: PNDataSyncProjectionScope{
				Entities: map[string]string{"user.A": "admin"},
			},
		}).opts

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Contains(string(body), `"pn-projections"`)
	assert.NotContains(meta, pnProjectionsMetaKey)
}

func TestGrantTokenDataSyncBuilderChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	scopes := PNDataSyncTokenScopes{
		Entities: map[string]DataSyncPermissions{"order-456": {Get: true}},
	}
	patterns := PNDataSyncTokenScopes{
		Entities: map[string]DataSyncPermissions{"order-*": {Get: true}},
	}
	projections := PNDataSyncProjections{
		Resources: PNDataSyncProjectionScope{
			Entities: map[string]string{"user.A": "admin"},
		},
	}

	root := newGrantTokenBuilder(pn)
	assert.Equal(root, root.DataSync(scopes).DataSyncPattern(patterns).DataSyncProjections(projections))
	assert.Equal(scopes, root.opts.DataSync)
	assert.Equal(patterns, root.opts.DataSyncPattern)
	assert.Equal(projections, root.opts.DataSyncProjections)

	objectsBuilder := newGrantTokenBuilder(pn).Channels(map[string]ChannelPermissions{
		"ch": {Read: true},
	})
	assert.Equal(objectsBuilder, objectsBuilder.DataSync(scopes).DataSyncPattern(patterns).DataSyncProjections(projections))
	assert.Equal(scopes, objectsBuilder.opts.DataSync)

	entitiesBuilder := newGrantTokenBuilder(pn).UsersPermissions(map[UserId]UserPermissions{
		"user": {Get: true},
	})
	assert.Equal(entitiesBuilder, entitiesBuilder.DataSync(scopes).DataSyncPattern(patterns).DataSyncProjections(projections))
	assert.Equal(scopes, entitiesBuilder.opts.DataSync)
}

func TestGrantTokenDataSyncProjectionMergeWinsOnCollision(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	opts := pn.GrantToken().
		TTL(60).
		Meta(map[string]interface{}{
			pnProjectionsMetaKey: map[string]interface{}{
				"res": map[string]interface{}{
					"datasync:entities:user.A": "stale",
					"datasync:entities:other":  "keep",
				},
			},
		}).
		DataSyncProjections(PNDataSyncProjections{
			Resources: PNDataSyncProjectionScope{
				Entities: map[string]string{"user.A": "admin"},
			},
		}).opts

	body, err := opts.buildBody()
	assert.Nil(err)
	bodyStr := string(body)
	assert.Contains(bodyStr, `"datasync:entities:user.A":"admin"`)
	assert.Contains(bodyStr, `"datasync:entities:other":"keep"`)
	assert.NotContains(bodyStr, `"stale"`)
}
