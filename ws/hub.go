package ws

import (
	"sync"
)

type WsHub struct {
	mu          sync.RWMutex
	connections map[int64]map[*WsClient]bool
}

func NewHub() *WsHub {
	return &WsHub{
		connections: make(map[int64]map[*WsClient]bool),
	}
}

func (h *WsHub) RegisterClient(c *WsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connections[c.CourseID] == nil {
		h.connections[c.CourseID] = make(map[*WsClient]bool)
	}
	h.connections[c.CourseID][c] = true
}

func (h *WsHub) UnregisterClient(c *WsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.connections[c.CourseID]; ok {
		if _, exists := clients[c]; exists {
			delete(clients, c)
			close(c.SendCh)
		}
	}
}

func (h *WsHub) BroadcastToCourse(courseID int64, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.connections[courseID]; ok {
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
