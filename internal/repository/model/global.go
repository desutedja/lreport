package model

type key string

const (
	ERROR_USER_EXIST     string = "username already exist"
	ERROR_WRONG_PASSWORD string = "wrong password"

	USER_LEVEL_DEFAULT string = "user"
	CONTEXT_KEY        key    = "userInfo"
)

type RespBody struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
