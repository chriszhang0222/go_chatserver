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
			fmt.Println("stop current goroutine")
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
	messageType := model.Type
	if messageType == "chat"{
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
	data := map[string]interface{}{
		"room_id": message.Roomid,
		"user_id": message.User,
		"identifier": message.Identifier,
		"message_content": message.MessageBody,
	}
	response, err = util.SendRequestWithAuth(message.GetPersistUri(), data, token.Access)
	if err != nil {
		c.Socket.WriteJSON(map[string]string{
			"type": "info",
			"success": "False",
			"message": "Persist to db failed",
		})
		return
	}
	persistModel := &model.PersistResponse{}
	err = json.Unmarshal(response, persistModel)
	if err != nil {
		return
	}
	if !persistModel.Success{
		c.Socket.WriteJSON(map[string]interface{}{
			"type": "info",
			"success": "false",
			"message": persistModel.Message,
		})
	}else{
		responseMessage := persistModel.Message
		c.Socket.WriteJSON(map[string]interface{}{
			"type": "info_chat",
			"success": true,
			"message": "Message sent out successfully",
			"id": responseMessage["id"],
			"identifier": message.Identifier,
			"discussion_date": persistModel.DiscussionDate,
			"discussion_time": persistModel.DiscussionTime,
		})
		responseMessage["type"] = "chat"
		responseMessage["identifier"] = message.Identifier
		c.PublishToAllUsersInRoom(persistModel.UserIds, message)

	}
}

func (c *Client) PublishToAllUsersInRoom(userIds []int, message *model.Message){
	for _, v := range userIds{
		user_message := buildMessage(v, message)
		if len(user_message) == 0{
			return
		}

	}
}

func buildMessage(id int, message *model.Message)map[string]interface{}{
	dict := map[string]interface{}{}
	return dict
}

