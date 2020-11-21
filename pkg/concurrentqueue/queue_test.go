package concurrentqueue

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	parallelCount = 100
	inputs        = genIputs(parallelCount)
)

func TestEnqueueAndDeque(t *testing.T) {
	var q Queue = &queue{}
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

func TestDequeueN(t *testing.T) {
	var q Queue = &queue{}
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

func BenchmarkLockerQueueConcurrentEnqueue(b *testing.B) {
	expectedLength := parallelCount * b.N

	var q Queue = &queue{}

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(parallelCount)
		for j := 0; j < parallelCount; j++ {
			go func(i int) {
				q.Enqueue(inputs[i])
				wg.Done()
			}(j)
		}
		wg.Wait()
	}
	if q.Len() != expectedLength {
		b.Errorf("Unexpected length: %d", q.Len())
	}
}

func BenchmarkLockerQueueConcurrentDequeue(b *testing.B) {
	dequeLen := 3
	total := parallelCount * b.N * dequeLen
	inputs := genIputs(total)
	q := NewLockerQueue(inputs)

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(parallelCount)
		for j := 0; j < parallelCount; j++ {
			go func(i int) {
				q.Dequeue(dequeLen)
				wg.Done()
			}(j)
		}
		wg.Wait()
	}
	if q.Len() != 0 {
		b.Errorf("Unexpected length: %d", q.Len())
	}
}

func BenchmarkLockerQueueConcurrentComplexOps(b *testing.B) {
	dequeLen := 2
	var q Queue = &queue{}

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(parallelCount)
		for j := 0; j < parallelCount; j++ {
			if (j+1)%(dequeLen+1) == 0 {
				go func(i int) {
					q.Dequeue(dequeLen)
					wg.Done()
				}(j)
			} else {
				go func(i int) {
					q.Enqueue(inputs[i])
					wg.Done()
				}(j)
			}
		}
		wg.Wait()
	}
}

func BenchmarkLockerQueueEnqueue(b *testing.B) {
	var q Queue = &queue{}

	for i := 0; i < b.N; i++ {
		q.Enqueue(fmt.Sprintf("a-%d", i))
	}
}

func BenchmarkLockerQueueDequeue(b *testing.B) {
	var q Queue = NewLockerQueue(genIputs(b.N * 2))

	for i := 0; i < b.N; i++ {
		q.Dequeue(2)
	}
}

func genIputs(n int) []string {
	inputs := make([]string, n)
	for i := 0; i < n; i++ {
		inputs[i] = fmt.Sprintf("item-%d", i)
	}
	return inputs
}
