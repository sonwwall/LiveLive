package ws

import (
	dao "LiveLive/dao/rdb"
	"LiveLive/model"
	"LiveLive/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type RegisterMessage struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	CourseID string `json:"course_id"`
}

func NewHandler(hub *WsHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		courseIDStr := r.URL.Query().Get("course_id")
		token := r.Header.Get("Authorization")

		if courseIDStr == "" {
			http.Error(w, "course_id required", http.StatusBadRequest)
			return
		}

		if token == "" {
			http.Error(w, "token required", http.StatusBadRequest)
		}

		courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid course_id", http.StatusBadRequest)
			return
		}

		user, err := ParseJWTFromHeader(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusBadRequest)
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrade error: %v", err)
			return
		}

		client := &WsClient{
			Conn:     conn,
			CourseID: courseID,
			UserId:   int64(user.ID),
			Role:     int8(user.Role - 1),
			SendCh:   make(chan []byte, 256),
		}

		hub.RegisterClient(client)

		go client.WritePump()
		go client.ReadPump(hub)
	}
}

var jwtSecretKey = []byte("secret key-sonwwall") // 与 middleware.InitJwt 中保持一致
var identityKey = "identity-sonwwall"            // 同 middleware.IdentityKey

func ParseJWTFromHeader(authHeader string) (*model.User, error) {
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, errors.New("invalid Authorization format")
	}
	tokenStr := parts[1]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	user := &model.User{
		Username: claims[identityKey].(string),
		Role:     int32(claims["role"].(float64)),
		Model: gorm.Model{
			ID: uint(claims["id"].(float64)),
		},
	}

	return user, nil
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
	case "true_or_false_answer":
		MessageTrueOrFalseQuestion(msg, client, hub, &wsMsg)
	case "chat":

		MessageChat(msg, client, hub, &wsMsg)
	case "register":
		MessageRegister(msg, client, hub, &wsMsg)

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
	kafkaWriter := kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "chat_messages",
		Balancer: &kafka.LeastBytes{},
	}
	//开一个进程写入kafka
	go func() {
		err := kafkaWriter.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(a.CourseID),
				Value: wsMsg.Data,
				Time:  time.Now(),
			})
		if err != nil {
			log.Println("写入 Kafka 失败:", err)
		}
	}()
	resultMsg := map[string]interface{}{
		"type": "chat_message",
		"data": map[string]interface{}{
			"user_id": utils.StringToInt64(a.UserID),
			"content": a.Content,
		},
	}
	payload, _ := json.Marshal(resultMsg)

	for client := range hub.Connections[utils.StringToInt64(a.CourseID)] {
		client.SendCh <- payload
	}
}

func MessageRegister(msg []byte, client *WsClient, hub *WsHub, wsMsg *WsMessage) {
	var a RegisterMessage
	if err := json.Unmarshal(wsMsg.Data, &a); err != nil {
		log.Println("签到消息解析失败:", err.Error())
		return
	}
	key := fmt.Sprintf("sign:course:%s", a.CourseID)
	//修改状态为1，表示已签到
	dao.Redis.HSet(context.Background(), key, a.Username, 1)

}

func MessageTrueOrFalseQuestion(msg []byte, client *WsClient, hub *WsHub, wsMsg *WsMessage) {
	var a AnswerMessage
	if err := json.Unmarshal(wsMsg.Data, &a); err != nil {
		log.Println("解析答题数据失败:", err)
		return
	}
	key := fmt.Sprintf("true_or_false_answer:%s", a.QuestionID)
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
