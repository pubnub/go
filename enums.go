package pubnub

import (
	"fmt"
	"reflect"
)

// StatusCategory is used as an enum to catgorize the various status events
// in the APIs lifecycle
type StatusCategory int

// OperationType is used as an enum to catgorize the various operations
// in the APIs lifecycle
type OperationType int

// ReconnectionPolicy is used as an enum to catgorize the reconnection policies
type ReconnectionPolicy int

// PNPushType is used as an enum to catgorize the available Push Types
type PNPushType int

// PNUUIDMetadataInclude is used as an enum to catgorize the available UUID include types
type PNUUIDMetadataInclude int

// PNChannelMetadataInclude is used as an enum to catgorize the available Channel include types
type PNChannelMetadataInclude int

// PNMembershipsInclude is used as an enum to catgorize the available Memberships include types
type PNMembershipsInclude int

// PNChannelMembersInclude is used as an enum to catgorize the available Members include types
type PNChannelMembersInclude int

// PNObjectsEvent is used as an enum to catgorize the available Object Events
type PNObjectsEvent string

// PNObjectsEventType is used as an enum to catgorize the available Object Event types
type PNObjectsEventType string

// PNMessageActionsEventType is used as an enum to catgorize the available Message Actions Event types
type PNMessageActionsEventType string

// PNPushEnvironment is used as an enum to catgorize the available Message Actions Event types
type PNPushEnvironment string

const (
	//PNPushEnvironmentDevelopment for development
	PNPushEnvironmentDevelopment PNPushEnvironment = "development"
	//PNPushEnvironmentProduction for production
	PNPushEnvironmentProduction = "production"
)

const (
	// PNMessageActionsAdded is the enum when the event of type `added` occurs
	PNMessageActionsAdded PNMessageActionsEventType = "added"
	// PNMessageActionsRemoved is the enum when the event of type `removed` occurs
	PNMessageActionsRemoved = "removed"
)

const (
	// PNObjectsMembershipEvent is the enum when the event of type `membership` occurs
	PNObjectsMembershipEvent PNObjectsEventType = "membership"
	// PNObjectsChannelEvent is the enum when the event of type `channel` occurs
	PNObjectsChannelEvent = "channel"
	// PNObjectsUUIDEvent is the enum when the event of type `uuid` occurs
	PNObjectsUUIDEvent = "uuid"
	// PNObjectsNoneEvent is used for error handling
	PNObjectsNoneEvent = "none"
)

const (
	// PNObjectsEventRemove is the enum when the event `delete` occurs
	PNObjectsEventRemove PNObjectsEvent = "delete"
	// PNObjectsEventSet is the enum when the event `set` occurs
	PNObjectsEventSet = "set"
)

const (
	// PNUUIDMetadataIncludeCustom is the enum equivalent to the value `custom` available UUID include types
	PNUUIDMetadataIncludeCustom PNUUIDMetadataInclude = 1 + iota
)

const (
	// PNChannelMetadataIncludeCustom is the enum equivalent to the value `custom` available Channel include types
	PNChannelMetadataIncludeCustom PNChannelMetadataInclude = 1 + iota
)

func (s PNUUIDMetadataInclude) String() string {
	return [...]string{"custom"}[s-1]
}

func (s PNChannelMetadataInclude) String() string {
	return [...]string{"custom"}[s-1]
}

const (
	// PNMembershipsIncludeCustom is the enum equivalent to the value `custom` available Memberships include types
	PNMembershipsIncludeCustom PNMembershipsInclude = 1 + iota
	// PNMembershipsIncludeChannel is the enum equivalent to the value `channel` available Memberships include types
	PNMembershipsIncludeChannel
	// PNMembershipsIncludeChannelCustom is the enum equivalent to the value `channel.custom` available Memberships include types
	PNMembershipsIncludeChannelCustom
)

func (s PNMembershipsInclude) String() string {
	return [...]string{"custom", "channel", "channel.custom"}[s-1]
}

