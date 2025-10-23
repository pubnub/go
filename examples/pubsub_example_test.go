// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"

	pubnub "github.com/pubnub/go/v7"
)

//snippet.end

// snippet.simple_publish
// Example_simplePublish demonstrates basic message publishing
func Example_simplePublish() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	// Override config with real keys for CI/CD testing. If you copy example from here, just skip this part.
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

// snippet.end

// snippet.publish_with_metadata
// Example_publishWithMetadata demonstrates publishing with custom metadata
func Example_publishWithMetadata() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	// Override config with real keys for CI/CD testing. If you copy example from here, just skip this part.
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

// snippet.end

// snippet.publish_with_ttl
// Example_publishWithTTL demonstrates publishing with message expiration
func Example_publishWithTTL() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	// Override config with real keys for CI/CD testing. If you copy example from here, just skip this part.
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

// snippet.end

// snippet.fire
// Example_fire demonstrates fire operation (publish without storage)
func Example_fire() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	// Override config with real keys for CI/CD testing. If you copy example from here, just skip this part.
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

// snippet.end

// snippet.signal
// Example_signal demonstrates sending signals (lightweight messages)
func Example_signal() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	// Override config with real keys for CI/CD testing. If you copy example from here, just skip this part.
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

// snippet.end

// snippet.subscribe
// Example_subscribe demonstrates basic channel subscription
func Example_subscribe() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	// Override config with real keys for CI/CD testing. If you copy example from here, just skip this part.
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Create listener to handle incoming messages
	listener := pubnub.NewListener()

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
			}
		}
	}()

	// Add listener and subscribe to channel
	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"my-channel"}).
		Execute()

	fmt.Println("Subscribed to channel")

	// Output:
	// Subscribed to channel
}

// snippet.end

// snippet.subscribe_multiple
// Example_subscribeMultiple demonstrates subscribing to multiple channels
func Example_subscribeMultiple() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	// Override config with real keys for CI/CD testing. If you copy example from here, just skip this part.
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()

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
			}
		}
	}()

	pn.AddListener(listener)

	// Subscribe to multiple channels at once
	pn.Subscribe().
		Channels([]string{"channel-1", "channel-2", "channel-3"}).
		Execute()

	fmt.Println("Subscribed to multiple channels")

	// Output:
	// Subscribed to multiple channels
}

// snippet.end

// snippet.subscribe_with_presence
// Example_subscribeWithPresence demonstrates subscription with presence enabled
func Example_subscribeWithPresence() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	// Override config with real keys for CI/CD testing. If you copy example from here, just skip this part.
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()

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

	// Output:
	// Subscribed with presence
}

// snippet.end

// snippet.subscribe_wildcard
// Example_subscribeWildcard demonstrates wildcard channel subscription
func Example_subscribeWildcard() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	// Override config with real keys for CI/CD testing. If you copy example from here, just skip this part.
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()

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
			}
		}
	}()

	pn.AddListener(listener)

	// Subscribe to wildcard pattern (matches room.*, room.1, room.2, etc.)
	pn.Subscribe().
		Channels([]string{"room.*"}).
		Execute()

	fmt.Println("Subscribed to wildcard channels")

	// Output:
	// Subscribed to wildcard channels
}

// snippet.end

// snippet.unsubscribe
// Example_unsubscribe demonstrates unsubscribing from channels
func Example_unsubscribe() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	// Override config with real keys for CI/CD testing. If you copy example from here, just skip this part.
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

// snippet.end

// snippet.unsubscribe_all
// Example_unsubscribeAll demonstrates unsubscribing from all channels
func Example_unsubscribeAll() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	// Override config with real keys for CI/CD testing. If you copy example from here, just skip this part.
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
