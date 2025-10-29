// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"

	pubnub "github.com/pubnub/go/v8"
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

// snippet.add_channels_to_group
// Example_addChannelsToGroup demonstrates adding channels to a channel group
func Example_addChannelsToGroup() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.DeleteChannelGroup().ChannelGroup("my-channel-group").Execute()
	// snippet.show

	// Add channels to a channel group
	// This allows you to subscribe to multiple channels at once using the group name
	_, status, err := pn.AddChannelToChannelGroup().
		Channels([]string{"channel-1", "channel-2", "channel-3"}). // Channels to add
		ChannelGroup("my-channel-group").                          // Channel group name
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Channels added to group successfully")
	}

	// Output:
	// Channels added to group successfully
}

// snippet.list_channels_in_group
// Example_listChannelsInGroup demonstrates listing all channels in a channel group
func Example_listChannelsInGroup() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.DeleteChannelGroup().ChannelGroup("my-channel-group").Execute()
	// snippet.show

	// First, add some channels to the group
	pn.AddChannelToChannelGroup().
		Channels([]string{"channel-1", "channel-2"}).
		ChannelGroup("my-channel-group").
		Execute()

	// List all channels in the channel group
	response, status, err := pn.ListChannelsInChannelGroup().
		ChannelGroup("my-channel-group").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Channels in group: %v\n", response.Channels)
	}

}

// snippet.remove_channels_from_group
// Example_removeChannelsFromGroup demonstrates removing channels from a channel group
func Example_removeChannelsFromGroup() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.DeleteChannelGroup().ChannelGroup("my-channel-group").Execute()
	// snippet.show

	// First, add channels to the group
	pn.AddChannelToChannelGroup().
		Channels([]string{"channel-1", "channel-2", "channel-3"}).
		ChannelGroup("my-channel-group").
		Execute()

	// Remove specific channels from the channel group
	_, status, err := pn.RemoveChannelFromChannelGroup().
		Channels([]string{"channel-1"}). // Channels to remove
		ChannelGroup("my-channel-group").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Channels removed from group successfully")
	}

	// Output:
	// Channels removed from group successfully
}

// snippet.delete_channel_group
// Example_deleteChannelGroup demonstrates deleting a channel group
func Example_deleteChannelGroup() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// First, create a channel group with some channels
	pn.AddChannelToChannelGroup().
		Channels([]string{"channel-1", "channel-2"}).
		ChannelGroup("temp-group").
		Execute()

	// Delete the entire channel group
	// This removes the group and all its channel associations
	response, status, err := pn.DeleteChannelGroup().
		ChannelGroup("temp-group").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Channel group deleted successfully")
	}

	_ = response

	// Output:
	// Channel group deleted successfully
}

// snippet.subscribe_to_channel_group
// Example_subscribeToChannelGroup demonstrates subscribing to a channel group
func Example_subscribeToChannelGroup() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.DeleteChannelGroup().ChannelGroup("my-channel-group").Execute()
	// snippet.show

	// First, add channels to the group
	pn.AddChannelToChannelGroup().
		Channels([]string{"channel-1", "channel-2", "channel-3"}).
		ChannelGroup("my-channel-group").
		Execute()

	// Create listener to handle incoming messages
	listener := pubnub.NewListener()

	// Create a done channel to stop the goroutine when needed
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-listener.Status:
				// Handle connection status changes

			case message := <-listener.Message:
				// Handle messages from any channel in the group
				fmt.Printf("Received message: %v from channel: %s\n",
					message.Message, message.Channel)

			case <-done:
				// Stop the goroutine when done signal is received
				return
			}
		}
	}()

	// Add listener and subscribe to the channel group
	pn.AddListener(listener)

	// Subscribe to channel group - receives messages from all channels in the group
	pn.Subscribe().
		ChannelGroups([]string{"my-channel-group"}).
		Execute()

	fmt.Println("Subscribed to channel group")

	// When done, unsubscribe and stop goroutine
	pn.UnsubscribeAll()
	close(done)

	// Output:
	// Subscribed to channel group
}

// snippet.subscribe_to_multiple_channel_groups
// Example_subscribeToMultipleChannelGroups demonstrates subscribing to multiple channel groups
func Example_subscribeToMultipleChannelGroups() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.DeleteChannelGroup().ChannelGroup("news-group").Execute()
	defer pn.DeleteChannelGroup().ChannelGroup("sports-group").Execute()
	// snippet.show

	// Create multiple channel groups
	pn.AddChannelToChannelGroup().
		Channels([]string{"news-1", "news-2"}).
		ChannelGroup("news-group").
		Execute()

	pn.AddChannelToChannelGroup().
		Channels([]string{"sports-1", "sports-2"}).
		ChannelGroup("sports-group").
		Execute()

	listener := pubnub.NewListener()

	// Create a done channel to stop the goroutine when needed
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-listener.Status:
				// Handle connection status changes

			case message := <-listener.Message:
				fmt.Printf("Received from %s: %v\n",
					message.Channel, message.Message)

			case <-done:
				// Stop the goroutine when done signal is received
				return
			}
		}
	}()

	pn.AddListener(listener)

	// Subscribe to multiple channel groups at once
	pn.Subscribe().
		ChannelGroups([]string{"news-group", "sports-group"}).
		Execute()

	fmt.Println("Subscribed to multiple channel groups")

	// When done, unsubscribe and stop goroutine
	pn.UnsubscribeAll()
	close(done)

	// Output:
	// Subscribed to multiple channel groups
}

// snippet.subscribe_to_channels_and_groups
// Example_subscribeToChannelsAndGroups demonstrates subscribing to both individual channels and channel groups
func Example_subscribeToChannelsAndGroups() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.DeleteChannelGroup().ChannelGroup("my-group").Execute()
	// snippet.show

	// Create a channel group
	pn.AddChannelToChannelGroup().
		Channels([]string{"group-channel-1", "group-channel-2"}).
		ChannelGroup("my-group").
		Execute()

	listener := pubnub.NewListener()

	// Create a done channel to stop the goroutine when needed
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-listener.Status:
				// Handle connection status changes

			case message := <-listener.Message:
				// Messages from both individual channels and channel group channels
				fmt.Printf("Message from %s: %v\n",
					message.Channel, message.Message)

			case <-done:
				// Stop the goroutine when done signal is received
				return
			}
		}
	}()

	pn.AddListener(listener)

	// Subscribe to individual channels AND channel groups simultaneously
	pn.Subscribe().
		Channels([]string{"standalone-channel-1", "standalone-channel-2"}).
		ChannelGroups([]string{"my-group"}).
		Execute()

	fmt.Println("Subscribed to channels and groups")

	// When done, unsubscribe and stop goroutine
	pn.UnsubscribeAll()
	close(done)

	// Output:
	// Subscribed to channels and groups
}

// snippet.unsubscribe_from_channel_group
// Example_unsubscribeFromChannelGroup demonstrates unsubscribing from a channel group
func Example_unsubscribeFromChannelGroup() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.DeleteChannelGroup().ChannelGroup("my-group").Execute()
	// snippet.show

	// First, add channels to the group and subscribe
	pn.AddChannelToChannelGroup().
		Channels([]string{"channel-1", "channel-2"}).
		ChannelGroup("my-group").
		Execute()

	pn.Subscribe().
		ChannelGroups([]string{"my-group"}).
		Execute()

	// Unsubscribe from the channel group
	pn.Unsubscribe().
		ChannelGroups([]string{"my-group"}).
		Execute()

	fmt.Println("Unsubscribed from channel group")

	// Output:
	// Unsubscribed from channel group
}

// snippet.end
