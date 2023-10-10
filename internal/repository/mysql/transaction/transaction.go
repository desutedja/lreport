package transaction

import (
	"context"
	"database/sql"
)

type TransactionStore struct {
	db *sql.DB
}

func NewTransactionStore(db *sql.DB) *TransactionStore {
	return &TransactionStore{
		db: db,
	}
}

func (s *TransactionStore) CreateCategory(ctx context.Context, name, description string) error {
	query := `
		INSERT INTO category (name, description)
		VALUES (?,?)
	`

	trx, err := s.db.BeginTx(ctx, nil)
	defer trx.Rollback()
	if err != nil {
		return err
	}

	_, err = trx.Exec(query, name, description)
	if err != nil {
		return err
	}

	trx.Commit()
	return nil
}
