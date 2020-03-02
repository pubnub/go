package pubnub

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sync"
)

// Default constants
const (
	// Version :the version of the SDK
	Version = "4.6.5"
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
}

//
func (pn *PubNub) Publish() *publishBuilder {
	return newPublishBuilder(pn)
}

func (pn *PubNub) PublishWithContext(ctx Context) *publishBuilder {
	return newPublishBuilderWithContext(pn, ctx)
}

func (pn *PubNub) Fire() *fireBuilder {
	return newFireBuilder(pn)
}

func (pn *PubNub) FireWithContext(ctx Context) *fireBuilder {
	return newFireBuilderWithContext(pn, ctx)
}

func (pn *PubNub) Subscribe() *subscribeBuilder {
	return newSubscribeBuilder(pn)
}

func (pn *PubNub) History() *historyBuilder {
	return newHistoryBuilder(pn)
}

func (pn *PubNub) HistoryWithContext(ctx Context) *historyBuilder {
	return newHistoryBuilderWithContext(pn, ctx)
}

func (pn *PubNub) Fetch() *fetchBuilder {
	return newFetchBuilder(pn)
}

func (pn *PubNub) FetchWithContext(ctx Context) *fetchBuilder {
	return newFetchBuilderWithContext(pn, ctx)
}

func (pn *PubNub) MessageCounts() *messageCountsBuilder {
	return newMessageCountsBuilder(pn)
}

func (pn *PubNub) MessageCountsWithContext(ctx Context) *messageCountsBuilder {
	return newMessageCountsBuilderWithContext(pn, ctx)
}

func (pn *PubNub) CreateUser() *createUserBuilder {
	return newCreateUserBuilder(pn)
}

func (pn *PubNub) CreateUserWithContext(ctx Context) *createUserBuilder {
	return newCreateUserBuilderWithContext(pn, ctx)
}

func (pn *PubNub) GetUsers() *getUsersBuilder {
	return newGetUsersBuilder(pn)
}

func (pn *PubNub) GetUsersWithContext(ctx Context) *getUsersBuilder {
	return newGetUsersBuilderWithContext(pn, ctx)
}

func (pn *PubNub) GetUser() *getUserBuilder {
	return newGetUserBuilder(pn)
}

func (pn *PubNub) GetUserWithContext(ctx Context) *getUserBuilder {
	return newGetUserBuilderWithContext(pn, ctx)
}

func (pn *PubNub) UpdateUser() *updateUserBuilder {
	return newUpdateUserBuilder(pn)
}

func (pn *PubNub) UpdateUserWithContext(ctx Context) *updateUserBuilder {
	return newUpdateUserBuilderWithContext(pn, ctx)
}

func (pn *PubNub) DeleteUser() *deleteUserBuilder {
	return newDeleteUserBuilder(pn)
}

func (pn *PubNub) DeleteUserWithContext(ctx Context) *deleteUserBuilder {
	return newDeleteUserBuilderWithContext(pn, ctx)
}

func (pn *PubNub) CreateSpace() *createSpaceBuilder {
	return newCreateSpaceBuilder(pn)
}

func (pn *PubNub) CreateSpaceWithContext(ctx Context) *createSpaceBuilder {
	return newCreateSpaceBuilderWithContext(pn, ctx)
}

func (pn *PubNub) GetSpaces() *getSpacesBuilder {
	return newGetSpacesBuilder(pn)
}

func (pn *PubNub) GetSpacesWithContext(ctx Context) *getSpacesBuilder {
	return newGetSpacesBuilderWithContext(pn, ctx)
}

func (pn *PubNub) GetSpace() *getSpaceBuilder {
	return newGetSpaceBuilder(pn)
}

func (pn *PubNub) GetSpaceWithContext(ctx Context) *getSpaceBuilder {
	return newGetSpaceBuilderWithContext(pn, ctx)
}

func (pn *PubNub) UpdateSpace() *updateSpaceBuilder {
	return newUpdateSpaceBuilder(pn)
}

func (pn *PubNub) UpdateSpaceWithContext(ctx Context) *updateSpaceBuilder {
	return newUpdateSpaceBuilderWithContext(pn, ctx)
}

