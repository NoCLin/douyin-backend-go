package utils

import (
	"errors"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(username string, userid string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: username,
		UserID:   userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	tokenString, err := claims.SignedString(G.TokenSecret)

	return tokenString, err
}

func CheckToken(tokenString string) (*Claims, error) {
	log.Println("CheckToken is calling")
	claim := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return G.TokenSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// TODO: 动态配置时间
	if time.Unix(claim.ExpiresAt, 0).Sub(time.Now()) <= 0 {
		return claim, errors.New("the token has expired,please login again")
	}

	log.Println("CheckToken successfully")
	return claim, nil
}