const (
	// PNChannelMembersIncludeCustom is the enum equivalent to the value `custom` available Members include types
	PNChannelMembersIncludeCustom PNChannelMembersInclude = 1 + iota
	// PNChannelMembersIncludeUUID is the enum equivalent to the value `uuid` available Members include types
	PNChannelMembersIncludeUUID
	// PNChannelMembersIncludeUUIDCustom is the enum equivalent to the value `uuid.custom` available Members include types
	PNChannelMembersIncludeUUIDCustom
)

func (s PNChannelMembersInclude) String() string {
	//return [...]string{"custom", "user", "user.custom", "uuid", "uuid.custom"}[s-1]
	return [...]string{"custom", "uuid", "uuid.custom"}[s-1]
}

// PNMessageType is used as an enum to catgorize the Subscribe response.
type PNMessageType int

const (
	// PNNonePolicy is to be used when selecting the no Reconnection Policy
	// ReconnectionPolicy is set in the config.
	PNNonePolicy ReconnectionPolicy = 1 + iota
	// PNLinearPolicy is to be used when selecting the Linear Reconnection Policy
	// ReconnectionPolicy is set in the config.
	PNLinearPolicy
	// PNExponentialPolicy is to be used when selecting the Exponential Reconnection Policy
	// ReconnectionPolicy is set in the config.
	PNExponentialPolicy
)

const (
	// PNMessageTypeSignal is to identify Signal the Subscribe response
	PNMessageTypeSignal PNMessageType = 1 + iota
	// PNMessageTypeObjects is to identify Objects the Subscribe response
	PNMessageTypeObjects
	// PNMessageTypeMessageActions is to identify Actions the Subscribe response
	PNMessageTypeMessageActions
	// PNMessageTypeFile is to identify Files the Subscribe response
	PNMessageTypeFile
)

const (
	// PNUnknownCategory as the StatusCategory means an unknown status category event occurred.
	PNUnknownCategory StatusCategory = 1 + iota
	// PNTimeoutCategory as the StatusCategory means the request timeout has reached.
	PNTimeoutCategory
	// PNConnectedCategory as the StatusCategory means the channel is subscribed to receive messages.
	PNConnectedCategory
	// PNDisconnectedCategory as the StatusCategory means a disconnection occurred due to network issues.
	PNDisconnectedCategory
	// PNCancelledCategory as the StatusCategory means the context was cancelled.
	PNCancelledCategory
	// PNLoopStopCategory as the StatusCategory means the subscribe loop was stopped.
	PNLoopStopCategory
	// PNAcknowledgmentCategory as the StatusCategory is the Acknowledgement of an operation (like Unsubscribe).
	PNAcknowledgmentCategory
	// PNBadRequestCategory as the StatusCategory means the request was malformed.
	PNBadRequestCategory
	// PNAccessDeniedCategory as the StatusCategory means that PAM is enabled and the channel is not granted R/W access.
	PNAccessDeniedCategory
	// PNNoStubMatchedCategory as the StatusCategory means an unknown status category event occurred.
	PNNoStubMatchedCategory
	// PNReconnectedCategory as the StatusCategory means that the network was reconnected (after a disconnection).
	// Applicable on for PNLinearPolicy and PNExponentialPolicy.
	PNReconnectedCategory
	// PNReconnectionAttemptsExhausted as the StatusCategory means that the reconnection attempts
	// to reconnect to the network were exhausted. All channels would be unsubscribed at this point.
	// Applicable on for PNLinearPolicy and PNExponentialPolicy.
	// Reconnection attempts are set in the config: MaximumReconnectionRetries.
	PNReconnectionAttemptsExhausted
	// PNRequestMessageCountExceededCategory is fired when the MessageQueueOverflowCount limit is exceeded by the number of messages received in a single subscribe request
	PNRequestMessageCountExceededCategory
)

