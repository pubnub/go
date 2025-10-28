// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"
	"time"

	pubnub "github.com/pubnub/go/v7"
)

// snippet.end

/*
IMPORTANT NOTE FOR COPYING EXAMPLES:

Throughout this file, you'll see code between "snippet.hide" and "snippet.show" comments.
These sections are used for CI/CD testing and should be SKIPPED if you're copying examples.

Example of what to skip:
	// snippet.hide
	config = setPubnubExampleConfigData(config)  // <- Skip this line (for testing only)
	defer pn.DeleteChannelGroup().Execute()      // <- Skip this line (cleanup for tests)
	// snippet.show

When copying examples to your own code:
- Use your own publish/subscribe keys instead of the "demo" keys
- Remove any statements that are between snippet.hide and snippet.show (they're only for testing purposes)
*/

// snippet.init
// Example_init demonstrates PubNub initialization
func Example_init() {
	// Create configuration with a unique user ID
	config := pubnub.NewConfigWithUserId(pubnub.UserId("my-user-id"))

	// Set your PubNub keys
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// Initialize PubNub client
	pn := pubnub.NewPubNub(config)

	if pn != nil {
		fmt.Println("PubNub client initialized successfully")
	}
}

// snippet.set_up_listeners
// Example_set_up_listeners demonstrates setting up listeners
func Example_setUpListeners() {
	// snippet.hide
	config := pubnub.NewConfigWithUserId(pubnub.UserId("my-user-id"))
	pn := pubnub.NewPubNub(config)
	// snippet.show

	// Create a listener
	listener := pubnub.NewListener()

	// Create channel to signal connection event
	doneConnect := make(chan bool)

	// Start a goroutine to process events
	go func() {
		for {
			select {
			case status := <-listener.Status:
				// Handle status events
				switch status.Category {
				case pubnub.PNConnectedCategory:
					// Connected to the PubNub network
					fmt.Println("Connected to PubNub!")
					doneConnect <- true
				case pubnub.PNDisconnectedCategory:
					// Disconnected from the PubNub network
					fmt.Println("Disconnected from PubNub!")
				case pubnub.PNReconnectedCategory:
					// Reconnected to the PubNub network
					fmt.Println("Reconnected to PubNub!")
				}

			case message := <-listener.Message:
				// Handle new message
				fmt.Println("Received message:", message.Message)

			case presence := <-listener.Presence:
				// Handle presence
				fmt.Println("Presence event:", presence.Event)
			}
		}
	}()

	// Add the listener to PubNub
	pn.AddListener(listener)
}

// snippet.create_subscription
// Example_createSubscription demonstrates creating a subscription
func Example_createSubscription() {
	// snippet.hide
	config := pubnub.NewConfigWithUserId(pubnub.UserId("my-user-id"))
	pn := pubnub.NewPubNub(config)
	doneConnect := make(chan bool)
	// snippet.show

	// Define the channel you want to subscribe to
	channel := "my-channel"

	// Subscribe to the channel
	pn.Subscribe().
		Channels([]string{channel}).
		Execute()

	<-doneConnect

	fmt.Println("Subscribed to channel:", channel)
}

// snippet.publish_message
// Example_publishMessage demonstrates publishing a message
func Example_publishMessage() {
	// snippet.hide
	config := pubnub.NewConfigWithUserId(pubnub.UserId("my-user-id"))
	pn := pubnub.NewPubNub(config)
	channel := "my-channel"
	// snippet.show

	// Create a message
	message := map[string]interface{}{
		"text":   "Hello, world!",
		"sender": "go-sdk",
	}

	fmt.Println("Publishing message:", message)

	// Publish the message to the channel
	response, _, err := pn.Publish().
		Channel(channel).
		Message(message).
		Execute()

	if err != nil {
		// Handle publish error
		fmt.Println("Publish error:", err)
	} else {
		// Handle successful publish
		fmt.Println("Publish successful! Timetoken:", response.Timestamp)
	}
}

// snippet.complete_example
// Example_completeExample demonstrates a complete PubNub workflow: initialize, subscribe, and publish
func Example_completeExample() {

	// Step 1: Initialize PubNub with configuration
	config := pubnub.NewConfigWithUserId(pubnub.UserId("go-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)
	fmt.Println("PubNub instance initialized")

	// Step 2: Create channel to signal connection event
	doneConnect := make(chan bool)
	listener := pubnub.NewListener()

	// Step 3: Set up listener for events
	go func() {
		for {
			select {
			case status := <-listener.Status:
				// Handle status events
				switch status.Category {
				case pubnub.PNConnectedCategory:
					// Connected to the PubNub network
					fmt.Println("Connected to PubNub!")
					doneConnect <- true
				case pubnub.PNDisconnectedCategory:
					// Disconnected from the PubNub network
					fmt.Println("Disconnected from PubNub!")
				case pubnub.PNReconnectedCategory:
					// Reconnected to the PubNub network
					fmt.Println("Reconnected to PubNub!")
				}

			case message := <-listener.Message:
				// Handle new message
				fmt.Println("Received message:", message.Message)

			case presence := <-listener.Presence:
				// Handle presence
				fmt.Println("Presence event:", presence.Event)
			}
		}
	}()

	// Add the listener to PubNub
	pn.AddListener(listener)

	// Step 4: Define the channel and subscribe
	channel := "my-channel"
	pn.Subscribe().
		Channels([]string{channel}).
		Execute()

	// Wait for connection to establish
	<-doneConnect
	fmt.Println("Subscribed to channel:", channel)

	// Optional: Wait a moment to ensure subscription is fully set up
	time.Sleep(1 * time.Second)

	// Step 5: Create and publish a message
	message := map[string]interface{}{
		"text":   "Hello, world!",
		"sender": "go-sdk",
	}

	fmt.Println("Publishing message:", message)

	// Publish the message
	_, _, err := pn.Publish().
		Channel(channel).
		Message(message).
		Execute()

	if err != nil {
		fmt.Println("Publish error:", err)
	} else {
		fmt.Println("Publish successful!")
	}

	// Wait for the message to be received
	time.Sleep(3 * time.Second)

	// Step 6: Clean up
	pn.UnsubscribeAll()
	fmt.Println("Example completed successfully!")

	// Output:
	// PubNub instance initialized
	// Connected to PubNub!
	// Subscribed to channel: my-channel
	// Publishing message: map[sender:go-sdk text:Hello, world!]
	// Publish successful!
	// Received message: map[sender:go-sdk text:Hello, world!]
	// Example completed successfully!
}

// snippet.end
