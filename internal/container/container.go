package container

import (
	"fmt"
	"time"

	"github.com/leeduyoung/fifo-parallelizer/internal/client"
	"github.com/leeduyoung/fifo-parallelizer/internal/config"
	"github.com/leeduyoung/fifo-parallelizer/internal/handler"
	"github.com/leeduyoung/fifo-parallelizer/internal/interfaces"
	"github.com/leeduyoung/fifo-parallelizer/internal/worker"
)

type Container struct {
	config         interfaces.Config
	sqsClient      interfaces.SQSClient
	messageHandler interfaces.MessageHandler
	processor      interfaces.MessageProcessor
	workerPool     interfaces.WorkerPool
}

func NewContainer() (*Container, error) {
	container := new(Container)

	// SET CONFIG
	container.config = config.NewConfig()

	// SET SQS CLIENT
	sqsClient, err := client.NewSQSClient(container.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqs client: %w", err)
	}
	container.sqsClient = sqsClient

	// SET MESSAGE HANDLER
	container.messageHandler = handler.NewDefaultMessageHandler(2 * time.Second)

	// SET PROCESSOR
	container.processor = worker.NewMessageProcessor(
		container.sqsClient,
		container.messageHandler,
		container.config,
	)

	// SET WORKER POOL
	container.workerPool = worker.NewWorkerPool(
		container.processor,
		container.config,
	)
	return container, nil
}

func (c *Container) WorkerPool() interfaces.WorkerPool {
	return c.workerPool
}

func (c *Container) Config() interfaces.Config {
	return c.config
}

func (c *Container) WithMessageHandler(handler interfaces.MessageHandler) *Container {
	c.messageHandler = handler

	// UPDATE PROCESSOR
	c.processor = worker.NewMessageProcessor(
		c.sqsClient,
		handler,
		c.config,
	)

	// UPDATE WORKER POOL
	c.workerPool = worker.NewWorkerPool(
		c.processor,
		c.config,
	)

	return c
}
