package core

type task struct {
	source int
	f      func()
}

func (c *Core) runTasks() {
	for _, x := range c.tasks[c.Frame] {
		x.f()
	}
	delete(c.tasks, c.Frame)
}

func (c *Core) Add(f func(), delay int) {
	c.tasks[c.Frame+delay] = append(c.tasks[c.Frame+delay], task{
		f:      f,
		source: c.Frame,
	})
}
