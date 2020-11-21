package concurrentqueue

type Queue interface {
	Len() int
	Dequeue(n int) ([]string, bool)
	Enqueue(string)
}
