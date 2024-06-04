package pubnub

import (
	"fmt"
	"github.com/pubnub/go/v7/crypto"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/pubnub/go/v7/utils"
)

// Default constants
const (
	// Version :the version of the SDK
	Version = "7.2.1"
	// MaxSequence for publish messages
	MaxSequence = 65535
)

const (
	// StrMissingPubKey shows Missing Publish Key message
	StrMissingPubKey = "Missing Publish Key"
	// StrMissingSubKey shows Missing Subscribe Key message
	StrMissingSubKey = "Missing Subscribe Key"
	// StrMissingChannel shows Channel message
	StrMissingChannel = "Missing Channel"
	// StrMissingChannelGroup shows Channel Group message
	StrMissingChannelGroup = "Missing Channel Group"
	// StrMissingMessage shows Missing Message message
	StrMissingMessage = "Missing Message"
	// StrMissingSecretKey shows Missing Secret Key message
	StrMissingSecretKey = "Missing Secret Key"
	// StrMissingUUID shows Missing UUID message
	StrMissingUUID = "Missing UUID"
	// StrMissingDeviceID shows Missing Device ID message
	StrMissingDeviceID = "Missing Device ID"
	// StrMissingPushType shows Missing Push Type message
	StrMissingPushType = "Missing Push Type"
	// StrMissingPushTopic shows Missing Push Topic message
	StrMissingPushTopic = "Missing Push Topic"
	// StrChannelsTimetoken shows Missing Channels Timetoken message
	StrChannelsTimetoken = "Missing Channels Timetoken"
	// StrChannelsTimetokenLength shows Length of Channels Timetoken message
	StrChannelsTimetokenLength = "Length of Channels Timetoken and Channels do not match"
	// StrInvalidTTL shows Invalid TTL message
	StrInvalidTTL = "Invalid TTL"
	// StrMissingPushTitle shows `Push title missing` message
	StrMissingPushTitle = "Push title missing"
	// StrMissingFileID shows `Missing File ID` message
	StrMissingFileID = "Missing File ID"
	// StrMissingFileName shows `Missing File Name` message
	StrMissingFileName = "Missing File Name"
	// StrMissingToken shows `Missing PAMv3 token` message
	StrMissingToken = "Missing PAMv3 token"
)

// PubNub No server connection will be established when you create a new PubNub object.
// To establish a new connection use Subscribe() function of PubNub type.
type PubNub struct {
	sync.RWMutex

	Config               *Config
	nextPublishSequence  int
	publishSequenceMutex sync.RWMutex
	subscriptionManager  *SubscriptionManager
	telemetryManager     *TelemetryManager
	heartbeatManager     *HeartbeatManager
	client               *http.Client
	subscribeClient      *http.Client
	requestWorkers       *RequestWorkers
	jobQueue             chan *JobQItem
	ctx                  Context
	cancel               func()
	tokenManager         *TokenManager
	previousCipherKey    string
	previousIvFlag       bool
}

// TODO this needs to be tested
func (pn *PubNub) getCryptoModule() crypto.CryptoModule {
	pn.Lock()
	defer pn.Unlock()
	if pn.previousCipherKey == pn.Config.CipherKey && pn.previousIvFlag == pn.Config.UseRandomInitializationVector {
		return pn.Config.CryptoModule
	}

	if pn.Config != nil && pn.Config.CipherKey != "" {
		pn.Config.CryptoModule, _ = crypto.NewLegacyCryptoModule(pn.Config.CipherKey, pn.Config.UseRandomInitializationVector)
		return pn.Config.CryptoModule
	} else if pn.Config != nil && pn.Config.CipherKey == "" {
		pn.Config.CryptoModule = nil
		return pn.Config.CryptoModule
	}
	return nil
}

// Publish is used to send a message to all subscribers of a channel.
func (pn *PubNub) Publish() *publishBuilder {
	return newPublishBuilder(pn)
}

// PublishWithContext function is used to send a message to all subscribers of a channel.
func (pn *PubNub) PublishWithContext(ctx Context) *publishBuilder {
	return newPublishBuilderWithContext(pn, ctx)
}

// Fire endpoint allows the client to send a message to PubNub Functions Event Handlers. These messages will go directly to any Event Handlers registered on the channel that you fire to and will trigger their execution.
func (pn *PubNub) Fire() *fireBuilder {
	return newFireBuilder(pn)
}

// FireWithContext endpoint allows the client to send a message to PubNub Functions Event Handlers. These messages will go directly to any Event Handlers registered on the channel that you fire to and will trigger their execution.
func (pn *PubNub) FireWithContext(ctx Context) *fireBuilder {
	return newFireBuilderWithContext(pn, ctx)
}

