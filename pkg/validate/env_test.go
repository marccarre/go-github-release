package validate_test

import (
	"os"
	"testing"

	"github.com/marccarre/go-github-release/pkg/validate"
	"github.com/stretchr/testify/assert"
)

func TestEnvNonExistingEnvVar(t *testing.T) {
	assert.EqualError(t, validate.Env("NOT_EXISTING_ENV_VAR"), "environment variable NOT_EXISTING_ENV_VAR is unset")
}

func TestEnvEmptyEnvVar(t *testing.T) {
	assert.NoError(t, os.Setenv("EMPTY_ENV_VAR", ""))
	defer os.Unsetenv("EMPTY_ENV_VAR")
	assert.EqualError(t, validate.Env("EMPTY_ENV_VAR"), "environment variable EMPTY_ENV_VAR is empty")
}

func TestEnvNonEmptyEnvVar(t *testing.T) {
	assert.NoError(t, os.Setenv("NON_EMPTY_ENV_VAR", "success!"))
	defer os.Unsetenv("NON_EMPTY_ENV_VAR")
	assert.NoError(t, validate.Env("NON_EMPTY_ENV_VAR"))
}
