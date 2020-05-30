package main

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	kafka "github.com/Shopify/sarama"
)

type AsyncProducer struct {
	producer  kafka.AsyncProducer
	errors    chan error
	successes chan *kafka.ProducerMessage
}

func NewAsyncProducer(brokers []string, config *kafka.Config) (*AsyncProducer, error) {
	producer, err := kafka.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	asyncProducer := &AsyncProducer{
		producer:  producer,
		errors:    make(chan error, 1),
		successes: make(chan *kafka.ProducerMessage, 1),
	}

	go func() {
		for {
			select {
			case err := <-asyncProducer.producer.Errors():
				asyncProducer.errors <- err
			case suc := <-asyncProducer.producer.Successes():
				asyncProducer.successes <- suc
			}
		}
	}()
	return asyncProducer, nil
}

func (k *AsyncProducer) Errors() <-chan error {
	return k.errors
}

func (k *AsyncProducer) Successes() <-chan *kafka.ProducerMessage {
	return k.successes
}

func (k *AsyncProducer) Send(topic string, message []byte) {
	k.producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: kafka.ByteEncoder(message)}
}

func (k *AsyncProducer) SendRawMsg(msg *sarama.ProducerMessage) {
	k.producer.Input() <- msg
}

func (k *AsyncProducer) Destroy() {
	k.producer.Close()
}

func asyncPublishMessage(message []byte, producer AsyncProducer) (partition int32, offset int64, err error) {
	// async send message
	producer.Send("topic_test", []byte(message))
	sendCount++
	fmt.Printf("send message count: %v,%v\n", sendCount, time.Now())
}
