// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"

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

// snippet.simple_publish
// Example_simplePublish demonstrates basic message publishing
func Example_simplePublish() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	response, status, err := pn.Publish().
		Channel("my-channel").
		Message("Hello World!").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if response.Timestamp > 0 && status.StatusCode == 200 {
		fmt.Println("Message published successfully")
	}

	// Output:
	// Message published successfully
}

// snippet.publish_with_metadata
// Example_publishWithMetadata demonstrates publishing with custom metadata
func Example_publishWithMetadata() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Create metadata object
	metadata := map[string]interface{}{
		"sender":   "user-123",
		"language": "english",
	}

	// Publish message with metadata
	response, status, err := pn.Publish().
		Channel("my-channel").
		Message("Hello with metadata!").
		Meta(metadata). // Add metadata
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if response.Timestamp > 0 && status.StatusCode == 200 {
		fmt.Println("Message with metadata published successfully")
	}

	// Output:
	// Message with metadata published successfully
}

// snippet.publish_with_ttl
// Example_publishWithTTL demonstrates publishing with message expiration
func Example_publishWithTTL() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Publish message with TTL (Time To Live)
	response, status, err := pn.Publish().
		Channel("my-channel").
		Message("This message expires in 24 hours").
		TTL(24).           // Message expires after 24 hours
		ShouldStore(true). // Ensure message is stored
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if response.Timestamp > 0 && status.StatusCode == 200 {
		fmt.Println("Message with TTL published successfully")
	}

	// Output:
	// Message with TTL published successfully
}

// snippet.publish_array
// Example_publishArray demonstrates publishing an array with metadata
func Example_publishArray() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Publish an array message with metadata and custom message type
	response, status, err := pn.Publish().
		Channel("my-channel").
		Message([]string{"Hello", "there"}).
		Meta([]string{"1a", "2b", "3c"}).
		CustomMessageType("text-message").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if response.Timestamp > 0 && status.StatusCode == 200 {
		fmt.Println("Array message published successfully")
	}

	// Output:
	// Array message published successfully
}

// snippet.fire
// Example_fire demonstrates fire operation (publish without storage)
func Example_fire() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Fire message (not stored in history, lower latency)
	response, status, err := pn.Fire().
		Channel("my-channel").
		Message("Ephemeral message").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if response.Timestamp > 0 && status.StatusCode == 200 {
		fmt.Println("Fire message sent successfully")
	}

	// Output:
	// Fire message sent successfully
}

// snippet.fire_with_metadata
// Example_fireWithMetadata demonstrates fire operation with custom metadata
func Example_fireWithMetadata() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Fire message with metadata (not stored in history)
	response, status, err := pn.Fire().
		Channel("my-channel").
		Message("test").
		Meta(map[string]interface{}{
			"name":      "important-event",
			"timestamp": 1234567890,
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if response.Timestamp > 0 && status.StatusCode == 200 {
		fmt.Println("Fire message with metadata sent successfully")
	}

	// Output:
	// Fire message with metadata sent successfully
}

// snippet.signal
// Example_signal demonstrates sending signals (lightweight messages)
func Example_signal() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Send signal (lightweight, not stored, up to 64 bytes)
	response, status, err := pn.Signal().
		Channel("my-channel").
		Message("typing...").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if response.Timestamp > 0 && status.StatusCode == 200 {
		fmt.Println("Signal sent successfully")
	}

	// Output:
	// Signal sent successfully
}

// snippet.subscribe_basic_with_logging
// Example_subscribeBasicWithLogging demonstrates basic subscription setup
func Example_subscribeBasicWithLogging() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.PublishKey = "demo"
	config.SubscribeKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Subscribe to a channel
	pn.Subscribe().
		Channels([]string{"my-channel"}).
		Execute()

	fmt.Println("Subscribed to my-channel")

	// Cleanup
	pn.UnsubscribeAll()

	// Output:
	// Subscribed to my-channel
}

// snippet.subscribe
// Example_subscribe demonstrates basic channel subscription
func Example_subscribe() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Create listener to handle incoming messages
	listener := pubnub.NewListener()

	// Create a done channel to stop the goroutine when needed
	done := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				// Handle connection status changes
				switch status.Category {
				case pubnub.PNConnectedCategory:
					fmt.Println("Connected to PubNub")
				case pubnub.PNReconnectedCategory:
					fmt.Println("Reconnected to PubNub")
				case pubnub.PNDisconnectedCategory:
					fmt.Println("Disconnected from PubNub")
				}

			case message := <-listener.Message:
				// Handle received messages
				fmt.Printf("Received message: %v on channel: %s\n",
					message.Message, message.Channel)

			case <-done:
				// Stop the goroutine when done signal is received
				return
			}
		}
	}()

	// Add listener and subscribe to channel
	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"my-channel"}).
		Execute()

	fmt.Println("Subscribed to channel")

	// When done, unsubscribe and stop goroutine
	pn.UnsubscribeAll()
	close(done)

	// Output:
	// Subscribed to channel
}

