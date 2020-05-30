package main

import (
	"bufio"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func inits3Client() *s3.S3 {
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
			Region:      optionalString(awsRegion),
			Credentials: creds,
		},
		Profile:           awsProfile,
		SharedConfigState: sharedConfig,
	})
	logError(err)
	s3Client := s3.New(awsSession)
	return s3Client
}

func getS3Key(message *sqs.Message) (string, error) {
	b := *message.Body
	var body Body
	err := json.Unmarshal([]byte(b), &body)
	if err != nil {
		logError(err)
	}
	key := body.Records[0].S3.Object.Key
	return key, err
}

func readObject(s3Client *s3.S3, producer Producer, bucket string, key string) (string, error) {
	file, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	logError(err)

	logInfo("Start reading s3 object")
	scanner := bufio.NewScanner(file.Body)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		val := []byte(scanner.Text())

		publishMessage(val, producer)
	}
	if err != nil {
		return "", err
	}
	return key, nil
}
