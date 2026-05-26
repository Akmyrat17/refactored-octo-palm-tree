package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)

	DBDebug(ctx context.Context, msg string, keysAndValues ...any)
	DBInfo(ctx context.Context, msg string, keysAndValues ...any)
	DBWarn(ctx context.Context, msg string, keysAndValues ...any)
	DBError(ctx context.Context, msg string, keysAndValues ...any)
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Time    string                 `json:"time"`
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
}

// ConsoleLogger is a JSON logger that writes structured logs to stdout
type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Debug(msg string, keysAndValues ...any) {
	l.logWithLevel("DEBUG", msg, keysAndValues...)
}

func (l *ConsoleLogger) Info(msg string, keysAndValues ...any) {
	l.logWithLevel("INFO", msg, keysAndValues...)
}

func (l *ConsoleLogger) Warn(msg string, keysAndValues ...any) {
	l.logWithLevel("WARN", msg, keysAndValues...)
}

func (l *ConsoleLogger) Error(msg string, keysAndValues ...any) {
	l.logWithLevel("ERROR", msg, keysAndValues...)
}

func (l *ConsoleLogger) logWithLevel(level, msg string, keysAndValues ...any) {
	fields := make(map[string]interface{})
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key := fmt.Sprintf("%v", keysAndValues[i])
			value := keysAndValues[i+1]
			// Convert UUID and other types to string for better readability
			fields[key] = value
		}
	}

	entry := LogEntry{
		Time:    time.Now().Format(time.RFC3339Nano),
		Level:   level,
		Message: msg,
		Fields:  fields,
	}

	jsonBytes, _ := json.Marshal(entry)
	fmt.Fprintln(os.Stdout, string(jsonBytes))
}

func (l *ConsoleLogger) DBDebug(ctx context.Context, msg string, keysAndValues ...any) {
	l.Debug(msg, keysAndValues...)
}

func (l *ConsoleLogger) DBInfo(ctx context.Context, msg string, keysAndValues ...any) {
	l.Info(msg, keysAndValues...)
}

func (l *ConsoleLogger) DBWarn(ctx context.Context, msg string, keysAndValues ...any) {
	l.Warn(msg, keysAndValues...)
}

func (l *ConsoleLogger) DBError(ctx context.Context, msg string, keysAndValues ...any) {
	l.Error(msg, keysAndValues...)
}
