package core

type TaskHandler interface {
	Add(f func(), delay int)
	Run()
}

type task struct {
	source int
	f      func()
}

type TaskCtrl struct {
	f     *int
	tasks map[int][]task
}

func NewTaskCtrl(f *int) *TaskCtrl {
	c := &TaskCtrl{f: f}
	c.tasks = make(map[int][]task)
	return c
}

func (c *TaskCtrl) Run() {
	for _, x := range c.tasks[*c.f] {
		x.f()
	}
	delete(c.tasks, *c.f)
}

func (c *TaskCtrl) Add(f func(), delay int) {
	c.tasks[*c.f+delay] = append(c.tasks[*c.f+delay], task{
		f:      f,
		source: *c.f,
	})
}
