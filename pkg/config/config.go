package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

// GetFlagOrEnvVar retrieves a value from either a provided flag value or an environment variable.
// It first checks if the flag value is non-empty, then falls back to the environment variable.
// Returns an error if neither source provides a value.
func GetFlagOrEnvVar(val string, key string) (string, error) {
	logger := log.Trace().Str("key", key)

	if val != "" {
		logger.Msg("using flag value")
		return val, nil
	}

	result := os.Getenv(key)
	if result != "" {
		logger.Msg("using environment variable")
		return result, nil
	}

	return "", fmt.Errorf("%s not provided and %s environment variable not set", key, key)
}
