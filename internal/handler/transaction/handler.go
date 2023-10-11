package transaction

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

type transactionService interface {
	CreateTransaction(ctx context.Context, userId string, req model.ReqTransaction) error
	GetTransaction(ctx context.Context, req model.BasicRequest) (data []model.DataTransaction, err error)
}

type Handler struct {
	render             *renderer.Render
	transactionService transactionService
}

func NewHandler(transactionService transactionService) *Handler {
	render := renderer.New()
	return &Handler{
		render:             render,
		transactionService: transactionService,
	}
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenPayload := ctx.Value(model.CONTEXT_KEY).(token.TokenClaims)

	body := model.ReqTransaction{}
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

	if body.Regis == 0 {
		h.render.JSON(w, http.StatusBadRequest, "regis is required")
		return
	}

	if body.RegisDp == 0 {
		h.render.JSON(w, http.StatusBadRequest, "regis dp is required")
		return
	}

	if body.ActivePlayer == 0 {
		h.render.JSON(w, http.StatusBadRequest, "active player is required")
		return
	}

	if body.TransDp == 0 {
		h.render.JSON(w, http.StatusBadRequest, "trans dp is required")
		return
	}

	if body.TransWd == 0 {
		h.render.JSON(w, http.StatusBadRequest, "trans wd is required")
		return
	}

	if body.TotalDp == 0 {
		h.render.JSON(w, http.StatusBadRequest, "total dp is required")
		return
	}

	if body.TotalWd == 0 {
		h.render.JSON(w, http.StatusBadRequest, "total wd is required")
		return
	}

	if body.Wl == 0 {
		h.render.JSON(w, http.StatusBadRequest, "wl is required")
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

	err = h.transactionService.CreateTransaction(ctx, tokenPayload.Id, body)
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

func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
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

	transactionData, err := h.transactionService.GetTransaction(ctx, body)
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