const (
	// PNSubscribeOperation is the enum used for the Subcribe operation.
	PNSubscribeOperation OperationType = 1 + iota
	// PNUnsubscribeOperation is the enum used for the Unsubcribe operation.
	PNUnsubscribeOperation
	// PNPublishOperation is the enum used for the Publish operation.
	PNPublishOperation
	// PNFireOperation is the enum used for the Fire operation.
	PNFireOperation
	// PNHistoryOperation is the enum used for the History operation.
	PNHistoryOperation
	// PNFetchMessagesOperation is the enum used for the Fetch operation.
	PNFetchMessagesOperation
	// PNWhereNowOperation is the enum used for the Where Now operation.
	PNWhereNowOperation
	// PNHereNowOperation is the enum used for the Here Now operation.
	PNHereNowOperation
	// PNHeartBeatOperation is the enum used for the Heartbeat operation.
	PNHeartBeatOperation
	// PNSetStateOperation is the enum used for the Set State operation.
	PNSetStateOperation
	// PNGetStateOperation is the enum used for the Get State operation.
	PNGetStateOperation
	// PNAddChannelsToChannelGroupOperation is the enum used for the Add Channels to Channel Group operation.
	PNAddChannelsToChannelGroupOperation
	// PNRemoveChannelFromChannelGroupOperation is the enum used for the Remove Channels from Channel Group operation.
	PNRemoveChannelFromChannelGroupOperation
	// PNRemoveGroupOperation is the enum used for the Remove Channel Group operation.
	PNRemoveGroupOperation
	// PNChannelsForGroupOperation is the enum used for the List Channels of Channel Group operation.
	PNChannelsForGroupOperation
	// PNPushNotificationsEnabledChannelsOperation is the enum used for the List Channels with Push Notifications enabled operation.
	PNPushNotificationsEnabledChannelsOperation
	// PNAddPushNotificationsOnChannelsOperation is the enum used for the Add Channels to Push Notifications operation.
	PNAddPushNotificationsOnChannelsOperation
	// PNRemovePushNotificationsFromChannelsOperation is the enum used for the Remove Channels from Push Notifications operation.
	PNRemovePushNotificationsFromChannelsOperation
	// PNRemoveAllPushNotificationsOperation is the enum used for the Remove All Channels from Push Notifications operation.
	PNRemoveAllPushNotificationsOperation
	// PNTimeOperation is the enum used for the Time operation.
	PNTimeOperation
	// PNAccessManagerGrant is the enum used for the Access Manager Grant operation.
	PNAccessManagerGrant
	// PNAccessManagerRevoke is the enum used for the Access Manager Revoke operation.
	PNAccessManagerRevoke
	// PNDeleteMessagesOperation is the enum used for the Delete Messages from History operation.
	PNDeleteMessagesOperation
	// PNMessageCountsOperation is the enum used for History with messages operation.
	PNMessageCountsOperation
	// PNSignalOperation is the enum used for Signal opertaion.
	PNSignalOperation
	// PNCreateUserOperation is the enum used to create users in the Object API.
	// ENUM ORDER needs to be maintained for Objects AIP
	PNCreateUserOperation
	// PNGetUsersOperation is the enum used to get users in the Object API.
	PNGetUsersOperation
	// PNGetUserOperation is the enum used to get user in the Object API.
	PNGetUserOperation
	// PNUpdateUserOperation is the enum used to update users in the Object API.
	PNUpdateUserOperation
	// PNDeleteUserOperation is the enum used to delete users in the Object API.
	PNDeleteUserOperation
	// PNGetSpaceOperation is the enum used to get space in the Object API.
	PNGetSpaceOperation
	// PNGetSpacesOperation is the enum used to get spaces in the Object API.
	PNGetSpacesOperation
	// PNCreateSpaceOperation is the enum used to create space in the Object API.
	PNCreateSpaceOperation
	// PNDeleteSpaceOperation is the enum used to delete space in the Object API.
	PNDeleteSpaceOperation
	// PNUpdateSpaceOperation is the enum used to update space in the Object API.
	PNUpdateSpaceOperation
	// PNGetMembershipsOperation is the enum used to get memberships in the Object API.
	PNGetMembershipsOperation
	// PNGetChannelMembersOperation is the enum used to get members in the Object API.
	PNGetChannelMembersOperation
	// PNManageMembershipsOperation is the enum used to manage memberships in the Object API.
	PNManageMembershipsOperation
	// PNManageMembersOperation is the enum used to manage members in the Object API.
	// ENUM ORDER needs to be maintained for Objects API.
	PNManageMembersOperation
	// PNSetChannelMembersOperation is the enum used to Set Members in the Object API.
	PNSetChannelMembersOperation
	// PNSetMembershipsOperation is the enum used to Set Memberships in the Object API.
	PNSetMembershipsOperation
	// PNRemoveChannelMetadataOperation is the enum used to Remove Channel Metadata in the Object API.
	PNRemoveChannelMetadataOperation
	// PNRemoveUUIDMetadataOperation is the enum used to Remove UUID Metadata in the Object API.
	PNRemoveUUIDMetadataOperation
	// PNGetAllChannelMetadataOperation is the enum used to Get All Channel Metadata in the Object API.
	PNGetAllChannelMetadataOperation
	// PNGetAllUUIDMetadataOperation is the enum used to Get All UUID Metadata in the Object API.
	PNGetAllUUIDMetadataOperation
	// PNGetUUIDMetadataOperation is the enum used to Get UUID Metadata in the Object API.
	PNGetUUIDMetadataOperation
	// PNRemoveMembershipsOperation is the enum used to Remove Memberships in the Object API.
	PNRemoveMembershipsOperation
	// PNRemoveChannelMembersOperation is the enum used to Remove Members in the Object API.
	PNRemoveChannelMembersOperation
	// PNSetUUIDMetadataOperation is the enum used to Set UUID Metadata in the Object API.
	PNSetUUIDMetadataOperation
	// PNSetChannelMetadataOperation is the enum used to Set Channel Metadata in the Object API.
	PNSetChannelMetadataOperation
	// PNGetChannelMetadataOperation is the enum used to Get Channel Metadata in the Object API.
	PNGetChannelMetadataOperation
	// PNAccessManagerGrantToken is the enum used for Grant v3 requests.
	PNAccessManagerGrantToken
	// PNGetMessageActionsOperation is the enum used for Message Actions Get requests.
	PNGetMessageActionsOperation
	// PNHistoryWithActionsOperation is the enum used for History with Actions requests.
	PNHistoryWithActionsOperation
	// PNAddMessageActionsOperation is the enum used for Message Actions Add requests.
	PNAddMessageActionsOperation
	// PNRemoveMessageActionsOperation is the enum used for Message Actions Remove requests.
	PNRemoveMessageActionsOperation
	// PNDeleteFileOperation is the enum used for DeleteFile requests.
	PNDeleteFileOperation
	// PNDownloadFileOperation is the enum used for DownloadFile requests.
	PNDownloadFileOperation
	// PNGetFileURLOperation is the enum used for GetFileURL requests.
	PNGetFileURLOperation
	// PNListFilesOperation is the enum used for ListFiles requests.
	PNListFilesOperation
	// PNSendFileOperation is the enum used for SendFile requests.
	PNSendFileOperation
	// PNSendFileToS3Operation is the enum used for v requests.
	PNSendFileToS3Operation
	// PNPublishFileMessageOperation is the enum used for PublishFileMessage requests.
	PNPublishFileMessageOperation
)

