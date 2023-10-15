package user

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/desutedja/lreport/internal/repository/model"
	"github.com/desutedja/lreport/pkg/token"
	"github.com/google/uuid"
	"github.com/thedevsaddam/renderer"
)

type userService interface {
	CreateUser(ctx context.Context, username, password, userLevel string) (uuid.UUID, error)
	Login(ctx context.Context, username, password, device, ipAddress string) (resp model.ResponseLogin, err error)
	ResetPassword(ctx context.Context, username, newPassword string) error
	ChangePassword(ctx context.Context, username, oldPassword, newPassword string) error
	GetLoginHistory(ctx context.Context, req model.BasicRequest) (data []model.LoginHistory, err error)
	GetUserList(ctx context.Context, req model.BasicRequest) (data []model.UserListData, err error)
}

type Handler struct {
	render      *renderer.Render
	userService userService
}

func NewHandler(userService userService) *Handler {
	render := renderer.New()
	return &Handler{
		render:      render,
		userService: userService,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body := model.UserData{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.render.JSON(w, http.StatusBadRequest, "bad request")
		return
	}

	// validate request
	if body.Username == "" {
		h.render.JSON(w, http.StatusBadRequest, "username is required")
		return
	}

	if body.Password == "" {
		h.render.JSON(w, http.StatusBadRequest, "password is required")
		return
	}

	// if user level is empty then default level is user
	if body.UserLevel == "" {
		body.UserLevel = model.USER_LEVEL_DEFAULT
	}

	_, err = h.userService.CreateUser(ctx, body.Username, body.Password, body.UserLevel)
	if err != nil {
		if err.Error() == model.ERROR_USER_EXIST {
			h.render.JSON(w, http.StatusBadRequest, model.ERROR_USER_EXIST)
			return
		}

		h.render.JSON(w, http.StatusInternalServerError, err)
		return
	}

	h.render.JSON(w, http.StatusOK, "User successfully created")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body := model.UserData{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.render.JSON(w, http.StatusBadRequest, "bad request")
		return
	}

	// validate request
	if body.Username == "" {
		h.render.JSON(w, http.StatusBadRequest, "username is required")
		return
	}

	if body.Password == "" {
		h.render.JSON(w, http.StatusBadRequest, "password is required")
		return
	}

	// Get the client's IP address
	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	fmt.Println("Client IP:", clientIP)

	// Get the User-Agent header
	userAgent := r.Header.Get("User-Agent")
	fmt.Println("User-Agent:", userAgent)

	userToken, err := h.userService.Login(ctx, body.Username, body.Password, userAgent, clientIP)
	if err != nil {
		if err.Error() == model.ERROR_WRONG_PASSWORD {
			h.render.JSON(w, http.StatusUnauthorized, model.RespBody{
				Message: "failed",
				Data:    err.Error(),
			})
			return
		}

		h.render.JSON(w, http.StatusInternalServerError, model.RespBody{
			Message: "failed",
			Data:    err.Error(),
		})
		return
	}

	h.render.JSON(w, http.StatusOK, model.RespBody{
		Message: "success",
		Data:    userToken,
	})
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := struct {
		Id          string `json:"id"`
		Username    string `json:"username"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.render.JSON(w, http.StatusBadRequest, "bad request")
		return
	}

	// validate request
	if body.Username == "" {
		h.render.JSON(w, http.StatusBadRequest, "username is required")
		return
	}

	if body.OldPassword == "" {
		h.render.JSON(w, http.StatusBadRequest, "password is required")
		return
	}

	if body.NewPassword == "" {
		h.render.JSON(w, http.StatusBadRequest, "new password is required")
		return
	}

	err = h.userService.ChangePassword(ctx, body.Username, body.OldPassword, body.NewPassword)
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

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenPayload := ctx.Value("userInfo").(token.TokenClaims)

	// kalau user levelnya user tidak bisa reset password
	if tokenPayload.Userlevel == model.USER_LEVEL_DEFAULT {
		h.render.JSON(w, http.StatusUnauthorized, "your level is not eligible to do reset password")
		return
	}

	body := model.UserData{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.render.JSON(w, http.StatusBadRequest, "bad request")
		return
	}

	// validate request
	if body.Username == "" {
		h.render.JSON(w, http.StatusBadRequest, "username is required")
		return
	}

	if body.Password == "" {
		h.render.JSON(w, http.StatusBadRequest, "password is required")
		return
	}

	err = h.userService.ResetPassword(ctx, body.Username, body.Password)
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

func (h *Handler) LoginHistory(w http.ResponseWriter, r *http.Request) {
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

	loginHistory, err := h.userService.GetLoginHistory(ctx, body)
	if err != nil {
		h.render.JSON(w, http.StatusInternalServerError, model.RespBody{
			Message: "failed",
			Data:    err.Error(),
		})
		return
	}

	h.render.JSON(w, http.StatusOK, model.RespBody{
		Message: "success",
		Data:    loginHistory,
	})
}

func (h *Handler) UserList(w http.ResponseWriter, r *http.Request) {
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

	userList, err := h.userService.GetUserList(ctx, body)
	if err != nil {
		h.render.JSON(w, http.StatusInternalServerError, model.RespBody{
			Message: "failed",
			Data:    err.Error(),
		})
		return
	}

	filteredPage := math.Ceil(float64(len(userList)) / float64(body.Limit))

	resp := model.RespListstruct{
		Items:        userList,
		TotalItems:   len(userList),
		FilteredPage: int(filteredPage),
	}

	h.render.JSON(w, http.StatusOK, model.RespBody{
		Message: "success",
		Data:    resp,
	})
}
