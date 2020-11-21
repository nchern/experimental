package concurrentqueue

import (
	"fmt"
	"sync"
	"testing"
)

var (
	parallelCount = 100
	inputs        = genIputs(parallelCount)

//	fullQueue     = NewLockerQueue()
)

func init() {
	return
}

func TestEnqueue(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(parallelCount)

	inputs := genIputs(parallelCount)

	var q Queue = &queue{}
	for j := 0; j < parallelCount; j++ {
		go func(i int) {
			q.Enqueue(inputs[i])
			wg.Done()
		}(j)
	}
	wg.Wait()
	if q.Len() != parallelCount {
		t.Errorf("Unexpected length: %d(!=%d", q.Len(), parallelCount)
	}
}

func BenchmarkLockerQueueEnqueue(b *testing.B) {
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

func BenchmarkLockerQueueDequeue(b *testing.B) {
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

func BenchmarkLockerQueueComplexOps(b *testing.B) {
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

func genIputs(n int) []string {
	inputs := make([]string, n)
	for i := 0; i < n; i++ {
		inputs[i] = fmt.Sprintf("item-%d", i)
	}
	return inputs
}
