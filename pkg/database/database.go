package database

import (
	"errors"
	"time"
)

const (
	CreateConnectionDefaultTimeout = 10 * time.Second
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string

	MaxOpenConnection     int
	MaxIdleConnection     int
	ConnectionMaxIdleTime time.Duration
	ConnectionMaxLifetime time.Duration
}

func (c *Config) Validate() error {
	switch "" {
	case c.Host:
		return errors.New("empty host")
	case c.User:
		return errors.New("empty username")
	case c.Password:
		return errors.New("empty password")
	case c.DBName:
		return errors.New("empty dbname")
	}

	switch 0 {
	case c.Port:
		return errors.New("empty port")
	}

	return nil
}
