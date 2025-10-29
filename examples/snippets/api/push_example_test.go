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
	// snippet.show

When copying examples to your own code:
- Use your own publish/subscribe keys instead of the "demo" keys
- Remove any statements that are between snippet.hide and snippet.show (they're only for testing purposes)
- Ensure Mobile Push Notifications add-on is enabled in your PubNub Admin Portal
*/

// snippet.add_channels_to_push_fcm
// Example_addChannelsToPushFCM demonstrates registering a device for FCM push notifications
func Example_addChannelsToPushFCM() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Add channels to push notifications for FCM (Firebase Cloud Messaging)
	_, status, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{"notifications-channel", "alerts-channel"}). // Channels to enable for push
		DeviceIDForPush("device-fcm-token").                           // FCM device token
		PushType(pubnub.PNPushTypeFCM).                                // Use FCM push type
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("FCM push notifications enabled for channels")
	}

	// Output:
	// FCM push notifications enabled for channels
}

// snippet.add_channels_to_push_apns2
// Example_addChannelsToPushAPNS2 demonstrates registering a device for APNs2 push notifications
func Example_addChannelsToPushAPNS2() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Add channels to push notifications for APNs2 (Apple Push Notification service)
	_, status, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{"notifications-channel", "alerts-channel"}).                       // Channels to enable for push
		DeviceIDForPush("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"). // APNs device token (64 hex chars)
		PushType(pubnub.PNPushTypeAPNS2).                                                    // Use APNs2 push type
		Topic("com.example.myapp").                                                          // APNs topic (bundle identifier)
		Environment(pubnub.PNPushEnvironmentDevelopment).                                    // APNs environment (development or production)
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("APNs2 push notifications enabled for channels")
	}

	// Output:
	// APNs2 push notifications enabled for channels
}

// snippet.add_channels_to_push_apns2_production
// Example_addChannelsToPushAPNS2Production demonstrates registering a device for production APNs2 push
func Example_addChannelsToPushAPNS2Production() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Add channels to push notifications for APNs2 production environment
	_, status, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{"notifications-channel"}).                                         // Channels to enable for push
		DeviceIDForPush("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"). // APNs device token (64 hex chars)
		PushType(pubnub.PNPushTypeAPNS2).                                                    // Use APNs2 push type
		Topic("com.example.myapp").                                                          // APNs topic (bundle identifier)
		Environment(pubnub.PNPushEnvironmentProduction).                                     // Production environment
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("APNs2 production push notifications enabled")
	}

	// Output:
	// APNs2 production push notifications enabled
}

// snippet.list_push_channels_fcm
// Example_listPushChannelsFCM demonstrates listing channels with push enabled for an FCM device
func Example_listPushChannelsFCM() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Setup: Register some channels first
	pn.AddPushNotificationsOnChannels().
		Channels([]string{"notifications-channel", "alerts-channel"}).
		DeviceIDForPush("device-fcm-token").
		PushType(pubnub.PNPushTypeFCM).
		Execute()

	// List all channels that have push notifications enabled for this FCM device
	response, status, err := pn.ListPushProvisions().
		DeviceIDForPush("device-fcm-token"). // FCM device token
		PushType(pubnub.PNPushTypeFCM).      // Use FCM push type
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		// Display the channels with push enabled
		if len(response.Channels) > 0 {
			fmt.Printf("Push enabled on %d channels\n", len(response.Channels))
		} else {
			fmt.Println("No channels with push enabled")
		}
	}

	// Output:
	// Push enabled on 2 channels
}

// snippet.list_push_channels_apns2
// Example_listPushChannelsAPNS2 demonstrates listing channels with push enabled for an APNs2 device
func Example_listPushChannelsAPNS2() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// List all channels that have push notifications enabled for this APNs2 device
	response, status, err := pn.ListPushProvisions().
		DeviceIDForPush("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"). // APNs device token (64 hex chars)
		PushType(pubnub.PNPushTypeAPNS2).                                                    // Use APNs2 push type
		Topic("com.example.myapp").                                                          // APNs topic (bundle identifier)
		Environment(pubnub.PNPushEnvironmentDevelopment).                                    // APNs environment
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		// Display the channels with push enabled
		if len(response.Channels) > 0 {
			fmt.Printf("Push enabled on %d channels\n", len(response.Channels))
		} else {
			fmt.Println("No channels with push enabled")
		}
	}
}

// snippet.remove_channels_from_push_fcm
// Example_removeChannelsFromPushFCM demonstrates removing specific channels from FCM push notifications
func Example_removeChannelsFromPushFCM() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Remove specific channels from FCM push notifications
	_, status, err := pn.RemovePushNotificationsFromChannels().
		Channels([]string{"alerts-channel"}). // Channels to remove from push
		DeviceIDForPush("device-fcm-token").  // FCM device token
		PushType(pubnub.PNPushTypeFCM).       // Use FCM push type
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Channels removed from FCM push notifications")
	}

	// Output:
	// Channels removed from FCM push notifications
}

// snippet.remove_channels_from_push_apns2
// Example_removeChannelsFromPushAPNS2 demonstrates removing specific channels from APNs2 push notifications
func Example_removeChannelsFromPushAPNS2() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Remove specific channels from APNs2 push notifications
	_, status, err := pn.RemovePushNotificationsFromChannels().
		Channels([]string{"alerts-channel"}).                                                // Channels to remove from push
		DeviceIDForPush("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"). // APNs device token (64 hex chars)
		PushType(pubnub.PNPushTypeAPNS2).                                                    // Use APNs2 push type
		Topic("com.example.myapp").                                                          // APNs topic (bundle identifier)
		Environment(pubnub.PNPushEnvironmentDevelopment).                                    // APNs environment
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Channels removed from APNs2 push notifications")
	}

	// Output:
	// Channels removed from APNs2 push notifications
}

// snippet.remove_all_push_channels_fcm
// Example_removeAllPushChannelsFCM demonstrates removing all channels from FCM push notifications
func Example_removeAllPushChannelsFCM() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Remove all push notification channels for this FCM device
	// This is useful when a user logs out or uninstalls the app
	_, status, err := pn.RemoveAllPushNotifications().
		DeviceIDForPush("device-fcm-token"). // FCM device token
		PushType(pubnub.PNPushTypeFCM).      // Use FCM push type
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("All FCM push notifications removed")
	}

	// Output:
	// All FCM push notifications removed
}

// snippet.remove_all_push_channels_apns2
// Example_removeAllPushChannelsAPNS2 demonstrates removing all channels from APNs2 push notifications
func Example_removeAllPushChannelsAPNS2() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Remove all push notification channels for this APNs2 device
	// This is useful when a user logs out or uninstalls the app
	_, status, err := pn.RemoveAllPushNotifications().
		DeviceIDForPush("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"). // APNs device token (64 hex chars)
		PushType(pubnub.PNPushTypeAPNS2).                                                    // Use APNs2 push type
		Topic("com.example.myapp").                                                          // APNs topic (bundle identifier)
		Environment(pubnub.PNPushEnvironmentDevelopment).                                    // APNs environment
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("All APNs2 push notifications removed")
	}

	// Output:
	// All APNs2 push notifications removed
}

// snippet.end
