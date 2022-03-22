package target

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (t *Tmpl) AddDefMod(key string, val float64, dur int) {
	mod := core.DefMod{
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

	// check if mod exists and has not expired
	if ind != -1 && (t.DefMod[ind].Expiry > t.Core.F || t.DefMod[ind].Expiry == -1) {
		t.Core.Log.NewEvent("enemy mod refreshed", core.LogStatusEvent, -1, "count", len(t.DefMod), "val", val, "target", t.TargetIndex, "key", mod.Key, "expiry", mod.Expiry)
		mod.Event = t.DefMod[ind].Event
	} else {
		mod.Event = t.Core.Log.NewEvent("enemy mod added", core.LogStatusEvent, -1, "count", len(t.DefMod), "val", val, "target", t.TargetIndex, "key", mod.Key, "expiry", mod.Expiry)
		// append empty mod if we can not reuse mods[ind]
		if ind == -1 {
			t.DefMod = append(t.DefMod, core.DefMod{})
			ind = len(t.DefMod) - 1
		}
	}
	mod.Event.SetEnded(mod.Expiry)
	t.DefMod[ind] = mod
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

func (t *Tmpl) AddResMod(key string, mod core.ResistMod) {
	mod.Expiry = t.Core.F + mod.Duration
	mod.Key = key
	//find if exists, if exists override, else append
	ind := -1
	for i, v := range t.ResMod {
		if v.Key == key {
			ind = i
		}
	}

	// check if mod exists and has not expired
	if ind != -1 && (t.ResMod[ind].Expiry > t.Core.F || t.ResMod[ind].Expiry == -1) {
		t.Core.Log.NewEvent("enemy mod refreshed", core.LogStatusEvent, -1, "count", len(t.ResMod), "val", mod, "target", t.TargetIndex, "key", mod.Key, "expiry", mod.Expiry)
		mod.Event = t.ResMod[ind].Event
	} else {
		mod.Event = t.Core.Log.NewEvent("enemy mod added", core.LogStatusEvent, -1, "count", len(t.ResMod), "val", mod, "target", t.TargetIndex, "key", mod.Key, "expiry", mod.Expiry)
		// append empty mod if we can not reuse mods[ind]
		if ind == -1 {
			t.ResMod = append(t.ResMod, core.ResistMod{})
			ind = len(t.ResMod) - 1
		}
	}
	mod.Event.SetEnded(mod.Expiry)
	t.ResMod[ind] = mod
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
