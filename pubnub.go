package pubnub

import (
	"context"
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
	ErrMissingPubKey  = pnerr.NewValidationError("pubnub: Missing Publish Key")
	ErrMissingSubKey  = pnerr.NewValidationError("pubnub: Missing Subscribe Key")
	ErrMissingChannel = pnerr.NewValidationError("pubnub: Missing Channel")
	ErrMissingMessage = pnerr.NewValidationError("pubnub: Missing Message")
)

// TODO: pn.UnsubscribeAll() to be deferred

// No server connection will be established when you create a new PubNub object.
// To establish a new connection use Subscribe() function of PubNub type.
type PubNub struct {
	Config          *Config
	publishSequence chan int
	client          *http.Client
}

func (pn *PubNub) Publish(opts *PublishOpts) (PublishResponse, error) {
	res, err := PublishRequest(pn, opts)
	return res, err
}

func (pn *PubNub) PublishWithContext(ctx context.Context,
	opts *PublishOpts) (PublishResponse, error) {

	return PublishRequestWithContext(ctx, pn, opts)
}

// func (pn *PubNub) SimplePublish(string channel, byte []msg) {
// 	return "", []byte(msg)
// }

func (pn *PubNub) SetClient(c *http.Client) {
	pn.client = c
}

func (pn *PubNub) GetClient() *http.Client {
	return pn.client
}

func NewPubNub(pnconf *Config) *PubNub {
	publishSequence := make(chan int)

	go runPublishSequenceManager(MaxSequence, publishSequence)

	return &PubNub{
		Config:          pnconf,
		publishSequence: publishSequence,
	}
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
