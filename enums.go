package pubnub

type StatusCategory int
type OperationType int
type ReconnectionPolicy int

// TODO: add prefix
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
)

var categories = [...]string{
	"Unknown",
	"Timeout",
	"Connected",
	"Disconnected",
	"Cancelled",
	"Loop Stop",
	"Acknowledgment",
	"Bad Request",
	"Access Denied",
	"No Stub Matched",
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
}

func (c StatusCategory) String() string {
	return categories[c-1]
}

func (t OperationType) String() string {
	return operations[t-1]
}
