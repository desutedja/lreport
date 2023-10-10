package config

import "time"

type DatabaseConfig struct {
	Host     string
	Port     int
	DBName   string
	Username string
	Password string

	MaxOpenConnection     int
	MaxIdleConnection     int
	ConnectionMaxIdleTime time.Duration
	ConnectionMaxLifetime time.Duration
}

func loadDatabaseConfig() *DatabaseConfig {
	connMaxIdleTimeInSecond := getIntWithDefault("DATABASE_CONNECTION_MAX_IDLE_TIME_IN_SECOND", 0)
	connMaxLifeTimeInSecond := getIntWithDefault("DATABASE_CONNECTION_MAX_LIFE_TIME_IN_SECOND", 0)
	return &DatabaseConfig{
		Host:     getStringOrPanic("DATABASE_HOST"),
		Port:     getIntOrPanic("DATABASE_PORT"),
		DBName:   getStringOrPanic("DATABASE_DBNAME"),
		Username: getStringOrPanic("DATABASE_USERNAME"),
		Password: getStringOrPanic("DATABASE_PASSWORD"),

		MaxOpenConnection:     getIntWithDefault("DATABASE_MAX_OPEN_CONNECTION", 0),
		MaxIdleConnection:     getIntWithDefault("DATABASE_MAX_IDLE_CONNECTION", 0),
		ConnectionMaxIdleTime: time.Duration(connMaxIdleTimeInSecond) * time.Second,
		ConnectionMaxLifetime: time.Duration(connMaxLifeTimeInSecond) * time.Second,
	}
}
