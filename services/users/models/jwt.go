package models

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenToken(user User, key string) string {
	jwt_token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	jwt_token.Claims = jwt.MapClaims{
		"id":       user.ID,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
		"roles":    user.Role,
		"username": user.Username,
	}
	// Sign and get the complete encoded token as a string
	token, err := jwt_token.SignedString([]byte(key))
	if err != nil {
		fmt.Println(err, key)
	}
	return token
}

func ParseToken(tokenString string, key string) (jwt.MapClaims, error) {
	tokenStr := tokenString[7:] //去掉 bearer 字符串
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
