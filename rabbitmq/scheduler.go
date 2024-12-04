package rabbitmq

import (
	"golang.org/x/time/rate"
)

type WorkerPool struct {
	QueueName   string
	RateLimiter *rate.Limiter
	Workers     int
	MaxWorkers  int
}

func NewWorkerPool(queueName string, rateLimiter *rate.Limiter, workers, maxworkers int) *WorkerPool {
	return &WorkerPool{
		QueueName:   queueName,
		RateLimiter: rateLimiter,
		Workers:     workers,
		MaxWorkers:  maxworkers,
	}
}
