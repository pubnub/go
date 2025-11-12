package pubnub

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// ============================================================================
// Logger Manager (Internal)
// ============================================================================

// loggerManager manages multiple loggers and provides convenient logging methods.
// This is an internal type used by the PubNub instance.
type loggerManager struct {
	loggers    []PNLogger
	instanceID string
	mu         sync.RWMutex
}

// newLoggerManager creates a new logger manager.
// instanceID: unique identifier for the PubNub instance
// loggers: slice of loggers to use (can be empty)
func newLoggerManager(instanceID string, loggers []PNLogger) *loggerManager {
	return &loggerManager{
		instanceID: instanceID,
		loggers:    loggers,
	}
}

// AddLogger adds a logger to the manager.
func (lm *loggerManager) AddLogger(logger PNLogger) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	lm.loggers = append(lm.loggers, logger)
}

// RemoveAllLoggers removes all loggers from the manager.
func (lm *loggerManager) RemoveAllLoggers() {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	lm.loggers = nil
}

// log dispatches a log message to all registered loggers.
// It checks each logger's minimum log level before dispatching.
func (lm *loggerManager) log(logMsg LogMessage) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if len(lm.loggers) == 0 {
		return
	}

	level := logMsg.GetLogLevel()
	for _, logger := range lm.loggers {
		if level >= logger.GetMinLogLevel() {
			logger.Log(logMsg)
		}
	}
}

// captureCallsite captures the caller's file and line number if requested.
// skip: number of stack frames to skip (typically 1 to skip the logging method itself)
// Returns formatted callsite string like "filename.go:123" or empty string if not captured.
func (lm *loggerManager) captureCallsite(includeCallsite bool, skip int) string {
	if !includeCallsite {
		return ""
	}
	if _, file, line, ok := runtime.Caller(skip); ok {
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return ""
}

// ============================================================================
// Simple Logging Method
// ============================================================================

// LogSimple logs a simple text message at the specified level.
// includeCallsite: if true, captures the file and line number of the caller
func (lm *loggerManager) LogSimple(level PNLogLevel, message string, includeCallsite bool) {
	lm.log(SimpleLogMessage{
		BaseLogMessage: BaseLogMessage{
			Timestamp:  time.Now(),
			InstanceID: lm.instanceID,
			LogLevel:   level,
			Message:    message,
			Callsite:   lm.captureCallsite(includeCallsite, 2),
		},
	})
}

// ============================================================================
// Specialized Logging Methods
// ============================================================================

// LogNetworkRequest logs an HTTP request.
// includeCallsite: if true, captures the file and line number of the caller
func (lm *loggerManager) LogNetworkRequest(level PNLogLevel, method, url string, headers map[string]string, body string, includeCallsite bool) {
	lm.log(NetworkRequestLogMessage{
		BaseLogMessage: BaseLogMessage{
			Timestamp:  time.Now(),
			InstanceID: lm.instanceID,
			LogLevel:   level,
			Message:    "HTTP Request",
			Callsite:   lm.captureCallsite(includeCallsite, 2),
		},
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}

// LogNetworkResponse logs an HTTP response.
// includeCallsite: if true, captures the file and line number of the caller
func (lm *loggerManager) LogNetworkResponse(level PNLogLevel, statusCode int, url, body string, includeCallsite bool) {
	lm.log(NetworkResponseLogMessage{
		BaseLogMessage: BaseLogMessage{
			Timestamp:  time.Now(),
			InstanceID: lm.instanceID,
			LogLevel:   level,
			Message:    "HTTP Response",
			Callsite:   lm.captureCallsite(includeCallsite, 2),
		},
		StatusCode: statusCode,
		URL:        url,
		Body:       body,
	})
}

// LogError logs an error with context.
// includeCallsite: if true, captures the file and line number of the caller
func (lm *loggerManager) LogError(err error, errorName string, operation OperationType, includeCallsite bool) {
	lm.log(ErrorLogMessage{
		BaseLogMessage: BaseLogMessage{
			Timestamp:  time.Now(),
			InstanceID: lm.instanceID,
			LogLevel:   PNLogLevelError,
			Message:    "Operation failed",
			Callsite:   lm.captureCallsite(includeCallsite, 2),
		},
		Error:     err,
		ErrorName: errorName,
		Operation: operation,
	})
}

// LogUserInput logs user-provided API parameters.
// includeCallsite: if true, captures the file and line number of the caller
func (lm *loggerManager) LogUserInput(level PNLogLevel, operation OperationType, parameters map[string]interface{}, includeCallsite bool) {
	lm.log(UserInputLogMessage{
		BaseLogMessage: BaseLogMessage{
			Timestamp:  time.Now(),
			InstanceID: lm.instanceID,
			LogLevel:   level,
			Message:    "API call",
			Callsite:   lm.captureCallsite(includeCallsite, 2),
		},
		Operation:  operation,
		Parameters: parameters,
	})
}
