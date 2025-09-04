package queue_processor

import (
	"AccountManagementSystem/internal/models"
	"AccountManagementSystem/internal/services"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

func StartConsumer(service *services.TransactionService, brokers, topic, group string) {
	go func() {
		// Initialize Kafka reader
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{brokers},
			Topic:    topic,
			GroupID:  group,
			MinBytes: 1,    // Fetch messages as soon as available
			MaxBytes: 10e6, // 10MB max per fetch
		})
		defer reader.Close()

		for {
			// Read message from Kafka
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("consumer error: %v\n", err)
				continue
			}

			// Deserialize the message
			var tmsg models.TransactionMessage
			if err := json.Unmarshal(msg.Value, &tmsg); err != nil {
				log.Printf("json decode error: %v\n", err)
				continue
			}

			// Process the message
			if err := service.ProcessMessage(context.Background(), tmsg); err != nil {
				log.Printf("failed to process transaction: %v", err)
				// In production: retry or send to dead-letter queue
			}
		}
	}()
}
