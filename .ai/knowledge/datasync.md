# DataSync

Conventions for the DataSync data-plane API. Confirm details against current code. Per-endpoint godoc already describes each builder — follow peers in `datasync_*.go` rather than inventing a new shape.

Shared endpoint, error, and test rules still apply: [endpoints.md](endpoints.md), [errors.md](errors.md), [testing.md](testing.md).

## Surface

Callers use `pn.DataSync` (`DataSyncService` in `datasync.go`), not methods on `*PubNub`. Five families — Entity, Relationship, User, Channel, Membership — each with a `WithContext` twin.

Verbs: Create, Get, Get-all, Set (PUT), Update (PATCH), Remove.

Paths, vendor content types, `entitiesDefaultLimit` (20, server range 1–100), and shared helpers live in `datasync_common.go`:

- `/v1/datasync/subkeys/%s/{entities|relationships|users|channels|memberships}`
- `conflictingDataSyncFilters`
- `marshalDataSyncSetBody`
- `PNJSONPatchOperation.MarshalJSON`

## Names that must stay split

- The public method is `Remove*`. The file and unexported builder are `delete*` (`datasync_delete_entity.go`, `deleteEntityBuilder`).
- Entity and Relationship operation constants have no `DataSync` infix: `PNCreateEntityOperation`, `PNCreateRelationshipOperation`, and their Get / Set / Update / Remove siblings. User, Channel, and Membership constants do: `PNCreateDataSyncUserOperation`, `PNCreateDataSyncChannelOperation`, `PNCreateDataSyncMembershipOperation`, and siblings. Do not rename or reorder either group. New operations still need the iota constant, `OperationType.String()`, and a telemetry case in `telemetry_manager.go`.
- `PNUser` and its response envelopes alias `PNEntity`. `PNDataSyncChannel` also aliases `PNEntity`, so it does not collide with the Objects `PNChannel` type. `PNDataSyncMembership` is its own struct (`channelId`, `userId`). Relationship resources use `entityAId` / `entityBId`.

## Set, update, and concurrency

Set is a full replacement. Class identity is immutable and omitted from the body. `marshalDataSyncSetBody` omits `payload` when `Payload` was never called or was passed `nil` (the server resets it). An explicit empty map is sent as `{}`.

Update is RFC 6902 JSON Patch (`application/json-patch+json`). Builders expose `Operations` plus `Add`, `Remove`, `Replace`, `Move`, `Copy`, and `Test`. `PNJSONPatchOperation.MarshalJSON` always emits `value` for add, replace, and test, including JSON null.

Set, Update, and Remove send `If-Match` only when `IfMatchETag` was called (`setIfMatchETag`).

## Lists

List calls require the class name. Class version is sent only when its setter was called (`setEntityClassVersion`), so an unset version is not written as `0`.

`FilterFast` (query `filter_fast`, strongly consistent, fewer conditions) and `Filter` (query `filter`, eventually consistent, richer expressions) cannot both be set. Validation uses `StrExclusiveDataSyncFilter`. `filter`, `filter_fast`, and `sort` are URL-encoded in `buildURL`.

Pagination is a forward-only cursor: pass `Meta.NextCursor` back as `Cursor` and stop when `HasNext` is false.

`EntityClassLevel` disambiguates classes that share a name: `PNEntityClassLevelGlobal` (`"Global"`) or `PNEntityClassLevelSubKey` (`"SubKey"`).

DataSync validation strings live next to the others in `pubnub.go` (`StrMissingEntityClass`, `StrMissingEntityID`, `StrInvalidEntityClassVersion`, and the relationship, user, channel, and membership variants).

## Events

Events are not a REST endpoint. Subscribe message type 5 with `source: "data-sync"` and `version: "1.0"` is parsed by `createPNDataSyncEventResult` in `subscription_manager.go` and delivered on `Listener.DataSyncEvent` as `PNDataSyncEventResult`. Other versions are ignored. Create and update fill `Entity` (entity, user, and channel kinds), `Relationship`, or `Membership`. Delete fills `ID` and `DeletedAt` and leaves the resource nil.

Subscribe to the resource id. The admin projection arrives on the channel `__admin__` plus that id. Every subscribed channel shares one `DataSyncEvent` stream, so match on `Channel`.

## PAM

Grant encoding lives in `grant_common.go`.

- CRUD flags are `DataSyncPermissions` (`Get`, `Create`, `Update`, `Delete`) on `PNDataSyncTokenScopes` for entities, relationships, and memberships.
- Projections are `PNDataSyncProjections`. The implicit projection is `PNDataSyncDefaultProjection` (`"__default__"`), stored under token meta `pn-projections`. Users and channels encode as `datasync:users:{id}` and `datasync:channels:{id}`.
- DataSync User CRUD uses `Users()` / `UsersPattern()` (token users). `UUIDs()` / `UUIDsPattern()` grant App Context uuid permissions and do not authorize DataSync User. Channel CRUD uses channel permissions.

Live DataSync tests use `DS_PUBLISH_KEY` / `DS_SUBSCRIBE_KEY` / `DS_SECRET_KEY` and skip when any is unset. See [testing.md](testing.md).
