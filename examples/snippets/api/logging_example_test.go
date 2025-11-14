// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"
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

// snippet.logging-simple
// Example_loggingSimple demonstrates how to enable simple console logging in the PubNub Go SDK
func Example_loggingSimple() {
	// Create a new PubNub configuration
	config := pubnub.NewConfigWithUserId(pubnub.UserId("loggingDemoUser"))

	// Set the subscribe and publish keys
	// Replace "demo" with your actual keys from the PubNub Admin Portal
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Create a simple console logger with INFO level (shows Info, Warn, Error)
	// Logs will be written to stdout with timestamps
	simpleLogger := pubnub.NewDefaultLogger(pubnub.PNLogLevelInfo)

	// Add the logger to the config
	config.Loggers = []pubnub.PNLogger{simpleLogger}

	// Initialize PubNub with the configured settings
	pn := pubnub.NewPubNub(config)

	// Perform operations - they will be logged to console
	_, _, err := pn.Time().Execute()
	if err != nil {
		fmt.Println("Error fetching time:", err)
	} else {
		fmt.Println("Time fetched successfully")
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
		fmt.Println("Check console output for detailed logging information")
	}
}

// snippet.end

// snippet.logging-file
// Example_loggingToFile demonstrates how to log PubNub operations to a file
func Example_loggingToFile() {
	// Create a new PubNub configuration
	config := pubnub.NewConfigWithUserId(pubnub.UserId("loggingDemoUser"))

	// Set the subscribe and publish keys
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

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
		return
	}

	// snippet.hide
	defer f.Close()
	// snippet.show

	fmt.Println("Logging enabled, writing to", logfileName)

	// Create a logger that writes to the file with DEBUG level
	// This will log Debug, Info, Warn, and Error messages
	fileLogger := pubnub.NewDefaultLoggerWithWriter(pubnub.PNLogLevelDebug, f)

	// Add the logger to the config
	config.Loggers = []pubnub.PNLogger{fileLogger}

	// Initialize PubNub with the configured settings
	pn := pubnub.NewPubNub(config)

	// Perform operations - they will be logged to the file
	_, _, err = pn.Time().Execute()
	if err != nil {
		fmt.Println("Error fetching time:", err)
	} else {
		fmt.Println("Time fetched successfully, check the log file for details")
	}

	// Publish a message to demonstrate logging
	_, status, err := pn.Publish().
		Channel("logging-demo-channel").
		Message("Hello from File Logging Example").
		Execute()

	if err != nil {
		fmt.Printf("Error publishing message: %v\n", err)
	} else {
		fmt.Printf("Publish status: %d\n", status.StatusCode)
		fmt.Println("Check the log file for detailed logging information")
	}

	fmt.Println("Example complete. Logging information has been saved to", logfileName)
}

// snippet.end

// snippet.logging-multiple
// Example_loggingMultiple demonstrates how to use multiple loggers with different log levels
func Example_loggingMultiple() {
	// Create a new PubNub configuration
	config := pubnub.NewConfigWithUserId(pubnub.UserId("loggingDemoUser"))

	// Set the subscribe and publish keys
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Create multiple loggers for different purposes

	// 1. Console logger for errors only (production-friendly)
	consoleLogger := pubnub.NewDefaultLogger(pubnub.PNLogLevelError)

	// 2. File logger for detailed debugging
	debugLogFile := "pubnub-debug.log"
	// snippet.hide
	defer os.Remove(debugLogFile)
	// snippet.show

	f, err := os.OpenFile(debugLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening debug log file:", err.Error())
		return
	}
	// snippet.hide
	defer f.Close()
	// snippet.show

	debugLogger := pubnub.NewDefaultLoggerWithWriter(pubnub.PNLogLevelDebug, f)

	// Add both loggers to the config
	// Errors will go to console, all debug info will go to file
	config.Loggers = []pubnub.PNLogger{consoleLogger, debugLogger}

	// Initialize PubNub with the configured settings
	pn := pubnub.NewPubNub(config)

	fmt.Println("Multiple loggers configured:")
	fmt.Println("- Console: ERROR level only")
	fmt.Println("- File: DEBUG level and above")

	// Perform operations
	_, _, err = pn.Time().Execute()
	if err != nil {
		fmt.Println("Error fetching time:", err)
	} else {
		fmt.Println("Time fetched successfully")
	}

	// Publish a message
	_, status, err := pn.Publish().
		Channel("logging-demo-channel").
		Message("Hello with multiple loggers").
		Execute()

	if err != nil {
		fmt.Printf("Error publishing message: %v\n", err)
	} else {
		fmt.Printf("Publish status: %d\n", status.StatusCode)
	}

	fmt.Println("Errors (if any) appear on console, all details in", debugLogFile)
}

// snippet.end

// snippet.logging-custom
// Example_loggingCustom demonstrates how to create a custom logger implementation

// CustomStructuredLogger is an example implementation of a custom logger
// that outputs structured logs in JSON-like format
type CustomStructuredLogger struct {
	minLevel pubnub.PNLogLevel
}

// Log implements the PNLogger interface
func (l *CustomStructuredLogger) Log(logMessage pubnub.LogMessage) {
	// Only log messages at or above the configured level
	if logMessage.GetLogLevel() < l.minLevel {
		return
	}

	// Output structured log in a custom format
	// In a real implementation, you might output JSON or send to a logging service
	fmt.Printf("[%s] level=%s instance=%s message=%q\n",
		logMessage.GetTimestamp().Format("2006-01-02 15:04:05.000"),
		logMessage.GetLogLevel().String(),
		logMessage.GetInstanceID(),
		logMessage.GetMessage(),
	)

	// Handle specific message types using type assertions
	switch msg := logMessage.(type) {
	case pubnub.ErrorLogMessage:
		// If it's an error message, include error details
		fmt.Printf("  error_name=%s operation=%s error_details=%q\n",
			msg.ErrorName,
			msg.Operation.String(),
			msg.Error.Error(),
		)
	case pubnub.UserInputLogMessage:
		// If it's a user input message, include parameters
		fmt.Printf("  operation=%s user_params=%v\n",
			msg.Operation.String(),
			msg.Parameters,
		)
	case pubnub.NetworkRequestLogMessage:
		// If it's a network request, include HTTP details
		fmt.Printf("  method=%s url=%s\n",
			msg.Method,
			msg.URL,
		)
	}
}

// GetMinLogLevel implements the PNLogger interface
func (l *CustomStructuredLogger) GetMinLogLevel() pubnub.PNLogLevel {
	return l.minLevel
}

func Example_loggingCustom() {
	// Create a new PubNub configuration
	config := pubnub.NewConfigWithUserId(pubnub.UserId("loggingDemoUser"))

	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Create an instance of our custom logger
	customLogger := &CustomStructuredLogger{
		minLevel: pubnub.PNLogLevelInfo,
	}

	// Add the custom logger to the config
	config.Loggers = []pubnub.PNLogger{customLogger}

	// Initialize PubNub with the configured settings
	pn := pubnub.NewPubNub(config)

	fmt.Println("Custom structured logger configured")

	// Perform operations
	_, status, err := pn.Publish().
		Channel("logging-demo-channel").
		Message("Hello from custom logger").
		Execute()

	if err != nil {
		fmt.Printf("Error publishing message: %v\n", err)
	} else {
		fmt.Printf("Publish status: %d\n", status.StatusCode)
	}
}

// snippet.end
