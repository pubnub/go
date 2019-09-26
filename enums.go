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

// PNUserSpaceInclude is used as an enum to catgorize the available User and Space include types
type PNUserSpaceInclude int

// PNMembershipsInclude is used as an enum to catgorize the available Memberships include types
type PNMembershipsInclude int

// PNMembersInclude is used as an enum to catgorize the available Members include types
type PNMembersInclude int

// PNObjectsEvent is used as an enum to catgorize the available Object Events
type PNObjectsEvent string

// PNObjectsEventType is used as an enum to catgorize the available Object Event types
type PNObjectsEventType string

// PNMessageActionsEventType is used as an enum to catgorize the available Message Actions Event types
type PNMessageActionsEventType string

const (
	// PNMessageActionsAdded is the enum when the event of type `added` occurs
	PNMessageActionsAdded PNMessageActionsEventType = "added"
	// PNMessageActionsRemoved is the enum when the event of type `removed` occurs
	PNMessageActionsRemoved = "removed"
)

const (
	// PNObjectsUserEvent is the enum when the event of type `user` occurs
	PNObjectsUserEvent PNObjectsEventType = "user"
	// PNObjectsSpaceEvent is the enum when the event of type `space` occurs
	PNObjectsSpaceEvent = "space"
	// PNObjectsMembershipEvent is the enum when the event of type `membership` occurs
	PNObjectsMembershipEvent = "membership"
	// PNObjectsNoneEvent is used for error handling
	PNObjectsNoneEvent = "none"
)

const (
	// PNObjectsEventCreate is the enum when the event `create` occurs
	PNObjectsEventCreate PNObjectsEvent = "create"
	// PNObjectsEventUpdate is the enum when the event `update` occurs
	PNObjectsEventUpdate = "update"
	// PNObjectsEventDelete is the enum when the event `delete` occurs
	PNObjectsEventDelete = "delete"
)

const (
	// PNUserSpaceCustom is the enum equivalent to the value `custom` available User and Space include types
	PNUserSpaceCustom PNUserSpaceInclude = 1 + iota
)

func (s PNUserSpaceInclude) String() string {
	return [...]string{"custom"}[s-1]
}

const (
	// PNMembershipsCustom is the enum equivalent to the value `custom` available Memberships include types
	PNMembershipsCustom PNMembershipsInclude = 1 + iota
	// PNMembershipsSpace is the enum equivalent to the value `space` available Memberships include types
	PNMembershipsSpace
	// PNMembershipsSpaceCustom is the enum equivalent to the value `space.custom` available Memberships include types
	PNMembershipsSpaceCustom
)

func (s PNMembershipsInclude) String() string {
	return [...]string{"custom", "space", "space.custom"}[s-1]
}

const (
	// PNMembersCustom is the enum equivalent to the value `custom` available Members include types
	PNMembersCustom PNMembersInclude = 1 + iota
	// PNMembersUser is the enum equivalent to the value `user` available Members include types
	PNMembersUser
	// PNMembersUserCustom is the enum equivalent to the value `user.custom` available Members include types
	PNMembersUserCustom
)

func (s PNMembersInclude) String() string {
	return [...]string{"custom", "user", "user.custom"}[s-1]
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
	// PNGetMembersOperation is the enum used to get members in the Object API.
	PNGetMembersOperation
	// PNManageMembershipsOperation is the enum used to manage memberships in the Object API.
	PNManageMembershipsOperation
	// PNManageMembersOperation is the enum used to manage members in the Object API.
	// ENUM ORDER needs to be maintained for Objects API
	PNManageMembersOperation
	// PNAccessManagerGrantToken is the enum used from Grant v3 requests
	PNAccessManagerGrantToken
	PNGetMessageActionsOperation
	PNHistoryWithActionsOperation
	PNAddMessageActionsOperation
	PNRemoveMessageActionsOperation
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
)

func (p PNPushType) String() string {
	switch p {
	case PNPushTypeAPNS:
		return "apns"

	case PNPushTypeGCM:
		return "gcm"

	case PNPushTypeMPNS:
		return "mpns"

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
	"PNGetMembershipsOperation",
	"PNGetMembersOperation",
	"PNManageMembershipsOperation",
	"PNManageMembersOperation",
	"GrantToken",
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
		return "Get Memberships"
	case PNGetMembersOperation:
		return "Get Members"
	case PNManageMembershipsOperation:
		return "Manage Memberships"
	case PNManageMembersOperation:
		return "Manage Members"
	case PNAccessManagerGrantToken:
		return "Grant Token"
	default:
		return "No Category Matched"
	}
}

// EnumArrayToStringArray converts a string enum to an array
func EnumArrayToStringArray(include interface{}) []string {
	s := []string{}
	switch fmt.Sprintf("%s", reflect.TypeOf(include)) {
	case "[]pubnub.PNMembersInclude":
		for _, v := range include.([]PNMembersInclude) {
			s = append(s, fmt.Sprintf("%s", v))
		}
	case "[]pubnub.PNMembershipsInclude":
		for _, v := range include.([]PNMembershipsInclude) {
			s = append(s, fmt.Sprintf("%s", v))
		}
	case "[]pubnub.PNUserSpaceInclude":
		for _, v := range include.([]PNUserSpaceInclude) {
			s = append(s, fmt.Sprintf("%s", v))
		}
	}
	return s
}
