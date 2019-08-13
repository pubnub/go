package pubnub

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

type PNUserSpaceInclude int
type PNSpaceMembershipsIncude int
type PNMembersInclude int

const (
	PNUserSpaceCustom PNUserSpaceInclude = 1 + iota
)

func (s PNUserSpaceInclude) String() string {
	return [...]string{"custom"}[s-1]
}

const (
	PNSpaceMembershipsCustom PNSpaceMembershipsIncude = 1 + iota
	PNSpaceMembershipsSpace
	PNSpaceMembershipsSpaceCustom
)

func (s PNSpaceMembershipsIncude) String() string {
	return [...]string{"custom", "space", "space.custom"}[s-1]
}

const (
	PNMembersCustom PNMembersInclude = 1 + iota
	PNMembersUser
	PNMembersUserCustom
)

func (s PNMembersInclude) String() string {
	return [...]string{"custom", "user", "user.custom"}[s-1]
}

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
	PNCreateUserOperation
	// PNGetUsersOperation is the enum used to get users in the Object API.
	PNGetUsersOperation
	// PNFetchUserOperation
	PNFetchUserOperation
	// PNUpdateUserOperation
	PNUpdateUserOperation
	// PNDeleteUserOperation
	PNDeleteUserOperation
	// PNGetSpaceOperation
	PNGetSpaceOperation
	// PNGetSpacesOperation
	PNGetSpacesOperation
	// PNCreateSpaceOperation
	PNCreateSpaceOperation
	// PNDeleteSpaceOperation
	PNDeleteSpaceOperation
	// PNUpdateSpaceOperation
	PNUpdateSpaceOperation
	// PNGetSpaceMembershipsOperation
	PNGetSpaceMembershipsOperation
	// PNGetMembersOperation
	PNGetMembersOperation
	// PNUpdateSpaceMembershipsOperation
	PNUpdateSpaceMembershipsOperation
	// PNUpdateUserSpaceMembershipsOperation
	PNUpdateUserSpaceMembershipsOperation
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
	"PNGetSpaceMembershipsOperation",
	"PNGetMembersOperation",
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
	case PNFetchUserOperation:
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
	case PNGetSpaceMembershipsOperation:
		return "Get Space Memberships"
	case PNGetMembersOperation:
		return "Get Members"
	default:
		return "No Category Matched"
	}
}
