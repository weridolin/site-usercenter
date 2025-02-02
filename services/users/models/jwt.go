package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenToken(user User, key string) string {
	jwt_token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	jwt_token.Claims = jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		// "roles": roles,
		// "menus":    menus,
		"username":   user.Username,
		"superAdmin": user.IsSuperAdmin,
	}
	// Sign and get the complete encoded token as a string
	token, err := jwt_token.SignedString([]byte(key))
	if err != nil {
		fmt.Println(err, key)
	}
	return token
}

func ParseToken(tokenString string, key string) (jwt.MapClaims, error) {
	var tokenStr = tokenString
	// fmt.Println(strings.HasPrefix(tokenString, "Bearer"))\
	fmt.Println("jwt token -> ", tokenString)
	if strings.HasPrefix(tokenString, "Bearer") {
		tokenStr = tokenString[7:] //去掉 bearer 字符串
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}

/*
oauth生成防伪token
key: 密钥
expire: 过期时间
*/
func GenStateToken(key string, expire time.Duration) string {
	jwt_token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	jwt_token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(expire).Unix(),
	}
	// Sign and get the complete encoded token as a string
	token, err := jwt_token.SignedString([]byte(key))
	if err != nil {
		fmt.Println("generate state token error ->", err, key)
	}
	return token
}

func ParseStateToken(tokenString string, key string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
