package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadDatabaseConfig(t *testing.T) {
	t.Parallel()

	os.Setenv("DATABASE_HOST", "localhost")
	os.Setenv("DATABASE_PORT", "1234")
	os.Setenv("DATABASE_DBNAME", "test")
	os.Setenv("DATABASE_USERNAME", "test")
	os.Setenv("DATABASE_PASSWORD", "test")
	os.Setenv("DATABASE_MAX_OPEN_CONNECTION", "1")
	os.Setenv("DATABASE_MAX_IDLE_CONNECTION", "2")
	os.Setenv("DATABASE_CONNECTION_MAX_IDLE_TIME_IN_SECOND", "3")
	os.Setenv("DATABASE_CONNECTION_MAX_LIFE_TIME_IN_SECOND", "4")

	expectedCfg := &DatabaseConfig{
		Host:                  "localhost",
		Port:                  1234,
		DBName:                "test",
		Username:              "test",
		Password:              "test",
		MaxOpenConnection:     1,
		MaxIdleConnection:     2,
		ConnectionMaxIdleTime: 3 * time.Second,
		ConnectionMaxLifetime: 4 * time.Second,
	}

	assert.Equal(t, expectedCfg, loadDatabaseConfig())
}
