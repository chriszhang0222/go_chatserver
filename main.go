package main

import (
	"encoding/json"
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
	"runtime"
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
		zap.S().Info("Close conn for " + channel)
		client.Pub.Unsubscribe(channel)
		client.Pub.Close()
		client.Socket.Close()
		return nil
	})
	go client.OnMessagePub()
	go client.Read()

}

func SystemHandler(w http.ResponseWriter, req *http.Request){
	numGoroutine := runtime.NumGoroutine()
	numCPU := runtime.NumCPU()
	data := map[string]interface{}{}
	data["cpu"] = numCPU
	data["goroutine"] = numGoroutine
	d, _ := json.Marshal(data)
	w.Write(d)
}

func main(){
	rtr := mux.NewRouter()
	rtr.HandleFunc("/chat/{id:\\d+}/{domain:[a-zA-Z0-9\\-\\_]*}", WebsocketHandler).Methods("GET")
	rtr.HandleFunc("/system/state", SystemHandler).Methods("GET")

	http.Handle("/", rtr)

	log.Println("Listening...")
	addr := fmt.Sprintf(":%d", global.Config.Port)
	http.ListenAndServe(addr, nil)
}