// Subscribe causes the client to create an open TCP socket to the PubNub Real-Time Network and begin listening for messages on a specified channel.
func (pn *PubNub) Subscribe() *subscribeBuilder {
	return newSubscribeBuilder(pn)
}

// History fetches historical messages of a channel.
func (pn *PubNub) History() *historyBuilder {
	return newHistoryBuilder(pn)
}

// HistoryWithContext fetches historical messages of a channel.
func (pn *PubNub) HistoryWithContext(ctx Context) *historyBuilder {
	return newHistoryBuilderWithContext(pn, ctx)
}

// Fetch fetches historical messages from multiple channels.
func (pn *PubNub) Fetch() *fetchBuilder {
	return newFetchBuilder(pn)
}

// FetchWithContext fetches historical messages from multiple channels.
func (pn *PubNub) FetchWithContext(ctx Context) *fetchBuilder {
	return newFetchBuilderWithContext(pn, ctx)
}

// MessageCounts Returns the number of messages published on one or more channels since a given time.
func (pn *PubNub) MessageCounts() *messageCountsBuilder {
	return newMessageCountsBuilder(pn)
}

// MessageCountsWithContext Returns the number of messages published on one or more channels since a given time.
func (pn *PubNub) MessageCountsWithContext(ctx Context) *messageCountsBuilder {
	return newMessageCountsBuilderWithContext(pn, ctx)
}

// GetAllUUIDMetadata Returns a paginated list of UUID Metadata objects, optionally including the custom data object for each.
func (pn *PubNub) GetAllUUIDMetadata() *getAllUUIDMetadataBuilder {
	return newGetAllUUIDMetadataBuilder(pn)
}

// GetAllUUIDMetadataWithContext Returns a paginated list of UUID Metadata objects, optionally including the custom data object for each.
func (pn *PubNub) GetAllUUIDMetadataWithContext(ctx Context) *getAllUUIDMetadataBuilder {
	return newGetAllUUIDMetadataBuilderWithContext(pn, ctx)
}

// GetUUIDMetadata Returns metadata for the specified UUID, optionally including the custom data object for each.
func (pn *PubNub) GetUUIDMetadata() *getUUIDMetadataBuilder {
	return newGetUUIDMetadataBuilder(pn)
}

// GetUUIDMetadataWithContext Returns metadata for the specified UUID, optionally including the custom data object for each.
func (pn *PubNub) GetUUIDMetadataWithContext(ctx Context) *getUUIDMetadataBuilder {
	return newGetUUIDMetadataBuilderWithContext(pn, ctx)
}

// SetUUIDMetadata Set metadata for a UUID in the database, optionally including the custom data object for each.
func (pn *PubNub) SetUUIDMetadata() *setUUIDMetadataBuilder {
	return newSetUUIDMetadataBuilder(pn)
}

// SetUUIDMetadataWithContext Set metadata for a UUID in the database, optionally including the custom data object for each.
func (pn *PubNub) SetUUIDMetadataWithContext(ctx Context) *setUUIDMetadataBuilder {
	return newSetUUIDMetadataBuilderWithContext(pn, ctx)
}

// RemoveUUIDMetadata Removes the metadata from a specified UUID.
func (pn *PubNub) RemoveUUIDMetadata() *removeUUIDMetadataBuilder {
	return newRemoveUUIDMetadataBuilder(pn)
}

// RemoveUUIDMetadataWithContext Removes the metadata from a specified UUID.
func (pn *PubNub) RemoveUUIDMetadataWithContext(ctx Context) *removeUUIDMetadataBuilder {
	return newRemoveUUIDMetadataBuilderWithContext(pn, ctx)
}

// GetAllChannelMetadata Returns a paginated list of Channel Metadata objects, optionally including the custom data object for each.
func (pn *PubNub) GetAllChannelMetadata() *getAllChannelMetadataBuilder {
	return newGetAllChannelMetadataBuilder(pn)
}

// GetAllChannelMetadataWithContext Returns a paginated list of Channel Metadata objects, optionally including the custom data object for each.
func (pn *PubNub) GetAllChannelMetadataWithContext(ctx Context) *getAllChannelMetadataBuilder {
	return newGetAllChannelMetadataBuilderWithContext(pn, ctx)
}

// GetChannelMetadata Returns metadata for the specified Channel, optionally including the custom data object for each.
func (pn *PubNub) GetChannelMetadata() *getChannelMetadataBuilder {
	return newGetChannelMetadataBuilder(pn)
}

// GetChannelMetadataWithContext Returns metadata for the specified Channel, optionally including the custom data object for each.
func (pn *PubNub) GetChannelMetadataWithContext(ctx Context) *getChannelMetadataBuilder {
	return newGetChannelMetadataBuilderWithContext(pn, ctx)
}

