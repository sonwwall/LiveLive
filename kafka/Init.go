package kafka

import "github.com/segmentio/kafka-go"

func InitKafkaWriter(topic string) *kafka.Writer {
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return kafkaWriter
}
