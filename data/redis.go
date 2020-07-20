package data

import (
	"github.com/go-redis/redis"
	"simple-sso/util"
)

var RedisCli *redis.Client

func init() {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     util.Redis.Host + ":" + util.Redis.Port,
		Password: util.Redis.Password,
		DB:       0,
	})
}



