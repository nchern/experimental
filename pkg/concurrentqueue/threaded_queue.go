package concurrentqueue

type taskLen struct {
	res chan int
}

type taskEnqueue struct {
	item string

	res chan bool
}

type taskDequeue struct {
	n int

	items []string
	res   chan bool
}

type threadedQueue struct {
	items []string

	tasks chan interface{}
}

func executor(q *threadedQueue) {
	for t := range q.tasks {
		switch tsk := t.(type) {
		case *taskLen:
			tsk.res <- len(q.items)
		case *taskEnqueue:
			q.items = append(q.items, tsk.item)
			tsk.res <- true
		case *taskDequeue:
			n := tsk.n
			if n > len(q.items) {
				tsk.res <- false
				continue
			}

			tsk.items = q.items[:n]
			q.items = q.items[n:]
			tsk.res <- true
		default:
			// handle default
		}
	}
}

func NewThreadedQueue(items []string) Queue {
	q := &threadedQueue{tasks: make(chan interface{}, 10)}
	copy(q.items, items)

	go executor(q)

	return q
}

func (q *threadedQueue) Len() int {
	t := &taskLen{make(chan int, 1)}
	q.tasks <- t

	return <-t.res
}

func (q *threadedQueue) Dequeue(n int) ([]string, bool) {
	t := &taskDequeue{n: n, res: make(chan bool, 1)}
	q.tasks <- t
	return t.items, <-t.res
}

func (q *threadedQueue) Enqueue(s string) {
	t := &taskEnqueue{res: make(chan bool, 1), item: s}
	q.tasks <- t
	<-t.res
}
