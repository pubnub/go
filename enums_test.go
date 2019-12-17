package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushString(t *testing.T) {
	assert := assert.New(t)

	pushAPNS := PNPushTypeAPNS
	pushAPNS2 := PNPushTypeAPNS2
	pushMPNS := PNPushTypeMPNS
	pushGCM := PNPushTypeGCM
	pushNONE := PNPushTypeNone

	assert.Equal("apns", pushAPNS.String())
	assert.Equal("apns2", pushAPNS2.String())
	assert.Equal("mpns", pushMPNS.String())
	assert.Equal("gcm", pushGCM.String())
	assert.Equal("none", pushNONE.String())
}

func TestStatusCategoryString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("Unknown", PNUnknownCategory.String())
	assert.Equal("Timeout", PNTimeoutCategory.String())
	assert.Equal("Connected", PNConnectedCategory.String())
	assert.Equal("Disconnected", PNDisconnectedCategory.String())
	assert.Equal("Cancelled", PNCancelledCategory.String())
	assert.Equal("Loop Stop", PNLoopStopCategory.String())
	assert.Equal("Acknowledgment", PNAcknowledgmentCategory.String())
	assert.Equal("Bad Request", PNBadRequestCategory.String())
	assert.Equal("Access Denied", PNAccessDeniedCategory.String())
	assert.Equal("Reconnected", PNReconnectedCategory.String())
	assert.Equal("Reconnection Attempts Exhausted", PNReconnectionAttemptsExhausted.String())
	assert.Equal("No Stub Matched", PNNoStubMatchedCategory.String())
}

func TestOperationTypeString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("Subscribe", PNSubscribeOperation.String())
	assert.Equal("Unsubscribe", PNUnsubscribeOperation.String())
	assert.Equal("Publish", PNPublishOperation.String())
	assert.Equal("Fire", PNFireOperation.String())
	assert.Equal("History", PNHistoryOperation.String())
	assert.Equal("Fetch Messages", PNFetchMessagesOperation.String())
	assert.Equal("Where Now", PNWhereNowOperation.String())
	assert.Equal("Here Now", PNHereNowOperation.String())
	assert.Equal("Heartbeat", PNHeartBeatOperation.String())
	assert.Equal("Set State", PNSetStateOperation.String())
	assert.Equal("Get State", PNGetStateOperation.String())
	assert.Equal("Add Channel To Channel Group", PNAddChannelsToChannelGroupOperation.String())
	assert.Equal("Remove Channel From Channel Group", PNRemoveChannelFromChannelGroupOperation.String())
	assert.Equal("Remove Channel Group", PNRemoveGroupOperation.String())
	assert.Equal("List Channels In Channel Group", PNChannelsForGroupOperation.String())
	assert.Equal("List Push Enabled Channels", PNPushNotificationsEnabledChannelsOperation.String())
	assert.Equal("Add Push From Channel", PNAddPushNotificationsOnChannelsOperation.String())
	assert.Equal("Remove Push From Channel", PNRemovePushNotificationsFromChannelsOperation.String())
	assert.Equal("Remove All Push Notifications", PNRemoveAllPushNotificationsOperation.String())
	assert.Equal("Time", PNTimeOperation.String())
	assert.Equal("Grant", PNAccessManagerGrant.String())
	assert.Equal("Revoke", PNAccessManagerRevoke.String())
	assert.Equal("Delete messages", PNDeleteMessagesOperation.String())
}