func (pn *PubNub) DeleteSpace() *deleteSpaceBuilder {
	return newDeleteSpaceBuilder(pn)
}

func (pn *PubNub) DeleteSpaceWithContext(ctx Context) *deleteSpaceBuilder {
	return newDeleteSpaceBuilderWithContext(pn, ctx)
}

func (pn *PubNub) GetMemberships() *getMembershipsBuilder {
	return newGetMembershipsBuilder(pn)
}

func (pn *PubNub) GetMembershipsWithContext(ctx Context) *getMembershipsBuilder {
	return newGetMembershipsBuilderWithContext(pn, ctx)
}

func (pn *PubNub) GetMembers() *getMembersBuilder {
	return newGetMembersBuilder(pn)
}

func (pn *PubNub) GetMembersWithContext(ctx Context) *getMembersBuilder {
	return newGetMembersBuilderWithContext(pn, ctx)
}

func (pn *PubNub) ManageMembers() *manageMembersBuilder {
	return newManageMembersBuilder(pn)
}

func (pn *PubNub) ManageMembersWithContext(ctx Context) *manageMembersBuilder {
	return newManageMembersBuilderWithContext(pn, ctx)
}

func (pn *PubNub) ManageMemberships() *manageMembershipsBuilder {
	return newManageMembershipsBuilder(pn)
}

func (pn *PubNub) ManageMembershipsWithContext(ctx Context) *manageMembershipsBuilder {
	return newManageMembershipsBuilderWithContext(pn, ctx)
}

func (pn *PubNub) Signal() *signalBuilder {
	return newSignalBuilder(pn)
}

func (pn *PubNub) SignalWithContext(ctx Context) *signalBuilder {
	return newSignalBuilderWithContext(pn, ctx)
}

func (pn *PubNub) SetState() *setStateBuilder {
	return newSetStateBuilder(pn)
}

func (pn *PubNub) SetStateWithContext(ctx Context) *setStateBuilder {
	return newSetStateBuilderWithContext(pn, ctx)
}

func (pn *PubNub) Grant() *grantBuilder {
	return newGrantBuilder(pn)
}

func (pn *PubNub) GrantWithContext(ctx Context) *grantBuilder {
	return newGrantBuilderWithContext(pn, ctx)
}

func (pn *PubNub) GrantToken() *grantTokenBuilder {
	return newGrantTokenBuilder(pn)
}

func (pn *PubNub) GrantTokenWithContext(ctx Context) *grantTokenBuilder {
	return newGrantTokenBuilderWithContext(pn, ctx)
}

func (pn *PubNub) AddMessageAction() *addMessageActionsBuilder {
	return newAddMessageActionsBuilder(pn)
}

func (pn *PubNub) AddMessageActionWithContext(ctx Context) *addMessageActionsBuilder {
	return newAddMessageActionsBuilderWithContext(pn, ctx)
}

func (pn *PubNub) GetMessageActions() *getMessageActionsBuilder {
	return newGetMessageActionsBuilder(pn)
}

func (pn *PubNub) GetMessageActionsWithContext(ctx Context) *getMessageActionsBuilder {
	return newGetMessageActionsBuilderWithContext(pn, ctx)
}

func (pn *PubNub) RemoveMessageAction() *removeMessageActionsBuilder {
	return newRemoveMessageActionsBuilder(pn)
}

func (pn *PubNub) RemoveMessageActionWithContext(ctx Context) *removeMessageActionsBuilder {
	return newRemoveMessageActionsBuilderWithContext(pn, ctx)
}

func (pn *PubNub) SetToken(token string) {
	pn.tokenManager.StoreToken(token)
}

func (pn *PubNub) SetTokens(tokens []string) {
	pn.tokenManager.StoreTokens(tokens)
}

func (pn *PubNub) GetTokens() GrantResourcesWithPermissions {
	return pn.tokenManager.GetAllTokens()
}

func (pn *PubNub) GetTokensByResource(resourceType PNResourceType) GrantResourcesWithPermissions {
	return pn.tokenManager.GetTokensByResource(resourceType)
}

