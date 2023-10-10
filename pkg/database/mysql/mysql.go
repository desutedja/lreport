package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/desutedja/lreport/pkg/database"
	_ "github.com/go-sql-driver/mysql"
)

// Connect creates new connection to postgres db with value from config parameter.
func Connect(ctx context.Context, cfg *database.Config) (*sql.DB, error) {
	dbCtx, cancel := context.WithTimeout(ctx, database.CreateConnectionDefaultTimeout)
	defer cancel()

	if cfg == nil {
		return nil, errors.New("nil config")
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	connSetting := "charset=utf8mb4&parseTime=true&checkConnLiveness=true"
	connString := cfg.User + ":" + cfg.Password + "@(" + cfg.Host + ")/" + cfg.DBName + "?" + connSetting
	// connString := fmt.Sprintf("host=%s port=%d user=%s pasword=%s dbname=%s sslmode=disable",
	// 	cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	// log.Println(connString)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	if cfg.MaxOpenConnection > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConnection)
	}
	if cfg.MaxIdleConnection > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConnection)
	}
	if cfg.ConnectionMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnectionMaxIdleTime)
	}
	if cfg.ConnectionMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	}

	// test connection to DB
	if err := db.PingContext(dbCtx); err != nil {
		return nil, err
	}

	return db, nil
}
