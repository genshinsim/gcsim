package character

import "github.com/genshinsim/gsim/pkg/core"

func (t *Tmpl) Tick() {

	//run tasks
	for _, x := range t.Tasks[t.Sim.Frame()] {
		x.F()
	}

	delete(t.Tasks, t.Sim.Frame())
}

func (t *Tmpl) AddTask(fun func(), name string, delay int) {
	f := t.Sim.Frame()
	t.Tasks[f+delay] = append(t.Tasks[f+delay], CharTask{
		Name:        name,
		F:           fun,
		Delay:       delay,
		originFrame: f,
	})
	t.Log.Debugw("task added: "+name, "frame", f, "event", core.LogTaskEvent, "name", name, "delay", delay)
}

func (t *Tmpl) QueueDmg(ds *core.Snapshot, delay int) {
	t.AddTask(func() {
		t.Sim.ApplyDamage(ds)
	}, "dmg", delay)
}
