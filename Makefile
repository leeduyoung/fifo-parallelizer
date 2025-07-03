.PHONY: test_setup
test_setup:
	docker compose up -d

.PHONY: test_down
test_down:
	docker compose down -v

AWS_ENV = AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test AWS_DEFAULT_REGION=ap-northeast-2

.PHONY: create_queue
create_queue:
	@$(AWS_ENV) aws --endpoint-url=http://localhost:4566 \
		sqs create-queue --queue-name \
		test-fifo-queue.fifo --attributes \
		'{"FifoQueue":"true","ContentBasedDeduplication":"true","MessageRetentionPeriod":"1209600","VisibilityTimeout":"30","ReceiveMessageWaitTimeSeconds":"20"}' \
		--region ap-northeast-2 

.PHONY: list_queues
list_queues:
	@$(AWS_ENV) aws --endpoint-url=http://localhost:4566 \
	sqs list-queues \
	--region ap-northeast-2 \
	--output table 2>/dev/null

.PHONY: run
run:
	@$(AWS_ENV) AWS_ENDPOINT_URL=http://localhost:4566 SQS_QUEUE_URL=http://localhost:4566/000000000000/test-fifo-queue.fifo go run cmd/main.go

.PHONY: send_same_message_group
send_same_message_group:
	@chmod +x scripts/send-same-message-group.sh
	@./scripts/send-same-message-group.sh

.PHONY: send_different_message_group
send_different_message_group:
	@chmod +x scripts/send-different-message-group.sh
	@./scripts/send-different-message-group.sh	