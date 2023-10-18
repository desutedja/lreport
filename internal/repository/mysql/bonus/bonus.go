package bonus

import (
	"context"
	"database/sql"
	"log"

	"github.com/desutedja/lreport/internal/repository/model"
)

type BonusStore struct {
	db *sql.DB
}

func NewBonusStore(db *sql.DB) *BonusStore {
	return &BonusStore{
		db: db,
	}
}

func (s *BonusStore) CreateBonus(ctx context.Context, req model.DataBonus) error {
	query := `
		INSERT INTO bonus (user_id, category_id, new_member, cb_sl, rb_sl, 
			cb_ca, roll_ca, cb_sp, rb_sp, refferal, promo, total, trans_date)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
	`

	trx, err := s.db.BeginTx(ctx, nil)
	defer trx.Rollback()
	if err != nil {
		return err
	}

	_, err = trx.Exec(query, req.UserId, req.CategoryId, req.NewMember, req.CbSl, req.RbSl,
		req.CbCa, req.RollCa, req.CbSp, req.RbSp, req.Refferal, req.Promo, req.Total, req.TransDate)
	if err != nil {
		return err
	}

	trx.Commit()
	return nil
}

// date must be on this format (yyyy-mm-dd)
func (s *BonusStore) IsBonusExist(ctx context.Context, date string) (bool, error) {
	data := struct {
		Id string
	}{}

	query := `
		SELECT id FROM bonus WHERE DATE_FORMAT(trans_date, '%Y-%m-%d') = ?
	`
	if err := s.db.QueryRowContext(ctx, query, date).Scan(&data.Id); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return false, nil
		}

		return true, err
	}

	return true, nil
}

func (s *BonusStore) GetBonus(ctx context.Context, req model.BasicRequest) (data []model.DataBonus, err error) {
	offset := (req.Page * req.Limit) - req.Limit

	query := `
		SELECT 
			b.id, b.user_id, u.username, b.category_id, c.name, b.new_member, b.cb_sl, b.rb_sl, 
			b.cb_ca, b.roll_ca, b.cb_sp, b.rb_sp, b.refferal, b.promo, b.total, b.trans_date
		FROM bonus b
		INNER JOIN users u ON b.user_id = u.id
		INNER JOIN category c on b.category_id = c.id
	`

	query = query + " ORDER BY b.created_on DESC LIMIT ? OFFSET ?"

	rows, err := s.db.QueryContext(ctx, query, req.Limit, offset)
	if err != nil {
		log.Println("error query: ", err)
		return data, err
	}

	for rows.Next() {

		dt := model.DataBonus{}

		if err := rows.Scan(
			&dt.Id, &dt.UserId, dt.Username, &dt.CategoryId, &dt.CategoryName, &dt.NewMember, &dt.CbSl, &dt.RbSl,
			&dt.CbCa, &dt.RollCa, &dt.CbSp, &dt.RbSp, &dt.Refferal, &dt.Promo, &dt.Total, &dt.TransDate,
		); err != nil {
			log.Println("error scan: ", err)
			return data, err
		}

		data = append(data, dt)
	}

	return data, nil
}
