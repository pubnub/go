package e2e

import (
	"fmt"
	"log"
	"sync"
	"testing"

	pubnub "github.com/pubnub/go"
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

// TODO
func xTestJoinLeave(t *testing.T) {
	assert := assert.New(t)

	// await both connected event on emitter and join presence event received
	var wgConnect sync.WaitGroup
	wgConnect.Add(2)

	donePresenceConnect := make(chan bool)
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
					wgConnect.Done()
				}
			case message := <-listenerEmitter.Message:
				errChan <- fmt.Sprintf("Unexpected message: %s",
					message.Message)
			case <-listenerEmitter.Presence:
				// ignore
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
				assert.Equal(presence.Channel, "ch")
				assert.Equal(presence.Event, "join")
				wgConnect.Done()
			}
		}
	}()

	pn.AddListener(listenerEmitter)
	pnPresenceListener.AddListener(listenerPresenceListener)

	pnPresenceListener.Subscribe(&pubnub.SubscribeOperation{
		Channels:        []string{"ch"},
		PresenceEnabled: true,
	})

	<-donePresenceConnect
	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch"},
	})

	done := make(chan bool)

	go func() {
		wgConnect.Wait()
		done <- true
	}()

	select {
	case <-done:
	case err := <-errChan:
		log.Println(err)
		assert.Fail(err)
	}
}

/////////////////////////////
// Channel Group Subscription
/////////////////////////////

// TODO
func aTestSubscribePresenceSingleGroup(t *testing.T) {
	assert := assert.New(t)

	var wgConnect sync.WaitGroup
	wgConnect.Add(2)

	donePresenceConnect := make(chan bool)
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
					log.Println(">>> emitter connected")
					wgConnect.Done()
				}
			case message := <-listenerEmitter.Message:
				errChan <- fmt.Sprintf("Unexpected message: %s",
					message.Message)
			case <-listenerEmitter.Presence:
				// ignore
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
				log.Println(">>> listener join")
				assert.Equal(presence.Channel, "ch")
				assert.Equal(presence.Event, "join")
				wgConnect.Done()
			}
		}
	}()

	pn.AddListener(listenerEmitter)
	pnPresenceListener.AddListener(listenerPresenceListener)

	// Channel has been to channel group from another SDK
	pnPresenceListener.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups:   []string{"cg"},
		PresenceEnabled: true,
	})

	<-donePresenceConnect

	// Channel has been to channel group from another SDK
	pn.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups: []string{"cg"},
	})

	done := make(chan bool)

	go func() {
		wgConnect.Wait()
		done <- true
	}()

	select {
	case <-done:
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

// TODO: add error handlers
// TODO: verify currect status event broadcasted
func TestSubscribe403Error(t *testing.T) {
	doneSubscribe := make(chan bool)

	pn := pubnub.NewPubNub(configCopy())
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
				fmt.Println(message)
			case presence := <-listener.Presence:
				fmt.Println(presence)
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch"},
	})

	<-doneSubscribe
}

// TODO: implement using request stubs
func xTestSubscribeWithMeta(t *testing.T) {
	assert := assert.New(t)

	doneSubscribe := make(chan bool)
	doneMeta := make(chan interface{})

	pn := pubnub.NewPubNub(configCopy())
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
				// Message has been published from pubnub console
				// Example of message:
				// {"message": "hello"}
				doneMeta <- message.UserMetadata
			case presence := <-listener.Presence:
				fmt.Println(presence)
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch"},
	})

	<-doneSubscribe
}

// TODO: use request stub to verify timetoken is set correctly
// TODO: add error handlers
func xTestSubscribeWithCustomTimetoken(t *testing.T) {
	doneSubscribe := make(chan bool)

	pn := pubnub.NewPubNub(configCopy())
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
				fmt.Println(message)
			case presence := <-listener.Presence:
				fmt.Println(presence)
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels:  []string{"ch"},
		Timetoken: int64(1337),
	})

	<-doneSubscribe
}

func TestSubscribeWithFilter(t *testing.T) {
	doneSubscribe := make(chan bool)

	pn := pubnub.NewPubNub(configCopy())
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
				fmt.Println(message)
			case presence := <-listener.Presence:
				fmt.Println(presence)
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels:         []string{"ch"},
		FilterExpression: "foo=bar",
	})

	<-doneSubscribe
}

// TODO: add publish, check for unencrypted message
func TestSubscribePublishUnsubscribeWithEncrypt(t *testing.T) {
	doneSubscribe := make(chan bool)

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.SecretKey = "my-secret"
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
				fmt.Println(message)
			case presence := <-listener.Presence:
				fmt.Println(presence)
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch"},
	})

	<-doneSubscribe
}

func TestSubscribeSuperCall(t *testing.T) {
	doneSubscribe := make(chan bool)
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
			case message := <-listener.Message:
				fmt.Println(message)
			case presence := <-listener.Presence:
				fmt.Println(presence)
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels:      []string{SPECIAL_CHANNEL},
		ChannelGroups: []string{SPECIAL_CHANNEL},
		Timetoken:     int64(1337),
	})

	<-doneSubscribe
}
