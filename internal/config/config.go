package config

import (
	"os"
	"strconv"

	"github.com/leeduyoung/fifo-parallelizer/internal/interfaces"
)

type Config struct {
	queueURL          string
	maxWorkers        int
	visibilityTimeout int32
	waitTimeSeconds   int32
	maxMessages       int32
	endpointURL       string
}

func NewConfig() interfaces.Config {
	return &Config{
		queueURL:          getEnvOrDefault("SQS_QUEUE_URL", "http://localhost:4566/000000000000/test-fifo-queue.fifo"),
		maxWorkers:        getEnvIntOrDefault("MAX_WORKERS", 5),
		visibilityTimeout: int32(getEnvIntOrDefault("VISIBILITY_TIMEOUT", 30)),
		waitTimeSeconds:   int32(getEnvIntOrDefault("WAIT_TIME_SECONDS", 20)),
		maxMessages:       int32(getEnvIntOrDefault("MAX_MESSAGES", 1)),
		endpointURL:       getEnvOrDefault("ENDPOINT_URL", "http://localhost:4566"),
	}
}

// MaxMessages implements interfaces.Config.
func (c *Config) MaxMessages() int32 {
	return c.maxMessages
}

// MaxWorkers implements interfaces.Config.
func (c *Config) MaxWorkers() int {
	return c.maxWorkers
}

// QueueURL implements interfaces.Config.
func (c *Config) QueueURL() string {
	return c.queueURL
}

// VisibilityTimeout implements interfaces.Config.
func (c *Config) VisibilityTimeout() int32 {
	return c.visibilityTimeout
}

// WaitTimeSeconds implements interfaces.Config.
func (c *Config) WaitTimeSeconds() int32 {
	return c.waitTimeSeconds
}

// EndPoint implements interfaces.Config.
func (c *Config) EndPointURL() string {
	return c.endpointURL
}

func getEnvIntOrDefault(value string, defaultValue int) int {
	if value := os.Getenv(value); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
