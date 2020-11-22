package concurrentqueue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreadedQueueEnqueueAndDeque(t *testing.T) {
	var q Queue = NewThreadedQueue([]string{})
	assert.Equal(t, 0, q.Len())

	_, ok := q.Dequeue(1)
	assert.False(t, ok)

	count := 100
	for i := 0; i < count; i++ {
		q.Enqueue(fmt.Sprintf("a-%d", i))
	}
	assert.Equal(t, count, q.Len())
	for i := 0; i < count; i++ {
		actual, ok := q.Dequeue(1)
		assert.True(t, ok)
		assert.Equal(t, []string{fmt.Sprintf("a-%d", i)}, actual)
	}
}

func TestThreadedDequeueN(t *testing.T) {
	var q Queue = NewThreadedQueue([]string{})
	input := []string{"a", "b", "c", "d", "e"}
	for _, s := range input {
		q.Enqueue(s)
	}
	actual, ok := q.Dequeue(100)
	assert.Empty(t, actual)
	assert.False(t, ok)

	actual, ok = q.Dequeue(3)
	assert.Equal(t, []string{"a", "b", "c"}, actual)
	assert.Equal(t, 2, q.Len())

	actual, ok = q.Dequeue(2)
	assert.Equal(t, []string{"d", "e"}, actual)
	assert.Equal(t, 0, q.Len())

	actual, ok = q.Dequeue(1)
	assert.Empty(t, actual)
	assert.False(t, ok)
}

func BenchmarkThreadedQueueConcurrentEnqueue(b *testing.B) {
	var q Queue = NewThreadedQueue([]string{})
	doBenchmarkConcurrentEnqueue(b, q)
}

func BenchmarkThreadedQueueConcurrentDequeue(b *testing.B) {
	dequeLen := 3
	total := b.N * dequeLen
	inputs := genIputs(total)
	q := NewThreadedQueue(inputs)

	doBenchmarkConcurrentDequeue(b, q, dequeLen)
}

func BenchmarkThreadedQueueConcurrentComplexOps(b *testing.B) {
	var q Queue = NewThreadedQueue([]string{})
	doBenchmarkConcurrentComplexOps(b, q)
}

func BenchmarkThreadedQueueEnqueue(b *testing.B) {
	var q Queue = NewThreadedQueue([]string{})

	for i := 0; i < b.N; i++ {
		q.Enqueue(fmt.Sprintf("a-%d", i))
	}
}

func BenchmarkThreadedQueueDequeue(b *testing.B) {
	var q Queue = NewThreadedQueue(genIputs(b.N * 2))

	for i := 0; i < b.N; i++ {
		q.Dequeue(2)
	}
}
