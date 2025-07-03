package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func main() {

	// queueURL := os.Getenv("SQS_QUEUE_URL")
	queueURL := "http://localhost:4566/000000000000/test-fifo-queue.fifo"
	// queueURL := "http://sqs.ap-northeast-2.localhost.localstack.cloud:4566/000000000000/test-fifo-queue.fifo"
	if queueURL == "" {
		log.Fatal("SQS_QUEUE_URL is not set")
	}

	worker, err := NewSQSWorker(queueURL, 5)
	if err != nil {
		log.Fatalf("Failed to create SQS worker: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := worker.Start(ctx); err != nil {
			log.Printf("Worker stopped: %v", err)
		}
	}()

	<-sigChan
	log.Println("Received signal to stop")

	cancel()
}

type SQSWorker struct {
	client     *sqs.Client
	queueURL   string
	maxWorkers int
}

func NewSQSWorker(queueURL string, maxWorkers int) (*SQSWorker, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return &SQSWorker{
		client: sqs.NewFromConfig(cfg, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String("http://localhost:4566")
		}),
		queueURL:   queueURL,
		maxWorkers: maxWorkers,
	}, nil
}

func (w *SQSWorker) Start(ctx context.Context) error {
	var wg sync.WaitGroup

	for i := 0; i < w.maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			w.pollMessages(ctx, workerID)
		}(i)
	}

	log.Printf("Started %d workers", w.maxWorkers)

	wg.Wait()
	log.Println("All workers finished")

	return nil
}

func (w *SQSWorker) pollMessages(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopped", workerID)
			return
		default:
			if err := w.receiveAndProcessMessages(ctx, workerID); err != nil {
				log.Printf("Worker %d error: %v", workerID, err)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (w *SQSWorker) receiveAndProcessMessages(ctx context.Context, workerID int) any {
	receiveInput := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(w.queueURL),
		MaxNumberOfMessages: 1,  // 한 번에 최대 1개 메시지 가져옴
		WaitTimeSeconds:     20, // Long polling
		VisibilityTimeout:   30, // 30초 동안 다른 워커에서 안보임
		MessageSystemAttributeNames: []types.MessageSystemAttributeName{
			types.MessageSystemAttributeNameAll,
		},
		MessageAttributeNames: []string{"All"},
	}

	receiveCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result, err := w.client.ReceiveMessage(receiveCtx, receiveInput)
	if err != nil {
		log.Printf("ReceiveMessage error: %v, workerID: %d", err, workerID)
		return err
	}

	if len(result.Messages) == 0 {
		log.Printf("No messages received, workerID: %d", workerID)
		return nil
	}

	log.Printf("Received %d messages, workerID: %d", len(result.Messages), workerID)

	// 사실상 메시지를 하나씩만 들고 오기 때문에 병렬처리 필요 없음
	var wg sync.WaitGroup
	for _, message := range result.Messages {
		wg.Add(1)
		go func(msg types.Message) {
			defer wg.Done()
			w.processMessage(ctx, msg, workerID)
		}(message)
	}

	wg.Wait()
	return nil
}

func (w *SQSWorker) processMessage(ctx context.Context, msg types.Message, workerID int) {
	// log.Printf("Processing message, workerID: %d, messageID: %s", workerID, *msg.MessageId)

	if err := w.handleMessage(ctx, msg); err != nil {
		log.Printf("Error processing message, workerID: %d, messageID: %s, error: %v", workerID, *msg.MessageId, err)
		return
	}

	if err := w.deleteMessage(ctx, msg); err != nil {
		log.Printf("Error deleting message, workerID: %d, messageID: %s, error: %v", workerID, *msg.MessageId, err)
	} else {
		// log.Printf("Deleted message, workerID: %d, messageID: %s", workerID, *msg.MessageId)
	}
}

func (w *SQSWorker) handleMessage(ctx context.Context, msg types.Message) error {
	// log.Printf("Message body: %s", *msg.Body)

	var (
		messageGroupId string
		messageDedupId string
	)

	if msg.Attributes != nil {
		if groupId, exists := msg.Attributes[string(types.MessageSystemAttributeNameMessageGroupId)]; exists {
			messageGroupId = groupId
		}

		if dedupId, exists := msg.Attributes[string(types.MessageSystemAttributeNameMessageDeduplicationId)]; exists {
			messageDedupId = dedupId
		}
	}

	fmt.Printf("Successfully handled message, messageID: %s, messageGroupId: %s, messageDedupId: %s\n", *msg.MessageId, messageGroupId, messageDedupId)
	time.Sleep(2 * time.Second)

	return nil
}

func (w *SQSWorker) deleteMessage(ctx context.Context, msg types.Message) error {
	deleteInput := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(w.queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	}

	deleteCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := w.client.DeleteMessage(deleteCtx, deleteInput)
	if err != nil {
		log.Printf("Error deleting message, messageID: %s, error: %v", *msg.MessageId, err)
		return err
	}

	return nil
}
