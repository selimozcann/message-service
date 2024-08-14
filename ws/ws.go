package ws

import (
	"log"
	"messageservice/rabbitmq"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	wsConn map[*websocket.Conn]bool
	send   chan []byte
}

func NewConnection(conn map[*websocket.Conn]bool) *Connection {
	return &Connection{wsConn: conn, send: make(chan []byte, 256)}
}

func WriteMessage(ch *amqp091.Channel, q amqp091.Queue, c *Connection) {
	for message := range c.send {
		err := rabbitmq.PublishMessage(ch, q, message)
		if err != nil {
			log.Printf("Failed to publish message:%v", err)
			continue
		}
		log.Printf("Sent message: %s\n", message)
		for conn := range c.wsConn {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Panic("Err", err)
				conn.Close()
				delete(c.wsConn, conn)
			}
		}
	}
}
func ReadMessage(c *Connection) {
	for conn := range c.wsConn {
		go func(conn *websocket.Conn) {
			defer conn.Close()
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Println("Read error:", err)
					delete(c.wsConn, conn)
					return
				}
				c.send <- message
			}
		}(conn)
	}
}
func HandleConnection(ch *amqp091.Channel, q amqp091.Queue, c *websocket.Conn) {

	connMap := map[*websocket.Conn]bool{c: true}
	connection := NewConnection(connMap)
	go WriteMessage(ch, q, connection)
	go ReadMessage(connection)
}
