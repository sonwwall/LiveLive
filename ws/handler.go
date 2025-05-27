package ws

import (
	dao "LiveLive/dao/rdb"
	"LiveLive/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WsMessage struct {
	Type string          `json:"type"` // "answer", "answer_result"
	Data json.RawMessage `json:"data"`
}

type AnswerMessage struct {
	QuestionID string `json:"question_id"`
	StudentID  string `json:"student_id"`
	Answer     string `json:"answer"`
	CourseID   string `json:"course_id"`
}

type ChatMessage struct {
	UserID   string `json:"user_id"`
	CourseID string `json:"course_id"`
	Content  string `json:"content"`
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

func HandleMessage(msg []byte, client *WsClient, hub *WsHub) {
	var wsMsg WsMessage
	if err := json.Unmarshal(msg, &wsMsg); err != nil {
		log.Println("解析消息失败:", err)
		return
	}

	switch wsMsg.Type {
	case "choice_answer":
		MessageChoiceQuestion(msg, client, hub, &wsMsg)
	case "chat":
		log.Println(wsMsg.Type, wsMsg.Data)
		MessageChat(msg, client, hub, &wsMsg)
	default:
		log.Println("未知类型:", wsMsg.Type)
	}
}

func MessageChoiceQuestion(msg []byte, client *WsClient, hub *WsHub, wsMsg *WsMessage) {
	var a AnswerMessage
	if err := json.Unmarshal(wsMsg.Data, &a); err != nil {
		log.Println("解析答题数据失败:", err)
		return
	}
	key := fmt.Sprintf("choice_answer:%s", a.QuestionID)
	//检查是否提交过答案
	exists, _ := dao.Redis.HExists(context.Background(), key, a.StudentID).Result()

	if exists {

		for client := range hub.Connections[utils.StringToInt64(a.CourseID)] {
			if client.UserId == utils.StringToInt64(a.StudentID) {
				client.SendCh <- []byte("您已提交过答案，无法再次提交")
				return
			}
		}
	}
	//设置redis缓存 hash类型
	err := dao.Redis.HSet(context.Background(), key, a.StudentID, a.Answer).Err()
	if err != nil {
		log.Println("写入 Redis 失败:", err)
	}
}

func MessageChat(msg []byte, client *WsClient, hub *WsHub, wsMsg *WsMessage) {
	var a ChatMessage
	if err := json.Unmarshal(wsMsg.Data, &a); err != nil {
		log.Println("聊天消息解析失败:", err)
		return
	}
	log.Println("内容:", a.Content)
	resultMsg := map[string]interface{}{
		"type": "chat_message",
		"data": map[string]interface{}{
			"user_id": utils.StringToInt64(a.UserID),
			"content": a.Content,
		},
	}
	log.Println("<UNK>:", resultMsg)
	payload, _ := json.Marshal(resultMsg)

	for client := range hub.Connections[utils.StringToInt64(a.CourseID)] {
		client.SendCh <- payload
	}
}
