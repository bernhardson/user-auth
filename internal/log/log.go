package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// NewDefaultLogger creates a logger with default settings
func NewDefaultLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// NewCustomLogger creates a logger with a custom output and level
func NewCustomLogger(output io.Writer, level string) (*zerolog.Logger, error) {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	logger := zerolog.New(output).Level(logLevel).With().Timestamp().Caller().Logger()
	return &logger, nil
}

// NewLoggerWithFields creates a logger with additional fields
func NewLoggerWithFields(fields map[string]interface{}) zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp()
	for key, value := range fields {
		logger = logger.Interface(key, value)
	}
	return logger.Logger()
}