// SetChannelMetadata Set metadata for a Channel in the database, optionally including the custom data object for each.
func (pn *PubNub) SetChannelMetadata() *setChannelMetadataBuilder {
	return newSetChannelMetadataBuilder(pn)
}

// SetChannelMetadataWithContext Set metadata for a Channel in the database, optionally including the custom data object for each.
func (pn *PubNub) SetChannelMetadataWithContext(ctx Context) *setChannelMetadataBuilder {
	return newSetChannelMetadataBuilderWithContext(pn, ctx)
}

// RemoveChannelMetadata Removes the metadata from a specified channel.
func (pn *PubNub) RemoveChannelMetadata() *removeChannelMetadataBuilder {
	return newRemoveChannelMetadataBuilder(pn)
}

// RemoveChannelMetadataWithContext Removes the metadata from a specified channel.
func (pn *PubNub) RemoveChannelMetadataWithContext(ctx Context) *removeChannelMetadataBuilder {
	return newRemoveChannelMetadataBuilderWithContext(pn, ctx)
}

// GetMemberships The method returns a list of channel memberships for a user. This method doesn't return a user's subscriptions.
func (pn *PubNub) GetMemberships() *getMembershipsBuilderV2 {
	return newGetMembershipsBuilderV2(pn)
}

// GetMembershipsWithContext The method returns a list of channel memberships for a user. This method doesn't return a user's subscriptions.
func (pn *PubNub) GetMembershipsWithContext(ctx Context) *getMembershipsBuilderV2 {
	return newGetMembershipsBuilderV2WithContext(pn, ctx)
}

// GetChannelMembers The method returns a list of members in a channel. The list will include user metadata for members that have additional metadata stored in the database.
func (pn *PubNub) GetChannelMembers() *getChannelMembersBuilderV2 {
	return newGetChannelMembersBuilderV2(pn)
}

// GetChannelMembersWithContext The method returns a list of members in a channel. The list will include user metadata for members that have additional metadata stored in the database.
func (pn *PubNub) GetChannelMembersWithContext(ctx Context) *getChannelMembersBuilderV2 {
	return newGetChannelMembersBuilderV2WithContext(pn, ctx)
}

// SetChannelMembers This method sets members in a channel.
func (pn *PubNub) SetChannelMembers() *setChannelMembersBuilder {
	return newSetChannelMembersBuilder(pn)
}

// SetChannelMembersWithContext This method sets members in a channel.
func (pn *PubNub) SetChannelMembersWithContext(ctx Context) *setChannelMembersBuilder {
	return newSetChannelMembersBuilderWithContext(pn, ctx)
}

// RemoveChannelMembers Remove members from a Channel.
func (pn *PubNub) RemoveChannelMembers() *removeChannelMembersBuilder {
	return newRemoveChannelMembersBuilder(pn)
}

// RemoveChannelMembersWithContext Remove members from a Channel.
func (pn *PubNub) RemoveChannelMembersWithContext(ctx Context) *removeChannelMembersBuilder {
	return newRemoveChannelMembersBuilderWithContext(pn, ctx)
}

// SetMemberships Set channel memberships for a UUID.
func (pn *PubNub) SetMemberships() *setMembershipsBuilder {
	return newSetMembershipsBuilder(pn)
}

// SetMembershipsWithContext Set channel memberships for a UUID.
func (pn *PubNub) SetMembershipsWithContext(ctx Context) *setMembershipsBuilder {
	return newSetMembershipsBuilderWithContext(pn, ctx)
}

// RemoveMemberships Remove channel memberships for a UUID.
func (pn *PubNub) RemoveMemberships() *removeMembershipsBuilder {
	return newRemoveMembershipsBuilder(pn)
}

// RemoveMembershipsWithContext Remove channel memberships for a UUID.
func (pn *PubNub) RemoveMembershipsWithContext(ctx Context) *removeMembershipsBuilder {
	return newRemoveMembershipsBuilderWithContext(pn, ctx)
}

// ManageChannelMembers The method Set and Remove channel memberships for a user.
func (pn *PubNub) ManageChannelMembers() *manageChannelMembersBuilderV2 {
	return newManageChannelMembersBuilderV2(pn)
}

// ManageChannelMembersWithContext The method Set and Remove channel memberships for a user.
func (pn *PubNub) ManageChannelMembersWithContext(ctx Context) *manageChannelMembersBuilderV2 {
	return newManageChannelMembersBuilderV2WithContext(pn, ctx)
}

// ManageMemberships Manage the specified UUID's memberships. You can Add, Remove, and Update a UUID's memberships.
func (pn *PubNub) ManageMemberships() *manageMembershipsBuilderV2 {
	return newManageMembershipsBuilderV2(pn)
}

