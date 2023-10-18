package bonus

import (
	"context"
	"errors"

	"github.com/desutedja/lreport/internal/repository/model"
)

type bonusStore interface {
	CreateBonus(ctx context.Context, req model.DataBonus) error
	IsBonusExist(ctx context.Context, date string) (bool, error)
	GetBonus(ctx context.Context, req model.BasicRequest) (data []model.DataBonus, err error)
}

type Service struct {
	bonusStore bonusStore
}

func NewService(bonusStore bonusStore) *Service {
	return &Service{
		bonusStore: bonusStore,
	}
}

func (s *Service) CreateBonus(ctx context.Context, userId string, req model.ReqBonus) error {
	// check is transaction on the given date is already exist
	isExist, err := s.bonusStore.IsBonusExist(ctx, req.TransDate)
	if err != nil {
		return err
	}

	if isExist {
		return errors.New("data on this date: " + req.TransDate + " is already exist")
	}

	// calculate data
	data := model.DataBonus{}
	data.UserId = userId
	data.ReqBonus = req
	data.Total = req.NewMember + req.CbSl + req.RbSl + req.CbCa + req.RollCa + req.CbSp + req.RbSp + req.Refferal + req.Promo

	// insert data to db
	if err := s.bonusStore.CreateBonus(ctx, data); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetBonus(ctx context.Context, req model.BasicRequest) (data []model.DataBonus, err error) {
	// get category from db
	data, err = s.bonusStore.GetBonus(ctx, req)
	if err != nil {
		return data, errors.New("data not found")
	}

	return data, nil
}
