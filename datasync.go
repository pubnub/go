package pubnub

// DataSyncService groups DataSync-related API entry points on a PubNub instance.
// Access via pn.DataSync (e.g. pn.DataSync.CreateEntity()).
type DataSyncService struct {
	pn *PubNub
}

// CreateEntity provisions a new generic entity of a registered entity class.
func (d *DataSyncService) CreateEntity() *createEntityBuilder {
	return newCreateEntityBuilder(d.pn)
}

// CreateEntityWithContext provisions a new generic entity of a registered entity class.
func (d *DataSyncService) CreateEntityWithContext(ctx Context) *createEntityBuilder {
	return newCreateEntityBuilderWithContext(d.pn, ctx)
}

// GetEntity reads a single generic entity by its identifier.
func (d *DataSyncService) GetEntity() *getEntityBuilder {
	return newGetEntityBuilder(d.pn)
}

// GetEntityWithContext reads a single generic entity by its identifier.
func (d *DataSyncService) GetEntityWithContext(ctx Context) *getEntityBuilder {
	return newGetEntityBuilderWithContext(d.pn, ctx)
}

// GetEntities returns a paginated list of generic entities of a given class.
func (d *DataSyncService) GetEntities() *getEntitiesBuilder {
	return newGetEntitiesBuilder(d.pn)
}

// GetEntitiesWithContext returns a paginated list of generic entities of a given class.
func (d *DataSyncService) GetEntitiesWithContext(ctx Context) *getEntitiesBuilder {
	return newGetEntitiesBuilderWithContext(d.pn, ctx)
}

// UpdateEntity fully replaces the mutable fields of an existing entity (PUT).
func (d *DataSyncService) UpdateEntity() *updateEntityBuilder {
	return newUpdateEntityBuilder(d.pn)
}

// UpdateEntityWithContext fully replaces the mutable fields of an existing entity (PUT).
func (d *DataSyncService) UpdateEntityWithContext(ctx Context) *updateEntityBuilder {
	return newUpdateEntityBuilderWithContext(d.pn, ctx)
}

// PatchEntity partially updates an existing entity using RFC 6902 JSON Patch.
func (d *DataSyncService) PatchEntity() *patchEntityBuilder {
	return newPatchEntityBuilder(d.pn)
}

// PatchEntityWithContext partially updates an existing entity using RFC 6902 JSON Patch.
func (d *DataSyncService) PatchEntityWithContext(ctx Context) *patchEntityBuilder {
	return newPatchEntityBuilderWithContext(d.pn, ctx)
}

// RemoveEntity removes a generic entity by its identifier.
func (d *DataSyncService) RemoveEntity() *deleteEntityBuilder {
	return newDeleteEntityBuilder(d.pn)
}

// RemoveEntityWithContext removes a generic entity by its identifier.
func (d *DataSyncService) RemoveEntityWithContext(ctx Context) *deleteEntityBuilder {
	return newDeleteEntityBuilderWithContext(d.pn, ctx)
}

// CreateRelationship provisions a new relationship linking two existing entities.
func (d *DataSyncService) CreateRelationship() *createRelationshipBuilder {
	return newCreateRelationshipBuilder(d.pn)
}

// CreateRelationshipWithContext provisions a new relationship linking two existing entities.
func (d *DataSyncService) CreateRelationshipWithContext(ctx Context) *createRelationshipBuilder {
	return newCreateRelationshipBuilderWithContext(d.pn, ctx)
}

// GetRelationship reads a single relationship by its identifier.
func (d *DataSyncService) GetRelationship() *getRelationshipBuilder {
	return newGetRelationshipBuilder(d.pn)
}

// GetRelationshipWithContext reads a single relationship by its identifier.
func (d *DataSyncService) GetRelationshipWithContext(ctx Context) *getRelationshipBuilder {
	return newGetRelationshipBuilderWithContext(d.pn, ctx)
}

// GetRelationships returns a paginated list of relationships of a given class.
func (d *DataSyncService) GetRelationships() *getRelationshipsBuilder {
	return newGetRelationshipsBuilder(d.pn)
}

// GetRelationshipsWithContext returns a paginated list of relationships of a given class.
func (d *DataSyncService) GetRelationshipsWithContext(ctx Context) *getRelationshipsBuilder {
	return newGetRelationshipsBuilderWithContext(d.pn, ctx)
}

// UpdateRelationship fully replaces the mutable fields of an existing relationship (PUT).
func (d *DataSyncService) UpdateRelationship() *updateRelationshipBuilder {
	return newUpdateRelationshipBuilder(d.pn)
}