// snippet.subscribe_multiple
// Example_subscribeMultiple demonstrates subscribing to multiple channels
func Example_subscribeMultiple() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()

	// Create a done channel to stop the goroutine when needed
	done := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				if status.Category == pubnub.PNConnectedCategory {
					fmt.Println("Connected to channels")
				}

			case message := <-listener.Message:
				fmt.Printf("Received on %s: %v\n",
					message.Channel, message.Message)

			case <-done:
				// Stop the goroutine when done signal is received
				return
			}
		}
	}()

	pn.AddListener(listener)

	// Subscribe to multiple channels at once
	pn.Subscribe().
		Channels([]string{"channel-1", "channel-2", "channel-3"}).
		Execute()

	fmt.Println("Subscribed to multiple channels")

	// When done, unsubscribe and stop goroutine
	pn.UnsubscribeAll()
	close(done)

	// Output:
	// Subscribed to multiple channels
}

// snippet.subscribe_with_presence
// Example_subscribeWithPresence demonstrates subscription with presence enabled
func Example_subscribeWithPresence() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()

	// Create a done channel to stop the goroutine when needed
	done := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				if status.Category == pubnub.PNConnectedCategory {
					fmt.Println("Connected with presence")
				}

			case message := <-listener.Message:
				// Handle regular messages
				fmt.Printf("Message: %v\n", message.Message)

			case presence := <-listener.Presence:
				// Handle presence events (join, leave, timeout)
				fmt.Printf("Presence event: %s for UUID: %s\n",
					presence.Event, presence.UUID)

			case <-done:
				// Stop the goroutine when done signal is received
				return
			}
		}
	}()

	pn.AddListener(listener)

	// Subscribe with presence events enabled
	pn.Subscribe().
		Channels([]string{"my-channel"}).
		WithPresence(true). // Enable presence events
		Execute()

	fmt.Println("Subscribed with presence")

	// When done, unsubscribe and stop goroutine
	pn.UnsubscribeAll()
	close(done)

	// Output:
	// Subscribed with presence
}

// snippet.subscribe_wildcard
// Example_subscribeWildcard demonstrates wildcard channel subscription
func Example_subscribeWildcard() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()

	// Create a done channel to stop the goroutine when needed
	done := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				if status.Category == pubnub.PNConnectedCategory {
					fmt.Println("Connected to wildcard channels")
				}

			case message := <-listener.Message:
				fmt.Printf("Received on %s: %v\n",
					message.Channel, message.Message)

			case <-done:
				// Stop the goroutine when done signal is received
				return
			}
		}
	}()

	pn.AddListener(listener)

	// Subscribe to wildcard pattern (matches room.*, room.1, room.2, etc.)
	pn.Subscribe().
		Channels([]string{"room.*"}).
		Execute()

	fmt.Println("Subscribed to wildcard channels")

	// When done, unsubscribe and stop goroutine
	pn.UnsubscribeAll()
	close(done)

	// Output:
	// Subscribed to wildcard channels
}

// snippet.subscribe_with_state
// Example_subscribeWithState demonstrates subscribing with state
func Example_subscribeWithState() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()
	done := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					// Set state after connection is established
					response, status, err := pn.SetState().
						Channels([]string{"ch"}).
						State(map[string]interface{}{
							"field_a": "cool",
							"field_b": 21,
						}).Execute()

					if err != nil {
						fmt.Printf("Error: %v\n", err)
					}

					if response != nil && status.StatusCode == 200 {
						fmt.Println("State set successfully")
					}
					done <- true
				}
			case <-listener.Message:
				// Handle messages
			case <-listener.Presence:
				// Handle presence
			case <-done:
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()

	<-done

	pn.UnsubscribeAll()

	// Output:
	// State set successfully
}

// snippet.subscribe_to_channel_group_presence
// subscribeToChannelGroupPresence demonstrates subscribing to presence of channel groups
func subscribeToChannelGroupPresence() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	pn := pubnub.NewPubNub(config)

	// Subscribe to channel groups with presence
	pn.Subscribe().
		ChannelGroups([]string{"cg1", "cg2"}). // subscribe to channel groups
		Timetoken(int64(1337)).                // optional, pass a timetoken
		WithPresence(true).                    // also subscribe to related presence information
		Execute()
}

// snippet.unsubscribe
// Example_unsubscribe demonstrates unsubscribing from channels
func Example_unsubscribe() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// First subscribe to a channel
	pn.Subscribe().
		Channels([]string{"my-channel"}).
		Execute()

	// Unsubscribe from specific channel(s)
	pn.Unsubscribe().
		Channels([]string{"my-channel"}).
		Execute()

	fmt.Println("Unsubscribed from channel")

	// Output:
	// Unsubscribed from channel
}

// snippet.unsubscribe_multiple
// Example_unsubscribeMultiple demonstrates unsubscribing from multiple channels
func Example_unsubscribeMultiple() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// First subscribe to multiple channels
	pn.Subscribe().
		Channels([]string{"my-channel", "my-channel2"}).
		Execute()

	// Unsubscribe from multiple channels
	pn.Unsubscribe().
		Channels([]string{"my-channel", "my-channel2"}).
		Execute()

	fmt.Println("Unsubscribed from multiple channels")

	// Output:
	// Unsubscribed from multiple channels
}

// snippet.unsubscribe_all
// Example_unsubscribeAll demonstrates unsubscribing from all channels
func Example_unsubscribeAll() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Subscribe to multiple channels
	pn.Subscribe().
		Channels([]string{"channel-1", "channel-2", "channel-3"}).
		Execute()

	// Unsubscribe from all channels at once
	pn.UnsubscribeAll()

	fmt.Println("Unsubscribed from all channels")

	// Output:
	// Unsubscribed from all channels
}

// snippet.end
