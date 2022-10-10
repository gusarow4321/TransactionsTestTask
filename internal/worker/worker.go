package worker

import (
	"TransactionsTestTask/internal/data"
	"TransactionsTestTask/internal/pkg/queue"
	"context"
	"log"
)

type Worker struct {
	repo  data.UserRepo
	queue *queue.Queue
}

func NewWorker(repo data.UserRepo, queue *queue.Queue) *Worker {
	return &Worker{
		repo:  repo,
		queue: queue,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			txs := w.queue.GetTxs()
			for i := 0; i < len(txs); i += 2 {
				_, err := w.repo.UpdateBalance(ctx, txs[i], txs[i+1])
				if err != nil {
					log.Printf("Balance update failed: %v", err)
					return err
				}
			}
		}
	}
}
