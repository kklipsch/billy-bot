package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

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
