package types

import "time"

// MessageMetadata is the metadata of a message.
type MessageMetadata struct {
	MessageID               string
	MessageGroupID          string
	DeduplicationID         string
	ReceiptHandle           string
	ApproximateReceiveCount int
}

type WorkerStats struct {
	WokerID         int
	ProcessedCount  int64
	ErrorCount      int64
	LastProcessedAt time.Time
}

type ProcessResult struct {
	Success   bool
	Error     error
	Duration  time.Duration
	MessageID string
	WokerID   int
}
