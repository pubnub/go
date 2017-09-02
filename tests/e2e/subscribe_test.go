package e2e

import (
	"fmt"
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
				case pubnub.ConnectedCategory:
					doneSubscribe <- true
				case pubnub.DisconnectedCategory:
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
				case pubnub.ConnectedCategory:
					doneSubscribe <- true
				case pubnub.DisconnectedCategory:
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

	pn := pubnub.NewPubNub(configCopy())

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.ConnectedCategory:
					go func() {
						pn.Publish().Channel("ch1").Message("hey").Execute()
					}()
					continue
				}

				if len(status.AffectedChannels) == 1 &&
					status.Operation == pubnub.PNUnsubscribeOperation {
					assert.Equal(status.AffectedChannels[0], "ch2")
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				if message.Message == "hey" {
					pn.Unsubscribe(&pubnub.UnsubscribeOperation{
						Channels: []string{"ch2"},
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
		Channels: []string{"ch1", "ch2"},
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
	donePresenceConnect := make(chan bool)
	doneConnect := make(chan bool)
	done := make(chan bool)
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
				case pubnub.ConnectedCategory:
					doneConnect <- true
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
				case pubnub.ConnectedCategory:
					donePresenceConnect <- true
				}
			case message := <-listenerPresenceListener.Message:
				errChan <- fmt.Sprintf("Unexpected message: %s",
					message.Message)
			case presence := <-listenerPresenceListener.Presence:
				if presence.Uuid == configPresenceListener.Uuid {
					continue
				}

				assert.Equal("ch-join", presence.Channel)

				if presence.Event == "leave" {
					assert.Equal(configEmitter.Uuid, presence.Uuid)
					done <- true
				} else {
					assert.Equal("join", presence.Event)
					assert.Equal(configEmitter.Uuid, presence.Uuid)
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
	}

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch-join"},
	})

	select {
	case <-doneConnect:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		Channels: []string{"ch-join"},
	})

	select {
	case <-done:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.RemoveChannelChannelGroup().
		Channels([]string{"ch-join"}).
		Group("cg").
		Execute()
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
				case pubnub.ConnectedCategory:
					doneSubscribe <- true
				case pubnub.DisconnectedCategory:
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

	pn.AddChannelChannelGroup().
		Channels([]string{"ch"}).
		Group("cg").
		Execute()

	pn.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups: []string{"cg"},
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		ChannelGroups: []string{"cg"},
	})

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}

	assert.Zero(len(pn.GetSubscribedChannels()))
	assert.Zero(len(pn.GetSubscribedGroups()))
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
				case pubnub.ConnectedCategory:
					doneSubscribe <- true
				case pubnub.DisconnectedCategory:
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				donePublish <- true
				assert.Equal("hey", message.Message)
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			}
		}
	}()

	pn.AddChannelChannelGroup().
		Channels([]string{"ch"}).
		Group("cg1").
		Execute()

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

	assert.Equal(len(pn.GetSubscribedGroups()), 0)
}

func TestSubscribeCGPublishUnsubscribe(t *testing.T) {
	assert := assert.New(t)
	doneSubscribe := make(chan bool)
	donePublish := make(chan bool)
	doneUnsubscribe := make(chan bool)
	errChan := make(chan string)

	pn := pubnub.NewPubNub(configCopy())
	listener := pubnub.NewListener()
	pn.AddListener(listener)
	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.ConnectedCategory:
					doneSubscribe <- true
				case pubnub.DisconnectedCategory:
					doneUnsubscribe <- true
				}
			case message := <-listener.Message:
				donePublish <- true
				assert.Equal("hey", message.Message)
			case <-listener.Presence:
				errChan <- "Got presence while awaiting for a status event"
			}
		}
	}()

	pn.AddChannelChannelGroup().
		Channels([]string{"ch"}).
		Group("cg1").
		Execute()

	pn.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups: []string{"cg1"},
	})

	assert.Equal(len(pn.GetSubscribedGroups()), 1)

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
		ChannelGroups: []string{"cg1"},
	})

	select {
	case <-doneUnsubscribe:
	case err := <-errChan:
		assert.Fail(err)
		return
	}

	assert.Equal(len(pn.GetSubscribedGroups()), 0)
}

func TestSubscribeJoinLeaveGroup(t *testing.T) {
	assert := assert.New(t)

	donePresenceConnect := make(chan bool)
	doneEmitterConnect := make(chan bool)
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
				case pubnub.ConnectedCategory:
					doneEmitterConnect <- true
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
				case pubnub.ConnectedCategory:
					donePresenceConnect <- true
				}
			case <-listenerPresenceListener.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case presence := <-listenerPresenceListener.Presence:
				if presence.Uuid == configPresenceListener.Uuid {
					continue
				}

				assert.Equal(presence.Channel, "ch-j")

				if presence.Event == "join" {
					doneJoinEvent <- true
				}

				if presence.Event == "leave" {
					doneLeaveEvent <- true
				}
			}
		}
	}()

	pn.AddListener(listenerEmitter)
	pnPresenceListener.AddListener(listenerPresenceListener)

	pnPresenceListener.AddChannelChannelGroup().
		Channels([]string{"ch-j"}).
		Group("cg").
		Execute()

	pnPresenceListener.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups:   []string{"cg"},
		PresenceEnabled: true,
	})

	select {
	case <-donePresenceConnect:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch-j"},
	})

	select {
	case <-doneJoinEvent:
	case err := <-errChan:
		assert.Fail(err)
	}

	select {
	case <-doneEmitterConnect:
	case err := <-errChan:
		assert.Fail(err)
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		Channels: []string{"ch-j"},
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
		Query:              "",
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
				case pubnub.ConnectedCategory:
					doneSubscribe <- true
				case pubnub.AccessDeniedCategory:
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
}

func TestSubscribeWithMeta(t *testing.T) {
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/subscribe/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/ch/0",
		Query:              "",
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
			case <-listener.Status:
				// ignore status messages
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
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/subscribe/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/ch/0",
		ResponseBody:       `{"t":{"t":"14607577960932487","r":1},"m":[{"a":"4","f":0,"i":"Client-g5d4g","p":{"t":"14607577960925503","r":1},"k":"sub-c-4cec9f8e-01fa-11e6-8180-0619f8945a4f","c":"coolChannel","d":{"text":"Enter Message Here"},"b":"coolChan-bnel"}]}`,
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
				case pubnub.ConnectedCategory:
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
	config.CipherKey = "hey"
	pn := pubnub.NewPubNub(config)
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.ConnectedCategory:
					doneConnect <- true
				}
			case message := <-listener.Message:
				assert.Equal(message.Message, "hey")
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
				case pubnub.ConnectedCategory:
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
		Channels:      []string{SPECIAL_CHANNEL},
		ChannelGroups: []string{SPECIAL_CHANNEL},
		Timetoken:     int64(1337),
	})

	select {
	case <-doneSubscribe:
	case err := <-errChan:
		assert.Fail(err)
	}
}

// TODO
func xTestSubscribeTimeoutError(t *testing.T) {
	doneSubscribe := make(chan bool)
	config := configCopy()
	pn := pubnub.NewPubNub(config)
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.ConnectedCategory:
					doneSubscribe <- true
				}
			case <-listener.Message:
				// ignore
			case <-listener.Presence:
				// ignore
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch"},
	})

	<-doneSubscribe
}
