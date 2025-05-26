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

func (c *WsClient) ReadPump(hub *WsHub) {
	defer func() {
		hub.UnregisterClient(c)
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}
		// 可扩展处理接收到的消息（如提交答题）
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
