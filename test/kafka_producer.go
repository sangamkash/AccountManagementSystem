package test

import (
	"AccountManagementSystem/internal/queue_producer"
	"AccountManagementSystem/log_color"
	"AccountManagementSystem/log_helper"
	"log"
	"log/slog"
)

func GetKafka() *queue_producer.KafkaQueue {
	slog.Info(log_color.BrightCyan("make sure Kafka and zookeeper is running in docker container"))
	slog.Info(log_helper.LogServiceInitializing("kafka producer"))
	kafkaQueue, err := queue_producer.NewKafkaQueue(kafkaBroker, kafkaTopic)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info(log_helper.LogServiceInitialized("kafka producer"))
	return kafkaQueue
}
