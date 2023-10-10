package transaction

import (
	"context"
	"errors"

	"github.com/desutedja/lreport/internal/repository/model"
)

type transactionStore interface {
	CreateTransaction(ctx context.Context, req model.DataTransaction) error
	IsTransactionExist(ctx context.Context, date string) (bool, error)
	GetTransaction(ctx context.Context, req model.BasicRequest) (data []model.DataTransaction, err error)
}

type Service struct {
	transactionStore transactionStore
}

func NewService(transactionStore transactionStore) *Service {
	return &Service{
		transactionStore: transactionStore,
	}
}

func (s *Service) CreateTransaction(ctx context.Context, userId string, req model.ReqTransaction) error {
	// check is transaction on the given date is already exist
	isExist, err := s.transactionStore.IsTransactionExist(ctx, req.TransDate)
	if err != nil {
		return err
	}

	if isExist {
		return errors.New("data on this date: " + req.TransDate + " is already exist")
	}

	// calculate data
	data := model.DataTransaction{}
	data.ReqTransaction = req
	data.UserId = userId
	data.ConvDp = float64(req.Regis) / float64(req.RegisDp)
	data.ConvTr = float64(req.TransWd) / float64(req.TransDp)
	data.SubTotal = req.TotalDp - req.TotalWd
	data.Ats = req.Wl * 0.2
	data.Total = req.Wl - data.Ats // TODO: masih harus dikurang sama bonus, tapi bonus belum tau dapet drmn

	// insert data to db
	if err := s.transactionStore.CreateTransaction(ctx, data); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetTransaction(ctx context.Context, req model.BasicRequest) (data []model.DataTransaction, err error) {
	// get category from db
	data, err = s.transactionStore.GetTransaction(ctx, req)
	if err != nil {
		return data, errors.New("data not found")
	}

	return data, nil
}