const (
	// PNPushTypeNone is used as an enum to for selecting `none` as the PNPushType
	PNPushTypeNone PNPushType = 1 + iota
	// PNPushTypeGCM is used as an enum to for selecting `GCM` as the PNPushType
	PNPushTypeGCM
	// PNPushTypeAPNS is used as an enum to for selecting `APNS` as the PNPushType
	PNPushTypeAPNS
	// PNPushTypeMPNS is used as an enum to for selecting `MPNS` as the PNPushType
	PNPushTypeMPNS
	// PNPushTypeAPNS2 is used as an enum to for selecting `APNS2` as the PNPushType
	PNPushTypeAPNS2
)

func (p PNPushType) String() string {
	switch p {
	case PNPushTypeAPNS:
		return "apns"

	case PNPushTypeGCM:
		return "gcm"

	case PNPushTypeMPNS:
		return "mpns"

	case PNPushTypeAPNS2:
		return "apns2"

	default:
		return "none"

	}
}

var operations = [...]string{
	"Subscribe",
	"Unsubscribe",
	"Publish",
	"History",
	"Fetch Messages",
	"Where Now",
	"Here Now",
	"Heartbeat",
	"Set State",
	"Get State",
	"Add Channel To Channel Group",
	"Remove Channel From Channel Group",
	"Remove Channel Group",
	"List Channels In Channel Group",
	"List Push Enabled Channels",
	"Add Push From Channel",
	"Remove Push From Channel",
	"Remove All Push Notifications",
	"Time",
	"Grant",
	"Revoke",
	"Delete messages",
	"Signal",
	"Create User",
	"Get Users",
	"Fetch User",
	"Update User",
	"Delete User",
	"Get Space",
	"Get Spaces",
	"Create Space",
	"Delete Space",
	"Update Space",
	"GetMemberships",
	"GetChannelMembers",
	"ManageMemberships",
	"ManageMembers",
	"SetChannelMembers",
	"SetMemberships",
	"RemoveChannelMetadata",
	"RemoveUUIDMetadata",
	"GetAllChannelMetadata",
	"GetAllUUIDMetadata",
	"GetUUIDMetadata",
	"RemoveMemberships",
	"RemoveChannelMembers",
	"SetUUIDMetadata",
	"SetChannelMetadata",
	"GetChannelMetadata",
	"Grant Token",
	"GetMessageActions",
	"HistoryWithActions",
	"AddMessageActions",
	"RemoveMessageActions",
	"DeleteFile",
	"DownloadFile",
	"GetFileURL",
	"ListFiles",
	"SendFile",
	"SendFileToS3",
	"PublishFile",
}

