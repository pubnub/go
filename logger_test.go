package pubnub

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Test Logger Implementation
// ============================================================================

type testLogger struct {
	logs     []LogMessage
	minLevel PNLogLevel
}

func (t *testLogger) Log(logMsg LogMessage) {
	t.logs = append(t.logs, logMsg)
}

func (t *testLogger) GetMinLogLevel() PNLogLevel {
	return t.minLevel
}

func (t *testLogger) reset() {
	t.logs = []LogMessage{}
}

func (t *testLogger) count() int {
	return len(t.logs)
}

func (t *testLogger) getLastLog() LogMessage {
	if len(t.logs) == 0 {
		return nil
	}
	return t.logs[len(t.logs)-1]
}

// ============================================================================
// LoggerManager Tests
// ============================================================================

func TestLoggerManager_AddLogger(t *testing.T) {
	mgr := newLoggerManager("test-instance", []PNLogger{})
	logger1 := &testLogger{minLevel: PNLogLevelInfo}
	logger2 := &testLogger{minLevel: PNLogLevelDebug}

	mgr.AddLogger(logger1)
	mgr.AddLogger(logger2)

	mgr.mu.RLock()
	loggerCount := len(mgr.loggers)
	mgr.mu.RUnlock()

	if loggerCount != 2 {
		t.Errorf("Expected 2 loggers, got %d", loggerCount)
	}
}

func TestLoggerManager_RemoveAllLoggers(t *testing.T) {
	logger1 := &testLogger{minLevel: PNLogLevelInfo}
	logger2 := &testLogger{minLevel: PNLogLevelDebug}
	mgr := newLoggerManager("test-instance", []PNLogger{logger1, logger2})

	mgr.RemoveAllLoggers()

	mgr.mu.RLock()
	loggerCount := len(mgr.loggers)
	mgr.mu.RUnlock()

	if loggerCount != 0 {
		t.Errorf("Expected 0 loggers after removal, got %d", loggerCount)
	}
}

func TestLoggerManager_LogSimple(t *testing.T) {
	logger := &testLogger{minLevel: PNLogLevelInfo}
	mgr := newLoggerManager("test-instance", []PNLogger{logger})

	mgr.LogSimple(PNLogLevelInfo, "test message", false)

	if logger.count() != 1 {
		t.Errorf("Expected 1 log, got %d", logger.count())
	}

	msg := logger.getLastLog()
	if msg == nil {
		t.Fatal("No log message received")
	}

	if msg.GetMessage() != "test message" {
		t.Errorf("Expected 'test message', got '%s'", msg.GetMessage())
	}

	if msg.GetLogLevel() != PNLogLevelInfo {
		t.Errorf("Expected PNLogLevelInfo, got %v", msg.GetLogLevel())
	}
}

func TestLoggerManager_LogSimpleWithCallsite(t *testing.T) {
	logger := &testLogger{minLevel: PNLogLevelInfo}
	mgr := newLoggerManager("test-instance", []PNLogger{logger})

	mgr.LogSimple(PNLogLevelInfo, "test message", true)

	if logger.count() != 1 {
		t.Errorf("Expected 1 log, got %d", logger.count())
	}

	msg := logger.getLastLog()
	if msg == nil {
		t.Fatal("No log message received")
	}

	if msg.GetCallsite() == "" {
		t.Error("Expected callsite to be set")
	}

	if !strings.Contains(msg.GetCallsite(), "logger_test.go") {
		t.Errorf("Expected callsite to contain 'logger_test.go', got '%s'", msg.GetCallsite())
	}
}

func TestLoggerManager_LogError(t *testing.T) {
	logger := &testLogger{minLevel: PNLogLevelError}
	mgr := newLoggerManager("test-instance", []PNLogger{logger})

	testErr := errors.New("test error")
	mgr.LogError(testErr, "TestError", PNPublishOperation, false)

	if logger.count() != 1 {
		t.Errorf("Expected 1 log, got %d", logger.count())
	}

	msg := logger.getLastLog()
	if msg == nil {
		t.Fatal("No log message received")
	}

	errorMsg, ok := msg.(ErrorLogMessage)
	if !ok {
		t.Fatal("Expected ErrorLogMessage type")
	}

	if errorMsg.Error != testErr {
		t.Errorf("Expected error to be %v, got %v", testErr, errorMsg.Error)
	}

	if errorMsg.ErrorName != "TestError" {
		t.Errorf("Expected ErrorName 'TestError', got '%s'", errorMsg.ErrorName)
	}

	if errorMsg.Operation != PNPublishOperation {
		t.Errorf("Expected PNPublishOperation, got %v", errorMsg.Operation)
	}
}

