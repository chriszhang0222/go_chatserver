package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"go_chatserver/global"
)


func WebsocketHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("hello world"))
}
func main(){
	rtr := mux.NewRouter()
	rtr.HandleFunc("/chat/{id:\\d+}/{company:\\d+}/{domain:[a-zA-Z0-9\\-\\_]*}", WebsocketHandler).Methods("GET")

	http.Handle("/", rtr)

	log.Println("Listening...")
	http.ListenAndServe(fmt.Sprintf(":%d", global.Config.Port), nil)
}
