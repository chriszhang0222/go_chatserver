package main

import (
	"github.com/go-redis/redis"
	"go_chatserver/global"
	"go_chatserver/redisconn"
	"log"
	"sync"
	"time"
)
var wg sync.WaitGroup
var Ok chan bool = make(chan bool)
func Simulate(pb *redis.PubSub){
	for {
		//if <- signal{
		//	wg.Done()
		//}
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

func stop(){
	time.Sleep(5* time.Second)
	global.Signal <- true
}
func main(){
	pub := redisconn.RedisClient.Subscribe("channel1")
	go Simulate(pub)
	go stop()
	<- global.Signal
	//wg.Done()
	//wg.Wait()


}
