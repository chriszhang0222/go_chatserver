package redisconn

import (
	"fmt"
	"go_chatserver/global"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func init(){
	redisConfig := global.Config.Redis
	strs := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: "", // no password set
		DB:       redisConfig.Db,  // use default DB
	})
}