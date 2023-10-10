package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Setenv("ENV", "test")
	if err := Load(); err != nil {
		t.Fatal(err)
	}

	assert.IsType(t, &config{}, appConfig)
}
