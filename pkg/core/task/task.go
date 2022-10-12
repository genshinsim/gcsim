package task

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
