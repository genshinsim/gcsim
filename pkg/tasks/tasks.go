package tasks

type task struct {
	source int
	f      func()
}

type Controller struct {
	f     *int
	tasks map[int][]task
}

func New(f *int) *Controller {
	c := &Controller{f: f}
	c.tasks = make(map[int][]task)
	return c
}

func (c *Controller) Run() {
	for _, x := range c.tasks[*c.f] {
		x.f()
	}
	delete(c.tasks, *c.f)
}

func (c *Controller) Add(f func(), delay int) {
	c.tasks[*c.f+delay] = append(c.tasks[*c.f+delay], task{
		f:      f,
		source: *c.f,
	})
}
