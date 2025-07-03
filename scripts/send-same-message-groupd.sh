#!/bin/bash

QUEUE_URL="http://localhost:4566/000000000000/test-fifo-queue.fifo"
AWS_ENDPOINT_URL="http://localhost:4566"

echo "동일한 MessageGroupId로 메시지 전송 (순서 보장)"
echo "================================================"

# 환경 변수 설정
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=ap-northeast-2

# 모든 메시지가 같은 그룹 ID를 사용 (순서 보장)
SINGLE_GROUP_ID="single-group"

aws --endpoint-url=$AWS_ENDPOINT_URL sqs send-message-batch \
    --queue-url "$QUEUE_URL" \
    --entries '[
        {
            "Id": "msg1",
            "MessageBody": "{\"id\": 1, \"message\": \"순서보장 메시지 1 ('$(date +%s%N)') \", \"timestamp\": \"'$(date -Iseconds)'\", \"order\": \"first\"}",
            "MessageGroupId": "'$SINGLE_GROUP_ID'",
            "MessageDeduplicationId": "same-group-'$(date +%s%N)'-1",
            "MessageAttributes": {
                "Type": {"StringValue": "order", "DataType": "String"},
                "Sequence": {"StringValue": "1", "DataType": "Number"}
            }
        },
        {
            "Id": "msg2",
            "MessageBody": "{\"id\": 2, \"message\": \"순서보장 메시지 2 ('$(date +%s%N)') \", \"timestamp\": \"'$(date -Iseconds)'\", \"order\": \"second\"}",
            "MessageGroupId": "'$SINGLE_GROUP_ID'",
            "MessageDeduplicationId": "same-group-'$(date +%s%N)'-2",
            "MessageAttributes": {
                "Type": {"StringValue": "order", "DataType": "String"},
                "Sequence": {"StringValue": "2", "DataType": "Number"}
            }
        },
        {
            "Id": "msg3",
            "MessageBody": "{\"id\": 3, \"message\": \"순서보장 메시지 3 ('$(date +%s%N)') \", \"timestamp\": \"'$(date -Iseconds)'\", \"order\": \"third\"}",
            "MessageGroupId": "'$SINGLE_GROUP_ID'",
            "MessageDeduplicationId": "same-group-'$(date +%s%N)'-3",
            "MessageAttributes": {
                "Type": {"StringValue": "order", "DataType": "String"},
                "Sequence": {"StringValue": "3", "DataType": "Number"}
            }
        },
        {
            "Id": "msg4",
            "MessageBody": "{\"id\": 4, \"message\": \"순서보장 메시지 4 ('$(date +%s%N)') \", \"timestamp\": \"'$(date -Iseconds)'\", \"order\": \"fourth\"}",
            "MessageGroupId": "'$SINGLE_GROUP_ID'",
            "MessageDeduplicationId": "same-group-'$(date +%s%N)'-4",
            "MessageAttributes": {
                "Type": {"StringValue": "order", "DataType": "String"},
                "Sequence": {"StringValue": "4", "DataType": "Number"}
            }
        },
        {
            "Id": "msg5",
            "MessageBody": "{\"id\": 5, \"message\": \"순서보장 메시지 5 ('$(date +%s%N)') \", \"timestamp\": \"'$(date -Iseconds)'\", \"order\": \"fifth\"}",
            "MessageGroupId": "'$SINGLE_GROUP_ID'",
            "MessageDeduplicationId": "same-group-'$(date +%s%N)'-5",
            "MessageAttributes": {
                "Type": {"StringValue": "order", "DataType": "String"},
                "Sequence": {"StringValue": "5", "DataType": "Number"}
            }
        }
    ]' \
    --region ap-northeast-2

if [ $? -eq 0 ]; then
    echo "✅ 동일한 그룹 ID로 10개 메시지 전송 완료!"
    echo "📝 MessageGroupId: $SINGLE_GROUP_ID"
    echo "🔄 이 메시지들은 순서대로 처리됩니다 (FIFO 순서 보장)"
else
    echo "❌ 메시지 전송 실패"
fi