// ManageMembershipsWithContext Manage the specified UUID's memberships. You can Add, Remove, and Update a UUID's memberships.
func (pn *PubNub) ManageMembershipsWithContext(ctx Context) *manageMembershipsBuilderV2 {
	return newManageMembershipsBuilderV2WithContext(pn, ctx)
}

// Signal The signal() function is used to send a signal to all subscribers of a channel.
func (pn *PubNub) Signal() *signalBuilder {
	return newSignalBuilder(pn)
}

// SignalWithContext The signal() function is used to send a signal to all subscribers of a channel.
func (pn *PubNub) SignalWithContext(ctx Context) *signalBuilder {
	return newSignalBuilderWithContext(pn, ctx)
}

// SetState The state API is used to set/get key/value pairs specific to a subscriber UUID. State information is supplied as a JSON object of key/value pairs.
func (pn *PubNub) SetState() *setStateBuilder {
	return newSetStateBuilder(pn)
}

// SetStateWithContext The state API is used to set/get key/value pairs specific to a subscriber UUID. State information is supplied as a JSON object of key/value pairs.
func (pn *PubNub) SetStateWithContext(ctx Context) *setStateBuilder {
	return newSetStateBuilderWithContext(pn, ctx)
}

// Grant This function establishes access permissions for PubNub Access Manager (PAM) by setting the read or write attribute to true. A grant with read or write set to false (or not included) will revoke any previous grants with read or write set to true.
func (pn *PubNub) Grant() *grantBuilder {
	return newGrantBuilder(pn)
}

// GrantWithContext This function establishes access permissions for PubNub Access Manager (PAM) by setting the read or write attribute to true. A grant with read or write set to false (or not included) will revoke any previous grants with read or write set to true.
func (pn *PubNub) GrantWithContext(ctx Context) *grantBuilder {
	return newGrantBuilderWithContext(pn, ctx)
}

// GrantToken Use the Grant Token method to generate an auth token with embedded access control lists. The client sends the auth token to PubNub along with each request.
func (pn *PubNub) GrantToken() *grantTokenBuilder {
	return newGrantTokenBuilder(pn)
}

// GrantTokenWithContext Use the Grant Token method to generate an auth token with embedded access control lists. The client sends the auth token to PubNub along with each request.
func (pn *PubNub) GrantTokenWithContext(ctx Context) *grantTokenBuilder {
	return newGrantTokenBuilderWithContext(pn, ctx)
}

// RevokeToken Use the Grant Token method to generate an auth token with embedded access control lists. The client sends the auth token to PubNub along with each request.
func (pn *PubNub) RevokeToken() *revokeTokenBuilder {
	return newRevokeTokenBuilder(pn)
}

// RevokeTokenWithContext Use the Grant Token method to generate an auth token with embedded access control lists. The client sends the auth token to PubNub along with each request.
func (pn *PubNub) RevokeTokenWithContext(ctx Context) *revokeTokenBuilder {
	return newRevokeTokenBuilderWithContext(pn, ctx)
}

// AddMessageAction Add an action on a published message. Returns the added action in the response.
func (pn *PubNub) AddMessageAction() *addMessageActionsBuilder {
	return newAddMessageActionsBuilder(pn)
}

// AddMessageActionWithContext Add an action on a published message. Returns the added action in the response.
func (pn *PubNub) AddMessageActionWithContext(ctx Context) *addMessageActionsBuilder {
	return newAddMessageActionsBuilderWithContext(pn, ctx)
}

// GetMessageActions Get a list of message actions in a channel. Returns a list of actions in the response.
func (pn *PubNub) GetMessageActions() *getMessageActionsBuilder {
	return newGetMessageActionsBuilder(pn)
}

// GetMessageActionsWithContext Get a list of message actions in a channel. Returns a list of actions in the response.
func (pn *PubNub) GetMessageActionsWithContext(ctx Context) *getMessageActionsBuilder {
	return newGetMessageActionsBuilderWithContext(pn, ctx)
}

// RemoveMessageAction Remove a peviously added action on a published message. Returns an empty response.
func (pn *PubNub) RemoveMessageAction() *removeMessageActionsBuilder {
	return newRemoveMessageActionsBuilder(pn)
}

// RemoveMessageActionWithContext Remove a peviously added action on a published message. Returns an empty response.
func (pn *PubNub) RemoveMessageActionWithContext(ctx Context) *removeMessageActionsBuilder {
	return newRemoveMessageActionsBuilderWithContext(pn, ctx)
}

// SetToken Stores a single token in the Token Management System for use in API calls.
func (pn *PubNub) SetToken(token string) {
	pn.tokenManager.StoreToken(token)
}

