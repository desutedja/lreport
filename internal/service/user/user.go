package user

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/desutedja/lreport/internal/repository/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userStore interface {
	Login(ctx context.Context, username string) (data model.UserData, err error)
	InsertLoginHistory(ctx context.Context, userID, device, ipAddress string) error
	CreateUser(ctx context.Context, username, password, userLevel string) (uuid.UUID, error)
	ChangePassword(ctx context.Context, userID, newPassword string) error
	LoginHistory(ctx context.Context, req model.BasicRequest) (data []model.LoginHistory, err error)
	UserList(ctx context.Context, req model.BasicRequest) (data []model.UserListData, err error)
}

type tokenStore interface {
	GenerateToken(id, userLevel string, timestamp time.Time) (string, error)
}

type Service struct {
	userStore  userStore
	tokenStore tokenStore
}

func NewService(userStore userStore, tokenStore tokenStore) *Service {
	return &Service{
		userStore:  userStore,
		tokenStore: tokenStore,
	}
}

func (s *Service) Login(ctx context.Context, username, password, device, ipAddress string) (resp model.ResponseLogin, err error) {
	// get password from db
	user, err := s.userStore.Login(ctx, username)
	if err != nil {
		return resp, errors.New("user not found")
	}

	// compare password to db password
	// assuming storedHashedPassword is the hashed password retrieved from the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// authentication failed
		log.Println("error validate password: ", err)
		return resp, errors.New("wrong password")
	}

	// authentication successful
	// generate token if password is valid
	timeNow := time.Now()
	token, err := s.tokenStore.GenerateToken(user.Id, user.UserLevel, timeNow)
	if err != nil {
		log.Println("error generate token: ", err)
		return resp, err
	}

	// insert log when success
	// insert userID, device, IP, timenow
	err = s.userStore.InsertLoginHistory(ctx, user.Id, device, ipAddress)
	if err != nil {
		// log info
		log.Println("error user login history: ", err)
	}

	resp = model.ResponseLogin{
		Id:        user.Id,
		Username:  username,
		UserLevel: user.UserLevel,
		Session:   3600, // second
		Token:     token,
	}

	return
}

func (s *Service) CreateUser(ctx context.Context, username, password, userLevel string) (uuid.UUID, error) {
	// check is username exist
	_, err := s.userStore.Login(ctx, username)
	if err == nil {
		// username is exist, log errornya
		return uuid.Nil, errors.New(model.ERROR_USER_EXIST)
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	// create user
	userID, err := s.userStore.CreateUser(ctx, username, string(hashedPassword), userLevel)
	if err != nil {
		// log info
		log.Println("error create user: ", err)
		return uuid.Nil, errors.New("error create user")
	}

	return userID, nil
}

func (s *Service) ChangePassword(ctx context.Context, username, oldPassword, newPassword string) error {
	// get password from db
	user, err := s.userStore.Login(ctx, username)
	if err != nil {
		return errors.New("data not valid, user not found")
	}

	// compare password to db password
	// assuming storedHashedPassword is the hashed password retrieved from the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		// authentication failed
		log.Println("error validate password: ", err)
		return errors.New("wrong old password")
	}

	// compare with new password
	// check if old password is same like new password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword))
	if err == nil {
		// authentication failed
		return errors.New("new password is same like old password")
	}

	// hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("new password error: ", err)
		return errors.New("new password error")
	}

	err = s.userStore.ChangePassword(ctx, user.Username, string(hashedPassword))
	if err != nil {
		// log info
		log.Println("error change password failed: ", err)
		return errors.New("change password failed")
	}

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, username, newPassword string) error {
	// get password from db
	user, err := s.userStore.Login(ctx, username)
	if err != nil {
		return errors.New("data not valid, user not found")
	}

	// hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("new password error: ", err)
		return errors.New("new password error")
	}

	err = s.userStore.ChangePassword(ctx, user.Username, string(hashedPassword))
	if err != nil {
		// log info
		log.Println("error reset password failed: ", err)
		return errors.New("reset password failed")
	}

	return nil
}

func (s *Service) GetLoginHistory(ctx context.Context, req model.BasicRequest) (data []model.LoginHistory, err error) {
	// get password from db
	data, err = s.userStore.LoginHistory(ctx, req)
	if err != nil {
		return data, errors.New("data not found")
	}

	return data, nil
}

func (s *Service) GetUserList(ctx context.Context, req model.BasicRequest) (data []model.UserListData, err error) {
	// get password from db
	data, err = s.userStore.UserList(ctx, req)
	if err != nil {
		return data, errors.New("data not found")
	}

	return data, nil
}
