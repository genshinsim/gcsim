package target

import "github.com/genshinsim/gcsim/pkg/core"

func (t *Tmpl) AddTask(fun func(), name string, delay int) {
	t.Core.Tasks.Add(fun, delay)
	t.Core.Log.NewEvent("task added: "+name, core.LogTaskEvent, -1, "name", name, "delay", delay)
}
