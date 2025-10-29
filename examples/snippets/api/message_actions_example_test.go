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
	defer pn.RemoveMessageAction().Execute()     // <- Skip this line (cleanup for tests)
	// snippet.show

When copying examples to your own code:
- Use your own publish/subscribe keys instead of the "demo" keys
- Remove any statements that are between snippet.hide and snippet.show (they're only for testing purposes)
*/

// snippet.add_message_action
// Example_addMessageAction demonstrates adding a reaction to a message
func Example_addMessageAction() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Publish a message first to get a timetoken
	publishResp, _, _ := pn.Publish().
		Channel("my-channel").
		Message("Hello, World!").
		Execute()

	messageTimetoken := fmt.Sprintf("%d", publishResp.Timestamp)

	response, status, err := pn.AddMessageAction().
		Channel("my-channel").
		MessageTimetoken(messageTimetoken). // Timetoken of the message to react to
		Action(pubnub.MessageAction{
			ActionType:  "reaction", // Type of action (e.g., "reaction", "receipt")
			ActionValue: "👍",        // Value (e.g., emoji, "read", "delivered")
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Data.ActionTimetoken != "" {
		fmt.Println("Message action added successfully")
	}

	// Output:
	// Message action added successfully
}

// snippet.add_multiple_reactions
// Example_addMultipleReactions demonstrates adding multiple reactions to a message
func Example_addMultipleReactions() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Publish a message first to get a timetoken
	publishResp, _, _ := pn.Publish().
		Channel("my-channel").
		Message("Great news!").
		Execute()

	messageTimetoken := fmt.Sprintf("%d", publishResp.Timestamp)

	// Add multiple reactions to the same message
	reactions := []string{"❤️", "👏", "🎉"}

	for _, emoji := range reactions {
		_, status, err := pn.AddMessageAction().
			Channel("my-channel").
			MessageTimetoken(messageTimetoken).
			Action(pubnub.MessageAction{
				ActionType:  "reaction",
				ActionValue: emoji,
			}).
			Execute()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if status.StatusCode == 200 {
			fmt.Printf("Added reaction: %s\n", emoji)
		}
	}

	// Output:
	// Added reaction: ❤️
	// Added reaction: 👏
	// Added reaction: 🎉
}

// snippet.get_message_actions
// Example_getMessageActions demonstrates retrieving all actions on a channel
func Example_getMessageActions() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Publish a message and add an action to it
	publishResp, _, _ := pn.Publish().
		Channel("get-actions-channel").
		Message("Test message").
		Execute()

	messageTimetoken := fmt.Sprintf("%d", publishResp.Timestamp)

	pn.AddMessageAction().
		Channel("get-actions-channel").
		MessageTimetoken(messageTimetoken).
		Action(pubnub.MessageAction{
			ActionType:  "reaction",
			ActionValue: "😊",
		}).
		Execute()

	// Get all message actions in a channel
	response, status, err := pn.GetMessageActions().
		Channel("get-actions-channel").
		Limit(25). // Limit number of results
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && len(response.Data) > 0 {
		fmt.Println("Message actions retrieved successfully")
		// Print the first message action as an example
		action := response.Data[0]
		fmt.Printf("Example action: %s = %s (by %s)\n",
			action.ActionType, action.ActionValue, action.UUID)
	}

	// Output:
	// Message actions retrieved successfully
	// Example action: reaction = 😊 (by GO_SDK_EXAMPLE_USER)
}

// snippet.get_message_actions_with_timerange
// Example_getMessageActionsWithTimeRange demonstrates retrieving actions within a time range
func Example_getMessageActionsWithTimeRange() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Publish a message and add an action
	publishResp, _, _ := pn.Publish().
		Channel("timerange-channel").
		Message("Test message").
		Execute()

	messageTimetoken := fmt.Sprintf("%d", publishResp.Timestamp)

	// Add a message action
	pn.AddMessageAction().
		Channel("timerange-channel").
		MessageTimetoken(messageTimetoken).
		Action(pubnub.MessageAction{
			ActionType:  "reaction",
			ActionValue: "⭐",
		}).
		Execute()

	// Get message actions within a specific time range
	// Note: Start should be a higher timetoken (more recent) than End (older)
	// Using a very wide range to ensure we capture our action
	response, status, err := pn.GetMessageActions().
		Channel("timerange-channel").
		Start(fmt.Sprintf("%d", publishResp.Timestamp+10000000)). // Start (more recent)
		End(fmt.Sprintf("%d", publishResp.Timestamp-10000000)).   // End (older)
		Limit(10).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Retrieved %d message action(s)\n", len(response.Data))
		// Print the message actions
		for _, action := range response.Data {
			fmt.Printf("Action: %s = %s (UUID: %s)\n",
				action.ActionType, action.ActionValue, action.UUID)
		}
	}

	// Output:
	// Retrieved 1 message action(s)
	// Action: reaction = ⭐ (UUID: GO_SDK_EXAMPLE_USER)
}

// snippet.remove_message_action
// Example_removeMessageAction demonstrates removing a specific action from a message
func Example_removeMessageAction() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Publish a message and add an action
	publishResp, _, _ := pn.Publish().
		Channel("my-channel").
		Message("Message to react to").
		Execute()

	messageTimetoken := fmt.Sprintf("%d", publishResp.Timestamp)

	addResp, _, _ := pn.AddMessageAction().
		Channel("my-channel").
		MessageTimetoken(messageTimetoken).
		Action(pubnub.MessageAction{
			ActionType:  "reaction",
			ActionValue: "👋",
		}).
		Execute()

	actionTimetoken := addResp.Data.ActionTimetoken

	// Remove a specific message action
	_, status, err := pn.RemoveMessageAction().
		Channel("my-channel").
		MessageTimetoken(messageTimetoken). // Timetoken of the original message
		ActionTimetoken(actionTimetoken).   // Timetoken of the action to remove
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Message action removed successfully")
	}

	// Output:
	// Message action removed successfully
}

// snippet.read_receipt_example
// Example_readReceiptTracking demonstrates using message actions for read receipts
func Example_readReceiptTracking() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Publish a message first
	publishResp, _, _ := pn.Publish().
		Channel("chat-channel").
		Message("Important message").
		Execute()

	messageTimetoken := fmt.Sprintf("%d", publishResp.Timestamp)

	// Mark a message as read using message actions
	_, status, err := pn.AddMessageAction().
		Channel("chat-channel").
		MessageTimetoken(messageTimetoken).
		Action(pubnub.MessageAction{
			ActionType:  "receipt", // Use "receipt" type for read tracking
			ActionValue: "read",    // Value indicates the message was read
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Message marked as read")
	}

	// Output:
	// Message marked as read
}

// snippet.end
