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
