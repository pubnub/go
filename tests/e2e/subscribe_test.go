package e2e

import (
	"fmt"
	"sync"
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

/////////////////////////////
/////////////////////////////
// Structure
// - Channel Subscription
// - Groups Subscription
// - Misc
/////////////////////////////
/////////////////////////////

/////////////////////////////
// Channel Subscription
/////////////////////////////

func TestSubscribeUnsubscribe(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-u-ch")

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/subscribe/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/%s/0", ch),
		Query:              "heartbeat=300",
		ResponseBody:       `{"t":{"t":"15079041051785708","r":12},"m":[]}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid", "tt"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(configCopy())
	pn.SetSubscribeClient(interceptor.GetClient())

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
					return
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{ch},
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		Channels: []string{ch},
	})

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
}

func TestSubscribePublishUnsubscribe(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)
	donePublish := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-pu-ch")

	pn := pubnub.NewPubNub(configCopy())

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				assert.Equal(message.Message, "hey")
				donePublish <- true
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{ch},
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Publish().Channel(ch).Message("hey").Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		Channels: []string{ch},
	})

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
}

// Also tests:
// - test operations like publish/unsubscribe invoked inside another goroutine
// - test unsubscribe all
func TestSubscribePublishPartialUnsubscribe(t *testing.T) {
	assert := assert.New(t)
	doneUnsubscribe := make(chan bool)
	errChan := make(chan string)
	var once sync.Once

	ch1 := randomized("sub-partialu-ch1")
	ch2 := randomized("sub-partialu-ch2")
	heyPub := heyIterator(3)
	heySub := heyIterator(3)

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.Uuid = randomized("sub-partialu-uuid")

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					once.Do(func() {
						pn.Publish().Channel(ch1).Message(<-heyPub).Execute()
					})
					continue
				}
				if len(status.AffectedChannels) == 1 && status.Operation == pubnub.PNUnsubscribeOperation {
					assert.Equal(status.AffectedChannels[0], ch2)
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				if message.Message == <-heySub {
					pn.Unsubscribe(&pubnub.UnsubscribeOperation{
						Channels: []string{ch2},
					})
				} else {
					errChan <- fmt.Sprintf("Unexpected message: %s",
						message.Message)
				}
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{ch1, ch2},
	})

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.RemoveListener(listener)
	pn.UnsubscribeAll()

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
}

func TestJoinLeaveChannel(t *testing.T) {
	assert := assert.New(t)

	// await both connected event on emitter and join presence event received
	var wg sync.WaitGroup
	wg.Add(2)

	donePresenceConnect := make(chan bool)
	doneJoin := make(chan bool)
	doneLeave := make(chan bool)
	errChan := make(chan string)
	ch := randomized("ch")

	configEmitter := configCopy()
	configPresenceListener := configCopy()

	configEmitter.Uuid = randomized("sub-lj-emitter")
	configPresenceListener.Uuid = randomized("sub-lj-listener")

	pn := pubnub.NewPubNub(configEmitter)
	pnPresenceListener := pubnub.NewPubNub(configPresenceListener)

	listenerEmitter := pubnub.NewListener()
	listenerPresenceListener := pubnub.NewListener()

	// emitter
	go func() {
		for {
			select {
			case status := <-listenerEmitter.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					wg.Done()
					return
				}
			case <-listenerEmitter.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listenerEmitter.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	// listener
	go func() {
		for {
			select {
			case status := <-listenerPresenceListener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					donePresenceConnect <- true
				}
			case message := <-listenerPresenceListener.Message:
				errChan <- fmt.Sprintf("Unexpected message: %s",
					message.Message)
			case presence := <-listenerPresenceListener.Presence:
				// ignore join event of presence listener
				if presence.Uuid == configPresenceListener.Uuid {
					continue
				}

				assert.Equal(ch, presence.Channel)

				if presence.Event == "leave" {
					assert.Equal(configEmitter.Uuid, presence.Uuid)
					doneLeave <- true
					return
				} else {
					assert.Equal("join", presence.Event)
					assert.Equal(configEmitter.Uuid, presence.Uuid)
					wg.Done()
				}
			}
		}
	}()

	pn.AddListener(listenerEmitter)
	pnPresenceListener.AddListener(listenerPresenceListener)

	pnPresenceListener.Subscribe(&pubnub.SubscribeOperation{
		Channels:        []string{ch},
		PresenceEnabled: true,
	})

	select {
	case <-donePresenceConnect:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{ch},
	})

	go func() {
		wg.Wait()
		doneJoin <- true
	}()

	select {
	case <-doneJoin:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		Channels: []string{ch},
	})

	select {
	case <-doneLeave:
	case err := <-errChan:
		assert.Fail(err)
		return
	}
}

/////////////////////////////
// Channel Group Subscription
/////////////////////////////

func TestSubscribeUnsubscribeGroup(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneUnsubscribe := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-sug-ch")
	cg := randomized("sub-sug-cg")

	pn := pubnub.NewPubNub(configCopy())

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
					return
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	pn.AddListener(listener)

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{ch}).
		Group(cg).
		Execute()

	assert.Nil(err)

	// await for adding channels
	time.Sleep(3 * time.Second)

	pn.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups: []string{cg},
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		ChannelGroups: []string{cg},
	})

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))

	_, _, err = pn.RemoveChannelFromChannelGroup().
		Channels([]string{ch}).
		Group(cg).
		Execute()
}

func TestSubscribePublishUnsubscribeAllGroup(t *testing.T) {
	assert := assert.New(t)
	pn := pubnub.NewPubNub(configCopy())
	listener := pubnub.NewListener()
	doneSubscribe := make(chan bool)
	donePublish := make(chan bool)
	doneUnsubscribe := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-spuag-ch")
	cg1 := randomized("sub-spuag-cg1")
	cg2 := randomized("sub-spuag-cg2")

	pn.AddListener(listener)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNDisconnectedCategory:
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				donePublish <- true
				assert.Equal("hey", message.Message)
				assert.Equal(ch, message.Channel)
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			}
		}
	}()

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{ch}).
		Group(cg1).
		Execute()

	assert.Nil(err)

	// await for adding channel
	time.Sleep(2 * time.Second)

	pn.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups: []string{cg1, cg2},
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Publish().Channel(ch).Message("hey").Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		ChannelGroups: []string{cg2},
	})

	assert.Equal(len(pn.GetSubscribedGroups()), 1)

	pn.UnsubscribeAll()

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	assert.Equal(len(pn.GetSubscribedGroups()), 0)

	_, _, err = pn.RemoveChannelFromChannelGroup().
		Channels([]string{ch}).
		Group(cg1).
		Execute()

	assert.Nil(err)
}

func TestSubscribeJoinLeaveGroup(t *testing.T) {
	assert := assert.New(t)

	// await both connected event on emitter and join presence event received
	var wg sync.WaitGroup
	wg.Add(2)

	donePresenceConnect := make(chan bool)
	doneJoinEvent := make(chan bool)
	doneLeaveEvent := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-jlg-ch")
	cg := randomized("sub-jlg-cg")

	configEmitter := configCopy()
	configPresenceListener := configCopy()

	configEmitter.Uuid = randomized("emitter")
	configPresenceListener.Uuid = randomized("listener")

	pn := pubnub.NewPubNub(configEmitter)
	pnPresenceListener := pubnub.NewPubNub(configPresenceListener)

	listenerEmitter := pubnub.NewListener()
	listenerPresenceListener := pubnub.NewListener()

	// emitter
	go func() {
		for {
			select {
			case status := <-listenerEmitter.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					wg.Done()
					return
				}
			case <-listenerEmitter.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listenerEmitter.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	// listener
	go func() {
		for {
			select {
			case status := <-listenerPresenceListener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					donePresenceConnect <- true
				}
			case <-listenerPresenceListener.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case presence := <-listenerPresenceListener.Presence:
				// ignore join event of presence listener
				if presence.Uuid == configPresenceListener.Uuid {
					continue
				}

				assert.Equal(presence.Channel, ch)

				if presence.Event == "leave" {
					assert.Equal(configEmitter.Uuid, presence.Uuid)
					doneLeaveEvent <- true
					return
				} else {
					assert.Equal("join", presence.Event)
					assert.Equal(configEmitter.Uuid, presence.Uuid)
					wg.Done()
				}
			}
		}
	}()

	pn.AddListener(listenerEmitter)
	pnPresenceListener.AddListener(listenerPresenceListener)

	pnPresenceListener.AddChannelToChannelGroup().
		Channels([]string{ch}).
		Group(cg).
		Execute()

	pnPresenceListener.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups:   []string{cg},
		PresenceEnabled: true,
	})

	select {
	case <-donePresenceConnect:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups: []string{cg},
	})

	go func() {
		wg.Wait()
		doneJoinEvent <- true
	}()

	select {
	case <-doneJoinEvent:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		ChannelGroups: []string{cg},
	})

	select {
	case <-doneLeaveEvent:
	case err := <-errChan:
		assert.Fail(err)
	}
}

/////////////////////////////
// Unsubscribe
/////////////////////////////

func TestUnsubscribeAll(t *testing.T) {
	assert := assert.New(t)
	pn := pubnub.NewPubNub(configCopy())
	channels := []string{
		randomized("sub-ua-ch1"),
		randomized("sub-ua-ch2"),
		randomized("sub-ua-ch3")}

	groups := []string{
		randomized("sub-ua-cg1"),
		randomized("sub-ua-cg2"),
		randomized("sub-ua-cg3")}

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels:        channels,
		ChannelGroups:   groups,
		PresenceEnabled: true,
	})

	assert.Equal(len(pn.GetSubscribedChannels()), 3)
	assert.Equal(len(pn.GetSubscribedGroups()), 3)

	pn.UnsubscribeAll()

	assert.Equal(len(pn.GetSubscribedChannels()), 0)
	assert.Equal(len(pn.GetSubscribedGroups()), 0)
}

/////////////////////////////
// Misc
/////////////////////////////

func TestSubscribe403Error(t *testing.T) {
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/subscribe/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/ch/0",
		Query:              "heartbeat=300",
		ResponseBody:       `{"message":"Forbidden","payload":{"channels":["ch1", "ch2"], "channel-groups":[":cg1", ":cg2"]},"error":true,"service":"Access Manager","status":403}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 403,
	})

	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	doneAccessDenied := make(chan bool)
	errChan := make(chan string)

	pn := pubnub.NewPubNub(configCopy())
	pn.SetSubscribeClient(interceptor.GetClient())
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				case pubnub.PNAccessDeniedCategory:
					doneAccessDenied <- true
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels:  []string{"ch"},
		Transport: interceptor.Transport,
	})

	select {
	case <-doneSubscribe:
		assert.Fail("Access denied expected")
	case <-doneAccessDenied:
	case err := <-errChan:
		assert.Fail(err)
	}

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
}

