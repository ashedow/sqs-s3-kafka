package main

import (
	"strings"

	kafka "github.com/Shopify/sarama"
)

type Producer kafka.SyncProducer

func initKafkaProducer() Producer {
	config := kafka.NewConfig()
	config.Producer.RequiredAcks = kafka.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	brokers := strings.Split(kafkaBrokers, ",")
	producer, err := kafka.NewSyncProducer(brokers, config)
	logError(err)
	return producer
}

func publishMessage(message []byte, producer Producer) (partition int32, offset int64, err error) {
	partition, offset, err = producer.SendMessage(&kafka.ProducerMessage{
		Topic: kafkaTopic,
		Value: kafka.ByteEncoder(message),
	})
	return
}
