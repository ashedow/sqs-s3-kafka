FROM alpine:3.8

RUN apk -U --no-cache upgrade \
 && apk -U --no-cache add ca-certificates

COPY sqs-s3-kafka /bin/sqs-s3-kafka
RUN chmod 755 /bin/sqs-s3-kafka \
 && adduser -s /bin/nologin -H -D sqs-s3-kafka sqs-s3-kafka

ENV METRICS_ADDRESS=:8080
EXPOSE 8080

USER sqs-s3-kafka
CMD sqs-s3-kafka
