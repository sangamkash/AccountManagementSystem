package queue_processor

import (
	"AccountManagementSystem/internal/models"
	"AccountManagementSystem/internal/services"
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

func StartConsumer(service *services.TransactionService, brokers, topic, group string) {
	go func() {
		c, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": brokers,
			"group.id":          group,
			"auto.offset.reset": "earliest",
		})
		if err != nil {
			log.Fatal("Kafka consumer error:", err)
		}
		defer c.Close()

		if err := c.Subscribe(topic, nil); err != nil {
			log.Fatal("Subscribe failed:", err)
		}

		for {
			msg, err := c.ReadMessage(-1)
			if err != nil {
				log.Printf("consumer error: %v\n", err)
				continue
			}

			var tmsg models.TransactionMessage
			if err := json.Unmarshal(msg.Value, &tmsg); err != nil {
				log.Printf("json decode error: %v\n", err)
				continue
			}

			if err := service.ProcessMessage(context.Background(), tmsg); err != nil {
				log.Printf("failed to process transaction: %v", err)
				// In production: retry or send to dead-letter
			}
		}
	}()
}
