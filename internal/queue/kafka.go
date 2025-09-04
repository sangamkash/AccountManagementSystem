package queue

import (
	"AccountManagementSystem/internal/models"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaQueue struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaQueue(broker, topic string) (*KafkaQueue, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return nil, err
	}
	return &KafkaQueue{producer: p, topic: topic}, nil
}

func (k *KafkaQueue) PublishMessage(msg models.TransactionMessage, key string) error {
	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &k.topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}, nil)
}
