package model

import "github.com/golang-jwt/jwt/v4"

type MyClaims struct {
	jwt.RegisteredClaims
	Id        string `json:"id"`
	Userlevel string `json:"userlevel"`
}
