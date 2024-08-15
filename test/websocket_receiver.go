package test

/*
import (
	"log"
	"messageservice/ws"

	"github.com/gorilla/websocket"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
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

	connection := ws.NewConnection(make(map[*websocket.Conn]bool))
	go ws.WriteMessage(ch, q, connection)
	select {}
}
*/
