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
		userIDStr := r.URL.Query().Get("user_id")
		roleStr := r.URL.Query().Get("role")
		if courseIDStr == "" {
			http.Error(w, "course_id required", http.StatusBadRequest)
			return
		}
		if userIDStr == "" {
			http.Error(w, "student_id required", http.StatusBadRequest)
			return
		}
		if roleStr == "" {
			http.Error(w, "student_role required", http.StatusBadRequest)
			return
		}

		courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid course_id", http.StatusBadRequest)
			return
		}
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid student_id", http.StatusBadRequest)
			return
		}
		roleID, err := strconv.ParseInt(roleStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid role", http.StatusBadRequest)
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
			UserId:   userID,
			Role:     int8(roleID),
			SendCh:   make(chan []byte, 256),
		}

		hub.RegisterClient(client)

		go client.WritePump()
		go client.ReadPump(hub)
	}
}