// ResetTokenManager resets the token manager.
func (pn *PubNub) ResetTokenManager() {
	pn.tokenManager.CleanUp()
}

// Unsubscribe When subscribed to a single channel, this function causes the client to issue a leave from the channel and close any open socket to the PubNub Network. For multiplexed channels, the specified channel(s) will be removed and the socket remains open until there are no more channels remaining in the list.
func (pn *PubNub) Unsubscribe() *unsubscribeBuilder {
	return newUnsubscribeBuilder(pn)
}

// AddListener lets you add a new listener.
func (pn *PubNub) AddListener(listener *Listener) {
	pn.subscriptionManager.AddListener(listener)
}

// RemoveListener lets you remove new listener.
func (pn *PubNub) RemoveListener(listener *Listener) {
	pn.subscriptionManager.RemoveListener(listener)
}

// GetListeners gets all the existing isteners.
func (pn *PubNub) GetListeners() map[*Listener]bool {
	return pn.subscriptionManager.GetListeners()
}

// Leave unsubscribes from a channel.
func (pn *PubNub) Leave() *leaveBuilder {
	return newLeaveBuilder(pn)
}

// LeaveWithContext unsubscribes from a channel.
func (pn *PubNub) LeaveWithContext(ctx Context) *leaveBuilder {
	return newLeaveBuilderWithContext(pn, ctx)
}

// Presence lets you subscribe to a presence channel.
func (pn *PubNub) Presence() *presenceBuilder {
	return newPresenceBuilder(pn)
}

// PresenceWithContext lets you subscribe to a presence channel.
func (pn *PubNub) PresenceWithContext(ctx Context) *presenceBuilder {
	return newPresenceBuilderWithContext(pn, ctx)
}

// Heartbeat You can send presence heartbeat notifications without subscribing to a channel. These notifications are sent periodically and indicate whether a client is connected or not.
func (pn *PubNub) Heartbeat() *heartbeatBuilder {
	return newHeartbeatBuilder(pn)
}

// HeartbeatWithContext You can send presence heartbeat notifications without subscribing to a channel. These notifications are sent periodically and indicate whether a client is connected or not.
func (pn *PubNub) HeartbeatWithContext(ctx Context) *heartbeatBuilder {
	return newHeartbeatBuilderWithContext(pn, ctx)
}

// SetClient Set a client for transactional requests (Non Subscribe).
func (pn *PubNub) SetClient(c *http.Client) {
	pn.Lock()
	pn.client = c
	pn.Unlock()
}

// GetClient Get a client for transactional requests (Non Subscribe).
func (pn *PubNub) GetClient() *http.Client {
	pn.Lock()
	defer pn.Unlock()

	if pn.client == nil {
		if pn.Config.UseHTTP2 {
			pn.client = NewHTTP2Client(pn.Config.ConnectTimeout,
				pn.Config.SubscribeRequestTimeout)
		} else {
			pn.client = NewHTTP1Client(pn.Config.ConnectTimeout,
				pn.Config.NonSubscribeRequestTimeout,
				pn.Config.MaxIdleConnsPerHost)
		}
	}

	return pn.client
}

// SetSubscribeClient Set a client for transactional requests.
func (pn *PubNub) SetSubscribeClient(client *http.Client) {
	pn.Lock()
	pn.subscribeClient = client
	pn.Unlock()
}

// GetSubscribeClient Get a client for transactional requests.
func (pn *PubNub) GetSubscribeClient() *http.Client {
	pn.Lock()
	defer pn.Unlock()
	if pn.subscribeClient == nil {

		if pn.Config.UseHTTP2 {
			pn.subscribeClient = NewHTTP2Client(pn.Config.ConnectTimeout,
				pn.Config.SubscribeRequestTimeout)
		} else {
			pn.subscribeClient = NewHTTP1Client(pn.Config.ConnectTimeout,
				pn.Config.SubscribeRequestTimeout, pn.Config.MaxIdleConnsPerHost)
		}

	}

	return pn.subscribeClient
}

// GetSubscribedChannels gets a list of all subscribed channels.
func (pn *PubNub) GetSubscribedChannels() []string {
	return pn.subscriptionManager.getSubscribedChannels()
}

// GetSubscribedGroups gets a list of all subscribed channel groups.
func (pn *PubNub) GetSubscribedGroups() []string {
	return pn.subscriptionManager.getSubscribedGroups()
}

// UnsubscribeAll Unsubscribe from all channels and all channel groups.
func (pn *PubNub) UnsubscribeAll() {
	pn.subscriptionManager.unsubscribeAll()
}

