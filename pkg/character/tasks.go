package character

import "github.com/genshinsim/gcsim/pkg/core"

func (t *Tmpl) Tick() {
}

func (t *Tmpl) AddTask(fun func(), name string, delay int) {
	t.Core.Tasks.Add(fun, delay)
	t.Core.Log.Debugw("task added: "+name, "frame", t.Core.F, "event", core.LogTaskEvent, "name", name, "delay", delay)
}

func (t *Tmpl) QueueDmg(ds *core.Snapshot, delay int) {
	t.AddTask(func() {
		t.Core.Combat.ApplyDamage(ds)
	}, "dmg", delay)
}

// Helper/descriptive function to create a snapshot instance and queue up the damage on the same frame.
// Best used for all abilities that do not snapshot in game (e.g. normal attacks, Hu Tao blood blossoms, etc.)
func (t *Tmpl) QueueDmgDynamic(generateSnapshot func() *core.Snapshot, delay int) {
	t.AddTask(func() {
		t.Core.Combat.ApplyDamage(generateSnapshot())
	}, "dmg", delay)
}
