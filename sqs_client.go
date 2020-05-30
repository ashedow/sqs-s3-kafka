package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Messages []*sqs.Message

func GetQueueURL(sqsClient *sqs.SQS, QueueName string) (*string, error) {
	queue, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(QueueName)})
	return queue.QueueUrl, err
}

func initSqsClient() *sqs.SQS {
	creds := (*credentials.Credentials)(nil)
	sharedConfig := session.SharedConfigStateFromEnv
	if awsReadConfig {
		sharedConfig = session.SharedConfigEnable
	}
	if awsAccessKey != "" || awsSecretKey != "" {
		creds = credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, "")
	}
	awsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint:    optionalString(awsEndpoint),
			Region:      optionalString(awsRegion),
			Credentials: creds,
		},
		Profile:           awsProfile,
		SharedConfigState: sharedConfig,
	})
	logError(err)
	sqsClient := sqs.New(awsSession)
	QueueURL, _ := GetQueueURL(sqsClient, sqsQueueName)
	_, err = sqsClient.GetQueueAttributes(&sqs.GetQueueAttributesInput{QueueUrl: QueueURL})
	logError(err)
	return sqsClient
}

func fetchMessages(sqsClient *sqs.SQS) Messages {
	QueueURL, _ := GetQueueURL(sqsClient, sqsQueueName)
	response, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            QueueURL,
		MaxNumberOfMessages: aws.Int64(10),
	})
	if err != nil {
		logError(err)
		return Messages{}
	}
	return response.Messages
}

func deleteMessageBatch(messages Messages, sqsClient *sqs.SQS) {
	if len(messages) == 0 {
		return
	}
	toDelete := []*sqs.DeleteMessageBatchRequestEntry{}
	for _, deleted := range messages {
		toDelete = append(toDelete, &sqs.DeleteMessageBatchRequestEntry{
			Id:            deleted.MessageId,
			ReceiptHandle: deleted.ReceiptHandle,
		})
	}
	QueueURL, _ := GetQueueURL(sqsClient, sqsQueueName)
	_, err := sqsClient.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
		QueueUrl: QueueURL,
		Entries:  toDelete,
	})
	logError(err)
}

func releaseMessageBatch(messages Messages, sqsClient *sqs.SQS) {
	if len(messages) == 0 {
		return
	}
	toRelease := []*sqs.ChangeMessageVisibilityBatchRequestEntry{}
	for _, released := range messages {
		toRelease = append(toRelease, &sqs.ChangeMessageVisibilityBatchRequestEntry{
			Id:                released.MessageId,
			ReceiptHandle:     released.ReceiptHandle,
			VisibilityTimeout: aws.Int64(0),
		})
	}
	QueueURL, _ := GetQueueURL(sqsClient, sqsQueueName)
	_, err := sqsClient.ChangeMessageVisibilityBatch(&sqs.ChangeMessageVisibilityBatchInput{
		QueueUrl: QueueURL,
		Entries:  toRelease,
	})
	logError(err)
}

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
