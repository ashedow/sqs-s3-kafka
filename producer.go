package main

import (
	"runtime"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	loadConfig()
	sqsClient := initSqsClient()
	s3Client := inits3Client()
	kafkaProducer := initKafkaProducer()

	goroutines := make(chan struct{}, 5)
	for {
		messages := fetchMessages(sqsClient)
		go func(goroutines chan<- struct{}) {
			go produceMessages(messages, sqsClient, s3Client, kafkaProducer, goroutines)
		}(goroutines)
		<-goroutines
	}
}

func produceMessages(messages Messages, sqsClient *sqs.SQS, s3Client *s3.S3, kafkaProducer Producer, goroutines chan<- struct{}) {
	forwardedMessages, skippedMessages := forwardToKafka(messages, kafkaProducer, s3Client)
	deleteMessageBatch(forwardedMessages, sqsClient)
	releaseMessageBatch(skippedMessages, sqsClient)
	runtime.GC()
	goroutines <- struct{}{}
}

func forwardToKafka(messages Messages, kafkaProducer Producer, s3Client *s3.S3) (forwarded Messages, skipped Messages) {
	forwarded = Messages{}
	skipped = Messages{}
	for _, message := range messages {
		key, err := getS3Key(message)
		if err != nil {
			logError(err)
		} else {
			logInfo("Got S3 key:", key)
		}

		result, err := readObject(s3Client, kafkaProducer, s3Bucket, key)
		if err != nil {
			logError(err)
			skipped = append(skipped, message)
		} else {
			logInfo(result)
			forwarded = append(forwarded, message)
		}
	}
	return
}
