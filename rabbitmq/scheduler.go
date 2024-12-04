package rabbitmq

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
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

func Scheduler(pools []*WorkerPool, channel *amqp091.Channel) {
	var wg sync.WaitGroup

	// Detailed Explanation
	// prefetchCount (1):

	// Only deliver 1 unacknowledged message at a time
	// Consumer must ACK current message before receiving next one
	// Helps prevent worker overload
	// prefetchSize (0):

	// No limit on total size of messages being processed
	// Usually left at 0
	// global (false):

	// Settings apply per-consumer
	// If true, would apply to all consumers on the channel

	channel.Qos(
		1,     // prefetchCount: number of messages to prefetch
		0,     // prefetchSize: size in bytes (0 means no limit)
		false, // global: applies to all consumers on the channel if true
	)

	for _, pool := range pools {
		wg.Add(1)
		go func(pool *WorkerPool) {
			defer wg.Done()

			for i := 0; i < pool.Workers; i++ {
				go Worker(pool.QueueName, pool.RateLimiter, channel)
			}
		}(pool)
	}

	//TODO: Dynamic priority Adjustment Loop

	go func() {
		for {
			time.Sleep(10 * time.Second) // Adjusts the priority of the workers every 10 seconds

			for _, pool := range pools {
				// channel.QueueInspect(pool.QueueName) //Deprecated function replaced with channel.QueueDeclarePassive
				q, err := channel.QueueDeclarePassive(pool.QueueName, false, false, false, false, nil)
				if err != nil {
					log.Printf("Failed to inspect queue %s: %v", pool.QueueName, err)
					continue
				}
				if q.Messages == 0 && pool.Workers > 0 {
					log.Printf("Queue %s is empty, reallocating workers", pool.QueueName)
					reallocateWorkers(pools, channel)
				} else {
					log.Println("Queue is not idle, scaling workers...")
					scaleWorkers(pool, channel)
				}
			}
			println("Total messages: ", getTotalMessages(getQueueStats(pools, channel)))
			println()
		}
	}()
}

func reallocateWorkers(pools []*WorkerPool, channel *amqp091.Channel) {
	queueStats := getQueueStats(pools, channel)

	totalMessages := getTotalMessages(queueStats)
	println("Total messages: ", totalMessages)

	println("Reallocating workers...")
	// Adjust worker allocation based on queue backlog and priority
	for _, pool := range pools {
		if totalMessages == 0 {
			pool.Workers = 1 // At least one worker per pool
		} else {
			// Adjust based on the proportion of messages in the queue
			// and the priority level of the pool
			priorityWeight := getPriorityWeight(pool.QueueName)
			pool.Workers = int(float64(pool.Workers) * (float64(queueStats[pool.QueueName]) / float64(totalMessages)) * priorityWeight)

			// Ensure at least one worker is always allocated to avoid starvation
			if pool.Workers < 1 {
				pool.Workers = 1
			}

			// Ensure a maximum of 10 workers per pool
			if pool.Workers > pool.MaxWorkers {
				pool.Workers = pool.MaxWorkers
			}

		}

		log.Printf("Reallocated %d workers to queue %s", pool.Workers, pool.QueueName)

	}

	println("Workers reallocated")
	println()

}

func scaleWorkers(pool *WorkerPool, channel *amqp091.Channel) {
	// q, err := channel.QueueInspect(pool.QueueName)
	q, err := channel.QueueDeclarePassive(pool.QueueName, false, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to inspect queue %s: %v", pool.QueueName, err)
		return
	}

	messagesInQueue := q.Messages
	// required_workers = messagesInQueue / 5
	newWorkers := 1 + (messagesInQueue / 5)

	if newWorkers > pool.MaxWorkers {
		newWorkers = pool.MaxWorkers
	}

	pool.Workers = newWorkers

	log.Printf("Scaled workers for queue %s to %d", pool.QueueName, pool.Workers)
}

func getTotalMessages(queueStats map[string]int) int {
	totalMessages := 0
	for _, messages := range queueStats {
		totalMessages += messages
	}
	return totalMessages

}

func getQueueStats(pools []*WorkerPool, channel *amqp091.Channel) map[string]int {
	// Gather queue statistics
	queueStats := make(map[string]int)
	for _, pool := range pools {
		// q, err := channel.QueueInspect(pool.QueueName)
		q, err := channel.QueueDeclarePassive(pool.QueueName, false, false, false, false, nil)
		if err != nil {
			log.Printf("Failed to inspect queue %s: %v", pool.QueueName, err)
			continue
		}
		queueStats[pool.QueueName] = q.Messages
	}

	return queueStats
}
func getPriorityWeight(queueName string) float64 {
	switch {
	case strings.Contains(queueName, "high"):
		return 1.5
	case strings.Contains(queueName, "medium"):
		return 1.0
	case strings.Contains(queueName, "low"):
		return 0.5
	default:
		return 1.0
	}
}
