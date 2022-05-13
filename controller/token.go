package controller

import (
	"errors"
	"fmt"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"time"
)

func GenerateToken(username string, userid string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
		Username: username,
		UserID:   userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	})
	fmt.Println("Token = ", claims)
	tokenString, err := claims.SignedString(G.TokenSecret)
	G.RedisDB.HSet(utils.PREFIX_TOKEN, tokenString, userid)
	return tokenString, err
}

func CheckToken(tokenString string) (*model.Claims, error) {
	//检验tokenString是否存在
	_, err := G.RedisDB.HGet(utils.PREFIX_TOKEN, tokenString).Result()
	if err == redis.Nil {
		return nil, errors.New("the tokenKey doesn't exist")
	} else if err != nil {
		return nil, errors.New(fmt.Sprintf("Redis HGet Error:%v", err))
	}
	//解析tokenString
	claim := &model.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return G.TokenSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	//如果距离过期时间小于 1 min,更新过期时间
	if time.Unix(claim.ExpiresAt, 0).Sub(time.Now()) > time.Minute {
		return claim, nil
	}
	claim.ExpiresAt = time.Now().Add(time.Hour * 2).Unix()
	return claim, nil
}
