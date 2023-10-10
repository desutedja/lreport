package token

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/desutedja/lreport/internal/repository/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/thedevsaddam/renderer"
)

type TokenGenerator struct {
	SaltKey             string
	ExpiredTimeInSecond int
	render              *renderer.Render
}

func NewTokenGenerator(saltkey string, expiredTimeInSecond int) *TokenGenerator {
	renderer := renderer.New()
	return &TokenGenerator{
		SaltKey:             saltkey,
		ExpiredTimeInSecond: expiredTimeInSecond,
		render:              renderer,
	}
}

func (t *TokenGenerator) GenerateToken(id, userLevel string, timestamp time.Time) (string, error) {
	claims := model.MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "lucky",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(t.ExpiredTimeInSecond) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(timestamp),
		},
		Id:        id,
		Userlevel: userLevel,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(t.SaltKey)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

type TokenClaims struct {
	jwt.StandardClaims
	Id        string
	Userlevel string
}

func (t *TokenGenerator) MiddlewareJWTAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			t.render.JSON(w, http.StatusUnauthorized, "Token Not Found")
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("signing method error")
			} else if method != jwt.SigningMethodHS256 {
				return nil, errors.New("signing method invalid")
			}

			secretKey := []byte(t.SaltKey)
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			t.render.JSON(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			t.render.JSON(w, http.StatusUnauthorized, "Wrong Token")
			return
		}

		b, err := json.Marshal(claims)
		if err != nil {
			t.render.JSON(w, http.StatusUnauthorized, "Error Token")
			return
		}

		requestHead := TokenClaims{}
		err = json.Unmarshal(b, &requestHead)
		if err != nil {
			t.render.JSON(w, http.StatusUnauthorized, "Error Token !")
			return
		}

		ctx := context.WithValue(context.Background(), model.CONTEXT_KEY, requestHead)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
