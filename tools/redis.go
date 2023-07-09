package tools

import "github.com/go-redis/redis"

func NewRedisClient() *redis.Client {
	opt, err := redis.ParseURL("redis://:werido@8.131.78.84:6379/4")
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)
	return rdb
}

func InvalidTokenKey(token string) string {
	return "invalid_token." + "." + token
}
