package logging

import (
	"log"
	"os"

	"go.uber.org/zap"
)

// MustZapLogger creates a new zap logger, of fails and log the error using the
// standard library's logger, with "fatal" level.
func MustZapLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.New(os.Stderr, "", 0).Fatalf("failed to initialise zap logger: %v", err)
	}
	return logger
}
