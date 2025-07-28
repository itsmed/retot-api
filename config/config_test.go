package config_test

import (
	"app/config"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConfig_ReturnsSetEnvValue(t *testing.T) {
	os.Setenv("SECRET", "testsecret")

	// First access will trigger .env load (noop in this case)
	secret := config.Config("SECRET")

	assert.Equal(t, "testsecret", secret)
}

func TestConfig_EmptyForUnsetEnv(t *testing.T) {
	os.Unsetenv("UNDEFINED_VAR")

	val := config.Config("UNDEFINED_VAR")

	assert.Equal(t, "", val)
}
