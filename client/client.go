package client

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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
			zap.S().Error(err)
			return
		}
		fmt.Println(string(message))
	}
}

func (c *Client) OnMessagePub(){
	if c.Pub != nil{
		for {
			select {
				case mg := <- c.Pub.Channel():
					fmt.Println(mg.Payload)
			default:

			}
		}
	}
}