func (pn *PubNub) GetToken(resourceId string, resourceType PNResourceType) string {
	return pn.tokenManager.GetToken(resourceId, resourceType)
}

func (pn *PubNub) ResetTokenManager() {
	pn.tokenManager.CleanUp()
}

func (pn *PubNub) Unsubscribe() *unsubscribeBuilder {
	return newUnsubscribeBuilder(pn)
}

func (pn *PubNub) AddListener(listener *Listener) {
	pn.subscriptionManager.AddListener(listener)
}

func (pn *PubNub) RemoveListener(listener *Listener) {
	pn.subscriptionManager.RemoveListener(listener)
}

func (pn *PubNub) GetListeners() map[*Listener]bool {
	return pn.subscriptionManager.GetListeners()
}

func (pn *PubNub) Leave() *leaveBuilder {
	return newLeaveBuilder(pn)
}

func (pn *PubNub) LeaveWithContext(ctx Context) *leaveBuilder {
	return newLeaveBuilderWithContext(pn, ctx)
}

func (pn *PubNub) Presence() *presenceBuilder {
	return newPresenceBuilder(pn)
}

func (pn *PubNub) PresenceWithContext(ctx Context) *presenceBuilder {
	return newPresenceBuilderWithContext(pn, ctx)
}

func (pn *PubNub) heartbeat() *heartbeatBuilder {
	return newHeartbeatBuilder(pn)
}

func (pn *PubNub) heartbeatWithContext(ctx Context) *heartbeatBuilder {
	return newHeartbeatBuilderWithContext(pn, ctx)
}

// SetClient Set a client for transactional requests
func (pn *PubNub) SetClient(c *http.Client) {
	pn.Lock()
	pn.client = c
	pn.Unlock()
}

// GetClient Get a client for transactional requests
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

func (pn *PubNub) SetSubscribeClient(client *http.Client) {
	pn.Lock()
	pn.subscribeClient = client
	pn.Unlock()
}

// GetSubscribeClient Get a client for transactional requests
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

func (pn *PubNub) GetSubscribedChannels() []string {
	return pn.subscriptionManager.getSubscribedChannels()
}

func (pn *PubNub) GetSubscribedGroups() []string {
	return pn.subscriptionManager.getSubscribedGroups()
}

func (pn *PubNub) UnsubscribeAll() {
	pn.subscriptionManager.unsubscribeAll()
}

func (pn *PubNub) ListPushProvisions() *listPushProvisionsRequestBuilder {
	return newListPushProvisionsRequestBuilder(pn)
}

func (pn *PubNub) ListPushProvisionsWithContext(
	ctx Context) *listPushProvisionsRequestBuilder {
	return newListPushProvisionsRequestBuilderWithContext(pn, ctx)
}

func (pn *PubNub) AddPushNotificationsOnChannels() *addPushNotificationsOnChannelsBuilder {
	return newAddPushNotificationsOnChannelsBuilder(pn)
}

func (pn *PubNub) AddPushNotificationsOnChannelsWithContext(
	ctx Context) *addPushNotificationsOnChannelsBuilder {
	return newAddPushNotificationsOnChannelsBuilderWithContext(pn, ctx)
}

func (pn *PubNub) RemovePushNotificationsFromChannels() *removeChannelsFromPushBuilder {
	return newRemoveChannelsFromPushBuilder(pn)
}

func (pn *PubNub) RemovePushNotificationsFromChannelsWithContext(
	ctx Context) *removeChannelsFromPushBuilder {
	return newRemoveChannelsFromPushBuilderWithContext(pn, ctx)
}

func (pn *PubNub) RemoveAllPushNotifications() *removeAllPushChannelsForDeviceBuilder {
	return newRemoveAllPushChannelsForDeviceBuilder(pn)
}

func (pn *PubNub) RemoveAllPushNotificationsWithContext(
	ctx Context) *removeAllPushChannelsForDeviceBuilder {
	return newRemoveAllPushChannelsForDeviceBuilderWithContext(pn, ctx)
}

func (pn *PubNub) AddChannelToChannelGroup() *addChannelToChannelGroupBuilder {
	return newAddChannelToChannelGroupBuilder(pn)
}

