package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"go_chatserver/redisconn"
	"log"
	"sync"
)
var wg sync.WaitGroup
func Simulate(pb *redis.PubSub){
	for {
		select {
		case mg := <-pb.Channel():
			// 等待从 channel 中发布 close 关闭服务
			if mg.Payload == "close" {
				// 当
				wg.Done()
			} else {
				log.Println("接channel信息", mg.Payload)
			}
		default:
		}
	}

}
func main(){
	//wg.Add(1)
	//pub := redisconn.RedisClient.Subscribe("channel1")
	//go Simulate(pub)
	//wg.Wait()
    fmt.Println(redisconn.RedisClient.Get("10"))

}
