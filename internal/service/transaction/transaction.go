package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/desutedja/lreport/internal/repository/model"
)

type transactionStore interface {
	CreateTransaction(ctx context.Context, req model.DataTransaction) error
	IsTransactionExist(ctx context.Context, date string) (bool, error)
	GetTransaction(ctx context.Context, req model.BasicRequest) (data []model.DataTransaction, err error)
	GetTransactionStatistic(ctx context.Context, categoryId, year, month int) (data model.RespReportTransaction, err error)
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

func (s *Service) GetTransactionStatistic(ctx context.Context, categoryId, year, month int) (resp model.RespReportTransactionAvg, err error) {
	// get category from db
	data, err := s.transactionStore.GetTransactionStatistic(ctx, categoryId, year, month)
	if err != nil {
		return resp, errors.New("data not found")
	}

	// Initialize the data for averaging
	avgData := []model.AverageDataReportTransaction{}
	divider := 6.0 // Initial divider value

	week := []string{
		"1st", "2nd", "3rd", "4th", "5th",
	}

	// Process data in 5-week sections
	for i := 0; i < 5; i++ {
		// Extract the current 6-week section
		startIndex := i * 6
		endIndex := (i + 1) * 6
		sectionData := data.DataReport[startIndex:endIndex]

		// Initialize the sum variables
		regis, regisDp, ap, trDp, trWd, ttlDp, ttlWd, wl, convDp, convTr, ats, bonus := 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0

		// Calculate the sum for the current section
		for _, dt := range sectionData {
			regis += float64(dt.Regis)
			regisDp += float64(dt.RegisDp)
			ap += float64(dt.ActivePlayer)
			trDp += float64(dt.TransDp)
			trWd += float64(dt.TransWd)
			ttlDp += dt.TotalDp
			ttlWd += dt.TotalWd
			wl += dt.Wl
			convDp += dt.ConvDp
			convTr += dt.ConvTr
			ats += dt.Ats
			bonus += dt.Bonus
		}

		avgData = append(avgData, model.AverageDataReportTransaction{
			Regis:        regis,
			RegisDp:      regisDp,
			ActivePlayer: ap,
			TransDp:      trDp,
			TransWd:      trWd,
			TotalDp:      ttlDp,
			TotalWd:      ttlWd,
			Wl:           wl,
			ConvDp:       convDp,
			ConvTr:       convTr,
			Ats:          ats,
			Bonus:        bonus,
			Period:       fmt.Sprintf("Total %s Week", week[i]),
		})

		// Calculate the averages and append to avgData
		avgData = append(avgData, model.AverageDataReportTransaction{
			Regis:        regis / divider,
			RegisDp:      regisDp / divider,
			ActivePlayer: ap / divider,
			TransDp:      trDp / divider,
			TransWd:      trWd / divider,
			TotalDp:      ttlDp / divider,
			TotalWd:      ttlWd / divider,
			Wl:           wl / divider,
			ConvDp:       convDp / divider,
			ConvTr:       convTr / divider,
			Ats:          ats / divider,
			Bonus:        bonus / divider,
			Period:       fmt.Sprintf("Avg %s Week", week[i]),
		})

		// Update the divider for the next section
		divider = float64(len(sectionData))
	}

	resp.RespReportTransaction = data
	resp.DataAverage = avgData

	return resp, nil
}
