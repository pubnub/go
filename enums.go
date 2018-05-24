package pubnub

type StatusCategory int
type OperationType int
type ReconnectionPolicy int
type PNPushType int

const (
	PNNonePolicy ReconnectionPolicy = 1 + iota
	PNLinearPolicy
	PNExponentialPolicy
)

const (
	PNUnknownCategory StatusCategory = 1 + iota
	// Request timeout reached
	PNTimeoutCategory
	// Subscribe received an initial timetoken
	PNConnectedCategory
	// Disconnected due network error
	PNDisconnectedCategory
	// Context cancelled
	PNCancelledCategory
	PNLoopStopCategory
	PNAcknowledgmentCategory
	PNBadRequestCategory
	PNAccessDeniedCategory
	PNNoStubMatchedCategory
	PNReconnectedCategory
	PNReconnectionAttemptsExhausted
)

const (
	PNSubscribeOperation OperationType = 1 + iota
	PNUnsubscribeOperation
	PNPublishOperation
	PNFireOperation
	PNHistoryOperation
	PNFetchMessagesOperation
	PNWhereNowOperation
	PNHereNowOperation
	PNHeartBeatOperation
	PNSetStateOperation
	PNGetStateOperation
	PNAddChannelsToChannelGroupOperation
	PNRemoveChannelFromChannelGroupOperation
	PNRemoveGroupOperation
	PNChannelsForGroupOperation
	PNPushNotificationsEnabledChannelsOperation
	PNAddPushNotificationsOnChannelsOperation
	PNRemovePushNotificationsFromChannelsOperation
	PNRemoveAllPushNotificationsOperation
	PNTimeOperation
	PNAccessManagerGrant
	PNAccessManagerRevoke
	PNDeleteMessagesOperation
)

const (
	PNPushTypeNone PNPushType = 1 + iota
	PNPushTypeGCM
	PNPushTypeAPNS
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
	return "none"
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
	return "No Stub Matched"
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

	default:
		return "No Category Matched"

	}
	return "No Category Matched"

}
