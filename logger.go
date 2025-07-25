package ermeslog

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	entry *logrus.Entry
}

// New creates a new logger instance
func New(config Config) *Logger {
	log := logrus.New()

	log.SetLevel(config.Level)

	// Configure output
	if config.Output != nil {
		log.SetOutput(config.Output)
	} else if config.FileOutput != nil {
		log.SetOutput(&lumberjack.Logger{
			Filename:   config.FileOutput.Filename,
			MaxSize:    config.FileOutput.MaxSize,
			MaxBackups: config.FileOutput.MaxBackups,
			MaxAge:     config.FileOutput.MaxAge,
			Compress:   config.FileOutput.Compress,
		})
	}

	// Configure formatter
	if config.Environment == "development" {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	} else {
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	// Add hooks
	for _, hook := range config.Hooks {
		log.AddHook(hook)
	}

	hostname, _ := os.Hostname()

	entry := log.WithFields(logrus.Fields{
		"business":    config.Business,
		"service":     config.Service,
		"version":     config.Version,
		"environment": config.Environment,
		"hostname":    hostname,
	})

	return &Logger{entry: entry}
}

// NewDefault creates a logger with default configuration
func NewDefault(business, service string) *Logger {
	config := DefaultConfig()
	config.Business = business
	config.Service = service
	return New(config)
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	entry := l.entry

	if reqID := ctx.Value("requestID"); reqID != nil {
		entry = entry.WithField("request_id", reqID)
	}
	if userID := ctx.Value("userID"); userID != nil {
		entry = entry.WithField("user_id", userID)
	}
	if traceID := ctx.Value("traceID"); traceID != nil {
		entry = entry.WithField("trace_id", traceID)
	}

	return &Logger{entry: entry}
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return &Logger{entry: l.entry.WithFields(fields)}
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{entry: l.entry.WithField(key, value)}
}

func (l *Logger) Error(err error, msg string) {
	l.entry.WithError(err).Error(msg)
}

func (l *Logger) ErrorWithContext(ctx context.Context, err error, msg string) {
	l.WithContext(ctx).Error(err, msg)
}

func (l *Logger) Info(msg string) {
	l.entry.Info(msg)
}

func (l *Logger) InfoWithContext(ctx context.Context, msg string) {
	l.WithContext(ctx).Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.entry.Warn(msg)
}

func (l *Logger) Debug(msg string) {
	l.entry.Debug(msg)
}

func (l *Logger) Fatal(msg string) {
	l.entry.Fatal(msg)
}

func (l *Logger) Panic(msg string) {
	l.entry.Panic(msg)
}

func (l *Logger) GetEntry() *logrus.Entry {
	return l.entry
}
