package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Job represents a job to be processed
type Job interface {
	Handle() error
	GetName() string
	GetPayload() map[string]interface{}
	SetPayload(map[string]interface{})
}

// BaseJob provides basic job functionality
type BaseJob struct {
	Name    string                 `json:"name"`
	Payload map[string]interface{} `json:"payload"`
}

func (bj *BaseJob) Handle() error {
	// Default implementation - should be overridden
	return nil
}

func (bj *BaseJob) GetName() string {
	return bj.Name
}

func (bj *BaseJob) GetPayload() map[string]interface{} {
	return bj.Payload
}

func (bj *BaseJob) SetPayload(payload map[string]interface{}) {
	bj.Payload = payload
}

// Queue interface defines queue operations
type Queue interface {
	Push(job Job, delay ...time.Duration) error
	Pop() (Job, error)
	Size() (int64, error)
	Clear() error
}

// QueueManager manages job queues
type QueueManager struct {
	queues   map[string]Queue
	workers  map[string]*Worker
	handlers map[string]func() Job
	default_ string
	mutex    sync.RWMutex
}

// NewQueueManager creates a new queue manager
func NewQueueManager() *QueueManager {
	return &QueueManager{
		queues:   make(map[string]Queue),
		workers:  make(map[string]*Worker),
		handlers: make(map[string]func() Job),
	}
}

// AddQueue adds a queue
func (qm *QueueManager) AddQueue(name string, queue Queue) {
	qm.mutex.Lock()
	defer qm.mutex.Unlock()
	
	qm.queues[name] = queue
	if qm.default_ == "" {
		qm.default_ = name
	}
}

// Queue returns a queue by name
func (qm *QueueManager) Queue(name ...string) Queue {
	qm.mutex.RLock()
	defer qm.mutex.RUnlock()
	
	queueName := qm.default_
	if len(name) > 0 {
		queueName = name[0]
	}
	return qm.queues[queueName]
}

// RegisterJob registers a job handler
func (qm *QueueManager) RegisterJob(name string, handler func() Job) {
	qm.mutex.Lock()
	defer qm.mutex.Unlock()
	qm.handlers[name] = handler
}

// Dispatch dispatches a job to queue
func (qm *QueueManager) Dispatch(job Job, queueName ...string) error {
	queue := qm.Queue(queueName...)
	return queue.Push(job)
}

// StartWorker starts a worker for a queue
func (qm *QueueManager) StartWorker(queueName string, concurrency int) {
	qm.mutex.Lock()
	defer qm.mutex.Unlock()
	
	if _, exists := qm.workers[queueName]; exists {
		return // Worker already running
	}
	
	worker := NewWorker(qm.queues[queueName], qm.handlers, concurrency)
	qm.workers[queueName] = worker
	go worker.Start()
	
	log.Printf("🚀 Started queue worker for '%s' with %d workers", queueName, concurrency)
}

// StopWorker stops a worker
func (qm *QueueManager) StopWorker(queueName string) {
	qm.mutex.Lock()
	defer qm.mutex.Unlock()
	
	if worker, exists := qm.workers[queueName]; exists {
		worker.Stop()
		delete(qm.workers, queueName)
		log.Printf("🛑 Stopped queue worker for '%s'", queueName)
	}
}

// RedisQueue implements Redis-based queue
type RedisQueue struct {
	client    *redis.Client
	queueName string
}

// NewRedisQueue creates a new Redis queue
func NewRedisQueue(client *redis.Client, queueName string) *RedisQueue {
	return &RedisQueue{
		client:    client,
		queueName: queueName,
	}
}

func (rq *RedisQueue) Push(job Job, delay ...time.Duration) error {
	ctx := context.Background()
	
	jobData := map[string]interface{}{
		"name":    job.GetName(),
		"payload": job.GetPayload(),
	}
	
	data, err := json.Marshal(jobData)
	if err != nil {
		return err
	}
	
	if len(delay) > 0 && delay[0] > 0 {
		// Delayed job
		score := float64(time.Now().Add(delay[0]).Unix())
		return rq.client.ZAdd(ctx, rq.queueName+":delayed", redis.Z{
			Score:  score,
			Member: data,
		}).Err()
	}
	
	return rq.client.LPush(ctx, rq.queueName, data).Err()
}

