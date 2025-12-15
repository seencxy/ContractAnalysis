package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with additional functionality
type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
}

// Config represents logger configuration
type Config struct {
	Level      string   // debug, info, warn, error
	Format     string   // json, console
	Output     []string // stdout, stderr, file
	FilePath   string
	MaxSize    int  // megabytes
	MaxBackups int  // number of backups
	MaxAge     int  // days
	Compress   bool
}

var globalLogger *Logger

// New creates a new logger instance
func New(cfg Config) (*Logger, error) {
	// Parse log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create encoder based on format
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create cores for different outputs
	var cores []zapcore.Core

	for _, output := range cfg.Output {
		switch output {
		case "stdout":
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
		case "stderr":
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), level))
		case "file":
			if cfg.FilePath != "" {
				// Ensure log directory exists
				logDir := filepath.Dir(cfg.FilePath)
				if err := os.MkdirAll(logDir, 0755); err != nil {
					return nil, fmt.Errorf("failed to create log directory: %w", err)
				}

				// Open log file
				logFile, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return nil, fmt.Errorf("failed to open log file: %w", err)
				}

				cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(logFile), level))
			}
		}
	}

	// Combine cores
	core := zapcore.NewTee(cores...)

	// Create logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	logger := &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}

	return logger, nil
}

// SetGlobal sets the global logger
func SetGlobal(logger *Logger) {
	globalLogger = logger
}

// GetGlobal returns the global logger
func GetGlobal() *Logger {
	if globalLogger == nil {
		// Create a default logger if none exists
		logger, _ := New(Config{
			Level:  "info",
			Format: "console",
			Output: []string{"stdout"},
		})
		globalLogger = logger
	}
	return globalLogger
}

// Sugar returns the sugared logger for easier logging
func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

// WithFields returns a new logger with additional fields
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		sugar:  l.Logger.With(fields...).Sugar(),
	}
}

// WithError adds an error field to the logger
func (l *Logger) WithError(err error) *Logger {
	return l.WithFields(zap.Error(err))
}

// WithComponent adds a component field to the logger
func (l *Logger) WithComponent(component string) *Logger {
	return l.WithFields(zap.String("component", component))
}

// WithRequestID adds a request ID field to the logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithFields(zap.String("request_id", requestID))
}

// WithSymbol adds a symbol field to the logger
func (l *Logger) WithSymbol(symbol string) *Logger {
	return l.WithFields(zap.String("symbol", symbol))
}

// WithSignalID adds a signal ID field to the logger
func (l *Logger) WithSignalID(signalID string) *Logger {
	return l.WithFields(zap.String("signal_id", signalID))
}

// WithStrategy adds a strategy field to the logger
func (l *Logger) WithStrategy(strategy string) *Logger {
	return l.WithFields(zap.String("strategy", strategy))
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Global logging functions

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	GetGlobal().Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	GetGlobal().Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	GetGlobal().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	GetGlobal().Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	GetGlobal().Fatal(msg, fields...)
}

// Debugf logs a debug message with formatting
func Debugf(template string, args ...interface{}) {
	GetGlobal().Sugar().Debugf(template, args...)
}

// Infof logs an info message with formatting
func Infof(template string, args ...interface{}) {
	GetGlobal().Sugar().Infof(template, args...)
}

// Warnf logs a warning message with formatting
func Warnf(template string, args ...interface{}) {
	GetGlobal().Sugar().Warnf(template, args...)
}

// Errorf logs an error message with formatting
func Errorf(template string, args ...interface{}) {
	GetGlobal().Sugar().Errorf(template, args...)
}

// Fatalf logs a fatal message with formatting and exits
func Fatalf(template string, args ...interface{}) {
	GetGlobal().Sugar().Fatalf(template, args...)
}

// WithComponent returns a logger with component field
func WithComponent(component string) *Logger {
	return GetGlobal().WithComponent(component)
}

// WithError returns a logger with error field
func WithError(err error) *Logger {
	return GetGlobal().WithError(err)
}

// WithSymbol returns a logger with symbol field
func WithSymbol(symbol string) *Logger {
	return GetGlobal().WithSymbol(symbol)
}

// WithSignalID returns a logger with signal ID field
func WithSignalID(signalID string) *Logger {
	return GetGlobal().WithSignalID(signalID)
}

// WithStrategy returns a logger with strategy field
func WithStrategy(strategy string) *Logger {
	return GetGlobal().WithStrategy(strategy)
}
