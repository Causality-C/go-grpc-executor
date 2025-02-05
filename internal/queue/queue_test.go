package queue

import (
	ej "executor/internal/executorjob"
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewQueue()
	if q == nil {
		t.Errorf("Queue is nil")
	}
	_, ok := q.Dequeue()
	if ok {
		t.Errorf("Queue should be empty")
	}
	q.Enqueue(ej.ExecutorJob{})
	item, ok := q.Dequeue()
	if !ok {
		t.Errorf("Queue should not be empty")
	}
	if item.TaskID != 0 {
		t.Errorf("TaskID should be 0")
	}
}
