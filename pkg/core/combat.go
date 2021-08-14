package core

type CombatHandler interface {
	ApplyDamage(ds *Snapshot) float64
	TargetHasResMod(debuff string, param int) bool
	TargetHasDefMod(debuff string, param int) bool
	TargetHasElement(ele EleType, param int) bool
}

type CombatCtrl struct {
	core *Core
}

func NewCombatCtrl(c *Core) *CombatCtrl {
	return &CombatCtrl{
		core: c,
	}
}

func (c *CombatCtrl) ApplyDamage(ds *Snapshot) float64 {
	died := false
	var total float64
	for i, t := range c.core.Targets {
		d := ds.Clone()
		dmg, crit := t.Attack(&d)
		total += dmg

		//check if target is dead
		if c.core.Flags.DamageMode && t.HP() <= 0 {
			died = true
			c.core.Events.Emit(OnTargetDied, t, ds)
			c.core.Targets[i] = nil
			// log.Println("target died", i, dmg)
		}

		amp := ""
		if d.IsMeltVape {
			amp = string(d.ReactionType)
		}

		c.core.Log.Debugw(
			d.Abil,
			"frame", c.core.F,
			"event", LogDamageEvent,
			"char", d.ActorIndex,
			"target", i,
			"attack_tag", d.AttackTag,
			"damage", dmg,
			"crit", crit,
			"amp", amp,
			"abil", d.Abil,
		)

	}
	if died {
		//wipe out nil entries
		n := 0
		for _, v := range c.core.Targets {
			if v != nil {
				c.core.Targets[n] = v
				c.core.Targets[n].SetIndex(n)
				n++
			}
		}
		c.core.Targets = c.core.Targets[:n]
	}
	return total
}

func (c *CombatCtrl) TargetHasResMod(key string, param int) bool {
	if param >= len(c.core.Targets) {
		return false
	}
	return c.core.Targets[param].HasResMod(key)
}
func (c *CombatCtrl) TargetHasDefMod(key string, param int) bool {
	if param >= len(c.core.Targets) {
		return false
	}
	return c.core.Targets[param].HasDefMod(key)
}

func (c *CombatCtrl) TargetHasElement(ele EleType, param int) bool {
	if param >= len(c.core.Targets) {
		return false
	}
	return c.core.Targets[param].AuraContains(ele)
}
