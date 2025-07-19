package logger

import (
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger initializes the global logger
func InitLogger(environment string) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	switch environment {
	case "production", "staging":
		logger.SetLevel(logrus.InfoLevel)
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	default: // development
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	logger.SetOutput(os.Stdout)

	// Set global logger
	Logger = logger
	return logger
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	if Logger == nil {
		InitLogger("development")
	}
	return Logger
}

// WithContext creates a logger with request context
func WithContext(c *gin.Context) *logrus.Entry {
	logger := GetLogger()
	
	entry := logger.WithFields(logrus.Fields{
		"request_id": c.GetString("X-Request-ID"),
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"client_ip":  c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
	})

	// Add user ID if available
	if userID, exists := c.Get("user_id"); exists {
		entry = entry.WithField("user_id", userID)
	}

	// Add tenant ID if available
	if tenantID, exists := c.Get("tenant_id"); exists {
		entry = entry.WithField("tenant_id", tenantID)
	}

	return entry
}

// WithFields creates a logger with custom fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// Debug logs a debug message
func Debug(msg string) {
	GetLogger().Debug(msg)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Info logs an info message
func Info(msg string) {
	GetLogger().Info(msg)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warn logs a warning message
func Warn(msg string) {
	GetLogger().Warn(msg)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Error logs an error message
func Error(msg string) {
	GetLogger().Error(msg)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string) {
	GetLogger().Fatal(msg)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// SetOutput sets the logger output
func SetOutput(output io.Writer) {
	GetLogger().SetOutput(output)
}

// SetLevel sets the logger level
func SetLevel(level logrus.Level) {
	GetLogger().SetLevel(level)
} 