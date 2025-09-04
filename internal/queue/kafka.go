package queue

import (
	"AccountManagementSystem/internal/models"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type KafkaQueue struct {
	writer *kafka.Writer
}

func NewKafkaQueue(broker, topic string) (*KafkaQueue, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaQueue{writer: writer}, nil
}

func (k *KafkaQueue) PublishMessage(msg models.TransactionMessage, key string) error {
	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return k.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
}
