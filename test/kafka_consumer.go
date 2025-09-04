package test

import (
	"AccountManagementSystem/internal/queue_processor"
)

func KafkaConsumer(transactionService queue_processor.IKafkaConsumer) {
	queue_processor.StartConsumer(transactionService, kafkaBroker, kafkaTopic, kafkaGroup)
}
