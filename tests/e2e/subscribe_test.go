package e2e

import (
	"fmt"
	"log"
	"sync"
	"testing"

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

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"blah"},
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		Channels: []string{"blah"},
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
		Channels: []string{"ch"},
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Publish().Channel("ch").Message("hey").Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		Channels: []string{"ch"},
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

	ch1 := randomized("sub-ch1")
	ch2 := randomized("sub-ch2")

	// hey2push := heyIterator(1)
	// hey2pull := heyIterator(1)

	pn := pubnub.NewPubNub(configCopy())

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					go func() {
						// pn.Publish().Channel(ch1).Message(<-hey2push).Execute()
						pn.Publish().Channel(ch1).Message("hey").Execute()
					}()
					continue
				}

				if len(status.AffectedChannels) == 1 &&
					status.Operation == pubnub.PNUnsubscribeOperation {
					assert.Equal(status.AffectedChannels[0], ch2)
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				// if message.Message == <-hey2pull {
				if message.Message == "hey" {
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
			case message := <-listenerPresenceListener.Message:
				errChan <- fmt.Sprintf("Unexpected message: %s",
					message.Message)
			case presence := <-listenerPresenceListener.Presence:
				// ignore join event of presence listener
				if presence.Uuid == configPresenceListener.Uuid {
					continue
				}

				assert.Equal("ch-join", presence.Channel)

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
		Channels:        []string{"ch-join"},
		PresenceEnabled: true,
	})

	select {
	case <-donePresenceConnect:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch-join"},
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
		Channels: []string{"ch-join"},
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

	_, _, err := pn.AddChannelChannelGroup().
		Channels([]string{"subscribe-ch"}).
		Group("subscribe-cg").
		Execute()

	assert.Nil(err)

	pn.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups: []string{"subscribe-cg"},
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		ChannelGroups: []string{"subscribe-cg"},
	})

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))

	_, _, err = pn.RemoveChannelChannelGroup().
		Channels([]string{"subscribe-ch"}).
		Group("subscribe-cg").
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
				assert.Equal("ch", message.Channel)
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			}
		}
	}()

	_, _, err := pn.AddChannelChannelGroup().
		Channels([]string{"ch"}).
		Group("cg1").
		Execute()

	assert.Nil(err)

	pn.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups: []string{"cg1", "cg2"},
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Publish().Channel("ch").Message("hey").Execute()

	select {
	case <-donePublish:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		ChannelGroups: []string{"cg2"},
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

	_, _, err = pn.RemoveChannelChannelGroup().
		Channels([]string{"ch"}).
		Group("cg1").
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

				assert.Equal(presence.Channel, "my-channel")

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

	pnPresenceListener.AddChannelChannelGroup().
		Channels([]string{"my-channel"}).
		Group("my-group").
		Execute()

	pnPresenceListener.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups:   []string{"my-group"},
		PresenceEnabled: true,
	})

	select {
	case <-donePresenceConnect:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups: []string{"my-group"},
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
		ChannelGroups: []string{"my-group"},
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

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels:        []string{"ch1", "ch2", "ch3"},
		ChannelGroups:   []string{"cg1", "cg2", "cg3"},
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
		Query:              "heartbeat=300&hey=123",
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
				log.Println(">>>>>>>>>>>>>>>status", status)
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

	pn.UnsubscribeAll()

	select {
	case <-doneMeta:
	case err := <-errChan:
		assert.Fail(err)
	}
}

func TestSubscribeWithCustomTimetoken(t *testing.T) {
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/subscribe/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/ch/0",
		ResponseBody:       `{"t":{"t":"14607577960932487","r":1},"m":[{"a":"4","f":0,"i":"Client-g5d4g","p":{"t":"14607577960925503","r":1},"k":"sub-c-4cec9f8e-01fa-11e6-8180-0619f8945a4f","c":"coolChannel","d":{"text":"Enter Message Here"},"b":"coolChan-bnel"}]}`,
		Query:              "heartbeat=300&tt=1337",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	errChan := make(chan string)

	pn := pubnub.NewPubNub(configCopy())
	pn.SetSubscribeClient(interceptor.GetClient())
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case <-listener.Status:
				errChan <- "Got status while awaiting for a message"
				return
			case <-listener.Message:
				doneSubscribe <- true
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a message"
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels:  []string{"ch"},
		Timetoken: int64(1337),
	})

	pn.UnsubscribeAll()

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}
}

func TestSubscribeWithFilter(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	donePublish := make(chan bool)
	errChan := make(chan string)

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
		Channels: []string{"ch"},
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
		Channel("ch").
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
		Channels: []string{"ch"},
	})

	select {
	case <-doneConnect:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Publish().
		UsePost(true).
		Channel("ch").
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
	config.Uuid = SPECIAL_CHARACTERS
	config.AuthKey = SPECIAL_CHARACTERS

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
		Channels:      []string{SPECIAL_CHANNEL + "channel"},
		ChannelGroups: []string{SPECIAL_CHANNEL + "cg"},
		Timetoken:     int64(1337),
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}
}
