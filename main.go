package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"go_chatserver/client"
	"go_chatserver/global"
	"go_chatserver/redisconn"
	"go_chatserver/util"
	"log"
	"net/http"
	"strconv"
	"strings"
)


func WebsocketHandler(w http.ResponseWriter, req *http.Request){
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}).Upgrade(w, req, nil)
	if err != nil{
		http.NotFound(w, req)
		return
	}
	uri := strings.Split(req.RequestURI, "/")
	userId, _ := strconv.Atoi(uri[2])
	domain := uri[3]
	channel := util.GetChannel(userId, domain)
	zap.S().Info("connect to "+ channel)
	pubsub := redisconn.RedisClient.Subscribe(channel)
	client := client.NewClient(req.RemoteAddr, conn, pubsub, userId, domain)
	conn.SetCloseHandler(func(code int, text string) error {
		zap.S().Info("Close conn for " + string(userId) + " " + domain)
		client.Pub.Unsubscribe(channel)
		return nil
	})
	go client.OnMessagePub()
	go client.Read()

}
func main(){
	rtr := mux.NewRouter()
	rtr.HandleFunc("/chat/{id:\\d+}/{domain:[a-zA-Z0-9\\-\\_]*}", WebsocketHandler).Methods("GET")

	http.Handle("/", rtr)

	log.Println("Listening...")
	addr := fmt.Sprintf(":%d", global.Config.Port)
	http.ListenAndServe(addr, nil)
}
