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

// snippet.fetch_history_basic
// Example_fetchHistoryBasic demonstrates fetching message history from a single channel
func Example_fetchHistoryBasic() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Fetch message history from a channel
	response, status, err := pn.Fetch().
		Channels([]string{"my-go-channel"}).
		Count(10). // Maximum number of messages to retrieve
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Retrieved messages from history:\n")
		// Print the messages
		// Log the fetched messages
		for channel, messages := range response.Messages {
			fmt.Printf("Channel: %s\n", channel)
			for _, message := range messages {
				fmt.Printf("Message: %v, Timetoken: %v\n", message.Message, message.Timetoken)
			}
		}
	}
}

// snippet.fetch_with_metadata
// Example_fetchWithMetadata demonstrates fetching history with message metadata
func Example_fetchWithMetadata() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Fetch history including metadata
	response, status, err := pn.Fetch().
		Channels([]string{"my-channel"}).
		Count(10).
		IncludeMeta(true). // Include metadata for each message
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		// Get messages for the channel
		if messages, ok := response.Messages["my-channel"]; ok {
			fmt.Printf("Retrieved %d messages\n", len(messages))
			// Access message metadata
			for _, msg := range messages {
				fmt.Printf("Message: %v, Metadata: %v\n", msg.Message, msg.Meta)
			}
		}
	}
}

// snippet.fetch_with_time_range
// Example_fetchWithTimeRange demonstrates fetching history within a specific time range
func Example_fetchWithTimeRange() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get a timetoken before publishing
	timeResp, _, _ := pn.Time().Execute()
	startTime := timeResp.Timetoken

	// Publish test messages
	pn.Publish().Channel("my-channel").Message("Message 1").Execute()
	pn.Publish().Channel("my-channel").Message("Message 2").Execute()
	pn.Publish().Channel("my-channel").Message("Message 3").Execute()

	// Get a timetoken after publishing
	timeResp2, _, _ := pn.Time().Execute()
	endTime := timeResp2.Timetoken
	time.Sleep(2 * time.Second) // Wait for messages to be stored

	// Fetch history within a specific time range
	response, status, err := pn.Fetch().
		Channels([]string{"my-channel"}).
		Start(int64(startTime)). // Start timetoken
		End(int64(endTime)).     // End timetoken
		Count(100).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		// Get messages for the channel
		if messages, ok := response.Messages["my-channel"]; ok {
			fmt.Printf("Retrieved %d messages within time range\n", len(messages))
		}
	}

	// Output:
	// Retrieved 3 messages within time range
}

// snippet.fetch_multiple_channels
// Example_fetchMultipleChannels demonstrates fetching history from multiple channels at once
func Example_fetchMultipleChannels() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Fetch history from multiple channels simultaneously
	response, status, err := pn.Fetch().
		Channels([]string{"channel-1", "channel-2", "channel-3"}). // Multiple channels
		Count(25).                                                 // Max 25 messages per channel when using multiple channels
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Retrieved messages from history:\n")
		// Print the messages
		// Log the fetched messages
		for channel, messages := range response.Messages {
			fmt.Printf("Channel: %s\n", channel)
			for _, message := range messages {
				fmt.Printf("Message: %v, Timetoken: %v\n", message.Message, message.Timetoken)
			}
		}
	}
}

// snippet.message_counts
// Example_messageCounts demonstrates getting message counts for channels
func Example_messageCounts() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get count of messages published since a specific timetoken
	response, status, err := pn.MessageCounts().
		Channels([]string{"my-channel"}).
		ChannelsTimetoken([]int64{17300000000000000}). // Count messages since this timetoken
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Message counts retrieved successfully")
		for channel, count := range response.Channels {
			fmt.Printf("Channel: %s, Count: %d\n", channel, count)
		}
	}
}

// snippet.message_counts_multiple
// Example_messageCountsMultiple demonstrates getting message counts for multiple channels
func Example_messageCountsMultiple() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get count of messages published since a specific timetoken
	response, status, err := pn.MessageCounts().
		Channels([]string{"my-channel", "my-channel-1"}).
		ChannelsTimetoken([]int64{17300000000000000}). // Count messages since this timetoken
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Message counts retrieved successfully")
		for channel, count := range response.Channels {
			fmt.Printf("Channel: %s, Count: %d\n", channel, count)
		}
	}
}

// snippet.delete_messages
// Example_deleteMessages demonstrates deleting messages from channel history
func Example_deleteMessages() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Delete all messages from a channel's history
	_, status, err := pn.DeleteMessages().
		Channel("my-channel").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Messages deleted successfully")
	}

	// Output:
	// Messages deleted successfully
}

// snippet.delete_messages_time_range
// Example_deleteMessagesTimeRange demonstrates deleting messages within a specific time range
func Example_deleteMessagesTimeRange() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Publish test messages

	timeBeforeFirstMessage := time.Now().UnixMilli()

	pn.Publish().Channel("my-channel").Message("Message 1").Execute()
	pn.Publish().Channel("my-channel").Message("Message 2").Execute()
	pn.Publish().Channel("my-channel").Message("Message 3").Execute()

	timeAftereLastMessage := time.Now().UnixMilli()

	// Delete messages within a specific time range
	_, status, err := pn.DeleteMessages().
		Channel("my-channel").
		Start(int64(timeAftereLastMessage)). // Start timetoken is the most recent time from where to delete messages
		End(int64(timeBeforeFirstMessage)).  // End timetoken is the oldest time from where to delete messages
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Messages in time range deleted successfully")
	}

	// Output:
	// Messages in time range deleted successfully
}

// snippet.history_deprecated
// Example_historyDeprecated demonstrates using the deprecated History API
// Note: For new code, use Fetch() instead as it provides more features
func Example_historyDeprecated() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Use the deprecated History API (single channel only)
	response, status, err := pn.History().
		Channel("history-channel").
		Count(10).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Retrieved %d messages using History API\n", len(response.Messages))
		for _, msg := range response.Messages {
			fmt.Printf("Message: %v\n", msg.Message)
		}
	}
}

// snippet.end
