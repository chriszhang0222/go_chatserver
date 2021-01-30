package main

import (
	context2 "context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"go_chatserver/client"
	"go_chatserver/global"
	"go_chatserver/redisconn"
	router2 "go_chatserver/router"
	"go_chatserver/util"
	"golang.org/x/net/context"
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
	ctx, cancel := context2.WithCancel(context.Background())
	conn.SetCloseHandler(func(code int, text string) error {
		zap.S().Info("Close conn for " + channel)
		cancel()
		client.Pub.Unsubscribe(channel)
		client.Pub.Close()
		client.Socket.Close()
		return nil
	})
	go client.OnMessagePub(ctx)
	go client.Read(ctx)

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
	Router := router2.InitRouter()
	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", 9000)); err != nil {
			zap.S().Panic("serve error", err.Error())
		}
		zap.S().Debugf("serve goods server at %d", 9000)
	}()
	addr := fmt.Sprintf(":%d", global.Config.Port)
	http.ListenAndServe(addr, nil)
}
