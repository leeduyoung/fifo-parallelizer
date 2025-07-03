package client

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/leeduyoung/fifo-parallelizer/internal/interfaces"
)

type SQSClientImpl struct {
	client *sqs.Client
	config interfaces.Config
}

func NewSQSClient(cfg interfaces.Config) (interfaces.SQSClient, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	client := sqs.NewFromConfig(awsConfig, func(o *sqs.Options) {
		if endpoint := cfg.EndPointURL(); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	return &SQSClientImpl{
		client: client,
		config: cfg,
	}, nil
}

// DeleteMessage implements interfaces.SQSClient.
func (s *SQSClientImpl) DeleteMessage(ctx context.Context, receiptHandle string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.config.QueueURL()),
		ReceiptHandle: aws.String(receiptHandle),
	}

	deleteCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.client.DeleteMessage(deleteCtx, input)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

// ReceiveMessages implements interfaces.SQSClient.
func (s *SQSClientImpl) ReceiveMessages(ctx context.Context, maxMessages int32) ([]types.Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.config.QueueURL()),
		MaxNumberOfMessages: maxMessages,
		VisibilityTimeout:   s.config.VisibilityTimeout(),
		WaitTimeSeconds:     s.config.WaitTimeSeconds(),
		MessageSystemAttributeNames: []types.MessageSystemAttributeName{
			types.MessageSystemAttributeNameAll,
		},
		MessageAttributeNames: []string{
			"All",
		},
	}

	receiveCtx, cancel := context.WithTimeout(ctx, 35*time.Second)
	defer cancel()

	output, err := s.client.ReceiveMessage(receiveCtx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to receive messages: %w", err)
	}

	return output.Messages, nil
}
