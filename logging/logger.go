package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     LogLevel              `json:"level"`
	Message   string                `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Service   string                `json:"service"`
	RequestID string                `json:"request_id,omitempty"`
	UserID    uint                  `json:"user_id,omitempty"`
}

type Logger struct {
	level   LogLevel
	service string
	fields  map[string]interface{}
}

func NewLogger(level LogLevel, service string) *Logger {
	return &Logger{
		level:   level,
		service: service,
		fields:  make(map[string]interface{}),
	}
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := &Logger{
		level:   l.level,
		service: l.service,
		fields:  make(map[string]interface{}),
	}
	
	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	
	// Add new field
	newLogger.fields[key] = value
	return newLogger
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newLogger := &Logger{
		level:   l.level,
		service: l.service,
		fields:  make(map[string]interface{}),
	}
	
	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	
	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	
	return newLogger
}

func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithField("request_id", requestID)
}

func (l *Logger) WithUserID(userID uint) *Logger {
	return l.WithField("user_id", userID)
}

func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	if l.shouldLog(level) {
		entry := LogEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Level:     level,
			Message:   message,
			Service:   l.service,
			Fields:    l.mergeFields(fields),
		}
		
		jsonData, err := json.Marshal(entry)
		if err != nil {
			log.Printf("Failed to marshal log entry: %v", err)
			return
		}
		
		fmt.Fprintln(os.Stdout, string(jsonData))
	}
}

func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		DEBUG: 0,
		INFO:  1,
		WARN:  2,
		ERROR: 3,
		FATAL: 4,
	}
	return levels[level] >= levels[l.level]
}

func (l *Logger) mergeFields(fields map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	
	// Add logger fields first
	for k, v := range l.fields {
		merged[k] = v
	}
	
	// Add entry fields (override logger fields if same key)
	for k, v := range fields {
		merged[k] = v
	}
	
	return merged
}

func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	l.log(DEBUG, message, l.mergeFieldsList(fields...))
}

func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	l.log(INFO, message, l.mergeFieldsList(fields...))
}

func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	l.log(WARN, message, l.mergeFieldsList(fields...))
}

func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	l.log(ERROR, message, l.mergeFieldsList(fields...))
}

func (l *Logger) Fatal(message string, fields ...map[string]interface{}) {
	l.log(FATAL, message, l.mergeFieldsList(fields...))
	os.Exit(1)
}

func (l *Logger) mergeFieldsList(fieldsList ...map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for _, fields := range fieldsList {
		for k, v := range fields {
			merged[k] = v
		}
	}
	return merged
}

// HTTP request logging
func (l *Logger) LogHTTPRequest(method, path, userAgent string, statusCode int, duration time.Duration, requestID string) {
	l.WithRequestID(requestID).Info("HTTP Request", map[string]interface{}{
		"method":      method,
		"path":        path,
		"user_agent":  userAgent,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	})
}

// Database operation logging
func (l *Logger) LogDatabaseOperation(operation, table string, duration time.Duration, err error) {
	fields := map[string]interface{}{
		"operation":   operation,
		"table":       table,
		"duration_ms": duration.Milliseconds(),
	}
	
	if err != nil {
		l.Error("Database operation failed", fields)
	} else {
		l.Debug("Database operation completed", fields)
	}
}

// Authentication logging
func (l *Logger) LogAuthEvent(event, username string, userID uint, success bool, ip string) {
	fields := map[string]interface{}{
		"event":    event,
		"username": username,
		"user_id":  userID,
		"success":  success,
		"ip":       ip,
	}
	
	if success {
		l.Info("Authentication event", fields)
	} else {
		l.Warn("Authentication failed", fields)
	}
}

// Global logger instance
var GlobalLogger *Logger

func InitGlobalLogger(level LogLevel, service string) {
	GlobalLogger = NewLogger(level, service)
}

func GetLogger() *Logger {
	if GlobalLogger == nil {
		GlobalLogger = NewLogger(INFO, "go-crud")
	}
	return GlobalLogger
}
