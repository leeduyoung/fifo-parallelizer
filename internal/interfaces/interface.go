package interfaces

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// MessageHandler is the interface for the message handler.
type MessageHandler interface {
	Handle(ctx context.Context, message types.Message) error
}

// SQSClient is the interface for the SQS client.
type SQSClient interface {
	ReceiveMessages(ctx context.Context, maxMessages int32) ([]types.Message, error)
	DeleteMessage(ctx context.Context, receiptHandle string) error
}

// MessageProcessor is the interface for the message processor.
type MessageProcessor interface {
	ProcessMessage(ctx context.Context, workerID int) error
}

// WorkerPool is the interface for the worker pool.
type WorkerPool interface {
	Start(ctx context.Context) error
	Stop() error
}

// Config is the configuration for the application.
type Config interface {
	QueueURL() string
	MaxWorkers() int
	VisibilityTimeout() int32
	WaitTimeSeconds() int32
	MaxMessages() int32
	EndPointURL() string
}
