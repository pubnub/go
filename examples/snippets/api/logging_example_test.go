// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"
	"log"
	"os"

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
	defer os.Remove(logfileName)                 // <- Skip this line (cleanup for tests)
	// snippet.show

When copying examples to your own code:
- Use your own publish/subscribe keys instead of the "demo" keys
- Remove any statements that are between snippet.hide and snippet.show (they're only for testing purposes)
*/

// snippet.logging
// Example_logging demonstrates how to enable logging in the PubNub Go SDK
func Example_logging() {
	// Create a new PubNub configuration
	config := pubnub.NewConfigWithUserId(pubnub.UserId("loggingDemoUser"))

	// Set the subscribe and publish keys
	// Replace "demo" with your actual keys from the PubNub Admin Portal
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Set up logging
	var infoLogger *log.Logger

	// Specify log file name
	logfileName := "pubnubMessaging.log"

	// snippet.hide
	defer os.Remove(logfileName)
	// snippet.show

	// Open log file, creating it if needed with append mode
	f, err := os.OpenFile(logfileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err.Error())
		fmt.Println("Logging disabled")
	} else {
		// snippet.hide
		defer f.Close()
		// snippet.show
		fmt.Println("Logging enabled, writing to", logfileName)

		// Create a new logger with timestamp and file information
		infoLogger = log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)

		// Set the logger in the PubNub config
		config.Log = infoLogger
		config.Log.SetPrefix("PubNub :=  ")
	}

	// Initialize PubNub with the configured settings
	pn := pubnub.NewPubNub(config)

	// Perform an operation to generate log entries
	_, _, err = pn.Time().Execute()
	if err != nil {
		fmt.Println("Error fetching time:", err)
	} else {
		fmt.Println("Time fetched successfully, check the log file for details")
	}

	// Publish a message to demonstrate logging
	_, status, err := pn.Publish().
		Channel("logging-demo-channel").
		Message("Hello from Logging Example").
		Execute()

	if err != nil {
		fmt.Printf("Error publishing message: %v\n", err)
	} else {
		fmt.Printf("Publish status: %d\n", status.StatusCode)
		fmt.Println("Check the log file for detailed logging information")
	}

	fmt.Println("Example complete. Logging information has been saved to", logfileName)

	// Output:
	// Logging enabled, writing to pubnubMessaging.log
	// Time fetched successfully, check the log file for details
	// Publish status: 200
	// Check the log file for detailed logging information
	// Example complete. Logging information has been saved to pubnubMessaging.log
}

// snippet.end
