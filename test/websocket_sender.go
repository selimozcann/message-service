package test

/*
package main

import (
	"log"
	"messageservice/ws"
	"net/http"

	"github.com/gorilla/websocket" // ws paketini import edin

	"github.com/rabbitmq/amqp091-go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// RabbitMQ bağlantısını oluştur
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"my_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error while upgrading connection:", err)
			return
		}
		connection := ws.NewConnection(map[*websocket.Conn]bool{conn: true})
		go ws.ReadMessage(ch, q, connection)
	})

	log.Println("WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
*/
