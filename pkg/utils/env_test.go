package utils_test

import (
	"github.com/jyotishp/kafka-test/pkg/utils"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	envKey := "TEST_KEY"
	actualVal := "value_something"
	defaultVal := "default_value"
	os.Setenv(envKey, actualVal)

	value := utils.GetEnv(envKey, defaultVal)
	if value != actualVal {
		t.Errorf("env values do not match")
	}

	value = utils.GetEnv("SOMEKEY", defaultVal)
	if value != defaultVal {
		t.Errorf("default env value does not match")
	}
}