// ListPushProvisions Request for all channels on which push notification has been enabled using specified pushToken.
func (pn *PubNub) ListPushProvisions() *listPushProvisionsRequestBuilder {
	return newListPushProvisionsRequestBuilder(pn)
}

// ListPushProvisionsWithContext Request for all channels on which push notification has been enabled using specified pushToken.
func (pn *PubNub) ListPushProvisionsWithContext(ctx Context) *listPushProvisionsRequestBuilder {
	return newListPushProvisionsRequestBuilderWithContext(pn, ctx)
}

// AddPushNotificationsOnChannels Enable push notifications on provided set of channels.
func (pn *PubNub) AddPushNotificationsOnChannels() *addPushNotificationsOnChannelsBuilder {
	return newAddPushNotificationsOnChannelsBuilder(pn)
}

// AddPushNotificationsOnChannelsWithContext Enable push notifications on provided set of channels.
func (pn *PubNub) AddPushNotificationsOnChannelsWithContext(ctx Context) *addPushNotificationsOnChannelsBuilder {
	return newAddPushNotificationsOnChannelsBuilderWithContext(pn, ctx)
}

// RemovePushNotificationsFromChannels Disable push notifications on provided set of channels.
func (pn *PubNub) RemovePushNotificationsFromChannels() *removeChannelsFromPushBuilder {
	return newRemoveChannelsFromPushBuilder(pn)
}

// RemovePushNotificationsFromChannelsWithContext Disable push notifications on provided set of channels.
func (pn *PubNub) RemovePushNotificationsFromChannelsWithContext(ctx Context) *removeChannelsFromPushBuilder {
	return newRemoveChannelsFromPushBuilderWithContext(pn, ctx)
}

// RemoveAllPushNotifications Disable push notifications from all channels registered with the specified pushToken.
func (pn *PubNub) RemoveAllPushNotifications() *removeAllPushChannelsForDeviceBuilder {
	return newRemoveAllPushChannelsForDeviceBuilder(pn)
}

// RemoveAllPushNotificationsWithContext Disable push notifications from all channels registered with the specified pushToken.
func (pn *PubNub) RemoveAllPushNotificationsWithContext(ctx Context) *removeAllPushChannelsForDeviceBuilder {
	return newRemoveAllPushChannelsForDeviceBuilderWithContext(pn, ctx)
}

// AddChannelToChannelGroup This function adds a channel to a channel group.
func (pn *PubNub) AddChannelToChannelGroup() *addChannelToChannelGroupBuilder {
	return newAddChannelToChannelGroupBuilder(pn)
}

// AddChannelToChannelGroupWithContext This function adds a channel to a channel group.
func (pn *PubNub) AddChannelToChannelGroupWithContext(ctx Context) *addChannelToChannelGroupBuilder {
	return newAddChannelToChannelGroupBuilderWithContext(pn, ctx)
}

// RemoveChannelFromChannelGroup This function removes the channels from the channel group.
func (pn *PubNub) RemoveChannelFromChannelGroup() *removeChannelFromChannelGroupBuilder {
	return newRemoveChannelFromChannelGroupBuilder(pn)
}

// RemoveChannelFromChannelGroupWithContext This function removes the channels from the channel group.
func (pn *PubNub) RemoveChannelFromChannelGroupWithContext(ctx Context) *removeChannelFromChannelGroupBuilder {
	return newRemoveChannelFromChannelGroupBuilderWithContext(pn, ctx)
}

// DeleteChannelGroup This function removes the channel group.
func (pn *PubNub) DeleteChannelGroup() *deleteChannelGroupBuilder {
	return newDeleteChannelGroupBuilder(pn)
}

// DeleteChannelGroupWithContext This function removes the channel group.
func (pn *PubNub) DeleteChannelGroupWithContext(ctx Context) *deleteChannelGroupBuilder {
	return newDeleteChannelGroupBuilderWithContext(pn, ctx)
}

// ListChannelsInChannelGroup This function lists all the channels of the channel group.
func (pn *PubNub) ListChannelsInChannelGroup() *allChannelGroupBuilder {
	return newAllChannelGroupBuilder(pn)
}

// ListChannelsInChannelGroupWithContext This function lists all the channels of the channel group.
func (pn *PubNub) ListChannelsInChannelGroupWithContext(ctx Context) *allChannelGroupBuilder {
	return newAllChannelGroupBuilderWithContext(pn, ctx)
}

// GetState The state API is used to set/get key/value pairs specific to a subscriber UUID. State information is supplied as a JSON object of key/value pairs.
func (pn *PubNub) GetState() *getStateBuilder {
	return newGetStateBuilder(pn)
}

