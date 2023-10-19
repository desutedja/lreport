package model

import "time"

type UserData struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserLevel string `json:"user_level"`
}

type UserListData struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	UserLevel string `json:"user_level"`
	CreatedOn string `json:"created_on"`
}

type LoginHistory struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Username  string    `json:"username"`
	Device    string    `json:"device"`
	IpAddress string    `json:"ip_address"`
	CreatedOn time.Time `json:"created_on"`
}

type BasicRequest struct {
	Search string
	Page   int
	Limit  int
}

type ResponseLogin struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	UserLevel string `json:"user_level"`
	Session   int    `json:"session"`
	Token     string `json:"token"`
}
