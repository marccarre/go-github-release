package validate

import (
	"fmt"
	"os"
)

// Env validates that the provided environment variable exists.
func Env(name string) error {
	value, exists := os.LookupEnv(name)
	if !exists {
		return fmt.Errorf("environment variable %s is unset", name)
	}
	if value == "" {
		return fmt.Errorf("environment variable %s is empty", name)
	}
	return nil
}
