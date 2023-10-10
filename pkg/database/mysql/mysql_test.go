package mysql

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/desutedja/lreport/pkg/database"
)

// to run the test. specify the DATABASE_INFO with this format:
// host={db_host};port={db_port};dbname={db_name};username={db_username};password={db_password}
func extractConnection(t *testing.T) *database.Config {
	dbInfo := os.Getenv("DATABASE_INFO")
	if dbInfo == "" {
		t.Skip("no db connection info found")
	}

	info := strings.Split(dbInfo, ";")
	if len(info) < 3 {
		t.Skip("info isn't in correct format")
	}

	host := info[0]
	portString := info[1]
	dbname := info[2]
	username := info[3]
	password := info[4]

	port, err := strconv.Atoi(strings.Split(portString, "=")[1])
	if err != nil {
		t.Skip(err)
	}

	cfg := &database.Config{
		Host:     strings.Split(host, "=")[1],
		Port:     port,
		DBName:   strings.Split(dbname, "=")[1],
		User:     strings.Split(username, "=")[1],
		Password: strings.Split(password, "=")[1],
	}
	return cfg
}

func TestConnect_Success(t *testing.T) {
	cfg := extractConnection(t)

	tests := []struct {
		name            string
		maxIdleConn     int
		maxOpenconn     int
		connMaxIdleTime time.Duration
		connMaxLifetime time.Duration
	}{
		{
			name: "Success",
		},
		{
			name:        "SuccessWithMaxIdleConn",
			maxIdleConn: 1,
		},
		{
			name:        "SuccessWithMaxOpenConn",
			maxOpenconn: 1,
		},
		{
			name:            "SuccessWithConnMaxIdleTime",
			connMaxIdleTime: 1 * time.Second,
		},
		{
			name:            "SuccessWithConnMaxLifetime",
			connMaxLifetime: 1 * time.Second,
		},
		{
			name:            "SuccessWithAllAtributSet",
			maxIdleConn:     1,
			maxOpenconn:     1,
			connMaxIdleTime: 1 * time.Second,
			connMaxLifetime: 1 * time.Second,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			conn, err := Connect(context.Background(), cfg)
			if err != nil {
				t.Fatal(err)
			}

			if conn == nil {
				t.Fatal("connection is nil")
			}
		})
	}
}

func TestConnect_Failed(t *testing.T) {
	tests := []struct {
		name   string
		config *database.Config
	}{
		{
			name:   "ConfigNil",
			config: nil,
		},
		{
			name: "EmptyHost",
			config: &database.Config{
				Host: "",
			},
		},
		{
			name: "EmptyUsername",
			config: &database.Config{
				Host: "localhost",
				User: "",
			},
		},
		{
			name: "EmptyPassword",
			config: &database.Config{
				Host:     "localhost",
				User:     "root",
				Password: "",
			},
		},
		{
			name: "EmptyDBName",
			config: &database.Config{
				Host:     "localhost",
				User:     "root",
				Password: "test",
				DBName:   "",
			},
		},
		{
			name: "EmptyPort",
			config: &database.Config{
				Host:     "localhost",
				User:     "root",
				Password: "test",
				DBName:   "test",
				Port:     0,
			},
		},
		{
			name: "OpenConnection",
			config: &database.Config{
				Host:     "localhost",
				User:     "root",
				Password: "test",
				DBName:   "test",
				Port:     1234,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if _, err := Connect(context.Background(), tt.config); err == nil {
				t.Fatal("should return error")
			}
		})
	}
}