// UpdateRelationshipWithContext fully replaces the mutable fields of an existing relationship (PUT).
func (d *DataSyncService) UpdateRelationshipWithContext(ctx Context) *updateRelationshipBuilder {
	return newUpdateRelationshipBuilderWithContext(d.pn, ctx)
}

// PatchRelationship partially updates an existing relationship using RFC 6902 JSON Patch.
func (d *DataSyncService) PatchRelationship() *patchRelationshipBuilder {
	return newPatchRelationshipBuilder(d.pn)
}

// PatchRelationshipWithContext partially updates an existing relationship using RFC 6902 JSON Patch.
func (d *DataSyncService) PatchRelationshipWithContext(ctx Context) *patchRelationshipBuilder {
	return newPatchRelationshipBuilderWithContext(d.pn, ctx)
}

// RemoveRelationship removes a relationship by its identifier.
func (d *DataSyncService) RemoveRelationship() *deleteRelationshipBuilder {
	return newDeleteRelationshipBuilder(d.pn)
}

// RemoveRelationshipWithContext removes a relationship by its identifier.
func (d *DataSyncService) RemoveRelationshipWithContext(ctx Context) *deleteRelationshipBuilder {
	return newDeleteRelationshipBuilderWithContext(d.pn, ctx)
}

// CreateUser provisions a new user profile (specialized Entity of class User).
func (d *DataSyncService) CreateUser() *createUserBuilder {
	return newCreateUserBuilder(d.pn)
}

// CreateUserWithContext provisions a new user profile (specialized Entity of class User).
func (d *DataSyncService) CreateUserWithContext(ctx Context) *createUserBuilder {
	return newCreateUserBuilderWithContext(d.pn, ctx)
}

// GetUser reads a single user profile by its identifier.
func (d *DataSyncService) GetUser() *getUserBuilder {
	return newGetUserBuilder(d.pn)
}

// GetUserWithContext reads a single user profile by its identifier.
func (d *DataSyncService) GetUserWithContext(ctx Context) *getUserBuilder {
	return newGetUserBuilderWithContext(d.pn, ctx)
}

// GetUsers returns a paginated list of user profiles.
func (d *DataSyncService) GetUsers() *getUsersBuilder {
	return newGetUsersBuilder(d.pn)
}

// GetUsersWithContext returns a paginated list of user profiles.
func (d *DataSyncService) GetUsersWithContext(ctx Context) *getUsersBuilder {
	return newGetUsersBuilderWithContext(d.pn, ctx)
}

// UpdateUser fully replaces the mutable fields of an existing user (PUT).
func (d *DataSyncService) UpdateUser() *updateUserBuilder {
	return newUpdateUserBuilder(d.pn)
}

// UpdateUserWithContext fully replaces the mutable fields of an existing user (PUT).
func (d *DataSyncService) UpdateUserWithContext(ctx Context) *updateUserBuilder {
	return newUpdateUserBuilderWithContext(d.pn, ctx)
}

// PatchUser partially updates an existing user using RFC 6902 JSON Patch.
func (d *DataSyncService) PatchUser() *patchUserBuilder {
	return newPatchUserBuilder(d.pn)
}

// PatchUserWithContext partially updates an existing user using RFC 6902 JSON Patch.
func (d *DataSyncService) PatchUserWithContext(ctx Context) *patchUserBuilder {
	return newPatchUserBuilderWithContext(d.pn, ctx)
}

// RemoveUser removes a user profile by its identifier.
func (d *DataSyncService) RemoveUser() *deleteUserBuilder {
	return newDeleteUserBuilder(d.pn)
}

// RemoveUserWithContext removes a user profile by its identifier.
func (d *DataSyncService) RemoveUserWithContext(ctx Context) *deleteUserBuilder {
	return newDeleteUserBuilderWithContext(d.pn, ctx)
}

// CreateChannel provisions a new channel record (specialized Entity of class Channel).
func (d *DataSyncService) CreateChannel() *createChannelBuilder {
	return newCreateChannelBuilder(d.pn)
}

// CreateChannelWithContext provisions a new channel record (specialized Entity of class Channel).
func (d *DataSyncService) CreateChannelWithContext(ctx Context) *createChannelBuilder {
	return newCreateChannelBuilderWithContext(d.pn, ctx)
}

// GetChannel reads a single channel record by its identifier.
func (d *DataSyncService) GetChannel() *getChannelBuilder {
	return newGetChannelBuilder(d.pn)
}

// GetChannelWithContext reads a single channel record by its identifier.
func (d *DataSyncService) GetChannelWithContext(ctx Context) *getChannelBuilder {
	return newGetChannelBuilderWithContext(d.pn, ctx)
}

