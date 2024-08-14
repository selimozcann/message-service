package main

import (
	"log"
	"messageservice/ws"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
)

const (
	PORT = ":3000"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"websocket_queue", // Queue name
		false,             // Durable?
		false,             // Delete when unused?
		false,             // Exclusive?
		false,             // No-wait?
		nil,               // Arguments
	)
	failOnError(err, "Failed to declare a queue")

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error check upgrading connection:", err)
			return
		}
		ws.HandleConnection(ch, q, conn)
	})

	log.Fatal(http.ListenAndServe(PORT, nil))
}
