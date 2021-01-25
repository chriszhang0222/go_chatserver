package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"go_chatserver/global"
	"go_chatserver/redisconn"

	"sync"
	"time"
)
var wg sync.WaitGroup
var Ok chan bool = make(chan bool)
type Message struct{
	Room int `json:"room_id"`
	Read bool `json:"read"`
}
func Simulate(pb *redis.PubSub){
	for {
		select {
		case mg := <-pb.Channel():
			// 等待从 channel 中发布 close 关闭服务
			if mg.Payload == "close" {
				// 当
				wg.Done()
			} else {
				a := &Message{}
				json.Unmarshal([]byte(mg.Payload), a)
				fmt.Println(a)

			}
		default:
		}
	}

}

func stop(){
	time.Sleep(5* time.Second)
	global.Signal <- true
}
func main(){
	wg.Add(1)
	pb := redisconn.RedisClient.Subscribe("channel1")
	go Simulate(pb)
	wg.Wait()

}
