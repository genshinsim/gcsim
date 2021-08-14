package character

import "github.com/genshinsim/gsim/pkg/core"

func (t *Tmpl) Tick() {
}

func (t *Tmpl) AddTask(fun func(), name string, delay int) {
	t.Core.Tasks.Add(fun, delay)
	t.Log.Debugw("task added: "+name, "frame", t.Core.F, "event", core.LogTaskEvent, "name", name, "delay", delay)
}

func (t *Tmpl) QueueDmg(ds *core.Snapshot, delay int) {
	t.AddTask(func() {
		t.Core.Combat.ApplyDamage(ds)
	}, "dmg", delay)
}
