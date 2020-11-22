package concurrentqueue

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

const parallelCount = 100

var (
	inputs = genIputs(parallelCount)
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

func doBenchmarkConcurrentEnqueue(b *testing.B, q Queue) {
	var expectedLength int64

	b.SetParallelism(runtime.NumCPU() * 10)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Enqueue("a")
			atomic.AddInt64(&expectedLength, 1)
		}
	})

	if q.Len() != int(expectedLength) {
		b.Errorf("Unexpected length: %d", q.Len())
	}

}

func doBenchmarkConcurrentDequeue(b *testing.B, q Queue, dequeLen int) {
	b.SetParallelism(runtime.NumCPU() * 10)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Dequeue(dequeLen)
		}
	})
	assert.Zero(b, q.Len())
}

func doBenchmarkConcurrentComplexOps(b *testing.B, q Queue) {
	var i int64
	var removed int64
	dequeLen := 2
	b.SetParallelism(runtime.NumCPU() * 10)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.AddInt64(&i, 0)
			if (int(i)+1)%(dequeLen+1) == 0 {
				_, ok := q.Dequeue(dequeLen)
				assert.True(b, ok)
				atomic.AddInt64(&removed, int64(dequeLen))
			} else {
				q.Enqueue("a")
			}
		}
	})
	//assert.Equal(b, int(i-removed+1), q.Len())
}

func BenchmarkLockerQueueConcurrentEnqueue(b *testing.B) {
	var q Queue = &queue{}
	doBenchmarkConcurrentEnqueue(b, q)
}

func BenchmarkLockerQueueConcurrentDequeue(b *testing.B) {
	dequeLen := 3
	total := b.N * dequeLen
	inputs := genIputs(total)
	q := NewLockerQueue(inputs)

	doBenchmarkConcurrentDequeue(b, q, dequeLen)
}

func BenchmarkLockerQueueConcurrentComplexOps(b *testing.B) {
	var q Queue = &queue{}
	doBenchmarkConcurrentComplexOps(b, q)
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
