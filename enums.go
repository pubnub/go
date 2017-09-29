package pubnub

type StatusCategory int
type OperationType int

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
	PNChannelGroupOperation
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

func (c StatusCategory) String() string {
	return categories[c-1]
}
