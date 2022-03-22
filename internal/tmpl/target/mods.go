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
		m.Event = t.DefMod[ind].Event
		m.Event.SetEnded(m.Expiry)
		t.DefMod[ind] = m
		return
	}
	m.Event = t.Core.Log.NewEvent("enemy mod added", core.LogStatusEvent, -1, "count", len(t.DefMod), "val", val, "target", t.TargetIndex, "key", m.Key, "expiry", m.Expiry)
	m.Event.SetEnded(m.Expiry)
	t.DefMod = append(t.DefMod, m)
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

func (t *Tmpl) AddResMod(key string, m core.ResistMod) {
	m.Expiry = t.Core.F + m.Duration
	m.Key = key
	//find if exists, if exists override, else append
	ind := -1
	for i, v := range t.ResMod {
		if v.Key == key {
			ind = i
		}
	}
	if ind != -1 {
		t.Core.Log.NewEvent("enemy mod refreshed", core.LogStatusEvent, -1, "count", len(t.ResMod), "val", m, "target", t.TargetIndex, "key", m.Key, "expiry", m.Expiry)
		m.Event = t.ResMod[ind].Event
		m.Event.SetEnded(m.Expiry)
		t.ResMod[ind] = m
		return
	}
	m.Event = t.Core.Log.NewEvent("enemy mod added", core.LogStatusEvent, -1, "count", len(t.ResMod), "val", m, "target", t.TargetIndex, "key", m.Key, "expiry", m.Expiry)
	m.Event.SetEnded(m.Expiry)
	t.ResMod = append(t.ResMod, m)
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
