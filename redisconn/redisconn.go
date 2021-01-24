package redisconn

import (
	"fmt"
	"go.uber.org/zap"
	"go_chatserver/global"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func init(){
	redisConfig := global.Config.Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: "", // no password set
		DB:       redisConfig.Db,  // use default DB
	})
	_, err := RedisClient.Ping().Result()
	if err != nil{
		zap.S().Error("Error when connected to redis" + fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port))
		panic(err)
	}
}