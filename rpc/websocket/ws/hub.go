package ws

import (
	"log"
	"sync"
)

type WsHub struct {
	mu          sync.RWMutex
	Connections map[int64]map[*WsClient]bool
}

func NewHub() *WsHub {
	return &WsHub{
		Connections: make(map[int64]map[*WsClient]bool),
	}
}

func (h *WsHub) RegisterClient(c *WsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	//这是一个按照courseId分类的连接池
	if h.Connections[c.CourseID] == nil {
		h.Connections[c.CourseID] = make(map[*WsClient]bool)
	}
	h.Connections[c.CourseID][c] = true
	if c.Role == 0 {
		log.Printf("老师%d加入了课堂%d", c.UserId, c.CourseID)
	}
	if c.Role == 1 {
		log.Printf("学生%d加入了课堂%d", c.UserId, c.CourseID)
	}

}

func (h *WsHub) UnregisterClient(c *WsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.Connections[c.CourseID]; ok {
		if _, exists := clients[c]; exists {
			delete(clients, c)
			close(c.SendCh)
		}
	}
	log.Printf("学生%d退出了课堂%d", c.UserId, c.CourseID)
}

func (h *WsHub) BroadcastToCourse(courseID int64, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.Connections[courseID]; ok {
		for client := range clients {
			select {
			case client.SendCh <- msg:
			default:
				// 如果缓冲满了，剔除这个客户端
				go h.UnregisterClient(client)
			}
		}
	}
}
