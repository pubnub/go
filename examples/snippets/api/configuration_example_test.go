package pubnub_samples_test

import (
	"fmt"
	"io"
	"log"
	"os"

	pubnub "github.com/pubnub/go/v7"
)

/*
//common includes for most of examples
// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"

	pubnub "github.com/pubnub/go/v7"
)

// snippet.end
*/

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

// snippet.init_basic
// Example_initBasic demonstrates basic PubNub initialization
func Example_initBasic() {
	// Create configuration with a unique user ID
	config := pubnub.NewConfigWithUserId(pubnub.UserId("my-user-id"))

	// Set your PubNub keys
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Initialize PubNub client
	pn := pubnub.NewPubNub(config)

	if pn != nil {
		fmt.Println("PubNub client initialized successfully")
	}

	// Output:
	// PubNub client initialized successfully
}

// snippet.init_read_only
// Example_initReadOnly demonstrates read only PubNub initialization
func Example_initReadOnly() {
	// Create configuration with a unique user ID
	config := pubnub.NewConfigWithUserId(pubnub.UserId("my-user-id"))

	// Set your PubNub keys
	config.SubscribeKey = "demo" // Replace with your subscribe key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Initialize PubNub client
	pn := pubnub.NewPubNub(config)

	if pn != nil {
		fmt.Println("PubNub client initialized successfully")
	}

	// Output:
	// PubNub client initialized successfully
}

// snippet.init_with_uuid
// Example_initWithUUID demonstrates initialization with a generated UUID
func Example_initWithUUID() {
	// Generate a unique UUID for the user
	uuid := pubnub.GenerateUUID()

	// Create configuration with the generated UUID
	config := pubnub.NewConfigWithUserId(pubnub.UserId(uuid))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	if pn != nil && config.GetUserId() != "" {
		fmt.Println("PubNub initialized with generated UUID")
	}

	// Output:
	// PubNub initialized with generated UUID
}

// snippet.set_timeouts
// Example_setTimeouts demonstrates configuring request timeouts
func Example_setTimeouts() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Set custom timeouts (in seconds)
	config.ConnectTimeout = 15             // Connection timeout
	config.NonSubscribeRequestTimeout = 20 // Timeout for publish, history, etc.
	config.SubscribeRequestTimeout = 300   // Timeout for long-poll subscribe
	config.FileUploadRequestTimeout = 120  // Timeout for file uploads

	pn := pubnub.NewPubNub(config)

	if pn != nil {
		fmt.Println("Timeouts configured successfully")
	}

	// Output:
	// Timeouts configured successfully
}

// snippet.set_presence_timeout
// Example_setPresenceTimeout demonstrates configuring presence timeout
func Example_setPresenceTimeout() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Set presence timeout (how long before user appears offline)
	// The heartbeat interval is automatically calculated
	config.SetPresenceTimeout(60) // 60 seconds

	pn := pubnub.NewPubNub(config)

	if pn != nil && config.PresenceTimeout == 60 {
		fmt.Println("Presence timeout configured to 60 seconds")
	}

	// Output:
	// Presence timeout configured to 60 seconds
}

// snippet.set_presence_custom_interval
// Example_setPresenceCustomInterval demonstrates custom presence timeout and interval
func Example_setPresenceCustomInterval() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Set custom presence timeout and heartbeat interval
	config.SetPresenceTimeoutWithCustomInterval(120, 55) // timeout: 120s, interval: 55s

	pn := pubnub.NewPubNub(config)

	if pn != nil && config.PresenceTimeout == 120 {
		fmt.Println("Custom presence configuration set")
	}

	// Output:
	// Custom presence configuration set
}

// snippet.enable_secure_connection
// Example_enableSecureConnection demonstrates enabling TLS/SSL
func Example_enableSecureConnection() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Enable secure connection (TLS/SSL) - enabled by default
	config.Secure = true

	pn := pubnub.NewPubNub(config)

	if pn != nil && config.Secure {
		fmt.Println("Secure connection enabled")
	}

	// Output:
	// Secure connection enabled
}

