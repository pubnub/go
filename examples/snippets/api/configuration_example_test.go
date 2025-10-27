// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"
	"io"
	"log"
	"os"

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

// snippet.end
