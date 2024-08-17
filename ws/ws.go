package ws

import (
	"database/sql"
	"log"
	handler "messageservice/handlers"
	"messageservice/rabbitmq"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	db     *handler.DB
	wsConn map[*websocket.Conn]bool
	send   chan []byte
}

func NewConnection(conn map[*websocket.Conn]bool) *Connection {
	return &Connection{wsConn: conn, send: make(chan []byte, 256)}
}

func WriteMessage(ch *amqp091.Channel, q amqp091.Queue, c *Connection) {
	msgs, err := rabbitmq.ConsumeMessage(ch, q)
	if err != nil {
		log.Printf("Failed to publish message:%v", err)
		return
	}
	go func() {
		for msg := range msgs {
			log.Printf("Received message: %s\n", msg.Body)
			for conn := range c.wsConn {
				log.Println("Sending message to websocket client")
				err := conn.WriteMessage(websocket.TextMessage, msg.Body)
				if err != nil {
					log.Panic("Err", err)
					conn.Close()
					delete(c.wsConn, conn)
				}
			}
		}
	}()
}
func ReadMessage(ch *amqp091.Channel, q amqp091.Queue, c *Connection, db *sql.DB) {
	for conn := range c.wsConn {
		go func(conn *websocket.Conn) {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Printf("ReadMessage Error: %v", err)
					delete(c.wsConn, conn)
					return
				}
				_, err = db.Exec("INSERT INTO messages (content, user_id) VALUES ($1, $2)", string(message), 1)
				if err != nil {
					log.Printf("Failed to insert message into PostgreSQL: %v", err)
					continue
				}
				err = rabbitmq.PublishMessage(ch, q, message)
				if err != nil {
					log.Printf("Failed to publish message: %v", err)
				} else {
					log.Printf("Message sent to RabbitMQ: %s", message)
				}
			}
		}(conn)
	}
}
func HandleConnection(ch *amqp091.Channel, q amqp091.Queue, c *websocket.Conn, db *sql.DB) {
	connMap := map[*websocket.Conn]bool{c: true}
	connection := NewConnection(connMap)

	go WriteMessage(ch, q, connection)
	go ReadMessage(ch, q, connection, db)

}
