package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Logger wraps the logrus Logger
type Logger struct {
	*log.Logger
}

// NewLogger initializes and returns a new Logger with the provided log level.
func NewLogger(level log.Level) *Logger {
	logger := log.New()
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&log.JSONFormatter{})

	return &Logger{logger}
}

// HandleError logs the error based on the provided log level: "I" for Info, "E" for Error, and "W" for Warning.
// If the error is not nil, it logs the error and the associated message at the given log level.
func (l *Logger) HandleError(logLevel, msg string, err error) {
	if err != nil {
		// Create an entry for the error
		entry := l.WithFields(log.Fields{
			"error": err,
		})

		// Log the error at the appropriate level
		switch logLevel {
		case "I":
			entry.Info(msg)
		case "E":
			entry.Error(msg)
		case "W":
			entry.Warn(msg)
		default:
			entry.Info(msg)
		}
	}
}
