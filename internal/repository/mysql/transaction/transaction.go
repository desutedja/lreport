package transaction

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/desutedja/lreport/internal/repository/model"
	"github.com/google/uuid"
)

type TransactionStore struct {
	db *sql.DB
}

func NewTransactionStore(db *sql.DB) *TransactionStore {
	return &TransactionStore{
		db: db,
	}
}

func (s *TransactionStore) CreateTransaction(ctx context.Context, req model.DataTransaction) error {
	query := `
		INSERT INTO transaction (id,user_id, category_id, regis, regis_dp, active_player, 
			conv_dp, trans_dp, trans_wd, conv_tr, total_dp, total_wd, sub_total, wl, ats, total, trans_date)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`

	trx, err := s.db.BeginTx(ctx, nil)
	defer trx.Rollback()
	if err != nil {
		return err
	}

	id := uuid.New()
	_, err = trx.Exec(query, id, req.UserId, req.CategoryId, req.Regis, req.RegisDp, req.ActivePlayer,
		req.ConvDp, req.TransDp, req.TransWd, req.ConvTr, req.TotalDp, req.TotalWd, req.SubTotal, req.Wl, req.Ats, req.Total, req.TransDate)
	if err != nil {
		return err
	}

	trx.Commit()
	return nil
}

// date must be on this format (yyyy-mm-dd)
func (s *TransactionStore) IsTransactionExist(ctx context.Context, date string) (bool, error) {
	data := struct {
		Id string
	}{}

	query := `
		SELECT id FROM transaction WHERE DATE_FORMAT(trans_date, '%Y-%m-%d') = ?
	`
	if err := s.db.QueryRowContext(ctx, query, date).Scan(&data.Id); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return false, nil
		}

		return true, err
	}

	return true, nil
}

func (s *TransactionStore) GetTransaction(ctx context.Context, req model.BasicRequest) (data []model.DataTransaction, err error) {
	offset := (req.Page * req.Limit) - req.Limit

	query := `
		SELECT 
			id, user_id, category_id, regis, regis_dp, active_player, 
			conv_dp, trans_dp, trans_wd, conv_tr, total_dp, total_wd, sub_total, wl, ats, total, trans_date
		FROM transaction
	`

	query = query + " ORDER BY created_on DESC LIMIT ? OFFSET ?"

	rows, err := s.db.QueryContext(ctx, query, req.Limit, offset)
	if err != nil {
		log.Println("error query: ", err)
		return data, err
	}

	for rows.Next() {

		dt := model.DataTransaction{}

		if err := rows.Scan(
			&dt.Id, &dt.UserId, &dt.CategoryId, &dt.Regis, &dt.RegisDp, &dt.ActivePlayer,
			&dt.ConvDp, &dt.TransDp, &dt.TransWd, &dt.ConvTr, &dt.TotalDp, &dt.TotalWd, &dt.SubTotal, &dt.Wl, &dt.Ats, &dt.Total, &dt.TransDate,
		); err != nil {
			log.Println("error scan: ", err)
			return data, err
		}

		data = append(data, dt)
	}

	return data, nil
}

func (s *TransactionStore) GetTransactionStatistic(ctx context.Context, categoryId int) (data model.RespReportTransaction, err error) {
	query := `
		SELECT 
			regis, regis_dp, active_player, trans_dp, trans_wd, total_dp,
			total_wd, wl, conv_dp, conv_tr,  sub_total, ats,
			IFNULL(bn.total,0) bonus, (tr.total - IFNULL(bn.total,0)) total,
			DAY(tr.trans_date) day, MONTH(tr.trans_date) month, YEAR(tr.trans_date) year,
			DATE_FORMAT(tr.trans_date, "%Y-%m-%d") period
		FROM transaction tr
		LEFT JOIN bonus bn ON tr.category_id = bn.category_id AND tr.trans_date = bn.trans_date
		WHERE tr.category_id = ?
	`

	query = query + " ORDER BY tr.created_on DESC LIMIT 30"

	rows, err := s.db.QueryContext(ctx, query, categoryId)
	if err != nil {
		log.Println("error query: ", err)
		return data, err
	}

	for rows.Next() {

		dt := model.DataReportTransaction{}

		if err := rows.Scan(
			&dt.Regis, &dt.RegisDp, &dt.ActivePlayer, &dt.TransDp, &dt.TransWd,
			&dt.TotalDp, &dt.TotalWd, &dt.Wl, &dt.ConvDp, &dt.ConvTr, &dt.SubTotal, &dt.Ats,
			&dt.Bonus, &dt.Total, &dt.Day, &dt.Month, &dt.Year, &dt.Period,
		); err != nil {
			log.Println("error scan: ", err)
			return data, err
		}

		data.DataReport = append(data.DataReport, dt)
	}

	dtkey := "regis,regis_dp,active_player,trans_dp,trans_wd,total_dp,total_wd,wl,conv_dp,conv_tr,sub_total,ats,bonus,total"
	data.DataKey = strings.Split(dtkey, ",")

	return data, nil
}
