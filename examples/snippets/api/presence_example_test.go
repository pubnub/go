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

// snippet.here_now
// Example_hereNow demonstrates getting presence information for a channel
func Example_hereNow() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get presence information for a channel
	response, status, err := pn.HereNow().
		Channels([]string{"my-channel"}).
		IncludeUUIDs(true).  // Include list of UUIDs
		IncludeState(false). // Don't include state information
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("HereNow request successful")
		// Print occupancy for each channel
		for _, channelData := range response.Channels {
			fmt.Printf("Channel: %s, Occupancy: %d\n", channelData.ChannelName, channelData.Occupancy)
		}
	}

	// Output:
	// HereNow request successful
	// Channel: my-channel, Occupancy: 0
}

// snippet.here_now_with_state
// Example_hereNowWithState demonstrates getting presence with state information
func Example_hereNowWithState() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get presence information including user state
	_, status, err := pn.HereNow().
		Channels([]string{"my-channel"}).
		IncludeUUIDs(true). // Include list of UUIDs
		IncludeState(true). // Include state information
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("HereNow with state successful")
	}

	// Output:
	// HereNow with state successful
}

// snippet.here_now_occupancy_only
// Example_hereNowOccupancyOnly demonstrates getting only occupancy count
func Example_hereNowOccupancyOnly() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get only occupancy count (no UUIDs or state)
	response, status, err := pn.HereNow().
		Channels([]string{"my-channel-100"}).
		IncludeUUIDs(false). // Don't include UUIDs
		IncludeState(false). // Don't include state
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Occupancy count retrieved")
		for _, channelData := range response.Channels {
			fmt.Printf("Occupancy: %d\n", channelData.Occupancy)
		}
	}

	// Output:
	// Occupancy count retrieved
	// Occupancy: 0
}

// snippet.here_now_multiple_channels
// Example_hereNowMultipleChannels demonstrates getting presence for multiple channels
func Example_hereNowMultipleChannels() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get presence for multiple channels at once
	response, status, err := pn.HereNow().
		Channels([]string{"channel-10", "channel-20", "channel-30"}).
		IncludeUUIDs(true).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Multi-channel HereNow successful")
		fmt.Printf("Total occupancy: %d\n", response.TotalOccupancy)
	}

	// Output:
	// Multi-channel HereNow successful
	// Total occupancy: 0
}

// snippet.here_now_channel_group
// Example_hereNowChannelGroup demonstrates getting presence for channel groups
func Example_hereNowChannelGroup() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get presence for channel groups
	response, status, err := pn.HereNow().
		ChannelGroups([]string{"my-channel-group"}).
		IncludeUUIDs(true).
		IncludeState(false).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Channel group HereNow successful")
		fmt.Printf("Total occupancy: %d\n", response.TotalOccupancy)
	}

	// Output:
	// Channel group HereNow successful
	// Total occupancy: 0
}

// snippet.where_now
// Example_whereNow demonstrates getting channels a UUID is subscribed to
func Example_whereNow() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get list of channels a specific UUID is subscribed to
	response, status, err := pn.WhereNow().
		UUID("some-user-uuid").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("WhereNow request successful")
		fmt.Printf("Channels: %v\n", response.Channels)
	}

	// Output:
	// WhereNow request successful
	// Channels: []
}

// snippet.where_now_current_user
// Example_whereNowCurrentUser demonstrates getting channels for current user
func Example_whereNowCurrentUser() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get channels for current user (omit UUID parameter)
	_, status, err := pn.WhereNow().
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("WhereNow for current user successful")
	}

	// Output:
	// WhereNow for current user successful
}

// snippet.set_state
// Example_setState demonstrates setting user state on channels
func Example_setState() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Set custom state for user on channel
	state := map[string]interface{}{
		"is_typing": true,
		"mood":      "happy",
		"status":    "online",
	}

	_, status, err := pn.SetState().
		Channels([]string{"my-channel"}).
		State(state).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("State set successfully")
	}

	// Output:
	// State set successfully
}

// snippet.set_state_multiple_channels
// Example_setStateMultipleChannels demonstrates setting state on multiple channels
func Example_setStateMultipleChannels() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Set state on multiple channels at once
	state := map[string]interface{}{
		"location": "New York",
		"active":   true,
	}

	_, status, err := pn.SetState().
		Channels([]string{"channel-1", "channel-2", "channel-3"}).
		State(state).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("State set on multiple channels")
	}

	// Output:
	// State set on multiple channels
}

// snippet.set_state_channel_group
// Example_setStateChannelGroup demonstrates setting state on channel groups
func Example_setStateChannelGroup() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Set state for channel groups
	state := map[string]interface{}{
		"role": "moderator",
		"age":  25,
	}

	pn.SetState().
		ChannelGroups([]string{"my-channel-group"}).
		State(state).
		Execute()

	fmt.Println("SetState called for channel group")

	// Output:
	// SetState called for channel group
}

// snippet.get_state
// Example_getState demonstrates retrieving user state from channels
func Example_getState() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get state for specific user on channels
	_, status, err := pn.GetState().
		Channels([]string{"my-channel"}).
		UUID("some-user-uuid").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("State retrieved successfully")
	}

	// Output:
	// State retrieved successfully
}

// snippet.get_state_multiple_channels
// Example_getStateMultipleChannels demonstrates getting state from multiple channels
func Example_getStateMultipleChannels() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get state from multiple channels
	_, status, err := pn.GetState().
		Channels([]string{"channel-1", "channel-2"}).
		UUID("some-user-uuid").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("State retrieved from multiple channels")
	}

	// Output:
	// State retrieved from multiple channels
}

// snippet.get_state_channel_group
// Example_getStateChannelGroup demonstrates getting state from channel groups
func Example_getStateChannelGroup() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get state from channel groups
	pn.GetState().
		ChannelGroups([]string{"my-channel-group"}).
		UUID("some-user-uuid").
		Execute()

	fmt.Println("GetState called for channel group")

	// Output:
	// GetState called for channel group
}

// snippet.heartbeat_join
// Example_heartbeatJoin demonstrates sending heartbeat to join presence
func Example_heartbeatJoin() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Start heartbeating on channels without subscribing
	pn.Presence().
		Connected(true). // true = join presence
		Channels([]string{"my-channel"}).
		Execute()

	fmt.Println("Heartbeat started")

	// Output:
	// Heartbeat started
}

// snippet.heartbeat_leave
// Example_heartbeatLeave demonstrates stopping heartbeat to leave presence
func Example_heartbeatLeave() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Stop heartbeating on channels (leave presence)
	pn.Presence().
		Connected(false). // false = leave presence
		Channels([]string{"my-channel"}).
		Execute()

	fmt.Println("Heartbeat stopped")

	// Output:
	// Heartbeat stopped
}

// snippet.heartbeat_channel_group
// Example_heartbeatChannelGroup demonstrates heartbeat for channel groups
func Example_heartbeatChannelGroup() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Start heartbeating on channel groups
	pn.Presence().
		Connected(true).
		ChannelGroups([]string{"my-channel-group"}).
		Execute()

	fmt.Println("Heartbeat started for channel group")

	// Output:
	// Heartbeat started for channel group
}
