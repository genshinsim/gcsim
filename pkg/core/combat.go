package core

import (
	"log"
)

type CombatHandler interface {
	ApplyDamage(*AttackEvent) float64
	QueueAttack(a AttackInfo, p AttackPattern, snapshotDelay int, dmgDelay int, callbacks ...AttackCBFunc)
	QueueAttackWithSnap(a AttackInfo, s Snapshot, p AttackPattern, dmgDelay int, callbacks ...AttackCBFunc)
	QueueAttackEvent(ae *AttackEvent, dmgDelay int)
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

func (c *CombatCtrl) QueueAttackWithSnap(a AttackInfo, s Snapshot, p AttackPattern, dmgDelay int, callbacks ...AttackCBFunc) {
	if dmgDelay < 0 {
		panic("dmgDelay cannot be less than 0")
	}
	ae := AttackEvent{
		Info:    a,
		Pattern: p,
		// Timing: AttackTiming{
		// 	SnapshotDelay: snapshotDelay,
		// 	DamageDelay:   dmgDelay,
		// },
		Snapshot:    s,
		SourceFrame: c.core.F,
	}
	//add callbacks only if not nil
	for _, f := range callbacks {
		if f != nil {
			ae.Callbacks = append(ae.Callbacks, f)
		}
	}
	c.queueDmg(&ae, dmgDelay)
}

func (c *CombatCtrl) QueueAttackEvent(ae *AttackEvent, dmgDelay int) {
	c.queueDmg(ae, dmgDelay)
}

func (c *CombatCtrl) QueueAttack(a AttackInfo, p AttackPattern, snapshotDelay int, dmgDelay int, callbacks ...AttackCBFunc) {
	//panic if dmgDelay > snapshotDelay; this should not happen. if it happens then there's something wrong with the
	//character's code
	if dmgDelay < snapshotDelay {
		panic("dmgDelay cannot be less than snapshotDelay")
	}
	if dmgDelay < 0 {
		panic("dmgDelay cannot be less than 0")
	}
	//create attackevent
	ae := AttackEvent{
		Info:    a,
		Pattern: p,
		// Timing: AttackTiming{
		// 	SnapshotDelay: snapshotDelay,
		// 	DamageDelay:   dmgDelay,
		// },
		SourceFrame: c.core.F,
	}
	//add callbacks only if not nil
	for _, f := range callbacks {
		if f != nil {
			ae.Callbacks = append(ae.Callbacks, f)
		}
	}
	// log.Println(ae)

	switch {
	case snapshotDelay < 0:
		//snapshotDelay < 0 means we don't need a snapshot; optimization for reaction
		//damage essentially
		c.queueDmg(&ae, dmgDelay)
	case snapshotDelay == 0:
		c.generateSnapshot(&ae)
		c.queueDmg(&ae, dmgDelay)
	default:
		//use add task ctrl to queue; no need to track here
		c.core.Tasks.Add(func() {
			c.generateSnapshot(&ae)
			c.queueDmg(&ae, dmgDelay-snapshotDelay)
		}, snapshotDelay)
	}

}

func (c *CombatCtrl) generateSnapshot(a *AttackEvent) {
	a.Snapshot = c.core.Chars[a.Info.ActorIndex].Snapshot(&a.Info)
}

func (c *CombatCtrl) queueDmg(a *AttackEvent, delay int) {
	if delay == 0 {
		c.ApplyDamage(a)
		return
	}
	c.core.Tasks.Add(func() {
		c.ApplyDamage(a)
	}, delay)
}

func willAttackLand(a *AttackEvent, t Target, index int) (bool, string) {
	//shape shouldn't be nil; panic here
	if a.Pattern.Shape == nil {
		panic("unexpected nil shape")
	}
	//shape can't be nil now, check if type matches
	if !a.Pattern.Targets[t.Type()] {
		return false, "wrong type"
	}
	//skip if self harm is false and dmg src == i
	if !a.Pattern.SelfHarm && a.Info.DamageSrc == index {
		return false, "no self harm"
	}

	//check if shape matches
	switch v := a.Pattern.Shape.(type) {
	case *Circle:
		return t.Shape().IntersectCircle(*v), "intersect circle"
	case *Rectangle:
		return t.Shape().IntersectRectangle(*v), "intersect rectangle"
	case *SingleTarget:
		//only true if
		return v.Target == index, "target"
	default:
		return false, "unknown shape"
	}
}

func (c *CombatCtrl) ApplyDamage(a *AttackEvent) float64 {
	died := false
	var total float64
	for i, t := range c.core.Targets {

		willHit, reason := willAttackLand(a, t, i)
		if !willHit {
			if c.core.Flags.LogDebug {
				c.core.Log.Debugw("skipped "+a.Info.Abil+" "+reason,
					"frame", c.core.F,
					"event", LogElementEvent,
					"char", a.Info.ActorIndex,
					"attack_tag", a.Info.AttackTag,
					"applied_ele", a.Info.Element,
					"dur", a.Info.Durability,
					"target", i,
					"shape", a.Pattern.Shape.String(),
					// "type", fmt.Sprintf("%T", a.Pattern.Shape),
				)
			}
			continue
		}

		//make a copy first
		cpy := *a

		//at this point attack will land
		c.core.Events.Emit(OnAttackWillLand, t, &cpy)

		//check to make sure it's not cancelled for w/e reason
		if a.Cancelled {
			continue
		}

		// if c.core.Flags.LogDebug {
		// 	c.core.Log.Debugw(a.Info.Abil+" will land",
		// 		"frame", c.core.F,
		// 		"event", LogElementEvent,
		// 		"char", a.Info.ActorIndex,
		// 		"attack_tag", a.Info.AttackTag,
		// 		"applied_ele", a.Info.Element,
		// 		"dur", a.Info.Durability,
		// 		"target", i,
		// 	)
		// }

		char := c.core.Chars[cpy.Info.ActorIndex]
		char.PreDamageSnapshotAdjust(&cpy, t)

		dmg, crit := t.Attack(&cpy)
		total += dmg

		c.core.Events.Emit(OnDamage, t, &cpy, dmg, crit)

		//callbacks
		cb := AttackCB{
			Target:      t,
			AttackEvent: &cpy,
			Damage:      dmg,
			IsCrit:      crit,
		}
		for _, f := range cpy.Callbacks {
			f(cb)
		}

		//check if target is dead; skip this for i = 0 since we don't want to
		//delete the player by accident
		if c.core.Flags.DamageMode && t.HP() <= 0 {
			log.Println("died")
			died = true
			t.Kill()
			c.core.Events.Emit(OnTargetDied, t, cpy)
			//this should be ok for stuff like guoba since they won't take damage
			c.core.Targets[i] = nil
			// log.Println("target died", i, dmg)
		}

		amp := ""
		if cpy.Info.Amped {
			amp = string(cpy.Info.AmpType)
		}

		c.core.Log.Debugw(
			cpy.Info.Abil,
			"frame", c.core.F,
			"event", LogDamageEvent,
			"char", cpy.Info.ActorIndex,
			"target", i,
			"attack_tag", cpy.Info.AttackTag,
			"damage", dmg,
			"crit", crit,
			"amp", amp,
			"abil", cpy.Info.Abil,
			"source", cpy.SourceFrame,
		)

	}
	if died {
		c.core.ReindexTargets()
	}
	c.core.TotalDamage += total
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
