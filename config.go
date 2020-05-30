package main

import (
	"flag"
	"os"
)

var (
	awsProfile    = ""
	awsAccessKey  = ""
	awsSecretKey  = ""
	awsRegion     = ""
	sqsQueueName  = ""
	awsEndpoint   = ""
	kafkaBrokers  = ""
	kafkaTopic    = ""
	s3Bucket      = ""
	awsReadConfig = false
)

func loadConfig() {
	flag.StringVar(&awsAccessKey, "aws-access-key", "", "AWS access key")
	flag.StringVar(&awsSecretKey, "aws-secret-key", "", "AWS secret key")
	flag.StringVar(&awsRegion, "aws-region", os.Getenv("AWS_REGION"), "AWS region")
	flag.StringVar(&awsProfile, "aws-profile", os.Getenv("AWS_PROFILE"), "AWS profile")
	flag.BoolVar(&awsReadConfig, "aws-read-config", false, "read AWS configuration from `~/.aws/config`")
	flag.StringVar(&sqsQueueName, "sqs-queue-name", os.Getenv("SQS_QUEUE_NAME"), "SQS queue name for incomming messages")
	flag.StringVar(&kafkaBrokers, "kafka-brokers", os.Getenv("KAFKA_BROKERS"), "list of Kafka brokers used for bootstrapping")
	flag.StringVar(&kafkaTopic, "kafka-topic", os.Getenv("KAFKA_TOPIC"), "Kafka topic for outgoing messages")
	flag.StringVar(&s3Bucket, "s3-bucket", os.Getenv("S3_BUCKET"), "Listening address to serve metrics")
	flag.Parse()

	if sqsQueueName == "" {
		panic("Required parameter `SQS_QUEUE_NAME` is missing.")
	}
	if kafkaBrokers == "" {
		panic("Required parameter `KAFKA_BROKERS` is missing.")
	}
	if kafkaTopic == "" {
		panic("Required parameter `KAFKA_TOPIC` is missing.")
	}
	if awsRegion == "" {
		panic("Required parameter `AWS_REGION` is missing.")
	}
	if s3Bucket == "" {
		panic("Required parameter `S3_BUCKET` is missing.")
	}

	logConfig()
}

func logConfig() {
	parameters := []interface{}{}
	appendIfDefined := func(name string, value string) {
		if value != "" {
			parameters = append(parameters, name, value)
		}
	}
	appendIfDefined("awsProfile", awsProfile)
	appendIfDefined("awsRegion", awsRegion)
	appendIfDefined("s3Bucket", s3Bucket)
	appendIfDefined("sqsQueueName", sqsQueueName)
	appendIfDefined("kafkaBrokers", kafkaBrokers)
	appendIfDefined("kafkaTopic", kafkaTopic)
	logInfo("starting with configuration", parameters...)
}
