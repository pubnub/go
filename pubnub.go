package pubnub

import (
	"net/http"

	"github.com/pubnub/go/pnerr"
)

// Default constants
const (
	Version     = "4.0.0-alpha"
	MaxSequence = 65535
)

// Errors
var (
	ErrMissingPubKey    = pnerr.NewValidationError("pubnub: Missing Publish Key")
	ErrMissingSubKey    = pnerr.NewValidationError("pubnub: Missing Subscribe Key")
	ErrMissingChannel   = pnerr.NewValidationError("pubnub: Missing Channel")
	ErrMissingMessage   = pnerr.NewValidationError("pubnub: Missing Message")
	ErrMissingSecretKey = pnerr.NewValidationError("pubnub: Missing Secret Key")
)

// No server connection will be established when you create a new PubNub object.
// To establish a new connection use Subscribe() function of PubNub type.
type PubNub struct {
	Config              *Config
	publishSequence     chan int
	subscriptionManager *SubscriptionManager
	client              *http.Client
	subscribeClient     *http.Client
}

// TODO: replace result with a pointer
func (pn *PubNub) Publish() *publishBuilder {
	return newPublishBuilder(pn)
}

// TODO: replace result with a pointer
func (pn *PubNub) PublishWithContext(ctx Context) *publishBuilder {
	return newPublishBuilderWithContext(pn, ctx)
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
func (pn *PubNub) Subscribe(operation *SubscribeOperation) {
	pn.subscriptionManager.adaptSubscribe(operation)
}

func (pn *PubNub) Unsubscribe(operation *UnsubscribeOperation) {
	pn.subscriptionManager.adaptUnsubscribe(operation)
}

func (pn *PubNub) AddListener(listener *Listener) {
	pn.subscriptionManager.AddListener(listener)
}

func (pn *PubNub) RemoveListener(listener *Listener) {
	pn.subscriptionManager.RemoveListener(listener)
}

func (pn *PubNub) Leave() *leaveBuilder {
	return newLeaveBuilder(pn)
}

func (pn *PubNub) LeaveWithContext(ctx Context) *leaveBuilder {
	return newLeaveBuilderWithContext(pn, ctx)
}

func (pn *PubNub) Heartbeat() *heartbeatBuilder {
	return newHeartbeatBuilder(pn)
}

func (pn *PubNub) HeartbeatWithContext(ctx Context) *heartbeatBuilder {
	return newHeartbeatBuilderWithContext(pn, ctx)
}

// Set a client for transactional requests
func (pn *PubNub) SetClient(c *http.Client) {
	pn.client = c
}

// Set a client for transactional requests
func (pn *PubNub) GetClient() *http.Client {
	if pn.client == nil {
		pn.client = NewHttpClient(pn.Config.ConnectTimeout,
			pn.Config.NonSubscribeRequestTimeout)
	}

	return pn.client
}

// Set a client for transactional requests
func (pn *PubNub) GetSubscribeClient() *http.Client {
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

func NewPubNub(pnconf *Config) *PubNub {
	publishSequence := make(chan int)

	go runPublishSequenceManager(MaxSequence, publishSequence)

	pn := &PubNub{
		Config:          pnconf,
		publishSequence: publishSequence,
	}

	pn.subscriptionManager = newSubscriptionManager(pn)

	return pn
}

func NewPubNubDemo() *PubNub {
	return &PubNub{
		Config: NewDemoConfig(),
	}
}

func runPublishSequenceManager(maxSequence int, ch chan int) {
	for i := 1; ; i++ {
		if i == maxSequence {
			i = 1
		}

		ch <- i
	}
}
