package utils

import (
	"errors"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/dgrijalva/jwt-go"
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
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	})

	tokenString, err := claims.SignedString(G.TokenSecret)

	// FIXME: jwt 不需要存储
	//G.RedisDB.HSet(utils.PREFIX_REFRESH_TOKEN, tokenString, userid)

	return tokenString, err
}

func CheckToken(tokenString string) (*Claims, error) {

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
	//如果距离过期时间小于 1 min,更新过期时间
	if time.Unix(claim.ExpiresAt, 0).Sub(time.Now()) > time.Minute {
		// FIXME: Token 过期为未授权
		return claim, nil
	}
	claim.ExpiresAt = time.Now().Add(time.Hour * 2).Unix()
	return claim, nil
}
