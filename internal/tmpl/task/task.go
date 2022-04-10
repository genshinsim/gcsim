package task

type task struct {
	source int
	f      func()
}

type Ctrl struct {
	f     *int
	tasks map[int][]task
}

func NewCtrl(f *int) *Ctrl {
	c := &Ctrl{f: f}
	c.tasks = make(map[int][]task)
	return c
}

func (c *Ctrl) Run() {
	for _, x := range c.tasks[*c.f] {
		x.f()
	}
	delete(c.tasks, *c.f)
}

func (c *Ctrl) Add(f func(), delay int) {
	c.tasks[*c.f+delay] = append(c.tasks[*c.f+delay], task{
		f:      f,
		source: *c.f,
	})
}