func (pn *PubNub) AddChannelToChannelGroupWithContext(
	ctx Context) *addChannelToChannelGroupBuilder {
	return newAddChannelToChannelGroupBuilderWithContext(pn, ctx)
}

func (pn *PubNub) RemoveChannelFromChannelGroup() *removeChannelFromChannelGroupBuilder {
	return newRemoveChannelFromChannelGroupBuilder(pn)
}

func (pn *PubNub) RemoveChannelFromChannelGroupWithContext(
	ctx Context) *removeChannelFromChannelGroupBuilder {
	return newRemoveChannelFromChannelGroupBuilderWithContext(pn, ctx)
}

func (pn *PubNub) DeleteChannelGroup() *deleteChannelGroupBuilder {
	return newDeleteChannelGroupBuilder(pn)
}

func (pn *PubNub) DeleteChannelGroupWithContext(
	ctx Context) *deleteChannelGroupBuilder {
	return newDeleteChannelGroupBuilderWithContext(pn, ctx)
}

func (pn *PubNub) ListChannelsInChannelGroup() *allChannelGroupBuilder {
	return newAllChannelGroupBuilder(pn)
}

func (pn *PubNub) ListChannelsInChannelGroupWithContext(
	ctx Context) *allChannelGroupBuilder {
	return newAllChannelGroupBuilderWithContext(pn, ctx)
}

func (pn *PubNub) GetState() *getStateBuilder {
	return newGetStateBuilder(pn)
}

func (pn *PubNub) GetStateWithContext(ctx Context) *getStateBuilder {
	return newGetStateBuilderWithContext(pn, ctx)
}

func (pn *PubNub) HereNow() *hereNowBuilder {
	return newHereNowBuilder(pn)
}

func (pn *PubNub) HereNowWithContext(ctx Context) *hereNowBuilder {
	return newHereNowBuilderWithContext(pn, ctx)
}

func (pn *PubNub) WhereNow() *whereNowBuilder {
	return newWhereNowBuilder(pn)
}

func (pn *PubNub) WhereNowWithContext(ctx Context) *whereNowBuilder {
	return newWhereNowBuilderWithContext(pn, ctx)
}

func (pn *PubNub) Time() *timeBuilder {
	return newTimeBuilder(pn)
}

func (pn *PubNub) TimeWithContext(ctx Context) *timeBuilder {
	return newTimeBuilderWithContext(pn, ctx)
}

func (pn *PubNub) CreatePushPayload() *publishPushHelperBuilder {
	return newPublishPushHelperBuilder(pn)
}

func (pn *PubNub) CreatePushPayloadWithContext(ctx Context) *publishPushHelperBuilder {
	return newPublishPushHelperBuilderWithContext(pn, ctx)
}

func (pn *PubNub) DeleteMessages() *historyDeleteBuilder {
	return newHistoryDeleteBuilder(pn)
}

func (pn *PubNub) DeleteMessagesWithContext(ctx Context) *historyDeleteBuilder {
	return newHistoryDeleteBuilderWithContext(pn, ctx)
}

func (pn *PubNub) Destroy() {
	pn.Config.Log.Println("Calling Destroy")
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
	close(pn.jobQueue)
	pn.Config.Log.Println("after close jobQueue")
	pn.requestWorkers.Close()
	pn.Config.Log.Println("after close requestWorkers")
	pn.tokenManager.CleanUp()

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

func NewPubNub(pnconf *Config) *PubNub {
	ctx, cancel := contextWithCancel(backgroundContext)

	if pnconf.Log == nil {
		pnconf.Log = log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	pnconf.Log.Println(fmt.Sprintf("PubNub Go v4 SDK: %s\npnconf: %v\n%s\n%s\n%s", Version, pnconf, runtime.Version(), runtime.GOARCH, runtime.GOOS))

	pn := &PubNub{
		Config:              pnconf,
		nextPublishSequence: 0,
		ctx:                 ctx,
		cancel:              cancel,
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

func NewPubNubDemo() *PubNub {
	return NewPubNub(NewDemoConfig())
}