// GetStateWithContext The state API is used to set/get key/value pairs specific to a subscriber UUID. State information is supplied as a JSON object of key/value pairs.
func (pn *PubNub) GetStateWithContext(ctx Context) *getStateBuilder {
	return newGetStateBuilderWithContext(pn, ctx)
}

// HereNow You can obtain information about the current state of a channel including a list of unique user-ids currently subscribed to the channel and the total occupancy count of the channel by calling the HereNow() function in your application.
func (pn *PubNub) HereNow() *hereNowBuilder {
	return newHereNowBuilder(pn)
}

// HereNowWithContext You can obtain information about the current state of a channel including a list of unique user-ids currently subscribed to the channel and the total occupancy count of the channel by calling the HereNow() function in your application.
func (pn *PubNub) HereNowWithContext(ctx Context) *hereNowBuilder {
	return newHereNowBuilderWithContext(pn, ctx)
}

// WhereNow You can obtain information about the current list of channels to which a UUID is subscribed to by calling the WhereNow() function in your application.
func (pn *PubNub) WhereNow() *whereNowBuilder {
	return newWhereNowBuilder(pn)
}

// WhereNowWithContext You can obtain information about the current list of channels to which a UUID is subscribed to by calling the WhereNow() function in your application.
func (pn *PubNub) WhereNowWithContext(ctx Context) *whereNowBuilder {
	return newWhereNowBuilderWithContext(pn, ctx)
}

// Time This function will return a 17 digit precision Unix epoch.
func (pn *PubNub) Time() *timeBuilder {
	return newTimeBuilder(pn)
}

// TimeWithContext This function will return a 17 digit precision Unix epoch.
func (pn *PubNub) TimeWithContext(ctx Context) *timeBuilder {
	return newTimeBuilderWithContext(pn, ctx)
}

// CreatePushPayload This method creates the push payload for use in the appropriate endpoint calls.
func (pn *PubNub) CreatePushPayload() *publishPushHelperBuilder {
	return newPublishPushHelperBuilder(pn)
}

// CreatePushPayloadWithContext This method creates the push payload for use in the appropriate endpoint calls.
func (pn *PubNub) CreatePushPayloadWithContext(ctx Context) *publishPushHelperBuilder {
	return newPublishPushHelperBuilderWithContext(pn, ctx)
}

// DeleteMessages Removes the messages from the history of a specific channel.
func (pn *PubNub) DeleteMessages() *historyDeleteBuilder {
	return newHistoryDeleteBuilder(pn)
}

// DeleteMessagesWithContext Removes the messages from the history of a specific channel.
func (pn *PubNub) DeleteMessagesWithContext(ctx Context) *historyDeleteBuilder {
	return newHistoryDeleteBuilderWithContext(pn, ctx)
}

// SendFile Clients can use this SDK method to upload a file and publish it to other users in a channel.
func (pn *PubNub) SendFile() *sendFileBuilder {
	return newSendFileBuilder(pn)
}

// SendFileWithContext Clients can use this SDK method to upload a file and publish it to other users in a channel.
func (pn *PubNub) SendFileWithContext(ctx Context) *sendFileBuilder {
	return newSendFileBuilderWithContext(pn, ctx)
}

// ListFiles Provides the ability to fetch all files in a channel.
func (pn *PubNub) ListFiles() *listFilesBuilder {
	return newListFilesBuilder(pn)
}

// ListFilesWithContext Provides the ability to fetch all files in a channel.
func (pn *PubNub) ListFilesWithContext(ctx Context) *listFilesBuilder {
	return newListFilesBuilderWithContext(pn, ctx)
}

// GetFileURL gets the URL of the file.
func (pn *PubNub) GetFileURL() *getFileURLBuilder {
	return newGetFileURLBuilder(pn)
}

// GetFileURLWithContext gets the URL of the file.
func (pn *PubNub) GetFileURLWithContext(ctx Context) *getFileURLBuilder {
	return newGetFileURLBuilderWithContext(pn, ctx)
}

// DownloadFile Provides the ability to fetch an individual file.
func (pn *PubNub) DownloadFile() *downloadFileBuilder {
	return newDownloadFileBuilder(pn)
}

// DownloadFileWithContext Provides the ability to fetch an individual file.
func (pn *PubNub) DownloadFileWithContext(ctx Context) *downloadFileBuilder {
	return newDownloadFileBuilderWithContext(pn, ctx)
}

// DeleteFile Provides the ability to delete an individual file.
func (pn *PubNub) DeleteFile() *deleteFileBuilder {
	return newDeleteFileBuilder(pn)
}

// DeleteFileWithContext Provides the ability to delete an individual file
func (pn *PubNub) DeleteFileWithContext(ctx Context) *deleteFileBuilder {
	return newDeleteFileBuilderWithContext(pn, ctx)
}

