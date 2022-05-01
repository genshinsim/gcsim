package task

type task struct {
	source int
	f      func()
}

type Tasker interface {
	Add(f func(), delay int)
}

type Handler struct {
	f     *int
	tasks map[int][]task
}

func New(f *int) *Handler {
	c := &Handler{f: f}
	c.tasks = make(map[int][]task)
	return c
}

func (c *Handler) Run() {
	for _, x := range c.tasks[*c.f] {
		x.f()
	}
	delete(c.tasks, *c.f)
}

func (c *Handler) Add(f func(), delay int) {
	c.tasks[*c.f+delay] = append(c.tasks[*c.f+delay], task{
		f:      f,
		source: *c.f,
	})
}
