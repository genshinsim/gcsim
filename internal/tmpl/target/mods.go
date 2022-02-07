package target

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (t *Tmpl) AddDefMod(key string, val float64, dur int) {
	m := core.DefMod{
		Key:    key,
		Value:  val,
		Expiry: t.Core.F + dur,
	}
	//find if exists, if exists override, else append
	ind := -1
	for i, v := range t.DefMod {
		if v.Key == key {
			ind = i
		}
	}
	if ind != -1 {
		t.Core.Log.NewEvent("enemy mod refreshed", core.LogStatusEvent, -1, "count", len(t.DefMod), "val", val, "target", t.TargetIndex, "key", m.Key, "expiry", m.Expiry)
		t.DefMod[ind] = m
		return
	}
	t.DefMod = append(t.DefMod, m)
	t.Core.Log.NewEvent("enemy mod added", core.LogStatusEvent, -1, "count", len(t.DefMod), "val", val, "target", t.TargetIndex, "key", m.Key, "expiry", m.Expiry)
	// e.mod[key] = val

	// Add task to check for mod expiry in debug instances
	if t.Core.Flags.LogDebug && m.Expiry > 0 {
		t.AddTask(func() {
			if t.HasDefMod(m.Key) {
				return
			}
			t.Core.Log.NewEvent("enemy mod expired", core.LogStatusEvent, -1, "count", len(t.DefMod), "val", val, "target", t.TargetIndex, "key", m.Key, "expiry", m.Expiry)
		}, "check-m-expiry", m.Expiry+1-t.Core.F)
	}
}

func (t *Tmpl) HasDefMod(key string) bool {
	ind := -1
	for i, v := range t.DefMod {
		if v.Key == key {
			ind = i
		}
	}
	return ind != -1 && t.DefMod[ind].Expiry > t.Core.F
}

func (t *Tmpl) AddResMod(key string, val core.ResistMod) {
	val.Expiry = t.Core.F + val.Duration
	val.Key = key
	//find if exists, if exists override, else append
	ind := -1
	for i, v := range t.ResMod {
		if v.Key == key {
			ind = i
		}
	}
	if ind != -1 {
		t.Core.Log.NewEvent("enemy mod refreshed", core.LogStatusEvent, -1, "count", len(t.ResMod), "val", val, "target", t.TargetIndex, "key", val.Key, "expiry", val.Expiry)
		t.ResMod[ind] = val
		return
	}
	t.ResMod = append(t.ResMod, val)
	t.Core.Log.NewEvent("enemy mod added", core.LogStatusEvent, -1, "count", len(t.ResMod), "val", val, "target", t.TargetIndex, "key", val.Key, "expiry", val.Expiry)
	// e.mod[key] = val

	// Add task to check for mod expiry in debug instances
	if t.Core.Flags.LogDebug && val.Expiry > -1 {
		t.AddTask(func() {
			if t.HasResMod(val.Key) {
				return
			}
			t.Core.Log.NewEvent("enemy mod expired", core.LogStatusEvent, -1, "count", len(t.ResMod), "val", val, "target", t.TargetIndex, "key", val.Key, "expiry", val.Expiry)
		}, "check-m-expiry", val.Expiry+1-t.Core.F)
	}
}

func (t *Tmpl) RemoveResMod(key string) {
	for i, v := range t.ResMod {
		if v.Key == key {
			t.ResMod[i].Expiry = 0
		}
	}
}

func (t *Tmpl) RemoveDefMod(key string) {
	for i, v := range t.DefMod {
		if v.Key == key {
			t.DefMod[i].Expiry = 0
		}
	}
}

func (t *Tmpl) HasResMod(key string) bool {
	ind := -1
	for i, v := range t.ResMod {
		if v.Key == key {
			ind = i
		}
	}
	return ind != -1 && t.ResMod[ind].Expiry > t.Core.F
}
