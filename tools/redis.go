package tools

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {

	opt, err := redis.ParseURL("redis://:werido@8.131.78.84:6379/4")
	if err != nil {
		panic(err)
	}
	// fmt.Println("redis url:", opt)
	rdb := redis.NewClient(opt)
	return rdb
}

func InvalidTokenKey(token string) string {
	return "invalid_token." + "." + token
}

func UserPermissionKey(userID int) string {
	return "permission" + ":" + fmt.Sprintf("%d", userID)
}

// func ResourceAuthenticatedCacheKey(resource string) string {
// 	return "resource_authenticated" + ":" + resource
// }

const ResourceAuthenticatedCacheKey = "site:resource_authenticated"
