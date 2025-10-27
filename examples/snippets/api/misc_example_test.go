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
	// snippet.show

When copying examples to your own code:
- Use your own publish/subscribe keys instead of the "demo" keys
- Remove any statements that are between snippet.hide and snippet.show (they're only for testing purposes)
*/

// snippet.time
// Example_getServerTime demonstrates getting the current time from the PubNub server
func Example_getServerTime() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Get the current time from PubNub server
	response, status, err := pn.Time().Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Timetoken > 0 {
		fmt.Println("Server time retrieved successfully")
	}

	// Output:
	// Server time retrieved successfully
}

// snippet.create_push_payload_apns
// Example_createPushPayloadAPNS demonstrates creating a push payload for Apple Push Notification service
func Example_createPushPayloadAPNS() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Create APNS payload with alert, badge, and sound
	apnsPayload := pubnub.PNAPNSData{
		APS: pubnub.PNAPSData{
			Alert: "New message received!",
			Badge: 1,
			Sound: "default",
		},
		Custom: map[string]interface{}{
			"sender": "user-123",
		},
	}

	// Build the push payload
	pushPayload := pn.CreatePushPayload().
		SetAPNSPayload(apnsPayload, nil).
		BuildPayload()

	// Publish the push payload to a channel
	response, status, err := pn.Publish().
		Channel("notification-channel").
		Message(pushPayload).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Timestamp > 0 {
		fmt.Println("APNS push payload published successfully")
	}

	// Output:
	// APNS push payload published successfully
}

// snippet.create_push_payload_fcm
// Example_createPushPayloadFCM demonstrates creating a push payload for Firebase Cloud Messaging
func Example_createPushPayloadFCM() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Create FCM payload with notification data
	fcmPayload := pubnub.PNFCMData{
		Data: pubnub.PNFCMDataFields{
			Summary: "You have a new notification",
			Custom: map[string]interface{}{
				"message_id": "12345",
			},
		},
		Custom: map[string]interface{}{
			"priority": "high",
		},
	}

	// Build the push payload
	pushPayload := pn.CreatePushPayload().
		SetFCMPayload(fcmPayload).
		BuildPayload()

	// Publish the push payload to a channel
	response, status, err := pn.Publish().
		Channel("notification-channel").
		Message(pushPayload).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Timestamp > 0 {
		fmt.Println("FCM push payload published successfully")
	}

	// Output:
	// FCM push payload published successfully
}

// snippet.create_push_payload_combined
// Example_createPushPayloadCombined demonstrates creating a combined push payload for both APNS and FCM
func Example_createPushPayloadCombined() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Create APNS payload
	apnsPayload := pubnub.PNAPNSData{
		APS: pubnub.PNAPSData{
			Alert: "New message!",
			Badge: 1,
			Sound: "default",
		},
	}

	// Create FCM payload
	fcmPayload := pubnub.PNFCMData{
		Data: pubnub.PNFCMDataFields{
			Summary: "New message!",
		},
	}

	// Create common payload for regular PubNub subscribers
	commonPayload := map[string]interface{}{
		"text":      "New message!",
		"timestamp": "2024-01-01T00:00:00Z",
	}

	// Build combined push payload
	pushPayload := pn.CreatePushPayload().
		SetAPNSPayload(apnsPayload, nil).
		SetFCMPayload(fcmPayload).
		SetCommonPayload(commonPayload).
		BuildPayload()

	// Publish the combined push payload
	response, status, err := pn.Publish().
		Channel("notification-channel").
		Message(pushPayload).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Timestamp > 0 {
		fmt.Println("Combined push payload published successfully")
	}

	// Output:
	// Combined push payload published successfully
}

// snippet.end
