package ws

import (
	"github.com/gorilla/websocket"
	"log"
)

type WsClient struct {
	Conn     *websocket.Conn
	CourseID int64
	UserId   int64
	Role     int8 //0是老师，1是学生
	SendCh   chan []byte
}

// ReadPump 读取输入的消息
func (c *WsClient) ReadPump(hub *WsHub) {
	defer func() {
		hub.UnregisterClient(c)
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}
		HandleMessage(msg, c, hub)
	}
}

func (c *WsClient) WritePump() {
	for msg := range c.SendCh {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("write error: %v", err)
			break
		}
	}
	c.Conn.Close()
}
