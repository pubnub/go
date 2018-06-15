package main

import (
	"fmt"
	"sync"

	pubnub "github.com/pubnub/go"
)

var pn *pubnub.PubNub

func init() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	pn = pubnub.NewPubNub(config)
}

func pubnubCopy() *pubnub.PubNub {
	_pn := new(pubnub.PubNub)
	*_pn = *pn
	return _pn
}

func gettingStarted() {
	listener := pubnub.NewListener()
	doneConnect := make(chan bool)
	donePublish := make(chan bool)

	msg := map[string]interface{}{
		"msg": "hello",
	}
	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNDisconnectedCategory:
					// This event happens when radio / connectivity is lost
				case pubnub.PNConnectedCategory:
					// Connect event. You can do stuff like publish, and know you'll get it.
					// Or just use the connected event to confirm you are subscribed for
					// UI / internal notifications, etc
					doneConnect <- true
				case pubnub.PNReconnectedCategory:
					// Happens as part of our regular operation. This event happens when
					// radio / connectivity is lost, then regained.
				}
			case message := <-listener.Message:
				// Handle new message stored in message.message
				if message.Channel != "" {
					// Message has been received on channel group stored in
					// message.Channel
				} else {
					// Message has been received on channel stored in
					// message.Subscription
				}
				if msg, ok := message.Message.(map[string]interface{}); ok {
					fmt.Println("msg:=====>", msg["msg"])
				}
				/*
				   log the following items with your favorite logger
				       - message.Message
				       - message.Subscription
				       - message.Timetoken
				*/

				donePublish <- true
			case <-listener.Presence:
				// handle presence
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"hello_world"}).
		Execute()

	<-doneConnect

	response, status, err := pn.Publish().
		Channel("hello_world").Message(msg).Execute()

	if err != nil {
		// Request processing failed.
		// Handle message publish error
	}

	fmt.Println(response, status, err)

	<-donePublish
}

func listeners() {
	listener := pubnub.NewListener()
	doneSubscribe := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
					return
				case pubnub.PNDisconnectedCategory:
					//
				case pubnub.PNReconnectedCategory:
					//
				case pubnub.PNAccessDeniedCategory:
					//
				case pubnub.PNUnknownCategory:
					//
				}
			case <-listener.Message:
				//
			case <-listener.Presence:
				//
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()

	<-doneSubscribe
}

func time() {
	res, status, err := pn.Time().Execute()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
	fmt.Println(res)
}

func publish() {
	res, status, err := pn.Publish().
		Channel("ch").
		Message("hey").
		Execute()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
	fmt.Println(res)
}

func hereNow() {
	res, status, err := pn.HereNow().
		Channels([]string{"ch"}).
		IncludeUUIDs(true).
		Execute()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
	fmt.Println(res)
}

func presence() {
	// await both connected event on emitter and join presence event received
	var wg sync.WaitGroup
	wg.Add(2)

	donePresenceConnect := make(chan bool)
	doneJoin := make(chan bool)
	doneLeave := make(chan bool)
	errChan := make(chan string)
	ch := "my-channel"

	configPresenceListener := pubnub.NewConfig()
	configPresenceListener.SubscribeKey = "demo"
	configPresenceListener.PublishKey = "demo"

	pnPresenceListener := pubnub.NewPubNub(configPresenceListener)

	pn.Config.UUID = "my-emitter"
	pnPresenceListener.Config.UUID = "my-listener"

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
				fmt.Println(presence, "\n", configPresenceListener)
				// ignore join event of presence listener
				if presence.UUID == configPresenceListener.UUID {
					continue
				}

				if presence.Event == "leave" {
					doneLeave <- true
					return
				}
				wg.Done()
			}
		}
	}()

	pn.AddListener(listenerEmitter)
	pnPresenceListener.AddListener(listenerPresenceListener)

	pnPresenceListener.Subscribe().
		Channels([]string{ch}).
		WithPresence(true).
		Execute()

	select {
	case <-donePresenceConnect:
	case err := <-errChan:
		panic(err)
	}

	pn.Subscribe().
		Channels([]string{ch}).
		Execute()

	go func() {
		wg.Wait()
		doneJoin <- true
	}()

	select {
	case <-doneJoin:
	case err := <-errChan:
		panic(err)
	}

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	select {
	case <-doneLeave:
	case err := <-errChan:
		panic(err)
	}
}

func history() {
	res, status, err := pn.History().
		Channel("ch").
		Count(2).
		IncludeTimetoken(true).
		Reverse(true).
		Start(int64(1)).
		End(int64(2)).
		Execute()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
	fmt.Println(res)
}

func unsubscribe() {
	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()

	// t.Sleep(3 * t.Second)

	pn.Unsubscribe().
		Channels([]string{"ch"}).
		Execute()
}

func main() {
	// gettingStarted()
	// listeners()
	// time()
	// publish()
	// hereNow()
	presence()
	// history()
	// unsubscribe()
}
