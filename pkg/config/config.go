package config

import (
	"fmt"
	"log"
	"os"
)

func GetFlagOrEnvVar(val string, key string) (string, error) {
	if val != "" {
		log.Printf("using %s from command line flag", key)
		return val, nil
	}

	result := os.Getenv(key)
	if result != "" {
		log.Printf("using %s from environment variable", key)
		return result, nil
	}

	return "", fmt.Errorf("%s not provided and %s environment variable not set", key, key)
}
