package task

// TODO: the behavior of delay<=0 is inconsistent
// TODO: consider merging all tasks into a single handler
// Currently tasks are executed in the following order: (enemy1, enemy2, ...), (char1, char2, ...), (core tasks)
// In order to replace, the core task queue must support the ability to update the position of a task in the queue.
// Also will need to consider order. Currently everything that is queued via QueueCharTask/QueueEnemyTask will
// always happen before all entries in the core task queue. If any implementations depend on this order,
// this will cause additional problems.

// minHeap is a binary min-heap ordered by (executeBy, id). It is typed rather than built on
// container/heap so that Add and Run don't box every task into an interface or dispatch
// Less/Swap through one.
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

// taskCap is the initial capacity of a task queue
const taskCap = 53

func newMinHeap() *minHeap {
	h := make(minHeap, 0, taskCap)
	return &h
}

func New(f *int) *Handler {
	return &Handler{
		f:     f,
		tasks: newMinHeap(),
	}
}

func (s *Handler) Run() {
	for len(*s.tasks) > 0 && (*s.tasks)[0].executeBy <= *s.f {
		s.tasks.pop().f()
	}
}

func (s *Handler) Add(f func(), delay int) {
	s.tasks.push(task{
		executeBy: *s.f + delay,
		f:         f,
		id:        s.counter,
	})
	s.counter += 1
}

func (s *Handler) Extend(delay int) {
	for i := range *s.tasks {
		(*s.tasks)[i].extend(delay)
	}
}

// min heap functions. push and pop make the same comparisons and leave the same layout as
// container/heap's Push and Pop, but shift elements into a hole instead of swapping.

func (t *task) less(o *task) bool {
	return t.executeBy < o.executeBy || (t.executeBy == o.executeBy && t.id < o.id)
}

func (h *minHeap) push(t task) {
	*h = append(*h, t)
	s := *h
	j := len(s) - 1
	for j > 0 {
		i := (j - 1) / 2 // parent
		if !t.less(&s[i]) {
			break
		}
		s[j] = s[i]
		j = i
	}
	s[j] = t
}

func (h *minHeap) pop() task {
	s := *h
	n := len(s) - 1
	top, last := s[0], s[n]
	s[n] = task{} // drop the closure so it can be collected
	s = s[:n]
	*h = s
	if n == 0 {
		return top
	}
	i := 0
	for {
		j := 2*i + 1 // left child
		if j >= n {
			break
		}
		if j2 := j + 1; j2 < n && s[j2].less(&s[j]) {
			j = j2 // right child
		}
		if !s[j].less(&last) {
			break
		}
		s[i] = s[j]
		i = j
	}
	s[i] = last
	return top
}

func (t *task) extend(delay int) {
	t.executeBy += delay
}
