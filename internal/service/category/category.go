package category

import (
	"context"
	"errors"
	"log"

	"github.com/desutedja/lreport/internal/repository/model"
)

type categoryStore interface {
	InsertCategory(ctx context.Context, dt model.ReqCategory) error
	GetCategory(ctx context.Context) (data []model.CategoryData, err error)
}

type Service struct {
	categoryStore categoryStore
}

func NewService(categoryStore categoryStore) *Service {
	return &Service{
		categoryStore: categoryStore,
	}
}

func (s *Service) GetCategory(ctx context.Context) (data []model.CategoryData, err error) {
	// get category from db
	data, err = s.categoryStore.GetCategory(ctx)
	if err != nil {
		return data, errors.New("data not found")
	}

	return data, nil
}

func (s *Service) CreateCategory(ctx context.Context, req model.ReqCategory) (data []model.CategoryData, err error) {
	// insert category to db
	if err = s.categoryStore.InsertCategory(ctx, req); err != nil {
		log.Println("create category failed: ", err)
		return data, errors.New("create category failed")
	}

	return data, nil
}