// PublishFileMessage Provides the ability to publish the asccociated messages with the uploaded file in case of failure to auto publish. To be used only in the case of failure.
func (pn *PubNub) PublishFileMessage() *publishFileMessageBuilder {
	return newPublishFileMessageBuilder(pn)
}

// PublishFileMessageWithContext Provides the ability to publish the asccociated messages with the uploaded file in case of failure to auto publish. To be used only in the case of failure.
func (pn *PubNub) PublishFileMessageWithContext(ctx Context) *publishFileMessageBuilder {
	return newPublishFileMessageBuilderWithContext(pn, ctx)
}

// Destroy stops all open requests, removes listeners, closes heartbeats, and cleans up.
func (pn *PubNub) Destroy() {
	pn.Config.Log.Println("Calling Destroy")
	pn.UnsubscribeAll()
	pn.cancel()

	if pn.subscriptionManager != nil {
		pn.subscriptionManager.Destroy()
		pn.Config.Log.Println("after subscription manager Destroy")
	}

	pn.Config.Log.Println("calling subscriptionManager Destroy")
	if pn.heartbeatManager != nil {
		pn.heartbeatManager.Destroy()
		pn.Config.Log.Println("after heartbeat manager Destroy")
	}

	pn.Config.Log.Println("After Destroy")
	pn.Config.Log.Println("calling RemoveAllListeners")
	pn.subscriptionManager.RemoveAllListeners()
	pn.Config.Log.Println("after RemoveAllListeners")
	// Check if jobQueue is already closed before attempting to close it
	select {
	case _, ok := <-pn.jobQueue:
		if !ok {
			pn.Config.Log.Println("jobQueue is already closed")
			break
		}
		// If the channel is open, proceed to close it
		close(pn.jobQueue)
		pn.Config.Log.Println("after close jobQueue")
	default:
		// If the channel is closed, no action is needed
		pn.Config.Log.Println("jobQueue is already closed")
	}
	pn.requestWorkers.Close()
	pn.Config.Log.Println("after close requestWorkers")
	pn.tokenManager.CleanUp()
	pn.client.CloseIdleConnections()

}

func (pn *PubNub) getPublishSequence() int {
	pn.publishSequenceMutex.Lock()
	defer pn.publishSequenceMutex.Unlock()

	if pn.nextPublishSequence == MaxSequence {
		pn.nextPublishSequence = 1
	} else {
		pn.nextPublishSequence++
	}

	return pn.nextPublishSequence
}

func GenerateUUID() string {
	return utils.UUID()
}

// NewPubNub instantiates a PubNub instance with default values.
func NewPubNub(pnconf *Config) *PubNub {
	ctx, cancel := contextWithCancel(backgroundContext)

	if pnconf.Log == nil {
		pnconf.Log = log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	pnconf.Log.Println(fmt.Sprintf("PubNub Go v7 SDK: %s\npnconf: %v\n%s\n%s\n%s", Version, pnconf, runtime.Version(), runtime.GOARCH, runtime.GOOS))

	utils.CheckUUID(pnconf.UUID)
	pn := &PubNub{
		Config:              pnconf,
		nextPublishSequence: 0,
		ctx:                 ctx,
		cancel:              cancel,
		previousIvFlag:      pnconf.UseRandomInitializationVector,
		previousCipherKey:   pnconf.CipherKey,
	}

	if pnconf.CipherKey != "" {
		var e error
		pn.Config.CryptoModule, e = crypto.NewLegacyCryptoModule(pnconf.CipherKey, pnconf.UseRandomInitializationVector)
		if e != nil {
			panic(e)
		}
	}
	pn.subscriptionManager = newSubscriptionManager(pn, ctx)
	pn.heartbeatManager = newHeartbeatManager(pn, ctx)
	pn.telemetryManager = newTelemetryManager(pnconf.MaximumLatencyDataAge, ctx)
	pn.jobQueue = make(chan *JobQItem)
	pn.requestWorkers = pn.newNonSubQueueProcessor(pnconf.MaxWorkers, ctx)
	pn.tokenManager = newTokenManager(pn, ctx)

	return pn
}

func (pn *PubNub) newNonSubQueueProcessor(maxWorkers int, ctx Context) *RequestWorkers {
	workers := make(chan chan *JobQItem, maxWorkers)

	pn.Config.Log.Printf("Init RequestWorkers: workers %d", maxWorkers)

	p := &RequestWorkers{
		WorkersChannel: workers,
		MaxWorkers:     maxWorkers,
	}
	p.Start(pn, ctx)
	return p
}

// NewPubNubDemo returns an instance with demo keys
func NewPubNubDemo() *PubNub {
	return NewPubNub(NewDemoConfig())
}
