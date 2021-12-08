package target

import "github.com/genshinsim/gcsim/pkg/core"

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
		t.Core.Log.Debugw("mod overwritten", "frame", t.Core.F, "event", core.LogEnemyEvent, "count", len(t.DefMod), "old", t.DefMod[ind], "next", val, "target", t.TargetIndex)
		// LogEnemyEvent
		t.DefMod[ind] = m
		return
	}
	t.DefMod = append(t.DefMod, m)
	t.Core.Log.Debugw("new def mod", "frame", t.Core.F, "event", core.LogEnemyEvent, "count", len(t.DefMod), "next", val, "target", t.TargetIndex)
	// e.mod[key] = val
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
		t.Core.Log.Debugw("mod overwritten", "frame", t.Core.F, "event", core.LogEnemyEvent, "count", len(t.ResMod), "old", t.ResMod[ind], "next", val)
		// LogEnemyEvent
		t.ResMod[ind] = val
		return
	}
	t.ResMod = append(t.ResMod, val)
	t.Core.Log.Debugw("new mod", "frame", t.Core.F, "event", core.LogEnemyEvent, "count", len(t.ResMod), "next", val)
	// e.mod[key] = val
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