func (c StatusCategory) String() string {
	switch c {
	case PNUnknownCategory:
		return "Unknown"

	case PNTimeoutCategory:
		return "Timeout"

	case PNConnectedCategory:
		return "Connected"

	case PNDisconnectedCategory:
		return "Disconnected"

	case PNCancelledCategory:
		return "Cancelled"

	case PNLoopStopCategory:
		return "Loop Stop"

	case PNAcknowledgmentCategory:
		return "Acknowledgment"

	case PNBadRequestCategory:
		return "Bad Request"

	case PNAccessDeniedCategory:
		return "Access Denied"

	case PNReconnectedCategory:
		return "Reconnected"

	case PNReconnectionAttemptsExhausted:
		return "Reconnection Attempts Exhausted"

	case PNNoStubMatchedCategory:
		return "No Stub Matched"

	default:
		return "No Stub Matched"

	}
}

func (t OperationType) String() string {
	switch t {
	case PNSubscribeOperation:
		return "Subscribe"

	case PNUnsubscribeOperation:
		return "Unsubscribe"

	case PNPublishOperation:
		return "Publish"

	case PNFireOperation:
		return "Fire"

	case PNHistoryOperation:
		return "History"

	case PNFetchMessagesOperation:
		return "Fetch Messages"

	case PNWhereNowOperation:
		return "Where Now"

	case PNHereNowOperation:
		return "Here Now"

	case PNHeartBeatOperation:
		return "Heartbeat"

	case PNSetStateOperation:
		return "Set State"

	case PNGetStateOperation:
		return "Get State"

	case PNAddChannelsToChannelGroupOperation:
		return "Add Channel To Channel Group"

	case PNRemoveChannelFromChannelGroupOperation:
		return "Remove Channel From Channel Group"

	case PNRemoveGroupOperation:
		return "Remove Channel Group"

	case PNChannelsForGroupOperation:
		return "List Channels In Channel Group"

	case PNPushNotificationsEnabledChannelsOperation:
		return "List Push Enabled Channels"

	case PNAddPushNotificationsOnChannelsOperation:
		return "Add Push From Channel"

	case PNRemovePushNotificationsFromChannelsOperation:
		return "Remove Push From Channel"

	case PNRemoveAllPushNotificationsOperation:
		return "Remove All Push Notifications"

	case PNTimeOperation:
		return "Time"

	case PNAccessManagerGrant:
		return "Grant"

	case PNAccessManagerRevoke:
		return "Revoke"

	case PNDeleteMessagesOperation:
		return "Delete messages"

	case PNSignalOperation:
		return "Signal"

	case PNCreateUserOperation:
		return "Create User"
	case PNGetUsersOperation:
		return "Get Users"
	case PNGetUserOperation:
		return "Fetch Users"
	case PNUpdateUserOperation:
		return "Update User"
	case PNDeleteUserOperation:
		return "Delete User"
	case PNGetSpaceOperation:
		return "Get Space"
	case PNGetSpacesOperation:
		return "Get Spaces"
	case PNCreateSpaceOperation:
		return "Create Space"
	case PNDeleteSpaceOperation:
		return "Delete Space"
	case PNUpdateSpaceOperation:
		return "Update Space"
	case PNGetMembershipsOperation:
		return "Get Memberships V2"
	case PNGetChannelMembersOperation:
		return "Get Members V2"
	case PNManageMembershipsOperation:
		return "Manage Memberships V2"
	case PNManageMembersOperation:
		return "Manage Members V2"
	case PNSetChannelMembersOperation:
		return "Set Members V2"
	case PNSetMembershipsOperation:
		return "Set Memberships V2"
	case PNRemoveChannelMetadataOperation:
		return "Remove Channel Metadata V2"
	case PNRemoveUUIDMetadataOperation:
		return "Remove Metadata V2"
	case PNGetAllChannelMetadataOperation:
		return "Get All Channel Metadata V2"
	case PNGetAllUUIDMetadataOperation:
		return "Get All UUID Metadata V2"
	case PNGetUUIDMetadataOperation:
		return "Get UUID Metadata V2"
	case PNRemoveMembershipsOperation:
		return "Remove Memberships V2"
	case PNRemoveChannelMembersOperation:
		return "Remove Members V2"
	case PNSetUUIDMetadataOperation:
		return "Set UUID Metadata V2"
	case PNSetChannelMetadataOperation:
		return "Set Channel Metadata V2"
	case PNGetChannelMetadataOperation:
		return "Get Channel Metadata V2"
	case PNAccessManagerGrantToken:
		return "Grant Token"
	case PNGetMessageActionsOperation:
		return "Get Message Actions"
	case PNHistoryWithActionsOperation:
		return "History With Actions"
	case PNAddMessageActionsOperation:
		return "Add Message Actions"
	case PNRemoveMessageActionsOperation:
		return "Remove Message Actions"
	case PNDeleteFileOperation:
		return "Delete File"
	case PNDownloadFileOperation:
		return "Download File"
	case PNGetFileURLOperation:
		return "Get File URL"
	case PNListFilesOperation:
		return "List Files"
	case PNSendFileOperation:
		return "Send File"
	case PNSendFileToS3Operation:
		return "Send File To S3"
	case PNPublishFileMessageOperation:
		return "Publish File"
	default:
		return "No Category Matched"
	}
}

// EnumArrayToStringArray converts a string enum to an array
func EnumArrayToStringArray(include interface{}) []string {
	s := []string{}
	switch fmt.Sprintf("%s", reflect.TypeOf(include)) {
	case "[]pubnub.PNChannelMembersInclude":
		for _, v := range include.([]PNChannelMembersInclude) {
			s = append(s, fmt.Sprintf("%s", v))
		}
	case "[]pubnub.PNMembershipsInclude":
		for _, v := range include.([]PNMembershipsInclude) {
			s = append(s, fmt.Sprintf("%s", v))
		}
	case "[]pubnub.PNUUIDMetadataInclude":
		for _, v := range include.([]PNUUIDMetadataInclude) {
			s = append(s, fmt.Sprintf("%s", v))
		}
	case "[]pubnub.PNChannelMetadataInclude":
		for _, v := range include.([]PNChannelMetadataInclude) {
			s = append(s, fmt.Sprintf("%s", v))
		}
	}
	return s
}