func (rq *RedisQueue) Pop() (Job, error) {
	ctx := context.Background()
	
	// First check for delayed jobs that are ready
	rq.processDelayedJobs()
	
	result, err := rq.client.BRPop(ctx, 1*time.Second, rq.queueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No job available
		}
		return nil, err
	}
	
	var jobData map[string]interface{}
	if err := json.Unmarshal([]byte(result[1]), &jobData); err != nil {
		return nil, err
	}
	
	job := &BaseJob{
		Name:    jobData["name"].(string),
		Payload: jobData["payload"].(map[string]interface{}),
	}
	
	return job, nil
}

func (rq *RedisQueue) Size() (int64, error) {
	ctx := context.Background()
	return rq.client.LLen(ctx, rq.queueName).Result()
}

func (rq *RedisQueue) Clear() error {
	ctx := context.Background()
	return rq.client.Del(ctx, rq.queueName).Err()
}

func (rq *RedisQueue) processDelayedJobs() {
	ctx := context.Background()
	now := float64(time.Now().Unix())
	
	// Get ready delayed jobs
	jobs, err := rq.client.ZRangeByScore(ctx, rq.queueName+":delayed", &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%f", now),
	}).Result()
	
	if err != nil || len(jobs) == 0 {
		return
	}
	
	// Move ready jobs to main queue
	for _, job := range jobs {
		rq.client.LPush(ctx, rq.queueName, job)
		rq.client.ZRem(ctx, rq.queueName+":delayed", job)
	}
}

// Worker processes jobs from queue
type Worker struct {
	queue       Queue
	handlers    map[string]func() Job
	concurrency int
	quit        chan bool
	wg          sync.WaitGroup
}

// NewWorker creates a new worker
func NewWorker(queue Queue, handlers map[string]func() Job, concurrency int) *Worker {
	return &Worker{
		queue:       queue,
		handlers:    handlers,
		concurrency: concurrency,
		quit:        make(chan bool),
	}
}

// Start starts the worker
func (w *Worker) Start() {
	for i := 0; i < w.concurrency; i++ {
		w.wg.Add(1)
		go w.work()
	}
	w.wg.Wait()
}

// Stop stops the worker
func (w *Worker) Stop() {
	close(w.quit)
	w.wg.Wait()
}

func (w *Worker) work() {
	defer w.wg.Done()
	
	for {
		select {
		case <-w.quit:
			return
		default:
			job, err := w.queue.Pop()
			if err != nil {
				log.Printf("Error popping job: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			
			if job == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			
			w.processJob(job)
		}
	}
}

func (w *Worker) processJob(job Job) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Job %s panicked: %v", job.GetName(), r)
		}
	}()
	
	// Get handler for job
	handler, exists := w.handlers[job.GetName()]
	if !exists {
		log.Printf("No handler found for job: %s", job.GetName())
		return
	}
	
	// Create job instance and set payload
	jobInstance := handler()
	jobInstance.SetPayload(job.GetPayload())
	
	// Execute job
	if err := jobInstance.Handle(); err != nil {
		log.Printf("Job %s failed: %v", job.GetName(), err)
		// TODO: Implement retry logic
	} else {
		log.Printf("Job %s completed successfully", job.GetName())
	}
}

// MemoryQueue implements in-memory queue for testing
type MemoryQueue struct {
	jobs []Job
	mutex sync.Mutex
}

func (mq *MemoryQueue) Push(job Job, delay ...time.Duration) error {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	mq.jobs = append(mq.jobs, job)
	return nil
}

func (mq *MemoryQueue) Pop() (Job, error) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	if len(mq.jobs) == 0 {
		return nil, nil
	}
	job := mq.jobs[0]
	mq.jobs = mq.jobs[1:]
	return job, nil
}

func (mq *MemoryQueue) Size() (int64, error) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	return int64(len(mq.jobs)), nil
}

func (mq *MemoryQueue) Clear() error {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	mq.jobs = nil
	return nil
}