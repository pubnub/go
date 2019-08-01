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
	Version = "4.2.5"
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
	// StrChannelsTimetoken shows Missing Channels Timetoken message
	StrChannelsTimetoken = "Missing Channels Timetoken"
	// StrChannelsTimetokenLength shows Length of Channels Timetoken message
	StrChannelsTimetokenLength = "Length of Channels Timetoken and Channels do not match"
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

func (pn *PubNub) DeleteMessages() *historyDeleteBuilder {
	return newHistoryDeleteBuilder(pn)
}

func (pn *PubNub) DeleteMessagesWithContext(ctx Context) *historyDeleteBuilder {
	return newHistoryDeleteBuilderWithContext(pn, ctx)
}

func (pn *PubNub) Destroy() {
	pn.requestWorkers.Close()

	close(pn.jobQueue)
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

	return pn
}

func (pn *PubNub) newNonSubQueueProcessor(maxWorkers int, ctx Context) *RequestWorkers {
	workers := make(chan chan *JobQItem, maxWorkers)

	pn.Config.Log.Printf("Init RequestWorkers: workers %d", maxWorkers)

	p := &RequestWorkers{
		Workers:    workers,
		MaxWorkers: maxWorkers,
	}
	p.Start(pn, ctx)
	return p
}

func NewPubNubDemo() *PubNub {
	return NewPubNub(NewDemoConfig())
}
