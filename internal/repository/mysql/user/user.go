package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/desutedja/lreport/internal/repository/model"
	"github.com/google/uuid"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) Login(ctx context.Context, username string) (data model.UserData, err error) {
	query := `
		SELECT
			id, username, password, user_level
		FROM users
		WHERE username = ?
			AND deleted = 0
	`

	if err := s.db.QueryRowContext(ctx, query, username).Scan(
		&data.Id,
		&data.Username,
		&data.Password,
		&data.UserLevel,
	); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return data, errors.New("record not found")
		}
		return data, err
	}

	return
}

func (s *UserStore) InsertLoginHistory(ctx context.Context, userID, device, ipAddress string) error {
	query := `
		INSERT INTO login_history (id, user_id, device, ip_address)
		VALUES (?,?,?,?)
	`

	trx, err := s.db.BeginTx(ctx, nil)
	defer trx.Rollback()
	if err != nil {
		return err
	}

	historyID := uuid.New()
	if err != nil {
		return err
	}

	_, err = trx.Exec(query, historyID, userID, device, ipAddress)
	if err != nil {
		return err
	}

	trx.Commit()
	return nil
}

func (s *UserStore) ChangePassword(ctx context.Context, userID, newPassword string) error {
	query := `
		UPDATE users
		SET password = ?
		WHERE username= ?
	`

	trx, err := s.db.BeginTx(ctx, nil)
	defer trx.Rollback()
	if err != nil {
		return err
	}

	dt, err := trx.Exec(query, newPassword, userID)
	if err != nil {
		return err
	}

	row, err := dt.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return errors.New("now row affected")
	}

	trx.Commit()
	return nil
}

func (s *UserStore) CreateUser(ctx context.Context, username, password, userLevel string) (uuid.UUID, error) {
	query := `
		INSERT INTO users (id, username, password, user_level)
		VALUES (?,?,?,?)
	`

	trx, err := s.db.BeginTx(ctx, nil)
	defer trx.Rollback()
	if err != nil {
		return uuid.Nil, err
	}

	userId := uuid.New()
	if err != nil {
		return uuid.Nil, err
	}

	_, err = trx.Exec(query, userId, username, password, userLevel)
	if err != nil {
		return uuid.Nil, err
	}

	trx.Commit()
	return userId, nil
}

func (s *UserStore) LoginHistory(ctx context.Context, req model.BasicRequest) (data []model.LoginHistory, err error) {
	offset := (req.Page * req.Limit) - req.Limit

	query := `
		SELECT user_id,device,ip_address,created_on
		FROM login_history
	`

	if req.Search != "" {
		query = query + fmt.Sprintf(" WHERE user_id like %%%s%% ", req.Search)
	}

	query = query + " ORDER BY created_on DESC LIMIT ? OFFSET ?"

	log.Println(query, "\nOFFSET: ", offset, "\nlimit", req.Limit)

	rows, err := s.db.QueryContext(ctx, query, req.Limit, offset)
	if err != nil {
		log.Println("error query: ", err)
		return data, err
	}

	for rows.Next() {

		dt := model.LoginHistory{}

		if err := rows.Scan(
			&dt.UserId,
			&dt.Device,
			&dt.IpAddress,
			&dt.CreatedOn,
		); err != nil {
			log.Println("error scan: ", err)
			return data, err
		}

		data = append(data, dt)
	}

	return data, nil

}

// func (s *UserStore) GetUUID(trx *sql.Tx, tablename string) (idNew uuid.UUID, err error) {
// 	dt := struct {
// 		ttl int
// 	}{}

// 	for {
// 		id := uuid.New()

// 		query := fmt.Sprintf("SELECT COUNT(*) ttl FROM %s WHERE id=?", tablename)
// 		if err := trx.QueryRow(query, id).Scan(
// 			&dt.ttl,
// 		); err != nil {
// 			log.Println("error query get UUID: ", err)
// 			break
// 		}

// 		if dt.ttl == 0 {
// 			idNew = id
// 			err = nil
// 			break
// 		}
// 	}

// 	return
// }
