package category

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/desutedja/lreport/internal/repository/model"
	"github.com/thedevsaddam/renderer"
)

type categoryService interface {
	CreateCategory(ctx context.Context, req model.ReqCategory) (data []model.CategoryData, err error)
	GetCategory(ctx context.Context) (data []model.CategoryData, err error)
}

type Handler struct {
	render          *renderer.Render
	categoryService categoryService
}

func NewHandler(categoryService categoryService) *Handler {
	render := renderer.New()
	return &Handler{
		render:          render,
		categoryService: categoryService,
	}
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body := model.ReqCategory{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.render.JSON(w, http.StatusBadRequest, "bad request")
		return
	}

	// validate request
	if body.Name == "" {
		h.render.JSON(w, http.StatusBadRequest, "name is required")
		return
	}

	_, err = h.categoryService.CreateCategory(ctx, body)
	if err != nil {
		if err.Error() == model.ERROR_USER_EXIST {
			h.render.JSON(w, http.StatusBadRequest, model.ERROR_USER_EXIST)
			return
		}

		h.render.JSON(w, http.StatusInternalServerError, err)
		return
	}

	h.render.JSON(w, http.StatusOK, "Category successfully created")
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categoryList, err := h.categoryService.GetCategory(ctx)
	if err != nil {
		h.render.JSON(w, http.StatusInternalServerError, model.RespBody{
			Message: "failed",
			Data:    err.Error(),
		})
		return
	}

	h.render.JSON(w, http.StatusOK, model.RespBody{
		Message: "success",
		Data:    categoryList,
	})
}
