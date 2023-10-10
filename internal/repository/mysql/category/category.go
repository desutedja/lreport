package category

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/desutedja/lreport/internal/repository/model"
)

type CategoryStore struct {
	db *sql.DB
}

func NewCategoryStore(db *sql.DB) *CategoryStore {
	return &CategoryStore{
		db: db,
	}
}

func (s *CategoryStore) InsertCategory(ctx context.Context, dt model.ReqCategory) error {
	query := `
		INSERT INTO category (name, description)
		VALUES (?,?)
	`

	trx, err := s.db.BeginTx(ctx, nil)
	defer trx.Rollback()
	if err != nil {
		return err
	}

	_, err = trx.Exec(query, dt.Name, dt.Description)
	if err != nil {
		return err
	}

	trx.Commit()
	return nil
}

func (s *CategoryStore) GetCategory(ctx context.Context) (data []model.CategoryData, err error) {
	query := `
		SELECT
			id, name, description
		FROM category
		WHERE deleted = 0
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return data, errors.New("record not found")
		}
		return data, err
	}

	for rows.Next() {

		dt := model.CategoryData{}

		if err := rows.Scan(
			&dt.Id,
			&dt.Name,
			&dt.Description,
		); err != nil {
			log.Println("error scan: ", err)
			return data, err
		}

		data = append(data, dt)
	}

	return
}
