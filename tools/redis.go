package tools

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(RedisUri string) *redis.Client {

	opt, err := redis.ParseURL(RedisUri)
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
