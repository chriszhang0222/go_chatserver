package client

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"go_chatserver/model"
	"go.uber.org/zap"
	"go_chatserver/util"
)

type Client struct{
	Addr string
	Socket *websocket.Conn
	Pub *redis.PubSub
	Send chan []byte
	UserId int
	Domain string
}


func NewClient(Addr string, Socket *websocket.Conn, Pub *redis.PubSub , UserId int, Domain string) *Client{
	return &Client{
		Addr: Addr,
		Socket: Socket,
		Pub: Pub,
		UserId: UserId,
		Domain: Domain,
	}
}

func (c *Client) Read(){
	defer func(){
		if r := recover(); r != nil{
			zap.S().Error(r)
			c.Socket.Close()
		}
	}()
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil{
			return
		}
		c.HandleMessage(message)
	}
}

func (c *Client) OnMessagePub(){
	defer func(){
		if r := recover(); r != nil{
			return
		}
	}()

	if c.Pub != nil{
		for {
			if c.Pub == nil{
				break
			}
			select {
				case mg := <- c.Pub.Channel():
					fmt.Println(mg.Payload)
			default:

			}
		}
	}else{
		return
	}
}

func (c *Client) HandleMessage(message []byte){
	model := &model.Message{}
	err := json.Unmarshal(message, model)
	if err != nil {
		return
	}
	message_type := model.Type
	if message_type == "chat"{
		c.PublishMessage(model)
	}else{
		return
	}

}

func (c *Client) PublishMessage(message *model.Message){
	uri := message.GetUri()
	body := map[string]interface{}{
		"user_id": message.User,
		"company_id": message.Company,
	}

	response, err := util.SendRequest(uri, body)
	if err != nil {
		c.Socket.WriteJSON(map[string]string{
			"type": "info",
			"success": "False",
			"message": "Authentication Failed!'",
		})
		return

	}
	token := &model.Token{}
	err = json.Unmarshal(response, token)
	if err != nil {
		c.Socket.WriteJSON(map[string]string{
			"type": "info",
			"success": "False",
			"message": "Authentication Failed!'",
		})
		return
	}


}

