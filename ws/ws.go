package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Connection struct {
	wsConn map[*websocket.Conn]bool
	send   chan []byte
}

func NewConnection(conn map[*websocket.Conn]bool) *Connection {
	return &Connection{wsConn: conn, send: make(chan []byte, 256)}
}

func WriteMessage(c *Connection) {
	for message := range c.send {
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
func HandleConnection(c *websocket.Conn) {
	connMap := map[*websocket.Conn]bool{c: true}
	connection := NewConnection(connMap)
	go WriteMessage(connection)
	go ReadMessage(connection)
}
