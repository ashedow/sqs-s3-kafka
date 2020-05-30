# SQS S3 Kafka produser

Triggered by SQS topic, getting s3 key from message, read file from s3 bucket and create new message in Kafka topic for each line from logs.

## Usage

### Build

For build run:
```
  go get -d -t .
  go build
```

The resulting binary - `sqs-s3-kafka`

### Run

`sqs-s3-kafka` takes the following commandline parameters:

  - `--aws-access-key=STRING`: AWS access key
  - `--aws-secret-key=STRING`: AWS secret key
  - `--aws-profile=STRING`: AWS profile
  - `--aws-region=STRING`: AWS region
  - `--aws-read-config`: read AWS configuration from `~/.aws/config`
  - `--sqs-query-name=STRING`: SQS queue name for incomming messages
  - `--kafka-topic=STRING`: Kafka topic for outgoing messages
  - `--kafka-brokers=STRING`: list of Kafka brokers used for bootstrapping
  - `--aws-s3-bucket=STRING`: s3 bucket for getting log file

`sqs-query-name`, `kafka-brokers`, `kafka-topic` and `aws-s3-bucket` are required. Parameters may set via environment
variables as well

  - `AWS_ACCESS_KEY_ID`: AWS access key
  - `AWS_SECRET_ACCESS_KEY`: AWS secret key
  - `AWS_PROFILE`: AWS profile
  - `AWS_REGION`: AWS region
  - `S3_BUCKET`: s3 bucket for getting log file
  - `SQS_QUEUE_NAME`: SQS queue name for incomming messages
  - `KAFKA_BROKERS`: list of Kafka brokers used for bootstrapping
  - `KAFKA_TOPIC`: Kafka topic for outgoing messages
  - `AWS_SDK_LOAD_CONFIG`: if set to `1` read AWS configuration from `~/.aws/config`

When both are specified, commandline parameters take
precedence over environment variables.

### Docker

```
docker run -P sqs-s3-kafka <additional commandline arguments>
```

