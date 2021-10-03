package env_test

import (
	"github.com/bernardosecades/feeder/pkg/tools/env"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnvOrFallback(t *testing.T) {
	val := env.GetEnvOrFallback("MY_DB", "lololo")
	assert.Equal(t, "lololo", val)

	err := os.Setenv("MY_DB", "bro")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	val = env.GetEnvOrFallback("MY_DB", "lololo")
	assert.Equal(t, "bro", val)
}