func TestLoggerManager_LogUserInput(t *testing.T) {
	logger := &testLogger{minLevel: PNLogLevelDebug}
	mgr := newLoggerManager("test-instance", []PNLogger{logger})

	params := map[string]interface{}{
		"Channel": "test-channel",
		"Message": "test-message",
	}

	mgr.LogUserInput(PNLogLevelDebug, PNPublishOperation, params, false)

	if logger.count() != 1 {
		t.Errorf("Expected 1 log, got %d", logger.count())
	}

	msg := logger.getLastLog()
	if msg == nil {
		t.Fatal("No log message received")
	}

	userInputMsg, ok := msg.(UserInputLogMessage)
	if !ok {
		t.Fatal("Expected UserInputLogMessage type")
	}

	if userInputMsg.Operation != PNPublishOperation {
		t.Errorf("Expected PNPublishOperation, got %v", userInputMsg.Operation)
	}

	if userInputMsg.Parameters["Channel"] != "test-channel" {
		t.Errorf("Expected Channel 'test-channel', got '%v'", userInputMsg.Parameters["Channel"])
	}
}

func TestLoggerManager_MultipleLoggers(t *testing.T) {
	logger1 := &testLogger{minLevel: PNLogLevelInfo}
	logger2 := &testLogger{minLevel: PNLogLevelDebug}
	mgr := newLoggerManager("test-instance", []PNLogger{logger1, logger2})

	mgr.LogSimple(PNLogLevelInfo, "test message", false)

	if logger1.count() != 1 {
		t.Errorf("Logger1: Expected 1 log, got %d", logger1.count())
	}

	if logger2.count() != 1 {
		t.Errorf("Logger2: Expected 1 log, got %d", logger2.count())
	}
}

func TestLoggerManager_LogLevelFiltering(t *testing.T) {
	logger := &testLogger{minLevel: PNLogLevelWarn}
	mgr := newLoggerManager("test-instance", []PNLogger{logger})

	// This should be logged (WARN >= WARN)
	mgr.LogSimple(PNLogLevelWarn, "warning message", false)

	// This should be logged (ERROR > WARN)
	mgr.LogSimple(PNLogLevelError, "error message", false)

	// This should NOT be logged (INFO < WARN)
	mgr.LogSimple(PNLogLevelInfo, "info message", false)

	// This should NOT be logged (DEBUG < WARN)
	mgr.LogSimple(PNLogLevelDebug, "debug message", false)

	if logger.count() != 2 {
		t.Errorf("Expected 2 logs, got %d", logger.count())
	}
}

// ============================================================================
// DefaultLogger Tests
// ============================================================================

func TestDefaultLogger_Log(t *testing.T) {
	var buf bytes.Buffer
	logger := NewDefaultLoggerWithWriter(PNLogLevelInfo, &buf)

	msg := SimpleLogMessage{
		BaseLogMessage: BaseLogMessage{
			Timestamp:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			InstanceID: "test-id",
			LogLevel:   PNLogLevelInfo,
			Message:    "test message",
		},
	}

	logger.Log(msg)

	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Errorf("Expected output to contain 'test message', got: %s", output)
	}

	if !strings.Contains(output, "PubNub-test-id") {
		t.Errorf("Expected output to contain 'PubNub-test-id', got: %s", output)
	}

	if !strings.Contains(output, "Info") {
		t.Errorf("Expected output to contain 'Info', got: %s", output)
	}
}

func TestDefaultLogger_GetMinLogLevel(t *testing.T) {
	logger := NewDefaultLogger(PNLogLevelDebug)

	if logger.GetMinLogLevel() != PNLogLevelDebug {
		t.Errorf("Expected PNLogLevelDebug, got %v", logger.GetMinLogLevel())
	}
}

// ============================================================================
// Log Level Tests
// ============================================================================

func TestPNLogLevel_String(t *testing.T) {
	tests := []struct {
		level    PNLogLevel
		expected string
	}{
		{PNLogLevelNone, "None"},
		{PNLogLevelError, "Error"},
		{PNLogLevelWarn, "Warn"},
		{PNLogLevelInfo, "Info"},
		{PNLogLevelDebug, "Debug"},
		{PNLogLevelTrace, "Trace"},
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, test.level.String())
		}
	}
}

// ============================================================================
// Log Message Tests
// ============================================================================

