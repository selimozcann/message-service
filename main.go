package main

import (
	"log"
	"messageservice/ws"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	PORT = ":3000"
)

func main() {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error check upgrading connection:", err)
			return
		}
		ws.HandleConnection(conn)
	})
	log.Fatal(http.ListenAndServe(PORT, nil))
}
