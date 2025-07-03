package handler

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/leeduyoung/fifo-parallelizer/internal/interfaces"
	apptypes "github.com/leeduyoung/fifo-parallelizer/internal/types"
)

type DefaultMessageHandler struct {
	processingDelay time.Duration
}

// Handle implements interfaces.MessageHandler.
func (d *DefaultMessageHandler) Handle(ctx context.Context, message types.Message) error {

	metadata := d.extractMetadata(message)

	// TODO: Do something with the message
	time.Sleep(3 * time.Second)

	log.Printf("Successfully processed message: MessageID=%s, MessageGroupID=%s, DeduplicationID=%s, ApproximateReceiveCount=%d",
		metadata.MessageID,
		metadata.MessageGroupID,
		metadata.DeduplicationID,
		metadata.ApproximateReceiveCount,
	)

	return nil
}

func (d *DefaultMessageHandler) extractMetadata(message types.Message) apptypes.MessageMetadata {

	metadata := apptypes.MessageMetadata{
		MessageID:     *message.MessageId,
		ReceiptHandle: *message.ReceiptHandle,
	}

	if message.Attributes != nil {
		if groupID, exists := message.Attributes[string(types.MessageSystemAttributeNameMessageGroupId)]; exists {
			metadata.MessageGroupID = groupID
		}

		if deduplicationID, exists := message.Attributes[string(types.MessageSystemAttributeNameMessageDeduplicationId)]; exists {
			metadata.DeduplicationID = deduplicationID
		}

		if approximateReceiveCount, exists := message.Attributes[string(types.MessageSystemAttributeNameApproximateReceiveCount)]; exists {
			count, err := strconv.Atoi(approximateReceiveCount)
			if err == nil {
				metadata.ApproximateReceiveCount = count
			}
		}
	}

	return metadata
}

func NewDefaultMessageHandler(processingDelay time.Duration) interfaces.MessageHandler {
	return &DefaultMessageHandler{
		processingDelay: processingDelay,
	}
}