func TestSubscribeParseUserMeta(t *testing.T) {
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/subscribe/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/ch/0",
		Query:              "heartbeat=300",
		ResponseBody:       `{"t":{"t":"14858178301085322","r":7},"m":[{"a":"4","f":512,"i":"02a7b822-220c-49b0-90c4-d9cbecc0fd85","s":1,"p":{"t":"14858178301075219","r":7},"k":"demo-36","c":"chTest","u":"my-data","d":{"City":"Goiania","Name":"Marcelo"}}]}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid", "tt"},
		ResponseStatusCode: 200,
	})

	assert := assert.New(t)
	doneMeta := make(chan bool)
	errChan := make(chan string)

	pn := pubnub.NewPubNub(configCopy())
	pn.SetSubscribeClient(interceptor.GetClient())
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				// ignore status messages
				if status.Error {
					errChan <- fmt.Sprintf("Status Error: %s", status.Category)
				}
			case message := <-listener.Message:
				meta, ok := message.UserMetadata.(string)
				if !ok {
					errChan <- "Invalid message type"
				}

				assert.Equal(meta, "my-data")
				doneMeta <- true
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch"},
	})

	select {
	case <-doneMeta:
	case err := <-errChan:
		assert.Fail(err)
	}
}

func TestSubscribeWithCustomTimetoken(t *testing.T) {
	ch := "ch"
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/subscribe/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/ch/0",
		ResponseBody:       `{"t":{"t":"15069659902324693","r":12},"m":[]}`,
		Query:              "heartbeat=300",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/subscribe/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/ch/0",
		ResponseBody:       `{"t":{"t":"14607577960932487","r":1},"m":[{"a":"4","f":0,"i":"Client-g5d4g","p":{"t":"14607577960925503","r":1},"k":"sub-c-4cec9f8e-01fa-11e6-8180-0619f8945a4f","c":"ch","d":{"text":"Enter Message Here"},"b":"ch"}]}`,
		Query:              "heartbeat=300&tt=1337",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
		Hang:               true,
	})

	assert := assert.New(t)
	doneConnected := make(chan bool)
	errChan := make(chan string)

	pn := pubnub.NewPubNub(configCopy())
	pn.SetSubscribeClient(interceptor.GetClient())
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				if status.Category == pubnub.PNConnectedCategory {
					doneConnected <- true
				} else {
					errChan <- fmt.Sprintf("Got status while awaiting for a message: %s",
						status.Category)
					return
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a message"
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a message"
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels:  []string{ch},
		Timetoken: int64(1337),
	})

	select {
	case <-doneConnected:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.UnsubscribeAll()
}

func TestSubscribeWithFilter(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	donePublish := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-wf-ch")

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.FilterExpression = "language!=spanish"
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				}
			case message := <-listener.Message:
				if message.Message == "Hello!" {
					donePublish <- true
				}
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{ch},
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	pnPublish := pubnub.NewPubNub(configCopy())

	meta := make(map[string]string)
	meta["language"] = "spanish"

	pnPublish.Publish().
		Channel("ch").
		Meta(meta).
		Message("Hola!").
		Execute()

	anotherMeta := make(map[string]string)
	anotherMeta["language"] = "english"

	pnPublish.Publish().
		Channel(ch).
		Meta(anotherMeta).
		Message("Hello!").
		Execute()

	<-donePublish
}

func TestSubscribePublishUnsubscribeWithEncrypt(t *testing.T) {
	assert := assert.New(t)
	doneConnect := make(chan bool)
	donePublish := make(chan bool)
	errChan := make(chan string)
	ch := randomized("sub-puwe-ch")

	config := configCopy()
	config.CipherKey = "my-key"
	pn := pubnub.NewPubNub(config)
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneConnect <- true
				}
			case message := <-listener.Message:
				assert.Equal("hey", message.Message)
				donePublish <- true
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{ch},
	})

	select {
	case <-doneConnect:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Publish().
		UsePost(true).
		Channel(ch).
		Message("hey").
		Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
	}
}

func TestSubscribeSuperCall(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	errChan := make(chan string)
	config := pamConfigCopy()
	// Not allowed characters:
	// .,:*
	validCharacters := "-_~?#[]@!$&'()+;=`|"
	config.Uuid = validCharacters
	config.AuthKey = validCharacters

	pn := pubnub.NewPubNub(config)
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels:      []string{validCharacters + "channel"},
		ChannelGroups: []string{validCharacters + "cg"},
		Timetoken:     int64(1337),
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}
}

func TestReconnectionExhaustion(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	errChan := make(chan string)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/subscribe/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/ch/0",
		ResponseBody:       "",
		Query:              "heartbeat=300",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 400,
	})

	config.MaximumReconnectionRetries = 1
	config.PNReconnectionPolicy = pubnub.PNLinearPolicy
	pn := pubnub.NewPubNub(config)
	pn.SetSubscribeClient(interceptor.GetClient())
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNReconnectedCategory:
					doneSubscribe <- true
				}
			case <-listener.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch"},
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}
}
