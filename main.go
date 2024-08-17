package main

import (
	"database/sql"
	"log"
	"messageservice/ws"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq" // PostgreSQL driver'ı
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

	connStr := "user=username password=password dbname=dbName sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	failOnError(err, "PostgreSQL connection is not successful")
	defer db.Close()

	err = db.Ping()
	failOnError(err, "PostgreSQL is not successful")

	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "RabbitMQ connection is not successful")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "RabbitMQ channel is not successful")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"websocket_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Queue is not successful")

	r := gin.Default()

	activeConnections := make(map[*websocket.Conn]bool)

	r.GET("/ws", func(c *gin.Context) {
		wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Fatal(err)
		}

		activeConnections[wsConn] = true

		go func() {
			ws.HandleConnection(ch, q, wsConn, db)
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

		log.Println("Server is shutting down...")
		os.Exit(0)
	}()

	r.Run(PORT)
}
