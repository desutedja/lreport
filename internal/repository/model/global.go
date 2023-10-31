package model

type key string

const (
	ERROR_USER_EXIST     string = "username already exist"
	ERROR_WRONG_PASSWORD string = "wrong password"

	USER_LEVEL_DEFAULT string = "user"
	CONTEXT_KEY        key    = "userInfo"
	TIME_YYYYMMDD      string = "2006-01-02"
)

type RespBody struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RespListstruct struct {
	Items        interface{} `json:"items"`
	TotalItems   int         `json:"total_page"`
	FilteredPage int         `json:"filtered_page"`
}
