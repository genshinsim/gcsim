package target

import "github.com/genshinsim/gcsim/pkg/core"

func (t *Tmpl) AddTask(fun func(), name string, delay int) {
	t.Core.Tasks.Add(fun, delay)
	t.Core.Log.Debugw("task added: "+name, "frame", t.Core.F, "event", core.LogTaskEvent, "name", name, "delay", delay)
}
