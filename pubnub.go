package pubnub

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

// Default constants
const (
	Version     = "4.0.0-beta.5"
	MaxSequence = 65535
)

const (
	StrMissingPubKey       = "Missing Publish Key"
	StrMissingSubKey       = "Missing Subscribe Key"
	StrMissingChannel      = "Missing Channel"
	StrMissingChannelGroup = "Missing Channel Group"
	StrMissingMessage      = "Missing Message"
	StrMissingSecretKey    = "Missing Secret Key"
	StrMissingUuid         = "Missing Uuid"
)

// No server connection will be established when you create a new PubNub object.
// To establish a new connection use Subscribe() function of PubNub type.
type PubNub struct {
	sync.RWMutex

	Config               *Config
	nextPublishSequence  int
	publishSequenceMutex sync.RWMutex
	subscriptionManager  *SubscriptionManager
	telemetryManager     *TelemetryManager
	client               *http.Client
	subscribeClient      *http.Client
	ctx                  Context
	cancel               func()
}

func (pn *PubNub) Publish() *publishBuilder {
	return newPublishBuilder(pn)
}

func (pn *PubNub) PublishWithContext(ctx Context) *publishBuilder {
	return newPublishBuilderWithContext(pn, ctx)
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

func (pn *PubNub) heartbeat() *heartbeatBuilder {
	return newHeartbeatBuilder(pn)
}

func (pn *PubNub) heartbeatWithContext(ctx Context) *heartbeatBuilder {
	return newHeartbeatBuilderWithContext(pn, ctx)
}

// Set a client for transactional requests
func (pn *PubNub) SetClient(c *http.Client) {
	pn.Lock()
	pn.client = c
	pn.Unlock()
}

// Set a client for transactional requests
func (pn *PubNub) GetClient() *http.Client {
	pn.Lock()
	defer pn.Unlock()

	if pn.client == nil {
		pn.client = NewHttpClient(pn.Config.ConnectTimeout,
			pn.Config.NonSubscribeRequestTimeout)
	}

	return pn.client
}

func (pn *PubNub) SetSubscribeClient(client *http.Client) {
	pn.Lock()
	pn.subscribeClient = client
	pn.Unlock()
}

// Set a client for transactional requests
func (pn *PubNub) GetSubscribeClient() *http.Client {
	pn.Lock()
	defer pn.Unlock()
	if pn.subscribeClient == nil {

		pn.subscribeClient = NewHttpClient(pn.Config.ConnectTimeout,
			pn.Config.SubscribeRequestTimeout)
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

func (pn *PubNub) AddChannelToChannelGroup() *AddChannelToChannelGroupBuilder {
	return newAddChannelToChannelGroupBuilder(pn)
}

func (pn *PubNub) AddChannelToChannelGroupWithContext(
	ctx Context) *AddChannelToChannelGroupBuilder {
	return newAddChannelToChannelGroupBuilderWithContext(pn, ctx)
}

func (pn *PubNub) RemoveChannelFromChannelGroup() *RemoveChannelFromChannelGroupBuilder {
	return newRemoveChannelFromChannelGroupBuilder(pn)
}

func (pn *PubNub) RemoveChannelFromChannelGroupWithContext(
	ctx Context) *RemoveChannelFromChannelGroupBuilder {
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

func (pn *PubNub) DeleteMessagesWithContext() *historyDeleteBuilder {
	return newHistoryDeleteBuilder(pn)
}

func (pn *PubNub) Destroy() {
	pn.cancel()
	pn.subscriptionManager.RemoveAllListeners()
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

	pn := &PubNub{
		Config:              pnconf,
		nextPublishSequence: 0,
		ctx:                 ctx,
		cancel:              cancel,
	}

	pn.subscriptionManager = newSubscriptionManager(pn, ctx)
	pn.telemetryManager = newTelemetryManager(
		pnconf.MaximumLatencyDataAge, ctx)

	return pn
}

func NewPubNubDemo() *PubNub {
	return NewPubNub(NewDemoConfig())
}
