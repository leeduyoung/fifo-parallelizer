package worker

import (
	"context"
	"log"
	"sync"

	"github.com/leeduyoung/fifo-parallelizer/internal/interfaces"
)

type WorkerPoolImpl struct {
	processor  interfaces.MessageProcessor
	maxWorkers int
	wg         sync.WaitGroup
	cancel     context.CancelFunc
}

func NewWorkerPool(processor interfaces.MessageProcessor, cfg interfaces.Config) interfaces.WorkerPool {
	return &WorkerPoolImpl{
		processor:  processor,
		maxWorkers: cfg.MaxWorkers(),
	}
}

// Start implements interfaces.WorkerPool.
func (w *WorkerPoolImpl) Start(ctx context.Context) error {
	workerCtx, cancel := context.WithCancel(ctx)
	w.cancel = cancel

	for i := 0; i < w.maxWorkers; i++ {
		w.wg.Add(1)

		go func(workerID int) {
			defer w.wg.Done()

			if err := w.processor.ProcessMessage(workerCtx, workerID); err != nil {
				log.Printf("Worker %d error: %v", workerID, err)
			}
		}(i)
	}

	log.Printf("Started %d workers", w.maxWorkers)

	w.wg.Wait()
	log.Println("All workers finished")

	return nil
}

// Stop implements interfaces.WorkerPool.
func (w *WorkerPoolImpl) Stop() error {
	if w.cancel != nil {
		w.cancel()
	}

	return nil
}
