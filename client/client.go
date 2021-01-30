package client

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"go_chatserver/model"
	"go.uber.org/zap"
	"go_chatserver/util"
	"go_chatserver/redisconn"
	"golang.org/x/net/context"
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

func (c *Client) Read(ctx context.Context){
	defer func(){
		if r := recover(); r != nil{
			zap.S().Error(r)
			c.Socket.Close()
		}
	}()
	LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		default:
			_, message, err := c.Socket.ReadMessage()
			if err != nil{
				return
			}
			c.HandleMessage(message)
		}
	}
	return
}

func (c *Client) OnMessagePub(ctx context.Context){
	defer func(){
		if r := recover(); r != nil{
			fmt.Println("stop current goroutine")
			return
		}
	}()

	if c.Pub != nil{
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			default:
				select {
				case mg := <- c.Pub.Channel():
					fmt.Println(mg.Payload)
				default:
				}
			}
		}
	}else{
		return
	}
	return
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
		c.PublishToAllUsersInRoom(persistModel.UserIds, responseMessage, c.Domain)

	}
}

func (c *Client) PublishToAllUsersInRoom(userIds []int, message map[string]interface{}, domain string){
	for _, id := range userIds{
		userMessage := buildMessage(id, message)
		if len(userMessage) == 0{
			return
		}else{
			data, _ := json.Marshal(userMessage)
			channel := util.GetChannel(id, domain)
			zap.S().Info("Publish message to " + channel)
			err := redisconn.RedisClient.Publish(channel, data).Err()
			if err != nil {
				zap.S().Error("Publish to " + channel + " failed")
			}
		}

	}
}

func buildMessage(id int, message map[string]interface{})map[string]interface{}{
	dict := message
	if message["from_user_id"].(int) == id{
		dict["read"] = true
		dict["mc"] = "fade-message"
	}else{
		dict["read"] = false
		dict["mc"] = "fade-message unread"
	}
	dict["body"] = message["body"]
	return dict
}

