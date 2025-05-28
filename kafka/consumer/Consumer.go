package consumer

import (
	"LiveLive/dao/db"
	"LiveLive/model"
	"LiveLive/utils"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func StartKafkaConsumerChat() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          []string{"localhost:9092"},
		Topic:            "chat_messages",
		GroupID:          "chat-consumer-group",
		MinBytes:         10e3,             // 10KB
		MaxBytes:         10e6,             // 10MB
		ReadBatchTimeout: 10 * time.Second, //十秒后自动读取
	})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("读取 Kafka 消息失败:", err)
			continue
		}
		type ChatMessage struct {
			UserID   string `json:"user_id"`
			CourseID string `json:"course_id"`
			Content  string `json:"content"`
		}

		var chat ChatMessage
		if err := json.Unmarshal(msg.Value, &chat); err != nil {
			log.Println("解析 Kafka 聊天消息失败:", err)
			continue
		}

		chatMessageRecord := &model.ChatMsgRecord{
			UserID:   utils.StringToInt64(chat.UserID),
			CourseID: utils.StringToInt64(chat.CourseID),
			Content:  chat.Content,
		}

		// 写入数据库
		err = db.AddChatMessageRecord(chatMessageRecord)
		if err != nil {
			log.Println("写入数据库失败:", err)
		}
	}
}
