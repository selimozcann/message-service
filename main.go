package main

import (
	"log"
	"messageservice/ws"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
)

const PORT = ":8080"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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

	r := gin.Default()

	activeConnections := make(map[*websocket.Conn]bool)

	r.GET("/ws", func(c *gin.Context) {
		wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Fatal(err)
		}

		activeConnections[wsConn] = true

		go func() {
			ws.HandleConnection(ch, q, wsConn)
			delete(activeConnections, wsConn)
		}()
	})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		for conn := range activeConnections {
			conn.Close()
		}

		log.Println("Server shutting down gracefully")
		os.Exit(0)
	}()

	r.Run(PORT)
}