// snippet.set_reconnection_policy
// Example_setReconnectionPolicy demonstrates configuring reconnection behavior
func Example_setReconnectionPolicy() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Set reconnection policy
	// Options: PNNonePolicy, PNLinearPolicy, PNExponentialPolicy
	config.PNReconnectionPolicy = pubnub.PNLinearPolicy
	config.MaximumReconnectionRetries = 10 // Maximum reconnection attempts

	pn := pubnub.NewPubNub(config)

	if pn != nil {
		fmt.Println("Reconnection policy configured")
	}

	// Output:
	// Reconnection policy configured
}

// snippet.end

/*
// snippet.includes_enable_logging
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"
	"log"
	"io"
	"os"

	pubnub "github.com/pubnub/go/v7"
)

// snippet.end
*/

// snippet.enable_logging
// Example_enableLogging demonstrates enabling SDK logging
func Example_enableLogging() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Enable logging to stdout
	config.Log = log.New(os.Stdout, "PubNub: ", log.Ldate|log.Ltime|log.Lshortfile)

	// snippet.hide
	// For testing, override with discard logger to avoid test output
	config.Log = log.New(io.Discard, "PubNub: ", log.Ldate|log.Ltime|log.Lshortfile)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	if pn != nil && config.Log != nil {
		fmt.Println("Logging enabled")
	}

	// Output:
	// Logging enabled
}

// snippet.set_filter_expression
// Example_setFilterExpression demonstrates using filter expressions for subscriptions
func Example_setFilterExpression() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Set filter expression to only receive messages matching the criteria
	config.FilterExpression = "language == 'english'"

	pn := pubnub.NewPubNub(config)

	if pn != nil && config.FilterExpression != "" {
		fmt.Println("Filter expression configured")
	}

	// Output:
	// Filter expression configured
}

// snippet.set_origin
// Example_setOrigin demonstrates setting a custom origin
func Example_setOrigin() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Set custom origin (default is "ps.pndsn.com")
	config.Origin = "ps.pndsn.com"

	pn := pubnub.NewPubNub(config)

	if pn != nil && config.Origin != "" {
		fmt.Println("Custom origin configured")
	}

	// Output:
	// Custom origin configured
}

// snippet.suppress_leave_events
// Example_suppressLeaveEvents demonstrates suppressing leave events
func Example_suppressLeaveEvents() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Suppress leave events when unsubscribing
	config.SuppressLeaveEvents = true

	pn := pubnub.NewPubNub(config)

	if pn != nil && config.SuppressLeaveEvents {
		fmt.Println("Leave events suppressed")
	}

	// Output:
	// Leave events suppressed
}

// snippet.set_max_workers
// Example_setMaxWorkers demonstrates configuring maximum concurrent workers
func Example_setMaxWorkers() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Set maximum number of workers for concurrent requests
	config.MaxWorkers = 50 // Default is 20

	pn := pubnub.NewPubNub(config)

	if pn != nil && config.MaxWorkers == 50 {
		fmt.Println("Max workers configured to 50")
	}

	// Output:
	// Max workers configured to 50
}

// snippet.set_user_id
// Example_setUserId demonstrates setting user ID
func Example_setUserId() {
	// Create a new configuration with initial user ID
	config := pubnub.NewConfigWithUserId(pubnub.UserId("initial-user"))

	// Set a new user ID after configuration creation
	config.SetUserId(pubnub.UserId("myUniqueUserId"))
}

// snippet.get_user_id
// Example_getUserId demonstrates getting user ID
func Example_getUserId() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("myUniqueUserId"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	pn := pubnub.NewPubNub(config)

	// Get the user ID from configuration
	userId := config.GetUserId()

	if pn != nil {
		fmt.Printf("User ID: %s\n", userId)
	}

	// Output:
	// User ID: myUniqueUserId
}

