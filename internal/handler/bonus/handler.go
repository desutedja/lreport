package bonus

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/desutedja/lreport/internal/repository/model"
	"github.com/desutedja/lreport/pkg/token"
	"github.com/thedevsaddam/renderer"
)

type bonusService interface {
	CreateBonus(ctx context.Context, userId string, req model.ReqBonus) error
	GetBonus(ctx context.Context, req model.BasicRequest) (data []model.DataBonus, err error)
}

type Handler struct {
	render       *renderer.Render
	bonusService bonusService
}

func NewHandler(bonusService bonusService) *Handler {
	render := renderer.New()
	return &Handler{
		render:       render,
		bonusService: bonusService,
	}
}

func (h *Handler) CreateBonus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenPayload := ctx.Value(model.CONTEXT_KEY).(token.TokenClaims)

	body := model.ReqBonus{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.render.JSON(w, http.StatusBadRequest, "bad request")
		return
	}

	// validate request
	if body.CategoryId == 0 {
		h.render.JSON(w, http.StatusBadRequest, "category is required")
		return
	}

	if body.TransDate == "" {
		h.render.JSON(w, http.StatusBadRequest, "trans date is required")
		return
	}

	// check is request body transdate is date
	transDate, err := time.Parse("2006-01-02", body.TransDate)
	if err != nil {
		h.render.JSON(w, http.StatusBadRequest, "trans date must date")
		return
	}

	// make transdate format to (yyyy-mm-dd)
	dateString := transDate.Format("2006-01-02")
	body.TransDate = dateString

	err = h.bonusService.CreateBonus(ctx, tokenPayload.Id, body)
	if err != nil {
		h.render.JSON(w, http.StatusBadRequest, model.RespBody{
			Message: "failed",
			Data:    err.Error(),
		})
		return
	}

	h.render.JSON(w, http.StatusOK, model.RespBody{
		Message: "success",
		Data:    "success",
	})
}

func (h *Handler) GetBonus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	param := r.URL.Query()
	page, _ := strconv.Atoi(param.Get("page"))
	limit, _ := strconv.Atoi(param.Get("limit"))
	search := param.Get("search")

	// validate request
	if page == 0 {
		h.render.JSON(w, http.StatusBadRequest, "page is required")
		return
	}

	if limit == 0 {
		h.render.JSON(w, http.StatusBadRequest, "limit is required")
		return
	}

	body := model.BasicRequest{
		Search: search,
		Page:   page,
		Limit:  limit,
	}

	transactionData, err := h.bonusService.GetBonus(ctx, body)
	if err != nil {
		h.render.JSON(w, http.StatusInternalServerError, model.RespBody{
			Message: "failed",
			Data:    err.Error(),
		})
		return
	}

	testStruct := struct {
		Items        interface{} `json:"items"`
		TotalItems   int         `json:"total_page"`
		FilteredPage int         `json:"filtered_page"`
	}{}
	testStruct.Items = transactionData
	h.render.JSON(w, http.StatusOK, model.RespBody{
		Message: "success",
		Data:    testStruct,
	})
}
