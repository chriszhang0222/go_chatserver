package main

import (
	"encoding/json"
	"go_chatserver/redisconn"
)

func main(){
	j, _ := json.Marshal(
		map[string]interface{}{
			"room_id": 1,
			"user_id": 3,
			"identifier": "message.Identifier",
			"message_content": "ok",
			"read": true,
		})
	redisconn.RedisClient.Publish("channel1",j)
}
