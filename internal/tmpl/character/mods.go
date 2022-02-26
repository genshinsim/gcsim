package character

import "github.com/genshinsim/gcsim/pkg/core"

func (c *Tmpl) AddPreDamageMod(mod core.PreDamageMod) {
	ind := -1
	for i, v := range c.PreDamageMods {
		if v.Key == mod.Key {
			ind = i
		}
	}
	if ind > -1 {
		ref := c.Core.Log.NewEvent("mod refreshed", core.LogStatusEvent, c.Index, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		//pull out existing event and update the ended
		mod.Evt = c.PreDamageMods[ind].Evt
		if mod.Expiry != -1 {
			mod.Evt.SetEnded(mod.Expiry)
			ref.SetEnded(mod.Expiry)
		}
		c.PreDamageMods[ind] = mod

	} else {
		//create new event
		mod.Evt = c.Core.Log.NewEvent("mod added", core.LogStatusEvent, c.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
		if mod.Expiry != -1 {
			mod.Evt.SetEnded(mod.Expiry)
		}
		c.PreDamageMods = append(c.PreDamageMods, mod)
	}

	// Add task to check for mod expiry in debug instances
	// if c.Core.Flags.LogDebug && mod.Expiry > -1 {
	// 	c.AddTask(func() {
	// 		if c.PreDamageModIsActive(mod.Key) {
	// 			return
	// 		}
	// 		c.Core.Log.NewEvent("mod expired", core.LogStatusEvent, c.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
	// 	}, "check-mod-expiry", mod.Expiry+1-c.Core.F)
	// }
}

func (c *Tmpl) AddMod(mod core.CharStatMod) {
	ind := -1
	for i, v := range c.Mods {
		if v.Key == mod.Key {
			ind = i
		}
	}
	if ind > -1 {
		ref := c.Core.Log.NewEvent("mod refreshed", core.LogStatusEvent, c.Index, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		//pull out existing event and update the ended
		mod.Evt = c.Mods[ind].Evt
		if mod.Expiry != -1 {
			mod.Evt.SetEnded(mod.Expiry)
			ref.SetEnded(mod.Expiry)
		}
		c.Mods[ind] = mod
	} else {
		mod.Evt = c.Core.Log.NewEvent("mod added", core.LogStatusEvent, c.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
		if mod.Expiry != -1 {
			mod.Evt.SetEnded(mod.Expiry)
		}
		c.Mods = append(c.Mods, mod)
	}

	// Add task to check for mod expiry in debug instances
	// if c.Core.Flags.LogDebug && mod.Expiry > -1 {
	// 	c.AddTask(func() {
	// 		if c.ModIsActive(mod.Key) {
	// 			return
	// 		}
	// 		c.Core.Log.NewEvent("mod expired", core.LogStatusEvent, c.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
	// 	}, "check-mod-expiry", mod.Expiry+1-c.Core.F)
	// }
}

func (t *Tmpl) AddReactBonusMod(mod core.ReactionBonusMod) {
	ind := -1
	for i, v := range t.ReactMod {
		if v.Key == mod.Key {
			ind = i
		}
	}
	if ind != -1 {
		ref := t.Core.Log.NewEvent("mod refreshed", core.LogStatusEvent, t.Index, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		mod.Evt = t.ReactMod[ind].Evt
		if mod.Expiry != -1 {
			mod.Evt.SetEnded(mod.Expiry)
			ref.SetEnded(mod.Expiry)
		}
		t.ReactMod[ind] = mod

		return
	}
	mod.Evt = t.Core.Log.NewEvent("mod added", core.LogStatusEvent, t.Index, "key", mod.Key, "expiry", mod.Expiry)
	if mod.Expiry != -1 {
		mod.Evt.SetEnded(mod.Expiry)
	}
	t.ReactMod = append(t.ReactMod, mod)

	// Add task to check for mod expiry in debug instances
	// if t.Core.Flags.LogDebug && mod.Expiry > -1 {
	// 	t.AddTask(func() {
	// 		if t.ReactBonusModIsActive(mod.Key) {
	// 			return
	// 		}
	// 		t.Core.Log.NewEvent("mod expired", core.LogStatusEvent, t.Index, "key", mod.Key, "expiry", mod.Expiry)
	// 	}, "check-mod-expiry", mod.Expiry+1-t.Core.F)
	// }
}
