package pubnub

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// Log Message Interface and Types
// ============================================================================

// LogMessage is the interface for all log messages.
// All log message types implement this interface to provide access to common fields.
type LogMessage interface {
	GetTimestamp() time.Time
	GetInstanceID() string
	GetLogLevel() PNLogLevel
	GetMessage() string
	GetCallsite() string
}

// BaseLogMessage contains common fields for all log messages.
// All specific message types embed this struct.
type BaseLogMessage struct {
	Timestamp  time.Time
	InstanceID string
	LogLevel   PNLogLevel
	Message    string
	Callsite   string
}

// Interface method implementations for BaseLogMessage
func (b BaseLogMessage) GetTimestamp() time.Time { return b.Timestamp }
func (b BaseLogMessage) GetInstanceID() string   { return b.InstanceID }
func (b BaseLogMessage) GetLogLevel() PNLogLevel { return b.LogLevel }
func (b BaseLogMessage) GetMessage() string      { return b.Message }
func (b BaseLogMessage) GetCallsite() string     { return b.Callsite }

// SimpleLogMessage is used for general purpose logging.
type SimpleLogMessage struct {
	BaseLogMessage
}

// NetworkRequestLogMessage contains HTTP request details.
type NetworkRequestLogMessage struct {
	BaseLogMessage
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

// NetworkResponseLogMessage contains HTTP response details.
type NetworkResponseLogMessage struct {
	BaseLogMessage
	StatusCode int
	URL        string
	Body       string
}

// ErrorLogMessage contains error information with context.
type ErrorLogMessage struct {
	BaseLogMessage
	Error     error
	ErrorName string
	Operation OperationType
}

// UserInputLogMessage contains API call parameters.
type UserInputLogMessage struct {
	BaseLogMessage
	Operation  OperationType
	Parameters map[string]interface{}
}

// ============================================================================
// Log Message Formatting (String methods)
// ============================================================================

// String implements fmt.Stringer for SimpleLogMessage
func (msg SimpleLogMessage) String() string {
	return fmt.Sprintf("%s %s", formatLogBase(msg), msg.Message)
}

// String implements fmt.Stringer for NetworkRequestLogMessage
func (msg NetworkRequestLogMessage) String() string {
	return fmt.Sprintf("%s Sending HTTP request %s %s with headers: %v body: %s",
		formatLogBase(msg), msg.Method, msg.URL, msg.Headers, msg.Body)
}

// String implements fmt.Stringer for NetworkResponseLogMessage
func (msg NetworkResponseLogMessage) String() string {
	return fmt.Sprintf("%s Received response with %d content %s for request url %s",
		formatLogBase(msg), msg.StatusCode, msg.Body, msg.URL)
}

// String implements fmt.Stringer for ErrorLogMessage
func (msg ErrorLogMessage) String() string {
	if msg.Operation != 0 {
		return fmt.Sprintf("%s Error %s in %s: %v",
			formatLogBase(msg), msg.ErrorName, msg.Operation.String(), msg.Error)
	}
	return fmt.Sprintf("%s Error %s: %v",
		formatLogBase(msg), msg.ErrorName, msg.Error)
}

// String implements fmt.Stringer for UserInputLogMessage
func (msg UserInputLogMessage) String() string {
	return fmt.Sprintf("%s %s with parameters:\n%s",
		formatLogBase(msg), msg.Operation.String(), formatParamsMap(msg.Parameters))
}

// formatLogBase creates the common prefix for all log messages
func formatLogBase(logMsg LogMessage) string {
	base := fmt.Sprintf("%s PubNub-%s %s",
		logMsg.GetTimestamp().Format("02/01/2006 15:04:05.000"),
		logMsg.GetInstanceID(),
		logMsg.GetLogLevel().String())

	if callsite := logMsg.GetCallsite(); callsite != "" {
		base = fmt.Sprintf("%s %s", base, callsite)
	}

	return base
}

// formatParamsMap formats user parameters as a readable string
func formatParamsMap(params map[string]interface{}) string {
	if len(params) == 0 {
		return "  (none)"
	}

	var lines []string
	for key, value := range params {
		lines = append(lines, fmt.Sprintf("  %s: %v", key, value))
	}
	return strings.Join(lines, "\n")
}

// ============================================================================
// Logger Interface
// ============================================================================

// PNLogger is the interface for custom logger implementations.
// Users can implement this interface to provide their own logging mechanism.
//
// The Log method receives different LogMessage implementations depending on
// the logging context. All messages implement the LogMessage interface.
//
// Available message types:
//   - SimpleLogMessage: General purpose logs
//   - NetworkRequestLogMessage: HTTP request details
//   - NetworkResponseLogMessage: HTTP response details
//   - ErrorLogMessage: Error information with context
//   - UserInputLogMessage: API call parameters
//
// Custom loggers can handle specific types using type switches, or use only
// the String() method or interface methods for generic handling.
// See DefaultLogger for a reference implementation.
type PNLogger interface {
	// Log outputs a log message
	Log(logMsg LogMessage)
	// GetMinLogLevel returns the minimum log level for this logger
	GetMinLogLevel() PNLogLevel
}

// Compile-time check to ensure DefaultLogger implements PNLogger
var _ PNLogger = (*DefaultLogger)(nil)

// ============================================================================
// Default Logger Implementation
// ============================================================================

// DefaultLogger is the default console logger implementation.
// It outputs formatted log messages to the specified writer (typically os.Stdout).
// It serves as a reference implementation for custom loggers.
type DefaultLogger struct {
	minLogLevel PNLogLevel
	output      io.Writer
	mu          sync.Mutex
}

// NewDefaultLogger creates a new default console logger that writes to os.Stdout.
// minLogLevel: the minimum log level to output
func NewDefaultLogger(minLogLevel PNLogLevel) *DefaultLogger {
	return &DefaultLogger{
		minLogLevel: minLogLevel,
		output:      os.Stdout,
	}
}

// NewDefaultLoggerWithWriter creates a new default logger with a custom output writer.
// This is useful for testing or redirecting logs to a file.
// minLogLevel: the minimum log level to output
// writer: the destination for log output
func NewDefaultLoggerWithWriter(minLogLevel PNLogLevel, writer io.Writer) *DefaultLogger {
	return &DefaultLogger{
		minLogLevel: minLogLevel,
		output:      writer,
	}
}

// Log implements PNLogger interface.
// It outputs the log message using its String() method if the log level is sufficient.
func (dl *DefaultLogger) Log(logMsg LogMessage) {
	if !dl.shouldLog(logMsg.GetLogLevel()) {
		return
	}

	dl.mu.Lock()
	defer dl.mu.Unlock()

	// String() method is called automatically
	fmt.Fprintln(dl.output, logMsg)
}

// GetMinLogLevel returns the minimum log level for this logger
func (dl *DefaultLogger) GetMinLogLevel() PNLogLevel {
	return dl.minLogLevel
}

// shouldLog checks if a message at the given level should be logged
func (dl *DefaultLogger) shouldLog(level PNLogLevel) bool {
	if dl.minLogLevel == PNLogLevelNone {
		return false
	}
	return level >= dl.minLogLevel
}
