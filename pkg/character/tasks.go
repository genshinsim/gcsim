package character

import "github.com/genshinsim/gcsim/pkg/core"

func (t *Tmpl) Tick() {
}

func (t *Tmpl) AddTask(fun func(), name string, delay int) {
	t.Core.Tasks.Add(fun, delay)
	t.Core.Log.Debugw("task added: "+name, "frame", t.Core.F, "event", core.LogTaskEvent, "name", name, "delay", delay)
}

// Main function used to deal damage. In the case of 0 delay, immediately procs damage on the same frame
// In the case of delay > 0, used for damage instances that snapshot prior to the damage occuring
// Examples include deployables like Guoba.
// func (t *Tmpl) QueueDmg(ds *core.Snapshot, delay int) {
// 	if delay == 0 {
// 		t.Core.Combat.ApplyDamage(ds)
// 	} else {
// 		t.AddTask(func() {
// 			t.Core.Combat.ApplyDamage(ds)
// 		}, "dmg", delay)
// 	}
// }

// Helper/descriptive function to create a snapshot instance and queue up the damage on the same frame.
// Best used for all abilities that do not snapshot in game (e.g. normal attacks, Hu Tao blood blossoms, etc.)
// func (t *Tmpl) QueueDmgDynamic(generateSnapshot func() *core.Snapshot, delay int) {
// 	t.AddTask(func() {
// 		t.Core.Combat.ApplyDamage(generateSnapshot())
// 	}, "dmg", delay)
// }

// Helper/descriptive function to create a snapshot instance at some point "delay" frames in the future
// Then queue up the damage on "snapshotDelay" frames after the creation of the snapshot
// Examples include bow charged attacks or other things that have a significant animation time before snapshotting
// func (t *Tmpl) QueueDmgDynamicSnapshotDelay(generateSnapshot func() *core.Snapshot, delay int, snapshotDelay int) {
// 	t.AddTask(func() {
// 		d := generateSnapshot()
// 		t.QueueDmg(d, snapshotDelay)
// 	}, "dmg", delay)
// }
