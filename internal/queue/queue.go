package queue

import (
	ej "executor/internal/executorjob"
	"sync"
)

type Queue struct {
	items []ej.ExecutorJob
	mu    sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		items: make([]ej.ExecutorJob, 0),
	}
}

func (q *Queue) Enqueue(item ej.ExecutorJob) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}
func (q *Queue) Dequeue() (ej.ExecutorJob, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		var empty ej.ExecutorJob
		return empty, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}
