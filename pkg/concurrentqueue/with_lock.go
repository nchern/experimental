package concurrentqueue

import "sync"

type queue struct {
	items []string

	lock sync.RWMutex
}

func NewLockerQueue(items []string) Queue {
	var q queue
	copy(q.items, items)
	return &q
}

func (q *queue) Len() int {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return len(q.items)
}

func (q *queue) Dequeue(n int) ([]string, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if n > len(q.items) {
		return nil, false
	}

	res := q.items[:n]
	q.items = q.items[n:]

	return res, true
}

func (q *queue) Enqueue(s string) {
	q.lock.Lock()
	q.items = append(q.items, s)
	q.lock.Unlock()
}
