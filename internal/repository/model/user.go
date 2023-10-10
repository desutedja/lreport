package model

import "time"

type UserData struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserLevel string `json:"user_level"`
}

type LoginHistory struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Device    string    `json:"device"`
	IpAddress string    `json:"ip_address"`
	CreatedOn time.Time `json:"created_on"`
}

type BasicRequest struct {
	Search string
	Page   int
	Limit  int
}