func TestSimpleLogMessage_String(t *testing.T) {
	msg := SimpleLogMessage{
		BaseLogMessage: BaseLogMessage{
			Timestamp:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			InstanceID: "test-id",
			LogLevel:   PNLogLevelInfo,
			Message:    "test message",
		},
	}

	str := msg.String()

	if !strings.Contains(str, "test message") {
		t.Errorf("Expected string to contain 'test message', got: %s", str)
	}

	if !strings.Contains(str, "PubNub-test-id") {
		t.Errorf("Expected string to contain 'PubNub-test-id', got: %s", str)
	}

	if !strings.Contains(str, "Info") {
		t.Errorf("Expected string to contain 'Info', got: %s", str)
	}
}

func TestErrorLogMessage_String(t *testing.T) {
	msg := ErrorLogMessage{
		BaseLogMessage: BaseLogMessage{
			Timestamp:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			InstanceID: "test-id",
			LogLevel:   PNLogLevelError,
		},
		Error:     errors.New("test error"),
		ErrorName: "TestError",
		Operation: PNPublishOperation,
	}

	str := msg.String()

	if !strings.Contains(str, "TestError") {
		t.Errorf("Expected string to contain 'TestError', got: %s", str)
	}

	if !strings.Contains(str, "test error") {
		t.Errorf("Expected string to contain 'test error', got: %s", str)
	}

	if !strings.Contains(str, "Publish") {
		t.Errorf("Expected string to contain 'Publish', got: %s", str)
	}
}

func TestUserInputLogMessage_String(t *testing.T) {
	msg := UserInputLogMessage{
		BaseLogMessage: BaseLogMessage{
			Timestamp:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			InstanceID: "test-id",
			LogLevel:   PNLogLevelDebug,
		},
		Operation: PNPublishOperation,
		Parameters: map[string]interface{}{
			"Channel": "test-channel",
			"Message": "test-message",
		},
	}

	str := msg.String()

	if !strings.Contains(str, "Publish") {
		t.Errorf("Expected string to contain 'Publish', got: %s", str)
	}

	if !strings.Contains(str, "Channel") {
		t.Errorf("Expected string to contain 'Channel', got: %s", str)
	}

	if !strings.Contains(str, "test-channel") {
		t.Errorf("Expected string to contain 'test-channel', got: %s", str)
	}
}

func TestNetworkRequestLogMessage_String(t *testing.T) {
	msg := NetworkRequestLogMessage{
		BaseLogMessage: BaseLogMessage{
			Timestamp:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			InstanceID: "test-id",
			LogLevel:   PNLogLevelDebug,
		},
		Method: "POST",
		URL:    "https://ps.pndsn.com/publish",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "test body",
	}

	str := msg.String()

	if !strings.Contains(str, "POST") {
		t.Errorf("Expected string to contain 'POST', got: %s", str)
	}

	if !strings.Contains(str, "https://ps.pndsn.com/publish") {
		t.Errorf("Expected string to contain URL, got: %s", str)
	}
}

func TestNetworkResponseLogMessage_String(t *testing.T) {
	msg := NetworkResponseLogMessage{
		BaseLogMessage: BaseLogMessage{
			Timestamp:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			InstanceID: "test-id",
			LogLevel:   PNLogLevelDebug,
		},
		StatusCode: 200,
		URL:        "https://ps.pndsn.com/publish",
		Body:       "test response",
	}

	str := msg.String()

	if !strings.Contains(str, "200") {
		t.Errorf("Expected string to contain '200', got: %s", str)
	}

	if !strings.Contains(str, "https://ps.pndsn.com/publish") {
		t.Errorf("Expected string to contain URL, got: %s", str)
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkLoggerManager_LogSimple(b *testing.B) {
	logger := NewDefaultLogger(PNLogLevelInfo)
	mgr := newLoggerManager("test-instance", []PNLogger{logger})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.LogSimple(PNLogLevelInfo, "benchmark test message", false)
	}
}

func BenchmarkLoggerManager_LogError(b *testing.B) {
	logger := NewDefaultLogger(PNLogLevelError)
	mgr := newLoggerManager("test-instance", []PNLogger{logger})

	testErr := errors.New("benchmark error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.LogError(testErr, "BenchmarkError", PNPublishOperation, false)
	}
}

func BenchmarkLoggerManager_LogUserInput(b *testing.B) {
	logger := NewDefaultLogger(PNLogLevelDebug)
	mgr := newLoggerManager("test-instance", []PNLogger{logger})

	params := map[string]interface{}{
		"Channel": "test-channel",
		"Message": "test-message",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.LogUserInput(PNLogLevelDebug, PNPublishOperation, params, false)
	}
}
