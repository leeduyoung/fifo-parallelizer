package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/leeduyoung/fifo-parallelizer/internal/interfaces"
	apptypes "github.com/leeduyoung/fifo-parallelizer/internal/types"
)

type MessageProcessorImpl struct {
	sqsClient      interfaces.SQSClient
	messageHandler interfaces.MessageHandler
	config         interfaces.Config
}

func NewMessageProcessor(
	sqsClient interfaces.SQSClient,
	messageHandler interfaces.MessageHandler,
	config interfaces.Config,
) interfaces.MessageProcessor {
	return &MessageProcessorImpl{
		sqsClient:      sqsClient,
		messageHandler: messageHandler,
		config:         config,
	}
}

// ProcessMessage implements interfaces.MessageProcessor.
func (m *MessageProcessorImpl) ProcessMessage(ctx context.Context, workerID int) error {

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopped", workerID)
			return ctx.Err()
		default:
			if err := m.processMessageBatch(ctx, workerID); err != nil {
				log.Printf("Worker %d error: %v", workerID, err)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

// processMessageBatch is the function that processes a batch of messages.
func (m *MessageProcessorImpl) processMessageBatch(ctx context.Context, workerID int) error {

	messages, err := m.sqsClient.ReceiveMessages(ctx, m.config.MaxMessages())
	if err != nil {
		return fmt.Errorf("failed to receive messages: %w", err)
	}

	if len(messages) == 0 {
		log.Printf("Worker %d received no messages", workerID)
		return nil
	}

	for _, message := range messages {

		result := m.processMessage(ctx, message, workerID)

		if !result.Success {
			log.Printf("Worker %d error processing message %s: %v", workerID, result.MessageID, result.Error)
		}

		if err := m.sqsClient.DeleteMessage(ctx, *message.ReceiptHandle); err != nil {
			log.Printf("Worker %d error deleting message %s: %v", workerID, result.MessageID, err)
		}
	}

	return nil
}

// processMessage is the function that processes a single message.
func (m *MessageProcessorImpl) processMessage(ctx context.Context, message types.Message, workerID int) apptypes.ProcessResult {

	log.Printf("Worker %d processing message %s", workerID, *message.MessageId)

	start := time.Now()
	messageID := *message.MessageId

	err := m.messageHandler.Handle(ctx, message)
	duration := time.Since(start)

	return apptypes.ProcessResult{
		Success:   err == nil,
		Error:     err,
		Duration:  duration,
		MessageID: messageID,
		WokerID:   workerID,
	}
}
