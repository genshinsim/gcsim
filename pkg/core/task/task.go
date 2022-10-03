package task

type Tasker interface {
	Add(f func(), delay int)
}

func New(f *int) *SliceHandler {
	return &SliceHandler{
		f: f,
	}
}

type SliceHandler struct {
	f     *int
	tasks []sliceTask
}

type sliceTask struct {
	source    int
	executeBy int
	f         func()
}

func (s *SliceHandler) Run() {
	//execute all tasks with executedBy <= f
	n := 0
	for i := 0; i < len(s.tasks); i++ {
		if s.tasks[i].executeBy <= *s.f {
			s.tasks[i].f()
		} else {
			s.tasks[n] = s.tasks[i]
			n++
		}
	}
	s.tasks = s.tasks[:n]
}

func (s *SliceHandler) Add(f func(), delay int) {
	s.tasks = append(s.tasks, sliceTask{
		source:    *s.f,
		executeBy: *s.f + delay,
		f:         f,
	})
}
