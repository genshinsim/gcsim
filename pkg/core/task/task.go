package task

// TODO: the behavior of delay<=0 is inconsistent
// TODO: consider merging all tasks into a single handler
// Currently tasks are executed in the following order: (enemy1, enemy2, ...), (char1, char2, ...), (core tasks)
// In order to replace, the core task queue must support the ability to update the position of a task in the queue.
// Also will need to consider order. Currently everything that is queued via QueueCharTask/QueueEnemyTask will
// always happen before all entries in the core task queue. If any implementations depend on this order,
// this will cause additional problems.

import "container/heap"

type minHeap []task

type task struct {
	executeBy int
	f         func()
	id        int
}

type Handler struct {
	f       *int
	tasks   *minHeap
	counter int
}

type Tasker interface {
	Add(f func(), delay int)
}

func New(f *int) *Handler {
	return &Handler{
		f:     f,
		tasks: &minHeap{},
	}
}

func (s *Handler) Run() {
	for s.tasks.Len() > 0 && s.tasks.Peek().executeBy <= *s.f {
		heap.Pop(s.tasks).(task).f()
	}
}

func (s *Handler) Add(f func(), delay int) {
	heap.Push(s.tasks, task{
		executeBy: *s.f + delay,
		f:         f,
		id:        s.counter,
	})
	s.counter += 1
}

// min heap functions

func (h minHeap) Len() int {
	return len(h)
}

func (h minHeap) Less(i, j int) bool {
	return h[i].executeBy < h[j].executeBy || (h[i].executeBy == h[j].executeBy && h[i].id < h[j].id)
}

func (h minHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *minHeap) Push(x any) {
	*h = append(*h, x.(task))
}

func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h minHeap) Peek() task {
	return h[0]
}
