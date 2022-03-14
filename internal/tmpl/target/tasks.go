package target

import "github.com/genshinsim/gcsim/pkg/coretype"

func (t *Tmpl) AddTask(fun func(), name string, delay int) {
	t.Core.Tasks.Add(fun, delay)
	t.coretype.Log.NewEvent("task added: "+name, coretype.LogTaskEvent, -1, "name", name, "delay", delay)
}
