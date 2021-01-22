package main


import (
	"fmt"
	"go_chatserver/global"
)

func main(){
	fmt.Println(global.Config.Redis.Host)
}
