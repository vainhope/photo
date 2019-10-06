package utils

import (
	"github.com/dgrijalva/jwt-go"
	"goweb/base"
	"goweb/config"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	Id       int64  `json:"id"`
	jwt.StandardClaims
}

//依赖 config 的先初始化
var jwtSecret = []byte(config.Setting.Jwt.Secret)
var expireTime = config.Setting.Jwt.ExpireTime

/**
生成token
*/
func GenerateToken(username string, id int64) (string, error) {
	if username == "" {
		return "", base.Err(-1, "加密参数不能为空")
	}
	now := time.Now()
	endTime := now.Add(time.Minute * time.Duration(expireTime))
	claims := Claims{username, id,
		jwt.StandardClaims{
			ExpiresAt: endTime.Unix()}}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokenClaims.SignedString(jwtSecret)
}

/**
解析token
*/

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecret, nil
	})
	if nil != tokenClaims {
		//时间没有过期
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
