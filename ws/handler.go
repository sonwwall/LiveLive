package ws

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewHandler(hub *WsHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		courseIDStr := r.URL.Query().Get("course_id")
		if courseIDStr == "" {
			http.Error(w, "course_id required", http.StatusBadRequest)
			return
		}
		courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid course_id", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrade error: %v", err)
			return
		}

		client := &WsClient{
			Conn:     conn,
			CourseID: courseID,
			SendCh:   make(chan []byte, 256),
		}

		hub.RegisterClient(client)

		go client.WritePump()
		go client.ReadPump(hub)
	}
}
