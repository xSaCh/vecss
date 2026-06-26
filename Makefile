PORT=9090

S3MOCK_CONTAINER_NAME=vecss_s3mock
MQ_CONTAINER_NAME=vecss_mq

BUCKET_NAME=bkt

run_s3mock:
	# docker rm -f $(shell docker ps -aq --filter "name=$(S3MOCK_CONTAINER_NAME)")
	docker run -d -p $(PORT):9090 -p 9191:9191 -e initialBuckets=$(BUCKET_NAME) --name $(S3MOCK_CONTAINER_NAME) adobe/s3mock

run_mq:
	docker run -d -p 5672:5672 -p 15672:15672 --name $(MQ_CONTAINER_NAME) rabbitmq:management-alpine

start_s3mock:
	docker start $(S3MOCK_CONTAINER_NAME)

start_mq:
	docker start $(MQ_CONTAINER_NAME)

.PHONY: run_s3mock run_mq start_s3mock start_mq

