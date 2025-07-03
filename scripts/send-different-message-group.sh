#!/bin/bash

QUEUE_URL="http://localhost:4566/000000000000/test-fifo-queue.fifo"
AWS_ENDPOINT_URL="http://localhost:4566"

echo "🚀 서로 다른 MessageGroupId로 메시지 전송 (병렬 처리 가능)"
echo "========================================================="

# 환경 변수 설정
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=ap-northeast-2

# 각 메시지마다 다른 그룹 ID 사용 (병렬 처리 가능)
aws --endpoint-url=$AWS_ENDPOINT_URL sqs send-message-batch \
    --queue-url "$QUEUE_URL" \
    --entries '[
        {
            "Id": "msg1",
            "MessageBody": "{\"id\": 1, \"message\": \"병렬처리 메시지 1\", \"timestamp\": \"'$(date -Iseconds)'\", \"task\": \"process-order\"}",
            "MessageGroupId": "group-order-001",
            "MessageDeduplicationId": "diff-group-'$(date +%s%N)'-1",
            "MessageAttributes": {
                "Type": {"StringValue": "order", "DataType": "String"},
                "Priority": {"StringValue": "high", "DataType": "String"}
            }
        },
        {
            "Id": "msg2",
            "MessageBody": "{\"id\": 2, \"message\": \"병렬처리 메시지 2\", \"timestamp\": \"'$(date -Iseconds)'\", \"task\": \"process-payment\"}",
            "MessageGroupId": "group-payment-001",
            "MessageDeduplicationId": "diff-group-'$(date +%s%N)'-2",
            "MessageAttributes": {
                "Type": {"StringValue": "payment", "DataType": "String"},
                "Priority": {"StringValue": "high", "DataType": "String"}
            }
        },
        {
            "Id": "msg3",
            "MessageBody": "{\"id\": 3, \"message\": \"병렬처리 메시지 3\", \"timestamp\": \"'$(date -Iseconds)'\", \"task\": \"send-notification\"}",
            "MessageGroupId": "group-notification-001",
            "MessageDeduplicationId": "diff-group-'$(date +%s%N)'-3",
            "MessageAttributes": {
                "Type": {"StringValue": "notification", "DataType": "String"},
                "Priority": {"StringValue": "normal", "DataType": "String"}
            }
        },
        {
            "Id": "msg4",
            "MessageBody": "{\"id\": 4, \"message\": \"병렬처리 메시지 4\", \"timestamp\": \"'$(date -Iseconds)'\", \"task\": \"update-inventory\"}",
            "MessageGroupId": "group-inventory-001",
            "MessageDeduplicationId": "diff-group-'$(date +%s%N)'-4",
            "MessageAttributes": {
                "Type": {"StringValue": "inventory", "DataType": "String"},
                "Priority": {"StringValue": "normal", "DataType": "String"}
            }
        },
        {
            "Id": "msg5",
            "MessageBody": "{\"id\": 5, \"message\": \"병렬처리 메시지 5\", \"timestamp\": \"'$(date -Iseconds)'\", \"task\": \"generate-report\"}",
            "MessageGroupId": "group-report-001",
            "MessageDeduplicationId": "diff-group-'$(date +%s%N)'-5",
            "MessageAttributes": {
                "Type": {"StringValue": "report", "DataType": "String"},
                "Priority": {"StringValue": "low", "DataType": "String"}
            }
        }
    ]' \
    --region ap-northeast-2

if [ $? -eq 0 ]; then
    echo "✅ 서로 다른 그룹 ID로 10개 메시지 전송 완료!"
    echo "📝 각 메시지가 다른 MessageGroupId를 가집니다:"
    echo "   - group-order-001, group-payment-001, group-notification-001, ..."
    echo "🚀 이 메시지들은 병렬로 처리될 수 있습니다"
else
    echo "❌ 메시지 전송 실패"
fi