// snippet.set_auth_key
// Example_setAuthKey demonstrates setting authentication key
func Example_setAuthKey() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// Set authentication key for Access Manager
	config.AuthKey = "my_auth_key"
}

// snippet.get_auth_key
// Example_getAuthKey demonstrates getting authentication key
func Example_getAuthKey() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.AuthKey = "my_auth_key"

	pn := pubnub.NewPubNub(config)

	// Get the authentication key from configuration
	authKey := config.AuthKey

	if pn != nil {
		fmt.Printf("Auth key: %s\n", authKey)
	}

	// Output:
	// Auth key: my_auth_key
}

// snippet.get_filter_expression
// Example_getFilterExpression demonstrates getting filter expression
func Example_getFilterExpression() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.FilterExpression = "language == 'english'"

	pn := pubnub.NewPubNub(config)

	// Get the filter expression from configuration
	filterExpr := config.FilterExpression

	if pn != nil {
		fmt.Printf("Filter expression: %s\n", filterExpr)
	}

	// Output:
	// Filter expression: language == 'english'
}

// snippet.add_listeners
// addListeners demonstrates adding event listeners
func addListeners() {
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case signal := <-listener.Signal:
				//Channel
				fmt.Println(signal.Channel)
				//Subscription
				fmt.Println(signal.Subscription)
				//Payload
				fmt.Println(signal.Message)
				//Publisher ID
				fmt.Println(signal.Publisher)
				//Timetoken
				fmt.Println(signal.Timetoken)
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNDisconnectedCategory:
					// this is the expected category for an unsubscribe. This means there
					// was no error in unsubscribing from everything
				case pubnub.PNConnectedCategory:
					// this is expected for a subscribe, this means there is no error or issue whatsoever
				case pubnub.PNReconnectedCategory:
					// this usually occurs if subscribe temporarily fails but reconnects. This means
					// there was an error but there is no longer any issue
				case pubnub.PNAccessDeniedCategory:
					// this means that Access Manager does allow this client to subscribe to this
					// channel and channel group configuration. This is another explicit error
				}
			case message := <-listener.Message:
				//Channel
				fmt.Println(message.Channel)
				//Subscription
				fmt.Println(message.Subscription)
				//Payload
				fmt.Println(message.Message)
				//Publisher ID
				fmt.Println(message.Publisher)
				//Timetoken
				fmt.Println(message.Timetoken)
			case presence := <-listener.Presence:
				fmt.Println(presence.Event)
				//Channel
				fmt.Println(presence.Channel)
				//Subscription
				fmt.Println(presence.Subscription)
				//Timetoken
				fmt.Println(presence.Timetoken)
				//Occupancy
				fmt.Println(presence.Occupancy)
			case uuidEvent := <-listener.UUIDEvent:
				fmt.Printf("uuidEvent.Channel: %s\n", uuidEvent.Channel)
				fmt.Printf("uuidEvent.SubscribedChannel: %s\n", uuidEvent.SubscribedChannel)
				fmt.Printf("uuidEvent.Event: %s\n", uuidEvent.Event)
				fmt.Printf("uuidEvent.UUID: %s\n", uuidEvent.UUID)
				fmt.Printf("uuidEvent.Description: %s\n", uuidEvent.Description)
				fmt.Printf("uuidEvent.Timestamp: %s\n", uuidEvent.Timestamp)
				fmt.Printf("uuidEvent.Name: %s\n", uuidEvent.Name)
				fmt.Printf("uuidEvent.ExternalID: %s\n", uuidEvent.ExternalID)
				fmt.Printf("uuidEvent.ProfileURL: %s\n", uuidEvent.ProfileURL)
				fmt.Printf("uuidEvent.Email: %s\n", uuidEvent.Email)
				fmt.Printf("uuidEvent.Updated: %s\n", uuidEvent.Updated)
				fmt.Printf("uuidEvent.ETag: %s\n", uuidEvent.ETag)
				fmt.Printf("uuidEvent.Custom: %v\n", uuidEvent.Custom)
			case channelEvent := <-listener.ChannelEvent:
				fmt.Printf("channelEvent.Channel: %s\n", channelEvent.Channel)
				fmt.Printf("channelEvent.SubscribedChannel: %s\n", channelEvent.SubscribedChannel)
				fmt.Printf("channelEvent.Event: %s\n", channelEvent.Event)
				fmt.Printf("channelEvent.Channel: %s\n", channelEvent.Channel)
				fmt.Printf("channelEvent.Description: %s\n", channelEvent.Description)
				fmt.Printf("channelEvent.Timestamp: %s\n", channelEvent.Timestamp)
				fmt.Printf("channelEvent.Updated: %s\n", channelEvent.Updated)
				fmt.Printf("channelEvent.ETag: %s\n", channelEvent.ETag)
				fmt.Printf("channelEvent.Custom: %v\n", channelEvent.Custom)
			case membershipEvent := <-listener.MembershipEvent:
				fmt.Printf("membershipEvent.Channel: %s\n", membershipEvent.Channel)
				fmt.Printf("membershipEvent.SubscribedChannel: %s\n", membershipEvent.SubscribedChannel)
				fmt.Printf("membershipEvent.Event: %s\n", membershipEvent.Event)
				fmt.Printf("membershipEvent.Channel: %s\n", membershipEvent.Channel)
				fmt.Printf("membershipEvent.UUID: %s\n", membershipEvent.UUID)
				fmt.Printf("membershipEvent.Description: %s\n", membershipEvent.Description)
				fmt.Printf("membershipEvent.Timestamp: %s\n", membershipEvent.Timestamp)
				fmt.Printf("membershipEvent.Custom: %v\n", membershipEvent.Custom)
			case messageActionsEvent := <-listener.MessageActionsEvent:
				fmt.Printf("messageActionsEvent.Channel: %s\n", messageActionsEvent.Channel)
				fmt.Printf("messageActionsEvent.SubscribedChannel: %s\n", messageActionsEvent.SubscribedChannel)
				fmt.Printf("messageActionsEvent.Event: %s\n", messageActionsEvent.Event)
				fmt.Printf("messageActionsEvent.Data.ActionType: %s\n", messageActionsEvent.Data.ActionType)
				fmt.Printf("messageActionsEvent.Data.ActionValue: %s\n", messageActionsEvent.Data.ActionValue)
				fmt.Printf("messageActionsEvent.Data.ActionTimetoken: %s\n", messageActionsEvent.Data.ActionTimetoken)
				fmt.Printf("messageActionsEvent.Data.MessageTimetoken: %s\n", messageActionsEvent.Data.MessageTimetoken)
			case file := <-listener.File:
				fmt.Printf("file.File.PNMessage.Text: %s\n", file.File.PNMessage.Text)
				fmt.Printf("file.File.PNFile.Name: %s\n", file.File.PNFile.Name)
				fmt.Printf("file.File.PNFile.ID: %s\n", file.File.PNFile.ID)
				fmt.Printf("file.File.PNFile.URL: %s\n", file.File.PNFile.URL)
				fmt.Printf("file.Channel: %s\n", file.Channel)
				fmt.Printf("file.Timetoken: %d\n", file.Timetoken)
				fmt.Printf("file.SubscribedChannel: %s\n", file.SubscribedChannel)
				fmt.Printf("file.Publisher: %s\n", file.Publisher)
			}
		}
	}()
}

// snippet.remove_listeners
// removeListeners demonstrates removing event listeners
func removeListeners() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("myUniqueUserId"))
	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()

	pn.AddListener(listener)

	// some time later
	pn.RemoveListener(listener)
}

// snippet.handling_disconnects
// handlingDisconnects demonstrates handling disconnect events
func handlingDisconnects() {
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNDisconnectedCategory:
					// handle disconnect here
				}
			case <-listener.Message:
			case <-listener.Presence:
			}
		}
	}()
}

// snippet.end
