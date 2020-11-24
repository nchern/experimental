package concurrentqueue

type Cmd int

const (
	cmdLen Cmd = iota
	cmdEnqueue
	cmdDequeue
)

type task struct {
	cmd Cmd
	res chan bool

	resLen int

	enqueueItem string

	dequeueN   int
	resDequeue []string
}

type threadedQueue struct {
	items []string

	tasks chan *task
}

func executor(q *threadedQueue) {
	for tsk := range q.tasks {
		switch tsk.cmd {
		case cmdLen:
			tsk.resLen = len(q.items)
			tsk.res <- true
		case cmdEnqueue:
			q.items = append(q.items, tsk.enqueueItem)
			tsk.res <- true
		case cmdDequeue:
			n := tsk.dequeueN
			if n > len(q.items) {
				tsk.res <- false
				continue
			}

			tsk.resDequeue = q.items[:n]
			q.items = q.items[n:]
			tsk.res <- true
		default:
			// handle default
		}
	}
}

func NewThreadedQueue(items []string) Queue {
	q := &threadedQueue{tasks: make(chan *task, 10)}
	copy(q.items, items)

	go executor(q)

	return q
}

func (q *threadedQueue) Len() int {
	t := &task{res: make(chan bool, 1), cmd: cmdLen}
	q.tasks <- t

	<-t.res
	return t.resLen

}

func (q *threadedQueue) Dequeue(n int) ([]string, bool) {
	t := &task{dequeueN: n, res: make(chan bool, 1), cmd: cmdDequeue}
	q.tasks <- t

	return t.resDequeue, <-t.res
}

func (q *threadedQueue) Enqueue(s string) {
	t := &task{res: make(chan bool, 1), enqueueItem: s, cmd: cmdEnqueue}
	q.tasks <- t
	<-t.res
}