// GetChannels returns a paginated list of channel records.
func (d *DataSyncService) GetChannels() *getChannelsBuilder {
	return newGetChannelsBuilder(d.pn)
}

// GetChannelsWithContext returns a paginated list of channel records.
func (d *DataSyncService) GetChannelsWithContext(ctx Context) *getChannelsBuilder {
	return newGetChannelsBuilderWithContext(d.pn, ctx)
}

// UpdateChannel fully replaces the mutable fields of an existing channel (PUT).
func (d *DataSyncService) UpdateChannel() *updateChannelBuilder {
	return newUpdateChannelBuilder(d.pn)
}

// UpdateChannelWithContext fully replaces the mutable fields of an existing channel (PUT).
func (d *DataSyncService) UpdateChannelWithContext(ctx Context) *updateChannelBuilder {
	return newUpdateChannelBuilderWithContext(d.pn, ctx)
}

// PatchChannel partially updates an existing channel using RFC 6902 JSON Patch.
func (d *DataSyncService) PatchChannel() *patchChannelBuilder {
	return newPatchChannelBuilder(d.pn)
}

// PatchChannelWithContext partially updates an existing channel using RFC 6902 JSON Patch.
func (d *DataSyncService) PatchChannelWithContext(ctx Context) *patchChannelBuilder {
	return newPatchChannelBuilderWithContext(d.pn, ctx)
}

// RemoveChannel removes a channel record by its identifier.
func (d *DataSyncService) RemoveChannel() *deleteChannelBuilder {
	return newDeleteChannelBuilder(d.pn)
}

// RemoveChannelWithContext removes a channel record by its identifier.
func (d *DataSyncService) RemoveChannelWithContext(ctx Context) *deleteChannelBuilder {
	return newDeleteChannelBuilderWithContext(d.pn, ctx)
}

// CreateMembership provisions a new membership linking a user to a channel.
func (d *DataSyncService) CreateMembership() *createMembershipBuilder {
	return newCreateMembershipBuilder(d.pn)
}

// CreateMembershipWithContext provisions a new membership linking a user to a channel.
func (d *DataSyncService) CreateMembershipWithContext(ctx Context) *createMembershipBuilder {
	return newCreateMembershipBuilderWithContext(d.pn, ctx)
}

// GetMembership reads a single membership by its identifier.
func (d *DataSyncService) GetMembership() *getMembershipBuilder {
	return newGetMembershipBuilder(d.pn)
}

// GetMembershipWithContext reads a single membership by its identifier.
func (d *DataSyncService) GetMembershipWithContext(ctx Context) *getMembershipBuilder {
	return newGetMembershipBuilderWithContext(d.pn, ctx)
}

// GetMemberships returns a paginated list of memberships.
func (d *DataSyncService) GetMemberships() *getMembershipsBuilder {
	return newGetMembershipsBuilder(d.pn)
}

// GetMembershipsWithContext returns a paginated list of memberships.
func (d *DataSyncService) GetMembershipsWithContext(ctx Context) *getMembershipsBuilder {
	return newGetMembershipsBuilderWithContext(d.pn, ctx)
}

// UpdateMembership fully replaces the mutable fields of an existing membership (PUT).
func (d *DataSyncService) UpdateMembership() *updateMembershipBuilder {
	return newUpdateMembershipBuilder(d.pn)
}

// UpdateMembershipWithContext fully replaces the mutable fields of an existing membership (PUT).
func (d *DataSyncService) UpdateMembershipWithContext(ctx Context) *updateMembershipBuilder {
	return newUpdateMembershipBuilderWithContext(d.pn, ctx)
}

// PatchMembership partially updates an existing membership using RFC 6902 JSON Patch.
func (d *DataSyncService) PatchMembership() *patchMembershipBuilder {
	return newPatchMembershipBuilder(d.pn)
}

// PatchMembershipWithContext partially updates an existing membership using RFC 6902 JSON Patch.
func (d *DataSyncService) PatchMembershipWithContext(ctx Context) *patchMembershipBuilder {
	return newPatchMembershipBuilderWithContext(d.pn, ctx)
}

// RemoveMembership removes a membership by its identifier.
func (d *DataSyncService) RemoveMembership() *deleteMembershipBuilder {
	return newDeleteMembershipBuilder(d.pn)
}

// RemoveMembershipWithContext removes a membership by its identifier.
func (d *DataSyncService) RemoveMembershipWithContext(ctx Context) *deleteMembershipBuilder {
	return newDeleteMembershipBuilderWithContext(d.pn, ctx)
}
